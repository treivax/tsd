# 📊 RAPPORT COMPLET - TESTS DE COUVERTURE ALPHA NODES

**Date d'exécution:** 2025-11-18 12:01:26
**Nombre de tests:** 26

## 🎯 RÉSUMÉ EXÉCUTIF

- ✅ **Tests réussis:** 26/26 (100.0%)
- 🎬 **Actions déclenchées:** 26
- ⚡ **Couverture:** Nœuds Alpha positifs et négatifs

## 🧪 TEST 1: alpha_abs_negative

### 📋 Informations générales

- **Description:** Test fonction ABS() (valeur absolue) - négation
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_abs_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_abs_negative.facts`
- **Temps d'exécution:** 742.808µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1 🚫:**
```constraint
{b: Balance} / NOT(ABS(b.amount) > 100) ==> small_balance_found(b.id, b.amount)
```
- **Action:** `small_balance_found`
- **Condition:** `NOT(map[args:[map[field:amount object:b type:fieldAccess]] name:ABS type:functionCall] > 100)`
- **Type:** Condition négative (NOT)

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `B001`
```json
Type: Balance
Champs:
  type: credit
  id: B001
  amount: 150
```

**Fait 2:** `B002`
```json
Type: Balance
Champs:
  id: B002
  amount: -25
  type: debit
```

**Fait 3:** `B003`
```json
Type: Balance
Champs:
  id: B003
  amount: 75
  type: credit
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Balance (type_Balance)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: NOT(...) [Négation]
│   │   └── Variable: b
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: small_balance_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `small_balance_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `B002` (Type: Balance)
  2. `B003` (Type: Balance)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| small_balance_found | 2 | 2 | B002, B003 | B002, B003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `small_balance_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `B002`
  2. `B003`


---

## 🧪 TEST 2: alpha_abs_positive

### 📋 Informations générales

- **Description:** Test fonction ABS() (valeur absolue)
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_abs_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_abs_positive.facts`
- **Temps d'exécution:** 439.101µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1:**
```constraint
{b: Balance} / ABS(b.amount) > 100 ==> significant_balance_found(b.id, b.amount)
```
- **Action:** `significant_balance_found`
- **Condition:** `map[args:[map[field:amount object:b type:fieldAccess]] name:ABS type:functionCall] > 100`
- **Type:** Condition positive

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `B001`
```json
Type: Balance
Champs:
  id: B001
  amount: 150
  type: credit
```

**Fait 2:** `B002`
```json
Type: Balance
Champs:
  id: B002
  amount: -200
  type: debit
```

**Fait 3:** `B003`
```json
Type: Balance
Champs:
  id: B003
  amount: 50
  type: credit
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Balance (type_Balance)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: b
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: significant_balance_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `significant_balance_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `B001` (Type: Balance)
  2. `B002` (Type: Balance)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| significant_balance_found | 2 | 2 | B001, B002 | B001, B002 | ✅ |

#### 📋 Détails des tuples attendus

**Action `significant_balance_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `B001`
  2. `B002`


---

## 🧪 TEST 3: alpha_boolean_negative

### 📋 Informations générales

- **Description:** Test condition booléenne négative
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_boolean_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_boolean_negative.facts`
- **Temps d'exécution:** 487.422µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1 🚫:**
```constraint
{a: Account} / NOT(a.active == true) ==> inactive_account_found(a.id, a.balance)
```
- **Action:** `inactive_account_found`
- **Condition:** `NOT(a.active == true)`
- **Type:** Condition négative (NOT)

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `ACC001`
```json
Type: Account
Champs:
  id: ACC001
  balance: 1000
  active: true
```

**Fait 2:** `ACC002`
```json
Type: Account
Champs:
  id: ACC002
  balance: 500
  active: false
```

**Fait 3:** `ACC003`
```json
Type: Account
Champs:
  id: ACC003
  balance: 2000
  active: true
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Account (type_Account)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: NOT(...) [Négation]
│   │   └── Variable: a
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: inactive_account_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `inactive_account_found`
- **Nombre de déclenchements:** 1
- **Faits concernés:**
  1. `ACC002` (Type: Account)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| inactive_account_found | 1 | 1 | ACC002 | ACC002 | ✅ |

