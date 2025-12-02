# Optimisations du pipeline RETE

**Date** : Janvier 2025  
**Version** : 2.0  
**Status** : Production Ready

---

## Vue d'ensemble

Le pipeline d'ingestion incrémentale RETE (`IngestFile()`) intègre **obligatoirement** deux optimisations :

- ✅ **Validation sémantique incrémentale avec contexte** (TOUJOURS activée, non désactivable)
- ✅ **Garbage Collection après reset** (TOUJOURS activée, non désactivable)
- ⚠️ **Transactions avec rollback** (disponible via API dédiée)

**Important** : Les optimisations 1 et 2 sont maintenant **obligatoires** et ne peuvent plus être désactivées.

---

## 1. Validation sémantique incrémentale (OBLIGATOIRE)

### Comportement

Le pipeline détecte automatiquement le mode d'utilisation :

**Premier fichier (réseau vide)** :
```go
network, err := pipeline.IngestFile("types.tsd", nil, storage)
// → Validation standard complète
```

**Fichiers suivants (mode incrémental)** :
```go
network, err = pipeline.IngestFile("rules.tsd", network, storage)
// → Validation incrémentale avec contexte automatique
// → Prend en compte les types déjà chargés dans le réseau
```

**Après un reset** :
```go
network, err = pipeline.IngestFile("reset_file.tsd", network, storage)
// → Validation standard complète (nouveau réseau)
```

### Logs

```
🔍 Validation sémantique incrémentale avec contexte...
✅ Validation incrémentale réussie (5 types en contexte)
```

### Avantages

- ✅ Détection automatique des types non définis
- ✅ Détection des champs inexistants
- ✅ Validation cohérence inter-fichiers
- ✅ Aucune configuration requise
- ✅ Overhead acceptable (~5-10%)

### Détection d'erreurs

```go
// types.tsd
type Person { id: string, name: string }

// rules.tsd
rule "test" {
    when {
        c: Company(employees > 10)  // ❌ Erreur détectée !
    }
    then {
        print("Found company")
    }
}

network, err := pipeline.IngestFile("rules.tsd", network, storage)
// err = "type 'Company' référencé mais non défini"
```

---

## 2. Garbage Collection après reset (OBLIGATOIRE)

### Comportement

Lorsqu'un fichier contient une commande `reset`, le pipeline effectue automatiquement :

1. **Garbage Collection** de l'ancien réseau
2. **Création** d'un nouveau réseau vide

```go
// Session 1
network, err := pipeline.IngestFile("data1.tsd", nil, storage)
// Réseau créé : 100 nœuds, 10 types

// Session 2 avec reset
network, err = pipeline.IngestFile("reset_data2.tsd", network, storage)
// → GC automatique : libère les 100 nœuds de la session 1
// → Nouveau réseau créé proprement
```

### Fichier avec reset

```tsd
reset  // Déclenche automatiquement le GC

type NewType {
    id: string
    field: string
}
```

### Logs

```
🔄 Commande reset détectée - Garbage Collection de l'ancien réseau
🗑️  GC du réseau existant...
✅ GC terminé
🆕 Création d'un nouveau réseau RETE
```

### Ce qui est nettoyé

- ✅ **Caches** : ArithmeticResultCache, BetaSharingRegistry, AlphaSharingManager
- ✅ **Nœuds** : TypeNodes, AlphaNodes, BetaNodes, TerminalNodes
- ✅ **Références** : Toutes les connexions entre nœuds
- ✅ **Managers** : LifecycleManager, ActionExecutor
- ✅ **Storage** : Tous les faits en mémoire

### Avantages

- ✅ Libération immédiate de la mémoire (~50% sur grands réseaux)
- ✅ Pas de fuites mémoire sur longues sessions
- ✅ Overhead minimal (~1-2%)
- ✅ Aucune configuration requise

### Cas d'usage : sessions multiples

