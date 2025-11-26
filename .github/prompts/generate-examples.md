# 🎯 Générer des Exemples RETE (Generate Examples)

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Tu veux créer des exemples de contraintes RETE (fichiers `.constraint` et `.facts`) pour démontrer des fonctionnalités, tester des patterns, documenter des cas d'usage, ou valider le comportement du moteur.

## Objectif

Générer des fichiers d'exemple `.constraint` et `.facts` bien structurés, commentés, et testables qui illustrent clairement des fonctionnalités ou des patterns spécifiques du moteur RETE.

## Types d'Exemples

### 1. **Exemples Pédagogiques**
- Démonstration de fonctionnalités basiques
- Introduction aux concepts RETE
- Tutoriels progressifs (simple → complexe)

### 2. **Exemples de Fonctionnalités**
- Opérateurs spécifiques (comparaison, chaînes, etc.)
- Types de nœuds (Alpha, Beta, Join)
- Patterns de règles

### 3. **Exemples de Cas d'Usage**
- Scénarios réels (business rules)
- Cas métier concrets
- Démonstrations pratiques

### 4. **Exemples de Tests**
- Validation de comportements
- Tests de régression
- Benchmarks

### 5. **Exemples de Documentation**
- README et guides
- Documentation API
- Exemples pour prompts réutilisables

## Instructions

### PHASE 1 : DÉFINITION (Qu'est-ce qu'on Veut Montrer)

#### 1.1 Identifier l'Objectif de l'Exemple

**Questions à se poser** :
- Quelle fonctionnalité démontrer ? (opérateurs, jointures, propagation, etc.)
- Quel niveau de complexité ? (débutant, intermédiaire, avancé)
- Pour quelle audience ? (utilisateurs, développeurs, documentation)
- Quel est le message principal ?

**Exemples d'objectifs** :
```
✨ Démontrer l'opérateur "startsWith" sur les chaînes
✨ Illustrer une jointure 3-way avec conditions multiples
✨ Montrer la propagation incrémentale avec ajout de faits
✨ Valider le comportement des AlphaNodes avec types multiples
✨ Créer un exemple de règle métier (détection de fraude)
```

#### 1.2 Définir le Scope et la Complexité

**Niveaux de complexité** :

**Niveau 1 - Débutant** :
- 1 type de fait
- 1-2 règles simples
- Opérateurs basiques
- Pas de jointures
- Résultats évidents

**Niveau 2 - Intermédiaire** :
- 2-3 types de faits
- 3-5 règles
- Quelques jointures
- Opérateurs variés
- Propagation simple

**Niveau 3 - Avancé** :
- 3+ types de faits
- 5+ règles complexes
- Jointures multiples (3-way+)
- Conditions complexes
- Propagation incrémentale

#### 1.3 Choisir le Domaine Métier

**Domaines courants** :
- 👥 **Personnes/Utilisateurs** : Gestion d'identités, permissions
- 🛒 **E-commerce** : Commandes, produits, inventaire
- 💼 **RH** : Employés, départements, salaires
- 🏦 **Finance** : Transactions, comptes, détection fraude
- 🚗 **IoT** : Capteurs, événements, alertes
- 📚 **Éducation** : Étudiants, cours, notes

**Choisir un domaine familier et compréhensible** pour que l'exemple soit accessible.

### PHASE 2 : CONCEPTION (Structure de l'Exemple)

#### 2.1 Concevoir le Modèle de Données

**Définir les types de faits** :

```
Exemple : E-commerce

Types de faits :
1. Customer (id, name, email, status, totalSpent)
2. Order (id, customerId, amount, date, status)
3. Product (id, name, price, category, stock)
4. OrderItem (orderId, productId, quantity)
```

**Règles de conception** :
- Types clairs et focalisés
- Relations explicites (IDs pour jointures)
- Propriétés pertinentes pour les règles
- Valeurs réalistes et variées

#### 2.2 Concevoir les Règles

**Pattern de règle** :
```
Nom descriptif : Ce que la règle détecte/fait

Conditions :
- Quels faits sont nécessaires ?
- Quelles propriétés tester ?
- Quelles jointures effectuer ?

Action (ce qui est produit) :
- Quel token terminal ?
- Quelles variables inclure ?
```

**Exemple de conception** :
```
Règle : "Clients VIP"
Objectif : Identifier les clients ayant dépensé plus de 10000

Conditions :
- {c: Customer} avec c.totalSpent > 10000

Action :
- ==> vipCustomer(c)

Résultat attendu :
- 1 token si 1 client satisfait
- Variables liées : c (Customer complet)
```

#### 2.3 Concevoir les Données de Test

**Principes** :
- **Cas positifs** : Données qui satisfont les règles
- **Cas négatifs** : Données qui ne satisfont PAS (tester les limites)
- **Cas limites** : Valeurs aux frontières (égalité, vide, etc.)
- **Diversité** : Différents chemins d'exécution

