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

1. [1] TestPerson{status=active, level=2, score=8.5, department=sales, age=25, active=true, name=Alice, tags=junior, salary=45000}
2. [1] TestPerson{level=5, active=true, score=9.2, department=engineering, name=Bob, age=35, tags=senior, salary=75000, status=active}
3. [1] TestPerson{active=false, status=inactive, level=1, age=16, score=6, department=hr, tags=intern, name=Charlie, salary=0}
4. [1] TestPerson{score=7.8, status=active, tags=manager, name=Diana, active=true, level=7, age=45, department=marketing, salary=85000}
5. [1] TestPerson{status=inactive, level=3, age=30, name=Eve, score=8, tags=employee, department=sales, salary=55000, active=false}
6. [1] TestPerson{level=1, name=Frank, active=true, tags=test, status=active, department=qa, score=0, age=0, salary=-5000}
7. [1] TestPerson{level=9, name=Grace, age=65, score=10, tags=executive, active=true, salary=95000, status=active, department=management}
8. [1] TestPerson{salary=25000, active=false, score=5.5, tags=junior, age=18, status=inactive, level=1, name=Henry, department=support}
9. [1] TestPerson{status=active, age=40, salary=68000, active=true, score=8.7, department=engineering, level=6, name=Ivy, tags=senior}
10. [1] TestPerson{tags=temp, level=1, age=22, active=true, status=active, department=intern, name=X, salary=28000, score=6.5}

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

1. [1] TestOrder{total=1999.98, priority=normal, customer_id=P001, product_id=PROD001, status=pending, amount=2, date=2024-01-15, discount=50, region=north}
2. [1] TestOrder{amount=1, priority=low, discount=0, region=south, status=confirmed, customer_id=P002, total=25.5, product_id=PROD002, date=2024-01-20}
3. [1] TestOrder{customer_id=P001, status=shipped, amount=3, date=2024-02-01, priority=high, discount=15, product_id=PROD003, total=225, region=north}
4. [1] TestOrder{status=delivered, priority=normal, date=2024-02-05, discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99}
5. [1] TestOrder{date=2024-02-10, status=confirmed, product_id=PROD001, amount=1, priority=high, discount=100, customer_id=P002, region=south, total=999.99}
6. [1] TestOrder{date=2024-02-15, status=cancelled, region=west, discount=0, customer_id=P005, product_id=PROD005, total=999.98, priority=low, amount=2}
7. [1] TestOrder{product_id=PROD006, status=shipped, customer_id=P007, date=2024-03-01, discount=50, amount=4, total=600, priority=urgent, region=north}
8. [1] TestOrder{amount=10, date=2024-03-05, status=pending, product_id=PROD002, total=255, discount=0, region=south, customer_id=P010, priority=normal}
9. [1] TestOrder{amount=1, date=2024-03-10, priority=low, region=north, product_id=PROD007, status=completed, discount=10, total=89.99, customer_id=P001}
10. [1] TestOrder{region=east, customer_id=P006, total=75000, status=refunded, discount=0, amount=1, date=2024-03-15, priority=urgent, product_id=PROD001}

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

1. [1] TestPerson{active=true, name=Alice, tags=junior, salary=45000, status=active, level=2, score=8.5, department=sales, age=25}
2. [1] TestPerson{age=35, tags=senior, salary=75000, status=active, level=5, active=true, score=9.2, department=engineering, name=Bob}
3. [1] TestPerson{level=1, age=16, score=6, department=hr, tags=intern, name=Charlie, salary=0, active=false, status=inactive}
4. [1] TestPerson{score=7.8, status=active, tags=manager, name=Diana, active=true, level=7, age=45, department=marketing, salary=85000}
5. [1] TestPerson{status=inactive, level=3, age=30, name=Eve, score=8, tags=employee, department=sales, salary=55000, active=false}
6. [1] TestPerson{active=true, tags=test, status=active, department=qa, score=0, age=0, salary=-5000, level=1, name=Frank}
7. [1] TestPerson{active=true, salary=95000, status=active, department=management, level=9, name=Grace, age=65, score=10, tags=executive}
8. [1] TestPerson{department=support, salary=25000, active=false, score=5.5, tags=junior, age=18, status=inactive, level=1, name=Henry}
9. [1] TestPerson{salary=68000, active=true, score=8.7, department=engineering, level=6, name=Ivy, tags=senior, status=active, age=40}
10. [1] TestPerson{department=intern, name=X, salary=28000, score=6.5, tags=temp, level=1, age=22, active=true, status=active}

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

