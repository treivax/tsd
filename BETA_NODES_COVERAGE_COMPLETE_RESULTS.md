# RAPPORT COMPLET DE COUVERTURE DES NŒUDS BETA
================================================

**📊 Tests exécutés:** 12
**✅ Tests réussis:** 12 (100.0%)
**🧠 Score sémantique moyen:** 100.0%
**📅 Date d'exécution:** 2025-01-18 19:45:00

## 🎯 NŒUDS BETA ANALYSÉS
| Type de Nœud | Tests | Succès | Score Sémantique |
|---------------|--------|--------|------------------|
| ExistsNode | 3 | 3 | 100.0% |
| NotNode | 3 | 3 | 100.0% |
| JoinNode | 8 | 8 | 100.0% |

## 🧪 TEST 1: complex_not_exists_combination
---

### 📋 Informations générales
- **Description:** Test combinaison NOT + EXISTS avec jointures
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/complex_not_exists_combination.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/complex_not_exists_combination.facts`
- **Temps d'exécution:** 3.395667ms
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
- **Texte original:** `{s: Student, c: Course} / s.active == true AND s.grade >= 60 AND NOT (s.grade < 60) AND EXISTS (e: Enrollment / e.student_id == s.id AND e.course_id == c.id AND e.status == "enrolled") ==> qualified_enrolled_student(s.id, c.id)`
- **Action:** qualified_enrolled_student
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - s (Student): primary
  - c (Course): secondary

#### Règle 2
- **Texte original:** `{s: Student, c: Course} / s.active == true AND s.grade >= 60 AND NOT (s.grade < 60) AND EXISTS (e: Enrollment / e.student_id == s.id AND e.course_id == c.id AND e.status == "enrolled") ==> qualified_enrolled_student(s.id, c.id)`
- **Action:** advanced_course_student
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - s (Student): primary
  - c (Course): secondary

#### Règle 3
- **Texte original:** `{s: Student, c: Course} / s.active == true AND s.grade >= 60 AND NOT (s.grade < 60) AND EXISTS (e: Enrollment / e.student_id == s.id AND e.course_id == c.id AND e.status == "enrolled") ==> qualified_enrolled_student(s.id, c.id)`
- **Action:** successful_student
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** simple
- **Variables:**
  - s (Student): primary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Student (type_Student)
│   ├── Course (type_Course)
│   ├── Enrollment (type_Enrollment)
│
├── 🔍 AlphaNodes
│   ├── rule_2_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: s
│
├── 🔗 BetaNodes (Jointures)
│   ├── rule_0_join
│   │   ├── Variables: s ⋈ c
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_1_join
│   │   ├── Variables: s ⋈ c
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_2_terminal
    │   └── Action: successful_student
    ├── rule_0_terminal
    │   └── Action: qualified_enrolled_student
    ├── rule_1_terminal
    │   └── Action: advanced_course_student
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Student[id=S001, name=Alice, grade=85, active=true]
Student[id=S002, name=Bob, grade=45, active=true]
Student[id=S003, name=Charlie, grade=92, active=false]
Student[id=S004, name=Diana, grade=78, active=true]
Course[id=C001, title=Advanced Math, level=advanced, credits=4]
Course[id=C002, title=Basic Physics, level=beginner, credits=2]
Course[id=C003, title=Computer Science, level=advanced, credits=3]
Enrollment[id=EN001, student_id=S001, course_id=C001, status=enrolled]
Enrollment[id=EN002, student_id=S002, course_id=C002, status=failed]
Enrollment[id=EN003, student_id=S003, course_id=C001, status=dropped]
Enrollment[id=EN004, student_id=S004, course_id=C003, status=enrolled]

```

**Total faits:** 11

- **Student:** 4 faits
- **Course:** 3 faits
- **Enrollment:** 4 faits

**📋 Détail des faits parsés:**
1. **Student[S001]** - `Student[id=S001, name=Alice, grade=85, active=true]`
2. **Student[S002]** - `Student[id=S002, name=Bob, grade=45, active=true]`
3. **Student[S003]** - `Student[id=S003, name=Charlie, grade=92, active=false]`
4. **Student[S004]** - `Student[active=true, id=S004, name=Diana, grade=78]`
5. **Course[C001]** - `Course[id=C001, title=Advanced Math, level=advanced, credits=4]`
6. **Course[C002]** - `Course[id=C002, title=Basic Physics, level=beginner, credits=2]`
7. **Course[C003]** - `Course[credits=3, id=C003, title=Computer Science, level=advanced]`
8. **Enrollment[EN001]** - `Enrollment[student_id=S001, course_id=C001, status=enrolled, id=EN001]`
9. **Enrollment[EN002]** - `Enrollment[id=EN002, student_id=S002, course_id=C002, status=failed]`
10. **Enrollment[EN003]** - `Enrollment[id=EN003, student_id=S003, course_id=C001, status=dropped]`
11. **Enrollment[EN004]** - `Enrollment[course_id=C003, status=enrolled, id=EN004, student_id=S004]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| qualified_enrolled_student | 12 | AlphaNode | ❌ |
| advanced_course_student | 12 | AlphaNode | ❌ |
| successful_student | 3 | AlphaNode | ❌ |

#### 🎯 Activation détaillée: `qualified_enrolled_student`
- **Nombre de déclenchements:** 12
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`s`**: Student[S001] - `Student[grade=85, active=true, id=S001, name=Alice]`
- **`c`**: Course[C003] - `Course[id=C003, title=Computer Science, level=advanced, credits=3]`
- **Association:** Student[S001] ⋈ Course[C003]

##### Token combiné 2
- **`s`**: Student[S002] - `Student[id=S002, name=Bob, grade=45, active=true]`
- **`c`**: Course[C003] - `Course[credits=3, id=C003, title=Computer Science, level=advanced]`
- **Association:** Student[S002] ⋈ Course[C003]

##### Token combiné 3
- **`s`**: Student[S003] - `Student[active=false, id=S003, name=Charlie, grade=92]`
- **`c`**: Course[C001] - `Course[id=C001, title=Advanced Math, level=advanced, credits=4]`
- **Association:** Student[S003] ⋈ Course[C001]

##### Token combiné 4
- **`s`**: Student[S004] - `Student[active=true, id=S004, name=Diana, grade=78]`
- **`c`**: Course[C001] - `Course[id=C001, title=Advanced Math, level=advanced, credits=4]`
- **Association:** Student[S004] ⋈ Course[C001]

##### Token combiné 5
- **`s`**: Student[S001] - `Student[id=S001, name=Alice, grade=85, active=true]`
- **`c`**: Course[C002] - `Course[id=C002, title=Basic Physics, level=beginner, credits=2]`
- **Association:** Student[S001] ⋈ Course[C002]

##### Token combiné 6
- **`s`**: Student[S002] - `Student[id=S002, name=Bob, grade=45, active=true]`
- **`c`**: Course[C002] - `Course[id=C002, title=Basic Physics, level=beginner, credits=2]`
- **Association:** Student[S002] ⋈ Course[C002]

##### Token combiné 7
- **`s`**: Student[S003] - `Student[id=S003, name=Charlie, grade=92, active=false]`
- **`c`**: Course[C003] - `Course[id=C003, title=Computer Science, level=advanced, credits=3]`
- **Association:** Student[S003] ⋈ Course[C003]

##### Token combiné 8
- **`s`**: Student[S004] - `Student[id=S004, name=Diana, grade=78, active=true]`
- **`c`**: Course[C003] - `Course[id=C003, title=Computer Science, level=advanced, credits=3]`
- **Association:** Student[S004] ⋈ Course[C003]

##### Token combiné 9
- **`s`**: Student[S001] - `Student[id=S001, name=Alice, grade=85, active=true]`
- **`c`**: Course[C001] - `Course[id=C001, title=Advanced Math, level=advanced, credits=4]`
- **Association:** Student[S001] ⋈ Course[C001]

##### Token combiné 10
- **`s`**: Student[S002] - `Student[grade=45, active=true, id=S002, name=Bob]`
- **`c`**: Course[C001] - `Course[title=Advanced Math, level=advanced, credits=4, id=C001]`
- **Association:** Student[S002] ⋈ Course[C001]

##### Token combiné 11
- **`s`**: Student[S003] - `Student[active=false, id=S003, name=Charlie, grade=92]`
- **`c`**: Course[C002] - `Course[title=Basic Physics, level=beginner, credits=2, id=C002]`
- **Association:** Student[S003] ⋈ Course[C002]

##### Token combiné 12
- **`s`**: Student[S004] - `Student[grade=78, active=true, id=S004, name=Diana]`
- **`c`**: Course[C002] - `Course[id=C002, title=Basic Physics, level=beginner, credits=2]`
- **Association:** Student[S004] ⋈ Course[C002]

#### 🎯 Activation détaillée: `advanced_course_student`
- **Nombre de déclenchements:** 12
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`s`**: Student[S003] - `Student[id=S003, name=Charlie, grade=92, active=false]`
- **`c`**: Course[C002] - `Course[id=C002, title=Basic Physics, level=beginner, credits=2]`
- **Association:** Student[S003] ⋈ Course[C002]

##### Token combiné 2
- **`s`**: Student[S004] - `Student[id=S004, name=Diana, grade=78, active=true]`
- **`c`**: Course[C002] - `Course[id=C002, title=Basic Physics, level=beginner, credits=2]`
- **Association:** Student[S004] ⋈ Course[C002]

##### Token combiné 3
- **`s`**: Student[S001] - `Student[name=Alice, grade=85, active=true, id=S001]`
- **`c`**: Course[C002] - `Course[id=C002, title=Basic Physics, level=beginner, credits=2]`
- **Association:** Student[S001] ⋈ Course[C002]

##### Token combiné 4
- **`s`**: Student[S002] - `Student[id=S002, name=Bob, grade=45, active=true]`
- **`c`**: Course[C002] - `Course[level=beginner, credits=2, id=C002, title=Basic Physics]`
- **Association:** Student[S002] ⋈ Course[C002]

##### Token combiné 5
- **`s`**: Student[S001] - `Student[active=true, id=S001, name=Alice, grade=85]`
- **`c`**: Course[C003] - `Course[id=C003, title=Computer Science, level=advanced, credits=3]`
- **Association:** Student[S001] ⋈ Course[C003]

##### Token combiné 6
- **`s`**: Student[S002] - `Student[active=true, id=S002, name=Bob, grade=45]`
- **`c`**: Course[C003] - `Course[credits=3, id=C003, title=Computer Science, level=advanced]`
- **Association:** Student[S002] ⋈ Course[C003]

##### Token combiné 7
- **`s`**: Student[S003] - `Student[id=S003, name=Charlie, grade=92, active=false]`
- **`c`**: Course[C003] - `Course[id=C003, title=Computer Science, level=advanced, credits=3]`
- **Association:** Student[S003] ⋈ Course[C003]

##### Token combiné 8
- **`s`**: Student[S004] - `Student[grade=78, active=true, id=S004, name=Diana]`
- **`c`**: Course[C003] - `Course[id=C003, title=Computer Science, level=advanced, credits=3]`
- **Association:** Student[S004] ⋈ Course[C003]

##### Token combiné 9
- **`s`**: Student[S004] - `Student[grade=78, active=true, id=S004, name=Diana]`
- **`c`**: Course[C001] - `Course[id=C001, title=Advanced Math, level=advanced, credits=4]`
- **Association:** Student[S004] ⋈ Course[C001]

##### Token combiné 10
- **`s`**: Student[S001] - `Student[id=S001, name=Alice, grade=85, active=true]`
- **`c`**: Course[C001] - `Course[id=C001, title=Advanced Math, level=advanced, credits=4]`
- **Association:** Student[S001] ⋈ Course[C001]

##### Token combiné 11
- **`s`**: Student[S002] - `Student[id=S002, name=Bob, grade=45, active=true]`
- **`c`**: Course[C001] - `Course[id=C001, title=Advanced Math, level=advanced, credits=4]`
- **Association:** Student[S002] ⋈ Course[C001]

##### Token combiné 12
- **`s`**: Student[S003] - `Student[active=false, id=S003, name=Charlie, grade=92]`
- **`c`**: Course[C001] - `Course[id=C001, title=Advanced Math, level=advanced, credits=4]`
- **Association:** Student[S003] ⋈ Course[C001]

#### 🎯 Activation détaillée: `successful_student`
- **Nombre de déclenchements:** 3
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`s`**: Student[S001] - `Student[id=S001, name=Alice, grade=85, active=true]`

##### Token combiné 2
- **`s`**: Student[S002] - `Student[grade=45, active=true, id=S002, name=Bob]`

##### Token combiné 3
- **`s`**: Student[S004] - `Student[name=Diana, grade=78, active=true, id=S004]`

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | s <-> c | 12 | inner | ✅ |
| join_1 | s <-> c | 12 | inner | ✅ |

---

## 🧪 TEST 2: exists_complex_operator
---

### 📋 Informations générales
- **Description:** Test opérateur EXISTS complexe dans nœuds beta
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/exists_complex_operator.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/exists_complex_operator.facts`
- **Temps d'exécution:** 3.179943ms
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
- **Texte original:** `{c: Customer, o: Order} / c.id == o.customer_id AND o.amount > 0 AND EXISTS (p: Payment / p.order_id == o.id AND p.status == "completed") ==> paid_customer_order(c.id, o.id)`
- **Action:** paid_customer_order
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - c (Customer): primary
  - o (Order): secondary

#### Règle 2
- **Texte original:** `{c: Customer, o: Order} / c.id == o.customer_id AND o.amount > 0 AND EXISTS (p: Payment / p.order_id == o.id AND p.status == "completed") ==> paid_customer_order(c.id, o.id)`
- **Action:** credit_customer_order
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - c (Customer): primary
  - o (Order): secondary

#### Règle 3
- **Texte original:** `{c: Customer, o: Order} / c.id == o.customer_id AND o.amount > 0 AND EXISTS (p: Payment / p.order_id == o.id AND p.status == "completed") ==> paid_customer_order(c.id, o.id)`
- **Action:** premium_high_value_customer
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** simple
- **Variables:**
  - c (Customer): primary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Payment (type_Payment)
│   ├── Customer (type_Customer)
│   ├── Order (type_Order)
│
├── 🔍 AlphaNodes
│   ├── rule_2_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: c
│
├── 🔗 BetaNodes (Jointures)
│   ├── rule_0_join
│   │   ├── Variables: c ⋈ o
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_1_join
│   │   ├── Variables: c ⋈ o
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: paid_customer_order
    ├── rule_1_terminal
    │   └── Action: credit_customer_order
    ├── rule_2_terminal
    │   └── Action: premium_high_value_customer
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Customer[id=C001, name=Alice, segment=premium, country=USA]
Customer[id=C002, name=Bob, segment=standard, country=France]
Customer[id=C003, name=Charlie, segment=premium, country=Germany]
Order[id=O001, customer_id=C001, amount=1500, date=2024-01-01]
Order[id=O002, customer_id=C002, amount=150, date=2024-01-02]
Order[id=O003, customer_id=C003, amount=800, date=2024-01-03]
Order[id=O004, customer_id=C001, amount=2000, date=2024-01-04]
Payment[id=PY001, order_id=O001, method=credit_card, status=completed]
Payment[id=PY002, order_id=O002, method=paypal, status=pending]
Payment[id=PY003, order_id=O003, method=credit_card, status=completed]
Payment[id=PY004, order_id=O004, method=bank_transfer, status=completed]

```

**Total faits:** 11

- **Customer:** 3 faits
- **Order:** 4 faits
- **Payment:** 4 faits

**📋 Détail des faits parsés:**
1. **Customer[C001]** - `Customer[country=USA, id=C001, name=Alice, segment=premium]`
2. **Customer[C002]** - `Customer[name=Bob, segment=standard, country=France, id=C002]`
3. **Customer[C003]** - `Customer[id=C003, name=Charlie, segment=premium, country=Germany]`
4. **Order[O001]** - `Order[id=O001, customer_id=C001, amount=1500, date=2024-01-01]`
5. **Order[O002]** - `Order[id=O002, customer_id=C002, amount=150, date=2024-01-02]`
6. **Order[O003]** - `Order[id=O003, customer_id=C003, amount=800, date=2024-01-03]`
7. **Order[O004]** - `Order[id=O004, customer_id=C001, amount=2000, date=2024-01-04]`
8. **Payment[PY001]** - `Payment[method=credit_card, status=completed, id=PY001, order_id=O001]`
9. **Payment[PY002]** - `Payment[id=PY002, order_id=O002, method=paypal, status=pending]`
10. **Payment[PY003]** - `Payment[id=PY003, order_id=O003, method=credit_card, status=completed]`
11. **Payment[PY004]** - `Payment[id=PY004, order_id=O004, method=bank_transfer, status=completed]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| paid_customer_order | 12 | AlphaNode | ❌ |
| credit_customer_order | 12 | AlphaNode | ❌ |
| premium_high_value_customer | 2 | AlphaNode | ❌ |

#### 🎯 Activation détaillée: `paid_customer_order`
- **Nombre de déclenchements:** 12
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`c`**: Customer[C002] - `Customer[id=C002, name=Bob, segment=standard, country=France]`
- **`o`**: Order[O002] - `Order[id=O002, customer_id=C002, amount=150, date=2024-01-02]`
- **Association:** Customer[C002] ⋈ Order[O002]

##### Token combiné 2
- **`c`**: Customer[C001] - `Customer[id=C001, name=Alice, segment=premium, country=USA]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=C003, amount=800, date=2024-01-03]`
- **Association:** Customer[C001] ⋈ Order[O003]

##### Token combiné 3
- **`c`**: Customer[C002] - `Customer[segment=standard, country=France, id=C002, name=Bob]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=C003, amount=800, date=2024-01-03]`
- **Association:** Customer[C002] ⋈ Order[O003]

##### Token combiné 4
- **`c`**: Customer[C001] - `Customer[id=C001, name=Alice, segment=premium, country=USA]`
- **`o`**: Order[O004] - `Order[customer_id=C001, amount=2000, date=2024-01-04, id=O004]`
- **Association:** Customer[C001] ⋈ Order[O004]

##### Token combiné 5
- **`c`**: Customer[C002] - `Customer[id=C002, name=Bob, segment=standard, country=France]`
- **`o`**: Order[O004] - `Order[customer_id=C001, amount=2000, date=2024-01-04, id=O004]`
- **Association:** Customer[C002] ⋈ Order[O004]

##### Token combiné 6
- **`c`**: Customer[C002] - `Customer[id=C002, name=Bob, segment=standard, country=France]`
- **`o`**: Order[O001] - `Order[amount=1500, date=2024-01-01, id=O001, customer_id=C001]`
- **Association:** Customer[C002] ⋈ Order[O001]

##### Token combiné 7
- **`c`**: Customer[C001] - `Customer[id=C001, name=Alice, segment=premium, country=USA]`
- **`o`**: Order[O001] - `Order[date=2024-01-01, id=O001, customer_id=C001, amount=1500]`
- **Association:** Customer[C001] ⋈ Order[O001]

