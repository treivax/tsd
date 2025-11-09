package rete

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// MonitoringServer gère le serveur de monitoring en temps réel
type MonitoringServer struct {
	server          *http.Server
	router          *mux.Router
	metrics         *MetricsCollector
	clients         map[*websocket.Conn]bool
	clientsMutex    sync.RWMutex
	upgrader        websocket.Upgrader
	reteNetwork     *ReteNetwork
	performanceData *PerformanceData
	alertSystem     *AlertSystem
}

// MetricsCollector collecte les métriques système et RETE
type MetricsCollector struct {
	sync.RWMutex

	// Métriques système
	SystemMetrics SystemMetrics `json:"system_metrics"`

	// Métriques RETE
	ReteMetrics ReteMetrics `json:"rete_metrics"`

	// Métriques de performance
	PerformanceMetrics PerformanceMetrics `json:"performance_metrics"`

	// Historique des métriques
	MetricsHistory []HistoricalMetrics `json:"metrics_history"`

	// Configuration
	MaxHistorySize int           `json:"max_history_size"`
	UpdateInterval time.Duration `json:"update_interval"`

	// État de collecte
	IsCollecting bool      `json:"is_collecting"`
	StartTime    time.Time `json:"start_time"`
}

// SystemMetrics représente les métriques système
type SystemMetrics struct {
	Timestamp       time.Time `json:"timestamp"`
	MemoryUsage     uint64    `json:"memory_usage_bytes"`
	MemoryAllocated uint64    `json:"memory_allocated_bytes"`
	MemorySystem    uint64    `json:"memory_system_bytes"`
	GoroutineCount  int       `json:"goroutine_count"`
	GCCount         uint32    `json:"gc_count"`
	CPUUsage        float64   `json:"cpu_usage_percent"`
	UptimeSeconds   int64     `json:"uptime_seconds"`
}

// ReteMetrics représente les métriques spécifiques au réseau RETE
type ReteMetrics struct {
	Timestamp       time.Time `json:"timestamp"`
	TotalNodes      int       `json:"total_nodes"`
	ActiveNodes     int       `json:"active_nodes"`
	TotalFacts      int64     `json:"total_facts"`
	FactsPerSecond  float64   `json:"facts_per_second"`
	TokensProcessed int64     `json:"tokens_processed"`
	TokensPerSecond float64   `json:"tokens_per_second"`
	RulesTriggered  int64     `json:"rules_triggered"`
	RulesPerSecond  float64   `json:"rules_per_second"`
	AverageLatency  float64   `json:"average_latency_ms"`
	ErrorCount      int64     `json:"error_count"`
	ErrorRate       float64   `json:"error_rate_percent"`

	// Métriques par type de nœud
	NodeTypeMetrics map[string]NodeMetrics `json:"node_type_metrics"`
}

// NodeMetrics représente les métriques d'un type de nœud
type NodeMetrics struct {
	Count               int     `json:"count"`
	ProcessingTime      float64 `json:"processing_time_ms"`
	ThroughputPerSecond float64 `json:"throughput_per_second"`
	ErrorCount          int64   `json:"error_count"`
	CacheHitRatio       float64 `json:"cache_hit_ratio"`
}

// PerformanceMetrics représente les métriques de performance
type PerformanceMetrics struct {
	Timestamp             time.Time              `json:"timestamp"`
	IndexedStorageStats   map[string]interface{} `json:"indexed_storage_stats"`
	HashJoinStats         map[string]interface{} `json:"hash_join_stats"`
	EvaluationCacheStats  map[string]interface{} `json:"evaluation_cache_stats"`
	TokenPropagationStats map[string]interface{} `json:"token_propagation_stats"`
	OverallThroughput     float64                `json:"overall_throughput_ops_per_sec"`
	ResponseTime95th      time.Duration          `json:"response_time_95th_percentile"`
	ResponseTime99th      time.Duration          `json:"response_time_99th_percentile"`
}

// HistoricalMetrics contient un snapshot historique des métriques
type HistoricalMetrics struct {
	Timestamp          time.Time          `json:"timestamp"`
	SystemMetrics      SystemMetrics      `json:"system_metrics"`
	ReteMetrics        ReteMetrics        `json:"rete_metrics"`
	PerformanceMetrics PerformanceMetrics `json:"performance_metrics"`
}

