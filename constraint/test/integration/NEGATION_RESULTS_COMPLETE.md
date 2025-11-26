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

1. [1] TestPerson{status=active, level=2, age=25, tags=junior, department=sales, active=true, score=8.5, name=Alice, salary=45000}
2. [1] TestPerson{status=active, level=5, score=9.2, age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true}
3. [1] TestPerson{active=false, department=hr, level=1, name=Charlie, age=16, score=6, tags=intern, status=inactive, salary=0}
4. [1] TestPerson{score=7.8, level=7, status=active, name=Diana, salary=85000, age=45, active=true, tags=manager, department=marketing}
5. [1] TestPerson{department=sales, age=30, salary=55000, status=inactive, active=false, score=8, level=3, tags=employee, name=Eve}
6. [1] TestPerson{department=qa, age=0, salary=-5000, name=Frank, score=0, status=active, level=1, active=true, tags=test}
7. [1] TestPerson{name=Grace, level=9, salary=95000, score=10, department=management, age=65, tags=executive, status=active, active=true}
8. [1] TestPerson{score=5.5, tags=junior, department=support, level=1, status=inactive, age=18, salary=25000, name=Henry, active=false}
9. [1] TestPerson{score=8.7, status=active, level=6, name=Ivy, age=40, department=engineering, salary=68000, active=true, tags=senior}
10. [1] TestPerson{score=6.5, status=active, department=intern, level=1, age=22, active=true, name=X, tags=temp, salary=28000}

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

1. [1] TestOrder{region=north, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, customer_id=P001}
2. [1] TestOrder{priority=low, region=south, discount=0, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5}
3. [1] TestOrder{discount=15, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, region=north, customer_id=P001, total=225, priority=high}
4. [1] TestOrder{date=2024-02-05, priority=normal, discount=0, region=east, amount=1, status=delivered, customer_id=P004, total=299.99, product_id=PROD004}
5. [1] TestOrder{customer_id=P002, date=2024-02-10, discount=100, region=south, product_id=PROD001, status=confirmed, priority=high, amount=1, total=999.99}
6. [1] TestOrder{customer_id=P005, date=2024-02-15, priority=low, region=west, amount=2, total=999.98, status=cancelled, product_id=PROD005, discount=0}
7. [1] TestOrder{total=600, date=2024-03-01, priority=urgent, amount=4, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006}
8. [1] TestOrder{customer_id=P010, discount=0, date=2024-03-05, status=pending, priority=normal, product_id=PROD002, total=255, amount=10, region=south}
9. [1] TestOrder{region=north, customer_id=P001, product_id=PROD007, discount=10, amount=1, total=89.99, status=completed, priority=low, date=2024-03-10}
10. [1] TestOrder{priority=urgent, product_id=PROD001, amount=1, customer_id=P006, date=2024-03-15, discount=0, total=75000, status=refunded, region=east}

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

1. [1] TestPerson{level=2, age=25, tags=junior, department=sales, active=true, score=8.5, name=Alice, salary=45000, status=active}
2. [1] TestPerson{level=5, score=9.2, age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true, status=active}
3. [1] TestPerson{salary=0, active=false, department=hr, level=1, name=Charlie, age=16, score=6, tags=intern, status=inactive}
4. [1] TestPerson{score=7.8, level=7, status=active, name=Diana, salary=85000, age=45, active=true, tags=manager, department=marketing}
5. [1] TestPerson{tags=employee, name=Eve, department=sales, age=30, salary=55000, status=inactive, active=false, score=8, level=3}
6. [1] TestPerson{department=qa, age=0, salary=-5000, name=Frank, score=0, status=active, level=1, active=true, tags=test}
7. [1] TestPerson{age=65, tags=executive, status=active, active=true, name=Grace, level=9, salary=95000, score=10, department=management}
8. [1] TestPerson{status=inactive, age=18, salary=25000, name=Henry, active=false, score=5.5, tags=junior, department=support, level=1}
9. [1] TestPerson{salary=68000, active=true, tags=senior, score=8.7, status=active, level=6, name=Ivy, age=40, department=engineering}
10. [1] TestPerson{name=X, tags=temp, salary=28000, score=6.5, status=active, department=intern, level=1, age=22, active=true}

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

