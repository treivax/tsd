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

1. [1] TestPerson{salary=45000, tags=junior, status=active, score=8.5, name=Alice, age=25, active=true, department=sales, level=2}
2. [1] TestPerson{name=Bob, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, salary=75000}
3. [1] TestPerson{status=inactive, age=16, department=hr, active=false, score=6, name=Charlie, salary=0, tags=intern, level=1}
4. [1] TestPerson{score=7.8, tags=manager, level=7, active=true, status=active, age=45, salary=85000, department=marketing, name=Diana}
5. [1] TestPerson{name=Eve, salary=55000, active=false, department=sales, level=3, age=30, score=8, tags=employee, status=inactive}
6. [1] TestPerson{score=0, status=active, level=1, tags=test, department=qa, age=0, salary=-5000, active=true, name=Frank}
7. [1] TestPerson{department=management, level=9, name=Grace, salary=95000, active=true, tags=executive, score=10, age=65, status=active}
8. [1] TestPerson{salary=25000, tags=junior, name=Henry, active=false, score=5.5, status=inactive, age=18, department=support, level=1}
9. [1] TestPerson{name=Ivy, active=true, department=engineering, level=6, age=40, status=active, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{status=active, department=intern, salary=28000, level=1, name=X, age=22, active=true, score=6.5, tags=temp}

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

1. [1] TestOrder{product_id=PROD001, amount=2, date=2024-01-15, priority=normal, status=pending, total=1999.98, region=north, customer_id=P001, discount=50}
2. [1] TestOrder{customer_id=P002, date=2024-01-20, status=confirmed, product_id=PROD002, amount=1, discount=0, total=25.5, priority=low, region=south}
3. [1] TestOrder{date=2024-02-01, status=shipped, amount=3, priority=high, customer_id=P001, discount=15, product_id=PROD003, region=north, total=225}
4. [1] TestOrder{discount=0, region=east, customer_id=P004, amount=1, date=2024-02-05, priority=normal, total=299.99, status=delivered, product_id=PROD004}
5. [1] TestOrder{product_id=PROD001, status=confirmed, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, region=south, customer_id=P002}
6. [1] TestOrder{amount=2, total=999.98, region=west, customer_id=P005, date=2024-02-15, discount=0, product_id=PROD005, status=cancelled, priority=low}
7. [1] TestOrder{priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped}
8. [1] TestOrder{product_id=PROD002, status=pending, priority=normal, region=south, total=255, date=2024-03-05, amount=10, customer_id=P010, discount=0}
9. [1] TestOrder{date=2024-03-10, status=completed, discount=10, product_id=PROD007, region=north, total=89.99, priority=low, amount=1, customer_id=P001}
10. [1] TestOrder{customer_id=P006, product_id=PROD001, status=refunded, discount=0, region=east, amount=1, date=2024-03-15, total=75000, priority=urgent}

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

1. [1] TestPerson{status=active, score=8.5, name=Alice, age=25, active=true, department=sales, level=2, salary=45000, tags=junior}
2. [1] TestPerson{status=active, department=engineering, level=5, salary=75000, name=Bob, age=35, active=true, score=9.2, tags=senior}
3. [1] TestPerson{tags=intern, level=1, status=inactive, age=16, department=hr, active=false, score=6, name=Charlie, salary=0}
4. [1] TestPerson{name=Diana, score=7.8, tags=manager, level=7, active=true, status=active, age=45, salary=85000, department=marketing}
5. [1] TestPerson{score=8, tags=employee, status=inactive, name=Eve, salary=55000, active=false, department=sales, level=3, age=30}
6. [1] TestPerson{salary=-5000, active=true, name=Frank, score=0, status=active, level=1, tags=test, department=qa, age=0}
7. [1] TestPerson{score=10, age=65, status=active, department=management, level=9, name=Grace, salary=95000, active=true, tags=executive}
8. [1] TestPerson{tags=junior, name=Henry, active=false, score=5.5, status=inactive, age=18, department=support, level=1, salary=25000}
9. [1] TestPerson{age=40, status=active, salary=68000, score=8.7, tags=senior, name=Ivy, active=true, department=engineering, level=6}
10. [1] TestPerson{score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1, name=X, age=22, active=true}

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

1. [1] TestOrder{product_id=PROD001, amount=2, date=2024-01-15, priority=normal, status=pending, total=1999.98, region=north, customer_id=P001, discount=50}
2. [1] TestOrder{date=2024-01-20, status=confirmed, product_id=PROD002, amount=1, discount=0, total=25.5, priority=low, region=south, customer_id=P002}
3. [1] TestOrder{discount=15, product_id=PROD003, region=north, total=225, date=2024-02-01, status=shipped, amount=3, priority=high, customer_id=P001}
4. [1] TestOrder{discount=0, region=east, customer_id=P004, amount=1, date=2024-02-05, priority=normal, total=299.99, status=delivered, product_id=PROD004}
5. [1] TestOrder{product_id=PROD001, status=confirmed, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, region=south, customer_id=P002}
6. [1] TestOrder{customer_id=P005, date=2024-02-15, discount=0, product_id=PROD005, status=cancelled, priority=low, amount=2, total=999.98, region=west}
7. [1] TestOrder{discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent}
8. [1] TestOrder{amount=10, customer_id=P010, discount=0, product_id=PROD002, status=pending, priority=normal, region=south, total=255, date=2024-03-05}
9. [1] TestOrder{region=north, total=89.99, priority=low, amount=1, customer_id=P001, date=2024-03-10, status=completed, discount=10, product_id=PROD007}
10. [1] TestOrder{amount=1, date=2024-03-15, total=75000, priority=urgent, customer_id=P006, product_id=PROD001, status=refunded, discount=0, region=east}

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

1. [1] TestProduct{stock=50, category=electronics, brand=TechCorp, rating=4.5, keywords=computer, supplier=TechSupply, name=Laptop, price=999.99, available=true}
2. [1] TestProduct{available=true, stock=200, keywords=peripheral, category=accessories, rating=4, brand=TechCorp, supplier=TechSupply, price=25.5, name=Mouse}
3. [1] TestProduct{price=75, available=false, rating=3.5, brand=KeyTech, stock=0, name=Keyboard, keywords=typing, supplier=KeySupply, category=accessories}
4. [1] TestProduct{name=Monitor, category=electronics, keywords=display, brand=ScreenPro, price=299.99, available=true, rating=4.8, supplier=ScreenSupply, stock=30}
5. [1] TestProduct{supplier=OldSupply, name=OldKeyboard, available=false, rating=2, keywords=obsolete, brand=OldTech, category=accessories, price=8.5, stock=0}
6. [1] TestProduct{available=true, rating=4.6, stock=75, supplier=AudioSupply, category=audio, keywords=sound, brand=AudioMax, name=Headphones, price=150}
7. [1] TestProduct{keywords=video, stock=25, price=89.99, available=true, rating=3.8, brand=CamTech, category=electronics, supplier=CamSupply, name=Webcam}

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

1. [1] TestPerson{salary=45000, tags=junior, status=active, score=8.5, name=Alice, age=25, active=true, department=sales, level=2}
2. [1] TestPerson{name=Bob, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, salary=75000}
3. [1] TestPerson{score=6, name=Charlie, salary=0, tags=intern, level=1, status=inactive, age=16, department=hr, active=false}
4. [1] TestPerson{active=true, status=active, age=45, salary=85000, department=marketing, name=Diana, score=7.8, tags=manager, level=7}
5. [1] TestPerson{tags=employee, status=inactive, name=Eve, salary=55000, active=false, department=sales, level=3, age=30, score=8}
6. [1] TestPerson{salary=-5000, active=true, name=Frank, score=0, status=active, level=1, tags=test, department=qa, age=0}
7. [1] TestPerson{score=10, age=65, status=active, department=management, level=9, name=Grace, salary=95000, active=true, tags=executive}
8. [1] TestPerson{tags=junior, name=Henry, active=false, score=5.5, status=inactive, age=18, department=support, level=1, salary=25000}
9. [1] TestPerson{name=Ivy, active=true, department=engineering, level=6, age=40, status=active, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{age=22, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1, name=X}

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

1. [1] TestOrder{region=north, customer_id=P001, discount=50, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, status=pending, total=1999.98}
2. [1] TestOrder{customer_id=P002, date=2024-01-20, status=confirmed, product_id=PROD002, amount=1, discount=0, total=25.5, priority=low, region=south}
3. [1] TestOrder{discount=15, product_id=PROD003, region=north, total=225, date=2024-02-01, status=shipped, amount=3, priority=high, customer_id=P001}
4. [1] TestOrder{customer_id=P004, amount=1, date=2024-02-05, priority=normal, total=299.99, status=delivered, product_id=PROD004, discount=0, region=east}
5. [1] TestOrder{date=2024-02-10, region=south, customer_id=P002, product_id=PROD001, status=confirmed, priority=high, discount=100, amount=1, total=999.99}
6. [1] TestOrder{customer_id=P005, date=2024-02-15, discount=0, product_id=PROD005, status=cancelled, priority=low, amount=2, total=999.98, region=west}
7. [1] TestOrder{date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600}
8. [1] TestOrder{amount=10, customer_id=P010, discount=0, product_id=PROD002, status=pending, priority=normal, region=south, total=255, date=2024-03-05}
9. [1] TestOrder{discount=10, product_id=PROD007, region=north, total=89.99, priority=low, amount=1, customer_id=P001, date=2024-03-10, status=completed}
10. [1] TestOrder{product_id=PROD001, status=refunded, discount=0, region=east, amount=1, date=2024-03-15, total=75000, priority=urgent, customer_id=P006}

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

1. [1] TestPerson{salary=45000, tags=junior, status=active, score=8.5, name=Alice, age=25, active=true, department=sales, level=2}
2. [1] TestPerson{department=engineering, level=5, salary=75000, name=Bob, age=35, active=true, score=9.2, tags=senior, status=active}
3. [1] TestPerson{status=inactive, age=16, department=hr, active=false, score=6, name=Charlie, salary=0, tags=intern, level=1}
4. [1] TestPerson{status=active, age=45, salary=85000, department=marketing, name=Diana, score=7.8, tags=manager, level=7, active=true}
5. [1] TestPerson{score=8, tags=employee, status=inactive, name=Eve, salary=55000, active=false, department=sales, level=3, age=30}
6. [1] TestPerson{status=active, level=1, tags=test, department=qa, age=0, salary=-5000, active=true, name=Frank, score=0}
7. [1] TestPerson{score=10, age=65, status=active, department=management, level=9, name=Grace, salary=95000, active=true, tags=executive}
8. [1] TestPerson{salary=25000, tags=junior, name=Henry, active=false, score=5.5, status=inactive, age=18, department=support, level=1}
9. [1] TestPerson{name=Ivy, active=true, department=engineering, level=6, age=40, status=active, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{age=22, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1, name=X}

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

1. [1] TestOrder{total=1999.98, region=north, customer_id=P001, discount=50, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, status=pending}
2. [1] TestOrder{discount=0, total=25.5, priority=low, region=south, customer_id=P002, date=2024-01-20, status=confirmed, product_id=PROD002, amount=1}
3. [1] TestOrder{amount=3, priority=high, customer_id=P001, discount=15, product_id=PROD003, region=north, total=225, date=2024-02-01, status=shipped}
4. [1] TestOrder{status=delivered, product_id=PROD004, discount=0, region=east, customer_id=P004, amount=1, date=2024-02-05, priority=normal, total=299.99}
5. [1] TestOrder{priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, region=south, customer_id=P002, product_id=PROD001, status=confirmed}
6. [1] TestOrder{customer_id=P005, date=2024-02-15, discount=0, product_id=PROD005, status=cancelled, priority=low, amount=2, total=999.98, region=west}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}
8. [1] TestOrder{product_id=PROD002, status=pending, priority=normal, region=south, total=255, date=2024-03-05, amount=10, customer_id=P010, discount=0}
9. [1] TestOrder{date=2024-03-10, status=completed, discount=10, product_id=PROD007, region=north, total=89.99, priority=low, amount=1, customer_id=P001}
10. [1] TestOrder{customer_id=P006, product_id=PROD001, status=refunded, discount=0, region=east, amount=1, date=2024-03-15, total=75000, priority=urgent}

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

1. [1] TestPerson{department=sales, level=2, salary=45000, tags=junior, status=active, score=8.5, name=Alice, age=25, active=true}
2. [1] TestPerson{name=Bob, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, salary=75000}
3. [1] TestPerson{active=false, score=6, name=Charlie, salary=0, tags=intern, level=1, status=inactive, age=16, department=hr}
4. [1] TestPerson{status=active, age=45, salary=85000, department=marketing, name=Diana, score=7.8, tags=manager, level=7, active=true}
5. [1] TestPerson{level=3, age=30, score=8, tags=employee, status=inactive, name=Eve, salary=55000, active=false, department=sales}
6. [1] TestPerson{active=true, name=Frank, score=0, status=active, level=1, tags=test, department=qa, age=0, salary=-5000}
7. [1] TestPerson{score=10, age=65, status=active, department=management, level=9, name=Grace, salary=95000, active=true, tags=executive}
8. [1] TestPerson{status=inactive, age=18, department=support, level=1, salary=25000, tags=junior, name=Henry, active=false, score=5.5}
9. [1] TestPerson{name=Ivy, active=true, department=engineering, level=6, age=40, status=active, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{age=22, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1, name=X}

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

1. [1] TestProduct{stock=50, category=electronics, brand=TechCorp, rating=4.5, keywords=computer, supplier=TechSupply, name=Laptop, price=999.99, available=true}
2. [1] TestProduct{brand=TechCorp, supplier=TechSupply, price=25.5, name=Mouse, available=true, stock=200, keywords=peripheral, category=accessories, rating=4}
3. [1] TestProduct{stock=0, name=Keyboard, keywords=typing, supplier=KeySupply, category=accessories, price=75, available=false, rating=3.5, brand=KeyTech}
4. [1] TestProduct{price=299.99, available=true, rating=4.8, supplier=ScreenSupply, stock=30, name=Monitor, category=electronics, keywords=display, brand=ScreenPro}
5. [1] TestProduct{available=false, rating=2, keywords=obsolete, brand=OldTech, category=accessories, price=8.5, stock=0, supplier=OldSupply, name=OldKeyboard}
6. [1] TestProduct{rating=4.6, stock=75, supplier=AudioSupply, category=audio, keywords=sound, brand=AudioMax, name=Headphones, price=150, available=true}
7. [1] TestProduct{available=true, rating=3.8, brand=CamTech, category=electronics, supplier=CamSupply, name=Webcam, keywords=video, stock=25, price=89.99}

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

1. [1] TestPerson{department=sales, level=2, salary=45000, tags=junior, status=active, score=8.5, name=Alice, age=25, active=true}
2. [1] TestPerson{name=Bob, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, salary=75000}
3. [1] TestPerson{department=hr, active=false, score=6, name=Charlie, salary=0, tags=intern, level=1, status=inactive, age=16}
4. [1] TestPerson{score=7.8, tags=manager, level=7, active=true, status=active, age=45, salary=85000, department=marketing, name=Diana}
5. [1] TestPerson{name=Eve, salary=55000, active=false, department=sales, level=3, age=30, score=8, tags=employee, status=inactive}
6. [1] TestPerson{score=0, status=active, level=1, tags=test, department=qa, age=0, salary=-5000, active=true, name=Frank}
7. [1] TestPerson{department=management, level=9, name=Grace, salary=95000, active=true, tags=executive, score=10, age=65, status=active}
8. [1] TestPerson{score=5.5, status=inactive, age=18, department=support, level=1, salary=25000, tags=junior, name=Henry, active=false}
9. [1] TestPerson{name=Ivy, active=true, department=engineering, level=6, age=40, status=active, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{age=22, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1, name=X}

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

1. [1] TestOrder{product_id=PROD001, amount=2, date=2024-01-15, priority=normal, status=pending, total=1999.98, region=north, customer_id=P001, discount=50}
2. [1] TestOrder{total=25.5, priority=low, region=south, customer_id=P002, date=2024-01-20, status=confirmed, product_id=PROD002, amount=1, discount=0}
3. [1] TestOrder{total=225, date=2024-02-01, status=shipped, amount=3, priority=high, customer_id=P001, discount=15, product_id=PROD003, region=north}
4. [1] TestOrder{discount=0, region=east, customer_id=P004, amount=1, date=2024-02-05, priority=normal, total=299.99, status=delivered, product_id=PROD004}
5. [1] TestOrder{product_id=PROD001, status=confirmed, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, region=south, customer_id=P002}
6. [1] TestOrder{customer_id=P005, date=2024-02-15, discount=0, product_id=PROD005, status=cancelled, priority=low, amount=2, total=999.98, region=west}
7. [1] TestOrder{region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50}
8. [1] TestOrder{status=pending, priority=normal, region=south, total=255, date=2024-03-05, amount=10, customer_id=P010, discount=0, product_id=PROD002}
9. [1] TestOrder{discount=10, product_id=PROD007, region=north, total=89.99, priority=low, amount=1, customer_id=P001, date=2024-03-10, status=completed}
10. [1] TestOrder{amount=1, date=2024-03-15, total=75000, priority=urgent, customer_id=P006, product_id=PROD001, status=refunded, discount=0, region=east}

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

1. [1] TestPerson{status=active, score=8.5, name=Alice, age=25, active=true, department=sales, level=2, salary=45000, tags=junior}
2. [1] TestPerson{name=Bob, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, salary=75000}
3. [1] TestPerson{salary=0, tags=intern, level=1, status=inactive, age=16, department=hr, active=false, score=6, name=Charlie}
4. [1] TestPerson{age=45, salary=85000, department=marketing, name=Diana, score=7.8, tags=manager, level=7, active=true, status=active}
5. [1] TestPerson{level=3, age=30, score=8, tags=employee, status=inactive, name=Eve, salary=55000, active=false, department=sales}
6. [1] TestPerson{score=0, status=active, level=1, tags=test, department=qa, age=0, salary=-5000, active=true, name=Frank}
7. [1] TestPerson{score=10, age=65, status=active, department=management, level=9, name=Grace, salary=95000, active=true, tags=executive}
8. [1] TestPerson{name=Henry, active=false, score=5.5, status=inactive, age=18, department=support, level=1, salary=25000, tags=junior}
9. [1] TestPerson{salary=68000, score=8.7, tags=senior, name=Ivy, active=true, department=engineering, level=6, age=40, status=active}
10. [1] TestPerson{salary=28000, level=1, name=X, age=22, active=true, score=6.5, tags=temp, status=active, department=intern}
11. [1] TestOrder{total=1999.98, region=north, customer_id=P001, discount=50, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, status=pending}
12. [1] TestOrder{status=confirmed, product_id=PROD002, amount=1, discount=0, total=25.5, priority=low, region=south, customer_id=P002, date=2024-01-20}
13. [1] TestOrder{discount=15, product_id=PROD003, region=north, total=225, date=2024-02-01, status=shipped, amount=3, priority=high, customer_id=P001}
14. [1] TestOrder{total=299.99, status=delivered, product_id=PROD004, discount=0, region=east, customer_id=P004, amount=1, date=2024-02-05, priority=normal}
15. [1] TestOrder{amount=1, total=999.99, date=2024-02-10, region=south, customer_id=P002, product_id=PROD001, status=confirmed, priority=high, discount=100}
16. [1] TestOrder{customer_id=P005, date=2024-02-15, discount=0, product_id=PROD005, status=cancelled, priority=low, amount=2, total=999.98, region=west}
17. [1] TestOrder{customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}
18. [1] TestOrder{amount=10, customer_id=P010, discount=0, product_id=PROD002, status=pending, priority=normal, region=south, total=255, date=2024-03-05}
19. [1] TestOrder{region=north, total=89.99, priority=low, amount=1, customer_id=P001, date=2024-03-10, status=completed, discount=10, product_id=PROD007}
20. [1] TestOrder{customer_id=P006, product_id=PROD001, status=refunded, discount=0, region=east, amount=1, date=2024-03-15, total=75000, priority=urgent}

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

1. [1] TestPerson{department=sales, level=2, salary=45000, tags=junior, status=active, score=8.5, name=Alice, age=25, active=true}
2. [1] TestPerson{status=active, department=engineering, level=5, salary=75000, name=Bob, age=35, active=true, score=9.2, tags=senior}
3. [1] TestPerson{department=hr, active=false, score=6, name=Charlie, salary=0, tags=intern, level=1, status=inactive, age=16}
4. [1] TestPerson{name=Diana, score=7.8, tags=manager, level=7, active=true, status=active, age=45, salary=85000, department=marketing}
5. [1] TestPerson{status=inactive, name=Eve, salary=55000, active=false, department=sales, level=3, age=30, score=8, tags=employee}
6. [1] TestPerson{score=0, status=active, level=1, tags=test, department=qa, age=0, salary=-5000, active=true, name=Frank}
7. [1] TestPerson{department=management, level=9, name=Grace, salary=95000, active=true, tags=executive, score=10, age=65, status=active}
8. [1] TestPerson{status=inactive, age=18, department=support, level=1, salary=25000, tags=junior, name=Henry, active=false, score=5.5}
9. [1] TestPerson{active=true, department=engineering, level=6, age=40, status=active, salary=68000, score=8.7, tags=senior, name=Ivy}
10. [1] TestPerson{department=intern, salary=28000, level=1, name=X, age=22, active=true, score=6.5, tags=temp, status=active}

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

1. [1] TestPerson{level=2, salary=45000, tags=junior, status=active, score=8.5, name=Alice, age=25, active=true, department=sales}
2. [1] TestPerson{status=active, department=engineering, level=5, salary=75000, name=Bob, age=35, active=true, score=9.2, tags=senior}
3. [1] TestPerson{department=hr, active=false, score=6, name=Charlie, salary=0, tags=intern, level=1, status=inactive, age=16}
4. [1] TestPerson{tags=manager, level=7, active=true, status=active, age=45, salary=85000, department=marketing, name=Diana, score=7.8}
5. [1] TestPerson{salary=55000, active=false, department=sales, level=3, age=30, score=8, tags=employee, status=inactive, name=Eve}
6. [1] TestPerson{age=0, salary=-5000, active=true, name=Frank, score=0, status=active, level=1, tags=test, department=qa}
7. [1] TestPerson{active=true, tags=executive, score=10, age=65, status=active, department=management, level=9, name=Grace, salary=95000}
8. [1] TestPerson{age=18, department=support, level=1, salary=25000, tags=junior, name=Henry, active=false, score=5.5, status=inactive}
9. [1] TestPerson{age=40, status=active, salary=68000, score=8.7, tags=senior, name=Ivy, active=true, department=engineering, level=6}
10. [1] TestPerson{level=1, name=X, age=22, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000}
11. [1] TestOrder{priority=normal, status=pending, total=1999.98, region=north, customer_id=P001, discount=50, product_id=PROD001, amount=2, date=2024-01-15}
12. [1] TestOrder{customer_id=P002, date=2024-01-20, status=confirmed, product_id=PROD002, amount=1, discount=0, total=25.5, priority=low, region=south}
13. [1] TestOrder{discount=15, product_id=PROD003, region=north, total=225, date=2024-02-01, status=shipped, amount=3, priority=high, customer_id=P001}
14. [1] TestOrder{priority=normal, total=299.99, status=delivered, product_id=PROD004, discount=0, region=east, customer_id=P004, amount=1, date=2024-02-05}
15. [1] TestOrder{priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, region=south, customer_id=P002, product_id=PROD001, status=confirmed}
16. [1] TestOrder{region=west, customer_id=P005, date=2024-02-15, discount=0, product_id=PROD005, status=cancelled, priority=low, amount=2, total=999.98}
17. [1] TestOrder{customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}
18. [1] TestOrder{priority=normal, region=south, total=255, date=2024-03-05, amount=10, customer_id=P010, discount=0, product_id=PROD002, status=pending}
19. [1] TestOrder{customer_id=P001, date=2024-03-10, status=completed, discount=10, product_id=PROD007, region=north, total=89.99, priority=low, amount=1}
20. [1] TestOrder{customer_id=P006, product_id=PROD001, status=refunded, discount=0, region=east, amount=1, date=2024-03-15, total=75000, priority=urgent}

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

1. [1] TestPerson{score=8.5, name=Alice, age=25, active=true, department=sales, level=2, salary=45000, tags=junior, status=active}
2. [1] TestPerson{status=active, department=engineering, level=5, salary=75000, name=Bob, age=35, active=true, score=9.2, tags=senior}
3. [1] TestPerson{tags=intern, level=1, status=inactive, age=16, department=hr, active=false, score=6, name=Charlie, salary=0}
4. [1] TestPerson{status=active, age=45, salary=85000, department=marketing, name=Diana, score=7.8, tags=manager, level=7, active=true}
5. [1] TestPerson{status=inactive, name=Eve, salary=55000, active=false, department=sales, level=3, age=30, score=8, tags=employee}
6. [1] TestPerson{score=0, status=active, level=1, tags=test, department=qa, age=0, salary=-5000, active=true, name=Frank}
7. [1] TestPerson{level=9, name=Grace, salary=95000, active=true, tags=executive, score=10, age=65, status=active, department=management}
8. [1] TestPerson{status=inactive, age=18, department=support, level=1, salary=25000, tags=junior, name=Henry, active=false, score=5.5}
9. [1] TestPerson{tags=senior, name=Ivy, active=true, department=engineering, level=6, age=40, status=active, salary=68000, score=8.7}
10. [1] TestPerson{age=22, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1, name=X}

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

1. [1] TestOrder{product_id=PROD001, amount=2, date=2024-01-15, priority=normal, status=pending, total=1999.98, region=north, customer_id=P001, discount=50}
2. [1] TestOrder{amount=1, discount=0, total=25.5, priority=low, region=south, customer_id=P002, date=2024-01-20, status=confirmed, product_id=PROD002}
3. [1] TestOrder{date=2024-02-01, status=shipped, amount=3, priority=high, customer_id=P001, discount=15, product_id=PROD003, region=north, total=225}
4. [1] TestOrder{priority=normal, total=299.99, status=delivered, product_id=PROD004, discount=0, region=east, customer_id=P004, amount=1, date=2024-02-05}
5. [1] TestOrder{date=2024-02-10, region=south, customer_id=P002, product_id=PROD001, status=confirmed, priority=high, discount=100, amount=1, total=999.99}
6. [1] TestOrder{product_id=PROD005, status=cancelled, priority=low, amount=2, total=999.98, region=west, customer_id=P005, date=2024-02-15, discount=0}
7. [1] TestOrder{date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600}
8. [1] TestOrder{product_id=PROD002, status=pending, priority=normal, region=south, total=255, date=2024-03-05, amount=10, customer_id=P010, discount=0}
9. [1] TestOrder{customer_id=P001, date=2024-03-10, status=completed, discount=10, product_id=PROD007, region=north, total=89.99, priority=low, amount=1}
10. [1] TestOrder{date=2024-03-15, total=75000, priority=urgent, customer_id=P006, product_id=PROD001, status=refunded, discount=0, region=east, amount=1}

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

1. [1] TestPerson{department=sales, level=2, salary=45000, tags=junior, status=active, score=8.5, name=Alice, age=25, active=true}
2. [1] TestPerson{status=active, department=engineering, level=5, salary=75000, name=Bob, age=35, active=true, score=9.2, tags=senior}
3. [1] TestPerson{tags=intern, level=1, status=inactive, age=16, department=hr, active=false, score=6, name=Charlie, salary=0}
4. [1] TestPerson{name=Diana, score=7.8, tags=manager, level=7, active=true, status=active, age=45, salary=85000, department=marketing}
5. [1] TestPerson{name=Eve, salary=55000, active=false, department=sales, level=3, age=30, score=8, tags=employee, status=inactive}
6. [1] TestPerson{score=0, status=active, level=1, tags=test, department=qa, age=0, salary=-5000, active=true, name=Frank}
7. [1] TestPerson{name=Grace, salary=95000, active=true, tags=executive, score=10, age=65, status=active, department=management, level=9}
8. [1] TestPerson{salary=25000, tags=junior, name=Henry, active=false, score=5.5, status=inactive, age=18, department=support, level=1}
9. [1] TestPerson{age=40, status=active, salary=68000, score=8.7, tags=senior, name=Ivy, active=true, department=engineering, level=6}
10. [1] TestPerson{level=1, name=X, age=22, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000}

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
