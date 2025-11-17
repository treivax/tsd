# 📊 RAPPORT COMPLET - TESTS DE COUVERTURE ALPHA ÉTENDUS

**Date d'exécution:** 2025-11-17 14:55:32
**Nombre de tests:** 26

## 🎯 RÉSUMÉ EXÉCUTIF

- ✅ **Tests réussis:** 15/26 (57.7%)
- 🎬 **Actions déclenchées:** 0
- 📊 **Tests originaux:** 10
- 🆕 **Tests étendus:** 16
- ⚡ **Couverture:** Nœuds Alpha complets avec tous opérateurs/fonctions

## 📈 MATRICE DE COUVERTURE

### Opérateurs testés

| Opérateur | Tests | Succès | Taux |
|-----------|-------|--------|------|
| `==` | 8 | 8 | 100.0% |
| `>` | 2 | 2 | 100.0% |
| `CONTAINS` | 2 | 2 | 100.0% |
| `=` | 2 | 0 | 0.0% |
| `IN` | 2 | 0 | 0.0% |
| `LIKE` | 2 | 0 | 0.0% |
| `MATCHES` | 2 | 0 | 0.0% |

### Fonctions testées

| Fonction | Tests | Succès | Taux |
|----------|-------|--------|------|
| `ABS()` | 2 | 2 | 100.0% |
| `LENGTH()` | 2 | 1 | 50.0% |
| `UPPER()` | 2 | 0 | 0.0% |

## 🧪 DÉTAILS DES TESTS

### 🧪 TEST 1: alpha_boolean_negative

#### 📋 Informations générales

- **Type:** ORIGINAL
- **Opérateur testé:** `==`
- **Temps d'exécution:** 388.486µs
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
- **Temps d'exécution:** 229.459µs
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
- **Temps d'exécution:** 239.658µs
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
- **Temps d'exécution:** 213.93µs
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
- **Temps d'exécution:** 195.165µs
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
- **Temps d'exécution:** 205.163µs
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
- **Temps d'exécution:** 231.262µs
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
- **Temps d'exécution:** 234.849µs
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
- **Temps d'exécution:** 223.888µs
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
- **Temps d'exécution:** 199.122µs
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
- **Temps d'exécution:** 200.445µs
- **Faits analysés:** 1
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Balance]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 12: alpha_abs_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `ABS()`
- **Temps d'exécution:** 160.67µs
- **Faits analysés:** 1
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Balance]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 13: alpha_contains_negative

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `CONTAINS`
- **Temps d'exécution:** 240.82µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Message]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 14: alpha_contains_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `CONTAINS`
- **Temps d'exécution:** 226.163µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Message]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 15: alpha_equal_sign_negative

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `=`
- **Temps d'exécution:** 161.151µs
- **Faits analysés:** 0
- **Statut:** ❌ Échec
- **Erreur:** Erreur construction réseau: erreur parsing faits: erreur fait ligne 1: erreur parsing champs: erreur conversion champ points: strconv.ParseInt: parsing "5000]Customer[id=\"C002\"": invalid syntax

#### 🕸️ Structure réseau RETE

- **TypeNodes:** []
- **AlphaNodes:** []
- **TerminalNodes:** []

---

### 🧪 TEST 16: alpha_equal_sign_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `=`
- **Temps d'exécution:** 150.952µs
- **Faits analysés:** 0
- **Statut:** ❌ Échec
- **Erreur:** Erreur construction réseau: erreur parsing faits: erreur fait ligne 1: erreur parsing champs: erreur conversion champ points: strconv.ParseInt: parsing "5000]Customer[id=\"C002\"": invalid syntax

#### 🕸️ Structure réseau RETE

- **TypeNodes:** []
- **AlphaNodes:** []
- **TerminalNodes:** []

---

### 🧪 TEST 17: alpha_in_negative

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `IN`
- **Temps d'exécution:** 176.871µs
- **Faits analysés:** 0
- **Statut:** ❌ Échec
- **Erreur:** Erreur construction réseau: erreur parsing faits: erreur fait ligne 1: erreur parsing champs: erreur conversion champ priority: strconv.ParseInt: parsing "1]Status[id=\"S002\"": invalid syntax

#### 🕸️ Structure réseau RETE

- **TypeNodes:** []
- **AlphaNodes:** []
- **TerminalNodes:** []

---

### 🧪 TEST 18: alpha_in_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `IN`
- **Temps d'exécution:** 198.351µs
- **Faits analysés:** 0
- **Statut:** ❌ Échec
- **Erreur:** Erreur construction réseau: erreur parsing faits: erreur fait ligne 1: erreur parsing champs: erreur conversion champ priority: strconv.ParseInt: parsing "1]Status[id=\"S002\"": invalid syntax

