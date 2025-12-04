// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package metrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// HTTPServer expose les métriques RETE via HTTP
type HTTPServer struct {
	exporter *PrometheusExporter
	mux      *http.ServeMux
	server   *http.Server
	addr     string
	mutex    sync.RWMutex
	started  bool
}

// NewHTTPServer crée un nouveau serveur HTTP de métriques
func NewHTTPServer(exporter *PrometheusExporter, addr string) *HTTPServer {
	if addr == "" {
		addr = ":9090"
	}

	mux := http.NewServeMux()
	server := &HTTPServer{
		exporter: exporter,
		mux:      mux,
		addr:     addr,
	}

	// Register handlers
	mux.Handle("/metrics", promhttp.HandlerFor(
		exporter.Registry(),
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	mux.HandleFunc("/health", server.handleHealth)
	mux.HandleFunc("/ready", server.handleReady)
	mux.HandleFunc("/", server.handleIndex)

	server.server = &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server
}

// Start démarre le serveur HTTP
func (s *HTTPServer) Start() error {
	s.mutex.Lock()
	if s.started {
		s.mutex.Unlock()
		return fmt.Errorf("server already started")
	}
	s.started = true
	s.mutex.Unlock()

	return s.server.ListenAndServe()
}

// Stop arrête le serveur HTTP
func (s *HTTPServer) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.started {
		return fmt.Errorf("server not started")
	}

	s.started = false
	return s.server.Close()
}

// Addr retourne l'adresse du serveur
func (s *HTTPServer) Addr() string {
	return s.addr
}

// handleHealth gère le endpoint de santé
func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"service":   "tsd-metrics",
	}

	json.NewEncoder(w).Encode(response)
}

// handleReady gère le endpoint de disponibilité
func (s *HTTPServer) handleReady(w http.ResponseWriter, r *http.Request) {
	s.mutex.RLock()
	ready := s.started
	s.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")

	if ready {
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"status":    "ready",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		json.NewEncoder(w).Encode(response)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		response := map[string]interface{}{
			"status":    "not ready",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		}
		json.NewEncoder(w).Encode(response)
	}
}

// handleIndex gère la page d'accueil
func (s *HTTPServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TSD RETE Metrics</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            max-width: 800px;
            margin: 40px auto;
            padding: 0 20px;
            line-height: 1.6;
            color: #333;
        }
        h1 {
            color: #2c3e50;
            border-bottom: 3px solid #3498db;
            padding-bottom: 10px;
        }
        h2 {
            color: #34495e;
            margin-top: 30px;
        }
        .endpoint {
            background: #f8f9fa;
            border-left: 4px solid #3498db;
            padding: 15px;
            margin: 15px 0;
            border-radius: 4px;
        }
        .endpoint a {
            color: #3498db;
            text-decoration: none;
            font-weight: 600;
        }
        .endpoint a:hover {
            text-decoration: underline;
        }
        .description {
            color: #7f8c8d;
            margin-top: 5px;
        }
        .footer {
            margin-top: 50px;
            padding-top: 20px;
            border-top: 1px solid #ecf0f1;
            text-align: center;
            color: #95a5a6;
            font-size: 0.9em;
        }
        code {
            background: #ecf0f1;
            padding: 2px 6px;
            border-radius: 3px;
            font-family: 'Courier New', monospace;
        }
        .metrics-list {
            background: #fff;
            border: 1px solid #e1e8ed;
            border-radius: 4px;
            padding: 20px;
            margin: 20px 0;
        }
        .metrics-list ul {
            margin: 10px 0;
            padding-left: 20px;
        }
        .metrics-list li {
            margin: 5px 0;
        }
        .badge {
            display: inline-block;
            background: #3498db;
            color: white;
            padding: 2px 8px;
            border-radius: 3px;
            font-size: 0.85em;
            margin-left: 5px;
        }
    </style>