1. [1] TestOrder{amount=2, date=2024-01-15, discount=50, region=north, total=1999.98, priority=normal, customer_id=P001, product_id=PROD001, status=pending}
2. [1] TestOrder{priority=low, discount=0, region=south, status=confirmed, customer_id=P002, total=25.5, product_id=PROD002, date=2024-01-20, amount=1}
3. [1] TestOrder{priority=high, discount=15, product_id=PROD003, total=225, region=north, customer_id=P001, status=shipped, amount=3, date=2024-02-01}
4. [1] TestOrder{customer_id=P004, product_id=PROD004, amount=1, total=299.99, status=delivered, priority=normal, date=2024-02-05, discount=0, region=east}
5. [1] TestOrder{priority=high, discount=100, customer_id=P002, region=south, total=999.99, date=2024-02-10, status=confirmed, product_id=PROD001, amount=1}
6. [1] TestOrder{date=2024-02-15, status=cancelled, region=west, discount=0, customer_id=P005, product_id=PROD005, total=999.98, priority=low, amount=2}
7. [1] TestOrder{amount=4, total=600, priority=urgent, region=north, product_id=PROD006, status=shipped, customer_id=P007, date=2024-03-01, discount=50}
8. [1] TestOrder{discount=0, region=south, customer_id=P010, priority=normal, amount=10, date=2024-03-05, status=pending, product_id=PROD002, total=255}
9. [1] TestOrder{region=north, product_id=PROD007, status=completed, discount=10, total=89.99, customer_id=P001, amount=1, date=2024-03-10, priority=low}
10. [1] TestOrder{region=east, customer_id=P006, total=75000, status=refunded, discount=0, amount=1, date=2024-03-15, priority=urgent, product_id=PROD001}

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

1. [1] TestProduct{available=true, stock=50, name=Laptop, rating=4.5, supplier=TechSupply, category=electronics, brand=TechCorp, price=999.99, keywords=computer}
2. [1] TestProduct{keywords=peripheral, supplier=TechSupply, category=accessories, brand=TechCorp, name=Mouse, available=true, stock=200, price=25.5, rating=4}
3. [1] TestProduct{name=Keyboard, price=75, supplier=KeySupply, stock=0, available=false, category=accessories, rating=3.5, keywords=typing, brand=KeyTech}
4. [1] TestProduct{category=electronics, rating=4.8, supplier=ScreenSupply, name=Monitor, available=true, stock=30, price=299.99, keywords=display, brand=ScreenPro}
5. [1] TestProduct{rating=2, available=false, brand=OldTech, name=OldKeyboard, keywords=obsolete, supplier=OldSupply, stock=0, category=accessories, price=8.5}
6. [1] TestProduct{category=audio, keywords=sound, stock=75, name=Headphones, price=150, available=true, brand=AudioMax, rating=4.6, supplier=AudioSupply}
7. [1] TestProduct{price=89.99, keywords=video, brand=CamTech, stock=25, supplier=CamSupply, name=Webcam, category=electronics, rating=3.8, available=true}

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

