# 📝 Prompts Réutilisables - Projet TSD

Ce répertoire contient des prompts réutilisables pour faciliter les interactions avec l'assistant IA lors du développement du projet TSD (Type System with Dependencies - Moteur de règles RETE).

## 📚 Prompts Disponibles

| Catégorie | Prompt | Description | Statut |
|-----------|--------|-------------|--------|
| 🧪 **Tests** | [`run-tests.md`](run-tests.md) | Lancer l'ensemble des tests (unitaires, intégration, runner universel) | ✅ |
| 🧪 **Tests** | [`add-test.md`](add-test.md) | Ajouter des tests (unitaires, intégration, RETE) | ✅ |
| 🧪 **Tests** | [`debug-test.md`](debug-test.md) | Débugger un test qui échoue avec analyse détaillée | ✅ |
| 🔧 **Dev** | [`add-feature.md`](add-feature.md) | Ajouter une nouvelle fonctionnalité au projet | ✅ |
| 🔧 **Dev** | [`modify-behavior.md`](modify-behavior.md) | Modifier un comportement ou une fonctionnalité existante | ✅ |
| 🔧 **Dev** | [`fix-bug.md`](fix-bug.md) | Corriger un bug identifié avec méthodologie complète | ✅ |
| 🔧 **Dev** | [`refactor.md`](refactor.md) | Refactoriser du code sans changer le comportement | ✅ |
| 🔧 **Dev** | [`deep-clean.md`](deep-clean.md) | Nettoyage approfondi du code (fichiers, code mort, refactoring, tests, doc) | ✅ |
| 🐛 **Debug** | [`analyze-error.md`](analyze-error.md) | Analyser une erreur ou un problème avec diagnostic | ✅ |
| 🐛 **Debug** | [`investigate.md`](investigate.md) | Investiguer un comportement étrange sans erreur explicite | ✅ |
| ⚡ **Performance** | [`optimize-performance.md`](optimize-performance.md) | Optimiser les performances avec profiling et benchmarks | ✅ |
| 👀 **Review** | [`code-review.md`](code-review.md) | Faire une revue de code complète et structurée | ✅ |
| 📖 **Docs** | [`explain-code.md`](explain-code.md) | Expliquer une partie du code en détail | ✅ |
| 📖 **Docs** | [`update-docs.md`](update-docs.md) | Mettre à jour la documentation (README, docs/, GoDoc, CHANGELOG) | ✅ |
| 📖 **Docs** | [`generate-examples.md`](generate-examples.md) | Générer des exemples RETE (.constraint, .facts, documentation) | ✅ |
| ✓ **Validation** | [`validate-network.md`](validate-network.md) | Valider un réseau RETE (structure, propagation, résultats) | ✅ |
| 🔄 **Migration** | [`migrate.md`](migrate.md) | Migrer version Go, dépendances, ou adapter à changements d'API | ✅ |
| 📊 **Analyse** | [`stats-code.md`](stats-code.md) | Générer statistiques du code (lignes, complexité, fichiers volumineux) | ✅ |

## ⚠️ RÈGLES STRICTES APPLIQUÉES

Tous les prompts suivent des règles strictes pour garantir la qualité du code :

### 🚫 Pour le Code Golang

**INTERDICTIONS ABSOLUES** :
- ❌ **AUCUN HARDCODING** : Pas de valeurs en dur, magic numbers, ou configurations hardcodées
- ❌ **CODE SPÉCIFIQUE** : Pas de code limité à un cas d'usage particulier
- ✅ **CODE GÉNÉRIQUE** : Toujours utiliser paramètres, constantes nommées, interfaces
- ✅ **BONNES PRATIQUES GO** : Respect strict de Effective Go, conventions, gestion d'erreurs

**Prompts concernés** : `add-feature.md`, `modify-behavior.md`, `debug-test.md`, `analyze-error.md`, `deep-clean.md`

### 🚫 Pour les Tests RETE

**INTERDICTIONS ABSOLUES** :
- ❌ **AUCUNE SIMULATION** : Pas de résultats hardcodés ou calculés manuellement
- ❌ **AUCUN MOCK RÉSEAU** : Pas de simulation du réseau RETE
- ✅ **EXTRACTION OBLIGATOIRE** : Toujours extraire résultats depuis le réseau RETE réel
- ✅ **VALIDATION RÉELLE** : Interroger TerminalNodes, inspecter mémoires (Left/Right/Result)

**Prompts concernés** : `debug-test.md`, `validate-network.md`, `analyze-error.md`, `modify-behavior.md`, `deep-clean.md`

### 📋 Exemples

❌ **MAUVAIS** :
```go
// Hardcoding interdit !
timeout := 30
expectedTokens := 5  // Simulé manuellement
```

