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

1. [1] TestPerson{department=sales, name=Alice, tags=junior, level=2, age=25, active=true, salary=45000, score=8.5, status=active}
2. [1] TestPerson{salary=75000, active=true, tags=senior, status=active, name=Bob, age=35, score=9.2, department=engineering, level=5}
3. [1] TestPerson{active=false, status=inactive, name=Charlie, salary=0, department=hr, level=1, age=16, score=6, tags=intern}
4. [1] TestPerson{tags=manager, score=7.8, status=active, name=Diana, age=45, salary=85000, active=true, department=marketing, level=7}
5. [1] TestPerson{name=Eve, age=30, salary=55000, score=8, department=sales, active=false, status=inactive, level=3, tags=employee}
6. [1] TestPerson{salary=-5000, department=qa, level=1, name=Frank, active=true, score=0, tags=test, status=active, age=0}
7. [1] TestPerson{tags=executive, age=65, salary=95000, status=active, department=management, level=9, name=Grace, score=10, active=true}
8. [1] TestPerson{salary=25000, name=Henry, age=18, department=support, level=1, active=false, score=5.5, tags=junior, status=inactive}
9. [1] TestPerson{active=true, department=engineering, name=Ivy, score=8.7, tags=senior, status=active, age=40, level=6, salary=68000}
10. [1] TestPerson{age=22, level=1, name=X, salary=28000, score=6.5, department=intern, active=true, tags=temp, status=active}

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

1. [1] TestOrder{product_id=PROD001, status=pending, date=2024-01-15, discount=50, customer_id=P001, region=north, amount=2, total=1999.98, priority=normal}
2. [1] TestOrder{customer_id=P002, region=south, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, discount=0, amount=1, priority=low}
3. [1] TestOrder{amount=3, status=shipped, priority=high, customer_id=P001, total=225, date=2024-02-01, region=north, product_id=PROD003, discount=15}
4. [1] TestOrder{total=299.99, priority=normal, region=east, customer_id=P004, date=2024-02-05, discount=0, product_id=PROD004, status=delivered, amount=1}
5. [1] TestOrder{customer_id=P002, product_id=PROD001, priority=high, total=999.99, discount=100, region=south, amount=1, date=2024-02-10, status=confirmed}
6. [1] TestOrder{discount=0, total=999.98, status=cancelled, region=west, customer_id=P005, product_id=PROD005, amount=2, priority=low, date=2024-02-15}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50, date=2024-03-01, total=600, priority=urgent, region=north}
8. [1] TestOrder{discount=0, amount=10, total=255, customer_id=P010, priority=normal, product_id=PROD002, date=2024-03-05, status=pending, region=south}
9. [1] TestOrder{discount=10, product_id=PROD007, priority=low, customer_id=P001, date=2024-03-10, total=89.99, status=completed, region=north, amount=1}
10. [1] TestOrder{discount=0, customer_id=P006, amount=1, priority=urgent, date=2024-03-15, status=refunded, region=east, product_id=PROD001, total=75000}

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

1. [1] TestPerson{active=true, salary=45000, score=8.5, status=active, department=sales, name=Alice, tags=junior, level=2, age=25}
2. [1] TestPerson{age=35, score=9.2, department=engineering, level=5, salary=75000, active=true, tags=senior, status=active, name=Bob}
3. [1] TestPerson{salary=0, department=hr, level=1, age=16, score=6, tags=intern, active=false, status=inactive, name=Charlie}
4. [1] TestPerson{department=marketing, level=7, tags=manager, score=7.8, status=active, name=Diana, age=45, salary=85000, active=true}
5. [1] TestPerson{age=30, salary=55000, score=8, department=sales, active=false, status=inactive, level=3, tags=employee, name=Eve}
6. [1] TestPerson{salary=-5000, department=qa, level=1, name=Frank, active=true, score=0, tags=test, status=active, age=0}
7. [1] TestPerson{status=active, department=management, level=9, name=Grace, score=10, active=true, tags=executive, age=65, salary=95000}
8. [1] TestPerson{level=1, active=false, score=5.5, tags=junior, status=inactive, salary=25000, name=Henry, age=18, department=support}
9. [1] TestPerson{age=40, level=6, salary=68000, active=true, department=engineering, name=Ivy, score=8.7, tags=senior, status=active}
10. [1] TestPerson{active=true, tags=temp, status=active, age=22, level=1, name=X, salary=28000, score=6.5, department=intern}

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

