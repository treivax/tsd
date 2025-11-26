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

1. [1] TestPerson{level=2, active=true, department=sales, salary=45000, tags=junior, name=Alice, status=active, age=25, score=8.5}
2. [1] TestPerson{tags=senior, status=active, department=engineering, age=35, active=true, score=9.2, level=5, name=Bob, salary=75000}
3. [1] TestPerson{age=16, status=inactive, department=hr, level=1, name=Charlie, active=false, score=6, tags=intern, salary=0}
4. [1] TestPerson{active=true, tags=manager, score=7.8, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45}
5. [1] TestPerson{status=inactive, level=3, name=Eve, score=8, department=sales, age=30, salary=55000, active=false, tags=employee}
6. [1] TestPerson{level=1, name=Frank, age=0, score=0, salary=-5000, active=true, status=active, department=qa, tags=test}
7. [1] TestPerson{status=active, department=management, active=true, tags=executive, level=9, salary=95000, age=65, score=10, name=Grace}
8. [1] TestPerson{name=Henry, age=18, score=5.5, tags=junior, active=false, status=inactive, department=support, salary=25000, level=1}
9. [1] TestPerson{level=6, age=40, score=8.7, name=Ivy, active=true, department=engineering, salary=68000, tags=senior, status=active}
10. [1] TestPerson{name=X, salary=28000, status=active, department=intern, age=22, score=6.5, level=1, active=true, tags=temp}

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

1. [1] TestOrder{product_id=PROD001, date=2024-01-15, status=pending, amount=2, total=1999.98, priority=normal, customer_id=P001, discount=50, region=north}
2. [1] TestOrder{customer_id=P002, status=confirmed, priority=low, amount=1, date=2024-01-20, product_id=PROD002, total=25.5, region=south, discount=0}
3. [1] TestOrder{status=shipped, discount=15, total=225, date=2024-02-01, priority=high, customer_id=P001, region=north, product_id=PROD003, amount=3}
4. [1] TestOrder{customer_id=P004, product_id=PROD004, total=299.99, date=2024-02-05, amount=1, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{amount=1, total=999.99, status=confirmed, discount=100, region=south, customer_id=P002, product_id=PROD001, date=2024-02-10, priority=high}
6. [1] TestOrder{date=2024-02-15, region=west, product_id=PROD005, discount=0, amount=2, status=cancelled, priority=low, customer_id=P005, total=999.98}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, date=2024-03-01, discount=50, region=north, status=shipped, amount=4, priority=urgent, total=600}
8. [1] TestOrder{status=pending, priority=normal, customer_id=P010, product_id=PROD002, discount=0, region=south, amount=10, date=2024-03-05, total=255}
9. [1] TestOrder{date=2024-03-10, status=completed, amount=1, customer_id=P001, total=89.99, discount=10, region=north, product_id=PROD007, priority=low}
10. [1] TestOrder{product_id=PROD001, amount=1, status=refunded, discount=0, date=2024-03-15, priority=urgent, region=east, total=75000, customer_id=P006}

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

1. [1] TestPerson{name=Alice, status=active, age=25, score=8.5, level=2, active=true, department=sales, salary=45000, tags=junior}
2. [1] TestPerson{level=5, name=Bob, salary=75000, tags=senior, status=active, department=engineering, age=35, active=true, score=9.2}
3. [1] TestPerson{name=Charlie, active=false, score=6, tags=intern, salary=0, age=16, status=inactive, department=hr, level=1}
4. [1] TestPerson{level=7, name=Diana, age=45, active=true, tags=manager, score=7.8, status=active, department=marketing, salary=85000}
5. [1] TestPerson{department=sales, age=30, salary=55000, active=false, tags=employee, status=inactive, level=3, name=Eve, score=8}
6. [1] TestPerson{age=0, score=0, salary=-5000, active=true, status=active, department=qa, tags=test, level=1, name=Frank}
7. [1] TestPerson{department=management, active=true, tags=executive, level=9, salary=95000, age=65, score=10, name=Grace, status=active}
8. [1] TestPerson{salary=25000, level=1, name=Henry, age=18, score=5.5, tags=junior, active=false, status=inactive, department=support}
9. [1] TestPerson{salary=68000, tags=senior, status=active, level=6, age=40, score=8.7, name=Ivy, active=true, department=engineering}
10. [1] TestPerson{tags=temp, name=X, salary=28000, status=active, department=intern, age=22, score=6.5, level=1, active=true}

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

