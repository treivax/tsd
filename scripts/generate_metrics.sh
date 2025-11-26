#!/bin/bash

# Génère un fichier JSON avec toutes les métriques du code

# Fonction pour obtenir la couverture d'un package
get_coverage() {
    local pkg=$1
    local coverage=$(go test -cover "./$pkg" 2>/dev/null | grep -o 'coverage: [0-9.]*%' | grep -o '[0-9.]*')
    if [ -z "$coverage" ]; then
        echo "0.0"
    else
        echo "$coverage"
    fi
}

# Calculer la couverture de chaque package
CMD_TSD_COV=$(get_coverage "cmd/tsd")
CMD_RUNNER_COV=$(get_coverage "cmd/universal-rete-runner")
CONSTRAINT_COV=$(get_coverage "constraint")
CONSTRAINT_VALIDATOR_COV=$(get_coverage "constraint/pkg/validator")
CONSTRAINT_DOMAIN_COV=$(get_coverage "constraint/pkg/domain")
CONSTRAINT_CMD_COV=$(get_coverage "constraint/cmd")
CONSTRAINT_CONFIG_COV=$(get_coverage "constraint/internal/config")
RETE_COV=$(get_coverage "rete")
RETE_DOMAIN_COV=$(get_coverage "rete/pkg/domain")
RETE_NETWORK_COV=$(get_coverage "rete/pkg/network")
RETE_NODES_COV=$(get_coverage "rete/pkg/nodes")
RETE_CONFIG_COV=$(get_coverage "rete/internal/config")
TEST_INTEGRATION_COV=$(get_coverage "test/integration")
TEST_TESTUTIL_COV=$(get_coverage "test/testutil")

# Calculer la couverture globale (moyenne pondérée approximative)
# On utilise les principaux packages
GLOBAL_COV=$(echo "scale=1; ($CONSTRAINT_COV + $RETE_COV + $CMD_TSD_COV) / 3" | bc)

cat > docs/reports/code_metrics.json << JSONEOF
{
  "generated_at": "$(date -Iseconds)",
  "commit": "$(git rev-parse HEAD)",
  "branch": "$(git branch --show-current)",
  "metrics": {
    "total_lines_go": $(find . -name "*.go" -not -path "*/vendor/*" -exec cat {} \; | wc -l),
    "test_lines": $(find . -name "*_test.go" -not -path "*/vendor/*" -exec cat {} \; | wc -l),
    "generated_lines": $(find . -name "parser.go" -path "*/constraint/*" -exec cat {} \; | wc -l),
    "manual_lines": $(($(find . -name "*.go" -not -path "*/vendor/*" -exec cat {} \; | wc -l) - $(find . -name "*_test.go" -not -path "*/vendor/*" -exec cat {} \; | wc -l) - $(find . -name "parser.go" -path "*/constraint/*" -exec cat {} \; | wc -l))),
    "total_files": $(find . -name "*.go" -not -path "*/vendor/*" | wc -l),
    "test_files": $(find . -name "*_test.go" -not -path "*/vendor/*" | wc -l),
    "production_files": $(($(find . -name "*.go" -not -path "*/vendor/*" | wc -l) - $(find . -name "*_test.go" -not -path "*/vendor/*" | wc -l)))
  },
  "coverage": {
    "global": "${GLOBAL_COV}%",
    "packages": {
      "cmd/tsd": "${CMD_TSD_COV}%",
      "cmd/universal-rete-runner": "${CMD_RUNNER_COV}%",
      "constraint": "${CONSTRAINT_COV}%",
      "constraint/pkg/validator": "${CONSTRAINT_VALIDATOR_COV}%",
      "constraint/pkg/domain": "${CONSTRAINT_DOMAIN_COV}%",
      "constraint/cmd": "${CONSTRAINT_CMD_COV}%",
      "constraint/internal/config": "${CONSTRAINT_CONFIG_COV}%",
      "rete": "${RETE_COV}%",
      "rete/pkg/domain": "${RETE_DOMAIN_COV}%",
      "rete/pkg/network": "${RETE_NETWORK_COV}%",
      "rete/pkg/nodes": "${RETE_NODES_COV}%",
      "rete/internal/config": "${RETE_CONFIG_COV}%",
      "test/integration": "${TEST_INTEGRATION_COV}%",
      "test/testutil": "${TEST_TESTUTIL_COV}%"
    }
  },
  "largest_files": {
    "production": [
      {"file": "constraint/parser.go", "lines": 5230, "type": "generated"},
      {"file": "rete/pkg/nodes/advanced_beta.go", "lines": 689, "type": "manual"},
      {"file": "rete/constraint_pipeline_builder.go", "lines": 617, "type": "manual"},
      {"file": "constraint/constraint_utils.go", "lines": 617, "type": "manual"},
      {"file": "rete/node_join.go", "lines": 445, "type": "manual"}
    ],
    "tests": [
      {"file": "cmd/tsd/main_test.go", "lines": $(wc -l < cmd/tsd/main_test.go 2>/dev/null || echo 0)},
      {"file": "constraint/coverage_test.go", "lines": $(wc -l < constraint/coverage_test.go 2>/dev/null || echo 0)},
      {"file": "rete/pkg/nodes/advanced_beta_test.go", "lines": $(wc -l < rete/pkg/nodes/advanced_beta_test.go 2>/dev/null || echo 0)},
      {"file": "constraint/pkg/validator/types_test.go", "lines": $(wc -l < constraint/pkg/validator/types_test.go 2>/dev/null || echo 0)},
      {"file": "constraint/pkg/validator/validator_test.go", "lines": $(wc -l < constraint/pkg/validator/validator_test.go 2>/dev/null || echo 0)}
    ]
  },
  "priorities": {
    "high": [
      "cmd/universal-rete-runner (${CMD_RUNNER_COV}% -> 70%)"
    ],
    "medium": [
      "rete (${RETE_COV}% -> 70%)",
      "constraint (${CONSTRAINT_COV}% -> 75%)",
      "rete/pkg/nodes (${RETE_NODES_COV}% -> 90%)"
    ],
    "low": [
      "constraint/internal/config (${CONSTRAINT_CONFIG_COV}% -> 80%)",
      "rete/internal/config (${RETE_CONFIG_COV}% -> 80%)",
      "test/integration (${TEST_INTEGRATION_COV}% -> 60%)"
    ],
    "completed": [
      "cmd/tsd (${CMD_TSD_COV}%) ✅"
    ]
  }
}
JSONEOF

echo "✅ Métriques générées dans docs/reports/code_metrics.json"
echo ""
echo "📊 Couverture par package:"
echo "  cmd/tsd: ${CMD_TSD_COV}%"
echo "  cmd/universal-rete-runner: ${CMD_RUNNER_COV}%"
echo "  constraint: ${CONSTRAINT_COV}%"
echo "  rete: ${RETE_COV}%"
echo "  Global: ${GLOBAL_COV}%"
