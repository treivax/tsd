# Chargement Incrémental Multi-Fichiers

## Vue d'Ensemble

Le système TSD supporte le **chargement incrémental de programmes répartis sur plusieurs fichiers**. Cette fonctionnalité permet d'organiser vos programmes TSD de manière modulaire en séparant :

- Les définitions de types (schémas)
- Les règles métier
- Les données (faits)
- Les déclarations de xuple-spaces

Chaque fichier peut être chargé séquentiellement, et le moteur RETE fusionne automatiquement les définitions pour créer un réseau cohérent.

## Concepts Clés

### Validation Incrémentale

Lors du chargement d'un fichier, le système :

1. **Analyse** le contenu du fichier courant
2. **Fusionne** les types avec ceux déjà chargés dans le réseau
3. **Valide** la cohérence avec le contexte existant
4. **Étend** le réseau RETE sans perdre l'état précédent

Les types définis dans les fichiers précédents sont **automatiquement disponibles** pour les fichiers suivants.

### Réseau RETE Partagé

Le réseau RETE est maintenu entre les chargements de fichiers :

- Les **TypeNodes** sont préservés et réutilisés
- Les **faits** restent en mémoire
- Les **règles** s'ajoutent de manière incrémentale
- Le **contexte de validation** est propagé

## Cas d'Usage

### 1. Séparation Schéma / Données

**Pattern le plus courant** : définir les types dans un fichier, charger les données dans un autre.

#### Fichier 1 : `schema.tsd`

```tsd
// Définitions des types avec clés primaires
type Person(#id: string, name: string, age: number, email: string)
type Department(#code: string, name: string, budget: number)
type Assignment(#person_id: string, #dept_code: string, role: string)
```

#### Fichier 2 : `data.tsd`

```tsd
// Faits utilisant les types définis précédemment
Person(id: "P001", name: "Alice Dupont", age: 30, email: "alice@example.com")
Person(id: "P002", name: "Bob Martin", age: 25, email: "bob@example.com")

Department(code: "ENG", name: "Engineering", budget: 500000)
Department(code: "HR", name: "Human Resources", budget: 200000)

Assignment(person_id: "P001", dept_code: "ENG", role: "Senior Developer")
Assignment(person_id: "P002", dept_code: "ENG", role: "Junior Developer")
```

#### Code Go

```go
package main

import (
    "log"
    "github.com/treivax/tsd/api"
)

func main() {
    pipeline := api.NewPipeline()
    
    // Charger le schéma
    _, err := pipeline.IngestFile("schema.tsd")
    if err != nil {
        log.Fatalf("Erreur chargement schéma: %v", err)
    }
    
    // Charger les données (les types sont automatiquement disponibles)
    result, err := pipeline.IngestFile("data.tsd")
    if err != nil {
        log.Fatalf("Erreur chargement données: %v", err)
    }
    
    log.Printf("Chargement réussi: %d faits soumis", result.FactsSubmitted)
}
```

### 2. Organisation Modulaire par Domaine

Séparer les types et règles par domaine métier.

#### Structure de Fichiers

```
project/
├── schemas/
│   ├── customers.tsd      # Types clients
│   ├── products.tsd       # Types produits
│   └── orders.tsd         # Types commandes
├── rules/
│   ├── pricing.tsd        # Règles de tarification
│   ├── inventory.tsd      # Règles de stock
│   └── promotions.tsd     # Règles promotionnelles
└── data/
    ├── customers.tsd      # Données clients
    ├── products.tsd       # Données produits
    └── orders.tsd         # Données commandes
```

#### Code de Chargement

```go
func loadModularProject() error {
    pipeline := api.NewPipeline()
    
    // Ordre de chargement recommandé
    files := []string{
        // 1. Schémas
        "schemas/customers.tsd",
        "schemas/products.tsd",
        "schemas/orders.tsd",
        // 2. Règles
        "rules/pricing.tsd",
        "rules/inventory.tsd",
        "rules/promotions.tsd",
        // 3. Données
        "data/customers.tsd",
        "data/products.tsd",
        "data/orders.tsd",
    }
    
    for _, file := range files {
        _, err := pipeline.IngestFile(file)
        if err != nil {
            return fmt.Errorf("échec chargement %s: %w", file, err)
        }
        log.Printf("✓ Chargé: %s", file)
    }
    
    return nil
}
```

### 3. Évolution Incrémentale

