# RÉSULTATS COMPLETS - ANALYSE RÈGLES DE NÉGATION TSD
=====================================================

**Date d'exécution**: 13 novembre 2025
**Fichier contraintes**: /home/resinsec/dev/tsd/constraint/test/integration/negation_rules.tsd
**Nombre de règles**: 19
**Nombre de faits**: 27

## 🎯 RÈGLE 0: not_zero_age

**Condition**: `NOT (p.age == 0)`
**Types concernés**: [TestPerson]
**Terminal**: rule_0_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{tags=junior, status=active, department=sales, salary=45000, active=true, name=Alice, age=25, score=8.5, level=2}
2. [1] TestPerson{salary=75000, active=true, score=9.2, tags=senior, department=engineering, age=35, status=active, level=5, name=Bob}
3. [1] TestPerson{active=false, tags=intern, status=inactive, department=hr, level=1, name=Charlie, score=6, age=16, salary=0}
4. [1] TestPerson{salary=85000, active=true, score=7.8, tags=manager, department=marketing, level=7, age=45, name=Diana, status=active}
5. [1] TestPerson{status=inactive, tags=employee, salary=55000, department=sales, name=Eve, level=3, age=30, active=false, score=8}
6. [1] TestPerson{age=0, salary=-5000, score=0, tags=test, name=Frank, status=active, department=qa, level=1, active=true}
7. [1] TestPerson{age=65, score=10, level=9, salary=95000, tags=executive, status=active, department=management, active=true, name=Grace}
8. [1] TestPerson{tags=junior, department=support, name=Henry, active=false, status=inactive, salary=25000, level=1, age=18, score=5.5}
9. [1] TestPerson{level=6, name=Ivy, age=40, status=active, department=engineering, active=true, score=8.7, tags=senior, salary=68000}
10. [1] TestPerson{status=active, level=1, active=true, department=intern, age=22, score=6.5, name=X, salary=28000, tags=temp}

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

1. [1] TestOrder{priority=normal, discount=50, date=2024-01-15, region=north, product_id=PROD001, status=pending, amount=2, total=1999.98, customer_id=P001}
2. [1] TestOrder{date=2024-01-20, discount=0, region=south, amount=1, priority=low, product_id=PROD002, customer_id=P002, total=25.5, status=confirmed}
3. [1] TestOrder{date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, status=shipped, total=225, region=north, amount=3, priority=high}
4. [1] TestOrder{product_id=PROD004, date=2024-02-05, discount=0, customer_id=P004, amount=1, status=delivered, total=299.99, priority=normal, region=east}
5. [1] TestOrder{status=confirmed, product_id=PROD001, region=south, customer_id=P002, amount=1, total=999.99, priority=high, discount=100, date=2024-02-10}
6. [1] TestOrder{product_id=PROD005, amount=2, date=2024-02-15, status=cancelled, customer_id=P005, total=999.98, priority=low, region=west, discount=0}
7. [1] TestOrder{priority=urgent, discount=50, region=north, status=shipped, amount=4, date=2024-03-01, customer_id=P007, product_id=PROD006, total=600}
8. [1] TestOrder{status=pending, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, priority=normal, region=south}
9. [1] TestOrder{status=completed, priority=low, region=north, amount=1, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, total=89.99}
10. [1] TestOrder{amount=1, total=75000, date=2024-03-15, discount=0, status=refunded, priority=urgent, customer_id=P006, product_id=PROD001, region=east}

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

1. [1] TestPerson{tags=junior, status=active, department=sales, salary=45000, active=true, name=Alice, age=25, score=8.5, level=2}
2. [1] TestPerson{name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, age=35, status=active, level=5}
3. [1] TestPerson{score=6, age=16, salary=0, active=false, tags=intern, status=inactive, department=hr, level=1, name=Charlie}
4. [1] TestPerson{department=marketing, level=7, age=45, name=Diana, status=active, salary=85000, active=true, score=7.8, tags=manager}
5. [1] TestPerson{status=inactive, tags=employee, salary=55000, department=sales, name=Eve, level=3, age=30, active=false, score=8}
6. [1] TestPerson{age=0, salary=-5000, score=0, tags=test, name=Frank, status=active, department=qa, level=1, active=true}
7. [1] TestPerson{salary=95000, tags=executive, status=active, department=management, active=true, name=Grace, age=65, score=10, level=9}
8. [1] TestPerson{age=18, score=5.5, tags=junior, department=support, name=Henry, active=false, status=inactive, salary=25000, level=1}
9. [1] TestPerson{status=active, department=engineering, active=true, score=8.7, tags=senior, salary=68000, level=6, name=Ivy, age=40}
10. [1] TestPerson{age=22, score=6.5, name=X, salary=28000, tags=temp, status=active, level=1, active=true, department=intern}

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

