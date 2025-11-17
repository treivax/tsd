# 📊 RAPPORT STRUCTURÉ - TESTS ALPHA COMPLETS

**Date de génération:** 2025-11-17 15:17:23
**Format:** Test par test avec structure complète

---

## 🔬 Tests Alpha Originaux

### 🧪 Test: `alpha_boolean_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_boolean_negative
**Description:** Test condition booléenne négative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
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
**Logique:** NOT(acc.active == true) → Comptes avec active=false

**Faits devant déclencher l'action:**
- ACC002 (active=false)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Succès

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[active=false, id=ACC002, balance=500])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[active=true, id=ACC003, balance=2000])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[category=electronics, id=PROD001, price=150])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[status=active, id=P002, age=30])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[status=completed, id=ORD003, total=300])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[name=Bob, role=user, id=U002])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[role=admin, id=U003, name=Charlie])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[name=Charlie, role=admin, id=U003])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 30 actions obtenues vs 1 attendues

⚠️ **Écart:** 30 actions obtenues vs 1 attendues

---

### 🧪 Test: `alpha_boolean_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_boolean_positive
**Description:** Test condition booléenne positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
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
**Logique:** acc.active == true → Comptes avec active=true

**Faits devant déclencher l'action:**
- ACC001 (active=true)
- ACC003 (active=true)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Succès

**Actions déclenchées:**
- ✅ inactive_account_found (Account[active=false, id=ACC002, balance=500])
- ✅ inactive_account_found (Account[balance=500, active=false, id=ACC002])
- ✅ active_account_found (Account[balance=1000, active=true, id=ACC001])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[balance=2000, active=true, id=ACC003])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[category=electronics, id=PROD003, price=200])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[status=inactive, id=P003, age=25])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[total=200, status=cancelled, id=ORD002])
- ✅ cancelled_order_found (Order[status=cancelled, id=ORD002, total=200])
- ✅ valid_order_found (Order[status=pending, id=ORD001, total=100])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[status=pending, id=ORD001, total=100])
- ✅ valid_order_found (Order[status=completed, id=ORD003, total=300])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[name=Bob, role=user, id=U002])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[name=Charlie, role=admin, id=U003])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 30 actions obtenues vs 2 attendues

⚠️ **Écart:** 30 actions obtenues vs 2 attendues

---

### 🧪 Test: `alpha_comparison_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_comparison_negative
**Description:** Test comparaison numérique négative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
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

**Statut:** ✅ Succès

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[price=50, category=books, id=PROD002])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[age=30, status=active, id=P002])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 30 actions obtenues vs 1 attendues

⚠️ **Écart:** 30 actions obtenues vs 1 attendues

---

### 🧪 Test: `alpha_comparison_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_comparison_positive
**Description:** Test comparaison numérique positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
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

**Statut:** ✅ Succès

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[active=false, id=ACC002, balance=500])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[active=true, id=ACC001, balance=1000])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[price=150, category=electronics, id=PROD001])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[status=inactive, id=P003, age=25])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[name=Bob, role=user, id=U002])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 30 actions obtenues vs 2 attendues

⚠️ **Écart:** 30 actions obtenues vs 2 attendues

---

### 🧪 Test: `alpha_equality_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_equality_negative
**Description:** Test égalité négative simple

#### 2️⃣ Règles Complètes (.constraint)

```constraint
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

**Statut:** ✅ Succès

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[balance=1000, active=true, id=ACC001])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[balance=2000, active=true, id=ACC003])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[status=inactive, id=P003, age=25])
- ✅ cancelled_order_found (Order[total=200, status=cancelled, id=ORD002])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[total=300, status=completed, id=ORD003])
- ✅ non_admin_user_found (User[role=user, id=U002, name=Bob])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 30 actions obtenues vs 1 attendues

⚠️ **Écart:** 30 actions obtenues vs 1 attendues

---

### 🧪 Test: `alpha_equality_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_equality_positive
**Description:** Test égalité positive simple

#### 2️⃣ Règles Complètes (.constraint)

