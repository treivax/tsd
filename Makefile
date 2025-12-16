# Makefile pour le projet TSD - Validation RETE et conventions Go

.PHONY: help build test clean lint format deps validate

# Variables
PROJECT_NAME := tsd
BINARY_NAME := tsd
UNIVERSAL_RUNNER := universal-rete-runner
GO_VERSION := 1.24
BUILD_DIR := ./bin
CMD_TSD_DIR := ./cmd/tsd
CMD_UNIVERSAL_DIR := ./cmd/universal-rete-runner
GO_FILES := $(shell find . -name "*.go" -not -path "./vendor/*")
TEST_TIMEOUT := 10m
TEST_PARALLEL := 4
BETA_TESTS_DIR := ./beta_coverage_tests

# Couleurs pour l'output
GREEN := \033[32m
RED := \033[31m
YELLOW := \033[33m
BLUE := \033[34m
CYAN := \033[36m
NC := \033[0m # No Color

help: ## Afficher cette aide
	@echo "$(BLUE)🛠️  PROJET TSD - VALIDATION RETE$(NC)"
	@echo "================================="
	@echo ""
	@echo "$(CYAN)🏗️  BUILD & INSTALL:$(NC)"
	@echo "$(GREEN)build$(NC)                - Compiler le binaire TSD unique"
	@echo "$(GREEN)build-tsd$(NC)            - Compiler le binaire TSD unique"
	@echo "$(GREEN)build-runners$(NC)        - Compiler les runners de test"
	@echo "$(GREEN)install$(NC)              - Installation complète"
	@echo "$(GREEN)clean$(NC)                - Nettoyer les artefacts"
	@echo ""
	@echo "$(CYAN)🔥 VALIDATION RETE:$(NC)"
	@echo "$(GREEN)rete-all$(NC)             - Valider tous les tests beta"
	@echo "$(GREEN)rete-unified$(NC)         - Exécuter TOUS les tests (Alpha+Beta+Intégration)"
	@echo ""
	@echo "$(CYAN)🧪 TESTS & QUALITÉ:$(NC)"
	@echo "$(GREEN)test-unit$(NC)            - Tests unitaires (rapides)"
	@echo "$(GREEN)test-fixtures$(NC)        - Tests fixtures partagées"
	@echo "$(GREEN)test-e2e$(NC)             - Tests E2E (fixtures TSD)"
	@echo "$(GREEN)test-integration$(NC)     - Tests d'intégration"
	@echo "$(GREEN)test-performance$(NC)     - Tests de performance"
	@echo "$(GREEN)test-all$(NC)             - Tous les tests standards"
	@echo "$(GREEN)test-complete$(NC)        - TOUS les tests (complet)"
	@echo "$(GREEN)coverage$(NC)             - Rapport de couverture"
	@echo "$(GREEN)bench$(NC)                - Benchmarks"
	@echo "$(GREEN)lint$(NC)                 - Analyse statique du code"
	@echo ""
	@echo "$(CYAN)🔒 SÉCURITÉ:$(NC)"
	@echo "$(GREEN)security-scan$(NC)        - Scan sécurité complet (gosec + govulncheck)"
	@echo "$(GREEN)security-gosec$(NC)       - Scan sécurité statique (gosec)"
	@echo "$(GREEN)security-vulncheck$(NC)   - Scan vulnérabilités (govulncheck)"
	@echo "$(GREEN)format$(NC)               - Formatage du code"
	@echo "$(GREEN)check-conventions$(NC)    - Vérifier conventions Go"
	@echo ""
	@echo "$(CYAN)🛠️  DÉVELOPPEMENT:$(NC)"
	@echo "$(GREEN)deps$(NC)                 - Installer les dépendances"
	@echo "$(GREEN)deps-dev$(NC)             - Installer outils de développement"
	@echo "$(GREEN)structure$(NC)            - Afficher la structure"
	@echo "$(GREEN)watch-test$(NC)           - Surveiller et relancer tests"
	@echo ""
	@echo "$(CYAN)✅ VALIDATION:$(NC)"
	@echo "$(GREEN)validate$(NC)             - Validation complète (format+lint+build+test)"
	@echo "$(GREEN)quick-check$(NC)          - Validation rapide (sans tests)"
	@echo "$(GREEN)ci$(NC)                   - Validation pour CI/CD"
	@echo ""
	@echo "$(CYAN)📊 INFORMATION:$(NC)"
	@echo "$(GREEN)info$(NC)                 - Informations sur le projet"
	@echo "$(GREEN)demo$(NC)                 - Démonstration rapide"
	@echo "$(GREEN)rete-unified$(NC)         - Runner unifié (Alpha+Beta+Intégration)"

# ================================
# BUILD & COMPILATION
# ================================

build: build-tsd ## BUILD - Compiler le binaire TSD unique

build-tsd: ## BUILD - Compiler le binaire TSD unique
	@echo "$(BLUE)🔨 Compilation de TSD (binaire unifié)...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_TSD_DIR)
	@echo "$(GREEN)✅ Binaire unifié créé: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"
	@echo "   Rôles disponibles: auth, client, server, compilateur (défaut)"

build-runners: ## BUILD - Compiler les runners de test (DEPRECATED - use go test)
	@echo "$(YELLOW)⚠️  DEPRECATED: Le runner universel n'existe plus$(NC)"
	@echo "$(YELLOW)    Utilisez 'make test-e2e' à la place$(NC)"

install: deps build ## BUILD - Installation complète
	@echo "$(GREEN)🚀 Installation terminée$(NC)"
	@echo "   Binaire unifié TSD: $(BUILD_DIR)/$(BINARY_NAME)"
	@echo "   Rôles disponibles:"
	@echo "     - tsd [fichier]      : Compilateur/Runner (défaut)"
	@echo "     - tsd auth ...       : Gestion authentification"
	@echo "     - tsd client ...     : Client HTTP"
	@echo "     - tsd server ...     : Serveur HTTP"

clean: ## BUILD - Nettoyer les artefacts
	@echo "$(BLUE)🧹 Nettoyage...$(NC)"
	@rm -rf $(BUILD_DIR)
	@go clean ./...
	@rm -f *.log *.out
	@echo "$(GREEN)✅ Nettoyage terminé$(NC)"

# ================================
# VALIDATION RETE
# ================================

rete-all: build ## RETE - Valider tous les tests beta
	@echo "$(BLUE)🔥 Validation de tous les tests RETE...$(NC)"
	@cd test/coverage/beta && ./run_all_rete_tests.sh

rete-unified: test-e2e ## RETE - Exécuter TOUS les tests (Alpha+Beta+Intégration) via go test
	@echo "$(BLUE)🚀 TOUS LES TESTS RETE via go test$(NC)"
	@echo "========================================"
	@echo "$(GREEN)✅ Tests exécutés via 'make test-e2e'$(NC)"

rete-unified-legacy: build-runners ## RETE - Ancien runner universel (DEPRECATED)
	@echo "$(YELLOW)⚠️  DEPRECATED: Utilisez 'make test-e2e'$(NC)"
	@$(BUILD_DIR)/$(UNIVERSAL_RUNNER) $(PWD)

rete-unified-report: build-runners ## RETE - Générer seulement le rapport universel
	@echo "$(CYAN)📄 Génération rapport universel...$(NC)"
	@$(BUILD_DIR)/$(UNIVERSAL_RUNNER) $(PWD) report

# ================================
# TESTS & QUALITÉ
# ================================

test: test-unit ## TEST - Alias pour tests unitaires (raccourci)

test-unit: ## TEST - Tests unitaires (rapides, sans build tags)
	@echo "$(BLUE)🧪 Exécution des tests unitaires...$(NC)"
	@go test -v -short -timeout=$(TEST_TIMEOUT) ./constraint/... ./rete/... ./cmd/...
	@echo "$(GREEN)✅ Tests unitaires terminés$(NC)"

test-fixtures: ## TEST - Tests des fixtures partagées
	@echo "$(YELLOW)⚠️  Les fixtures sont des fichiers de données (.tsd), pas des tests Go$(NC)"
	@echo "$(BLUE)📦 Utiliser 'make test-e2e' pour tester les fixtures$(NC)"
	@echo "$(CYAN)📊 Fixtures disponibles:$(NC)"
	@echo "   - Alpha: $$(find ./tests/fixtures/alpha -name '*.tsd' 2>/dev/null | wc -l) fichiers"
	@echo "   - Beta: $$(find ./tests/fixtures/beta -name '*.tsd' 2>/dev/null | wc -l) fichiers"
	@echo "   - Integration: $$(find ./tests/fixtures/integration -name '*.tsd' 2>/dev/null | wc -l) fichiers"

test-e2e: ## TEST - Tests E2E (fixtures TSD)
	@echo "$(BLUE)🎯 Exécution des tests E2E...$(NC)"
	@go test -v -tags=e2e -timeout=$(TEST_TIMEOUT) ./tests/e2e/...
	@echo "$(GREEN)✅ Tests E2E terminés$(NC)"

test-e2e-alpha: ## TEST - Tests fixtures alpha uniquement
	@echo "$(BLUE)🎯 Tests fixtures alpha...$(NC)"
	@go test -v -tags=e2e -run=TestAlphaFixtures -timeout=$(TEST_TIMEOUT) ./tests/e2e/...

test-e2e-beta: ## TEST - Tests fixtures beta uniquement
	@echo "$(BLUE)🎯 Tests fixtures beta...$(NC)"
	@go test -v -tags=e2e -run=TestBetaFixtures -timeout=$(TEST_TIMEOUT) ./tests/e2e/...

test-e2e-integration: ## TEST - Tests fixtures integration uniquement
	@echo "$(BLUE)🎯 Tests fixtures integration...$(NC)"
	@go test -v -tags=e2e -run=TestIntegrationFixtures -timeout=$(TEST_TIMEOUT) ./tests/e2e/...

test-integration: ## TEST - Tests d'intégration (modules)
	@echo "$(BLUE)🔗 Exécution des tests d'intégration...$(NC)"
	@go test -v -tags=integration -timeout=$(TEST_TIMEOUT) ./tests/integration/...
	@echo "$(GREEN)✅ Tests d'intégration terminés$(NC)"

test-integration-verbose: ## TEST - Tests d'intégration avec logs détaillés
	@echo "$(BLUE)🔗 Tests d'intégration (verbose)...$(NC)"
	@go test -v -tags=integration -count=1 -timeout=$(TEST_TIMEOUT) ./tests/integration/... 2>&1 | tee integration-test.log
	@echo "$(GREEN)✅ Logs sauvegardés: integration-test.log$(NC)"

test-integration-coverage: ## TEST - Tests d'intégration avec couverture
	@echo "$(BLUE)🔗 Tests d'intégration avec couverture...$(NC)"
	@go test -v -tags=integration -timeout=$(TEST_TIMEOUT) -coverprofile=integration-coverage.out ./tests/integration/...
	@go tool cover -html=integration-coverage.out -o integration-coverage.html
	@echo "$(GREEN)📊 Rapport de couverture: integration-coverage.html$(NC)"
	@go tool cover -func=integration-coverage.out | grep total


test-performance: ## TEST - Tests de performance et load
	@echo "$(BLUE)⚡ Exécution des tests de performance...$(NC)"
	@go test -v -tags=performance -timeout=1h ./tests/performance/...
	@echo "$(GREEN)✅ Tests de performance terminés$(NC)"

test-load: ## TEST - Tests de charge avec profiling
	@echo "$(BLUE)📈 Tests de charge avec profiling...$(NC)"
	@go test -v -tags=performance -run=TestLoad -cpuprofile=cpu.prof -memprofile=mem.prof ./tests/performance/...
	@echo "$(GREEN)✅ Profiles générés: cpu.prof, mem.prof$(NC)"

test-all: test-unit test-integration test-e2e test-performance ## TEST - Tous les tests standards
	@echo ""
	@echo "$(GREEN)🎉 TOUS LES TESTS STANDARDS RÉUSSIS$(NC)"

test-complete: ## TEST - TOUS les tests (tous les sous-répertoires de tests/)
	@echo "$(BLUE)🚀 Exécution COMPLÈTE de tous les tests...$(NC)"
	@echo "$(CYAN)📂 Tests unitaires...$(NC)"
	@go test -v -short -timeout=$(TEST_TIMEOUT) ./constraint/... ./rete/... ./cmd/...
	@echo ""
	@echo "$(CYAN)🔗 Tests intégration...$(NC)"
	@go test -v -tags=integration -timeout=$(TEST_TIMEOUT) ./tests/integration/...
	@echo ""
	@echo "$(CYAN)🎯 Tests E2E...$(NC)"
	@go test -v -tags=e2e -timeout=$(TEST_TIMEOUT) ./tests/e2e/...
	@echo ""
	@echo "$(CYAN)⚡ Tests performance...$(NC)"
	@go test -v -tags=performance -timeout=1h ./tests/performance/...
	@echo ""
	@echo "$(GREEN)🎉 VALIDATION COMPLÈTE - TOUS LES TESTS RÉUSSIS$(NC)"

test-race: ## TEST - Tests avec race detector
	@echo "$(BLUE)🏁 Tests avec race detector...$(NC)"
	@go test -race -tags=e2e,integration ./...
	@echo "$(GREEN)✅ Tests race terminés$(NC)"

test-parallel: ## TEST - Tests en parallèle
	@echo "$(BLUE)⚡ Tests en parallèle ($(TEST_PARALLEL) workers)...$(NC)"
	@go test -v -tags=e2e,integration -parallel=$(TEST_PARALLEL) ./tests/...

coverage: ## TEST - Rapport de couverture complet
	@echo "$(BLUE)📊 Génération du rapport de couverture...$(NC)"
	@go test -tags=e2e,integration -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Rapport généré: coverage.html$(NC)"

coverage-prod: ## TEST - Couverture code production (sans exemples)
	@echo "$(BLUE)📊 Génération couverture code production...$(NC)"
	@echo "$(YELLOW)⚠️  Exclusion: examples/, rete/examples/, tests/shared/testutil$(NC)"
	@go test -tags=e2e,integration -coverprofile=coverage-prod.out \
		$$(go list ./... | grep -v '/examples' | grep -v '/testutil')
	@go tool cover -html=coverage-prod.out -o coverage-prod.html
	@echo ""
	@echo "$(CYAN)📊 Couverture Globale Production:$(NC)"
	@go tool cover -func=coverage-prod.out | grep total
	@echo ""
	@echo "$(GREEN)✅ Rapport production: coverage-prod.html$(NC)"

coverage-unit: ## TEST - Couverture tests unitaires uniquement
	@echo "$(BLUE)📊 Couverture tests unitaires...$(NC)"
	@go test -short -coverprofile=coverage-unit.out ./constraint/... ./rete/...
	@go tool cover -html=coverage-unit.out -o coverage-unit.html
	@echo "$(GREEN)✅ Rapport: coverage-unit.html$(NC)"

coverage-e2e: ## TEST - Couverture tests E2E uniquement
	@echo "$(BLUE)📊 Couverture tests E2E...$(NC)"
	@go test -tags=e2e -coverprofile=coverage-e2e.out ./tests/e2e/...
	@go tool cover -html=coverage-e2e.out -o coverage-e2e.html
	@echo "$(GREEN)✅ Rapport: coverage-e2e.html$(NC)"

coverage-report: coverage-prod ## TEST - Rapport détaillé couverture production
	@echo ""
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
	@echo "$(CYAN)📊 RAPPORT DE COUVERTURE - CODE PRODUCTION$(NC)"
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"
	@echo ""
	@echo "$(BLUE)📈 Couverture Globale:$(NC)"
	@go tool cover -func=coverage-prod.out | grep total
	@echo ""
	@echo "$(BLUE)📋 Couverture par Module (>80%):$(NC)"
	@go tool cover -func=coverage-prod.out | grep -E "github.com/treivax/tsd/(auth|cmd|constraint|internal|rete|tsdio)/" | grep -v "_test.go" | awk '{print $$1, $$NF}' | sort -t: -k2 -rn | head -20
	@echo ""
	@echo "$(GREEN)✅ Fichier HTML: coverage-prod.html$(NC)"
	@echo "$(CYAN)═══════════════════════════════════════════════════════════$(NC)"

bench: ## TEST - Benchmarks standards
	@echo "$(BLUE)⏱️  Exécution des benchmarks...$(NC)"
	@go test -bench=. -benchmem -run=^$$ ./...

bench-performance: ## TEST - Benchmarks de performance
	@echo "$(BLUE)⚡ Benchmarks de performance...$(NC)"
	@go test -tags=performance -bench=. -benchmem -run=^$$ ./tests/performance/...

bench-profile: ## TEST - Benchmarks avec profiling
	@echo "$(BLUE)📊 Benchmarks avec profiling...$(NC)"
	@go test -bench=. -benchmem -cpuprofile=bench-cpu.prof -memprofile=bench-mem.prof ./...
	@echo "$(GREEN)✅ Profiles: bench-cpu.prof, bench-mem.prof$(NC)"

profile-cpu: ## TEST - Visualiser profile CPU
	@echo "$(BLUE)🔍 Ouverture du profile CPU sur :8080...$(NC)"
	@go tool pprof -http=:8080 cpu.prof

profile-mem: ## TEST - Visualiser profile mémoire
	@echo "$(BLUE)🔍 Ouverture du profile mémoire sur :8080...$(NC)"
	@go tool pprof -http=:8080 mem.prof

test-verbose: ## TEST - Tests avec sortie verbose
	@echo "$(BLUE)📢 Tests en mode verbose...$(NC)"
	@go test -v -tags=e2e,integration ./...

test-smoke: ## TEST - Tests rapides (smoke test)
	@echo "$(BLUE)💨 Smoke test...$(NC)"
	@go test -short -run=TestAlphaFixtures ./tests/e2e/... 2>&1 | head -20

clean-test: ## TEST - Nettoyer artefacts de test
	@echo "$(BLUE)🧹 Nettoyage des artefacts de test...$(NC)"
	@rm -f coverage*.out coverage*.html
	@rm -f *.prof
	@rm -f *.test
	@echo "$(GREEN)✅ Artefacts nettoyés$(NC)"

lint: ## TEST - Analyse statique du code
	@echo "$(BLUE)🔍 Analyse statique...$(NC)"
	@go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)⚠️ golangci-lint non installé$(NC)"; \
	fi
	@echo "$(GREEN)✅ Analyse statique terminée$(NC)"

security-gosec: ## SECURITY - Scanner sécurité statique avec gosec
	@echo "$(BLUE)🔍 Analyse sécurité statique (gosec)...$(NC)"
	@if command -v gosec >/dev/null 2>&1; then \
		gosec -quiet ./...; \
	else \
		echo "$(YELLOW)⚠️ gosec non installé$(NC)"; \
		echo "   Installation: go install github.com/securego/gosec/v2/cmd/gosec@latest"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Analyse gosec terminée$(NC)"

security-vulncheck: ## SECURITY - Scanner vulnérabilités dépendances avec govulncheck
	@echo "$(BLUE)🛡️  Scan vulnérabilités (govulncheck)...$(NC)"
	@if command -v govulncheck >/dev/null 2>&1; then \
		govulncheck ./...; \
	else \
		echo "$(YELLOW)⚠️ govulncheck non installé$(NC)"; \
		echo "   Installation: make deps-dev"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Scan govulncheck terminé$(NC)"

security-scan: security-gosec security-vulncheck ## SECURITY - Scan sécurité complet (gosec + govulncheck)
	@echo ""
	@echo "$(GREEN)🎉 SCAN SÉCURITÉ COMPLET RÉUSSI$(NC)"
	@echo "================================="
	@echo "$(GREEN)✅ Analyse statique (gosec)$(NC)"
	@echo "$(GREEN)✅ Scan vulnérabilités (govulncheck)$(NC)"

format: ## TEST - Formatage du code
	@echo "$(BLUE)✨ Formatage du code...$(NC)"
	@go fmt ./...
	@if command -v goimports >/dev/null 2>&1; then \
		find . -name "*.go" -not -path "./vendor/*" | xargs goimports -w; \
	fi
	@echo "$(GREEN)✅ Code formaté$(NC)"

check-conventions: ## TEST - Vérifier conventions Go
	@echo "$(BLUE)🔍 Validation des conventions...$(NC)"
	@if [ -f scripts/validate_conventions.sh ]; then \
		./scripts/validate_conventions.sh; \
	else \
		echo "$(YELLOW)⚠️ Script de validation non trouvé$(NC)"; \
	fi
	@echo "$(GREEN)✅ Conventions vérifiées$(NC)"

# ================================
# DÉVELOPPEMENT
# ================================

deps: ## DEV - Installer les dépendances
	@echo "$(BLUE)📦 Installation des dépendances...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✅ Dépendances installées$(NC)"

deps-dev: ## DEV - Installer outils de développement
	@echo "$(BLUE)🛠️ Installation des outils...$(NC)"
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install golang.org/x/vuln/cmd/govulncheck@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "$(CYAN)Installation de golangci-lint...$(NC)"; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin latest; \
	fi
	@echo "$(GREEN)✅ Outils installés$(NC)"

structure: ## DEV - Afficher la structure
	@echo "$(BLUE)📁 Structure du projet:$(NC)"
	@if command -v tree >/dev/null 2>&1; then \
		tree -I 'vendor|node_modules|.git|bin' -L 3; \
	else \
		find . -type d -not -path "./vendor*" -not -path "./.git*" -not -path "./bin*" | head -20; \
	fi

watch-test: ## DEV - Surveiller et relancer tests
	@echo "$(BLUE)👀 Surveillance des tests (Ctrl+C pour arrêter)...$(NC)"
	@while true; do \
		if command -v inotifywait >/dev/null 2>&1; then \
			inotifywait -q -r -e modify,create,delete --include='.*\.go$$' . 2>/dev/null && \
			echo "$(YELLOW)🔄 Relance des tests...$(NC)" && \
			make test || true; \
		else \
			echo "$(RED)❌ inotifywait non installé$(NC)"; \
			break; \
		fi \
	done

# ================================
# VALIDATION COMPLÈTE
# ================================

validate: format lint security-scan build test-complete ## VALIDATION COMPLÈTE (tous les tests)
	@echo ""
	@echo "$(GREEN)🎉 VALIDATION COMPLÈTE RÉUSSIE$(NC)"
	@echo "==============================="
	@echo "$(GREEN)✅ Formatage$(NC)"
	@echo "$(GREEN)✅ Analyse statique$(NC)"
	@echo "$(GREEN)✅ Scan de sécurité$(NC)"
	@echo "$(GREEN)✅ Compilation$(NC)"
	@echo "$(GREEN)✅ Tests unitaires$(NC)"
	@echo "$(GREEN)✅ Tests fixtures$(NC)"
	@echo "$(GREEN)✅ Tests d'intégration$(NC)"
	@echo "$(GREEN)✅ Tests E2E$(NC)"
	@echo "$(GREEN)✅ Tests performance$(NC)"
	@echo ""
	@echo "$(BLUE)🚀 Projet prêt pour la production !$(NC)"

quick-check: format lint build ## Validation rapide sans tests
	@echo "$(GREEN)✅ Validation rapide terminée$(NC)"

ci: clean deps lint test-complete build ## Validation pour CI/CD
	@echo "$(GREEN)🤖 Validation CI/CD terminée$(NC)"

# ================================
# MÉTRIQUES & INFORMATION
# ================================

info: ## Informations sur le projet
	@echo "$(BLUE)📊 INFORMATIONS PROJET TSD$(NC)"
	@echo "=========================="
	@echo "$(YELLOW)Nom:$(NC) $(PROJECT_NAME)"
	@echo "$(YELLOW)CLI:$(NC) $(BINARY_NAME)"
	@echo "$(YELLOW)Go version:$(NC) $(GO_VERSION)"
	@echo "$(YELLOW)Fichiers Go:$(NC) $(shell echo $(GO_FILES) | wc -w)"
	@echo "$(YELLOW)Packages:$(NC) $(shell find . -name "*.go" -not -path "./vendor/*" -exec dirname {} \; | sort -u | wc -l)"
	@echo ""
	@echo "$(CYAN)🏗️  ARCHITECTURE:$(NC)"
	@echo "• cmd/tsd/              - CLI principal"
	@echo "• cmd/*-runner/         - Runners de tests"
	@echo "• constraint/           - Parseur de contraintes"
	@echo "• rete/                 - Moteur RETE"
	@echo "• test/                 - Tests et validation"

demo: rete-quick ## Démonstration rapide
	@echo ""
	@echo "$(CYAN)✨ DÉMONSTRATION TERMINÉE$(NC)"
	@echo "Pour plus de tests: make rete-all"
