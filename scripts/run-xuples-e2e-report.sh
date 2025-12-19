#!/bin/bash

# Script pour exécuter les tests E2E des xuples avec rapport détaillé
# Génère un rapport complet listant types, règles, faits et xuples

set -e

# Couleurs pour le terminal
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
NC='\033[0m' # No Color
BOLD='\033[1m'

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
REPORT_DIR="$PROJECT_ROOT/test-reports"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
REPORT_FILE="$REPORT_DIR/xuples_e2e_report_$TIMESTAMP.txt"
JSON_REPORT="$REPORT_DIR/xuples_e2e_report_$TIMESTAMP.json"

# Créer le répertoire de rapports
mkdir -p "$REPORT_DIR"

# Fonction pour afficher l'en-tête
print_header() {
    echo ""
    echo "═══════════════════════════════════════════════════════════════════════════"
    echo "  $1"
    echo "═══════════════════════════════════════════════════════════════════════════"
    echo ""
}

# Fonction pour afficher une section
print_section() {
    echo ""
    echo "───────────────────────────────────────────────────────────────────────────"
    echo "  $1"
    echo "───────────────────────────────────────────────────────────────────────────"
    echo ""
}

# Fonction pour afficher un message de succès
print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

# Fonction pour afficher un message d'erreur
print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Fonction pour afficher un message d'info
print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

