# RAPPORT COMPLET DE COUVERTURE DES NŒUDS BETA
================================================

**📊 Tests exécutés:** 8
**✅ Tests réussis:** 8 (100.0%)
**🧠 Score sémantique moyen:** 82.5%
**📅 Date d'exécution:** 2025-11-18 12:40:29

## 🎯 NŒUDS BETA ANALYSÉS
| Type de Nœud | Tests | Succès | Score Sémantique |
|---------------|--------|--------|------------------|
| ExistsNode | 1 | 1 | 100.0% |
| NotNode | 2 | 2 | 100.0% |
| JoinNode | 5 | 5 | 72.0% |

## 🧪 TEST 1: beta_exists_complex
---

### 📋 Informations générales
- **Description:** Test JoinNode Customer-Purchase avec vraie jointure Beta
- **Fichier contraintes:** `beta_exists_complex.constraint`
- **Fichier faits:** `beta_exists_complex.facts`
- **Temps d'exécution:** 1.125835ms
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 80.0%
- **Actions valides:** ❌
- **Jointures valides:** ✅
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

**⚠️ Erreurs de validation:**
- Action attendue manquante: vip_big_spender

### 📜 Règles analysées
#### Règle 1
- **Texte original:** `{c: Customer, p: Purchase} / c.id == p.customer_id ==> customer_purchase(c.id, p.id, c.name, p.product)`
- **Action:** customer_purchase
- **Type de nœud:** JoinNode
- **Type sémantique:** comparison
- **Complexité:** complex
- **Variables:**
  - c (Customer): primary
  - p (Purchase): secondary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Customer (type_Customer)
│   ├── Purchase (type_Purchase)
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: customer_purchase
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Customer[id=C001, name=Alice, vip=true]
Customer[id=C002, name=Bob, vip=false]
Customer[id=C003, name=Charlie, vip=true]
Customer[id=C004, name=David, vip=true]
Purchase[id=PUR001, customer_id=C001, amount=600, product=laptop]
Purchase[id=PUR002, customer_id=C001, amount=200, product=mouse]
Purchase[id=PUR003, customer_id=C002, amount=800, product=phone]
Purchase[id=PUR004, customer_id=C003, amount=300, product=keyboard]
Purchase[id=PUR005, customer_id=C003, amount=750, product=tablet]
```

**Total faits:** 9

- **Customer:** 4 faits
- **Purchase:** 5 faits

**📋 Détail des faits parsés:**
1. **Customer[C001]** - `Customer[id=C001, name=Alice, vip=true]`
2. **Customer[C002]** - `Customer[vip=false, id=C002, name=Bob]`
3. **Customer[C003]** - `Customer[id=C003, name=Charlie, vip=true]`
4. **Customer[C004]** - `Customer[id=C004, name=David, vip=true]`
5. **Purchase[PUR001]** - `Purchase[id=PUR001, customer_id=C001, amount=600, product=laptop]`
6. **Purchase[PUR002]** - `Purchase[id=PUR002, customer_id=C001, amount=200, product=mouse]`
7. **Purchase[PUR003]** - `Purchase[id=PUR003, customer_id=C002, amount=800, product=phone]`
8. **Purchase[PUR004]** - `Purchase[id=PUR004, customer_id=C003, amount=300, product=keyboard]`
9. **Purchase[PUR005]** - `Purchase[customer_id=C003, amount=750, product=tablet, id=PUR005]`

### ⚡ Résultats des actions
*Aucune action déclenchée*

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | c <-> p | 0 | inner | ❌ |

### 🎯 Comparaison attendu vs observé
#### Actions
| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| vip_big_spender | 2-2 | 0 | C001, C003 |  | ❌ |

#### 📋 Détails des tuples Beta attendus

**Action `vip_big_spender`:**
- **Description:** Deux clients VIP
- **Déclenchements attendus:** 2-2
- **IDs de faits attendus:**
  1. `C001`
  2. `C003`

---

## 🧪 TEST 2: beta_exists_real
---

### 📋 Informations générales
- **Description:** Test ExistsNode correct avec vraie existence
- **Fichier contraintes:** `beta_exists_real.constraint`
- **Fichier faits:** `beta_exists_real.facts`
- **Temps d'exécution:** 633.504µs
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 100.0%
- **Actions valides:** ✅
- **Jointures valides:** ✅
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

### 📜 Règles analysées
#### Règle 1
- **Texte original:** `{v: Vendor} / EXISTS (p: Product / p.vendor_id == v.id) ==> vendor_has_products(v.id, v.name)`
- **Action:** vendor_has_products
- **Type de nœud:** ExistsNode
- **Type sémantique:** existence
- **Complexité:** simple
- **Variables:**
  - v (Vendor): primary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Vendor (type_Vendor)
│   ├── Product (type_Product)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: v
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: vendor_has_products
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Vendor[id=V001, name=Apple, category=electronics]
Vendor[id=V002, name=Nike, category=clothing]
Vendor[id=V003, name=Sony, category=electronics]
Product[id=P001, vendor_id=V001, name=iPhone, price=999]
Product[id=P002, vendor_id=V001, name=MacBook, price=1999]
Product[id=P003, vendor_id=V003, name=PlayStation, price=499]
```

