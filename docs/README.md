# Documentation TSD

Documentation centralisée du projet TSD - Moteur de règles RETE avec système de contraintes.

**Version** : 2.0.0

---

## 🎯 Par Où Commencer ?

### Nouveaux Utilisateurs

1. [Démarrage Rapide](../README.md#démarrage-rapide) - Installation et premier programme
2. [Affectations de Faits](user-guide/fact-assignments.md) - Nommer et réutiliser des faits
3. [Comparaisons de Faits](user-guide/fact-comparisons.md) - Relations entre faits
4. [Exemples](../examples/) - Programmes d'exemple complets

### Migration depuis v1.x

⚠️ **Important** : La v2.0 introduit des breaking changes.

1. **[Guide de Migration v1.x → v2.0](migration/from-v1.x.md)** - ⚠️ **OBLIGATOIRE**
2. [Identifiants Internes](internal-ids.md) - Nouveau système `_id_`
3. [Nouveautés v2.0](../README.md#nouveautés-v20) - Résumé des changements

---

## 📚 Documentation Utilisateur

### Guides Essentiels

| Guide | Description |
|-------|-------------|
| **[Affectations de Faits](user-guide/fact-assignments.md)** | Créer et nommer des faits avec `variable = Type(...)` |
| **[Comparaisons de Faits](user-guide/fact-comparisons.md)** | Comparer des faits directement avec `==` |
| **[Système de Types](user-guide/type-system.md)** | Types primitifs et types de faits dans les champs |
| **[Clés Primaires](primary-keys.md)** | Génération automatique d'IDs avec `#` |
| **[Actions](actions/README.md)** | Actions par défaut et personnalisées |
| **[Changements de Syntaxe](syntax-changes.md)** | Évolution de la syntaxe du langage |

### Référence Technique

| Document | Contenu |
|----------|---------|
| **[Identifiants Internes](internal-ids.md)** | Système `_id_` : génération, format, règles |
| **[Référence Syntaxe](reference.md)** | Grammaire complète du langage TSD |
| **[API Publique](api.md)** | Interface Go pour intégration |
| **[Architecture](architecture.md)** | Algorithme RETE et architecture interne |

### Tutoriels

| Tutoriel | Niveau | Sujets |
|----------|--------|--------|
| **[Clés Primaires](tutorials/primary-keys-tutorial.md)** | Débutant | Clés simples, composites, hash |
| *Plus à venir* | | |

---

## 🔧 Documentation Technique

### Architecture

| Document | Description |
|----------|-------------|
| **[Vue d'Ensemble](architecture.md)** | Architecture globale du système |
| **[Génération d'IDs](architecture/id-generation.md)** | Algorithme de génération des identifiants |
| **[Diagrammes](architecture/diagrams/)** | Diagrammes d'architecture |

### API

| Package | Documentation |
|---------|---------------|
| **[constraint](api/constraint.md)** | Parser et validation |
| **[rete](api/rete.md)** | Moteur RETE |
| **[tsdio](api/tsdio.md)** | Structures I/O |

---

## 📦 Migration et Mises à Jour

| Document | Description |
|----------|-------------|
| **[Guide de Migration v1.x → v2.0](migration/from-v1.x.md)** | ⚠️ **Breaking changes** et migration complète |
| **[CHANGELOG](../CHANGELOG.md)** | Historique des versions |

---

## 💡 Exemples

| Répertoire | Description |
|------------|-------------|
| **[examples/](../examples/)** | Programmes TSD complets |
| **[tests/fixtures/](../tests/fixtures/)** | Fixtures de test (cas d'usage) |

---

## 🤝 Contribution

| Document | Description |
|----------|-------------|
| **[CONTRIBUTING.md](../CONTRIBUTING.md)** | Guide de contribution |
| **[.github/prompts/common.md](../.github/prompts/common.md)** | Standards de code |
| **[.github/prompts/develop.md](../.github/prompts/develop.md)** | Standards de développement |

---

## 🔍 Index par Fonctionnalité

### Identifiants et Clés Primaires

- [Identifiants Internes](internal-ids.md) - Système `_id_` complet
- [Clés Primaires](primary-keys.md) - Syntaxe `#field`
- [Guide de Migration](migration/from-v1.x.md) - Ancien système `id` → nouveau `_id_`

### Affectations et Comparaisons

- [Affectations](user-guide/fact-assignments.md) - `variable = Type(...)`
- [Comparaisons](user-guide/fact-comparisons.md) - `fact1 == fact2`

### Types

- [Système de Types](user-guide/type-system.md) - Primitifs et types de faits
- [Référence](reference.md) - Grammaire complète

### Règles

- [Guides](guides.md) - Syntaxe des règles
- [Référence](reference.md) - Conditions, actions, opérateurs

---

## 📊 Index par Niveau

### Débutant

- [README Principal](../README.md) - Vue d'ensemble
- [Installation](installation.md) - Démarrage
- [Affectations](user-guide/fact-assignments.md) - Bases
- [Tutoriels](tutorials/) - Apprentissage guidé

### Intermédiaire

- [Comparaisons](user-guide/fact-comparisons.md) - Relations
- [Système de Types](user-guide/type-system.md) - Types avancés
- [Clés Primaires](primary-keys.md) - IDs personnalisés
- [Exemples](../examples/) - Cas d'usage réels

### Avancé

- [Architecture](architecture.md) - RETE et internals
- [API](api.md) - Intégration Go
- [Identifiants Internes](internal-ids.md) - Détails techniques

---

## 📞 Support

| Ressource | Description |
|-----------|-------------|
| **[Issues GitHub](https://github.com/chrlesur/tsd/issues)** | Rapporter des bugs et demander de l'aide |
| **[Guide de Migration](migration/from-v1.x.md)** | Aide pour migration v1.x → v2.0 |

---

## 📖 Documentation par Module

### Modules Principaux

| Module | README | Documentation |
|--------|--------|---------------|
| **constraint** | [constraint/README.md](../constraint/README.md) | Parser, validation, types |
| **rete** | [rete/README.md](../rete/README.md) | Moteur RETE |
| **tsdio** | [tsdio/README.md](../tsdio/README.md) | I/O et structures |
| **xuples** | [xuples/README.md](../xuples/README.md) | Espace de tuples |

### Modules Spécialisés

| Module | Documentation |
|--------|---------------|
| **[Actions](actions/)** | Actions CRUD et Xuple |

### Archives

Les anciennes documentations sont archivées dans :
- [docs/archive/](archive/) - Documentation pré-v2.0
- [docs/archive/constraint/](archive/constraint/) - Anciennes docs constraint
- [docs/archive/rete/](archive/rete/) - Anciennes docs RETE
- [ARCHIVES/cleanup-20260102/](../ARCHIVES/cleanup-20260102/) - Fichiers nettoyés (2025-01-02)

---

## 🎯 Résumé v2.0

### Fonctionnalités Principales

✅ **Affectations** : `alice = User("alice", "alice@example.com")`  
✅ **Comparaisons** : `{u: User, o: Order} / o.customer == u`  
✅ **Types de faits** : `Order(customer: Customer, ...)`  
✅ **IDs cachés** : `_id_` interne, jamais accessible  
✅ **Type-safety** : Validation complète au parsing  

### Breaking Changes

❌ `id` → `_id_` (caché, inaccessible)  
❌ Pas d'affectation manuelle d'ID  
❌ Pas d'accès à `_id_` dans expressions  

**Voir** : [Guide de Migration](migration/from-v1.x.md)

---

**Version** : 2.0.0  
**Dernière mise à jour** : 2025-12-19  
**Mainteneur** : Équipe TSD

### 3. Suite

→ [Installation complète](installation.md)  
→ [Tutoriel détaillé](guides.md#guide-débutant)

---

## 🎯 Parcours d'Apprentissage

### Débutant (2-4 heures)

1. [Installation](installation.md) - Installer TSD
2. [Démarrage rapide](installation.md#démarrage-rapide-5-minutes) - Premier programme
3. [Guide débutant](guides.md#guide-débutant) - Apprendre les bases

**Vous saurez :** Créer types, règles, actions et exécuter des programmes simples.

### Développeur (1-2 jours)

1. [Guide développeur](guides.md#guide-développeur) - Syntaxe avancée
2. [API Go](api.md) - Intégration programmatique
3. [Configuration](configuration.md) - Optimiser les performances
4. [API HTTP](reference.md#api-httprest) - Mode serveur

**Vous saurez :** Intégrer TSD dans vos applications Go, configurer pour production, utiliser l'API REST.

### Avancé (1 semaine)

1. [Architecture](architecture.md) - Comprendre l'algorithme RETE
2. [Guide avancé](guides.md#guide-avancé) - Patterns complexes
3. [Optimisations](architecture.md#optimisations) - Performance maximale
4. [Contribution](reference.md#contribution) - Contribuer au projet

**Vous saurez :** Optimiser les règles, comprendre les performances, contribuer au code.

---

## 🔍 Navigation Rapide

### Je veux...

- **Installer TSD** → [Installation](installation.md)
- **Apprendre la syntaxe** → [Guides](guides.md)
- **Utiliser les clés primaires** → [Clés Primaires](primary-keys.md)
- **Utiliser les actions** → [Actions](actions/README.md)
- **Comprendre les changements de syntaxe** → [Changements de Syntaxe](syntax-changes.md)
- **Migrer vers les clés primaires** → [Migration Guide](MIGRATION_IDS.md)
- **Configurer le système** → [Configuration](configuration.md)
- **Intégrer dans mon app Go** → [API](api.md)
- **Utiliser le serveur HTTP** → [Référence API HTTP](reference.md#api-httprest)
- **Comprendre la grammaire** → [Référence Grammaire](reference.md#grammaire-tsd)
- **Sécuriser avec auth** → [Référence Auth](reference.md#authentification)
- **Débugger** → [Référence Logging](reference.md#logging)
- **Contribuer** → [Référence Contribution](reference.md#contribution)
- **Comprendre l'architecture** → [Architecture](architecture.md)

---

## 📖 Ressources Additionnelles

### Fichiers du Projet

- [README Principal](../README.md) - Vue d'ensemble du projet
- [CHANGELOG](../CHANGELOG.md) - Historique des versions
- [TODO Actifs](../TODO_ACTIFS.md) - Améliorations futures
- [Archives](../ARCHIVES/README.md) - Documentation archivée
- [Reports](../REPORTS/README.md) - Rapports techniques

### Tutoriels

- [Tutoriel Clés Primaires](tutorials/primary-keys-tutorial.md) - Système de blog complet (30 min)

### Référence API

- [API ID Generator](api/id-generator.md) - Référence complète API génération d'IDs
- [Architecture ID Generation](architecture/id-generation.md) - Architecture interne

### Exemples

```bash
ls examples/          # Explorer les exemples
tsd examples/*.tsd    # Exécuter les exemples
```

### Aide

- **GitHub Issues** : [https://github.com/treivax/tsd/issues](https://github.com/treivax/tsd/issues)
- **Debug** : `TSD_LOG_LEVEL=debug tsd program.tsd`
- **Help** : `tsd --help`

---

## 🎓 Concepts Clés

### Types

Définissent la structure des données :

```tsd
type Product(#sku: string, name: string, price: number, inStock: bool)
```

**Note :** Le préfixe `#` marque `sku` comme clé primaire. L'ID sera généré automatiquement : `Product~LAPTOP-001`

### Faits

Instances de types :

```tsd
Product(sku: "LAPTOP-001", name: "Laptop", price: 999.99, inStock: true)
// ID généré automatiquement: Product~LAPTOP-001
```

### Règles

Logique métier avec pattern matching :

```tsd
rule expensive : {p: Product} / p.price > 500 ==> markAsPremium(p.name, p.id)
```

**Note :** Le champ `id` est toujours disponible et contient l'ID généré automatiquement.

### Actions

Déclenchées par les règles :

```tsd
action markAsPremium(name: string)
```

---

## 🏗️ Architecture

TSD utilise l'**algorithme RETE** pour une évaluation efficace des règles :

- **Alpha Network** : Filtrage des faits par conditions simples
- **Beta Network** : Jointures entre faits multiples  
- **Optimisations** : Partage de nœuds (alpha/beta sharing)
- **Performance** : Cache LRU, normalisation, passthrough

→ [Architecture complète](architecture.md)

---

## ⚙️ Configuration

Profils prédéfinis pour différents usages :

| Profil | Usage | Performance | Mémoire |
|--------|-------|-------------|---------|
| **Développement** | Debug, tests | Normale | Normale |
| **Production** | Déploiement | Maximale | Optimisée |
| **Test** | CI/CD | Déterministe | Contrôlée |
| **Embarqué** | IoT, edge | Réduite | Minimale |

→ [Configuration complète](configuration.md)

---

## 🔐 Sécurité

- **API Keys** : Authentification simple pour scripts/CI
- **JWT** : Authentification avancée pour applications
- **TLS/HTTPS** : Transport chiffré obligatoire en production
- **Validation** : Entrées validées strictement

→ [Authentification](reference.md#authentification)

---

## 📊 Monitoring

- **Métriques Prometheus** : Exposition sur `/metrics`
- **Logging Structuré** : Niveaux ERROR/WARN/INFO/DEBUG/TRACE
- **Health Check** : Endpoint `/health`

→ [Logging](reference.md#logging)

---

## 🤝 Contribution

Contributions bienvenues ! 

1. Fork le projet
2. Créer une branche feature
3. Committer vos changements
4. Pousser et créer une Pull Request

→ [Guide de contribution](reference.md#contribution)

---

## 📝 Conventions

### Documentation Code (GoDoc)

- **Langue** : Anglais
- **Format** : GoDoc standard
- **Cible** : Développeurs utilisant l'API Go

### Documentation Technique

- **Langue** : Français
- **Format** : Markdown
- **Cible** : Utilisateurs et contributeurs

### Exemples

Tous les exemples de code doivent être :
- ✅ Testables
- ✅ Fonctionnels
- ✅ Documentés

---

## 📅 Versions

- **Version actuelle** : 1.0.0
- **Statut** : ✅ Production Ready
- **Go minimum** : 1.21+

---

**Bon développement avec TSD ! 🚀**