1. [1] TestOrder{product_id=PROD001, status=pending, amount=2, total=1999.98, customer_id=P001, priority=normal, discount=50, date=2024-01-15, region=north}
2. [1] TestOrder{priority=low, product_id=PROD002, customer_id=P002, total=25.5, status=confirmed, date=2024-01-20, discount=0, region=south, amount=1}
3. [1] TestOrder{date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, status=shipped, total=225, region=north, amount=3, priority=high}
4. [1] TestOrder{customer_id=P004, amount=1, status=delivered, total=299.99, priority=normal, region=east, product_id=PROD004, date=2024-02-05, discount=0}
5. [1] TestOrder{customer_id=P002, amount=1, total=999.99, priority=high, discount=100, date=2024-02-10, status=confirmed, product_id=PROD001, region=south}
6. [1] TestOrder{amount=2, date=2024-02-15, status=cancelled, customer_id=P005, total=999.98, priority=low, region=west, discount=0, product_id=PROD005}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, total=600, priority=urgent, discount=50, region=north, status=shipped, amount=4, date=2024-03-01}
8. [1] TestOrder{amount=10, total=255, date=2024-03-05, discount=0, priority=normal, region=south, status=pending, customer_id=P010, product_id=PROD002}
9. [1] TestOrder{status=completed, priority=low, region=north, amount=1, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, total=89.99}
10. [1] TestOrder{amount=1, total=75000, date=2024-03-15, discount=0, status=refunded, priority=urgent, customer_id=P006, product_id=PROD001, region=east}

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

1. [1] TestProduct{brand=TechCorp, name=Laptop, category=electronics, supplier=TechSupply, stock=50, price=999.99, available=true, rating=4.5, keywords=computer}
2. [1] TestProduct{supplier=TechSupply, rating=4, keywords=peripheral, price=25.5, stock=200, category=accessories, available=true, name=Mouse, brand=TechCorp}
3. [1] TestProduct{category=accessories, price=75, supplier=KeySupply, name=Keyboard, available=false, keywords=typing, brand=KeyTech, rating=3.5, stock=0}
4. [1] TestProduct{rating=4.8, keywords=display, category=electronics, stock=30, name=Monitor, price=299.99, available=true, brand=ScreenPro, supplier=ScreenSupply}
5. [1] TestProduct{rating=2, keywords=obsolete, stock=0, supplier=OldSupply, category=accessories, price=8.5, brand=OldTech, name=OldKeyboard, available=false}
6. [1] TestProduct{rating=4.6, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply, name=Headphones, category=audio, price=150, available=true}
7. [1] TestProduct{supplier=CamSupply, name=Webcam, category=electronics, rating=3.8, available=true, brand=CamTech, price=89.99, keywords=video, stock=25}

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

1. [1] TestPerson{level=2, tags=junior, status=active, department=sales, salary=45000, active=true, name=Alice, age=25, score=8.5}
2. [1] TestPerson{level=5, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, age=35, status=active}
3. [1] TestPerson{score=6, age=16, salary=0, active=false, tags=intern, status=inactive, department=hr, level=1, name=Charlie}
4. [1] TestPerson{score=7.8, tags=manager, department=marketing, level=7, age=45, name=Diana, status=active, salary=85000, active=true}
5. [1] TestPerson{name=Eve, level=3, age=30, active=false, score=8, status=inactive, tags=employee, salary=55000, department=sales}
6. [1] TestPerson{age=0, salary=-5000, score=0, tags=test, name=Frank, status=active, department=qa, level=1, active=true}
7. [1] TestPerson{score=10, level=9, salary=95000, tags=executive, status=active, department=management, active=true, name=Grace, age=65}
8. [1] TestPerson{age=18, score=5.5, tags=junior, department=support, name=Henry, active=false, status=inactive, salary=25000, level=1}
9. [1] TestPerson{active=true, score=8.7, tags=senior, salary=68000, level=6, name=Ivy, age=40, status=active, department=engineering}
10. [1] TestPerson{tags=temp, status=active, level=1, active=true, department=intern, age=22, score=6.5, name=X, salary=28000}

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

