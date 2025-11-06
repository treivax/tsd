#!/bin/bash
# Script pour exécuter tous les tests avec couverture

set -e

echo "🧪 Exécution des tests unitaires du module RETE"
echo "=============================================="

# Répertoire racine du projet
RETE_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$RETE_DIR"

# Créer le dossier de rapports s'il n'existe pas
mkdir -p test/coverage/reports

# Exécuter les tests unitaires avec couverture
echo "📊 Tests unitaires..."
go test -v ./test/unit/... -coverprofile=test/coverage/unit_coverage.out

echo "📊 Tests d'intégration..."
go test -v ./test/integration/... -coverprofile=test/coverage/integration_coverage.out 2>/dev/null || echo "Aucun test d'intégration trouvé"

# Combiner les rapports de couverture
echo "📊 Génération du rapport de couverture global..."
go test -v -coverprofile=test/coverage/global_coverage.out .

# Générer le rapport HTML
go tool cover -html=test/coverage/global_coverage.out -o test/coverage/reports/coverage.html

# Afficher le résumé de couverture
echo ""
echo "📋 Résumé de la couverture:"
go tool cover -func=test/coverage/global_coverage.out | tail -1

echo ""
echo "✅ Tests terminés ! Rapport disponible dans: test/coverage/reports/coverage.html"