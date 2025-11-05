package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// new_tuple writes a new value to etcd at the next available position (max+1).
// This simplified approach always adds tuples at the end and increments the max counter.
// No holes management is needed with this strategy.
//
// Parameters:
//   - client: an etcd v3 client instance.
//   - state: a string representing the logical group or namespace for the tuple.
//   - value: the value to be stored.
//
// Returns:
//   - int64: the position (index) where the value was written.
//   - error: any error encountered during the operation.
func new_tuple(client *clientv3.Client, state, value string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	maxKey := fmt.Sprintf("/%s/max", state)

	// Get current max value
	getMaxResp, err := client.Get(ctx, maxKey)
	if err != nil {
		return 0, fmt.Errorf("failed to get max key: %w", err)
	}
	var max int64 = 0
	if len(getMaxResp.Kvs) > 0 {
		fmt.Sscanf(string(getMaxResp.Kvs[0].Value), "%d", &max)
	}

	// New position is always max + 1
	position := max + 1
	tupleKey := fmt.Sprintf("/%s/%d", state, position)

	// Atomic transaction: update max and create tuple
	txn := client.Txn(ctx)
	txn = txn.Then(clientv3.OpPut(maxKey, fmt.Sprintf("%d", position)))
	txn = txn.Then(clientv3.OpPut(tupleKey, value))

	_, err = txn.Commit()
	if err != nil {
		return 0, fmt.Errorf("transaction failed: %w", err)
	}
	return position, nil
}

// La fonction vérifie d'abord si la clé existe à la position spécifiée. Si c'est le cas, elle met à jour la valeur.
// Si la clé n'existe pas, une erreur est retournée.
//
// Paramètres :
//   - client : une instance du client etcd v3.
//   - state : une chaîne représentant le groupe logique ou l'espace de noms.
//   - pos : la position (index) du tuple à mettre à jour.
//   - value : la nouvelle valeur à stocker.
//
// Retourne :
//   - error : toute erreur rencontrée lors de l'opération.
func update_tuple(client *clientv3.Client, state string, pos int64, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tupleKey := fmt.Sprintf("/%s/%d", state, pos)

	// Check if the key exists
	getResp, err := client.Get(ctx, tupleKey)
	if err != nil {
		return fmt.Errorf("failed to get tuple key: %w", err)
	}
	if len(getResp.Kvs) == 0 {
		return fmt.Errorf("tuple at position %d does not exist", pos)
	}

	// Update the value
	_, err = client.Put(ctx, tupleKey, value)
	if err != nil {
		return fmt.Errorf("failed to update tuple: %w", err)
	}
	return nil
}

// delete_tuple supprime une valeur à une position donnée dans etcd pour un état (state) spécifique.
// Nouvelle stratégie : si la position supprimée n'est pas la dernière, on déplace le dernier tuple
// à cette position pour éviter les trous. Ensuite, on décrémente toujours max.
// L'opération est atomique grâce à une transaction etcd.
//
// Paramètres :
//   - client : une instance du client etcd v3.
//   - state : une chaîne représentant le groupe logique ou l'espace de noms.
//   - position : la position (index) du tuple à supprimer.
//
// Retourne :
//   - error : toute erreur rencontrée lors de l'opération.
func delete_tuple(client *clientv3.Client, state string, position int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	maxKey := fmt.Sprintf("/%s/max", state)
	tupleKey := fmt.Sprintf("/%s/%d", state, position)

	// Get current max value
	getMaxResp, err := client.Get(ctx, maxKey)
	if err != nil {
		return fmt.Errorf("failed to get max key: %w", err)
	}
	if len(getMaxResp.Kvs) == 0 {
		return fmt.Errorf("max key does not exist")
	}
	var max int64
	fmt.Sscanf(string(getMaxResp.Kvs[0].Value), "%d", &max)

	if max == 0 {
		return fmt.Errorf("no tuples to delete")
	}

	// Check if tuple exists at position
	getTupleResp, err := client.Get(ctx, tupleKey)
	if err != nil {
		return fmt.Errorf("failed to get tuple key: %w", err)
	}
	if len(getTupleResp.Kvs) == 0 {
		return fmt.Errorf("no value at position %d", position)
	}

	txn := client.Txn(ctx)

	if position == max {
		// If deleting the last position, just delete it and decrement max
		txn = txn.Then(clientv3.OpDelete(tupleKey))
		txn = txn.Then(clientv3.OpPut(maxKey, fmt.Sprintf("%d", max-1)))
	} else {
		// If not the last position, move the last tuple to the deleted position
		lastTupleKey := fmt.Sprintf("/%s/%d", state, max)

		// Get the last tuple value
		getLastResp, err := client.Get(ctx, lastTupleKey)
		if err != nil {
			return fmt.Errorf("failed to get last tuple: %w", err)
		}
		if len(getLastResp.Kvs) == 0 {
			return fmt.Errorf("last tuple not found at position %d", max)
		}
		lastValue := string(getLastResp.Kvs[0].Value)

		// Move last tuple to deleted position, delete last position, decrement max
		txn = txn.Then(clientv3.OpPut(tupleKey, lastValue))
		txn = txn.Then(clientv3.OpDelete(lastTupleKey))
		txn = txn.Then(clientv3.OpPut(maxKey, fmt.Sprintf("%d", max-1)))
	}

	_, err = txn.Commit()
	if err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}