```constraint
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

**Statut:** ✅ Succès

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[active=true, id=ACC001, balance=1000])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[price=50, category=books, id=PROD002])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[price=150, category=electronics, id=PROD001])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[age=30, status=active, id=P002])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[role=user, id=U002, name=Bob])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 30 actions obtenues vs 2 attendues

⚠️ **Écart:** 30 actions obtenues vs 2 attendues

---

### 🧪 Test: `alpha_inequality_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_inequality_negative
**Description:** Test inégalité négative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
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
**Logique:** NOT(ord.status != "cancelled") → Commandes avec status = cancelled

**Faits devant déclencher l'action:**
- ORD002 (status=cancelled)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Succès

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[balance=2000, active=true, id=ACC003])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[status=active, id=P002, age=30])
- ✅ age_is_twenty_five (Person[age=25, status=active, id=P001])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[status=active, id=P001, age=25])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[status=pending, id=ORD001, total=100])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[status=completed, id=ORD003, total=300])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[name=Alice, role=admin, id=U001])
- ✅ admin_user_found (User[role=admin, id=U003, name=Charlie])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 30 actions obtenues vs 1 attendues

⚠️ **Écart:** 30 actions obtenues vs 1 attendues

---

### 🧪 Test: `alpha_inequality_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_inequality_positive
**Description:** Test inégalité positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
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
**Logique:** ord.status != "cancelled" → Commandes avec status ≠ cancelled

**Faits devant déclencher l'action:**
- ORD001 (status=pending)
- ORD003 (status=completed)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Succès

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[price=50, category=books, id=PROD002])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[category=electronics, id=PROD003, price=200])
- ✅ age_is_not_twenty_five (Person[age=30, status=active, id=P002])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[age=25, status=active, id=P001])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[status=pending, id=ORD001, total=100])
- ✅ valid_order_found (Order[total=300, status=completed, id=ORD003])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[name=Charlie, role=admin, id=U003])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 30 actions obtenues vs 2 attendues

⚠️ **Écart:** 30 actions obtenues vs 2 attendues

---

### 🧪 Test: `alpha_string_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_string_negative
**Description:** Test condition string négative

#### 2️⃣ Règles Complètes (.constraint)

```constraint
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
**Logique:** NOT(u.role == "admin") → Utilisateurs avec role ≠ admin

**Faits devant déclencher l'action:**
- U002 (role=user)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Succès

**Actions déclenchées:**
- ✅ inactive_account_found (Account[balance=500, active=false, id=ACC002])
- ✅ inactive_account_found (Account[balance=500, active=false, id=ACC002])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[active=true, id=ACC003, balance=2000])
- ✅ active_account_found (Account[active=true, id=ACC001, balance=1000])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[total=300, status=completed, id=ORD003])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 30 actions obtenues vs 1 attendues

⚠️ **Écart:** 30 actions obtenues vs 1 attendues

---

### 🧪 Test: `alpha_string_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_string_positive
**Description:** Test condition string positive

#### 2️⃣ Règles Complètes (.constraint)

```constraint
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
**Logique:** u.role == "admin" → Utilisateurs avec role = admin

**Faits devant déclencher l'action:**
- U001 (role=admin)
- U003 (role=admin)

#### 5️⃣ Résultat Obtenu

**Statut:** ✅ Succès

**Actions déclenchées:**
- ✅ inactive_account_found (Account[balance=500, active=false, id=ACC002])
- ✅ inactive_account_found (Account[balance=500, active=false, id=ACC002])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[balance=2000, active=true, id=ACC003])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[price=200, category=electronics, id=PROD003])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[age=25, status=active, id=P001])
- ✅ age_is_twenty_five (Person[age=25, status=inactive, id=P003])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[total=100, status=pending, id=ORD001])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[name=Charlie, role=admin, id=U003])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[role=admin, id=U003, name=Charlie])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 30 actions obtenues vs 2 attendues

⚠️ **Écart:** 30 actions obtenues vs 2 attendues

---

## 🔬 Tests Alpha Étendus

### 🧪 Test: `alpha_abs_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_abs_negative
**Description:** Test fonction ABS() (valeur absolue) - négation

#### 2️⃣ Règles Complètes (.constraint)

```constraint
type Balance : <id: string, amount: number, type: string>
{b: Balance} / NOT(ABS(b.amount) > 100) ==> small_balance_found(b.id, b.amount)

```

#### 3️⃣ Faits Soumis (.facts)

