package rete

import (
	"container/heap"
	"fmt"
	"sync"
	"time"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// TokenPropagationEngine optimise la propagation des tokens dans le réseau RETE
type TokenPropagationEngine struct {
	// Queue de priorité pour les tokens
	priorityQueue *TokenPriorityQueue

	// Pool de workers pour la propagation parallèle
	workerPool *TokenWorkerPool

	// Buffer pour le batching des tokens
	batchBuffer *TokenBatchBuffer

	// Configuration de propagation
	config PropagationConfig

	// Statistiques de performance
	stats PropagationStats

	// Verrou pour la concurrence
	mutex sync.RWMutex
}

// PropagationConfig configure les paramètres de propagation
type PropagationConfig struct {
	// Nombre de workers pour la propagation parallèle
	NumWorkers int

	// Taille du buffer de batching
	BatchSize int

	// Timeout pour le batching
	BatchTimeout time.Duration

	// Activer la propagation par priorité
	EnablePrioritization bool

	// Facteur de priorité basé sur le temps
	TimePriorityFactor float64

	// Facteur de priorité basé sur la complexité
	ComplexityPriorityFactor float64

	// Taille maximale de la queue
	MaxQueueSize int
}

// PropagationStats contient les statistiques de propagation
type PropagationStats struct {
	TokensProcessed    int64
	BatchesProcessed   int64
	AverageBatchSize   float64
	AverageProcTime    time.Duration
	ParallelEfficiency float64
	QueueOverflows     int64
	WorkerUtilization  []float64
}

// TokenPriorityItem représente un token avec sa priorité
type TokenPriorityItem struct {
	Token    *domain.Token
	Priority float64
	NodeID   string
	Created  time.Time
	Index    int // Index dans la heap
}

// TokenPriorityQueue implémente une queue de priorité pour les tokens
type TokenPriorityQueue struct {
	items []*TokenPriorityItem
	mutex sync.Mutex
}

// TokenWorkerPool gère un pool de workers pour la propagation
type TokenWorkerPool struct {
	workers    []*TokenWorker
	workChan   chan *TokenPriorityItem
	numWorkers int
	stats      []WorkerStats
}

// TokenWorker représente un worker individuel
type TokenWorker struct {
	id       int
	workChan chan *TokenPriorityItem
	engine   *TokenPropagationEngine
	stats    WorkerStats
}

// WorkerStats contient les statistiques d'un worker
type WorkerStats struct {
	TokensProcessed int64
	TotalTime       time.Duration
	IdleTime        time.Duration
	ActiveTime      time.Duration
}

// TokenBatchBuffer buffer les tokens pour le traitement par batch
type TokenBatchBuffer struct {
	buffer    []*TokenPriorityItem
	maxSize   int
	timeout   time.Duration
	lastFlush time.Time
	mutex     sync.Mutex
	flushChan chan []*TokenPriorityItem
}

// NewTokenPropagationEngine crée un nouveau moteur de propagation optimisé
func NewTokenPropagationEngine(config PropagationConfig) *TokenPropagationEngine {
	if config.NumWorkers <= 0 {
		config.NumWorkers = 4
	}
	if config.BatchSize <= 0 {
		config.BatchSize = 100
	}
	if config.BatchTimeout <= 0 {
		config.BatchTimeout = 10 * time.Millisecond
	}
	if config.MaxQueueSize <= 0 {
		config.MaxQueueSize = 10000
	}

	engine := &TokenPropagationEngine{
		priorityQueue: NewTokenPriorityQueue(),
		config:        config,
		stats:         PropagationStats{WorkerUtilization: make([]float64, config.NumWorkers)},
	}

	// Initialiser le pool de workers
	engine.workerPool = NewTokenWorkerPool(config.NumWorkers, engine)

	// Initialiser le buffer de batching
	engine.batchBuffer = NewTokenBatchBuffer(config.BatchSize, config.BatchTimeout)

	return engine
}

// NewTokenPriorityQueue crée une nouvelle queue de priorité pour tokens
func NewTokenPriorityQueue() *TokenPriorityQueue {
	pq := &TokenPriorityQueue{
		items: make([]*TokenPriorityItem, 0),
	}
	heap.Init(pq)
	return pq
}

// NewTokenWorkerPool crée un nouveau pool de workers
func NewTokenWorkerPool(numWorkers int, engine *TokenPropagationEngine) *TokenWorkerPool {
	pool := &TokenWorkerPool{
		workers:    make([]*TokenWorker, numWorkers),
		workChan:   make(chan *TokenPriorityItem, numWorkers*2),
		numWorkers: numWorkers,
		stats:      make([]WorkerStats, numWorkers),
	}

	// Créer et démarrer les workers
	for i := 0; i < numWorkers; i++ {
		worker := &TokenWorker{
			id:       i,
			workChan: pool.workChan,
			engine:   engine,
		}
		pool.workers[i] = worker
		go worker.Start()
	}

	return pool
}

// NewTokenBatchBuffer crée un nouveau buffer de batching
func NewTokenBatchBuffer(maxSize int, timeout time.Duration) *TokenBatchBuffer {
	buffer := &TokenBatchBuffer{
		buffer:    make([]*TokenPriorityItem, 0, maxSize),
		maxSize:   maxSize,
		timeout:   timeout,
		lastFlush: time.Now(),
		flushChan: make(chan []*TokenPriorityItem, 10),
	}

	// Démarrer le timer de flush automatique
	go buffer.autoFlush()

	return buffer
}

// EnqueueToken ajoute un token à la queue de propagation
func (tpe *TokenPropagationEngine) EnqueueToken(
	token *domain.Token,
	nodeID string,
	priority float64,
) error {
	tpe.mutex.Lock()
	defer tpe.mutex.Unlock()

	// Vérifier si la queue n'est pas pleine
	if tpe.priorityQueue.Len() >= tpe.config.MaxQueueSize {
		tpe.stats.QueueOverflows++
		return fmt.Errorf("queue overflow: maximum size %d reached", tpe.config.MaxQueueSize)
	}

	// Calculer la priorité si la prioritisation est activée
	if tpe.config.EnablePrioritization {
		priority = tpe.calculatePriority(token, priority)
	}

	// Créer l'item de priorité
	item := &TokenPriorityItem{
		Token:    token,
		Priority: priority,
		NodeID:   nodeID,
		Created:  time.Now(),
	}

	// Ajouter à la queue de priorité
	heap.Push(tpe.priorityQueue, item)

	return nil
}

// ProcessTokens traite les tokens en utilisant la stratégie optimisée
func (tpe *TokenPropagationEngine) ProcessTokens() error {
	tpe.mutex.RLock()
	config := tpe.config
	tpe.mutex.RUnlock()

	if config.BatchSize > 1 {
		// Traitement par batch
		return tpe.processBatched()
	} else {
		// Traitement individuel optimisé
		return tpe.processIndividual()
	}
}

// processBatched traite les tokens par batch pour une meilleure efficacité
func (tpe *TokenPropagationEngine) processBatched() error {
	for {
		// Récupérer un batch de tokens
		batch := tpe.getBatch()
		if len(batch) == 0 {
			time.Sleep(time.Millisecond) // Attendre de nouveaux tokens
			continue
		}

		// Traiter le batch en parallèle
		tpe.processBatch(batch)

		// Mettre à jour les statistiques
		tpe.updateBatchStats(batch)
	}
}

// processIndividual traite les tokens individuellement avec optimisations
func (tpe *TokenPropagationEngine) processIndividual() error {
	for {
		// Récupérer le token de plus haute priorité
		item := tpe.dequeueHighestPriority()
		if item == nil {
			time.Sleep(time.Millisecond) // Attendre de nouveaux tokens
			continue
		}

		// Envoyer au pool de workers
		select {
		case tpe.workerPool.workChan <- item:
			// Token envoyé au worker
		default:
			// Workers occupés, traiter directement
			tpe.processTokenItem(item)
		}
	}
}

// getBatch récupère un batch de tokens optimisé
func (tpe *TokenPropagationEngine) getBatch() []*TokenPriorityItem {
	tpe.mutex.Lock()
	defer tpe.mutex.Unlock()

	batchSize := tpe.config.BatchSize
	if tpe.priorityQueue.Len() < batchSize {
		batchSize = tpe.priorityQueue.Len()
	}

	if batchSize == 0 {
		return nil
	}

	batch := make([]*TokenPriorityItem, batchSize)
	for i := 0; i < batchSize; i++ {
		if tpe.priorityQueue.Len() > 0 {
			batch[i] = heap.Pop(tpe.priorityQueue).(*TokenPriorityItem)
		}
	}

	return batch
}

// processBatch traite un batch de tokens en parallèle
func (tpe *TokenPropagationEngine) processBatch(batch []*TokenPriorityItem) {
	// Diviser le batch entre les workers disponibles
	numWorkers := tpe.config.NumWorkers
	chunkSize := len(batch) / numWorkers
	if chunkSize == 0 {
		chunkSize = 1
	}

	var wg sync.WaitGroup

	for i := 0; i < len(batch); i += chunkSize {
		end := i + chunkSize
		if end > len(batch) {
			end = len(batch)
		}

		wg.Add(1)
		go func(chunk []*TokenPriorityItem) {
			defer wg.Done()
			for _, item := range chunk {
				tpe.processTokenItem(item)
			}
		}(batch[i:end])
	}

	wg.Wait()
}

// dequeueHighestPriority récupère le token de plus haute priorité
func (tpe *TokenPropagationEngine) dequeueHighestPriority() *TokenPriorityItem {
	tpe.mutex.Lock()
	defer tpe.mutex.Unlock()

	if tpe.priorityQueue.Len() == 0 {
		return nil
	}

	return heap.Pop(tpe.priorityQueue).(*TokenPriorityItem)
}

// processTokenItem traite un item de token individuel
func (tpe *TokenPropagationEngine) processTokenItem(item *TokenPriorityItem) {
	startTime := time.Now()

	// Ici, on appellerait la logique de propagation réelle du nœud
	// Pour cet exemple, on simule le traitement

	// Simuler le temps de traitement basé sur la complexité du token
	processingTime := tpe.estimateProcessingTime(item.Token)
	time.Sleep(processingTime)

	// Mettre à jour les statistiques
	tpe.updateProcessingStats(time.Since(startTime))
}

// calculatePriority calcule la priorité d'un token
func (tpe *TokenPropagationEngine) calculatePriority(token *domain.Token, basePriority float64) float64 {
	priority := basePriority

	// Facteur basé sur l'âge du token (utiliser l'âge du plus ancien fait)
	var oldestTime time.Time
	for _, fact := range token.Facts {
		if oldestTime.IsZero() || fact.Timestamp.Before(oldestTime) {
			oldestTime = fact.Timestamp
		}
	}
	if !oldestTime.IsZero() {
		age := time.Since(oldestTime)
		timeFactor := float64(age.Nanoseconds()) * tpe.config.TimePriorityFactor
		priority += timeFactor
	}

	// Facteur basé sur la complexité (nombre de faits dans le token)
	complexityFactor := float64(len(token.Facts)) * tpe.config.ComplexityPriorityFactor
	priority += complexityFactor

	// Facteur basé sur le type de nœud (les nœuds terminaux ont une priorité plus élevée)
	if token.NodeID != "" {
		// Logique spécifique au type de nœud
		priority *= 1.1 // Exemple
	}

	return priority
}

// estimateProcessingTime estime le temps de traitement d'un token
func (tpe *TokenPropagationEngine) estimateProcessingTime(token *domain.Token) time.Duration {
	// Base sur la complexité du token
	baseTime := 100 * time.Microsecond
	complexityMultiplier := len(token.Facts)

	return time.Duration(int64(baseTime) * int64(complexityMultiplier))
}

// updateProcessingStats met à jour les statistiques de traitement
func (tpe *TokenPropagationEngine) updateProcessingStats(duration time.Duration) {
	tpe.mutex.Lock()
	defer tpe.mutex.Unlock()

	tpe.stats.TokensProcessed++

	// Mise à jour de la moyenne du temps de traitement
	totalTokens := tpe.stats.TokensProcessed
	if totalTokens == 1 {
		tpe.stats.AverageProcTime = duration
	} else {
		tpe.stats.AverageProcTime = time.Duration(
			(int64(tpe.stats.AverageProcTime)*int64(totalTokens-1) + int64(duration)) / int64(totalTokens),
		)
	}
}

// updateBatchStats met à jour les statistiques de batching
func (tpe *TokenPropagationEngine) updateBatchStats(batch []*TokenPriorityItem) {
	tpe.mutex.Lock()
	defer tpe.mutex.Unlock()

	tpe.stats.BatchesProcessed++
	batchSize := float64(len(batch))

	// Mise à jour de la taille moyenne des batches
	totalBatches := float64(tpe.stats.BatchesProcessed)
	if totalBatches == 1 {
		tpe.stats.AverageBatchSize = batchSize
	} else {
		tpe.stats.AverageBatchSize = (tpe.stats.AverageBatchSize*(totalBatches-1) + batchSize) / totalBatches
	}
}

// Start démarre un worker
func (tw *TokenWorker) Start() {
	for item := range tw.workChan {
		startTime := time.Now()

		// Traiter le token
		tw.engine.processTokenItem(item)

		// Mettre à jour les statistiques du worker
		processingTime := time.Since(startTime)
		tw.stats.TokensProcessed++
		tw.stats.ActiveTime += processingTime
		tw.stats.TotalTime += processingTime
	}
}

// autoFlush flush automatiquement le buffer selon le timeout
func (tbb *TokenBatchBuffer) autoFlush() {
	ticker := time.NewTicker(tbb.timeout)
	defer ticker.Stop()

	for range ticker.C {
		tbb.mutex.Lock()
		if len(tbb.buffer) > 0 && time.Since(tbb.lastFlush) >= tbb.timeout {
			// Copier le buffer pour flush
			batch := make([]*TokenPriorityItem, len(tbb.buffer))
			copy(batch, tbb.buffer)

			// Vider le buffer
			tbb.buffer = tbb.buffer[:0]
			tbb.lastFlush = time.Now()

			// Envoyer le batch
			select {
			case tbb.flushChan <- batch:
			default:
				// Canal plein, ignorer ce flush
			}
		}
		tbb.mutex.Unlock()
	}
}

// GetStats retourne les statistiques de propagation
func (tpe *TokenPropagationEngine) GetStats() PropagationStats {
	tpe.mutex.RLock()
	defer tpe.mutex.RUnlock()

	// Calculer l'utilisation des workers
	for i, worker := range tpe.workerPool.workers {
		if worker.stats.TotalTime > 0 {
			tpe.stats.WorkerUtilization[i] = float64(worker.stats.ActiveTime) / float64(worker.stats.TotalTime)
		}
	}

	return tpe.stats
}

// Implémentation de l'interface heap.Interface pour TokenPriorityQueue

func (pq *TokenPriorityQueue) Len() int {
	return len(pq.items)
}

func (pq *TokenPriorityQueue) Less(i, j int) bool {
	// Plus haute priorité en premier
	return pq.items[i].Priority > pq.items[j].Priority
}

func (pq *TokenPriorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].Index = i
	pq.items[j].Index = j
}

func (pq *TokenPriorityQueue) Push(x interface{}) {
	item := x.(*TokenPriorityItem)
	item.Index = len(pq.items)
	pq.items = append(pq.items, item)
}

func (pq *TokenPriorityQueue) Pop() interface{} {
	old := pq.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.Index = -1
	pq.items = old[0 : n-1]
	return item
}
