# 📊 RAPPORT COMPLET - TESTS DE COUVERTURE ALPHA ÉTENDUS

**Date d'exécution:** 2025-11-17 11:37:31
**Nombre de tests:** 26

## 🎯 RÉSUMÉ EXÉCUTIF

- ✅ **Tests réussis:** 26/26 (100.0%)
- 🎬 **Actions déclenchées:** 42
- 📊 **Tests originaux:** 10
- 🆕 **Tests étendus:** 16
- ⚡ **Couverture:** Nœuds Alpha complets avec tous opérateurs/fonctions

## 📈 MATRICE DE COUVERTURE

### Opérateurs testés

| Opérateur | Tests | Succès | Taux |
|-----------|-------|--------|------|
| `IN` | 2 | 2 | 100.0% |
| `LIKE` | 2 | 2 | 100.0% |
| `MATCHES` | 2 | 2 | 100.0% |
| `==` | 8 | 8 | 100.0% |
| `>` | 2 | 2 | 100.0% |
| `CONTAINS` | 2 | 2 | 100.0% |
| `=` | 2 | 2 | 100.0% |

### Fonctions testées

| Fonction | Tests | Succès | Taux |
|----------|-------|--------|------|
| `UPPER()` | 2 | 2 | 100.0% |
| `ABS()` | 2 | 2 | 100.0% |
| `LENGTH()` | 2 | 2 | 100.0% |

## 🧪 DÉTAILS DES TESTS

### 🧪 TEST 1: alpha_boolean_negative

#### 📋 Informations générales

- **Type:** ORIGINAL
- **Opérateur testé:** `==`
- **Temps d'exécution:** 617.826µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Account]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 2: alpha_boolean_positive

#### 📋 Informations générales

- **Type:** ORIGINAL
- **Opérateur testé:** `==`
- **Temps d'exécution:** 395.78µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Account]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 3: alpha_comparison_negative

#### 📋 Informations générales

- **Type:** ORIGINAL
- **Opérateur testé:** `>`
- **Temps d'exécution:** 362.407µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Product]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 4: alpha_comparison_positive

#### 📋 Informations générales

- **Type:** ORIGINAL
- **Opérateur testé:** `>`
- **Temps d'exécution:** 352.068µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Product]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 5: alpha_equality_negative

#### 📋 Informations générales

- **Type:** ORIGINAL
- **Opérateur testé:** `==`
- **Temps d'exécution:** 254.526µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Person]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 6: alpha_equality_positive

#### 📋 Informations générales

- **Type:** ORIGINAL
- **Opérateur testé:** `==`
- **Temps d'exécution:** 748.57µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Person]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 7: alpha_inequality_negative

#### 📋 Informations générales

- **Type:** ORIGINAL
- **Opérateur testé:** `==`
- **Temps d'exécution:** 433.311µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Order]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 8: alpha_inequality_positive

#### 📋 Informations générales

- **Type:** ORIGINAL
- **Opérateur testé:** `==`
- **Temps d'exécution:** 464.088µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Order]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 9: alpha_string_negative

#### 📋 Informations générales

- **Type:** ORIGINAL
- **Opérateur testé:** `==`
- **Temps d'exécution:** 515.895µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [User]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 10: alpha_string_positive

#### 📋 Informations générales

- **Type:** ORIGINAL
- **Opérateur testé:** `==`
- **Temps d'exécution:** 353.381µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [User]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 11: alpha_abs_negative

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `ABS()`
- **Temps d'exécution:** 329.405µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 3

#### ⚡ Actions déclenchées

**Action:** `action_for_B001` (1 fois)
- Fact{ID:B001, Type:Balance, Fields:map[amount:150 id:B001 type:credit]}

**Action:** `action_for_B002` (1 fois)
- Fact{ID:B002, Type:Balance, Fields:map[amount:-25 id:B002 type:debit]}

**Action:** `action_for_B003` (1 fois)
- Fact{ID:B003, Type:Balance, Fields:map[amount:75 id:B003 type:credit]}

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Balance]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 12: alpha_abs_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `ABS()`
- **Temps d'exécution:** 306.463µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 3