#### 📋 Détails des tuples attendus

**Action `inactive_account_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 1
- **IDs attendus:**
  1. `ACC002`


---

## 🧪 TEST 4: alpha_boolean_positive

### 📋 Informations générales

- **Description:** Test condition booléenne positive
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_boolean_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_boolean_positive.facts`
- **Temps d'exécution:** 516.256µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1:**
```constraint
{a: Account} / a.active == true ==> active_account_found(a.id, a.balance)
```
- **Action:** `active_account_found`
- **Condition:** `a.active == true`
- **Type:** Condition positive

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `ACC001`
```json
Type: Account
Champs:
  balance: 1000
  active: true
  id: ACC001
```

**Fait 2:** `ACC002`
```json
Type: Account
Champs:
  id: ACC002
  balance: 500
  active: false
```

**Fait 3:** `ACC003`
```json
Type: Account
Champs:
  id: ACC003
  balance: 2000
  active: true
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Account (type_Account)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: a
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: active_account_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `active_account_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `ACC001` (Type: Account)
  2. `ACC003` (Type: Account)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| active_account_found | 2 | 2 | ACC001, ACC003 | ACC001, ACC003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `active_account_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `ACC001`
  2. `ACC003`


---

## 🧪 TEST 5: alpha_comparison_negative

### 📋 Informations générales

- **Description:** Test comparaison numérique négative
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_comparison_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_comparison_negative.facts`
- **Temps d'exécution:** 720.417µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1 🚫:**
```constraint
{prod: Product} / NOT(prod.price > 100) ==> affordable_product(prod.id, prod.price)
```
- **Action:** `affordable_product`
- **Condition:** `NOT(prod.price > 100)`
- **Type:** Condition négative (NOT)

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `PROD001`
```json
Type: Product
Champs:
  id: PROD001
  price: 150
  category: electronics
```

**Fait 2:** `PROD002`
```json
Type: Product
Champs:
  id: PROD002
  price: 50
  category: books
```

**Fait 3:** `PROD003`
```json
Type: Product
Champs:
  id: PROD003
  price: 200
  category: electronics
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Product (type_Product)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: NOT(...) [Négation]
│   │   └── Variable: prod
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: affordable_product
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `affordable_product`
- **Nombre de déclenchements:** 1
- **Faits concernés:**
  1. `PROD002` (Type: Product)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| affordable_product | 1 | 1 | PROD002 | PROD002 | ✅ |

#### 📋 Détails des tuples attendus

**Action `affordable_product`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 1
- **IDs attendus:**
  1. `PROD002`


---

## 🧪 TEST 6: alpha_comparison_positive

### 📋 Informations générales

- **Description:** Test comparaison numérique positive
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_comparison_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_comparison_positive.facts`
- **Temps d'exécution:** 558.414µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1:**
```constraint
{prod: Product} / prod.price > 100 ==> expensive_product(prod.id, prod.price)
```
- **Action:** `expensive_product`
- **Condition:** `prod.price > 100`
- **Type:** Condition positive

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `PROD001`
```json
Type: Product
Champs:
  id: PROD001
  price: 150
  category: electronics
```

**Fait 2:** `PROD002`
```json
Type: Product
Champs:
  id: PROD002
  price: 50
  category: books
```

**Fait 3:** `PROD003`
```json
Type: Product
Champs:
  id: PROD003
  price: 200
  category: electronics
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Product (type_Product)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: prod
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: expensive_product
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `expensive_product`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `PROD001` (Type: Product)
  2. `PROD003` (Type: Product)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| expensive_product | 2 | 2 | PROD001, PROD003 | PROD001, PROD003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `expensive_product`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `PROD001`
  2. `PROD003`


---

## 🧪 TEST 7: alpha_contains_negative

### 📋 Informations générales

- **Description:** Test opérateur CONTAINS (contenance) - négation
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_contains_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_contains_negative.facts`
- **Temps d'exécution:** 400.158µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1 🚫:**
```constraint
{m: Message} / NOT(m.content CONTAINS "urgent") ==> normal_message_found(m.id, m.content)
```
- **Action:** `normal_message_found`
- **Condition:** `NOT(m.content CONTAINS "urgent")`
- **Type:** Condition négative (NOT)

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `M001`
```json
Type: Message
Champs:
  id: M001
  content: This is urgent please respond
  urgent: true
```

