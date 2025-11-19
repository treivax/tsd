# RAPPORT DÉTAILLÉ - TESTS BETA NODES

Date: 2025-11-19 23:45:28

## Test: complex_not_exists_combination

### Fichiers de Test
- **Contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/complex_not_exists_combination.constraint`
- **Faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/complex_not_exists_combination.facts`

**Statut: ✅ VALIDÉE**

- **Type:** Jointure (Student, Course)
- **Temps d'exécution:** 829.04µs
- **Tokens attendus:** 0
- **Tokens observés:** 0
- **Correspondances:** 0
- **Mismatches:** 0
- **Taux de succès:** 0.0%

### Règles de Contraintes
```
1. // Test combinaison NOT + EXISTS avec jointures
2. type Student : <id: string, name: string, grade: number, active: bool>
3. type Course : <id: string, title: string, level: string, credits: number>
4. type Enrollment : <id: string, student_id: string, course_id: string, status: string>
5. {s: Student, c: Course} / s.active == true AND s.grade >= 60 AND NOT (s.grade < 60) AND EXISTS (e: Enrollment / e.student_id == s.id AND e.course_id == c.id AND e.status == "enrolled") ==> qualified_enrolled_student(s.id, c.id)
6. {s: Student, c: Course} / c.level == "advanced" AND c.credits >= 3 AND NOT (c.credits < 3) AND EXISTS (e: Enrollment / e.student_id == s.id AND e.course_id == c.id AND e.status != "dropped") ==> advanced_course_student(s.id, c.id)
7. {s: Student} / s.active == true AND s.grade > 0 AND NOT (EXISTS (e: Enrollment / e.student_id == s.id AND e.status == "failed")) ==> successful_student(s.id)
```

### Faits Injectés
```
1. Student(id:S001, name:Alice, grade:85, active:true)
2. Student(id:S002, name:Bob, grade:45, active:true)
3. Student(id:S003, name:Charlie, grade:92, active:false)
4. Student(id:S004, name:Diana, grade:78, active:true)
5. Course(id:C001, title:"Advanced Math", level:advanced, credits:4)
6. Course(id:C002, title:"Basic Physics", level:beginner, credits:2)
7. Course(id:C003, title:"Computer Science", level:advanced, credits:3)
8. Enrollment(id:EN001, student_id:S001, course_id:C001, status:enrolled)
9. Enrollment(id:EN002, student_id:S002, course_id:C002, status:failed)
10. Enrollment(id:EN003, student_id:S003, course_id:C001, status:dropped)
11. Enrollment(id:EN004, student_id:S004, course_id:C003, status:enrolled)
```

---

## Test: exists_complex_operator

### Fichiers de Test
- **Contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/exists_complex_operator.constraint`
- **Faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/exists_complex_operator.facts`

**Statut: ✅ VALIDÉE**

- **Type:** Jointure (Customer, Order)
- **Temps d'exécution:** 596.105µs
- **Tokens attendus:** 0
- **Tokens observés:** 0
- **Correspondances:** 0
- **Mismatches:** 0
- **Taux de succès:** 0.0%

### Règles de Contraintes
```
1. // Test opérateur EXISTS complexe dans nœuds beta
2. type Customer : <id: string, name: string, segment: string, country: string>
3. type Order : <id: string, customer_id: string, amount: number, date: string>
4. type Payment : <id: string, order_id: string, method: string, status: string>
5. {c: Customer, o: Order} / c.id == o.customer_id AND o.amount > 0 AND EXISTS (p: Payment / p.order_id == o.id AND p.status == "completed") ==> paid_customer_order(c.id, o.id)
6. {c: Customer, o: Order} / c.id == o.customer_id AND o.amount >= 100 AND EXISTS (p: Payment / p.order_id == o.id AND p.method == "credit_card") ==> credit_customer_order(c.id, o.id)
7. {c: Customer} / c.segment == "premium" AND c.country != "" AND EXISTS (o: Order / o.customer_id == c.id AND o.amount > 1000) ==> premium_high_value_customer(c.id)
```

### Faits Injectés
```
1. Customer(id:C001, name:Alice, segment:premium, country:USA)
2. Customer(id:C002, name:Bob, segment:standard, country:France)
3. Customer(id:C003, name:Charlie, segment:premium, country:Germany)
4. Order(id:O001, customer_id:C001, amount:1500, date:"2024-01-01")
5. Order(id:O002, customer_id:C002, amount:150, date:"2024-01-02")
6. Order(id:O003, customer_id:C003, amount:800, date:"2024-01-03")
7. Order(id:O004, customer_id:C001, amount:2000, date:"2024-01-04")
8. Payment(id:PY001, order_id:O001, method:"credit_card", status:completed)
9. Payment(id:PY002, order_id:O002, method:paypal, status:pending)
10. Payment(id:PY003, order_id:O003, method:"credit_card", status:completed)
11. Payment(id:PY004, order_id:O004, method:"bank_transfer", status:completed)
```