1. [1] TestOrder{date=2024-01-15, status=pending, priority=normal, discount=50, customer_id=P001, region=north, product_id=PROD001, total=1999.98, amount=2}
2. [1] TestOrder{amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, region=south, discount=0}
3. [1] TestOrder{status=shipped, region=north, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01}
4. [1] TestOrder{status=delivered, customer_id=P004, total=299.99, product_id=PROD004, date=2024-02-05, priority=normal, discount=0, region=east, amount=1}
5. [1] TestOrder{amount=1, total=999.99, customer_id=P002, date=2024-02-10, discount=100, region=south, product_id=PROD001, status=confirmed, priority=high}
6. [1] TestOrder{discount=0, customer_id=P005, date=2024-02-15, priority=low, region=west, amount=2, total=999.98, status=cancelled, product_id=PROD005}
7. [1] TestOrder{amount=4, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, priority=urgent}
8. [1] TestOrder{discount=0, date=2024-03-05, status=pending, priority=normal, product_id=PROD002, total=255, amount=10, region=south, customer_id=P010}
9. [1] TestOrder{product_id=PROD007, discount=10, amount=1, total=89.99, status=completed, priority=low, date=2024-03-10, region=north, customer_id=P001}
10. [1] TestOrder{amount=1, customer_id=P006, date=2024-03-15, discount=0, total=75000, status=refunded, region=east, priority=urgent, product_id=PROD001}

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

1. [1] TestProduct{supplier=TechSupply, name=Laptop, available=true, keywords=computer, category=electronics, rating=4.5, brand=TechCorp, price=999.99, stock=50}
2. [1] TestProduct{rating=4, stock=200, keywords=peripheral, supplier=TechSupply, category=accessories, price=25.5, available=true, brand=TechCorp, name=Mouse}
3. [1] TestProduct{category=accessories, price=75, rating=3.5, brand=KeyTech, name=Keyboard, available=false, keywords=typing, stock=0, supplier=KeySupply}
4. [1] TestProduct{available=true, stock=30, category=electronics, brand=ScreenPro, keywords=display, supplier=ScreenSupply, rating=4.8, name=Monitor, price=299.99}
5. [1] TestProduct{supplier=OldSupply, name=OldKeyboard, price=8.5, stock=0, category=accessories, brand=OldTech, keywords=obsolete, available=false, rating=2}
6. [1] TestProduct{name=Headphones, price=150, brand=AudioMax, category=audio, keywords=sound, stock=75, supplier=AudioSupply, rating=4.6, available=true}
7. [1] TestProduct{supplier=CamSupply, name=Webcam, available=true, rating=3.8, stock=25, category=electronics, keywords=video, price=89.99, brand=CamTech}

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

1. [1] TestPerson{level=2, age=25, tags=junior, department=sales, active=true, score=8.5, name=Alice, salary=45000, status=active}
2. [1] TestPerson{age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true, status=active, level=5, score=9.2}
3. [1] TestPerson{active=false, department=hr, level=1, name=Charlie, age=16, score=6, tags=intern, status=inactive, salary=0}
4. [1] TestPerson{score=7.8, level=7, status=active, name=Diana, salary=85000, age=45, active=true, tags=manager, department=marketing}
5. [1] TestPerson{department=sales, age=30, salary=55000, status=inactive, active=false, score=8, level=3, tags=employee, name=Eve}
6. [1] TestPerson{name=Frank, score=0, status=active, level=1, active=true, tags=test, department=qa, age=0, salary=-5000}
7. [1] TestPerson{status=active, active=true, name=Grace, level=9, salary=95000, score=10, department=management, age=65, tags=executive}
8. [1] TestPerson{active=false, score=5.5, tags=junior, department=support, level=1, status=inactive, age=18, salary=25000, name=Henry}
9. [1] TestPerson{status=active, level=6, name=Ivy, age=40, department=engineering, salary=68000, active=true, tags=senior, score=8.7}
10. [1] TestPerson{score=6.5, status=active, department=intern, level=1, age=22, active=true, name=X, tags=temp, salary=28000}

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