**Template** :
```json
{
  "facts": [
    // CAS POSITIF : Devrait matcher
    {
      "type": "Customer",
      "data": {"id": 1, "totalSpent": 15000} // > 10000 ✓
    },
    
    // CAS NÉGATIF : Ne devrait PAS matcher
    {
      "type": "Customer",
      "data": {"id": 2, "totalSpent": 5000}  // < 10000 ✗
    },
    
    // CAS LIMITE : Frontière
    {
      "type": "Customer",
      "data": {"id": 3, "totalSpent": 10000} // = 10000 ✗ (car >)
    }
  ]
}
```

### PHASE 3 : IMPLÉMENTATION (Écrire les Fichiers)

#### 3.1 Créer le Fichier `.constraint`

**Structure recommandée** :
```constraint
# =============================================================================
# TITRE DE L'EXEMPLE
# =============================================================================
#
# DESCRIPTION :
# Description détaillée de ce que cet exemple démontre, pourquoi il est utile,
# et ce qu'on peut apprendre en l'étudiant.
#
# OBJECTIFS D'APPRENTISSAGE :
# - Objectif 1 : Comprendre X
# - Objectif 2 : Maîtriser Y
# - Objectif 3 : Voir Z en action
#
# NIVEAU : [Débutant / Intermédiaire / Avancé]
#
# DOMAINE : [E-commerce / Finance / etc.]
#
# =============================================================================

# -----------------------------------------------------------------------------
# RÈGLE 1 : Nom Descriptif
# -----------------------------------------------------------------------------
# Description : Ce que cette règle fait et pourquoi
#
# Conditions :
#   - Condition 1 : Explication
#   - Condition 2 : Explication
#
# Action : Ce qui est produit
#
# Résultats attendus avec les données de test :
#   - Nombre de tokens : X
#   - Variables liées : v1, v2
# -----------------------------------------------------------------------------

{c: Customer} / c.totalSpent > 10000 ==> vipCustomer(c)

# -----------------------------------------------------------------------------
# RÈGLE 2 : Autre Règle
# -----------------------------------------------------------------------------
# [Même structure de commentaires]
# -----------------------------------------------------------------------------

{o: Order}, {c: Customer} /
    o.customerId == c.id,
    o.amount > 1000
==> largeOrder(o, c)

# =============================================================================
# RÉSULTATS ATTENDUS (RÉSUMÉ)
# =============================================================================
#
# Avec les données dans le fichier .facts associé :
#
# 1. vipCustomer(c) :
#    - 2 tokens attendus
#    - Customers : Alice (id=1), Charlie (id=3)
#
# 2. largeOrder(o, c) :
#    - 3 tokens attendus
#    - Orders : 101, 103, 105
#
# TOTAL : 5 tokens terminaux
#
# =============================================================================
```

**Bonnes pratiques** :
- ✅ Commentaires abondants et structurés
- ✅ En-tête explicite avec contexte
- ✅ Chaque règle documentée
- ✅ Résultats attendus clairement indiqués
- ✅ Niveau et domaine mentionnés
- ✅ Séparateurs visuels (lignes de =, -)

#### 3.2 Créer le Fichier `.facts`

**Structure JSON** :
```json
{
  "facts": [
    {
      "type": "TypeName",
      "data": {
        "property1": "value1",
        "property2": 123
      }
    }
  ]
}
```

**Exemple complet avec commentaires** :
```json
{
  "description": "Données de test pour l'exemple E-commerce VIP Customers",
  "version": "1.0",
  "author": "TSD Team",
  "date": "2025-11-26",
  
  "facts": [
    {
      "comment": "CAS POSITIF - Alice est VIP (totalSpent > 10000)",
      "type": "Customer",
      "data": {
        "id": 1,
        "name": "Alice Martin",
        "email": "alice@example.com",
        "status": "active",
        "totalSpent": 15000
      }
    },
    
    {
      "comment": "CAS NÉGATIF - Bob n'est pas VIP (totalSpent < 10000)",
      "type": "Customer",
      "data": {
        "id": 2,
        "name": "Bob Smith",
        "email": "bob@example.com",
        "status": "active",
        "totalSpent": 5000
      }
    },
    
    {
      "comment": "CAS LIMITE - Charlie est VIP (totalSpent = 15000)",
      "type": "Customer",
      "data": {
        "id": 3,
        "name": "Charlie Brown",
        "email": "charlie@example.com",
        "status": "active",
        "totalSpent": 15000
      }
    },
    
    {
      "comment": "Large order de Alice (devrait matcher largeOrder)",
      "type": "Order",
      "data": {
        "id": 101,
        "customerId": 1,
        "amount": 2500,
        "date": "2025-11-26",
        "status": "completed"
      }
    },
    
    {
      "comment": "Small order de Bob (ne devrait PAS matcher largeOrder)",
      "type": "Order",
      "data": {
        "id": 102,
        "customerId": 2,
        "amount": 500,
        "date": "2025-11-27",
        "status": "completed"
      }
    },
    
    {
      "comment": "Large order de Charlie",
      "type": "Order",
      "data": {
        "id": 103,
        "customerId": 3,
        "amount": 3000,
        "date": "2025-11-28",
        "status": "pending"
      }
    }
  ],
  
  "expectedResults": {
    "vipCustomer": {
      "count": 2,
      "customers": ["Alice Martin", "Charlie Brown"]
    },
    "largeOrder": {
      "count": 2,
      "orders": [101, 103]
    }
  }
}
```

