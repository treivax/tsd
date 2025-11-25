# Guide de Migration - IDs Internes (Type_ID)

## Vue d'ensemble

Ce document décrit les changements importants apportés à l'API de gestion des faits dans le moteur RETE. Ces changements visent à éliminer les ambiguïtés et garantir des performances O(1) constantes.

## Motivation

Avant cette modification, l'API `GetFact()` acceptait à la fois des IDs simples ("f1") et des IDs internes ("Type_f1"), avec un fallback O(n) en cas d'échec de la recherche directe. Cela créait plusieurs problèmes :

1. **Ambiguïté** : Si deux types différents avaient le même ID ("Person_P1" et "Company_P1"), la recherche par ID simple "P1" pouvait retourner n'importe lequel
2. **Performance incohérente** : O(1) pour les IDs internes, O(n) pour les IDs simples
3. **API non déterministe** : Le comportement changeait selon l'ordre d'insertion des faits

## Changements d'API

### Avant

```go
// Ambiguë - acceptait ID simple ou ID interne
func (wm *WorkingMemory) GetFact(factID string) (*Fact, bool) {
    // Essayer ID interne direct (O(1))
    if fact, exists := wm.Facts[factID]; exists {
        return fact, true
    }
    // Fallback: chercher par ID simple (O(n))
    for _, fact := range wm.Facts {
        if fact.ID == factID {
            return fact, true
        }
    }
    return nil, false
}

// Ambiguë
func (wm *WorkingMemory) RemoveFact(factID string) {
    // Recherche similaire avec fallback
}
```

### Après

```go
// Claire - accepte uniquement ID interne
// Pour rechercher par type et ID séparément, utiliser GetFactByTypeAndID
func (wm *WorkingMemory) GetFact(internalID string) (*Fact, bool) {
    fact, exists := wm.Facts[internalID]
    return fact, exists  // O(1) garanti, pas de fallback
}

// Simple et direct
func (wm *WorkingMemory) RemoveFact(internalID string) {
    delete(wm.Facts, internalID)
}
```

## Fonctions Disponibles

### Recherche par ID Interne

```go
// Recherche directe par ID interne (Type_ID)
fact, exists := wm.GetFact("Person_P1")
if !exists {
    return fmt.Errorf("fait introuvable")
}
```

### Recherche par Type et ID

```go
// Si vous connaissez le type et l'ID séparément
fact, exists := wm.GetFactByTypeAndID("Person", "P1")
if !exists {
    return fmt.Errorf("fait introuvable")
}
```

### Génération d'ID Interne

```go
// Depuis un fait
internalID := fact.GetInternalID()  // Retourne "Person_P1"

// Depuis type et ID
internalID := MakeInternalID("Person", "P1")  // Retourne "Person_P1"
```

### Décomposition d'ID Interne

```go
// Parser un ID interne
factType, factID, ok := ParseInternalID("Person_P1")
if !ok {
    return fmt.Errorf("format d'ID interne invalide")
}
// factType = "Person", factID = "P1"
```

## Migration du Code Existant

### Scénario 1 : Vous avez l'objet Fact

```go
// Avant
wm.GetFact(fact.ID)
wm.RemoveFact(fact.ID)

// Après
wm.GetFact(fact.GetInternalID())
wm.RemoveFact(fact.GetInternalID())
```

### Scénario 2 : Vous avez Type et ID

```go
// Avant
wm.GetFact(id)

// Après - Option 1 : Utiliser GetFactByTypeAndID
fact, exists := wm.GetFactByTypeAndID(factType, factID)

// Après - Option 2 : Construire l'ID interne
internalID := MakeInternalID(factType, factID)
fact, exists := wm.GetFact(internalID)
```

### Scénario 3 : Rétractation de Fait

```go
// Avant
network.RetractFact("P1")

// Après
network.RetractFact("Person_P1")

// Ou si vous avez le fait
network.RetractFact(fact.GetInternalID())
```

## Impact sur les Nœuds RETE

Tous les nœuds ont été mis à jour pour utiliser `fact.GetInternalID()` dans leurs méthodes `ActivateRetract()` :