1. [1] TestOrder{product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, customer_id=P001, region=north}
2. [1] TestOrder{product_id=PROD002, total=25.5, priority=low, region=south, discount=0, amount=1, date=2024-01-20, status=confirmed, customer_id=P002}
3. [1] TestOrder{status=shipped, region=north, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01}
4. [1] TestOrder{product_id=PROD004, date=2024-02-05, priority=normal, discount=0, region=east, amount=1, status=delivered, customer_id=P004, total=299.99}
5. [1] TestOrder{customer_id=P002, date=2024-02-10, discount=100, region=south, product_id=PROD001, status=confirmed, priority=high, amount=1, total=999.99}
6. [1] TestOrder{discount=0, customer_id=P005, date=2024-02-15, priority=low, region=west, amount=2, total=999.98, status=cancelled, product_id=PROD005}
7. [1] TestOrder{discount=50, region=north, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, priority=urgent, amount=4, status=shipped}
8. [1] TestOrder{region=south, customer_id=P010, discount=0, date=2024-03-05, status=pending, priority=normal, product_id=PROD002, total=255, amount=10}
9. [1] TestOrder{priority=low, date=2024-03-10, region=north, customer_id=P001, product_id=PROD007, discount=10, amount=1, total=89.99, status=completed}
10. [1] TestOrder{discount=0, total=75000, status=refunded, region=east, priority=urgent, product_id=PROD001, amount=1, customer_id=P006, date=2024-03-15}

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

1. [1] TestPerson{active=true, score=8.5, name=Alice, salary=45000, status=active, level=2, age=25, tags=junior, department=sales}
2. [1] TestPerson{tags=senior, department=engineering, name=Bob, salary=75000, active=true, status=active, level=5, score=9.2, age=35}
3. [1] TestPerson{name=Charlie, age=16, score=6, tags=intern, status=inactive, salary=0, active=false, department=hr, level=1}
4. [1] TestPerson{name=Diana, salary=85000, age=45, active=true, tags=manager, department=marketing, score=7.8, level=7, status=active}
5. [1] TestPerson{level=3, tags=employee, name=Eve, department=sales, age=30, salary=55000, status=inactive, active=false, score=8}
6. [1] TestPerson{age=0, salary=-5000, name=Frank, score=0, status=active, level=1, active=true, tags=test, department=qa}
7. [1] TestPerson{department=management, age=65, tags=executive, status=active, active=true, name=Grace, level=9, salary=95000, score=10}
8. [1] TestPerson{age=18, salary=25000, name=Henry, active=false, score=5.5, tags=junior, department=support, level=1, status=inactive}
9. [1] TestPerson{active=true, tags=senior, score=8.7, status=active, level=6, name=Ivy, age=40, department=engineering, salary=68000}
10. [1] TestPerson{score=6.5, status=active, department=intern, level=1, age=22, active=true, name=X, tags=temp, salary=28000}

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

1. [1] TestOrder{date=2024-01-15, status=pending, priority=normal, discount=50, customer_id=P001, region=north, product_id=PROD001, total=1999.98, amount=2}
2. [1] TestOrder{region=south, discount=0, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low}
3. [1] TestOrder{status=shipped, region=north, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01}
4. [1] TestOrder{product_id=PROD004, date=2024-02-05, priority=normal, discount=0, region=east, amount=1, status=delivered, customer_id=P004, total=299.99}
5. [1] TestOrder{priority=high, amount=1, total=999.99, customer_id=P002, date=2024-02-10, discount=100, region=south, product_id=PROD001, status=confirmed}
6. [1] TestOrder{product_id=PROD005, discount=0, customer_id=P005, date=2024-02-15, priority=low, region=west, amount=2, total=999.98, status=cancelled}
7. [1] TestOrder{amount=4, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, priority=urgent}
8. [1] TestOrder{discount=0, date=2024-03-05, status=pending, priority=normal, product_id=PROD002, total=255, amount=10, region=south, customer_id=P010}
9. [1] TestOrder{status=completed, priority=low, date=2024-03-10, region=north, customer_id=P001, product_id=PROD007, discount=10, amount=1, total=89.99}
10. [1] TestOrder{customer_id=P006, date=2024-03-15, discount=0, total=75000, status=refunded, region=east, priority=urgent, product_id=PROD001, amount=1}

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