1. [1] TestOrder{total=1999.98, priority=normal, customer_id=P001, discount=50, region=north, product_id=PROD001, date=2024-01-15, status=pending, amount=2}
2. [1] TestOrder{customer_id=P002, status=confirmed, priority=low, amount=1, date=2024-01-20, product_id=PROD002, total=25.5, region=south, discount=0}
3. [1] TestOrder{customer_id=P001, region=north, product_id=PROD003, amount=3, status=shipped, discount=15, total=225, date=2024-02-01, priority=high}
4. [1] TestOrder{status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, total=299.99, date=2024-02-05, amount=1}
5. [1] TestOrder{product_id=PROD001, date=2024-02-10, priority=high, amount=1, total=999.99, status=confirmed, discount=100, region=south, customer_id=P002}
6. [1] TestOrder{priority=low, customer_id=P005, total=999.98, date=2024-02-15, region=west, product_id=PROD005, discount=0, amount=2, status=cancelled}
7. [1] TestOrder{status=shipped, amount=4, priority=urgent, total=600, customer_id=P007, product_id=PROD006, date=2024-03-01, discount=50, region=north}
8. [1] TestOrder{date=2024-03-05, total=255, status=pending, priority=normal, customer_id=P010, product_id=PROD002, discount=0, region=south, amount=10}
9. [1] TestOrder{date=2024-03-10, status=completed, amount=1, customer_id=P001, total=89.99, discount=10, region=north, product_id=PROD007, priority=low}
10. [1] TestOrder{product_id=PROD001, amount=1, status=refunded, discount=0, date=2024-03-15, priority=urgent, region=east, total=75000, customer_id=P006}

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

1. [1] TestProduct{stock=50, supplier=TechSupply, name=Laptop, available=true, brand=TechCorp, price=999.99, category=electronics, keywords=computer, rating=4.5}
2. [1] TestProduct{name=Mouse, category=accessories, available=true, rating=4, price=25.5, brand=TechCorp, keywords=peripheral, stock=200, supplier=TechSupply}
3. [1] TestProduct{category=accessories, price=75, available=false, keywords=typing, supplier=KeySupply, name=Keyboard, rating=3.5, brand=KeyTech, stock=0}
4. [1] TestProduct{supplier=ScreenSupply, name=Monitor, category=electronics, price=299.99, rating=4.8, keywords=display, brand=ScreenPro, available=true, stock=30}
5. [1] TestProduct{rating=2, stock=0, name=OldKeyboard, category=accessories, price=8.5, available=false, keywords=obsolete, brand=OldTech, supplier=OldSupply}
6. [1] TestProduct{rating=4.6, available=true, keywords=sound, stock=75, name=Headphones, price=150, brand=AudioMax, supplier=AudioSupply, category=audio}
7. [1] TestProduct{keywords=video, supplier=CamSupply, name=Webcam, rating=3.8, brand=CamTech, stock=25, category=electronics, price=89.99, available=true}

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