```facts
Balance[id="B001", amount=150.0, type="credit"]Balance[id="B002", amount=-25.0, type="debit"]Balance[id="B003", amount=75.0, type="credit"]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `small_balance_found`
**Logique:** NOT(ABS(b.amount) > 100) → Soldes absolus ≤ 100

**Faits devant déclencher l'action:**
- B003 (|50| ≤ 100)

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[active=false, id=ACC002, balance=500])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[balance=2000, active=true, id=ACC003])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[price=200, category=electronics, id=PROD003])
- ✅ age_is_not_twenty_five (Person[age=30, status=active, id=P002])
- ✅ age_is_not_twenty_five (Person[age=30, status=active, id=P002])
- ✅ age_is_twenty_five (Person[status=active, id=P001, age=25])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[age=25, status=active, id=P001])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[status=cancelled, id=ORD002, total=200])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[role=admin, id=U003, name=Charlie])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ small_balance_found (Balance[amount=75, type=credit, id=B001])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 1 attendues

⚠️ **Écart:** 44 actions obtenues vs 1 attendues

---

### 🧪 Test: `alpha_abs_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_abs_positive
**Description:** Test fonction ABS() (valeur absolue)

#### 2️⃣ Règles Complètes (.constraint)

```constraint
type Balance : <id: string, amount: number, type: string>
{b: Balance} / ABS(b.amount) > 100 ==> significant_balance_found(b.id, b.amount)

```

#### 3️⃣ Faits Soumis (.facts)

```facts
Balance[id="B001", amount=150.0, type="credit"]Balance[id="B002", amount=-200.0, type="debit"]Balance[id="B003", amount=50.0, type="credit"]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `significant_balance_found`
**Logique:** ABS(b.amount) > 100 → Soldes absolus > 100

**Faits devant déclencher l'action:**
- B001 (|150| > 100)
- B002 (|-200| > 100)

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[category=books, id=PROD002, price=50])
- ✅ affordable_product (Product[price=50, category=books, id=PROD002])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[status=pending, id=ORD001, total=100])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[total=300, status=completed, id=ORD003])
- ✅ non_admin_user_found (User[name=Bob, role=user, id=U002])
- ✅ non_admin_user_found (User[role=user, id=U002, name=Bob])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ small_balance_found (Balance[amount=75, type=credit, id=B001])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[content=Simple notification, urgent=false, id=M003])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ secure_password_found (Password[secure=true, id=P001, value=password123])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 2 attendues

⚠️ **Écart:** 44 actions obtenues vs 2 attendues

---

### 🧪 Test: `alpha_contains_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_contains_negative
**Description:** Test opérateur CONTAINS (contenance) - négation

#### 2️⃣ Règles Complètes (.constraint)

```constraint
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
**Logique:** NOT(m.content CONTAINS "urgent") → Messages sans 'urgent'

**Faits devant déclencher l'action:**
- M002 (content sans 'urgent')
- M003 (content sans 'urgent')

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[category=electronics, id=PROD001, price=150])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[status=active, id=P002, age=30])
- ✅ age_is_not_twenty_five (Person[age=30, status=active, id=P002])
- ✅ age_is_twenty_five (Person[age=25, status=active, id=P001])
- ✅ age_is_twenty_five (Person[age=25, status=inactive, id=P003])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[status=pending, id=ORD001, total=100])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[status=pending, id=ORD001, total=100])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[name=Alice, role=admin, id=U001])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[secure=true, id=P003, value=verysecurepass])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 2 attendues

⚠️ **Écart:** 44 actions obtenues vs 2 attendues

---

### 🧪 Test: `alpha_contains_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_contains_positive
**Description:** Test opérateur CONTAINS (contenance)

#### 2️⃣ Règles Complètes (.constraint)

```constraint
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
**Logique:** m.content CONTAINS "urgent" → Messages contenant 'urgent'

**Faits devant déclencher l'action:**
- M001 (content avec 'urgent')
- M003 (content avec 'urgent')

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[active=false, id=ACC002, balance=500])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[category=electronics, id=PROD003, price=200])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[status=inactive, id=P003, age=25])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[status=completed, id=ORD003, total=300])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[name=Bob, role=user, id=U002])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[name=Alice, role=admin, id=U001])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[name=Alice, role=admin, id=U001])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ normal_message_found (Message[urgent=false, id=M002, content=Regular message content])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[urgent=false, id=M003, content=Simple notification])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[urgent=true, id=M003, content=Very urgent matter!])
- ✅ urgent_message_found (Message[content=This is urgent please respond, urgent=true, id=M001])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 2 attendues

⚠️ **Écart:** 44 actions obtenues vs 2 attendues

---

### 🧪 Test: `alpha_equal_sign_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_equal_sign_negative
**Description:** Test opérateur = (égalité alternative) - négation

