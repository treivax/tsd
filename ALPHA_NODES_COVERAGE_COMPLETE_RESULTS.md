# 📊 RAPPORT COMPLET - TESTS DE COUVERTURE ALPHA NODES

**Date d'exécution:** 2025-11-17 11:24:48
**Nombre de tests:** 10

## 🎯 RÉSUMÉ EXÉCUTIF

- ✅ **Tests réussis:** 10/10 (100.0%)
- 🎬 **Actions déclenchées:** 10
- ⚡ **Couverture:** Nœuds Alpha positifs et négatifs

## 🧪 TEST 1: alpha_boolean_negative

### 📋 Informations générales

- **Description:** Test condition booléenne négative
- **Fichier contraintes:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_boolean_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_boolean_negative.facts`
- **Temps d'exécution:** 871.631µs
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

---

## 🧪 TEST 2: alpha_boolean_positive

### 📋 Informations générales

- **Description:** Test condition booléenne positive
- **Fichier contraintes:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_boolean_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_boolean_positive.facts`
- **Temps d'exécution:** 663.812µs
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

---

## 🧪 TEST 3: alpha_comparison_negative

### 📋 Informations générales

- **Description:** Test comparaison numérique négative
- **Fichier contraintes:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_comparison_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_comparison_negative.facts`
- **Temps d'exécution:** 618.988µs
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

---

## 🧪 TEST 4: alpha_comparison_positive

### 📋 Informations générales

- **Description:** Test comparaison numérique positive
- **Fichier contraintes:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_comparison_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_comparison_positive.facts`
- **Temps d'exécution:** 669.483µs
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

---

## 🧪 TEST 5: alpha_equality_negative

### 📋 Informations générales

- **Description:** Test égalité négative simple
- **Fichier contraintes:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_equality_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_equality_negative.facts`
- **Temps d'exécution:** 668.49µs
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

---

## 🧪 TEST 6: alpha_equality_positive

### 📋 Informations générales

- **Description:** Test égalité positive simple
- **Fichier contraintes:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_equality_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_equality_positive.facts`
- **Temps d'exécution:** 634.306µs
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
  age: 25
  status: inactive
  id: P003
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

---

## 🧪 TEST 7: alpha_inequality_negative

### 📋 Informations générales

- **Description:** Test inégalité négative
- **Fichier contraintes:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_inequality_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_inequality_negative.facts`
- **Temps d'exécution:** 666.717µs
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
  total: 100
  status: pending
  id: ORD001
```

**Fait 2:** `ORD002`
```json
Type: Order
Champs:
  status: cancelled
  id: ORD002
  total: 200
```

**Fait 3:** `ORD003`
```json
Type: Order
Champs:
  total: 300
  status: completed
  id: ORD003
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

---

## 🧪 TEST 8: alpha_inequality_positive

### 📋 Informations générales

- **Description:** Test inégalité positive
- **Fichier contraintes:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_inequality_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_inequality_positive.facts`
- **Temps d'exécution:** 570.016µs
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
  status: cancelled
  id: ORD002
  total: 200
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

---

## 🧪 TEST 9: alpha_string_negative

### 📋 Informations générales

- **Description:** Test condition string négative
- **Fichier contraintes:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_string_negative.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_string_negative.facts`
- **Temps d'exécution:** 428.792µs
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
  role: admin
  id: U001
  name: Alice
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

---

## 🧪 TEST 10: alpha_string_positive

### 📋 Informations générales

- **Description:** Test condition string positive
- **Fichier contraintes:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_string_positive.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/alpha_coverage_tests/alpha_string_positive.facts`
- **Temps d'exécution:** 409.647µs
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
  id: U002
  name: Bob
  role: user
```

**Fait 3:** `U003`
```json
Type: User
Champs:
  role: admin
  id: U003
  name: Charlie
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

---

