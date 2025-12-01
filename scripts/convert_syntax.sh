#!/bin/bash
# Copyright (c) 2025 TSD Contributors
# Licensed under the MIT License
# See LICENSE file in the project root for full license text

# Script to convert old type syntax to new syntax
# Old: type Person : <name: string, age: number>
# New: type Person(name: string, age: number)

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "🔄 Converting TSD files to new syntax..."
echo "========================================"

# Counter for files processed
total_files=0
converted_files=0
error_files=0

# Find all .tsd and .constraint files
while IFS= read -r file; do
    total_files=$((total_files + 1))
    echo -n "📝 Processing: $file ... "

    # Create backup
    cp "$file" "$file.backup"

    # Convert type syntax: type Name : <field: type, ...> -> type Name(field: type, ...)
    # Using perl for better regex support
    if perl -i -pe 's/type\s+(\w+)\s*:\s*<\s*(.*?)\s*>/type $1($2)/g' "$file"; then
        echo -e "${GREEN}✓${NC}"
        converted_files=$((converted_files + 1))
    else
        echo -e "${RED}✗${NC}"
        error_files=$((error_files + 1))
        # Restore backup on error
        mv "$file.backup" "$file"
    fi
done < <(find . -type f \( -name "*.tsd" -o -name "*.constraint" \) ! -path "./.git/*" ! -path "*/vendor/*")

echo ""
echo "========================================"
echo "📊 Conversion Summary:"
echo "   Total files: $total_files"
echo -e "   ${GREEN}Converted: $converted_files${NC}"
if [ $error_files -gt 0 ]; then
    echo -e "   ${RED}Errors: $error_files${NC}"
fi

# Clean up backup files if successful
if [ $error_files -eq 0 ]; then
    echo ""
    echo "🧹 Cleaning up backup files..."
    find . -type f \( -name "*.tsd.backup" -o -name "*.constraint.backup" \) -delete
    echo -e "${GREEN}✅ All files converted successfully!${NC}"
else
    echo ""
    echo -e "${YELLOW}⚠️  Some files had errors. Backups kept as .backup files.${NC}"
    exit 1
fi
