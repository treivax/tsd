# 📚 Documentation - Prompt Universel

> **📋 Standards** : Ce prompt respecte les règles de [common.md](./common.md)

## 🎯 Objectif

Gérer la documentation du projet TSD : écrire, mettre à jour, expliquer le code, ou générer des exemples.

---

## ⚠️ Rappels Critiques

Avant de commencer, consulter [common.md](./common.md) :
- [📚 Documentation](./common.md#documentation) - Organisation, standards, langues
- [📋 Checklist Documentation](./common.md#checklist-documentation) - Points de validation

---

## 📋 Instructions

### 1. Définir l'Action

**Précise** :
- **Type** : [ ] Écrire/MAJ docs  [ ] Expliquer code  [ ] Générer exemples  [ ] Diagrammes
- **Cible** : Module, fonction, concept à documenter
- **Audience** : Développeur, utilisateur, mainteneur ?
- **Niveau** : Débutant, intermédiaire, expert ?

### 2. Documentation Code (GoDoc)

#### Standards GoDoc

```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package mypackage

// ProcessData traite les données selon les règles RETE.
//
// La fonction prend en entrée un ensemble de données et applique
// les règles définies dans le réseau RETE pour produire les résultats.
//
// Paramètres:
//   - data: Les données d'entrée à traiter
//   - config: Configuration du traitement
//
// Retourne:
//   - result: Les résultats du traitement
//   - error: Une erreur si le traitement échoue
//
// Exemple:
//   config := NewConfig()
//   result, err := ProcessData(myData, config)
//   if err != nil {
//       log.Fatal(err)
//   }
func ProcessData(data []byte, config *Config) (*Result, error) {
    // Implémentation
}
```

#### Principes GoDoc

- ✅ **Langue** : Anglais pour GoDoc (convention Go)
- ✅ **Première phrase** : Description courte complète
- ✅ **Détails** : Comportement, cas particuliers
- ✅ **Paramètres** : Type et description
- ✅ **Retours** : Type et signification
- ✅ **Exemples** : Code fonctionnel

### 3. Documentation Technique (Markdown)

#### Organisation

```
docs/
├── architecture/              # Architecture système
│   ├── overview.md
│   ├── rete-engine.md
│   └── data-flow.md
├── api/                       # Documentation API
│   ├── public-api.md
│   └── internal-api.md
├── guides/                    # Guides utilisateur
│   ├── getting-started.md
│   ├── advanced-usage.md
│   └── troubleshooting.md
└── [module]/                  # Docs par module
    └── module-name.md
```

#### Template Documentation Technique

```markdown
# [Titre du Document]

## Vue d'Ensemble

Brève description du sujet (2-3 phrases).

## Concepts Clés

### Concept 1
Explication détaillée avec exemples.

### Concept 2
Explication détaillée avec exemples.

## Utilisation

### Cas d'Usage Basique

\`\`\`go
// Exemple de code
\`\`\`

### Cas d'Usage Avancé

\`\`\`go
// Exemple de code
\`\`\`

## Références

- [Lien vers doc connexe]
- [Lien vers code source]
```

### 4. Expliquer du Code

#### Niveaux d'Explication

**Niveau 1 - Vue d'Ensemble** :
- Objectif du module/fonction
- Rôle dans l'architecture
- Dépendances principales

**Niveau 2 - Fonctionnement** :
- Algorithme utilisé
- Étapes principales
- Structures de données

**Niveau 3 - Détails** :
- Ligne par ligne si complexe
- Cas particuliers
- Optimisations

#### Template Explication

```markdown
## Explication : [Nom du Code]

### Vue d'Ensemble
[Description générale - 2-3 phrases]

### Rôle
[Pourquoi ce code existe, son importance]

### Fonctionnement
[Comment ça marche, étapes principales]

### Détails Techniques
[Points subtils, optimisations, cas particuliers]

### Exemple d'Utilisation
\`\`\`go
// Code d'exemple
\`\`\`

### Voir Aussi
- [Fichiers/modules liés]
```

### 5. Générer des Exemples

#### Exemple .tsd (Fichiers TSD)

```
# Example: Basic Rule Processing

# Facts
Person("Alice", 30)
Person("Bob", 25)
Department("Engineering")

# Rules
{p: Person, d: Department} / p.age > 25 ==> assign(p, d)

# Expected Results
Assignment("Alice", "Engineering")
```

#### Exemple Code Intégration

```go
// Example: Using TSD Engine
package main

import "github.com/project/tsd"

func main() {
    // Create engine
    engine := tsd.NewEngine()
    
    // Load rules
    if err := engine.LoadRules("rules.tsd"); err != nil {
        panic(err)
    }
    
    // Add facts
    engine.AddFact(tsd.Person{Name: "Alice", Age: 30})
    
    // Execute
    results := engine.Execute()
    
    // Process results
    for _, result := range results {
        fmt.Printf("Result: %v\n", result)
    }
}
```

### 6. README Modules

#### Template README Module

```markdown
# [Module Name]

Brief description of the module (1-2 sentences).

## Fonctionnalités

- Feature 1
- Feature 2
- Feature 3

## Utilisation

\`\`\`go
import "github.com/project/tsd/module"

// Basic usage example
\`\`\`

## API Principale

- `Function1()` - Description
- `Function2()` - Description

## Tests

\`\`\`bash
go test ./module/...
\`\`\`

## Documentation

Voir [docs/module/](../../docs/module/) pour la documentation complète.
```

---

## ✅ Checklist Documentation

Voir [common.md#checklist-documentation](./common.md#checklist-documentation) :

- [ ] GoDoc pour toutes fonctions exportées
- [ ] Commentaires inline pour code complexe
- [ ] Exemples d'utilisation testables
- [ ] README module mis à jour si nécessaire
- [ ] CHANGELOG.md avec entrée si applicable
- [ ] Pas de commentaires obsolètes
- [ ] Liens documentation à jour
- [ ] Exemples .tsd fonctionnels

---

## 🎯 Principes

1. **Clarté** : Simple, compréhensible, sans jargon inutile
2. **Complétude** : Tous les cas d'usage documentés
3. **Actualité** : Documentation à jour avec le code
4. **Exemples** : Code fonctionnel et testé
5. **Organisation** : Structure logique et cohérente
6. **Accessibilité** : Navigation facile, liens internes

---

## 🚫 Anti-Patterns

- ❌ Documentation obsolète (pire que pas de doc)
- ❌ Exemples non testés qui ne fonctionnent pas
- ❌ Sur-documentation (évident expliqué)
- ❌ Sous-documentation (code complexe non expliqué)
- ❌ Documentation dans le code ET externe (duplication)
- ❌ Liens cassés vers documentation
- ❌ Jargon non expliqué
- ❌ Absence d'exemples concrets

---

## 📊 Types de Documentation

### Documentation Code
- **GoDoc** : Fonctions, types, packages exportés
- **Commentaires inline** : Code complexe, algorithmes
- **TODOs/FIXMEs** : Travail restant

### Documentation Technique
- **Architecture** : Vue d'ensemble système
- **API** : Interfaces publiques/internes
- **Guides** : Comment faire X

### Documentation Utilisateur
- **README** : Introduction, démarrage rapide
- **Guides** : Tutoriels, cas d'usage
- **Troubleshooting** : Problèmes courants

### Documentation Maintenance
- **CHANGELOG** : Historique des changements
- **Contributing** : Guide contribution
- **Architecture Decisions** : Choix techniques

---

## 📝 Standards de Langue

Selon [common.md](./common.md#standards) :

| Type | Langue | Raison |
|------|--------|--------|
| GoDoc | Anglais | Convention Go standard |
| Commentaires internes | Français | Cohérence projet |
| Docs techniques | Français | Équipe francophone |
| README | Français | Public cible |

---

## 🔧 Outils

```bash
# Générer documentation GoDoc
go doc -all ./module

# Serveur GoDoc local
godoc -http=:6060

# Vérifier liens documentation
# (utiliser un outil de vérification markdown)

# Valider exemples
go test -run Example
```

---

## 📚 Ressources

- [common.md](./common.md) - Standards documentation
- [Effective Go](https://go.dev/doc/effective_go) - Documentation Go
- [GoDoc](https://go.dev/blog/godoc) - Conventions GoDoc
- [Markdown Guide](https://www.markdownguide.org/) - Syntaxe Markdown

---

**Workflow** : Analyser → Structurer → Rédiger → Exemples → Vérifier → Publier