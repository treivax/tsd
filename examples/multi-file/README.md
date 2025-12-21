# Exemples de Chargement Multi-Fichiers

Ce dossier contient des exemples concrets d'utilisation du **chargement incrémental multi-fichiers** avec TSD.

## Vue d'Ensemble

Le système TSD permet de répartir vos programmes sur plusieurs fichiers et de les charger de manière incrémentale. Les types, règles et faits définis dans les fichiers précédents sont automatiquement disponibles pour les fichiers suivants.

## Exemples Disponibles

### 1. Basic - Séparation Schéma/Données

Le pattern le plus simple et le plus courant : types dans un fichier, données dans un autre.

- **Fichiers** :
  - `01-basic/schema.tsd` - Définitions de types
  - `01-basic/data.tsd` - Faits utilisant ces types
  - `01-basic/main.go` - Code de chargement

- **Démo** :
  ```bash
  cd 01-basic
  go run main.go
  ```

### 2. Modular - Organisation par Domaine

Organisation modulaire avec plusieurs domaines métier.

- **Fichiers** :
  - `02-modular/schemas/` - Types par domaine
  - `02-modular/rules/` - Règles métier
  - `02-modular/data/` - Données
  - `02-modular/main.go` - Chargement orchestré

- **Démo** :
  ```bash
  cd 02-modular
  go run main.go
  ```

### 3. Events - Système Complet

Exemple réaliste d'un système de gestion d'événements.

- **Fichiers** :
  - `03-events/schemas/events.tsd` - Types événements
  - `03-events/rules/events.tsd` - Règles métier
  - `03-events/data/conference-2025.tsd` - Données événement
  - `03-events/main.go` - Application complète

- **Démo** :
  ```bash
  cd 03-events
  go run main.go
  ```

### 4. Incremental - Extension Progressive

Démontre l'ajout progressif de types et règles.

- **Fichiers** :
  - `04-incremental/base.tsd` - Système de base
  - `04-incremental/extension-1.tsd` - Première extension
  - `04-incremental/extension-2.tsd` - Deuxième extension
  - `04-incremental/main.go` - Chargement progressif

- **Démo** :
  ```bash
  cd 04-incremental
  go run main.go
  ```

## Guide de Démarrage Rapide

### Exemple Minimal

1. **Créer un schéma** (`schema.tsd`) :
   ```tsd
   type Person(#id: string, name: string, age: number)
   ```

2. **Créer des données** (`data.tsd`) :
   ```tsd
   Person(id: "P001", name: "Alice", age: 30)
   Person(id: "P002", name: "Bob", age: 25)
   ```

3. **Charger les fichiers** (`main.go`) :
   ```go
   package main

   import (
       "log"
       "github.com/treivax/tsd/api"
   )

   func main() {
       pipeline := api.NewPipeline()
       
       // Charger le schéma en premier
       _, err := pipeline.IngestFile("schema.tsd")
       if err != nil {
           log.Fatal(err)
       }
       
       // Charger les données (les types sont disponibles)
       result, err := pipeline.IngestFile("data.tsd")
       if err != nil {
           log.Fatal(err)
       }
       
       log.Printf("✅ Succès: %d faits chargés", result.FactsSubmitted)
   }
   ```

## Bonnes Pratiques

### Ordre de Chargement Recommandé

```
1. Types (schémas)
2. Actions et Xuple-Spaces
3. Règles métier
4. Faits (données)
```

### Nommage des Fichiers

- Utiliser des préfixes numériques pour forcer l'ordre : `01-types.tsd`, `02-rules.tsd`
- Noms descriptifs : `customer-types.tsd`, `pricing-rules.tsd`
- Grouper par domaine dans des dossiers

### Structure de Projet Recommandée

```
project/
├── schemas/
│   ├── core.tsd          # Types de base
│   ├── customers.tsd     # Types clients
│   └── products.tsd      # Types produits
├── rules/
│   ├── validation.tsd    # Règles de validation
│   └── business.tsd      # Règles métier
└── data/
    ├── dev/             # Données de dev
    ├── staging/         # Données de staging
    └── prod/            # Données de production
```

## Points Importants

### ✅ Ce qui Fonctionne

- ✅ Types définis dans un fichier disponibles dans les suivants
- ✅ Clés primaires (`#field`) préservées automatiquement
- ✅ Fusion intelligente sans duplication
- ✅ Rollback automatique en cas d'erreur
- ✅ Réseau RETE maintenu entre les chargements

### ⚠️ Points d'Attention

- ⚠️ Charger les types **avant** les faits qui les référencent
- ⚠️ Ne pas redéfinir un type déjà chargé (détecté mais inefficace)
- ⚠️ Vérifier les erreurs après chaque `IngestFile()`

## Exécuter Tous les Exemples

```bash
# Depuis ce dossier
for dir in */; do
    echo "=== Exécution de $dir ==="
    (cd "$dir" && go run main.go)
    echo ""
done
```

## Dépannage

### Erreur : "type X non défini"

**Problème** : Le type est référencé avant d'être défini.

**Solution** : Charger le fichier de schéma avant le fichier de données.

```go
// ❌ Mauvais ordre
pipeline.IngestFile("data.tsd")    // Erreur: Person non défini
pipeline.IngestFile("schema.tsd")

// ✅ Bon ordre
pipeline.IngestFile("schema.tsd")
pipeline.IngestFile("data.tsd")
```

### Performance Dégradée

**Problème** : Trop de petits fichiers chargés séquentiellement.

**Solution** : Regrouper les fichiers par domaine ou fonction.

## Documentation Complète

Pour plus de détails, consultez :

- [Guide Complet Multi-Fichiers](../../docs/user-guide/multi-file-loading.md)
- [Architecture RETE](../../docs/architecture/rete-engine.md)
- [API Pipeline](../../docs/api/pipeline.md)

## Contribuer

Pour ajouter un nouvel exemple :

1. Créer un dossier `XX-nom/` (XX = numéro séquentiel)
2. Ajouter les fichiers `.tsd` et `main.go`
3. Documenter dans ce README
4. Tester avec `go run main.go`

## Licence

Copyright (c) 2025 TSD Contributors  
Licensed under the MIT License