# 🏗️ Architecture Globale TSD

**Date** : 2025-12-16  
**Version** : 1.0.0  
**Statut** : Documentation officielle

---

## Vue d'Ensemble du Système

Ce diagramme présente l'architecture globale du système TSD avec ses modules principaux et leurs interactions.

```mermaid
graph TB
    subgraph "User Interface"
        CLI[👤 CLI User]
        HTTP[🌐 HTTP Client]
    end

    subgraph "TSD Binary - cmd/tsd/"
        DISPATCHER[🎯 Main Dispatcher<br/>main.go]
    end

    subgraph "Internal Commands - internal/"
        COMPILER[📦 CompilerCmd<br/>Local Execution]
        SERVER[🖥️ ServerCmd<br/>HTTPS Server]
        CLIENT[📡 ClientCmd<br/>HTTPS Client]
        AUTH[🔐 AuthCmd<br/>Auth Management]
    end

    subgraph "Core Modules"
        CONSTRAINT[📝 Constraint<br/>Parser + AST]
        RETE[🧠 RETE Engine<br/>Rule Inference]
        AUTHMOD[🔒 Auth Module<br/>Auth + JWT]
        TSDIO[📄 TSDIO<br/>Thread-Safe I/O]
        TLSCONF[🔐 TLS Config<br/>Security]
    end

    subgraph "Storage"
        MEMORY[💾 In-Memory Store<br/>Facts + Rules]
    end

    CLI -->|tsd program.tsd| DISPATCHER
    CLI -->|tsd auth generate-key| DISPATCHER
    CLI -->|tsd server start| DISPATCHER
    HTTP -->|POST /execute| SERVER

    DISPATCHER -->|Default| COMPILER
    DISPATCHER -->|auth cmd| AUTH
    DISPATCHER -->|server cmd| SERVER
    DISPATCHER -->|client cmd| CLIENT

    COMPILER --> CONSTRAINT
    COMPILER --> RETE
    
    SERVER --> AUTHMOD
    SERVER --> TLSCONF
    SERVER --> CONSTRAINT
    SERVER --> RETE
    
    CLIENT --> AUTHMOD
    CLIENT --> TLSCONF
    
    AUTH --> AUTHMOD
    AUTH --> TLSCONF

    CONSTRAINT --> RETE
    RETE --> MEMORY
    
    COMPILER -.-> TSDIO
    SERVER -.-> TSDIO
    CLIENT -.-> TSDIO

    style DISPATCHER fill:#4A90E2,stroke:#333,stroke-width:3px,color:#fff
    style RETE fill:#E74C3C,stroke:#333,stroke-width:2px,color:#fff
    style CONSTRAINT fill:#F39C12,stroke:#333,stroke-width:2px,color:#fff
    style AUTHMOD fill:#27AE60,stroke:#333,stroke-width:2px,color:#fff
    style MEMORY fill:#9B59B6,stroke:#333,stroke-width:2px,color:#fff
```

---

## Architecture en Couches

```mermaid
graph LR
    subgraph "Layer 1: Entry Point"
        E[cmd/tsd/main.go]
    end

    subgraph "Layer 2: Commands"
        C1[internal/compilercmd]
        C2[internal/servercmd]
        C3[internal/clientcmd]
        C4[internal/authcmd]
    end

    subgraph "Layer 3: Core Logic"
        M1[constraint/]
        M2[rete/]
        M3[auth/]
        M4[tsdio/]
        M5[internal/tlsconfig/]
    end

    subgraph "Layer 4: Storage"
        S[In-Memory Storage]
    end

    E --> C1
    E --> C2
    E --> C3
    E --> C4

    C1 --> M1
    C1 --> M2
    C1 --> M4
    
    C2 --> M1
    C2 --> M2
    C2 --> M3
    C2 --> M4
    C2 --> M5
    
    C3 --> M3
    C3 --> M4
    C3 --> M5
    
    C4 --> M3
    C4 --> M5

    M2 --> S

    style E fill:#4A90E2,stroke:#333,stroke-width:3px
    style S fill:#9B59B6,stroke:#333,stroke-width:2px
```

---

## Graphe de Dépendances

Ce diagramme montre les dépendances entre modules (DAG - Directed Acyclic Graph).

