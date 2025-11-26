#!/bin/bash

echo "=== Fonctions longues (>50 lignes, hors tests) ==="

find . -name "*.go" -not -name "*_test.go" -not -path "*/vendor/*" -not -path "*/parser.go" | while read file; do
    awk -v file="$file" '
    /^func / {
        if (in_func && func_name != "") {
            lines = NR - start_line
            if (lines > 50) {
                printf "%d:%s:%s\n", lines, file, func_name
            }
        }
        in_func = 1
        func_name = $2
        gsub(/\(.*/, "", func_name)
        start_line = NR
    }
    /^}$/ && NF == 1 {
        if (in_func && func_name != "") {
            lines = NR - start_line + 1
            if (lines > 50) {
                printf "%d:%s:%s\n", lines, file, func_name
            }
        }
        in_func = 0
        func_name = ""
    }
    END {
        if (in_func && func_name != "") {
            lines = NR - start_line + 1
            if (lines > 50) {
                printf "%d:%s:%s\n", lines, file, func_name
            }
        }
    }
    ' "$file"
done | sort -t: -k1 -rn | head -20
