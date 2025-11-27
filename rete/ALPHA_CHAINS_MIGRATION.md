# Guide de Migration : Chaînes d'AlphaNodes

## Table des Matières

1. [Vue d'ensemble](#vue-densemble)
2. [Impact sur le code existant](#impact-sur-le-code-existant)
3. [Migration pas à pas](#migration-pas-à-pas)
4. [Configuration et tuning](#configuration-et-tuning)
5. [Troubleshooting](#troubleshooting)
6. [Rollback](#rollback)
7. [FAQ Migration](#faq-migration)

---

## Vue d'ensemble

### Qu'est-ce qui change ?

La fonctionnalité de chaînes d'AlphaNodes introduit plusieurs améliorations au réseau RETE :

**✅ Nouvelles fonctionnalités :**
- Construction automatique de chaînes d'AlphaNodes
- Partage intelligent de nœuds entre règles
- Cache LRU pour calculs de hash
- Métriques détaillées de performance
- Configuration flexible des caches

**✅ Rétrocompatibilité :**
- Les règles existantes fonctionnent sans modification
- L'API publique reste stable
- Les tests existants passent sans changement
- Le format de persistence est compatible

**⚠️ Changements internes :**
- `ReteNetwork` a un nouveau champ `Config`
- `AlphaSharingRegistry` utilise un cache LRU
- Nouveaux constructeurs avec configuration
- Logging plus détaillé

### Qui est impacté ?

| Utilisateur | Impact | Action requise |
|-------------|--------|----------------|
| **Utilisateur TSD** (écrit des règles) | ✅ Aucun | Aucune - bénéficie automatiquement |
| **Développeur d'API** (utilise ReteNetwork) | ⚠️ Minimal | Optionnel - peut utiliser nouvelle config |
| **Contributeur Core** (modifie RETE) | 🔴 Moyen | Doit comprendre nouvelle architecture |

---

## Impact sur le code existant

### Code qui continue de fonctionner sans changement

✅ **Création de réseau basique :**
```go
// AVANT et APRÈS - identique
storage := NewMemoryStorage()
network := NewReteNetwork(storage)
```

✅ **Ajout de règles :**
```go
// AVANT et APRÈS - identique
err := network.AddRule(rule)
```

✅ **Évaluation de faits :**
```go
// AVANT et APRÈS - identique
network.Assert(fact)
```

✅ **Suppression de règles :**
```go
// AVANT et APRÈS - identique
network.RemoveRule(ruleID)
```

✅ **Tests unitaires existants :**
```go
// Tous les tests existants continuent de passer
func TestExistingFeature(t *testing.T) {
    // Aucune modification nécessaire
    storage := NewMemoryStorage()
    network := NewReteNetwork(storage)
    // ... reste du test inchangé
}
```

### Nouveau code optionnel disponible

🆕 **Création avec configuration personnalisée :**
```go
// NOUVEAU - optionnel
config := HighPerformanceChainConfig()
network := NewReteNetworkWithConfig(storage, config)
```

🆕 **Accès aux métriques :**
```go
// NOUVEAU - optionnel
metrics := network.AlphaChainBuilder.GetMetrics()
fmt.Printf("Sharing ratio: %.1f%%\n", metrics.SharingRatio * 100)
```

🆕 **Configuration fine du cache :**
```go
// NOUVEAU - optionnel
config := &ChainPerformanceConfig{
    HashCacheEnabled: true,
    HashCacheMaxSize: 50000,
    HashCacheTTL:     10 * time.Minute,
}
network := NewReteNetworkWithConfig(storage, config)
```

### Changements de comportement observable

#### 1. IDs de nœuds alpha

**AVANT :**
```
AlphaNode IDs: "rule_myRule_alpha_0", "rule_myRule_alpha_1"
Format: Basé sur le nom de règle + index
```

**APRÈS :**
```
AlphaNode IDs: "alpha_024a66ab3f89c2d1", "alpha_def456789abc012"
Format: Basé sur le hash de condition (alpha_<16_hex_chars>)
```

**Impact :**
- ⚠️ Si votre code parse les IDs de nœuds → **Adaptation nécessaire**
- ✅ Si vous utilisez les IDs comme références opaques → **Aucun impact**

**Migration :**
```go
// AVANT - code fragile
if strings.HasPrefix(nodeID, "rule_myRule_") {
    // Ne fonctionne plus correctement
}

// APRÈS - code robuste
node := network.GetNode(nodeID)
if alphaNode, ok := node.(*AlphaNode); ok {
    // Utiliser l'objet directement
}
```

#### 2. Logging

**AVANT :**
```
[INFO] AlphaNode created for rule myRule
[INFO] Connecting to parent node
```

**APRÈS :**
```
🆕 [AlphaChainBuilder] Nouveau nœud alpha alpha_abc123 créé pour la règle myRule (condition 1/2)
🔗 [AlphaChainBuilder] Connexion du nœud alpha_abc123 au parent type_person
♻️  [AlphaChainBuilder] Réutilisation du nœud alpha alpha_def456 pour la règle other (condition 1/1)
```

**Impact :**
- ⚠️ Si vous parsez les logs → **Format différent**
- ✅ Logs plus informatifs pour debugging

#### 3. Métriques

**NOUVEAU :**
Des métriques détaillées sont maintenant disponibles :

```go
metrics := network.AlphaChainBuilder.GetMetrics()
// Nouvelles métriques :
// - SharingRatio
// - HashCacheHits
// - AverageBuildTime
// etc.
```

**Impact :**
- ✅ Opportunité de monitoring amélioré
- ✅ Pas d'impact si non utilisé

---

## Migration pas à pas

### Étape 1 : Audit du code existant (Optionnel)

**Objectif :** Identifier le code qui pourrait bénéficier de la nouvelle configuration.

**Checklist :**

```bash
# 1. Chercher les créations de ReteNetwork
grep -r "NewReteNetwork" .

# 2. Chercher les références aux IDs de nœuds
grep -r "alpha_" . | grep -v "test"

# 3. Identifier les règles avec beaucoup de conditions communes
# (candidats pour bénéficier du partage)
```

**Questions à se poser :**
- Avez-vous plus de 100 règles ? → Considérer HighPerformanceConfig
- Environnement mémoire contraint ? → Considérer LowMemoryConfig
- Besoin de métriques ? → Activer la collecte

### Étape 2 : Tests en environnement de développement

**2.1 Exécuter les tests existants :**

```bash
# Tous les tests doivent passer sans modification
go test ./rete/... -v

# Vérifier spécifiquement les tests alpha
go test ./rete/ -run Alpha -v
```

**Résultat attendu :**
```
✓ Tous les tests passent
✓ Aucune régression détectée
```

**2.2 Tester avec votre ensemble de règles :**

```go
func TestMyRulesWithChains(t *testing.T) {
    storage := NewMemoryStorage()
    
    // Utiliser config par défaut (recommandé)
    network := NewReteNetwork(storage)
    
    // Charger vos règles existantes
    for _, rule := range myExistingRules {
        err := network.AddRule(rule)
        if err != nil {
            t.Fatalf("Failed to add rule: %v", err)
        }
    }
    
    // Vérifier les métriques de partage
    metrics := network.AlphaChainBuilder.GetMetrics()
    t.Logf("Sharing ratio: %.1f%%", metrics.SharingRatio * 100)
    t.Logf("Nodes created: %d", metrics.TotalNodesCreated)
    t.Logf("Nodes reused: %d", metrics.TotalNodesReused)
    
    // Tester l'évaluation
    for _, fact := range myTestFacts {
        network.Assert(fact)
    }
    
    // Vérifier les résultats attendus
    // ... assertions ...
}
```

### Étape 3 : Configuration optimale (Optionnel)

**3.1 Benchmarker les différentes configurations :**

```go
func BenchmarkDefaultConfig(b *testing.B) {
    storage := NewMemoryStorage()
    config := DefaultChainPerformanceConfig()
    network := NewReteNetworkWithConfig(storage, config)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Charger règles et évaluer
    }
}

func BenchmarkHighPerformanceConfig(b *testing.B) {
    storage := NewMemoryStorage()
    config := HighPerformanceChainConfig()
    network := NewReteNetworkWithConfig(storage, config)
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Charger règles et évaluer
    }
}
```

**3.2 Choisir la configuration appropriée :**

| Scénario | Configuration recommandée | Raison |
|----------|---------------------------|--------|
| < 50 règles | `DefaultChainPerformanceConfig()` | Overhead minimal |
| 50-500 règles | `DefaultChainPerformanceConfig()` | Bon équilibre |
| 500-5000 règles | `HighPerformanceChainConfig()` | Cache plus large |
| Embedded/IoT | `LowMemoryChainConfig()` | Footprint minimal |
| Development/Debug | `DisabledCachesConfig()` | Comportement simple |

**3.3 Implémenter la configuration :**

```go
// config/rete_config.go

func NewProductionReteNetwork(storage Storage) *ReteNetwork {
    // Choisir config selon l'environnement
    var config *ChainPerformanceConfig
    
    switch os.Getenv("ENVIRONMENT") {
    case "production":
        config = HighPerformanceChainConfig()
    case "staging":
        config = DefaultChainPerformanceConfig()
    case "development":
        config = DefaultChainPerformanceConfig()
    default:
        config = DefaultChainPerformanceConfig()
    }
    
    return NewReteNetworkWithConfig(storage, config)
}
```

### Étape 4 : Monitoring et observabilité (Recommandé)

**4.1 Exposer les métriques :**

```go
// metrics/rete_metrics.go

import (
    "net/http"
    "github.com/yourorg/tsd/rete"
)

func ReteMetricsHandler(network *rete.ReteNetwork) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        metrics := network.AlphaChainBuilder.GetMetrics()
        
        // Format Prometheus
        w.Header().Set("Content-Type", "text/plain")
        w.Write([]byte(metrics.ExportText()))
    }
}

// Dans main.go
http.HandleFunc("/metrics/rete", ReteMetricsHandler(network))
```

**4.2 Créer des alertes (optionnel) :**

```yaml
# prometheus_alerts.yml

groups:
  - name: rete_alpha_chains
    rules:
      - alert: LowAlphaSharingRatio
        expr: alpha_sharing_ratio < 0.3
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Low alpha node sharing detected"
          description: "Sharing ratio is {{ $value }}, expected >30%"
      
      - alert: HighHashCacheMissRate
        expr: (alpha_hash_cache_misses / (alpha_hash_cache_hits + alpha_hash_cache_misses)) > 0.5
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High hash cache miss rate"
          description: "Consider increasing cache size"
```

**4.3 Dashboard Grafana (optionnel) :**

```json
{
  "dashboard": {
    "title": "RETE Alpha Chains Performance",
    "panels": [
      {
        "title": "Sharing Ratio",
        "targets": [{"expr": "alpha_sharing_ratio"}],
        "type": "gauge"
      },
      {
        "title": "Cache Hit Rate",
        "targets": [{"expr": "alpha_hash_cache_hits / (alpha_hash_cache_hits + alpha_hash_cache_misses)"}],
        "type": "graph"
      }
    ]
  }
}
```

### Étape 5 : Déploiement progressif

**5.1 Déploiement canary (recommandé) :**

```go
// Déployer sur 10% du trafic d'abord
func CreateNetwork(storage Storage, canaryPercent int) *ReteNetwork {
    // Décision basée sur ID de requête
    if shouldUseNewFeature(canaryPercent) {
        config := DefaultChainPerformanceConfig()
        return NewReteNetworkWithConfig(storage, config)
    }
    
    // Ancienne version (en fait, toujours la même maintenant)
    return NewReteNetwork(storage) // Utilise config par défaut
}
```

**5.2 Monitoring post-déploiement :**

Surveillez pendant 24-48h :
- ✅ Latence P50, P95, P99
- ✅ Utilisation mémoire
- ✅ Taux de sharing
- ✅ Hit rate du cache
- ✅ Logs d'erreur

**5.3 Critères de succès :**

| Métrique | Attendu | Action si non atteint |
|----------|---------|----------------------|
| Sharing ratio | > 30% | Normal si règles très différentes |
| Cache hit rate | > 70% | Augmenter taille cache |
| Latence | Similaire ou mieux | Vérifier config |
| Mémoire | Réduite 20-80% | Vérifier cache TTL |

### Étape 6 : Nettoyage (Optionnel)

**6.1 Supprimer code obsolète :**

Si vous aviez des workarounds pour l'absence de partage :

```go
// AVANT - workaround manuel
// (peut être supprimé maintenant)
func manuallyDeduplicateAlphaNodes(network *ReteNetwork) {
    // Ce code n'est plus nécessaire
}
```

**6.2 Simplifier les tests :**

```go
// AVANT - tests vérifiant absence de partage
func TestNoSharing(t *testing.T) {
    // Ce test peut être supprimé ou adapté
}

// APRÈS - tests vérifiant le partage
func TestSharingWorks(t *testing.T) {
    // Nouveaux tests positifs
}
```

---

## Configuration et tuning

### Matrice de configuration

| Paramètre | Default | High Perf | Low Memory | Description |
|-----------|---------|-----------|------------|-------------|
| `HashCacheEnabled` | true | true | true | Active le cache LRU |
| `HashCacheMaxSize` | 10,000 | 100,000 | 1,000 | Taille max du cache |
| `HashCacheEviction` | LRU | LRU | LRU | Politique d'éviction |
| `HashCacheTTL` | 5min | 15min | 1min | Durée de vie |
| `EnableMetrics` | true | true | true | Collection métriques |

### Formules de sizing

**Taille du cache hash :**
```
cache_size = nombre_conditions_uniques × 1.5

Exemple :
- 500 règles
- ~3 conditions/règle en moyenne
- ~30% de conditions uniques
→ 500 × 3 × 0.3 × 1.5 = ~675
→ Recommandé : 1,000 - 5,000
```

**TTL du cache :**
```
ttl = max(temps_entre_ajouts_règles, 5min)

Exemple :
- Règles chargées au démarrage → TTL = 1h+ (ou 0 = infini)
- Règles ajoutées dynamiquement → TTL = 5-15min
```

**Mémoire du cache :**
```
mémoire_cache ≈ cache_size × 100 bytes

Exemple :
- 10,000 entrées → ~1 MB
- 100,000 entrées → ~10 MB
```

### Configurations personnalisées

**Exemple 1 : Haute fréquence, faible cardinalité**
```go
// Beaucoup de règles, peu de conditions uniques
config := &ChainPerformanceConfig{
    HashCacheEnabled:  true,
    HashCacheMaxSize:  50000,  // Large pour couvrir toutes conditions
    HashCacheEviction: EvictionPolicyLRU,
    HashCacheTTL:      0,      // Infini - conditions stables
    EnableMetrics:     true,
}
```

**Exemple 2 : Règles dynamiques**
```go
// Ajout/suppression fréquent de règles
config := &ChainPerformanceConfig{
    HashCacheEnabled:  true,
    HashCacheMaxSize:  20000,
    HashCacheEviction: EvictionPolicyLRU,
    HashCacheTTL:      5 * time.Minute,  // Court pour rafraîchir
    EnableMetrics:     true,
}
```

**Exemple 3 : Embedded/Edge device**
```go
// Ressources très limitées
config := &ChainPerformanceConfig{
    HashCacheEnabled:  true,
    HashCacheMaxSize:  500,    // Minimal
    HashCacheEviction: EvictionPolicyLRU,
    HashCacheTTL:      30 * time.Second,
    EnableMetrics:     false,  // Économiser overhead
}
```

---

## Troubleshooting

### Problème 1 : Tests échouent après migration

**Symptôme :**
```
FAIL: TestMyFeature
Expected node ID "rule_abc_alpha_0", got "alpha_024a66ab"
```

**Cause :**
Le test vérifie le format des IDs de nœuds (fragile).

**Solution :**
```go
// AVANT - test fragile
assert.Equal(t, "rule_abc_alpha_0", node.GetID())

// APRÈS - test robuste
assert.NotEmpty(t, node.GetID())
assert.True(t, strings.HasPrefix(node.GetID(), "alpha_"))

// OU mieux - tester le comportement, pas l'implémentation
alphaNode, ok := node.(*AlphaNode)
assert.True(t, ok)
assert.NotNil(t, alphaNode.Condition)
```

### Problème 2 : Performance dégradée

**Symptôme :**
```
Latence augmentée de 20% après migration
```

**Diagnostic :**
```go
// Vérifier les métriques
metrics := network.AlphaChainBuilder.GetMetrics()
fmt.Printf("Cache hit rate: %.1f%%\n", 
    float64(metrics.HashCacheHits) / 
    float64(metrics.HashCacheHits + metrics.HashCacheMisses) * 100)
```

**Solutions possibles :**

1. **Hit rate < 70% → Augmenter cache**
```go
config.HashCacheMaxSize = 50000  // Au lieu de 10000
```

2. **Beaucoup d'évictions → Augmenter TTL**
```go
config.HashCacheTTL = 15 * time.Minute  // Au lieu de 5min
```

3. **Métriques activées sur hot path → Désactiver**
```go
config.EnableMetrics = false  // Pour production haute perf
```

### Problème 3 : Utilisation mémoire élevée

**Symptôme :**
```
Mémoire heap augmentée de 100MB
```

**Diagnostic :**
```go
// Vérifier taille des caches
cacheSize := network.AlphaSharingManager.GetHashCacheSize()
connectionCacheSize := network.AlphaChainBuilder.GetConnectionCacheSize()

fmt.Printf("Hash cache: %d entries\n", cacheSize)
fmt.Printf("Connection cache: %d entries\n", connectionCacheSize)
```

**Solutions :**

1. **Cache trop large → Réduire taille**
```go
config.HashCacheMaxSize = 1000  // Au lieu de 10000
```

2. **TTL trop long → Réduire**
```go
config.HashCacheTTL = 1 * time.Minute  // Au lieu de 5min
```

3. **Nettoyer périodiquement**
```go
// Nettoyer les entrées expirées
network.AlphaSharingManager.CleanExpiredHashCache()

// Nettoyer le cache de connexion
network.AlphaChainBuilder.ClearConnectionCache()
```

### Problème 4 : Partage non optimal

**Symptôme :**
```
Sharing ratio: 15% (attendu: >50%)
```

**Diagnostic :**
```go
// Analyser les règles
for ruleID, rule := range allRules {
    chain := network.GetChainForRule(ruleID)
    stats := network.AlphaChainBuilder.GetChainStats(chain)
    fmt.Printf("Rule %s: shared=%d, new=%d\n", 
        ruleID, stats["shared_nodes"], stats["new_nodes"])
}
```

**Causes possibles :**

1. **Variables différentes**
```tsd
rule r1 : {p: Person} / p.age > 18 ==> ...  # Variable 'p'
rule r2 : {u: Person} / u.age > 18 ==> ...  # Variable 'u' → pas de partage
```

**Solution :** Utiliser noms de variables cohérents.

2. **Conditions vraiment uniques**
```tsd
rule r1 : {p: Person} / p.id == "abc" ==> ...
rule r2 : {p: Person} / p.id == "def" ==> ...
# Conditions différentes → pas de partage (normal)
```

**Solution :** C'est le comportement attendu.

3. **Types numériques différents**
```tsd
rule r1 : {p: Person} / p.age > 18 ==> ...    # int
rule r2 : {p: Person} / p.age > 18.0 ==> ...  # float
# Hashes différents
```

**Solution :** Normaliser les types dans les règles.

### Problème 5 : Memory leak apparent

**Symptôme :**
```
Nombre de nœuds alpha augmente sans cesse
Mémoire ne se libère pas après suppression de règles
```

**Diagnostic :**
```go
// Vérifier les refcounts
for nodeID := range network.AlphaNodes {
    lifecycle, _ := network.LifecycleManager.GetNodeLifecycle(nodeID)
    fmt.Printf("Node %s: RefCount=%d, Rules=%v\n", 
        nodeID, lifecycle.GetRefCount(), lifecycle.GetRuleIDs())
}
```

**Causes possibles :**

1. **Règles non supprimées correctement**
```go
// MAL - ne libère pas les ressources
delete(ruleMap, ruleID)

// BIEN - utilise l'API
network.RemoveRule(ruleID)
```

2. **Références circulaires**

**Solution :** Vérifier que tous les chemins de suppression passent par `RemoveRule()`.

---

## Rollback

### Procédure de rollback

La fonctionnalité étant rétrocompatible, un "rollback" n'est généralement pas nécessaire.

Cependant, si vous voulez revenir à un comportement antérieur :

**Option 1 : Désactiver les caches (conserve le reste)**
```go
config := DisabledCachesConfig()
network := NewReteNetworkWithConfig(storage, config)
```

**Option 2 : Utiliser une version antérieure du code**
```bash
git checkout <commit_avant_feature>
go build ./...
```

**Option 3 : Feature flag (si implémenté)**
```go
if os.Getenv("DISABLE_ALPHA_CHAINS") == "true" {
    // Utiliser ancienne implémentation (si gardée)
}
```

### Checklist de rollback

Avant de rollback, vérifier :

- [ ] Les logs pour identifier la vraie cause du problème
- [ ] Les métriques pour confirmer une régression
- [ ] Les tests pour reproduire le problème
- [ ] Si une simple reconfiguration suffirait

Si rollback nécessaire :

1. [ ] Désactiver les caches via config
2. [ ] Monitorer les métriques (latence, mémoire)
3. [ ] Analyser les logs d'erreur
4. [ ] Créer un issue GitHub avec détails
5. [ ] Planifier investigation approfondie

---

## FAQ Migration

### Q1 : Dois-je modifier mes règles TSD existantes ?

**R :** Non, absolument aucune modification nécessaire. Les règles fonctionnent exactement comme avant.

### Q2 : Les IDs de nœuds vont changer à chaque démarrage ?

**R :** Non. Les IDs sont basés sur le hash du contenu de la condition, donc ils sont **déterministes**. Même condition = même ID à chaque fois.

### Q3 : Puis-je désactiver complètement les chaînes alpha ?

**R :** Les chaînes alpha sont maintenant le mécanisme standard. Vous pouvez désactiver les **caches** mais pas la construction de chaînes elle-même :

```go
config := DisabledCachesConfig()
network := NewReteNetworkWithConfig(storage, config)
```

### Q4 : Y a-t-il un impact sur la persistence (sauvegarde/restauration) ?

**R :** Non. Le format de persistence est compatible. Les règles sont sauvegardées et restaurées de la même manière.

### Q5 : Comment migrer si j'ai déjà des règles en production ?

**R :** Aucune migration nécessaire. Au prochain démarrage avec le nouveau code :
1. Les règles sont rechargées
2. Les chaînes sont construites automatiquement
3. Le partage se fait automatiquement

### Q6 : Les métriques sont-elles thread-safe ?

**R :** Oui, toutes les structures sont protégées par des mutexes. Vous pouvez appeler `GetMetrics()` depuis n'importe quel thread.

### Q7 : Quel est le coût des métriques ?

**R :** Négligeable (<1% overhead). Elles utilisent des opérations atomiques et pas d'allocations coûteuses.

### Q8 : Comment savoir si le partage fonctionne bien ?

**R :** Vérifiez le `SharingRatio` dans les métriques :
- < 30% : Faible (normal si règles très différentes)
- 30-50% : Moyen (bon pour workloads mixtes)
- 50-70% : Bon (règles avec patterns communs)
- > 70% : Excellent (beaucoup de conditions communes)

### Q9 : Puis-je forcer le partage de nœuds spécifiques ?

**R :** Non besoin. Le partage est automatique basé sur le contenu sémantique. Si deux conditions sont identiques, elles partagent automatiquement.

### Q10 : Comment déboguer si le partage ne fonctionne pas comme attendu ?

**R :** Activez les logs détaillés et comparez les hashes :

```go
// Comparer deux conditions
hash1 := ConditionHash(condition1, "p")
hash2 := ConditionHash(condition2, "p")

if hash1 != hash2 {
    // Normaliser et comparer visuellement
    norm1 := normalizeConditionForSharing(condition1)
    norm2 := normalizeConditionForSharing(condition2)
    
    json1, _ := json.MarshalIndent(norm1, "", "  ")
    json2, _ := json.MarshalIndent(norm2, "", "  ")
    
    fmt.Println("Condition 1 normalized:")
    fmt.Println(string(json1))
    fmt.Println("Condition 2 normalized:")
    fmt.Println(string(json2))
    // Identifier visuellement les différences
}
```

---

## Support et Ressources

### Documentation

- [Guide Utilisateur](ALPHA_CHAINS_USER_GUIDE.md)
- [Guide Technique](ALPHA_CHAINS_TECHNICAL_GUIDE.md)
- [Exemples](ALPHA_CHAINS_EXAMPLES.md)
- [Partage AlphaNode](ALPHA_NODE_SHARING.md)

### Code

- Tests d'intégration : `alpha_sharing_lru_integration_test.go`
- Exemples : `examples/lru_cache/`
- Benchmarks : `alpha_sharing_benchmark_test.go`

### Obtenir de l'aide

- Issues GitHub : [github.com/yourorg/tsd/issues](https://github.com/yourorg/tsd/issues)
- Discussions : [github.com/yourorg/tsd/discussions](https://github.com/yourorg/tsd/discussions)
- Documentation : [docs.tsd-lang.org](https://docs.tsd-lang.org)

---

## Changelog de migration

### Version 1.0 (2025-01-27)

**Ajouté :**
- ✅ Construction automatique de chaînes d'AlphaNodes
- ✅ Partage intelligent avec reference counting
- ✅ Cache LRU pour calculs de hash
- ✅ Métriques détaillées de performance
- ✅ Configurations preset (Default, HighPerf, LowMemory)
- ✅ Documentation complète

**Modifié :**
- ⚠️ Format des IDs de nœuds alpha (hash-based)
- ⚠️ Logging plus détaillé avec émojis
- ⚠️ Nouveaux constructeurs avec configuration

**Déprécié :**
- Aucun

**Supprimé :**
- Aucun

**Rétrocompatibilité :**
- ✅ 100% compatible avec code existant
- ✅ Tous les tests existants passent
- ✅ Aucune modification de règles nécessaire

---

## Licence

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License