1. [1] TestPerson{score=8.5, name=Alice, salary=45000, status=active, level=2, age=25, tags=junior, department=sales, active=true}
2. [1] TestPerson{age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true, status=active, level=5, score=9.2}
3. [1] TestPerson{score=6, tags=intern, status=inactive, salary=0, active=false, department=hr, level=1, name=Charlie, age=16}
4. [1] TestPerson{salary=85000, age=45, active=true, tags=manager, department=marketing, score=7.8, level=7, status=active, name=Diana}
5. [1] TestPerson{name=Eve, department=sales, age=30, salary=55000, status=inactive, active=false, score=8, level=3, tags=employee}
6. [1] TestPerson{name=Frank, score=0, status=active, level=1, active=true, tags=test, department=qa, age=0, salary=-5000}
7. [1] TestPerson{age=65, tags=executive, status=active, active=true, name=Grace, level=9, salary=95000, score=10, department=management}
8. [1] TestPerson{score=5.5, tags=junior, department=support, level=1, status=inactive, age=18, salary=25000, name=Henry, active=false}
9. [1] TestPerson{score=8.7, status=active, level=6, name=Ivy, age=40, department=engineering, salary=68000, active=true, tags=senior}
10. [1] TestPerson{score=6.5, status=active, department=intern, level=1, age=22, active=true, name=X, tags=temp, salary=28000}

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

1. [1] TestProduct{price=999.99, stock=50, supplier=TechSupply, name=Laptop, available=true, keywords=computer, category=electronics, rating=4.5, brand=TechCorp}
2. [1] TestProduct{stock=200, keywords=peripheral, supplier=TechSupply, category=accessories, price=25.5, available=true, brand=TechCorp, name=Mouse, rating=4}
3. [1] TestProduct{category=accessories, price=75, rating=3.5, brand=KeyTech, name=Keyboard, available=false, keywords=typing, stock=0, supplier=KeySupply}
4. [1] TestProduct{supplier=ScreenSupply, rating=4.8, name=Monitor, price=299.99, available=true, stock=30, category=electronics, brand=ScreenPro, keywords=display}
5. [1] TestProduct{brand=OldTech, keywords=obsolete, available=false, rating=2, supplier=OldSupply, name=OldKeyboard, price=8.5, stock=0, category=accessories}
6. [1] TestProduct{name=Headphones, price=150, brand=AudioMax, category=audio, keywords=sound, stock=75, supplier=AudioSupply, rating=4.6, available=true}
7. [1] TestProduct{brand=CamTech, supplier=CamSupply, name=Webcam, available=true, rating=3.8, stock=25, category=electronics, keywords=video, price=89.99}

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

1. [1] TestPerson{salary=45000, status=active, level=2, age=25, tags=junior, department=sales, active=true, score=8.5, name=Alice}
2. [1] TestPerson{department=engineering, name=Bob, salary=75000, active=true, status=active, level=5, score=9.2, age=35, tags=senior}
3. [1] TestPerson{salary=0, active=false, department=hr, level=1, name=Charlie, age=16, score=6, tags=intern, status=inactive}
4. [1] TestPerson{tags=manager, department=marketing, score=7.8, level=7, status=active, name=Diana, salary=85000, age=45, active=true}
5. [1] TestPerson{tags=employee, name=Eve, department=sales, age=30, salary=55000, status=inactive, active=false, score=8, level=3}
6. [1] TestPerson{name=Frank, score=0, status=active, level=1, active=true, tags=test, department=qa, age=0, salary=-5000}
7. [1] TestPerson{status=active, active=true, name=Grace, level=9, salary=95000, score=10, department=management, age=65, tags=executive}
8. [1] TestPerson{active=false, score=5.5, tags=junior, department=support, level=1, status=inactive, age=18, salary=25000, name=Henry}
9. [1] TestPerson{department=engineering, salary=68000, active=true, tags=senior, score=8.7, status=active, level=6, name=Ivy, age=40}
10. [1] TestPerson{score=6.5, status=active, department=intern, level=1, age=22, active=true, name=X, tags=temp, salary=28000}

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

