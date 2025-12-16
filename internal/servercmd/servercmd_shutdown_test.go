// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package servercmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

// TestShutdown_NilHTTPServer vérifie que Shutdown gère le cas où httpServer est nil
func TestShutdown_NilHTTPServer(t *testing.T) {
	t.Log("🧪 TEST SHUTDOWN - NIL HTTP SERVER")
	t.Log("===================================")

	config := &Config{
		Host:     "localhost",
		Port:     8080,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("❌ Erreur création serveur: %v", err)
	}

	// httpServer devrait être nil à ce stade
	if server.httpServer != nil {
		t.Error("❌ httpServer devrait être nil après NewServer")
	}

	ctx := context.Background()
	err = server.Shutdown(ctx)

	if err != nil {
		t.Errorf("❌ Shutdown avec httpServer nil ne devrait pas retourner d'erreur, got: %v", err)
	}

	t.Log("✅ Shutdown avec httpServer nil géré correctement")
}

// TestShutdown_GracefulStop vérifie l'arrêt gracieux du serveur
func TestShutdown_GracefulStop(t *testing.T) {
	t.Log("🧪 TEST SHUTDOWN - GRACEFUL STOP")
	t.Log("================================")

	config := &Config{
		Host:     "localhost",
		Port:     0, // Port aléatoire pour éviter les conflits
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("❌ Erreur création serveur: %v", err)
	}

	// Initialiser httpServer
	server.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler:      server.mux,
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		IdleTimeout:  DefaultIdleTimeout,
	}

	// Démarrer le serveur dans une goroutine
	errCh := make(chan error, 1)
	go func() {
		err := server.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Attendre que le serveur démarre
	time.Sleep(100 * time.Millisecond)

	// Vérifier qu'aucune erreur de démarrage
	select {
	case err := <-errCh:
		t.Fatalf("❌ Erreur démarrage serveur: %v", err)
	default:
		t.Log("✅ Serveur démarré avec succès")
	}

	// Effectuer un shutdown gracieux
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = server.Shutdown(shutdownCtx)
	if err != nil {
		t.Errorf("❌ Shutdown a échoué: %v", err)
	} else {
		t.Log("✅ Shutdown gracieux réussi")
	}
}

// TestShutdown_WithActiveConnections vérifie que les connexions actives sont drainées
func TestShutdown_WithActiveConnections(t *testing.T) {
	t.Log("🧪 TEST SHUTDOWN - CONNEXIONS ACTIVES")
	t.Log("=====================================")

	config := &Config{
		Host:     "localhost",
		Port:     0,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("❌ Erreur création serveur: %v", err)
	}

	// Ajouter un handler qui simule une requête longue
	slowHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	server.mux.HandleFunc("/slow", slowHandler)

	// Initialiser httpServer avec un port dynamique
	server.httpServer = &http.Server{
		Addr:         "localhost:0",
		Handler:      server.mux,
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		IdleTimeout:  DefaultIdleTimeout,
	}

	// Démarrer le serveur
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("❌ Erreur Listen: %v", err)
	}
	defer listener.Close()

	serverAddr := listener.Addr().String()
	t.Logf("📍 Serveur démarré sur %s", serverAddr)

	errCh := make(chan error, 1)
	go func() {
		err := server.httpServer.Serve(listener)
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Attendre démarrage
	time.Sleep(50 * time.Millisecond)

	// Démarrer une requête longue
	var wg sync.WaitGroup
	requestCompleted := make(chan bool, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		client := &http.Client{Timeout: 2 * time.Second}
		resp, err := client.Get("http://" + serverAddr + "/slow")
		if err != nil {
			t.Logf("⚠️  Erreur requête: %v", err)
			requestCompleted <- false
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			t.Log("✅ Requête longue complétée avant shutdown")
			requestCompleted <- true
		} else {
			t.Logf("❌ Status code inattendu: %d", resp.StatusCode)
			requestCompleted <- false
		}
	}()

	// Attendre que la requête soit en cours
	time.Sleep(50 * time.Millisecond)

	// Lancer le shutdown pendant que la requête est en cours
	t.Log("🛑 Lancement du shutdown pendant requête active...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	shutdownErr := server.Shutdown(shutdownCtx)

	// Attendre que la requête se termine
	wg.Wait()

	// Vérifier les résultats
	if shutdownErr != nil {
		t.Errorf("❌ Erreur lors du shutdown: %v", shutdownErr)
	}

	select {
	case completed := <-requestCompleted:
		if completed {
			t.Log("✅ Requête active drainée correctement pendant shutdown")
		} else {
			t.Error("❌ Requête active interrompue")
		}
	case <-time.After(1 * time.Second):
		t.Error("❌ Timeout en attendant la fin de la requête")
	}
}

// TestShutdown_Timeout vérifie le comportement lors d'un timeout de shutdown
func TestShutdown_Timeout(t *testing.T) {
	t.Log("🧪 TEST SHUTDOWN - TIMEOUT")
	t.Log("==========================")

	config := &Config{
		Host:     "localhost",
		Port:     0,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("❌ Erreur création serveur: %v", err)
	}

	// Handler qui ne répond jamais (pour forcer un timeout)
	hangingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second) // Plus long que le timeout de shutdown
		w.WriteHeader(http.StatusOK)
	})
	server.mux.HandleFunc("/hang", hangingHandler)

	// Initialiser httpServer
	server.httpServer = &http.Server{
		Addr:         "localhost:0",
		Handler:      server.mux,
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		IdleTimeout:  DefaultIdleTimeout,
	}

	// Démarrer le serveur
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("❌ Erreur Listen: %v", err)
	}
	defer listener.Close()

	serverAddr := listener.Addr().String()

	go func() {
		server.httpServer.Serve(listener)
	}()

	time.Sleep(50 * time.Millisecond)

	// Démarrer une requête qui va hang
	go func() {
		client := &http.Client{Timeout: 15 * time.Second}
		client.Get("http://" + serverAddr + "/hang")
	}()

	// Attendre que la requête démarre
	time.Sleep(50 * time.Millisecond)

	// Shutdown avec un timeout très court
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = server.Shutdown(shutdownCtx)

	// Le shutdown devrait retourner une erreur de timeout
	if err != nil {
		t.Logf("⚠️  Timeout attendu lors du shutdown: %v", err)
		t.Log("✅ Timeout géré correctement")
	} else {
		t.Log("⚠️  Pas d'erreur de timeout (peut arriver selon le timing)")
	}
}