#### 🕸️ Structure réseau RETE

- **TypeNodes:** []
- **AlphaNodes:** []
- **TerminalNodes:** []

---

### 🧪 TEST 19: alpha_length_negative

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `LENGTH()`
- **Temps d'exécution:** 154.659µs
- **Faits analysés:** 0
- **Statut:** ❌ Échec
- **Erreur:** Erreur construction réseau: erreur parsing faits: erreur fait ligne 1: erreur parsing champs: erreur conversion champ secure: strconv.ParseBool: parsing "true]Password[id=\"P002\"": invalid syntax

#### 🕸️ Structure réseau RETE

- **TypeNodes:** []
- **AlphaNodes:** []
- **TerminalNodes:** []

---

### 🧪 TEST 20: alpha_length_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `LENGTH()`
- **Temps d'exécution:** 229.529µs
- **Faits analysés:** 3
- **Statut:** ✅ Succès
- **Actions déclenchées:** 0

#### 🕸️ Structure réseau RETE

- **TypeNodes:** [Password]
- **AlphaNodes:** [rule_0_alpha]
- **TerminalNodes:** [rule_0_terminal]

---

### 🧪 TEST 21: alpha_like_negative

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `LIKE`
- **Temps d'exécution:** 161.672µs
- **Faits analysés:** 0
- **Statut:** ❌ Échec
- **Erreur:** Erreur construction réseau: erreur parsing faits: erreur fait ligne 1: erreur parsing champs: erreur conversion champ verified: strconv.ParseBool: parsing "true]Email[id=\"E002\"": invalid syntax

#### 🕸️ Structure réseau RETE

- **TypeNodes:** []
- **AlphaNodes:** []
- **TerminalNodes:** []

---

### 🧪 TEST 22: alpha_like_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `LIKE`
- **Temps d'exécution:** 155.802µs
- **Faits analysés:** 0
- **Statut:** ❌ Échec
- **Erreur:** Erreur construction réseau: erreur parsing faits: erreur fait ligne 1: erreur parsing champs: erreur conversion champ verified: strconv.ParseBool: parsing "true]Email[id=\"E002\"": invalid syntax

#### 🕸️ Structure réseau RETE

- **TypeNodes:** []
- **AlphaNodes:** []
- **TerminalNodes:** []

---

### 🧪 TEST 23: alpha_matches_negative

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `MATCHES`
- **Temps d'exécution:** 148.047µs
- **Faits analysés:** 0
- **Statut:** ❌ Échec
- **Erreur:** Erreur construction réseau: erreur parsing faits: erreur fait ligne 1: erreur parsing champs: erreur conversion champ active: strconv.ParseBool: parsing "true]Code[id=\"C002\"": invalid syntax

#### 🕸️ Structure réseau RETE

- **TypeNodes:** []
- **AlphaNodes:** []
- **TerminalNodes:** []

---

### 🧪 TEST 24: alpha_matches_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `MATCHES`
- **Temps d'exécution:** 154.068µs
- **Faits analysés:** 0
- **Statut:** ❌ Échec
- **Erreur:** Erreur construction réseau: erreur parsing faits: erreur fait ligne 1: erreur parsing champs: erreur conversion champ active: strconv.ParseBool: parsing "true]Code[id=\"C002\"": invalid syntax

#### 🕸️ Structure réseau RETE

- **TypeNodes:** []
- **AlphaNodes:** []
- **TerminalNodes:** []

---

### 🧪 TEST 25: alpha_upper_negative

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `UPPER()`
- **Temps d'exécution:** 172.292µs
- **Faits analysés:** 0
- **Statut:** ❌ Échec
- **Erreur:** Erreur construction réseau: erreur parsing faits: erreur fait ligne 1: erreur parsing champs: erreur conversion champ active: strconv.ParseBool: parsing "true]Department[id=\"D002\"": invalid syntax

#### 🕸️ Structure réseau RETE

- **TypeNodes:** []
- **AlphaNodes:** []
- **TerminalNodes:** []

---

### 🧪 TEST 26: alpha_upper_positive

#### 📋 Informations générales

- **Type:** EXTENDED
- **Opérateur testé:** `UPPER()`
- **Temps d'exécution:** 146.093µs
- **Faits analysés:** 0
- **Statut:** ❌ Échec
- **Erreur:** Erreur construction réseau: erreur parsing faits: erreur fait ligne 1: erreur parsing champs: erreur conversion champ active: strconv.ParseBool: parsing "true]Department[id=\"D002\"": invalid syntax

#### 🕸️ Structure réseau RETE

- **TypeNodes:** []
- **AlphaNodes:** []
- **TerminalNodes:** []

---

