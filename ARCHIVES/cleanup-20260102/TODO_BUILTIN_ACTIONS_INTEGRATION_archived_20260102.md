# TODO - Intégration des Actions Builtin (Update/Insert/Retract)

**Priorité** : Moyenne  
**Difficulté** : Faible  
**Impact** : Élevé  
**Statut** : À faire

---

## 📋 Contexte

Les actions natives `Update`, `Insert` et `Retract` sont implémentées dans `rete/actions/builtin.go` 
via le `BuiltinActionExecutor`, mais elles ne sont **pas intégrées** dans le pipeline API.

### Symptôme Observé
Lors de l'exécution du test `TestRelationshipStatusE2E_ThreeSteps` :
```
📋 ACTION: Update(...)
📋 ACTION NON DÉFINIE (log uniquement): Update(Personne{Personne_1})
```

Les règles se déclenchent correctement, mais les actions ne sont pas exécutées.

### Cause Racine
1. Le `BuiltinActionExecutor` existe et fonctionne ✅
2. Mais il n'est **pas enregistré** dans l'`ActionExecutor` du réseau RETE ❌
3. L'`ActionExecutor` n'enregistre que `Print` et `Xuple` par défaut
4. Les actions non enregistrées sont loguées sans être exécutées

---

## 🎯 Objectif

Intégrer les actions `Update`, `Insert` et `Retract` dans le pipeline API 
pour qu'elles soient automatiquement disponibles lors de l'ingestion de fichiers TSD.

---

## 📝 Tâches à Réaliser

### 1. Créer des Wrappers ActionHandler ⭐⭐⭐

**Fichier** : `tsd/rete/actions/builtin_handlers.go` (nouveau)

Créer des wrappers qui implémentent l'interface `ActionHandler` pour chaque action builtin :

```go
package actions

import "github.com/treivax/tsd/rete"

// UpdateActionHandler est un wrapper pour l'action Update
type UpdateActionHandler struct {
    executor *BuiltinActionExecutor
}

func NewUpdateActionHandler(executor *BuiltinActionExecutor) *UpdateActionHandler {
    return &UpdateActionHandler{executor: executor}
}

func (h *UpdateActionHandler) Name() string {
    return "Update"
}

func (h *UpdateActionHandler) Execute(args []interface{}, token *rete.Token) error {
    return h.executor.Execute("Update", args, token)
}

// Idem pour Insert et Retract
type InsertActionHandler struct { ... }
type RetractActionHandler struct { ... }
```

**Temps estimé** : 1h

---

### 2. Modifier le Pipeline API ⭐⭐

**Fichier** : `tsd/api/pipeline.go`

Intégrer le `BuiltinActionExecutor` dans la création du pipeline :

```go
func NewPipelineWithConfig(config *Config) *Pipeline {
    // ... code existant ...
    
    storage := rete.NewMemoryStorage()
    network := rete.NewReteNetwork(storage)
    xupleManager := xuples.NewXupleManager()
    
    // Créer le BuiltinActionExecutor
    logger := createLogger(config.LogLevel)
    builtinExecutor := actions.NewBuiltinActionExecutor(
        network, 
        xupleManager, 
        os.Stdout, 
        logger.ToStdLogger(), // Convertir rete.Logger en *log.Logger
    )
    
    // Enregistrer les actions builtin dans le réseau
    actionExecutor := network.GetActionExecutor()
    actionExecutor.RegisterAction(actions.NewUpdateActionHandler(builtinExecutor))
    actionExecutor.RegisterAction(actions.NewInsertActionHandler(builtinExecutor))
    actionExecutor.RegisterAction(actions.NewRetractActionHandler(builtinExecutor))
    
    // ... reste du code ...
}
```

**Points d'attention** :
- Vérifier la compatibilité des types de logger
- S'assurer que le `BuiltinActionExecutor` est créé après le `XupleManager`
- Tester que les actions sont bien enregistrées

**Temps estimé** : 2h

---

### 3. Ajouter une Méthode ToStdLogger ⭐

**Fichier** : `tsd/rete/logger.go`

Si elle n'existe pas déjà, ajouter une méthode pour convertir `rete.Logger` en `*log.Logger` :

```go
func (l *Logger) ToStdLogger() *log.Logger {
    return log.New(l.writer, "", log.LstdFlags)
}
```

**Temps estimé** : 30min

---

### 4. Tests d'Intégration ⭐⭐⭐

**Fichier** : `tsd/tests/integration/builtin_actions_test.go` (nouveau)

Créer des tests pour vérifier que les actions fonctionnent via le pipeline API :

```go
func TestBuiltinActions_Update_Integration(t *testing.T) {
    program := `