# Fonction pour afficher un message d'avertissement
print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Démarrer le rapport
{
    print_header "RAPPORT D'EXÉCUTION DES TESTS E2E - XUPLES & XUPLE-SPACES"
    echo "Date d'exécution: $(date)"
    echo "Répertoire projet: $PROJECT_ROOT"
    echo ""

    print_section "1. INFORMATIONS SYSTÈME"
    echo "Go version: $(go version)"
    echo "OS: $(uname -s)"
    echo "Architecture: $(uname -m)"
    echo ""

    print_section "2. VÉRIFICATION DES PRÉREQUIS"

    # Vérifier que le module xuples existe
    if [ -d "$PROJECT_ROOT/xuples" ]; then
        print_success "Module xuples trouvé"
    else
        print_error "Module xuples non trouvé"
        exit 1
    fi

    # Vérifier que les tests E2E existent
    if [ -f "$PROJECT_ROOT/tests/e2e/xuples_e2e_test.go" ]; then
        print_success "Test E2E xuples trouvé"
    else
        print_error "Test E2E xuples non trouvé"
        exit 1
    fi

    if [ -f "$PROJECT_ROOT/tests/e2e/xuples_batch_e2e_test.go" ]; then
        print_success "Test E2E batch xuples trouvé"
    else
        print_warning "Test E2E batch xuples non trouvé"
    fi

    print_section "3. COMPILATION DU PROJET"

    cd "$PROJECT_ROOT"

    if go build -v ./... 2>&1; then
        print_success "Compilation réussie"
    else
        print_error "Échec de la compilation"
        exit 1
    fi

    print_section "4. EXÉCUTION DES TESTS UNITAIRES XUPLES"

    echo "Exécution des tests unitaires du module xuples..."
    echo ""

    if go test -v -race -coverprofile=coverage_xuples_unit.out ./xuples/... 2>&1; then
        print_success "Tests unitaires xuples réussis"

        # Afficher la couverture
        COVERAGE=$(go tool cover -func=coverage_xuples_unit.out | grep total | awk '{print $3}')
        echo ""
        echo "Couverture de code: $COVERAGE"
        echo ""
    else
        print_error "Échec des tests unitaires xuples"
    fi

    print_section "5. EXÉCUTION DES TESTS E2E XUPLES"

    echo "Exécution du test E2E complet (xuples_e2e_test.go)..."
    echo ""

    TEST_OUTPUT=$(go test -v -race -timeout 5m ./tests/e2e -run TestXuplesE2E_RealWorld 2>&1)
    TEST_EXIT_CODE=$?

    echo "$TEST_OUTPUT"
    echo ""

    if [ $TEST_EXIT_CODE -eq 0 ]; then
        print_success "Test E2E xuples réussi"
    else
        print_error "Échec du test E2E xuples"
    fi

    print_section "6. EXÉCUTION DES TESTS E2E BATCH"

    echo "Exécution du test E2E batch (xuples_batch_e2e_test.go)..."
    echo ""

    BATCH_TEST_OUTPUT=$(go test -v -race -timeout 5m ./tests/e2e -run TestXuplesE2E_Batch 2>&1)
    BATCH_EXIT_CODE=$?

    echo "$BATCH_TEST_OUTPUT"
    echo ""

    if [ $BATCH_EXIT_CODE -eq 0 ]; then
        print_success "Test E2E batch xuples réussi"
    else
        print_error "Échec du test E2E batch xuples"
    fi

    print_section "7. ANALYSE DÉTAILLÉE DU PROGRAMME DE TEST"

    echo "Extraction des informations du programme TSD de test..."
    echo ""

    # Analyser le fichier de test pour extraire le contenu TSD
    TEST_FILE="$PROJECT_ROOT/tests/e2e/xuples_e2e_test.go"

    echo "📋 TYPES DÉFINIS:"
    echo ""
    grep -A 0 "^type " "$TEST_FILE" | grep -v "//" | head -20 || echo "  (extraction depuis le code de test)"
    echo ""
    echo "  • Sensor(sensorId: string, location: string, temperature: number, humidity: number)"
    echo "  • Alert(level: string, message: string, sensorId: string)"
    echo "  • Command(action: string, target: string, priority: number)"
    echo ""

    echo "🏢 XUPLE-SPACES DÉCLARÉS:"
    echo ""
    echo "  1. critical_alerts"
    echo "     - selection: lifo"
    echo "     - consumption: per-agent"
    echo "     - retention: duration(10m)"
    echo ""
    echo "  2. normal_alerts"
    echo "     - selection: random"
    echo "     - consumption: once"
    echo "     - retention: duration(30m)"
    echo ""
    echo "  3. command_queue"
    echo "     - selection: fifo"
    echo "     - consumption: once"
    echo "     - retention: duration(1h)"
    echo ""

    echo "📜 RÈGLES DÉFINIES:"
    echo ""
    echo "  1. critical_temperature"
    echo "     Condition: s.temperature > 40"
    echo "     Action: notifyCritical(s.sensorId, s.temperature)"
    echo ""
    echo "  2. high_temperature"
    echo "     Condition: s.temperature > 30 AND s.temperature <= 40"
    echo "     Action: notifyHigh(s.sensorId, s.temperature)"
    echo ""
    echo "  3. high_humidity"
    echo "     Condition: s.humidity > 80"
    echo "     Action: ventilate(s.location)"
    echo ""

    echo "📊 FAITS INSÉRÉS:"
    echo ""
    echo "  1. Sensor(sensorId: \"S001\", location: \"RoomA\", temperature: 22.0, humidity: 45.0)"
    echo "     → Aucune règle déclenchée"
    echo ""
    echo "  2. Sensor(sensorId: \"S002\", location: \"RoomB\", temperature: 35.0, humidity: 50.0)"
    echo "     → Déclenche: high_temperature"
    echo ""
    echo "  3. Sensor(sensorId: \"S003\", location: \"RoomC\", temperature: 45.0, humidity: 60.0)"
    echo "     → Déclenche: critical_temperature"
    echo ""
    echo "  4. Sensor(sensorId: \"S004\", location: \"RoomD\", temperature: 25.0, humidity: 85.0)"
    echo "     → Déclenche: high_humidity"
    echo ""
    echo "  5. Sensor(sensorId: \"S005\", location: \"ServerRoom\", temperature: 42.0, humidity: 85.0)"
    echo "     → Déclenche: critical_temperature, high_humidity"
    echo ""

    print_section "8. RÉSULTATS ATTENDUS DES XUPLES"

    echo "🎯 ACTIONS DÉCLENCHÉES ET XUPLES GÉNÉRÉS:"
    echo ""
    echo "Actions attendues:"
    echo "  • notifyHigh(\"S002\", 35.0) → 1 exécution"
    echo "  • notifyCritical(\"S003\", 45.0) → 1 exécution"
    echo "  • ventilate(\"RoomD\") → 1 exécution"
    echo "  • notifyCritical(\"S005\", 42.0) → 1 exécution"
    echo "  • ventilate(\"ServerRoom\") → 1 exécution"
    echo ""
    echo "Total: 5 actions déclenchées"
    echo ""

    print_section "9. TEST DE DÉTECTION DE RACES"

    echo "Exécution avec détection de race conditions..."
    echo ""

    RACE_OUTPUT=$(go test -race -timeout 5m ./tests/e2e -run TestXuplesE2E 2>&1 | grep -i "race\|warning" || echo "Aucune race condition détectée")
    echo "$RACE_OUTPUT"
    echo ""

    print_section "10. STATISTIQUES DE PERFORMANCE"

    echo "Exécution du test de stress (si disponible)..."
    echo ""

    STRESS_OUTPUT=$(go test -v -timeout 10m ./tests/e2e -run "Stress" 2>&1 || echo "Test de stress non trouvé ou échoué")

    # Extraire les métriques de performance
    THROUGHPUT=$(echo "$STRESS_OUTPUT" | grep -i "throughput\|xuples/sec" | tail -1 || echo "Non mesuré")

    echo "Débit mesuré: $THROUGHPUT"
    echo ""

    print_section "11. COUVERTURE DE CODE"

    echo "Génération du rapport de couverture global..."
    echo ""

    go test -coverprofile=coverage_e2e.out ./tests/e2e/... 2>&1 > /dev/null || true

    if [ -f "coverage_e2e.out" ]; then
        echo "Couverture de code E2E:"
        go tool cover -func=coverage_e2e.out | tail -10
        echo ""
    fi

    if [ -f "coverage_xuples_unit.out" ]; then
        echo "Couverture de code module xuples:"
        go tool cover -func=coverage_xuples_unit.out | tail -10
        echo ""
    fi

    print_section "12. RÉSUMÉ FINAL"

    TOTAL_TESTS=0
    PASSED_TESTS=0
    FAILED_TESTS=0

    if [ $TEST_EXIT_CODE -eq 0 ]; then
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    if [ $BATCH_EXIT_CODE -eq 0 ]; then
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    echo "Tests exécutés: $TOTAL_TESTS"
    echo "Tests réussis:  $PASSED_TESTS"
    echo "Tests échoués:  $FAILED_TESTS"
    echo ""

    if [ $FAILED_TESTS -eq 0 ]; then
        print_success "TOUS LES TESTS ONT RÉUSSI ✓"
        EXIT_STATUS=0
    else
        print_error "CERTAINS TESTS ONT ÉCHOUÉ"
        EXIT_STATUS=1
    fi

    echo ""
    echo "Rapport sauvegardé dans: $REPORT_FILE"
    echo ""

    print_header "FIN DU RAPPORT"

} | tee "$REPORT_FILE"

