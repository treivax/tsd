package main

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// put déplace un tuple d'un état initial vers un nouvel état et nettoie l'entrée /taken
// Si l'état de destination est "failed", la fonction ne fait rien pour permettre un nouveau take
//
// Paramètres:
//   - client: client etcd
//   - initialState: l'état d'origine du tuple
//   - position: la position du tuple dans l'état initial
//   - value: la valeur à déplacer
//   - newState: l'état de destination
//
// Retourne:
//   - int64: la nouvelle position dans l'état de destination (ou 0 si failed)
//   - error: toute erreur rencontrée
func put(client *clientv3.Client, initialState string, position int64, value string, newState string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Si l'état de destination est "failed", ne rien faire
	// Cela permet à un nouvel appel de take de retenter l'exécution
	if newState == "failed" {
		return 0, nil
	}

	// Clés impliquées
	initialTupleKey := fmt.Sprintf("/%s/%d", initialState, position)
	takenKey := fmt.Sprintf("/taken/%d", position)

	// Vérifier que le tuple existe dans l'état initial
	getTupleResp, err := client.Get(ctx, initialTupleKey)
	if err != nil {
		return 0, fmt.Errorf("failed to check initial tuple: %w", err)
	}
	if len(getTupleResp.Kvs) == 0 {
		return 0, fmt.Errorf("no tuple found at position %d in state %s", position, initialState)
	}

	// Utiliser la fonction move_tuple pour déplacer atomiquement le tuple
	newPosition, err := move_tuple(client, initialState, position, newState)
	if err != nil {
		return 0, fmt.Errorf("failed to move tuple from %s to %s: %w", initialState, newState, err)
	}

	// Supprimer l'entrée /taken/$position de manière séparée
	// (move_tuple ne gère que le déplacement des tuples, pas les entrées /taken)
	_, err = client.Delete(ctx, takenKey)
	if err != nil {
		// Ne pas échouer si l'entrée /taken n'existe pas
		// Cela peut arriver si put est appelé sans take préalable
		fmt.Printf("Warning: could not delete taken key %s: %v\n", takenKey, err)
	}

	return newPosition, nil
}

// put_batch déplace plusieurs tuples atomiquement d'un état vers un autre
// Utile pour traiter plusieurs résultats en une seule transaction
func put_batch(client *clientv3.Client, moves []TupleMove) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var allOps []clientv3.Op

	// Préparer toutes les opérations de déplacement
	for _, move := range moves {
		// Ignorer les états "failed"
		if move.NewState == "failed" {
			continue
		}

		// Préparer les opérations pour ce déplacement
		deleteOps, err := delete_tuple_raw(ctx, client, move.InitialState, move.Position)
		if err != nil {
			return fmt.Errorf("failed to prepare delete for position %d in state %s: %w",
				move.Position, move.InitialState, err)
		}

		addOps, _, err := new_tuple_raw(ctx, client, move.NewState, move.Value)
		if err != nil {
			return fmt.Errorf("failed to prepare add to state %s: %w", move.NewState, err)
		}

		// Ajouter les opérations à la liste globale
		allOps = append(allOps, deleteOps...)
		allOps = append(allOps, addOps...)

		// Ajouter la suppression de l'entrée /taken
		takenKey := fmt.Sprintf("/taken/%d", move.Position)
		allOps = append(allOps, clientv3.OpDelete(takenKey))
	}

	// Si aucune opération à effectuer (tous les états étaient "failed")
	if len(allOps) == 0 {
		return nil
	}

	// Exécuter toutes les opérations dans une seule transaction
	txn := client.Txn(ctx)
	txn = txn.Then(allOps...)

	_, err := txn.Commit()
	if err != nil {
		return fmt.Errorf("batch put transaction failed: %w", err)
	}

	return nil
}

// TupleMove représente un déplacement de tuple
type TupleMove struct {
	InitialState string
	Position     int64
	Value        string
	NewState     string
}

// put_with_cleanup effectue un put avec nettoyage complet des métadonnées
// Supprime non seulement /taken mais aussi d'éventuelles autres métadonnées
func put_with_cleanup(client *clientv3.Client, initialState string, position int64, value string, newState string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Si l'état de destination est "failed", ne rien faire
	if newState == "failed" {
		return 0, nil
	}

	// Déplacer le tuple
	newPosition, err := put(client, initialState, position, value, newState)
	if err != nil {
		return 0, err
	}

	// Nettoyage supplémentaire : supprimer d'éventuelles métadonnées liées à cette position
	// Par exemple, des logs, des metrics, etc.
	metadataPrefix := fmt.Sprintf("/metadata/%d", position)
	_, err = client.Delete(ctx, metadataPrefix, clientv3.WithPrefix())
	if err != nil {
		fmt.Printf("Warning: could not clean metadata for position %d: %v\n", position, err)
	}

	return newPosition, nil
}

// put_conditional effectue un put conditionnel basé sur l'état actuel
// Utile pour éviter les conflits de concurrence
func put_conditional(client *clientv3.Client, initialState string, position int64, value string, newState string, expectedValue string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Si l'état de destination est "failed", ne rien faire
	if newState == "failed" {
		return 0, nil
	}

	initialTupleKey := fmt.Sprintf("/%s/%d", initialState, position)
	takenKey := fmt.Sprintf("/taken/%d", position)

	// Préparer les opérations de déplacement
	deleteOps, err := delete_tuple_raw(ctx, client, initialState, position)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare delete operations: %w", err)
	}

	addOps, newPosition, err := new_tuple_raw(ctx, client, newState, value)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare add operations: %w", err)
	}

	// Créer une transaction conditionnelle
	txn := client.Txn(ctx)

	// Condition : la valeur actuelle doit correspondre à expectedValue
	txn = txn.If(clientv3.Compare(clientv3.Value(initialTupleKey), "=", expectedValue))

	// Si la condition est vraie, effectuer le déplacement
	var allOps []clientv3.Op
	allOps = append(allOps, deleteOps...)
	allOps = append(allOps, addOps...)
	allOps = append(allOps, clientv3.OpDelete(takenKey))
	txn = txn.Then(allOps...)

	// Si la condition est fausse, ne rien faire
	txn = txn.Else()

	resp, err := txn.Commit()
	if err != nil {
		return 0, fmt.Errorf("conditional put transaction failed: %w", err)
	}

	if !resp.Succeeded {
		return 0, fmt.Errorf("conditional put failed: current value does not match expected value %s", expectedValue)
	}

	return newPosition, nil
}
