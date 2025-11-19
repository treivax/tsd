# RAPPORT ALPHA NODES - DÉTAILLÉ
==================================

**Tests:** 26
**Succès:** 26 (100.0%)
**Date:** 2025-11-19 22:49:28

## Test 1: alpha_abs_negative
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_abs_negative.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_abs_negative.facts
- **Temps:** 748.309µs

### Règles Analysées
- `{b: Balance} / NOT(ABS(b.amount) > 100) ==> small_balance_found(b.id, b.amount)`

### Faits Traités
- `Balance(id:"B001", amount:150.0, type:"credit")`
- `Balance(id:"B002", amount:-25.0, type:"debit")`
- `Balance(id:"B003", amount:75.0, type:"credit")`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: small_balance_found**
Faits:
  - Balance {id: B002, amount: -25, type: debit}

**Token 2 (Node: rule_0_expected) → Action: small_balance_found**
Faits:
  - Balance {id: B003, amount: 75, type: credit}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: small_balance_found**
Faits:
  - Balance {type: debit, id: B002, amount: -25}

**Token 2 (Node: rule_0_terminal) → Action: small_balance_found**
Faits:
  - Balance {id: B003, amount: 75, type: credit}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: small_balance_found(b.id, b.amount)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 748.309µs

---

## Test 2: alpha_abs_positive
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_abs_positive.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_abs_positive.facts
- **Temps:** 520.964µs

### Règles Analysées
- `{b: Balance} / ABS(b.amount) > 100 ==> significant_balance_found(b.id, b.amount)`

### Faits Traités
- `Balance(id:"B001", amount:150.0, type:"credit")`
- `Balance(id:"B002", amount:-200.0, type:"debit")`
- `Balance(id:"B003", amount:50.0, type:"credit")`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: significant_balance_found**
Faits:
  - Balance {id: B001, amount: 150, type: credit}

**Token 2 (Node: rule_0_expected) → Action: significant_balance_found**
Faits:
  - Balance {type: debit, id: B002, amount: -200}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: significant_balance_found**
Faits:
  - Balance {type: credit, id: B001, amount: 150}

**Token 2 (Node: rule_0_terminal) → Action: significant_balance_found**
Faits:
  - Balance {id: B002, amount: -200, type: debit}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: significant_balance_found(b.id, b.amount)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 520.964µs

---

## Test 3: alpha_boolean_negative
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_boolean_negative.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_boolean_negative.facts
- **Temps:** 511.806µs

### Règles Analysées
- `{a: Account} / NOT(a.active == true) ==> inactive_account_found(a.id, a.balance)`

### Faits Traités
- `Account(id:ACC001, balance:1000, active:true)`
- `Account(id:ACC002, balance:500, active:false)`
- `Account(id:ACC003, balance:2000, active:true)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 1

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: inactive_account_found**
Faits:
  - Account {balance: 500, active: false, id: ACC002}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: inactive_account_found**
Faits:
  - Account {balance: 500, active: false, id: ACC002}


#### Comparaison
- **Matches:** 1
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: inactive_account_found(a.id, a.balance)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 1 trouvés
- Analyse tokens: 1 attendus, 1 observés, 1 matches, 0 mismatches
- Test terminé avec succès en 511.806µs

---

## Test 4: alpha_boolean_positive
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_boolean_positive.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_boolean_positive.facts
- **Temps:** 543.456µs

### Règles Analysées
- `{a: Account} / a.active == true ==> active_account_found(a.id, a.balance)`

### Faits Traités
- `Account(id:ACC001, balance:1000, active:true)`
- `Account(id:ACC002, balance:500, active:false)`
- `Account(id:ACC003, balance:2000, active:true)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: active_account_found**
Faits:
  - Account {id: ACC001, balance: 1000, active: true}

**Token 2 (Node: rule_0_expected) → Action: active_account_found**
Faits:
  - Account {id: ACC003, balance: 2000, active: true}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: active_account_found**