**Fait 2:** `M002`
```json
Type: Message
Champs:
  id: M002
  content: Regular message content
  urgent: false
```

**Fait 3:** `M003`
```json
Type: Message
Champs:
  id: M003
  content: Simple notification
  urgent: false
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Message (type_Message)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: NOT(...) [Négation]
│   │   └── Variable: m
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: normal_message_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `normal_message_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `M002` (Type: Message)
  2. `M003` (Type: Message)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| normal_message_found | 2 | 2 | M002, M003 | M002, M003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `normal_message_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `M002`
  2. `M003`


---

## 🧪 TEST 8: alpha_contains_positive

### 📋 Informations générales

- **Description:** Test opérateur CONTAINS (contenance)
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_contains_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_contains_positive.facts`
- **Temps d'exécution:** 377.155µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1:**
```constraint
{m: Message} / m.content CONTAINS "urgent" ==> urgent_message_found(m.id, m.content)
```
- **Action:** `urgent_message_found`
- **Condition:** `m.content CONTAINS "urgent"`
- **Type:** Condition positive

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `M001`
```json
Type: Message
Champs:
  content: This is urgent please respond
  urgent: true
  id: M001
```

**Fait 2:** `M002`
```json
Type: Message
Champs:
  id: M002
  content: Regular message content
  urgent: false
```

**Fait 3:** `M003`
```json
Type: Message
Champs:
  id: M003
  content: Very urgent matter!
  urgent: true
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Message (type_Message)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: m
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: urgent_message_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `urgent_message_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `M001` (Type: Message)
  2. `M003` (Type: Message)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| urgent_message_found | 2 | 2 | M001, M003 | M001, M003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `urgent_message_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `M001`
  2. `M003`


---

## 🧪 TEST 9: alpha_equal_sign_negative

### 📋 Informations générales

- **Description:** Test opérateur == (égalité) - négation
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_equal_sign_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_equal_sign_negative.facts`
- **Temps d'exécution:** 436.155µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1 🚫:**
```constraint
{c: Customer} / NOT(c.tier == "gold") ==> non_gold_customer_found(c.id, c.tier)
```
- **Action:** `non_gold_customer_found`
- **Condition:** `NOT(c.tier == "gold")`
- **Type:** Condition négative (NOT)

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `C001`
```json
Type: Customer
Champs:
  id: C001
  tier: gold
  points: 5000
```

**Fait 2:** `C002`
```json
Type: Customer
Champs:
  id: C002
  tier: silver
  points: 2000
```

**Fait 3:** `C003`
```json
Type: Customer
Champs:
  points: 1000
  id: C003
  tier: bronze
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Customer (type_Customer)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: NOT(...) [Négation]
│   │   └── Variable: c
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: non_gold_customer_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `non_gold_customer_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `C002` (Type: Customer)
  2. `C003` (Type: Customer)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| non_gold_customer_found | 2 | 2 | C002, C003 | C002, C003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `non_gold_customer_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `C002`
  2. `C003`


---

## 🧪 TEST 10: alpha_equal_sign_positive

### 📋 Informations générales

- **Description:** Test opérateur == (égalité)
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_equal_sign_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_equal_sign_positive.facts`
- **Temps d'exécution:** 361.977µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1:**
```constraint
{c: Customer} / c.tier == "gold" ==> gold_customer_found(c.id, c.points)
```
- **Action:** `gold_customer_found`
- **Condition:** `c.tier == "gold"`
- **Type:** Condition positive

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `C001`
```json
Type: Customer
Champs:
  id: C001
  tier: gold
  points: 5000
```

**Fait 2:** `C002`
```json
Type: Customer
Champs:
  id: C002
  tier: silver
  points: 2000
```

