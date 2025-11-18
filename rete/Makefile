# Makefile pour le module RETE
# Usage: make [target]

.PHONY: help test coverage clean validate build format lint

# Affichage de l'aide par défaut
help:
	@echo "🎯 Module RETE - Commandes disponibles:"
	@echo ""
	@echo "  build      - Compiler tous les packages"
	@echo "  test       - Exécuter tous les tests avec couverture"
	@echo "  coverage   - Générer le rapport de couverture HTML"
	@echo "  validate   - Valider la structure du projet"
	@echo "  format     - Formater le code Go"
	@echo "  lint       - Analyser le code avec golint"
	@echo "  clean      - Nettoyer les artefacts temporaires"
	@echo ""
	@echo "📁 Structure des répertoires:"
	@echo "  pkg/       - Packages publics"
	@echo "  internal/  - Packages internes"
	@echo "  test/      - Tests et couverture"
	@echo "  docs/      - Documentation"
	@echo "  scripts/   - Scripts utilitaires"

# Compilation
build:
	@echo "🔨 Compilation des packages..."
	@go build ./pkg/...
	@go build ./internal/...
	@echo "✅ Compilation terminée"

# Tests complets avec couverture
test:
	@./scripts/run_tests.sh

# Génération de couverture uniquement
coverage:
	@echo "📊 Génération de la couverture..."
	@go test -coverprofile=test/coverage/coverage.out .
	@go tool cover -html=test/coverage/coverage.out -o test/coverage/reports/coverage.html
	@echo "✅ Rapport disponible: test/coverage/reports/coverage.html"

# Validation de la structure
validate:
	@./scripts/validate_structure.sh

# Formatage du code
format:
	@echo "✨ Formatage du code Go..."
	@go fmt ./...
	@echo "✅ Code formaté"

# Analyse statique (si golint est disponible)
lint:
	@echo "🔍 Analyse du code..."
	@if command -v golint >/dev/null 2>&1; then \
		golint ./pkg/... ; \
		golint ./internal/... ; \
	else \
		echo "⚠️  golint non installé, utilisation de go vet..."; \
		go vet ./... ; \
	fi
	@echo "✅ Analyse terminée"

# Nettoyage
clean:
	@./scripts/clean.sh

# Tests rapides (sans couverture)
test-quick:
	@echo "⚡ Tests rapides..."
	@go test -short ./test/unit/...
	@echo "✅ Tests rapides terminés"

# Installation des dépendances de développement
dev-deps:
	@echo "📦 Installation des outils de développement..."
	@go install golang.org/x/lint/golint@latest || echo "⚠️  golint installation échouée"
	@echo "✅ Outils installés"
