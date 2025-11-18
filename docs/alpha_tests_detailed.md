# 📊 RAPPORT STRUCTURÉ - TESTS ALPHA COMPLETS (ACTIONS FILTRÉES)

**Date de génération:** 2025-11-17 16:04:05
**Format:** Test par test avec actions filtrées par règle spécifique

---

## 🔬 Tests Alpha Originaux

### 🧪 Test: `alpha_boolean_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_boolean_negative
**Description:** Test Boolean Negative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test condition booléenne négative
type Account : <id: string, balance: number, active: bool>

{a: Account} / NOT(a.active == true) ==> inactive_account_found(a.id, a.balance)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Account[id=ACC001, balance=1000, active=true]
Account[id=ACC002, balance=500, active=false]
Account[id=ACC003, balance=2000, active=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `inactive_account_found`
**Logique:** NOT(a.active == true) → Comptes avec active=false

**Faits devant déclencher l'action:**
- ACC002 (active=false)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 1 actions déclenchées correspondent aux 1 attendues

---

### 🧪 Test: `alpha_boolean_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_boolean_positive
**Description:** Test Boolean Positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test condition booléenne positive
type Account : <id: string, balance: number, active: bool>

{a: Account} / a.active == true ==> active_account_found(a.id, a.balance)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Account[id=ACC001, balance=1000, active=true]
Account[id=ACC002, balance=500, active=false]
Account[id=ACC003, balance=2000, active=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `active_account_found`
**Logique:** a.active == true → Comptes avec active=true

**Faits devant déclencher l'action:**
- ACC001 (active=true)
- ACC003 (active=true)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_comparison_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_comparison_negative
**Description:** Test Comparison Negative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test comparaison numérique négative
type Product : <id: string, price: number, category: string>

{prod: Product} / NOT(prod.price > 100) ==> affordable_product(prod.id, prod.price)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Product[id=PROD001, price=150, category=electronics]
Product[id=PROD002, price=50, category=books]
Product[id=PROD003, price=200, category=electronics]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `affordable_product`
**Logique:** NOT(prod.price > 100) → Produits avec price ≤ 100

**Faits devant déclencher l'action:**
- PROD002 (price=50)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 1 actions déclenchées correspondent aux 1 attendues

---

### 🧪 Test: `alpha_comparison_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_comparison_positive
**Description:** Test Comparison Positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test comparaison numérique positive
type Product : <id: string, price: number, category: string>

{prod: Product} / prod.price > 100 ==> expensive_product(prod.id, prod.price)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Product[id=PROD001, price=150, category=electronics]
Product[id=PROD002, price=50, category=books]
Product[id=PROD003, price=200, category=electronics]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `expensive_product`
**Logique:** prod.price > 100 → Produits avec price > 100

**Faits devant déclencher l'action:**
- PROD001 (price=150)
- PROD003 (price=200)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_equality_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_equality_negative
**Description:** Test Equality Negative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test égalité négative simple
type Person : <id: string, age: number, status: string>

{p: Person} / NOT(p.age == 25) ==> age_is_not_twenty_five(p.id, p.age)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Person[id=P001, age=25, status=active]
Person[id=P002, age=30, status=active]
Person[id=P003, age=25, status=inactive]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `age_is_not_twenty_five`
**Logique:** NOT(p.age == 25) → Personnes avec age ≠ 25

**Faits devant déclencher l'action:**
- P002 (age=30)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 1 actions déclenchées correspondent aux 1 attendues

---

### 🧪 Test: `alpha_equality_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_equality_positive
**Description:** Test Equality Positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test égalité positive simple
type Person : <id: string, age: number, status: string>

{p: Person} / p.age == 25 ==> age_is_twenty_five(p.id, p.age)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Person[id=P001, age=25, status=active]
Person[id=P002, age=30, status=active]
Person[id=P003, age=25, status=inactive]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `age_is_twenty_five`
**Logique:** p.age == 25 → Personnes avec age = 25

**Faits devant déclencher l'action:**
- P001 (age=25)
- P003 (age=25)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_inequality_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_inequality_negative
**Description:** Test Inequality Negative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test inégalité négative
type Order : <id: string, total: number, status: string>

{o: Order} / NOT(o.status != "cancelled") ==> cancelled_order_found(o.id, o.total)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Order[id=ORD001, total=100, status=pending]
Order[id=ORD002, total=200, status=cancelled]
Order[id=ORD003, total=300, status=completed]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `cancelled_order_found`
**Logique:** NOT(o.status != 'cancelled') → Commandes avec status = cancelled

**Faits devant déclencher l'action:**
- ORD002 (status=cancelled)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 1 actions déclenchées correspondent aux 1 attendues

---

### 🧪 Test: `alpha_inequality_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_inequality_positive
**Description:** Test Inequality Positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test inégalité positive
type Order : <id: string, total: number, status: string>

{o: Order} / o.status != "cancelled" ==> valid_order_found(o.id, o.total)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Order[id=ORD001, total=100, status=pending]
Order[id=ORD002, total=200, status=cancelled]
Order[id=ORD003, total=300, status=completed]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `valid_order_found`
**Logique:** o.status != 'cancelled' → Commandes avec status ≠ cancelled

**Faits devant déclencher l'action:**
- ORD001 (status=pending)
- ORD003 (status=completed)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_string_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_string_negative
**Description:** Test String Negative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test condition string négative
type User : <id: string, name: string, role: string>

{u: User} / NOT(u.role == "admin") ==> non_admin_user_found(u.id, u.name)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
User[id=U001, name=Alice, role=admin]
User[id=U002, name=Bob, role=user]
User[id=U003, name=Charlie, role=admin]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `non_admin_user_found`
**Logique:** NOT(u.role == 'admin') → Utilisateurs avec role ≠ admin

**Faits devant déclencher l'action:**
- U002 (role=user)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ non_admin_user_found (User[name=Bob, role=user, id=U002])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 1 actions déclenchées correspondent aux 1 attendues

---

### 🧪 Test: `alpha_string_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_string_positive
**Description:** Test String Positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test condition string positive
type User : <id: string, name: string, role: string>

{u: User} / u.role == "admin" ==> admin_user_found(u.id, u.name)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
User[id=U001, name=Alice, role=admin]
User[id=U002, name=Bob, role=user]
User[id=U003, name=Charlie, role=admin]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `admin_user_found`
**Logique:** u.role == 'admin' → Utilisateurs avec role = admin

**Faits devant déclencher l'action:**
- U001 (role=admin)
- U003 (role=admin)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ admin_user_found (User[name=Alice, role=admin, id=U001])
- ✅ admin_user_found (User[role=admin, id=U003, name=Charlie])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---


## 🧪 Tests Alpha Étendus

### 🧪 Test: `alpha_abs_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_abs_negative
**Description:** Test Abs Negative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test fonction ABS() (valeur absolue) - négation
type Balance : <id: string, amount: number, type: string>

{b: Balance} / NOT(ABS(b.amount) > 100) ==> small_balance_found(b.id, b.amount)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Balance[id="B001", amount=150.0, type="credit"]
Balance[id="B002", amount=-25.0, type="debit"]
Balance[id="B003", amount=75.0, type="credit"]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `small_balance_found`
**Logique:** NOT(ABS(b.amount) > 100) → Montants absolus ≤ 100

**Faits devant déclencher l'action:**
- B002 (|-25|=25)
- B003 (|75|=75)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ small_balance_found (Balance[id=B002, amount=-25, type=debit])
- ✅ small_balance_found (Balance[amount=75, type=credit, id=B003])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_abs_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_abs_positive
**Description:** Test Abs Positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test fonction ABS() (valeur absolue)
type Balance : <id: string, amount: number, type: string>

{b: Balance} / ABS(b.amount) > 100 ==> significant_balance_found(b.id, b.amount)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Balance[id="B001", amount=150.0, type="credit"]
Balance[id="B002", amount=-200.0, type="debit"]
Balance[id="B003", amount=50.0, type="credit"]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `significant_balance_found`
**Logique:** ABS(b.amount) > 100 → Montants absolus > 100

**Faits devant déclencher l'action:**
- B001 (|150|=150)
- B002 (|-200|=200)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ significant_balance_found (Balance[amount=150, type=credit, id=B001])
- ✅ significant_balance_found (Balance[id=B002, amount=-200, type=debit])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_contains_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_contains_negative
**Description:** Test Contains Negative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test opérateur CONTAINS (contenance) - négation
type Message : <id: string, content: string, urgent: bool>

{m: Message} / NOT(m.content CONTAINS "urgent") ==> normal_message_found(m.id, m.content)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Message[id="M001", content="This is urgent please respond", urgent=true]
Message[id="M002", content="Regular message content", urgent=false]
Message[id="M003", content="Simple notification", urgent=false]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `normal_message_found`
**Logique:** NOT(m.content CONTAINS 'urgent') → Messages ne contenant pas 'urgent'

**Faits devant déclencher l'action:**
- M002 (sans 'urgent')
- M003 (sans 'urgent')

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_contains_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_contains_positive
**Description:** Test Contains Positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test opérateur CONTAINS (contenance)
type Message : <id: string, content: string, urgent: bool>

{m: Message} / m.content CONTAINS "urgent" ==> urgent_message_found(m.id, m.content)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Message[id="M001", content="This is urgent please respond", urgent=true]
Message[id="M002", content="Regular message content", urgent=false]
Message[id="M003", content="Very urgent matter!", urgent=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `urgent_message_found`
**Logique:** m.content CONTAINS 'urgent' → Messages contenant 'urgent'

**Faits devant déclencher l'action:**
- M001 (contient 'urgent')
- M003 (contient 'urgent')

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ urgent_message_found (Message[urgent=true, id=M001, content=This is urgent please respond])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_equal_sign_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_equal_sign_negative
**Description:** Test Equal Sign Negative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test opérateur = (égalité alternative) - négation
type Customer : <id: string, tier: string, points: number>

{c: Customer} / NOT(c.tier = "gold") ==> non_gold_customer_found(c.id, c.tier)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Customer[id="C001", tier="gold", points=5000]
Customer[id="C002", tier="silver", points=2000]
Customer[id="C003", tier="bronze", points=1000]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `non_gold_customer_found`
**Logique:** NOT(c.tier = 'gold') → Clients non gold

**Faits devant déclencher l'action:**
- C002 (tier=silver)
- C003 (tier=bronze)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ non_gold_customer_found (Customer[id=C002, tier=silver, points=2000])
- ✅ non_gold_customer_found (Customer[id=C003, tier=bronze, points=1000])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_equal_sign_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_equal_sign_positive
**Description:** Test Equal Sign Positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test opérateur = (égalité alternative)
type Customer : <id: string, tier: string, points: number>

{c: Customer} / c.tier = "gold" ==> gold_customer_found(c.id, c.points)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Customer[id="C001", tier="gold", points=5000]
Customer[id="C002", tier="silver", points=2000]
Customer[id="C003", tier="gold", points=7500]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `gold_customer_found`
**Logique:** c.tier = 'gold' → Clients gold

**Faits devant déclencher l'action:**
- C001 (tier=gold)
- C003 (tier=gold)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ gold_customer_found (Customer[tier=gold, points=5000, id=C001])
- ✅ gold_customer_found (Customer[id=C003, tier=gold, points=7500])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_in_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_in_negative
**Description:** Test In Negative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test opérateur IN (appartenance) - négation
type Status : <id: string, state: string, priority: number>

{s: Status} / NOT(s.state IN ["active", "pending"]) ==> invalid_state_found(s.id, s.state)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Status[id="S001", state="active", priority=1]
Status[id="S002", state="inactive", priority=3]
Status[id="S003", state="archived", priority=5]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `invalid_state_found`
**Logique:** NOT(s.state IN ['active', 'pending']) → États inactifs

**Faits devant déclencher l'action:**
- S002 (state='inactive')
- S003 (state='archived')

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ invalid_state_found (Status[id=S002, state=inactive, priority=3])
- ✅ invalid_state_found (Status[state=archived, priority=5, id=S003])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_in_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_in_positive
**Description:** Test In Positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test opérateur IN (appartenance)
type Status : <id: string, state: string, priority: number>

{s: Status} / s.state IN ["active", "pending", "review"] ==> valid_state_found(s.id, s.state)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Status[id="S001", state="active", priority=1]
Status[id="S002", state="inactive", priority=3]
Status[id="S003", state="pending", priority=2]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `valid_state_found`
**Logique:** s.state IN ['active', 'pending'] → États actifs

**Faits devant déclencher l'action:**
- S001 (state='active')
- S003 (state='pending')

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ valid_state_found (Status[id=S001, state=active, priority=1])
- ✅ valid_state_found (Status[id=S003, state=pending, priority=2])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_length_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_length_negative
**Description:** Test Length Negative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test fonction LENGTH() - négation
type Password : <id: string, value: string, secure: bool>

{p: Password} / NOT(LENGTH(p.value) >= 8) ==> weak_password_found(p.id, p.value)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Password[id="P001", value="password123", secure=true]
Password[id="P002", value="123", secure=false]
Password[id="P003", value="pass", secure=false]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `weak_password_found`
**Logique:** NOT(LENGTH(p.value) >= 8) → Mots de passe courts

**Faits devant déclencher l'action:**
- P002 (length=3)
- P003 (length=4)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ weak_password_found (Password[id=P002, value=123, secure=false])
- ✅ weak_password_found (Password[id=P003, value=pass, secure=false])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_length_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_length_positive
**Description:** Test Length Positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test fonction LENGTH()
type Password : <id: string, value: string, secure: bool>

{p: Password} / LENGTH(p.value) >= 8 ==> secure_password_found(p.id, p.value)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Password[id="P001", value="password123", secure=true]
Password[id="P002", value="123", secure=false]
Password[id="P003", value="verysecurepass", secure=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `secure_password_found`
**Logique:** LENGTH(p.value) >= 8 → Mots de passe longs

**Faits devant déclencher l'action:**
- P001 (length=11)
- P003 (length=14)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_like_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_like_negative
**Description:** Test Like Negative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test opérateur LIKE (motif) - négation
type Email : <id: string, address: string, verified: bool>

{e: Email} / NOT(e.address LIKE "%@company.com") ==> external_email_found(e.id, e.address)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Email[id="E001", address="john@company.com", verified=true]
Email[id="E002", address="jane@external.org", verified=false]
Email[id="E003", address="user@other.net", verified=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `external_email_found`
**Logique:** NOT(e.address LIKE '%@company.com') → Emails externes

**Faits devant déclencher l'action:**
- E002 (jane@external.org)
- E003 (user@other.net)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ external_email_found (Email[id=E002, address=jane@external.org, verified=false])
- ✅ external_email_found (Email[address=user@other.net, verified=true, id=E003])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_like_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_like_positive
**Description:** Test Like Positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test opérateur LIKE (motif)
type Email : <id: string, address: string, verified: bool>

{e: Email} / e.address LIKE "%@company.com" ==> company_email_found(e.id, e.address)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Email[id="E001", address="john@company.com", verified=true]
Email[id="E002", address="jane@external.org", verified=false]
Email[id="E003", address="admin@company.com", verified=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `company_email_found`
**Logique:** e.address LIKE '%@company.com' → Emails company

**Faits devant déclencher l'action:**
- E001 (john@company.com)
- E003 (admin@company.com)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ company_email_found (Email[verified=true, id=E001, address=john@company.com])
- ✅ company_email_found (Email[id=E003, address=admin@company.com, verified=true])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_matches_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_matches_negative
**Description:** Test Matches Negative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test opérateur MATCHES (regex) - négation
type Code : <id: string, value: string, active: bool>

{c: Code} / NOT(c.value MATCHES "CODE[0-9]+") ==> invalid_code_found(c.id, c.value)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Code[id="C001", value="CODE123", active=true]
Code[id="C002", value="INVALID", active=false]
Code[id="C003", value="BADFORMAT", active=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `invalid_code_found`
**Logique:** NOT(c.value MATCHES 'CODE[0-9]+') → Codes invalides

**Faits devant déclencher l'action:**
- C002 (INVALID)
- C003 (BADFORMAT)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ invalid_code_found (Code[active=false, id=C002, value=INVALID])
- ✅ invalid_code_found (Code[id=C003, value=BADFORMAT, active=true])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_matches_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_matches_positive
**Description:** Test Matches Positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test opérateur MATCHES (regex)
type Code : <id: string, value: string, active: bool>

{c: Code} / c.value MATCHES "CODE[0-9]+" ==> valid_code_found(c.id, c.value)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Code[id="C001", value="CODE123", active=true]
Code[id="C002", value="INVALID", active=false]
Code[id="C003", value="CODE999", active=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `valid_code_found`
**Logique:** c.value MATCHES 'CODE[0-9]+' → Codes valides

**Faits devant déclencher l'action:**
- C001 (CODE123)
- C003 (CODE999)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ valid_code_found (Code[id=C001, value=CODE123, active=true])
- ✅ valid_code_found (Code[active=true, id=C003, value=CODE999])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_upper_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_upper_negative
**Description:** Test Upper Negative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test fonction UPPER() (majuscules) - négation
type Department : <id: string, name: string, active: bool>

{d: Department} / NOT(UPPER(d.name) == "FINANCE") ==> non_finance_dept_found(d.id, d.name)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Department[id="D001", name="finance", active=true]
Department[id="D002", name="IT", active=true]
Department[id="D003", name="HR", active=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `non_finance_dept_found`
**Logique:** NOT(UPPER(d.name) == 'FINANCE') → Départements non finance

**Faits devant déclencher l'action:**
- D002 (UPPER('IT'))
- D003 (UPPER('HR'))

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ non_finance_dept_found (Department[id=D002, name=IT, active=true])
- ✅ non_finance_dept_found (Department[id=D003, name=HR, active=true])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---

### 🧪 Test: `alpha_upper_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_upper_positive
**Description:** Test Upper Positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
// Test fonction UPPER() (majuscules)
type Department : <id: string, name: string, active: bool>

{d: Department} / UPPER(d.name) == "FINANCE" ==> finance_dept_found(d.id, d.name)
```

#### 3️⃣ Faits Soumis (.facts)

```facts
Department[id="D001", name="finance", active=true]
Department[id="D002", name="IT", active=true]
Department[id="D003", name="Finance", active=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `finance_dept_found`
**Logique:** UPPER(d.name) == 'FINANCE' → Départements finance

**Faits devant déclencher l'action:**
- D001 (UPPER('finance'))
- D003 (UPPER('Finance'))

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Conforme

**Actions déclenchées (filtrées):**
- ✅ finance_dept_found (Department[id=D001, name=finance, active=true])
- ✅ finance_dept_found (Department[id=D003, name=Finance, active=true])

#### 6️⃣ Analyse du Test

**Résultat:** ✅ **Conforme:** 2 actions déclenchées correspondent aux 2 attendues

---


## 📈 Bilan Final

**Tests exécutés:** 26
**Tests conformes:** 26
**Taux de conformité:** 100.0%

**Note:** Les actions affichées sont filtrées pour ne montrer que celles correspondant exactement à l'action définie dans la règle de chaque test.
