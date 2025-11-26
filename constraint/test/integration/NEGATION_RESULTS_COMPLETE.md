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

1. [1] TestPerson{salary=45000, active=true, department=sales, level=2, status=active, name=Alice, age=25, score=8.5, tags=junior}
2. [1] TestPerson{tags=senior, name=Bob, status=active, level=5, age=35, salary=75000, active=true, department=engineering, score=9.2}
3. [1] TestPerson{status=inactive, department=hr, level=1, age=16, score=6, name=Charlie, active=false, tags=intern, salary=0}
4. [1] TestPerson{score=7.8, level=7, name=Diana, department=marketing, salary=85000, tags=manager, status=active, age=45, active=true}
5. [1] TestPerson{age=30, active=false, name=Eve, tags=employee, level=3, salary=55000, department=sales, score=8, status=inactive}
6. [1] TestPerson{status=active, department=qa, age=0, salary=-5000, active=true, name=Frank, score=0, tags=test, level=1}
7. [1] TestPerson{salary=95000, tags=executive, age=65, active=true, status=active, department=management, level=9, name=Grace, score=10}
8. [1] TestPerson{name=Henry, status=inactive, level=1, age=18, active=false, tags=junior, salary=25000, score=5.5, department=support}
9. [1] TestPerson{active=true, score=8.7, tags=senior, name=Ivy, salary=68000, status=active, department=engineering, level=6, age=40}
10. [1] TestPerson{salary=28000, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, status=active, level=1}

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

1. [1] TestOrder{discount=50, region=north, product_id=PROD001, amount=2, total=1999.98, status=pending, priority=normal, customer_id=P001, date=2024-01-15}
2. [1] TestOrder{date=2024-01-20, region=south, amount=1, total=25.5, status=confirmed, discount=0, product_id=PROD002, priority=low, customer_id=P002}
3. [1] TestOrder{date=2024-02-01, priority=high, product_id=PROD003, total=225, status=shipped, customer_id=P001, discount=15, region=north, amount=3}
4. [1] TestOrder{customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, status=delivered, region=east, amount=1, total=299.99, discount=0}
5. [1] TestOrder{customer_id=P002, total=999.99, date=2024-02-10, discount=100, status=confirmed, region=south, product_id=PROD001, amount=1, priority=high}
6. [1] TestOrder{amount=2, total=999.98, status=cancelled, date=2024-02-15, discount=0, priority=low, region=west, product_id=PROD005, customer_id=P005}
7. [1] TestOrder{customer_id=P007, amount=4, discount=50, region=north, product_id=PROD006, date=2024-03-01, status=shipped, priority=urgent, total=600}
8. [1] TestOrder{customer_id=P010, total=255, amount=10, date=2024-03-05, region=south, priority=normal, discount=0, product_id=PROD002, status=pending}
9. [1] TestOrder{product_id=PROD007, status=completed, priority=low, amount=1, discount=10, region=north, total=89.99, customer_id=P001, date=2024-03-10}
10. [1] TestOrder{region=east, total=75000, status=refunded, customer_id=P006, product_id=PROD001, priority=urgent, amount=1, date=2024-03-15, discount=0}

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

1. [1] TestPerson{salary=45000, active=true, department=sales, level=2, status=active, name=Alice, age=25, score=8.5, tags=junior}
2. [1] TestPerson{status=active, level=5, age=35, salary=75000, active=true, department=engineering, score=9.2, tags=senior, name=Bob}
3. [1] TestPerson{salary=0, status=inactive, department=hr, level=1, age=16, score=6, name=Charlie, active=false, tags=intern}
4. [1] TestPerson{age=45, active=true, score=7.8, level=7, name=Diana, department=marketing, salary=85000, tags=manager, status=active}
5. [1] TestPerson{score=8, status=inactive, age=30, active=false, name=Eve, tags=employee, level=3, salary=55000, department=sales}
6. [1] TestPerson{salary=-5000, active=true, name=Frank, score=0, tags=test, level=1, status=active, department=qa, age=0}
7. [1] TestPerson{name=Grace, score=10, salary=95000, tags=executive, age=65, active=true, status=active, department=management, level=9}
8. [1] TestPerson{tags=junior, salary=25000, score=5.5, department=support, name=Henry, status=inactive, level=1, age=18, active=false}
9. [1] TestPerson{active=true, score=8.7, tags=senior, name=Ivy, salary=68000, status=active, department=engineering, level=6, age=40}
10. [1] TestPerson{score=6.5, tags=temp, department=intern, name=X, age=22, status=active, level=1, salary=28000, active=true}

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