1. [1] TestPerson{salary=45000, status=active, level=2, score=8.5, department=sales, age=25, active=true, name=Alice, tags=junior}
2. [1] TestPerson{name=Bob, age=35, tags=senior, salary=75000, status=active, level=5, active=true, score=9.2, department=engineering}
3. [1] TestPerson{department=hr, tags=intern, name=Charlie, salary=0, active=false, status=inactive, level=1, age=16, score=6}
4. [1] TestPerson{name=Diana, active=true, level=7, age=45, department=marketing, salary=85000, score=7.8, status=active, tags=manager}
5. [1] TestPerson{age=30, name=Eve, score=8, tags=employee, department=sales, salary=55000, active=false, status=inactive, level=3}
6. [1] TestPerson{status=active, department=qa, score=0, age=0, salary=-5000, level=1, name=Frank, active=true, tags=test}
7. [1] TestPerson{score=10, tags=executive, active=true, salary=95000, status=active, department=management, level=9, name=Grace, age=65}
8. [1] TestPerson{salary=25000, active=false, score=5.5, tags=junior, age=18, status=inactive, level=1, name=Henry, department=support}
9. [1] TestPerson{salary=68000, active=true, score=8.7, department=engineering, level=6, name=Ivy, tags=senior, status=active, age=40}
10. [1] TestPerson{name=X, salary=28000, score=6.5, tags=temp, level=1, age=22, active=true, status=active, department=intern}

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

1. [1] TestOrder{priority=normal, customer_id=P001, product_id=PROD001, status=pending, amount=2, date=2024-01-15, discount=50, region=north, total=1999.98}
2. [1] TestOrder{total=25.5, product_id=PROD002, date=2024-01-20, amount=1, priority=low, discount=0, region=south, status=confirmed, customer_id=P002}
3. [1] TestOrder{customer_id=P001, status=shipped, amount=3, date=2024-02-01, priority=high, discount=15, product_id=PROD003, total=225, region=north}
4. [1] TestOrder{amount=1, total=299.99, status=delivered, priority=normal, date=2024-02-05, discount=0, region=east, customer_id=P004, product_id=PROD004}
5. [1] TestOrder{customer_id=P002, region=south, total=999.99, date=2024-02-10, status=confirmed, product_id=PROD001, amount=1, priority=high, discount=100}
6. [1] TestOrder{status=cancelled, region=west, discount=0, customer_id=P005, product_id=PROD005, total=999.98, priority=low, amount=2, date=2024-02-15}
7. [1] TestOrder{discount=50, amount=4, total=600, priority=urgent, region=north, product_id=PROD006, status=shipped, customer_id=P007, date=2024-03-01}
8. [1] TestOrder{priority=normal, amount=10, date=2024-03-05, status=pending, product_id=PROD002, total=255, discount=0, region=south, customer_id=P010}
9. [1] TestOrder{total=89.99, customer_id=P001, amount=1, date=2024-03-10, priority=low, region=north, product_id=PROD007, status=completed, discount=10}
10. [1] TestOrder{region=east, customer_id=P006, total=75000, status=refunded, discount=0, amount=1, date=2024-03-15, priority=urgent, product_id=PROD001}

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

1. [1] TestPerson{age=25, active=true, name=Alice, tags=junior, salary=45000, status=active, level=2, score=8.5, department=sales}
2. [1] TestPerson{level=5, active=true, score=9.2, department=engineering, name=Bob, age=35, tags=senior, salary=75000, status=active}
3. [1] TestPerson{department=hr, tags=intern, name=Charlie, salary=0, active=false, status=inactive, level=1, age=16, score=6}
4. [1] TestPerson{name=Diana, active=true, level=7, age=45, department=marketing, salary=85000, score=7.8, status=active, tags=manager}
5. [1] TestPerson{salary=55000, active=false, status=inactive, level=3, age=30, name=Eve, score=8, tags=employee, department=sales}
6. [1] TestPerson{department=qa, score=0, age=0, salary=-5000, level=1, name=Frank, active=true, tags=test, status=active}
7. [1] TestPerson{age=65, score=10, tags=executive, active=true, salary=95000, status=active, department=management, level=9, name=Grace}
8. [1] TestPerson{level=1, name=Henry, department=support, salary=25000, active=false, score=5.5, tags=junior, age=18, status=inactive}
9. [1] TestPerson{salary=68000, active=true, score=8.7, department=engineering, level=6, name=Ivy, tags=senior, status=active, age=40}
10. [1] TestPerson{tags=temp, level=1, age=22, active=true, status=active, department=intern, name=X, salary=28000, score=6.5}

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

