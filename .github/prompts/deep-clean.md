# 🧹 Nettoyage Approfondi du Code (Deep Clean)

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Tu veux effectuer un nettoyage complet et approfondi du code pour éliminer tout le superflu, améliorer la structure, garantir la qualité, et maintenir un code propre et maintenable.

## Objectif

Effectuer un audit complet du projet et nettoyer systématiquement :
- Fichiers inutilisés, doublons, obsolètes
- Code mort, redondant, non utilisé
- Structure des modules et organisation globale
- Documentation obsolète ou manquante
- Tests insuffisants, incorrects ou obsolètes
- Non-conformités aux bonnes pratiques Go

## ⚠️ RÈGLES STRICTES

### 🚫 INTERDICTIONS ABSOLUES

1. **CODE GOLANG** :
   - ❌ AUCUN HARDCODING introduit ou laissé
   - ❌ AUCUNE fonction/variable non utilisée
   - ❌ AUCUN code mort ou commenté
   - ❌ AUCUNE duplication de code
   - ✅ Code générique avec paramètres/interfaces
   - ✅ Constantes nommées pour toutes les valeurs
   - ✅ Respect strict Effective Go

2. **TESTS RETE** :
   - ❌ AUCUNE simulation de résultats
   - ❌ AUCUN test obsolète ou cassé
   - ✅ Extraction depuis réseau RETE réel uniquement
   - ✅ Couverture de tests maximale
   - ✅ Tests déterministes et isolés

3. **FICHIERS** :
   - ❌ AUCUN fichier inutilisé ou en double
   - ❌ AUCUN fichier temporaire ou de backup
   - ❌ AUCUN fichier de rapport en dehors du dossier `tsd/REPORTS`
   - ✅ Organisation claire et logique
   - ✅ Nommage cohérent

## Instructions

### PHASE 1 : AUDIT COMPLET (Analyse)

#### 1.1 Scanner les Fichiers

**Identifier les fichiers problématiques** :

```bash
# Fichiers Go non référencés
find . -name "*.go" -type f | while read file; do
    if ! grep -r "$(basename $file .go)" --include="*.go" . > /dev/null; then
        echo "Potentiellement non utilisé: $file"
    fi
done

# Fichiers temporaires
find . -name "*~" -o -name "*.swp" -o -name "*.bak" -o -name ".DS_Store"

# Fichiers en double (même contenu)
find . -type f -exec md5sum {} + | sort | uniq -w32 -dD

# Fichiers obsolètes (pas modifiés depuis 6+ mois et non utilisés)
find . -name "*.go" -mtime +180 -type f
```

**Questions à poser** :
- Ce fichier est-il importé quelque part ?
- Ce fichier contient-il du code actif ?
- Y a-t-il un doublon de ce fichier ?
- Ce fichier est-il documenté dans le README ?

#### 1.2 Analyser le Code

**Détecter le code mort** :

```bash
# Variables/fonctions non utilisées
go vet ./...
staticcheck ./...
golangci-lint run --enable unused,deadcode,varcheck,structcheck

# Code commenté (suspect)
grep -r "^[[:space:]]*//.*func\|^[[:space:]]*//.*type" --include="*.go" .

# Imports non utilisés
goimports -l .

# Code dupliqué
dupl -threshold 15 ./...
```

**Vérifier** :
- Fonctions/méthodes jamais appelées
- Variables/constantes jamais utilisées
- Types/structs jamais instanciés
- Imports non utilisés
- Code commenté (à supprimer ou documenter)
- Duplication de code (DRY)

#### 1.3 Auditer la Structure

**Analyser l'organisation** :

```bash
# Structure des packages
go list -f '{{.ImportPath}} {{.Imports}}' ./...

# Dépendances cycliques
go list -f '{{.ImportPath}} {{.Imports}}' ./... | grep cycle

# Complexité du code
gocyclo -over 15 .

# Taille des fichiers (> 500 lignes = suspect)
find . -name "*.go" -exec wc -l {} + | awk '$1 > 500'
```

**Questions** :
- Les packages sont-ils bien organisés ?
- Y a-t-il des dépendances circulaires ?
- Les fichiers sont-ils trop gros ?
- La hiérarchie est-elle logique ?

#### 1.4 Vérifier la Documentation

**Audit documentation** :

```bash
# Fonctions exportées sans GoDoc
grep -r "^func [A-Z]" --include="*.go" . | while read line; do
    file=$(echo "$line" | cut -d: -f1)
    func=$(echo "$line" | cut -d: -f2 | awk '{print $2}' | cut -d'(' -f1)
    if ! grep -B1 "^func $func" "$file" | grep -q "^//"; then
        echo "Sans GoDoc: $file:$func"
    fi
done

# README obsolète
git log -1 --format=%ai README.md
git log -1 --format=%ai $(find . -name "*.go" | head -1)

# Documentation vs code
godoc -http=:6060 &
# Vérifier manuellement la cohérence
```

