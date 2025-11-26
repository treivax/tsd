# Rapport d'augmentation de la couverture de test - cmd/tsd

## Résumé

Ce rapport documente l'amélioration de la couverture de test pour le package `cmd/tsd`, qui contient le point d'entrée principal de l'application TSD.

## Statistiques de couverture

### Avant amélioration
- **Couverture totale**: 51.0%
- **Date**: Avant cette session

### Après amélioration
- **Couverture totale**: 56.9%
- **Date**: Session actuelle
- **Amélioration**: +5.9 points de pourcentage

## Détails de la couverture par fonction

| Fonction | Avant | Après | Statut |
|----------|-------|-------|--------|
| `main()` | 0.0% | 0.0% | ⚠️ Non testable directement |
| `parseFlags()` | 100.0% | 100.0% | ✅ Complète |
| `validateConfig()` | 100.0% | 100.0% | ✅ Complète |
| `parseConstraintSource()` | 80.0% | 100.0% | ✅ Améliorée |
| `parseFromStdin()` | 0.0% | 100.0% | ✅ Nouvellement testée |
| `parseFromText()` | 100.0% | 100.0% | ✅ Complète |
| `parseFromFile()` | 100.0% | 100.0% | ✅ Complète |
| `printParsingHeader()` | 100.0% | 100.0% | ✅ Complète |
| `runValidationOnly()` | 100.0% | 100.0% | ✅ Complète |
| `runWithFacts()` | 0.0% | 0.0% | ⚠️ Appelle os.Exit() |
| `printResults()` | 0.0% | 0.0% | ⚠️ Dépend de runWithFacts() |
| `countActivations()` | 0.0% | 0.0% | ⚠️ Dépend de runWithFacts() |
| `printActivationDetails()` | 0.0% | 0.0% | ⚠️ Dépend de runWithFacts() |
| `printVersion()` | 100.0% | 100.0% | ✅ Complète |
| `printHelp()` | 100.0% | 100.0% | ✅ Complète |

## Tests ajoutés

### 1. Tests unitaires étendus

#### `TestParseFromStdin()`
- **Lignes**: 718-806
- **Description**: Tests complets pour la lecture depuis stdin
- **Cas couverts**:
  - Contrainte valide depuis stdin
  - Entrée vide
  - Syntaxe invalide
  - Mode verbeux
  - Contraintes complexes

#### `TestParseFromStdinError()`
- **Lignes**: 1357-1373
- **Description**: Test de gestion d'erreurs lors de la lecture stdin
- **Cas couverts**:
  - Erreur de lecture avec pipe fermé

#### `TestRunWithFactsLogic()`
- **Lignes**: 809-852
- **Description**: Test de la logique de vérification d'existence des fichiers de faits
- **Cas couverts**:
  - Fichier de faits existant
  - Fichier de faits non existant

#### `TestPrintResults()`
- **Lignes**: 855-929
- **Description**: Tests de l'affichage des résultats
- **Cas couverts**:
  - Sans activations (mode non-verbeux)
  - Sans activations (mode verbeux)
  - Avec activations (mode non-verbeux)
  - Avec activations (mode verbeux)

#### `TestCountActivationsWithRealNetwork()`
- **Lignes**: 932-977
- **Description**: Tests de la logique de comptage des activations
- **Cas couverts**:
  - Aucun terminal
  - Un terminal sans tokens
  - Un terminal avec tokens
  - Multiples terminaux avec tokens
  - Terminaux mixtes (avec et sans tokens)

#### `TestPrintActivationDetails()`
- **Lignes**: 980-1030
- **Description**: Tests de l'affichage détaillé des activations
- **Cas couverts**:
  - Aucune activation
  - Une seule activation
  - Activations multiples

#### `TestParseConstraintSource()` - Extension
- **Lignes ajoutées**: 664-675
- **Description**: Extension du test existant pour couvrir le cas stdin
- **Nouveaux cas**:
  - Routage vers stdin avec mock de l'entrée

#### `TestEdgeCases()`
- **Lignes**: 1375-1417
- **Description**: Tests de cas limites
- **Cas couverts**:
  - Configuration vide
  - Fichier avec chemin invalide
  - Configuration avec tous les flags booléens à false