1. [1] TestOrder{total=1999.98, priority=normal, customer_id=P001, product_id=PROD001, status=pending, amount=2, date=2024-01-15, discount=50, region=north}
2. [1] TestOrder{total=25.5, product_id=PROD002, date=2024-01-20, amount=1, priority=low, discount=0, region=south, status=confirmed, customer_id=P002}
3. [1] TestOrder{discount=15, product_id=PROD003, total=225, region=north, customer_id=P001, status=shipped, amount=3, date=2024-02-01, priority=high}
4. [1] TestOrder{discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, status=delivered, priority=normal, date=2024-02-05}
5. [1] TestOrder{product_id=PROD001, amount=1, priority=high, discount=100, customer_id=P002, region=south, total=999.99, date=2024-02-10, status=confirmed}
6. [1] TestOrder{date=2024-02-15, status=cancelled, region=west, discount=0, customer_id=P005, product_id=PROD005, total=999.98, priority=low, amount=2}
7. [1] TestOrder{amount=4, total=600, priority=urgent, region=north, product_id=PROD006, status=shipped, customer_id=P007, date=2024-03-01, discount=50}
8. [1] TestOrder{amount=10, date=2024-03-05, status=pending, product_id=PROD002, total=255, discount=0, region=south, customer_id=P010, priority=normal}
9. [1] TestOrder{customer_id=P001, amount=1, date=2024-03-10, priority=low, region=north, product_id=PROD007, status=completed, discount=10, total=89.99}
10. [1] TestOrder{product_id=PROD001, region=east, customer_id=P006, total=75000, status=refunded, discount=0, amount=1, date=2024-03-15, priority=urgent}

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

1. [1] TestPerson{tags=junior, salary=45000, status=active, level=2, score=8.5, department=sales, age=25, active=true, name=Alice}
2. [1] TestPerson{score=9.2, department=engineering, name=Bob, age=35, tags=senior, salary=75000, status=active, level=5, active=true}
3. [1] TestPerson{level=1, age=16, score=6, department=hr, tags=intern, name=Charlie, salary=0, active=false, status=inactive}
4. [1] TestPerson{department=marketing, salary=85000, score=7.8, status=active, tags=manager, name=Diana, active=true, level=7, age=45}
5. [1] TestPerson{tags=employee, department=sales, salary=55000, active=false, status=inactive, level=3, age=30, name=Eve, score=8}
6. [1] TestPerson{level=1, name=Frank, active=true, tags=test, status=active, department=qa, score=0, age=0, salary=-5000}
7. [1] TestPerson{salary=95000, status=active, department=management, level=9, name=Grace, age=65, score=10, tags=executive, active=true}
8. [1] TestPerson{age=18, status=inactive, level=1, name=Henry, department=support, salary=25000, active=false, score=5.5, tags=junior}
9. [1] TestPerson{status=active, age=40, salary=68000, active=true, score=8.7, department=engineering, level=6, name=Ivy, tags=senior}
10. [1] TestPerson{department=intern, name=X, salary=28000, score=6.5, tags=temp, level=1, age=22, active=true, status=active}

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

1. [1] TestProduct{available=true, stock=50, name=Laptop, rating=4.5, supplier=TechSupply, category=electronics, brand=TechCorp, price=999.99, keywords=computer}
2. [1] TestProduct{name=Mouse, available=true, stock=200, price=25.5, rating=4, keywords=peripheral, supplier=TechSupply, category=accessories, brand=TechCorp}
3. [1] TestProduct{supplier=KeySupply, stock=0, available=false, category=accessories, rating=3.5, keywords=typing, brand=KeyTech, name=Keyboard, price=75}
4. [1] TestProduct{price=299.99, keywords=display, brand=ScreenPro, category=electronics, rating=4.8, supplier=ScreenSupply, name=Monitor, available=true, stock=30}
5. [1] TestProduct{supplier=OldSupply, stock=0, category=accessories, price=8.5, rating=2, available=false, brand=OldTech, name=OldKeyboard, keywords=obsolete}
6. [1] TestProduct{keywords=sound, stock=75, name=Headphones, price=150, available=true, brand=AudioMax, rating=4.6, supplier=AudioSupply, category=audio}
7. [1] TestProduct{name=Webcam, category=electronics, rating=3.8, available=true, price=89.99, keywords=video, brand=CamTech, stock=25, supplier=CamSupply}

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

