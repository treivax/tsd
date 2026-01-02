# 🐛 Debug Report - builtin_integration_test.go

**Date**: 2025-12-20  
**Test**: `TestBuiltinActions_EndToEnd_XupleAction`  
**Package**: `github.com/treivax/tsd/rete/actions`  
**Status**: ✅ **RÉSOLU**

---

## 📋 Résumé Exécutif

Le test `TestBuiltinActions_EndToEnd_XupleAction` échouait avec l'erreur :
```
❌ Failed to mark consumed: xuple not available for consumption
```

**Cause racine** : Le test tentait d'appeler manuellement `MarkConsumed()` après `Retrieve()`, mais `Retrieve()` **marque déjà automatiquement** le xuple comme consommé.

**Solution** : Suppression de l'appel manuel à `MarkConsumed()` et ajout de tests de validation pour vérifier le comportement correct de la politique `per-agent`.

---

## 🔍 Analyse Détaillée

### Comportement Attendu vs Observé

#### Flux du Code Original (Erroné)
```go
// 1. Récupérer le xuple
retrievedAlert, err := criticalSpace.Retrieve("agent1")

// 2. ❌ ERREUR : Tentative de marquer comme consommé manuellement
err = criticalSpace.MarkConsumed(retrievedAlert.ID, "agent1")
// -> ÉCHEC : xuple déjà marqué consommé par agent1 (via Retrieve)
```

#### Flux Correct
```go
// 1. Récupérer le xuple - MARQUE AUTOMATIQUEMENT COMME CONSOMMÉ
retrievedAlert, err := criticalSpace.Retrieve("agent1")
// Le xuple est déjà marqué consommé par agent1 à ce stade

// 2. ✅ Vérifier que le xuple a bien été consommé
if _, consumed := retrievedAlert.Metadata.ConsumedBy["agent1"]; !consumed {
    t.Error("agent1 devrait être dans ConsumedBy")
}
```

### Mécanisme Interne

#### Dans `xuplespace.go::Retrieve()` (ligne 120-125)
```go
// CORRECTION DU BUG CRITIQUE : Marquer automatiquement comme consommé
// Cela évite que l'appelant oublie d'appeler MarkConsumed() et garantit
// que la politique de consommation 'once' fonctionne correctement
selected.markConsumedBy(agentID)

if xs.config.ConsumptionPolicy.OnConsumed(selected, agentID) {
    selected.Metadata.State = XupleStateConsumed
}
```

#### Dans `xuplespace.go::MarkConsumed()` (ligne 254)
```go
if !xuple.CanBeConsumedBy(agentID, xs.config.ConsumptionPolicy) {
    return ErrXupleNotAvailable  // ← Erreur retournée
}
```

#### Dans `policy_consumption.go::PerAgentConsumptionPolicy.CanConsume()` (ligne 45-49)
```go
func (p *PerAgentConsumptionPolicy) CanConsume(xuple *Xuple, agentID string) bool {
    if xuple.Metadata.ConsumedBy == nil {
        return true
    }
    _, alreadyConsumed := xuple.Metadata.ConsumedBy[agentID]
    return !alreadyConsumed  // ← false si déjà consommé par cet agent
}
```

### Séquence d'Événements

1. **Retrieve("agent1")** appelé
   - Trouve xuple A002 (LIFO)
   - Appelle `markConsumedBy("agent1")`
   - Ajoute `agent1` dans `ConsumedBy` map
   - Retourne le xuple

2. **MarkConsumed(xupleID, "agent1")** appelé ❌
   - Vérifie `CanBeConsumedBy("agent1", policy)`
   - Policy vérifie : agent1 dans `ConsumedBy`? **OUI**
   - Retourne `!alreadyConsumed` = **false**
   - MarkConsumed retourne `ErrXupleNotAvailable`

---

## 🔧 Modifications Apportées

### Fichier : `rete/actions/builtin_integration_test.go`

#### Changement 1 : Suppression de l'appel erroné à MarkConsumed

**Avant** (lignes 549-558) :
```go
// Marquer comme consommée
if retrievedAlert != nil {
    err = criticalSpace.MarkConsumed(retrievedAlert.ID, "agent1")
    if err != nil {
        t.Errorf("❌ Failed to mark consumed: %v", err)
    } else {
        t.Log("✅ Alerte marquée comme consommée par agent1")
    }
```

**Après** (lignes 549-557) :
```go
// Note: Retrieve() marque automatiquement le xuple comme consommé
// Pas besoin d'appeler MarkConsumed() manuellement
if retrievedAlert != nil {
    t.Log("✅ Alerte automatiquement marquée comme consommée par agent1 (via Retrieve)")

    // Vérifier que le xuple a bien été marqué comme consommé par agent1
    if _, consumed := retrievedAlert.Metadata.ConsumedBy["agent1"]; !consumed {
        t.Errorf("❌ agent1 devrait être dans ConsumedBy après Retrieve")
    }
```

#### Changement 2 : Amélioration de la validation per-agent policy

**Avant** (lignes 560-567) :
```go
// Per-agent policy: un autre agent devrait pouvoir récupérer le même xuple
retrievedAlert2, err := criticalSpace.Retrieve("agent2")
if err != nil {
    t.Errorf("❌ Failed to retrieve for agent2: %v", err)
} else {
    t.Logf("✅ Agent2 a récupéré alerte: %s (per-agent policy fonctionne)", 
           retrievedAlert2.Fact.ID)
}
```

