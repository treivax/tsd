#!/bin/bash

echo "=== Analyse de complexité ==="
echo

echo "=== Fonctions avec plus de 50 lignes (hors tests) ==="
for file in $(find . -name "*.go" -not -name "*_test.go" -not -path "*/vendor/*"); do
    gawk '
    /^func / {
        if (in_func && func_name != "") {
            lines = NR - start_line
            if (lines > 50) {
                printf "%s:%s:%d lignes\n", FILENAME, func_name, lines
            }
        }
        in_func = 1
        func_name = $0
        start_line = NR
    }
    /^}$/ && in_func && NF == 1 {
        lines = NR - start_line
        if (lines > 50) {
            printf "%s:%s:%d lignes\n", FILENAME, func_name, lines
        }
        in_func = 0
        func_name = ""
    }
    END {
        if (in_func && func_name != "") {
            lines = NR - start_line
            if (lines > 50) {
                printf "%s:%s:%d lignes\n", FILENAME, func_name, lines
            }
        }
    }
    ' "$file"
done | sort -t: -k3 -rn | head -20

echo
echo "=== Packages et leur taille ==="
find . -name "*.go" -not -name "*_test.go" -not -path "*/vendor/*" -not -path "*/parser.go" | \
    xargs -I {} dirname {} | sort | uniq -c | sort -rn | head -15