1. [1] TestOrder{region=north, product_id=PROD001, status=pending, amount=2, total=1999.98, customer_id=P001, priority=normal, discount=50, date=2024-01-15}
2. [1] TestOrder{total=25.5, status=confirmed, date=2024-01-20, discount=0, region=south, amount=1, priority=low, product_id=PROD002, customer_id=P002}
3. [1] TestOrder{date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, status=shipped, total=225, region=north, amount=3, priority=high}
4. [1] TestOrder{status=delivered, total=299.99, priority=normal, region=east, product_id=PROD004, date=2024-02-05, discount=0, customer_id=P004, amount=1}
5. [1] TestOrder{status=confirmed, product_id=PROD001, region=south, customer_id=P002, amount=1, total=999.99, priority=high, discount=100, date=2024-02-10}
6. [1] TestOrder{product_id=PROD005, amount=2, date=2024-02-15, status=cancelled, customer_id=P005, total=999.98, priority=low, region=west, discount=0}
7. [1] TestOrder{total=600, priority=urgent, discount=50, region=north, status=shipped, amount=4, date=2024-03-01, customer_id=P007, product_id=PROD006}
8. [1] TestOrder{amount=10, total=255, date=2024-03-05, discount=0, priority=normal, region=south, status=pending, customer_id=P010, product_id=PROD002}
9. [1] TestOrder{status=completed, priority=low, region=north, amount=1, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, total=89.99}
10. [1] TestOrder{customer_id=P006, product_id=PROD001, region=east, amount=1, total=75000, date=2024-03-15, discount=0, status=refunded, priority=urgent}

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

1. [1] TestPerson{tags=junior, status=active, department=sales, salary=45000, active=true, name=Alice, age=25, score=8.5, level=2}
2. [1] TestPerson{department=engineering, age=35, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2, tags=senior}
3. [1] TestPerson{name=Charlie, score=6, age=16, salary=0, active=false, tags=intern, status=inactive, department=hr, level=1}
4. [1] TestPerson{level=7, age=45, name=Diana, status=active, salary=85000, active=true, score=7.8, tags=manager, department=marketing}
5. [1] TestPerson{tags=employee, salary=55000, department=sales, name=Eve, level=3, age=30, active=false, score=8, status=inactive}
6. [1] TestPerson{department=qa, level=1, active=true, age=0, salary=-5000, score=0, tags=test, name=Frank, status=active}
7. [1] TestPerson{name=Grace, age=65, score=10, level=9, salary=95000, tags=executive, status=active, department=management, active=true}
8. [1] TestPerson{age=18, score=5.5, tags=junior, department=support, name=Henry, active=false, status=inactive, salary=25000, level=1}
9. [1] TestPerson{age=40, status=active, department=engineering, active=true, score=8.7, tags=senior, salary=68000, level=6, name=Ivy}
10. [1] TestPerson{score=6.5, name=X, salary=28000, tags=temp, status=active, level=1, active=true, department=intern, age=22}

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

1. [1] TestOrder{customer_id=P001, priority=normal, discount=50, date=2024-01-15, region=north, product_id=PROD001, status=pending, amount=2, total=1999.98}
2. [1] TestOrder{total=25.5, status=confirmed, date=2024-01-20, discount=0, region=south, amount=1, priority=low, product_id=PROD002, customer_id=P002}
3. [1] TestOrder{region=north, amount=3, priority=high, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, status=shipped, total=225}
4. [1] TestOrder{product_id=PROD004, date=2024-02-05, discount=0, customer_id=P004, amount=1, status=delivered, total=299.99, priority=normal, region=east}
5. [1] TestOrder{status=confirmed, product_id=PROD001, region=south, customer_id=P002, amount=1, total=999.99, priority=high, discount=100, date=2024-02-10}
6. [1] TestOrder{status=cancelled, customer_id=P005, total=999.98, priority=low, region=west, discount=0, product_id=PROD005, amount=2, date=2024-02-15}
7. [1] TestOrder{amount=4, date=2024-03-01, customer_id=P007, product_id=PROD006, total=600, priority=urgent, discount=50, region=north, status=shipped}
8. [1] TestOrder{customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, priority=normal, region=south, status=pending}
9. [1] TestOrder{status=completed, priority=low, region=north, amount=1, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, total=89.99}
10. [1] TestOrder{product_id=PROD001, region=east, amount=1, total=75000, date=2024-03-15, discount=0, status=refunded, priority=urgent, customer_id=P006}

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