**Fait 3:** `C003`
```json
Type: Customer
Champs:
  id: C003
  tier: gold
  points: 7500
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Customer (type_Customer)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: c
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: gold_customer_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `gold_customer_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `C003` (Type: Customer)
  2. `C001` (Type: Customer)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| gold_customer_found | 2 | 2 | C001, C003 | C003, C001 | ✅ |

#### 📋 Détails des tuples attendus

**Action `gold_customer_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `C001`
  2. `C003`


---

## 🧪 TEST 11: alpha_equality_negative

### 📋 Informations générales

- **Description:** Test égalité négative simple
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_equality_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_equality_negative.facts`
- **Temps d'exécution:** 382.765µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1 🚫:**
```constraint
{p: Person} / NOT(p.age == 25) ==> age_is_not_twenty_five(p.id, p.age)
```
- **Action:** `age_is_not_twenty_five`
- **Condition:** `NOT(p.age == 25)`
- **Type:** Condition négative (NOT)

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `P001`
```json
Type: Person
Champs:
  id: P001
  age: 25
  status: active
```

**Fait 2:** `P002`
```json
Type: Person
Champs:
  id: P002
  age: 30
  status: active
```

**Fait 3:** `P003`
```json
Type: Person
Champs:
  id: P003
  age: 25
  status: inactive
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Person (type_Person)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: NOT(...) [Négation]
│   │   └── Variable: p
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: age_is_not_twenty_five
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `age_is_not_twenty_five`
- **Nombre de déclenchements:** 1
- **Faits concernés:**
  1. `P002` (Type: Person)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| age_is_not_twenty_five | 1 | 1 | P002 | P002 | ✅ |

#### 📋 Détails des tuples attendus

**Action `age_is_not_twenty_five`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 1
- **IDs attendus:**
  1. `P002`


---

## 🧪 TEST 12: alpha_equality_positive

### 📋 Informations générales

- **Description:** Test égalité positive simple
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_equality_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_equality_positive.facts`
- **Temps d'exécution:** 362.027µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1:**
```constraint
{p: Person} / p.age == 25 ==> age_is_twenty_five(p.id, p.age)
```
- **Action:** `age_is_twenty_five`
- **Condition:** `p.age == 25`
- **Type:** Condition positive

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `P001`
```json
Type: Person
Champs:
  age: 25
  status: active
  id: P001
```

**Fait 2:** `P002`
```json
Type: Person
Champs:
  id: P002
  age: 30
  status: active
```

**Fait 3:** `P003`
```json
Type: Person
Champs:
  id: P003
  age: 25
  status: inactive
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Person (type_Person)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: p
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: age_is_twenty_five
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `age_is_twenty_five`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `P001` (Type: Person)
  2. `P003` (Type: Person)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| age_is_twenty_five | 2 | 2 | P001, P003 | P001, P003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `age_is_twenty_five`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `P001`
  2. `P003`


---

## 🧪 TEST 13: alpha_in_negative

### 📋 Informations générales

- **Description:** Test opérateur IN (appartenance) - négation
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_in_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_in_negative.facts`
- **Temps d'exécution:** 423.152µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1 🚫:**
```constraint
{s: Status} / NOT(s.state IN ["active", "pending"]) ==> invalid_state_found(s.id, s.state)
```
- **Action:** `invalid_state_found`
- **Condition:** `NOT(s.state IN map[elements:[map[type:string value:active] map[type:string value:pending]] type:arrayLiteral])`
- **Type:** Condition négative (NOT)

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `S001`
```json
Type: Status
Champs:
  id: S001
  state: active
  priority: 1
```

**Fait 2:** `S002`
```json
Type: Status
Champs:
  id: S002
  state: inactive
  priority: 3
```

