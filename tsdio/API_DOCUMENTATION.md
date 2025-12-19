# API TSD - Documentation Complète

## 🎯 Vue d'Ensemble

L'API TSD permet d'interagir avec le moteur de règles via des structures de données sérialisables en JSON.

**⚠️ IMPORTANT** : L'identifiant interne des faits (`_id_`) est géré automatiquement par le système et n'est **JAMAIS** exposé via l'API publique JSON.

---

## 📋 Structures de Données

### Fact

Représente un fait dans le système TSD.

```go
type Fact struct {
    // internalID - ID interne du fait (CACHÉ, non sérialisé en JSON)
    Type   string                 `json:"type"`
    Fields map[string]interface{} `json:"fields"`
}
```

**Exemple JSON** :
```json
{
  "type": "User",
  "fields": {
    "name": "Alice",
    "age": 30,
    "email": "alice@example.com"
  }
}
```

**Restrictions** :
- ❌ Le champ `_id_` ne peut **jamais** être défini manuellement
- ❌ Le champ `_id_` n'apparaît **jamais** dans le JSON
- ✅ L'ID est généré automatiquement lors de l'insertion dans le moteur
- ✅ L'ID peut être récupéré en interne via `GetInternalID()` pour usage système uniquement

**Méthodes** :
- `GetInternalID() string` - Retourne l'ID interne (usage interne uniquement)
- `SetInternalID(id string)` - Définit l'ID interne (usage interne uniquement)

---

### ExecuteRequest

Représente une requête d'exécution de programme TSD.

```go
type ExecuteRequest struct {
    Source     string `json:"source"`
    SourceName string `json:"source_name,omitempty"`
    Verbose    bool   `json:"verbose,omitempty"`
}
```

**Champs** :
- `source` (string, requis) : Code TSD à exécuter (types, facts, rules)
- `source_name` (string, optionnel) : Nom de la source pour messages d'erreur (défaut: `<request>`)
- `verbose` (bool, optionnel) : Active le mode verbeux (défaut: `false`)

**Exemple** :
```json
{
  "source": "type User : <name: string, age: number>\nfact alice : User <name: \"Alice\", age: 30>",
  "source_name": "example.tsd",
  "verbose": true
}
```

---

### ExecuteResponse

Représente la réponse d'une exécution.

```go
type ExecuteResponse struct {
    Success         bool              `json:"success"`
    Error           string            `json:"error,omitempty"`
    ErrorType       string            `json:"error_type,omitempty"`
    Results         *ExecutionResults `json:"results,omitempty"`
    ExecutionTimeMs int64             `json:"execution_time_ms"`
}
```

**Cas de succès** :
```json
{
  "success": true,
  "results": {
    "facts_count": 5,
    "activations_count": 2,
    "activations": [...]
  },
  "execution_time_ms": 150
}
```

**Cas d'erreur** :
```json
{
  "success": false,
  "error": "syntax error at line 5",
  "error_type": "parsing_error",
  "execution_time_ms": 25
}
```

**Types d'erreurs** :
- `parsing_error` : Erreur de syntaxe TSD
- `validation_error` : Erreur de validation de types ou contraintes
- `execution_error` : Erreur lors de l'exécution
- `server_error` : Erreur interne du serveur

---

### ExecutionResults

Contient les détails des résultats d'exécution.

```go
type ExecutionResults struct {
    FactsCount       int          `json:"facts_count"`
    ActivationsCount int          `json:"activations_count"`
    Activations      []Activation `json:"activations"`
}
```

**Exemple** :
```json
{
  "facts_count": 5,
  "activations_count": 3,
  "activations": [
    {
      "action_name": "greet",
      "arguments": [
        {"position": 0, "value": "Alice", "type": "string"}
      ],
      "triggering_facts": [
        {
          "type": "User",
          "fields": {"name": "Alice", "age": 30}
        }
      ],
      "bindings_count": 1
    }
  ]
}
```

**Note** : Les `triggering_facts` ne contiennent **jamais** le champ `_id_`.

---

### Activation

Représente une action déclenchée avec ses détails.

```go
type Activation struct {
    ActionName      string          `json:"action_name"`
    Arguments       []ArgumentValue `json:"arguments"`
    TriggeringFacts []Fact          `json:"triggering_facts"`
    BindingsCount   int             `json:"bindings_count"`
}
```

**Champs** :
- `action_name` : Nom de l'action déclenchée
- `arguments` : Arguments évalués de l'action
- `triggering_facts` : Faits ayant déclenché l'action (sans `_id_`)
- `bindings_count` : Nombre de bindings dans le token

---

### ArgumentValue

Représente un argument évalué d'une action.

```go
type ArgumentValue struct {
    Position int         `json:"position"`
    Value    interface{} `json:"value"`
    Type     string      `json:"type"`
}
```

**Types possibles** :
- `string` : Chaîne de caractères
- `number` : Nombre (int ou float)
- `bool` : Booléen
- `identifier` : Identifiant
- `variable` : Variable

---

## 🔒 Sécurité et Validation

### Champ Réservé `_id_`

Le champ `_id_` est **strictement réservé** au système :