#### 2️⃣ Règles Complètes (.constraint)

```constraint
type Customer : <id: string, tier: string, points: number>
{c: Customer} / NOT(c.tier = "gold") ==> non_gold_customer_found(c.id, c.tier)

```

#### 3️⃣ Faits Soumis (.facts)

```facts
Customer[id="C001", tier="gold", points=5000]Customer[id="C002", tier="silver", points=2000]Customer[id="C003", tier="bronze", points=1000]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `non_gold_customer_found`
**Logique:** NOT(cust.tier = "gold") → Clients avec tier ≠ gold

**Faits devant déclencher l'action:**
- C002 (tier=silver)
- C003 (tier=bronze)

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[active=true, id=ACC001, balance=1000])
- ✅ active_account_found (Account[balance=2000, active=true, id=ACC003])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[balance=2000, active=true, id=ACC003])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[category=electronics, id=PROD001, price=150])
- ✅ expensive_product (Product[price=200, category=electronics, id=PROD003])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[status=cancelled, id=ORD002, total=200])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[total=300, status=completed, id=ORD003])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[role=user, id=U002, name=Bob])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[name=Alice, role=admin, id=U001])
- ✅ admin_user_found (User[role=admin, id=U003, name=Charlie])
- ✅ small_balance_found (Balance[type=credit, id=B001, amount=75])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[urgent=false, id=M003, content=Simple notification])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[urgent=false, id=M003, content=Simple notification])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ secure_password_found (Password[secure=true, id=P001, value=password123])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])
- ✅ secure_password_found (Password[value=password123, secure=true, id=P001])
- ✅ secure_password_found (Password[value=verysecurepass, secure=true, id=P003])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 2 attendues

⚠️ **Écart:** 44 actions obtenues vs 2 attendues

---

### 🧪 Test: `alpha_equal_sign_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_equal_sign_positive
**Description:** Test opérateur = (égalité alternative)

#### 2️⃣ Règles Complètes (.constraint)

```constraint
type Customer : <id: string, tier: string, points: number>
{c: Customer} / c.tier = "gold" ==> gold_customer_found(c.id, c.points)

```

#### 3️⃣ Faits Soumis (.facts)

```facts
Customer[id="C001", tier="gold", points=5000]Customer[id="C002", tier="silver", points=2000]Customer[id="C003", tier="gold", points=7500]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `gold_customer_found`
**Logique:** cust.tier = "gold" → Clients avec tier = gold

**Faits devant déclencher l'action:**
- C001 (tier=gold)
- C003 (tier=gold)

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[total=200, status=cancelled, id=ORD002])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ normal_message_found (Message[content=Regular message content, urgent=false, id=M002])
- ✅ normal_message_found (Message[urgent=false, id=M003, content=Simple notification])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ urgent_message_found (Message[content=This is urgent please respond, urgent=true, id=M001])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 2 attendues

⚠️ **Écart:** 44 actions obtenues vs 2 attendues

---

### 🧪 Test: `alpha_in_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_in_negative
**Description:** Test opérateur IN (appartenance) - négation

#### 2️⃣ Règles Complètes (.constraint)

```constraint
type Status : <id: string, state: string, priority: number>
{s: Status} / NOT(s.state IN ["active", "pending"]) ==> invalid_state_found(s.id, s.state)

```

#### 3️⃣ Faits Soumis (.facts)