#### ⚡ Actions déclenchées

**Action:** `action_for_B002` (1 fois)
- Fact{ID:B002, Type:Balance, Fields:map[amount:-200 id:B002 type:debit]}

**Action:** `action_for_B003` (1 fois)
- Fact{ID:B003, Type:Balance, Fields:map[amount:50 id:B003 type:credit]}

**Action:** `action_for_B001` (1 fois)
- Fact{ID:B001, Type:Balance, Fields:map[amount:150 id:B001 type:credit]}

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Balance]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 13: alpha_contains_negative

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `CONTAINS`
- **Temps d'exécution:** 332.552µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 3

#### ⚡ Actions déclenchées

**Action:** `action_for_M001` (1 fois)
- Fact{ID:M001, Type:Message, Fields:map[content:This is urgent please respond id:M001 urgent:true]}

**Action:** `action_for_M002` (1 fois)
- Fact{ID:M002, Type:Message, Fields:map[content:Regular message content id:M002 urgent:false]}

**Action:** `action_for_M003` (1 fois)
- Fact{ID:M003, Type:Message, Fields:map[content:Simple notification id:M003 urgent:false]}

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Message]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 14: alpha_contains_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `CONTAINS`
- **Temps d'exécution:** 326.03µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 3

#### ⚡ Actions déclenchées

**Action:** `action_for_M003` (1 fois)
- Fact{ID:M003, Type:Message, Fields:map[content:Very urgent matter! id:M003 urgent:true]}

**Action:** `action_for_M001` (1 fois)
- Fact{ID:M001, Type:Message, Fields:map[content:This is urgent please respond id:M001 urgent:true]}

**Action:** `action_for_M002` (1 fois)
- Fact{ID:M002, Type:Message, Fields:map[content:Regular message content id:M002 urgent:false]}

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Message]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 15: alpha_equal_sign_negative

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `=`
- **Temps d'exécution:** 409.215µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Customer]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 16: alpha_equal_sign_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `=`
- **Temps d'exécution:** 454.881µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Customer]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 17: alpha_in_negative

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `IN`
- **Temps d'exécution:** 437.107µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 3

#### ⚡ Actions déclenchées

**Action:** `action_for_S001` (1 fois)
- Fact{ID:S001, Type:Status, Fields:map[id:S001 priority:1 state:active]}

**Action:** `action_for_S002` (1 fois)
- Fact{ID:S002, Type:Status, Fields:map[id:S002 priority:3 state:inactive]}

**Action:** `action_for_S003` (1 fois)
- Fact{ID:S003, Type:Status, Fields:map[id:S003 priority:5 state:archived]}

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Status]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 18: alpha_in_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `IN`
- **Temps d'exécution:** 700.51µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 3

#### ⚡ Actions déclenchées

**Action:** `action_for_S001` (1 fois)
- Fact{ID:S001, Type:Status, Fields:map[id:S001 priority:1 state:active]}

**Action:** `action_for_S002` (1 fois)
- Fact{ID:S002, Type:Status, Fields:map[id:S002 priority:3 state:inactive]}

**Action:** `action_for_S003` (1 fois)
- Fact{ID:S003, Type:Status, Fields:map[id:S003 priority:2 state:pending]}

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Status]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 19: alpha_length_negative

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `LENGTH()`
- **Temps d'exécution:** 483.374µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 3

#### ⚡ Actions déclenchées

**Action:** `action_for_P001` (1 fois)
- Fact{ID:P001, Type:Password, Fields:map[id:P001 secure:true value:password123]}

**Action:** `action_for_P002` (1 fois)
- Fact{ID:P002, Type:Password, Fields:map[id:P002 secure:false value:123]}

**Action:** `action_for_P003` (1 fois)
- Fact{ID:P003, Type:Password, Fields:map[id:P003 secure:false value:abc]}

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Password]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 20: alpha_length_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `LENGTH()`
- **Temps d'exécution:** 386.894µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 3

#### ⚡ Actions déclenchées