Faits:
  - Account {active: true, id: ACC001, balance: 1000}

**Token 2 (Node: rule_0_terminal) → Action: active_account_found**
Faits:
  - Account {id: ACC003, balance: 2000, active: true}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: active_account_found(a.id, a.balance)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 543.456µs

---

## Test 5: alpha_comparison_negative
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_comparison_negative.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_comparison_negative.facts
- **Temps:** 522.737µs

### Règles Analysées
- `{prod: Product} / NOT(prod.price > 100) ==> affordable_product(prod.id, prod.price)`

### Faits Traités
- `Product(id:PROD001, price:150, category:electronics)`
- `Product(id:PROD002, price:50, category:books)`
- `Product(id:PROD003, price:200, category:electronics)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 1

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: affordable_product**
Faits:
  - Product {id: PROD002, price: 50, category: books}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: affordable_product**
Faits:
  - Product {id: PROD002, price: 50, category: books}


#### Comparaison
- **Matches:** 1
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: affordable_product(prod.id, prod.price)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 1 trouvés
- Analyse tokens: 1 attendus, 1 observés, 1 matches, 0 mismatches
- Test terminé avec succès en 522.737µs

---

## Test 6: alpha_comparison_positive
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_comparison_positive.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_comparison_positive.facts
- **Temps:** 501.778µs

### Règles Analysées
- `{prod: Product} / prod.price > 100 ==> expensive_product(prod.id, prod.price)`

### Faits Traités
- `Product(id:PROD001, price:150, category:electronics)`
- `Product(id:PROD002, price:50, category:books)`
- `Product(id:PROD003, price:200, category:electronics)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: expensive_product**
Faits:
  - Product {id: PROD001, price: 150, category: electronics}

**Token 2 (Node: rule_0_expected) → Action: expensive_product**
Faits:
  - Product {category: electronics, id: PROD003, price: 200}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: expensive_product**
Faits:
  - Product {id: PROD001, price: 150, category: electronics}

**Token 2 (Node: rule_0_terminal) → Action: expensive_product**
Faits:
  - Product {category: electronics, id: PROD003, price: 200}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: expensive_product(prod.id, prod.price)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 501.778µs

---

## Test 7: alpha_contains_negative
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_contains_negative.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_contains_negative.facts
- **Temps:** 578.942µs

### Règles Analysées
- `{m: Message} / NOT(m.content CONTAINS "urgent") ==> normal_message_found(m.id, m.content)`

### Faits Traités
- `Message(id:"M001", content:"This is urgent please respond", urgent:true)`
- `Message(id:"M002", content:"Regular message content", urgent:false)`
- `Message(id:"M003", content:"Simple notification", urgent:false)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: normal_message_found**
Faits:
  - Message {content: Regular message content, urgent: false, id: M002}

**Token 2 (Node: rule_0_expected) → Action: normal_message_found**
Faits:
  - Message {id: M003, content: Simple notification, urgent: false}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: normal_message_found**
Faits:
  - Message {id: M002, content: Regular message content, urgent: false}

**Token 2 (Node: rule_0_terminal) → Action: normal_message_found**
Faits:
  - Message {id: M003, content: Simple notification, urgent: false}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: normal_message_found(m.id, m.content)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 578.942µs

---

## Test 8: alpha_contains_positive
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_contains_positive.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_contains_positive.facts
- **Temps:** 557.603µs

### Règles Analysées
- `{m: Message} / m.content CONTAINS "urgent" ==> urgent_message_found(m.id, m.content)`

### Faits Traités
- `Message(id:"M001", content:"This is urgent please respond", urgent:true)`
- `Message(id:"M002", content:"Regular message content", urgent:false)`
- `Message(id:"M003", content:"Very urgent matter!", urgent:true)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: urgent_message_found**
Faits:
  - Message {id: M001, content: This is urgent please respond, urgent: true}