1. [1] TestPerson{name=Alice, age=25, score=8.5, level=2, tags=junior, status=active, department=sales, salary=45000, active=true}
2. [1] TestPerson{level=5, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, age=35, status=active}
3. [1] TestPerson{level=1, name=Charlie, score=6, age=16, salary=0, active=false, tags=intern, status=inactive, department=hr}
4. [1] TestPerson{department=marketing, level=7, age=45, name=Diana, status=active, salary=85000, active=true, score=7.8, tags=manager}
5. [1] TestPerson{score=8, status=inactive, tags=employee, salary=55000, department=sales, name=Eve, level=3, age=30, active=false}
6. [1] TestPerson{age=0, salary=-5000, score=0, tags=test, name=Frank, status=active, department=qa, level=1, active=true}
7. [1] TestPerson{active=true, name=Grace, age=65, score=10, level=9, salary=95000, tags=executive, status=active, department=management}
8. [1] TestPerson{tags=junior, department=support, name=Henry, active=false, status=inactive, salary=25000, level=1, age=18, score=5.5}
9. [1] TestPerson{salary=68000, level=6, name=Ivy, age=40, status=active, department=engineering, active=true, score=8.7, tags=senior}
10. [1] TestPerson{age=22, score=6.5, name=X, salary=28000, tags=temp, status=active, level=1, active=true, department=intern}

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

1. [1] TestProduct{supplier=TechSupply, stock=50, price=999.99, available=true, rating=4.5, keywords=computer, brand=TechCorp, name=Laptop, category=electronics}
2. [1] TestProduct{name=Mouse, brand=TechCorp, supplier=TechSupply, rating=4, keywords=peripheral, price=25.5, stock=200, category=accessories, available=true}
3. [1] TestProduct{name=Keyboard, available=false, keywords=typing, brand=KeyTech, rating=3.5, stock=0, category=accessories, price=75, supplier=KeySupply}
4. [1] TestProduct{price=299.99, available=true, brand=ScreenPro, supplier=ScreenSupply, rating=4.8, keywords=display, category=electronics, stock=30, name=Monitor}
5. [1] TestProduct{rating=2, keywords=obsolete, stock=0, supplier=OldSupply, category=accessories, price=8.5, brand=OldTech, name=OldKeyboard, available=false}
6. [1] TestProduct{name=Headphones, category=audio, price=150, available=true, rating=4.6, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply}
7. [1] TestProduct{price=89.99, keywords=video, stock=25, supplier=CamSupply, name=Webcam, category=electronics, rating=3.8, available=true, brand=CamTech}

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

1. [1] TestPerson{level=2, tags=junior, status=active, department=sales, salary=45000, active=true, name=Alice, age=25, score=8.5}
2. [1] TestPerson{salary=75000, active=true, score=9.2, tags=senior, department=engineering, age=35, status=active, level=5, name=Bob}
3. [1] TestPerson{name=Charlie, score=6, age=16, salary=0, active=false, tags=intern, status=inactive, department=hr, level=1}
4. [1] TestPerson{age=45, name=Diana, status=active, salary=85000, active=true, score=7.8, tags=manager, department=marketing, level=7}
5. [1] TestPerson{name=Eve, level=3, age=30, active=false, score=8, status=inactive, tags=employee, salary=55000, department=sales}
6. [1] TestPerson{name=Frank, status=active, department=qa, level=1, active=true, age=0, salary=-5000, score=0, tags=test}
7. [1] TestPerson{active=true, name=Grace, age=65, score=10, level=9, salary=95000, tags=executive, status=active, department=management}
8. [1] TestPerson{age=18, score=5.5, tags=junior, department=support, name=Henry, active=false, status=inactive, salary=25000, level=1}
9. [1] TestPerson{score=8.7, tags=senior, salary=68000, level=6, name=Ivy, age=40, status=active, department=engineering, active=true}
10. [1] TestPerson{age=22, score=6.5, name=X, salary=28000, tags=temp, status=active, level=1, active=true, department=intern}

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

