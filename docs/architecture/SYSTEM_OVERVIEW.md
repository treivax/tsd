# 🏗️ Architecture Système TSD - Vue d'Ensemble

**Date** : 2025-12-15  
**Version** : 1.0.0  
**Statut** : Documentation officielle

---

## 📋 Table des Matières

1. [Vue Globale](#vue-globale)
2. [Architecture Modulaire](#architecture-modulaire)
3. [Flux de Données](#flux-de-données)
4. [Composants Principaux](#composants-principaux)
5. [Sécurité](#sécurité)

---

## 🎯 Vue Globale

TSD est un système de règles métier basé sur l'algorithme RETE, avec architecture client-serveur HTTPS et authentification intégrée.

### Principes Architecturaux

1. **Binaire Unique Multi-Rôles** : Un seul binaire `tsd` avec dispatcher intelligent
2. **Séparation CLI / Logique** : `cmd/` dispatche, `internal/` implémente
3. **Modules Réutilisables** : `auth/`, `tsdio/`, `constraint/`, `rete/` publics
4. **Isolation Interne** : `internal/` empêche utilisation hors projet
5. **Centralisation Configuration** : TLS centralisé dans `tlsconfig/`

---

## 🧩 Architecture Modulaire

### Structure des Packages

```
tsd/
├── cmd/tsd/                    # Point d'entrée unique (177 lignes)
│   ├── main.go                 # Dispatcher multi-rôles
│   └── unified_test.go         # Tests du dispatcher
│
├── internal/                   # Packages internes non exportables
│   ├── authcmd/               # Commandes auth (génération clés, JWT, certificats)
│   ├── clientcmd/             # Client HTTPS/TLS
│   ├── servercmd/             # Serveur HTTPS/TLS avec endpoints
│   ├── compilercmd/           # Compilateur/Runner TSD local
│   └── tlsconfig/             # Configuration TLS centralisée
│
├── auth/                       # Module authentification public (313 lignes)
├── tsdio/                      # I/O thread-safe + types API (400 lignes)
├── constraint/                 # Parser de programmes TSD
├── rete/                       # Moteur RETE
│
└── tests/                      # Tests organisés
    ├── e2e/                   # Tests end-to-end
    ├── integration/           # Tests d'intégration
    ├── performance/           # Benchmarks
    └── shared/testutil/       # Utilitaires de test
```

**Points clés** :
- ✅ Graphe **acyclique** (DAG) - Pas de cycles de dépendances
- ✅ Dépendances **unidirectionnelles** : `cmd/` → `internal/` → modules publics
- ✅ Modules **sans dépendances** : `auth/`, `tsdio/`, `tlsconfig/` (excellent découplage)

---

## 🔄 Flux de Données

### Compilation Locale

1. User exécute `tsd program.tsd`
2. CLI dispatche vers `compilercmd`
3. Parser (`constraint`) analyse le code source
4. Construction du réseau RETE
5. Injection des faits
6. Propagation des tokens
7. Retour des activations
8. Affichage résultats

### Exécution Client-Serveur

1. Client lit fichier `.tsd`
2. Client envoie POST `/execute` avec code source
3. Serveur valide token d'authentification
4. Serveur parse et exécute via RETE
5. Serveur retourne résultats JSON
6. Client affiche résultats

📊 **Voir les diagrammes détaillés** : [Flux de Données](diagrams/02-data-flow.md)

---

## 🔧 Composants Principaux

### 1. Dispatcher (`cmd/tsd/main.go`)

**Responsabilité** : Router vers le rôle approprié

**Rôles supportés** :
- `auth` → `internal/authcmd`
- `client` → `internal/clientcmd`
- `server` → `internal/servercmd`
- (défaut) → `internal/compilercmd`

### 2. Auth Manager (`auth/`)

**Responsabilité** : Gestion authentification (Auth Key + JWT)

**Types d'authentification** :
- `none` : Pas d'authentification
- `key` : Clés API statiques
- `jwt` : JSON Web Tokens

### 3. TLS Config (`internal/tlsconfig/`)

**Responsabilité** : Configuration TLS centralisée

**Avantages** :
- ✅ Évite duplication code
- ✅ Standards de sécurité cohérents
- ✅ Configuration par défaut sécurisée

### 4. Logger Thread-Safe (`tsdio/`)

**Responsabilité** : I/O sécurisé pour utilisation concurrente

**Solution** :
- Mutex global sur toutes les écritures
- API simple : `tsdio.Printf()`, `tsdio.Println()`
- Support redirection pour tests

---

## 🔒 Sécurité

### Bonnes Pratiques Implémentées

1. **Constant-time Comparison** : Protection contre timing attacks
2. **TLS par défaut** : HTTPS activé par défaut
3. **JWT Standards** : Signature HMAC-SHA256
4. **Validation stricte** : Rejet tokens malformés
5. **Pas de credentials hardcodés** : Configuration externe

---

## 📊 Métriques Architecture

| Métrique | Valeur |
|----------|--------|
| **Packages totaux** | 8 (hors constraint/rete) |
| **Lignes code production** | ~4540 |
| **Lignes code tests** | ~10534 |
| **Ratio tests/production** | 2.3:1 |
| **Couverture tests** | 81.3% |
| **Dépendances directes** | 5 |
| **Cycles de dépendances** | 0 |

---

## 📚 Documentation Visuelle Complète

Ce document fournit une vue d'ensemble textuelle. Pour une compréhension approfondie avec diagrammes visuels :

### 🎯 Diagrammes d'Architecture
Consultez le répertoire [diagrams/](diagrams/) qui contient :

1. **[Architecture Globale](diagrams/01-global-architecture.md)** - Vue d'ensemble système, couches, dépendances
2. **[Flux de Données](diagrams/02-data-flow.md)** - Séquences, propagation tokens, compilation
3. **[Moteur RETE](diagrams/03-rete-architecture.md)** - Nœuds Alpha/Beta, optimisations
4. **[Sécurité](diagrams/04-security-flow.md)** - Authentification, TLS, JWT
5. **[Modèle de Données](diagrams/05-data-model.md)** - Types, règles, contraintes

📋 **Index complet** : [diagrams/README.md](diagrams/README.md)

### 🎓 Guide par Profil

| Profil | Documents Recommandés |
|--------|----------------------|
| **Nouveau contributeur** | [SYSTEM_OVERVIEW.md](SYSTEM_OVERVIEW.md) → [Architecture Globale](diagrams/01-global-architecture.md) → [Flux de Données](diagrams/02-data-flow.md) |
| **Développeur RETE** | [Moteur RETE](diagrams/03-rete-architecture.md) → [Flux de Données](diagrams/02-data-flow.md) |
| **DevOps / Sécurité** | [Architecture Globale](diagrams/01-global-architecture.md) → [Sécurité](diagrams/04-security-flow.md) |
| **Utilisateur API** | [Modèle de Données](diagrams/05-data-model.md) → [Flux de Données](diagrams/02-data-flow.md) |

---

**Maintenu par** : TSD Contributors  
**Dernière mise à jour** : 2025-12-16
