#!/bin/bash

# Génère un fichier JSON avec toutes les métriques du code

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
    "global": "48.7%",
    "packages": {
      "rete/pkg/domain": "100.0%",
      "rete/pkg/network": "100.0%",
      "constraint/pkg/validator": "96.5%",
      "constraint/pkg/domain": "90.0%",
      "rete/pkg/nodes": "71.6%",
      "constraint": "59.6%",
      "rete": "39.7%",
      "test/integration": "29.4%",
      "cmd/tsd": "0.0%",
      "cmd/universal-rete-runner": "0.0%",
      "constraint/cmd": "0.0%",
      "constraint/internal/config": "0.0%",
      "rete/internal/config": "0.0%",
      "scripts": "0.0%",
      "test/testutil": "0.0%"
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
      {"file": "constraint/coverage_test.go", "lines": 1395},
      {"file": "rete/pkg/nodes/advanced_beta_test.go", "lines": 1292},
      {"file": "constraint/pkg/validator/types_test.go", "lines": 886},
      {"file": "constraint/pkg/validator/validator_test.go", "lines": 880},
      {"file": "constraint/pkg/domain/types_test.go", "lines": 743}
    ]
  },
  "priorities": {
    "high": [
      "cmd/tsd (0% -> 80%)",
      "cmd/universal-rete-runner (0% -> 70%)"
    ],
    "medium": [
      "rete (39.7% -> 70%)",
      "constraint (59.6% -> 75%)",
      "rete/pkg/nodes (71.6% -> 90%)"
    ],
    "low": [
      "constraint/internal/config (0% -> 80%)",
      "rete/internal/config (0% -> 80%)",
      "test/integration (29.4% -> 60%)"
    ]
  }
}
JSONEOF

echo "Métriques générées dans docs/reports/code_metrics.json"