1. [1] TestOrder{status=pending, priority=normal, discount=50, customer_id=P001, region=north, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15}
2. [1] TestOrder{product_id=PROD002, total=25.5, priority=low, region=south, discount=0, amount=1, date=2024-01-20, status=confirmed, customer_id=P002}
3. [1] TestOrder{status=shipped, region=north, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01}
4. [1] TestOrder{date=2024-02-05, priority=normal, discount=0, region=east, amount=1, status=delivered, customer_id=P004, total=299.99, product_id=PROD004}
5. [1] TestOrder{amount=1, total=999.99, customer_id=P002, date=2024-02-10, discount=100, region=south, product_id=PROD001, status=confirmed, priority=high}
6. [1] TestOrder{product_id=PROD005, discount=0, customer_id=P005, date=2024-02-15, priority=low, region=west, amount=2, total=999.98, status=cancelled}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, priority=urgent, amount=4, status=shipped, discount=50, region=north}
8. [1] TestOrder{region=south, customer_id=P010, discount=0, date=2024-03-05, status=pending, priority=normal, product_id=PROD002, total=255, amount=10}
9. [1] TestOrder{product_id=PROD007, discount=10, amount=1, total=89.99, status=completed, priority=low, date=2024-03-10, region=north, customer_id=P001}
10. [1] TestOrder{total=75000, status=refunded, region=east, priority=urgent, product_id=PROD001, amount=1, customer_id=P006, date=2024-03-15, discount=0}

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

1. [1] TestPerson{name=Alice, salary=45000, status=active, level=2, age=25, tags=junior, department=sales, active=true, score=8.5}
2. [1] TestPerson{score=9.2, age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true, status=active, level=5}
3. [1] TestPerson{name=Charlie, age=16, score=6, tags=intern, status=inactive, salary=0, active=false, department=hr, level=1}
4. [1] TestPerson{level=7, status=active, name=Diana, salary=85000, age=45, active=true, tags=manager, department=marketing, score=7.8}
5. [1] TestPerson{score=8, level=3, tags=employee, name=Eve, department=sales, age=30, salary=55000, status=inactive, active=false}
6. [1] TestPerson{name=Frank, score=0, status=active, level=1, active=true, tags=test, department=qa, age=0, salary=-5000}
7. [1] TestPerson{age=65, tags=executive, status=active, active=true, name=Grace, level=9, salary=95000, score=10, department=management}
8. [1] TestPerson{score=5.5, tags=junior, department=support, level=1, status=inactive, age=18, salary=25000, name=Henry, active=false}
9. [1] TestPerson{level=6, name=Ivy, age=40, department=engineering, salary=68000, active=true, tags=senior, score=8.7, status=active}
10. [1] TestPerson{tags=temp, salary=28000, score=6.5, status=active, department=intern, level=1, age=22, active=true, name=X}
11. [1] TestOrder{priority=normal, discount=50, customer_id=P001, region=north, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending}
12. [1] TestOrder{product_id=PROD002, total=25.5, priority=low, region=south, discount=0, amount=1, date=2024-01-20, status=confirmed, customer_id=P002}
13. [1] TestOrder{total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, region=north, customer_id=P001}
14. [1] TestOrder{customer_id=P004, total=299.99, product_id=PROD004, date=2024-02-05, priority=normal, discount=0, region=east, amount=1, status=delivered}
15. [1] TestOrder{amount=1, total=999.99, customer_id=P002, date=2024-02-10, discount=100, region=south, product_id=PROD001, status=confirmed, priority=high}
16. [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, discount=0, customer_id=P005, date=2024-02-15, priority=low, region=west, amount=2}
17. [1] TestOrder{amount=4, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, priority=urgent}
18. [1] TestOrder{status=pending, priority=normal, product_id=PROD002, total=255, amount=10, region=south, customer_id=P010, discount=0, date=2024-03-05}
19. [1] TestOrder{product_id=PROD007, discount=10, amount=1, total=89.99, status=completed, priority=low, date=2024-03-10, region=north, customer_id=P001}
20. [1] TestOrder{total=75000, status=refunded, region=east, priority=urgent, product_id=PROD001, amount=1, customer_id=P006, date=2024-03-15, discount=0}

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

