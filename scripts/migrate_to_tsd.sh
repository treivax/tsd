#!/bin/bash

# Script to migrate .constraint and .facts files to unified .tsd extension
# This script:
# 1. Finds all .constraint and .facts files
# 2. Merges paired files (same basename, different extensions)
# 3. Renames standalone files
# 4. Creates .tsd files with the merged/renamed content

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
total_constraint_files=0
total_facts_files=0
merged_files=0
renamed_constraint_files=0
renamed_facts_files=0
errors=0

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}TSD File Extension Migration${NC}"
echo -e "${BLUE}================================${NC}"
echo ""
echo "Project root: $PROJECT_ROOT"
echo ""

# Function to merge files
merge_files() {
    local constraint_file="$1"
    local facts_file="$2"
    local tsd_file="$3"

    echo -e "${YELLOW}  Merging:${NC}"
    echo "    - $constraint_file"
    echo "    - $facts_file"
    echo -e "    ${GREEN}→${NC} $tsd_file"

    # Create merged file: constraint content + newline + facts content
    cat "$constraint_file" > "$tsd_file"
    echo "" >> "$tsd_file"  # Add blank line separator
    cat "$facts_file" >> "$tsd_file"

    # Remove original files
    rm "$constraint_file"
    rm "$facts_file"

    merged_files=$((merged_files + 1))
}

# Function to rename single file
rename_file() {
    local source_file="$1"
    local dest_file="$2"
    local file_type="$3"

    echo -e "${YELLOW}  Renaming:${NC} $source_file ${GREEN}→${NC} $dest_file"

    mv "$source_file" "$dest_file"

    if [ "$file_type" = "constraint" ]; then
        renamed_constraint_files=$((renamed_constraint_files + 1))
    else
        renamed_facts_files=$((renamed_facts_files + 1))
    fi
}

# Find all unique basenames
declare -A basenames

echo -e "${BLUE}Scanning for .constraint and .facts files...${NC}"
echo ""

# Collect all .constraint files
while IFS= read -r -d '' file; do
    total_constraint_files=$((total_constraint_files + 1))
    basename="${file%.constraint}"
    basenames["$basename"]=1
done < <(find "$PROJECT_ROOT" -type f -name "*.constraint" -print0)

# Collect all .facts files
while IFS= read -r -d '' file; do
    total_facts_files=$((total_facts_files + 1))
    basename="${file%.facts}"
    basenames["$basename"]=1
done < <(find "$PROJECT_ROOT" -type f -name "*.facts" -print0)

echo "Found:"
echo "  - $total_constraint_files .constraint files"
echo "  - $total_facts_files .facts files"
echo ""

# Process each unique basename
echo -e "${BLUE}Processing files...${NC}"
echo ""

for basename in "${!basenames[@]}"; do
    constraint_file="${basename}.constraint"
    facts_file="${basename}.facts"
    tsd_file="${basename}.tsd"

    # Skip if .tsd already exists
    if [ -f "$tsd_file" ]; then
        echo -e "${YELLOW}⚠️  Skipping $basename - .tsd file already exists${NC}"
        continue
    fi

    # Check which files exist
    has_constraint=false
    has_facts=false

    [ -f "$constraint_file" ] && has_constraint=true
    [ -f "$facts_file" ] && has_facts=true

    if [ "$has_constraint" = true ] && [ "$has_facts" = true ]; then
        # Both files exist - merge them
        echo "📦 $(basename "$basename")"
        merge_files "$constraint_file" "$facts_file" "$tsd_file"
    elif [ "$has_constraint" = true ]; then
        # Only constraint file exists - rename it
        echo "📝 $(basename "$basename")"
        rename_file "$constraint_file" "$tsd_file" "constraint"
    elif [ "$has_facts" = true ]; then
        # Only facts file exists - rename it
        echo "📝 $(basename "$basename")"
        rename_file "$facts_file" "$tsd_file" "facts"
    fi
done

echo ""
echo -e "${BLUE}================================${NC}"
echo -e "${GREEN}Migration Complete!${NC}"
echo -e "${BLUE}================================${NC}"
echo ""
echo "Summary:"
echo "  - Merged file pairs: $merged_files"
echo "  - Renamed .constraint files: $renamed_constraint_files"
echo "  - Renamed .facts files: $renamed_facts_files"
echo "  - Total .tsd files created: $((merged_files + renamed_constraint_files + renamed_facts_files))"
echo ""

if [ $errors -gt 0 ]; then
    echo -e "${RED}⚠️  $errors errors occurred during migration${NC}"
    exit 1
else
    echo -e "${GREEN}✓ All files migrated successfully!${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. Run tests: make test"
    echo "  2. Update code references to use .tsd extension"
    echo "  3. Update documentation"
    echo "  4. Commit changes"
fi