**Token 2 (Node: rule_0_expected) → Action: urgent_message_found**
Faits:
  - Message {id: M003, content: Very urgent matter!, urgent: true}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: urgent_message_found**
Faits:
  - Message {id: M001, content: This is urgent please respond, urgent: true}

**Token 2 (Node: rule_0_terminal) → Action: urgent_message_found**
Faits:
  - Message {id: M003, content: Very urgent matter!, urgent: true}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: urgent_message_found(m.id, m.content)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 557.603µs

---

## Test 9: alpha_equal_sign_negative
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_equal_sign_negative.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_equal_sign_negative.facts
- **Temps:** 494.945µs

### Règles Analysées
- `{c: Customer} / NOT(c.tier == "gold") ==> non_gold_customer_found(c.id, c.tier)`

### Faits Traités
- `Customer(id:"C001", tier:"gold", points:5000)`
- `Customer(id:"C002", tier:"silver", points:2000)`
- `Customer(id:"C003", tier:"bronze", points:1000)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: non_gold_customer_found**
Faits:
  - Customer {id: C002, tier: silver, points: 2000}

**Token 2 (Node: rule_0_expected) → Action: non_gold_customer_found**
Faits:
  - Customer {id: C003, tier: bronze, points: 1000}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: non_gold_customer_found**
Faits:
  - Customer {points: 2000, id: C002, tier: silver}

**Token 2 (Node: rule_0_terminal) → Action: non_gold_customer_found**
Faits:
  - Customer {tier: bronze, points: 1000, id: C003}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: non_gold_customer_found(c.id, c.tier)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 494.945µs

---

## Test 10: alpha_equal_sign_positive
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_equal_sign_positive.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_equal_sign_positive.facts
- **Temps:** 510.825µs

### Règles Analysées
- `{c: Customer} / c.tier == "gold" ==> gold_customer_found(c.id, c.points)`

### Faits Traités
- `Customer(id:"C001", tier:"gold", points:5000)`
- `Customer(id:"C002", tier:"silver", points:2000)`
- `Customer(id:"C003", tier:"gold", points:7500)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: gold_customer_found**
Faits:
  - Customer {id: C001, tier: gold, points: 5000}

**Token 2 (Node: rule_0_expected) → Action: gold_customer_found**
Faits:
  - Customer {id: C003, tier: gold, points: 7500}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: gold_customer_found**
Faits:
  - Customer {id: C001, tier: gold, points: 5000}

**Token 2 (Node: rule_0_terminal) → Action: gold_customer_found**
Faits:
  - Customer {points: 7500, id: C003, tier: gold}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: gold_customer_found(c.id, c.points)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 510.825µs

---

## Test 11: alpha_equality_negative
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_equality_negative.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_equality_negative.facts
- **Temps:** 459.118µs

### Règles Analysées
- `{p: Person} / NOT(p.age == 25) ==> age_is_not_twenty_five(p.id, p.age)`

### Faits Traités
- `Person(id:P001, age:25, status:active)`
- `Person(id:P002, age:30, status:active)`
- `Person(id:P003, age:25, status:inactive)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 1

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: age_is_not_twenty_five**
Faits:
  - Person {id: P002, age: 30, status: active}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: age_is_not_twenty_five**
Faits:
  - Person {id: P002, age: 30, status: active}


#### Comparaison
- **Matches:** 1
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: age_is_not_twenty_five(p.id, p.age)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 1 trouvés
- Analyse tokens: 1 attendus, 1 observés, 1 matches, 0 mismatches
- Test terminé avec succès en 459.118µs

---

## Test 12: alpha_equality_positive
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_equality_positive.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_equality_positive.facts
- **Temps:** 530.271µs

### Règles Analysées
- `{p: Person} / p.age == 25 ==> age_is_twenty_five(p.id, p.age)`

### Faits Traités
- `Person(id:P001, age:25, status:active)`
- `Person(id:P002, age:30, status:active)`
- `Person(id:P003, age:25, status:inactive)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: age_is_twenty_five**
Faits:
  - Person {id: P001, age: 25, status: active}