**Bonnes pratiques** :
- ✅ Métadonnées en en-tête (description, version, date)
- ✅ Commentaires pour chaque fait (cas positif/négatif/limite)
- ✅ Données réalistes et variées
- ✅ Résultats attendus en fin de fichier
- ✅ Nommage cohérent avec le domaine
- ✅ Valeurs aux frontières testées

#### 3.3 Créer la Documentation Associée (Optionnel)

**Fichier : `exemple_name.md`**

```markdown
# Exemple : VIP Customers et Large Orders

## 📋 Vue d'Ensemble

**Niveau** : Intermédiaire  
**Domaine** : E-commerce  
**Fichiers** :
- Contraintes : [`vip_customers.constraint`](vip_customers.constraint)
- Données : [`vip_customers.facts`](vip_customers.facts)

**Objectif** : Démontrer les jointures simples et la détection de patterns
basée sur des seuils numériques.

## 🎯 Ce que Vous Allez Apprendre

1. **Filtrage par seuils** : Utiliser les opérateurs de comparaison (`>`)
2. **Jointures** : Lier des faits via des propriétés communes
3. **Tokens multiples** : Comprendre quand plusieurs tokens sont générés
4. **Variables liées** : Propager les données à travers les jointures

## 📊 Modèle de Données

### Types de Faits

```
Customer
├── id: number (identifiant unique)
├── name: string (nom complet)
├── email: string (email)
├── status: string (active/inactive)
└── totalSpent: number (montant total dépensé)

Order
├── id: number (identifiant unique)
├── customerId: number (référence vers Customer)
├── amount: number (montant de la commande)
├── date: string (date ISO)
└── status: string (completed/pending/cancelled)
```

### Relations

```
Customer --< Order (1:N via customerId)
```

## 📝 Règles Implémentées

### Règle 1 : VIP Customers

**Objectif** : Identifier les clients ayant dépensé plus de 10 000€

**Contrainte** :
```constraint
{c: Customer} / c.totalSpent > 10000 ==> vipCustomer(c)
```

**Explication** :
- **Pattern** : `{c: Customer}` - Capture chaque fait Customer
- **Condition** : `c.totalSpent > 10000` - Filtre sur le montant
- **Action** : `vipCustomer(c)` - Génère un token avec le customer

**Avec les données de test** :
- ✅ Alice (15000€) → Match
- ❌ Bob (5000€) → No match
- ✅ Charlie (15000€) → Match

**Résultat** : 2 tokens

### Règle 2 : Large Orders

**Objectif** : Identifier les commandes importantes (> 1000€) avec leur client

**Contrainte** :
```constraint
{o: Order}, {c: Customer} /
    o.customerId == c.id,
    o.amount > 1000
==> largeOrder(o, c)
```

**Explication** :
- **Patterns** : `{o: Order}, {c: Customer}` - Jointure entre Order et Customer
- **Conditions** :
  - `o.customerId == c.id` - Lie l'ordre au client
  - `o.amount > 1000` - Filtre sur le montant
- **Action** : `largeOrder(o, c)` - Token avec ordre ET client

**Avec les données de test** :
- ✅ Order 101 (2500€) + Alice → Match
- ❌ Order 102 (500€) + Bob → No match (amount)
- ✅ Order 103 (3000€) + Charlie → Match

**Résultat** : 2 tokens

## 🚀 Exécution

### Commande

```bash
make rete-run CONSTRAINT=docs/examples/vip_customers.constraint \
              FACTS=docs/examples/vip_customers.facts
```

### Résultats Attendus

```
✅ RÈGLE: vipCustomer
   Tokens générés : 2
   
   Token 1:
   - Customer: Alice Martin (id=1)
   - totalSpent: 15000
   
   Token 2:
   - Customer: Charlie Brown (id=3)
   - totalSpent: 15000

✅ RÈGLE: largeOrder
   Tokens générés : 2
   
   Token 1:
   - Order: 101 (amount=2500)
   - Customer: Alice Martin (id=1)
   
   Token 2:
   - Order: 103 (amount=3000)
   - Customer: Charlie Brown (id=3)

📊 TOTAL : 4 tokens terminaux
```

## 🔍 Analyse Détaillée

### Propagation dans le Réseau RETE

```
Customer Facts (3) → AlphaNode[Customer]
                     ↓
                     FilterNode[totalSpent > 10000]
                     ↓
                     TerminalNode[vipCustomer] → 2 tokens

Order Facts (3) → AlphaNode[Order]
                  ↓