// PerformanceData gère les données de performance en temps réel
type PerformanceData struct {
	sync.RWMutex

	FactProcessingTimes   []time.Duration    `json:"fact_processing_times"`
	TokenPropagationTimes []time.Duration    `json:"token_propagation_times"`
	RuleExecutionTimes    []time.Duration    `json:"rule_execution_times"`
	ThroughputSamples     []ThroughputSample `json:"throughput_samples"`
	LatencySamples        []LatencySample    `json:"latency_samples"`
	ErrorEvents           []ErrorEvent       `json:"error_events"`

	MaxSamplesSize int `json:"max_samples_size"`
}

// ThroughputSample représente un échantillon de débit
type ThroughputSample struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Type      string    `json:"type"` // "facts", "tokens", "rules"
}

// LatencySample représente un échantillon de latence
type LatencySample struct {
	Timestamp time.Time     `json:"timestamp"`
	Duration  time.Duration `json:"duration"`
	Operation string        `json:"operation"`
}

// ErrorEvent représente un événement d'erreur
type ErrorEvent struct {
	Timestamp time.Time              `json:"timestamp"`
	Error     string                 `json:"error"`
	Component string                 `json:"component"`
	Severity  string                 `json:"severity"` // "low", "medium", "high", "critical"
	Context   map[string]interface{} `json:"context"`
}

// AlertSystem gère les alertes et notifications
type AlertSystem struct {
	sync.RWMutex

	Rules            []AlertRule `json:"rules"`
	ActiveAlerts     []Alert     `json:"active_alerts"`
	AlertHistory     []Alert     `json:"alert_history"`
	NotificationChan chan Alert  `json:"-"`
	IsEnabled        bool        `json:"is_enabled"`
	MaxHistorySize   int         `json:"max_history_size"`
}