**Token 2 (Node: rule_0_expected) → Action: age_is_twenty_five**
Faits:
  - Person {id: P003, age: 25, status: inactive}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: age_is_twenty_five**
Faits:
  - Person {id: P001, age: 25, status: active}

**Token 2 (Node: rule_0_terminal) → Action: age_is_twenty_five**
Faits:
  - Person {id: P003, age: 25, status: inactive}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: age_is_twenty_five(p.id, p.age)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 530.271µs

---

## Test 13: alpha_in_negative
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_in_negative.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_in_negative.facts
- **Temps:** 520.323µs

### Règles Analysées
- `{s: Status} / NOT(s.state IN ["active", "pending"]) ==> invalid_state_found(s.id, s.state)`

### Faits Traités
- `Status(id:"S001", state:"active", priority:1)`
- `Status(id:"S002", state:"inactive", priority:3)`
- `Status(id:"S003", state:"archived", priority:5)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: invalid_state_found**
Faits:
  - Status {id: S002, state: inactive, priority: 3}

**Token 2 (Node: rule_0_expected) → Action: invalid_state_found**
Faits:
  - Status {id: S003, state: archived, priority: 5}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: invalid_state_found**
Faits:
  - Status {priority: 5, id: S003, state: archived}

**Token 2 (Node: rule_0_terminal) → Action: invalid_state_found**
Faits:
  - Status {state: inactive, priority: 3, id: S002}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: invalid_state_found(s.id, s.state)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 520.323µs

---

## Test 14: alpha_in_positive
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_in_positive.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_in_positive.facts
- **Temps:** 574.034µs

### Règles Analysées
- `{s: Status} / s.state IN ["active", "pending", "review"] ==> valid_state_found(s.id, s.state)`

### Faits Traités
- `Status(id:"S001", state:"active", priority:1)`
- `Status(id:"S002", state:"inactive", priority:3)`
- `Status(id:"S003", state:"pending", priority:2)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: valid_state_found**
Faits:
  - Status {id: S001, state: active, priority: 1}

**Token 2 (Node: rule_0_expected) → Action: valid_state_found**
Faits:
  - Status {id: S003, state: pending, priority: 2}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: valid_state_found**
Faits:
  - Status {id: S001, state: active, priority: 1}

**Token 2 (Node: rule_0_terminal) → Action: valid_state_found**
Faits:
  - Status {id: S003, state: pending, priority: 2}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: valid_state_found(s.id, s.state)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 574.034µs

---

## Test 15: alpha_inequality_negative
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_inequality_negative.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_inequality_negative.facts
- **Temps:** 1.602907ms

### Règles Analysées
- `{o: Order} / NOT(o.status != "cancelled") ==> cancelled_order_found(o.id, o.total)`

### Faits Traités
- `Order(id:ORD001, total:100, status:pending)`
- `Order(id:ORD002, total:200, status:cancelled)`
- `Order(id:ORD003, total:300, status:completed)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 1

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: cancelled_order_found**
Faits:
  - Order {id: ORD002, total: 200, status: cancelled}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: cancelled_order_found**
Faits:
  - Order {id: ORD002, total: 200, status: cancelled}


#### Comparaison
- **Matches:** 1
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: cancelled_order_found(o.id, o.total)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 1 trouvés
- Analyse tokens: 1 attendus, 1 observés, 1 matches, 0 mismatches
- Test terminé avec succès en 1.602907ms

---

## Test 16: alpha_inequality_positive
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_inequality_positive.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_inequality_positive.facts
- **Temps:** 811.277µs

### Règles Analysées
- `{o: Order} / o.status != "cancelled" ==> valid_order_found(o.id, o.total)`

### Faits Traités
- `Order(id:ORD001, total:100, status:pending)`
- `Order(id:ORD002, total:200, status:cancelled)`
- `Order(id:ORD003, total:300, status:completed)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: valid_order_found**
Faits:
  - Order {id: ORD001, total: 100, status: pending}