**Fait 3:** `S003`
```json
Type: Status
Champs:
  id: S003
  state: archived
  priority: 5
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Status (type_Status)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: NOT(...) [Négation]
│   │   └── Variable: s
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: invalid_state_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `invalid_state_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `S002` (Type: Status)
  2. `S003` (Type: Status)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| invalid_state_found | 2 | 2 | S002, S003 | S002, S003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `invalid_state_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `S002`
  2. `S003`


---

## 🧪 TEST 14: alpha_in_positive

### 📋 Informations générales

- **Description:** Test opérateur IN (appartenance)
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_in_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_in_positive.facts`
- **Temps d'exécution:** 427.7µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1:**
```constraint
{s: Status} / s.state IN ["active", "pending", "review"] ==> valid_state_found(s.id, s.state)
```
- **Action:** `valid_state_found`
- **Condition:** `s.state IN map[elements:[map[type:string value:active] map[type:string value:pending] map[type:string value:review]] type:arrayLiteral]`
- **Type:** Condition positive

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `S001`
```json
Type: Status
Champs:
  id: S001
  state: active
  priority: 1
```

**Fait 2:** `S002`
```json
Type: Status
Champs:
  id: S002
  state: inactive
  priority: 3
```

**Fait 3:** `S003`
```json
Type: Status
Champs:
  id: S003
  state: pending
  priority: 2
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Status (type_Status)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: s
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: valid_state_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `valid_state_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `S001` (Type: Status)
  2. `S003` (Type: Status)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| valid_state_found | 2 | 2 | S001, S003 | S001, S003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `valid_state_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `S001`
  2. `S003`


---

## 🧪 TEST 15: alpha_inequality_negative

### 📋 Informations générales

- **Description:** Test inégalité négative
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_inequality_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_inequality_negative.facts`
- **Temps d'exécution:** 366.766µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1 🚫:**
```constraint
{o: Order} / NOT(o.status != "cancelled") ==> cancelled_order_found(o.id, o.total)
```
- **Action:** `cancelled_order_found`
- **Condition:** `NOT(o.status != "cancelled")`
- **Type:** Condition négative (NOT)

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `ORD001`
```json
Type: Order
Champs:
  id: ORD001
  total: 100
  status: pending
```

**Fait 2:** `ORD002`
```json
Type: Order
Champs:
  id: ORD002
  total: 200
  status: cancelled
```

**Fait 3:** `ORD003`
```json
Type: Order
Champs:
  id: ORD003
  total: 300
  status: completed
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Order (type_Order)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: NOT(...) [Négation]
│   │   └── Variable: o
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: cancelled_order_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `cancelled_order_found`
- **Nombre de déclenchements:** 1
- **Faits concernés:**
  1. `ORD002` (Type: Order)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| cancelled_order_found | 1 | 1 | ORD002 | ORD002 | ✅ |

#### 📋 Détails des tuples attendus

**Action `cancelled_order_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 1
- **IDs attendus:**
  1. `ORD002`


---

## 🧪 TEST 16: alpha_inequality_positive

### 📋 Informations générales

- **Description:** Test inégalité positive
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_inequality_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_inequality_positive.facts`
- **Temps d'exécution:** 491.129µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1:**
```constraint
{o: Order} / o.status != "cancelled" ==> valid_order_found(o.id, o.total)
```
- **Action:** `valid_order_found`
- **Condition:** `o.status != "cancelled"`
- **Type:** Condition positive

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `ORD001`
```json
Type: Order
Champs:
  id: ORD001
  total: 100
  status: pending
```

**Fait 2:** `ORD002`
```json
Type: Order
Champs:
  id: ORD002
  total: 200
  status: cancelled
```

**Fait 3:** `ORD003`
```json
Type: Order
Champs:
  status: completed
  id: ORD003
  total: 300
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Order (type_Order)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: o
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: valid_order_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `valid_order_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `ORD001` (Type: Order)
  2. `ORD003` (Type: Order)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| valid_order_found | 2 | 2 | ORD001, ORD003 | ORD001, ORD003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `valid_order_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `ORD001`
  2. `ORD003`


---

## 🧪 TEST 17: alpha_length_negative

### 📋 Informations générales

- **Description:** Test fonction LENGTH() - négation
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_length_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_length_negative.facts`
- **Temps d'exécution:** 379.149µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1 🚫:**
```constraint
{p: Password} / NOT(LENGTH(p.value) >= 8) ==> weak_password_found(p.id, p.value)
```
- **Action:** `weak_password_found`
- **Condition:** `NOT(map[args:[map[field:value object:p type:fieldAccess]] name:LENGTH type:functionCall] >= 8)`
- **Type:** Condition négative (NOT)

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `P001`
```json
Type: Password
Champs:
  id: P001
  value: password123
  secure: true