1. [1] TestOrder{priority=normal, product_id=PROD001, status=pending, date=2024-01-15, discount=50, customer_id=P001, region=north, amount=2, total=1999.98}
2. [1] TestOrder{total=25.5, date=2024-01-20, status=confirmed, discount=0, amount=1, priority=low, customer_id=P002, region=south, product_id=PROD002}
3. [1] TestOrder{amount=3, status=shipped, priority=high, customer_id=P001, total=225, date=2024-02-01, region=north, product_id=PROD003, discount=15}
4. [1] TestOrder{customer_id=P004, date=2024-02-05, discount=0, product_id=PROD004, status=delivered, amount=1, total=299.99, priority=normal, region=east}
5. [1] TestOrder{date=2024-02-10, status=confirmed, customer_id=P002, product_id=PROD001, priority=high, total=999.99, discount=100, region=south, amount=1}
6. [1] TestOrder{total=999.98, status=cancelled, region=west, customer_id=P005, product_id=PROD005, amount=2, priority=low, date=2024-02-15, discount=0}
7. [1] TestOrder{total=600, priority=urgent, region=north, customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50, date=2024-03-01}
8. [1] TestOrder{discount=0, amount=10, total=255, customer_id=P010, priority=normal, product_id=PROD002, date=2024-03-05, status=pending, region=south}
9. [1] TestOrder{region=north, amount=1, discount=10, product_id=PROD007, priority=low, customer_id=P001, date=2024-03-10, total=89.99, status=completed}
10. [1] TestOrder{product_id=PROD001, total=75000, discount=0, customer_id=P006, amount=1, priority=urgent, date=2024-03-15, status=refunded, region=east}

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

1. [1] TestProduct{keywords=computer, brand=TechCorp, stock=50, category=electronics, price=999.99, name=Laptop, available=true, rating=4.5, supplier=TechSupply}
2. [1] TestProduct{stock=200, rating=4, brand=TechCorp, name=Mouse, price=25.5, available=true, keywords=peripheral, supplier=TechSupply, category=accessories}
3. [1] TestProduct{available=false, stock=0, category=accessories, rating=3.5, keywords=typing, brand=KeyTech, supplier=KeySupply, name=Keyboard, price=75}
4. [1] TestProduct{keywords=display, supplier=ScreenSupply, available=true, brand=ScreenPro, stock=30, category=electronics, rating=4.8, name=Monitor, price=299.99}
5. [1] TestProduct{rating=2, supplier=OldSupply, available=false, keywords=obsolete, brand=OldTech, category=accessories, stock=0, name=OldKeyboard, price=8.5}
6. [1] TestProduct{brand=AudioMax, stock=75, name=Headphones, rating=4.6, keywords=sound, supplier=AudioSupply, category=audio, price=150, available=true}
7. [1] TestProduct{brand=CamTech, stock=25, name=Webcam, category=electronics, supplier=CamSupply, price=89.99, available=true, rating=3.8, keywords=video}

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