</head>
<body>
    <h1>🚀 TSD RETE Metrics Server</h1>

    <p>Welcome to the TSD RETE metrics exporter. This service exposes runtime metrics about the RETE engine in Prometheus format.</p>

    <h2>📊 Available Endpoints</h2>

    <div class="endpoint">
        <a href="/metrics">/metrics</a>
        <div class="description">Prometheus metrics endpoint in OpenMetrics format</div>
    </div>

    <div class="endpoint">
        <a href="/health">/health</a>
        <div class="description">Health check endpoint (always returns 200 OK)</div>
    </div>

    <div class="endpoint">
        <a href="/ready">/ready</a>
        <div class="description">Readiness check endpoint (returns 200 when ready)</div>
    </div>

    <h2>📈 Metrics Categories</h2>

    <div class="metrics-list">
        <h3>🔧 Ingestion Metrics <span class="badge">Core</span></h3>
        <ul>
            <li><code>tsd_rete_ingestion_duration_seconds</code> - Duration of ingestion operations by phase</li>
            <li><code>tsd_rete_ingestion_total</code> - Total number of ingestion operations</li>
            <li><code>tsd_rete_types_added_total</code> - Types added to network</li>
            <li><code>tsd_rete_rules_added_total</code> - Rules added to network</li>
            <li><code>tsd_rete_facts_submitted_total</code> - Facts submitted</li>
            <li><code>tsd_rete_facts_propagated_total</code> - Facts propagated through network</li>
        </ul>

        <h3>🏗️ Network State <span class="badge">Structure</span></h3>
        <ul>
            <li><code>tsd_rete_type_nodes</code> - Current number of type nodes</li>
            <li><code>tsd_rete_terminal_nodes</code> - Current number of terminal nodes</li>
            <li><code>tsd_rete_alpha_nodes</code> - Current number of alpha nodes</li>
            <li><code>tsd_rete_beta_nodes</code> - Current number of beta nodes</li>
        </ul>

        <h3>⚡ Performance <span class="badge">Speed</span></h3>
        <ul>
            <li><code>tsd_rete_efficiency_score</code> - Network efficiency score (0-1)</li>
            <li><code>tsd_rete_bottleneck_phase_percentage</code> - Time percentage per phase</li>
            <li><code>tsd_performance_rule_evaluations_total</code> - Total rule evaluations</li>
            <li><code>tsd_performance_rule_matches_total</code> - Successful rule matches</li>
            <li><code>tsd_performance_rule_evaluation_duration_seconds</code> - Rule evaluation time</li>
        </ul>

        <h3>💾 Storage <span class="badge">Data</span></h3>
        <ul>
            <li><code>tsd_storage_facts_total</code> - Facts currently in storage</li>
            <li><code>tsd_storage_operations_total</code> - Storage operations by type</li>
            <li><code>tsd_storage_errors_total</code> - Storage errors by type</li>
        </ul>

        <h3>🔄 Transactions <span class="badge">Safety</span></h3>
        <ul>
            <li><code>tsd_rete_transaction_total</code> - Transaction operations</li>
            <li><code>tsd_rete_transaction_duration_seconds</code> - Transaction duration</li>
            <li><code>tsd_rete_commands_executed_total</code> - Commands executed</li>
            <li><code>tsd_rete_rollbacks_total</code> - Transaction rollbacks</li>
        </ul>

        <h3>✅ Coherence <span class="badge">Quality</span></h3>
        <ul>
            <li><code>tsd_coherence_violations_total</code> - Coherence violations by type</li>
            <li><code>tsd_coherence_checks_total</code> - Coherence checks performed</li>
            <li><code>tsd_coherence_mode</code> - Current coherence mode (1=strong, 2=relaxed, 3=eventual)</li>
        </ul>

        <h3>🔢 Arithmetic <span class="badge">Advanced</span></h3>
        <ul>
            <li><code>tsd_arithmetic_decompositions_total</code> - Arithmetic decompositions</li>
            <li><code>tsd_arithmetic_optimizations_total</code> - Arithmetic optimizations</li>
        </ul>
    </div>

    <h2>🔌 Prometheus Configuration</h2>

    <p>Add this job to your <code>prometheus.yml</code>:</p>

    <div class="metrics-list">
        <pre><code>scrape_configs:
  - job_name: 'tsd-rete'
    scrape_interval: 15s
    static_configs:
      - targets: ['localhost:9090']
        labels:
          service: 'tsd-rete-engine'
          env: 'production'</code></pre>
    </div>

    <h2>📚 Documentation</h2>

    <p>For more information about RETE metrics and monitoring:</p>
    <ul>
        <li>See <code>LOGGING_GUIDE.md</code> for logging and metrics best practices</li>
        <li>See <code>PHASE3_COMPLETION.md</code> for metrics implementation details</li>
        <li>Visit the <a href="https://prometheus.io/docs/introduction/overview/">Prometheus documentation</a></li>
    </ul>

    <div class="footer">
        <p>TSD RETE Engine - Thread-Safe Production Rules System</p>
        <p>© 2025 TSD Contributors | Licensed under MIT</p>
    </div>
</body>
</html>`

	fmt.Fprint(w, html)
}