```facts
Status[id="S001", state="active", priority=1]Status[id="S002", state="inactive", priority=3]Status[id="S003", state="archived", priority=5]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `invalid_state_found`
**Logique:** NOT(s.state IN ["active", "pending", "review"]) → États non valides

**Faits devant déclencher l'action:**
- S002 (state=inactive)

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[balance=2000, active=true, id=ACC003])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[category=books, id=PROD002, price=50])
- ✅ expensive_product (Product[category=electronics, id=PROD001, price=150])
- ✅ expensive_product (Product[category=electronics, id=PROD003, price=200])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[age=30, status=active, id=P002])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[status=inactive, id=P003, age=25])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[status=cancelled, id=ORD002, total=200])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[status=completed, id=ORD003, total=300])
- ✅ valid_order_found (Order[total=100, status=pending, id=ORD001])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[name=Alice, role=admin, id=U001])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ small_balance_found (Balance[amount=75, type=credit, id=B001])
- ✅ small_balance_found (Balance[amount=75, type=credit, id=B001])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[urgent=false, id=M003, content=Simple notification])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[value=verysecurepass, secure=true, id=P003])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 1 attendues

⚠️ **Écart:** 44 actions obtenues vs 1 attendues

---

### 🧪 Test: `alpha_in_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_in_positive
**Description:** Test opérateur IN (appartenance)

#### 2️⃣ Règles Complètes (.constraint)

```constraint
type Status : <id: string, state: string, priority: number>
{s: Status} / s.state IN ["active", "pending", "review"] ==> valid_state_found(s.id, s.state)

```

#### 3️⃣ Faits Soumis (.facts)

```facts
Status[id="S001", state="active", priority=1]Status[id="S002", state="inactive", priority=3]Status[id="S003", state="pending", priority=2]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `valid_state_found`
**Logique:** s.state IN ["active", "pending", "review"] → États valides

**Faits devant déclencher l'action:**
- S001 (state=active)
- S003 (state=pending)

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[balance=1000, active=true, id=ACC001])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[balance=2000, active=true, id=ACC003])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[price=50, category=books, id=PROD002])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[category=electronics, id=PROD003, price=200])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[status=active, id=P002, age=30])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[age=25, status=inactive, id=P003])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[total=200, status=cancelled, id=ORD002])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[status=completed, id=ORD003, total=300])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[role=user, id=U002, name=Bob])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[role=admin, id=U003, name=Charlie])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ normal_message_found (Message[content=Regular message content, urgent=false, id=M002])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ urgent_message_found (Message[content=This is urgent please respond, urgent=true, id=M001])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[secure=true, id=P003, value=verysecurepass])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 2 attendues

⚠️ **Écart:** 44 actions obtenues vs 2 attendues

---

### 🧪 Test: `alpha_length_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_length_negative
**Description:** Test fonction LENGTH() - négation

#### 2️⃣ Règles Complètes (.constraint)

```constraint
type Password : <id: string, value: string, secure: bool>
{p: Password} / NOT(LENGTH(p.value) >= 8) ==> weak_password_found(p.id, p.value)

```

#### 3️⃣ Faits Soumis (.facts)

```facts
Password[id="P001", value="password123", secure=true]Password[id="P002", value="123", secure=false]Password[id="P003", value="pass", secure=false]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `weak_password_found`
**Logique:** NOT(LENGTH(p.value) >= 8) → Mots de passe < 8 caractères

**Faits devant déclencher l'action:**
- P002 (length < 8)

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[balance=500, active=false, id=ACC002])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[category=books, id=PROD002, price=50])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[category=electronics, id=PROD003, price=200])
- ✅ expensive_product (Product[price=150, category=electronics, id=PROD001])
- ✅ expensive_product (Product[price=200, category=electronics, id=PROD003])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ small_balance_found (Balance[type=credit, id=B001, amount=75])
- ✅ small_balance_found (Balance[amount=75, type=credit, id=B001])
- ✅ normal_message_found (Message[content=Regular message content, urgent=false, id=M002])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[urgent=false, id=M003, content=Simple notification])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[content=Very urgent matter!, urgent=true, id=M003])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 1 attendues

⚠️ **Écart:** 44 actions obtenues vs 1 attendues

---

### 🧪 Test: `alpha_length_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_length_positive
**Description:** Test fonction LENGTH()

#### 2️⃣ Règles Complètes (.constraint)