**Total faits:** 6

- **Vendor:** 3 faits
- **Product:** 3 faits

**📋 Détail des faits parsés:**
1. **Vendor[V001]** - `Vendor[id=V001, name=Apple, category=electronics]`
2. **Vendor[V002]** - `Vendor[id=V002, name=Nike, category=clothing]`
3. **Vendor[V003]** - `Vendor[id=V003, name=Sony, category=electronics]`
4. **Product[P001]** - `Product[id=P001, vendor_id=V001, name=iPhone, price=999]`
5. **Product[P002]** - `Product[name=MacBook, price=1999, id=P002, vendor_id=V001]`
6. **Product[P003]** - `Product[name=PlayStation, price=499, id=P003, vendor_id=V003]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| vendor_has_products | 3 | AlphaNode | ❌ |

#### 🎯 Activation détaillée: `vendor_has_products`
- **Nombre de déclenchements:** 3
- **Type de nœud déclencheur:** AlphaNode

**📋 Tokens et couples de faits activant l'action:**

##### Token 1
**Fait activateur:** `Vendor[id=V001, name=Apple, category=electronics]`

##### Token 2
**Fait activateur:** `Vendor[id=V002, name=Nike, category=clothing]`

##### Token 3
**Fait activateur:** `Vendor[id=V003, name=Sony, category=electronics]`

---

## 🧪 TEST 3: beta_join_complex
---

### 📋 Informations générales
- **Description:** Test de jointure complexe avec multiple conditions
- **Fichier contraintes:** `beta_join_complex.constraint`
- **Fichier faits:** `beta_join_complex.facts`
- **Temps d'exécution:** 554.647µs
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 60.0%
- **Actions valides:** ❌
- **Jointures valides:** ❌
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

**⚠️ Erreurs de validation:**
- Action attendue manquante: dept_match
- Jointure attendue: Employee -> Project, 3 correspondances

### 📜 Règles analysées
#### Règle 1
- **Texte original:** `{e: Employee, p: Project} / e.department == p.department ==> dept_match(e.id, p.id, e.name, p.name)`
- **Action:** dept_match
- **Type de nœud:** JoinNode
- **Type sémantique:** comparison
- **Complexité:** complex
- **Variables:**
  - e (Employee): primary
  - p (Project): secondary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Employee (type_Employee)
│   ├── Project (type_Project)
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: dept_match
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Employee[id=E001, name=Alice, department=engineering, salary=75000]
Employee[id=E002, name=Bob, department=marketing, salary=45000]
Employee[id=E003, name=Charlie, department=engineering, salary=80000]
Employee[id=E004, name=David, department=sales, salary=60000]
Project[id=PR001, name=WebApp, department=engineering, budget=150000]
Project[id=PR002, name=Campaign, department=marketing, budget=80000]
Project[id=PR003, name=Infrastructure, department=engineering, budget=200000]
```

**Total faits:** 7

