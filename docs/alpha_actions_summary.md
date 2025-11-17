# 📊 RAPPORT ALPHA NODES - ACTIONS FILTRÉES

**Date de génération:** 2025-11-17 16:03:53  
**Méthode:** Filtrage par action spécifique de chaque test

---

## 📋 Résumé Exécutif

| Test | Action Testée | Attendu | Obtenu | Statut |
|------|---------------|---------|---------|---------|
| `alpha_abs_negative` | `small_balance_found` | 2 | 2 | ✅ |
| `alpha_abs_positive` | `significant_balance_found` | 2 | 2 | ✅ |
| `alpha_boolean_negative` | `inactive_account_found` | 1 | 1 | ✅ |
| `alpha_boolean_positive` | `active_account_found` | 2 | 2 | ✅ |
| `alpha_comparison_negative` | `affordable_product` | 1 | 1 | ✅ |
| `alpha_comparison_positive` | `expensive_product` | 2 | 2 | ✅ |
| `alpha_contains_negative` | `normal_message_found` | 2 | 2 | ✅ |
| `alpha_contains_positive` | `urgent_message_found` | 2 | 2 | ✅ |
| `alpha_equal_sign_negative` | `non_gold_customer_found` | 2 | 2 | ✅ |
| `alpha_equal_sign_positive` | `gold_customer_found` | 2 | 2 | ✅ |
| `alpha_equality_negative` | `age_is_not_twenty_five` | 1 | 1 | ✅ |
| `alpha_equality_positive` | `age_is_twenty_five` | 2 | 2 | ✅ |
| `alpha_in_negative` | `invalid_state_found` | 2 | 2 | ✅ |
| `alpha_in_positive` | `valid_state_found` | 2 | 2 | ✅ |
| `alpha_inequality_negative` | `cancelled_order_found` | 1 | 1 | ✅ |
| `alpha_inequality_positive` | `valid_order_found` | 2 | 2 | ✅ |
| `alpha_length_negative` | `weak_password_found` | 2 | 2 | ✅ |
| `alpha_length_positive` | `secure_password_found` | 2 | 2 | ✅ |
| `alpha_like_negative` | `external_email_found` | 2 | 2 | ✅ |
| `alpha_like_positive` | `company_email_found` | 2 | 2 | ✅ |
| `alpha_matches_negative` | `invalid_code_found` | 2 | 2 | ✅ |
| `alpha_matches_positive` | `valid_code_found` | 2 | 2 | ✅ |
| `alpha_string_negative` | `non_admin_user_found` | 1 | 1 | ✅ |
| `alpha_string_positive` | `admin_user_found` | 2 | 2 | ✅ |
| `alpha_upper_negative` | `non_finance_dept_found` | 2 | 2 | ✅ |
| `alpha_upper_positive` | `finance_dept_found` | 2 | 2 | ✅ |

## 📈 Statistiques Globales

- **Tests exécutés:** 26
- **Tests conformes:** 26
- **Taux de conformité:** 100.0%

---

## 🔬 Détails par Test

### 🧪 Test: `alpha_abs_negative`

**Action testée:** `small_balance_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test fonction ABS() (valeur absolue) - négation
type Balance : <id: string, amount: number, type: string>

{b: Balance} / NOT(ABS(b.amount) > 100) ==> small_balance_found(b.id, b.amount)
```

#### 📊 Faits (.facts)
```facts
Balance[id="B001", amount=150.0, type="credit"]
Balance[id="B002", amount=-25.0, type="debit"]
Balance[id="B003", amount=75.0, type="credit"]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ small_balance_found (Balance[amount=-25, type=debit, id=B002])
- ✅ small_balance_found (Balance[id=B003, amount=75, type=credit])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `small_balance_found`:** 2

---


### 🧪 Test: `alpha_abs_positive`

**Action testée:** `significant_balance_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test fonction ABS() (valeur absolue)
type Balance : <id: string, amount: number, type: string>