---

## Test: exists_simple

### Fichiers de Test
- **Contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/exists_simple.constraint`
- **Faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/exists_simple.facts`

**Statut: ✅ VALIDÉE**

- **Type:** Jointure (Person)
- **Temps d'exécution:** 61.475µs
- **Tokens attendus:** 0
- **Tokens observés:** 0
- **Correspondances:** 0
- **Mismatches:** 0
- **Taux de succès:** 0.0%

### Règles de Contraintes
```
1. // Test existence simple
2. type Person : <id: string, name: string>
3. type Order : <customer_id: string, amount: number>
4. {p: Person} / EXISTS (o: Order / o.customer_id == p.id) ==> person_has_orders(p.id)
```

### Faits Injectés
```
1. Person(id:P001, name:Alice)
2. Person(id:P002, name:Bob)
3. Order(customer_id:P001, amount:100)
```

---

## Test: join_and_operator

### Fichiers de Test
- **Contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_and_operator.constraint`
- **Faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_and_operator.facts`

**Statut: ✅ VALIDÉE**

- **Type:** Jointure (Person, Order)
- **Temps d'exécution:** 1.16074ms
- **Tokens attendus:** 11
- **Tokens observés:** 11
- **Correspondances:** 11
- **Mismatches:** 0
- **Taux de succès:** 100.0%

### Règles de Contraintes
```
1. // Test opérateur AND dans jointures beta
2. type Person : <id: string, name: string, age: number, status: string>
3. type Order : <id: string, customer_id: string, amount: number, status: string>
4. {p: Person, o: Order} / p.id == o.customer_id AND o.amount > 100 AND p.status == "active" ==> high_value_order(p.id, o.id)
5. {p: Person, o: Order} / p.id == o.customer_id AND p.age >= 18 AND o.status == "confirmed" AND p.status == "active" ==> adult_confirmed_order(p.id, o.id)
6. {p: Person, o: Order} / p.id == o.customer_id AND o.amount >= 50 AND o.amount <= 500 AND p.age > 21 ==> medium_value_order(p.id, o.id)
```

### Faits Injectés
```
1. Person(id:P001, name:Alice, age:25, status:active)
2. Person(id:P002, name:Bob, age:30, status:inactive)
3. Person(id:P003, name:Charlie, age:16, status:active)
4. Person(id:P004, name:Diana, age:22, status:active)
5. Order(id:O001, customer_id:P001, amount:150, status:confirmed)
6. Order(id:O002, customer_id:P002, amount:75, status:pending)
7. Order(id:O003, customer_id:P001, amount:200, status:confirmed)
8. Order(id:O004, customer_id:P003, amount:300, status:confirmed)
9. Order(id:O005, customer_id:P004, amount:125, status:confirmed)
```

### Tokens Attendus:
1. `Order(amount:150,customer_id:P001,id:O001,status:confirmed)+Person(age:25,id:P001,name:Alice,status:active)`
2. `Order(amount:200,customer_id:P001,id:O003,status:confirmed)+Person(age:25,id:P001,name:Alice,status:active)`
3. `Order(amount:300,customer_id:P003,id:O004,status:confirmed)+Person(age:16,id:P003,name:Charlie,status:active)`
4. `Order(amount:125,customer_id:P004,id:O005,status:confirmed)+Person(age:22,id:P004,name:Diana,status:active)`
5. `Order(amount:150,customer_id:P001,id:O001,status:confirmed)+Person(age:25,id:P001,name:Alice,status:active)`
6. `Order(amount:200,customer_id:P001,id:O003,status:confirmed)+Person(age:25,id:P001,name:Alice,status:active)`
7. `Order(amount:125,customer_id:P004,id:O005,status:confirmed)+Person(age:22,id:P004,name:Diana,status:active)`
8. `Order(amount:150,customer_id:P001,id:O001,status:confirmed)+Person(age:25,id:P001,name:Alice,status:active)`
9. `Order(amount:200,customer_id:P001,id:O003,status:confirmed)+Person(age:25,id:P001,name:Alice,status:active)`
10. `Order(amount:75,customer_id:P002,id:O002,status:pending)+Person(age:30,id:P002,name:Bob,status:inactive)`
11. `Order(amount:125,customer_id:P004,id:O005,status:confirmed)+Person(age:22,id:P004,name:Diana,status:active)`