1. [1] TestPerson{active=true, name=Alice, tags=junior, salary=45000, status=active, level=2, score=8.5, department=sales, age=25}
2. [1] TestPerson{name=Bob, age=35, tags=senior, salary=75000, status=active, level=5, active=true, score=9.2, department=engineering}
3. [1] TestPerson{level=1, age=16, score=6, department=hr, tags=intern, name=Charlie, salary=0, active=false, status=inactive}
4. [1] TestPerson{salary=85000, score=7.8, status=active, tags=manager, name=Diana, active=true, level=7, age=45, department=marketing}
5. [1] TestPerson{status=inactive, level=3, age=30, name=Eve, score=8, tags=employee, department=sales, salary=55000, active=false}
6. [1] TestPerson{level=1, name=Frank, active=true, tags=test, status=active, department=qa, score=0, age=0, salary=-5000}
7. [1] TestPerson{active=true, salary=95000, status=active, department=management, level=9, name=Grace, age=65, score=10, tags=executive}
8. [1] TestPerson{level=1, name=Henry, department=support, salary=25000, active=false, score=5.5, tags=junior, age=18, status=inactive}
9. [1] TestPerson{tags=senior, status=active, age=40, salary=68000, active=true, score=8.7, department=engineering, level=6, name=Ivy}
10. [1] TestPerson{status=active, department=intern, name=X, salary=28000, score=6.5, tags=temp, level=1, age=22, active=true}

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

1. [1] TestOrder{amount=2, date=2024-01-15, discount=50, region=north, total=1999.98, priority=normal, customer_id=P001, product_id=PROD001, status=pending}
2. [1] TestOrder{total=25.5, product_id=PROD002, date=2024-01-20, amount=1, priority=low, discount=0, region=south, status=confirmed, customer_id=P002}
3. [1] TestOrder{priority=high, discount=15, product_id=PROD003, total=225, region=north, customer_id=P001, status=shipped, amount=3, date=2024-02-01}
4. [1] TestOrder{discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, status=delivered, priority=normal, date=2024-02-05}
5. [1] TestOrder{total=999.99, date=2024-02-10, status=confirmed, product_id=PROD001, amount=1, priority=high, discount=100, customer_id=P002, region=south}
6. [1] TestOrder{customer_id=P005, product_id=PROD005, total=999.98, priority=low, amount=2, date=2024-02-15, status=cancelled, region=west, discount=0}
7. [1] TestOrder{date=2024-03-01, discount=50, amount=4, total=600, priority=urgent, region=north, product_id=PROD006, status=shipped, customer_id=P007}
8. [1] TestOrder{amount=10, date=2024-03-05, status=pending, product_id=PROD002, total=255, discount=0, region=south, customer_id=P010, priority=normal}
9. [1] TestOrder{status=completed, discount=10, total=89.99, customer_id=P001, amount=1, date=2024-03-10, priority=low, region=north, product_id=PROD007}
10. [1] TestOrder{region=east, customer_id=P006, total=75000, status=refunded, discount=0, amount=1, date=2024-03-15, priority=urgent, product_id=PROD001}

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