{b: Balance} / ABS(b.amount) > 100 ==> significant_balance_found(b.id, b.amount)
```

#### 📊 Faits (.facts)
```facts
Balance[id="B001", amount=150.0, type="credit"]
Balance[id="B002", amount=-200.0, type="debit"]
Balance[id="B003", amount=50.0, type="credit"]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ significant_balance_found (Balance[id=B001, amount=150, type=credit])
- ✅ significant_balance_found (Balance[id=B002, amount=-200, type=debit])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `significant_balance_found`:** 2

---


### 🧪 Test: `alpha_boolean_negative`

**Action testée:** `inactive_account_found`  
**Résultat:** ✅ Conforme (1/1 actions)

#### 📋 Règle (.constraint)
```constraint
// Test condition booléenne négative
type Account : <id: string, balance: number, active: bool>

{a: Account} / NOT(a.active == true) ==> inactive_account_found(a.id, a.balance)
```

#### 📊 Faits (.facts)
```facts
Account[id=ACC001, balance=1000, active=true]
Account[id=ACC002, balance=500, active=false]
Account[id=ACC003, balance=2000, active=true]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])

**Toutes les actions du test:** 2 total
**Actions filtrées pour `inactive_account_found`:** 1

---


### 🧪 Test: `alpha_boolean_positive`

**Action testée:** `active_account_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test condition booléenne positive
type Account : <id: string, balance: number, active: bool>

{a: Account} / a.active == true ==> active_account_found(a.id, a.balance)
```

#### 📊 Faits (.facts)
```facts
Account[id=ACC001, balance=1000, active=true]
Account[id=ACC002, balance=500, active=false] 
Account[id=ACC003, balance=2000, active=true]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `active_account_found`:** 2

---


### 🧪 Test: `alpha_comparison_negative`

**Action testée:** `affordable_product`  
**Résultat:** ✅ Conforme (1/1 actions)

#### 📋 Règle (.constraint)
```constraint
// Test comparaison numérique négative
type Product : <id: string, price: number, category: string>

{prod: Product} / NOT(prod.price > 100) ==> affordable_product(prod.id, prod.price)
```

#### 📊 Faits (.facts)
```facts
Product[id=PROD001, price=150, category=electronics]
Product[id=PROD002, price=50, category=books]
Product[id=PROD003, price=200, category=electronics]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ affordable_product (Product[id=PROD002, price=50, category=books])

**Toutes les actions du test:** 2 total
**Actions filtrées pour `affordable_product`:** 1

---


### 🧪 Test: `alpha_comparison_positive`

**Action testée:** `expensive_product`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test comparaison numérique positive
type Product : <id: string, price: number, category: string>

{prod: Product} / prod.price > 100 ==> expensive_product(prod.id, prod.price)
```

#### 📊 Faits (.facts)
```facts
Product[id=PROD001, price=150, category=electronics]
Product[id=PROD002, price=50, category=books]
Product[id=PROD003, price=200, category=electronics]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[price=200, category=electronics, id=PROD003])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `expensive_product`:** 2

---


### 🧪 Test: `alpha_contains_negative`

**Action testée:** `normal_message_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test opérateur CONTAINS (contenance) - négation
type Message : <id: string, content: string, urgent: bool>

{m: Message} / NOT(m.content CONTAINS "urgent") ==> normal_message_found(m.id, m.content)
```

#### 📊 Faits (.facts)
```facts
Message[id="M001", content="This is urgent please respond", urgent=true]
Message[id="M002", content="Regular message content", urgent=false]
Message[id="M003", content="Simple notification", urgent=false]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `normal_message_found`:** 2

---


### 🧪 Test: `alpha_contains_positive`

**Action testée:** `urgent_message_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test opérateur CONTAINS (contenance)
type Message : <id: string, content: string, urgent: bool>

{m: Message} / m.content CONTAINS "urgent" ==> urgent_message_found(m.id, m.content)
```

#### 📊 Faits (.facts)
```facts
Message[id="M001", content="This is urgent please respond", urgent=true]
Message[id="M002", content="Regular message content", urgent=false]
Message[id="M003", content="Very urgent matter!", urgent=true]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `urgent_message_found`:** 2

---


### 🧪 Test: `alpha_equal_sign_negative`

**Action testée:** `non_gold_customer_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test opérateur = (égalité alternative) - négation
type Customer : <id: string, tier: string, points: number>

