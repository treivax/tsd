# 📁 Structure Organisée du Module RETE

## 🎯 Organisation Logique

Le module RETE est maintenant organisé avec des **préfixes de fichiers logiques** pour une meilleure navigation et maintenance, tout en conservant un **package unique** pour éviter les problèmes de dépendances circulaires.

## 📂 Structure des Fichiers

### **🏗️ Core RETE (Architecture principale)**
- `network.go` - Réseau RETE principal et gestion des nœuds
- `rete.go` - Types et structures fondamentales  
- `evaluator.go` - Évaluateur d'expressions et conditions
- `alpha_builder.go` - Construction des nœuds Alpha
- `converter.go` - Conversion AST vers RETE

### **📊 Monitoring (Préfixe `monitor_`)**
- `monitor_server.go` - Serveur HTTP avec API REST et WebSocket
- `monitor_integrator.go` - Collecteur de métriques en temps réel
- `monitor_network.go` - Réseau RETE avec monitoring intégré

### **⚡ Performance (Préfixe `perf_`)**
- `perf_hash_joins.go` - Moteur de jointures hash optimisé
- `perf_eval_cache.go` - Cache LRU intelligent pour évaluations
- `perf_token_propagation.go` - Propagation parallèle de tokens
- `perf_profiler.go` - Profileur de performance et métriques

### **💾 Storage (Préfixe `store_`)**
- `store_base.go` - Interface de stockage de base
- `store_indexed.go` - Stockage indexé multi-niveaux

### **🧪 Tests (Préfixe `test_`)**
- `test_integration.go` - Tests d'intégration avancés
- `test_performance.go` - Benchmarks de performance
- `test_perf_integration.go` - Tests d'intégration performance

### **📁 Assets Organisés**
- `assets/web/` - Interface web de monitoring (HTML, CSS, JS)
- `cmd/` - Applications et démos
- `docs/` - Documentation technique
- `pkg/` - Packages utilitaires
- `scripts/` - Scripts d'automatisation

## 🔗 Avantages de cette Organisation

### **✅ Navigabilité Améliorée**
- **Regroupement logique** : Fichiers similaires groupés par préfixe
- **Recherche facilitée** : `monitor_*`, `perf_*`, `store_*`, `test_*`
- **Responsabilités claires** : Chaque préfixe a un rôle défini

### **✅ Maintenabilité**
- **Package unique** : Pas de dépendances circulaires
- **API cohérente** : Tous les types accessibles directement
- **Structure claire** : Organisation évidente par nom de fichier

### **✅ Développement Efficace**
- **Édition ciblée** : Facilite la modification de composants spécifiques
- **Tests organisés** : Tests regroupés et identifiables
- **Assets séparés** : Interface web dans un dossier dédié

## 📋 Convention de Nommage

| Préfixe | Responsabilité | Exemples |
|---------|---------------|-----------|
| `monitor_` | Surveillance et observabilité | `monitor_server.go` |
| `perf_` | Optimisations de performance | `perf_hash_joins.go` |
| `store_` | Persistance et stockage | `store_indexed.go` |
| `test_` | Tests et benchmarks | `test_integration.go` |
| *(sans préfixe)* | Core RETE | `network.go`, `rete.go` |

## 🚀 Utilisation

### **Import Simple**
```go
import "github.com/treivax/tsd/rete"

// Tous les types disponibles directement
network := rete.NewReteNetwork(storage)
monitor := rete.NewMonitoringServer(config, network)
cache := rete.NewEvaluationCache(config)
```

### **Navigation Fichiers**
```bash
# Fichiers de monitoring
ls monitor_*.go

# Fichiers de performance  
ls perf_*.go

# Fichiers de stockage
ls store_*.go

# Tests d'intégration
ls test_*.go
```

## 📊 Comparaison Avant/Après

### **Avant (Non organisé)**
```
rete/
├── monitoring_server.go       # ❌ Mélange confus
├── hash_join_engine.go        # ❌ Responsabilités mélangées  
├── evaluation_cache.go        # ❌ Difficile à naviguer
├── network.go                 # ❌ Pas de logique claire
└── [30+ fichiers mélangés]    # ❌ Navigation difficile
```

### **Après (Organisé)**
```
rete/
├── network.go rete.go evaluator.go         # ✅ Core RETE
├── monitor_*.go                             # ✅ Monitoring
├── perf_*.go                               # ✅ Performance  
├── store_*.go                              # ✅ Stockage
├── test_*.go                               # ✅ Tests
└── assets/web/                             # ✅ Interface web
```

---

## 🎉 Résultat

Le module RETE est maintenant **parfaitement organisé** avec une structure claire, une navigation facilitée et une maintenabilité améliorée, tout en conservant une API simple et cohérente ! ✨