1. [1] TestPerson{salary=45000, status=active, level=2, score=8.5, department=sales, age=25, active=true, name=Alice, tags=junior}
2. [1] TestPerson{level=5, active=true, score=9.2, department=engineering, name=Bob, age=35, tags=senior, salary=75000, status=active}
3. [1] TestPerson{active=false, status=inactive, level=1, age=16, score=6, department=hr, tags=intern, name=Charlie, salary=0}
4. [1] TestPerson{tags=manager, name=Diana, active=true, level=7, age=45, department=marketing, salary=85000, score=7.8, status=active}
5. [1] TestPerson{name=Eve, score=8, tags=employee, department=sales, salary=55000, active=false, status=inactive, level=3, age=30}
6. [1] TestPerson{level=1, name=Frank, active=true, tags=test, status=active, department=qa, score=0, age=0, salary=-5000}
7. [1] TestPerson{department=management, level=9, name=Grace, age=65, score=10, tags=executive, active=true, salary=95000, status=active}
8. [1] TestPerson{tags=junior, age=18, status=inactive, level=1, name=Henry, department=support, salary=25000, active=false, score=5.5}
9. [1] TestPerson{level=6, name=Ivy, tags=senior, status=active, age=40, salary=68000, active=true, score=8.7, department=engineering}
10. [1] TestPerson{age=22, active=true, status=active, department=intern, name=X, salary=28000, score=6.5, tags=temp, level=1}
11. [1] TestOrder{amount=2, date=2024-01-15, discount=50, region=north, total=1999.98, priority=normal, customer_id=P001, product_id=PROD001, status=pending}
12. [1] TestOrder{discount=0, region=south, status=confirmed, customer_id=P002, total=25.5, product_id=PROD002, date=2024-01-20, amount=1, priority=low}
13. [1] TestOrder{customer_id=P001, status=shipped, amount=3, date=2024-02-01, priority=high, discount=15, product_id=PROD003, total=225, region=north}
14. [1] TestOrder{discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, status=delivered, priority=normal, date=2024-02-05}
15. [1] TestOrder{product_id=PROD001, amount=1, priority=high, discount=100, customer_id=P002, region=south, total=999.99, date=2024-02-10, status=confirmed}
16. [1] TestOrder{total=999.98, priority=low, amount=2, date=2024-02-15, status=cancelled, region=west, discount=0, customer_id=P005, product_id=PROD005}
17. [1] TestOrder{priority=urgent, region=north, product_id=PROD006, status=shipped, customer_id=P007, date=2024-03-01, discount=50, amount=4, total=600}
18. [1] TestOrder{product_id=PROD002, total=255, discount=0, region=south, customer_id=P010, priority=normal, amount=10, date=2024-03-05, status=pending}
19. [1] TestOrder{product_id=PROD007, status=completed, discount=10, total=89.99, customer_id=P001, amount=1, date=2024-03-10, priority=low, region=north}
20. [1] TestOrder{region=east, customer_id=P006, total=75000, status=refunded, discount=0, amount=1, date=2024-03-15, priority=urgent, product_id=PROD001}

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

1. [1] TestPerson{age=25, active=true, name=Alice, tags=junior, salary=45000, status=active, level=2, score=8.5, department=sales}
2. [1] TestPerson{status=active, level=5, active=true, score=9.2, department=engineering, name=Bob, age=35, tags=senior, salary=75000}
3. [1] TestPerson{name=Charlie, salary=0, active=false, status=inactive, level=1, age=16, score=6, department=hr, tags=intern}
4. [1] TestPerson{name=Diana, active=true, level=7, age=45, department=marketing, salary=85000, score=7.8, status=active, tags=manager}
5. [1] TestPerson{name=Eve, score=8, tags=employee, department=sales, salary=55000, active=false, status=inactive, level=3, age=30}
6. [1] TestPerson{status=active, department=qa, score=0, age=0, salary=-5000, level=1, name=Frank, active=true, tags=test}
7. [1] TestPerson{active=true, salary=95000, status=active, department=management, level=9, name=Grace, age=65, score=10, tags=executive}
8. [1] TestPerson{level=1, name=Henry, department=support, salary=25000, active=false, score=5.5, tags=junior, age=18, status=inactive}
9. [1] TestPerson{department=engineering, level=6, name=Ivy, tags=senior, status=active, age=40, salary=68000, active=true, score=8.7}
10. [1] TestPerson{score=6.5, tags=temp, level=1, age=22, active=true, status=active, department=intern, name=X, salary=28000}

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