**Token 2 (Node: rule_0_expected) → Action: valid_order_found**
Faits:
  - Order {total: 300, status: completed, id: ORD003}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: valid_order_found**
Faits:
  - Order {id: ORD001, total: 100, status: pending}

**Token 2 (Node: rule_0_terminal) → Action: valid_order_found**
Faits:
  - Order {id: ORD003, total: 300, status: completed}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: valid_order_found(o.id, o.total)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 811.277µs

---

## Test 17: alpha_length_negative
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_length_negative.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_length_negative.facts
- **Temps:** 508.461µs

### Règles Analysées
- `{p: Password} / NOT(LENGTH(p.value) >= 8) ==> weak_password_found(p.id, p.value)`

### Faits Traités
- `Password(id:"P001", value:"password123", secure:true)`
- `Password(id:"P002", value:"123", secure:false)`
- `Password(id:"P003", value:"pass", secure:false)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: weak_password_found**
Faits:
  - Password {secure: false, id: P002, value: 123}

**Token 2 (Node: rule_0_expected) → Action: weak_password_found**
Faits:
  - Password {secure: false, id: P003, value: pass}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: weak_password_found**
Faits:
  - Password {secure: false, id: P002, value: 123}

**Token 2 (Node: rule_0_terminal) → Action: weak_password_found**
Faits:
  - Password {id: P003, value: pass, secure: false}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: weak_password_found(p.id, p.value)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 508.461µs

---

## Test 18: alpha_length_positive
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_length_positive.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_length_positive.facts
- **Temps:** 470.029µs

### Règles Analysées
- `{p: Password} / LENGTH(p.value) >= 8 ==> secure_password_found(p.id, p.value)`

### Faits Traités
- `Password(id:"P001", value:"password123", secure:true)`
- `Password(id:"P002", value:"123", secure:false)`
- `Password(id:"P003", value:"verysecurepass", secure:true)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: secure_password_found**
Faits:
  - Password {id: P001, value: password123, secure: true}

**Token 2 (Node: rule_0_expected) → Action: secure_password_found**
Faits:
  - Password {id: P003, value: verysecurepass, secure: true}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: secure_password_found**
Faits:
  - Password {id: P001, value: password123, secure: true}

**Token 2 (Node: rule_0_terminal) → Action: secure_password_found**
Faits:
  - Password {id: P003, value: verysecurepass, secure: true}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: secure_password_found(p.id, p.value)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 470.029µs

---

## Test 19: alpha_like_negative
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_like_negative.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_like_negative.facts
- **Temps:** 568.363µs

### Règles Analysées
- `{e: Email} / NOT(e.address LIKE "%@company.com") ==> external_email_found(e.id, e.address)`

### Faits Traités
- `Email(id:"E001", address:"john@company.com", verified:true)`
- `Email(id:"E002", address:"jane@external.org", verified:false)`
- `Email(id:"E003", address:"user@other.net", verified:true)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: external_email_found**
Faits:
  - Email {id: E002, address: jane@external.org, verified: false}

**Token 2 (Node: rule_0_expected) → Action: external_email_found**
Faits:
  - Email {id: E003, address: user@other.net, verified: true}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: external_email_found**
Faits:
  - Email {verified: true, id: E003, address: user@other.net}

**Token 2 (Node: rule_0_terminal) → Action: external_email_found**
Faits:
  - Email {id: E002, address: jane@external.org, verified: false}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: external_email_found(e.id, e.address)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 568.363µs

---

## Test 20: alpha_like_positive
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_like_positive.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_like_positive.facts
- **Temps:** 542.294µs

### Règles Analysées
- `{e: Email} / e.address LIKE "%@company.com" ==> company_email_found(e.id, e.address)`

### Faits Traités
- `Email(id:"E001", address:"john@company.com", verified:true)`
- `Email(id:"E002", address:"jane@external.org", verified:false)`
- `Email(id:"E003", address:"admin@company.com", verified:true)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: company_email_found**
Faits:
  - Email {id: E001, address: john@company.com, verified: true}