Customer Facts → AlphaNode[Customer]
                  ↓
                  JoinNode[customerId == id]
                  ↓
                  FilterNode[amount > 1000]
                  ↓
                  TerminalNode[largeOrder] → 2 tokens
```

### Pourquoi Bob Ne Match Pas ?

**Pour vipCustomer** :
- ❌ `totalSpent = 5000` qui n'est PAS `> 10000`

**Pour largeOrder** :
- ❌ Son ordre (102) a `amount = 500` qui n'est PAS `> 1000`

### Cas Limites Testés

1. **Égalité vs Supériorité** :
   - Seuil : 10000
   - Charlie a exactement 15000 → Match (> 10000)
   - Si seuil était `>= 10000`, un client à 10000 matcherait

2. **Jointures avec données manquantes** :
   - Si un Order n'avait pas de Customer correspondant → Pas de match
   - Toutes les jointures nécessitent les deux faits

## 💡 Variations Possibles

### Ajouter une Règle

```constraint
# VIP avec large order
{c: Customer}, {o: Order} /
    c.totalSpent > 10000,
    o.customerId == c.id,
    o.amount > 1000
==> vipWithLargeOrder(c, o)
```

**Résultat attendu** : 2 tokens (Alice+101, Charlie+103)

### Modifier les Seuils

```constraint
# Super VIP (> 20000)
{c: Customer} / c.totalSpent > 20000 ==> superVip(c)
```

**Résultat attendu** : 0 tokens (aucun client > 20000)

### Ajouter des Conditions

```constraint
# VIP actifs seulement
{c: Customer} /
    c.totalSpent > 10000,
    c.status == "active"
==> activeVip(c)
```

## 🎓 Points Clés à Retenir

1. **Opérateurs de comparaison** : `>`, `<`, `>=`, `<=`, `==`, `!=`
2. **Jointures** : Lier faits via propriétés communes (ex: `id`)
3. **Filtrage** : Réduire les résultats avec conditions
4. **Tokens multiples** : Une règle peut générer N tokens si N faits matchent
5. **Variables liées** : Les tokens contiennent toutes les variables capturées

## 🔗 Voir Aussi

- [Guide des Opérateurs](../operators.md)
- [Architecture RETE](../architecture.md)
- [Exemples Avancés](advanced_examples.md)

## 📝 Exercices

1. **Modifier le seuil VIP à 12000** : Combien de tokens ?
2. **Ajouter une règle pour "small orders" (< 1000)** : Combien de tokens ?
3. **Créer une règle "inactive VIP"** : Clients VIP mais status = "inactive"

---

**Auteur** : TSD Team  
**Version** : 1.0  
**Date** : 2025-11-26
```

### PHASE 4 : VALIDATION (Tester l'Exemple)

#### 4.1 Exécuter l'Exemple

**Commandes de test** :
```bash
# Méthode 1 : Via Makefile (si défini)
make rete-run CONSTRAINT=docs/examples/vip_customers.constraint \
              FACTS=docs/examples/vip_customers.facts

# Méthode 2 : Via runner universel
./bin/rete-runner docs/examples/vip_customers.constraint \
                  docs/examples/vip_customers.facts

# Méthode 3 : Via tests Go
go test -v -run TestExample_VIPCustomers ./test/integration
```

#### 4.2 Vérifier les Résultats

**Checklist de validation** :
```
✅ Le programme s'exécute sans erreur
✅ Le nombre de tokens correspond aux attentes
✅ Les tokens contiennent les bonnes variables
✅ Les valeurs des variables sont correctes
✅ Les cas positifs matchent
✅ Les cas négatifs ne matchent PAS
✅ Les cas limites se comportent comme prévu
```

**Comparaison attendu vs réel** :
```
Attendu (dans .constraint ou .facts) :
- vipCustomer: 2 tokens (Alice, Charlie)
- largeOrder: 2 tokens (101, 103)

Réel (output du runner) :
- vipCustomer: 2 tokens ✅
  - Token 1: Alice ✅
  - Token 2: Charlie ✅
- largeOrder: 2 tokens ✅
  - Token 1: Order 101 + Alice ✅
  - Token 2: Order 103 + Charlie ✅

RÉSULTAT : ✅ Validation réussie
```

#### 4.3 Débugger si Nécessaire

**Problèmes courants** :

**Problème 1 : Pas de tokens générés**
```
Causes possibles :
- Syntaxe .constraint incorrecte
- Conditions trop restrictives
- Données .facts ne matchent pas les types
- Erreur dans les jointures (IDs incorrects)

Debug :
1. Vérifier la syntaxe avec le parseur
2. Simplifier les conditions une par une
3. Ajouter des logs dans le code
4. Vérifier les types de données
```

**Problème 2 : Trop de tokens générés**
```
Causes possibles :
- Conditions pas assez restrictives
- Jointures manquantes
- Données en doublon dans .facts

Debug :
1. Vérifier le nombre de faits
2. Ajouter des conditions pour filtrer
3. Vérifier l'unicité des IDs
```

