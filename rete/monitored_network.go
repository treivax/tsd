package rete

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/treivax/tsd/rete/pkg/domain"
)

// MonitoredRETENetwork est un wrapper qui ajoute le monitoring au réseau RETE existant
type MonitoredRETENetwork struct {
	*ReteNetwork

	// Composants de monitoring
	monitoringServer  *MonitoringServer
	metricsIntegrator *MetricsIntegrator

	// Configuration du monitoring
	monitoringConfig    *MonitoredNetworkConfig
	isMonitoringEnabled bool
}

// MonitoredNetworkConfig contient la configuration du monitoring
type MonitoredNetworkConfig struct {
	ServerPort         int              `json:"server_port"`
	MetricsInterval    time.Duration    `json:"metrics_interval"`
	EnableWebInterface bool             `json:"enable_web_interface"`
	EnableAlerts       bool             `json:"enable_alerts"`
	WebInterfaceDir    string           `json:"web_interface_dir"`
	MaxHistoryPoints   int              `json:"max_history_points"`
	AlertThresholds    *AlertThresholds `json:"alert_thresholds"`
}

// AlertThresholds définit les seuils d'alerte
type AlertThresholds struct {
	MaxLatencyMs        float64 `json:"max_latency_ms"`
	MaxErrorRate        float64 `json:"max_error_rate"`
	MinThroughput       float64 `json:"min_throughput"`
	MaxMemoryUsageBytes int64   `json:"max_memory_usage_bytes"`
	MinCacheHitRatio    float64 `json:"min_cache_hit_ratio"`
}

// DefaultMonitoredNetworkConfig retourne une configuration par défaut
func DefaultMonitoredNetworkConfig() *MonitoredNetworkConfig {
	return &MonitoredNetworkConfig{
		ServerPort:         8080,
		MetricsInterval:    5 * time.Second,
		EnableWebInterface: true,
		EnableAlerts:       true,
		WebInterfaceDir:    "./rete/web",
		MaxHistoryPoints:   1000,
		AlertThresholds: &AlertThresholds{
			MaxLatencyMs:        100.0,
			MaxErrorRate:        5.0,
			MinThroughput:       100.0,
			MaxMemoryUsageBytes: 500 * 1024 * 1024, // 500 MB
			MinCacheHitRatio:    70.0,
		},
	}
}

// NewMonitoredRETENetwork crée un nouveau réseau RETE avec monitoring
func NewMonitoredRETENetwork(storage Storage, config *MonitoredNetworkConfig) *MonitoredRETENetwork {
	if config == nil {
		config = DefaultMonitoredNetworkConfig()
	}

	// Créer le réseau RETE de base
	network := NewReteNetwork(storage)

	// Créer l'intégrateur de métriques
	metricsIntegrator := NewMetricsIntegrator(config.MetricsInterval)

	// Créer la configuration du serveur de monitoring
	monitoringConfig := MonitoringConfig{
		Port:            config.ServerPort,
		UpdateInterval:  config.MetricsInterval,
		MaxHistorySize:  config.MaxHistoryPoints,
		EnableProfiling: true,
		EnableAlerts:    config.EnableAlerts,
		LogLevel:        "info",
		MaxConnections:  100,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
	}

	// Créer le serveur de monitoring
	monitoringServer := NewMonitoringServer(monitoringConfig, network)

	// Créer le réseau monitoré
	monitoredNetwork := &MonitoredRETENetwork{
		ReteNetwork:         network,
		monitoringServer:    monitoringServer,
		metricsIntegrator:   metricsIntegrator,
		monitoringConfig:    config,
		isMonitoringEnabled: false,
	}

	// Connecter l'intégrateur au serveur de monitoring
	metricsIntegrator.RegisterUpdateCallback(func(metrics *AggregatedMetrics) {
		// Le serveur de monitoring collecte automatiquement les métriques
		// Pas besoin d'updateMetrics explicite
	})

	return monitoredNetwork
}

