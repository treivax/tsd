# 01 - Structures de Données Core du Module Xuples

**Date** : 2025-12-17  
**Status** : ✅ CONCEPTION COMPLÈTE  

---

## 🎯 Objectif

Définir les structures de données fondamentales du module xuples sans hardcoding, permettant une architecture flexible et extensible.

---

## 📋 Structures Principales

### 1. Xuple - Activation Disponible

```go
// Xuple représente une activation de règle disponible dans un xuple-space.
//
// Un xuple encapsule :
//   - L'action à exécuter (issue de la règle RETE)
//   - Le token déclencheur avec tous ses bindings
//   - Les faits ayant déclenché la règle
//   - Les métadonnées de tracking (création, consommation, etc.)
//   - Le statut dans son cycle de vie
//
// Thread-Safety:
//   - Xuple est immutable après création (sauf Status et ConsumedBy)
//   - Les modifications de statut doivent être synchronisées par XupleSpace
type Xuple struct {
    // ID unique du xuple (format: "xuple_<counter>")
    ID string
    
    // Référence à l'action RETE déclenchée
    Action *rete.Action
    
    // Token RETE contenant tous les bindings
    Token *rete.Token
    
    // Faits ayant déclenché cette activation (subset de Token.Facts)
    TriggeringFacts []*rete.Fact
    
    // Statut dans le cycle de vie
    Status XupleStatus
    
    // Timestamp de création
    CreatedAt time.Time
    
    // Timestamp de dernière modification de statut
    UpdatedAt time.Time
    
    // Timestamp d'expiration (nil si pas d'expiration)
    ExpiresAt *time.Time
    
    // IDs des agents ayant consommé ce xuple
    // (peut être vide ou contenir plusieurs IDs selon la politique)
    ConsumedBy []string
    
    // Nombre de fois que ce xuple a été consommé
    ConsumptionCount int
    
    // Métadonnées additionnelles (extensibilité)
    Metadata map[string]interface{}
}

// XupleStatus représente l'état d'un xuple dans son cycle de vie
type XupleStatus string

const (
    // StatusPending : xuple créé et disponible pour consommation
    StatusPending XupleStatus = "pending"
    
    // StatusConsumed : xuple consommé par au moins un agent
    StatusConsumed XupleStatus = "consumed"
    
    // StatusExpired : xuple expiré (dépassé sa durée de vie)
    StatusExpired XupleStatus = "expired"
    
    // StatusArchived : xuple archivé (conservé pour historique mais inactif)
    StatusArchived XupleStatus = "archived"
)

// IsTerminal retourne true si le xuple est dans un état terminal
// (ne peut plus changer de statut)
func (s XupleStatus) IsTerminal() bool {
    return s == StatusExpired || s == StatusArchived
}

// IsAvailable retourne true si le xuple est disponible pour consommation
func (s XupleStatus) IsAvailable() bool {
    return s == StatusPending || s == StatusConsumed
}
```

**Justification des champs** :
- `ID` : Identifiant unique thread-safe (compteur atomique)
- `Action` : Référence immuable vers action RETE
- `Token` : Contient tous les bindings nécessaires (immuable via BindingChain)
- `TriggeringFacts` : Permet de filtrer/indexer par faits
- `Status` : Machine à états claire
- `CreatedAt/UpdatedAt/ExpiresAt` : Traçabilité et rétention
- `ConsumedBy` : Supporte politiques multi-consommation
- `ConsumptionCount` : Limite de consommation
- `Metadata` : Extensibilité sans changer structure

### 2. XupleSpace - Espace de Stockage

```go
// XupleSpace gère un espace nommé de xuples avec politiques configurables.
//
// Responsabilités :
//   - Stockage thread-safe des xuples
//   - Application des politiques de sélection/consommation/rétention
//   - Indexation pour recherche efficace
//   - Gestion du cycle de vie (expiration, archivage)
//
// Thread-Safety:
//   - Toutes les opérations sont thread-safe via sync.RWMutex
//   - Les politiques sont appliquées de manière atomique
type XupleSpace struct {
    // Nom unique du xuple-space
    name string
    
    // Stockage principal des xuples par ID
    xuples map[string]*Xuple
    
    // Index par action (nom) pour recherche rapide
    // map[actionName]map[xupleID]*Xuple
    xuplesByAction map[string]map[string]*Xuple
    
    // Index par statut pour recherche rapide
    // map[status]map[xupleID]*Xuple
    xuplesByStatus map[XupleStatus]map[string]*Xuple
    
    // Politiques configurables
    selectionPolicy   SelectionPolicy
    consumptionPolicy ConsumptionPolicy
    retentionPolicy   RetentionPolicy
    
    // Statistiques
    stats XupleSpaceStats
    
    // Synchronisation
    mu sync.RWMutex
    
    // Timestamp de création
    createdAt time.Time
}

// XupleSpaceStats contient les statistiques d'un xuple-space
type XupleSpaceStats struct {
    TotalCreated     int64 // Nombre total de xuples créés
    TotalConsumed    int64 // Nombre total de consommations
    TotalExpired     int64 // Nombre total de xuples expirés
    CurrentPending   int64 // Nombre actuel de xuples pending
    CurrentConsumed  int64 // Nombre actuel de xuples consumed
    CurrentExpired   int64 // Nombre actuel de xuples expired
    LastCleanupAt    time.Time
    LastConsumptionAt time.Time
}
```

**Justification** :
- `name` : Identifie le xuple-space (peut y avoir plusieurs espaces)
- Indexation multiple : Performances O(1) pour recherches fréquentes
- Politiques injectées : Découplage, extensibilité
- Stats : Observabilité et monitoring
- `sync.RWMutex` : Thread-safety avec lecture parallèle

---

## 💾 Considérations Mémoire

**Estimation** : 240 bytes/xuple + références  
**10,000 xuples** : ~2.4 MB + données

---

## 📚 Références

- [common.md](../../../.github/prompts/common.md)
- [02-design-xuples-architecture.md](../../../scripts/xuples/02-design-xuples-architecture.md)

---

**Prochaine étape** : [02-interfaces.md](./02-interfaces.md)