```constraint
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
**Logique:** LENGTH(p.value) >= 8 → Mots de passe ≥ 8 caractères

**Faits devant déclencher l'action:**
- P001 (length ≥ 8)
- P003 (length ≥ 8)

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[balance=500, active=false, id=ACC002])
- ✅ active_account_found (Account[active=true, id=ACC001, balance=1000])
- ✅ active_account_found (Account[active=true, id=ACC003, balance=2000])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[active=true, id=ACC003, balance=2000])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[category=electronics, id=PROD003, price=200])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[status=active, id=P002, age=30])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[status=active, id=P001, age=25])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[role=admin, id=U001, name=Alice])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[urgent=false, id=M003, content=Simple notification])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ urgent_message_found (Message[content=This is urgent please respond, urgent=true, id=M001])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[urgent=true, id=M003, content=Very urgent matter!])
- ✅ secure_password_found (Password[value=password123, secure=true, id=P001])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 2 attendues

⚠️ **Écart:** 44 actions obtenues vs 2 attendues

---

### 🧪 Test: `alpha_like_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_like_negative
**Description:** Test opérateur LIKE (motif) - négation

#### 2️⃣ Règles Complètes (.constraint)

```constraint
type Email : <id: string, address: string, verified: bool>
{e: Email} / NOT(e.address LIKE "%@company.com") ==> external_email_found(e.id, e.address)

```

#### 3️⃣ Faits Soumis (.facts)

```facts
Email[id="E001", address="john@company.com", verified=true]Email[id="E002", address="jane@external.org", verified=false]Email[id="E003", address="user@other.net", verified=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `non_company_email_found`
**Logique:** NOT(e.address LIKE "%@company.com") → Emails non-entreprise

**Faits devant déclencher l'action:**
- E002 (@gmail.com)

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[price=200, category=electronics, id=PROD003])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[status=active, id=P002, age=30])
- ✅ age_is_twenty_five (Person[status=active, id=P001, age=25])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[role=user, id=U002, name=Bob])
- ✅ non_admin_user_found (User[role=user, id=U002, name=Bob])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[role=admin, id=U001, name=Alice])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ secure_password_found (Password[value=password123, secure=true, id=P001])
- ✅ secure_password_found (Password[secure=true, id=P003, value=verysecurepass])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[secure=true, id=P003, value=verysecurepass])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 1 attendues

⚠️ **Écart:** 44 actions obtenues vs 1 attendues

---

### 🧪 Test: `alpha_like_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_like_positive
**Description:** Test opérateur LIKE (motif)

#### 2️⃣ Règles Complètes (.constraint)

```constraint
type Email : <id: string, address: string, verified: bool>
{e: Email} / e.address LIKE "%@company.com" ==> company_email_found(e.id, e.address)

```

#### 3️⃣ Faits Soumis (.facts)

```facts
Email[id="E001", address="john@company.com", verified=true]Email[id="E002", address="jane@external.org", verified=false]Email[id="E003", address="admin@company.com", verified=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `company_email_found`
**Logique:** e.address LIKE "%@company.com" → Emails d'entreprise

**Faits devant déclencher l'action:**
- E001 (@company.com)
- E003 (@company.com)

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[balance=2000, active=true, id=ACC003])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[price=200, category=electronics, id=PROD003])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[age=30, status=active, id=P002])
- ✅ age_is_twenty_five (Person[status=active, id=P001, age=25])
- ✅ age_is_twenty_five (Person[status=inactive, id=P003, age=25])
- ✅ age_is_twenty_five (Person[status=active, id=P001, age=25])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[total=300, status=completed, id=ORD003])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[name=Bob, role=user, id=U002])
- ✅ non_admin_user_found (User[name=Bob, role=user, id=U002])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[content=Simple notification, urgent=false, id=M003])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ urgent_message_found (Message[urgent=true, id=M001, content=This is urgent please respond])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 2 attendues

⚠️ **Écart:** 44 actions obtenues vs 2 attendues

---

### 🧪 Test: `alpha_matches_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_matches_negative
**Description:** Test opérateur MATCHES (regex) - négation

#### 2️⃣ Règles Complètes (.constraint)

```constraint
type Code : <id: string, value: string, active: bool>
{c: Code} / NOT(c.value MATCHES "CODE[0-9]+") ==> invalid_code_found(c.id, c.value)

```

#### 3️⃣ Faits Soumis (.facts)

