# TODO - Validation Layer Improvements

**Date** : 2025-12-11
**Module** : constraint/validation
**Priorité** : MOYENNE à BASSE
**Status** : Post-refactoring security fixes

---

## ✅ Completed (Session 2)

### Sécurité
- [x] BASE64_DECODE_NO_VALIDATION - Ajout safeBase64Decode avec limite 1MB
- [x] MISSING_NIL_CHECKS - Ajout validateInputNotNil
- [x] RECURSIVE_VALIDATION_NO_DEPTH_LIMIT - MaxValidationDepth = 100
- [x] NO_INPUT_SANITIZATION - Ajout sanitizeForLog

### Maintenabilité
- [x] HARDCODED_FUNCTION_LIST - FunctionRegistry extensible
- [x] MAGIC_STRING_COMPARISONS - Constantes pour tous types/opérateurs
- [x] DUPLICATE_PATTERN_LOGIC - Helper findVariableType
- [x] UNUSED_FUNCTION - Suppression validateBinaryOpConstraint
- [x] INCONSISTENT_NAMING - Messages d'erreur en anglais
- [x] FRENCH_COMMENTS_IN_ERRORS - Traduction complète

---

## 📋 Remaining TODO Items

### HIGH Priority (2 semaines)

#### 1. Type Compatibility Rules Enhancement
**File**: `constraint/action_validator.go`, `constraint/constraint_type_checking.go`
**Issue**: `INCOMPLETE_TYPE_COMPATIBILITY` (#8)
**Description**: Vérification de compatibilité limitée aux correspondances exactes

**Current Behavior**:
```go
func (av *ActionValidator) isTypeCompatible(argType, paramType string) bool {
    // Exact match only
    if argType == paramType {
        return true
    }
    // User-defined types require exact match
    // No coercion for primitives
    return false
}
```

**Proposed Enhancement**:
```go
type TypeCompatibilityRule struct {
    From       string
    To         string
    Coercible  bool
    Validator  func(value interface{}) error
}

type TypeCompatibilityRegistry struct {
    rules map[string][]TypeCompatibilityRule
}

// Allow number -> int, float -> number, etc.
// With proper validation of value ranges
```

**Tasks**:
- [ ] Define compatibility matrix (number/int/float/decimal)
- [ ] Implement TypeCompatibilityRegistry
- [ ] Add coercion validation
- [ ] Update isTypeCompatible to use registry
- [ ] Add comprehensive tests
- [ ] Document compatibility rules

**Estimated Effort**: 8h

---

#### 2. Error Context Enhancement
**File**: `constraint/errors.go`, all validation files
**Issue**: `ERROR_CONTEXT_INCONSISTENT` (#9)
**Description**: Messages d'erreur parfois sans numéro de ligne

**Current Behavior**:
```go
return fmt.Errorf("action '%s' is not defined", actionName)
// No file, line, or context information
```

**Proposed Enhancement**:
```go
type ValidationContext struct {
    File           string
    Line           int
    Column         int
    ExpressionID   string
    RuleName       string
}

type ContextualError struct {
    Context ValidationContext
    Message string
    Cause   error
}

func (ctx *ValidationContext) Errorf(format string, args ...interface{}) error {
    return &ContextualError{
        Context: *ctx,
        Message: fmt.Sprintf(format, args...),
    }
}
```

**Tasks**:
- [ ] Define ContextualError type
- [ ] Add ValidationContext to all validation functions
- [ ] Extract line/column from AST nodes
- [ ] Update error construction everywhere
- [ ] Improve error messages format
- [ ] Add tests for error formatting

**Estimated Effort**: 6h

---

#### 3. Documentation - Validation Assumptions
**Files**: `constraint/docs/validation.md` (create)
**Issue**: Missing documentation (#14)
**Description**: Documenter hypothèses et invariants de validation

**Tasks**:
- [ ] Document validation order and dependencies
- [ ] Document type system assumptions
- [ ] Document recursion limits rationale
- [ ] Document security considerations
- [ ] Add validation flow diagrams
- [ ] Document extension points (registry)
- [ ] Add examples of validation failures
- [ ] Document error codes/categories

**Estimated Effort**: 4h

---

### MEDIUM Priority (1 mois)

#### 4. Validation Result Caching
**File**: `constraint/validation_cache.go` (create)
**Issue**: `NO_VALIDATION_CACHING` (#10)
**Description**: Pas de cache pour résultats de validation répétés

**Proposed Implementation**:
```go
type ValidationCache struct {
    mu              sync.RWMutex
    typeFieldCache  map[string][]Field       // typeName -> fields
    typeValidCache  map[string]bool          // typeName -> isValid
    fieldTypeCache  map[string]string        // "typeName.fieldName" -> type
    maxSize         int
    enabled         bool
}

func (vc *ValidationCache) GetTypeFields(program Program, typeName string) ([]Field, bool) {
    vc.mu.RLock()
    defer vc.mu.RUnlock()
    fields, ok := vc.typeFieldCache[typeName]
    return fields, ok
}
```

**Tasks**:
- [ ] Design cache key strategy
- [ ] Implement LRU or size-limited cache
- [ ] Add cache to ActionValidator
- [ ] Invalidation strategy for dynamic updates
- [ ] Benchmarks cache hit/miss
- [ ] Performance comparison with/without cache
- [ ] Configuration for cache size/enable

**Estimated Effort**: 6h

---

#### 5. Performance Benchmarks
**File**: `constraint/validation_bench_test.go` (create)
**Issue**: No performance baseline
**Description**: Benchmarks pour optimisations futures

**Tasks**:
- [ ] Benchmark ValidateActionCall (simple)
- [ ] Benchmark ValidateActionCall (complex, nested)
- [ ] Benchmark ValidateFieldAccess (shallow)
- [ ] Benchmark ValidateFieldAccess (deep recursion)
- [ ] Benchmark ValidateTypeCompatibility
- [ ] Benchmark with different cache strategies
- [ ] Profile CPU/memory with pprof
- [ ] Document performance characteristics

**Estimated Effort**: 4h

---

#### 6. Integration Tests E2E
**File**: `tests/integration/validation_test.go` (create)
**Issue**: Tests unitaires seulement
**Description**: Tests validation end-to-end

**Tasks**:
- [ ] Test complet fichier TSD avec erreurs multiples
- [ ] Test validation multi-fichiers
- [ ] Test validation avec imports
- [ ] Test performance gros fichiers (>1000 rules)
- [ ] Test cas edge combinés
- [ ] Test robustesse avec fichiers malformés
- [ ] Automatic generation of test cases

**Estimated Effort**: 8h

---

### LOW Priority (Opportuniste)

#### 7. Fuzzing Tests
**File**: `constraint/validation_fuzz_test.go` (create)
**Issue**: Pas de fuzzing
**Description**: Tests fuzzing pour robustesse extrême

**Tasks**:
- [ ] Setup go-fuzz or native fuzzing
- [ ] Fuzz inferArgumentType with random AST
- [ ] Fuzz ValidateFieldAccess with random programs
- [ ] Fuzz base64 decode with random payloads
- [ ] Fuzz deeply nested structures
- [ ] Analyze crash reports
- [ ] Fix discovered edge cases

**Estimated Effort**: 8h

---

#### 8. Metrics and Observability
**File**: `constraint/validation_metrics.go` (create)
**Issue**: Pas de métriques
**Description**: Métriques de validation pour monitoring

**Proposed**:
```go
type ValidationMetrics struct {
    TotalValidations      int64
    FailedValidations     int64
    ValidationErrors      map[string]int64  // errorType -> count
    AverageDepth          float64
    MaxDepthReached       int
    CacheHitRate          float64
    ValidationDuration    time.Duration
}

func (vm *ValidationMetrics) RecordValidation(err error, depth int, duration time.Duration)
func (vm *ValidationMetrics) Report() string
```

**Tasks**:
- [ ] Define metric types
- [ ] Instrument validation functions
- [ ] Export metrics (Prometheus format?)
- [ ] Dashboard queries examples
- [ ] Performance impact measurement

**Estimated Effort**: 6h

---

## 🎯 Recommended Roadmap

### Sprint 1 (2 semaines) - High Priority
1. Type Compatibility Rules Enhancement (8h)
2. Error Context Enhancement (6h)
3. Documentation - Validation Assumptions (4h)

**Total**: 18h (~2.5 jours)

### Sprint 2 (1 mois) - Medium Priority
4. Validation Result Caching (6h)
5. Performance Benchmarks (4h)
6. Integration Tests E2E (8h)

**Total**: 18h (~2.5 jours)

### Sprint 3 (Opportuniste) - Low Priority
7. Fuzzing Tests (8h)
8. Metrics and Observability (6h)

**Total**: 14h (~2 jours)

**Total Effort**: ~50h (~7 jours)

---

## 📝 Notes

### Dependencies
- Item #4 (cache) dépend de #5 (benchmarks) pour valider gains
- Item #2 (context) requis pour #6 (integration tests)

### Risk Assessment
- **LOW** : Tous items sont des améliorations, pas de breaking changes
- Compatibilité backward garantie
- Refactoring progressif possible

### Testing Strategy
- Chaque item doit avoir tests
- Coverage minimum 80% pour nouveau code
- Tests de non-régression systématiques

---

**Maintainer**: GitHub Copilot CLI
**Last Updated**: 2025-12-11