Ajouter progressivement des types et règles à un système existant.

#### Étape 1 : Base Initiale

```tsd
// base.tsd
type User(#id: string, username: string)

rule user_validation : {u: User} / len(u.username) > 3 ==> validate(u)
```

#### Étape 2 : Extension

```tsd
// extensions.tsd
// Nouveau type (s'ajoute aux types existants)
type UserProfile(#user_id: string, bio: string, avatar: string)

// Nouvelle règle (s'ajoute aux règles existantes)
rule profile_completeness : {p: UserProfile} / len(p.bio) > 50 ==> award_badge(p)
```

#### Code

```go
pipeline := api.NewPipeline()

// Charger la base
pipeline.IngestFile("base.tsd")

// Ajouter dynamiquement des extensions
pipeline.IngestFile("extensions.tsd")

// Le réseau contient maintenant :
// - Types: User, UserProfile
// - Règles: user_validation, profile_completeness
```

### 4. Configuration par Environnement

Charger différentes données selon l'environnement.

```go
func loadEnvironment(env string) error {
    pipeline := api.NewPipeline()
    
    // Schéma commun
    pipeline.IngestFile("schema.tsd")
    
    // Règles communes
    pipeline.IngestFile("rules.tsd")
    
    // Données spécifiques à l'environnement
    switch env {
    case "dev":
        pipeline.IngestFile("data/dev/test-users.tsd")
    case "staging":
        pipeline.IngestFile("data/staging/sample-data.tsd")
    case "prod":
        pipeline.IngestFile("data/prod/initial-data.tsd")
    }
    
    return nil
}
```

## Bonnes Pratiques

### ✅ Recommandations

1. **Ordre de Chargement**
   - Charger les **types** en premier
   - Puis les **actions** et **xuple-spaces**
   - Ensuite les **règles**
   - Enfin les **faits**

2. **Nommage de Fichiers**
   - Utiliser des noms descriptifs : `customer-types.tsd`, `pricing-rules.tsd`
   - Préfixer par numéro pour forcer l'ordre : `01-types.tsd`, `02-rules.tsd`
   - Grouper par domaine dans des dossiers

3. **Modularité**
   - Un fichier = un domaine ou une responsabilité
   - Éviter les fichiers trop volumineux (>500 lignes)
   - Séparer les types stables des données volatiles

4. **Clés Primaires**
   - Toujours définir les clés primaires dans les types (`#field`)
   - Les clés primaires sont **préservées** lors de la fusion incrémentale
   - Essentiel pour la cohérence multi-fichiers

5. **Gestion d'Erreurs**
   - Vérifier les erreurs après chaque `IngestFile()`
   - En cas d'erreur, le réseau reste dans l'état précédent (rollback automatique)
   - Logger le fichier qui pose problème pour faciliter le debug

### ⚠️ Pièges à Éviter

1. **Duplication de Types**
   - ❌ Ne pas redéfinir un type déjà chargé
   - Le système détecte et ignore les doublons, mais c'est inefficace
   - Garder les définitions dans un seul fichier de schéma

2. **Ordre Incorrect**
   - ❌ Ne pas charger les faits avant les types
   - Résultat : erreur "type X non défini"
   - Toujours charger les dépendances avant

3. **Références Cassées**
   - ❌ Faits référençant des types non chargés
   - Vérifier que tous les types nécessaires sont disponibles

4. **Fichiers Trop Couplés**
   - ❌ Fichiers interdépendants difficiles à charger séparément
   - Concevoir des modules indépendants quand possible

## Exemples Complets

### Exemple : Système de Gestion d'Événements

#### 1. Types (`schemas/events.tsd`)

```tsd
type Event(#id: string, name: string, date: string, capacity: number)
type Attendee(#email: string, name: string, company: string)
type Registration(#event_id: string, #attendee_email: string, status: string)
```

#### 2. Règles (`rules/events.tsd`)

```tsd
// Règle : événement complet
rule event_full : {e: Event, r: Registration} 
    / count(r, e.id) >= e.capacity 
    ==> mark_full(e)

// Règle : rappel de confirmation
rule confirmation_reminder : {r: Registration} 
    / r.status == "pending" 
    ==> send_reminder(r)
```

#### 3. Données (`data/conference-2025.tsd`)