**Après** (lignes 560-585) :
```go
// Per-agent policy: un autre agent devrait pouvoir récupérer le même xuple
retrievedAlert2, err := criticalSpace.Retrieve("agent2")
if err != nil {
    t.Errorf("❌ Failed to retrieve for agent2: %v", err)
} else {
    t.Logf("✅ Agent2 a récupéré alerte: %s (per-agent policy fonctionne)", 
           retrievedAlert2.Fact.ID)

    // Vérifier que c'est bien le même xuple (per-agent permet ça)
    if retrievedAlert2.ID != retrievedAlert.ID {
        t.Errorf("❌ Agent2 devrait obtenir le même xuple que agent1, got %s vs %s",
            retrievedAlert2.ID, retrievedAlert.ID)
    }
}

// Agent1 peut récupérer UN AUTRE xuple (A001) car il y en a 2 dans l'espace
// mais ne peut pas récupérer A002 à nouveau (déjà consommé par agent1)
retrievedAlert3, err := criticalSpace.Retrieve("agent1")
if err != nil {
    t.Errorf("❌ agent1 devrait pouvoir récupérer l'autre xuple disponible: %v", err)
} else {
    t.Logf("✅ agent1 a récupéré un autre xuple: %s", retrievedAlert3.Fact.ID)

    // Devrait être un xuple différent (l'autre alerte)
    if retrievedAlert3.ID == retrievedAlert.ID {
        t.Errorf("❌ agent1 a récupéré le même xuple deux fois (violation per-agent policy)")
    }
}
```

---

## ✅ Validation

### Tests Exécutés

```bash
# Test spécifique
go test -v ./rete/actions -run TestBuiltinActions_EndToEnd_XupleAction
```

**Résultat** : ✅ **PASS**

```
=== RUN   TestBuiltinActions_EndToEnd_XupleAction
    builtin_integration_test.go:552: ✅ Alerte automatiquement marquée comme consommée par agent1 (via Retrieve)
    builtin_integration_test.go:564: ✅ Agent2 a récupéré alerte: A002 (per-agent policy fonctionne)
    builtin_integration_test.go:579: ✅ agent1 a récupéré un autre xuple: A001
    builtin_integration_test.go:626: 🎉 Tests de l'action Xuple validés avec succès!
--- PASS: TestBuiltinActions_EndToEnd_XupleAction (0.00s)
PASS
```

### Suite de Tests Complète

```bash
# Tous les tests du package actions
go test ./rete/actions/...
```

**Résultat** : ✅ **PASS** (18/18 tests)

---

## 📚 Leçons Apprises

### 1. Comportement Auto-Consommation

**Documentation claire** : `Retrieve()` marque **automatiquement** les xuples comme consommés. L'appelant ne doit **jamais** appeler `MarkConsumed()` manuellement après `Retrieve()`.

**Raison de conception** (commentaire dans le code) :
> "CORRECTION DU BUG CRITIQUE : Marquer automatiquement comme consommé
> Cela évite que l'appelant oublie d'appeler MarkConsumed() et garantit
> que la politique de consommation 'once' fonctionne correctement"

### 2. Politique Per-Agent

Avec `PerAgentConsumptionPolicy` :
- ✅ Plusieurs agents **différents** peuvent récupérer le **même xuple**
- ❌ Le **même agent** ne peut **pas** récupérer le même xuple deux fois
- ✅ Un agent peut récupérer des **xuples différents** plusieurs fois

### 3. Tests Existants à Consulter

Le fichier `xuples/xuplespace_consumption_test.go` contient des exemples **corrects** :
- `TestRetrieveAutomaticallyMarksConsumed` (ligne 54-148)
- `TestRetrievePerAgentPolicy` (ligne 150-260)

Ces tests montrent le pattern correct d'utilisation de `Retrieve()`.

---

## 🎯 Pattern Correct d'Utilisation

### ✅ BON

```go
// Récupérer un xuple (auto-consomme)
xuple, err := space.Retrieve("agent1")
if err != nil {
    return err
}

// Vérifier la consommation si nécessaire
if _, consumed := xuple.Metadata.ConsumedBy["agent1"]; !consumed {
    t.Error("xuple devrait être marqué consommé")
}
```

### ❌ MAUVAIS

```go
// Récupérer un xuple
xuple, err := space.Retrieve("agent1")

// ❌ NE PAS FAIRE : Double consommation
err = space.MarkConsumed(xuple.ID, "agent1")
// -> ERREUR : xuple not available for consumption
```

---

## 📖 Références

- **Fichier de test fixé** : `rete/actions/builtin_integration_test.go`
- **Implémentation** : `xuples/xuplespace.go::Retrieve()` (ligne 89-130)
- **Politiques** : `xuples/policy_consumption.go`
- **Tests de référence** : `xuples/xuplespace_consumption_test.go`

---

## 🔒 Checklist Post-Fix

- [x] Test spécifique passe
- [x] Suite complète du package passe
- [x] Validation de la politique per-agent ajoutée
- [x] Commentaires explicatifs ajoutés dans le code
- [x] Pas de régression sur les autres tests
- [x] Documentation du comportement auto-consommation

---

**Statut Final** : ✅ **RÉSOLU ET VALIDÉ**