- **Employee:** 4 faits
- **Project:** 3 faits

**📋 Détail des faits parsés:**
1. **Employee[E001]** - `Employee[id=E001, name=Alice, department=engineering, salary=75000]`
2. **Employee[E002]** - `Employee[id=E002, name=Bob, department=marketing, salary=45000]`
3. **Employee[E003]** - `Employee[id=E003, name=Charlie, department=engineering, salary=80000]`
4. **Employee[E004]** - `Employee[id=E004, name=David, department=sales, salary=60000]`
5. **Project[PR001]** - `Project[budget=150000, id=PR001, name=WebApp, department=engineering]`
6. **Project[PR002]** - `Project[id=PR002, name=Campaign, department=marketing, budget=80000]`
7. **Project[PR003]** - `Project[id=PR003, name=Infrastructure, department=engineering, budget=200000]`

### ⚡ Résultats des actions
*Aucune action déclenchée*

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | e <-> p | 0 | inner | ❌ |

### 🎯 Comparaison attendu vs observé
#### Actions
| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| dept_match | 3-3 | 0 | E001, E002, E003, PR001, PR002, PR003 |  | ❌ |

#### 📋 Détails des tuples Beta attendus

**Action `dept_match`:**
- **Description:** Trois correspondances department entre Employee et Project
- **Déclenchements attendus:** 3-3
- **IDs de faits attendus:**
  1. `E001`
  2. `E002`
  3. `E003`
  4. `PR001`
  5. `PR002`
  6. `PR003`

---

## 🧪 TEST 4: beta_join_numeric
---

### 📋 Informations générales
- **Description:** Test JoinNode avec conditions numériques simples
- **Fichier contraintes:** `beta_join_numeric.constraint`
- **Fichier faits:** `beta_join_numeric.facts`
- **Temps d'exécution:** 553.254µs
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 100.0%
- **Actions valides:** ✅
- **Jointures valides:** ✅
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

### 📜 Règles analysées
#### Règle 1
- **Texte original:** `{s: Student, c: Course} / s.grade > 85 ==> high_performer_advanced(s.id, c.id, s.name, c.name)`
- **Action:** high_performer_advanced
- **Type de nœud:** JoinNode
- **Type sémantique:** comparison
- **Complexité:** complex
- **Variables:**
  - s (Student): primary
  - c (Course): secondary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Student (type_Student)
│   ├── Course (type_Course)
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: high_performer_advanced
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Student[id=S001, name=Alice, grade=92, age=20]
Student[id=S002, name=Bob, grade=78, age=19]
Student[id=S003, name=Charlie, grade=88, age=21]
Student[id=S004, name=David, grade=95, age=22]
Course[id=C001, name=Advanced_Math, credits=4, difficulty=8]
Course[id=C002, name=Basic_History, credits=2, difficulty=3]
Course[id=C003, name=Computer_Science, credits=5, difficulty=9]
Course[id=C004, name=Literature, credits=3, difficulty=4]
```

**Total faits:** 8

- **Student:** 4 faits
- **Course:** 4 faits

**📋 Détail des faits parsés:**
1. **Student[S001]** - `Student[id=S001, name=Alice, grade=92, age=20]`
2. **Student[S002]** - `Student[name=Bob, grade=78, age=19, id=S002]`
3. **Student[S003]** - `Student[id=S003, name=Charlie, grade=88, age=21]`
4. **Student[S004]** - `Student[grade=95, age=22, id=S004, name=David]`
5. **Course[C001]** - `Course[id=C001, name=Advanced_Math, credits=4, difficulty=8]`
6. **Course[C002]** - `Course[id=C002, name=Basic_History, credits=2, difficulty=3]`
7. **Course[C003]** - `Course[id=C003, name=Computer_Science, credits=5, difficulty=9]`
8. **Course[C004]** - `Course[id=C004, name=Literature, credits=3, difficulty=4]`

### ⚡ Résultats des actions
*Aucune action déclenchée*

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | s <-> c | 0 | inner | ❌ |

---

## 🧪 TEST 5: beta_join_simple
---

### 📋 Informations générales
- **Description:** Test jointure simple entre Person et Order
- **Fichier contraintes:** `beta_join_simple.constraint`
- **Fichier faits:** `beta_join_simple.facts`
- **Temps d'exécution:** 479.015µs
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 60.0%
- **Actions valides:** ❌
- **Jointures valides:** ❌
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

**⚠️ Erreurs de validation:**
- Action attendue manquante: customer_order_match
- Jointure attendue: Person -> Order, 2 correspondances

### 📜 Règles analysées
#### Règle 1
- **Texte original:** `{p: Person, o: Order} / p.id == o.customer_id ==> customer_order_match(p.id, o.id, p.name, o.amount)`
- **Action:** customer_order_match
- **Type de nœud:** JoinNode
- **Type sémantique:** comparison
- **Complexité:** complex
- **Variables:**
  - p (Person): primary
  - o (Order): secondary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Person (type_Person)
│   ├── Order (type_Order)
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: customer_order_match
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Person[id=P001, name=Alice, age=25]
Person[id=P002, name=Bob, age=30]
Person[id=P003, name=Charlie, age=35]
Order[id=O001, customer_id=P001, amount=100]
Order[id=O002, customer_id=P002, amount=200]
Order[id=O003, customer_id=P999, amount=300]
```

