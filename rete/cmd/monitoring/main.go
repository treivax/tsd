package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/treivax/tsd/rete"
)

// MonitoringExample démontre l'utilisation du serveur de monitoring RETE
func main() {
	log.Println("🚀 Démarrage du serveur de monitoring RETE...")

	// Configuration du monitoring
	config := rete.MonitoringConfig{
		Port:            8080,
		UpdateInterval:  2 * time.Second,
		MaxHistorySize:  100,
		EnableProfiling: true,
		EnableAlerts:    true,
		LogLevel:        "info",
		MaxConnections:  100,
		ReadTimeout:     30 * time.Second,
		WriteTimeout:    30 * time.Second,
	}

	// Créer un réseau RETE de base pour le monitoring
	// En production, ceci serait votre réseau RETE existant
	reteNetwork := createSampleReteNetwork()

	// Créer le serveur de monitoring
	monitoringServer := rete.NewMonitoringServer(config, reteNetwork)

	// Configurer des règles d'alerte par défaut
	setupDefaultAlertRules(monitoringServer)

	// Context pour l'arrêt gracieux
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Gérer les signaux d'arrêt
	go handleShutdown(cancel)

	// Démarrer la simulation de données en arrière-plan
	go simulateReteActivity(ctx, reteNetwork)

	log.Printf("🌐 Interface web disponible sur: http://localhost:%d", config.Port)
	log.Printf("📊 API métriques: http://localhost:%d/api/metrics", config.Port)
	log.Printf("🔌 WebSocket: ws://localhost:%d/ws/metrics", config.Port)
	log.Println("💡 Utilisez Ctrl+C pour arrêter le serveur")

	// Démarrer le serveur de monitoring
	if err := monitoringServer.Start(ctx); err != nil {
		log.Fatalf("❌ Erreur lors du démarrage du serveur: %v", err)
	}
}

// createSampleReteNetwork crée un réseau RETE de démonstration
func createSampleReteNetwork() *rete.ReteNetwork {
	// En production, vous utiliseriez votre réseau RETE existant
	// Ici, nous créons un réseau minimal pour la démonstration
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)

	log.Println("🏗️ Réseau RETE de démonstration créé")
	return network
}

// setupDefaultAlertRules configure des règles d'alerte par défaut
func setupDefaultAlertRules(server *rete.MonitoringServer) {
	defaultRules := []rete.AlertRule{
		{
			ID:          "high_memory_usage",
			Name:        "High Memory Usage",
			Description: "Alert when memory usage exceeds 80%",
			Condition:   "memory_usage_percent > 80",
			Threshold:   80.0,
			Duration:    5 * time.Minute,
			Severity:    "high",
			IsEnabled:   true,
			Metadata: map[string]interface{}{
				"category": "system",
				"action":   "scale_up",
			},
		},
		{
			ID:          "high_error_rate",
			Name:        "High Error Rate",
			Description: "Alert when error rate exceeds 5%",
			Condition:   "error_rate > 5",
			Threshold:   5.0,
			Duration:    2 * time.Minute,
			Severity:    "critical",
			IsEnabled:   true,
			Metadata: map[string]interface{}{
				"category": "performance",
				"action":   "investigate",
			},
		},
		{
			ID:          "low_throughput",
			Name:        "Low Throughput",
			Description: "Alert when facts/sec drops below 100",
			Condition:   "facts_per_second < 100",
			Threshold:   100.0,
			Duration:    3 * time.Minute,
			Severity:    "medium",
			IsEnabled:   true,
			Metadata: map[string]interface{}{
				"category": "performance",
				"action":   "optimize",
			},
		},
		{
			ID:          "cache_miss_high",
			Name:        "High Cache Miss Rate",
			Description: "Alert when cache hit ratio drops below 70%",
			Condition:   "cache_hit_ratio < 70",
			Threshold:   70.0,
			Duration:    5 * time.Minute,
			Severity:    "medium",
			IsEnabled:   true,
			Metadata: map[string]interface{}{
				"category": "cache",
				"action":   "tune_cache",
			},
		},
	}

	// En production, vous ajouteriez ces règles via l'API du serveur
	log.Printf("📋 %d règles d'alerte par défaut configurées", len(defaultRules))
}

// simulateReteActivity simule l'activité du réseau RETE pour la démonstration
func simulateReteActivity(ctx context.Context, network *rete.ReteNetwork) {
	log.Println("🎭 Démarrage de la simulation d'activité RETE...")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	factCounter := 0

	for {
		select {
		case <-ctx.Done():
			log.Println("🎭 Simulation d'activité arrêtée")
			return
		case <-ticker.C:
			// Simuler l'ajout de faits
			for i := 0; i < 5; i++ {
				fact := &rete.Fact{
					ID:        generateFactID(factCounter),
					Type:      selectRandomType(),
					Fields:    generateRandomFields(),
					Timestamp: time.Now(),
				}

				// En mode démonstration, nous n'ajoutons pas vraiment les faits
				// network.SubmitFact(fact)
				_ = fact
				factCounter++
			}

			// Log périodique d'activité
			if factCounter%50 == 0 {
				log.Printf("📊 Simulation: %d faits traités", factCounter)
			}
		}
	}
}

