# Makefile pour le projet TSD - Validation RETE et conventions Go

.PHONY: help build test clean lint format deps validate rete-validate

# Variables
PROJECT_NAME := tsd
BINARY_NAME := tsd
RETE_VALIDATE := rete-validate
UNIVERSAL_RUNNER := universal-rete-runner
GO_VERSION := 1.21
BUILD_DIR := ./bin
CMD_TSD_DIR := ./cmd/tsd
CMD_RETE_VALIDATE_DIR := ./cmd/rete-validate
CMD_UNIVERSAL_DIR := ./cmd/universal-rete-runner
GO_FILES := $(shell find . -name "*.go" -not -path "./vendor/*")
TEST_TIMEOUT := 300s
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
	@echo "$(GREEN)build$(NC)                - Compiler tous les binaires"
	@echo "$(GREEN)build-tsd$(NC)            - Compiler l'outil CLI principal"
	@echo "$(GREEN)build-runners$(NC)        - Compiler les runners de test"
	@echo "$(GREEN)install$(NC)              - Installation complète"
	@echo "$(GREEN)clean$(NC)                - Nettoyer les artefacts"
	@echo ""
	@echo "$(CYAN)🔥 VALIDATION RETE:$(NC)"
	@echo "$(GREEN)rete-validate$(NC)        - Valider un test (make rete-validate TEST=join_simple)"
	@echo "$(GREEN)rete-all$(NC)             - Valider tous les tests beta"
	@echo "$(GREEN)rete-quick$(NC)           - Test rapide (join_simple)"
	@echo "$(GREEN)rete-unified$(NC)         - Exécuter TOUS les tests (Alpha+Beta+Intégration)"
	@echo "$(GREEN)rete-dev$(NC)             - Interface développeur"
	@echo ""
	@echo "$(CYAN)🧪 TESTS & QUALITÉ:$(NC)"
	@echo "$(GREEN)test$(NC)                 - Tests unitaires"
	@echo "$(GREEN)test-coverage$(NC)        - Tests avec couverture"
	@echo "$(GREEN)test-integration$(NC)     - Tests d'intégration"
	@echo "$(GREEN)lint$(NC)                 - Analyse statique du code"
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

build: build-tsd build-runners ## BUILD - Compiler tous les binaires

build-tsd: ## BUILD - Compiler l'outil CLI principal
	@echo "$(BLUE)🔨 Compilation de TSD CLI...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_TSD_DIR)
	@echo "$(GREEN)✅ Binaire créé: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

build-runners: ## BUILD - Compiler les runners de test
	@echo "$(BLUE)🔨 Compilation des runners...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(RETE_VALIDATE) $(CMD_RETE_VALIDATE_DIR)
	@go build -o $(BUILD_DIR)/$(UNIVERSAL_RUNNER) $(CMD_UNIVERSAL_DIR)
	@echo "$(GREEN)✅ Runners compilés:$(NC)"
	@echo "   - $(BUILD_DIR)/$(RETE_VALIDATE)"
	@echo "   - $(BUILD_DIR)/$(UNIVERSAL_RUNNER)"

install: deps build ## BUILD - Installation complète
	@echo "$(GREEN)🚀 Installation terminée$(NC)"
	@echo "   TSD CLI: $(BUILD_DIR)/$(BINARY_NAME)"
	@echo "   Test Runners: $(BUILD_DIR)/$(RETE_VALIDATE), $(BUILD_DIR)/$(UNIVERSAL_RUNNER)"

clean: ## BUILD - Nettoyer les artefacts
	@echo "$(BLUE)🧹 Nettoyage...$(NC)"
	@rm -rf $(BUILD_DIR)
	@go clean ./...
	@rm -f *.log *.out
	@echo "$(GREEN)✅ Nettoyage terminé$(NC)"

# ================================
# VALIDATION RETE
# ================================