**Token 2 (Node: rule_0_expected) → Action: company_email_found**
Faits:
  - Email {verified: true, id: E003, address: admin@company.com}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: company_email_found**
Faits:
  - Email {id: E001, address: john@company.com, verified: true}

**Token 2 (Node: rule_0_terminal) → Action: company_email_found**
Faits:
  - Email {verified: true, id: E003, address: admin@company.com}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: company_email_found(e.id, e.address)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 542.294µs

---

## Test 21: alpha_matches_negative
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_matches_negative.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_matches_negative.facts
- **Temps:** 519.171µs

### Règles Analysées
- `{c: Code} / NOT(c.value MATCHES "CODE[0-9]+") ==> invalid_code_found(c.id, c.value)`

### Faits Traités
- `Code(id:"C001", value:"CODE123", active:true)`
- `Code(id:"C002", value:"INVALID", active:false)`
- `Code(id:"C003", value:"BADFORMAT", active:true)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: invalid_code_found**
Faits:
  - Code {id: C002, value: INVALID, active: false}

**Token 2 (Node: rule_0_expected) → Action: invalid_code_found**
Faits:
  - Code {active: true, id: C003, value: BADFORMAT}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: invalid_code_found**
Faits:
  - Code {value: INVALID, active: false, id: C002}

**Token 2 (Node: rule_0_terminal) → Action: invalid_code_found**
Faits:
  - Code {active: true, id: C003, value: BADFORMAT}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: invalid_code_found(c.id, c.value)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 519.171µs

---

## Test 22: alpha_matches_positive
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_matches_positive.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_matches_positive.facts
- **Temps:** 476.1µs

### Règles Analysées
- `{c: Code} / c.value MATCHES "CODE[0-9]+" ==> valid_code_found(c.id, c.value)`

### Faits Traités
- `Code(id:"C001", value:"CODE123", active:true)`
- `Code(id:"C002", value:"INVALID", active:false)`
- `Code(id:"C003", value:"CODE999", active:true)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: valid_code_found**
Faits:
  - Code {id: C001, value: CODE123, active: true}

**Token 2 (Node: rule_0_expected) → Action: valid_code_found**
Faits:
  - Code {value: CODE999, active: true, id: C003}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: valid_code_found**
Faits:
  - Code {id: C001, value: CODE123, active: true}

**Token 2 (Node: rule_0_terminal) → Action: valid_code_found**
Faits:
  - Code {id: C003, value: CODE999, active: true}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: valid_code_found(c.id, c.value)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 476.1µs

---

## Test 23: alpha_string_negative
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_string_negative.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_string_negative.facts
- **Temps:** 438.52µs

### Règles Analysées
- `{u: User} / NOT(u.role == "admin") ==> non_admin_user_found(u.id, u.name)`

### Faits Traités
- `User(id:U001, name:Alice, role:admin)`
- `User(id:U002, name:Bob, role:user)`
- `User(id:U003, name:Charlie, role:admin)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 1

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: non_admin_user_found**
Faits:
  - User {name: Bob, role: user, id: U002}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: non_admin_user_found**
Faits:
  - User {id: U002, name: Bob, role: user}


#### Comparaison
- **Matches:** 1
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: non_admin_user_found(u.id, u.name)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 1 trouvés
- Analyse tokens: 1 attendus, 1 observés, 1 matches, 0 mismatches
- Test terminé avec succès en 438.52µs

---

## Test 24: alpha_string_positive
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_string_positive.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_string_positive.facts
- **Temps:** 425.105µs

### Règles Analysées
- `{u: User} / u.role == "admin" ==> admin_user_found(u.id, u.name)`

### Faits Traités
- `User(id:U001, name:Alice, role:admin)`
- `User(id:U002, name:Bob, role:user)`
- `User(id:U003, name:Charlie, role:admin)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: admin_user_found**
Faits:
  - User {id: U001, name: Alice, role: admin}

