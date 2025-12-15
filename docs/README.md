# Documentation TSD

**Guide complet du moteur de règles TSD**

---

## 📚 Documentation Principale

| Document | Description |
|----------|-------------|
| **[Installation](installation.md)** | Installation et démarrage rapide (5 min) |
| **[Guides](guides.md)** | Tutoriels et guides utilisateur complets |
| **[Architecture](architecture.md)** | Architecture interne et algorithme RETE |
| **[Configuration](configuration.md)** | Configuration système complète |
| **[API](api.md)** | API publique Go |
| **[Référence](reference.md)** | API HTTP, grammaire, auth, logging, contribution |

---

## 🚀 Démarrage Rapide

### 1. Installation

```bash
git clone https://github.com/treivax/tsd.git
cd tsd
make build
./bin/tsd --version
```

### 2. Premier Programme

Créez `hello.tsd` :

```tsd
type Person(name: string, age: number)
action greet(name: string)

rule welcome : {p: Person} / p.age >= 18 ==> greet(p.name)

Person(name: "Alice", age: 25)
Person(name: "Bob", age: 16)
```

Exécutez :

```bash
tsd hello.tsd
```

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
type Product(name: string, price: number, inStock: bool)
```

### Faits

Instances de types :

```tsd
Product(name: "Laptop", price: 999.99, inStock: true)
```

### Règles

Logique métier avec pattern matching :

```tsd
rule expensive : {p: Product} / p.price > 500 ==> markAsPremium(p.name)
```

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