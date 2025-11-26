#!/bin/bash

echo "=== Détail de la couverture par package ==="
echo

go tool cover -func=coverage.out | grep -E 'github.com/treivax/tsd/(cmd|constraint|rete)' | \
    awk -F: '{print $1}' | \
    sed 's|github.com/treivax/tsd/||' | \
    sort -u | while read pkg; do
    coverage=$(go tool cover -func=coverage.out | grep "^$pkg" | tail -1 | awk '{print $NF}')
    if [ -n "$coverage" ]; then
        echo "$pkg: $coverage"
    fi
done

echo
echo "=== Fichiers sans couverture (0.0%) ==="
go tool cover -func=coverage.out | grep '0.0%' | grep -v test | head -20
