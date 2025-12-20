# Correctif de propagation des faits dans RETE

**Date:** 2025-12-20  
**Problème:** Les faits ne déclenchaient pas les règles avec références fact-to-fact (ex: `c.produit == p`)

---

## 🔍 Diagnostic

### Symptômes observés
- Les faits étaient correctement stockés dans les TypeNodes
- Les JoinNodes avaient des mémoires Left/Right **apparemment vides**
- Aucun token n'était généré par les jointures
- Les règles utilisant des références de faits (`c.produit == p`) ne se déclenchaient jamais

### Investigation
En ajoutant des logs de diagnostic, nous avons découvert que :
1. ✅ Les faits SE PROPAGEAIENT bien depuis TypeNodes → PassthroughAlpha → JoinNodes
2. ✅ Les tokens ÉTAIENT ajoutés aux mémoires Left/Right des JoinNodes
3. ❌ Les jointures ÉCHOUAIENT systématiquement lors de la comparaison des valeurs

---

## 🐛 Cause racine : Incohérence des formats d'ID

### Le problème

Le système utilisait **trois formats d'ID différents** de manière incohérente :

1. **Format de génération** (`ConvertFactsToReteFormat`) :
   - Génère des IDs au format `"Type~Value"` (ex: `"Produit~PROD001"`)
   - Stocke dans `factMap["_id_"]`

2. **Format de lecture** (`submitFactsFromGrammarWithMetrics`) :
   - Lisait `factMap["id"]` au lieu de `factMap["_id_"]` ❌
   - Résultat : `fact.ID = "PROD001"` (manque le préfixe de type)

3. **Format interne** (`Fact.GetInternalID()`) :
   - Retournait `fmt.Sprintf("%s_%s", f.Type, f.ID)` ❌
   - Avec `f.ID = "PROD001"`, cela donnait `"Produit_PROD001"` (underscore)
   - Mais le système attendait `"Produit~PROD001"` (tilde)

### Exemple de l'échec

Pour une commande avec `produit: p1` où `p1 = Produit(id: "PROD001", ...)` :

```
Condition de jointure : c.produit == p._id_

Valeurs comparées :
  - c.produit = "Produit~PROD001"  (correct, stocké dans le champ)
  - p._id_    = "PROD001"           (incorrect, manque le préfixe)

Résultat : "Produit~PROD001" != "PROD001" → jointure échoue ❌
```

---

## ✅ Solution appliquée

### 1. Corriger `Fact.GetInternalID()` (rete/fact_token.go)

**Avant :**
```go
func (f *Fact) GetInternalID() string {
    return fmt.Sprintf("%s_%s", f.Type, f.ID)
}
```

**Après :**
```go
func (f *Fact) GetInternalID() string {
    // L'ID est déjà au format "Type~Value" donc on le retourne tel quel
    return f.ID
}
```

**Justification :** L'ID du fait est déjà au format complet `"Type~Value"` généré par `GenerateFactID()`. Pas besoin de reconstruire.

### 2. Corriger la lecture de l'ID (rete/network_manager.go)

**Avant :**
```go
factID := fmt.Sprintf("fact_%d", i)
if id, ok := factMap["id"].(string); ok {
    factID = id
}
```

**Après :**
```go
factID := fmt.Sprintf("fact_%d", i)
// Utiliser _id_ qui contient l'ID interne complet (Type~Value)
if id, ok := factMap["_id_"].(string); ok {
    factID = id
} else if id, ok := factMap["id"].(string); ok {
    // Fallback pour compatibilité
    factID = id
}
```

**Justification :** `ConvertFactsToReteFormat` stocke l'ID complet dans `factMap["_id_"]`, pas dans `factMap["id"]`.

### 3. Utiliser GetInternalID() pour le champ "_id_" (rete/node_join.go)

**Avant :**
```go
if cond.RightField == "_id_" {
    rightValue = rightFact.ID  // ID incomplet
    rightExists = true
}
```

**Après :**
```go
if cond.RightField == "_id_" {
    // Utiliser l'ID interne complet (Type~Value) pour la comparaison
    rightValue = rightFact.GetInternalID()
    rightExists = true
}
```

**Justification :** La comparaison doit utiliser le format complet cohérent avec les valeurs stockées dans les champs.

---

## 📊 Résultats

### Tests qui passent maintenant

✅ **TestFactReferenceJoin** : Test de jointure fact-to-fact basique
- 2 Produits × 2 Commandes
- 2 jointures réussies (c.produit == p)

### Exemple de propagation réussie

```
TypeNode(Produit)
  ↓ ActivateRight(fact=Produit~PROD001)
PassthroughAlpha(left)
  ↓ ActivateLeft(token vars=[p])
JoinNode
  ├─ LeftMemory: 2 tokens (p=PROD001, p=PROD002)
  ↓ ActivateRight(fact=Commande~CMD001)
PassthroughAlpha(right)
  ↓ ActivateRight(fact=CMD001)
JoinNode
  ├─ RightMemory: 2 tokens (c=CMD001, c=CMD002)
  ├─ Comparaison: c.produit="Produit~PROD001" == p._id_="Produit~PROD001" ✅
  └─ ResultMemory: 2 tokens (jointures réussies)
TerminalNode
  └─ 2 exécutions ✅
```

---

## 🎯 Points clés à retenir

1. **Format unique d'ID** : `"Type~Value"` (tilde, pas underscore)
2. **Source de vérité** : `factMap["_id_"]` contient l'ID complet
3. **GetInternalID()** : Retourne directement `f.ID` (pas de reconstruction)
4. **Comparaisons** : Toujours utiliser `GetInternalID()` pour le champ spécial `"_id_"`

---

## 📝 Notes pour les développeurs

### Tests à mettre à jour

Les anciens tests qui créent manuellement des faits doivent utiliser le format correct :

**Avant :**
```go
fact := &Fact{ID: "f1", Type: "Person"}
```

**Après :**
```go
fact := &Fact{ID: "Person~f1", Type: "Person"}
```

### Convention de nommage

- **Séparateur type/valeur** : `~` (tilde) défini dans `constraint.IDSeparatorType`
- **Séparateur multi-valeurs** : `_` (underscore) défini dans `constraint.IDSeparatorValue`
- **Exemple clé composite** : `"Order~User123_Product456"`

---

## 🔗 Fichiers modifiés

1. `tsd/rete/fact_token.go` - Simplification de GetInternalID()
2. `tsd/rete/network_manager.go` - Lecture correcte de _id_
3. `tsd/rete/node_join.go` - Utilisation de GetInternalID() pour comparaisons

---

**Auteur:** Claude (assistance au debugging)  
**Reviewer:** [À compléter]