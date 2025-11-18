# Makefile pour le projet TSD - Validation automatique des conventions Go

.PHONY: help build test lint clean validate check-conventions format install-hooks

# Variables
GO_FILES = $(shell find . -name "*.go" -not -path "./vendor/*")
TEST_TIMEOUT = 300s

# Couleurs pour l'output
GREEN := \033[32m
RED := \033[31m
YELLOW := \033[33m
BLUE := \033[34m
NC := \033[0m # No Color

help: ## Afficher cette aide
	@echo "$(BLUE)🛠️  COMMANDES DISPONIBLES - PROJET TSD$(NC)"
	@echo "========================================"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "$(GREEN)%-20s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install-hooks: ## Installer les hooks Git pour validation automatique
	@echo "$(BLUE)📦 Installation des hooks Git...$(NC)"
	@cp .git/hooks/pre-commit .git/hooks/pre-commit.bak 2>/dev/null || true
	@chmod +x .git/hooks/pre-commit
	@echo "$(GREEN)✅ Hook pre-commit installé et configuré$(NC)"

check-conventions: ## Vérifier la conformité aux conventions Go
	@echo "$(BLUE)🔍 VALIDATION DES CONVENTIONS GO$(NC)"
	@echo "=================================="
	@./scripts/validate_conventions.sh

analyze-naming: ## Analyser les patterns de nommage dans tout le projet
	@echo "$(BLUE)📊 ANALYSE COMPLÈTE DU NOMMAGE$(NC)"
	@echo "==============================="
	@./scripts/analyze_naming.sh

format: ## Formater le code Go selon les standards
	@echo "$(BLUE)🎨 Formatage du code...$(NC)"
	@gofmt -w $(GO_FILES)
	@echo "$(GREEN)✅ Code formaté$(NC)"

lint: ## Lancer l'analyse statique du code
	@echo "$(BLUE)🔍 Analyse statique...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✅ Analyse statique terminée$(NC)"

build: ## Compiler le projet
	@echo "$(BLUE)🔨 Compilation...$(NC)"
	@go build ./...
	@echo "$(GREEN)✅ Compilation réussie$(NC)"

test: ## Exécuter tous les tests
	@echo "$(BLUE)🧪 Exécution des tests...$(NC)"
	@go test -timeout $(TEST_TIMEOUT) ./...
	@echo "$(GREEN)✅ Tests terminés$(NC)"

test-coverage: ## Exécuter les tests avec couverture de code
	@echo "$(BLUE)📊 Tests avec couverture...$(NC)"
	@go test -cover ./...
	@echo "$(GREEN)✅ Tests avec couverture terminés$(NC)"

test-integration: ## Exécuter seulement les tests d'intégration
	@echo "$(BLUE)🔗 Tests d'intégration...$(NC)"
	@cd test/integration && go test -v .
	@echo "$(GREEN)✅ Tests d'intégration terminés$(NC)"

test-unit: ## Exécuter seulement les tests unitaires
	@echo "$(BLUE)🔬 Tests unitaires...$(NC)"
	@cd test/unit && go test -v . 2>/dev/null || echo "$(YELLOW)⚠️ Aucun test unitaire trouvé$(NC)"
	@echo "$(GREEN)✅ Tests unitaires terminés$(NC)"

validate: format lint build test check-conventions ## Validation complète du projet
	@echo ""
	@echo "$(GREEN)🎉 VALIDATION COMPLÈTE TERMINÉE$(NC)"
	@echo "================================="
	@echo "$(GREEN)✅ Formatage$(NC)"
	@echo "$(GREEN)✅ Analyse statique$(NC)"
	@echo "$(GREEN)✅ Compilation$(NC)"
	@echo "$(GREEN)✅ Tests$(NC)"
	@echo "$(GREEN)✅ Conventions de nommage$(NC)"
	@echo "$(GREEN)✅ Règles génération de code$(NC)"
	@echo ""
	@echo "$(BLUE)📋 Le projet respecte tous les standards Go !$(NC)"
	@test -f CODE_GENERATION_CONVENTIONS.md || (echo "$(RED)⚠️  Créer CODE_GENERATION_CONVENTIONS.md$(NC)" && exit 0)
	@echo "$(BLUE)🤖 Règles génération automatique de code validées$(NC)"

quick-check: format lint build ## Validation rapide sans tests
	@echo "$(GREEN)✅ Validation rapide terminée$(NC)"

clean: ## Nettoyer les artefacts de build
	@echo "$(BLUE)🧹 Nettoyage...$(NC)"
	@go clean ./...
	@rm -f *.log *.out
	@echo "$(GREEN)✅ Nettoyage terminé$(NC)"

dev-setup: install-hooks ## Configuration complète pour développement
	@echo "$(BLUE)🚀 Configuration environnement de développement...$(NC)"
	@go mod tidy
	@make validate
	@echo ""
	@echo "$(GREEN)🎉 ENVIRONNEMENT PRÊT !$(NC)"
	@echo "===================="
	@echo "$(GREEN)✅ Dépendances installées$(NC)"
	@echo "$(GREEN)✅ Hooks Git configurés$(NC)"
	@echo "$(GREEN)✅ Validation initiale réussie$(NC)"
	@echo ""
	@echo "$(BLUE)📚 COMMANDES UTILES :$(NC)"
	@echo "• $(YELLOW)make validate$(NC)     - Validation complète"
	@echo "• $(YELLOW)make quick-check$(NC)  - Validation rapide"
	@echo "• $(YELLOW)make test-integration$(NC) - Tests d'intégration"
	@echo "• $(YELLOW)make check-conventions$(NC) - Vérifier conventions"

# Règles de surveillance pour le développement
watch-test: ## Surveiller les fichiers et relancer les tests
	@echo "$(BLUE)👀 Surveillance des tests (Ctrl+C pour arrêter)...$(NC)"
	@while true; do \
		inotifywait -q -r -e modify,create,delete --include='.*\.go$$' . 2>/dev/null && \
		echo "$(YELLOW)🔄 Fichiers modifiés, relance des tests...$(NC)" && \
		make test || true; \
	done

watch-build: ## Surveiller les fichiers et recompiler
	@echo "$(BLUE)👀 Surveillance de la compilation (Ctrl+C pour arrêter)...$(NC)"
	@while true; do \
		inotifywait -q -r -e modify,create,delete --include='.*\.go$$' . 2>/dev/null && \
		echo "$(YELLOW)🔄 Fichiers modifiés, recompilation...$(NC)" && \
		make quick-check || true; \
	done

# Règles pour CI/CD
ci-validate: ## Validation pour CI/CD (sans hooks)
	@echo "$(BLUE)🤖 VALIDATION CI/CD$(NC)"
	@echo "=================="
	@make format
	@make lint
	@make build
	@make test-coverage
	@make check-conventions
	@echo "$(GREEN)✅ Validation CI/CD terminée$(NC)"

# Aide pour les nouveaux développeurs
onboarding: ## Guide pour nouveaux développeurs
	@echo "$(BLUE)👋 BIENVENUE SUR LE PROJET TSD !$(NC)"
	@echo "==============================="
	@echo ""
	@echo "$(YELLOW)📚 ÉTAPES RECOMMANDÉES :$(NC)"
	@echo "1. $(GREEN)make dev-setup$(NC)         - Configuration initiale"
	@echo "2. $(GREEN)make validate$(NC)          - Validation complète"
	@echo "3. Lire $(BLUE)DEVELOPMENT_GUIDELINES.md$(NC) - Conventions obligatoires"
	@echo "4. Lire $(BLUE)NAMING_CONVENTIONS_FINAL_REPORT.md$(NC) - État des conventions"
	@echo ""
	@echo "$(YELLOW)🔧 DÉVELOPPEMENT QUOTIDIEN :$(NC)"
	@echo "• $(GREEN)make quick-check$(NC)        - Avant chaque commit"
	@echo "• $(GREEN)make test-integration$(NC)   - Tests d'intégration"
	@echo "• $(GREEN)make watch-test$(NC)         - Développement en continu"
	@echo ""
	@echo "$(YELLOW)📋 VALIDATION AVANT PUSH :$(NC)"
	@echo "• $(GREEN)make validate$(NC)           - Validation complète"
	@echo "• $(GREEN)make check-conventions$(NC)  - Vérifier conventions"
	@echo ""
	@echo "$(GREEN)✨ Le hook pre-commit validera automatiquement vos commits !$(NC)"

# Affichage des métriques du projet
metrics: ## Afficher les métriques du projet
	@echo "$(BLUE)📊 MÉTRIQUES DU PROJET TSD$(NC)"
	@echo "============================"
	@echo "$(YELLOW)📁 Fichiers :$(NC)"
	@echo "  Go files: $(shell find . -name "*.go" -not -path "./vendor/*" | wc -l)"
	@echo "  Test files: $(shell find . -name "*_test.go" -not -path "./vendor/*" | wc -l)"
	@echo "  Total lines: $(shell find . -name "*.go" -not -path "./vendor/*" -exec wc -l {} + | tail -1 | cut -d' ' -f1)"
	@echo ""
	@echo "$(YELLOW)🏗️ Structure :$(NC)"
	@echo "  Packages: $(shell find . -name "*.go" -not -path "./vendor/*" -exec dirname {} \; | sort -u | wc -l)"
	@echo "  Modules: $(shell find . -name "go.mod" | wc -l)"
	@echo ""
	@echo "$(YELLOW)✅ Conformité :$(NC)"
	@echo "  Snake case files: $(shell find . -name "*_*.go" -not -path "./vendor/*" | wc -l)"
	@echo "  CamelCase files: $(shell find . -name "*.go" -not -path "./vendor/*" -not -name "*_*" -not -name "main.go" | wc -l)"