1. [1] TestPerson{age=25, active=true, salary=45000, score=8.5, status=active, department=sales, name=Alice, tags=junior, level=2}
2. [1] TestPerson{salary=75000, active=true, tags=senior, status=active, name=Bob, age=35, score=9.2, department=engineering, level=5}
3. [1] TestPerson{active=false, status=inactive, name=Charlie, salary=0, department=hr, level=1, age=16, score=6, tags=intern}
4. [1] TestPerson{score=7.8, status=active, name=Diana, age=45, salary=85000, active=true, department=marketing, level=7, tags=manager}
5. [1] TestPerson{age=30, salary=55000, score=8, department=sales, active=false, status=inactive, level=3, tags=employee, name=Eve}
6. [1] TestPerson{salary=-5000, department=qa, level=1, name=Frank, active=true, score=0, tags=test, status=active, age=0}
7. [1] TestPerson{age=65, salary=95000, status=active, department=management, level=9, name=Grace, score=10, active=true, tags=executive}
8. [1] TestPerson{name=Henry, age=18, department=support, level=1, active=false, score=5.5, tags=junior, status=inactive, salary=25000}
9. [1] TestPerson{department=engineering, name=Ivy, score=8.7, tags=senior, status=active, age=40, level=6, salary=68000, active=true}
10. [1] TestPerson{age=22, level=1, name=X, salary=28000, score=6.5, department=intern, active=true, tags=temp, status=active}

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

1. [1] TestOrder{date=2024-01-15, discount=50, customer_id=P001, region=north, amount=2, total=1999.98, priority=normal, product_id=PROD001, status=pending}
2. [1] TestOrder{customer_id=P002, region=south, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, discount=0, amount=1, priority=low}
3. [1] TestOrder{customer_id=P001, total=225, date=2024-02-01, region=north, product_id=PROD003, discount=15, amount=3, status=shipped, priority=high}
4. [1] TestOrder{total=299.99, priority=normal, region=east, customer_id=P004, date=2024-02-05, discount=0, product_id=PROD004, status=delivered, amount=1}
5. [1] TestOrder{date=2024-02-10, status=confirmed, customer_id=P002, product_id=PROD001, priority=high, total=999.99, discount=100, region=south, amount=1}
6. [1] TestOrder{total=999.98, status=cancelled, region=west, customer_id=P005, product_id=PROD005, amount=2, priority=low, date=2024-02-15, discount=0}
7. [1] TestOrder{date=2024-03-01, total=600, priority=urgent, region=north, customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50}
8. [1] TestOrder{discount=0, amount=10, total=255, customer_id=P010, priority=normal, product_id=PROD002, date=2024-03-05, status=pending, region=south}
9. [1] TestOrder{priority=low, customer_id=P001, date=2024-03-10, total=89.99, status=completed, region=north, amount=1, discount=10, product_id=PROD007}
10. [1] TestOrder{total=75000, discount=0, customer_id=P006, amount=1, priority=urgent, date=2024-03-15, status=refunded, region=east, product_id=PROD001}

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

1. [1] TestPerson{age=25, active=true, salary=45000, score=8.5, status=active, department=sales, name=Alice, tags=junior, level=2}
2. [1] TestPerson{score=9.2, department=engineering, level=5, salary=75000, active=true, tags=senior, status=active, name=Bob, age=35}
3. [1] TestPerson{department=hr, level=1, age=16, score=6, tags=intern, active=false, status=inactive, name=Charlie, salary=0}
4. [1] TestPerson{department=marketing, level=7, tags=manager, score=7.8, status=active, name=Diana, age=45, salary=85000, active=true}
5. [1] TestPerson{department=sales, active=false, status=inactive, level=3, tags=employee, name=Eve, age=30, salary=55000, score=8}
6. [1] TestPerson{salary=-5000, department=qa, level=1, name=Frank, active=true, score=0, tags=test, status=active, age=0}
7. [1] TestPerson{age=65, salary=95000, status=active, department=management, level=9, name=Grace, score=10, active=true, tags=executive}
8. [1] TestPerson{age=18, department=support, level=1, active=false, score=5.5, tags=junior, status=inactive, salary=25000, name=Henry}
9. [1] TestPerson{age=40, level=6, salary=68000, active=true, department=engineering, name=Ivy, score=8.7, tags=senior, status=active}
10. [1] TestPerson{tags=temp, status=active, age=22, level=1, name=X, salary=28000, score=6.5, department=intern, active=true}

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