**Problème 3 : Erreurs de parsing**
```
Causes possibles :
- Syntaxe .constraint invalide
- Caractères spéciaux non échappés
- JSON .facts malformé

Debug :
1. Valider JSON avec jsonlint
2. Vérifier syntaxe .constraint ligne par ligne
3. Consulter docs/grammar.md
```

### PHASE 5 : DOCUMENTATION (Finaliser)

#### 5.1 Organiser les Fichiers

**Structure recommandée** :
```
docs/examples/
├── README.md                      # Index des exemples
├── basic/
│   ├── simple_filter.constraint
│   ├── simple_filter.facts
│   └── simple_filter.md
├── intermediate/
│   ├── vip_customers.constraint
│   ├── vip_customers.facts
│   └── vip_customers.md
└── advanced/
    ├── fraud_detection.constraint
    ├── fraud_detection.facts
    └── fraud_detection.md
```

#### 5.2 Créer l'Index des Exemples

**Fichier : docs/examples/README.md**

```markdown
# Exemples RETE

Collection d'exemples progressifs pour apprendre et maîtriser le moteur RETE TSD.

## 📚 Par Niveau

### 🟢 Débutant

| Exemple | Description | Concepts |
|---------|-------------|----------|
| [Simple Filter](basic/simple_filter.md) | Filtrage basique | AlphaNode, Comparaisons |
| [Multiple Rules](basic/multiple_rules.md) | Plusieurs règles | TerminalNodes multiples |

### 🟡 Intermédiaire

| Exemple | Description | Concepts |
|---------|-------------|----------|
| [VIP Customers](intermediate/vip_customers.md) | Jointures simples | JoinNode, Seuils |
| [String Operators](intermediate/string_ops.md) | Opérateurs chaînes | startsWith, contains |

### 🔴 Avancé

| Exemple | Description | Concepts |
|---------|-------------|----------|
| [Fraud Detection](advanced/fraud_detection.md) | Détection fraude | Jointures 3-way, Agrégation |
| [Incremental](advanced/incremental.md) | Propagation incrémentale | Ajout dynamique |

## 🎯 Par Fonctionnalité

### Opérateurs de Comparaison
- [Numeric Comparisons](basic/numeric_comparisons.md) - `>`, `<`, `>=`, `<=`
- [Equality](basic/equality.md) - `==`, `!=`

### Opérateurs de Chaînes
- [String Operators](intermediate/string_ops.md) - `startsWith`, `endsWith`, `contains`

### Jointures
- [Simple Join](intermediate/simple_join.md) - Jointure 2 faits
- [Three-Way Join](advanced/three_way_join.md) - Jointure 3+ faits
- [Multiple Conditions](advanced/multi_conditions.md) - Conditions multiples

### Patterns Avancés
- [Negation](advanced/negation.md) - Conditions négatives
- [Aggregation](advanced/aggregation.md) - Compteurs, sommes
- [Temporal](advanced/temporal.md) - Fenêtres temporelles

## 🏢 Par Domaine Métier

- 🛒 **E-commerce** : [VIP Customers](intermediate/vip_customers.md)
- 🏦 **Finance** : [Fraud Detection](advanced/fraud_detection.md)
- 👥 **RH** : [Employee Management](intermediate/employees.md)
- 🚗 **IoT** : [Sensor Alerts](advanced/sensors.md)

## 🚀 Démarrage Rapide

### Exécuter un Exemple

```bash
# Via Makefile
make rete-run CONSTRAINT=docs/examples/basic/simple_filter.constraint \
              FACTS=docs/examples/basic/simple_filter.facts

# Via runner
./bin/rete-runner docs/examples/basic/simple_filter.constraint \
                  docs/examples/basic/simple_filter.facts
```

### Créer Votre Propre Exemple

1. Copier un exemple similaire
2. Modifier les règles et données
3. Tester avec le runner
4. Ajuster jusqu'à obtenir les résultats souhaités

## 📖 Guide de Lecture

**Parcours Débutant** :
1. [Simple Filter](basic/simple_filter.md)
2. [Multiple Rules](basic/multiple_rules.md)
3. [Simple Join](intermediate/simple_join.md)

**Parcours Développeur** :
1. [VIP Customers](intermediate/vip_customers.md)
2. [String Operators](intermediate/string_ops.md)
3. [Three-Way Join](advanced/three_way_join.md)

**Parcours Architecte** :
1. [Incremental Propagation](advanced/incremental.md)
2. [Fraud Detection](advanced/fraud_detection.md)
3. [Performance Benchmarks](advanced/benchmarks.md)

## 🤝 Contribuer

Pour ajouter un nouvel exemple :

1. Choisir catégorie (basic/intermediate/advanced)
2. Créer les fichiers `.constraint`, `.facts`, et `.md`
3. Suivre les templates dans [generate-examples.md](../../.github/prompts/generate-examples.md)
4. Tester l'exemple
5. Mettre à jour cet index

---

**Total d'exemples** : 15  
**Dernière mise à jour** : 2025-11-26
```