1. [1] TestPerson{salary=45000, tags=junior, name=Alice, status=active, age=25, score=8.5, level=2, active=true, department=sales}
2. [1] TestPerson{tags=senior, status=active, department=engineering, age=35, active=true, score=9.2, level=5, name=Bob, salary=75000}
3. [1] TestPerson{department=hr, level=1, name=Charlie, active=false, score=6, tags=intern, salary=0, age=16, status=inactive}
4. [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, tags=manager, score=7.8, status=active, department=marketing}
5. [1] TestPerson{status=inactive, level=3, name=Eve, score=8, department=sales, age=30, salary=55000, active=false, tags=employee}
6. [1] TestPerson{salary=-5000, active=true, status=active, department=qa, tags=test, level=1, name=Frank, age=0, score=0}
7. [1] TestPerson{score=10, name=Grace, status=active, department=management, active=true, tags=executive, level=9, salary=95000, age=65}
8. [1] TestPerson{department=support, salary=25000, level=1, name=Henry, age=18, score=5.5, tags=junior, active=false, status=inactive}
9. [1] TestPerson{salary=68000, tags=senior, status=active, level=6, age=40, score=8.7, name=Ivy, active=true, department=engineering}
10. [1] TestPerson{score=6.5, level=1, active=true, tags=temp, name=X, salary=28000, status=active, department=intern, age=22}

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

1. [1] TestOrder{amount=2, total=1999.98, priority=normal, customer_id=P001, discount=50, region=north, product_id=PROD001, date=2024-01-15, status=pending}
2. [1] TestOrder{customer_id=P002, status=confirmed, priority=low, amount=1, date=2024-01-20, product_id=PROD002, total=25.5, region=south, discount=0}
3. [1] TestOrder{amount=3, status=shipped, discount=15, total=225, date=2024-02-01, priority=high, customer_id=P001, region=north, product_id=PROD003}
4. [1] TestOrder{customer_id=P004, product_id=PROD004, total=299.99, date=2024-02-05, amount=1, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{region=south, customer_id=P002, product_id=PROD001, date=2024-02-10, priority=high, amount=1, total=999.99, status=confirmed, discount=100}
6. [1] TestOrder{total=999.98, date=2024-02-15, region=west, product_id=PROD005, discount=0, amount=2, status=cancelled, priority=low, customer_id=P005}
7. [1] TestOrder{total=600, customer_id=P007, product_id=PROD006, date=2024-03-01, discount=50, region=north, status=shipped, amount=4, priority=urgent}
8. [1] TestOrder{total=255, status=pending, priority=normal, customer_id=P010, product_id=PROD002, discount=0, region=south, amount=10, date=2024-03-05}
9. [1] TestOrder{priority=low, date=2024-03-10, status=completed, amount=1, customer_id=P001, total=89.99, discount=10, region=north, product_id=PROD007}
10. [1] TestOrder{customer_id=P006, product_id=PROD001, amount=1, status=refunded, discount=0, date=2024-03-15, priority=urgent, region=east, total=75000}

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

1. [1] TestPerson{age=25, score=8.5, level=2, active=true, department=sales, salary=45000, tags=junior, name=Alice, status=active}
2. [1] TestPerson{score=9.2, level=5, name=Bob, salary=75000, tags=senior, status=active, department=engineering, age=35, active=true}
3. [1] TestPerson{status=inactive, department=hr, level=1, name=Charlie, active=false, score=6, tags=intern, salary=0, age=16}
4. [1] TestPerson{level=7, name=Diana, age=45, active=true, tags=manager, score=7.8, status=active, department=marketing, salary=85000}
5. [1] TestPerson{status=inactive, level=3, name=Eve, score=8, department=sales, age=30, salary=55000, active=false, tags=employee}
6. [1] TestPerson{age=0, score=0, salary=-5000, active=true, status=active, department=qa, tags=test, level=1, name=Frank}
7. [1] TestPerson{salary=95000, age=65, score=10, name=Grace, status=active, department=management, active=true, tags=executive, level=9}
8. [1] TestPerson{name=Henry, age=18, score=5.5, tags=junior, active=false, status=inactive, department=support, salary=25000, level=1}
9. [1] TestPerson{active=true, department=engineering, salary=68000, tags=senior, status=active, level=6, age=40, score=8.7, name=Ivy}
10. [1] TestPerson{name=X, salary=28000, status=active, department=intern, age=22, score=6.5, level=1, active=true, tags=temp}

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

1. [1] TestOrder{status=pending, amount=2, total=1999.98, priority=normal, customer_id=P001, discount=50, region=north, product_id=PROD001, date=2024-01-15}
2. [1] TestOrder{region=south, discount=0, customer_id=P002, status=confirmed, priority=low, amount=1, date=2024-01-20, product_id=PROD002, total=25.5}
3. [1] TestOrder{amount=3, status=shipped, discount=15, total=225, date=2024-02-01, priority=high, customer_id=P001, region=north, product_id=PROD003}
4. [1] TestOrder{status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, total=299.99, date=2024-02-05, amount=1}
5. [1] TestOrder{region=south, customer_id=P002, product_id=PROD001, date=2024-02-10, priority=high, amount=1, total=999.99, status=confirmed, discount=100}
6. [1] TestOrder{total=999.98, date=2024-02-15, region=west, product_id=PROD005, discount=0, amount=2, status=cancelled, priority=low, customer_id=P005}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, date=2024-03-01, discount=50, region=north, status=shipped, amount=4, priority=urgent, total=600}
8. [1] TestOrder{date=2024-03-05, total=255, status=pending, priority=normal, customer_id=P010, product_id=PROD002, discount=0, region=south, amount=10}
9. [1] TestOrder{discount=10, region=north, product_id=PROD007, priority=low, date=2024-03-10, status=completed, amount=1, customer_id=P001, total=89.99}
10. [1] TestOrder{amount=1, status=refunded, discount=0, date=2024-03-15, priority=urgent, region=east, total=75000, customer_id=P006, product_id=PROD001}

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