- `RootNode.ActivateRetract(factID)` - factID doit être un ID interne
- `TypeNode.ActivateRetract(factID)` - factID doit être un ID interne
- `AlphaNode.ActivateRetract(factID)` - factID doit être un ID interne
- `JoinNode.ActivateRetract(factID)` - factID doit être un ID interne
- `ExistsNode.ActivateRetract(factID)` - factID doit être un ID interne
- `TerminalNode.ActivateRetract(factID)` - factID doit être un ID interne

## Structure WorkingMemory

La structure interne de `WorkingMemory` utilise maintenant exclusivement des IDs internes comme clés :

```go
type WorkingMemory struct {
    NodeID string
    Facts  map[string]*Fact  // Clé : Type_ID (ex: "Person_P1")
}
```

## Avantages

1. **Pas d'ambiguïté** : Chaque ID interne est unique globalement
2. **Performance constante** : O(1) garanti pour toutes les opérations
3. **API claire** : Une seule façon de faire les choses
4. **Type-safety** : Le format Type_ID encode le type dans l'identifiant

## Breaking Changes

⚠️ **Important** : Cette modification casse la compatibilité avec le code existant qui utilisait des IDs simples pour `GetFact()`, `RemoveFact()`, ou `RetractFact()`.

### Checklist de Migration

- [ ] Remplacer tous les appels à `GetFact(id)` par `GetFact(Type_id)` ou `GetFactByTypeAndID(type, id)`
- [ ] Remplacer tous les appels à `RemoveFact(id)` par `RemoveFact(Type_id)`
- [ ] Remplacer tous les appels à `RetractFact(id)` par `RetractFact(Type_id)`
- [ ] Mettre à jour les tests pour utiliser des IDs internes
- [ ] Vérifier que tous les nœuds personnalisés utilisent `fact.GetInternalID()` dans `ActivateRetract()`

## Exemple Complet

```go
package main

import (
    "fmt"
    "github.com/treivax/tsd/rete"
)

func main() {
    wm := &rete.WorkingMemory{
        NodeID: "test",
        Facts:  make(map[string]*rete.Fact),
    }

    // Ajouter des faits
    person := &rete.Fact{ID: "P1", Type: "Person"}
    company := &rete.Fact{ID: "C1", Type: "Company"}

    wm.AddFact(person)
    wm.AddFact(company)

    // Recherche par ID interne
    p, exists := wm.GetFact("Person_P1")
    if exists {
        fmt.Printf("Trouvé: %s\n", p.GetInternalID())
    }

    // Recherche par type et ID
    c, exists := wm.GetFactByTypeAndID("Company", "C1")
    if exists {
        fmt.Printf("Trouvé: %s\n", c.GetInternalID())
    }

    // Suppression
    wm.RemoveFact("Person_P1")

    // Vérification
    _, exists = wm.GetFact("Person_P1")
    fmt.Printf("Person_P1 existe: %v\n", exists)  // false
}
```

## Questions Fréquentes

### Q: Puis-je encore ajouter deux faits avec le même ID mais de types différents ?

R: Oui ! C'est justement l'intérêt des IDs internes. "Person_P1" et "Company_P1" coexistent sans conflit.

### Q: Comment gérer les faits venant de sources externes qui n'ont pas de type ?

R: Vous devez assigner un type au fait avant de l'ajouter à la WorkingMemory. Le système nécessite un type pour générer l'ID interne.

### Q: Que se passe-t-il si j'utilise un mauvais format d'ID interne ?

R: `GetFact()` retournera simplement `(nil, false)` car la clé n'existera pas dans la map. Utilisez `ParseInternalID()` pour valider le format si nécessaire.

### Q: Comment migrer une grande base de code ?

R: Cherchez tous les appels à `GetFact(`, `RemoveFact(`, et `RetractFact(` et vérifiez si le paramètre est un ID simple ou interne. Utilisez `git grep` ou votre IDE pour trouver toutes les occurrences.

## Support

Pour toute question ou problème lié à cette migration, veuillez consulter la documentation ou créer une issue sur le dépôt GitHub.

---

**Date de Migration** : 2025-01-XX  
**Version** : 1.0.0  
**Breaking Change** : Oui