{c: Customer} / NOT(c.tier = "gold") ==> non_gold_customer_found(c.id, c.tier)
```

#### 📊 Faits (.facts)
```facts
Customer[id="C001", tier="gold", points=5000]
Customer[id="C002", tier="silver", points=2000]
Customer[id="C003", tier="bronze", points=1000]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ non_gold_customer_found (Customer[id=C002, tier=silver, points=2000])
- ✅ non_gold_customer_found (Customer[id=C003, tier=bronze, points=1000])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `non_gold_customer_found`:** 2

---


### 🧪 Test: `alpha_equal_sign_positive`

**Action testée:** `gold_customer_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test opérateur = (égalité alternative)
type Customer : <id: string, tier: string, points: number>

{c: Customer} / c.tier = "gold" ==> gold_customer_found(c.id, c.points)
```

#### 📊 Faits (.facts)
```facts
Customer[id="C001", tier="gold", points=5000]
Customer[id="C002", tier="silver", points=2000]
Customer[id="C003", tier="gold", points=7500]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ gold_customer_found (Customer[points=5000, id=C001, tier=gold])
- ✅ gold_customer_found (Customer[id=C003, tier=gold, points=7500])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `gold_customer_found`:** 2

---


### 🧪 Test: `alpha_equality_negative`

**Action testée:** `age_is_not_twenty_five`  
**Résultat:** ✅ Conforme (1/1 actions)

#### 📋 Règle (.constraint)
```constraint
// Test égalité négative simple
type Person : <id: string, age: number, status: string>

{p: Person} / NOT(p.age == 25) ==> age_is_not_twenty_five(p.id, p.age)
```

#### 📊 Faits (.facts)
```facts
Person[id=P001, age=25, status=active]
Person[id=P002, age=30, status=active]  
Person[id=P003, age=25, status=inactive]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ age_is_not_twenty_five (Person[status=active, id=P002, age=30])

**Toutes les actions du test:** 2 total
**Actions filtrées pour `age_is_not_twenty_five`:** 1

---


### 🧪 Test: `alpha_equality_positive`

**Action testée:** `age_is_twenty_five`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test égalité positive simple
type Person : <id: string, age: number, status: string>

{p: Person} / p.age == 25 ==> age_is_twenty_five(p.id, p.age)
```

#### 📊 Faits (.facts)
```facts
Person[id=P001, age=25, status=active]
Person[id=P002, age=30, status=active]
Person[id=P003, age=25, status=inactive]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `age_is_twenty_five`:** 2

---


### 🧪 Test: `alpha_in_negative`

**Action testée:** `invalid_state_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test opérateur IN (appartenance) - négation
type Status : <id: string, state: string, priority: number>

{s: Status} / NOT(s.state IN ["active", "pending"]) ==> invalid_state_found(s.id, s.state)
```

#### 📊 Faits (.facts)
```facts
Status[id="S001", state="active", priority=1]
Status[id="S002", state="inactive", priority=3]
Status[id="S003", state="archived", priority=5]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ invalid_state_found (Status[id=S002, state=inactive, priority=3])
- ✅ invalid_state_found (Status[state=archived, priority=5, id=S003])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `invalid_state_found`:** 2

---


### 🧪 Test: `alpha_in_positive`

**Action testée:** `valid_state_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test opérateur IN (appartenance)
type Status : <id: string, state: string, priority: number>

{s: Status} / s.state IN ["active", "pending", "review"] ==> valid_state_found(s.id, s.state)
```

#### 📊 Faits (.facts)
```facts
Status[id="S001", state="active", priority=1]
Status[id="S002", state="inactive", priority=3]
Status[id="S003", state="pending", priority=2]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ valid_state_found (Status[state=active, priority=1, id=S001])
- ✅ valid_state_found (Status[id=S003, state=pending, priority=2])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `valid_state_found`:** 2

---


### 🧪 Test: `alpha_inequality_negative`

**Action testée:** `cancelled_order_found`  
**Résultat:** ✅ Conforme (1/1 actions)

