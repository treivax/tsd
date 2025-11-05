package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Constantes
const (
	max_retry = 3
	timeout   = 30 * time.Second
)

// exec_state représente l'état d'exécution d'un tuple
type exec_state struct {
	Retry     int   `json:"retry"`
	Timestamp int64 `json:"timestamp"`
}

// call_job fonction fake qui simule l'appel d'un job
// Prend en argument la position et la valeur, retourne une nouvelle valeur ou une erreur
func call_job(position int64, value string) (string, error) {
	// TODO: Remplacer par la vraie logique métier
	// Pour l'instant, simulation simple :

	// Simuler un traitement qui peut échouer parfois
	if position%7 == 0 {
		return "", fmt.Errorf("job failed for position %d with value %s", position, value)
	}

	// Simuler une transformation de la valeur
	newValue := fmt.Sprintf("processed_%s_at_pos_%d", value, position)
	return newValue, nil
}

// take récupère un tuple d'un état donné et gère sa réservation avec retry et timeout
func take(client *clientv3.Client, state string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Étape 1: Récupérer un tuple via one_tuple
	position, value, err := one_tuple(client, state)
	if err != nil {
		return "", fmt.Errorf("no tuple available in state %s: %w", state, err)
	}

	takenKey := fmt.Sprintf("/taken/%d", position)
	stateKey := fmt.Sprintf("/%s/%d", state, position)
	errorKey := fmt.Sprintf("/error/%d", position)

	// Étape 2: Vérifier si la clé /taken/$position existe
	getTakenResp, err := client.Get(ctx, takenKey)
	if err != nil {
		return "", fmt.Errorf("failed to check taken key: %w", err)
	}

	var execState exec_state
	now := time.Now().Unix()

	if len(getTakenResp.Kvs) == 0 {
		// La clé n'existe pas, initialiser une nouvelle exec_state
		// IMPORTANT: retry=0 car c'est la première tentative
		execState = exec_state{
			Retry:     0,
			Timestamp: now,
		}
		// Ne pas écrire maintenant, on le fera dans l'étape 4
	} else {
		// La clé existe, récupérer sa valeur
		err = json.Unmarshal(getTakenResp.Kvs[0].Value, &execState)
		if err != nil {
			return "", fmt.Errorf("failed to unmarshal exec_state: %w", err)
		}
	}

	// Étape 3: Vérifier si retry > max_retry
	if execState.Retry > max_retry {
		// Transaction pour déplacer vers /error et nettoyer
		txn := client.Txn(ctx)
		txn = txn.Then(
			clientv3.OpPut(errorKey, value), // Créer /error/$position
			clientv3.OpDelete(stateKey),     // Supprimer /$state/$position
			clientv3.OpDelete(takenKey),     // Supprimer /taken/$position
		)

		_, err = txn.Commit()
		if err != nil {
			return "", fmt.Errorf("failed to move to error state: %w", err)
		}

		return "", fmt.Errorf("tuple at position %d exceeded max retries (%d), moved to error state", position, max_retry)
	}

	// Étape 4: Vérifier le timeout
	if now < execState.Timestamp+int64(timeout.Seconds()) {
		// Pas encore de timeout, retourner une erreur
		remainingTime := time.Duration(execState.Timestamp+int64(timeout.Seconds())-now) * time.Second
		return "", fmt.Errorf("tuple at position %d is still locked for %v (retry %d/%d)",
			position, remainingTime, execState.Retry, max_retry)
	}

	// Le timeout est dépassé, mettre à jour la clé /taken/$position
	execState.Retry++
	execState.Timestamp = now

	// Sérialiser la structure mise à jour
	execStateJSON, err := json.Marshal(execState)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated exec_state: %w", err)
	}

	// Mettre à jour la clé /taken/$position
	_, err = client.Put(ctx, takenKey, string(execStateJSON))
	if err != nil {
		return "", fmt.Errorf("failed to update taken key: %w", err)
	}

	// Étape 5: Appeler call_job pour traiter la valeur
	newValue, err := call_job(position, value)
	if err != nil {
		return "", fmt.Errorf("job execution failed for position %d: %w", position, err)
	}

	return newValue, nil
}

// release_tuple libère un tuple en supprimant sa clé /taken/$position
// Cette fonction utilitaire peut être appelée quand le traitement est terminé avec succès
func release_tuple(client *clientv3.Client, state string, position int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	takenKey := fmt.Sprintf("/taken/%d", position)
	stateKey := fmt.Sprintf("/%s/%d", state, position)

	// Transaction pour supprimer le tuple de l'état et libérer le verrou
	txn := client.Txn(ctx)
	txn = txn.Then(
		clientv3.OpDelete(stateKey), // Supprimer le tuple de l'état
		clientv3.OpDelete(takenKey), // Libérer le verrou
	)

	_, err := txn.Commit()
	if err != nil {
		return fmt.Errorf("failed to release tuple: %w", err)
	}

	return nil
}

// get_taken_info récupère les informations d'un tuple en cours de traitement
func get_taken_info(client *clientv3.Client, position int64) (*exec_state, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	takenKey := fmt.Sprintf("/taken/%d", position)

	getTakenResp, err := client.Get(ctx, takenKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get taken key: %w", err)
	}

	if len(getTakenResp.Kvs) == 0 {
		return nil, fmt.Errorf("no taken info for position %d", position)
	}

	var execState exec_state
	err = json.Unmarshal(getTakenResp.Kvs[0].Value, &execState)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal exec_state: %w", err)
	}

	return &execState, nil
}