```facts
Code[id="C001", value="CODE123", active=true]Code[id="C002", value="INVALID", active=false]Code[id="C003", value="BADFORMAT", active=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `invalid_code_found`
**Logique:** NOT(c.value MATCHES "[A-Z]+[0-9]+") → Codes ne matchant pas

**Faits devant déclencher l'action:**
- C002 (pattern invalide)

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[active=false, id=ACC002, balance=500])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[balance=2000, active=true, id=ACC003])
- ✅ active_account_found (Account[active=true, id=ACC001, balance=1000])
- ✅ active_account_found (Account[active=true, id=ACC003, balance=2000])
- ✅ affordable_product (Product[price=50, category=books, id=PROD002])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[age=25, status=inactive, id=P003])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[total=300, status=completed, id=ORD003])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[name=Bob, role=user, id=U002])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[name=Alice, role=admin, id=U001])
- ✅ admin_user_found (User[name=Charlie, role=admin, id=U003])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ small_balance_found (Balance[amount=75, type=credit, id=B001])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ urgent_message_found (Message[urgent=true, id=M001, content=This is urgent please respond])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[content=Very urgent matter!, urgent=true, id=M003])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 1 attendues

⚠️ **Écart:** 44 actions obtenues vs 1 attendues

---

### 🧪 Test: `alpha_matches_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_matches_positive
**Description:** Test opérateur MATCHES (regex)

#### 2️⃣ Règles Complètes (.constraint)

```constraint
type Code : <id: string, value: string, active: bool>
{c: Code} / c.value MATCHES "CODE[0-9]+" ==> valid_code_found(c.id, c.value)

```

#### 3️⃣ Faits Soumis (.facts)

```facts
Code[id="C001", value="CODE123", active=true]Code[id="C002", value="INVALID", active=false]Code[id="C003", value="CODE999", active=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `valid_code_found`
**Logique:** c.value MATCHES "[A-Z]+[0-9]+" → Codes matchant

**Faits devant déclencher l'action:**
- C001 (CODE123)
- C003 (PROD456)

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[category=books, id=PROD002, price=50])
- ✅ expensive_product (Product[price=150, category=electronics, id=PROD001])
- ✅ expensive_product (Product[price=200, category=electronics, id=PROD003])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[status=active, id=P001, age=25])
- ✅ age_is_twenty_five (Person[age=25, status=inactive, id=P003])
- ✅ age_is_twenty_five (Person[age=25, status=active, id=P001])
- ✅ age_is_twenty_five (Person[age=25, status=inactive, id=P003])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[role=admin, id=U003, name=Charlie])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ urgent_message_found (Message[urgent=true, id=M001, content=This is urgent please respond])
- ✅ urgent_message_found (Message[content=Very urgent matter!, urgent=true, id=M003])
- ✅ urgent_message_found (Message[urgent=true, id=M001, content=This is urgent please respond])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ secure_password_found (Password[secure=true, id=P001, value=password123])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 2 attendues

⚠️ **Écart:** 44 actions obtenues vs 2 attendues

---

### 🧪 Test: `alpha_upper_negative`

#### 1️⃣ Description du Test

**Nom:** alpha_upper_negative
**Description:** Test fonction UPPER() (majuscules) - négation

#### 2️⃣ Règles Complètes (.constraint)

```constraint
type Department : <id: string, name: string, active: bool>
{d: Department} / NOT(UPPER(d.name) == "FINANCE") ==> non_finance_dept_found(d.id, d.name)

```

#### 3️⃣ Faits Soumis (.facts)

```facts
Department[id="D001", name="finance", active=true]Department[id="D002", name="IT", active=true]Department[id="D003", name="HR", active=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `lowercase_department_found`
**Logique:** NOT(UPPER(d.name) = d.name) → Noms non en majuscules

**Faits devant déclencher l'action:**
- D002 (sales ≠ SALES)

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[category=books, id=PROD002, price=50])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[price=150, category=electronics, id=PROD001])
- ✅ expensive_product (Product[price=200, category=electronics, id=PROD003])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[status=active, id=P002, age=30])
- ✅ age_is_twenty_five (Person[age=25, status=active, id=P001])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[age=25, status=active, id=P001])
- ✅ age_is_twenty_five (Person[age=25, status=inactive, id=P003])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[status=cancelled, id=ORD002, total=200])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ small_balance_found (Balance[type=credit, id=B001, amount=75])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[content=Very urgent matter!, urgent=true, id=M003])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[value=verysecurepass, secure=true, id=P003])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 1 attendues

⚠️ **Écart:** 44 actions obtenues vs 1 attendues

---