// StartMonitoring démarre le système de monitoring
func (mrn *MonitoredRETENetwork) StartMonitoring() error {
	if mrn.isMonitoringEnabled {
		return fmt.Errorf("monitoring is already enabled")
	}

	// Enregistrer les composants optimisés si disponibles
	mrn.registerOptimizedComponents()

	// Démarrer l'intégrateur de métriques
	mrn.metricsIntegrator.Start()

	// Démarrer le serveur de monitoring
	go func() {
		log.Printf("Starting monitoring server on port %d", mrn.monitoringConfig.ServerPort)
		ctx := context.Background()
		if err := mrn.monitoringServer.Start(ctx); err != nil && err != http.ErrServerClosed {
			log.Printf("Monitoring server error: %v", err)
		}
	}()

	mrn.isMonitoringEnabled = true
	log.Println("RETE monitoring system started successfully")

	return nil
}

// StopMonitoring arrête le système de monitoring
func (mrn *MonitoredRETENetwork) StopMonitoring() error {
	if !mrn.isMonitoringEnabled {
		return fmt.Errorf("monitoring is not enabled")
	}

	// Arrêter l'intégrateur de métriques
	mrn.metricsIntegrator.Stop()

	// Arrêter le serveur de monitoring
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := mrn.monitoringServer.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop monitoring server: %v", err)
	}

	mrn.isMonitoringEnabled = false
	log.Println("RETE monitoring system stopped")

	return nil
}

// registerOptimizedComponents enregistre les composants optimisés pour le monitoring
func (mrn *MonitoredRETENetwork) registerOptimizedComponents() {
	// Vérifier si les composants optimisés sont disponibles et les enregistrer
	// Note: Ceci dépend de l'implémentation spécifique de votre réseau RETE

	// Exemple d'enregistrement (à adapter selon votre architecture)
	/*
		if storage, ok := mrn.ReteNetwork.Storage.(*IndexedFactStorage); ok {
			if joinEngine, ok := mrn.ReteNetwork.GetJoinEngine().(*HashJoinEngine); ok {
				if evalCache, ok := mrn.ReteNetwork.GetEvaluationCache().(*EvaluationCache); ok {
					if tokenProp, ok := mrn.ReteNetwork.GetTokenPropagation().(*TokenPropagationEngine); ok {
						mrn.metricsIntegrator.RegisterComponents(storage, joinEngine, evalCache, tokenProp)
					}
				}
			}
		}
	*/

	log.Println("Optimized components registered for monitoring")
}

// AddFact ajoute un fait au réseau avec monitoring
func (mrn *MonitoredRETENetwork) AddFact(fact *domain.Fact) error {
	startTime := time.Now()

	// Convertir le fait domain vers le type Fact du réseau
	networkFact := &Fact{
		ID:   fact.ID,
		Type: fact.Type,
		// Adapter les champs selon le besoin
	}

	// Utiliser la mémoire de travail du noeud racine pour ajouter le fait
	mrn.ReteNetwork.RootNode.Memory.AddFact(networkFact)

	processingTime := time.Since(startTime)

	// Enregistrer les métriques si le monitoring est activé
	if mrn.isMonitoringEnabled {
		mrn.metricsIntegrator.RecordFactProcessed(fact, processingTime)
	}

	return nil
}

// ProcessToken traite un token avec monitoring
func (mrn *MonitoredRETENetwork) ProcessToken(token *domain.Token) error {
	startTime := time.Now()

	// Traiter le token (cette méthode devrait exister dans votre réseau RETE)
	// err := mrn.ReteNetwork.ProcessToken(token)
	err := mrn.processTokenWithNetwork(token)

	processingTime := time.Since(startTime)

	// Enregistrer les métriques si le monitoring est activé
	if mrn.isMonitoringEnabled {
		mrn.metricsIntegrator.RecordTokenProcessed(token, processingTime)

		if err != nil {
			mrn.metricsIntegrator.RecordError("token_processing", err)
		}
	}

	return err
}

// processTokenWithNetwork traite un token avec le réseau (méthode adaptée)
func (mrn *MonitoredRETENetwork) processTokenWithNetwork(token *domain.Token) error {
	// Implémentation adaptée selon votre architecture RETE
	// Cette méthode devrait déléguer au réseau RETE existant

	// Exemple simplifié - traitement basique du token
	return nil
}