1. [1] TestPerson{tags=junior, name=Alice, status=active, age=25, score=8.5, level=2, active=true, department=sales, salary=45000}
2. [1] TestPerson{tags=senior, status=active, department=engineering, age=35, active=true, score=9.2, level=5, name=Bob, salary=75000}
3. [1] TestPerson{name=Charlie, active=false, score=6, tags=intern, salary=0, age=16, status=inactive, department=hr, level=1}
4. [1] TestPerson{level=7, name=Diana, age=45, active=true, tags=manager, score=7.8, status=active, department=marketing, salary=85000}
5. [1] TestPerson{status=inactive, level=3, name=Eve, score=8, department=sales, age=30, salary=55000, active=false, tags=employee}
6. [1] TestPerson{status=active, department=qa, tags=test, level=1, name=Frank, age=0, score=0, salary=-5000, active=true}
7. [1] TestPerson{tags=executive, level=9, salary=95000, age=65, score=10, name=Grace, status=active, department=management, active=true}
8. [1] TestPerson{salary=25000, level=1, name=Henry, age=18, score=5.5, tags=junior, active=false, status=inactive, department=support}
9. [1] TestPerson{score=8.7, name=Ivy, active=true, department=engineering, salary=68000, tags=senior, status=active, level=6, age=40}
10. [1] TestPerson{department=intern, age=22, score=6.5, level=1, active=true, tags=temp, name=X, salary=28000, status=active}

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

1. [1] TestProduct{category=electronics, keywords=computer, rating=4.5, stock=50, supplier=TechSupply, name=Laptop, available=true, brand=TechCorp, price=999.99}
2. [1] TestProduct{price=25.5, brand=TechCorp, keywords=peripheral, stock=200, supplier=TechSupply, name=Mouse, category=accessories, available=true, rating=4}
3. [1] TestProduct{category=accessories, price=75, available=false, keywords=typing, supplier=KeySupply, name=Keyboard, rating=3.5, brand=KeyTech, stock=0}
4. [1] TestProduct{keywords=display, brand=ScreenPro, available=true, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics, price=299.99, rating=4.8}
5. [1] TestProduct{price=8.5, available=false, keywords=obsolete, brand=OldTech, supplier=OldSupply, rating=2, stock=0, name=OldKeyboard, category=accessories}
6. [1] TestProduct{category=audio, rating=4.6, available=true, keywords=sound, stock=75, name=Headphones, price=150, brand=AudioMax, supplier=AudioSupply}
7. [1] TestProduct{keywords=video, supplier=CamSupply, name=Webcam, rating=3.8, brand=CamTech, stock=25, category=electronics, price=89.99, available=true}

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