1. [1] TestOrder{discount=50, region=north, product_id=PROD001, amount=2, total=1999.98, status=pending, priority=normal, customer_id=P001, date=2024-01-15}
2. [1] TestOrder{amount=1, total=25.5, status=confirmed, discount=0, product_id=PROD002, priority=low, customer_id=P002, date=2024-01-20, region=south}
3. [1] TestOrder{date=2024-02-01, priority=high, product_id=PROD003, total=225, status=shipped, customer_id=P001, discount=15, region=north, amount=3}
4. [1] TestOrder{priority=normal, product_id=PROD004, status=delivered, region=east, amount=1, total=299.99, discount=0, customer_id=P004, date=2024-02-05}
5. [1] TestOrder{status=confirmed, region=south, product_id=PROD001, amount=1, priority=high, customer_id=P002, total=999.99, date=2024-02-10, discount=100}
6. [1] TestOrder{region=west, product_id=PROD005, customer_id=P005, amount=2, total=999.98, status=cancelled, date=2024-02-15, discount=0, priority=low}
7. [1] TestOrder{customer_id=P007, amount=4, discount=50, region=north, product_id=PROD006, date=2024-03-01, status=shipped, priority=urgent, total=600}
8. [1] TestOrder{customer_id=P010, total=255, amount=10, date=2024-03-05, region=south, priority=normal, discount=0, product_id=PROD002, status=pending}
9. [1] TestOrder{customer_id=P001, date=2024-03-10, product_id=PROD007, status=completed, priority=low, amount=1, discount=10, region=north, total=89.99}
10. [1] TestOrder{priority=urgent, amount=1, date=2024-03-15, discount=0, region=east, total=75000, status=refunded, customer_id=P006, product_id=PROD001}

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

1. [1] TestProduct{rating=4.5, brand=TechCorp, keywords=computer, supplier=TechSupply, name=Laptop, category=electronics, available=true, stock=50, price=999.99}
2. [1] TestProduct{price=25.5, stock=200, supplier=TechSupply, category=accessories, brand=TechCorp, available=true, keywords=peripheral, name=Mouse, rating=4}
3. [1] TestProduct{available=false, price=75, rating=3.5, stock=0, category=accessories, keywords=typing, brand=KeyTech, name=Keyboard, supplier=KeySupply}
4. [1] TestProduct{rating=4.8, supplier=ScreenSupply, category=electronics, available=true, brand=ScreenPro, price=299.99, keywords=display, stock=30, name=Monitor}
5. [1] TestProduct{available=false, brand=OldTech, stock=0, name=OldKeyboard, price=8.5, rating=2, keywords=obsolete, category=accessories, supplier=OldSupply}
6. [1] TestProduct{rating=4.6, stock=75, brand=AudioMax, name=Headphones, category=audio, available=true, keywords=sound, supplier=AudioSupply, price=150}
7. [1] TestProduct{category=electronics, available=true, brand=CamTech, supplier=CamSupply, name=Webcam, rating=3.8, stock=25, price=89.99, keywords=video}

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