1. [1] TestOrder{discount=50, customer_id=P001, region=north, amount=2, total=1999.98, priority=normal, product_id=PROD001, status=pending, date=2024-01-15}
2. [1] TestOrder{total=25.5, date=2024-01-20, status=confirmed, discount=0, amount=1, priority=low, customer_id=P002, region=south, product_id=PROD002}
3. [1] TestOrder{region=north, product_id=PROD003, discount=15, amount=3, status=shipped, priority=high, customer_id=P001, total=225, date=2024-02-01}
4. [1] TestOrder{total=299.99, priority=normal, region=east, customer_id=P004, date=2024-02-05, discount=0, product_id=PROD004, status=delivered, amount=1}
5. [1] TestOrder{priority=high, total=999.99, discount=100, region=south, amount=1, date=2024-02-10, status=confirmed, customer_id=P002, product_id=PROD001}
6. [1] TestOrder{amount=2, priority=low, date=2024-02-15, discount=0, total=999.98, status=cancelled, region=west, customer_id=P005, product_id=PROD005}
7. [1] TestOrder{discount=50, date=2024-03-01, total=600, priority=urgent, region=north, customer_id=P007, product_id=PROD006, amount=4, status=shipped}
8. [1] TestOrder{customer_id=P010, priority=normal, product_id=PROD002, date=2024-03-05, status=pending, region=south, discount=0, amount=10, total=255}
9. [1] TestOrder{total=89.99, status=completed, region=north, amount=1, discount=10, product_id=PROD007, priority=low, customer_id=P001, date=2024-03-10}
10. [1] TestOrder{priority=urgent, date=2024-03-15, status=refunded, region=east, product_id=PROD001, total=75000, discount=0, customer_id=P006, amount=1}

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

1. [1] TestPerson{active=true, salary=45000, score=8.5, status=active, department=sales, name=Alice, tags=junior, level=2, age=25}
2. [1] TestPerson{salary=75000, active=true, tags=senior, status=active, name=Bob, age=35, score=9.2, department=engineering, level=5}
3. [1] TestPerson{status=inactive, name=Charlie, salary=0, department=hr, level=1, age=16, score=6, tags=intern, active=false}
4. [1] TestPerson{department=marketing, level=7, tags=manager, score=7.8, status=active, name=Diana, age=45, salary=85000, active=true}
5. [1] TestPerson{score=8, department=sales, active=false, status=inactive, level=3, tags=employee, name=Eve, age=30, salary=55000}
6. [1] TestPerson{salary=-5000, department=qa, level=1, name=Frank, active=true, score=0, tags=test, status=active, age=0}
7. [1] TestPerson{score=10, active=true, tags=executive, age=65, salary=95000, status=active, department=management, level=9, name=Grace}
8. [1] TestPerson{name=Henry, age=18, department=support, level=1, active=false, score=5.5, tags=junior, status=inactive, salary=25000}
9. [1] TestPerson{score=8.7, tags=senior, status=active, age=40, level=6, salary=68000, active=true, department=engineering, name=Ivy}
10. [1] TestPerson{tags=temp, status=active, age=22, level=1, name=X, salary=28000, score=6.5, department=intern, active=true}

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