1. [1] TestOrder{total=1999.98, customer_id=P001, priority=normal, discount=50, date=2024-01-15, region=north, product_id=PROD001, status=pending, amount=2}
2. [1] TestOrder{total=25.5, status=confirmed, date=2024-01-20, discount=0, region=south, amount=1, priority=low, product_id=PROD002, customer_id=P002}
3. [1] TestOrder{date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, status=shipped, total=225, region=north, amount=3, priority=high}
4. [1] TestOrder{date=2024-02-05, discount=0, customer_id=P004, amount=1, status=delivered, total=299.99, priority=normal, region=east, product_id=PROD004}
5. [1] TestOrder{status=confirmed, product_id=PROD001, region=south, customer_id=P002, amount=1, total=999.99, priority=high, discount=100, date=2024-02-10}
6. [1] TestOrder{priority=low, region=west, discount=0, product_id=PROD005, amount=2, date=2024-02-15, status=cancelled, customer_id=P005, total=999.98}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, total=600, priority=urgent, discount=50, region=north, status=shipped, amount=4, date=2024-03-01}
8. [1] TestOrder{date=2024-03-05, discount=0, priority=normal, region=south, status=pending, customer_id=P010, product_id=PROD002, amount=10, total=255}
9. [1] TestOrder{discount=10, customer_id=P001, product_id=PROD007, total=89.99, status=completed, priority=low, region=north, amount=1, date=2024-03-10}
10. [1] TestOrder{customer_id=P006, product_id=PROD001, region=east, amount=1, total=75000, date=2024-03-15, discount=0, status=refunded, priority=urgent}

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

1. [1] TestPerson{status=active, department=sales, salary=45000, active=true, name=Alice, age=25, score=8.5, level=2, tags=junior}
2. [1] TestPerson{name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, age=35, status=active, level=5}
3. [1] TestPerson{status=inactive, department=hr, level=1, name=Charlie, score=6, age=16, salary=0, active=false, tags=intern}
4. [1] TestPerson{salary=85000, active=true, score=7.8, tags=manager, department=marketing, level=7, age=45, name=Diana, status=active}
5. [1] TestPerson{salary=55000, department=sales, name=Eve, level=3, age=30, active=false, score=8, status=inactive, tags=employee}
6. [1] TestPerson{age=0, salary=-5000, score=0, tags=test, name=Frank, status=active, department=qa, level=1, active=true}
7. [1] TestPerson{name=Grace, age=65, score=10, level=9, salary=95000, tags=executive, status=active, department=management, active=true}
8. [1] TestPerson{level=1, age=18, score=5.5, tags=junior, department=support, name=Henry, active=false, status=inactive, salary=25000}
9. [1] TestPerson{status=active, department=engineering, active=true, score=8.7, tags=senior, salary=68000, level=6, name=Ivy, age=40}
10. [1] TestPerson{age=22, score=6.5, name=X, salary=28000, tags=temp, status=active, level=1, active=true, department=intern}
11. [1] TestOrder{customer_id=P001, priority=normal, discount=50, date=2024-01-15, region=north, product_id=PROD001, status=pending, amount=2, total=1999.98}
12. [1] TestOrder{total=25.5, status=confirmed, date=2024-01-20, discount=0, region=south, amount=1, priority=low, product_id=PROD002, customer_id=P002}
13. [1] TestOrder{product_id=PROD003, status=shipped, total=225, region=north, amount=3, priority=high, date=2024-02-01, discount=15, customer_id=P001}
14. [1] TestOrder{customer_id=P004, amount=1, status=delivered, total=299.99, priority=normal, region=east, product_id=PROD004, date=2024-02-05, discount=0}
15. [1] TestOrder{amount=1, total=999.99, priority=high, discount=100, date=2024-02-10, status=confirmed, product_id=PROD001, region=south, customer_id=P002}
16. [1] TestOrder{product_id=PROD005, amount=2, date=2024-02-15, status=cancelled, customer_id=P005, total=999.98, priority=low, region=west, discount=0}
17. [1] TestOrder{priority=urgent, discount=50, region=north, status=shipped, amount=4, date=2024-03-01, customer_id=P007, product_id=PROD006, total=600}
18. [1] TestOrder{status=pending, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, priority=normal, region=south}
19. [1] TestOrder{customer_id=P001, product_id=PROD007, total=89.99, status=completed, priority=low, region=north, amount=1, date=2024-03-10, discount=10}
20. [1] TestOrder{customer_id=P006, product_id=PROD001, region=east, amount=1, total=75000, date=2024-03-15, discount=0, status=refunded, priority=urgent}

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