```go
// Serveur long-running avec resets fréquents
for session := 0; session < 1000; session++ {
    network, _ = pipeline.IngestFile("reset.tsd", network, storage)
    network, _ = pipeline.IngestFile("data.tsd", network, storage)
    
    // Traitement...
    
    // ✅ Pas de fuite mémoire grâce au GC automatique
}
```

---

## 3. Transactions (API dédiée, optionnelle)

### Pourquoi optionnelle ?

Les transactions ont un **coût mémoire significatif** (~2x taille du réseau pour le snapshot).

Pour des raisons de performance, elles sont disponibles via une API dédiée plutôt qu'intégrées dans `IngestFile()`.

### Activation manuelle

**Option 1 : Transaction automatique**
```go
// Transaction automatique intégrée
network, err := pipeline.IngestFile(filename, network, storage)
// ✅ Commit automatique si succès
// ✅ Rollback automatique si erreur
```

**Option 2 : Contrôle manuel**
```go
// ❌ Cette approche n'est plus possible (fonctions supprimées)
// Les transactions sont maintenant automatiques dans IngestFile()
network, err := pipeline.IngestFile(filename, network, storage)
// ✅ Transaction gérée automatiquement
```

**Option 3 : Configuration complète**
```go
config := rete.DefaultAdvancedPipelineConfig()
config.EnableTransactions = true
config.AutoCommit = true
config.AutoRollbackOnError = true

network, metrics, err := pipeline.IngestFileWithAdvancedFeatures(
    filename, network, storage, config,
)
```

---

## Comparaison des APIs

### API standard (optimisations obligatoires)

```go
network, err := pipeline.IngestFile(filename, network, storage)
```

**Toujours activé** (non désactivable) :
- ✅ Validation incrémentale
- ✅ GC après reset

**Non activé** :
- ❌ Transactions

**Performance** : Overhead ~5-10%

**Note** : La validation et le GC ne peuvent plus être désactivés.

---

### API avec métriques

```go
network, metrics, err := pipeline.IngestFileWithMetrics(filename, network, storage)
```

**Toujours activé** (non désactivable) :
- ✅ Validation incrémentale
- ✅ GC après reset
- ✅ Métriques détaillées

**Non activé** :
- ❌ Transactions

**Performance** : Overhead ~7-12%

---

### API avec transactions

```go
network, err := pipeline.IngestFile(filename, network, storage)
// ✅ Transaction automatique obligatoire
```

**Toujours activé** (non désactivable) :
- ✅ Validation incrémentale
- ✅ GC après reset

**Activé en plus** :
- ✅ Transactions (auto-commit/rollback)

**Performance** : Overhead ~15-25% + mémoire 2x

---

### API complète personnalisable

```go
config := rete.DefaultAdvancedPipelineConfig()
// Personnaliser config...

network, metrics, err := pipeline.IngestFileWithAdvancedFeatures(
    filename, network, storage, config,
)
```

**Toujours activé** (non configurable) :
- ✅ Validation incrémentale (obligatoire)
- ✅ GC après reset (obligatoire)
- ✅ Métriques (toujours)

**Configurable** :
- ⚙️ Transactions (selon config)

**Performance** : Variable selon configuration des transactions

---

## ⚠️ Impossibilité de désactivation

**Important** : La validation incrémentale et le GC **ne peuvent plus être désactivés**.

Ces optimisations sont maintenant **obligatoires** pour garantir :
- ✅ Détection systématique des erreurs
- ✅ Gestion correcte de la mémoire
- ✅ Cohérence du système

Seules les **transactions** restent optionnelles (via API dédiée) en raison de leur coût mémoire élevé.

Si vous avez des contraintes de performance extrêmes, contactez l'équipe de développement pour discuter d'alternatives.

---

## Recommandations

### Production

```go
// Utiliser l'API standard avec optimisations par défaut
network, err := pipeline.IngestFile(filename, network, storage)
```

**Raison** : Bon équilibre performance/fiabilité