1. [1] TestProduct{rating=4.5, supplier=TechSupply, keywords=computer, brand=TechCorp, stock=50, category=electronics, price=999.99, name=Laptop, available=true}
2. [1] TestProduct{stock=200, rating=4, brand=TechCorp, name=Mouse, price=25.5, available=true, keywords=peripheral, supplier=TechSupply, category=accessories}
3. [1] TestProduct{brand=KeyTech, supplier=KeySupply, name=Keyboard, price=75, available=false, stock=0, category=accessories, rating=3.5, keywords=typing}
4. [1] TestProduct{available=true, brand=ScreenPro, stock=30, category=electronics, rating=4.8, name=Monitor, price=299.99, keywords=display, supplier=ScreenSupply}
5. [1] TestProduct{keywords=obsolete, brand=OldTech, category=accessories, stock=0, name=OldKeyboard, price=8.5, rating=2, supplier=OldSupply, available=false}
6. [1] TestProduct{supplier=AudioSupply, category=audio, price=150, available=true, brand=AudioMax, stock=75, name=Headphones, rating=4.6, keywords=sound}
7. [1] TestProduct{price=89.99, available=true, rating=3.8, keywords=video, brand=CamTech, stock=25, name=Webcam, category=electronics, supplier=CamSupply}

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

1. [1] TestPerson{status=active, department=sales, name=Alice, tags=junior, level=2, age=25, active=true, salary=45000, score=8.5}
2. [1] TestPerson{age=35, score=9.2, department=engineering, level=5, salary=75000, active=true, tags=senior, status=active, name=Bob}
3. [1] TestPerson{department=hr, level=1, age=16, score=6, tags=intern, active=false, status=inactive, name=Charlie, salary=0}
4. [1] TestPerson{name=Diana, age=45, salary=85000, active=true, department=marketing, level=7, tags=manager, score=7.8, status=active}
5. [1] TestPerson{active=false, status=inactive, level=3, tags=employee, name=Eve, age=30, salary=55000, score=8, department=sales}
6. [1] TestPerson{salary=-5000, department=qa, level=1, name=Frank, active=true, score=0, tags=test, status=active, age=0}
7. [1] TestPerson{status=active, department=management, level=9, name=Grace, score=10, active=true, tags=executive, age=65, salary=95000}
8. [1] TestPerson{age=18, department=support, level=1, active=false, score=5.5, tags=junior, status=inactive, salary=25000, name=Henry}
9. [1] TestPerson{age=40, level=6, salary=68000, active=true, department=engineering, name=Ivy, score=8.7, tags=senior, status=active}
10. [1] TestPerson{department=intern, active=true, tags=temp, status=active, age=22, level=1, name=X, salary=28000, score=6.5}

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

1. [1] TestOrder{customer_id=P001, region=north, amount=2, total=1999.98, priority=normal, product_id=PROD001, status=pending, date=2024-01-15, discount=50}
2. [1] TestOrder{amount=1, priority=low, customer_id=P002, region=south, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, discount=0}
3. [1] TestOrder{discount=15, amount=3, status=shipped, priority=high, customer_id=P001, total=225, date=2024-02-01, region=north, product_id=PROD003}
4. [1] TestOrder{product_id=PROD004, status=delivered, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, date=2024-02-05, discount=0}
5. [1] TestOrder{region=south, amount=1, date=2024-02-10, status=confirmed, customer_id=P002, product_id=PROD001, priority=high, total=999.99, discount=100}
6. [1] TestOrder{amount=2, priority=low, date=2024-02-15, discount=0, total=999.98, status=cancelled, region=west, customer_id=P005, product_id=PROD005}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50, date=2024-03-01, total=600, priority=urgent, region=north}
8. [1] TestOrder{product_id=PROD002, date=2024-03-05, status=pending, region=south, discount=0, amount=10, total=255, customer_id=P010, priority=normal}
9. [1] TestOrder{status=completed, region=north, amount=1, discount=10, product_id=PROD007, priority=low, customer_id=P001, date=2024-03-10, total=89.99}
10. [1] TestOrder{priority=urgent, date=2024-03-15, status=refunded, region=east, product_id=PROD001, total=75000, discount=0, customer_id=P006, amount=1}

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