1. [1] TestPerson{level=2, tags=junior, status=active, department=sales, salary=45000, active=true, name=Alice, age=25, score=8.5}
2. [1] TestPerson{level=5, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, age=35, status=active}
3. [1] TestPerson{age=16, salary=0, active=false, tags=intern, status=inactive, department=hr, level=1, name=Charlie, score=6}
4. [1] TestPerson{department=marketing, level=7, age=45, name=Diana, status=active, salary=85000, active=true, score=7.8, tags=manager}
5. [1] TestPerson{name=Eve, level=3, age=30, active=false, score=8, status=inactive, tags=employee, salary=55000, department=sales}
6. [1] TestPerson{age=0, salary=-5000, score=0, tags=test, name=Frank, status=active, department=qa, level=1, active=true}
7. [1] TestPerson{salary=95000, tags=executive, status=active, department=management, active=true, name=Grace, age=65, score=10, level=9}
8. [1] TestPerson{age=18, score=5.5, tags=junior, department=support, name=Henry, active=false, status=inactive, salary=25000, level=1}
9. [1] TestPerson{status=active, department=engineering, active=true, score=8.7, tags=senior, salary=68000, level=6, name=Ivy, age=40}
10. [1] TestPerson{age=22, score=6.5, name=X, salary=28000, tags=temp, status=active, level=1, active=true, department=intern}

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

1. [1] TestPerson{status=active, department=sales, salary=45000, active=true, name=Alice, age=25, score=8.5, level=2, tags=junior}
2. [1] TestPerson{name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, age=35, status=active, level=5}
3. [1] TestPerson{score=6, age=16, salary=0, active=false, tags=intern, status=inactive, department=hr, level=1, name=Charlie}
4. [1] TestPerson{salary=85000, active=true, score=7.8, tags=manager, department=marketing, level=7, age=45, name=Diana, status=active}
5. [1] TestPerson{tags=employee, salary=55000, department=sales, name=Eve, level=3, age=30, active=false, score=8, status=inactive}
6. [1] TestPerson{name=Frank, status=active, department=qa, level=1, active=true, age=0, salary=-5000, score=0, tags=test}
7. [1] TestPerson{department=management, active=true, name=Grace, age=65, score=10, level=9, salary=95000, tags=executive, status=active}
8. [1] TestPerson{age=18, score=5.5, tags=junior, department=support, name=Henry, active=false, status=inactive, salary=25000, level=1}
9. [1] TestPerson{department=engineering, active=true, score=8.7, tags=senior, salary=68000, level=6, name=Ivy, age=40, status=active}
10. [1] TestPerson{active=true, department=intern, age=22, score=6.5, name=X, salary=28000, tags=temp, status=active, level=1}
11. [1] TestOrder{amount=2, total=1999.98, customer_id=P001, priority=normal, discount=50, date=2024-01-15, region=north, product_id=PROD001, status=pending}
12. [1] TestOrder{priority=low, product_id=PROD002, customer_id=P002, total=25.5, status=confirmed, date=2024-01-20, discount=0, region=south, amount=1}
13. [1] TestOrder{date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, status=shipped, total=225, region=north, amount=3, priority=high}
14. [1] TestOrder{product_id=PROD004, date=2024-02-05, discount=0, customer_id=P004, amount=1, status=delivered, total=299.99, priority=normal, region=east}
15. [1] TestOrder{priority=high, discount=100, date=2024-02-10, status=confirmed, product_id=PROD001, region=south, customer_id=P002, amount=1, total=999.99}
16. [1] TestOrder{product_id=PROD005, amount=2, date=2024-02-15, status=cancelled, customer_id=P005, total=999.98, priority=low, region=west, discount=0}
17. [1] TestOrder{amount=4, date=2024-03-01, customer_id=P007, product_id=PROD006, total=600, priority=urgent, discount=50, region=north, status=shipped}
18. [1] TestOrder{total=255, date=2024-03-05, discount=0, priority=normal, region=south, status=pending, customer_id=P010, product_id=PROD002, amount=10}
19. [1] TestOrder{status=completed, priority=low, region=north, amount=1, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, total=89.99}
20. [1] TestOrder{product_id=PROD001, region=east, amount=1, total=75000, date=2024-03-15, discount=0, status=refunded, priority=urgent, customer_id=P006}

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