#### 📋 Règle (.constraint)
```constraint
// Test inégalité négative
type Order : <id: string, total: number, status: string>

{o: Order} / NOT(o.status != "cancelled") ==> cancelled_order_found(o.id, o.total)
```

#### 📊 Faits (.facts)
```facts
Order[id=ORD001, total=100, status=pending]
Order[id=ORD002, total=200, status=cancelled]
Order[id=ORD003, total=300, status=completed]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])

**Toutes les actions du test:** 2 total
**Actions filtrées pour `cancelled_order_found`:** 1

---


### 🧪 Test: `alpha_inequality_positive`

**Action testée:** `valid_order_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test inégalité positive
type Order : <id: string, total: number, status: string>

{o: Order} / o.status != "cancelled" ==> valid_order_found(o.id, o.total)
```

#### 📊 Faits (.facts)
```facts
Order[id=ORD001, total=100, status=pending]
Order[id=ORD002, total=200, status=cancelled]
Order[id=ORD003, total=300, status=completed]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ valid_order_found (Order[status=pending, id=ORD001, total=100])
- ✅ valid_order_found (Order[total=300, status=completed, id=ORD003])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `valid_order_found`:** 2

---


### 🧪 Test: `alpha_length_negative`

**Action testée:** `weak_password_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test fonction LENGTH() - négation
type Password : <id: string, value: string, secure: bool>

{p: Password} / NOT(LENGTH(p.value) >= 8) ==> weak_password_found(p.id, p.value)
```

#### 📊 Faits (.facts)
```facts
Password[id="P001", value="password123", secure=true]
Password[id="P002", value="123", secure=false]
Password[id="P003", value="pass", secure=false]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ weak_password_found (Password[id=P002, value=123, secure=false])
- ✅ weak_password_found (Password[id=P003, value=pass, secure=false])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `weak_password_found`:** 2

---


### 🧪 Test: `alpha_length_positive`

**Action testée:** `secure_password_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test fonction LENGTH()
type Password : <id: string, value: string, secure: bool>

{p: Password} / LENGTH(p.value) >= 8 ==> secure_password_found(p.id, p.value)
```

#### 📊 Faits (.facts)
```facts
Password[id="P001", value="password123", secure=true]
Password[id="P002", value="123", secure=false]
Password[id="P003", value="verysecurepass", secure=true]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `secure_password_found`:** 2

---


### 🧪 Test: `alpha_like_negative`

**Action testée:** `external_email_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test opérateur LIKE (motif) - négation
type Email : <id: string, address: string, verified: bool>

{e: Email} / NOT(e.address LIKE "%@company.com") ==> external_email_found(e.id, e.address)
```

#### 📊 Faits (.facts)
```facts
Email[id="E001", address="john@company.com", verified=true]
Email[id="E002", address="jane@external.org", verified=false]
Email[id="E003", address="user@other.net", verified=true]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ external_email_found (Email[verified=false, id=E002, address=jane@external.org])
- ✅ external_email_found (Email[id=E003, address=user@other.net, verified=true])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `external_email_found`:** 2

---


### 🧪 Test: `alpha_like_positive`

**Action testée:** `company_email_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test opérateur LIKE (motif)
type Email : <id: string, address: string, verified: bool>

{e: Email} / e.address LIKE "%@company.com" ==> company_email_found(e.id, e.address)
```

#### 📊 Faits (.facts)
```facts
Email[id="E001", address="john@company.com", verified=true]
Email[id="E002", address="jane@external.org", verified=false]
Email[id="E003", address="admin@company.com", verified=true]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ company_email_found (Email[id=E001, address=john@company.com, verified=true])
- ✅ company_email_found (Email[id=E003, address=admin@company.com, verified=true])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `company_email_found`:** 2

---


### 🧪 Test: `alpha_matches_negative`

**Action testée:** `invalid_code_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test opérateur MATCHES (regex) - négation
type Code : <id: string, value: string, active: bool>

{c: Code} / NOT(c.value MATCHES "CODE[0-9]+") ==> invalid_code_found(c.id, c.value)
```