```tsd
// Événements
Event(id: "E001", name: "TSD Conference 2025", date: "2025-06-15", capacity: 100)
Event(id: "E002", name: "Workshop RETE", date: "2025-06-16", capacity: 30)

// Participants
Attendee(email: "alice@tech.com", name: "Alice Dupont", company: "TechCorp")
Attendee(email: "bob@dev.io", name: "Bob Martin", company: "DevStudio")

// Inscriptions
Registration(event_id: "E001", attendee_email: "alice@tech.com", status: "confirmed")
Registration(event_id: "E002", attendee_email: "bob@dev.io", status: "pending")
```

#### 4. Code de Chargement

```go
package main

import (
    "log"
    "github.com/treivax/tsd/api"
)

func main() {
    pipeline := api.NewPipeline()
    
    files := []string{
        "schemas/events.tsd",
        "rules/events.tsd",
        "data/conference-2025.tsd",
    }
    
    for _, file := range files {
        result, err := pipeline.IngestFile(file)
        if err != nil {
            log.Fatalf("❌ Erreur %s: %v", file, err)
        }
        log.Printf("✅ %s: %d types, %d règles, %d faits", 
            file, 
            len(result.GetTypes()),
            len(result.GetRules()),
            result.FactsSubmitted,
        )
    }
    
    // Accéder aux xuple-spaces générés par les règles
    spaces := pipeline.GetXupleSpaces()
    for name, space := range spaces {
        log.Printf("📦 Xuple-space '%s': %d xuples", name, len(space.GetAll()))
    }
}
```

## Détails Techniques

### Mécanisme de Fusion

Lors de l'appel à `pipeline.IngestFile(file)` avec un réseau existant :

1. **Parsing** : Le fichier est parsé en AST
2. **Extraction** : Les types, règles, faits sont extraits
3. **Enrichissement** : Les types du réseau existant sont fusionnés avec ceux du fichier
4. **Validation** : Le programme enrichi est validé (cohérence des types)
5. **Soumission** : Les nouveaux éléments sont ajoutés au réseau
6. **Commit** : La transaction est validée

### Gestion des Types

```go
// Pseudo-code interne
func enrichProgramWithNetworkTypes(program, network) {
    enrichedProgram = copy(program)
    
    // Construire la map des types existants
    existingTypes = map[typeName]bool
    for type in program.Types {
        existingTypes[type.Name] = true
    }
    
    // Ajouter les types du réseau non présents dans le programme
    for networkType in network.Types {
        if !existingTypes[networkType.Name] {
            enrichedProgram.Types.append(networkType)
        }
    }
    
    return enrichedProgram
}
```

Cette fusion garantit que les faits d'un fichier peuvent référencer des types définis dans des fichiers précédents.

### Préservation des Clés Primaires

Les champs marqués avec `#` (clés primaires) sont **systématiquement préservés** lors de la fusion :

```tsd
// schema.tsd
type User(#id: string, name: string)

// data.tsd (chargé après)
User(id: "U001", name: "Alice")  // ✅ La clé primaire 'id' est reconnue
```

Le système sait que `id` est la clé primaire car cette information a été sauvegardée dans `network.Types` lors du chargement de `schema.tsd`.

## Dépannage

### Erreur : "type X non défini"

**Cause** : Le type référencé n'a pas été chargé avant le fichier courant.

**Solution** :
```go
// ❌ Mauvais ordre
pipeline.IngestFile("data.tsd")    // Référence Person
pipeline.IngestFile("schema.tsd")  // Définit Person

// ✅ Bon ordre
pipeline.IngestFile("schema.tsd")  // Définit Person
pipeline.IngestFile("data.tsd")    // Référence Person
```

### Erreur : "Primary key mismatch"

**Cause** : Tentative de redéfinir un type avec une clé primaire différente.

**Solution** : Garder une seule définition canonique du type dans un fichier de schéma.

### Performances Dégradées

**Cause** : Chargement de trop nombreux fichiers séquentiellement.

**Optimisation** :
- Regrouper les fichiers petits
- Charger les données en batch
- Utiliser des transactions explicites si disponible

## Références

- [Architecture RETE](../architecture/rete-engine.md)
- [Validation Incrémentale](../architecture/incremental-validation.md)
- [API Pipeline](../api/pipeline.md)
- [Exemples Multi-Fichiers](../../examples/multi-file/)

## Changelog

- **2025-12-21** : Ajout du support multi-fichiers avec fusion automatique des types
- **2025-12-21** : Documentation initiale du pattern