1. [1] TestPerson{score=8.5, status=active, department=sales, name=Alice, tags=junior, level=2, age=25, active=true, salary=45000}
2. [1] TestPerson{score=9.2, department=engineering, level=5, salary=75000, active=true, tags=senior, status=active, name=Bob, age=35}
3. [1] TestPerson{department=hr, level=1, age=16, score=6, tags=intern, active=false, status=inactive, name=Charlie, salary=0}
4. [1] TestPerson{score=7.8, status=active, name=Diana, age=45, salary=85000, active=true, department=marketing, level=7, tags=manager}
5. [1] TestPerson{age=30, salary=55000, score=8, department=sales, active=false, status=inactive, level=3, tags=employee, name=Eve}
6. [1] TestPerson{score=0, tags=test, status=active, age=0, salary=-5000, department=qa, level=1, name=Frank, active=true}
7. [1] TestPerson{age=65, salary=95000, status=active, department=management, level=9, name=Grace, score=10, active=true, tags=executive}
8. [1] TestPerson{department=support, level=1, active=false, score=5.5, tags=junior, status=inactive, salary=25000, name=Henry, age=18}
9. [1] TestPerson{name=Ivy, score=8.7, tags=senior, status=active, age=40, level=6, salary=68000, active=true, department=engineering}
10. [1] TestPerson{name=X, salary=28000, score=6.5, department=intern, active=true, tags=temp, status=active, age=22, level=1}
11. [1] TestOrder{discount=50, customer_id=P001, region=north, amount=2, total=1999.98, priority=normal, product_id=PROD001, status=pending, date=2024-01-15}
12. [1] TestOrder{total=25.5, date=2024-01-20, status=confirmed, discount=0, amount=1, priority=low, customer_id=P002, region=south, product_id=PROD002}
13. [1] TestOrder{amount=3, status=shipped, priority=high, customer_id=P001, total=225, date=2024-02-01, region=north, product_id=PROD003, discount=15}
14. [1] TestOrder{product_id=PROD004, status=delivered, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, date=2024-02-05, discount=0}
15. [1] TestOrder{amount=1, date=2024-02-10, status=confirmed, customer_id=P002, product_id=PROD001, priority=high, total=999.99, discount=100, region=south}
16. [1] TestOrder{region=west, customer_id=P005, product_id=PROD005, amount=2, priority=low, date=2024-02-15, discount=0, total=999.98, status=cancelled}
17. [1] TestOrder{discount=50, date=2024-03-01, total=600, priority=urgent, region=north, customer_id=P007, product_id=PROD006, amount=4, status=shipped}
18. [1] TestOrder{customer_id=P010, priority=normal, product_id=PROD002, date=2024-03-05, status=pending, region=south, discount=0, amount=10, total=255}
19. [1] TestOrder{customer_id=P001, date=2024-03-10, total=89.99, status=completed, region=north, amount=1, discount=10, product_id=PROD007, priority=low}
20. [1] TestOrder{priority=urgent, date=2024-03-15, status=refunded, region=east, product_id=PROD001, total=75000, discount=0, customer_id=P006, amount=1}

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

1. [1] TestPerson{age=25, active=true, salary=45000, score=8.5, status=active, department=sales, name=Alice, tags=junior, level=2}
2. [1] TestPerson{level=5, salary=75000, active=true, tags=senior, status=active, name=Bob, age=35, score=9.2, department=engineering}
3. [1] TestPerson{department=hr, level=1, age=16, score=6, tags=intern, active=false, status=inactive, name=Charlie, salary=0}
4. [1] TestPerson{age=45, salary=85000, active=true, department=marketing, level=7, tags=manager, score=7.8, status=active, name=Diana}
5. [1] TestPerson{level=3, tags=employee, name=Eve, age=30, salary=55000, score=8, department=sales, active=false, status=inactive}
6. [1] TestPerson{salary=-5000, department=qa, level=1, name=Frank, active=true, score=0, tags=test, status=active, age=0}
7. [1] TestPerson{name=Grace, score=10, active=true, tags=executive, age=65, salary=95000, status=active, department=management, level=9}
8. [1] TestPerson{salary=25000, name=Henry, age=18, department=support, level=1, active=false, score=5.5, tags=junior, status=inactive}
9. [1] TestPerson{active=true, department=engineering, name=Ivy, score=8.7, tags=senior, status=active, age=40, level=6, salary=68000}
10. [1] TestPerson{score=6.5, department=intern, active=true, tags=temp, status=active, age=22, level=1, name=X, salary=28000}

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