### 🧪 Test: `alpha_upper_positive`

#### 1️⃣ Description du Test

**Nom:** alpha_upper_positive
**Description:** Test fonction UPPER() (majuscules)

#### 2️⃣ Règles Complètes (.constraint)

```constraint
type Department : <id: string, name: string, active: bool>
{d: Department} / UPPER(d.name) == "FINANCE" ==> finance_dept_found(d.id, d.name)

```

#### 3️⃣ Faits Soumis (.facts)

```facts
Department[id="D001", name="finance", active=true]Department[id="D002", name="IT", active=true]Department[id="D003", name="Finance", active=true]
```

#### 4️⃣ Résultat Attendu

**Action attendue:** `uppercase_department_found`
**Logique:** UPPER(d.name) = d.name → Noms en majuscules

**Faits devant déclencher l'action:**
- D001 (FINANCE)
- D003 (HR)

#### 5️⃣ Résultat Obtenu

**Statut:** ❌ Échec

**Actions déclenchées:**
- ✅ inactive_account_found (Account[id=ACC002, balance=500, active=false])
- ✅ inactive_account_found (Account[active=false, id=ACC002, balance=500])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[active=true, id=ACC003, balance=2000])
- ✅ active_account_found (Account[id=ACC001, balance=1000, active=true])
- ✅ active_account_found (Account[id=ACC003, balance=2000, active=true])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ affordable_product (Product[id=PROD002, price=50, category=books])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ expensive_product (Product[id=PROD001, price=150, category=electronics])
- ✅ expensive_product (Product[id=PROD003, price=200, category=electronics])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_not_twenty_five (Person[id=P002, age=30, status=active])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ age_is_twenty_five (Person[id=P001, age=25, status=active])
- ✅ age_is_twenty_five (Person[id=P003, age=25, status=inactive])
- ✅ cancelled_order_found (Order[id=ORD002, total=200, status=cancelled])
- ✅ cancelled_order_found (Order[status=cancelled, id=ORD002, total=200])
- ✅ valid_order_found (Order[id=ORD001, total=100, status=pending])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ valid_order_found (Order[status=pending, id=ORD001, total=100])
- ✅ valid_order_found (Order[id=ORD003, total=300, status=completed])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ non_admin_user_found (User[id=U002, name=Bob, role=user])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ admin_user_found (User[id=U001, name=Alice, role=admin])
- ✅ admin_user_found (User[id=U003, name=Charlie, role=admin])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ small_balance_found (Balance[id=B001, amount=75, type=credit])
- ✅ normal_message_found (Message[content=Regular message content, urgent=false, id=M002])
- ✅ normal_message_found (Message[id=M003, content=Simple notification, urgent=false])
- ✅ normal_message_found (Message[id=M002, content=Regular message content, urgent=false])
- ✅ normal_message_found (Message[urgent=false, id=M003, content=Simple notification])
- ✅ urgent_message_found (Message[id=M001, content=This is urgent please respond, urgent=true])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ urgent_message_found (Message[urgent=true, id=M001, content=This is urgent please respond])
- ✅ urgent_message_found (Message[id=M003, content=Very urgent matter!, urgent=true])
- ✅ secure_password_found (Password[value=password123, secure=true, id=P001])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])
- ✅ secure_password_found (Password[id=P001, value=password123, secure=true])
- ✅ secure_password_found (Password[id=P003, value=verysecurepass, secure=true])

#### 6️⃣ Analyse du Test

**Résultat:** ⚠️ PARTIEL - 44 actions obtenues vs 2 attendues

⚠️ **Écart:** 44 actions obtenues vs 2 attendues

---

## 🎯 Analyse Globale

### 📊 Statistiques Générales

- **Tests exécutés:** 26
- **Tests réussis:** 10
- **Taux de conformité:** 38.5%

### 📈 Analyse par Catégorie

**Tests Alpha Originaux:**
- Succès: 10/10 (100.0%)

**Tests Alpha Étendus:**
- Succès: 0/16 (0.0%)

### 🏁 Conclusions

⚠️ **AMÉLIORATION REQUISE:** Plusieurs tests échouent

🔧 Des corrections importantes sont nécessaires

---

**Rapport généré par:** `generate_structured_test_report.py`
**Horodatage:** 2025-11-17T15:17:24.788071