### Tokens Observés:
1. `Order(amount:150,customer_id:P001,id:O001,status:confirmed)+Person(age:25,id:P001,name:Alice,status:active)`
2. `Order(amount:200,customer_id:P001,id:O003,status:confirmed)+Person(age:25,id:P001,name:Alice,status:active)`
3. `Order(amount:300,customer_id:P003,id:O004,status:confirmed)+Person(age:16,id:P003,name:Charlie,status:active)`
4. `Order(amount:125,customer_id:P004,id:O005,status:confirmed)+Person(age:22,id:P004,name:Diana,status:active)`
5. `Order(amount:150,customer_id:P001,id:O001,status:confirmed)+Person(age:25,id:P001,name:Alice,status:active)`
6. `Order(amount:200,customer_id:P001,id:O003,status:confirmed)+Person(age:25,id:P001,name:Alice,status:active)`
7. `Order(amount:125,customer_id:P004,id:O005,status:confirmed)+Person(age:22,id:P004,name:Diana,status:active)`
8. `Order(amount:150,customer_id:P001,id:O001,status:confirmed)+Person(age:25,id:P001,name:Alice,status:active)`
9. `Order(amount:200,customer_id:P001,id:O003,status:confirmed)+Person(age:25,id:P001,name:Alice,status:active)`
10. `Order(amount:75,customer_id:P002,id:O002,status:pending)+Person(age:30,id:P002,name:Bob,status:inactive)`
11. `Order(amount:125,customer_id:P004,id:O005,status:confirmed)+Person(age:22,id:P004,name:Diana,status:active)`

---

## Test: join_arithmetic_operators

### Fichiers de Test
- **Contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_arithmetic_operators.constraint`
- **Faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_arithmetic_operators.facts`

**Statut: ✅ VALIDÉE**

- **Type:** Jointure (Account, Transaction)
- **Temps d'exécution:** 728.893µs
- **Tokens attendus:** 8
- **Tokens observés:** 8
- **Correspondances:** 8
- **Mismatches:** 0
- **Taux de succès:** 100.0%

### Règles de Contraintes
```
1. // Test opérateurs arithmétiques dans jointures beta
2. type Account : <id: string, balance: number, credit_limit: number, fees: number>
3. type Transaction : <id: string, account_id: string, amount: number, type: string>
4. {a: Account, t: Transaction} / a.id == t.account_id AND t.amount + a.fees <= a.balance ==> valid_transaction(a.id, t.id)
5. {a: Account, t: Transaction} / a.id == t.account_id AND a.balance + a.credit_limit >= t.amount * 2 ==> safe_transaction(a.id, t.id)
6. {a: Account, t: Transaction} / a.id == t.account_id AND (a.balance - t.amount) / a.credit_limit > 0.1 ==> conservative_transaction(a.id, t.id)
```

### Faits Injectés
```
1. Account(id:ACC001, balance:1000, credit_limit:500, fees:10)
2. Account(id:ACC002, balance:2500, credit_limit:1000, fees:5)
3. Account(id:ACC003, balance:500, credit_limit:200, fees:15)
4. Transaction(id:TXN001, account_id:ACC001, amount:900, type:debit)
5. Transaction(id:TXN002, account_id:ACC002, amount:1200, type:debit)
6. Transaction(id:TXN003, account_id:ACC003, amount:600, type:debit)
7. Transaction(id:TXN004, account_id:ACC001, amount:50, type:debit)
```

### Tokens Attendus:
1. `Account(balance:1000,credit_limit:500,fees:10,id:ACC001)+Transaction(account_id:ACC001,amount:900,id:TXN001,type:debit)`
2. `Account(balance:1000,credit_limit:500,fees:10,id:ACC001)+Transaction(account_id:ACC001,amount:50,id:TXN004,type:debit)`
3. `Account(balance:2500,credit_limit:1000,fees:5,id:ACC002)+Transaction(account_id:ACC002,amount:1200,id:TXN002,type:debit)`
4. `Account(balance:500,credit_limit:200,fees:15,id:ACC003)+Transaction(account_id:ACC003,amount:600,id:TXN003,type:debit)`
5. `Account(balance:1000,credit_limit:500,fees:10,id:ACC001)+Transaction(account_id:ACC001,amount:900,id:TXN001,type:debit)`
6. `Account(balance:1000,credit_limit:500,fees:10,id:ACC001)+Transaction(account_id:ACC001,amount:50,id:TXN004,type:debit)`
7. `Account(balance:2500,credit_limit:1000,fees:5,id:ACC002)+Transaction(account_id:ACC002,amount:1200,id:TXN002,type:debit)`
8. `Account(balance:500,credit_limit:200,fees:15,id:ACC003)+Transaction(account_id:ACC003,amount:600,id:TXN003,type:debit)`