```

**Fait 2:** `P002`
```json
Type: Password
Champs:
  secure: false
  id: P002
  value: 123
```

**Fait 3:** `P003`
```json
Type: Password
Champs:
  id: P003
  value: pass
  secure: false
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Password (type_Password)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: NOT(...) [Négation]
│   │   └── Variable: p
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: weak_password_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `weak_password_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `P002` (Type: Password)
  2. `P003` (Type: Password)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| weak_password_found | 2 | 2 | P002, P003 | P002, P003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `weak_password_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `P002`
  2. `P003`


---

## 🧪 TEST 18: alpha_length_positive

### 📋 Informations générales

- **Description:** Test fonction LENGTH()
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_length_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_length_positive.facts`
- **Temps d'exécution:** 348.11µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1:**
```constraint
{p: Password} / LENGTH(p.value) >= 8 ==> secure_password_found(p.id, p.value)
```
- **Action:** `secure_password_found`
- **Condition:** `map[args:[map[field:value object:p type:fieldAccess]] name:LENGTH type:functionCall] >= 8`
- **Type:** Condition positive

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `P001`
```json
Type: Password
Champs:
  id: P001
  value: password123
  secure: true
```

**Fait 2:** `P002`
```json
Type: Password
Champs:
  id: P002
  value: 123
  secure: false
```

**Fait 3:** `P003`
```json
Type: Password
Champs:
  id: P003
  value: verysecurepass
  secure: true
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Password (type_Password)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: p
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: secure_password_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `secure_password_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `P001` (Type: Password)
  2. `P003` (Type: Password)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| secure_password_found | 2 | 2 | P001, P003 | P001, P003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `secure_password_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `P001`
  2. `P003`


---

## 🧪 TEST 19: alpha_like_negative

### 📋 Informations générales

- **Description:** Test opérateur LIKE (motif) - négation
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_like_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_like_negative.facts`
- **Temps d'exécution:** 483.935µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1 🚫:**
```constraint
{e: Email} / NOT(e.address LIKE "%@company.com") ==> external_email_found(e.id, e.address)
```
- **Action:** `external_email_found`
- **Condition:** `NOT(e.address LIKE "%@company.com")`
- **Type:** Condition négative (NOT)

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `E001`
```json
Type: Email
Champs:
  id: E001
  address: john@company.com
  verified: true
```

**Fait 2:** `E002`
```json
Type: Email
Champs:
  verified: false
  id: E002
  address: jane@external.org
```

**Fait 3:** `E003`
```json
Type: Email
Champs:
  id: E003
  address: user@other.net
  verified: true
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Email (type_Email)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: NOT(...) [Négation]
│   │   └── Variable: e
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: external_email_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `external_email_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `E003` (Type: Email)
  2. `E002` (Type: Email)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| external_email_found | 2 | 2 | E002, E003 | E003, E002 | ✅ |

#### 📋 Détails des tuples attendus

**Action `external_email_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `E002`
  2. `E003`


---

## 🧪 TEST 20: alpha_like_positive

### 📋 Informations générales

- **Description:** Test opérateur LIKE (motif)
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_like_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_like_positive.facts`
- **Temps d'exécution:** 755.473µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1:**
```constraint
{e: Email} / e.address LIKE "%@company.com" ==> company_email_found(e.id, e.address)
```
- **Action:** `company_email_found`
- **Condition:** `e.address LIKE "%@company.com"`
- **Type:** Condition positive

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `E001`
```json
Type: Email
Champs:
  id: E001
  address: john@company.com
  verified: true
```

**Fait 2:** `E002`
```json
Type: Email
Champs:
  address: jane@external.org
  verified: false
  id: E002
```

**Fait 3:** `E003`
```json
Type: Email
Champs:
  id: E003
  address: admin@company.com
  verified: true
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Email (type_Email)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: e
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: company_email_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `company_email_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `E001` (Type: Email)
  2. `E003` (Type: Email)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| company_email_found | 2 | 2 | E001, E003 | E001, E003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `company_email_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `E001`
  2. `E003`