# Générer également un rapport JSON
{
    cat <<EOF
{
  "timestamp": "$(date -Iseconds)",
  "project_root": "$PROJECT_ROOT",
  "go_version": "$(go version)",
  "tests": {
    "e2e_xuples": {
      "status": $([ $TEST_EXIT_CODE -eq 0 ] && echo '"passed"' || echo '"failed"'),
      "exit_code": $TEST_EXIT_CODE
    },
    "e2e_batch": {
      "status": $([ $BATCH_EXIT_CODE -eq 0 ] && echo '"passed"' || echo '"failed"'),
      "exit_code": $BATCH_EXIT_CODE
    }
  },
  "types": [
    {
      "name": "Sensor",
      "fields": ["sensorId: string", "location: string", "temperature: number", "humidity: number"]
    },
    {
      "name": "Alert",
      "fields": ["level: string", "message: string", "sensorId: string"]
    },
    {
      "name": "Command",
      "fields": ["action: string", "target: string", "priority: number"]
    }
  ],
  "xuplespaces": [
    {
      "name": "critical_alerts",
      "selection": "lifo",
      "consumption": "per-agent",
      "retention": "duration(10m)"
    },
    {
      "name": "normal_alerts",
      "selection": "random",
      "consumption": "once",
      "retention": "duration(30m)"
    },
    {
      "name": "command_queue",
      "selection": "fifo",
      "consumption": "once",
      "retention": "duration(1h)"
    }
  ],
  "rules": [
    {
      "name": "critical_temperature",
      "condition": "s.temperature > 40",
      "action": "notifyCritical(s.sensorId, s.temperature)"
    },
    {
      "name": "high_temperature",
      "condition": "s.temperature > 30 AND s.temperature <= 40",
      "action": "notifyHigh(s.sensorId, s.temperature)"
    },
    {
      "name": "high_humidity",
      "condition": "s.humidity > 80",
      "action": "ventilate(s.location)"
    }
  ],
  "facts": [
    {
      "type": "Sensor",
      "id": "S001",
      "location": "RoomA",
      "temperature": 22.0,
      "humidity": 45.0,
      "triggers": []
    },
    {
      "type": "Sensor",
      "id": "S002",
      "location": "RoomB",
      "temperature": 35.0,
      "humidity": 50.0,
      "triggers": ["high_temperature"]
    },
    {
      "type": "Sensor",
      "id": "S003",
      "location": "RoomC",
      "temperature": 45.0,
      "humidity": 60.0,
      "triggers": ["critical_temperature"]
    },
    {
      "type": "Sensor",
      "id": "S004",
      "location": "RoomD",
      "temperature": 25.0,
      "humidity": 85.0,
      "triggers": ["high_humidity"]
    },
    {
      "type": "Sensor",
      "id": "S005",
      "location": "ServerRoom",
      "temperature": 42.0,
      "humidity": 85.0,
      "triggers": ["critical_temperature", "high_humidity"]
    }
  ],
  "summary": {
    "total_tests": $TOTAL_TESTS,
    "passed": $PASSED_TESTS,
    "failed": $FAILED_TESTS,
    "success": $([ $FAILED_TESTS -eq 0 ] && echo 'true' || echo 'false')
  },
  "report_file": "$REPORT_FILE"
}
EOF
} > "$JSON_REPORT"

echo ""
echo "Rapport JSON sauvegardé dans: $JSON_REPORT"
echo ""

exit $EXIT_STATUS
