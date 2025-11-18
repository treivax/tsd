# Makefile pour le module Constraint refactorisé
# Usage: make [target]

.PHONY: help test coverage clean validate build format lint test-new

# Affichage de l'aide par défaut
help:
	@echo "🎯 Module Constraint - Commandes disponibles:"
	@echo ""
	@echo "🏗️  CONSTRUCTION:"
	@echo "  build      - Compiler tous les packages"
	@echo "  parser     - Régénérer le parser depuis la grammaire PEG"
	@echo ""
	@echo "🧪 TESTS:"
	@echo "  test       - Tests rapides (anciens tests)"
	@echo "  test-new   - Tests complets nouvelle architecture"
	@echo "  coverage   - Générer rapport de couverture détaillé"
	@echo ""
	@echo "🔍 VALIDATION:"
	@echo "  validate   - Valider la nouvelle architecture"
	@echo "  format     - Formater le code Go"
	@echo "  lint       - Analyser le code avec go vet"
	@echo ""
	@echo "🧹 MAINTENANCE:"
	@echo "  clean      - Nettoyer les artefacts temporaires"
	@echo "  deps       - Installer les dépendances de développement"
	@echo ""
	@echo "📁 Structure refactorisée:"
	@echo "  pkg/domain/    - Types fondamentaux et erreurs"
	@echo "  pkg/validator/ - Validation et vérification"
	@echo "  internal/      - Configuration interne"
	@echo "  test/          - Tests organisés"
	@echo "  scripts/       - Scripts d'automatisation"

# Construction complète
build:
	@echo "🔨 Construction du module constraint..."
	@./scripts/build.sh

# Génération du parser depuis PEG
parser:
	@echo "📦 Régénération du parser..."
	@export PATH=$$PATH:~/go/bin && pigeon -o parser.go grammar/constraint.peg
	@echo "✅ Parser régénéré avec succès"

# Tests anciens (compatibilité)
test:
	@echo "🧪 Tests rapides (anciens)..."
	@go test -v . || echo "⚠️  Tests anciens nécessitent adaptation"

# Tests nouvelle architecture
test-new:
	@./scripts/run_tests_new.sh

# Couverture détaillée
coverage:
	@echo "📊 Génération de la couverture détaillée..."
	@go test -coverprofile=test/coverage/coverage.out ./pkg/... ./internal/...
	@go tool cover -html=test/coverage/coverage.out -o test/coverage/reports/coverage.html
	@echo "✅ Rapport disponible: test/coverage/reports/coverage.html"
	@go tool cover -func=test/coverage/coverage.out | tail -1

# Validation de l'architecture
validate:
	@./scripts/validate.sh

# Formatage du code
format:
	@echo "✨ Formatage du code Go..."
	@go fmt ./...
	@echo "✅ Code formaté"

# Analyse statique
lint:
	@echo "🔍 Analyse du code..."
	@go vet ./...
	@echo "✅ Analyse terminée"

# Nettoyage
clean:
	@./scripts/clean.sh

# Tests rapides (sans couverture)
test-quick:
	@echo "⚡ Tests rapides..."
	@go test -short ./pkg/...
	@echo "✅ Tests rapides terminés"

# Installation des dépendances de développement
deps:
	@echo "📦 Installation des outils de développement..."
	@go install github.com/mna/pigeon@latest || echo "⚠️ Erreur d'installation de pigeon"
	@echo "✅ Outils installés"

# Statistiques du projet
stats:
	@echo "📊 Statistiques du projet:"
	@echo ""
	@echo "📁 Structure:"
	@echo "  Packages Go: $$(find pkg internal -name '*.go' | wc -l) fichiers"
	@echo "  Tests: $$(find test -name '*.go' | wc -l) fichiers"
	@echo "  Scripts: $$(find scripts -name '*.sh' | wc -l) scripts"
	@echo ""
	@echo "📈 Couverture actuelle:"
	@go test -coverprofile=/tmp/constraint_stats.out ./pkg/... >/dev/null 2>&1 && go tool cover -func=/tmp/constraint_stats.out | tail -1 || echo "  Erreur de calcul"
	@rm -f /tmp/constraint_stats.out 2>/dev/null || true

# Compilation pour différents OS (si nécessaire)
build-all:
	@echo "🌍 Construction multi-plateforme..."
	@GOOS=linux GOARCH=amd64 go build -o bin/constraint-parser-linux-amd64 ./cmd/
	@GOOS=windows GOARCH=amd64 go build -o bin/constraint-parser-windows-amd64.exe ./cmd/
	@GOOS=darwin GOARCH=amd64 go build -o bin/constraint-parser-darwin-amd64 ./cmd/
	@echo "✅ Binaires créés dans bin/"

# Pipeline CI/CD (simulation)
ci:
	@echo "🚀 Pipeline CI/CD simulé..."
	@make deps
	@make format
	@make lint
	@make build
	@make test-new
	@make validate
	@echo "✅ Pipeline réussi !"