---

## 🧪 TEST 21: alpha_matches_negative

### 📋 Informations générales

- **Description:** Test opérateur MATCHES (regex) - négation
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_matches_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_matches_negative.facts`
- **Temps d'exécution:** 1.227194ms
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1 🚫:**
```constraint
{c: Code} / NOT(c.value MATCHES "CODE[0-9]+") ==> invalid_code_found(c.id, c.value)
```
- **Action:** `invalid_code_found`
- **Condition:** `NOT(c.value MATCHES "CODE[0-9]+")`
- **Type:** Condition négative (NOT)

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `C001`
```json
Type: Code
Champs:
  id: C001
  value: CODE123
  active: true
```

**Fait 2:** `C002`
```json
Type: Code
Champs:
  id: C002
  value: INVALID
  active: false
```

**Fait 3:** `C003`
```json
Type: Code
Champs:
  value: BADFORMAT
  active: true
  id: C003
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Code (type_Code)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: NOT(...) [Négation]
│   │   └── Variable: c
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: invalid_code_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `invalid_code_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `C003` (Type: Code)
  2. `C002` (Type: Code)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| invalid_code_found | 2 | 2 | C002, C003 | C003, C002 | ✅ |

#### 📋 Détails des tuples attendus

**Action `invalid_code_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `C002`
  2. `C003`


---

## 🧪 TEST 22: alpha_matches_positive

### 📋 Informations générales

- **Description:** Test opérateur MATCHES (regex)
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_matches_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_matches_positive.facts`
- **Temps d'exécution:** 364.671µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1:**
```constraint
{c: Code} / c.value MATCHES "CODE[0-9]+" ==> valid_code_found(c.id, c.value)
```
- **Action:** `valid_code_found`
- **Condition:** `c.value MATCHES "CODE[0-9]+"`
- **Type:** Condition positive

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `C001`
```json
Type: Code
Champs:
  id: C001
  value: CODE123
  active: true
```

**Fait 2:** `C002`
```json
Type: Code
Champs:
  id: C002
  value: INVALID
  active: false
```

**Fait 3:** `C003`
```json
Type: Code
Champs:
  id: C003
  value: CODE999
  active: true
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Code (type_Code)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: c
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: valid_code_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `valid_code_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `C001` (Type: Code)
  2. `C003` (Type: Code)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| valid_code_found | 2 | 2 | C001, C003 | C001, C003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `valid_code_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `C001`
  2. `C003`


---

## 🧪 TEST 23: alpha_string_negative

### 📋 Informations générales

- **Description:** Test condition string négative
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_string_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_string_negative.facts`
- **Temps d'exécution:** 313.046µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1 🚫:**
```constraint
{u: User} / NOT(u.role == "admin") ==> non_admin_user_found(u.id, u.name)
```
- **Action:** `non_admin_user_found`
- **Condition:** `NOT(u.role == "admin")`
- **Type:** Condition négative (NOT)

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `U001`
```json
Type: User
Champs:
  id: U001
  name: Alice
  role: admin
```

**Fait 2:** `U002`
```json
Type: User
Champs:
  id: U002
  name: Bob
  role: user
```

**Fait 3:** `U003`
```json
Type: User
Champs:
  id: U003
  name: Charlie
  role: admin
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── User (type_User)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: NOT(...) [Négation]
│   │   └── Variable: u
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: non_admin_user_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `non_admin_user_found`
- **Nombre de déclenchements:** 1
- **Faits concernés:**
  1. `U002` (Type: User)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| non_admin_user_found | 1 | 1 | U002 | U002 | ✅ |

#### 📋 Détails des tuples attendus

**Action `non_admin_user_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 1
- **IDs attendus:**
  1. `U002`


---

## 🧪 TEST 24: alpha_string_positive

### 📋 Informations générales

- **Description:** Test condition string positive
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_string_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_string_positive.facts`
- **Temps d'exécution:** 312.224µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1:**
```constraint
{u: User} / u.role == "admin" ==> admin_user_found(u.id, u.name)
```
- **Action:** `admin_user_found`
- **Condition:** `u.role == "admin"`
- **Type:** Condition positive

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `U001`
```json
Type: User
Champs:
  role: admin
  id: U001
  name: Alice
```