1. [1] TestPerson{salary=45000, status=active, level=2, age=25, tags=junior, department=sales, active=true, score=8.5, name=Alice}
2. [1] TestPerson{tags=senior, department=engineering, name=Bob, salary=75000, active=true, status=active, level=5, score=9.2, age=35}
3. [1] TestPerson{salary=0, active=false, department=hr, level=1, name=Charlie, age=16, score=6, tags=intern, status=inactive}
4. [1] TestPerson{name=Diana, salary=85000, age=45, active=true, tags=manager, department=marketing, score=7.8, level=7, status=active}
5. [1] TestPerson{status=inactive, active=false, score=8, level=3, tags=employee, name=Eve, department=sales, age=30, salary=55000}
6. [1] TestPerson{active=true, tags=test, department=qa, age=0, salary=-5000, name=Frank, score=0, status=active, level=1}
7. [1] TestPerson{name=Grace, level=9, salary=95000, score=10, department=management, age=65, tags=executive, status=active, active=true}
8. [1] TestPerson{score=5.5, tags=junior, department=support, level=1, status=inactive, age=18, salary=25000, name=Henry, active=false}
9. [1] TestPerson{level=6, name=Ivy, age=40, department=engineering, salary=68000, active=true, tags=senior, score=8.7, status=active}
10. [1] TestPerson{score=6.5, status=active, department=intern, level=1, age=22, active=true, name=X, tags=temp, salary=28000}

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

1. [1] TestPerson{name=Alice, salary=45000, status=active, level=2, age=25, tags=junior, department=sales, active=true, score=8.5}
2. [1] TestPerson{score=9.2, age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true, status=active, level=5}
3. [1] TestPerson{score=6, tags=intern, status=inactive, salary=0, active=false, department=hr, level=1, name=Charlie, age=16}
4. [1] TestPerson{name=Diana, salary=85000, age=45, active=true, tags=manager, department=marketing, score=7.8, level=7, status=active}
5. [1] TestPerson{level=3, tags=employee, name=Eve, department=sales, age=30, salary=55000, status=inactive, active=false, score=8}
6. [1] TestPerson{name=Frank, score=0, status=active, level=1, active=true, tags=test, department=qa, age=0, salary=-5000}
7. [1] TestPerson{age=65, tags=executive, status=active, active=true, name=Grace, level=9, salary=95000, score=10, department=management}
8. [1] TestPerson{department=support, level=1, status=inactive, age=18, salary=25000, name=Henry, active=false, score=5.5, tags=junior}
9. [1] TestPerson{active=true, tags=senior, score=8.7, status=active, level=6, name=Ivy, age=40, department=engineering, salary=68000}
10. [1] TestPerson{score=6.5, status=active, department=intern, level=1, age=22, active=true, name=X, tags=temp, salary=28000}
11. [1] TestOrder{date=2024-01-15, status=pending, priority=normal, discount=50, customer_id=P001, region=north, product_id=PROD001, total=1999.98, amount=2}
12. [1] TestOrder{amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, region=south, discount=0}
13. [1] TestOrder{status=shipped, region=north, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01}
14. [1] TestOrder{total=299.99, product_id=PROD004, date=2024-02-05, priority=normal, discount=0, region=east, amount=1, status=delivered, customer_id=P004}
15. [1] TestOrder{amount=1, total=999.99, customer_id=P002, date=2024-02-10, discount=100, region=south, product_id=PROD001, status=confirmed, priority=high}
16. [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, discount=0, customer_id=P005, date=2024-02-15, priority=low, region=west, amount=2}
17. [1] TestOrder{customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, priority=urgent, amount=4, status=shipped, discount=50, region=north}
18. [1] TestOrder{date=2024-03-05, status=pending, priority=normal, product_id=PROD002, total=255, amount=10, region=south, customer_id=P010, discount=0}
19. [1] TestOrder{product_id=PROD007, discount=10, amount=1, total=89.99, status=completed, priority=low, date=2024-03-10, region=north, customer_id=P001}
20. [1] TestOrder{amount=1, customer_id=P006, date=2024-03-15, discount=0, total=75000, status=refunded, region=east, priority=urgent, product_id=PROD001}

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