##### Token combiné 8
- **`c`**: Customer[C003] - `Customer[segment=premium, country=Germany, id=C003, name=Charlie]`
- **`o`**: Order[O002] - `Order[id=O002, customer_id=C002, amount=150, date=2024-01-02]`
- **Association:** Customer[C003] ⋈ Order[O002]

##### Token combiné 9
- **`c`**: Customer[C003] - `Customer[id=C003, name=Charlie, segment=premium, country=Germany]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=C003, amount=800, date=2024-01-03]`
- **Association:** Customer[C003] ⋈ Order[O003]

##### Token combiné 10
- **`c`**: Customer[C003] - `Customer[id=C003, name=Charlie, segment=premium, country=Germany]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=C001, amount=2000, date=2024-01-04]`
- **Association:** Customer[C003] ⋈ Order[O004]

##### Token combiné 11
- **`c`**: Customer[C003] - `Customer[id=C003, name=Charlie, segment=premium, country=Germany]`
- **`o`**: Order[O001] - `Order[date=2024-01-01, id=O001, customer_id=C001, amount=1500]`
- **Association:** Customer[C003] ⋈ Order[O001]

##### Token combiné 12
- **`c`**: Customer[C001] - `Customer[id=C001, name=Alice, segment=premium, country=USA]`
- **`o`**: Order[O002] - `Order[id=O002, customer_id=C002, amount=150, date=2024-01-02]`
- **Association:** Customer[C001] ⋈ Order[O002]

#### 🎯 Activation détaillée: `credit_customer_order`
- **Nombre de déclenchements:** 12
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`c`**: Customer[C003] - `Customer[id=C003, name=Charlie, segment=premium, country=Germany]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=C001, amount=2000, date=2024-01-04]`
- **Association:** Customer[C003] ⋈ Order[O004]

##### Token combiné 2
- **`c`**: Customer[C003] - `Customer[id=C003, name=Charlie, segment=premium, country=Germany]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=C001, amount=1500, date=2024-01-01]`
- **Association:** Customer[C003] ⋈ Order[O001]

##### Token combiné 3
- **`c`**: Customer[C001] - `Customer[id=C001, name=Alice, segment=premium, country=USA]`
- **`o`**: Order[O002] - `Order[amount=150, date=2024-01-02, id=O002, customer_id=C002]`
- **Association:** Customer[C001] ⋈ Order[O002]

##### Token combiné 4
- **`c`**: Customer[C002] - `Customer[id=C002, name=Bob, segment=standard, country=France]`
- **`o`**: Order[O002] - `Order[id=O002, customer_id=C002, amount=150, date=2024-01-02]`
- **Association:** Customer[C002] ⋈ Order[O002]

##### Token combiné 5
- **`c`**: Customer[C001] - `Customer[country=USA, id=C001, name=Alice, segment=premium]`
- **`o`**: Order[O003] - `Order[amount=800, date=2024-01-03, id=O003, customer_id=C003]`
- **Association:** Customer[C001] ⋈ Order[O003]

##### Token combiné 6
- **`c`**: Customer[C002] - `Customer[country=France, id=C002, name=Bob, segment=standard]`
- **`o`**: Order[O003] - `Order[customer_id=C003, amount=800, date=2024-01-03, id=O003]`
- **Association:** Customer[C002] ⋈ Order[O003]

##### Token combiné 7
- **`c`**: Customer[C001] - `Customer[id=C001, name=Alice, segment=premium, country=USA]`
- **`o`**: Order[O004] - `Order[customer_id=C001, amount=2000, date=2024-01-04, id=O004]`
- **Association:** Customer[C001] ⋈ Order[O004]

##### Token combiné 8
- **`c`**: Customer[C002] - `Customer[id=C002, name=Bob, segment=standard, country=France]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=C001, amount=2000, date=2024-01-04]`
- **Association:** Customer[C002] ⋈ Order[O004]

##### Token combiné 9
- **`c`**: Customer[C001] - `Customer[id=C001, name=Alice, segment=premium, country=USA]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=C001, amount=1500, date=2024-01-01]`
- **Association:** Customer[C001] ⋈ Order[O001]

##### Token combiné 10
- **`c`**: Customer[C002] - `Customer[id=C002, name=Bob, segment=standard, country=France]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=C001, amount=1500, date=2024-01-01]`
- **Association:** Customer[C002] ⋈ Order[O001]

##### Token combiné 11
- **`c`**: Customer[C003] - `Customer[id=C003, name=Charlie, segment=premium, country=Germany]`
- **`o`**: Order[O002] - `Order[amount=150, date=2024-01-02, id=O002, customer_id=C002]`
- **Association:** Customer[C003] ⋈ Order[O002]

##### Token combiné 12
- **`c`**: Customer[C003] - `Customer[segment=premium, country=Germany, id=C003, name=Charlie]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=C003, amount=800, date=2024-01-03]`
- **Association:** Customer[C003] ⋈ Order[O003]

#### 🎯 Activation détaillée: `premium_high_value_customer`
- **Nombre de déclenchements:** 2
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`c`**: Customer[C001] - `Customer[segment=premium, country=USA, id=C001, name=Alice]`

##### Token combiné 2
- **`c`**: Customer[C003] - `Customer[id=C003, name=Charlie, segment=premium, country=Germany]`

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | c <-> o | 12 | inner | ✅ |
| join_1 | c <-> o | 12 | inner | ✅ |

---

## 🧪 TEST 3: exists_simple
---

### 📋 Informations générales
- **Description:** Test existence simple
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/exists_simple.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/exists_simple.facts`
- **Temps d'exécution:** 626.36µs
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 20.0%
- **Actions valides:** ✅
- **Jointures valides:** ✅
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

### 📜 Règles analysées
#### Règle 1
- **Texte original:** `{p: Person} / EXISTS (o: Order / o.customer_id == p.id) ==> person_has_orders(p.id)`
- **Action:** person_has_orders
- **Type de nœud:** ExistsNode
- **Type sémantique:** existence
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
│   ├── Order (type_Order)
│   ├── Person (type_Person)
│
├── 🔗 BetaNodes (Jointures)
│   ├── rule_0_exists
│   │   └── Type: *rete.ExistsNode
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: person_has_orders
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Person[id=P001, name=Alice]
Person[id=P002, name=Bob]
Order[customer_id=P001, amount=100]

```

**Total faits:** 3

- **Person:** 2 faits
- **Order:** 1 faits

**📋 Détail des faits parsés:**
1. **Person[P001]** - `Person[id=P001, name=Alice]`
2. **Person[P002]** - `Person[id=P002, name=Bob]`
3. **Order[fact_Order_3]** - `Order[customer_id=P001, amount=100, id=fact_Order_3]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| person_has_orders | 1 | AlphaNode | ✅ |

#### 🎯 Activation détaillée: `person_has_orders`
- **Nombre de déclenchements:** 1
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`p`**: Person[P001] - `Person[id=P001, name=Alice]`

### 🎯 Comparaison attendu vs observé
#### Actions
| Action | Attendu | Observé | Statut |
|--------|---------|---------|--------|
| person_has_orders | 1-1 | 1 | ✅ |

#### 📋 TOKENS COMBINÉS ATTENDUS vs OBTENUS

**🎯 Action `person_has_orders`:**
- **Description:** Une personne (Alice) a une commande existante
- **Variables de la règle:** p

**📍 TOKENS COMBINÉS ATTENDUS:**
- **Nombre de tokens attendus:** 1-1
- **Token attendu 1:**
  * `p`: Person[P001] - `Person[name=Alice, id=P001]`

**📊 TOKENS COMBINÉS OBTENUS:**
- **Nombre de tokens obtenus:** 1
- **Token obtenu 1:**
  * `p`: Person[P001] - `Person[id=P001, name=Alice]`

**🎯 RÉSULTAT:** ✅ SUCCÈS
- ✅ Nombre de tokens correct

---

## 🧪 TEST 4: join_and_operator
---

### 📋 Informations générales
- **Description:** Test opérateur AND dans jointures beta
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_and_operator.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_and_operator.facts`
- **Temps d'exécution:** 5.627208ms
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
- **Texte original:** `{p: Person, o: Order} / p.id == o.customer_id AND o.amount > 100 AND p.status == "active" ==> high_value_order(p.id, o.id)`
- **Action:** high_value_order
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - p (Person): primary
  - o (Order): secondary

#### Règle 2
- **Texte original:** `{p: Person, o: Order} / p.id == o.customer_id AND o.amount > 100 AND p.status == "active" ==> high_value_order(p.id, o.id)`
- **Action:** adult_confirmed_order
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - p (Person): primary
  - o (Order): secondary

#### Règle 3
- **Texte original:** `{p: Person, o: Order} / p.id == o.customer_id AND o.amount > 100 AND p.status == "active" ==> high_value_order(p.id, o.id)`
- **Action:** medium_value_order
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
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
│   ├── Order (type_Order)
│   ├── Person (type_Person)
│
├── 🔗 BetaNodes (Jointures)
│   ├── rule_2_join
│   │   ├── Variables: p ⋈ o
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_0_join
│   │   ├── Variables: p ⋈ o
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_1_join
│   │   ├── Variables: p ⋈ o
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: high_value_order
    ├── rule_1_terminal
    │   └── Action: adult_confirmed_order
    ├── rule_2_terminal
    │   └── Action: medium_value_order
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Person[id=P001, name=Alice, age=25, status=active]
Person[id=P002, name=Bob, age=30, status=inactive]
Person[id=P003, name=Charlie, age=16, status=active]
Person[id=P004, name=Diana, age=22, status=active]
Order[id=O001, customer_id=P001, amount=150, status=confirmed]
Order[id=O002, customer_id=P002, amount=75, status=pending]
Order[id=O003, customer_id=P001, amount=200, status=confirmed]
Order[id=O004, customer_id=P003, amount=300, status=confirmed]
Order[id=O005, customer_id=P004, amount=125, status=confirmed]

```

**Total faits:** 9

- **Person:** 4 faits
- **Order:** 5 faits

**📋 Détail des faits parsés:**
1. **Person[P001]** - `Person[name=Alice, age=25, status=active, id=P001]`
2. **Person[P002]** - `Person[id=P002, name=Bob, age=30, status=inactive]`
3. **Person[P003]** - `Person[id=P003, name=Charlie, age=16, status=active]`
4. **Person[P004]** - `Person[status=active, id=P004, name=Diana, age=22]`
5. **Order[O001]** - `Order[amount=150, status=confirmed, id=O001, customer_id=P001]`
6. **Order[O002]** - `Order[id=O002, customer_id=P002, amount=75, status=pending]`
7. **Order[O003]** - `Order[amount=200, status=confirmed, id=O003, customer_id=P001]`
8. **Order[O004]** - `Order[id=O004, customer_id=P003, amount=300, status=confirmed]`
9. **Order[O005]** - `Order[id=O005, customer_id=P004, amount=125, status=confirmed]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| high_value_order | 20 | AlphaNode | ❌ |
| adult_confirmed_order | 20 | AlphaNode | ❌ |
| medium_value_order | 20 | AlphaNode | ❌ |

#### 🎯 Activation détaillée: `high_value_order`
- **Nombre de déclenchements:** 20
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`p`**: Person[P002] - `Person[name=Bob, age=30, status=inactive, id=P002]`
- **`o`**: Order[O002] - `Order[status=pending, id=O002, customer_id=P002, amount=75]`
- **Association:** Person[P002] ⋈ Order[O002]

##### Token combiné 2
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=16, status=active]`
- **`o`**: Order[O005] - `Order[customer_id=P004, amount=125, status=confirmed, id=O005]`
- **Association:** Person[P003] ⋈ Order[O005]

##### Token combiné 3
- **`p`**: Person[P004] - `Person[name=Diana, age=22, status=active, id=P004]`
- **`o`**: Order[O005] - `Order[id=O005, customer_id=P004, amount=125, status=confirmed]`
- **Association:** Person[P004] ⋈ Order[O005]

##### Token combiné 4
- **`p`**: Person[P001] - `Person[name=Alice, age=25, status=active, id=P001]`
- **`o`**: Order[O005] - `Order[id=O005, customer_id=P004, amount=125, status=confirmed]`
- **Association:** Person[P001] ⋈ Order[O005]

##### Token combiné 5
- **`p`**: Person[P002] - `Person[age=30, status=inactive, id=P002, name=Bob]`
- **`o`**: Order[O003] - `Order[amount=200, status=confirmed, id=O003, customer_id=P001]`
- **Association:** Person[P002] ⋈ Order[O003]

##### Token combiné 6
- **`p`**: Person[P001] - `Person[name=Alice, age=25, status=active, id=P001]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P003, amount=300, status=confirmed]`
- **Association:** Person[P001] ⋈ Order[O004]

##### Token combiné 7
- **`p`**: Person[P003] - `Person[status=active, id=P003, name=Charlie, age=16]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P003, amount=300, status=confirmed]`
- **Association:** Person[P003] ⋈ Order[O004]

##### Token combiné 8
- **`p`**: Person[P004] - `Person[status=active, id=P004, name=Diana, age=22]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P003, amount=300, status=confirmed]`
- **Association:** Person[P004] ⋈ Order[O004]

##### Token combiné 9
- **`p`**: Person[P004] - `Person[id=P004, name=Diana, age=22, status=active]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=150, status=confirmed]`
- **Association:** Person[P004] ⋈ Order[O001]

##### Token combiné 10
- **`p`**: Person[P001] - `Person[age=25, status=active, id=P001, name=Alice]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=150, status=confirmed]`
- **Association:** Person[P001] ⋈ Order[O001]

##### Token combiné 11
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=16, status=active]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=150, status=confirmed]`
- **Association:** Person[P003] ⋈ Order[O001]

##### Token combiné 12
- **`p`**: Person[P001] - `Person[id=P001, name=Alice, age=25, status=active]`
- **`o`**: Order[O002] - `Order[customer_id=P002, amount=75, status=pending, id=O002]`
- **Association:** Person[P001] ⋈ Order[O002]

##### Token combiné 13
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=16, status=active]`
- **`o`**: Order[O002] - `Order[amount=75, status=pending, id=O002, customer_id=P002]`
- **Association:** Person[P003] ⋈ Order[O002]

##### Token combiné 14
- **`p`**: Person[P004] - `Person[name=Diana, age=22, status=active, id=P004]`
- **`o`**: Order[O002] - `Order[id=O002, customer_id=P002, amount=75, status=pending]`
- **Association:** Person[P004] ⋈ Order[O002]

##### Token combiné 15
- **`p`**: Person[P002] - `Person[id=P002, name=Bob, age=30, status=inactive]`
- **`o`**: Order[O005] - `Order[id=O005, customer_id=P004, amount=125, status=confirmed]`
- **Association:** Person[P002] ⋈ Order[O005]

##### Token combiné 16
- **`p`**: Person[P001] - `Person[id=P001, name=Alice, age=25, status=active]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=P001, amount=200, status=confirmed]`
- **Association:** Person[P001] ⋈ Order[O003]

##### Token combiné 17
- **`p`**: Person[P003] - `Person[name=Charlie, age=16, status=active, id=P003]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=P001, amount=200, status=confirmed]`
- **Association:** Person[P003] ⋈ Order[O003]

##### Token combiné 18
- **`p`**: Person[P004] - `Person[id=P004, name=Diana, age=22, status=active]`
- **`o`**: Order[O003] - `Order[amount=200, status=confirmed, id=O003, customer_id=P001]`
- **Association:** Person[P004] ⋈ Order[O003]

##### Token combiné 19
- **`p`**: Person[P002] - `Person[id=P002, name=Bob, age=30, status=inactive]`
- **`o`**: Order[O004] - `Order[status=confirmed, id=O004, customer_id=P003, amount=300]`
- **Association:** Person[P002] ⋈ Order[O004]

##### Token combiné 20
- **`p`**: Person[P002] - `Person[id=P002, name=Bob, age=30, status=inactive]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=150, status=confirmed]`
- **Association:** Person[P002] ⋈ Order[O001]

#### 🎯 Activation détaillée: `adult_confirmed_order`
- **Nombre de déclenchements:** 20
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`p`**: Person[P001] - `Person[age=25, status=active, id=P001, name=Alice]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=150, status=confirmed]`
- **Association:** Person[P001] ⋈ Order[O001]

##### Token combiné 2
- **`p`**: Person[P004] - `Person[name=Diana, age=22, status=active, id=P004]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=150, status=confirmed]`
- **Association:** Person[P004] ⋈ Order[O001]

##### Token combiné 3
- **`p`**: Person[P002] - `Person[id=P002, name=Bob, age=30, status=inactive]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=P001, amount=200, status=confirmed]`
- **Association:** Person[P002] ⋈ Order[O003]

##### Token combiné 4
- **`p`**: Person[P001] - `Person[id=P001, name=Alice, age=25, status=active]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P003, amount=300, status=confirmed]`
- **Association:** Person[P001] ⋈ Order[O004]

##### Token combiné 5
- **`p`**: Person[P004] - `Person[name=Diana, age=22, status=active, id=P004]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P003, amount=300, status=confirmed]`
- **Association:** Person[P004] ⋈ Order[O004]

##### Token combiné 6
- **`p`**: Person[P002] - `Person[id=P002, name=Bob, age=30, status=inactive]`
- **`o`**: Order[O005] - `Order[id=O005, customer_id=P004, amount=125, status=confirmed]`
- **Association:** Person[P002] ⋈ Order[O005]

##### Token combiné 7
- **`p`**: Person[P002] - `Person[age=30, status=inactive, id=P002, name=Bob]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=150, status=confirmed]`
- **Association:** Person[P002] ⋈ Order[O001]

##### Token combiné 8
- **`p`**: Person[P003] - `Person[status=active, id=P003, name=Charlie, age=16]`
- **`o`**: Order[O002] - `Order[status=pending, id=O002, customer_id=P002, amount=75]`
- **Association:** Person[P003] ⋈ Order[O002]

##### Token combiné 9
- **`p`**: Person[P004] - `Person[status=active, id=P004, name=Diana, age=22]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=P001, amount=200, status=confirmed]`
- **Association:** Person[P004] ⋈ Order[O003]

##### Token combiné 10
- **`p`**: Person[P001] - `Person[id=P001, name=Alice, age=25, status=active]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=P001, amount=200, status=confirmed]`
- **Association:** Person[P001] ⋈ Order[O003]

##### Token combiné 11
- **`p`**: Person[P002] - `Person[id=P002, name=Bob, age=30, status=inactive]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P003, amount=300, status=confirmed]`
- **Association:** Person[P002] ⋈ Order[O004]

##### Token combiné 12
- **`p`**: Person[P004] - `Person[id=P004, name=Diana, age=22, status=active]`
- **`o`**: Order[O005] - `Order[customer_id=P004, amount=125, status=confirmed, id=O005]`
- **Association:** Person[P004] ⋈ Order[O005]

##### Token combiné 13
- **`p`**: Person[P001] - `Person[id=P001, name=Alice, age=25, status=active]`
- **`o`**: Order[O005] - `Order[id=O005, customer_id=P004, amount=125, status=confirmed]`
- **Association:** Person[P001] ⋈ Order[O005]

##### Token combiné 14
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=16, status=active]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=150, status=confirmed]`
- **Association:** Person[P003] ⋈ Order[O001]

