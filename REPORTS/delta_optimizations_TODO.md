# TODO et Améliorations Futures - Optimisations Delta

Ce document liste les points d'amélioration identifiés lors de l'optimisation du système de propagation delta.

---

## ✅ Complété (Prompt 09)

- [x] Object pooling (FactDelta, NodeReference, StringBuilder, Map)
- [x] Cache LRU optimisé avec éviction automatique
- [x] Optimisations comparaisons (fast paths pour types simples)
- [x] Batch processing des nœuds par type
- [x] Scripts de profiling et benchmarking
- [x] Tests complets des optimisations
- [x] Rapport de performance détaillé

---

## 🔧 Améliorations Techniques

### 1. Intégration avec le Réseau RETE

**Priorité**: Haute  
**Complexité**: Moyenne  

**Description**:
Les optimisations ont été implémentées au niveau du package `delta` mais ne sont pas encore intégrées avec le réseau RETE complet.

**Actions nécessaires**:
- [ ] Modifier `DeltaPropagator.executeDeltaPropagation()` pour utiliser le pooling
- [ ] Implémenter le callback `classicPropagation()` avec Retract+Insert
- [ ] Intégrer `BatchNodeReferences` dans la propagation vers nœuds
- [ ] Ajouter des tests end-to-end avec réseau RETE complet

**Code affecté**:
```go
// rete/delta/delta_propagator.go
func (dp *DeltaPropagator) executeDeltaPropagation(...) error {
    // TODO: Utiliser BatchNodeReferences pour groupage par type
    // TODO: Release delta après propagation si pas en cache
}

func (dp *DeltaPropagator) classicPropagation(...) error {
    // TODO: Implémenter via callback Retract+Insert
    // Cette implémentation doit être faite par l'appelant
}
```

**Estimation**: 2-3 heures

---

### 2. Gestion du Cycle de Vie des FactDelta

**Priorité**: Haute  
**Complexité**: Faible  

**Description**:
Les `FactDelta` acquis depuis le pool doivent être relâchés par l'appelant. Actuellement, cette responsabilité n'est pas clairement définie.

**Actions nécessaires**:
- [ ] Documenter clairement qui est responsable de `ReleaseFactDelta()`
- [ ] Ajouter des helpers pour gérer automatiquement le cycle de vie
- [ ] Considérer `defer ReleaseFactDelta(delta)` comme pattern recommandé
- [ ] Ajouter des warnings/checks en mode debug pour fuites

**Exemple pattern recommandé**:
```go
func ProcessUpdate(...) error {
    delta, err := detector.DetectDelta(...)
    if err != nil {
        return err
    }
    defer ReleaseFactDelta(delta) // Automatic cleanup
    
    // Use delta...
    return propagator.Propagate(delta)
}
```

**Estimation**: 1 heure

---

### 3. Métriques et Monitoring

**Priorité**: Moyenne  
**Complexité**: Faible  

**Description**:
Bien que des métriques existent (cache stats, detector metrics), elles ne sont pas exposées pour monitoring en production.

**Actions nécessaires**:
- [ ] Exposer métriques via interface standard (Prometheus, etc.)
- [ ] Ajouter métriques sur pool usage (acquisitions, releases, GC)
- [ ] Dashboard pour visualisation temps réel
- [ ] Alerting sur anomalies (hit rate trop bas, trop d'évictions, etc.)

**Métriques à exposer**:
```go
// Pool metrics
pool_acquisitions_total{type="FactDelta"}
pool_releases_total{type="FactDelta"}
pool_capacity_current{type="FactDelta"}

// Cache metrics
cache_hits_total
cache_misses_total
cache_hit_rate
cache_evictions_total
cache_size_current

// Detection metrics
detector_comparisons_total
detector_delta_size_avg
detector_change_ratio_avg
```

**Estimation**: 3-4 heures

---

### 4. Configuration Dynamique

**Priorité**: Basse  
**Complexité**: Faible  

**Description**:
Les tailles de pool et cache sont actuellement hardcodées (constantes). Permettre configuration dynamique.

**Actions nécessaires**:
- [ ] Ajouter `PoolConfig` pour tailles initiales et max
- [ ] Permettre ajustement runtime du cache size
- [ ] Hot-reload de configuration sans redémarrage
- [ ] Validation des valeurs de config (min/max sensibles)

**Exemple**:
```go
type PoolConfig struct {
    FactDeltaInitialCap int
    FactDeltaMaxSize    int
    SliceInitialCap     int
    SliceMaxCap         int
}

func ConfigurePools(config PoolConfig) error {
    // Apply configuration
}
```

**Estimation**: 2 heures

---

### 5. Profiling en Production

**Priorité**: Moyenne  
**Complexité**: Moyenne  

**Description**:
Les scripts de profiling sont adaptés pour développement mais pas production.

**Actions nécessaires**:
- [ ] Profiling continu avec sampling faible overhead
- [ ] Capture automatique lors de pics de CPU/mémoire
- [ ] Stockage et rotation des profiles
- [ ] Analyse automatique des hotspots

**Tools**:
- `pprof` continuous profiling
- Datadog/New Relic APM integration
- Custom profiling triggers

**Estimation**: 4-6 heures

---

## 🚀 Optimisations Avancées (Futur)

### 6. SIMD pour Comparaisons de Masse

**Priorité**: Basse  
**Complexité**: Élevée  

**Description**:
Pour faits très larges (>100 champs), utiliser SIMD pour comparaisons parallèles.

**Pré-requis**:
- Analyse profiling montrant que comparaisons sont hotspot
- Faits homogènes (mêmes types de champs)

**Estimation**: 1-2 semaines (R&D + implémentation)

---

### 7. Compression Delta

**Priorité**: Basse  
**Complexité**: Élevée  

**Description**:
Compresser les deltas pour réduire footprint mémoire en cache.

**Trade-offs**:
- ➕ Réduction mémoire (potentiellement 50-70%)
- ➖ Overhead CPU (compression/décompression)
- ➖ Complexité accrue

**Estimation**: 1-2 semaines

---

### 8. Batch Processing au Niveau Réseau

**Priorité**: Moyenne  
**Complexité**: Élevée  

**Description**:
Plutôt que propager delta-par-delta, accumuler et propager par batches.

**Avantages**:
- Meilleure utilisation cache
- Moins de context switches
- Optimisations vectorielles possibles

**Estimation**: 2-3 semaines

---

## 📋 Checklist Intégration

Avant de considérer les optimisations "production-ready":

- [ ] Intégration avec réseau RETE complet
- [ ] Tests end-to-end avec données réelles
- [ ] Profiling sous charge réaliste
- [ ] Documentation API complète
- [ ] Métriques exposées pour monitoring
- [ ] Configuration ajustable
- [ ] Gestion erreurs robuste
- [ ] Logs appropriés (niveaux, contexte)
- [ ] Benchmarks regression dans CI/CD
- [ ] Performance baseline documentée

---

## 🎯 Priorités Recommandées

1. **Court terme** (Sprint actuel):
   - Intégration RETE (#1)
   - Cycle de vie FactDelta (#2)

2. **Moyen terme** (1-2 sprints):
   - Métriques monitoring (#3)
   - Profiling production (#5)

3. **Long terme** (Backlog):
   - Configuration dynamique (#4)
   - Optimisations SIMD (#6)
   - Compression delta (#7)
   - Batch processing réseau (#8)

---

**Dernière mise à jour**: 2026-01-02  
**Mainteneur**: TSD Team  
**Contact**: optimization@tsd-project.dev