// one_tuple lit un tuple au hasard pour un état donné.
// Avec la nouvelle stratégie sans trous, on choisit simplement un nombre aléatoire entre 1 et max.
// Retourne la position et la valeur, ou une erreur si aucun tuple n'est trouvé.
func one_tuple(client *clientv3.Client, state string) (int64, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	maxKey := fmt.Sprintf("/%s/max", state)

	// Récupérer la valeur max
	getMaxResp, err := client.Get(ctx, maxKey)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get max key: %w", err)
	}
	if len(getMaxResp.Kvs) == 0 {
		return 0, "", fmt.Errorf("no tuples for state %s", state)
	}
	var max int64
	fmt.Sscanf(string(getMaxResp.Kvs[0].Value), "%d", &max)
	if max == 0 {
		return 0, "", fmt.Errorf("no tuples for state %s", state)
	}

	// Choisir une position au hasard entre 1 et max (plus de trous à gérer)
	pos := rand.Int63n(max) + 1
	tupleKey := fmt.Sprintf("/%s/%d", state, pos)

	// Récupérer le tuple à cette position
	getTupleResp, err := client.Get(ctx, tupleKey)
	if err != nil {
		return 0, "", fmt.Errorf("failed to get tuple: %w", err)
	}
	if len(getTupleResp.Kvs) == 0 {
		return 0, "", fmt.Errorf("tuple at position %d not found", pos)
	}
	return pos, string(getTupleResp.Kvs[0].Value), nil
}

// ===================== VERSIONS RAW POUR TRANSACTIONS COMPOSÉES =====================

// new_tuple_raw prépare les opérations pour ajouter un tuple sans exécuter la transaction.
// Utilisée pour composer des transactions plus complexes.
func new_tuple_raw(ctx context.Context, client *clientv3.Client, state, value string) ([]clientv3.Op, int64, error) {
	maxKey := fmt.Sprintf("/%s/max", state)

	// Get current max value
	getMaxResp, err := client.Get(ctx, maxKey)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get max key: %w", err)
	}
	var max int64 = 0
	if len(getMaxResp.Kvs) > 0 {
		fmt.Sscanf(string(getMaxResp.Kvs[0].Value), "%d", &max)
	}

	// New position is always max + 1
	position := max + 1
	tupleKey := fmt.Sprintf("/%s/%d", state, position)

	// Return operations without executing
	ops := []clientv3.Op{
		clientv3.OpPut(maxKey, fmt.Sprintf("%d", position)),
		clientv3.OpPut(tupleKey, value),
	}

	return ops, position, nil
}