##### Token combiné 15
- **`p`**: Person[P002] - `Person[id=P002, name=Bob, age=30, status=inactive]`
- **`o`**: Order[O002] - `Order[id=O002, customer_id=P002, amount=75, status=pending]`
- **Association:** Person[P002] ⋈ Order[O002]

##### Token combiné 16
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=16, status=active]`
- **`o`**: Order[O004] - `Order[customer_id=P003, amount=300, status=confirmed, id=O004]`
- **Association:** Person[P003] ⋈ Order[O004]

##### Token combiné 17
- **`p`**: Person[P001] - `Person[id=P001, name=Alice, age=25, status=active]`
- **`o`**: Order[O002] - `Order[id=O002, customer_id=P002, amount=75, status=pending]`
- **Association:** Person[P001] ⋈ Order[O002]

##### Token combiné 18
- **`p`**: Person[P004] - `Person[status=active, id=P004, name=Diana, age=22]`
- **`o`**: Order[O002] - `Order[status=pending, id=O002, customer_id=P002, amount=75]`
- **Association:** Person[P004] ⋈ Order[O002]

##### Token combiné 19
- **`p`**: Person[P003] - `Person[age=16, status=active, id=P003, name=Charlie]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=P001, amount=200, status=confirmed]`
- **Association:** Person[P003] ⋈ Order[O003]

##### Token combiné 20
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=16, status=active]`
- **`o`**: Order[O005] - `Order[id=O005, customer_id=P004, amount=125, status=confirmed]`
- **Association:** Person[P003] ⋈ Order[O005]

#### 🎯 Activation détaillée: `medium_value_order`
- **Nombre de déclenchements:** 20
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`p`**: Person[P001] - `Person[status=active, id=P001, name=Alice, age=25]`
- **`o`**: Order[O004] - `Order[customer_id=P003, amount=300, status=confirmed, id=O004]`
- **Association:** Person[P001] ⋈ Order[O004]

##### Token combiné 2
- **`p`**: Person[P002] - `Person[age=30, status=inactive, id=P002, name=Bob]`
- **`o`**: Order[O005] - `Order[amount=125, status=confirmed, id=O005, customer_id=P004]`
- **Association:** Person[P002] ⋈ Order[O005]

##### Token combiné 3
- **`p`**: Person[P002] - `Person[id=P002, name=Bob, age=30, status=inactive]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=150, status=confirmed]`
- **Association:** Person[P002] ⋈ Order[O001]

##### Token combiné 4
- **`p`**: Person[P001] - `Person[status=active, id=P001, name=Alice, age=25]`
- **`o`**: Order[O002] - `Order[status=pending, id=O002, customer_id=P002, amount=75]`
- **Association:** Person[P001] ⋈ Order[O002]

##### Token combiné 5
- **`p`**: Person[P001] - `Person[name=Alice, age=25, status=active, id=P001]`
- **`o`**: Order[O003] - `Order[status=confirmed, id=O003, customer_id=P001, amount=200]`
- **Association:** Person[P001] ⋈ Order[O003]

##### Token combiné 6
- **`p`**: Person[P004] - `Person[id=P004, name=Diana, age=22, status=active]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P003, amount=300, status=confirmed]`
- **Association:** Person[P004] ⋈ Order[O004]

##### Token combiné 7
- **`p`**: Person[P003] - `Person[status=active, id=P003, name=Charlie, age=16]`
- **`o`**: Order[O005] - `Order[customer_id=P004, amount=125, status=confirmed, id=O005]`
- **Association:** Person[P003] ⋈ Order[O005]

##### Token combiné 8
- **`p`**: Person[P001] - `Person[age=25, status=active, id=P001, name=Alice]`
- **`o`**: Order[O001] - `Order[customer_id=P001, amount=150, status=confirmed, id=O001]`
- **Association:** Person[P001] ⋈ Order[O001]

##### Token combiné 9
- **`p`**: Person[P002] - `Person[name=Bob, age=30, status=inactive, id=P002]`
- **`o`**: Order[O002] - `Order[customer_id=P002, amount=75, status=pending, id=O002]`
- **Association:** Person[P002] ⋈ Order[O002]

##### Token combiné 10
- **`p`**: Person[P002] - `Person[status=inactive, id=P002, name=Bob, age=30]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=P001, amount=200, status=confirmed]`
- **Association:** Person[P002] ⋈ Order[O003]

##### Token combiné 11
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=16, status=active]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P003, amount=300, status=confirmed]`
- **Association:** Person[P003] ⋈ Order[O004]

##### Token combiné 12
- **`p`**: Person[P004] - `Person[age=22, status=active, id=P004, name=Diana]`
- **`o`**: Order[O005] - `Order[id=O005, customer_id=P004, amount=125, status=confirmed]`
- **Association:** Person[P004] ⋈ Order[O005]

##### Token combiné 13
- **`p`**: Person[P004] - `Person[age=22, status=active, id=P004, name=Diana]`
- **`o`**: Order[O001] - `Order[amount=150, status=confirmed, id=O001, customer_id=P001]`
- **Association:** Person[P004] ⋈ Order[O001]

##### Token combiné 14
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=16, status=active]`
- **`o`**: Order[O002] - `Order[id=O002, customer_id=P002, amount=75, status=pending]`
- **Association:** Person[P003] ⋈ Order[O002]

##### Token combiné 15
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=16, status=active]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=P001, amount=200, status=confirmed]`
- **Association:** Person[P003] ⋈ Order[O003]

##### Token combiné 16
- **`p`**: Person[P002] - `Person[id=P002, name=Bob, age=30, status=inactive]`
- **`o`**: Order[O004] - `Order[customer_id=P003, amount=300, status=confirmed, id=O004]`
- **Association:** Person[P002] ⋈ Order[O004]

##### Token combiné 17
- **`p`**: Person[P001] - `Person[name=Alice, age=25, status=active, id=P001]`
- **`o`**: Order[O005] - `Order[id=O005, customer_id=P004, amount=125, status=confirmed]`
- **Association:** Person[P001] ⋈ Order[O005]

##### Token combiné 18
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=16, status=active]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=150, status=confirmed]`
- **Association:** Person[P003] ⋈ Order[O001]

##### Token combiné 19
- **`p`**: Person[P004] - `Person[id=P004, name=Diana, age=22, status=active]`
- **`o`**: Order[O002] - `Order[customer_id=P002, amount=75, status=pending, id=O002]`
- **Association:** Person[P004] ⋈ Order[O002]

##### Token combiné 20
- **`p`**: Person[P004] - `Person[status=active, id=P004, name=Diana, age=22]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=P001, amount=200, status=confirmed]`
- **Association:** Person[P004] ⋈ Order[O003]

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | p <-> o | 20 | inner | ✅ |
| join_1 | p <-> o | 20 | inner | ✅ |
| join_2 | p <-> o | 20 | inner | ✅ |

---

## 🧪 TEST 5: join_arithmetic_operators
---

### 📋 Informations générales
- **Description:** Test opérateurs arithmétiques dans jointures beta
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_arithmetic_operators.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_arithmetic_operators.facts`
- **Temps d'exécution:** 22.445052ms
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
- **Texte original:** `{a: Account, t: Transaction} / a.id == t.account_id AND t.amount + a.fees <= a.balance ==> valid_transaction(a.id, t.id)`
- **Action:** valid_transaction
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - a (Account): primary
  - t (Transaction): secondary

#### Règle 2
- **Texte original:** `{a: Account, t: Transaction} / a.id == t.account_id AND t.amount + a.fees <= a.balance ==> valid_transaction(a.id, t.id)`
- **Action:** safe_transaction
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - a (Account): primary
  - t (Transaction): secondary

#### Règle 3
- **Texte original:** `{a: Account, t: Transaction} / a.id == t.account_id AND t.amount + a.fees <= a.balance ==> valid_transaction(a.id, t.id)`
- **Action:** conservative_transaction
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
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
├── 🔗 BetaNodes (Jointures)
│   ├── rule_0_join
│   │   ├── Variables: a ⋈ t
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_1_join
│   │   ├── Variables: a ⋈ t
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_2_join
│   │   ├── Variables: a ⋈ t
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: valid_transaction
    ├── rule_1_terminal
    │   └── Action: safe_transaction
    ├── rule_2_terminal
    │   └── Action: conservative_transaction
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Account[id=ACC001, balance=1000, credit_limit=500, fees=10]
Account[id=ACC002, balance=2500, credit_limit=1000, fees=5]
Account[id=ACC003, balance=500, credit_limit=200, fees=15]
Transaction[id=TXN001, account_id=ACC001, amount=900, type=debit]
Transaction[id=TXN002, account_id=ACC002, amount=1200, type=debit]
Transaction[id=TXN003, account_id=ACC003, amount=600, type=debit]
Transaction[id=TXN004, account_id=ACC001, amount=50, type=debit]

```

**Total faits:** 7

- **Account:** 3 faits
- **Transaction:** 4 faits

**📋 Détail des faits parsés:**
1. **Account[ACC001]** - `Account[id=ACC001, balance=1000, credit_limit=500, fees=10]`
2. **Account[ACC002]** - `Account[id=ACC002, balance=2500, credit_limit=1000, fees=5]`
3. **Account[ACC003]** - `Account[id=ACC003, balance=500, credit_limit=200, fees=15]`
4. **Transaction[TXN001]** - `Transaction[id=TXN001, account_id=ACC001, amount=900, type=debit]`
5. **Transaction[TXN002]** - `Transaction[type=debit, id=TXN002, account_id=ACC002, amount=1200]`
6. **Transaction[TXN003]** - `Transaction[account_id=ACC003, amount=600, type=debit, id=TXN003]`
7. **Transaction[TXN004]** - `Transaction[type=debit, id=TXN004, account_id=ACC001, amount=50]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| safe_transaction | 12 | AlphaNode | ❌ |
| conservative_transaction | 12 | AlphaNode | ❌ |
| valid_transaction | 12 | AlphaNode | ❌ |

#### 🎯 Activation détaillée: `safe_transaction`
- **Nombre de déclenchements:** 12
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`a`**: Account[ACC001] - `Account[id=ACC001, balance=1000, credit_limit=500, fees=10]`
- **`t`**: Transaction[TXN004] - `Transaction[account_id=ACC001, amount=50, type=debit, id=TXN004]`
- **Association:** Account[ACC001] ⋈ Transaction[TXN004]

##### Token combiné 2
- **`a`**: Account[ACC001] - `Account[balance=1000, credit_limit=500, fees=10, id=ACC001]`
- **`t`**: Transaction[TXN001] - `Transaction[id=TXN001, account_id=ACC001, amount=900, type=debit]`
- **Association:** Account[ACC001] ⋈ Transaction[TXN001]

##### Token combiné 3
- **`a`**: Account[ACC002] - `Account[fees=5, id=ACC002, balance=2500, credit_limit=1000]`
- **`t`**: Transaction[TXN002] - `Transaction[id=TXN002, account_id=ACC002, amount=1200, type=debit]`
- **Association:** Account[ACC002] ⋈ Transaction[TXN002]

##### Token combiné 4
- **`a`**: Account[ACC003] - `Account[fees=15, id=ACC003, balance=500, credit_limit=200]`
- **`t`**: Transaction[TXN002] - `Transaction[id=TXN002, account_id=ACC002, amount=1200, type=debit]`
- **Association:** Account[ACC003] ⋈ Transaction[TXN002]

##### Token combiné 5
- **`a`**: Account[ACC001] - `Account[fees=10, id=ACC001, balance=1000, credit_limit=500]`
- **`t`**: Transaction[TXN003] - `Transaction[id=TXN003, account_id=ACC003, amount=600, type=debit]`
- **Association:** Account[ACC001] ⋈ Transaction[TXN003]

##### Token combiné 6
- **`a`**: Account[ACC002] - `Account[id=ACC002, balance=2500, credit_limit=1000, fees=5]`
- **`t`**: Transaction[TXN004] - `Transaction[id=TXN004, account_id=ACC001, amount=50, type=debit]`
- **Association:** Account[ACC002] ⋈ Transaction[TXN004]

##### Token combiné 7
- **`a`**: Account[ACC003] - `Account[id=ACC003, balance=500, credit_limit=200, fees=15]`
- **`t`**: Transaction[TXN004] - `Transaction[amount=50, type=debit, id=TXN004, account_id=ACC001]`
- **Association:** Account[ACC003] ⋈ Transaction[TXN004]

##### Token combiné 8
- **`a`**: Account[ACC002] - `Account[id=ACC002, balance=2500, credit_limit=1000, fees=5]`
- **`t`**: Transaction[TXN001] - `Transaction[amount=900, type=debit, id=TXN001, account_id=ACC001]`
- **Association:** Account[ACC002] ⋈ Transaction[TXN001]

##### Token combiné 9
- **`a`**: Account[ACC003] - `Account[id=ACC003, balance=500, credit_limit=200, fees=15]`
- **`t`**: Transaction[TXN001] - `Transaction[type=debit, id=TXN001, account_id=ACC001, amount=900]`
- **Association:** Account[ACC003] ⋈ Transaction[TXN001]

##### Token combiné 10
- **`a`**: Account[ACC001] - `Account[id=ACC001, balance=1000, credit_limit=500, fees=10]`
- **`t`**: Transaction[TXN002] - `Transaction[id=TXN002, account_id=ACC002, amount=1200, type=debit]`
- **Association:** Account[ACC001] ⋈ Transaction[TXN002]

##### Token combiné 11
- **`a`**: Account[ACC002] - `Account[id=ACC002, balance=2500, credit_limit=1000, fees=5]`
- **`t`**: Transaction[TXN003] - `Transaction[type=debit, id=TXN003, account_id=ACC003, amount=600]`
- **Association:** Account[ACC002] ⋈ Transaction[TXN003]

##### Token combiné 12
- **`a`**: Account[ACC003] - `Account[id=ACC003, balance=500, credit_limit=200, fees=15]`
- **`t`**: Transaction[TXN003] - `Transaction[id=TXN003, account_id=ACC003, amount=600, type=debit]`
- **Association:** Account[ACC003] ⋈ Transaction[TXN003]

#### 🎯 Activation détaillée: `conservative_transaction`
- **Nombre de déclenchements:** 12
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`a`**: Account[ACC001] - `Account[id=ACC001, balance=1000, credit_limit=500, fees=10]`
- **`t`**: Transaction[TXN001] - `Transaction[amount=900, type=debit, id=TXN001, account_id=ACC001]`
- **Association:** Account[ACC001] ⋈ Transaction[TXN001]

##### Token combiné 2
- **`a`**: Account[ACC002] - `Account[id=ACC002, balance=2500, credit_limit=1000, fees=5]`
- **`t`**: Transaction[TXN001] - `Transaction[id=TXN001, account_id=ACC001, amount=900, type=debit]`
- **Association:** Account[ACC002] ⋈ Transaction[TXN001]

##### Token combiné 3
- **`a`**: Account[ACC001] - `Account[balance=1000, credit_limit=500, fees=10, id=ACC001]`
- **`t`**: Transaction[TXN002] - `Transaction[id=TXN002, account_id=ACC002, amount=1200, type=debit]`
- **Association:** Account[ACC001] ⋈ Transaction[TXN002]

##### Token combiné 4
- **`a`**: Account[ACC002] - `Account[id=ACC002, balance=2500, credit_limit=1000, fees=5]`
- **`t`**: Transaction[TXN002] - `Transaction[type=debit, id=TXN002, account_id=ACC002, amount=1200]`
- **Association:** Account[ACC002] ⋈ Transaction[TXN002]

##### Token combiné 5
- **`a`**: Account[ACC003] - `Account[id=ACC003, balance=500, credit_limit=200, fees=15]`
- **`t`**: Transaction[TXN003] - `Transaction[amount=600, type=debit, id=TXN003, account_id=ACC003]`
- **Association:** Account[ACC003] ⋈ Transaction[TXN003]

##### Token combiné 6
- **`a`**: Account[ACC001] - `Account[credit_limit=500, fees=10, id=ACC001, balance=1000]`
- **`t`**: Transaction[TXN004] - `Transaction[id=TXN004, account_id=ACC001, amount=50, type=debit]`
- **Association:** Account[ACC001] ⋈ Transaction[TXN004]

##### Token combiné 7
- **`a`**: Account[ACC002] - `Account[id=ACC002, balance=2500, credit_limit=1000, fees=5]`
- **`t`**: Transaction[TXN004] - `Transaction[amount=50, type=debit, id=TXN004, account_id=ACC001]`
- **Association:** Account[ACC002] ⋈ Transaction[TXN004]

##### Token combiné 8
- **`a`**: Account[ACC003] - `Account[id=ACC003, balance=500, credit_limit=200, fees=15]`
- **`t`**: Transaction[TXN001] - `Transaction[id=TXN001, account_id=ACC001, amount=900, type=debit]`
- **Association:** Account[ACC003] ⋈ Transaction[TXN001]

##### Token combiné 9
- **`a`**: Account[ACC003] - `Account[id=ACC003, balance=500, credit_limit=200, fees=15]`
- **`t`**: Transaction[TXN002] - `Transaction[amount=1200, type=debit, id=TXN002, account_id=ACC002]`
- **Association:** Account[ACC003] ⋈ Transaction[TXN002]

##### Token combiné 10
- **`a`**: Account[ACC001] - `Account[id=ACC001, balance=1000, credit_limit=500, fees=10]`
- **`t`**: Transaction[TXN003] - `Transaction[id=TXN003, account_id=ACC003, amount=600, type=debit]`
- **Association:** Account[ACC001] ⋈ Transaction[TXN003]

##### Token combiné 11
- **`a`**: Account[ACC002] - `Account[balance=2500, credit_limit=1000, fees=5, id=ACC002]`
- **`t`**: Transaction[TXN003] - `Transaction[id=TXN003, account_id=ACC003, amount=600, type=debit]`
- **Association:** Account[ACC002] ⋈ Transaction[TXN003]

##### Token combiné 12
- **`a`**: Account[ACC003] - `Account[id=ACC003, balance=500, credit_limit=200, fees=15]`
- **`t`**: Transaction[TXN004] - `Transaction[id=TXN004, account_id=ACC001, amount=50, type=debit]`
- **Association:** Account[ACC003] ⋈ Transaction[TXN004]

#### 🎯 Activation détaillée: `valid_transaction`
- **Nombre de déclenchements:** 12
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`a`**: Account[ACC002] - `Account[id=ACC002, balance=2500, credit_limit=1000, fees=5]`
- **`t`**: Transaction[TXN003] - `Transaction[id=TXN003, account_id=ACC003, amount=600, type=debit]`
- **Association:** Account[ACC002] ⋈ Transaction[TXN003]

##### Token combiné 2
- **`a`**: Account[ACC003] - `Account[fees=15, id=ACC003, balance=500, credit_limit=200]`
- **`t`**: Transaction[TXN003] - `Transaction[id=TXN003, account_id=ACC003, amount=600, type=debit]`
- **Association:** Account[ACC003] ⋈ Transaction[TXN003]

##### Token combiné 3
- **`a`**: Account[ACC001] - `Account[id=ACC001, balance=1000, credit_limit=500, fees=10]`
- **`t`**: Transaction[TXN001] - `Transaction[id=TXN001, account_id=ACC001, amount=900, type=debit]`
- **Association:** Account[ACC001] ⋈ Transaction[TXN001]

##### Token combiné 4
- **`a`**: Account[ACC002] - `Account[fees=5, id=ACC002, balance=2500, credit_limit=1000]`
- **`t`**: Transaction[TXN001] - `Transaction[account_id=ACC001, amount=900, type=debit, id=TXN001]`
- **Association:** Account[ACC002] ⋈ Transaction[TXN001]

##### Token combiné 5
- **`a`**: Account[ACC003] - `Account[id=ACC003, balance=500, credit_limit=200, fees=15]`
- **`t`**: Transaction[TXN001] - `Transaction[id=TXN001, account_id=ACC001, amount=900, type=debit]`
- **Association:** Account[ACC003] ⋈ Transaction[TXN001]

##### Token combiné 6
- **`a`**: Account[ACC001] - `Account[id=ACC001, balance=1000, credit_limit=500, fees=10]`
- **`t`**: Transaction[TXN004] - `Transaction[id=TXN004, account_id=ACC001, amount=50, type=debit]`
- **Association:** Account[ACC001] ⋈ Transaction[TXN004]

##### Token combiné 7
- **`a`**: Account[ACC002] - `Account[balance=2500, credit_limit=1000, fees=5, id=ACC002]`
- **`t`**: Transaction[TXN004] - `Transaction[amount=50, type=debit, id=TXN004, account_id=ACC001]`
- **Association:** Account[ACC002] ⋈ Transaction[TXN004]

##### Token combiné 8
- **`a`**: Account[ACC003] - `Account[id=ACC003, balance=500, credit_limit=200, fees=15]`
- **`t`**: Transaction[TXN004] - `Transaction[amount=50, type=debit, id=TXN004, account_id=ACC001]`
- **Association:** Account[ACC003] ⋈ Transaction[TXN004]

##### Token combiné 9
- **`a`**: Account[ACC001] - `Account[id=ACC001, balance=1000, credit_limit=500, fees=10]`
- **`t`**: Transaction[TXN002] - `Transaction[type=debit, id=TXN002, account_id=ACC002, amount=1200]`
- **Association:** Account[ACC001] ⋈ Transaction[TXN002]

##### Token combiné 10
- **`a`**: Account[ACC002] - `Account[id=ACC002, balance=2500, credit_limit=1000, fees=5]`
- **`t`**: Transaction[TXN002] - `Transaction[account_id=ACC002, amount=1200, type=debit, id=TXN002]`
- **Association:** Account[ACC002] ⋈ Transaction[TXN002]

##### Token combiné 11
- **`a`**: Account[ACC003] - `Account[id=ACC003, balance=500, credit_limit=200, fees=15]`
- **`t`**: Transaction[TXN002] - `Transaction[id=TXN002, account_id=ACC002, amount=1200, type=debit]`
- **Association:** Account[ACC003] ⋈ Transaction[TXN002]

##### Token combiné 12
- **`a`**: Account[ACC001] - `Account[id=ACC001, balance=1000, credit_limit=500, fees=10]`
- **`t`**: Transaction[TXN003] - `Transaction[id=TXN003, account_id=ACC003, amount=600, type=debit]`
- **Association:** Account[ACC001] ⋈ Transaction[TXN003]

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | a <-> t | 12 | inner | ✅ |
| join_1 | a <-> t | 12 | inner | ✅ |
| join_2 | a <-> t | 12 | inner | ✅ |

---

## 🧪 TEST 6: join_comparison_operators
---

### 📋 Informations générales
- **Description:** Test opérateurs de comparaison dans jointures beta
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_comparison_operators.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_comparison_operators.facts`
- **Temps d'exécution:** 5.404562ms
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
- **Texte original:** `{u: User, a: Activity} / u.id == a.user_id AND a.points > u.score ==> improvement_activity(u.id, a.id)`
- **Action:** improvement_activity
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - u (User): primary
  - a (Activity): secondary