✅ **BON** :
```go
// Code générique avec constantes
const DefaultTimeout = 30 * time.Second

// Extraction depuis le réseau réel
actualTokens := 0
for _, terminal := range network.TerminalNodes {
    actualTokens += len(terminal.Memory.GetTokens())
}
```

### 🚧 Prompts Potentiels Futurs

| Catégorie | Prompt | Description |
|-----------|--------|-------------|
| 🏗️ Architecture | `design-decision.md` | Documenter une décision architecturale |
| 👀 Review | `security-audit.md` | Audit de sécurité complet |

## 🚀 Utilisation

### Méthode 1 : Copier-Coller
1. Ouvrir le fichier prompt souhaité (ex: `run-tests.md`)
2. Copier le contenu
3. Coller dans votre conversation avec l'assistant

### Méthode 2 : Référence
Simplement dire à l'assistant :
```
Utilise le prompt "run-tests"
```
ou
```
Applique le prompt dans .github/prompts/run-tests.md
```

### Méthode 3 : Adaptation
Copier le prompt et modifier les paramètres selon vos besoins :
```
[Prompt run-tests.md]
Mais seulement pour les tests du module rete/
```

## 📋 Convention de Nommage

- **Nom du fichier** : `action-cible.md` (kebab-case)
- **Titre** : Description claire de l'action
- **Sections** :
  - `## Contexte` - Informations sur le projet
  - `## Objectif` - Ce que tu veux accomplir
  - `## Instructions` - Étapes à suivre
  - `## Critères de Succès` - Comment valider le résultat

## 🎯 Bonnes Pratiques

1. **Sois spécifique** : Plus le prompt est précis, meilleure sera la réponse
2. **Fournis le contexte** : Indique le module, le fichier ou la fonction concernée
3. **Définis le succès** : Explique ce que tu attends comme résultat
4. **Itère** : N'hésite pas à ajuster le prompt si nécessaire

## 🆕 Ajouter un Nouveau Prompt

Pour ajouter un nouveau prompt :

1. Créer un fichier `.md` dans ce répertoire
2. Suivre la structure standard :
```markdown
# [Titre du Prompt]

## Contexte
[Description du contexte du projet]

## Objectif
[Ce que tu veux accomplir]

## Instructions
[Étapes détaillées]

## Critères de Succès
[Comment vérifier que c'est réussi]

## Exemple
[Exemple concret d'utilisation]
```
3. Mettre à jour ce README avec le nouveau prompt

## 🔗 Liens Utiles

- [Makefile du projet](../../Makefile) - Commandes disponibles
- [Architecture RETE](../../docs/) - Documentation technique
- [Tests](../../test/) - Répertoire des tests

## 🎓 Garanties de Qualité

Chaque prompt garantit :
- ✅ **Code sans hardcoding** : Valeurs en constantes, paramètres configurables
- ✅ **Code générique** : Réutilisable, extensible, interfaces appropriées
- ✅ **Tests RETE authentiques** : Extraction réelle du réseau, pas de simulation
- ✅ **Bonnes pratiques Go** : go vet, golangci-lint, Effective Go
- ✅ **Documentation** : GoDoc, commentaires explicites, exemples

## 💡 Exemples d'Utilisation

### Exemple 1 : Lancer tous les tests ✅
```
Relance moi l'ensemble des tests, dont le runner universel
```
ou simplement :
```
Utilise le prompt "run-tests"
```

### Exemple 2 : Débugger un test spécifique ✅
```
Le test TestIncrementalPropagation échoue avec l'erreur "variable non liée: p". 
Utilise le prompt "debug-test" pour identifier et corriger le problème.
```

### Exemple 3 : Ajouter une fonctionnalité ✅
```
Je veux ajouter le support des opérateurs de comparaison de chaînes 
(startsWith, endsWith, contains) dans les AlphaNodes. 
Utilise le prompt "add-feature".
```

### Exemple 4 : Valider un réseau RETE ✅
```
J'ai créé un nouveau réseau dans beta_coverage_tests/join_complex.constraint.
Utilise le prompt "validate-network" pour vérifier qu'il fonctionne correctement.
```

### Exemple 5 : Code Review ✅
```
Je viens de terminer la correction du bug de propagation incrémentale 
dans rete/node_join.go. Utilise le prompt "code-review".
```

### Exemple 6 : Expliquer du Code ✅
```
Je ne comprends pas comment fonctionne evaluateJoinConditions dans 
rete/node_join.go. Utilise le prompt "explain-code" niveau intermédiaire.
```

### Exemple 7 : Modifier un Comportement ✅
```
Je veux modifier le comportement de evaluateJoinConditions pour gérer
intelligemment les variables manquantes au lieu de générer une erreur.
Utilise le prompt "modify-behavior".
```