type Person(#id: string, name: string, age: number)

rule update_age : {p: Person} / p.age < 18 ==>
    Update(Person(id: p.id, name: p.name, age: 18))

Person(id: "p1", name: "Alice", age: 15)
`
    
    pipeline := api.NewPipeline()
    result, err := pipeline.IngestString(program)
    require.NoError(t, err)
    
    // Vérifier que le fait a été modifié
    network := result.Network()
    facts := network.Storage.GetAllFacts()
    // ... assertions ...
}
```

**Tests à créer** :
- `TestBuiltinActions_Update_Integration` ✅
- `TestBuiltinActions_Insert_Integration` ✅
- `TestBuiltinActions_Retract_Integration` ✅
- `TestBuiltinActions_Combined_Integration` ✅

**Temps estimé** : 3h

---

### 5. Activer les Assertions du Test E2E ⭐

**Fichier** : `tsd/tests/e2e/relationship_status_e2e_test.go`

Une fois l'intégration terminée, activer les assertions commentées :

```go
// AVANT (actuel)
require.Equal(t, "", alain.Fields["statut"],
    "LIMITATION: Le statut d'Alain reste vierge (action Update non exécutée)")

// APRÈS (une fois intégré)
require.Equal(t, "en couple", alain.Fields["statut"],
    "Le statut d'Alain doit avoir été modifié à 'en couple' par la règle")
```

**Temps estimé** : 30min

---

### 6. Documentation ⭐

**Fichiers à mettre à jour** :
- `README.md` : Mentionner les actions Update/Insert/Retract
- `docs/actions.md` : Documenter l'utilisation de ces actions
- `CHANGELOG.md` : Ajouter l'entrée pour cette fonctionnalité

**Temps estimé** : 1h

---

## 🔍 Points de Vigilance

### 1. Thread-Safety
- Le `BuiltinActionExecutor` doit être thread-safe
- Vérifier que le réseau RETE ne subit pas de modifications concurrentes

### 2. Gestion des Erreurs
- Les actions Update/Insert/Retract peuvent échouer (fait inexistant, etc.)
- S'assurer que les erreurs remontent correctement au test

### 3. Inline Facts
- Les actions utilisent des "inline facts" : `Update(Person(id: "p1", ...))`
- Vérifier que l'ID est correctement extrait et utilisé

### 4. Compatibilité
- Vérifier que les tests existants ne sont pas cassés
- Tester avec différentes configurations de pipeline

---

## ✅ Critères de Complétion

- [ ] Les wrappers `ActionHandler` sont créés
- [ ] Le `BuiltinActionExecutor` est intégré dans le pipeline
- [ ] Les actions Update/Insert/Retract sont enregistrées automatiquement
- [ ] Les tests d'intégration passent
- [ ] Le test e2e `TestRelationshipStatusE2E_ThreeSteps` passe avec les assertions activées
- [ ] La documentation est mise à jour
- [ ] Tous les tests existants passent toujours
- [ ] Code review effectuée

---

## 📊 Estimation Totale

**Temps de développement** : 8h  
**Temps de tests** : 2h  
**Temps de documentation** : 1h  
**Total** : ~11h (environ 1.5 jour)

---

## 🔗 Références

### Code Existant
- `tsd/rete/actions/builtin.go` - Implémentation des actions
- `tsd/rete/actions/builtin_test.go` - Tests unitaires
- `tsd/rete/action_executor.go` - Gestionnaire d'actions
- `tsd/api/pipeline.go` - Pipeline API

### Tests
- `tsd/tests/e2e/relationship_status_e2e_test.go` - Test e2e à activer
- `tsd/tests/shared/testutil.go` - Utilitaires de test

### Documentation
- `tsd/tests/e2e/testdata/README_relationship_test.md` - Doc du test e2e
- `tsd/RAPPORT_TEST_E2E_RELATIONS.md` - Rapport détaillé
- `.github/prompts/test.md` - Standards de tests

---

## 🚀 Exemple d'Utilisation Future

Une fois l'intégration terminée, les utilisateurs pourront écrire :

```tsd
type Product(#id: string, name: string, stock: number, status: string)

// Règle : Marquer les produits en rupture
rule mark_out_of_stock : {p: Product} / p.stock == 0 AND p.status != "out_of_stock" ==>
    Update(Product(id: p.id, name: p.name, stock: 0, status: "out_of_stock"))

// Règle : Supprimer les produits obsolètes
rule remove_obsolete : {p: Product} / p.status == "obsolete" ==>
    Retract(p.id)

// Règle : Créer une alerte pour stock faible
rule low_stock_alert : {p: Product} / p.stock < 10 AND p.stock > 0 ==>
    Insert(Alert(productId: p.id, level: "warning", message: "Stock faible"))
```

Et les actions seront **automatiquement exécutées** ! 🎉

---

## 📞 Contact

Pour toute question ou assistance sur cette intégration :
- Voir le rapport complet : `RAPPORT_TEST_E2E_RELATIONS.md`
- Consulter les tests existants dans `rete/actions/`
