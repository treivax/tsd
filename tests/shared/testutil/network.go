// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package testutil

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"
)

// SlowServer crée un serveur de test qui répond lentement.
// Utile pour tester les timeouts et la patience du client.
func SlowServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
}

// UnreliableServer crée un serveur qui échoue de manière aléatoire.
// Le paramètre failRate définit la probabilité d'échec (0.0-1.0).
// Par exemple, 0.3 signifie 30% de chance d'échec par requête.
func UnreliableServer(failRate float64) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rand.Float64() < failRate {
			time.Sleep(10 * time.Second)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))
}

// ClosingServer crée un serveur qui ferme brutalement les connexions.
// Simule une perte de connexion réseau ou un crash serveur.
func ClosingServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
			return
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		conn.Close()
	}))
}

// TimeoutServer crée un serveur qui ne répond jamais.
// Force un timeout côté client.
func TimeoutServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(30 * time.Second)
	}))
}

// IncompleteResponseServer crée un serveur qui envoie une réponse incomplète.
// Simule une connexion interrompue pendant la transmission de données.
func IncompleteResponseServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "data": `))

		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}

		hj, ok := w.(http.Hijacker)
		if !ok {
			return
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			return
		}
		conn.Close()
	}))
}

// FlakyServer crée un serveur qui alterne entre succès et échec.
// Utile pour tester la logique de retry.
type FlakyServer struct {
	server        *httptest.Server
	requestCount  int
	failureCount  int
	successAfter  int
}

// NewFlakyServer crée un serveur qui réussit après un nombre spécifié de tentatives.
// Les premières tentatives échouent, puis le serveur commence à répondre correctement.
func NewFlakyServer(successAfter int) *FlakyServer {
	fs := &FlakyServer{
		successAfter: successAfter,
	}

	fs.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.requestCount++

		if fs.requestCount <= fs.successAfter {
			fs.failureCount++
			hj, ok := w.(http.Hijacker)
			if ok {
				conn, _, _ := hj.Hijack()
				conn.Close()
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}))

	return fs
}

// URL retourne l'URL du serveur de test.
func (fs *FlakyServer) URL() string {
	return fs.server.URL
}

// Close ferme le serveur de test.
func (fs *FlakyServer) Close() {
	fs.server.Close()
}

// RequestCount retourne le nombre total de requêtes reçues.
func (fs *FlakyServer) RequestCount() int {
	return fs.requestCount
}

// FailureCount retourne le nombre de requêtes ayant échoué.
func (fs *FlakyServer) FailureCount() int {
	return fs.failureCount
}

// Reset réinitialise les compteurs du serveur.
func (fs *FlakyServer) Reset() {
	fs.requestCount = 0
	fs.failureCount = 0
}