### Tokens Observés:
1. `Account(balance:1000,credit_limit:500,fees:10,id:ACC001)+Transaction(account_id:ACC001,amount:900,id:TXN001,type:debit)`
2. `Account(balance:1000,credit_limit:500,fees:10,id:ACC001)+Transaction(account_id:ACC001,amount:50,id:TXN004,type:debit)`
3. `Account(balance:2500,credit_limit:1000,fees:5,id:ACC002)+Transaction(account_id:ACC002,amount:1200,id:TXN002,type:debit)`
4. `Account(balance:500,credit_limit:200,fees:15,id:ACC003)+Transaction(account_id:ACC003,amount:600,id:TXN003,type:debit)`
5. `Account(balance:1000,credit_limit:500,fees:10,id:ACC001)+Transaction(account_id:ACC001,amount:900,id:TXN001,type:debit)`
6. `Account(balance:1000,credit_limit:500,fees:10,id:ACC001)+Transaction(account_id:ACC001,amount:50,id:TXN004,type:debit)`
7. `Account(balance:2500,credit_limit:1000,fees:5,id:ACC002)+Transaction(account_id:ACC002,amount:1200,id:TXN002,type:debit)`
8. `Account(balance:500,credit_limit:200,fees:15,id:ACC003)+Transaction(account_id:ACC003,amount:600,id:TXN003,type:debit)`

---

## Test: join_comparison_operators

### Fichiers de Test
- **Contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_comparison_operators.constraint`
- **Faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_comparison_operators.facts`

**Statut: ✅ VALIDÉE**

- **Type:** Jointure (User, Activity)
- **Temps d'exécution:** 958.702µs
- **Tokens attendus:** 8
- **Tokens observés:** 8
- **Correspondances:** 8
- **Mismatches:** 0
- **Taux de succès:** 100.0%

### Règles de Contraintes
```
1. // Test opérateurs de comparaison dans jointures beta
2. type User : <id: string, name: string, score: number, created: number>
3. type Activity : <id: string, user_id: string, points: number, timestamp: number>
4. {u: User, a: Activity} / u.id == a.user_id AND a.points > u.score ==> improvement_activity(u.id, a.id)
5. {u: User, a: Activity} / u.id == a.user_id AND a.timestamp >= u.created ==> valid_activity(u.id, a.id)
6. {u: User, a: Activity} / u.id == a.user_id AND a.points <= 50 ==> low_activity(u.id, a.id)
7. {u: User, a: Activity} / u.id == a.user_id AND a.points != u.score ==> different_score_activity(u.id, a.id)
```

### Faits Injectés
```
1. User(id:U001, name:Alice, score:75, created:1700000000)
2. User(id:U002, name:Bob, score:60, created:1700100000)
3. User(id:U003, name:Carol, score:90, created:1700200000)
4. Activity(id:A001, user_id:U001, points:80, timestamp:1700000100)
5. Activity(id:A002, user_id:U002, points:45, timestamp:1700050000)
6. Activity(id:A003, user_id:U001, points:75, timestamp:1699900000)
7. Activity(id:A004, user_id:U003, points:95, timestamp:1700300000)
```

### Tokens Attendus:
1. `Activity(id:A001,points:80,timestamp:1700000100,user_id:U001)+User(created:1700000000,id:U001,name:Alice,score:75)`
2. `Activity(id:A004,points:95,timestamp:1700300000,user_id:U003)+User(created:1700200000,id:U003,name:Carol,score:90)`
3. `Activity(id:A001,points:80,timestamp:1700000100,user_id:U001)+User(created:1700000000,id:U001,name:Alice,score:75)`
4. `Activity(id:A004,points:95,timestamp:1700300000,user_id:U003)+User(created:1700200000,id:U003,name:Carol,score:90)`
5. `Activity(id:A002,points:45,timestamp:1700050000,user_id:U002)+User(created:1700100000,id:U002,name:Bob,score:60)`
6. `Activity(id:A001,points:80,timestamp:1700000100,user_id:U001)+User(created:1700000000,id:U001,name:Alice,score:75)`
7. `Activity(id:A002,points:45,timestamp:1700050000,user_id:U002)+User(created:1700100000,id:U002,name:Bob,score:60)`
8. `Activity(id:A004,points:95,timestamp:1700300000,user_id:U003)+User(created:1700200000,id:U003,name:Carol,score:90)`