// ExecuteRule exécute une règle avec monitoring
func (mrn *MonitoredRETENetwork) ExecuteRule(ruleName string) error {
	startTime := time.Now()

	// Exécuter la règle (cette méthode devrait exister dans votre réseau RETE)
	err := mrn.executeRuleWithNetwork(ruleName)

	executionTime := time.Since(startTime)

	// Enregistrer les métriques si le monitoring est activé
	if mrn.isMonitoringEnabled {
		mrn.metricsIntegrator.RecordRuleTriggered(ruleName, executionTime)

		if err != nil {
			mrn.metricsIntegrator.RecordError("rule_execution", err)
		}
	}

	return err
}

// executeRuleWithNetwork exécute une règle avec le réseau (méthode adaptée)
func (mrn *MonitoredRETENetwork) executeRuleWithNetwork(ruleName string) error {
	// Implémentation adaptée selon votre architecture RETE
	// Cette méthode devrait déléguer au réseau RETE existant

	// Exemple simplifié
	return nil
}

// GetMonitoringURL retourne l'URL de l'interface de monitoring
func (mrn *MonitoredRETENetwork) GetMonitoringURL() string {
	if !mrn.isMonitoringEnabled {
		return ""
	}

	return fmt.Sprintf("http://localhost:%d", mrn.monitoringConfig.ServerPort)
}

// GetCurrentMetrics retourne les métriques actuelles
func (mrn *MonitoredRETENetwork) GetCurrentMetrics() *AggregatedMetrics {
	if !mrn.isMonitoringEnabled {
		return nil
	}

	return mrn.metricsIntegrator.GetCurrentMetrics()
}

// IsMonitoringEnabled retourne true si le monitoring est activé
func (mrn *MonitoredRETENetwork) IsMonitoringEnabled() bool {
	return mrn.isMonitoringEnabled
}

// UpdateMonitoringConfig met à jour la configuration du monitoring
func (mrn *MonitoredRETENetwork) UpdateMonitoringConfig(config *MonitoredNetworkConfig) error {
	if mrn.isMonitoringEnabled {
		return fmt.Errorf("cannot update config while monitoring is enabled")
	}

	mrn.monitoringConfig = config
	return nil
}

// Example démonstration d'utilisation du réseau RETE monitoré
func ExampleMonitoredRETEUsage() {
	// Créer un stockage nil (sera géré par le réseau)
	var storage Storage = nil

	// Créer un réseau RETE monitoré avec configuration par défaut
	config := DefaultMonitoredNetworkConfig()
	config.ServerPort = 8080
	config.MetricsInterval = 2 * time.Second

	network := NewMonitoredRETENetwork(storage, config)

	// Démarrer le monitoring
	if err := network.StartMonitoring(); err != nil {
		log.Fatalf("Failed to start monitoring: %v", err)
	}

	fmt.Printf("Monitoring interface available at: %s\n", network.GetMonitoringURL())

	// Simuler l'utilisation du réseau RETE
	go func() {
		for i := 0; i < 100; i++ {
			fact := &domain.Fact{
				ID:   fmt.Sprintf("fact_%d", i),
				Type: "test",
				Fields: map[string]interface{}{
					"value": i,
					"name":  fmt.Sprintf("test_fact_%d", i),
				},
			}

			if err := network.AddFact(fact); err != nil {
				log.Printf("Error adding fact: %v", err)
			}

			time.Sleep(100 * time.Millisecond)
		}
	}()

	// Attendre quelques secondes pour voir les métriques
	time.Sleep(10 * time.Second)

	// Obtenir les métriques actuelles
	metrics := network.GetCurrentMetrics()
	if metrics != nil {
		fmt.Printf("Current metrics - Facts processed: %d, Average latency: %.2f ms\n",
			metrics.GlobalMetrics.TotalFactsProcessed,
			metrics.GlobalMetrics.AverageLatencyMs)
	}

	// Arrêter le monitoring
	if err := network.StopMonitoring(); err != nil {
		log.Printf("Error stopping monitoring: %v", err)
	}
}