1. [1] TestPerson{name=Alice, tags=junior, salary=45000, status=active, level=2, score=8.5, department=sales, age=25, active=true}
2. [1] TestPerson{level=5, active=true, score=9.2, department=engineering, name=Bob, age=35, tags=senior, salary=75000, status=active}
3. [1] TestPerson{active=false, status=inactive, level=1, age=16, score=6, department=hr, tags=intern, name=Charlie, salary=0}
4. [1] TestPerson{name=Diana, active=true, level=7, age=45, department=marketing, salary=85000, score=7.8, status=active, tags=manager}
5. [1] TestPerson{score=8, tags=employee, department=sales, salary=55000, active=false, status=inactive, level=3, age=30, name=Eve}
6. [1] TestPerson{department=qa, score=0, age=0, salary=-5000, level=1, name=Frank, active=true, tags=test, status=active}
7. [1] TestPerson{active=true, salary=95000, status=active, department=management, level=9, name=Grace, age=65, score=10, tags=executive}
8. [1] TestPerson{salary=25000, active=false, score=5.5, tags=junior, age=18, status=inactive, level=1, name=Henry, department=support}
9. [1] TestPerson{salary=68000, active=true, score=8.7, department=engineering, level=6, name=Ivy, tags=senior, status=active, age=40}
10. [1] TestPerson{tags=temp, level=1, age=22, active=true, status=active, department=intern, name=X, salary=28000, score=6.5}
11. [1] TestOrder{region=north, total=1999.98, priority=normal, customer_id=P001, product_id=PROD001, status=pending, amount=2, date=2024-01-15, discount=50}
12. [1] TestOrder{amount=1, priority=low, discount=0, region=south, status=confirmed, customer_id=P002, total=25.5, product_id=PROD002, date=2024-01-20}
13. [1] TestOrder{status=shipped, amount=3, date=2024-02-01, priority=high, discount=15, product_id=PROD003, total=225, region=north, customer_id=P001}
14. [1] TestOrder{discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, status=delivered, priority=normal, date=2024-02-05}
15. [1] TestOrder{customer_id=P002, region=south, total=999.99, date=2024-02-10, status=confirmed, product_id=PROD001, amount=1, priority=high, discount=100}
16. [1] TestOrder{date=2024-02-15, status=cancelled, region=west, discount=0, customer_id=P005, product_id=PROD005, total=999.98, priority=low, amount=2}
17. [1] TestOrder{product_id=PROD006, status=shipped, customer_id=P007, date=2024-03-01, discount=50, amount=4, total=600, priority=urgent, region=north}
18. [1] TestOrder{total=255, discount=0, region=south, customer_id=P010, priority=normal, amount=10, date=2024-03-05, status=pending, product_id=PROD002}
19. [1] TestOrder{amount=1, date=2024-03-10, priority=low, region=north, product_id=PROD007, status=completed, discount=10, total=89.99, customer_id=P001}
20. [1] TestOrder{region=east, customer_id=P006, total=75000, status=refunded, discount=0, amount=1, date=2024-03-15, priority=urgent, product_id=PROD001}

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