**Fait 2:** `U002`
```json
Type: User
Champs:
  name: Bob
  role: user
  id: U002
```

**Fait 3:** `U003`
```json
Type: User
Champs:
  id: U003
  name: Charlie
  role: admin
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── User (type_User)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: u
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: admin_user_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `admin_user_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `U001` (Type: User)
  2. `U003` (Type: User)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| admin_user_found | 2 | 2 | U001, U003 | U001, U003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `admin_user_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `U001`
  2. `U003`


---

## 🧪 TEST 25: alpha_upper_negative

### 📋 Informations générales

- **Description:** Test fonction UPPER() (majuscules) - négation
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_upper_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_upper_negative.facts`
- **Temps d'exécution:** 425.436µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1 🚫:**
```constraint
{d: Department} / NOT(UPPER(d.name) == "FINANCE") ==> non_finance_dept_found(d.id, d.name)
```
- **Action:** `non_finance_dept_found`
- **Condition:** `NOT(map[args:[map[field:name object:d type:fieldAccess]] name:UPPER type:functionCall] == "FINANCE")`
- **Type:** Condition négative (NOT)

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `D001`
```json
Type: Department
Champs:
  id: D001
  name: finance
  active: true
```

**Fait 2:** `D002`
```json
Type: Department
Champs:
  id: D002
  name: IT
  active: true
```

**Fait 3:** `D003`
```json
Type: Department
Champs:
  id: D003
  name: HR
  active: true
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Department (type_Department)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: NOT(...) [Négation]
│   │   └── Variable: d
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: non_finance_dept_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `non_finance_dept_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `D002` (Type: Department)
  2. `D003` (Type: Department)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| non_finance_dept_found | 2 | 2 | D002, D003 | D002, D003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `non_finance_dept_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `D002`
  2. `D003`


---

## 🧪 TEST 26: alpha_upper_positive

### 📋 Informations générales

- **Description:** Test fonction UPPER() (majuscules)
- **Fichier contraintes:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_upper_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/test/coverage/alpha/alpha_upper_positive.facts`
- **Temps d'exécution:** 382.596µs
- **Statut:** ✅ Succès

### 📏 Règles du test

**Règle 1:**
```constraint
{d: Department} / UPPER(d.name) == "FINANCE" ==> finance_dept_found(d.id, d.name)
```
- **Action:** `finance_dept_found`
- **Condition:** `map[args:[map[field:name object:d type:fieldAccess]] name:UPPER type:functionCall] == "FINANCE"`
- **Type:** Condition positive

### 📦 Faits du test

**Nombre total:** 3 faits

**Fait 1:** `D001`
```json
Type: Department
Champs:
  id: D001
  name: finance
  active: true
```

**Fait 2:** `D002`
```json
Type: Department
Champs:
  id: D002
  name: IT
  active: true
```

**Fait 3:** `D003`
```json
Type: Department
Champs:
  active: true
  id: D003
  name: Finance
```

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE - STRUCTURE HIÉRARCHIQUE
=====================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Department (type_Department)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: d
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: finance_dept_found
```

### ⚡ Résultats d'exécution

**1 actions déclenchées:**

#### 🎯 Action: `finance_dept_found`
- **Nombre de déclenchements:** 2
- **Faits concernés:**
  1. `D001` (Type: Department)
  2. `D003` (Type: Department)

### 🧠 Validation sémantique

- **Score de validation:** 100.0%
- **✅ Validation parfaite**

### 📊 Comparaison Attendu vs Observé

| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| finance_dept_found | 2 | 2 | D001, D003 | D001, D003 | ✅ |

#### 📋 Détails des tuples attendus

**Action `finance_dept_found`:**
- **Description:** Action basée sur règle 1
- **Faits attendus:** 2
- **IDs attendus:**
  1. `D001`
  2. `D003`


---