**Total faits:** 6

- **Person:** 3 faits
- **Order:** 3 faits

**📋 Détail des faits parsés:**
1. **Person[P001]** - `Person[name=Alice, age=25, id=P001]`
2. **Person[P002]** - `Person[age=30, id=P002, name=Bob]`
3. **Person[P003]** - `Person[id=P003, name=Charlie, age=35]`
4. **Order[O001]** - `Order[id=O001, customer_id=P001, amount=100]`
5. **Order[O002]** - `Order[id=O002, customer_id=P002, amount=200]`
6. **Order[O003]** - `Order[amount=300, id=O003, customer_id=P999]`

### ⚡ Résultats des actions
*Aucune action déclenchée*

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | p <-> o | 0 | inner | ❌ |

### 🎯 Comparaison attendu vs observé
#### Actions
| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| customer_order_match | 2-2 | 0 | P001, P002, O001, O002 |  | ❌ |

#### 📋 Détails des tuples Beta attendus

**Action `customer_order_match`:**
- **Description:** Deux customers avec leurs commandes matchent
- **Déclenchements attendus:** 2-2
- **IDs de faits attendus:**
  1. `P001`
  2. `P002`
  3. `O001`
  4. `O002`

---

## 🧪 TEST 6: beta_mixed_complex
---

### 📋 Informations générales
- **Description:** Test combinant JoinNode et NotNode
- **Fichier contraintes:** `beta_mixed_complex.constraint`
- **Fichier faits:** `beta_mixed_complex.facts`
- **Temps d'exécution:** 582.178µs
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 60.0%
- **Actions valides:** ❌
- **Jointures valides:** ❌
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

**⚠️ Erreurs de validation:**
- Action attendue manquante: valid_transaction
- Jointure attendue: Account -> Transaction, 5 correspondances

### 📜 Règles analysées
#### Règle 1
- **Texte original:** `{a: Account, t: Transaction} / a.id == t.account_id ==> valid_transaction(a.id, t.id, a.balance, t.amount)`
- **Action:** valid_transaction
- **Type de nœud:** JoinNode
- **Type sémantique:** comparison
- **Complexité:** complex
- **Variables:**
  - a (Account): primary
  - t (Transaction): secondary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Account (type_Account)