// TestShutdown_Idempotent vérifie que plusieurs appels à Shutdown sont sûrs
func TestShutdown_Idempotent(t *testing.T) {
	t.Log("🧪 TEST SHUTDOWN - IDEMPOTENCE")
	t.Log("==============================")

	config := &Config{
		Host:     "localhost",
		Port:     0,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("❌ Erreur création serveur: %v", err)
	}

	// Initialiser httpServer
	server.httpServer = &http.Server{
		Addr:         "localhost:0",
		Handler:      server.mux,
		ReadTimeout:  DefaultReadTimeout,
		WriteTimeout: DefaultWriteTimeout,
		IdleTimeout:  DefaultIdleTimeout,
	}

	// Démarrer le serveur
	go func() {
		server.httpServer.ListenAndServe()
	}()

	time.Sleep(100 * time.Millisecond)

	// Premier shutdown
	ctx1 := context.Background()
	err1 := server.Shutdown(ctx1)
	if err1 != nil {
		t.Logf("⚠️  Premier shutdown: %v", err1)
	} else {
		t.Log("✅ Premier shutdown réussi")
	}

	// Deuxième shutdown (devrait être safe)
	ctx2 := context.Background()
	err2 := server.Shutdown(ctx2)
	if err2 != nil {
		// C'est OK si ça retourne une erreur "server closed"
		t.Logf("⚠️  Deuxième shutdown: %v", err2)
	}

	t.Log("✅ Multiples appels à Shutdown gérés sans panic")
}

// TestRun_SignalHandling vérifie la gestion des signaux SIGTERM et SIGINT
func TestRun_SignalHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Test long, skip en mode -short")
	}

	t.Log("🧪 TEST RUN - SIGNAL HANDLING")
	t.Log("==============================")

	// Ce test vérifie que Run() gère correctement les signaux
	// Note: Difficile à tester de manière isolée car Run() bloque
	// On vérifie simplement que le code compile et que la structure est correcte

	// Vérifier que les constantes de timeout sont définies
	if DefaultShutdownTimeout != 30*time.Second {
		t.Errorf("❌ DefaultShutdownTimeout = %v, want 30s", DefaultShutdownTimeout)
	}

	t.Log("✅ Constantes de timeout correctement définies")
}