rete-validate: build-runners ## RETE - Valider un test spécifique (make rete-validate TEST=join_simple)
	@if [ -z "$(TEST)" ]; then \
		echo "$(RED)❌ Erreur: Spécifiez un test avec TEST=nom_du_test$(NC)"; \
		echo "   $(YELLOW)Exemple: make rete-validate TEST=join_simple$(NC)"; \
		exit 1; \
	fi
	@echo "$(CYAN)🎯 Validation RETE du test: $(TEST)$(NC)"
	@$(BUILD_DIR)/$(RETE_VALIDATE) $(BETA_TESTS_DIR)/$(TEST).constraint $(BETA_TESTS_DIR)/$(TEST).facts

rete-all: build ## RETE - Valider tous les tests beta
	@echo "$(BLUE)🔥 Validation de tous les tests RETE...$(NC)"
	@cd test/coverage/beta && ./run_all_rete_tests.sh

rete-quick: ## RETE - Test rapide avec runner.go (join_simple)
	@echo "$(CYAN)⚡ Test RETE rapide...$(NC)"
	@cd test/coverage/beta && go run runner.go /home/resinsec/dev/tsd/beta_coverage_tests/join_simple.constraint /home/resinsec/dev/tsd/beta_coverage_tests/join_simple.facts

rete-dev: ## RETE - Interface développeur (cd test/coverage/beta)
	@echo "$(YELLOW)🛠️  Interface développeur activée$(NC)"
	@echo "   Répertoire: test/coverage/beta/"
	@echo "   Commande: go run runner.go [constraint] [facts]"
	@cd test/coverage/beta && bash

rete-unified: build-runners ## RETE - Exécuter TOUS les tests (Alpha+Beta+Intégration)
	@echo "$(BLUE)🚀 RUNNER UNIVERSEL - TOUS LES TESTS RETE$(NC)"
	@echo "========================================"
	@$(BUILD_DIR)/$(UNIVERSAL_RUNNER) $(PWD)

rete-unified-report: build-runners ## RETE - Générer seulement le rapport universel
	@echo "$(CYAN)📄 Génération rapport universel...$(NC)"
	@$(BUILD_DIR)/$(UNIVERSAL_RUNNER) $(PWD) report

# ================================
# TESTS & QUALITÉ
# ================================

test: ## TEST - Tests unitaires
	@echo "$(BLUE)🧪 Exécution des tests unitaires...$(NC)"
	@go test -timeout $(TEST_TIMEOUT) ./...
	@echo "$(GREEN)✅ Tests unitaires terminés$(NC)"

test-coverage: ## TEST - Tests avec couverture
	@echo "$(BLUE)📊 Tests avec couverture...$(NC)"
	@go test -cover ./...
	@echo "$(GREEN)✅ Tests avec couverture terminés$(NC)"

test-integration: ## TEST - Tests d'intégration
	@echo "$(BLUE)🔗 Tests d'intégration...$(NC)"
	@cd test/integration && go test -v .
	@echo "$(GREEN)✅ Tests d'intégration terminés$(NC)"

lint: ## TEST - Analyse statique du code
	@echo "$(BLUE)🔍 Analyse statique...$(NC)"
	@go vet ./...
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)⚠️ golangci-lint non installé$(NC)"; \
	fi
	@echo "$(GREEN)✅ Analyse statique terminée$(NC)"

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

validate: format lint build test ## VALIDATION COMPLÈTE
	@echo ""
	@echo "$(GREEN)🎉 VALIDATION COMPLÈTE RÉUSSIE$(NC)"
	@echo "==============================="
	@echo "$(GREEN)✅ Formatage$(NC)"
	@echo "$(GREEN)✅ Analyse statique$(NC)"
	@echo "$(GREEN)✅ Compilation$(NC)"
	@echo "$(GREEN)✅ Tests$(NC)"
	@echo ""
	@echo "$(BLUE)🚀 Projet prêt pour la production !$(NC)"

quick-check: format lint build ## Validation rapide sans tests
	@echo "$(GREEN)✅ Validation rapide terminée$(NC)"

ci: clean deps lint test build ## Validation pour CI/CD
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
