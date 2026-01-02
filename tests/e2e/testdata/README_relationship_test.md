# Test E2E - Modification de Faits via Relations

## 📋 Description

Ce test end-to-end démontre le workflow automatique complet de TSD en 3 étapes itératives, 
sans accès direct aux fonctions internes. Il teste la modification automatique de faits 
via des règles RETE déclenchées par l'ajout de relations.

## 🎯 Objectif

Vérifier que le système TSD peut :
1. Parser et ingérer des fichiers TSD successifs
2. Maintenir l'état du réseau RETE entre les ingestions
3. Déclencher automatiquement des règles lors de l'ajout de nouveaux faits
4. Afficher les résultats de manière lisible

## 📁 Fichiers

### `relationship_step1_types_rules.tsd`
**Contenu** : Définition des types et des règles

- **Type Personne** : Avec clé primaire sur `nom` et un champ `statut`
  - statut peut être : vide, 'célibataire' ou 'en couple'

- **Type Relation** : Pour mettre en relation deux personnes
  - personne1 : nom de la première personne
  - personne2 : nom de la deuxième personne
  - lien : type de relation ('pacs', 'mariage', 'union-libre', 'ennemis', etc.)

- **Règles** : Deux règles qui modifient automatiquement le statut des personnes
  - `mettre_en_couple_personne1` : Si personne1 est liée via une relation de couple
  - `mettre_en_couple_personne2` : Si personne2 est liée via une relation de couple

### `relationship_step2_persons.tsd`
**Contenu** : Ajout de 3 personnes avec statut vierge

- Alain (statut : "")
- Catherine (statut : "")
- Chantal (statut : "")

### `relationship_step3_relation.tsd`
**Contenu** : Ajout d'une relation de couple

- Relation entre Alain et Chantal (lien : "mariage")

## 🔄 Déroulement du Test

### Étape 1 : Définition des Types et Règles
```tsd
type Personne(#nom: string, statut: string)
type Relation(personne1: string, personne2: string, lien: string)

rule mettre_en_couple_personne1 : {p: Personne, r: Relation} /
    p.nom == r.personne1 AND
    (r.lien == "pacs" OR r.lien == "mariage" OR r.lien == "union-libre") AND
    p.statut != "en couple" ==>
    Update(Personne(nom: p.nom, statut: "en couple"))
```

**Résultat attendu** :
- ✅ 2 types définis
- ✅ 2 règles actives
- ✅ 0 fait

### Étape 2 : Ajout des Personnes
```tsd
Personne(nom: "Alain", statut: "")
Personne(nom: "Catherine", statut: "")
Personne(nom: "Chantal", statut: "")
```

**Résultat attendu** :
- ✅ 3 faits de type Personne
- ✅ Tous avec statut vierge
- ✅ IDs générés : Personne~Alain, Personne~Catherine, Personne~Chantal

### Étape 3 : Ajout d'une Relation
```tsd
Relation(personne1: "Alain", personne2: "Chantal", lien: "mariage")
```

**Résultat attendu** :
- ✅ 1 fait de type Relation
- ✅ Règles déclenchées pour Alain et Chantal
- ⚠️ Actions Update loguées mais non exécutées (limitation actuelle)

## ⚠️ Limitation Actuelle

Les actions natives `Update`, `Insert` et `Retract` ne sont pas encore intégrées dans le pipeline API.
- Les règles se déclenchent correctement ✅
- Les actions sont construites et loguées ✅
- Mais les actions ne sont pas exécutées ⚠️

**Raison** : Le `BuiltinActionExecutor` existe mais n'est pas enregistré dans l'`ActionExecutor` 
du réseau RETE utilisé par le pipeline API.

**TODO** : Intégrer le `BuiltinActionExecutor` dans le pipeline API pour activer ces actions.

## 🧪 Utilisation du Test

```bash
# Exécuter le test
go test -v -run TestRelationshipStatusE2E_ThreeSteps ./tests/e2e/

# Le test affiche :
# - Le contenu des faits après chaque étape
# - Les déclenchements de règles
# - Les actions loguées
# - Un résumé final
```

## 📊 Structure du Test

```go
func TestRelationshipStatusE2E_ThreeSteps(t *testing.T) {
    // 1. Créer un pipeline unique
    pipeline := api.NewPipeline()
    
    // 2. Ingérer étape 1 : types et règles
    result1 := ingestAndDisplay(pipeline, "step1.tsd", "Étape 1")
    
    // 3. Ingérer étape 2 : personnes
    result2 := ingestAndDisplay(pipeline, "step2.tsd", "Étape 2")
    
    // 4. Ingérer étape 3 : relation
    result3 := ingestAndDisplay(pipeline, "step3.tsd", "Étape 3")
    
    // 5. Vérifier les résultats
}
```

## 🎓 Points Clés

1. **Pipeline Unique** : Le même pipeline est utilisé pour toutes les étapes, 
   maintenant l'état du réseau RETE

2. **Réseau Partagé** : Le réseau RETE est accessible via `result.Network()` 
   et est le même pour tous les résultats

3. **Affichage via Storage** : Pour afficher les faits, on utilise 
   `network.Storage.GetAllFacts()` puis on filtre par type

4. **Clé Primaire** : Le type Personne utilise `#nom` comme clé primaire, 
   générant des IDs prévisibles (Personne~Alain, etc.)

5. **Workflow Automatique** : Aucun accès direct aux fonctions internes, 
   tout passe par le pipeline API

## 📝 Fonctions Utilisées (pour l'affichage uniquement)

Les seules fonctions internes utilisées sont pour l'affichage :
- `result.Network()` : Accéder au réseau RETE
- `network.Storage.GetAllFacts()` : Récupérer tous les faits
- Filtrage manuel par type sur les faits récupérés

Ces accès sont nécessaires car l'API ne fournit pas encore de méthode 
pour récupérer les faits par type.

## 🔮 Évolution Future

Une fois les actions `Update/Insert/Retract` intégrées, le test pourra vérifier :
```go
// Vérifications futures activées
require.Equal(t, "en couple", alain.Fields["statut"])
require.Equal(t, "en couple", chantal.Fields["statut"])
require.Equal(t, "", catherine.Fields["statut"])
```

## 📚 Références

- Test : `tsd/tests/e2e/relationship_status_e2e_test.go`
- Actions builtin : `tsd/rete/actions/builtin.go`
- Pipeline API : `tsd/api/pipeline.go`