1. [1] TestPerson{name=Alice, tags=junior, salary=45000, status=active, level=2, score=8.5, department=sales, age=25, active=true}
2. [1] TestPerson{tags=senior, salary=75000, status=active, level=5, active=true, score=9.2, department=engineering, name=Bob, age=35}
3. [1] TestPerson{tags=intern, name=Charlie, salary=0, active=false, status=inactive, level=1, age=16, score=6, department=hr}
4. [1] TestPerson{active=true, level=7, age=45, department=marketing, salary=85000, score=7.8, status=active, tags=manager, name=Diana}
5. [1] TestPerson{name=Eve, score=8, tags=employee, department=sales, salary=55000, active=false, status=inactive, level=3, age=30}
6. [1] TestPerson{age=0, salary=-5000, level=1, name=Frank, active=true, tags=test, status=active, department=qa, score=0}
7. [1] TestPerson{active=true, salary=95000, status=active, department=management, level=9, name=Grace, age=65, score=10, tags=executive}
8. [1] TestPerson{salary=25000, active=false, score=5.5, tags=junior, age=18, status=inactive, level=1, name=Henry, department=support}
9. [1] TestPerson{age=40, salary=68000, active=true, score=8.7, department=engineering, level=6, name=Ivy, tags=senior, status=active}
10. [1] TestPerson{age=22, active=true, status=active, department=intern, name=X, salary=28000, score=6.5, tags=temp, level=1}

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

1. [1] TestOrder{total=1999.98, priority=normal, customer_id=P001, product_id=PROD001, status=pending, amount=2, date=2024-01-15, discount=50, region=north}
2. [1] TestOrder{priority=low, discount=0, region=south, status=confirmed, customer_id=P002, total=25.5, product_id=PROD002, date=2024-01-20, amount=1}
3. [1] TestOrder{customer_id=P001, status=shipped, amount=3, date=2024-02-01, priority=high, discount=15, product_id=PROD003, total=225, region=north}
4. [1] TestOrder{amount=1, total=299.99, status=delivered, priority=normal, date=2024-02-05, discount=0, region=east, customer_id=P004, product_id=PROD004}
5. [1] TestOrder{product_id=PROD001, amount=1, priority=high, discount=100, customer_id=P002, region=south, total=999.99, date=2024-02-10, status=confirmed}
6. [1] TestOrder{total=999.98, priority=low, amount=2, date=2024-02-15, status=cancelled, region=west, discount=0, customer_id=P005, product_id=PROD005}
7. [1] TestOrder{priority=urgent, region=north, product_id=PROD006, status=shipped, customer_id=P007, date=2024-03-01, discount=50, amount=4, total=600}
8. [1] TestOrder{amount=10, date=2024-03-05, status=pending, product_id=PROD002, total=255, discount=0, region=south, customer_id=P010, priority=normal}
9. [1] TestOrder{priority=low, region=north, product_id=PROD007, status=completed, discount=10, total=89.99, customer_id=P001, amount=1, date=2024-03-10}
10. [1] TestOrder{customer_id=P006, total=75000, status=refunded, discount=0, amount=1, date=2024-03-15, priority=urgent, product_id=PROD001, region=east}

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

1. [1] TestPerson{salary=45000, status=active, level=2, score=8.5, department=sales, age=25, active=true, name=Alice, tags=junior}
2. [1] TestPerson{name=Bob, age=35, tags=senior, salary=75000, status=active, level=5, active=true, score=9.2, department=engineering}
3. [1] TestPerson{level=1, age=16, score=6, department=hr, tags=intern, name=Charlie, salary=0, active=false, status=inactive}
4. [1] TestPerson{salary=85000, score=7.8, status=active, tags=manager, name=Diana, active=true, level=7, age=45, department=marketing}
5. [1] TestPerson{department=sales, salary=55000, active=false, status=inactive, level=3, age=30, name=Eve, score=8, tags=employee}
6. [1] TestPerson{level=1, name=Frank, active=true, tags=test, status=active, department=qa, score=0, age=0, salary=-5000}
7. [1] TestPerson{name=Grace, age=65, score=10, tags=executive, active=true, salary=95000, status=active, department=management, level=9}
8. [1] TestPerson{department=support, salary=25000, active=false, score=5.5, tags=junior, age=18, status=inactive, level=1, name=Henry}
9. [1] TestPerson{status=active, age=40, salary=68000, active=true, score=8.7, department=engineering, level=6, name=Ivy, tags=senior}
10. [1] TestPerson{status=active, department=intern, name=X, salary=28000, score=6.5, tags=temp, level=1, age=22, active=true}

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