1. [1] TestPerson{level=2, active=true, department=sales, salary=45000, tags=junior, name=Alice, status=active, age=25, score=8.5}
2. [1] TestPerson{tags=senior, status=active, department=engineering, age=35, active=true, score=9.2, level=5, name=Bob, salary=75000}
3. [1] TestPerson{score=6, tags=intern, salary=0, age=16, status=inactive, department=hr, level=1, name=Charlie, active=false}
4. [1] TestPerson{score=7.8, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45, active=true, tags=manager}
5. [1] TestPerson{department=sales, age=30, salary=55000, active=false, tags=employee, status=inactive, level=3, name=Eve, score=8}
6. [1] TestPerson{name=Frank, age=0, score=0, salary=-5000, active=true, status=active, department=qa, tags=test, level=1}
7. [1] TestPerson{score=10, name=Grace, status=active, department=management, active=true, tags=executive, level=9, salary=95000, age=65}
8. [1] TestPerson{department=support, salary=25000, level=1, name=Henry, age=18, score=5.5, tags=junior, active=false, status=inactive}
9. [1] TestPerson{level=6, age=40, score=8.7, name=Ivy, active=true, department=engineering, salary=68000, tags=senior, status=active}
10. [1] TestPerson{age=22, score=6.5, level=1, active=true, tags=temp, name=X, salary=28000, status=active, department=intern}

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