1. [1] TestPerson{age=25, active=true, salary=45000, score=8.5, status=active, department=sales, name=Alice, tags=junior, level=2}
2. [1] TestPerson{salary=75000, active=true, tags=senior, status=active, name=Bob, age=35, score=9.2, department=engineering, level=5}
3. [1] TestPerson{department=hr, level=1, age=16, score=6, tags=intern, active=false, status=inactive, name=Charlie, salary=0}
4. [1] TestPerson{score=7.8, status=active, name=Diana, age=45, salary=85000, active=true, department=marketing, level=7, tags=manager}
5. [1] TestPerson{department=sales, active=false, status=inactive, level=3, tags=employee, name=Eve, age=30, salary=55000, score=8}
6. [1] TestPerson{age=0, salary=-5000, department=qa, level=1, name=Frank, active=true, score=0, tags=test, status=active}
7. [1] TestPerson{status=active, department=management, level=9, name=Grace, score=10, active=true, tags=executive, age=65, salary=95000}
8. [1] TestPerson{level=1, active=false, score=5.5, tags=junior, status=inactive, salary=25000, name=Henry, age=18, department=support}
9. [1] TestPerson{tags=senior, status=active, age=40, level=6, salary=68000, active=true, department=engineering, name=Ivy, score=8.7}
10. [1] TestPerson{name=X, salary=28000, score=6.5, department=intern, active=true, tags=temp, status=active, age=22, level=1}
11. [1] TestOrder{status=pending, date=2024-01-15, discount=50, customer_id=P001, region=north, amount=2, total=1999.98, priority=normal, product_id=PROD001}
12. [1] TestOrder{customer_id=P002, region=south, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, discount=0, amount=1, priority=low}
13. [1] TestOrder{date=2024-02-01, region=north, product_id=PROD003, discount=15, amount=3, status=shipped, priority=high, customer_id=P001, total=225}
14. [1] TestOrder{total=299.99, priority=normal, region=east, customer_id=P004, date=2024-02-05, discount=0, product_id=PROD004, status=delivered, amount=1}
15. [1] TestOrder{customer_id=P002, product_id=PROD001, priority=high, total=999.99, discount=100, region=south, amount=1, date=2024-02-10, status=confirmed}
16. [1] TestOrder{amount=2, priority=low, date=2024-02-15, discount=0, total=999.98, status=cancelled, region=west, customer_id=P005, product_id=PROD005}
17. [1] TestOrder{total=600, priority=urgent, region=north, customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50, date=2024-03-01}
18. [1] TestOrder{product_id=PROD002, date=2024-03-05, status=pending, region=south, discount=0, amount=10, total=255, customer_id=P010, priority=normal}
19. [1] TestOrder{date=2024-03-10, total=89.99, status=completed, region=north, amount=1, discount=10, product_id=PROD007, priority=low, customer_id=P001}
20. [1] TestOrder{product_id=PROD001, total=75000, discount=0, customer_id=P006, amount=1, priority=urgent, date=2024-03-15, status=refunded, region=east}

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