```mermaid
graph TD
    CMD[cmd/tsd/main.go]
    
    COMPILERCMD[internal/compilercmd]
    SERVERCMD[internal/servercmd]
    CLIENTCMD[internal/clientcmd]
    AUTHCMD[internal/authcmd]
    
    CONSTRAINT[constraint/]
    RETE[rete/]
    AUTH[auth/]
    TSDIO[tsdio/]
    TLSCONFIG[internal/tlsconfig/]
    
    CMD --> COMPILERCMD
    CMD --> SERVERCMD
    CMD --> CLIENTCMD
    CMD --> AUTHCMD
    
    COMPILERCMD --> CONSTRAINT
    COMPILERCMD --> RETE
    COMPILERCMD --> TSDIO
    
    SERVERCMD --> CONSTRAINT
    SERVERCMD --> RETE
    SERVERCMD --> AUTH
    SERVERCMD --> TSDIO
    SERVERCMD --> TLSCONFIG
    
    CLIENTCMD --> AUTH
    CLIENTCMD --> TSDIO
    CLIENTCMD --> TLSCONFIG
    
    AUTHCMD --> AUTH
    AUTHCMD --> TLSCONFIG
    
    CONSTRAINT --> RETE

    style CMD fill:#4A90E2,color:#fff
    style RETE fill:#E74C3C,color:#fff
    style CONSTRAINT fill:#F39C12,color:#fff
    style AUTH fill:#27AE60,color:#fff
    
    classDef independent fill:#95A5A6,stroke:#333,stroke-width:2px
    class TSDIO,TLSCONFIG independent
```

**✅ Points clés :**
- Graphe **acyclique** (pas de cycles)
- Dépendances **unidirectionnelles**
- Modules indépendants : `auth/`, `tsdio/`, `internal/tlsconfig/`
- Réutilisabilité maximale

---

## Modules Principaux

### 1. cmd/tsd/ - Point d'Entrée Unique
- **Lignes** : ~177
- **Rôle** : Dispatcher intelligent multi-rôles
- **Responsabilité** : Router vers la commande appropriée

### 2. internal/compilercmd/ - Compilateur Local
- **Rôle** : Exécution locale de programmes TSD
- **Flux** : Fichier TSD → Parser → RETE → Résultats

### 3. internal/servercmd/ - Serveur HTTPS
- **Rôle** : Serveur HTTPS avec authentification
- **Endpoints** :
  - `POST /execute` : Exécuter programme TSD
  - `GET /health` : Health check
  - `GET /metrics` : Métriques Prometheus

### 4. internal/clientcmd/ - Client HTTPS
- **Rôle** : Client pour exécution distante
- **Fonctionnalités** : Envoie code TSD au serveur

### 5. internal/authcmd/ - Gestion Authentification
- **Rôle** : Génération clés, JWT, certificats TLS
- **Commandes** :
  - `generate-key` : Génère clé API
  - `generate-jwt` : Génère token JWT
  - `generate-certs` : Génère certificats TLS

### 6. constraint/ - Parser
- **Rôle** : Analyse syntaxique du langage TSD
- **Sortie** : AST (Abstract Syntax Tree)

### 7. rete/ - Moteur d'Inférence
- **Rôle** : Exécution des règles (algorithme RETE)
- **Optimisations** :
  - Alpha sharing
  - Beta sharing
  - Result caching
  - Token pooling

### 8. auth/ - Module Authentification
- **Lignes** : ~313
- **Types** : Auth Key, JWT
- **Indépendant** : Aucune dépendance interne

### 9. tsdio/ - I/O Thread-Safe
- **Lignes** : ~400
- **Rôle** : Logging sécurisé pour concurrence
- **Indépendant** : Aucune dépendance interne

### 10. internal/tlsconfig/ - Configuration TLS
- **Rôle** : Configuration TLS centralisée
- **Avantages** : Standards sécurité uniformes

---

## Métriques Architecture

| Métrique | Valeur |
|----------|--------|
| **Packages totaux** | 10 |
| **Lignes code production** | ~4540 |
| **Lignes code tests** | ~10534 |
| **Ratio tests/production** | 2.3:1 |
| **Couverture tests** | 81.3% |
| **Cycles de dépendances** | 0 |
| **Dépendances externes** | 5 |

---

## Références

- [Architecture Détaillée](../architecture.md)
- [Vue d'Ensemble Système](../SYSTEM_OVERVIEW.md)
- [Flux de Données](02-data-flow.md)
- [RETE Engine](03-rete-architecture.md)

---

**Maintenu par** : TSD Contributors  
**Dernière mise à jour** : 2025-12-16