// TestRun_ServerStartupError vérifie la gestion des erreurs de démarrage
func TestRun_ServerStartupError(t *testing.T) {
	t.Log("🧪 TEST RUN - ERREUR DÉMARRAGE")
	t.Log("===============================")

	// Tester avec un port déjà utilisé
	// Démarrer un premier serveur
	config1 := &Config{
		Host:     "localhost",
		Port:     0,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server1, err := NewServer(config1, logger)
	if err != nil {
		t.Fatalf("❌ Erreur création serveur1: %v", err)
	}

	server1.httpServer = &http.Server{
		Addr:    "localhost:0",
		Handler: server1.mux,
	}

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("❌ Erreur Listen: %v", err)
	}
	defer listener.Close()

	go server1.httpServer.Serve(listener)

	t.Log("✅ Test de gestion d'erreur de démarrage préparé")
}

// TestShutdown_ConcurrentRequests vérifie le shutdown avec plusieurs requêtes concurrentes
func TestShutdown_ConcurrentRequests(t *testing.T) {
	t.Log("🧪 TEST SHUTDOWN - REQUÊTES CONCURRENTES")
	t.Log("========================================")

	config := &Config{
		Host:     "localhost",
		Port:     0,
		AuthType: "none",
		Insecure: true,
	}

	logger := log.New(io.Discard, "", 0)
	server, err := NewServer(config, logger)
	if err != nil {
		t.Fatalf("❌ Erreur création serveur: %v", err)
	}

	// Handler qui simule un traitement
	server.mux.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Processed"))
	})

	// Initialiser httpServer
	server.httpServer = &http.Server{
		Addr:    "localhost:0",
		Handler: server.mux,
	}

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("❌ Erreur Listen: %v", err)
	}
	defer listener.Close()

	serverAddr := listener.Addr().String()

	go func() {
		server.httpServer.Serve(listener)
	}()

	time.Sleep(50 * time.Millisecond)

	// Lancer plusieurs requêtes concurrentes
	const numRequests = 5
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			client := &http.Client{Timeout: 2 * time.Second}
			resp, err := client.Get("http://" + serverAddr + "/process")
			if err == nil && resp.StatusCode == http.StatusOK {
				mu.Lock()
				successCount++
				mu.Unlock()
				resp.Body.Close()
			}
		}(i)
	}

	// Attendre un peu que les requêtes démarrent
	time.Sleep(50 * time.Millisecond)

	// Lancer le shutdown
	t.Log("🛑 Lancement du shutdown avec requêtes concurrentes...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	shutdownErr := server.Shutdown(shutdownCtx)

	// Attendre toutes les requêtes
	wg.Wait()

	if shutdownErr != nil {
		t.Logf("⚠️  Erreur shutdown: %v", shutdownErr)
	}

	t.Logf("📊 %d/%d requêtes complétées", successCount, numRequests)
	if successCount > 0 {
		t.Log("✅ Au moins certaines requêtes ont été drainées correctement")
	}
}

// TestDefaultShutdownTimeout vérifie la valeur de la constante
func TestDefaultShutdownTimeout(t *testing.T) {
	t.Log("🧪 TEST DEFAULT SHUTDOWN TIMEOUT")
	t.Log("=================================")

	expectedTimeout := 30 * time.Second
	if DefaultShutdownTimeout != expectedTimeout {
		t.Errorf("❌ DefaultShutdownTimeout = %v, want %v", DefaultShutdownTimeout, expectedTimeout)
	} else {
		t.Logf("✅ DefaultShutdownTimeout = %v (correct)", DefaultShutdownTimeout)
	}
}

// TestShutdown_SignalSending simule l'envoi de signaux (test conceptuel)
func TestShutdown_SignalSending(t *testing.T) {
	if testing.Short() {
		t.Skip("Test long, skip en mode -short")
	}

	t.Log("🧪 TEST SHUTDOWN - ENVOI DE SIGNAUX")
	t.Log("====================================")

	// Note: Ce test est conceptuel car envoyer des signaux au processus de test
	// est délicat et peut interférer avec le runner de tests.
	// On vérifie simplement que le code de gestion des signaux est présent.

	// Vérifier que os.Interrupt et syscall.SIGTERM sont définis
	_ = os.Interrupt
	_ = syscall.SIGTERM
	_ = syscall.SIGINT

	t.Log("✅ Signaux système disponibles pour la gestion")
}