1. [1] TestPerson{score=8.5, tags=junior, salary=45000, active=true, department=sales, level=2, status=active, name=Alice, age=25}
2. [1] TestPerson{score=9.2, tags=senior, name=Bob, status=active, level=5, age=35, salary=75000, active=true, department=engineering}
3. [1] TestPerson{active=false, tags=intern, salary=0, status=inactive, department=hr, level=1, age=16, score=6, name=Charlie}
4. [1] TestPerson{department=marketing, salary=85000, tags=manager, status=active, age=45, active=true, score=7.8, level=7, name=Diana}
5. [1] TestPerson{name=Eve, tags=employee, level=3, salary=55000, department=sales, score=8, status=inactive, age=30, active=false}
6. [1] TestPerson{name=Frank, score=0, tags=test, level=1, status=active, department=qa, age=0, salary=-5000, active=true}
7. [1] TestPerson{department=management, level=9, name=Grace, score=10, salary=95000, tags=executive, age=65, active=true, status=active}
8. [1] TestPerson{active=false, tags=junior, salary=25000, score=5.5, department=support, name=Henry, status=inactive, level=1, age=18}
9. [1] TestPerson{age=40, active=true, score=8.7, tags=senior, name=Ivy, salary=68000, status=active, department=engineering, level=6}
10. [1] TestPerson{name=X, age=22, status=active, level=1, salary=28000, active=true, score=6.5, tags=temp, department=intern}

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

1. [1] TestOrder{status=pending, priority=normal, customer_id=P001, date=2024-01-15, discount=50, region=north, product_id=PROD001, amount=2, total=1999.98}
2. [1] TestOrder{priority=low, customer_id=P002, date=2024-01-20, region=south, amount=1, total=25.5, status=confirmed, discount=0, product_id=PROD002}
3. [1] TestOrder{status=shipped, customer_id=P001, discount=15, region=north, amount=3, date=2024-02-01, priority=high, product_id=PROD003, total=225}
4. [1] TestOrder{priority=normal, product_id=PROD004, status=delivered, region=east, amount=1, total=299.99, discount=0, customer_id=P004, date=2024-02-05}
5. [1] TestOrder{customer_id=P002, total=999.99, date=2024-02-10, discount=100, status=confirmed, region=south, product_id=PROD001, amount=1, priority=high}
6. [1] TestOrder{priority=low, region=west, product_id=PROD005, customer_id=P005, amount=2, total=999.98, status=cancelled, date=2024-02-15, discount=0}
7. [1] TestOrder{status=shipped, priority=urgent, total=600, customer_id=P007, amount=4, discount=50, region=north, product_id=PROD006, date=2024-03-01}
8. [1] TestOrder{customer_id=P010, total=255, amount=10, date=2024-03-05, region=south, priority=normal, discount=0, product_id=PROD002, status=pending}
9. [1] TestOrder{priority=low, amount=1, discount=10, region=north, total=89.99, customer_id=P001, date=2024-03-10, product_id=PROD007, status=completed}
10. [1] TestOrder{product_id=PROD001, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east, total=75000, status=refunded, customer_id=P006}

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

1. [1] TestPerson{salary=45000, active=true, department=sales, level=2, status=active, name=Alice, age=25, score=8.5, tags=junior}
2. [1] TestPerson{age=35, salary=75000, active=true, department=engineering, score=9.2, tags=senior, name=Bob, status=active, level=5}
3. [1] TestPerson{age=16, score=6, name=Charlie, active=false, tags=intern, salary=0, status=inactive, department=hr, level=1}
4. [1] TestPerson{department=marketing, salary=85000, tags=manager, status=active, age=45, active=true, score=7.8, level=7, name=Diana}
5. [1] TestPerson{name=Eve, tags=employee, level=3, salary=55000, department=sales, score=8, status=inactive, age=30, active=false}
6. [1] TestPerson{status=active, department=qa, age=0, salary=-5000, active=true, name=Frank, score=0, tags=test, level=1}
7. [1] TestPerson{name=Grace, score=10, salary=95000, tags=executive, age=65, active=true, status=active, department=management, level=9}
8. [1] TestPerson{department=support, name=Henry, status=inactive, level=1, age=18, active=false, tags=junior, salary=25000, score=5.5}
9. [1] TestPerson{status=active, department=engineering, level=6, age=40, active=true, score=8.7, tags=senior, name=Ivy, salary=68000}
10. [1] TestPerson{salary=28000, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, status=active, level=1}

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