### Tokens Observés:
1. `Activity(id:A001,points:80,timestamp:1700000100,user_id:U001)+User(created:1700000000,id:U001,name:Alice,score:75)`
2. `Activity(id:A004,points:95,timestamp:1700300000,user_id:U003)+User(created:1700200000,id:U003,name:Carol,score:90)`
3. `Activity(id:A001,points:80,timestamp:1700000100,user_id:U001)+User(created:1700000000,id:U001,name:Alice,score:75)`
4. `Activity(id:A004,points:95,timestamp:1700300000,user_id:U003)+User(created:1700200000,id:U003,name:Carol,score:90)`
5. `Activity(id:A002,points:45,timestamp:1700050000,user_id:U002)+User(created:1700100000,id:U002,name:Bob,score:60)`
6. `Activity(id:A001,points:80,timestamp:1700000100,user_id:U001)+User(created:1700000000,id:U001,name:Alice,score:75)`
7. `Activity(id:A002,points:45,timestamp:1700050000,user_id:U002)+User(created:1700100000,id:U002,name:Bob,score:60)`
8. `Activity(id:A004,points:95,timestamp:1700300000,user_id:U003)+User(created:1700200000,id:U003,name:Carol,score:90)`

---

## Test: join_in_contains_operators

### Fichiers de Test
- **Contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_in_contains_operators.constraint`
- **Faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_in_contains_operators.facts`

**Statut: ✅ VALIDÉE**

- **Type:** Jointure (Product, Review)
- **Temps d'exécution:** 703.415µs
- **Tokens attendus:** 9
- **Tokens observés:** 9
- **Correspondances:** 9
- **Mismatches:** 0
- **Taux de succès:** 100.0%

### Règles de Contraintes
```
1. // Test opérateurs IN et CONTAINS dans jointures beta
2. type Product : <id: string, name: string, categories: string, keywords: string>
3. type Review : <id: string, product_id: string, rating: number, status: string>
4. {p: Product, r: Review} / p.id == r.product_id AND r.status IN ["approved", "verified"] ==> approved_review(p.id, r.id)
5. {p: Product, r: Review} / p.id == r.product_id AND p.keywords CONTAINS "premium" ==> premium_product_review(p.id, r.id)
6. {p: Product, r: Review} / p.id == r.product_id AND p.categories IN ["electronics", "tech"] AND r.rating >= 4 ==> tech_high_rating(p.id, r.id)
```

### Faits Injectés
```
1. Product(id:PROD001, name:Laptop, categories:electronics, keywords:"premium computer high-end")
2. Product(id:PROD002, name:Mouse, categories:accessories, keywords:"basic cheap")
3. Product(id:PROD003, name:Phone, categories:tech, keywords:"premium mobile smartphone")
4. Review(id:R001, product_id:PROD001, rating:5, status:approved)
5. Review(id:R002, product_id:PROD002, rating:3, status:pending)
6. Review(id:R003, product_id:PROD001, rating:4, status:verified)
7. Review(id:R004, product_id:PROD003, rating:5, status:approved)
```

### Tokens Attendus:
1. `Product(categories:electronics,id:PROD001,keywords:"premium computer high-end",name:Laptop)+Review(id:R001,product_id:PROD001,rating:5,status:approved)`
2. `Product(categories:electronics,id:PROD001,keywords:"premium computer high-end",name:Laptop)+Review(id:R003,product_id:PROD001,rating:4,status:verified)`
3. `Product(categories:tech,id:PROD003,keywords:"premium mobile smartphone",name:Phone)+Review(id:R004,product_id:PROD003,rating:5,status:approved)`
4. `Product(categories:electronics,id:PROD001,keywords:"premium computer high-end",name:Laptop)+Review(id:R001,product_id:PROD001,rating:5,status:approved)`
5. `Product(categories:electronics,id:PROD001,keywords:"premium computer high-end",name:Laptop)+Review(id:R003,product_id:PROD001,rating:4,status:verified)`
6. `Product(categories:tech,id:PROD003,keywords:"premium mobile smartphone",name:Phone)+Review(id:R004,product_id:PROD003,rating:5,status:approved)`
7. `Product(categories:electronics,id:PROD001,keywords:"premium computer high-end",name:Laptop)+Review(id:R001,product_id:PROD001,rating:5,status:approved)`
8. `Product(categories:electronics,id:PROD001,keywords:"premium computer high-end",name:Laptop)+Review(id:R003,product_id:PROD001,rating:4,status:verified)`
9. `Product(categories:tech,id:PROD003,keywords:"premium mobile smartphone",name:Phone)+Review(id:R004,product_id:PROD003,rating:5,status:approved)`