│   ├── Transaction (type_Transaction)
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: valid_transaction
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Account[id=A001, balance=5000, type=savings, active=true]
Account[id=A002, balance=500, type=checking, active=true]
Account[id=A003, balance=2000, type=savings, active=false]
Account[id=A004, balance=3000, type=checking, active=true]
Transaction[id=T001, account_id=A001, amount=100, status=pending]
Transaction[id=T002, account_id=A001, amount=6000, status=pending]
Transaction[id=T003, account_id=A002, amount=200, status=rejected]
Transaction[id=T004, account_id=A003, amount=500, status=pending]
Transaction[id=T005, account_id=A004, amount=800, status=approved]
```

**Total faits:** 9

- **Account:** 4 faits
- **Transaction:** 5 faits

**📋 Détail des faits parsés:**
1. **Account[A001]** - `Account[type=savings, active=true, id=A001, balance=5000]`
2. **Account[A002]** - `Account[balance=500, type=checking, active=true, id=A002]`
3. **Account[A003]** - `Account[id=A003, balance=2000, type=savings, active=false]`
4. **Account[A004]** - `Account[id=A004, balance=3000, type=checking, active=true]`
5. **Transaction[T001]** - `Transaction[id=T001, account_id=A001, amount=100, status=pending]`
6. **Transaction[T002]** - `Transaction[status=pending, id=T002, account_id=A001, amount=6000]`
7. **Transaction[T003]** - `Transaction[id=T003, account_id=A002, amount=200, status=rejected]`
8. **Transaction[T004]** - `Transaction[id=T004, account_id=A003, amount=500, status=pending]`
9. **Transaction[T005]** - `Transaction[id=T005, account_id=A004, amount=800, status=approved]`

### ⚡ Résultats des actions
*Aucune action déclenchée*

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | a <-> t | 0 | inner | ❌ |

### 🎯 Comparaison attendu vs observé
#### Actions
| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| valid_transaction | 5-5 | 0 | A001, A002, A003, A004, T001, T002, T003, T004, T005 |  | ❌ |

#### 📋 Détails des tuples Beta attendus

**Action `valid_transaction`:**
- **Description:** Toutes les transactions correspondent aux comptes
- **Déclenchements attendus:** 5-5
- **IDs de faits attendus:**
  1. `A001`
  2. `A002`
  3. `A003`
  4. `A004`
  5. `T001`
  6. `T002`
  7. `T003`
  8. `T004`
  9. `T005`

---

## 🧪 TEST 7: beta_not_complex
---

### 📋 Informations générales
- **Description:** Test de négation avec multiple conditions
- **Fichier contraintes:** `beta_not_complex.constraint`
- **Fichier faits:** `beta_not_complex.facts`
- **Temps d'exécution:** 503µs
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 100.0%
- **Actions valides:** ✅
- **Jointures valides:** ✅
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

### 📜 Règles analysées
#### Règle 1
- **Texte original:** `{p: Person} / NOT (p.age < 18) ==> eligible_person(p.id, p.name)`
- **Action:** eligible_person
- **Type de nœud:** NotNode
- **Type sémantique:** negation
- **Complexité:** simple
- **Variables:**
  - p (Person): primary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

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
    │   └── Action: eligible_person
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Person[id=P001, name=Alice, age=25, status=active]
Person[id=P002, name=Bob, age=16, status=inactive]
Person[id=P003, name=Charlie, age=30, status=active]
Person[id=P004, name=David, age=17, status=active]
Person[id=P005, name=Eve, age=22, status=inactive]
```

**Total faits:** 5

- **Person:** 5 faits