#### 5.3 Commit et Publication

**Commits séparés** :
```bash
# Fichier contrainte
git add docs/examples/vip_customers.constraint
git commit -m "docs(examples): add VIP customers constraint example"

# Fichier facts
git add docs/examples/vip_customers.facts
git commit -m "docs(examples): add VIP customers test data"

# Documentation
git add docs/examples/vip_customers.md
git commit -m "docs(examples): add VIP customers documentation"

# Mise à jour index
git add docs/examples/README.md
git commit -m "docs(examples): update index with VIP customers example"
```

## Critères de Succès

### ✅ Qualité de l'Exemple

- [ ] Objectif clairement défini et atteint
- [ ] Commentaires abondants et utiles
- [ ] Structure logique et progressive
- [ ] Cas positifs, négatifs, et limites couverts
- [ ] Données réalistes et pertinentes

### ✅ Fonctionnalité

- [ ] L'exemple s'exécute sans erreur
- [ ] Les résultats correspondent aux attentes
- [ ] Tous les cas de test passent
- [ ] La syntaxe est correcte (.constraint et .facts)

### ✅ Documentation

- [ ] En-tête descriptif complet
- [ ] Chaque règle documentée
- [ ] Résultats attendus indiqués
- [ ] Fichier .md complet (si créé)
- [ ] Index mis à jour

### ✅ Utilité Pédagogique

- [ ] L'exemple enseigne clairement le concept visé
- [ ] Niveau de difficulté approprié
- [ ] Progression logique si série d'exemples
- [ ] Utilisable dans documentation/formation

## Format de Réponse

```markdown
# 🎯 EXEMPLE RETE GÉNÉRÉ

## 📋 Résumé

**Nom** : VIP Customers et Large Orders  
**Niveau** : Intermédiaire  
**Domaine** : E-commerce  
**Objectif** : Démontrer jointures simples et filtrage par seuils

## 📁 Fichiers Créés

1. **vip_customers.constraint** (125 lignes)
   - 2 règles commentées
   - En-tête descriptif
   - Résultats attendus documentés

2. **vip_customers.facts** (85 lignes JSON)
   - 3 Customers (2 VIP, 1 non-VIP)
   - 3 Orders (2 large, 1 small)
   - Cas positifs, négatifs, limites

3. **vip_customers.md** (250 lignes)
   - Guide complet
   - Explications détaillées
   - Diagrammes de propagation
   - Exercices suggérés

## 🎯 Concepts Démontrés

1. **Filtrage par seuils** : `c.totalSpent > 10000`
2. **Jointures** : `o.customerId == c.id`
3. **Tokens multiples** : 2 tokens pour vipCustomer
4. **Variables liées** : Propagation de `c` et `o`

## 📊 Modèle de Données

### Types
- `Customer` : id, name, email, status, totalSpent
- `Order` : id, customerId, amount, date, status

### Relations
- Customer --< Order (1:N via customerId)

## 📝 Règles Implémentées

### Règle 1 : VIP Customers
```constraint
{c: Customer} / c.totalSpent > 10000 ==> vipCustomer(c)
```
**Résultat** : 2 tokens (Alice, Charlie)

### Règle 2 : Large Orders
```constraint
{o: Order}, {c: Customer} /
    o.customerId == c.id,
    o.amount > 1000
==> largeOrder(o, c)
```
**Résultat** : 2 tokens (Order 101+Alice, Order 103+Charlie)

## 🧪 Données de Test

### Cas Positifs (Doivent Matcher)
- Alice : totalSpent = 15000 ✅ VIP
- Charlie : totalSpent = 15000 ✅ VIP
- Order 101 : amount = 2500 ✅ Large
- Order 103 : amount = 3000 ✅ Large

### Cas Négatifs (Ne Doivent PAS Matcher)
- Bob : totalSpent = 5000 ❌ Non-VIP
- Order 102 : amount = 500 ❌ Small

### Cas Limites
- Seuil VIP : > 10000 (égalité exclue)
- Seuil Large : > 1000 (égalité exclue)

## ✅ Validation

### Exécution
```bash
$ make rete-run CONSTRAINT=docs/examples/vip_customers.constraint \
                FACTS=docs/examples/vip_customers.facts

✅ Parsing réussi
✅ Réseau RETE construit
✅ Faits soumis
✅ Propagation terminée
```

### Résultats Attendus vs Réels

| Règle | Attendu | Réel | Status |
|-------|---------|------|--------|
| vipCustomer | 2 tokens | 2 tokens | ✅ |
| largeOrder | 2 tokens | 2 tokens | ✅ |
| **TOTAL** | **4 tokens** | **4 tokens** | ✅ |

### Détail des Tokens

**vipCustomer** :
- ✅ Token 1 : Customer id=1 (Alice, 15000)
- ✅ Token 2 : Customer id=3 (Charlie, 15000)

**largeOrder** :
- ✅ Token 1 : Order 101 (2500) + Customer 1 (Alice)
- ✅ Token 2 : Order 103 (3000) + Customer 3 (Charlie)

## 📈 Structure du Réseau RETE

```
Customer Facts (3)
  ↓