**Token 2 (Node: rule_0_expected) → Action: admin_user_found**
Faits:
  - User {id: U003, name: Charlie, role: admin}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: admin_user_found**
Faits:
  - User {role: admin, id: U001, name: Alice}

**Token 2 (Node: rule_0_terminal) → Action: admin_user_found**
Faits:
  - User {id: U003, name: Charlie, role: admin}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: admin_user_found(u.id, u.name)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 425.105µs

---

## Test 25: alpha_upper_negative
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_upper_negative.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_upper_negative.facts
- **Temps:** 468.145µs

### Règles Analysées
- `{d: Department} / NOT(UPPER(d.name) == "FINANCE") ==> non_finance_dept_found(d.id, d.name)`

### Faits Traités
- `Department(id:"D001", name:"finance", active:true)`
- `Department(id:"D002", name:"IT", active:true)`
- `Department(id:"D003", name:"HR", active:true)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: non_finance_dept_found**
Faits:
  - Department {active: true, id: D002, name: IT}

**Token 2 (Node: rule_0_expected) → Action: non_finance_dept_found**
Faits:
  - Department {id: D003, name: HR, active: true}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: non_finance_dept_found**
Faits:
  - Department {id: D002, name: IT, active: true}

**Token 2 (Node: rule_0_terminal) → Action: non_finance_dept_found**
Faits:
  - Department {name: HR, active: true, id: D003}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: non_finance_dept_found(d.id, d.name)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 468.145µs

---

## Test 26: alpha_upper_positive
- **Contraintes:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_upper_positive.constraint
- **Faits:** /home/resinsec/dev/tsd/test/coverage/alpha/alpha_upper_positive.facts
- **Temps:** 435.364µs

### Règles Analysées
- `{d: Department} / UPPER(d.name) == "FINANCE" ==> finance_dept_found(d.id, d.name)`

### Faits Traités
- `Department(id:"D001", name:"finance", active:true)`
- `Department(id:"D002", name:"IT", active:true)`
- `Department(id:"D003", name:"Finance", active:true)`

### État du Réseau
- **Total nœuds:** 4
- **Nœuds Alpha:** 1
- **Nœuds Terminal:** 1
- **Tokens générés:** 2

### Analyse Détaillée des Tokens

#### Tokens Attendus
**Token 1 (Node: rule_0_expected) → Action: finance_dept_found**
Faits:
  - Department {id: D001, name: finance, active: true}

**Token 2 (Node: rule_0_expected) → Action: finance_dept_found**
Faits:
  - Department {active: true, id: D003, name: Finance}


#### Tokens Observés
**Token 1 (Node: rule_0_terminal) → Action: finance_dept_found**
Faits:
  - Department {id: D003, name: Finance, active: true}

**Token 2 (Node: rule_0_terminal) → Action: finance_dept_found**
Faits:
  - Department {id: D001, name: finance, active: true}


#### Comparaison
- **Matches:** 2
- **Mismatches:** 0
- **Cohérence sémantique:** ✅ VALIDÉE

### Résultats
- **Attendu:** Action attendue: finance_dept_found(d.id, d.name)
- **Observé:** Réseau construit: 4 nœuds, 1 actions possibles
- **Status:** ✅ Succès

### Log d'Inférence
- Début analyse: 1 règles, 3 faits
- Pipeline construit avec succès
- Analyse de l'état du réseau...
- Nœuds créés: 4 total (1 alpha, 1 terminal, 1 type, 0 beta)
- Tokens analysés: 2 trouvés
- Analyse tokens: 2 attendus, 2 observés, 2 matches, 0 mismatches
- Test terminé avec succès en 435.364µs

---