// Fonctions utilitaires pour la simulation
func generateFactID(counter int) string {
	return fmt.Sprintf("fact_%d_%d", time.Now().Unix(), counter)
}

func selectRandomType() string {
	types := []string{"Person", "Order", "Product", "Invoice", "Customer", "Event"}
	return types[time.Now().UnixNano()%int64(len(types))]
}

func generateRandomFields() map[string]interface{} {
	fields := []map[string]interface{}{
		{
			"id":     time.Now().UnixNano() % 10000,
			"name":   "Sample",
			"active": true,
			"score":  float64(time.Now().UnixNano()%100) + 0.5,
		},
		{
			"user_id":   time.Now().UnixNano() % 1000,
			"action":    "click",
			"timestamp": time.Now().Unix(),
			"duration":  time.Now().UnixNano() % 5000,
		},
		{
			"order_id": time.Now().UnixNano() % 50000,
			"amount":   float64(time.Now().UnixNano()%10000) / 100.0,
			"currency": "EUR",
			"status":   "pending",
		},
	}

	return fields[time.Now().UnixNano()%int64(len(fields))]
}

// handleShutdown gère l'arrêt gracieux du serveur
func handleShutdown(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("🛑 Signal reçu: %v. Arrêt du serveur...", sig)
	cancel()
}

// Fonctions de démonstration pour tester les API

// DemonstrateAPIs montre comment utiliser les APIs du serveur
func DemonstrateAPIs() {
	log.Println("📡 Démonstration des APIs du serveur de monitoring...")

	// Ces exemples montrent comment vous pourriez interagir avec le serveur
	// En pratique, ces appels seraient faits depuis votre application

	examples := []string{
		"GET /api/metrics - Toutes les métriques",
		"GET /api/metrics/system - Métriques système uniquement",
		"GET /api/metrics/rete - Métriques RETE uniquement",
		"GET /api/metrics/performance - Métriques de performance",
		"GET /api/metrics/history - Historique des métriques",
		"GET /api/alerts - État des alertes",
		"POST /api/alerts/rules - Créer une règle d'alerte",
		"GET /api/network/status - État du réseau RETE",
		"GET /api/network/nodes - Liste des nœuds",
		"WS /ws/metrics - Flux temps réel via WebSocket",
	}

	log.Println("📋 APIs disponibles:")
	for _, example := range examples {
		log.Printf("   - %s", example)
	}
}

// ShowMonitoringFeatures présente les fonctionnalités du monitoring
func ShowMonitoringFeatures() {
	log.Println("🎯 Fonctionnalités du système de monitoring RETE:")
	log.Println("")
	log.Println("📊 MÉTRIQUES EN TEMPS RÉEL:")
	log.Println("   • Throughput (facts/sec, tokens/sec, rules/sec)")
	log.Println("   • Latence (P50, P75, P90, P95, P99)")
	log.Println("   • Usage mémoire et ressources système")
	log.Println("   • Statistiques des nœuds RETE")
	log.Println("")
	log.Println("⚡ OPTIMISATIONS DE PERFORMANCE:")
	log.Println("   • IndexedStorage: Cache hit ratio, temps de lookup")
	log.Println("   • HashJoinEngine: Statistiques de jointures, cache")
	log.Println("   • EvaluationCache: Hit/miss ratio, évictions")
	log.Println("   • TokenPropagation: Efficacité parallèle, utilisation workers")
	log.Println("")
	log.Println("🚨 SYSTÈME D'ALERTES:")
	log.Println("   • Règles d'alerte configurables")
	log.Println("   • Seuils et conditions personnalisés")
	log.Println("   • Notifications temps réel")
	log.Println("   • Historique des alertes")
	log.Println("")
	log.Println("🌐 INTERFACE WEB:")
	log.Println("   • Dashboard interactif responsive")
	log.Println("   • Graphiques temps réel avec Chart.js")
	log.Println("   • Visualisation de la topologie réseau")
	log.Println("   • Gestion des alertes et règles")
	log.Println("")
	log.Println("🔌 INTÉGRATION:")
	log.Println("   • API REST complète")
	log.Println("   • WebSocket pour mises à jour temps réel")
	log.Println("   • Export des données de monitoring")
	log.Println("   • Compatible avec systèmes de monitoring existants")
}

func init() {
	// Afficher les fonctionnalités au démarrage
	ShowMonitoringFeatures()
	DemonstrateAPIs()
}