// delete_tuple_raw prépare les opérations pour supprimer un tuple sans exécuter la transaction.
// Utilisée pour composer des transactions plus complexes.
func delete_tuple_raw(ctx context.Context, client *clientv3.Client, state string, position int64) ([]clientv3.Op, error) {
	maxKey := fmt.Sprintf("/%s/max", state)
	tupleKey := fmt.Sprintf("/%s/%d", state, position)

	// Get current max value
	getMaxResp, err := client.Get(ctx, maxKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get max key: %w", err)
	}
	if len(getMaxResp.Kvs) == 0 {
		return nil, fmt.Errorf("max key does not exist")
	}
	var max int64
	fmt.Sscanf(string(getMaxResp.Kvs[0].Value), "%d", &max)

	if max == 0 {
		return nil, fmt.Errorf("no tuples to delete")
	}

	// Check if tuple exists at position
	getTupleResp, err := client.Get(ctx, tupleKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get tuple key: %w", err)
	}
	if len(getTupleResp.Kvs) == 0 {
		return nil, fmt.Errorf("no value at position %d", position)
	}

	var ops []clientv3.Op

	if position == max {
		// If deleting the last position, just delete it and decrement max
		ops = []clientv3.Op{
			clientv3.OpDelete(tupleKey),
			clientv3.OpPut(maxKey, fmt.Sprintf("%d", max-1)),
		}
	} else {
		// If not the last position, move the last tuple to the deleted position
		lastTupleKey := fmt.Sprintf("/%s/%d", state, max)

		// Get the last tuple value
		getLastResp, err := client.Get(ctx, lastTupleKey)
		if err != nil {
			return nil, fmt.Errorf("failed to get last tuple: %w", err)
		}
		if len(getLastResp.Kvs) == 0 {
			return nil, fmt.Errorf("last tuple not found at position %d", max)
		}
		lastValue := string(getLastResp.Kvs[0].Value)

		// Move last tuple to deleted position, delete last position, decrement max
		ops = []clientv3.Op{
			clientv3.OpPut(tupleKey, lastValue),
			clientv3.OpDelete(lastTupleKey),
			clientv3.OpPut(maxKey, fmt.Sprintf("%d", max-1)),
		}
	}

	return ops, nil
}

// update_tuple_raw prépare les opérations pour mettre à jour un tuple sans exécuter la transaction.
func update_tuple_raw(ctx context.Context, client *clientv3.Client, state string, pos int64, value string) ([]clientv3.Op, error) {
	tupleKey := fmt.Sprintf("/%s/%d", state, pos)

	// Check if the key exists
	getResp, err := client.Get(ctx, tupleKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get tuple key: %w", err)
	}
	if len(getResp.Kvs) == 0 {
		return nil, fmt.Errorf("tuple at position %d does not exist", pos)
	}

	// Return update operation
	ops := []clientv3.Op{
		clientv3.OpPut(tupleKey, value),
	}

	return ops, nil
}

// ===================== FONCTIONS COMPOSÉES ATOMIQUES =====================

// move_tuple déplace atomiquement un tuple d'un état vers un autre.
// Cette opération combine une lecture, une suppression et un ajout dans une seule transaction.
func move_tuple(client *clientv3.Client, fromState string, position int64, toState string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// D'abord, récupérer la valeur du tuple à déplacer
	fromTupleKey := fmt.Sprintf("/%s/%d", fromState, position)
	getTupleResp, err := client.Get(ctx, fromTupleKey)
	if err != nil {
		return 0, fmt.Errorf("failed to get tuple to move: %w", err)
	}
	if len(getTupleResp.Kvs) == 0 {
		return 0, fmt.Errorf("no tuple at position %d in state %s", position, fromState)
	}
	value := string(getTupleResp.Kvs[0].Value)

	// Préparer les opérations de suppression depuis fromState
	deleteOps, err := delete_tuple_raw(ctx, client, fromState, position)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare delete operations: %w", err)
	}

	// Préparer les opérations d'ajout vers toState
	addOps, newPosition, err := new_tuple_raw(ctx, client, toState, value)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare add operations: %w", err)
	}

	// Combiner toutes les opérations dans une seule transaction
	var allOps []clientv3.Op
	allOps = append(allOps, deleteOps...)
	allOps = append(allOps, addOps...)

	// Exécuter la transaction atomique
	txn := client.Txn(ctx)
	txn = txn.Then(allOps...)

	_, err = txn.Commit()
	if err != nil {
		return 0, fmt.Errorf("move transaction failed: %w", err)
	}

	return newPosition, nil
}

// swap_tuples échange atomiquement deux tuples entre deux états.
func swap_tuples(client *clientv3.Client, state1 string, pos1 int64, state2 string, pos2 int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Récupérer les deux valeurs
	tuple1Key := fmt.Sprintf("/%s/%d", state1, pos1)
	tuple2Key := fmt.Sprintf("/%s/%d", state2, pos2)

	get1Resp, err := client.Get(ctx, tuple1Key)
	if err != nil {
		return fmt.Errorf("failed to get first tuple: %w", err)
	}
	if len(get1Resp.Kvs) == 0 {
		return fmt.Errorf("no tuple at position %d in state %s", pos1, state1)
	}

	get2Resp, err := client.Get(ctx, tuple2Key)
	if err != nil {
		return fmt.Errorf("failed to get second tuple: %w", err)
	}
	if len(get2Resp.Kvs) == 0 {
		return fmt.Errorf("no tuple at position %d in state %s", pos2, state2)
	}

	value1 := string(get1Resp.Kvs[0].Value)
	value2 := string(get2Resp.Kvs[0].Value)

	// Créer une transaction atomique pour échanger les valeurs
	txn := client.Txn(ctx)
	txn = txn.Then(
		clientv3.OpPut(tuple1Key, value2),
		clientv3.OpPut(tuple2Key, value1),
	)

	_, err = txn.Commit()
	if err != nil {
		return fmt.Errorf("swap transaction failed: %w", err)
	}

	return nil
}