### Tokens Observés:
1. `Product(categories:electronics,id:PROD001,keywords:"premium computer high-end",name:Laptop)+Review(id:R001,product_id:PROD001,rating:5,status:approved)`
2. `Product(categories:electronics,id:PROD001,keywords:"premium computer high-end",name:Laptop)+Review(id:R003,product_id:PROD001,rating:4,status:verified)`
3. `Product(categories:tech,id:PROD003,keywords:"premium mobile smartphone",name:Phone)+Review(id:R004,product_id:PROD003,rating:5,status:approved)`
4. `Product(categories:electronics,id:PROD001,keywords:"premium computer high-end",name:Laptop)+Review(id:R001,product_id:PROD001,rating:5,status:approved)`
5. `Product(categories:electronics,id:PROD001,keywords:"premium computer high-end",name:Laptop)+Review(id:R003,product_id:PROD001,rating:4,status:verified)`
6. `Product(categories:tech,id:PROD003,keywords:"premium mobile smartphone",name:Phone)+Review(id:R004,product_id:PROD003,rating:5,status:approved)`
7. `Product(categories:electronics,id:PROD001,keywords:"premium computer high-end",name:Laptop)+Review(id:R001,product_id:PROD001,rating:5,status:approved)`
8. `Product(categories:electronics,id:PROD001,keywords:"premium computer high-end",name:Laptop)+Review(id:R003,product_id:PROD001,rating:4,status:verified)`
9. `Product(categories:tech,id:PROD003,keywords:"premium mobile smartphone",name:Phone)+Review(id:R004,product_id:PROD003,rating:5,status:approved)`

---

## Test: join_multi_variable_complex

### Fichiers de Test
- **Contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_multi_variable_complex.constraint`
- **Faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_multi_variable_complex.facts`

**Statut: ✅ VALIDÉE**

- **Type:** Jointure (User, Team, Task)
- **Temps d'exécution:** 434.112µs
- **Tokens attendus:** 0
- **Tokens observés:** 0
- **Correspondances:** 0
- **Mismatches:** 0
- **Taux de succès:** 0.0%

### Règles de Contraintes
```
1. // Test jointures complexes multi-variables
2. type User : <id: string, name: string, role: string, team_id: string>
3. type Team : <id: string, name: string, budget: number, manager_id: string>
4. type Task : <id: string, assignee_id: string, team_id: string, priority: string, effort: number>
5. {u: User, t: Team, task: Task} / u.id == t.manager_id AND t.id == task.team_id AND task.priority == "high" ==> manager_high_priority_task(u.id, t.id, task.id)
6. {u: User, t: Team, task: Task} / u.team_id == t.id AND u.id == task.assignee_id AND t.budget > task.effort * 100 ==> affordable_task_assignment(u.id, t.id, task.id)
7. {u: User, t: Team, task: Task} / u.role == "lead" AND u.team_id == t.id AND task.team_id == t.id AND task.effort >= 40 ==> lead_complex_task(u.id, t.id, task.id)
```

### Faits Injectés
```
1. User(id:U001, name:Alice, role:manager, team_id:T001)
2. User(id:U002, name:Bob, role:lead, team_id:T001)
3. User(id:U003, name:Carol, role:developer, team_id:T002)
4. Team(id:T001, name:Alpha, budget:10000, manager_id:U001)
5. Team(id:T002, name:Beta, budget:5000, manager_id:U003)
6. Task(id:TASK001, assignee_id:U002, team_id:T001, priority:high, effort:50)
7. Task(id:TASK002, assignee_id:U003, team_id:T002, priority:medium, effort:20)
8. Task(id:TASK003, assignee_id:U001, team_id:T001, priority:high, effort:30)
```

---

## Test: join_or_operator

### Fichiers de Test
- **Contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_or_operator.constraint`
- **Faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_or_operator.facts`

**Statut: ✅ VALIDÉE**

- **Type:** Jointure (Person, Order)
- **Temps d'exécution:** 955.066µs
- **Tokens attendus:** 7
- **Tokens observés:** 7
- **Correspondances:** 7
- **Mismatches:** 0
- **Taux de succès:** 100.0%