#### Règle 2
- **Texte original:** `{u: User, a: Activity} / u.id == a.user_id AND a.points > u.score ==> improvement_activity(u.id, a.id)`
- **Action:** valid_activity
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - u (User): primary
  - a (Activity): secondary

#### Règle 3
- **Texte original:** `{u: User, a: Activity} / u.id == a.user_id AND a.points > u.score ==> improvement_activity(u.id, a.id)`
- **Action:** low_activity
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - u (User): primary
  - a (Activity): secondary

#### Règle 4
- **Texte original:** `{u: User, a: Activity} / u.id == a.user_id AND a.points > u.score ==> improvement_activity(u.id, a.id)`
- **Action:** different_score_activity
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - u (User): primary
  - a (Activity): secondary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── User (type_User)
│   ├── Activity (type_Activity)
│
├── 🔗 BetaNodes (Jointures)
│   ├── rule_0_join
│   │   ├── Variables: u ⋈ a
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_1_join
│   │   ├── Variables: u ⋈ a
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_2_join
│   │   ├── Variables: u ⋈ a
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_3_join
│   │   ├── Variables: u ⋈ a
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_3_terminal
    │   └── Action: different_score_activity
    ├── rule_0_terminal
    │   └── Action: improvement_activity
    ├── rule_1_terminal
    │   └── Action: valid_activity
    ├── rule_2_terminal
    │   └── Action: low_activity
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
User[id=U001, name=Alice, score=75, created=1700000000]
User[id=U002, name=Bob, score=60, created=1700100000]
User[id=U003, name=Carol, score=90, created=1700200000]
Activity[id=A001, user_id=U001, points=80, timestamp=1700000100]
Activity[id=A002, user_id=U002, points=45, timestamp=1700050000]
Activity[id=A003, user_id=U001, points=75, timestamp=1699900000]
Activity[id=A004, user_id=U003, points=95, timestamp=1700300000]

```

**Total faits:** 7

- **User:** 3 faits
- **Activity:** 4 faits

**📋 Détail des faits parsés:**
1. **User[U001]** - `User[id=U001, name=Alice, score=75, created=1.7e+09]`
2. **User[U002]** - `User[score=60, created=1.7001e+09, id=U002, name=Bob]`
3. **User[U003]** - `User[id=U003, name=Carol, score=90, created=1.7002e+09]`
4. **Activity[A001]** - `Activity[timestamp=1.7000001e+09, id=A001, user_id=U001, points=80]`
5. **Activity[A002]** - `Activity[id=A002, user_id=U002, points=45, timestamp=1.70005e+09]`
6. **Activity[A003]** - `Activity[id=A003, user_id=U001, points=75, timestamp=1.6999e+09]`
7. **Activity[A004]** - `Activity[user_id=U003, points=95, timestamp=1.7003e+09, id=A004]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| different_score_activity | 12 | AlphaNode | ❌ |
| improvement_activity | 12 | AlphaNode | ❌ |
| valid_activity | 12 | AlphaNode | ❌ |
| low_activity | 12 | AlphaNode | ❌ |

#### 🎯 Activation détaillée: `different_score_activity`
- **Nombre de déclenchements:** 12
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`u`**: User[U003] - `User[name=Carol, score=90, created=1.7002e+09, id=U003]`
- **`a`**: Activity[A001] - `Activity[id=A001, user_id=U001, points=80, timestamp=1.7000001e+09]`
- **Association:** User[U003] ⋈ Activity[A001]

##### Token combiné 2
- **`u`**: User[U001] - `User[id=U001, name=Alice, score=75, created=1.7e+09]`
- **`a`**: Activity[A002] - `Activity[timestamp=1.70005e+09, id=A002, user_id=U002, points=45]`
- **Association:** User[U001] ⋈ Activity[A002]

##### Token combiné 3
- **`u`**: User[U002] - `User[name=Bob, score=60, created=1.7001e+09, id=U002]`
- **`a`**: Activity[A002] - `Activity[id=A002, user_id=U002, points=45, timestamp=1.70005e+09]`
- **Association:** User[U002] ⋈ Activity[A002]

##### Token combiné 4
- **`u`**: User[U003] - `User[id=U003, name=Carol, score=90, created=1.7002e+09]`
- **`a`**: Activity[A003] - `Activity[id=A003, user_id=U001, points=75, timestamp=1.6999e+09]`
- **Association:** User[U003] ⋈ Activity[A003]

##### Token combiné 5
- **`u`**: User[U003] - `User[score=90, created=1.7002e+09, id=U003, name=Carol]`
- **`a`**: Activity[A004] - `Activity[id=A004, user_id=U003, points=95, timestamp=1.7003e+09]`
- **Association:** User[U003] ⋈ Activity[A004]

##### Token combiné 6
- **`u`**: User[U001] - `User[id=U001, name=Alice, score=75, created=1.7e+09]`
- **`a`**: Activity[A001] - `Activity[user_id=U001, points=80, timestamp=1.7000001e+09, id=A001]`
- **Association:** User[U001] ⋈ Activity[A001]

##### Token combiné 7
- **`u`**: User[U002] - `User[score=60, created=1.7001e+09, id=U002, name=Bob]`
- **`a`**: Activity[A001] - `Activity[user_id=U001, points=80, timestamp=1.7000001e+09, id=A001]`
- **Association:** User[U002] ⋈ Activity[A001]

##### Token combiné 8
- **`u`**: User[U003] - `User[score=90, created=1.7002e+09, id=U003, name=Carol]`
- **`a`**: Activity[A002] - `Activity[id=A002, user_id=U002, points=45, timestamp=1.70005e+09]`
- **Association:** User[U003] ⋈ Activity[A002]

##### Token combiné 9
- **`u`**: User[U001] - `User[id=U001, name=Alice, score=75, created=1.7e+09]`
- **`a`**: Activity[A003] - `Activity[id=A003, user_id=U001, points=75, timestamp=1.6999e+09]`
- **Association:** User[U001] ⋈ Activity[A003]

##### Token combiné 10
- **`u`**: User[U002] - `User[id=U002, name=Bob, score=60, created=1.7001e+09]`
- **`a`**: Activity[A003] - `Activity[id=A003, user_id=U001, points=75, timestamp=1.6999e+09]`
- **Association:** User[U002] ⋈ Activity[A003]

##### Token combiné 11
- **`u`**: User[U001] - `User[id=U001, name=Alice, score=75, created=1.7e+09]`
- **`a`**: Activity[A004] - `Activity[points=95, timestamp=1.7003e+09, id=A004, user_id=U003]`
- **Association:** User[U001] ⋈ Activity[A004]

##### Token combiné 12
- **`u`**: User[U002] - `User[score=60, created=1.7001e+09, id=U002, name=Bob]`
- **`a`**: Activity[A004] - `Activity[id=A004, user_id=U003, points=95, timestamp=1.7003e+09]`
- **Association:** User[U002] ⋈ Activity[A004]

#### 🎯 Activation détaillée: `improvement_activity`
- **Nombre de déclenchements:** 12
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`u`**: User[U001] - `User[id=U001, name=Alice, score=75, created=1.7e+09]`
- **`a`**: Activity[A001] - `Activity[points=80, timestamp=1.7000001e+09, id=A001, user_id=U001]`
- **Association:** User[U001] ⋈ Activity[A001]

##### Token combiné 2
- **`u`**: User[U002] - `User[id=U002, name=Bob, score=60, created=1.7001e+09]`
- **`a`**: Activity[A001] - `Activity[id=A001, user_id=U001, points=80, timestamp=1.7000001e+09]`
- **Association:** User[U002] ⋈ Activity[A001]

##### Token combiné 3
- **`u`**: User[U003] - `User[id=U003, name=Carol, score=90, created=1.7002e+09]`
- **`a`**: Activity[A001] - `Activity[user_id=U001, points=80, timestamp=1.7000001e+09, id=A001]`
- **Association:** User[U003] ⋈ Activity[A001]

##### Token combiné 4
- **`u`**: User[U001] - `User[id=U001, name=Alice, score=75, created=1.7e+09]`
- **`a`**: Activity[A003] - `Activity[points=75, timestamp=1.6999e+09, id=A003, user_id=U001]`
- **Association:** User[U001] ⋈ Activity[A003]

##### Token combiné 5
- **`u`**: User[U002] - `User[score=60, created=1.7001e+09, id=U002, name=Bob]`
- **`a`**: Activity[A003] - `Activity[id=A003, user_id=U001, points=75, timestamp=1.6999e+09]`
- **Association:** User[U002] ⋈ Activity[A003]

##### Token combiné 6
- **`u`**: User[U003] - `User[id=U003, name=Carol, score=90, created=1.7002e+09]`
- **`a`**: Activity[A003] - `Activity[id=A003, user_id=U001, points=75, timestamp=1.6999e+09]`
- **Association:** User[U003] ⋈ Activity[A003]

##### Token combiné 7
- **`u`**: User[U002] - `User[id=U002, name=Bob, score=60, created=1.7001e+09]`
- **`a`**: Activity[A004] - `Activity[id=A004, user_id=U003, points=95, timestamp=1.7003e+09]`
- **Association:** User[U002] ⋈ Activity[A004]

##### Token combiné 8
- **`u`**: User[U003] - `User[id=U003, name=Carol, score=90, created=1.7002e+09]`
- **`a`**: Activity[A004] - `Activity[user_id=U003, points=95, timestamp=1.7003e+09, id=A004]`
- **Association:** User[U003] ⋈ Activity[A004]

##### Token combiné 9
- **`u`**: User[U001] - `User[name=Alice, score=75, created=1.7e+09, id=U001]`
- **`a`**: Activity[A002] - `Activity[user_id=U002, points=45, timestamp=1.70005e+09, id=A002]`
- **Association:** User[U001] ⋈ Activity[A002]

##### Token combiné 10
- **`u`**: User[U002] - `User[id=U002, name=Bob, score=60, created=1.7001e+09]`
- **`a`**: Activity[A002] - `Activity[id=A002, user_id=U002, points=45, timestamp=1.70005e+09]`
- **Association:** User[U002] ⋈ Activity[A002]

##### Token combiné 11
- **`u`**: User[U003] - `User[id=U003, name=Carol, score=90, created=1.7002e+09]`
- **`a`**: Activity[A002] - `Activity[id=A002, user_id=U002, points=45, timestamp=1.70005e+09]`
- **Association:** User[U003] ⋈ Activity[A002]

##### Token combiné 12
- **`u`**: User[U001] - `User[id=U001, name=Alice, score=75, created=1.7e+09]`
- **`a`**: Activity[A004] - `Activity[id=A004, user_id=U003, points=95, timestamp=1.7003e+09]`
- **Association:** User[U001] ⋈ Activity[A004]

#### 🎯 Activation détaillée: `valid_activity`
- **Nombre de déclenchements:** 12
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`u`**: User[U002] - `User[id=U002, name=Bob, score=60, created=1.7001e+09]`
- **`a`**: Activity[A002] - `Activity[points=45, timestamp=1.70005e+09, id=A002, user_id=U002]`
- **Association:** User[U002] ⋈ Activity[A002]

##### Token combiné 2
- **`u`**: User[U003] - `User[id=U003, name=Carol, score=90, created=1.7002e+09]`
- **`a`**: Activity[A002] - `Activity[id=A002, user_id=U002, points=45, timestamp=1.70005e+09]`
- **Association:** User[U003] ⋈ Activity[A002]

##### Token combiné 3
- **`u`**: User[U001] - `User[score=75, created=1.7e+09, id=U001, name=Alice]`
- **`a`**: Activity[A003] - `Activity[points=75, timestamp=1.6999e+09, id=A003, user_id=U001]`
- **Association:** User[U001] ⋈ Activity[A003]

##### Token combiné 4
- **`u`**: User[U001] - `User[id=U001, name=Alice, score=75, created=1.7e+09]`
- **`a`**: Activity[A004] - `Activity[id=A004, user_id=U003, points=95, timestamp=1.7003e+09]`
- **Association:** User[U001] ⋈ Activity[A004]

##### Token combiné 5
- **`u`**: User[U001] - `User[score=75, created=1.7e+09, id=U001, name=Alice]`
- **`a`**: Activity[A001] - `Activity[points=80, timestamp=1.7000001e+09, id=A001, user_id=U001]`
- **Association:** User[U001] ⋈ Activity[A001]

##### Token combiné 6
- **`u`**: User[U001] - `User[id=U001, name=Alice, score=75, created=1.7e+09]`
- **`a`**: Activity[A002] - `Activity[id=A002, user_id=U002, points=45, timestamp=1.70005e+09]`
- **Association:** User[U001] ⋈ Activity[A002]

##### Token combiné 7
- **`u`**: User[U002] - `User[id=U002, name=Bob, score=60, created=1.7001e+09]`
- **`a`**: Activity[A003] - `Activity[id=A003, user_id=U001, points=75, timestamp=1.6999e+09]`
- **Association:** User[U002] ⋈ Activity[A003]

##### Token combiné 8
- **`u`**: User[U003] - `User[id=U003, name=Carol, score=90, created=1.7002e+09]`
- **`a`**: Activity[A003] - `Activity[id=A003, user_id=U001, points=75, timestamp=1.6999e+09]`
- **Association:** User[U003] ⋈ Activity[A003]

##### Token combiné 9
- **`u`**: User[U002] - `User[name=Bob, score=60, created=1.7001e+09, id=U002]`
- **`a`**: Activity[A004] - `Activity[id=A004, user_id=U003, points=95, timestamp=1.7003e+09]`
- **Association:** User[U002] ⋈ Activity[A004]

##### Token combiné 10
- **`u`**: User[U003] - `User[id=U003, name=Carol, score=90, created=1.7002e+09]`
- **`a`**: Activity[A004] - `Activity[id=A004, user_id=U003, points=95, timestamp=1.7003e+09]`
- **Association:** User[U003] ⋈ Activity[A004]

##### Token combiné 11
- **`u`**: User[U002] - `User[id=U002, name=Bob, score=60, created=1.7001e+09]`
- **`a`**: Activity[A001] - `Activity[user_id=U001, points=80, timestamp=1.7000001e+09, id=A001]`
- **Association:** User[U002] ⋈ Activity[A001]

##### Token combiné 12
- **`u`**: User[U003] - `User[id=U003, name=Carol, score=90, created=1.7002e+09]`
- **`a`**: Activity[A001] - `Activity[id=A001, user_id=U001, points=80, timestamp=1.7000001e+09]`
- **Association:** User[U003] ⋈ Activity[A001]

#### 🎯 Activation détaillée: `low_activity`
- **Nombre de déclenchements:** 12
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`u`**: User[U001] - `User[id=U001, name=Alice, score=75, created=1.7e+09]`
- **`a`**: Activity[A002] - `Activity[id=A002, user_id=U002, points=45, timestamp=1.70005e+09]`
- **Association:** User[U001] ⋈ Activity[A002]