#### 1.5 Auditer les Tests

**Analyse des tests** :

```bash
# Couverture actuelle
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out | grep total

# Tests qui échouent
go test ./... -v 2>&1 | grep FAIL

# Tests obsolètes (fichiers de test sans test)
find . -name "*_test.go" -exec grep -L "func Test" {} \;

# Tests RETE avec simulation (INTERDIT)
grep -r "expectedTokens.*:=.*[0-9]" --include="*_test.go" test/

# Durée des tests
go test -v ./... 2>&1 | grep -E "PASS|FAIL" | awk '{print $NF}'
```

### PHASE 2 : NETTOYAGE (Action)

#### 2.1 Supprimer les Fichiers Inutiles

**Plan d'action** :

1. **Sauvegarder avant suppression** :
   ```bash
   git checkout -b deep-clean-backup
   git add .
   git commit -m "Backup avant nettoyage"
   git checkout -b deep-clean
   ```

2. **Supprimer progressivement** :
   ```bash
   # Fichiers temporaires
   find . -name "*~" -o -name "*.swp" -delete
   
   # Fichiers de backup
   find . -name "*.bak" -o -name "*.backup" -delete
   
   # Fichiers système
   find . -name ".DS_Store" -delete
   ```

3. **Vérifier après chaque suppression** :
   ```bash
   make test
   make rete-unified
   ```

**Checklist par fichier** :
- [ ] Fichier importé nulle part → Supprimer
- [ ] Doublon d'un autre fichier → Supprimer
- [ ] Fichier de test vide → Supprimer
- [ ] Fichier obsolète documenté comme tel → Supprimer

#### 2.2 Éliminer le Code Mort

**Stratégie** :

1. **Fonctions/variables non utilisées** :
   ```bash
   # Identifier avec go vet
   go vet ./...
   
   # Supprimer avec prudence
   # Vérifier que ce n'est pas une API publique
   ```

2. **Code commenté** :
   - Si utile → Convertir en documentation
   - Si obsolète → Supprimer
   - Si exemple → Déplacer dans docs/

3. **Imports inutilisés** :
   ```bash
   goimports -w .
   ```

**Processus** :
1. Lister les éléments non utilisés
2. Vérifier un par un (pas une API publique ?)
3. Supprimer et tester
4. Commit par lot cohérent

#### 2.3 Refactoriser et Déduplication

**Identifier les duplications** :

```bash
dupl -threshold 15 ./... > duplications.txt
```

**Refactoring** :

1. **Extraire fonctions communes** :
   ```go
   // Avant (dupliqué)
   func ProcessA() {
       // 20 lignes de code
   }
   func ProcessB() {
       // 20 mêmes lignes
   }
   
   // Après (factorisé)
   func commonProcess() {
       // 20 lignes une seule fois
   }
   func ProcessA() { commonProcess() }
   func ProcessB() { commonProcess() }
   ```

2. **Simplifier les fonctions complexes** :
   - Fonctions > 50 lignes → Découper
   - Complexité cyclomatique > 15 → Simplifier
   - Imbrication > 4 niveaux → Refactorer

3. **Utiliser des interfaces** :
   - Remplacer code dupliqué par interface
   - Dependency injection

**⚠️ Important** :
- ✅ Code générique avec paramètres
- ❌ Aucun hardcoding introduit
- ✅ Tests passent après chaque refactoring

#### 2.4 Restructurer les Modules

**Organisation cible** :

```
tsd/
├── cmd/              # Binaires (mains)
├── pkg/              # Packages réutilisables
│   ├── rete/        # Moteur RETE
│   ├── constraint/  # Parseur contraintes
│   └── common/      # Utilitaires communs
├── internal/         # Code privé au projet
├── test/            # Tests d'intégration
├── docs/            # Documentation
└── scripts/         # Scripts utilitaires
```

**Actions** :

1. **Regrouper par fonctionnalité** :
   - Tout le RETE dans rete/
   - Tout le parsing dans constraint/
   - Utilitaires dans common/

2. **Séparer public/privé** :
   - API publique → pkg/
   - Implémentation interne → internal/

3. **Éliminer cycles de dépendances** :
   - Utiliser interfaces
   - Inverser les dépendances

#### 2.5 Mettre à Jour la Documentation

**Plan d'action** :

1. **README.md** :
   - [ ] Architecture à jour
   - [ ] Exemples fonctionnels
   - [ ] Installation valide
   - [ ] Commandes correctes

