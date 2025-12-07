# Documentation Cleanup and Consolidation Plan

## New Consolidated Documentation Structure

### Core Documents (COMPLETED)
✅ `INSTALLATION.md` - Complete installation guide
✅ `QUICK_START.md` - 5-minute quick start
✅ `USER_GUIDE.md` - Comprehensive user guide (all features)

### To Create
- `ARCHITECTURE.md` - Technical architecture
- `CONTRIBUTING.md` - Contribution guidelines

### Existing to Keep
- `README.md` - Project overview
- `GRAMMAR_GUIDE.md` - Language grammar reference
- `API_REFERENCE.md` - HTTP API documentation
- `AUTHENTICATION.md` - Authentication guide
- `TUTORIAL.md` - Step-by-step tutorial

### To Remove/Consolidate
- `AUTHENTICATION_DIAGRAMS.md` → Merge into AUTHENTICATION.md
- `AUTHENTICATION_QUICKSTART.md` → Covered in QUICK_START.md
- `AUTHENTICATION_TUTORIAL.md` → Merge into AUTHENTICATION.md
- `EXAMPLES.md` → Examples directory is self-documenting
- `FEATURES.md` → Covered in USER_GUIDE.md
- `LOGGING_GUIDE.md` → Keep as separate reference
- `OPTIMIZATIONS.md` → Merge into ARCHITECTURE.md
- `PROMETHEUS_INTEGRATION.md` → Merge into ARCHITECTURE.md
- `STRONG_MODE_TUNING_GUIDE.md` → Merge into USER_GUIDE.md
- `TLS_CONFIGURATION.md` → Merge into AUTHENTICATION.md
- `TRANSACTION_ARCHITECTURE.md` → Merge into ARCHITECTURE.md
- `TRANSACTION_README.md` → Merge into ARCHITECTURE.md
- `UNIFIED_BINARY.md` → Covered in QUICK_START.md and INSTALLATION.md
- `development_guidelines.md` → Merge into CONTRIBUTING.md
- `feature-type-casting.md` → Covered in USER_GUIDE.md
- `fix-case-insensitive-keywords.md` → Implementation detail, remove
- `quick-start-case-insensitive.md` → Covered in QUICK_START.md
- `type-casting.md` → Covered in USER_GUIDE.md
- `utf8-and-identifier-styles.md` → Covered in GRAMMAR_GUIDE.md

## Cleanup Actions

### 1. Remove obsolete scripts
✅ `validate_advanced_features.sh` - REMOVED

### 2. Remove duplicate/scattered documentation
Find and remove:
```bash
find . -name "*README*.md" -not -path "./docs/*" -not -path "./README.md"
find . -name "*GUIDE*.md" -not -path "./docs/*"
```

### 3. Consolidate scattered examples
- Keep `examples/` directory
- Remove example READMEs that duplicate main docs

### 4. Remove build artifacts
```bash
make clean
go clean -cache -testcache -modcache
```

### 5. Remove temporary files
```bash
find . -name "*.tmp" -o -name "*.bak" -o -name "*~"
```

## Final Documentation Structure

```
docs/
├── README.md                   # Documentation index
├── INSTALLATION.md             # How to install
├── QUICK_START.md              # 5-minute start
├── USER_GUIDE.md               # Complete feature guide
├── TUTORIAL.md                 # Step-by-step learning
├── GRAMMAR_GUIDE.md            # Language syntax
├── API_REFERENCE.md            # HTTP API
├── AUTHENTICATION.md           # Auth complete guide
├── ARCHITECTURE.md             # Technical design
├── LOGGING_GUIDE.md            # Logging reference
├── CONTRIBUTING.md             # How to contribute
```

All other documents will be removed or consolidated into the above.