##### Token combiné 2
- **`u`**: User[U002] - `User[id=U002, name=Bob, score=60, created=1.7001e+09]`
- **`a`**: Activity[A002] - `Activity[timestamp=1.70005e+09, id=A002, user_id=U002, points=45]`
- **Association:** User[U002] ⋈ Activity[A002]

##### Token combiné 3
- **`u`**: User[U003] - `User[id=U003, name=Carol, score=90, created=1.7002e+09]`
- **`a`**: Activity[A002] - `Activity[id=A002, user_id=U002, points=45, timestamp=1.70005e+09]`
- **Association:** User[U003] ⋈ Activity[A002]

##### Token combiné 4
- **`u`**: User[U003] - `User[id=U003, name=Carol, score=90, created=1.7002e+09]`
- **`a`**: Activity[A004] - `Activity[timestamp=1.7003e+09, id=A004, user_id=U003, points=95]`
- **Association:** User[U003] ⋈ Activity[A004]

##### Token combiné 5
- **`u`**: User[U001] - `User[id=U001, name=Alice, score=75, created=1.7e+09]`
- **`a`**: Activity[A001] - `Activity[timestamp=1.7000001e+09, id=A001, user_id=U001, points=80]`
- **Association:** User[U001] ⋈ Activity[A001]

##### Token combiné 6
- **`u`**: User[U002] - `User[id=U002, name=Bob, score=60, created=1.7001e+09]`
- **`a`**: Activity[A001] - `Activity[id=A001, user_id=U001, points=80, timestamp=1.7000001e+09]`
- **Association:** User[U002] ⋈ Activity[A001]

##### Token combiné 7
- **`u`**: User[U003] - `User[id=U003, name=Carol, score=90, created=1.7002e+09]`
- **`a`**: Activity[A001] - `Activity[id=A001, user_id=U001, points=80, timestamp=1.7000001e+09]`
- **Association:** User[U003] ⋈ Activity[A001]

##### Token combiné 8
- **`u`**: User[U003] - `User[score=90, created=1.7002e+09, id=U003, name=Carol]`
- **`a`**: Activity[A003] - `Activity[timestamp=1.6999e+09, id=A003, user_id=U001, points=75]`
- **Association:** User[U003] ⋈ Activity[A003]

##### Token combiné 9
- **`u`**: User[U001] - `User[id=U001, name=Alice, score=75, created=1.7e+09]`
- **`a`**: Activity[A003] - `Activity[user_id=U001, points=75, timestamp=1.6999e+09, id=A003]`
- **Association:** User[U001] ⋈ Activity[A003]

##### Token combiné 10
- **`u`**: User[U002] - `User[id=U002, name=Bob, score=60, created=1.7001e+09]`
- **`a`**: Activity[A003] - `Activity[timestamp=1.6999e+09, id=A003, user_id=U001, points=75]`
- **Association:** User[U002] ⋈ Activity[A003]

##### Token combiné 11
- **`u`**: User[U001] - `User[id=U001, name=Alice, score=75, created=1.7e+09]`
- **`a`**: Activity[A004] - `Activity[id=A004, user_id=U003, points=95, timestamp=1.7003e+09]`
- **Association:** User[U001] ⋈ Activity[A004]

##### Token combiné 12
- **`u`**: User[U002] - `User[created=1.7001e+09, id=U002, name=Bob, score=60]`
- **`a`**: Activity[A004] - `Activity[id=A004, user_id=U003, points=95, timestamp=1.7003e+09]`
- **Association:** User[U002] ⋈ Activity[A004]

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | u <-> a | 12 | inner | ✅ |
| join_1 | u <-> a | 12 | inner | ✅ |
| join_2 | u <-> a | 12 | inner | ✅ |
| join_3 | u <-> a | 12 | inner | ✅ |

---

## 🧪 TEST 7: join_in_contains_operators
---

### 📋 Informations générales
- **Description:** Test opérateurs IN et CONTAINS dans jointures beta
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_in_contains_operators.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_in_contains_operators.facts`
- **Temps d'exécution:** 5.085716ms
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
- **Texte original:** `{p: Product, r: Review} / p.id == r.product_id AND r.status IN ["approved", "verified"] ==> approved_review(p.id, r.id)`
- **Action:** approved_review
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - p (Product): primary
  - r (Review): secondary

#### Règle 2
- **Texte original:** `{p: Product, r: Review} / p.id == r.product_id AND r.status IN ["approved", "verified"] ==> approved_review(p.id, r.id)`
- **Action:** premium_product_review
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - p (Product): primary
  - r (Review): secondary

#### Règle 3
- **Texte original:** `{p: Product, r: Review} / p.id == r.product_id AND r.status IN ["approved", "verified"] ==> approved_review(p.id, r.id)`
- **Action:** tech_high_rating
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - p (Product): primary
  - r (Review): secondary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Product (type_Product)
│   ├── Review (type_Review)
│
├── 🔗 BetaNodes (Jointures)
│   ├── rule_0_join
│   │   ├── Variables: p ⋈ r
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_1_join
│   │   ├── Variables: p ⋈ r
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_2_join
│   │   ├── Variables: p ⋈ r
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: approved_review
    ├── rule_1_terminal
    │   └── Action: premium_product_review
    ├── rule_2_terminal
    │   └── Action: tech_high_rating
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Product[id=PROD001, name=Laptop, categories=electronics, keywords=premium computer high-end]
Product[id=PROD002, name=Mouse, categories=accessories, keywords=basic cheap]
Product[id=PROD003, name=Phone, categories=tech, keywords=premium mobile smartphone]
Review[id=R001, product_id=PROD001, rating=5, status=approved]
Review[id=R002, product_id=PROD002, rating=3, status=pending]
Review[id=R003, product_id=PROD001, rating=4, status=verified]
Review[id=R004, product_id=PROD003, rating=5, status=approved]

```

**Total faits:** 7

- **Product:** 3 faits
- **Review:** 4 faits

**📋 Détail des faits parsés:**
1. **Product[PROD001]** - `Product[id=PROD001, name=Laptop, categories=electronics, keywords=premium computer high-end]`
2. **Product[PROD002]** - `Product[id=PROD002, name=Mouse, categories=accessories, keywords=basic cheap]`
3. **Product[PROD003]** - `Product[id=PROD003, name=Phone, categories=tech, keywords=premium mobile smartphone]`
4. **Review[R001]** - `Review[id=R001, product_id=PROD001, rating=5, status=approved]`
5. **Review[R002]** - `Review[rating=3, status=pending, id=R002, product_id=PROD002]`
6. **Review[R003]** - `Review[id=R003, product_id=PROD001, rating=4, status=verified]`
7. **Review[R004]** - `Review[product_id=PROD003, rating=5, status=approved, id=R004]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| premium_product_review | 12 | AlphaNode | ❌ |
| tech_high_rating | 12 | AlphaNode | ❌ |
| approved_review | 12 | AlphaNode | ❌ |

#### 🎯 Activation détaillée: `premium_product_review`
- **Nombre de déclenchements:** 12
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`p`**: Product[PROD003] - `Product[id=PROD003, name=Phone, categories=tech, keywords=premium mobile smartphone]`
- **`r`**: Review[R001] - `Review[id=R001, product_id=PROD001, rating=5, status=approved]`
- **Association:** Product[PROD003] ⋈ Review[R001]

##### Token combiné 2
- **`p`**: Product[PROD001] - `Product[id=PROD001, name=Laptop, categories=electronics, keywords=premium computer high-end]`
- **`r`**: Review[R002] - `Review[id=R002, product_id=PROD002, rating=3, status=pending]`
- **Association:** Product[PROD001] ⋈ Review[R002]

##### Token combiné 3
- **`p`**: Product[PROD003] - `Product[keywords=premium mobile smartphone, id=PROD003, name=Phone, categories=tech]`
- **`r`**: Review[R003] - `Review[status=verified, id=R003, product_id=PROD001, rating=4]`
- **Association:** Product[PROD003] ⋈ Review[R003]

##### Token combiné 4
- **`p`**: Product[PROD002] - `Product[id=PROD002, name=Mouse, categories=accessories, keywords=basic cheap]`
- **`r`**: Review[R003] - `Review[product_id=PROD001, rating=4, status=verified, id=R003]`
- **Association:** Product[PROD002] ⋈ Review[R003]

##### Token combiné 5
- **`p`**: Product[PROD001] - `Product[name=Laptop, categories=electronics, keywords=premium computer high-end, id=PROD001]`
- **`r`**: Review[R004] - `Review[status=approved, id=R004, product_id=PROD003, rating=5]`
- **Association:** Product[PROD001] ⋈ Review[R004]

##### Token combiné 6
- **`p`**: Product[PROD001] - `Product[name=Laptop, categories=electronics, keywords=premium computer high-end, id=PROD001]`
- **`r`**: Review[R001] - `Review[rating=5, status=approved, id=R001, product_id=PROD001]`
- **Association:** Product[PROD001] ⋈ Review[R001]

##### Token combiné 7
- **`p`**: Product[PROD003] - `Product[id=PROD003, name=Phone, categories=tech, keywords=premium mobile smartphone]`
- **`r`**: Review[R002] - `Review[rating=3, status=pending, id=R002, product_id=PROD002]`
- **Association:** Product[PROD003] ⋈ Review[R002]

##### Token combiné 8
- **`p`**: Product[PROD002] - `Product[id=PROD002, name=Mouse, categories=accessories, keywords=basic cheap]`
- **`r`**: Review[R002] - `Review[id=R002, product_id=PROD002, rating=3, status=pending]`
- **Association:** Product[PROD002] ⋈ Review[R002]

##### Token combiné 9
- **`p`**: Product[PROD001] - `Product[id=PROD001, name=Laptop, categories=electronics, keywords=premium computer high-end]`
- **`r`**: Review[R003] - `Review[id=R003, product_id=PROD001, rating=4, status=verified]`
- **Association:** Product[PROD001] ⋈ Review[R003]

##### Token combiné 10
- **`p`**: Product[PROD002] - `Product[categories=accessories, keywords=basic cheap, id=PROD002, name=Mouse]`
- **`r`**: Review[R004] - `Review[id=R004, product_id=PROD003, rating=5, status=approved]`
- **Association:** Product[PROD002] ⋈ Review[R004]

##### Token combiné 11
- **`p`**: Product[PROD003] - `Product[keywords=premium mobile smartphone, id=PROD003, name=Phone, categories=tech]`
- **`r`**: Review[R004] - `Review[id=R004, product_id=PROD003, rating=5, status=approved]`
- **Association:** Product[PROD003] ⋈ Review[R004]

##### Token combiné 12
- **`p`**: Product[PROD002] - `Product[id=PROD002, name=Mouse, categories=accessories, keywords=basic cheap]`
- **`r`**: Review[R001] - `Review[product_id=PROD001, rating=5, status=approved, id=R001]`
- **Association:** Product[PROD002] ⋈ Review[R001]

#### 🎯 Activation détaillée: `tech_high_rating`
- **Nombre de déclenchements:** 12
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`p`**: Product[PROD003] - `Product[id=PROD003, name=Phone, categories=tech, keywords=premium mobile smartphone]`
- **`r`**: Review[R001] - `Review[status=approved, id=R001, product_id=PROD001, rating=5]`
- **Association:** Product[PROD003] ⋈ Review[R001]

##### Token combiné 2
- **`p`**: Product[PROD003] - `Product[id=PROD003, name=Phone, categories=tech, keywords=premium mobile smartphone]`
- **`r`**: Review[R002] - `Review[product_id=PROD002, rating=3, status=pending, id=R002]`
- **Association:** Product[PROD003] ⋈ Review[R002]

##### Token combiné 3
- **`p`**: Product[PROD001] - `Product[id=PROD001, name=Laptop, categories=electronics, keywords=premium computer high-end]`
- **`r`**: Review[R003] - `Review[id=R003, product_id=PROD001, rating=4, status=verified]`
- **Association:** Product[PROD001] ⋈ Review[R003]

##### Token combiné 4
- **`p`**: Product[PROD002] - `Product[keywords=basic cheap, id=PROD002, name=Mouse, categories=accessories]`
- **`r`**: Review[R003] - `Review[id=R003, product_id=PROD001, rating=4, status=verified]`
- **Association:** Product[PROD002] ⋈ Review[R003]

##### Token combiné 5
- **`p`**: Product[PROD003] - `Product[keywords=premium mobile smartphone, id=PROD003, name=Phone, categories=tech]`
- **`r`**: Review[R004] - `Review[id=R004, product_id=PROD003, rating=5, status=approved]`
- **Association:** Product[PROD003] ⋈ Review[R004]

##### Token combiné 6
- **`p`**: Product[PROD001] - `Product[categories=electronics, keywords=premium computer high-end, id=PROD001, name=Laptop]`
- **`r`**: Review[R001] - `Review[id=R001, product_id=PROD001, rating=5, status=approved]`
- **Association:** Product[PROD001] ⋈ Review[R001]

##### Token combiné 7
- **`p`**: Product[PROD002] - `Product[id=PROD002, name=Mouse, categories=accessories, keywords=basic cheap]`
- **`r`**: Review[R001] - `Review[id=R001, product_id=PROD001, rating=5, status=approved]`
- **Association:** Product[PROD002] ⋈ Review[R001]

##### Token combiné 8
- **`p`**: Product[PROD001] - `Product[id=PROD001, name=Laptop, categories=electronics, keywords=premium computer high-end]`
- **`r`**: Review[R002] - `Review[id=R002, product_id=PROD002, rating=3, status=pending]`
- **Association:** Product[PROD001] ⋈ Review[R002]

##### Token combiné 9
- **`p`**: Product[PROD002] - `Product[keywords=basic cheap, id=PROD002, name=Mouse, categories=accessories]`
- **`r`**: Review[R002] - `Review[rating=3, status=pending, id=R002, product_id=PROD002]`
- **Association:** Product[PROD002] ⋈ Review[R002]

##### Token combiné 10
- **`p`**: Product[PROD003] - `Product[id=PROD003, name=Phone, categories=tech, keywords=premium mobile smartphone]`
- **`r`**: Review[R003] - `Review[rating=4, status=verified, id=R003, product_id=PROD001]`
- **Association:** Product[PROD003] ⋈ Review[R003]

##### Token combiné 11
- **`p`**: Product[PROD001] - `Product[id=PROD001, name=Laptop, categories=electronics, keywords=premium computer high-end]`
- **`r`**: Review[R004] - `Review[id=R004, product_id=PROD003, rating=5, status=approved]`
- **Association:** Product[PROD001] ⋈ Review[R004]

##### Token combiné 12
- **`p`**: Product[PROD002] - `Product[id=PROD002, name=Mouse, categories=accessories, keywords=basic cheap]`
- **`r`**: Review[R004] - `Review[id=R004, product_id=PROD003, rating=5, status=approved]`
- **Association:** Product[PROD002] ⋈ Review[R004]

#### 🎯 Activation détaillée: `approved_review`
- **Nombre de déclenchements:** 12
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`p`**: Product[PROD002] - `Product[name=Mouse, categories=accessories, keywords=basic cheap, id=PROD002]`
- **`r`**: Review[R002] - `Review[id=R002, product_id=PROD002, rating=3, status=pending]`
- **Association:** Product[PROD002] ⋈ Review[R002]

##### Token combiné 2
- **`p`**: Product[PROD003] - `Product[categories=tech, keywords=premium mobile smartphone, id=PROD003, name=Phone]`
- **`r`**: Review[R002] - `Review[id=R002, product_id=PROD002, rating=3, status=pending]`
- **Association:** Product[PROD003] ⋈ Review[R002]

##### Token combiné 3
- **`p`**: Product[PROD001] - `Product[id=PROD001, name=Laptop, categories=electronics, keywords=premium computer high-end]`
- **`r`**: Review[R003] - `Review[id=R003, product_id=PROD001, rating=4, status=verified]`
- **Association:** Product[PROD001] ⋈ Review[R003]

##### Token combiné 4
- **`p`**: Product[PROD002] - `Product[keywords=basic cheap, id=PROD002, name=Mouse, categories=accessories]`
- **`r`**: Review[R003] - `Review[id=R003, product_id=PROD001, rating=4, status=verified]`
- **Association:** Product[PROD002] ⋈ Review[R003]

##### Token combiné 5
- **`p`**: Product[PROD003] - `Product[id=PROD003, name=Phone, categories=tech, keywords=premium mobile smartphone]`
- **`r`**: Review[R003] - `Review[id=R003, product_id=PROD001, rating=4, status=verified]`
- **Association:** Product[PROD003] ⋈ Review[R003]

##### Token combiné 6
- **`p`**: Product[PROD001] - `Product[id=PROD001, name=Laptop, categories=electronics, keywords=premium computer high-end]`
- **`r`**: Review[R001] - `Review[id=R001, product_id=PROD001, rating=5, status=approved]`
- **Association:** Product[PROD001] ⋈ Review[R001]

##### Token combiné 7
- **`p`**: Product[PROD002] - `Product[id=PROD002, name=Mouse, categories=accessories, keywords=basic cheap]`
- **`r`**: Review[R001] - `Review[id=R001, product_id=PROD001, rating=5, status=approved]`
- **Association:** Product[PROD002] ⋈ Review[R001]

##### Token combiné 8
- **`p`**: Product[PROD003] - `Product[id=PROD003, name=Phone, categories=tech, keywords=premium mobile smartphone]`
- **`r`**: Review[R001] - `Review[id=R001, product_id=PROD001, rating=5, status=approved]`
- **Association:** Product[PROD003] ⋈ Review[R001]

##### Token combiné 9
- **`p`**: Product[PROD001] - `Product[id=PROD001, name=Laptop, categories=electronics, keywords=premium computer high-end]`
- **`r`**: Review[R004] - `Review[product_id=PROD003, rating=5, status=approved, id=R004]`
- **Association:** Product[PROD001] ⋈ Review[R004]

##### Token combiné 10
- **`p`**: Product[PROD002] - `Product[id=PROD002, name=Mouse, categories=accessories, keywords=basic cheap]`
- **`r`**: Review[R004] - `Review[status=approved, id=R004, product_id=PROD003, rating=5]`
- **Association:** Product[PROD002] ⋈ Review[R004]

##### Token combiné 11
- **`p`**: Product[PROD003] - `Product[id=PROD003, name=Phone, categories=tech, keywords=premium mobile smartphone]`
- **`r`**: Review[R004] - `Review[id=R004, product_id=PROD003, rating=5, status=approved]`
- **Association:** Product[PROD003] ⋈ Review[R004]

##### Token combiné 12
- **`p`**: Product[PROD001] - `Product[name=Laptop, categories=electronics, keywords=premium computer high-end, id=PROD001]`
- **`r`**: Review[R002] - `Review[id=R002, product_id=PROD002, rating=3, status=pending]`
- **Association:** Product[PROD001] ⋈ Review[R002]

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | p <-> r | 12 | inner | ✅ |
| join_1 | p <-> r | 12 | inner | ✅ |
| join_2 | p <-> r | 12 | inner | ✅ |

---

## 🧪 TEST 8: join_multi_variable_complex
---

### 📋 Informations générales
- **Description:** Test jointures complexes multi-variables
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_multi_variable_complex.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_multi_variable_complex.facts`
- **Temps d'exécution:** 14.013513ms
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
- **Texte original:** `{u: User, t: Team, task: Task} / u.id == t.manager_id AND t.id == task.team_id AND task.priority == "high" ==> manager_high_priority_task(u.id, t.id, task.id)`
- **Action:** manager_high_priority_task
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** multi-variable
- **Variables:**
  - u (User): primary
  - t (Team): secondary
  - task (Task): secondary