AlphaNode[Customer]
  ↓
  ├─→ FilterNode[totalSpent > 10000]
  │     ↓
  │   TerminalNode[vipCustomer] → 2 tokens
  │
  └─→ JoinNode (avec Order)
        ↓
      FilterNode[amount > 1000]
        ↓
      TerminalNode[largeOrder] → 2 tokens

Order Facts (3)
  ↓
AlphaNode[Order]
  ↓
JoinNode[customerId == id]
```

## 💡 Points Clés

### Ce que Démontre Cet Exemple

1. **Filtrage simple** : Une condition sur une propriété numérique
2. **Jointure 2-way** : Lier deux types de faits
3. **Conditions multiples** : Combiner jointure ET filtre
4. **Cas négatifs** : Importance des données qui ne matchent pas

### Ce qu'On Peut Apprendre

- Comment écrire des règles avec seuils
- Comment joindre des faits via propriétés communes
- Comment un token contient plusieurs variables
- Pourquoi certains faits ne produisent pas de tokens

### Variations Possibles

1. **Ajouter une règle** : VIP avec large order (intersection)
2. **Modifier seuils** : Tester différentes valeurs
3. **Ajouter statuts** : Filtrer sur Customer.status
4. **Dates** : Filtrer Orders par période

## 📚 Utilisation

### Dans la Documentation
- README.md : Exemple de jointure simple
- Guide opérateurs : Illustration de `>`
- Tutoriel RETE : Progression niveau intermédiaire

### Dans les Tests
- Tests d'intégration : Valider jointures
- Tests de régression : Garantir stabilité
- Benchmarks : Mesurer performance jointures

### Pour la Formation
- Atelier débutant : Introduction aux jointures
- TP étudiant : Modifier et expérimenter
- Démonstration : Visualiser propagation

## 🔗 Liens

**Fichiers** :
- [vip_customers.constraint](docs/examples/vip_customers.constraint)
- [vip_customers.facts](docs/examples/vip_customers.facts)
- [vip_customers.md](docs/examples/vip_customers.md)

**Documentation** :
- [Index des Exemples](docs/examples/README.md)
- [Guide des Opérateurs](docs/operators.md)
- [Architecture RETE](docs/architecture.md)

## ✅ Checklist Complétée

- [x] Objectif clairement défini
- [x] Fichier .constraint avec commentaires abondants
- [x] Fichier .facts avec cas positifs/négatifs/limites
- [x] Documentation .md complète
- [x] Exemple testé et validé
- [x] Résultats correspondent aux attentes
- [x] Index mis à jour
- [x] Commits effectués
```

## Exemple d'Utilisation

```
Je veux créer un exemple RETE pour démontrer l'utilisation des opérateurs
de chaînes (startsWith, endsWith, contains) dans un contexte de filtrage
d'utilisateurs par email et nom.

Niveau : Intermédiaire
Domaine : Gestion d'utilisateurs
Objectif pédagogique : Montrer comment combiner plusieurs opérateurs de chaînes

Utilise le prompt "generate-examples" pour créer les fichiers .constraint,
.facts, et la documentation associée.
```

## Checklist de Génération

### Avant de Commencer

- [ ] Objectif de l'exemple clairement défini
- [ ] Niveau de complexité choisi (débutant/intermédiaire/avancé)
- [ ] Domaine métier identifié
- [ ] Audience cible connue (utilisateurs/développeurs/docs)
- [ ] Concepts à démontrer listés

### Pendant la Création

- [ ] Modèle de données conçu (types, propriétés)
- [ ] Règles conçues (conditions, actions)
- [ ] Données de test planifiées (positif/négatif/limite)
- [ ] Commentaires ajoutés au fur et à mesure
- [ ] Structure logique et progressive

### Fichier .constraint

- [ ] En-tête descriptif complet
- [ ] Chaque règle documentée
- [ ] Commentaires abondants
- [ ] Résultats attendus indiqués
- [ ] Syntaxe correcte (testée)
- [ ] Séparateurs visuels (lisibilité)

### Fichier .facts

- [ ] JSON valide (testé avec jsonlint)
- [ ] Métadonnées en en-tête
- [ ] Commentaires pour chaque fait
- [ ] Cas positifs inclus
- [ ] Cas négatifs inclus
- [ ] Cas limites testés
- [ ] Résultats attendus documentés
- [ ] Données réalistes

### Documentation .md (si créée)

- [ ] Vue d'ensemble claire
- [ ] Objectifs d'apprentissage listés
- [ ] Modèle de données expliqué
- [ ] Chaque règle détaillée
- [ ] Instructions d'exécution
- [ ] Résultats attendus
- [ ] Analyse détaillée
- [ ] Variations suggérées
- [ ] Exercices proposés