✅ **Ce qui est autorisé** :
- Créer des faits sans spécifier `_id_`
- Récupérer l'ID en interne via `GetInternalID()` (usage système uniquement)

❌ **Ce qui est interdit** :
- Définir manuellement le champ `_id_` dans un fait
- Inclure `_id_` dans les champs d'un fait
- Accéder à `_id_` dans les expressions TSD

**Exemple d'erreur** :
```json
// ❌ INTERDIT - Génèrera une erreur
{
  "type": "User",
  "fields": {
    "_id_": "manual-id",  // ❌ Champ réservé
    "name": "Alice"
  }
}
```

**Message d'erreur** :
```
le champ '_id_' est réservé et ne peut pas être défini manuellement
```

---

## 📊 Exemples d'Utilisation

### Exécuter un Programme Simple

**Requête** :
```json
{
  "source": "type Person : <name: string>\nfact alice : Person <name: \"Alice\">\naction greet(p: Person) { print(p.name) }",
  "verbose": false
}
```

**Réponse** :
```json
{
  "success": true,
  "results": {
    "facts_count": 1,
    "activations_count": 1,
    "activations": [
      {
        "action_name": "greet",
        "arguments": [
          {"position": 0, "value": "Alice", "type": "string"}
        ],
        "triggering_facts": [
          {
            "type": "Person",
            "fields": {"name": "Alice"}
          }
        ],
        "bindings_count": 1
      }
    ]
  },
  "execution_time_ms": 50
}
```

**Note** : Le fait dans `triggering_facts` ne contient **pas** `_id_`.

---

### Gestion d'Erreur

**Requête avec erreur de syntaxe** :
```json
{
  "source": "type User : <name: string\nfact bob : User <name: \"Bob\">",
  "source_name": "invalid.tsd"
}
```

**Réponse** :
```json
{
  "success": false,
  "error": "syntax error: expected '>' at line 1",
  "error_type": "parsing_error",
  "execution_time_ms": 15
}
```

---

## 🔄 Migration depuis l'Ancienne API

### Changements Breaking

#### Avant (❌ Ancien) :
```go
fact := tsdio.Fact{
    ID:   "user-1",  // ❌ Champ public exposé
    Type: "User",
    Fields: map[string]interface{}{
        "name": "Alice",
    },
}
```

#### Après (✅ Nouveau) :
```go
fact := tsdio.Fact{
    Type: "User",
    Fields: map[string]interface{}{
        "name": "Alice",
    },
}
// L'ID interne est défini automatiquement par le système
// ou manuellement en interne via :
fact.SetInternalID("user-1")  // Usage interne uniquement
```

### JSON Serialization

#### Avant (❌ Ancien) :
```json
{
  "_id_": "user-1",  // ❌ Exposé publiquement
  "type": "User",
  "fields": {"name": "Alice"}
}
```

#### Après (✅ Nouveau) :
```json
{
  "type": "User",
  "fields": {"name": "Alice"}
}
```

**Note** : `_id_` est **complètement caché** de l'API JSON publique.

---

## 🧪 Tests

### Test de Sérialisation JSON

```go
func TestFact_JSONSerialization(t *testing.T) {
    fact := tsdio.Fact{
        Type: "User",
        Fields: map[string]interface{}{
            "name": "Alice",
            "age":  30,
        },
    }
    fact.SetInternalID("User~Alice")
    
    jsonData, _ := json.Marshal(fact)
    jsonStr := string(jsonData)
    
    // ✅ Vérifier que _id_ n'est PAS dans le JSON
    if strings.Contains(jsonStr, "_id_") {
        t.Error("_id_ should be hidden from JSON")
    }
    
    // ✅ Vérifier que les champs sont présents
    if !strings.Contains(jsonStr, "Alice") {
        t.Error("fields should be in JSON")
    }
}
```

### Test des Méthodes ID

```go
func TestFact_InternalIDMethods(t *testing.T) {
    fact := tsdio.Fact{
        Type: "User",
        Fields: map[string]interface{}{"name": "Alice"},
    }
    
    // ID initial vide
    assert.Equal(t, "", fact.GetInternalID())
    
    // Définir un ID
    fact.SetInternalID("User~Alice")
    
    // Vérifier l'ID
    assert.Equal(t, "User~Alice", fact.GetInternalID())
}
```

---

## 📚 Références

- **Code source** : `tsdio/api.go`
- **Tests** : `tsdio/api_test.go`
- **Constantes** : `constraint/constraint_constants.go`
- **Validation** : `constraint/parser.go`

---

## ⚠️ Notes Importantes

1. **Sécurité** : Le champ `_id_` ne doit **JAMAIS** être accessible aux utilisateurs finaux
2. **Cohérence** : Tous les faits générés par le système utilisent `SetInternalID()` en interne
3. **Validation** : Le parser TSD rejette toute tentative d'utiliser `_id_` dans les expressions
4. **Sérialisation** : Le JSON ne contient **jamais** le champ `_id_`
5. **Migration** : Tout code utilisant `fact.ID` doit être mis à jour pour utiliser `fact.GetInternalID()`

---

**Version** : 1.2.0  
**Date** : 2025-12-19  
**Auteur** : TSD Contributors  
**Licence** : MIT