#### 📊 Faits (.facts)
```facts
Code[id="C001", value="CODE123", active=true]
Code[id="C002", value="INVALID", active=false]
Code[id="C003", value="BADFORMAT", active=true]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ invalid_code_found (Code[id=C002, value=INVALID, active=false])
- ✅ invalid_code_found (Code[id=C003, value=BADFORMAT, active=true])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `invalid_code_found`:** 2

---


### 🧪 Test: `alpha_matches_positive`

**Action testée:** `valid_code_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test opérateur MATCHES (regex)
type Code : <id: string, value: string, active: bool>

{c: Code} / c.value MATCHES "CODE[0-9]+" ==> valid_code_found(c.id, c.value)
```

#### 📊 Faits (.facts)
```facts
Code[id="C001", value="CODE123", active=true]
Code[id="C002", value="INVALID", active=false]
Code[id="C003", value="CODE999", active=true]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ valid_code_found (Code[id=C001, value=CODE123, active=true])
- ✅ valid_code_found (Code[id=C003, value=CODE999, active=true])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `valid_code_found`:** 2

---


### 🧪 Test: `alpha_string_negative`

**Action testée:** `non_admin_user_found`  
**Résultat:** ✅ Conforme (1/1 actions)

#### 📋 Règle (.constraint)
```constraint
// Test condition string négative
type User : <id: string, name: string, role: string>

{u: User} / NOT(u.role == "admin") ==> non_admin_user_found(u.id, u.name)
```

#### 📊 Faits (.facts)
```facts
User[id=U001, name=Alice, role=admin]
User[id=U002, name=Bob, role=user]
User[id=U003, name=Charlie, role=admin]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])

**Toutes les actions du test:** 2 total
**Actions filtrées pour `non_admin_user_found`:** 1

---


### 🧪 Test: `alpha_string_positive`

**Action testée:** `admin_user_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test condition string positive
type User : <id: string, name: string, role: string>

{u: User} / u.role == "admin" ==> admin_user_found(u.id, u.name)
```

#### 📊 Faits (.facts)
```facts
User[id=U001, name=Alice, role=admin]
User[id=U002, name=Bob, role=user]
User[id=U003, name=Charlie, role=admin]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ admin_user_found (User[name=Alice, role=admin, id=U001])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `admin_user_found`:** 2

---


### 🧪 Test: `alpha_upper_negative`

**Action testée:** `non_finance_dept_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test fonction UPPER() (majuscules) - négation
type Department : <id: string, name: string, active: bool>

{d: Department} / NOT(UPPER(d.name) == "FINANCE") ==> non_finance_dept_found(d.id, d.name)
```

#### 📊 Faits (.facts)
```facts
Department[id="D001", name="finance", active=true]
Department[id="D002", name="IT", active=true]
Department[id="D003", name="HR", active=true]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ non_finance_dept_found (Department[name=IT, active=true, id=D002])
- ✅ non_finance_dept_found (Department[id=D003, name=HR, active=true])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `non_finance_dept_found`:** 2

---


### 🧪 Test: `alpha_upper_positive`

**Action testée:** `finance_dept_found`  
**Résultat:** ✅ Conforme (2/2 actions)

#### 📋 Règle (.constraint)
```constraint
// Test fonction UPPER() (majuscules)
type Department : <id: string, name: string, active: bool>

{d: Department} / UPPER(d.name) == "FINANCE" ==> finance_dept_found(d.id, d.name)
```

#### 📊 Faits (.facts)
```facts
Department[id="D001", name="finance", active=true]
Department[id="D002", name="IT", active=true]
Department[id="D003", name="Finance", active=true]
```

#### 🎯 Actions Déclenchées (Filtrées)

- ✅ finance_dept_found (Department[id=D001, name=finance, active=true])
- ✅ finance_dept_found (Department[id=D003, name=Finance, active=true])

**Toutes les actions du test:** 4 total
**Actions filtrées pour `finance_dept_found`:** 2

---


## 📝 Notes

- **Filtrage:** Seules les actions correspondant exactement à l'action définie dans chaque règle sont affichées
- **Conformité:** Un test est conforme si le nombre d'actions obtenues correspond au nombre attendu
- **Actions totales:** Chaque test peut déclencher d'autres actions du réseau global, mais seules les actions spécifiques sont comptabilisées
