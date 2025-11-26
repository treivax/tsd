# RÉSULTATS COMPLETS - ANALYSE RÈGLES DE NÉGATION TSD
=====================================================

**Date d'exécution**: 13 novembre 2025
**Fichier contraintes**: /home/resinsec/dev/tsd/constraint/test/integration/negation_rules.constraint
**Nombre de règles**: 19
**Nombre de faits**: 27

## 🎯 RÈGLE 0: not_zero_age

**Condition**: `NOT (p.age == 0)`
**Types concernés**: [TestPerson]
**Terminal**: rule_0_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{tags=junior, name=Alice, salary=45000, score=8.5, status=active, department=sales, level=2, age=25, active=true}
2. [1] TestPerson{age=35, salary=75000, active=true, score=9.2, department=engineering, status=active, name=Bob, tags=senior, level=5}
3. [1] TestPerson{age=16, salary=0, active=false, department=hr, tags=intern, status=inactive, score=6, name=Charlie, level=1}
4. [1] TestPerson{status=active, salary=85000, active=true, tags=manager, department=marketing, age=45, score=7.8, level=7, name=Diana}
5. [1] TestPerson{age=30, tags=employee, salary=55000, level=3, active=false, status=inactive, department=sales, score=8, name=Eve}
6. [1] TestPerson{tags=test, status=active, name=Frank, score=0, department=qa, level=1, active=true, age=0, salary=-5000}
7. [1] TestPerson{status=active, level=9, age=65, tags=executive, active=true, department=management, name=Grace, salary=95000, score=10}
8. [1] TestPerson{tags=junior, level=1, age=18, active=false, score=5.5, status=inactive, name=Henry, salary=25000, department=support}
9. [1] TestPerson{salary=68000, status=active, name=Ivy, age=40, department=engineering, active=true, score=8.7, tags=senior, level=6}
10. [1] TestPerson{salary=28000, tags=temp, department=intern, status=active, active=true, name=X, age=22, score=6.5, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{customer_id=P001, date=2024-01-15, priority=normal, discount=50, amount=2, total=1999.98, region=north, product_id=PROD001, status=pending}
2. [1] TestOrder{priority=low, discount=0, region=south, customer_id=P002, amount=1, total=25.5, status=confirmed, product_id=PROD002, date=2024-01-20}
3. [1] TestOrder{product_id=PROD003, amount=3, total=225, priority=high, discount=15, status=shipped, customer_id=P001, date=2024-02-01, region=north}
4. [1] TestOrder{customer_id=P004, product_id=PROD004, total=299.99, status=delivered, discount=0, priority=normal, amount=1, date=2024-02-05, region=east}
5. [1] TestOrder{customer_id=P002, date=2024-02-10, priority=high, product_id=PROD001, status=confirmed, discount=100, amount=1, total=999.99, region=south}
6. [1] TestOrder{total=999.98, discount=0, region=west, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, status=cancelled, priority=low}
7. [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, status=shipped, total=600, product_id=PROD006, amount=4, priority=urgent, discount=50}
8. [1] TestOrder{priority=normal, customer_id=P010, date=2024-03-05, region=south, discount=0, product_id=PROD002, amount=10, total=255, status=pending}
9. [1] TestOrder{priority=low, discount=10, product_id=PROD007, region=north, customer_id=P001, total=89.99, status=completed, amount=1, date=2024-03-10}
10. [1] TestOrder{product_id=PROD001, date=2024-03-15, customer_id=P006, region=east, amount=1, total=75000, priority=urgent, discount=0, status=refunded}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{score=8.5, status=active, department=sales, level=2, age=25, active=true, tags=junior, name=Alice, salary=45000}
2. [1] TestPerson{status=active, name=Bob, tags=senior, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering}
3. [1] TestPerson{age=16, salary=0, active=false, department=hr, tags=intern, status=inactive, score=6, name=Charlie, level=1}
4. [1] TestPerson{level=7, name=Diana, status=active, salary=85000, active=true, tags=manager, department=marketing, age=45, score=7.8}
5. [1] TestPerson{active=false, status=inactive, department=sales, score=8, name=Eve, age=30, tags=employee, salary=55000, level=3}
6. [1] TestPerson{score=0, department=qa, level=1, active=true, age=0, salary=-5000, tags=test, status=active, name=Frank}
7. [1] TestPerson{name=Grace, salary=95000, score=10, status=active, level=9, age=65, tags=executive, active=true, department=management}
8. [1] TestPerson{tags=junior, level=1, age=18, active=false, score=5.5, status=inactive, name=Henry, salary=25000, department=support}
9. [1] TestPerson{name=Ivy, age=40, department=engineering, active=true, score=8.7, tags=senior, level=6, salary=68000, status=active}
10. [1] TestPerson{salary=28000, tags=temp, department=intern, status=active, active=true, name=X, age=22, score=6.5, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{total=1999.98, region=north, product_id=PROD001, status=pending, customer_id=P001, date=2024-01-15, priority=normal, discount=50, amount=2}
2. [1] TestOrder{discount=0, region=south, customer_id=P002, amount=1, total=25.5, status=confirmed, product_id=PROD002, date=2024-01-20, priority=low}
3. [1] TestOrder{priority=high, discount=15, status=shipped, customer_id=P001, date=2024-02-01, region=north, product_id=PROD003, amount=3, total=225}
4. [1] TestOrder{customer_id=P004, product_id=PROD004, total=299.99, status=delivered, discount=0, priority=normal, amount=1, date=2024-02-05, region=east}
5. [1] TestOrder{product_id=PROD001, status=confirmed, discount=100, amount=1, total=999.99, region=south, customer_id=P002, date=2024-02-10, priority=high}
6. [1] TestOrder{total=999.98, discount=0, region=west, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, status=cancelled, priority=low}
7. [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, status=shipped, total=600, product_id=PROD006, amount=4, priority=urgent, discount=50}
8. [1] TestOrder{status=pending, priority=normal, customer_id=P010, date=2024-03-05, region=south, discount=0, product_id=PROD002, amount=10, total=255}
9. [1] TestOrder{customer_id=P001, total=89.99, status=completed, amount=1, date=2024-03-10, priority=low, discount=10, product_id=PROD007, region=north}
10. [1] TestOrder{product_id=PROD001, date=2024-03-15, customer_id=P006, region=east, amount=1, total=75000, priority=urgent, discount=0, status=refunded}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{keywords=computer, stock=50, supplier=TechSupply, name=Laptop, available=true, rating=4.5, brand=TechCorp, category=electronics, price=999.99}
2. [1] TestProduct{available=true, keywords=peripheral, name=Mouse, brand=TechCorp, stock=200, price=25.5, rating=4, supplier=TechSupply, category=accessories}
3. [1] TestProduct{brand=KeyTech, category=accessories, stock=0, name=Keyboard, price=75, available=false, keywords=typing, supplier=KeySupply, rating=3.5}
4. [1] TestProduct{keywords=display, rating=4.8, name=Monitor, brand=ScreenPro, stock=30, supplier=ScreenSupply, category=electronics, price=299.99, available=true}
5. [1] TestProduct{price=8.5, stock=0, rating=2, name=OldKeyboard, category=accessories, available=false, keywords=obsolete, brand=OldTech, supplier=OldSupply}
6. [1] TestProduct{supplier=AudioSupply, brand=AudioMax, name=Headphones, category=audio, price=150, available=true, keywords=sound, rating=4.6, stock=75}
7. [1] TestProduct{name=Webcam, keywords=video, supplier=CamSupply, category=electronics, available=true, rating=3.8, brand=CamTech, stock=25, price=89.99}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/7 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{tags=junior, name=Alice, salary=45000, score=8.5, status=active, department=sales, level=2, age=25, active=true}
2. [1] TestPerson{status=active, name=Bob, tags=senior, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering}
3. [1] TestPerson{level=1, age=16, salary=0, active=false, department=hr, tags=intern, status=inactive, score=6, name=Charlie}
4. [1] TestPerson{active=true, tags=manager, department=marketing, age=45, score=7.8, level=7, name=Diana, status=active, salary=85000}
5. [1] TestPerson{status=inactive, department=sales, score=8, name=Eve, age=30, tags=employee, salary=55000, level=3, active=false}
6. [1] TestPerson{age=0, salary=-5000, tags=test, status=active, name=Frank, score=0, department=qa, level=1, active=true}
7. [1] TestPerson{level=9, age=65, tags=executive, active=true, department=management, name=Grace, salary=95000, score=10, status=active}
8. [1] TestPerson{age=18, active=false, score=5.5, status=inactive, name=Henry, salary=25000, department=support, tags=junior, level=1}
9. [1] TestPerson{salary=68000, status=active, name=Ivy, age=40, department=engineering, active=true, score=8.7, tags=senior, level=6}
10. [1] TestPerson{active=true, name=X, age=22, score=6.5, level=1, salary=28000, tags=temp, department=intern, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{product_id=PROD001, status=pending, customer_id=P001, date=2024-01-15, priority=normal, discount=50, amount=2, total=1999.98, region=north}
2. [1] TestOrder{product_id=PROD002, date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002, amount=1, total=25.5, status=confirmed}
3. [1] TestOrder{region=north, product_id=PROD003, amount=3, total=225, priority=high, discount=15, status=shipped, customer_id=P001, date=2024-02-01}
4. [1] TestOrder{region=east, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, discount=0, priority=normal, amount=1, date=2024-02-05}
5. [1] TestOrder{amount=1, total=999.99, region=south, customer_id=P002, date=2024-02-10, priority=high, product_id=PROD001, status=confirmed, discount=100}
6. [1] TestOrder{total=999.98, discount=0, region=west, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, status=cancelled, priority=low}
7. [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, status=shipped, total=600, product_id=PROD006, amount=4, priority=urgent, discount=50}
8. [1] TestOrder{customer_id=P010, date=2024-03-05, region=south, discount=0, product_id=PROD002, amount=10, total=255, status=pending, priority=normal}
9. [1] TestOrder{date=2024-03-10, priority=low, discount=10, product_id=PROD007, region=north, customer_id=P001, total=89.99, status=completed, amount=1}
10. [1] TestOrder{region=east, amount=1, total=75000, priority=urgent, discount=0, status=refunded, product_id=PROD001, date=2024-03-15, customer_id=P006}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{department=sales, level=2, age=25, active=true, tags=junior, name=Alice, salary=45000, score=8.5, status=active}
2. [1] TestPerson{score=9.2, department=engineering, status=active, name=Bob, tags=senior, level=5, age=35, salary=75000, active=true}
3. [1] TestPerson{age=16, salary=0, active=false, department=hr, tags=intern, status=inactive, score=6, name=Charlie, level=1}
4. [1] TestPerson{status=active, salary=85000, active=true, tags=manager, department=marketing, age=45, score=7.8, level=7, name=Diana}
5. [1] TestPerson{name=Eve, age=30, tags=employee, salary=55000, level=3, active=false, status=inactive, department=sales, score=8}
6. [1] TestPerson{age=0, salary=-5000, tags=test, status=active, name=Frank, score=0, department=qa, level=1, active=true}
7. [1] TestPerson{name=Grace, salary=95000, score=10, status=active, level=9, age=65, tags=executive, active=true, department=management}
8. [1] TestPerson{name=Henry, salary=25000, department=support, tags=junior, level=1, age=18, active=false, score=5.5, status=inactive}
9. [1] TestPerson{age=40, department=engineering, active=true, score=8.7, tags=senior, level=6, salary=68000, status=active, name=Ivy}
10. [1] TestPerson{score=6.5, level=1, salary=28000, tags=temp, department=intern, status=active, active=true, name=X, age=22}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{customer_id=P001, date=2024-01-15, priority=normal, discount=50, amount=2, total=1999.98, region=north, product_id=PROD001, status=pending}
2. [1] TestOrder{date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002, amount=1, total=25.5, status=confirmed, product_id=PROD002}
3. [1] TestOrder{priority=high, discount=15, status=shipped, customer_id=P001, date=2024-02-01, region=north, product_id=PROD003, amount=3, total=225}
4. [1] TestOrder{total=299.99, status=delivered, discount=0, priority=normal, amount=1, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004}
5. [1] TestOrder{amount=1, total=999.99, region=south, customer_id=P002, date=2024-02-10, priority=high, product_id=PROD001, status=confirmed, discount=100}
6. [1] TestOrder{total=999.98, discount=0, region=west, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, status=cancelled, priority=low}
7. [1] TestOrder{product_id=PROD006, amount=4, priority=urgent, discount=50, region=north, customer_id=P007, date=2024-03-01, status=shipped, total=600}
8. [1] TestOrder{discount=0, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, customer_id=P010, date=2024-03-05, region=south}
9. [1] TestOrder{region=north, customer_id=P001, total=89.99, status=completed, amount=1, date=2024-03-10, priority=low, discount=10, product_id=PROD007}
10. [1] TestOrder{amount=1, total=75000, priority=urgent, discount=0, status=refunded, product_id=PROD001, date=2024-03-15, customer_id=P006, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{active=true, tags=junior, name=Alice, salary=45000, score=8.5, status=active, department=sales, level=2, age=25}
2. [1] TestPerson{tags=senior, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering, status=active, name=Bob}
3. [1] TestPerson{score=6, name=Charlie, level=1, age=16, salary=0, active=false, department=hr, tags=intern, status=inactive}
4. [1] TestPerson{score=7.8, level=7, name=Diana, status=active, salary=85000, active=true, tags=manager, department=marketing, age=45}
5. [1] TestPerson{level=3, active=false, status=inactive, department=sales, score=8, name=Eve, age=30, tags=employee, salary=55000}
6. [1] TestPerson{score=0, department=qa, level=1, active=true, age=0, salary=-5000, tags=test, status=active, name=Frank}
7. [1] TestPerson{age=65, tags=executive, active=true, department=management, name=Grace, salary=95000, score=10, status=active, level=9}
8. [1] TestPerson{name=Henry, salary=25000, department=support, tags=junior, level=1, age=18, active=false, score=5.5, status=inactive}
9. [1] TestPerson{name=Ivy, age=40, department=engineering, active=true, score=8.7, tags=senior, level=6, salary=68000, status=active}
10. [1] TestPerson{salary=28000, tags=temp, department=intern, status=active, active=true, name=X, age=22, score=6.5, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{name=Laptop, available=true, rating=4.5, brand=TechCorp, category=electronics, price=999.99, keywords=computer, stock=50, supplier=TechSupply}
2. [1] TestProduct{rating=4, supplier=TechSupply, category=accessories, available=true, keywords=peripheral, name=Mouse, brand=TechCorp, stock=200, price=25.5}
3. [1] TestProduct{name=Keyboard, price=75, available=false, keywords=typing, supplier=KeySupply, rating=3.5, brand=KeyTech, category=accessories, stock=0}
4. [1] TestProduct{stock=30, supplier=ScreenSupply, category=electronics, price=299.99, available=true, keywords=display, rating=4.8, name=Monitor, brand=ScreenPro}
5. [1] TestProduct{stock=0, rating=2, name=OldKeyboard, category=accessories, available=false, keywords=obsolete, brand=OldTech, supplier=OldSupply, price=8.5}
6. [1] TestProduct{supplier=AudioSupply, brand=AudioMax, name=Headphones, category=audio, price=150, available=true, keywords=sound, rating=4.6, stock=75}
7. [1] TestProduct{category=electronics, available=true, rating=3.8, brand=CamTech, stock=25, price=89.99, name=Webcam, keywords=video, supplier=CamSupply}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/7 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{active=true, tags=junior, name=Alice, salary=45000, score=8.5, status=active, department=sales, level=2, age=25}
2. [1] TestPerson{tags=senior, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering, status=active, name=Bob}
3. [1] TestPerson{score=6, name=Charlie, level=1, age=16, salary=0, active=false, department=hr, tags=intern, status=inactive}
4. [1] TestPerson{status=active, salary=85000, active=true, tags=manager, department=marketing, age=45, score=7.8, level=7, name=Diana}
5. [1] TestPerson{level=3, active=false, status=inactive, department=sales, score=8, name=Eve, age=30, tags=employee, salary=55000}
6. [1] TestPerson{age=0, salary=-5000, tags=test, status=active, name=Frank, score=0, department=qa, level=1, active=true}
7. [1] TestPerson{age=65, tags=executive, active=true, department=management, name=Grace, salary=95000, score=10, status=active, level=9}
8. [1] TestPerson{department=support, tags=junior, level=1, age=18, active=false, score=5.5, status=inactive, name=Henry, salary=25000}
9. [1] TestPerson{score=8.7, tags=senior, level=6, salary=68000, status=active, name=Ivy, age=40, department=engineering, active=true}
10. [1] TestPerson{tags=temp, department=intern, status=active, active=true, name=X, age=22, score=6.5, level=1, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{amount=2, total=1999.98, region=north, product_id=PROD001, status=pending, customer_id=P001, date=2024-01-15, priority=normal, discount=50}
2. [1] TestOrder{priority=low, discount=0, region=south, customer_id=P002, amount=1, total=25.5, status=confirmed, product_id=PROD002, date=2024-01-20}
3. [1] TestOrder{discount=15, status=shipped, customer_id=P001, date=2024-02-01, region=north, product_id=PROD003, amount=3, total=225, priority=high}
4. [1] TestOrder{status=delivered, discount=0, priority=normal, amount=1, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, total=299.99}
5. [1] TestOrder{customer_id=P002, date=2024-02-10, priority=high, product_id=PROD001, status=confirmed, discount=100, amount=1, total=999.99, region=south}
6. [1] TestOrder{discount=0, region=west, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, status=cancelled, priority=low, total=999.98}
7. [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, status=shipped, total=600, product_id=PROD006, amount=4, priority=urgent, discount=50}
8. [1] TestOrder{date=2024-03-05, region=south, discount=0, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, customer_id=P010}
9. [1] TestOrder{date=2024-03-10, priority=low, discount=10, product_id=PROD007, region=north, customer_id=P001, total=89.99, status=completed, amount=1}
10. [1] TestOrder{discount=0, status=refunded, product_id=PROD001, date=2024-03-15, customer_id=P006, region=east, amount=1, total=75000, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{status=active, department=sales, level=2, age=25, active=true, tags=junior, name=Alice, salary=45000, score=8.5}
2. [1] TestPerson{tags=senior, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering, status=active, name=Bob}
3. [1] TestPerson{age=16, salary=0, active=false, department=hr, tags=intern, status=inactive, score=6, name=Charlie, level=1}
4. [1] TestPerson{active=true, tags=manager, department=marketing, age=45, score=7.8, level=7, name=Diana, status=active, salary=85000}
5. [1] TestPerson{name=Eve, age=30, tags=employee, salary=55000, level=3, active=false, status=inactive, department=sales, score=8}
6. [1] TestPerson{department=qa, level=1, active=true, age=0, salary=-5000, tags=test, status=active, name=Frank, score=0}
7. [1] TestPerson{age=65, tags=executive, active=true, department=management, name=Grace, salary=95000, score=10, status=active, level=9}
8. [1] TestPerson{score=5.5, status=inactive, name=Henry, salary=25000, department=support, tags=junior, level=1, age=18, active=false}
9. [1] TestPerson{salary=68000, status=active, name=Ivy, age=40, department=engineering, active=true, score=8.7, tags=senior, level=6}
10. [1] TestPerson{active=true, name=X, age=22, score=6.5, level=1, salary=28000, tags=temp, department=intern, status=active}
11. [1] TestOrder{amount=2, total=1999.98, region=north, product_id=PROD001, status=pending, customer_id=P001, date=2024-01-15, priority=normal, discount=50}
12. [1] TestOrder{customer_id=P002, amount=1, total=25.5, status=confirmed, product_id=PROD002, date=2024-01-20, priority=low, discount=0, region=south}
13. [1] TestOrder{product_id=PROD003, amount=3, total=225, priority=high, discount=15, status=shipped, customer_id=P001, date=2024-02-01, region=north}
14. [1] TestOrder{region=east, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, discount=0, priority=normal, amount=1, date=2024-02-05}
15. [1] TestOrder{region=south, customer_id=P002, date=2024-02-10, priority=high, product_id=PROD001, status=confirmed, discount=100, amount=1, total=999.99}
16. [1] TestOrder{product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, status=cancelled, priority=low, total=999.98, discount=0, region=west}
17. [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, status=shipped, total=600, product_id=PROD006, amount=4, priority=urgent, discount=50}
18. [1] TestOrder{discount=0, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, customer_id=P010, date=2024-03-05, region=south}
19. [1] TestOrder{product_id=PROD007, region=north, customer_id=P001, total=89.99, status=completed, amount=1, date=2024-03-10, priority=low, discount=10}
20. [1] TestOrder{customer_id=P006, region=east, amount=1, total=75000, priority=urgent, discount=0, status=refunded, product_id=PROD001, date=2024-03-15}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/20 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{status=active, department=sales, level=2, age=25, active=true, tags=junior, name=Alice, salary=45000, score=8.5}
2. [1] TestPerson{tags=senior, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering, status=active, name=Bob}
3. [1] TestPerson{age=16, salary=0, active=false, department=hr, tags=intern, status=inactive, score=6, name=Charlie, level=1}
4. [1] TestPerson{status=active, salary=85000, active=true, tags=manager, department=marketing, age=45, score=7.8, level=7, name=Diana}
5. [1] TestPerson{name=Eve, age=30, tags=employee, salary=55000, level=3, active=false, status=inactive, department=sales, score=8}
6. [1] TestPerson{score=0, department=qa, level=1, active=true, age=0, salary=-5000, tags=test, status=active, name=Frank}
7. [1] TestPerson{name=Grace, salary=95000, score=10, status=active, level=9, age=65, tags=executive, active=true, department=management}
8. [1] TestPerson{name=Henry, salary=25000, department=support, tags=junior, level=1, age=18, active=false, score=5.5, status=inactive}
9. [1] TestPerson{salary=68000, status=active, name=Ivy, age=40, department=engineering, active=true, score=8.7, tags=senior, level=6}
10. [1] TestPerson{active=true, name=X, age=22, score=6.5, level=1, salary=28000, tags=temp, department=intern, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{tags=junior, name=Alice, salary=45000, score=8.5, status=active, department=sales, level=2, age=25, active=true}
2. [1] TestPerson{status=active, name=Bob, tags=senior, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering}
3. [1] TestPerson{age=16, salary=0, active=false, department=hr, tags=intern, status=inactive, score=6, name=Charlie, level=1}
4. [1] TestPerson{status=active, salary=85000, active=true, tags=manager, department=marketing, age=45, score=7.8, level=7, name=Diana}
5. [1] TestPerson{status=inactive, department=sales, score=8, name=Eve, age=30, tags=employee, salary=55000, level=3, active=false}
6. [1] TestPerson{score=0, department=qa, level=1, active=true, age=0, salary=-5000, tags=test, status=active, name=Frank}
7. [1] TestPerson{active=true, department=management, name=Grace, salary=95000, score=10, status=active, level=9, age=65, tags=executive}
8. [1] TestPerson{age=18, active=false, score=5.5, status=inactive, name=Henry, salary=25000, department=support, tags=junior, level=1}
9. [1] TestPerson{level=6, salary=68000, status=active, name=Ivy, age=40, department=engineering, active=true, score=8.7, tags=senior}
10. [1] TestPerson{tags=temp, department=intern, status=active, active=true, name=X, age=22, score=6.5, level=1, salary=28000}
11. [1] TestOrder{region=north, product_id=PROD001, status=pending, customer_id=P001, date=2024-01-15, priority=normal, discount=50, amount=2, total=1999.98}
12. [1] TestOrder{amount=1, total=25.5, status=confirmed, product_id=PROD002, date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002}
13. [1] TestOrder{customer_id=P001, date=2024-02-01, region=north, product_id=PROD003, amount=3, total=225, priority=high, discount=15, status=shipped}
14. [1] TestOrder{priority=normal, amount=1, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, discount=0}
15. [1] TestOrder{customer_id=P002, date=2024-02-10, priority=high, product_id=PROD001, status=confirmed, discount=100, amount=1, total=999.99, region=south}
16. [1] TestOrder{region=west, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, status=cancelled, priority=low, total=999.98, discount=0}
17. [1] TestOrder{date=2024-03-01, status=shipped, total=600, product_id=PROD006, amount=4, priority=urgent, discount=50, region=north, customer_id=P007}
18. [1] TestOrder{priority=normal, customer_id=P010, date=2024-03-05, region=south, discount=0, product_id=PROD002, amount=10, total=255, status=pending}
19. [1] TestOrder{customer_id=P001, total=89.99, status=completed, amount=1, date=2024-03-10, priority=low, discount=10, product_id=PROD007, region=north}
20. [1] TestOrder{total=75000, priority=urgent, discount=0, status=refunded, product_id=PROD001, date=2024-03-15, customer_id=P006, region=east, amount=1}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/20 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{status=active, department=sales, level=2, age=25, active=true, tags=junior, name=Alice, salary=45000, score=8.5}
2. [1] TestPerson{name=Bob, tags=senior, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering, status=active}
3. [1] TestPerson{score=6, name=Charlie, level=1, age=16, salary=0, active=false, department=hr, tags=intern, status=inactive}
4. [1] TestPerson{score=7.8, level=7, name=Diana, status=active, salary=85000, active=true, tags=manager, department=marketing, age=45}
5. [1] TestPerson{name=Eve, age=30, tags=employee, salary=55000, level=3, active=false, status=inactive, department=sales, score=8}
6. [1] TestPerson{age=0, salary=-5000, tags=test, status=active, name=Frank, score=0, department=qa, level=1, active=true}
7. [1] TestPerson{name=Grace, salary=95000, score=10, status=active, level=9, age=65, tags=executive, active=true, department=management}
8. [1] TestPerson{name=Henry, salary=25000, department=support, tags=junior, level=1, age=18, active=false, score=5.5, status=inactive}
9. [1] TestPerson{status=active, name=Ivy, age=40, department=engineering, active=true, score=8.7, tags=senior, level=6, salary=68000}
10. [1] TestPerson{active=true, name=X, age=22, score=6.5, level=1, salary=28000, tags=temp, department=intern, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{product_id=PROD001, status=pending, customer_id=P001, date=2024-01-15, priority=normal, discount=50, amount=2, total=1999.98, region=north}
2. [1] TestOrder{discount=0, region=south, customer_id=P002, amount=1, total=25.5, status=confirmed, product_id=PROD002, date=2024-01-20, priority=low}
3. [1] TestOrder{customer_id=P001, date=2024-02-01, region=north, product_id=PROD003, amount=3, total=225, priority=high, discount=15, status=shipped}
4. [1] TestOrder{customer_id=P004, product_id=PROD004, total=299.99, status=delivered, discount=0, priority=normal, amount=1, date=2024-02-05, region=east}
5. [1] TestOrder{discount=100, amount=1, total=999.99, region=south, customer_id=P002, date=2024-02-10, priority=high, product_id=PROD001, status=confirmed}
6. [1] TestOrder{total=999.98, discount=0, region=west, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, status=cancelled, priority=low}
7. [1] TestOrder{date=2024-03-01, status=shipped, total=600, product_id=PROD006, amount=4, priority=urgent, discount=50, region=north, customer_id=P007}
8. [1] TestOrder{date=2024-03-05, region=south, discount=0, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, customer_id=P010}
9. [1] TestOrder{region=north, customer_id=P001, total=89.99, status=completed, amount=1, date=2024-03-10, priority=low, discount=10, product_id=PROD007}
10. [1] TestOrder{status=refunded, product_id=PROD001, date=2024-03-15, customer_id=P006, region=east, amount=1, total=75000, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{name=Alice, salary=45000, score=8.5, status=active, department=sales, level=2, age=25, active=true, tags=junior}
2. [1] TestPerson{tags=senior, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering, status=active, name=Bob}
3. [1] TestPerson{age=16, salary=0, active=false, department=hr, tags=intern, status=inactive, score=6, name=Charlie, level=1}
4. [1] TestPerson{status=active, salary=85000, active=true, tags=manager, department=marketing, age=45, score=7.8, level=7, name=Diana}
5. [1] TestPerson{level=3, active=false, status=inactive, department=sales, score=8, name=Eve, age=30, tags=employee, salary=55000}
6. [1] TestPerson{age=0, salary=-5000, tags=test, status=active, name=Frank, score=0, department=qa, level=1, active=true}
7. [1] TestPerson{level=9, age=65, tags=executive, active=true, department=management, name=Grace, salary=95000, score=10, status=active}
8. [1] TestPerson{salary=25000, department=support, tags=junior, level=1, age=18, active=false, score=5.5, status=inactive, name=Henry}
9. [1] TestPerson{name=Ivy, age=40, department=engineering, active=true, score=8.7, tags=senior, level=6, salary=68000, status=active}
10. [1] TestPerson{active=true, name=X, age=22, score=6.5, level=1, salary=28000, tags=temp, department=intern, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 0 (0.0%)
- **Tokens générés**: 0
- **Faits traités**: 27
