package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/treivax/tsd/rete"
)

// MonitoringExample démontre l'utilisation du serveur de monitoring RETE
func main() {
	log.Println("🚀 Démarrage du serveur de monitoring RETE...")

	// Configuration simplifiée et stable
	config := rete.MonitoringConfig{
		Port:            8080,
		UpdateInterval:  3 * time.Second, // Plus lent pour éviter la surcharge
		MaxHistorySize:  50,              // Plus petit pour les tests
		EnableProfiling: false,           // Désactivé pour simplifier
		EnableAlerts:    false,           // Désactivé pour éviter les complexités
		LogLevel:        "info",
		MaxConnections:  20,
		ReadTimeout:     15 * time.Second,
		WriteTimeout:    15 * time.Second,
	}

	// Créer un réseau RETE simple pour la démonstration
	storage := rete.NewMemoryStorage()
	reteNetwork := rete.NewReteNetwork(storage)

	// Créer le serveur de monitoring
	monitoringServer := rete.NewMonitoringServer(config, reteNetwork)

	log.Println("🏗️ Réseau RETE et serveur de monitoring créés")

	// Context pour l'arrêt gracieux
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Gérer les signaux d'arrêt proprement
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-c
		log.Printf("🛑 Signal reçu: %v. Arrêt du serveur...", sig)
		cancel()

		// Donner du temps pour l'arrêt gracieux
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := monitoringServer.Stop(shutdownCtx); err != nil {
			log.Printf("❌ Erreur lors de l'arrêt: %v", err)
		} else {
			log.Println("✅ Serveur arrêté proprement")
		}

		os.Exit(0)
	}()

	log.Printf("🌐 Interface web disponible sur: http://localhost:%d", config.Port)
	log.Printf("📊 API métriques: http://localhost:%d/api/metrics", config.Port)
	log.Printf("🔌 WebSocket: ws://localhost:%d/ws/metrics", config.Port)
	log.Println("💡 Utilisez Ctrl+C pour arrêter le serveur")

	// Démarrer le serveur de monitoring
	if err := monitoringServer.Start(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatalf("❌ Erreur lors du démarrage du serveur: %v", err)
	}

	log.Println("🎉 Serveur de monitoring terminé")
}