1. [1] TestOrder{amount=2, total=1999.98, priority=normal, customer_id=P001, discount=50, region=north, product_id=PROD001, date=2024-01-15, status=pending}
2. [1] TestOrder{total=25.5, region=south, discount=0, customer_id=P002, status=confirmed, priority=low, amount=1, date=2024-01-20, product_id=PROD002}
3. [1] TestOrder{discount=15, total=225, date=2024-02-01, priority=high, customer_id=P001, region=north, product_id=PROD003, amount=3, status=shipped}
4. [1] TestOrder{customer_id=P004, product_id=PROD004, total=299.99, date=2024-02-05, amount=1, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{status=confirmed, discount=100, region=south, customer_id=P002, product_id=PROD001, date=2024-02-10, priority=high, amount=1, total=999.99}
6. [1] TestOrder{region=west, product_id=PROD005, discount=0, amount=2, status=cancelled, priority=low, customer_id=P005, total=999.98, date=2024-02-15}
7. [1] TestOrder{total=600, customer_id=P007, product_id=PROD006, date=2024-03-01, discount=50, region=north, status=shipped, amount=4, priority=urgent}
8. [1] TestOrder{customer_id=P010, product_id=PROD002, discount=0, region=south, amount=10, date=2024-03-05, total=255, status=pending, priority=normal}
9. [1] TestOrder{discount=10, region=north, product_id=PROD007, priority=low, date=2024-03-10, status=completed, amount=1, customer_id=P001, total=89.99}
10. [1] TestOrder{total=75000, customer_id=P006, product_id=PROD001, amount=1, status=refunded, discount=0, date=2024-03-15, priority=urgent, region=east}

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

1. [1] TestPerson{department=sales, salary=45000, tags=junior, name=Alice, status=active, age=25, score=8.5, level=2, active=true}
2. [1] TestPerson{tags=senior, status=active, department=engineering, age=35, active=true, score=9.2, level=5, name=Bob, salary=75000}
3. [1] TestPerson{name=Charlie, active=false, score=6, tags=intern, salary=0, age=16, status=inactive, department=hr, level=1}
4. [1] TestPerson{age=45, active=true, tags=manager, score=7.8, status=active, department=marketing, salary=85000, level=7, name=Diana}
5. [1] TestPerson{salary=55000, active=false, tags=employee, status=inactive, level=3, name=Eve, score=8, department=sales, age=30}
6. [1] TestPerson{name=Frank, age=0, score=0, salary=-5000, active=true, status=active, department=qa, tags=test, level=1}
7. [1] TestPerson{status=active, department=management, active=true, tags=executive, level=9, salary=95000, age=65, score=10, name=Grace}
8. [1] TestPerson{active=false, status=inactive, department=support, salary=25000, level=1, name=Henry, age=18, score=5.5, tags=junior}
9. [1] TestPerson{name=Ivy, active=true, department=engineering, salary=68000, tags=senior, status=active, level=6, age=40, score=8.7}
10. [1] TestPerson{name=X, salary=28000, status=active, department=intern, age=22, score=6.5, level=1, active=true, tags=temp}
11. [1] TestOrder{amount=2, total=1999.98, priority=normal, customer_id=P001, discount=50, region=north, product_id=PROD001, date=2024-01-15, status=pending}
12. [1] TestOrder{product_id=PROD002, total=25.5, region=south, discount=0, customer_id=P002, status=confirmed, priority=low, amount=1, date=2024-01-20}
13. [1] TestOrder{total=225, date=2024-02-01, priority=high, customer_id=P001, region=north, product_id=PROD003, amount=3, status=shipped, discount=15}
14. [1] TestOrder{status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, total=299.99, date=2024-02-05, amount=1}
15. [1] TestOrder{priority=high, amount=1, total=999.99, status=confirmed, discount=100, region=south, customer_id=P002, product_id=PROD001, date=2024-02-10}
16. [1] TestOrder{product_id=PROD005, discount=0, amount=2, status=cancelled, priority=low, customer_id=P005, total=999.98, date=2024-02-15, region=west}
17. [1] TestOrder{customer_id=P007, product_id=PROD006, date=2024-03-01, discount=50, region=north, status=shipped, amount=4, priority=urgent, total=600}
18. [1] TestOrder{region=south, amount=10, date=2024-03-05, total=255, status=pending, priority=normal, customer_id=P010, product_id=PROD002, discount=0}
19. [1] TestOrder{discount=10, region=north, product_id=PROD007, priority=low, date=2024-03-10, status=completed, amount=1, customer_id=P001, total=89.99}
20. [1] TestOrder{amount=1, status=refunded, discount=0, date=2024-03-15, priority=urgent, region=east, total=75000, customer_id=P006, product_id=PROD001}

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

1. [1] TestPerson{level=2, active=true, department=sales, salary=45000, tags=junior, name=Alice, status=active, age=25, score=8.5}
2. [1] TestPerson{active=true, score=9.2, level=5, name=Bob, salary=75000, tags=senior, status=active, department=engineering, age=35}
3. [1] TestPerson{name=Charlie, active=false, score=6, tags=intern, salary=0, age=16, status=inactive, department=hr, level=1}
4. [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, tags=manager, score=7.8, status=active, department=marketing}
5. [1] TestPerson{status=inactive, level=3, name=Eve, score=8, department=sales, age=30, salary=55000, active=false, tags=employee}
6. [1] TestPerson{status=active, department=qa, tags=test, level=1, name=Frank, age=0, score=0, salary=-5000, active=true}
7. [1] TestPerson{status=active, department=management, active=true, tags=executive, level=9, salary=95000, age=65, score=10, name=Grace}
8. [1] TestPerson{tags=junior, active=false, status=inactive, department=support, salary=25000, level=1, name=Henry, age=18, score=5.5}
9. [1] TestPerson{salary=68000, tags=senior, status=active, level=6, age=40, score=8.7, name=Ivy, active=true, department=engineering}
10. [1] TestPerson{name=X, salary=28000, status=active, department=intern, age=22, score=6.5, level=1, active=true, tags=temp}

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

1. [1] TestPerson{level=2, active=true, department=sales, salary=45000, tags=junior, name=Alice, status=active, age=25, score=8.5}
2. [1] TestPerson{level=5, name=Bob, salary=75000, tags=senior, status=active, department=engineering, age=35, active=true, score=9.2}
3. [1] TestPerson{status=inactive, department=hr, level=1, name=Charlie, active=false, score=6, tags=intern, salary=0, age=16}
4. [1] TestPerson{level=7, name=Diana, age=45, active=true, tags=manager, score=7.8, status=active, department=marketing, salary=85000}
5. [1] TestPerson{score=8, department=sales, age=30, salary=55000, active=false, tags=employee, status=inactive, level=3, name=Eve}
6. [1] TestPerson{age=0, score=0, salary=-5000, active=true, status=active, department=qa, tags=test, level=1, name=Frank}
7. [1] TestPerson{salary=95000, age=65, score=10, name=Grace, status=active, department=management, active=true, tags=executive, level=9}
8. [1] TestPerson{tags=junior, active=false, status=inactive, department=support, salary=25000, level=1, name=Henry, age=18, score=5.5}
9. [1] TestPerson{level=6, age=40, score=8.7, name=Ivy, active=true, department=engineering, salary=68000, tags=senior, status=active}
10. [1] TestPerson{name=X, salary=28000, status=active, department=intern, age=22, score=6.5, level=1, active=true, tags=temp}
11. [1] TestOrder{amount=2, total=1999.98, priority=normal, customer_id=P001, discount=50, region=north, product_id=PROD001, date=2024-01-15, status=pending}
12. [1] TestOrder{product_id=PROD002, total=25.5, region=south, discount=0, customer_id=P002, status=confirmed, priority=low, amount=1, date=2024-01-20}
13. [1] TestOrder{product_id=PROD003, amount=3, status=shipped, discount=15, total=225, date=2024-02-01, priority=high, customer_id=P001, region=north}
14. [1] TestOrder{customer_id=P004, product_id=PROD004, total=299.99, date=2024-02-05, amount=1, status=delivered, priority=normal, discount=0, region=east}
15. [1] TestOrder{region=south, customer_id=P002, product_id=PROD001, date=2024-02-10, priority=high, amount=1, total=999.99, status=confirmed, discount=100}
16. [1] TestOrder{date=2024-02-15, region=west, product_id=PROD005, discount=0, amount=2, status=cancelled, priority=low, customer_id=P005, total=999.98}
17. [1] TestOrder{customer_id=P007, product_id=PROD006, date=2024-03-01, discount=50, region=north, status=shipped, amount=4, priority=urgent, total=600}
18. [1] TestOrder{product_id=PROD002, discount=0, region=south, amount=10, date=2024-03-05, total=255, status=pending, priority=normal, customer_id=P010}
19. [1] TestOrder{discount=10, region=north, product_id=PROD007, priority=low, date=2024-03-10, status=completed, amount=1, customer_id=P001, total=89.99}
20. [1] TestOrder{priority=urgent, region=east, total=75000, customer_id=P006, product_id=PROD001, amount=1, status=refunded, discount=0, date=2024-03-15}

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

1. [1] TestPerson{level=2, active=true, department=sales, salary=45000, tags=junior, name=Alice, status=active, age=25, score=8.5}
2. [1] TestPerson{name=Bob, salary=75000, tags=senior, status=active, department=engineering, age=35, active=true, score=9.2, level=5}
3. [1] TestPerson{level=1, name=Charlie, active=false, score=6, tags=intern, salary=0, age=16, status=inactive, department=hr}
4. [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, tags=manager, score=7.8, status=active, department=marketing}
5. [1] TestPerson{status=inactive, level=3, name=Eve, score=8, department=sales, age=30, salary=55000, active=false, tags=employee}
6. [1] TestPerson{age=0, score=0, salary=-5000, active=true, status=active, department=qa, tags=test, level=1, name=Frank}
7. [1] TestPerson{status=active, department=management, active=true, tags=executive, level=9, salary=95000, age=65, score=10, name=Grace}
8. [1] TestPerson{salary=25000, level=1, name=Henry, age=18, score=5.5, tags=junior, active=false, status=inactive, department=support}
9. [1] TestPerson{name=Ivy, active=true, department=engineering, salary=68000, tags=senior, status=active, level=6, age=40, score=8.7}
10. [1] TestPerson{name=X, salary=28000, status=active, department=intern, age=22, score=6.5, level=1, active=true, tags=temp}

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

1. [1] TestOrder{customer_id=P001, discount=50, region=north, product_id=PROD001, date=2024-01-15, status=pending, amount=2, total=1999.98, priority=normal}
2. [1] TestOrder{customer_id=P002, status=confirmed, priority=low, amount=1, date=2024-01-20, product_id=PROD002, total=25.5, region=south, discount=0}
3. [1] TestOrder{product_id=PROD003, amount=3, status=shipped, discount=15, total=225, date=2024-02-01, priority=high, customer_id=P001, region=north}
4. [1] TestOrder{priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, total=299.99, date=2024-02-05, amount=1, status=delivered}
5. [1] TestOrder{status=confirmed, discount=100, region=south, customer_id=P002, product_id=PROD001, date=2024-02-10, priority=high, amount=1, total=999.99}
6. [1] TestOrder{total=999.98, date=2024-02-15, region=west, product_id=PROD005, discount=0, amount=2, status=cancelled, priority=low, customer_id=P005}
7. [1] TestOrder{status=shipped, amount=4, priority=urgent, total=600, customer_id=P007, product_id=PROD006, date=2024-03-01, discount=50, region=north}
8. [1] TestOrder{region=south, amount=10, date=2024-03-05, total=255, status=pending, priority=normal, customer_id=P010, product_id=PROD002, discount=0}
9. [1] TestOrder{status=completed, amount=1, customer_id=P001, total=89.99, discount=10, region=north, product_id=PROD007, priority=low, date=2024-03-10}
10. [1] TestOrder{date=2024-03-15, priority=urgent, region=east, total=75000, customer_id=P006, product_id=PROD001, amount=1, status=refunded, discount=0}

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

1. [1] TestPerson{level=2, active=true, department=sales, salary=45000, tags=junior, name=Alice, status=active, age=25, score=8.5}
2. [1] TestPerson{level=5, name=Bob, salary=75000, tags=senior, status=active, department=engineering, age=35, active=true, score=9.2}
3. [1] TestPerson{name=Charlie, active=false, score=6, tags=intern, salary=0, age=16, status=inactive, department=hr, level=1}
4. [1] TestPerson{level=7, name=Diana, age=45, active=true, tags=manager, score=7.8, status=active, department=marketing, salary=85000}
5. [1] TestPerson{status=inactive, level=3, name=Eve, score=8, department=sales, age=30, salary=55000, active=false, tags=employee}
6. [1] TestPerson{status=active, department=qa, tags=test, level=1, name=Frank, age=0, score=0, salary=-5000, active=true}
7. [1] TestPerson{status=active, department=management, active=true, tags=executive, level=9, salary=95000, age=65, score=10, name=Grace}
8. [1] TestPerson{name=Henry, age=18, score=5.5, tags=junior, active=false, status=inactive, department=support, salary=25000, level=1}
9. [1] TestPerson{age=40, score=8.7, name=Ivy, active=true, department=engineering, salary=68000, tags=senior, status=active, level=6}
10. [1] TestPerson{name=X, salary=28000, status=active, department=intern, age=22, score=6.5, level=1, active=true, tags=temp}

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