1. [1] TestPerson{status=active, department=sales, salary=45000, active=true, name=Alice, age=25, score=8.5, level=2, tags=junior}
2. [1] TestPerson{score=9.2, tags=senior, department=engineering, age=35, status=active, level=5, name=Bob, salary=75000, active=true}
3. [1] TestPerson{score=6, age=16, salary=0, active=false, tags=intern, status=inactive, department=hr, level=1, name=Charlie}
4. [1] TestPerson{name=Diana, status=active, salary=85000, active=true, score=7.8, tags=manager, department=marketing, level=7, age=45}
5. [1] TestPerson{name=Eve, level=3, age=30, active=false, score=8, status=inactive, tags=employee, salary=55000, department=sales}
6. [1] TestPerson{department=qa, level=1, active=true, age=0, salary=-5000, score=0, tags=test, name=Frank, status=active}
7. [1] TestPerson{age=65, score=10, level=9, salary=95000, tags=executive, status=active, department=management, active=true, name=Grace}
8. [1] TestPerson{status=inactive, salary=25000, level=1, age=18, score=5.5, tags=junior, department=support, name=Henry, active=false}
9. [1] TestPerson{department=engineering, active=true, score=8.7, tags=senior, salary=68000, level=6, name=Ivy, age=40, status=active}
10. [1] TestPerson{age=22, score=6.5, name=X, salary=28000, tags=temp, status=active, level=1, active=true, department=intern}

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

1. [1] TestOrder{discount=50, date=2024-01-15, region=north, product_id=PROD001, status=pending, amount=2, total=1999.98, customer_id=P001, priority=normal}
2. [1] TestOrder{product_id=PROD002, customer_id=P002, total=25.5, status=confirmed, date=2024-01-20, discount=0, region=south, amount=1, priority=low}
3. [1] TestOrder{date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, status=shipped, total=225, region=north, amount=3, priority=high}
4. [1] TestOrder{discount=0, customer_id=P004, amount=1, status=delivered, total=299.99, priority=normal, region=east, product_id=PROD004, date=2024-02-05}
5. [1] TestOrder{date=2024-02-10, status=confirmed, product_id=PROD001, region=south, customer_id=P002, amount=1, total=999.99, priority=high, discount=100}
6. [1] TestOrder{region=west, discount=0, product_id=PROD005, amount=2, date=2024-02-15, status=cancelled, customer_id=P005, total=999.98, priority=low}
7. [1] TestOrder{discount=50, region=north, status=shipped, amount=4, date=2024-03-01, customer_id=P007, product_id=PROD006, total=600, priority=urgent}
8. [1] TestOrder{amount=10, total=255, date=2024-03-05, discount=0, priority=normal, region=south, status=pending, customer_id=P010, product_id=PROD002}
9. [1] TestOrder{status=completed, priority=low, region=north, amount=1, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, total=89.99}
10. [1] TestOrder{amount=1, total=75000, date=2024-03-15, discount=0, status=refunded, priority=urgent, customer_id=P006, product_id=PROD001, region=east}

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

1. [1] TestPerson{name=Alice, age=25, score=8.5, level=2, tags=junior, status=active, department=sales, salary=45000, active=true}
2. [1] TestPerson{name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, age=35, status=active, level=5}
3. [1] TestPerson{tags=intern, status=inactive, department=hr, level=1, name=Charlie, score=6, age=16, salary=0, active=false}
4. [1] TestPerson{status=active, salary=85000, active=true, score=7.8, tags=manager, department=marketing, level=7, age=45, name=Diana}
5. [1] TestPerson{name=Eve, level=3, age=30, active=false, score=8, status=inactive, tags=employee, salary=55000, department=sales}
6. [1] TestPerson{department=qa, level=1, active=true, age=0, salary=-5000, score=0, tags=test, name=Frank, status=active}
7. [1] TestPerson{active=true, name=Grace, age=65, score=10, level=9, salary=95000, tags=executive, status=active, department=management}
8. [1] TestPerson{age=18, score=5.5, tags=junior, department=support, name=Henry, active=false, status=inactive, salary=25000, level=1}
9. [1] TestPerson{department=engineering, active=true, score=8.7, tags=senior, salary=68000, level=6, name=Ivy, age=40, status=active}
10. [1] TestPerson{age=22, score=6.5, name=X, salary=28000, tags=temp, status=active, level=1, active=true, department=intern}

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