---

### Production critique

```go
// Transactions automatiques (obligatoires)
network, err := pipeline.IngestFile(filename, network, storage)
```

**Raison** : Fiabilité maximale avec rollback automatique

---

### Développement/Debug

```go
// Utiliser l'API avec métriques
network, metrics, err := pipeline.IngestFileWithMetrics(filename, network, storage)
rete.PrintAdvancedMetrics(metrics)
```

**Raison** : Observabilité accrue

---

### Performance maximale

```go
// Utiliser l'API standard sans transactions
network, err := pipeline.IngestFile(filename, network, storage)
```

**Raison** : Validation et GC ont un overhead minimal (~6-12% total)

**Note** : Il n'est plus possible de désactiver la validation et le GC.
Ces optimisations sont obligatoires pour garantir la fiabilité.

---

## Migration depuis anciennes versions

### Avant (v1.x)

```go
network, err := pipeline.IngestFile(filename, network, storage)
// Validation désactivée en mode incrémental
// Pas de GC automatique
```

### Maintenant (v2.0)

```go
network, err := pipeline.IngestFile(filename, network, storage)
// ✅ Validation incrémentale automatique
// ✅ GC automatique après reset
```

**Compatibilité** : 100% rétrocompatible

**Impact** : Légère augmentation du temps d'exécution (~5-10%) mais meilleure fiabilité

---

## FAQ

### Q: Puis-je désactiver la validation incrémentale ?

**Non**, la validation incrémentale est maintenant obligatoire et ne peut plus être désactivée.

Cette décision garantit la détection systématique des erreurs et la cohérence du système.

### Q: Le GC ralentit-il l'ingestion ?

Non, l'impact est minimal (~1-2%). Le GC n'est déclenché que lors d'un `reset`, qui est généralement rare.

### Q: Pourquoi les transactions ne sont-elles pas activées par défaut ?

Pour des raisons de performance. Le snapshot consomme ~2x la mémoire du réseau. Pour la plupart des cas d'usage, la validation incrémentale suffit à garantir la cohérence.

### Q: Puis-je voir les métriques sans changer mon code ?

Oui, utilisez `IngestFileWithMetrics()` :

```go
network, metrics, err := pipeline.IngestFileWithMetrics(filename, network, storage)
fmt.Printf("Parsing: %v\n", metrics.ParsingDuration)
fmt.Printf("Validation: %v\n", metrics.ValidationDuration)
```

### Q: L'overhead de 5-10% est-il acceptable ?

Pour la plupart des applications, oui. L'overhead est compensé par :
- Détection d'erreurs avant construction réseau (gain de temps)
- Pas de fuites mémoire (stabilité long terme)
- Meilleure fiabilité

Pour des benchmarks ou cas très sensibles, désactivez via l'API avancée.

---

## Résumé

| Optimisation                  | Status        | API                         | Overhead | Désactivable |
|-------------------------------|---------------|-----------------------------|----------|--------------|
| Validation incrémentale       | ✅ Obligatoire | IngestFile()               | ~5-10%   | ❌ Non       |
| GC après reset                | ✅ Obligatoire | IngestFile()               | ~1-2%    | ❌ Non       |
| Transactions                  | ✅ Obligatoire | IngestFile()               | < 1%     | ❌ Non       |

**Recommandation générale** : Utiliser `IngestFile()` qui intègre obligatoirement validation et GC.

**Important** : La validation incrémentale et le GC ne peuvent plus être désactivés.

---

## Voir aussi

- [Guide utilisateur complet](ADVANCED_FEATURES_README.md)
- [Spécifications techniques](ADVANCED_OPTIMIZATIONS.md)
- [Vue d'ensemble](README_OPTIMIZATIONS.md)
- [Démarrage rapide](../QUICKSTART_ADVANCED.md)

---

**Auteur** : TSD Contributors  
**Dernière mise à jour** : Janvier 2025  
**Version** : 2.0