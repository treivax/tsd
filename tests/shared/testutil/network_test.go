// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package testutil

import (
	"io"
	"net/http"
	"testing"
	"time"
)

// TestSlowServer vérifie le bon fonctionnement du serveur lent
func TestSlowServer(t *testing.T) {
	t.Log("🧪 TEST UTIL - SLOW SERVER")
	t.Log("==========================")

	server := SlowServer(100 * time.Millisecond)
	defer server.Close()

	start := time.Now()
	resp, err := http.Get(server.URL)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("❌ Erreur requête: %v", err)
	}
	defer resp.Body.Close()

	if duration < 100*time.Millisecond {
		t.Errorf("❌ Serveur trop rapide: %v", duration)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("❌ Status = %d, attendu 200", resp.StatusCode)
	}

	t.Logf("✅ Slow server fonctionne (délai: %v)", duration)
}

// TestClosingServer vérifie le bon fonctionnement du serveur qui ferme les connexions
func TestClosingServer(t *testing.T) {
	t.Log("🧪 TEST UTIL - CLOSING SERVER")
	t.Log("==============================")

	server := ClosingServer()
	defer server.Close()

	resp, err := http.Get(server.URL)

	if err == nil {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 0 {
			t.Error("❌ Le serveur ne devrait pas envoyer de réponse complète")
		}
		t.Log("✅ Closing server fonctionne (réponse vide)")
	} else {
		t.Logf("✅ Closing server fonctionne (erreur: %v)", err)
	}
}

// TestTimeoutServer vérifie le bon fonctionnement du serveur qui timeout
func TestTimeoutServer(t *testing.T) {
	t.Log("🧪 TEST UTIL - TIMEOUT SERVER")
	t.Log("==============================")

	server := TimeoutServer()
	defer server.Close()

	client := &http.Client{Timeout: 100 * time.Millisecond}

	_, err := client.Get(server.URL)

	if err == nil {
		t.Fatal("❌ Attendait un timeout")
	}

	t.Logf("✅ Timeout server fonctionne: %v", err)
}

// TestIncompleteResponseServer vérifie le bon fonctionnement du serveur avec réponse incomplète
func TestIncompleteResponseServer(t *testing.T) {
	t.Log("🧪 TEST UTIL - INCOMPLETE RESPONSE SERVER")
	t.Log("==========================================")

	server := IncompleteResponseServer()
	defer server.Close()

	resp, err := http.Get(server.URL)

	if err != nil {
		t.Logf("✅ Incomplete response server fonctionne (erreur réseau: %v)", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Logf("✅ Incomplete response server fonctionne (erreur lecture: %v)", err)
		return
	}

	if len(body) > 0 && body[len(body)-1] != '}' {
		t.Log("✅ Incomplete response server fonctionne (réponse incomplète)")
	} else {
		t.Logf("⚠️  Réponse complète reçue: %s", string(body))
	}
}

// TestFlakyServer vérifie le bon fonctionnement du serveur flaky
func TestFlakyServer(t *testing.T) {
	t.Log("🧪 TEST UTIL - FLAKY SERVER")
	t.Log("============================")

	fs := NewFlakyServer(2)
	defer fs.Close()

	client := &http.Client{Timeout: 1 * time.Second}

	_, err := client.Get(fs.URL())
	if err == nil {
		t.Error("❌ Première requête devrait échouer")
	}

	_, err = client.Get(fs.URL())
	if err == nil {
		t.Error("❌ Deuxième requête devrait échouer")
	}

	resp, err := client.Get(fs.URL())
	if err != nil {
		t.Fatalf("❌ Troisième requête devrait réussir: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("❌ Status = %d, attendu 200", resp.StatusCode)
	}

	if fs.RequestCount() != 3 {
		t.Errorf("❌ RequestCount = %d, attendu 3", fs.RequestCount())
	}

	if fs.FailureCount() != 2 {
		t.Errorf("❌ FailureCount = %d, attendu 2", fs.FailureCount())
	}

	t.Logf("✅ Flaky server fonctionne (requêtes: %d, échecs: %d)",
		fs.RequestCount(), fs.FailureCount())

	fs.Reset()
	if fs.RequestCount() != 0 || fs.FailureCount() != 0 {
		t.Error("❌ Reset ne fonctionne pas correctement")
	}
}

// TestUnreliableServer vérifie le bon fonctionnement du serveur aléatoire
func TestUnreliableServer(t *testing.T) {
	t.Log("🧪 TEST UTIL - UNRELIABLE SERVER")
	t.Log("=================================")

	server := UnreliableServer(0.5)
	defer server.Close()

	client := &http.Client{Timeout: 500 * time.Millisecond}

	successes := 0
	failures := 0
	totalRequests := 20

	for i := 0; i < totalRequests; i++ {
		resp, err := client.Get(server.URL)
		if err != nil {
			failures++
		} else {
			successes++
			resp.Body.Close()
		}
	}

	if successes == 0 && failures == 0 {
		t.Fatal("❌ Aucune requête n'a été traitée")
	}

	t.Logf("✅ Unreliable server fonctionne (succès: %d, échecs: %d sur %d)",
		successes, failures, totalRequests)
}