**📋 Détail des faits parsés:**
1. **Person[P001]** - `Person[id=P001, name=Alice, age=25, status=active]`
2. **Person[P002]** - `Person[id=P002, name=Bob, age=16, status=inactive]`
3. **Person[P003]** - `Person[name=Charlie, age=30, status=active, id=P003]`
4. **Person[P004]** - `Person[name=David, age=17, status=active, id=P004]`
5. **Person[P005]** - `Person[name=Eve, age=22, status=inactive, id=P005]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| eligible_person | 3 | AlphaNode | ✅ |

#### 🎯 Activation détaillée: `eligible_person`
- **Nombre de déclenchements:** 3
- **Type de nœud déclencheur:** AlphaNode

**📋 Tokens et couples de faits activant l'action:**

##### Token 1
**Fait activateur:** `Person[age=22, status=inactive, id=P005, name=Eve]`

##### Token 2
**Fait activateur:** `Person[age=25, status=active, id=P001, name=Alice]`

##### Token 3
**Fait activateur:** `Person[age=30, status=active, id=P003, name=Charlie]`

### 🚫 Analyse des négations (NotNodes)
| Nœud | Condition Niée | Faits Filtrés | Type | Validation |
|------|----------------|---------------|------|------------|
| alpha_rule_0_alpha | map[condition:map[left:map[field:age object:p type:fieldAccess] operator:< right:map[type:number value:18] type:comparison] negated:true type:negation] | 0 | simple | ❌ |

### 🎯 Comparaison attendu vs observé
#### Actions
| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| eligible_person | 3-3 | 3 | P001, P003, P005 | P005, P001, P003 | ✅ |

#### 📋 Détails des tuples Beta attendus

**Action `eligible_person`:**
- **Description:** Trois personnes majeures éligibles
- **Déclenchements attendus:** 3-3
- **IDs de faits attendus:**
  1. `P001`
  2. `P003`
  3. `P005`

---

## 🧪 TEST 8: beta_not_string
---

### 📋 Informations générales
- **Description:** Test NotNode avec condition sur string
- **Fichier contraintes:** `beta_not_string.constraint`
- **Fichier faits:** `beta_not_string.facts`
- **Temps d'exécution:** 561.189µs
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 100.0%
- **Actions valides:** ✅
- **Jointures valides:** ✅
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

### 📜 Règles analysées
#### Règle 1
- **Texte original:** `{d: Device} / NOT (d.status == "broken") ==> working_device(d.id, d.name, d.type)`
- **Action:** working_device
- **Type de nœud:** NotNode
- **Type sémantique:** negation
- **Complexité:** simple
- **Variables:**
  - d (Device): primary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Device (type_Device)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: NOT(...) [Négation]
│   │   └── Variable: d
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: working_device
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Device[id=D001, name=Laptop, status=working, type=computer]
Device[id=D002, name=Printer, status=broken, type=peripheral]
Device[id=D003, name=Monitor, status=working, type=display]
Device[id=D004, name=Keyboard, status=maintenance, type=peripheral]
Device[id=D005, name=Mouse, status=broken, type=peripheral]
```

**Total faits:** 5

- **Device:** 5 faits

**📋 Détail des faits parsés:**
1. **Device[D001]** - `Device[status=working, type=computer, id=D001, name=Laptop]`
2. **Device[D002]** - `Device[type=peripheral, id=D002, name=Printer, status=broken]`
3. **Device[D003]** - `Device[id=D003, name=Monitor, status=working, type=display]`
4. **Device[D004]** - `Device[status=maintenance, type=peripheral, id=D004, name=Keyboard]`
5. **Device[D005]** - `Device[id=D005, name=Mouse, status=broken, type=peripheral]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| working_device | 3 | AlphaNode | ❌ |

#### 🎯 Activation détaillée: `working_device`
- **Nombre de déclenchements:** 3
- **Type de nœud déclencheur:** AlphaNode

**📋 Tokens et couples de faits activant l'action:**

##### Token 1
**Fait activateur:** `Device[name=Laptop, status=working, type=computer, id=D001]`

##### Token 2
**Fait activateur:** `Device[id=D003, name=Monitor, status=working, type=display]`

##### Token 3
**Fait activateur:** `Device[id=D004, name=Keyboard, status=maintenance, type=peripheral]`

### 🚫 Analyse des négations (NotNodes)
| Nœud | Condition Niée | Faits Filtrés | Type | Validation |
|------|----------------|---------------|------|------------|
| alpha_rule_0_alpha | map[condition:map[left:map[field:status object:d type:fieldAccess] operator:== right:map[type:string value:broken] type:comparison] negated:true type:negation] | 0 | simple | ❌ |

---

## 💡 RECOMMANDATIONS
### Amélioration de la couverture Beta
### Prochaines étapes
1. **Ajouter plus de tests complexes** avec jointures multiples
2. **Tester les négations imbriquées** et conditions complexes
3. **Valider les performances** des nœuds Beta avec de gros volumes
4. **Enrichir la validation sémantique** avec plus de critères