2. **GoDoc** :
   ```bash
   # Ajouter GoDoc manquant
   for file in $(find . -name "*.go"); do
       # Vérifier exports sans doc
   done
   ```

3. **CHANGELOG.md** :
   - Documenter le nettoyage
   - Version et date

4. **Commentaires inline** :
   - Code complexe expliqué
   - TODO/FIXME traités ou supprimés

5. **docs/** :
   - Architecture mise à jour
   - Diagrammes à jour
   - Exemples fonctionnels

#### 2.6 Améliorer les Tests

**Plan d'amélioration** :

1. **Couverture de tests** :
   ```bash
   # Identifier zones non couvertes
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out
   
   # Ajouter tests manquants
   # Objectif : > 80% de couverture
   ```

2. **Tests RETE - Correction stricte** :
   ```go
   // ❌ SUPPRIMER - Simulation interdite
   expectedTokens := 5
   
   // ✅ REMPLACER - Extraction réseau réel
   actualTokens := 0
   for _, terminal := range network.TerminalNodes {
       actualTokens += len(terminal.Memory.GetTokens())
   }
   ```

3. **Tests obsolètes** :
   - Supprimer tests qui ne testent plus rien
   - Mettre à jour tests avec anciennes API
   - Supprimer fichiers *_test.go vides

4. **Tests flaky** :
   - Identifier : `go test -count=10 ./...`
   - Corriger : rendre déterministes
   - Race conditions : `go test -race ./...`

5. **Organisation** :
   - Tests unitaires avec le code
   - Tests d'intégration dans test/
   - Benchmarks dans *_bench_test.go

### PHASE 3 : VALIDATION (Vérification)

#### 3.1 Validation Complète

**Checklist obligatoire** :

```bash
# 1. Formatage
go fmt ./...
goimports -w .

# 2. Analyse statique
go vet ./...
staticcheck ./...
golangci-lint run

# 3. Tests
go test ./...
go test -race ./...
go test -cover ./...

# 4. Build
make build
make build-runners

# 5. Tests d'intégration
make test-integration

# 6. Runner universel
make rete-unified

# 7. Validation complète
make validate
```

**Tous doivent passer** ✅

#### 3.2 Métriques de Qualité

**Vérifier les métriques** :

```bash
# Couverture de tests
go test -cover ./... | grep -E "coverage:"

# Complexité
gocyclo -over 15 . | wc -l  # Doit être 0

# Duplication
dupl -threshold 15 ./... | grep "found" | wc -l  # Doit être 0

# Taille des fonctions
grep -r "^func " --include="*.go" . -A 100 | \
    awk '/^func /{count=0} {count++} /^}$/{if(count>50) print}' | wc -l

# Dette technique
goreportcard-cli  # Si installé
```

**Cibles** :
- Couverture : > 80%
- Complexité : < 15
- Duplication : 0
- Fonctions : < 50 lignes

#### 3.3 Revue de Code

**Auto-revue structurée** :

1. **Architecture** :
   - [ ] Structure des packages logique
   - [ ] Pas de dépendances circulaires
   - [ ] Séparation public/privé claire

2. **Code** :
   - [ ] Aucun hardcoding
   - [ ] Code générique et réutilisable
   - [ ] Pas de duplication
   - [ ] Conventions Go respectées

3. **Tests** :
   - [ ] Couverture > 80%
   - [ ] Tests RETE avec extraction réseau réel
   - [ ] Pas de tests simulés
   - [ ] Tous les tests passent

4. **Documentation** :
   - [ ] README à jour
   - [ ] GoDoc complet
   - [ ] CHANGELOG mis à jour
   - [ ] Exemples fonctionnels

## Critères de Succès

### ✅ Code Nettoyé

- [ ] **AUCUN fichier inutilisé**
- [ ] **AUCUN code mort ou commenté**
- [ ] **AUCUNE duplication**
- [ ] **AUCUN hardcoding**
- [ ] Structure claire et logique
- [ ] Pas de dépendances circulaires

### ✅ Tests Améliorés

- [ ] Couverture > 80%
- [ ] **Tests RETE avec extraction réseau réel uniquement**
- [ ] Tous les tests passent
- [ ] Aucun test flaky
- [ ] Tests déterministes

### ✅ Documentation À Jour

- [ ] README fonctionnel
- [ ] GoDoc complet
- [ ] CHANGELOG mis à jour
- [ ] Architecture documentée
- [ ] Exemples valides

### ✅ Qualité Maximale

- [ ] go vet : 0 erreur
- [ ] golangci-lint : 0 erreur
- [ ] gocyclo < 15 partout
- [ ] Aucune duplication
- [ ] Conventions Go respectées

## Format de Réponse

```
=== NETTOYAGE APPROFONDI DU CODE ===

📊 AUDIT INITIAL

Fichiers :
  • Total : X fichiers Go
  • Non utilisés : X fichiers
  • Doublons : X fichiers
  • Temporaires : X fichiers

Code :
  • Fonctions non utilisées : X
  • Variables non utilisées : X
  • Code commenté : X lignes
  • Duplication : X blocs

Tests :
  • Couverture actuelle : X%
  • Tests qui échouent : X
  • Tests avec simulation : X
  • Tests obsolètes : X

Documentation :
  • GoDoc manquant : X fonctions
  • README obsolète : Oui/Non
  • CHANGELOG à jour : Oui/Non

🧹 ACTIONS DE NETTOYAGE

Phase 1 - Fichiers :
  ✅ Supprimé X fichiers inutilisés
  ✅ Supprimé X fichiers temporaires
  ✅ Éliminé X doublons

Phase 2 - Code :
  ✅ Supprimé X fonctions mortes
  ✅ Nettoyé X lignes de code commenté
  ✅ Refactorisé X blocs dupliqués
  ✅ Éliminé X hardcodings

Phase 3 - Tests :
  ✅ Ajouté X tests manquants
  ✅ Corrigé X tests RETE (extraction réseau réel)
  ✅ Supprimé X tests obsolètes
  ✅ Couverture : X% → Y%

Phase 4 - Documentation :
  ✅ Ajouté GoDoc pour X fonctions
  ✅ Mis à jour README
  ✅ Mis à jour CHANGELOG
  ✅ Corrigé X exemples

✅ VALIDATION FINALE

Tests :
  ✅ go test ./... : PASS
  ✅ go test -race ./... : PASS
  ✅ make rete-unified : 58/58 ✅

Qualité :
  ✅ go vet : 0 erreur
  ✅ golangci-lint : 0 erreur
  ✅ Couverture : Y% (> 80%)
  ✅ Complexité < 15 partout
  ✅ Aucune duplication

Structure :
  ✅ Packages bien organisés
  ✅ Aucune dépendance circulaire
  ✅ Séparation public/privé claire

📈 RÉSULTATS

Avant → Après :
  • Fichiers Go : X → Y (−Z)
  • Lignes de code : X → Y (−Z)
  • Couverture tests : X% → Y%
  • Complexité max : X → Y
  • Duplication : X → 0
  • Dette technique : Haute → Basse

🎯 VERDICT : CODE PROPRE ET MAINTENABLE ✅
```

## Exemple d'Utilisation

```
Le projet TSD s'est accumulé du code au fil du temps et je veux
faire un grand nettoyage de printemps.

Utilise le prompt "deep-clean" pour faire un audit complet et nettoyer :
- Fichiers inutilisés
- Code mort
- Duplication
- Tests obsolètes
- Documentation

Objectif : Code propre, maintenable, et performant.
```

## Commandes Utiles

```bash
# Audit complet
make validate
go vet ./...
staticcheck ./...
golangci-lint run --enable-all

# Détection code mort
go-cleanarch
deadcode ./...

# Duplication
dupl -threshold 15 ./...

# Complexité
gocyclo -over 15 .

# Couverture
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Dépendances
go mod tidy
go mod verify
```

## Checklist Finale

### Avant de Commencer
- [ ] Backup complet (git commit + branch)
- [ ] Tests passent actuellement
- [ ] Documentation des objectifs

### Pendant le Nettoyage
- [ ] Travailler par petits commits
- [ ] Tester après chaque modification
- [ ] Documenter les suppressions importantes

### Après le Nettoyage
- [ ] **Tous les tests passent** ✅
- [ ] **Aucun hardcoding introduit** ✅
- [ ] **Tests RETE avec extraction réseau réel** ✅
- [ ] go vet et golangci-lint sans erreur ✅
- [ ] Couverture > 80% ✅
- [ ] Documentation à jour ✅
- [ ] Code review effectuée ✅

## Bonnes Pratiques

1. **Progressif** : Nettoyer par petites étapes
2. **Testé** : Valider après chaque modification
3. **Documenté** : Expliquer les changements importants
4. **Réversible** : Commits fréquents, branches dédiées
5. **Complet** : Suivre toutes les phases
6. **Rigoureux** : Respecter les règles strictes

## Avertissements

⚠️ **ATTENTION** :
- Ne jamais supprimer sans backup
- Ne jamais supposer qu'un fichier est inutile sans vérifier
- Toujours tester après suppression
- Documenter les suppressions non-évidentes
- Préserver les API publiques

## Ressources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Clean Code](https://www.amazon.com/Clean-Code-Handbook-Software-Craftsmanship/dp/0132350882)

---

**Version** : 1.0  
**Dernière mise à jour** : Novembre 2025  
**Durée estimée** : 4-8 heures selon taille projet