### Validation

- [ ] Exemple exécuté sans erreur
- [ ] Résultats correspondent aux attentes
- [ ] Nombre de tokens correct
- [ ] Variables liées correctes
- [ ] Cas positifs matchent
- [ ] Cas négatifs ne matchent PAS
- [ ] Cas limites se comportent comme prévu

### Finalisation

- [ ] Fichiers organisés (dossier approprié)
- [ ] Index mis à jour
- [ ] Commits effectués
- [ ] Documentation projet mise à jour si nécessaire

## Commandes Utiles

```bash
# Valider JSON
jsonlint docs/examples/example.facts

# Exécuter exemple
make rete-run CONSTRAINT=docs/examples/example.constraint \
              FACTS=docs/examples/example.facts

# Tester avec runner
./bin/rete-runner docs/examples/example.constraint \
                  docs/examples/example.facts

# Créer test d'intégration
go test -v -run TestExample_Name ./test/integration

# Lister tous les exemples
find docs/examples -name "*.constraint"

# Vérifier syntaxe .constraint (si outil disponible)
./bin/constraint-validator docs/examples/example.constraint

# Compter tokens attendus vs réels
grep "==> " docs/examples/example.constraint | wc -l

# Formater JSON
jq '.' docs/examples/example.facts > /tmp/formatted.json
mv /tmp/formatted.json docs/examples/example.facts
```

## Bonnes Pratiques

### Conception

- **Focalisé** : Un exemple = un ou deux concepts, pas tout à la fois
- **Progressif** : Complexité croissante dans une série d'exemples
- **Réaliste** : Domaine familier et cas d'usage crédibles
- **Complet** : Cas positifs ET négatifs ET limites

### Commentaires

- **Abondants** : Ne jamais sous-estimer l'importance des commentaires
- **Structurés** : Utiliser sections, séparateurs, hiérarchie
- **Explicatifs** : Expliquer le "pourquoi", pas juste le "quoi"
- **Résultats** : Toujours indiquer ce qui est attendu

### Données

- **Variées** : Différents cas pour tester tous les chemins
- **Réalistes** : Valeurs crédibles, noms cohérents
- **Nommées** : Utiliser vrais noms, pas "User1", "User2"
- **Documentées** : Commenter chaque fait (pourquoi il est là)

### Documentation

- **Complète** : Vue d'ensemble, détails, variations, exercices
- **Visuelle** : Diagrammes si utile (propagation RETE)
- **Pédagogique** : Focus sur l'apprentissage
- **Liens** : Connecter aux autres ressources

## Anti-Patterns à Éviter

### ❌ Exemple Trop Complexe
```
❌ Mélanger 10 concepts dans un seul exemple
✅ Un exemple = 1-2 concepts clés
```

### ❌ Pas de Commentaires
```
❌ Code brut sans explication
✅ Commentaires abondants partout
```

### ❌ Données Artificielles
```
❌ {id: 1, name: "User1", value: 99999}
✅ {id: 1, name: "Alice Martin", totalSpent: 15000}
```

### ❌ Que des Cas Positifs
```
❌ Tous les faits matchent toutes les règles
✅ Mix de cas positifs, négatifs, limites
```

### ❌ Résultats Non Documentés
```
❌ Exécuter et espérer que ça marche
✅ Documenter précisément les résultats attendus
```

### ❌ Pas de Test
```
❌ Créer l'exemple sans l'exécuter
✅ Toujours tester avant de committer
```

## Templates

### Template .constraint Minimal

```constraint
# =============================================================================
# [TITRE]
# =============================================================================
#
# DESCRIPTION : [Ce que démontre cet exemple]
# NIVEAU : [Débutant/Intermédiaire/Avancé]
# DOMAINE : [E-commerce/Finance/etc.]
#
# =============================================================================

# -----------------------------------------------------------------------------
# RÈGLE : [Nom]
# Description : [Explication]
# Résultat attendu : [X tokens]
# -----------------------------------------------------------------------------

[règle ici]

# =============================================================================
# RÉSULTATS ATTENDUS : [X tokens au total]
# =============================================================================
```

### Template .facts Minimal

```json
{
  "description": "[Description de cet ensemble de données]",
  "facts": [
    {
      "comment": "[CAS POSITIF/NÉGATIF/LIMITE] - [Explication]",
      "type": "TypeName",
      "data": {
        "property": "value"
      }
    }
  ],
  "expectedResults": {
    "ruleName": {
      "count": 0,
      "details": "Description"
    }
  }
}
```

## Ressources

- [Makefile](../../Makefile) - Commandes disponibles
- [Grammaire PEG](../../constraint/grammar.peg) - Syntaxe .constraint
- [Documentation RETE](../../docs/) - Architecture et concepts
- [Exemples Existants](../../docs/examples/) - S'inspirer
- [Tests](../../test/) - Exemples de tests

---

**Version** : 1.0  
**Dernière mise à jour** : Novembre 2025  
**Mainteneur** : Équipe TSD