// batch_operations exécute plusieurs opérations sur les tuples de manière atomique.
// Exemple d'utilisation pour des opérations complexes personnalisées.
func batch_operations(client *clientv3.Client, operations []func(context.Context, *clientv3.Client) ([]clientv3.Op, error)) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var allOps []clientv3.Op

	// Préparer toutes les opérations
	for _, op := range operations {
		ops, err := op(ctx, client)
		if err != nil {
			return fmt.Errorf("failed to prepare operation: %w", err)
		}
		allOps = append(allOps, ops...)
	}

	// Exécuter toutes les opérations dans une seule transaction
	txn := client.Txn(ctx)
	txn = txn.Then(allOps...)

	_, err := txn.Commit()
	if err != nil {
		return fmt.Errorf("batch transaction failed: %w", err)
	}

	return nil
}

// ===================== APPROCHE HYBRIDE POUR LOGIQUE CONDITIONNELLE =====================

// ConditionalOperation représente une opération qui peut être conditionnelle
type ConditionalOperation struct {
	Condition func(map[string]interface{}) bool // Condition basée sur des résultats précédents
	Operation func(context.Context, *clientv3.Client) ([]clientv3.Op, error)
	ResultKey string // Clé pour stocker le résultat dans le contexte
}

// ExecuteConditionalBatch exécute une série d'opérations avec logique conditionnelle
// en utilisant plusieurs transactions si nécessaire
func ExecuteConditionalBatch(client *clientv3.Client, operations []ConditionalOperation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	results := make(map[string]interface{})

	for _, op := range operations {
		// Vérifier la condition si elle existe
		if op.Condition != nil && !op.Condition(results) {
			continue // Passer cette opération
		}

		// Préparer les opérations
		ops, err := op.Operation(ctx, client)
		if err != nil {
			return fmt.Errorf("failed to prepare operation: %w", err)
		}

		// Exécuter la transaction pour cette étape
		txn := client.Txn(ctx)
		txn = txn.Then(ops...)

		resp, err := txn.Commit()
		if err != nil {
			return fmt.Errorf("operation failed: %w", err)
		}

		// Stocker les résultats si nécessaire
		if op.ResultKey != "" {
			results[op.ResultKey] = resp.Succeeded
		}
	}

	return nil
}

// WorkflowStep représente une étape dans un workflow complexe
type WorkflowStep struct {
	Name     string
	Execute  func(ctx context.Context, client *clientv3.Client, state *WorkflowState) error
	Rollback func(ctx context.Context, client *clientv3.Client, state *WorkflowState) error
}

// WorkflowState maintient l'état entre les étapes
type WorkflowState struct {
	Data           map[string]interface{}
	CompletedSteps []string
}

// ExecuteWorkflow exécute un workflow avec possibilité de rollback
func ExecuteWorkflow(client *clientv3.Client, steps []WorkflowStep) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	state := &WorkflowState{
		Data:           make(map[string]interface{}),
		CompletedSteps: make([]string, 0),
	}

	// Exécuter les étapes
	for i, step := range steps {
		err := step.Execute(ctx, client, state)
		if err != nil {
			// Rollback des étapes déjà exécutées
			fmt.Printf("Step %s failed, rolling back...\n", step.Name)
			for j := i - 1; j >= 0; j-- {
				if steps[j].Rollback != nil {
					rollbackErr := steps[j].Rollback(ctx, client, state)
					if rollbackErr != nil {
						fmt.Printf("Rollback of step %s failed: %v\n", steps[j].Name, rollbackErr)
					}
				}
			}
			return fmt.Errorf("workflow failed at step %s: %w", step.Name, err)
		}
		state.CompletedSteps = append(state.CompletedSteps, step.Name)
	}

	return nil
}