1. [1] TestOrder{amount=2, total=1999.98, status=pending, priority=normal, customer_id=P001, date=2024-01-15, discount=50, region=north, product_id=PROD001}
2. [1] TestOrder{priority=low, customer_id=P002, date=2024-01-20, region=south, amount=1, total=25.5, status=confirmed, discount=0, product_id=PROD002}
3. [1] TestOrder{status=shipped, customer_id=P001, discount=15, region=north, amount=3, date=2024-02-01, priority=high, product_id=PROD003, total=225}
4. [1] TestOrder{region=east, amount=1, total=299.99, discount=0, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, status=delivered}
5. [1] TestOrder{status=confirmed, region=south, product_id=PROD001, amount=1, priority=high, customer_id=P002, total=999.99, date=2024-02-10, discount=100}
6. [1] TestOrder{amount=2, total=999.98, status=cancelled, date=2024-02-15, discount=0, priority=low, region=west, product_id=PROD005, customer_id=P005}
7. [1] TestOrder{total=600, customer_id=P007, amount=4, discount=50, region=north, product_id=PROD006, date=2024-03-01, status=shipped, priority=urgent}
8. [1] TestOrder{customer_id=P010, total=255, amount=10, date=2024-03-05, region=south, priority=normal, discount=0, product_id=PROD002, status=pending}
9. [1] TestOrder{priority=low, amount=1, discount=10, region=north, total=89.99, customer_id=P001, date=2024-03-10, product_id=PROD007, status=completed}
10. [1] TestOrder{total=75000, status=refunded, customer_id=P006, product_id=PROD001, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

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

1. [1] TestPerson{tags=junior, salary=45000, active=true, department=sales, level=2, status=active, name=Alice, age=25, score=8.5}
2. [1] TestPerson{department=engineering, score=9.2, tags=senior, name=Bob, status=active, level=5, age=35, salary=75000, active=true}
3. [1] TestPerson{age=16, score=6, name=Charlie, active=false, tags=intern, salary=0, status=inactive, department=hr, level=1}
4. [1] TestPerson{score=7.8, level=7, name=Diana, department=marketing, salary=85000, tags=manager, status=active, age=45, active=true}
5. [1] TestPerson{department=sales, score=8, status=inactive, age=30, active=false, name=Eve, tags=employee, level=3, salary=55000}
6. [1] TestPerson{status=active, department=qa, age=0, salary=-5000, active=true, name=Frank, score=0, tags=test, level=1}
7. [1] TestPerson{name=Grace, score=10, salary=95000, tags=executive, age=65, active=true, status=active, department=management, level=9}
8. [1] TestPerson{salary=25000, score=5.5, department=support, name=Henry, status=inactive, level=1, age=18, active=false, tags=junior}
9. [1] TestPerson{status=active, department=engineering, level=6, age=40, active=true, score=8.7, tags=senior, name=Ivy, salary=68000}
10. [1] TestPerson{department=intern, name=X, age=22, status=active, level=1, salary=28000, active=true, score=6.5, tags=temp}

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

1. [1] TestProduct{brand=TechCorp, keywords=computer, supplier=TechSupply, name=Laptop, category=electronics, available=true, stock=50, price=999.99, rating=4.5}
2. [1] TestProduct{keywords=peripheral, name=Mouse, rating=4, price=25.5, stock=200, supplier=TechSupply, category=accessories, brand=TechCorp, available=true}
3. [1] TestProduct{rating=3.5, stock=0, category=accessories, keywords=typing, brand=KeyTech, name=Keyboard, supplier=KeySupply, available=false, price=75}
4. [1] TestProduct{name=Monitor, rating=4.8, supplier=ScreenSupply, category=electronics, available=true, brand=ScreenPro, price=299.99, keywords=display, stock=30}
5. [1] TestProduct{available=false, brand=OldTech, stock=0, name=OldKeyboard, price=8.5, rating=2, keywords=obsolete, category=accessories, supplier=OldSupply}
6. [1] TestProduct{available=true, keywords=sound, supplier=AudioSupply, price=150, rating=4.6, stock=75, brand=AudioMax, name=Headphones, category=audio}
7. [1] TestProduct{name=Webcam, rating=3.8, stock=25, price=89.99, keywords=video, category=electronics, available=true, brand=CamTech, supplier=CamSupply}

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

1. [1] TestPerson{name=Alice, age=25, score=8.5, tags=junior, salary=45000, active=true, department=sales, level=2, status=active}
2. [1] TestPerson{score=9.2, tags=senior, name=Bob, status=active, level=5, age=35, salary=75000, active=true, department=engineering}
3. [1] TestPerson{tags=intern, salary=0, status=inactive, department=hr, level=1, age=16, score=6, name=Charlie, active=false}
4. [1] TestPerson{score=7.8, level=7, name=Diana, department=marketing, salary=85000, tags=manager, status=active, age=45, active=true}
5. [1] TestPerson{score=8, status=inactive, age=30, active=false, name=Eve, tags=employee, level=3, salary=55000, department=sales}
6. [1] TestPerson{active=true, name=Frank, score=0, tags=test, level=1, status=active, department=qa, age=0, salary=-5000}
7. [1] TestPerson{level=9, name=Grace, score=10, salary=95000, tags=executive, age=65, active=true, status=active, department=management}
8. [1] TestPerson{name=Henry, status=inactive, level=1, age=18, active=false, tags=junior, salary=25000, score=5.5, department=support}
9. [1] TestPerson{status=active, department=engineering, level=6, age=40, active=true, score=8.7, tags=senior, name=Ivy, salary=68000}
10. [1] TestPerson{salary=28000, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, status=active, level=1}

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

1. [1] TestOrder{status=pending, priority=normal, customer_id=P001, date=2024-01-15, discount=50, region=north, product_id=PROD001, amount=2, total=1999.98}
2. [1] TestOrder{region=south, amount=1, total=25.5, status=confirmed, discount=0, product_id=PROD002, priority=low, customer_id=P002, date=2024-01-20}
3. [1] TestOrder{date=2024-02-01, priority=high, product_id=PROD003, total=225, status=shipped, customer_id=P001, discount=15, region=north, amount=3}
4. [1] TestOrder{discount=0, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, status=delivered, region=east, amount=1, total=299.99}
5. [1] TestOrder{status=confirmed, region=south, product_id=PROD001, amount=1, priority=high, customer_id=P002, total=999.99, date=2024-02-10, discount=100}
6. [1] TestOrder{amount=2, total=999.98, status=cancelled, date=2024-02-15, discount=0, priority=low, region=west, product_id=PROD005, customer_id=P005}
7. [1] TestOrder{amount=4, discount=50, region=north, product_id=PROD006, date=2024-03-01, status=shipped, priority=urgent, total=600, customer_id=P007}
8. [1] TestOrder{amount=10, date=2024-03-05, region=south, priority=normal, discount=0, product_id=PROD002, status=pending, customer_id=P010, total=255}
9. [1] TestOrder{discount=10, region=north, total=89.99, customer_id=P001, date=2024-03-10, product_id=PROD007, status=completed, priority=low, amount=1}
10. [1] TestOrder{region=east, total=75000, status=refunded, customer_id=P006, product_id=PROD001, priority=urgent, amount=1, date=2024-03-15, discount=0}

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

1. [1] TestPerson{name=Alice, age=25, score=8.5, tags=junior, salary=45000, active=true, department=sales, level=2, status=active}
2. [1] TestPerson{name=Bob, status=active, level=5, age=35, salary=75000, active=true, department=engineering, score=9.2, tags=senior}
3. [1] TestPerson{age=16, score=6, name=Charlie, active=false, tags=intern, salary=0, status=inactive, department=hr, level=1}
4. [1] TestPerson{status=active, age=45, active=true, score=7.8, level=7, name=Diana, department=marketing, salary=85000, tags=manager}
5. [1] TestPerson{age=30, active=false, name=Eve, tags=employee, level=3, salary=55000, department=sales, score=8, status=inactive}
6. [1] TestPerson{status=active, department=qa, age=0, salary=-5000, active=true, name=Frank, score=0, tags=test, level=1}
7. [1] TestPerson{salary=95000, tags=executive, age=65, active=true, status=active, department=management, level=9, name=Grace, score=10}
8. [1] TestPerson{name=Henry, status=inactive, level=1, age=18, active=false, tags=junior, salary=25000, score=5.5, department=support}
9. [1] TestPerson{department=engineering, level=6, age=40, active=true, score=8.7, tags=senior, name=Ivy, salary=68000, status=active}
10. [1] TestPerson{department=intern, name=X, age=22, status=active, level=1, salary=28000, active=true, score=6.5, tags=temp}
11. [1] TestOrder{region=north, product_id=PROD001, amount=2, total=1999.98, status=pending, priority=normal, customer_id=P001, date=2024-01-15, discount=50}
12. [1] TestOrder{region=south, amount=1, total=25.5, status=confirmed, discount=0, product_id=PROD002, priority=low, customer_id=P002, date=2024-01-20}
13. [1] TestOrder{priority=high, product_id=PROD003, total=225, status=shipped, customer_id=P001, discount=15, region=north, amount=3, date=2024-02-01}
14. [1] TestOrder{discount=0, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, status=delivered, region=east, amount=1, total=299.99}
15. [1] TestOrder{status=confirmed, region=south, product_id=PROD001, amount=1, priority=high, customer_id=P002, total=999.99, date=2024-02-10, discount=100}
16. [1] TestOrder{priority=low, region=west, product_id=PROD005, customer_id=P005, amount=2, total=999.98, status=cancelled, date=2024-02-15, discount=0}
17. [1] TestOrder{customer_id=P007, amount=4, discount=50, region=north, product_id=PROD006, date=2024-03-01, status=shipped, priority=urgent, total=600}
18. [1] TestOrder{status=pending, customer_id=P010, total=255, amount=10, date=2024-03-05, region=south, priority=normal, discount=0, product_id=PROD002}
19. [1] TestOrder{customer_id=P001, date=2024-03-10, product_id=PROD007, status=completed, priority=low, amount=1, discount=10, region=north, total=89.99}
20. [1] TestOrder{product_id=PROD001, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east, total=75000, status=refunded, customer_id=P006}

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

1. [1] TestPerson{salary=45000, active=true, department=sales, level=2, status=active, name=Alice, age=25, score=8.5, tags=junior}
2. [1] TestPerson{name=Bob, status=active, level=5, age=35, salary=75000, active=true, department=engineering, score=9.2, tags=senior}
3. [1] TestPerson{age=16, score=6, name=Charlie, active=false, tags=intern, salary=0, status=inactive, department=hr, level=1}
4. [1] TestPerson{level=7, name=Diana, department=marketing, salary=85000, tags=manager, status=active, age=45, active=true, score=7.8}
5. [1] TestPerson{score=8, status=inactive, age=30, active=false, name=Eve, tags=employee, level=3, salary=55000, department=sales}
6. [1] TestPerson{tags=test, level=1, status=active, department=qa, age=0, salary=-5000, active=true, name=Frank, score=0}
7. [1] TestPerson{name=Grace, score=10, salary=95000, tags=executive, age=65, active=true, status=active, department=management, level=9}
8. [1] TestPerson{name=Henry, status=inactive, level=1, age=18, active=false, tags=junior, salary=25000, score=5.5, department=support}
9. [1] TestPerson{name=Ivy, salary=68000, status=active, department=engineering, level=6, age=40, active=true, score=8.7, tags=senior}
10. [1] TestPerson{age=22, status=active, level=1, salary=28000, active=true, score=6.5, tags=temp, department=intern, name=X}

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

1. [1] TestPerson{department=sales, level=2, status=active, name=Alice, age=25, score=8.5, tags=junior, salary=45000, active=true}
2. [1] TestPerson{name=Bob, status=active, level=5, age=35, salary=75000, active=true, department=engineering, score=9.2, tags=senior}
3. [1] TestPerson{age=16, score=6, name=Charlie, active=false, tags=intern, salary=0, status=inactive, department=hr, level=1}
4. [1] TestPerson{tags=manager, status=active, age=45, active=true, score=7.8, level=7, name=Diana, department=marketing, salary=85000}
5. [1] TestPerson{status=inactive, age=30, active=false, name=Eve, tags=employee, level=3, salary=55000, department=sales, score=8}
6. [1] TestPerson{active=true, name=Frank, score=0, tags=test, level=1, status=active, department=qa, age=0, salary=-5000}
7. [1] TestPerson{tags=executive, age=65, active=true, status=active, department=management, level=9, name=Grace, score=10, salary=95000}
8. [1] TestPerson{tags=junior, salary=25000, score=5.5, department=support, name=Henry, status=inactive, level=1, age=18, active=false}
9. [1] TestPerson{tags=senior, name=Ivy, salary=68000, status=active, department=engineering, level=6, age=40, active=true, score=8.7}
10. [1] TestPerson{name=X, age=22, status=active, level=1, salary=28000, active=true, score=6.5, tags=temp, department=intern}
11. [1] TestOrder{amount=2, total=1999.98, status=pending, priority=normal, customer_id=P001, date=2024-01-15, discount=50, region=north, product_id=PROD001}
12. [1] TestOrder{customer_id=P002, date=2024-01-20, region=south, amount=1, total=25.5, status=confirmed, discount=0, product_id=PROD002, priority=low}
13. [1] TestOrder{discount=15, region=north, amount=3, date=2024-02-01, priority=high, product_id=PROD003, total=225, status=shipped, customer_id=P001}
14. [1] TestOrder{discount=0, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, status=delivered, region=east, amount=1, total=299.99}
15. [1] TestOrder{amount=1, priority=high, customer_id=P002, total=999.99, date=2024-02-10, discount=100, status=confirmed, region=south, product_id=PROD001}
16. [1] TestOrder{amount=2, total=999.98, status=cancelled, date=2024-02-15, discount=0, priority=low, region=west, product_id=PROD005, customer_id=P005}
17. [1] TestOrder{product_id=PROD006, date=2024-03-01, status=shipped, priority=urgent, total=600, customer_id=P007, amount=4, discount=50, region=north}
18. [1] TestOrder{product_id=PROD002, status=pending, customer_id=P010, total=255, amount=10, date=2024-03-05, region=south, priority=normal, discount=0}
19. [1] TestOrder{customer_id=P001, date=2024-03-10, product_id=PROD007, status=completed, priority=low, amount=1, discount=10, region=north, total=89.99}
20. [1] TestOrder{product_id=PROD001, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east, total=75000, status=refunded, customer_id=P006}

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

1. [1] TestPerson{age=25, score=8.5, tags=junior, salary=45000, active=true, department=sales, level=2, status=active, name=Alice}
2. [1] TestPerson{status=active, level=5, age=35, salary=75000, active=true, department=engineering, score=9.2, tags=senior, name=Bob}
3. [1] TestPerson{age=16, score=6, name=Charlie, active=false, tags=intern, salary=0, status=inactive, department=hr, level=1}
4. [1] TestPerson{age=45, active=true, score=7.8, level=7, name=Diana, department=marketing, salary=85000, tags=manager, status=active}
5. [1] TestPerson{name=Eve, tags=employee, level=3, salary=55000, department=sales, score=8, status=inactive, age=30, active=false}
6. [1] TestPerson{level=1, status=active, department=qa, age=0, salary=-5000, active=true, name=Frank, score=0, tags=test}
7. [1] TestPerson{status=active, department=management, level=9, name=Grace, score=10, salary=95000, tags=executive, age=65, active=true}
8. [1] TestPerson{salary=25000, score=5.5, department=support, name=Henry, status=inactive, level=1, age=18, active=false, tags=junior}
9. [1] TestPerson{status=active, department=engineering, level=6, age=40, active=true, score=8.7, tags=senior, name=Ivy, salary=68000}
10. [1] TestPerson{salary=28000, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, status=active, level=1}

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

1. [1] TestOrder{discount=50, region=north, product_id=PROD001, amount=2, total=1999.98, status=pending, priority=normal, customer_id=P001, date=2024-01-15}
2. [1] TestOrder{amount=1, total=25.5, status=confirmed, discount=0, product_id=PROD002, priority=low, customer_id=P002, date=2024-01-20, region=south}
3. [1] TestOrder{discount=15, region=north, amount=3, date=2024-02-01, priority=high, product_id=PROD003, total=225, status=shipped, customer_id=P001}
4. [1] TestOrder{discount=0, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, status=delivered, region=east, amount=1, total=299.99}
5. [1] TestOrder{customer_id=P002, total=999.99, date=2024-02-10, discount=100, status=confirmed, region=south, product_id=PROD001, amount=1, priority=high}
6. [1] TestOrder{amount=2, total=999.98, status=cancelled, date=2024-02-15, discount=0, priority=low, region=west, product_id=PROD005, customer_id=P005}
7. [1] TestOrder{customer_id=P007, amount=4, discount=50, region=north, product_id=PROD006, date=2024-03-01, status=shipped, priority=urgent, total=600}
8. [1] TestOrder{product_id=PROD002, status=pending, customer_id=P010, total=255, amount=10, date=2024-03-05, region=south, priority=normal, discount=0}
9. [1] TestOrder{priority=low, amount=1, discount=10, region=north, total=89.99, customer_id=P001, date=2024-03-10, product_id=PROD007, status=completed}
10. [1] TestOrder{product_id=PROD001, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east, total=75000, status=refunded, customer_id=P006}

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

1. [1] TestPerson{salary=45000, active=true, department=sales, level=2, status=active, name=Alice, age=25, score=8.5, tags=junior}
2. [1] TestPerson{level=5, age=35, salary=75000, active=true, department=engineering, score=9.2, tags=senior, name=Bob, status=active}
3. [1] TestPerson{salary=0, status=inactive, department=hr, level=1, age=16, score=6, name=Charlie, active=false, tags=intern}
4. [1] TestPerson{department=marketing, salary=85000, tags=manager, status=active, age=45, active=true, score=7.8, level=7, name=Diana}
5. [1] TestPerson{name=Eve, tags=employee, level=3, salary=55000, department=sales, score=8, status=inactive, age=30, active=false}
6. [1] TestPerson{tags=test, level=1, status=active, department=qa, age=0, salary=-5000, active=true, name=Frank, score=0}
7. [1] TestPerson{salary=95000, tags=executive, age=65, active=true, status=active, department=management, level=9, name=Grace, score=10}
8. [1] TestPerson{salary=25000, score=5.5, department=support, name=Henry, status=inactive, level=1, age=18, active=false, tags=junior}
9. [1] TestPerson{age=40, active=true, score=8.7, tags=senior, name=Ivy, salary=68000, status=active, department=engineering, level=6}
10. [1] TestPerson{active=true, score=6.5, tags=temp, department=intern, name=X, age=22, status=active, level=1, salary=28000}

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