1. [1] TestPerson{age=25, active=true, salary=45000, score=8.5, status=active, department=sales, name=Alice, tags=junior, level=2}
2. [1] TestPerson{age=35, score=9.2, department=engineering, level=5, salary=75000, active=true, tags=senior, status=active, name=Bob}
3. [1] TestPerson{department=hr, level=1, age=16, score=6, tags=intern, active=false, status=inactive, name=Charlie, salary=0}
4. [1] TestPerson{name=Diana, age=45, salary=85000, active=true, department=marketing, level=7, tags=manager, score=7.8, status=active}
5. [1] TestPerson{age=30, salary=55000, score=8, department=sales, active=false, status=inactive, level=3, tags=employee, name=Eve}
6. [1] TestPerson{score=0, tags=test, status=active, age=0, salary=-5000, department=qa, level=1, name=Frank, active=true}
7. [1] TestPerson{score=10, active=true, tags=executive, age=65, salary=95000, status=active, department=management, level=9, name=Grace}
8. [1] TestPerson{salary=25000, name=Henry, age=18, department=support, level=1, active=false, score=5.5, tags=junior, status=inactive}
9. [1] TestPerson{score=8.7, tags=senior, status=active, age=40, level=6, salary=68000, active=true, department=engineering, name=Ivy}
10. [1] TestPerson{name=X, salary=28000, score=6.5, department=intern, active=true, tags=temp, status=active, age=22, level=1}

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

1. [1] TestOrder{date=2024-01-15, discount=50, customer_id=P001, region=north, amount=2, total=1999.98, priority=normal, product_id=PROD001, status=pending}
2. [1] TestOrder{date=2024-01-20, status=confirmed, discount=0, amount=1, priority=low, customer_id=P002, region=south, product_id=PROD002, total=25.5}
3. [1] TestOrder{date=2024-02-01, region=north, product_id=PROD003, discount=15, amount=3, status=shipped, priority=high, customer_id=P001, total=225}
4. [1] TestOrder{product_id=PROD004, status=delivered, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, date=2024-02-05, discount=0}
5. [1] TestOrder{priority=high, total=999.99, discount=100, region=south, amount=1, date=2024-02-10, status=confirmed, customer_id=P002, product_id=PROD001}
6. [1] TestOrder{priority=low, date=2024-02-15, discount=0, total=999.98, status=cancelled, region=west, customer_id=P005, product_id=PROD005, amount=2}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50, date=2024-03-01, total=600, priority=urgent, region=north}
8. [1] TestOrder{date=2024-03-05, status=pending, region=south, discount=0, amount=10, total=255, customer_id=P010, priority=normal, product_id=PROD002}
9. [1] TestOrder{status=completed, region=north, amount=1, discount=10, product_id=PROD007, priority=low, customer_id=P001, date=2024-03-10, total=89.99}
10. [1] TestOrder{region=east, product_id=PROD001, total=75000, discount=0, customer_id=P006, amount=1, priority=urgent, date=2024-03-15, status=refunded}

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

1. [1] TestPerson{age=25, active=true, salary=45000, score=8.5, status=active, department=sales, name=Alice, tags=junior, level=2}
2. [1] TestPerson{salary=75000, active=true, tags=senior, status=active, name=Bob, age=35, score=9.2, department=engineering, level=5}
3. [1] TestPerson{level=1, age=16, score=6, tags=intern, active=false, status=inactive, name=Charlie, salary=0, department=hr}
4. [1] TestPerson{score=7.8, status=active, name=Diana, age=45, salary=85000, active=true, department=marketing, level=7, tags=manager}
5. [1] TestPerson{age=30, salary=55000, score=8, department=sales, active=false, status=inactive, level=3, tags=employee, name=Eve}
6. [1] TestPerson{salary=-5000, department=qa, level=1, name=Frank, active=true, score=0, tags=test, status=active, age=0}
7. [1] TestPerson{salary=95000, status=active, department=management, level=9, name=Grace, score=10, active=true, tags=executive, age=65}
8. [1] TestPerson{level=1, active=false, score=5.5, tags=junior, status=inactive, salary=25000, name=Henry, age=18, department=support}
9. [1] TestPerson{score=8.7, tags=senior, status=active, age=40, level=6, salary=68000, active=true, department=engineering, name=Ivy}
10. [1] TestPerson{name=X, salary=28000, score=6.5, department=intern, active=true, tags=temp, status=active, age=22, level=1}

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
