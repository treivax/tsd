# TSD - Système de Traitement de Contraintes et Réseau RETE

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Build Status](https://img.shields.io/badge/Build-Passing-green.svg)](https://github.com/treivax/tsd)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

**TSD** est un système avancé de traitement de contraintes métier intégré à un réseau RETE haute performance. Il permet de définir des règles métier complexes via une syntaxe déclarative et de les exécuter efficacement grâce à l'algorithme RETE optimisé.

## 🎯 Fonctionnalités Principales

### � **Module Constraint**
- **Grammaire PEG complète** pour définir des contraintes métier
- **Parser robuste** généré automatiquement avec validation syntaxique
- **Actions obligatoires** garantissant des règles métier complètes
- **Validation sémantique** avec vérification de types
- **Support complet** : négation, quantification existentielle, agrégation

### ⚡ **Module RETE**
- **Implémentation optimisée** de l'algorithme RETE
- **Architecture modulaire** : AlphaNode, BetaNode, NotNode, ExistsNode, AccumulateNode
- **Monitoring en temps réel** avec interface web intégrée
- **Performance élevée** avec cache d'évaluation et optimisations
- **Compatibilité complète** entre grammaire PEG et réseau RETE

### 🌐 **Interface de Monitoring**
- **Dashboard web** en temps réel pour visualiser l'état du réseau
- **Métriques système** : mémoire, CPU, goroutines
- **Métriques RETE** : nœuds actifs, faits traités, latence
- **WebSocket** pour mises à jour en temps réel
- **API REST** complète pour intégration

## 📦 Installation

### Prérequis
- **Go 1.21+**
- **pigeon** (générateur PEG) : `go install github.com/mna/pigeon@latest`

### Installation
```bash
git clone https://github.com/treivax/tsd.git
cd tsd
go mod tidy
go build ./...
```

## 🚀 Démarrage Rapide

### 1. Définir des Contraintes

Créez un fichier `rules.constraint` :

```constraint
// Définition des types métier
type Customer : <id: string, age: number, vip: bool>
type Order : <id: string, customer_id: string, total: number>
type Transaction : <id: string, amount: number, status: string>

// Règles métier avec actions obligatoires
{c: Customer} / c.age >= 18 AND c.vip == true ==> apply_vip_benefits(c.id)

{o: Order} / o.total > 1000 ==> flag_large_order(o.id, o.total)

{c: Customer, o: Order} / c.id == o.customer_id AND o.total > 500 ==> process_order(c.id, o.id)

{t: Transaction} / t.amount > 10000 AND t.status == "pending" ==> require_approval(t.id)

// Règles complexes avec négation et quantification
{c: Customer} / NOT (c.age < 18) AND EXISTS (o: Order / o.customer_id == c.id AND o.total > 100) ==> activate_premium_account(c.id)
```

### 2. Lancer le Monitoring

```bash
cd rete
go run cmd/monitoring/main.go
```

Accédez à l'interface web : **http://localhost:8082**

### 3. Intégration Programmatique

```go
package main

import (
    "fmt"
    "github.com/treivax/tsd/constraint"
    "github.com/treivax/tsd/rete"
)

func main() {
    // Parser les contraintes
    content, _ := os.ReadFile("rules.constraint")
    ast, err := constraint.Parse("rules.constraint", content)
    if err != nil {
        log.Fatal("Erreur parsing:", err)
    }

    // Créer le réseau RETE
    storage := rete.NewMemoryStorage()
    network := rete.NewReteNetwork(storage)

    // Convertir et charger les règles
    converter := rete.NewASTConverter()
    expressions, _ := converter.ConvertProgram(ast)

    for _, expr := range expressions {
        network.AddRule(expr)
    }

    // Ajouter des faits
    customerFact := rete.NewFact("Customer", map[string]interface{}{
        "id": "C001",
        "age": 25,
        "vip": true,
    })
    network.AddFact(customerFact)

    fmt.Println("✅ Réseau RETE configuré et opérationnel")
}
```

## 📊 Architecture

```
tsd/
├── constraint/              # Module de traitement des contraintes
│   ├── pkg/                 # Packages internes
│   │   ├── domain/          # Types fondamentaux et erreurs
│   │   └── validator/       # Validation et vérification
│   ├── grammar/             # Grammaire PEG et parser
│   ├── docs/                # Documentation utilisateur
│   └── test/                # Tests d'intégration
│
├── rete/                    # Module réseau RETE
│   ├── pkg/                 # Packages internes
│   │   ├── domain/          # Types et interfaces RETE
│   │   ├── nodes/           # Implémentation des nœuds
│   │   └── network/         # Logique réseau et constructeurs
│   ├── cmd/                 # Commandes exécutables
│   │   └── monitoring/      # Serveur de monitoring
│   ├── assets/web/          # Interface web de monitoring
│   └── test/                # Tests unitaires et intégration
│
└── tests/                   # Tests système globaux
```

## 🔧 Syntaxe des Contraintes

### Types de Base
```constraint
type TypeName : <field1: string, field2: number, field3: bool>
```

### Règles Simples
```constraint
{variable: TypeName} / condition ==> action(args)
```

### Règles Complexes
```constraint
// Jointures
{a: TypeA, b: TypeB} / a.id == b.ref_id ==> process_link(a.id, b.id)

// Négation
{user: User} / NOT (user.status == "banned") ==> allow_access(user.id)

// Quantification existentielle
{account: Account} / EXISTS (tx: Transaction / tx.account_id == account.id AND tx.amount > 1000) ==> flag_high_activity(account.id)

// Agrégation
{portfolio: Portfolio, asset: Asset} / portfolio.id == asset.portfolio_id AND SUM(asset.value) > 100000 ==> apply_portfolio_tax(portfolio.id)
```

### Opérateurs Supportés
- **Comparaison** : `==`, `!=`, `<`, `>`, `<=`, `>=`, `IN`, `LIKE`, `CONTAINS`
- **Logiques** : `AND`, `OR`, `NOT`
- **Agrégation** : `SUM()`, `COUNT()`, `AVG()`, `MIN()`, `MAX()`
- **Fonctions** : `LENGTH()`, `UPPER()`, `LOWER()`, `ABS()`, `ROUND()`

## 📈 Monitoring et Performance

### Interface Web
- **Temps réel** : Graphiques mis à jour automatiquement
- **Métriques système** : Usage mémoire, CPU, goroutines
- **Métriques RETE** : Nœuds actifs, débit de faits, latence moyenne
- **API REST** : Accès programmatique aux métriques

### Optimisations Performance
- **Cache d'évaluation** pour conditions complexes
- **Jointures par hash** optimisées
- **Propagation de tokens** asynchrone
- **Stockage indexé** pour recherche rapide

## 🧪 Tests et Validation

### Tests Complets
```bash
# Tests unitaires tous modules
go test ./...

# Tests d'intégration système
go test ./tests/... -v

# Tests de cohérence PEG ↔ RETE
go test ./tests/rete_coherence_test.go -v
```

### Cas d'Usage Testés
- **Domaine financier** : Détection de fraude, évaluation de risque
- **E-commerce** : Gestion commandes, promotion automatique
- **Ressources humaines** : Validation accès, calcul permissions
- **Banking** : Anti-blanchiment, conformité réglementaire

## 🔄 Développement

### Régénération du Parser
```bash
cd constraint
pigeon -o parser.go grammar/constraint.peg
```

### Build Scripts
```bash
# Module constraint
cd constraint && ./scripts/build.sh

# Module rete
cd rete && ./scripts/run_tests.sh

# Nettoyage global
./scripts/clean.sh
```

## 📚 Documentation Avancée

- **Guide des Contraintes** : `constraint/docs/GUIDE_CONTRAINTES.md`
- **Tutoriel Utilisateur** : `constraint/docs/TUTORIEL_CONTRAINTES.md`
- **Grammaire Complète** : `constraint/docs/GRAMMAR_COMPLETE.md`
- **Guide d'Usage Nœuds** : `rete/docs/ADVANCED_NODES_USAGE_GUIDE.md`

## 🤝 Contribution

1. Fork le projet
2. Créer une branche feature (`git checkout -b feature/amazing-feature`)
3. Commit les changements (`git commit -m 'Add amazing feature'`)
4. Push sur la branche (`git push origin feature/amazing-feature`)
5. Ouvrir une Pull Request

## 📄 Licence

Ce projet est sous licence MIT. Voir le fichier `LICENSE` pour plus de détails.

## 🏆 Status du Projet

✅ **Production Ready**
✅ **Tests Complets**
✅ **Documentation Complète**
✅ **Performance Optimisée**
✅ **Monitoring Intégré**

---

**Développé avec ❤️ pour des systèmes de règles métier haute performance**