### Exemple 8 : Nettoyage Approfondi ✅
```
Le projet a accumulé du code au fil du temps. Je veux faire un grand
nettoyage : fichiers inutilisés, code mort, duplication, tests obsolètes.
Utilise le prompt "deep-clean".
```

### Exemple 9 : Refactoriser du Code ✅
```
La fonction evaluateJoinConditions fait 150 lignes avec une complexité
cyclomatique de 25. Je veux la refactoriser pour améliorer sa lisibilité
sans changer son comportement. Utilise le prompt "refactor".
```

### Exemple 10 : Investiguer un Comportement ✅
```
Les tokens se propagent bizarrement dans certains cas mais il n'y a pas
d'erreur. Le comportement change selon l'ordre des faits dans le fichier.
Utilise le prompt "investigate".
```

### Exemple 11 : Mettre à Jour Documentation ✅
```
Suite à l'ajout des opérateurs de chaînes, je veux mettre à jour toute
la documentation : README, CHANGELOG, docs/, GoDoc, et créer des exemples.
Utilise le prompt "update-docs".
```

### Exemple 12 : Générer des Exemples RETE ✅
```
Je veux créer un exemple complet pour démontrer les opérateurs de chaînes
(startsWith, endsWith, contains) avec fichiers .constraint, .facts, et
documentation. Utilise le prompt "generate-examples".
Utilise le prompt "generate-examples".
```

### Exemple 13 : Statistiques du Code ✅
```
Je veux connaître les statistiques complètes du code du projet :
combien de lignes de code Go fonctionnel, quels sont les fichiers
et fonctions les plus volumineux, et identifier le code nécessitant
du refactoring. Utilise le prompt "stats-code".
```

## 📊 Statistiques

- **Prompts disponibles** : 18
- **Prompts potentiels futurs** : 2
- **Catégories** : 8 (Tests, Dev, Debug, Performance, Review, Docs, Validation, Migration, Analyse)
- **Taux de complétion** : 85%

## 🎓 Guide Rapide

### Pour les Débutants
1. Commencez par `explain-code.md` pour comprendre le projet
2. Utilisez `run-tests.md` pour valider votre environnement
3. Explorez avec `validate-network.md` pour voir comment RETE fonctionne

### Pour les Développeurs
1. `add-feature.md` pour ajouter des fonctionnalités
2. `modify-behavior.md` pour modifier une fonctionnalité existante
3. `fix-bug.md` pour corriger un bug identifié
4. `refactor.md` pour améliorer la structure du code
5. `add-test.md` pour ajouter des tests
6. `debug-test.md` quand un test échoue
7. `code-review.md` avant de merger
8. `deep-clean.md` → Nettoyer régulièrement

### Pour le Debugging
1. `analyze-error.md` pour comprendre les erreurs
2. `debug-test.md` pour les problèmes de tests
3. `fix-bug.md` pour corriger un bug identifié
4. `investigate.md` pour comportements étranges sans erreur
5. `validate-network.md` pour les problèmes RETE

### Pour l'Optimisation
1. `optimize-performance.md` pour améliorer les performances
2. `investigate.md` pour profiler et analyser
3. `refactor.md` pour améliorer la structure

### Pour la Documentation
1. `update-docs.md` pour maintenir la documentation à jour
2. `generate-examples.md` pour créer des exemples RETE
3. `explain-code.md` pour documenter et expliquer

### Pour la Maintenance
1. `migrate.md` pour migrer versions/dépendances
2. `deep-clean.md` pour nettoyer le projet
3. `refactor.md` pour améliorer la qualité du code
4. `update-docs.md` pour maintenir la documentation

### Pour l'Analyse
1. `stats-code.md` pour analyser les statistiques du code
2. `code-review.md` pour reviewer la qualité
3. `optimize-performance.md` pour identifier optimisations
4. `investigate.md` pour explorer comportements

## 🤝 Contribution

Pour ajouter un nouveau prompt :
1. Créer le fichier `.md` dans ce dossier
2. Suivre la structure standard (voir section "Ajouter un Nouveau Prompt")
3. **Inclure les règles strictes** :
   - Si code Golang : Ajouter section "AUCUN HARDCODING" + "CODE GÉNÉRIQUE"
   - Si tests RETE : Ajouter section "AUCUNE SIMULATION" + "EXTRACTION RÉSEAU RÉEL"
4. Mettre à jour ce README avec le nouveau prompt
5. Mettre à jour INDEX.md avec le nouveau prompt
6. Tester le prompt avec l'assistant IA
7. Documenter des exemples d'utilisation

---

**Dernière mise à jour** : Novembre 2025
**Mainteneur** : Équipe TSD  
**Version** : 2.0 (18 prompts disponibles)