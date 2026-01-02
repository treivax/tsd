# TODO - Package API

## ✅ Fait

- [x] Créer la structure du package `api`
- [x] Implémenter `pipeline.go` avec NewPipeline() et IngestFile()
- [x] Implémenter `config.go` avec toutes les politiques
- [x] Implémenter `errors.go` avec tous les types d'erreur
- [x] Implémenter `result.go` avec accès aux métriques et xuples
- [x] Créer `doc.go` avec documentation GoDoc complète
- [x] Créer `README.md` avec guide utilisateur
- [x] Tests de base : config (14 tests) ✅
- [x] Tests de base : errors (7 tests) ✅
- [x] Test d'ingestion simple (1 test) ✅
- [x] Validation : `go build ./api` ✅
- [x] Validation : `go vet ./api` ✅

## 📝 À Faire - Tests

### Priorité Haute

- [ ] **Corriger `pipeline_test.go.bak`**
  - Convertir tous les programmes TSD en syntaxe correcte
  - Utiliser `type Name(field: type)` au lieu de `type Name { field: type }`
  - Tester tous les cas : ingestion simple, incrémentale, reset, erreurs
  - Restaurer comme `pipeline_test.go`

- [ ] **Corriger `result_test.go.bak`**
  - Convertir en syntaxe TSD correcte
  - Tester l'accès aux métriques
  - Tester l'accès aux xuples
  - Tester le Summary()
  - Restaurer comme `result_test.go`

- [ ] **Corriger `examples_test.go.bak`**
  - Convertir tous les exemples en syntaxe TSD correcte
  - Vérifier que les Output: commentaires sont corrects
  - Restaurer comme `examples_test.go`

### Priorité Moyenne

- [ ] **Tests de xuple-spaces**
  - Créer un test avec xuple-space-def dans le fichier TSD
  - Tester GetXuples() avec des xuples réels
  - Tester Retrieve() avec différents agents
  - Tester les différentes politiques (FIFO, LIFO, Random)

- [ ] **Tests de concurrence**
  - Test d'ingestion parallèle (plusieurs goroutines)
  - Vérifier le thread-safety du Pipeline
  - Tester les race conditions potentielles

- [ ] **Tests d'intégration**
  - Tester avec de vrais fichiers TSD du projet
  - Tester l'ingestion incrémentale sur cas réels
  - Tester les cas d'erreur (fichiers invalides, etc.)

### Priorité Basse

- [ ] **Benchmarks**
  - Benchmark d'ingestion
  - Benchmark de création de pipeline
  - Benchmark d'accès aux xuples

- [ ] **Tests de couverture**
  - Atteindre > 80% de couverture
  - Identifier les branches non testées
  - Ajouter tests manquants

## 📚 À Faire - Documentation

- [ ] Ajouter plus d'exemples dans `examples_test.go`
  - Exemple avec xuple-spaces
  - Exemple avec politiques personnalisées
  - Exemple d'erreur handling avancé

- [ ] Créer un guide de migration
  - Comment migrer du code existant vers l'API
  - Exemples avant/après

- [ ] Documenter les cas d'usage avancés
  - Configuration complexe
  - Intégration avec applications externes
  - Best practices

## 🚀 À Faire - Fonctionnalités

### Phase 1 (Essentiel)

- [ ] Support complet des xuple-spaces
  - Parsing des définitions depuis TSD
  - Création automatique via factory
  - Configuration des politiques

### Phase 2 (Améliorations)

- [ ] Support de IngestBytes() (depuis []byte)
- [ ] Support de IngestReader() (depuis io.Reader)
- [ ] Méthode GetResult() pour récupérer le dernier résultat
- [ ] Support du reset partiel (types, règles, faits séparément)

### Phase 3 (Avancé)

- [ ] Métriques Prometheus optionnelles
- [ ] Support de plusieurs formats (JSON, YAML en plus de TSD)
- [ ] Système de plugins pour actions personnalisées
- [ ] API REST optionnelle autour du pipeline

## 🔧 À Faire - Maintenance

- [ ] Vérifier la compatibilité avec les autres packages
- [ ] S'assurer qu'il n'y a pas de cycles d'importation
- [ ] Vérifier que les métriques sont cohérentes
- [ ] Vérifier la gestion de la mémoire (pas de leaks)

## ✅ Validation Finale

- [ ] Tous les tests passent (`make test`)
- [ ] Couverture > 80% (`make test-coverage`)
- [ ] Linting propre (`make lint`)
- [ ] Build réussi (`make build`)
- [ ] Validation complète (`make validate`)
- [ ] Documentation à jour

## 📊 Métriques Cibles

- **Tests** : 100% des fichiers testés
- **Couverture** : > 80% des lignes
- **Complexité** : < 15 par fonction
- **Taille** : < 50 lignes par fonction
- **Documentation** : GoDoc pour 100% des exports

---

**Note** : Ce TODO sera mis à jour au fur et à mesure de l'avancement.