### Règles de Contraintes
```
1. // Test opérateur OR dans jointures beta
2. type Person : <id: string, name: string, age: number, priority: string>
3. type Order : <id: string, customer_id: string, amount: number, urgent: bool>
4. {p: Person, o: Order} / p.id == o.customer_id AND (o.amount > 500 OR o.urgent == true) AND p.priority != "" ==> priority_order(p.id, o.id)
5. {p: Person, o: Order} / p.id == o.customer_id AND (p.priority == "high" OR p.age >= 65) AND o.amount > 0 ==> special_customer_order(p.id, o.id)
6. {p: Person, o: Order} / p.id == o.customer_id AND (o.urgent == false OR o.amount < 100) AND p.age >= 18 ==> standard_order(p.id, o.id)
```

### Faits Injectés
```
1. Person(id:P001, name:Alice, age:25, priority:high)
2. Person(id:P002, name:Bob, age:67, priority:normal)
3. Person(id:P003, name:Charlie, age:35, priority:low)
4. Person(id:P004, name:Diana, age:19, priority:normal)
5. Order(id:O001, customer_id:P001, amount:600, urgent:false)
6. Order(id:O002, customer_id:P002, amount:300, urgent:false)
7. Order(id:O003, customer_id:P003, amount:100, urgent:true)
8. Order(id:O004, customer_id:P004, amount:50, urgent:false)
```

### Tokens Attendus:
1. `Order(amount:600,customer_id:P001,id:O001,urgent:false)+Person(age:25,id:P001,name:Alice,priority:high)`
2. `Order(amount:100,customer_id:P003,id:O003,urgent:true)+Person(age:35,id:P003,name:Charlie,priority:low)`
3. `Order(amount:600,customer_id:P001,id:O001,urgent:false)+Person(age:25,id:P001,name:Alice,priority:high)`
4. `Order(amount:300,customer_id:P002,id:O002,urgent:false)+Person(age:67,id:P002,name:Bob,priority:normal)`
5. `Order(amount:600,customer_id:P001,id:O001,urgent:false)+Person(age:25,id:P001,name:Alice,priority:high)`
6. `Order(amount:300,customer_id:P002,id:O002,urgent:false)+Person(age:67,id:P002,name:Bob,priority:normal)`
7. `Order(amount:50,customer_id:P004,id:O004,urgent:false)+Person(age:19,id:P004,name:Diana,priority:normal)`

### Tokens Observés:
1. `Order(amount:600,customer_id:P001,id:O001,urgent:false)+Person(age:25,id:P001,name:Alice,priority:high)`
2. `Order(amount:100,customer_id:P003,id:O003,urgent:true)+Person(age:35,id:P003,name:Charlie,priority:low)`
3. `Order(amount:600,customer_id:P001,id:O001,urgent:false)+Person(age:25,id:P001,name:Alice,priority:high)`
4. `Order(amount:300,customer_id:P002,id:O002,urgent:false)+Person(age:67,id:P002,name:Bob,priority:normal)`
5. `Order(amount:600,customer_id:P001,id:O001,urgent:false)+Person(age:25,id:P001,name:Alice,priority:high)`
6. `Order(amount:300,customer_id:P002,id:O002,urgent:false)+Person(age:67,id:P002,name:Bob,priority:normal)`
7. `Order(amount:50,customer_id:P004,id:O004,urgent:false)+Person(age:19,id:P004,name:Diana,priority:normal)`

---

## Test: join_simple

### Fichiers de Test
- **Contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_simple.constraint`
- **Faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_simple.facts`

**Statut: ✅ VALIDÉE**

- **Type:** Jointure (Person, Order)
- **Temps d'exécution:** 129.311µs
- **Tokens attendus:** 2
- **Tokens observés:** 2
- **Correspondances:** 2
- **Mismatches:** 0
- **Taux de succès:** 100.0%

### Règles de Contraintes
```
1. // Test jointure simple entre deux faits
2. type Person : <id: string, name: string, age: number>
3. type Order : <id: string, customer_id: string, amount: number>
4. {p: Person, o: Order} / p.id == o.customer_id ==> customer_order_match(p.id, o.id)
```

### Faits Injectés
```
1. Person(id:P001, name:Alice, age:25)
2. Person(id:P002, name:Bob, age:30)
3. Order(id:O001, customer_id:P001, amount:100)
4. Order(id:O002, customer_id:P002, amount:200)
```

### Tokens Attendus:
1. `Order(amount:100,customer_id:P001,id:O001)+Person(age:25,id:P001,name:Alice)`
2. `Order(amount:200,customer_id:P002,id:O002)+Person(age:30,id:P002,name:Bob)`

### Tokens Observés:
1. `Order(amount:100,customer_id:P001,id:O001)+Person(age:25,id:P001,name:Alice)`
2. `Order(amount:200,customer_id:P002,id:O002)+Person(age:30,id:P002,name:Bob)`