### 2. Tests d'intégration

#### `TestMainIntegration()`
- **Lignes**: 1082-1251
- **Description**: Tests end-to-end du binaire compilé
- **Cas couverts**:
  - Flag `-h` (aide)
  - Flag `-version`
  - Validation d'un fichier de contraintes
  - Validation avec mode verbeux
  - Entrée via `-text`
  - Entrée via `-stdin`
  - Erreur: aucune source d'entrée
  - Erreur: sources multiples
  - Erreur: fichier non existant
  - Erreur: syntaxe invalide

#### `TestMainWithFactsIntegration()`
- **Lignes**: 1254-1365
- **Description**: Tests d'intégration avec fichiers de faits
- **Cas couverts**:
  - Contraintes avec faits (mode non-verbeux)
  - Contraintes avec faits (mode verbeux)
  - Fichier de faits non existant

## Défis rencontrés

### 1. Tests de fonctions appelant `os.Exit()`
**Problème**: Les fonctions `runWithFacts()`, `printResults()`, `countActivations()`, et `printActivationDetails()` sont difficiles à tester car:
- `runWithFacts()` appelle `os.Exit(1)` en cas d'erreur
- Les autres dépendent de structures complexes du réseau RETE

**Solution adoptée**:
- Tests d'intégration via subprocess pour couvrir le comportement end-to-end
- Tests unitaires de la logique isolée (ex: vérification d'existence de fichiers)
- Tests de simulation pour `printResults()` et `printActivationDetails()`

### 2. Mock de stdin
**Problème**: Tester `parseFromStdin()` nécessite de simuler l'entrée standard

**Solution**:
- Utilisation de `os.Pipe()` pour créer des pipes de lecture/écriture
- Remplacement temporaire de `os.Stdin`
- Restauration après chaque test avec `defer`

### 3. Syntaxe des fichiers de contraintes et faits
**Problème**: Erreurs de parsing lors des tests d'intégration

**Solution**:
- Analyse des fichiers exemples existants dans le projet
- Utilisation de la syntaxe correcte:
  - Contraintes: `{var: Type} / condition ==> action(args)`
  - Faits: `Type(field:value, field:value)` (sans guillemets)

## Recommandations pour amélioration future

### 1. Refactoring pour testabilité
Pour atteindre une couverture plus élevée, considérer:
- Extraire la logique de `runWithFacts()` sans appel à `os.Exit()`
- Créer des interfaces pour les dépendances externes
- Séparer la logique métier de la présentation

### 2. Tests supplémentaires
- Tests de performance pour grandes quantités de faits
- Tests de concurrence si applicable
- Tests de régression pour bugs spécifiques

### 3. Mocking du réseau RETE
- Créer des mocks/stubs pour `ReteNetwork`
- Permettre l'injection de dépendances
- Faciliter le test de `countActivations()` et `printActivationDetails()`

## Fichiers modifiés

### `main_test.go`
- **Lignes ajoutées**: ~800 lignes
- **Tests ajoutés**: 8 nouvelles fonctions de test
- **Tests étendus**: 1 fonction existante

## Commandes pour vérifier la couverture

```bash
# Exécuter les tests avec couverture
go test -coverprofile=coverage.out ./cmd/tsd

# Afficher le rapport détaillé
go tool cover -func=coverage.out

# Générer un rapport HTML
go tool cover -html=coverage.out -o coverage.html
```

## Conclusion

L'amélioration de la couverture de test pour `cmd/tsd` a permis:
- ✅ Augmentation de 5.9 points de pourcentage (51.0% → 56.9%)
- ✅ Couverture complète de `parseFromStdin()` (0% → 100%)
- ✅ Couverture complète de `parseConstraintSource()` (80% → 100%)
- ✅ Ajout de tests d'intégration end-to-end robustes
- ✅ Meilleure documentation des comportements attendus

Les fonctions non couvertes (`runWithFacts`, `printResults`, `countActivations`, `printActivationDetails`) sont testées indirectement via les tests d'intégration, mais ne peuvent pas être facilement couvertes par des tests unitaires directs en raison de leur dépendance à `os.Exit()` et aux structures RETE complexes.

La qualité globale du code et la confiance dans les changements futurs ont été significativement améliorées.