#### Règle 2
- **Texte original:** `{u: User, t: Team, task: Task} / u.id == t.manager_id AND t.id == task.team_id AND task.priority == "high" ==> manager_high_priority_task(u.id, t.id, task.id)`
- **Action:** affordable_task_assignment
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** multi-variable
- **Variables:**
  - u (User): primary
  - t (Team): secondary
  - task (Task): secondary

#### Règle 3
- **Texte original:** `{u: User, t: Team, task: Task} / u.id == t.manager_id AND t.id == task.team_id AND task.priority == "high" ==> manager_high_priority_task(u.id, t.id, task.id)`
- **Action:** lead_complex_task
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** multi-variable
- **Variables:**
  - u (User): primary
  - t (Team): secondary
  - task (Task): secondary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── User (type_User)
│   ├── Team (type_Team)
│   ├── Task (type_Task)
│
├── 🔗 BetaNodes (Jointures)
│   ├── rule_0_join
│   │   ├── Variables: u ⋈ t ⋈ task
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_1_join
│   │   ├── Variables: u ⋈ t ⋈ task
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_2_join
│   │   ├── Variables: u ⋈ t ⋈ task
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_2_terminal
    │   └── Action: lead_complex_task
    ├── rule_0_terminal
    │   └── Action: manager_high_priority_task
    ├── rule_1_terminal
    │   └── Action: affordable_task_assignment
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
User[id=U001, name=Alice, role=manager, team_id=T001]
User[id=U002, name=Bob, role=lead, team_id=T001]
User[id=U003, name=Carol, role=developer, team_id=T002]
Team[id=T001, name=Alpha, budget=10000, manager_id=U001]
Team[id=T002, name=Beta, budget=5000, manager_id=U003]
Task[id=TASK001, assignee_id=U002, team_id=T001, priority=high, effort=50]
Task[id=TASK002, assignee_id=U003, team_id=T002, priority=medium, effort=20]
Task[id=TASK003, assignee_id=U001, team_id=T001, priority=high, effort=30]

```

**Total faits:** 8

- **Task:** 3 faits
- **User:** 3 faits
- **Team:** 2 faits

**📋 Détail des faits parsés:**
1. **User[U001]** - `User[id=U001, name=Alice, role=manager, team_id=T001]`
2. **User[U002]** - `User[name=Bob, role=lead, team_id=T001, id=U002]`
3. **User[U003]** - `User[id=U003, name=Carol, role=developer, team_id=T002]`
4. **Team[T001]** - `Team[id=T001, name=Alpha, budget=10000, manager_id=U001]`
5. **Team[T002]** - `Team[manager_id=U003, id=T002, name=Beta, budget=5000]`
6. **Task[TASK001]** - `Task[id=TASK001, assignee_id=U002, team_id=T001, priority=high, effort=50]`
7. **Task[TASK002]** - `Task[assignee_id=U003, team_id=T002, priority=medium, effort=20, id=TASK002]`
8. **Task[TASK003]** - `Task[assignee_id=U001, team_id=T001, priority=high, effort=30, id=TASK003]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| manager_high_priority_task | 15 | AlphaNode | ❌ |
| affordable_task_assignment | 15 | AlphaNode | ❌ |
| lead_complex_task | 15 | AlphaNode | ❌ |

#### 🎯 Activation détaillée: `manager_high_priority_task`
- **Nombre de déclenchements:** 15
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`u`**: User[U001] - `User[name=Alice, role=manager, team_id=T001, id=U001]`
- **`t`**: Task[TASK002] - `Task[id=TASK002, assignee_id=U003, team_id=T002, priority=medium, effort=20]`
- **Association:** User[U001] ⋈ Task[TASK002]

##### Token combiné 2
- **`u`**: User[U002] - `User[team_id=T001, id=U002, name=Bob, role=lead]`
- **`t`**: Task[TASK002] - `Task[id=TASK002, assignee_id=U003, team_id=T002, priority=medium, effort=20]`
- **Association:** User[U002] ⋈ Task[TASK002]

##### Token combiné 3
- **`u`**: User[U002] - `User[id=U002, name=Bob, role=lead, team_id=T001]`
- **`t`**: Team[T001] - `Team[id=T001, name=Alpha, budget=10000, manager_id=U001]`
- **Association:** User[U002] ⋈ Team[T001]

##### Token combiné 4
- **`u`**: User[U002] - `User[id=U002, name=Bob, role=lead, team_id=T001]`
- **`t`**: Team[T002] - `Team[id=T002, name=Beta, budget=5000, manager_id=U003]`
- **Association:** User[U002] ⋈ Team[T002]

##### Token combiné 5
- **`u`**: User[U003] - `User[id=U003, name=Carol, role=developer, team_id=T002]`
- **`t`**: Task[TASK001] - `Task[effort=50, id=TASK001, assignee_id=U002, team_id=T001, priority=high]`
- **Association:** User[U003] ⋈ Task[TASK001]

##### Token combiné 6
- **`u`**: User[U003] - `User[team_id=T002, id=U003, name=Carol, role=developer]`
- **`t`**: Task[TASK003] - `Task[assignee_id=U001, team_id=T001, priority=high, effort=30, id=TASK003]`
- **Association:** User[U003] ⋈ Task[TASK003]

##### Token combiné 7
- **`u`**: User[U001] - `User[id=U001, name=Alice, role=manager, team_id=T001]`
- **`t`**: Task[TASK001] - `Task[assignee_id=U002, team_id=T001, priority=high, effort=50, id=TASK001]`
- **Association:** User[U001] ⋈ Task[TASK001]

##### Token combiné 8
- **`u`**: User[U002] - `User[team_id=T001, id=U002, name=Bob, role=lead]`
- **`t`**: Task[TASK001] - `Task[assignee_id=U002, team_id=T001, priority=high, effort=50, id=TASK001]`
- **Association:** User[U002] ⋈ Task[TASK001]

##### Token combiné 9
- **`u`**: User[U001] - `User[team_id=T001, id=U001, name=Alice, role=manager]`
- **`t`**: Task[TASK003] - `Task[id=TASK003, assignee_id=U001, team_id=T001, priority=high, effort=30]`
- **Association:** User[U001] ⋈ Task[TASK003]

##### Token combiné 10
- **`u`**: User[U002] - `User[id=U002, name=Bob, role=lead, team_id=T001]`
- **`t`**: Task[TASK003] - `Task[priority=high, effort=30, id=TASK003, assignee_id=U001, team_id=T001]`
- **Association:** User[U002] ⋈ Task[TASK003]

##### Token combiné 11
- **`u`**: User[U003] - `User[name=Carol, role=developer, team_id=T002, id=U003]`
- **`t`**: Team[T001] - `Team[id=T001, name=Alpha, budget=10000, manager_id=U001]`
- **Association:** User[U003] ⋈ Team[T001]

##### Token combiné 12
- **`u`**: User[U003] - `User[id=U003, name=Carol, role=developer, team_id=T002]`
- **`t`**: Team[T002] - `Team[id=T002, name=Beta, budget=5000, manager_id=U003]`
- **Association:** User[U003] ⋈ Team[T002]

##### Token combiné 13
- **`u`**: User[U003] - `User[id=U003, name=Carol, role=developer, team_id=T002]`
- **`t`**: Task[TASK002] - `Task[assignee_id=U003, team_id=T002, priority=medium, effort=20, id=TASK002]`
- **Association:** User[U003] ⋈ Task[TASK002]

##### Token combiné 14
- **`u`**: User[U001] - `User[team_id=T001, id=U001, name=Alice, role=manager]`
- **`t`**: Team[T001] - `Team[id=T001, name=Alpha, budget=10000, manager_id=U001]`
- **Association:** User[U001] ⋈ Team[T001]

##### Token combiné 15
- **`u`**: User[U001] - `User[team_id=T001, id=U001, name=Alice, role=manager]`
- **`t`**: Team[T002] - `Team[budget=5000, manager_id=U003, id=T002, name=Beta]`
- **Association:** User[U001] ⋈ Team[T002]

#### 🎯 Activation détaillée: `affordable_task_assignment`
- **Nombre de déclenchements:** 15
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`u`**: User[U003] - `User[team_id=T002, id=U003, name=Carol, role=developer]`
- **`t`**: Team[T001] - `Team[id=T001, name=Alpha, budget=10000, manager_id=U001]`
- **Association:** User[U003] ⋈ Team[T001]

##### Token combiné 2
- **`u`**: User[U003] - `User[id=U003, name=Carol, role=developer, team_id=T002]`
- **`t`**: Task[TASK001] - `Task[team_id=T001, priority=high, effort=50, id=TASK001, assignee_id=U002]`
- **Association:** User[U003] ⋈ Task[TASK001]

##### Token combiné 3
- **`u`**: User[U001] - `User[role=manager, team_id=T001, id=U001, name=Alice]`
- **`t`**: Task[TASK003] - `Task[priority=high, effort=30, id=TASK003, assignee_id=U001, team_id=T001]`
- **Association:** User[U001] ⋈ Task[TASK003]

##### Token combiné 4
- **`u`**: User[U002] - `User[name=Bob, role=lead, team_id=T001, id=U002]`
- **`t`**: Team[T001] - `Team[id=T001, name=Alpha, budget=10000, manager_id=U001]`
- **Association:** User[U002] ⋈ Team[T001]

##### Token combiné 5
- **`u`**: User[U001] - `User[id=U001, name=Alice, role=manager, team_id=T001]`
- **`t`**: Team[T002] - `Team[id=T002, name=Beta, budget=5000, manager_id=U003]`
- **Association:** User[U001] ⋈ Team[T002]

##### Token combiné 6
- **`u`**: User[U003] - `User[id=U003, name=Carol, role=developer, team_id=T002]`
- **`t`**: Task[TASK002] - `Task[effort=20, id=TASK002, assignee_id=U003, team_id=T002, priority=medium]`
- **Association:** User[U003] ⋈ Task[TASK002]

##### Token combiné 7
- **`u`**: User[U002] - `User[id=U002, name=Bob, role=lead, team_id=T001]`
- **`t`**: Task[TASK003] - `Task[id=TASK003, assignee_id=U001, team_id=T001, priority=high, effort=30]`
- **Association:** User[U002] ⋈ Task[TASK003]

##### Token combiné 8
- **`u`**: User[U001] - `User[role=manager, team_id=T001, id=U001, name=Alice]`
- **`t`**: Team[T001] - `Team[id=T001, name=Alpha, budget=10000, manager_id=U001]`
- **Association:** User[U001] ⋈ Team[T001]

##### Token combiné 9
- **`u`**: User[U002] - `User[id=U002, name=Bob, role=lead, team_id=T001]`
- **`t`**: Team[T002] - `Team[name=Beta, budget=5000, manager_id=U003, id=T002]`
- **Association:** User[U002] ⋈ Team[T002]

##### Token combiné 10
- **`u`**: User[U002] - `User[role=lead, team_id=T001, id=U002, name=Bob]`
- **`t`**: Task[TASK001] - `Task[team_id=T001, priority=high, effort=50, id=TASK001, assignee_id=U002]`
- **Association:** User[U002] ⋈ Task[TASK001]

##### Token combiné 11
- **`u`**: User[U001] - `User[role=manager, team_id=T001, id=U001, name=Alice]`
- **`t`**: Task[TASK002] - `Task[id=TASK002, assignee_id=U003, team_id=T002, priority=medium, effort=20]`
- **Association:** User[U001] ⋈ Task[TASK002]

##### Token combiné 12
- **`u`**: User[U003] - `User[team_id=T002, id=U003, name=Carol, role=developer]`
- **`t`**: Team[T002] - `Team[budget=5000, manager_id=U003, id=T002, name=Beta]`
- **Association:** User[U003] ⋈ Team[T002]

##### Token combiné 13
- **`u`**: User[U001] - `User[id=U001, name=Alice, role=manager, team_id=T001]`
- **`t`**: Task[TASK001] - `Task[priority=high, effort=50, id=TASK001, assignee_id=U002, team_id=T001]`
- **Association:** User[U001] ⋈ Task[TASK001]

##### Token combiné 14
- **`u`**: User[U002] - `User[id=U002, name=Bob, role=lead, team_id=T001]`
- **`t`**: Task[TASK002] - `Task[assignee_id=U003, team_id=T002, priority=medium, effort=20, id=TASK002]`
- **Association:** User[U002] ⋈ Task[TASK002]

##### Token combiné 15
- **`u`**: User[U003] - `User[id=U003, name=Carol, role=developer, team_id=T002]`
- **`t`**: Task[TASK003] - `Task[team_id=T001, priority=high, effort=30, id=TASK003, assignee_id=U001]`
- **Association:** User[U003] ⋈ Task[TASK003]

#### 🎯 Activation détaillée: `lead_complex_task`
- **Nombre de déclenchements:** 15
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`u`**: User[U001] - `User[id=U001, name=Alice, role=manager, team_id=T001]`
- **`t`**: Task[TASK002] - `Task[effort=20, id=TASK002, assignee_id=U003, team_id=T002, priority=medium]`
- **Association:** User[U001] ⋈ Task[TASK002]

##### Token combiné 2
- **`u`**: User[U001] - `User[team_id=T001, id=U001, name=Alice, role=manager]`
- **`t`**: Team[T002] - `Team[name=Beta, budget=5000, manager_id=U003, id=T002]`
- **Association:** User[U001] ⋈ Team[T002]

##### Token combiné 3
- **`u`**: User[U003] - `User[id=U003, name=Carol, role=developer, team_id=T002]`
- **`t`**: Team[T002] - `Team[id=T002, name=Beta, budget=5000, manager_id=U003]`
- **Association:** User[U003] ⋈ Team[T002]

##### Token combiné 4
- **`u`**: User[U002] - `User[id=U002, name=Bob, role=lead, team_id=T001]`
- **`t`**: Task[TASK001] - `Task[id=TASK001, assignee_id=U002, team_id=T001, priority=high, effort=50]`
- **Association:** User[U002] ⋈ Task[TASK001]

##### Token combiné 5
- **`u`**: User[U001] - `User[name=Alice, role=manager, team_id=T001, id=U001]`
- **`t`**: Task[TASK003] - `Task[team_id=T001, priority=high, effort=30, id=TASK003, assignee_id=U001]`
- **Association:** User[U001] ⋈ Task[TASK003]

##### Token combiné 6
- **`u`**: User[U002] - `User[team_id=T001, id=U002, name=Bob, role=lead]`
- **`t`**: Team[T001] - `Team[id=T001, name=Alpha, budget=10000, manager_id=U001]`
- **Association:** User[U002] ⋈ Team[T001]

##### Token combiné 7
- **`u`**: User[U001] - `User[name=Alice, role=manager, team_id=T001, id=U001]`
- **`t`**: Task[TASK001] - `Task[assignee_id=U002, team_id=T001, priority=high, effort=50, id=TASK001]`
- **Association:** User[U001] ⋈ Task[TASK001]

##### Token combiné 8
- **`u`**: User[U003] - `User[role=developer, team_id=T002, id=U003, name=Carol]`
- **`t`**: Task[TASK002] - `Task[effort=20, id=TASK002, assignee_id=U003, team_id=T002, priority=medium]`
- **Association:** User[U003] ⋈ Task[TASK002]

##### Token combiné 9
- **`u`**: User[U002] - `User[name=Bob, role=lead, team_id=T001, id=U002]`
- **`t`**: Task[TASK003] - `Task[id=TASK003, assignee_id=U001, team_id=T001, priority=high, effort=30]`
- **Association:** User[U002] ⋈ Task[TASK003]

##### Token combiné 10
- **`u`**: User[U003] - `User[team_id=T002, id=U003, name=Carol, role=developer]`
- **`t`**: Team[T001] - `Team[id=T001, name=Alpha, budget=10000, manager_id=U001]`
- **Association:** User[U003] ⋈ Team[T001]

##### Token combiné 11
- **`u`**: User[U001] - `User[id=U001, name=Alice, role=manager, team_id=T001]`
- **`t`**: Team[T001] - `Team[manager_id=U001, id=T001, name=Alpha, budget=10000]`
- **Association:** User[U001] ⋈ Team[T001]

##### Token combiné 12
- **`u`**: User[U002] - `User[id=U002, name=Bob, role=lead, team_id=T001]`
- **`t`**: Task[TASK002] - `Task[effort=20, id=TASK002, assignee_id=U003, team_id=T002, priority=medium]`
- **Association:** User[U002] ⋈ Task[TASK002]

##### Token combiné 13
- **`u`**: User[U003] - `User[id=U003, name=Carol, role=developer, team_id=T002]`
- **`t`**: Task[TASK003] - `Task[priority=high, effort=30, id=TASK003, assignee_id=U001, team_id=T001]`
- **Association:** User[U003] ⋈ Task[TASK003]

##### Token combiné 14
- **`u`**: User[U002] - `User[role=lead, team_id=T001, id=U002, name=Bob]`
- **`t`**: Team[T002] - `Team[id=T002, name=Beta, budget=5000, manager_id=U003]`
- **Association:** User[U002] ⋈ Team[T002]

##### Token combiné 15
- **`u`**: User[U003] - `User[name=Carol, role=developer, team_id=T002, id=U003]`
- **`t`**: Task[TASK001] - `Task[priority=high, effort=50, id=TASK001, assignee_id=U002, team_id=T001]`
- **Association:** User[U003] ⋈ Task[TASK001]

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | u <-> t, t <-> task | 15 | inner | ✅ |
| join_1 | u <-> t, t <-> task | 15 | inner | ✅ |
| join_2 | u <-> t, t <-> task | 15 | inner | ✅ |

---

## 🧪 TEST 9: join_or_operator
---

### 📋 Informations générales
- **Description:** Test opérateur OR dans jointures beta
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_or_operator.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_or_operator.facts`
- **Temps d'exécution:** 6.43077ms
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
- **Texte original:** `{p: Person, o: Order} / p.id == o.customer_id AND (o.amount > 500 OR o.urgent == true) AND p.priority != "" ==> priority_order(p.id, o.id)`
- **Action:** priority_order
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - p (Person): primary
  - o (Order): secondary

#### Règle 2
- **Texte original:** `{p: Person, o: Order} / p.id == o.customer_id AND (o.amount > 500 OR o.urgent == true) AND p.priority != "" ==> priority_order(p.id, o.id)`
- **Action:** special_customer_order
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - p (Person): primary
  - o (Order): secondary

#### Règle 3
- **Texte original:** `{p: Person, o: Order} / p.id == o.customer_id AND (o.amount > 500 OR o.urgent == true) AND p.priority != "" ==> priority_order(p.id, o.id)`
- **Action:** standard_order
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
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
├── 🔗 BetaNodes (Jointures)
│   ├── rule_0_join
│   │   ├── Variables: p ⋈ o
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_1_join
│   │   ├── Variables: p ⋈ o
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_2_join
│   │   ├── Variables: p ⋈ o
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: priority_order
    ├── rule_1_terminal
    │   └── Action: special_customer_order
    ├── rule_2_terminal
    │   └── Action: standard_order
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Person[id=P001, name=Alice, age=25, priority=high]
Person[id=P002, name=Bob, age=67, priority=normal]
Person[id=P003, name=Charlie, age=35, priority=low]
Person[id=P004, name=Diana, age=19, priority=normal]
Order[id=O001, customer_id=P001, amount=600, urgent=false]
Order[id=O002, customer_id=P002, amount=300, urgent=false]
Order[id=O003, customer_id=P003, amount=100, urgent=true]
Order[id=O004, customer_id=P004, amount=50, urgent=false]

```

