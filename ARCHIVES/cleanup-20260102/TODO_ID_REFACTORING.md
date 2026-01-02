# TODO: Mise à jour suite au changement de `id` vers `_id_`

## 📋 Résumé du changement

Le champ d'identifiant interne a été renommé de `id` à `_id_` pour:
- ✅ Le cacher complètement des expressions TSD
- ✅ Le rendre inaccessible aux utilisateurs
- ✅ Éviter la confusion avec des champs utilisateur nommés `id`

## ✅ Modifications effectuées

### Contraintes (constraint/)
- ✅ Constante `FieldNameInternalID = "_id_"` créée
- ✅ Validation interdisant `_id_` dans les types
- ✅ Validation interdisant `_id_` dans les faits
- ✅ Validation interdisant l'accès à `_id_` dans les expressions
- ✅ Génération automatique de `_id_` (jamais manuel)
- ✅ Tous les tests constraint/ passent

### RETE (rete/)
- ✅ Constante `FieldNameID = "_id_"` mise à jour
- ✅ Structure `Fact` avec tag JSON `json:"_id_"`
- ✅ Évaluateur interdit l'accès à `_id_`

### API (tsdio/)
- ✅ Structure `Fact` avec tag JSON `json:"_id_"`

## ⚠️ Actions requises

### Tests RETE (rete/)

De nombreux tests RETE utilisent encore l'accès au champ "id" dans les expressions TSD.
Ceci est maintenant **interdit** et provoque les échecs de tests.

**Fichiers à corriger:**
```bash
# Tests qui échouent actuellement
./rete/action_arithmetic_e2e_test.go
./rete/aggregation_test.go
./rete/alpha_chain_test.go
# ... et autres tests RETE
```

**Stratégies de correction:**

1. **Option A - Comparaisons via champs de type Fait (recommandé)**
   
   Attendre l'implémentation du prompt suivant qui ajoutera le support des champs de type Fait
   permettant les comparaisons comme `p.user == u` sans accès direct à `_id_`.

2. **Option B - Utiliser d'autres champs pour les tests**
   
   Remplacer les tests qui utilisent `id` par des tests utilisant d'autres champs:
   ```tsd
   # Avant (INTERDIT)
   p.id == c.produit_id
   
   # Après (VALIDE - si les champs existent)
   p.nom == c.produit_nom
   ```

3. **Option C - Ajuster les types de tests**
   
   Déclarer `id` comme champ explicite dans les types de test:
   ```tsd
   type Produit(#id: string, nom: string, prix: number)
   ```
   
   Note: Ceci crée un champ utilisateur `id` différent du champ interne `_id_`.

### Code appelant hors tests

Tout code qui utilise `fact["id"]` ou `fact.ID` doit être mis à jour:

```go
// ❌ Ancien code (peut ne plus fonctionner)
id := fact["id"]
if fact.ID == "expected" { ... }

// ✅ Nouveau code
id := fact[constraint.FieldNameInternalID]  // ou rete.FieldNameID
if fact.ID == "expected" { ... }  // OK - accès interne
```

**IMPORTANT:** Le champ `_id_` ne doit **JAMAIS** être exposé dans les expressions TSD ou l'API publique.

## 📝 Exemples de corrections

### Exemple 1: Test avec jointure sur ID

```go
// ❌ AVANT (NE COMPILE PLUS)
input := `
type Produit(#nom: string, prix: number)
type Commande(produit_id: string, qte: number)

rule test:
  p: Produit, c: Commande
  where c.produit_id == p.id  // ❌ p.id interdit
  then log("Match")
`

// ✅ APRÈS - Solution temporaire avec champ explicite
input := `
type Produit(#id: string, #nom: string, prix: number)
type Commande(produit_id: string, qte: number)

Produit(id: "PROD1", nom: "Laptop", prix: 1000)
Commande(produit_id: "PROD1", qte: 2)

rule test:
  p: Produit, c: Commande
  where c.produit_id == p.id  // ✅ OK - p.id est un champ utilisateur
  then log("Match")
`

// 🎯 APRÈS - Solution future (prompt 02-parser-syntax.md)
input := `
type Produit(#nom: string, prix: number)
type Commande(produit: Produit, qte: number)  // Type Fait

Produit(nom: "Laptop", prix: 1000) as p1
Commande(produit: p1, qte: 2)  // Référence directe

rule test:
  p: Produit, c: Commande
  where c.produit == p  // Comparaison d'objets Fait
  then log("Match")
`
```

### Exemple 2: Test accédant à ID en interne

```go
// Dans les tests, pour vérifier les IDs générés:

// ❌ AVANT
factID := reteFact["id"].(string)

// ✅ APRÈS
factID := reteFact[constraint.FieldNameInternalID].(string)
// ou
factID := reteFact[rete.FieldNameID].(string)
```

## 🚀 Prochaines étapes

1. **Prompt 02** : Ajouter le support des champs de type Fait dans la grammaire
   - Permettra `c.produit: Produit` au lieu de `c.produit_id: string`
   - Permettra les comparaisons `c.produit == p`
   - Rendra les tests plus naturels et expressifs

2. **Tests RETE** : Une fois le prompt 02 implémenté
   - Corriger tous les tests RETE pour utiliser la nouvelle syntaxe
   - Ou déclarer explicitement les champs `id` dans les types de test

3. **Documentation** : Mettre à jour
   - README avec exemples de la nouvelle syntaxe
   - Guide de migration pour les utilisateurs
   - Documentation API

## 📊 Impact estimé

- **constraint/** : ✅ 100% complété et testé
- **rete/** : ⚠️ Tests à adapter (~40 tests)
- **tsdio/** : ✅ Structure mise à jour
- **api/** : ⚠️ Vérifier les sérialisations JSON
- **tests/integration/** : ⚠️ Tests E2E à vérifier

## 🔒 Règles à respecter

1. **JAMAIS** exposer `_id_` dans les expressions TSD
2. **JAMAIS** permettre l'accès utilisateur à `_id_`
3. **TOUJOURS** générer `_id_` automatiquement
4. **TOUJOURS** utiliser les constantes (`FieldNameInternalID`, `FieldNameID`)
5. Les champs utilisateur peuvent s'appeler `id` (différent de `_id_`)

## 📞 Contact

En cas de questions sur ce refactoring, se référer à:
- `.github/prompts/review.md` - Processus de revue
- `scripts/new_ids/01-prompt-structures-base.md` - Spécifications
- `.github/prompts/common.md` - Standards du projet
