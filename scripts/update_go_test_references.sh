#!/bin/bash

# Script to update Go test files to use .tsd extension instead of .constraint and .facts
# This script updates string references in test files

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}Update Go Test References${NC}"
echo -e "${BLUE}================================${NC}"
echo ""
echo "Project root: $PROJECT_ROOT"
echo ""

# Counter
files_updated=0

echo -e "${BLUE}Finding and updating Go test files...${NC}"
echo ""

# Find all Go test files
while IFS= read -r -d '' file; do
    # Check if file contains .constraint or .facts references
    if grep -q '\.constraint\|\.facts' "$file"; then
        echo -e "${YELLOW}Updating:${NC} $file"

        # Create backup
        cp "$file" "${file}.bak"

        # Update references:
        # 1. .constraint" -> .tsd"
        # 2. .facts" -> .tsd"
        # 3. Update field names (ConstraintFile -> File where appropriate)
        sed -i \
            -e 's/\.constraint"/\.tsd"/g' \
            -e 's/\.facts"/\.tsd"/g' \
            "$file"

        # Check if changes were successful
        if [ $? -eq 0 ]; then
            rm "${file}.bak"
            files_updated=$((files_updated + 1))
            echo -e "  ${GREEN}✓${NC} Updated"
        else
            # Restore backup on error
            mv "${file}.bak" "$file"
            echo -e "  ${RED}✗${NC} Failed to update, restored backup"
        fi
    fi
done < <(find "$PROJECT_ROOT" -type f -name "*_test.go" -print0)

echo ""
echo -e "${BLUE}================================${NC}"
echo -e "${GREEN}Update Complete!${NC}"
echo -e "${BLUE}================================${NC}"
echo ""
echo "Summary:"
echo "  - Go test files updated: $files_updated"
echo ""

if [ $files_updated -gt 0 ]; then
    echo -e "${GREEN}✓ Test files updated successfully!${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. Review changes: git diff"
    echo "  2. Run tests: go test ./..."
    echo "  3. Fix any remaining issues manually"
else
    echo -e "${YELLOW}No test files needed updating${NC}"
fi