**Total faits:** 8

- **Person:** 4 faits
- **Order:** 4 faits

**📋 Détail des faits parsés:**
1. **Person[P001]** - `Person[id=P001, name=Alice, age=25, priority=high]`
2. **Person[P002]** - `Person[name=Bob, age=67, priority=normal, id=P002]`
3. **Person[P003]** - `Person[id=P003, name=Charlie, age=35, priority=low]`
4. **Person[P004]** - `Person[id=P004, name=Diana, age=19, priority=normal]`
5. **Order[O001]** - `Order[amount=600, urgent=false, id=O001, customer_id=P001]`
6. **Order[O002]** - `Order[id=O002, customer_id=P002, amount=300, urgent=false]`
7. **Order[O003]** - `Order[id=O003, customer_id=P003, amount=100, urgent=true]`
8. **Order[O004]** - `Order[id=O004, customer_id=P004, amount=50, urgent=false]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| special_customer_order | 16 | AlphaNode | ❌ |
| standard_order | 16 | AlphaNode | ❌ |
| priority_order | 16 | AlphaNode | ❌ |

#### 🎯 Activation détaillée: `special_customer_order`
- **Nombre de déclenchements:** 16
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`p`**: Person[P001] - `Person[age=25, priority=high, id=P001, name=Alice]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=600, urgent=false]`
- **Association:** Person[P001] ⋈ Order[O001]

##### Token combiné 2
- **`p`**: Person[P003] - `Person[name=Charlie, age=35, priority=low, id=P003]`
- **`o`**: Order[O002] - `Order[customer_id=P002, amount=300, urgent=false, id=O002]`
- **Association:** Person[P003] ⋈ Order[O002]

##### Token combiné 3
- **`p`**: Person[P002] - `Person[id=P002, name=Bob, age=67, priority=normal]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=P003, amount=100, urgent=true]`
- **Association:** Person[P002] ⋈ Order[O003]

##### Token combiné 4
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=35, priority=low]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P004, amount=50, urgent=false]`
- **Association:** Person[P003] ⋈ Order[O004]

##### Token combiné 5
- **`p`**: Person[P002] - `Person[priority=normal, id=P002, name=Bob, age=67]`
- **`o`**: Order[O001] - `Order[customer_id=P001, amount=600, urgent=false, id=O001]`
- **Association:** Person[P002] ⋈ Order[O001]

##### Token combiné 6
- **`p`**: Person[P004] - `Person[id=P004, name=Diana, age=19, priority=normal]`
- **`o`**: Order[O002] - `Order[id=O002, customer_id=P002, amount=300, urgent=false]`
- **Association:** Person[P004] ⋈ Order[O002]

##### Token combiné 7
- **`p`**: Person[P001] - `Person[id=P001, name=Alice, age=25, priority=high]`
- **`o`**: Order[O003] - `Order[customer_id=P003, amount=100, urgent=true, id=O003]`
- **Association:** Person[P001] ⋈ Order[O003]

##### Token combiné 8
- **`p`**: Person[P004] - `Person[id=P004, name=Diana, age=19, priority=normal]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P004, amount=50, urgent=false]`
- **Association:** Person[P004] ⋈ Order[O004]

##### Token combiné 9
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=35, priority=low]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=600, urgent=false]`
- **Association:** Person[P003] ⋈ Order[O001]

##### Token combiné 10
- **`p`**: Person[P001] - `Person[name=Alice, age=25, priority=high, id=P001]`
- **`o`**: Order[O002] - `Order[customer_id=P002, amount=300, urgent=false, id=O002]`
- **Association:** Person[P001] ⋈ Order[O002]

##### Token combiné 11
- **`p`**: Person[P004] - `Person[priority=normal, id=P004, name=Diana, age=19]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=P003, amount=100, urgent=true]`
- **Association:** Person[P004] ⋈ Order[O003]

##### Token combiné 12
- **`p`**: Person[P001] - `Person[priority=high, id=P001, name=Alice, age=25]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P004, amount=50, urgent=false]`
- **Association:** Person[P001] ⋈ Order[O004]

##### Token combiné 13
- **`p`**: Person[P004] - `Person[priority=normal, id=P004, name=Diana, age=19]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=600, urgent=false]`
- **Association:** Person[P004] ⋈ Order[O001]

##### Token combiné 14
- **`p`**: Person[P002] - `Person[id=P002, name=Bob, age=67, priority=normal]`
- **`o`**: Order[O002] - `Order[customer_id=P002, amount=300, urgent=false, id=O002]`
- **Association:** Person[P002] ⋈ Order[O002]

##### Token combiné 15
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=35, priority=low]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=P003, amount=100, urgent=true]`
- **Association:** Person[P003] ⋈ Order[O003]

##### Token combiné 16
- **`p`**: Person[P002] - `Person[age=67, priority=normal, id=P002, name=Bob]`
- **`o`**: Order[O004] - `Order[amount=50, urgent=false, id=O004, customer_id=P004]`
- **Association:** Person[P002] ⋈ Order[O004]

#### 🎯 Activation détaillée: `standard_order`
- **Nombre de déclenchements:** 16
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`p`**: Person[P002] - `Person[id=P002, name=Bob, age=67, priority=normal]`
- **`o`**: Order[O001] - `Order[urgent=false, id=O001, customer_id=P001, amount=600]`
- **Association:** Person[P002] ⋈ Order[O001]

##### Token combiné 2
- **`p`**: Person[P002] - `Person[name=Bob, age=67, priority=normal, id=P002]`
- **`o`**: Order[O003] - `Order[customer_id=P003, amount=100, urgent=true, id=O003]`
- **Association:** Person[P002] ⋈ Order[O003]

##### Token combiné 3
- **`p`**: Person[P002] - `Person[id=P002, name=Bob, age=67, priority=normal]`
- **`o`**: Order[O004] - `Order[amount=50, urgent=false, id=O004, customer_id=P004]`
- **Association:** Person[P002] ⋈ Order[O004]

##### Token combiné 4
- **`p`**: Person[P001] - `Person[id=P001, name=Alice, age=25, priority=high]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=600, urgent=false]`
- **Association:** Person[P001] ⋈ Order[O001]

##### Token combiné 5
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=35, priority=low]`
- **`o`**: Order[O002] - `Order[id=O002, customer_id=P002, amount=300, urgent=false]`
- **Association:** Person[P003] ⋈ Order[O002]

##### Token combiné 6
- **`p`**: Person[P004] - `Person[id=P004, name=Diana, age=19, priority=normal]`
- **`o`**: Order[O002] - `Order[customer_id=P002, amount=300, urgent=false, id=O002]`
- **Association:** Person[P004] ⋈ Order[O002]

##### Token combiné 7
- **`p`**: Person[P001] - `Person[id=P001, name=Alice, age=25, priority=high]`
- **`o`**: Order[O003] - `Order[customer_id=P003, amount=100, urgent=true, id=O003]`
- **Association:** Person[P001] ⋈ Order[O003]

##### Token combiné 8
- **`p`**: Person[P001] - `Person[priority=high, id=P001, name=Alice, age=25]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P004, amount=50, urgent=false]`
- **Association:** Person[P001] ⋈ Order[O004]

##### Token combiné 9
- **`p`**: Person[P004] - `Person[id=P004, name=Diana, age=19, priority=normal]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P004, amount=50, urgent=false]`
- **Association:** Person[P004] ⋈ Order[O004]

##### Token combiné 10
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=35, priority=low]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=600, urgent=false]`
- **Association:** Person[P003] ⋈ Order[O001]

##### Token combiné 11
- **`p`**: Person[P004] - `Person[id=P004, name=Diana, age=19, priority=normal]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=600, urgent=false]`
- **Association:** Person[P004] ⋈ Order[O001]

##### Token combiné 12
- **`p`**: Person[P001] - `Person[id=P001, name=Alice, age=25, priority=high]`
- **`o`**: Order[O002] - `Order[id=O002, customer_id=P002, amount=300, urgent=false]`
- **Association:** Person[P001] ⋈ Order[O002]

##### Token combiné 13
- **`p`**: Person[P004] - `Person[priority=normal, id=P004, name=Diana, age=19]`
- **`o`**: Order[O003] - `Order[customer_id=P003, amount=100, urgent=true, id=O003]`
- **Association:** Person[P004] ⋈ Order[O003]

##### Token combiné 14
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=35, priority=low]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=P003, amount=100, urgent=true]`
- **Association:** Person[P003] ⋈ Order[O003]

##### Token combiné 15
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=35, priority=low]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P004, amount=50, urgent=false]`
- **Association:** Person[P003] ⋈ Order[O004]

##### Token combiné 16
- **`p`**: Person[P002] - `Person[name=Bob, age=67, priority=normal, id=P002]`
- **`o`**: Order[O002] - `Order[urgent=false, id=O002, customer_id=P002, amount=300]`
- **Association:** Person[P002] ⋈ Order[O002]

#### 🎯 Activation détaillée: `priority_order`
- **Nombre de déclenchements:** 16
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`p`**: Person[P001] - `Person[name=Alice, age=25, priority=high, id=P001]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=P003, amount=100, urgent=true]`
- **Association:** Person[P001] ⋈ Order[O003]

##### Token combiné 2
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=35, priority=low]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P004, amount=50, urgent=false]`
- **Association:** Person[P003] ⋈ Order[O004]

##### Token combiné 3
- **`p`**: Person[P001] - `Person[priority=high, id=P001, name=Alice, age=25]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=600, urgent=false]`
- **Association:** Person[P001] ⋈ Order[O001]

##### Token combiné 4
- **`p`**: Person[P001] - `Person[name=Alice, age=25, priority=high, id=P001]`
- **`o`**: Order[O002] - `Order[urgent=false, id=O002, customer_id=P002, amount=300]`
- **Association:** Person[P001] ⋈ Order[O002]

##### Token combiné 5
- **`p`**: Person[P004] - `Person[age=19, priority=normal, id=P004, name=Diana]`
- **`o`**: Order[O003] - `Order[id=O003, customer_id=P003, amount=100, urgent=true]`
- **Association:** Person[P004] ⋈ Order[O003]

##### Token combiné 6
- **`p`**: Person[P002] - `Person[name=Bob, age=67, priority=normal, id=P002]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P004, amount=50, urgent=false]`
- **Association:** Person[P002] ⋈ Order[O004]

##### Token combiné 7
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=35, priority=low]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=600, urgent=false]`
- **Association:** Person[P003] ⋈ Order[O001]

##### Token combiné 8
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=35, priority=low]`
- **`o`**: Order[O002] - `Order[amount=300, urgent=false, id=O002, customer_id=P002]`
- **Association:** Person[P003] ⋈ Order[O002]

##### Token combiné 9
- **`p`**: Person[P002] - `Person[id=P002, name=Bob, age=67, priority=normal]`
- **`o`**: Order[O003] - `Order[amount=100, urgent=true, id=O003, customer_id=P003]`
- **Association:** Person[P002] ⋈ Order[O003]

##### Token combiné 10
- **`p`**: Person[P004] - `Person[name=Diana, age=19, priority=normal, id=P004]`
- **`o`**: Order[O004] - `Order[amount=50, urgent=false, id=O004, customer_id=P004]`
- **Association:** Person[P004] ⋈ Order[O004]

##### Token combiné 11
- **`p`**: Person[P002] - `Person[name=Bob, age=67, priority=normal, id=P002]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=600, urgent=false]`
- **Association:** Person[P002] ⋈ Order[O001]

##### Token combiné 12
- **`p`**: Person[P002] - `Person[priority=normal, id=P002, name=Bob, age=67]`
- **`o`**: Order[O002] - `Order[id=O002, customer_id=P002, amount=300, urgent=false]`
- **Association:** Person[P002] ⋈ Order[O002]

##### Token combiné 13
- **`p`**: Person[P003] - `Person[id=P003, name=Charlie, age=35, priority=low]`
- **`o`**: Order[O003] - `Order[amount=100, urgent=true, id=O003, customer_id=P003]`
- **Association:** Person[P003] ⋈ Order[O003]

##### Token combiné 14
- **`p`**: Person[P001] - `Person[age=25, priority=high, id=P001, name=Alice]`
- **`o`**: Order[O004] - `Order[id=O004, customer_id=P004, amount=50, urgent=false]`
- **Association:** Person[P001] ⋈ Order[O004]

##### Token combiné 15
- **`p`**: Person[P004] - `Person[name=Diana, age=19, priority=normal, id=P004]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=600, urgent=false]`
- **Association:** Person[P004] ⋈ Order[O001]

##### Token combiné 16
- **`p`**: Person[P004] - `Person[id=P004, name=Diana, age=19, priority=normal]`
- **`o`**: Order[O002] - `Order[id=O002, customer_id=P002, amount=300, urgent=false]`
- **Association:** Person[P004] ⋈ Order[O002]

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | p <-> o | 16 | inner | ✅ |
| join_1 | p <-> o | 16 | inner | ✅ |
| join_2 | p <-> o | 16 | inner | ✅ |

---

## 🧪 TEST 10: join_simple
---

### 📋 Informations générales
- **Description:** Test jointure simple entre deux faits
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_simple.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_simple.facts`
- **Temps d'exécution:** 1.216373ms
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 20.0%
- **Actions valides:** ✅
- **Jointures valides:** ✅
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

### 📜 Règles analysées
#### Règle 1
- **Texte original:** `{p: Person, o: Order} / p.id == o.customer_id ==> customer_order_match(p.id, o.id)`
- **Action:** customer_order_match
- **Type de nœud:** JoinNode
- **Type sémantique:** equality
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
├── 🔗 BetaNodes (Jointures)
│   ├── rule_0_join
│   │   ├── Variables: p ⋈ o
│   │   ├── Conditions: 1
│   │   └── Type: JoinNode
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
Order[id=O001, customer_id=P001, amount=100]
Order[id=O002, customer_id=P002, amount=200]

```

**Total faits:** 4

- **Person:** 2 faits
- **Order:** 2 faits

**📋 Détail des faits parsés:**
1. **Person[P001]** - `Person[id=P001, name=Alice, age=25]`
2. **Person[P002]** - `Person[id=P002, name=Bob, age=30]`
3. **Order[O001]** - `Order[id=O001, customer_id=P001, amount=100]`
4. **Order[O002]** - `Order[amount=200, id=O002, customer_id=P002]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| customer_order_match | 2 | AlphaNode | ✅ |

#### 🎯 Activation détaillée: `customer_order_match`
- **Nombre de déclenchements:** 2
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`p`**: Person[P001] - `Person[id=P001, name=Alice, age=25]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=100]`
- **Association:** Person[P001] ⋈ Order[O001]

##### Token combiné 2
- **`p`**: Person[P002] - `Person[name=Bob, age=30, id=P002]`
- **`o`**: Order[O002] - `Order[id=O002, customer_id=P002, amount=200]`
- **Association:** Person[P002] ⋈ Order[O002]

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | p <-> o | 2 | inner | ✅ |

### 🎯 Comparaison attendu vs observé
#### Actions
| Action | Attendu | Observé | Statut |
|--------|---------|---------|--------|
| customer_order_match | 2-2 | 2 | ✅ |

#### 📋 TOKENS COMBINÉS ATTENDUS vs OBTENUS

**🎯 Action `customer_order_match`:**
- **Description:** Deux customers avec leurs commandes matchent
- **Variables de la règle:** p, o

**📍 TOKENS COMBINÉS ATTENDUS:**
- **Nombre de tokens attendus:** 2-2
- **Token attendu 1:**
  * `p`: Person[P001] - `Person[name=Alice, age=25, id=P001]`
  * `o`: Order[O001] - `Order[id=O001, customer_id=P001, amount=100]`
- **Token attendu 2:**
  * `p`: Person[P002] - `Person[id=P002, name=Bob, age=30]`
  * `o`: Order[O002] - `Order[id=O002, customer_id=P002, amount=200]`

**📊 TOKENS COMBINÉS OBTENUS:**
- **Nombre de tokens obtenus:** 2
- **Token obtenu 1:**
  * `p`: Person[P001] - `Person[id=P001, name=Alice, age=25]`
  * `o`: Order[O001] - `Order[amount=100, id=O001, customer_id=P001]`
- **Token obtenu 2:**
  * `p`: Person[P002] - `Person[id=P002, name=Bob, age=30]`
  * `o`: Order[O002] - `Order[id=O002, customer_id=P002, amount=200]`

**🎯 RÉSULTAT:** ✅ SUCCÈS
- ✅ Nombre de tokens correct

---

## 🧪 TEST 11: not_complex_operator
---

### 📋 Informations générales
- **Description:** Test opérateur NOT complexe dans nœuds beta
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/not_complex_operator.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/not_complex_operator.facts`
- **Temps d'exécution:** 19.376907ms
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
- **Texte original:** `{e: Employee, p: Project} / e.id == p.lead_id AND e.active == true AND NOT (p.status == "cancelled") ==> active_project_lead(e.id, p.id)`
- **Action:** active_project_lead
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - e (Employee): primary
  - p (Project): secondary

#### Règle 2
- **Texte original:** `{e: Employee, p: Project} / e.id == p.lead_id AND e.active == true AND NOT (p.status == "cancelled") ==> active_project_lead(e.id, p.id)`
- **Action:** qualified_project_lead
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
- **Complexité:** complex
- **Variables:**
  - e (Employee): primary
  - p (Project): secondary

#### Règle 3
- **Texte original:** `{e: Employee, p: Project} / e.id == p.lead_id AND e.active == true AND NOT (p.status == "cancelled") ==> active_project_lead(e.id, p.id)`
- **Action:** permanent_project_lead
- **Type de nœud:** JoinNode
- **Type sémantique:** logical
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
├── 🔗 BetaNodes (Jointures)
│   ├── rule_0_join
│   │   ├── Variables: e ⋈ p
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_1_join
│   │   ├── Variables: e ⋈ p
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│   ├── rule_2_join
│   │   ├── Variables: e ⋈ p
│   │   ├── Conditions: 0
│   │   └── Type: JoinNode
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: active_project_lead
    ├── rule_1_terminal
    │   └── Action: qualified_project_lead
    ├── rule_2_terminal
    │   └── Action: permanent_project_lead
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Employee[id=E001, name=Alice, department=engineering, active=true]
Employee[id=E002, name=Bob, department=temp, active=true]
Employee[id=E003, name=Charlie, department=marketing, active=false]
Employee[id=E004, name=Diana, department=sales, active=true]
Project[id=P001, lead_id=E001, status=active, budget=5000]
Project[id=P002, lead_id=E002, status=cancelled, budget=2000]
Project[id=P003, lead_id=E003, status=completed, budget=800]
Project[id=P004, lead_id=E004, status=active, budget=1500]

```