1. [1] TestPerson{name=Alice, salary=45000, status=active, level=2, age=25, tags=junior, department=sales, active=true, score=8.5}
2. [1] TestPerson{status=active, level=5, score=9.2, age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true}
3. [1] TestPerson{name=Charlie, age=16, score=6, tags=intern, status=inactive, salary=0, active=false, department=hr, level=1}
4. [1] TestPerson{name=Diana, salary=85000, age=45, active=true, tags=manager, department=marketing, score=7.8, level=7, status=active}
5. [1] TestPerson{age=30, salary=55000, status=inactive, active=false, score=8, level=3, tags=employee, name=Eve, department=sales}
6. [1] TestPerson{name=Frank, score=0, status=active, level=1, active=true, tags=test, department=qa, age=0, salary=-5000}
7. [1] TestPerson{name=Grace, level=9, salary=95000, score=10, department=management, age=65, tags=executive, status=active, active=true}
8. [1] TestPerson{salary=25000, name=Henry, active=false, score=5.5, tags=junior, department=support, level=1, status=inactive, age=18}
9. [1] TestPerson{score=8.7, status=active, level=6, name=Ivy, age=40, department=engineering, salary=68000, active=true, tags=senior}
10. [1] TestPerson{tags=temp, salary=28000, score=6.5, status=active, department=intern, level=1, age=22, active=true, name=X}

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

1. [1] TestOrder{date=2024-01-15, status=pending, priority=normal, discount=50, customer_id=P001, region=north, product_id=PROD001, total=1999.98, amount=2}
2. [1] TestOrder{amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, region=south, discount=0}
3. [1] TestOrder{status=shipped, region=north, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01}
4. [1] TestOrder{product_id=PROD004, date=2024-02-05, priority=normal, discount=0, region=east, amount=1, status=delivered, customer_id=P004, total=299.99}
5. [1] TestOrder{customer_id=P002, date=2024-02-10, discount=100, region=south, product_id=PROD001, status=confirmed, priority=high, amount=1, total=999.99}
6. [1] TestOrder{region=west, amount=2, total=999.98, status=cancelled, product_id=PROD005, discount=0, customer_id=P005, date=2024-02-15, priority=low}
7. [1] TestOrder{amount=4, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, priority=urgent}
8. [1] TestOrder{region=south, customer_id=P010, discount=0, date=2024-03-05, status=pending, priority=normal, product_id=PROD002, total=255, amount=10}
9. [1] TestOrder{product_id=PROD007, discount=10, amount=1, total=89.99, status=completed, priority=low, date=2024-03-10, region=north, customer_id=P001}
10. [1] TestOrder{total=75000, status=refunded, region=east, priority=urgent, product_id=PROD001, amount=1, customer_id=P006, date=2024-03-15, discount=0}

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

1. [1] TestPerson{age=25, tags=junior, department=sales, active=true, score=8.5, name=Alice, salary=45000, status=active, level=2}
2. [1] TestPerson{age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true, status=active, level=5, score=9.2}
3. [1] TestPerson{active=false, department=hr, level=1, name=Charlie, age=16, score=6, tags=intern, status=inactive, salary=0}
4. [1] TestPerson{score=7.8, level=7, status=active, name=Diana, salary=85000, age=45, active=true, tags=manager, department=marketing}
5. [1] TestPerson{name=Eve, department=sales, age=30, salary=55000, status=inactive, active=false, score=8, level=3, tags=employee}
6. [1] TestPerson{age=0, salary=-5000, name=Frank, score=0, status=active, level=1, active=true, tags=test, department=qa}
7. [1] TestPerson{name=Grace, level=9, salary=95000, score=10, department=management, age=65, tags=executive, status=active, active=true}
8. [1] TestPerson{department=support, level=1, status=inactive, age=18, salary=25000, name=Henry, active=false, score=5.5, tags=junior}
9. [1] TestPerson{score=8.7, status=active, level=6, name=Ivy, age=40, department=engineering, salary=68000, active=true, tags=senior}
10. [1] TestPerson{score=6.5, status=active, department=intern, level=1, age=22, active=true, name=X, tags=temp, salary=28000}

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