// AlertRule définit une règle d'alerte
type AlertRule struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Condition   string                 `json:"condition"` // Expression logique
	Threshold   float64                `json:"threshold"`
	Duration    time.Duration          `json:"duration"` // Durée avant déclenchement
	Severity    string                 `json:"severity"`
	IsEnabled   bool                   `json:"is_enabled"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Alert représente une alerte active ou historique
type Alert struct {
	ID         string                 `json:"id"`
	RuleID     string                 `json:"rule_id"`
	RuleName   string                 `json:"rule_name"`
	Timestamp  time.Time              `json:"timestamp"`
	Severity   string                 `json:"severity"`
	Message    string                 `json:"message"`
	Value      float64                `json:"value"`
	Threshold  float64                `json:"threshold"`
	IsActive   bool                   `json:"is_active"`
	ResolvedAt *time.Time             `json:"resolved_at,omitempty"`
	Context    map[string]interface{} `json:"context"`
}

// MonitoringConfig représente la configuration du monitoring
type MonitoringConfig struct {
	Port            int           `json:"port"`
	UpdateInterval  time.Duration `json:"update_interval"`
	MaxHistorySize  int           `json:"max_history_size"`
	EnableProfiling bool          `json:"enable_profiling"`
	EnableAlerts    bool          `json:"enable_alerts"`
	LogLevel        string        `json:"log_level"`
	MaxConnections  int           `json:"max_connections"`
	ReadTimeout     time.Duration `json:"read_timeout"`
	WriteTimeout    time.Duration `json:"write_timeout"`
}

// NewMonitoringServer crée un nouveau serveur de monitoring
func NewMonitoringServer(config MonitoringConfig, reteNetwork *ReteNetwork) *MonitoringServer {
	router := mux.NewRouter()

	metrics := &MetricsCollector{
		MaxHistorySize: config.MaxHistorySize,
		UpdateInterval: config.UpdateInterval,
		StartTime:      time.Now(),
		MetricsHistory: make([]HistoricalMetrics, 0, config.MaxHistorySize),
	}

	performanceData := &PerformanceData{
		MaxSamplesSize:        1000,
		FactProcessingTimes:   make([]time.Duration, 0, 1000),
		TokenPropagationTimes: make([]time.Duration, 0, 1000),
		RuleExecutionTimes:    make([]time.Duration, 0, 1000),
		ThroughputSamples:     make([]ThroughputSample, 0, 1000),
		LatencySamples:        make([]LatencySample, 0, 1000),
		ErrorEvents:           make([]ErrorEvent, 0, 1000),
	}

	alertSystem := &AlertSystem{
		Rules:            make([]AlertRule, 0),
		ActiveAlerts:     make([]Alert, 0),
		AlertHistory:     make([]Alert, 0, config.MaxHistorySize),
		NotificationChan: make(chan Alert, 100),
		IsEnabled:        config.EnableAlerts,
		MaxHistorySize:   config.MaxHistorySize,
	}

	server := &MonitoringServer{
		router:  router,
		metrics: metrics,
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Permettre toutes les origines en développement
			},
		},
		reteNetwork:     reteNetwork,
		performanceData: performanceData,
		alertSystem:     alertSystem,
	}

	// Configurer les routes
	server.setupRoutes()

	// Créer le serveur HTTP
	server.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      router,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}

	return server
}

// setupRoutes configure les routes du serveur
func (ms *MonitoringServer) setupRoutes() {
	// API REST pour les métriques
	ms.router.HandleFunc("/api/metrics", ms.handleMetrics).Methods("GET")
	ms.router.HandleFunc("/api/metrics/system", ms.handleSystemMetrics).Methods("GET")
	ms.router.HandleFunc("/api/metrics/rete", ms.handleReteMetrics).Methods("GET")
	ms.router.HandleFunc("/api/metrics/performance", ms.handlePerformanceMetrics).Methods("GET")
	ms.router.HandleFunc("/api/metrics/history", ms.handleMetricsHistory).Methods("GET")

	// API pour les alertes
	ms.router.HandleFunc("/api/alerts", ms.handleAlerts).Methods("GET")
	ms.router.HandleFunc("/api/alerts/rules", ms.handleAlertRules).Methods("GET", "POST")
	ms.router.HandleFunc("/api/alerts/rules/{id}", ms.handleAlertRule).Methods("GET", "PUT", "DELETE")

	// API pour les données de performance en temps réel
	ms.router.HandleFunc("/api/performance/live", ms.handleLivePerformance).Methods("GET")
	ms.router.HandleFunc("/api/network/status", ms.handleNetworkStatus).Methods("GET")
	ms.router.HandleFunc("/api/network/nodes", ms.handleNetworkNodes).Methods("GET")

	// WebSocket pour les mises à jour en temps réel
	ms.router.HandleFunc("/ws/metrics", ms.handleWebSocketMetrics)

	// Interface web statique - chemin relatif au répertoire de travail
	ms.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./rete/assets/web/"))).Methods("GET")
}

// Start démarre le serveur de monitoring
func (ms *MonitoringServer) Start(ctx context.Context) error {
	// Démarrer la collecte de métriques
	go ms.startMetricsCollection(ctx)

	// Démarrer le système d'alertes
	if ms.alertSystem.IsEnabled {
		go ms.startAlertProcessing(ctx)
	}

	log.Printf("🚀 Serveur de monitoring RETE démarré sur le port %s", ms.server.Addr)
	log.Printf("📊 Interface web disponible sur http://localhost%s", ms.server.Addr)
	log.Printf("🔌 WebSocket endpoint: ws://localhost%s/ws/metrics", ms.server.Addr)

	// Démarrer le serveur HTTP
	if err := ms.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("failed to start monitoring server: %w", err)
	}

	return nil
}

// Stop arrête le serveur de monitoring
func (ms *MonitoringServer) Stop(ctx context.Context) error {
	log.Printf("🛑 Arrêt du serveur de monitoring...")

	// Fermer toutes les connexions WebSocket
	ms.clientsMutex.Lock()
	for client := range ms.clients {
		client.Close()
	}
	ms.clients = make(map[*websocket.Conn]bool)
	ms.clientsMutex.Unlock()

	// Arrêter le serveur HTTP
	return ms.server.Shutdown(ctx)
}

// startMetricsCollection démarre la collecte périodique de métriques
func (ms *MonitoringServer) startMetricsCollection(ctx context.Context) {
	ticker := time.NewTicker(ms.metrics.UpdateInterval)
	defer ticker.Stop()

	ms.metrics.IsCollecting = true
	log.Printf("📈 Collecte de métriques démarrée (intervalle: %v)", ms.metrics.UpdateInterval)

	for {
		select {
		case <-ctx.Done():
			ms.metrics.IsCollecting = false
			log.Printf("📈 Collecte de métriques arrêtée")
			return
		case <-ticker.C:
			ms.collectMetrics()
			ms.broadcastMetrics()
		}
	}
}

// collectMetrics collecte toutes les métriques système et RETE
func (ms *MonitoringServer) collectMetrics() {
	ms.metrics.Lock()
	defer ms.metrics.Unlock()

	now := time.Now()

	// Collecter les métriques système
	ms.collectSystemMetrics(now)

	// Collecter les métriques RETE
	ms.collectReteMetrics(now)

	// Collecter les métriques de performance
	ms.collectPerformanceMetrics(now)

	// Ajouter à l'historique
	historical := HistoricalMetrics{
		Timestamp:          now,
		SystemMetrics:      ms.metrics.SystemMetrics,
		ReteMetrics:        ms.metrics.ReteMetrics,
		PerformanceMetrics: ms.metrics.PerformanceMetrics,
	}

	ms.metrics.MetricsHistory = append(ms.metrics.MetricsHistory, historical)

	// Limiter la taille de l'historique
	if len(ms.metrics.MetricsHistory) > ms.metrics.MaxHistorySize {
		ms.metrics.MetricsHistory = ms.metrics.MetricsHistory[1:]
	}
}

// collectSystemMetrics collecte les métriques système
func (ms *MonitoringServer) collectSystemMetrics(timestamp time.Time) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	ms.metrics.SystemMetrics = SystemMetrics{
		Timestamp:       timestamp,
		MemoryUsage:     memStats.Alloc,
		MemoryAllocated: memStats.TotalAlloc,
		MemorySystem:    memStats.Sys,
		GoroutineCount:  runtime.NumGoroutine(),
		GCCount:         memStats.NumGC,
		CPUUsage:        ms.getCPUUsage(), // À implémenter
		UptimeSeconds:   int64(time.Since(ms.metrics.StartTime).Seconds()),
	}
}

// collectReteMetrics collecte les métriques du réseau RETE
func (ms *MonitoringServer) collectReteMetrics(timestamp time.Time) {
	// Cette méthode doit être étendue selon l'API disponible du ReteNetwork
	ms.metrics.ReteMetrics = ReteMetrics{
		Timestamp:       timestamp,
		TotalNodes:      ms.getReteNodeCount(),
		ActiveNodes:     ms.getActiveNodeCount(),
		TotalFacts:      ms.getTotalFactCount(),
		FactsPerSecond:  ms.calculateFactsPerSecond(),
		TokensProcessed: ms.getTokensProcessed(),
		TokensPerSecond: ms.calculateTokensPerSecond(),
		RulesTriggered:  ms.getRulesTriggered(),
		RulesPerSecond:  ms.calculateRulesPerSecond(),
		AverageLatency:  ms.calculateAverageLatency(),
		ErrorCount:      ms.getErrorCount(),
		ErrorRate:       ms.calculateErrorRate(),
		NodeTypeMetrics: ms.getNodeTypeMetrics(),
	}
}

// collectPerformanceMetrics collecte les métriques de performance des composants optimisés
func (ms *MonitoringServer) collectPerformanceMetrics(timestamp time.Time) {
	ms.metrics.PerformanceMetrics = PerformanceMetrics{
		Timestamp:             timestamp,
		IndexedStorageStats:   ms.getIndexedStorageStats(),
		HashJoinStats:         ms.getHashJoinStats(),
		EvaluationCacheStats:  ms.getEvaluationCacheStats(),
		TokenPropagationStats: ms.getTokenPropagationStats(),
		OverallThroughput:     ms.calculateOverallThroughput(),
		ResponseTime95th:      ms.calculateResponseTime95th(),
		ResponseTime99th:      ms.calculateResponseTime99th(),
	}
}

// Méthodes utilitaires pour calculer les métriques (à implémenter selon l'API RETE)
func (ms *MonitoringServer) getCPUUsage() float64 {
	// Implémentation simplifiée - dans un vrai système, utiliser une lib comme gopsutil
	return 0.0
}

func (ms *MonitoringServer) getReteNodeCount() int {
	// À implémenter selon l'API ReteNetwork
	return 0
}

func (ms *MonitoringServer) getActiveNodeCount() int {
	return 0
}

func (ms *MonitoringServer) getTotalFactCount() int64 {
	return 0
}

func (ms *MonitoringServer) calculateFactsPerSecond() float64 {
	return 0.0
}

func (ms *MonitoringServer) getTokensProcessed() int64 {
	return 0
}

func (ms *MonitoringServer) calculateTokensPerSecond() float64 {
	return 0.0
}

func (ms *MonitoringServer) getRulesTriggered() int64 {
	return 0
}

func (ms *MonitoringServer) calculateRulesPerSecond() float64 {
	return 0.0
}

func (ms *MonitoringServer) calculateAverageLatency() float64 {
	return 0.0
}

func (ms *MonitoringServer) getErrorCount() int64 {
	return 0
}

func (ms *MonitoringServer) calculateErrorRate() float64 {
	return 0.0
}

func (ms *MonitoringServer) getNodeTypeMetrics() map[string]NodeMetrics {
	return make(map[string]NodeMetrics)
}

func (ms *MonitoringServer) getIndexedStorageStats() map[string]interface{} {
	return make(map[string]interface{})
}

func (ms *MonitoringServer) getHashJoinStats() map[string]interface{} {
	return make(map[string]interface{})
}

func (ms *MonitoringServer) getEvaluationCacheStats() map[string]interface{} {
	return make(map[string]interface{})
}

func (ms *MonitoringServer) getTokenPropagationStats() map[string]interface{} {
	return make(map[string]interface{})
}

func (ms *MonitoringServer) calculateOverallThroughput() float64 {
	return 0.0
}

func (ms *MonitoringServer) calculateResponseTime95th() time.Duration {
	return 0
}

func (ms *MonitoringServer) calculateResponseTime99th() time.Duration {
	return 0
}

// broadcastMetrics diffuse les métriques à tous les clients WebSocket connectés
func (ms *MonitoringServer) broadcastMetrics() {
	ms.clientsMutex.RLock()
	defer ms.clientsMutex.RUnlock()

	if len(ms.clients) == 0 {
		return
	}

	ms.metrics.RLock()
	metricsJSON, err := json.Marshal(map[string]interface{}{
		"type": "metrics_update",
		"data": map[string]interface{}{
			"system":      ms.metrics.SystemMetrics,
			"rete":        ms.metrics.ReteMetrics,
			"performance": ms.metrics.PerformanceMetrics,
			"timestamp":   time.Now(),
		},
	})
	ms.metrics.RUnlock()

	if err != nil {
		log.Printf("❌ Erreur lors de la sérialisation des métriques: %v", err)
		return
	}

	// Diffuser à tous les clients connectés
	for client := range ms.clients {
		if err := client.WriteMessage(websocket.TextMessage, metricsJSON); err != nil {
			log.Printf("❌ Erreur lors de l'envoi des métriques WebSocket: %v", err)
			client.Close()
			delete(ms.clients, client)
		}
	}
}

// Handlers HTTP pour l'API REST
func (ms *MonitoringServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	ms.metrics.RLock()
	defer ms.metrics.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(ms.metrics); err != nil {
		http.Error(w, "Failed to encode metrics", http.StatusInternalServerError)
		return
	}
}

func (ms *MonitoringServer) handleSystemMetrics(w http.ResponseWriter, r *http.Request) {
	ms.metrics.RLock()
	defer ms.metrics.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(ms.metrics.SystemMetrics); err != nil {
		http.Error(w, "Failed to encode system metrics", http.StatusInternalServerError)
		return
	}
}

func (ms *MonitoringServer) handleReteMetrics(w http.ResponseWriter, r *http.Request) {
	ms.metrics.RLock()
	defer ms.metrics.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(ms.metrics.ReteMetrics); err != nil {
		http.Error(w, "Failed to encode RETE metrics", http.StatusInternalServerError)
		return
	}
}

func (ms *MonitoringServer) handlePerformanceMetrics(w http.ResponseWriter, r *http.Request) {
	ms.metrics.RLock()
	defer ms.metrics.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(ms.metrics.PerformanceMetrics); err != nil {
		http.Error(w, "Failed to encode performance metrics", http.StatusInternalServerError)
		return
	}
}

func (ms *MonitoringServer) handleMetricsHistory(w http.ResponseWriter, r *http.Request) {
	ms.metrics.RLock()
	defer ms.metrics.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(ms.metrics.MetricsHistory); err != nil {
		http.Error(w, "Failed to encode metrics history", http.StatusInternalServerError)
		return
	}
}

func (ms *MonitoringServer) handleAlerts(w http.ResponseWriter, r *http.Request) {
	ms.alertSystem.RLock()
	defer ms.alertSystem.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response := map[string]interface{}{
		"active_alerts": ms.alertSystem.ActiveAlerts,
		"alert_history": ms.alertSystem.AlertHistory,
		"rules_count":   len(ms.alertSystem.Rules),
		"is_enabled":    ms.alertSystem.IsEnabled,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode alerts", http.StatusInternalServerError)
		return
	}
}

func (ms *MonitoringServer) handleAlertRules(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	switch r.Method {
	case "GET":
		ms.alertSystem.RLock()
		rules := ms.alertSystem.Rules
		ms.alertSystem.RUnlock()

		if err := json.NewEncoder(w).Encode(rules); err != nil {
			http.Error(w, "Failed to encode alert rules", http.StatusInternalServerError)
		}

	case "POST":
		var rule AlertRule
		if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ms.alertSystem.Lock()
		rule.ID = fmt.Sprintf("rule_%d", time.Now().Unix())
		ms.alertSystem.Rules = append(ms.alertSystem.Rules, rule)
		ms.alertSystem.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(rule)
	}
}

func (ms *MonitoringServer) handleAlertRule(w http.ResponseWriter, r *http.Request) {
	// Implémentation des opérations CRUD sur les règles d'alerte
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	vars := mux.Vars(r)
	ruleID := vars["id"]

	// TODO: Implémenter les opérations CRUD
	_ = ruleID // Utiliser la variable pour éviter l'erreur

	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{"message": "Not implemented yet"})
}

func (ms *MonitoringServer) handleLivePerformance(w http.ResponseWriter, r *http.Request) {
	ms.performanceData.RLock()
	defer ms.performanceData.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(ms.performanceData); err != nil {
		http.Error(w, "Failed to encode performance data", http.StatusInternalServerError)
		return
	}
}

func (ms *MonitoringServer) handleNetworkStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// À implémenter selon l'API ReteNetwork
	status := map[string]interface{}{
		"status":    "running",
		"timestamp": time.Now(),
		"uptime":    time.Since(ms.metrics.StartTime).String(),
		"version":   "1.0.0",
	}

	if err := json.NewEncoder(w).Encode(status); err != nil {
		http.Error(w, "Failed to encode network status", http.StatusInternalServerError)
		return
	}
}

func (ms *MonitoringServer) handleNetworkNodes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// À implémenter selon l'API ReteNetwork
	nodes := []map[string]interface{}{
		{
			"id":     "root_node",
			"type":   "RootNode",
			"status": "active",
			"facts":  0,
		},
	}

	if err := json.NewEncoder(w).Encode(nodes); err != nil {
		http.Error(w, "Failed to encode network nodes", http.StatusInternalServerError)
		return
	}
}

// handleWebSocketMetrics gère les connexions WebSocket pour les mises à jour en temps réel
func (ms *MonitoringServer) handleWebSocketMetrics(w http.ResponseWriter, r *http.Request) {
	conn, err := ms.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ Erreur lors de l'upgrade WebSocket: %v", err)
		return
	}
	defer conn.Close()

	// Ajouter le client à la liste
	ms.clientsMutex.Lock()
	ms.clients[conn] = true
	clientCount := len(ms.clients)
	ms.clientsMutex.Unlock()

	log.Printf("🔌 Nouveau client WebSocket connecté (total: %d)", clientCount)

	// Envoyer les métriques actuelles immédiatement
	ms.metrics.RLock()
	initialData := map[string]interface{}{
		"type": "initial_data",
		"data": map[string]interface{}{
			"system":      ms.metrics.SystemMetrics,
			"rete":        ms.metrics.ReteMetrics,
			"performance": ms.metrics.PerformanceMetrics,
			"history":     ms.metrics.MetricsHistory,
		},
	}
	ms.metrics.RUnlock()

	if initialJSON, err := json.Marshal(initialData); err == nil {
		conn.WriteMessage(websocket.TextMessage, initialJSON)
	}

	// Garder la connexion ouverte et gérer les messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("🔌 Client WebSocket déconnecté: %v", err)
			break
		}
	}

	// Supprimer le client de la liste
	ms.clientsMutex.Lock()
	delete(ms.clients, conn)
	clientCount = len(ms.clients)
	ms.clientsMutex.Unlock()

	log.Printf("🔌 Client WebSocket supprimé (total: %d)", clientCount)
}

// startAlertProcessing démarre le traitement des alertes
func (ms *MonitoringServer) startAlertProcessing(ctx context.Context) {
	log.Printf("🚨 Système d'alertes démarré")

	for {
		select {
		case <-ctx.Done():
			log.Printf("🚨 Système d'alertes arrêté")
			return
		case alert := <-ms.alertSystem.NotificationChan:
			ms.processAlert(alert)
		default:
			// Évaluer les règles d'alerte périodiquement
			ms.evaluateAlertRules()
			time.Sleep(5 * time.Second)
		}
	}
}

// evaluateAlertRules évalue toutes les règles d'alerte actives
func (ms *MonitoringServer) evaluateAlertRules() {
	ms.alertSystem.RLock()
	rules := ms.alertSystem.Rules
	ms.alertSystem.RUnlock()

	for _, rule := range rules {
		if !rule.IsEnabled {
			continue
		}

		if ms.evaluateRule(rule) {
			alert := Alert{
				ID:        fmt.Sprintf("alert_%d", time.Now().Unix()),
				RuleID:    rule.ID,
				RuleName:  rule.Name,
				Timestamp: time.Now(),
				Severity:  rule.Severity,
				Message:   fmt.Sprintf("Rule '%s' triggered", rule.Name),
				Threshold: rule.Threshold,
				IsActive:  true,
				Context:   rule.Metadata,
			}

			ms.alertSystem.NotificationChan <- alert
		}
	}
}

// evaluateRule évalue une règle d'alerte spécifique
func (ms *MonitoringServer) evaluateRule(rule AlertRule) bool {
	// Implémentation simplifiée - dans un vrai système, implémenter un évaluateur d'expressions
	// Pour l'exemple, on retourne false
	return false
}

// processAlert traite une alerte
func (ms *MonitoringServer) processAlert(alert Alert) {
	ms.alertSystem.Lock()
	defer ms.alertSystem.Unlock()

	// Ajouter à la liste des alertes actives
	ms.alertSystem.ActiveAlerts = append(ms.alertSystem.ActiveAlerts, alert)

	// Ajouter à l'historique
	ms.alertSystem.AlertHistory = append(ms.alertSystem.AlertHistory, alert)

	// Limiter la taille de l'historique
	if len(ms.alertSystem.AlertHistory) > ms.alertSystem.MaxHistorySize {
		ms.alertSystem.AlertHistory = ms.alertSystem.AlertHistory[1:]
	}

	log.Printf("🚨 ALERTE [%s]: %s", alert.Severity, alert.Message)

	// Diffuser l'alerte via WebSocket
	alertJSON, _ := json.Marshal(map[string]interface{}{
		"type": "alert",
		"data": alert,
	})

	ms.clientsMutex.RLock()
	for client := range ms.clients {
		client.WriteMessage(websocket.TextMessage, alertJSON)
	}
	ms.clientsMutex.RUnlock()
}