**Total faits:** 8

- **Employee:** 4 faits
- **Project:** 4 faits

**📋 Détail des faits parsés:**
1. **Employee[E001]** - `Employee[id=E001, name=Alice, department=engineering, active=true]`
2. **Employee[E002]** - `Employee[name=Bob, department=temp, active=true, id=E002]`
3. **Employee[E003]** - `Employee[id=E003, name=Charlie, department=marketing, active=false]`
4. **Employee[E004]** - `Employee[id=E004, name=Diana, department=sales, active=true]`
5. **Project[P001]** - `Project[id=P001, lead_id=E001, status=active, budget=5000]`
6. **Project[P002]** - `Project[status=cancelled, budget=2000, id=P002, lead_id=E002]`
7. **Project[P003]** - `Project[id=P003, lead_id=E003, status=completed, budget=800]`
8. **Project[P004]** - `Project[budget=1500, id=P004, lead_id=E004, status=active]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| active_project_lead | 16 | AlphaNode | ❌ |
| qualified_project_lead | 16 | AlphaNode | ❌ |
| permanent_project_lead | 16 | AlphaNode | ❌ |

#### 🎯 Activation détaillée: `active_project_lead`
- **Nombre de déclenchements:** 16
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`e`**: Employee[E004] - `Employee[id=E004, name=Diana, department=sales, active=true]`
- **`p`**: Project[P004] - `Project[status=active, budget=1500, id=P004, lead_id=E004]`
- **Association:** Employee[E004] ⋈ Project[P004]

##### Token combiné 2
- **`e`**: Employee[E001] - `Employee[id=E001, name=Alice, department=engineering, active=true]`
- **`p`**: Project[P001] - `Project[id=P001, lead_id=E001, status=active, budget=5000]`
- **Association:** Employee[E001] ⋈ Project[P001]

##### Token combiné 3
- **`e`**: Employee[E002] - `Employee[department=temp, active=true, id=E002, name=Bob]`
- **`p`**: Project[P001] - `Project[status=active, budget=5000, id=P001, lead_id=E001]`
- **Association:** Employee[E002] ⋈ Project[P001]

##### Token combiné 4
- **`e`**: Employee[E003] - `Employee[active=false, id=E003, name=Charlie, department=marketing]`
- **`p`**: Project[P001] - `Project[id=P001, lead_id=E001, status=active, budget=5000]`
- **Association:** Employee[E003] ⋈ Project[P001]

##### Token combiné 5
- **`e`**: Employee[E004] - `Employee[id=E004, name=Diana, department=sales, active=true]`
- **`p`**: Project[P002] - `Project[budget=2000, id=P002, lead_id=E002, status=cancelled]`
- **Association:** Employee[E004] ⋈ Project[P002]

##### Token combiné 6
- **`e`**: Employee[E003] - `Employee[id=E003, name=Charlie, department=marketing, active=false]`
- **`p`**: Project[P003] - `Project[id=P003, lead_id=E003, status=completed, budget=800]`
- **Association:** Employee[E003] ⋈ Project[P003]

##### Token combiné 7
- **`e`**: Employee[E001] - `Employee[id=E001, name=Alice, department=engineering, active=true]`
- **`p`**: Project[P003] - `Project[id=P003, lead_id=E003, status=completed, budget=800]`
- **Association:** Employee[E001] ⋈ Project[P003]

##### Token combiné 8
- **`e`**: Employee[E002] - `Employee[id=E002, name=Bob, department=temp, active=true]`
- **`p`**: Project[P003] - `Project[status=completed, budget=800, id=P003, lead_id=E003]`
- **Association:** Employee[E002] ⋈ Project[P003]

##### Token combiné 9
- **`e`**: Employee[E001] - `Employee[name=Alice, department=engineering, active=true, id=E001]`
- **`p`**: Project[P004] - `Project[status=active, budget=1500, id=P004, lead_id=E004]`
- **Association:** Employee[E001] ⋈ Project[P004]

##### Token combiné 10
- **`e`**: Employee[E002] - `Employee[active=true, id=E002, name=Bob, department=temp]`
- **`p`**: Project[P004] - `Project[budget=1500, id=P004, lead_id=E004, status=active]`
- **Association:** Employee[E002] ⋈ Project[P004]

##### Token combiné 11
- **`e`**: Employee[E003] - `Employee[id=E003, name=Charlie, department=marketing, active=false]`
- **`p`**: Project[P004] - `Project[id=P004, lead_id=E004, status=active, budget=1500]`
- **Association:** Employee[E003] ⋈ Project[P004]

##### Token combiné 12
- **`e`**: Employee[E004] - `Employee[id=E004, name=Diana, department=sales, active=true]`
- **`p`**: Project[P001] - `Project[id=P001, lead_id=E001, status=active, budget=5000]`
- **Association:** Employee[E004] ⋈ Project[P001]

##### Token combiné 13
- **`e`**: Employee[E001] - `Employee[id=E001, name=Alice, department=engineering, active=true]`
- **`p`**: Project[P002] - `Project[id=P002, lead_id=E002, status=cancelled, budget=2000]`
- **Association:** Employee[E001] ⋈ Project[P002]

##### Token combiné 14
- **`e`**: Employee[E002] - `Employee[id=E002, name=Bob, department=temp, active=true]`
- **`p`**: Project[P002] - `Project[status=cancelled, budget=2000, id=P002, lead_id=E002]`
- **Association:** Employee[E002] ⋈ Project[P002]

##### Token combiné 15
- **`e`**: Employee[E003] - `Employee[id=E003, name=Charlie, department=marketing, active=false]`
- **`p`**: Project[P002] - `Project[id=P002, lead_id=E002, status=cancelled, budget=2000]`
- **Association:** Employee[E003] ⋈ Project[P002]

##### Token combiné 16
- **`e`**: Employee[E004] - `Employee[id=E004, name=Diana, department=sales, active=true]`
- **`p`**: Project[P003] - `Project[id=P003, lead_id=E003, status=completed, budget=800]`
- **Association:** Employee[E004] ⋈ Project[P003]

#### 🎯 Activation détaillée: `qualified_project_lead`
- **Nombre de déclenchements:** 16
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`e`**: Employee[E003] - `Employee[id=E003, name=Charlie, department=marketing, active=false]`
- **`p`**: Project[P001] - `Project[id=P001, lead_id=E001, status=active, budget=5000]`
- **Association:** Employee[E003] ⋈ Project[P001]

##### Token combiné 2
- **`e`**: Employee[E002] - `Employee[active=true, id=E002, name=Bob, department=temp]`
- **`p`**: Project[P002] - `Project[status=cancelled, budget=2000, id=P002, lead_id=E002]`
- **Association:** Employee[E002] ⋈ Project[P002]

##### Token combiné 3
- **`e`**: Employee[E001] - `Employee[department=engineering, active=true, id=E001, name=Alice]`
- **`p`**: Project[P002] - `Project[id=P002, lead_id=E002, status=cancelled, budget=2000]`
- **Association:** Employee[E001] ⋈ Project[P002]

##### Token combiné 4
- **`e`**: Employee[E003] - `Employee[id=E003, name=Charlie, department=marketing, active=false]`
- **`p`**: Project[P003] - `Project[id=P003, lead_id=E003, status=completed, budget=800]`
- **Association:** Employee[E003] ⋈ Project[P003]

##### Token combiné 5
- **`e`**: Employee[E004] - `Employee[id=E004, name=Diana, department=sales, active=true]`
- **`p`**: Project[P004] - `Project[id=P004, lead_id=E004, status=active, budget=1500]`
- **Association:** Employee[E004] ⋈ Project[P004]

##### Token combiné 6
- **`e`**: Employee[E004] - `Employee[id=E004, name=Diana, department=sales, active=true]`
- **`p`**: Project[P002] - `Project[lead_id=E002, status=cancelled, budget=2000, id=P002]`
- **Association:** Employee[E004] ⋈ Project[P002]

##### Token combiné 7
- **`e`**: Employee[E002] - `Employee[id=E002, name=Bob, department=temp, active=true]`
- **`p`**: Project[P004] - `Project[id=P004, lead_id=E004, status=active, budget=1500]`
- **Association:** Employee[E002] ⋈ Project[P004]

##### Token combiné 8
- **`e`**: Employee[E001] - `Employee[department=engineering, active=true, id=E001, name=Alice]`
- **`p`**: Project[P004] - `Project[id=P004, lead_id=E004, status=active, budget=1500]`
- **Association:** Employee[E001] ⋈ Project[P004]

##### Token combiné 9
- **`e`**: Employee[E004] - `Employee[id=E004, name=Diana, department=sales, active=true]`
- **`p`**: Project[P001] - `Project[id=P001, lead_id=E001, status=active, budget=5000]`
- **Association:** Employee[E004] ⋈ Project[P001]

##### Token combiné 10
- **`e`**: Employee[E004] - `Employee[id=E004, name=Diana, department=sales, active=true]`
- **`p`**: Project[P003] - `Project[status=completed, budget=800, id=P003, lead_id=E003]`
- **Association:** Employee[E004] ⋈ Project[P003]

##### Token combiné 11
- **`e`**: Employee[E003] - `Employee[id=E003, name=Charlie, department=marketing, active=false]`
- **`p`**: Project[P004] - `Project[id=P004, lead_id=E004, status=active, budget=1500]`
- **Association:** Employee[E003] ⋈ Project[P004]

##### Token combiné 12
- **`e`**: Employee[E001] - `Employee[active=true, id=E001, name=Alice, department=engineering]`
- **`p`**: Project[P001] - `Project[id=P001, lead_id=E001, status=active, budget=5000]`
- **Association:** Employee[E001] ⋈ Project[P001]

##### Token combiné 13
- **`e`**: Employee[E002] - `Employee[id=E002, name=Bob, department=temp, active=true]`
- **`p`**: Project[P001] - `Project[id=P001, lead_id=E001, status=active, budget=5000]`
- **Association:** Employee[E002] ⋈ Project[P001]

##### Token combiné 14
- **`e`**: Employee[E003] - `Employee[id=E003, name=Charlie, department=marketing, active=false]`
- **`p`**: Project[P002] - `Project[id=P002, lead_id=E002, status=cancelled, budget=2000]`
- **Association:** Employee[E003] ⋈ Project[P002]

##### Token combiné 15
- **`e`**: Employee[E001] - `Employee[id=E001, name=Alice, department=engineering, active=true]`
- **`p`**: Project[P003] - `Project[lead_id=E003, status=completed, budget=800, id=P003]`
- **Association:** Employee[E001] ⋈ Project[P003]

##### Token combiné 16
- **`e`**: Employee[E002] - `Employee[name=Bob, department=temp, active=true, id=E002]`
- **`p`**: Project[P003] - `Project[id=P003, lead_id=E003, status=completed, budget=800]`
- **Association:** Employee[E002] ⋈ Project[P003]

#### 🎯 Activation détaillée: `permanent_project_lead`
- **Nombre de déclenchements:** 16
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`e`**: Employee[E002] - `Employee[active=true, id=E002, name=Bob, department=temp]`
- **`p`**: Project[P002] - `Project[status=cancelled, budget=2000, id=P002, lead_id=E002]`
- **Association:** Employee[E002] ⋈ Project[P002]

##### Token combiné 2
- **`e`**: Employee[E001] - `Employee[id=E001, name=Alice, department=engineering, active=true]`
- **`p`**: Project[P004] - `Project[status=active, budget=1500, id=P004, lead_id=E004]`
- **Association:** Employee[E001] ⋈ Project[P004]

##### Token combiné 3
- **`e`**: Employee[E004] - `Employee[name=Diana, department=sales, active=true, id=E004]`
- **`p`**: Project[P004] - `Project[id=P004, lead_id=E004, status=active, budget=1500]`
- **Association:** Employee[E004] ⋈ Project[P004]

##### Token combiné 4
- **`e`**: Employee[E001] - `Employee[id=E001, name=Alice, department=engineering, active=true]`
- **`p`**: Project[P001] - `Project[budget=5000, id=P001, lead_id=E001, status=active]`
- **Association:** Employee[E001] ⋈ Project[P001]

##### Token combiné 5
- **`e`**: Employee[E004] - `Employee[id=E004, name=Diana, department=sales, active=true]`
- **`p`**: Project[P001] - `Project[id=P001, lead_id=E001, status=active, budget=5000]`
- **Association:** Employee[E004] ⋈ Project[P001]

##### Token combiné 6
- **`e`**: Employee[E003] - `Employee[id=E003, name=Charlie, department=marketing, active=false]`
- **`p`**: Project[P002] - `Project[id=P002, lead_id=E002, status=cancelled, budget=2000]`
- **Association:** Employee[E003] ⋈ Project[P002]

##### Token combiné 7
- **`e`**: Employee[E001] - `Employee[id=E001, name=Alice, department=engineering, active=true]`
- **`p`**: Project[P003] - `Project[lead_id=E003, status=completed, budget=800, id=P003]`
- **Association:** Employee[E001] ⋈ Project[P003]

##### Token combiné 8
- **`e`**: Employee[E004] - `Employee[name=Diana, department=sales, active=true, id=E004]`
- **`p`**: Project[P003] - `Project[id=P003, lead_id=E003, status=completed, budget=800]`
- **Association:** Employee[E004] ⋈ Project[P003]

##### Token combiné 9
- **`e`**: Employee[E003] - `Employee[id=E003, name=Charlie, department=marketing, active=false]`
- **`p`**: Project[P001] - `Project[status=active, budget=5000, id=P001, lead_id=E001]`
- **Association:** Employee[E003] ⋈ Project[P001]

##### Token combiné 10
- **`e`**: Employee[E004] - `Employee[department=sales, active=true, id=E004, name=Diana]`
- **`p`**: Project[P002] - `Project[id=P002, lead_id=E002, status=cancelled, budget=2000]`
- **Association:** Employee[E004] ⋈ Project[P002]

##### Token combiné 11
- **`e`**: Employee[E001] - `Employee[id=E001, name=Alice, department=engineering, active=true]`
- **`p`**: Project[P002] - `Project[status=cancelled, budget=2000, id=P002, lead_id=E002]`
- **Association:** Employee[E001] ⋈ Project[P002]

##### Token combiné 12
- **`e`**: Employee[E003] - `Employee[id=E003, name=Charlie, department=marketing, active=false]`
- **`p`**: Project[P003] - `Project[id=P003, lead_id=E003, status=completed, budget=800]`
- **Association:** Employee[E003] ⋈ Project[P003]

##### Token combiné 13
- **`e`**: Employee[E002] - `Employee[id=E002, name=Bob, department=temp, active=true]`
- **`p`**: Project[P004] - `Project[id=P004, lead_id=E004, status=active, budget=1500]`
- **Association:** Employee[E002] ⋈ Project[P004]

##### Token combiné 14
- **`e`**: Employee[E002] - `Employee[id=E002, name=Bob, department=temp, active=true]`
- **`p`**: Project[P001] - `Project[lead_id=E001, status=active, budget=5000, id=P001]`
- **Association:** Employee[E002] ⋈ Project[P001]

##### Token combiné 15
- **`e`**: Employee[E002] - `Employee[id=E002, name=Bob, department=temp, active=true]`
- **`p`**: Project[P003] - `Project[id=P003, lead_id=E003, status=completed, budget=800]`
- **Association:** Employee[E002] ⋈ Project[P003]

##### Token combiné 16
- **`e`**: Employee[E003] - `Employee[name=Charlie, department=marketing, active=false, id=E003]`
- **`p`**: Project[P004] - `Project[lead_id=E004, status=active, budget=1500, id=P004]`
- **Association:** Employee[E003] ⋈ Project[P004]

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | e <-> p | 16 | inner | ✅ |
| join_1 | e <-> p | 16 | inner | ✅ |
| join_2 | e <-> p | 16 | inner | ✅ |

---

## 🧪 TEST 12: not_simple
---

### 📋 Informations générales
- **Description:** Test négation simple
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/not_simple.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/not_simple.facts`
- **Temps d'exécution:** 482.873µs
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 20.0%
- **Actions valides:** ✅
- **Jointures valides:** ✅
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

### 📜 Règles analysées
#### Règle 1
- **Texte original:** `{p: Person} / NOT (p.active == false) ==> active_person(p.id)`
- **Action:** active_person
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
    │   └── Action: active_person
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Person[id=P001, name=Alice, active=true]
Person[id=P002, name=Bob, active=false]

```

**Total faits:** 2

- **Person:** 2 faits

**📋 Détail des faits parsés:**
1. **Person[P001]** - `Person[id=P001, name=Alice, active=true]`
2. **Person[P002]** - `Person[id=P002, name=Bob, active=false]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| active_person | 1 | AlphaNode | ✅ |

#### 🎯 Activation détaillée: `active_person`
- **Nombre de déclenchements:** 1
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`p`**: Person[P001] - `Person[id=P001, name=Alice, active=true]`

### 🚫 Analyse des négations (NotNodes)
| Nœud | Condition Niée | Faits Filtrés | Type | Validation |
|------|----------------|---------------|------|------------|
| alpha_rule_0_alpha | map[condition:map[left:map[field:active object:p type:fieldAccess] operator:== right:map[type:boolean value:false] type:comparison] negated:true type:negation] | 0 | simple | ❌ |

### 🎯 Comparaison attendu vs observé
#### Actions
| Action | Attendu | Observé | Statut |
|--------|---------|---------|--------|
| active_person | 1-1 | 1 | ✅ |

#### 📋 TOKENS COMBINÉS ATTENDUS vs OBTENUS

**🎯 Action `active_person`:**
- **Description:** Une personne active (Alice) passe le filtre NOT
- **Variables de la règle:** p

**📍 TOKENS COMBINÉS ATTENDUS:**
- **Nombre de tokens attendus:** 1-1
- **Token attendu 1:**
  * `p`: Person[P001] - `Person[name=Alice, active=true, id=P001]`

**📊 TOKENS COMBINÉS OBTENUS:**
- **Nombre de tokens obtenus:** 1
- **Token obtenu 1:**
  * `p`: Person[P001] - `Person[id=P001, name=Alice, active=true]`

**🎯 RÉSULTAT:** ✅ SUCCÈS
- ✅ Nombre de tokens correct

---

## 💡 RECOMMANDATIONS
### Amélioration de la couverture Beta
### Prochaines étapes
1. **Ajouter plus de tests complexes** avec jointures multiples
2. **Tester les négations imbriquées** et conditions complexes
3. **Valider les performances** des nœuds Beta avec de gros volumes
4. **Enrichir la validation sémantique** avec plus de critères