**Action:** `action_for_P002` (1 fois)
- Fact{ID:P002, Type:Password, Fields:map[id:P002 secure:false value:123]}

**Action:** `action_for_P003` (1 fois)
- Fact{ID:P003, Type:Password, Fields:map[id:P003 secure:true value:verysecurepass]}

**Action:** `action_for_P001` (1 fois)
- Fact{ID:P001, Type:Password, Fields:map[id:P001 secure:true value:password123]}

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Password]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 21: alpha_like_negative

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `LIKE`
- **Temps d'exécution:** 409.887µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 3

#### ⚡ Actions déclenchées

**Action:** `action_for_E002` (1 fois)
- Fact{ID:E002, Type:Email, Fields:map[address:jane@external.org id:E002 verified:false]}

**Action:** `action_for_E003` (1 fois)
- Fact{ID:E003, Type:Email, Fields:map[address:user@other.net id:E003 verified:true]}

**Action:** `action_for_E001` (1 fois)
- Fact{ID:E001, Type:Email, Fields:map[address:john@company.com id:E001 verified:true]}

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Email]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 22: alpha_like_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `LIKE`
- **Temps d'exécution:** 350.114µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 3

#### ⚡ Actions déclenchées

**Action:** `action_for_E001` (1 fois)
- Fact{ID:E001, Type:Email, Fields:map[address:john@company.com id:E001 verified:true]}

**Action:** `action_for_E002` (1 fois)
- Fact{ID:E002, Type:Email, Fields:map[address:jane@external.org id:E002 verified:false]}

**Action:** `action_for_E003` (1 fois)
- Fact{ID:E003, Type:Email, Fields:map[address:bob@company.com id:E003 verified:true]}

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Email]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 23: alpha_matches_negative

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `MATCHES`
- **Temps d'exécution:** 271.878µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 3

#### ⚡ Actions déclenchées

**Action:** `action_for_C003` (1 fois)
- Fact{ID:C003, Type:Code, Fields:map[active:false id:C003 value:BADFORMAT]}

**Action:** `action_for_C001` (1 fois)
- Fact{ID:C001, Type:Code, Fields:map[active:true id:C001 value:CODE123]}

**Action:** `action_for_C002` (1 fois)
- Fact{ID:C002, Type:Code, Fields:map[active:false id:C002 value:INVALID]}

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Code]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 24: alpha_matches_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `MATCHES`
- **Temps d'exécution:** 247.182µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 3

#### ⚡ Actions déclenchées

**Action:** `action_for_C002` (1 fois)
- Fact{ID:C002, Type:Code, Fields:map[active:false id:C002 value:INVALID]}

**Action:** `action_for_C003` (1 fois)
- Fact{ID:C003, Type:Code, Fields:map[active:true id:C003 value:CODE999]}

**Action:** `action_for_C001` (1 fois)
- Fact{ID:C001, Type:Code, Fields:map[active:true id:C001 value:CODE123]}

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Code]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 25: alpha_upper_negative

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `UPPER()`
- **Temps d'exécution:** 356.416µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 3

#### ⚡ Actions déclenchées

**Action:** `action_for_D001` (1 fois)
- Fact{ID:D001, Type:Department, Fields:map[active:true id:D001 name:finance]}

**Action:** `action_for_D002` (1 fois)
- Fact{ID:D002, Type:Department, Fields:map[active:true id:D002 name:IT]}

**Action:** `action_for_D003` (1 fois)
- Fact{ID:D003, Type:Department, Fields:map[active:true id:D003 name:HR]}

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Department]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 26: alpha_upper_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `UPPER()`
- **Temps d'exécution:** 262.501µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 3

#### ⚡ Actions déclenchées

**Action:** `action_for_D001` (1 fois)
- Fact{ID:D001, Type:Department, Fields:map[active:true id:D001 name:finance]}

**Action:** `action_for_D002` (1 fois)
- Fact{ID:D002, Type:Department, Fields:map[active:true id:D002 name:IT]}

**Action:** `action_for_D003` (1 fois)
- Fact{ID:D003, Type:Department, Fields:map[active:true id:D003 name:Finance]}

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Department]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