---

## Test: not_complex_operator

### Fichiers de Test
- **Contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/not_complex_operator.constraint`
- **Faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/not_complex_operator.facts`

**Statut: ✅ VALIDÉE**

- **Type:** Jointure (Employee, Project)
- **Temps d'exécution:** 929.698µs
- **Tokens attendus:** 6
- **Tokens observés:** 6
- **Correspondances:** 6
- **Mismatches:** 0
- **Taux de succès:** 100.0%

### Règles de Contraintes
```
1. // Test opérateur NOT complexe dans nœuds beta
2. type Employee : <id: string, name: string, department: string, active: bool>
3. type Project : <id: string, lead_id: string, status: string, budget: number>
4. {e: Employee, p: Project} / e.id == p.lead_id AND e.active == true AND NOT (p.status == "cancelled") ==> active_project_lead(e.id, p.id)
5. {e: Employee, p: Project} / e.id == p.lead_id AND NOT (e.active == false OR p.budget < 1000) AND p.status == "active" ==> qualified_project_lead(e.id, p.id)
6. {e: Employee, p: Project} / e.id == p.lead_id AND e.active == true AND NOT (e.department IN ["temp", "intern"]) ==> permanent_project_lead(e.id, p.id)
```

### Faits Injectés
```
1. Employee(id:E001, name:Alice, department:engineering, active:true)
2. Employee(id:E002, name:Bob, department:temp, active:true)
3. Employee(id:E003, name:Charlie, department:marketing, active:false)
4. Employee(id:E004, name:Diana, department:sales, active:true)
5. Project(id:P001, lead_id:E001, status:active, budget:5000)
6. Project(id:P002, lead_id:E002, status:cancelled, budget:2000)
7. Project(id:P003, lead_id:E003, status:completed, budget:800)
8. Project(id:P004, lead_id:E004, status:active, budget:1500)
```

### Tokens Attendus:
1. `Employee(active:true,department:engineering,id:E001,name:Alice)+Project(budget:5000,id:P001,lead_id:E001,status:active)`
2. `Employee(active:true,department:sales,id:E004,name:Diana)+Project(budget:1500,id:P004,lead_id:E004,status:active)`
3. `Employee(active:true,department:engineering,id:E001,name:Alice)+Project(budget:5000,id:P001,lead_id:E001,status:active)`
4. `Employee(active:true,department:sales,id:E004,name:Diana)+Project(budget:1500,id:P004,lead_id:E004,status:active)`
5. `Employee(active:true,department:engineering,id:E001,name:Alice)+Project(budget:5000,id:P001,lead_id:E001,status:active)`
6. `Employee(active:true,department:sales,id:E004,name:Diana)+Project(budget:1500,id:P004,lead_id:E004,status:active)`

### Tokens Observés:
1. `Employee(active:true,department:engineering,id:E001,name:Alice)+Project(budget:5000,id:P001,lead_id:E001,status:active)`
2. `Employee(active:true,department:sales,id:E004,name:Diana)+Project(budget:1500,id:P004,lead_id:E004,status:active)`
3. `Employee(active:true,department:engineering,id:E001,name:Alice)+Project(budget:5000,id:P001,lead_id:E001,status:active)`
4. `Employee(active:true,department:sales,id:E004,name:Diana)+Project(budget:1500,id:P004,lead_id:E004,status:active)`
5. `Employee(active:true,department:engineering,id:E001,name:Alice)+Project(budget:5000,id:P001,lead_id:E001,status:active)`
6. `Employee(active:true,department:sales,id:E004,name:Diana)+Project(budget:1500,id:P004,lead_id:E004,status:active)`

---

## Test: not_simple

### Fichiers de Test
- **Contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/not_simple.constraint`
- **Faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/not_simple.facts`

**Statut: ✅ VALIDÉE**

- **Type:** Jointure (Person)
- **Temps d'exécution:** 61.204µs
- **Tokens attendus:** 0
- **Tokens observés:** 0
- **Correspondances:** 0
- **Mismatches:** 0
- **Taux de succès:** 0.0%

### Règles de Contraintes
```
1. // Test négation simple
2. type Person : <id: string, name: string, active: bool>
3. {p: Person} / NOT (p.active == false) ==> active_person(p.id)
```

### Faits Injectés
```
1. Person(id:P001, name:Alice, active:true)
2. Person(id:P002, name:Bob, active:false)
```

---

## RÉSUMÉ GLOBAL

- **Tests réussis:** 12/12
- **Taux de réussite global:** 100.0%
- **Total mismatches:** 0
