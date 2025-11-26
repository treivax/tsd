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

1. [1] TestPerson{salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, age=25, active=true, department=sales}
2. [1] TestPerson{age=35, tags=senior, name=Bob, active=true, score=9.2, department=engineering, salary=75000, level=5, status=active}
3. [1] TestPerson{level=1, name=Charlie, salary=0, active=false, tags=intern, age=16, department=hr, score=6, status=inactive}
4. [1] TestPerson{status=active, department=marketing, level=7, age=45, salary=85000, name=Diana, tags=manager, active=true, score=7.8}
5. [1] TestPerson{age=30, score=8, department=sales, level=3, salary=55000, tags=employee, status=inactive, name=Eve, active=false}
6. [1] TestPerson{department=qa, age=0, active=true, tags=test, status=active, salary=-5000, level=1, name=Frank, score=0}
7. [1] TestPerson{level=9, name=Grace, tags=executive, status=active, department=management, age=65, active=true, salary=95000, score=10}
8. [1] TestPerson{salary=25000, score=5.5, department=support, name=Henry, active=false, tags=junior, level=1, age=18, status=inactive}
9. [1] TestPerson{tags=senior, level=6, age=40, score=8.7, status=active, salary=68000, active=true, name=Ivy, department=engineering}
10. [1] TestPerson{department=intern, level=1, tags=temp, name=X, salary=28000, score=6.5, status=active, age=22, active=true}

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

1. [1] TestOrder{product_id=PROD001, priority=normal, discount=50, status=pending, customer_id=P001, amount=2, date=2024-01-15, total=1999.98, region=north}
2. [1] TestOrder{discount=0, customer_id=P002, region=south, date=2024-01-20, status=confirmed, product_id=PROD002, amount=1, priority=low, total=25.5}
3. [1] TestOrder{status=shipped, priority=high, total=225, date=2024-02-01, region=north, product_id=PROD003, amount=3, discount=15, customer_id=P001}
4. [1] TestOrder{discount=0, total=299.99, date=2024-02-05, amount=1, status=delivered, priority=normal, region=east, product_id=PROD004, customer_id=P004}
5. [1] TestOrder{discount=100, region=south, total=999.99, status=confirmed, product_id=PROD001, amount=1, date=2024-02-10, priority=high, customer_id=P002}
6. [1] TestOrder{amount=2, total=999.98, date=2024-02-15, discount=0, status=cancelled, customer_id=P005, priority=low, region=west, product_id=PROD005}
7. [1] TestOrder{total=600, priority=urgent, region=north, status=shipped, discount=50, product_id=PROD006, amount=4, date=2024-03-01, customer_id=P007}
8. [1] TestOrder{discount=0, total=255, customer_id=P010, amount=10, date=2024-03-05, priority=normal, region=south, product_id=PROD002, status=pending}
9. [1] TestOrder{region=north, product_id=PROD007, amount=1, discount=10, customer_id=P001, total=89.99, date=2024-03-10, status=completed, priority=low}
10. [1] TestOrder{discount=0, amount=1, total=75000, priority=urgent, customer_id=P006, region=east, product_id=PROD001, status=refunded, date=2024-03-15}

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

1. [1] TestPerson{age=25, active=true, department=sales, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice}
2. [1] TestPerson{department=engineering, salary=75000, level=5, status=active, age=35, tags=senior, name=Bob, active=true, score=9.2}
3. [1] TestPerson{level=1, name=Charlie, salary=0, active=false, tags=intern, age=16, department=hr, score=6, status=inactive}
4. [1] TestPerson{tags=manager, active=true, score=7.8, status=active, department=marketing, level=7, age=45, salary=85000, name=Diana}
5. [1] TestPerson{name=Eve, active=false, age=30, score=8, department=sales, level=3, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{department=qa, age=0, active=true, tags=test, status=active, salary=-5000, level=1, name=Frank, score=0}
7. [1] TestPerson{tags=executive, status=active, department=management, age=65, active=true, salary=95000, score=10, level=9, name=Grace}
8. [1] TestPerson{level=1, age=18, status=inactive, salary=25000, score=5.5, department=support, name=Henry, active=false, tags=junior}
9. [1] TestPerson{tags=senior, level=6, age=40, score=8.7, status=active, salary=68000, active=true, name=Ivy, department=engineering}
10. [1] TestPerson{age=22, active=true, department=intern, level=1, tags=temp, name=X, salary=28000, score=6.5, status=active}

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

1. [1] TestOrder{customer_id=P001, amount=2, date=2024-01-15, total=1999.98, region=north, product_id=PROD001, priority=normal, discount=50, status=pending}
2. [1] TestOrder{total=25.5, discount=0, customer_id=P002, region=south, date=2024-01-20, status=confirmed, product_id=PROD002, amount=1, priority=low}
3. [1] TestOrder{region=north, product_id=PROD003, amount=3, discount=15, customer_id=P001, status=shipped, priority=high, total=225, date=2024-02-01}
4. [1] TestOrder{customer_id=P004, discount=0, total=299.99, date=2024-02-05, amount=1, status=delivered, priority=normal, region=east, product_id=PROD004}
5. [1] TestOrder{customer_id=P002, discount=100, region=south, total=999.99, status=confirmed, product_id=PROD001, amount=1, date=2024-02-10, priority=high}
6. [1] TestOrder{customer_id=P005, priority=low, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, discount=0, status=cancelled}
7. [1] TestOrder{amount=4, date=2024-03-01, customer_id=P007, total=600, priority=urgent, region=north, status=shipped, discount=50, product_id=PROD006}
8. [1] TestOrder{product_id=PROD002, status=pending, discount=0, total=255, customer_id=P010, amount=10, date=2024-03-05, priority=normal, region=south}
9. [1] TestOrder{customer_id=P001, total=89.99, date=2024-03-10, status=completed, priority=low, region=north, product_id=PROD007, amount=1, discount=10}
10. [1] TestOrder{status=refunded, date=2024-03-15, discount=0, amount=1, total=75000, priority=urgent, customer_id=P006, region=east, product_id=PROD001}

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

1. [1] TestProduct{stock=50, supplier=TechSupply, price=999.99, available=true, category=electronics, brand=TechCorp, name=Laptop, rating=4.5, keywords=computer}
2. [1] TestProduct{supplier=TechSupply, category=accessories, price=25.5, available=true, keywords=peripheral, brand=TechCorp, stock=200, name=Mouse, rating=4}
3. [1] TestProduct{supplier=KeySupply, name=Keyboard, category=accessories, price=75, stock=0, available=false, rating=3.5, keywords=typing, brand=KeyTech}
4. [1] TestProduct{category=electronics, rating=4.8, supplier=ScreenSupply, keywords=display, available=true, brand=ScreenPro, stock=30, name=Monitor, price=299.99}
5. [1] TestProduct{supplier=OldSupply, available=false, rating=2, price=8.5, stock=0, name=OldKeyboard, keywords=obsolete, category=accessories, brand=OldTech}
6. [1] TestProduct{category=audio, brand=AudioMax, stock=75, price=150, available=true, supplier=AudioSupply, rating=4.6, name=Headphones, keywords=sound}
7. [1] TestProduct{price=89.99, available=true, rating=3.8, name=Webcam, category=electronics, brand=CamTech, stock=25, keywords=video, supplier=CamSupply}

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

1. [1] TestPerson{level=2, name=Alice, age=25, active=true, department=sales, salary=45000, score=8.5, tags=junior, status=active}
2. [1] TestPerson{department=engineering, salary=75000, level=5, status=active, age=35, tags=senior, name=Bob, active=true, score=9.2}
3. [1] TestPerson{name=Charlie, salary=0, active=false, tags=intern, age=16, department=hr, score=6, status=inactive, level=1}
4. [1] TestPerson{active=true, score=7.8, status=active, department=marketing, level=7, age=45, salary=85000, name=Diana, tags=manager}
5. [1] TestPerson{status=inactive, name=Eve, active=false, age=30, score=8, department=sales, level=3, salary=55000, tags=employee}
6. [1] TestPerson{score=0, department=qa, age=0, active=true, tags=test, status=active, salary=-5000, level=1, name=Frank}
7. [1] TestPerson{department=management, age=65, active=true, salary=95000, score=10, level=9, name=Grace, tags=executive, status=active}
8. [1] TestPerson{active=false, tags=junior, level=1, age=18, status=inactive, salary=25000, score=5.5, department=support, name=Henry}
9. [1] TestPerson{age=40, score=8.7, status=active, salary=68000, active=true, name=Ivy, department=engineering, tags=senior, level=6}
10. [1] TestPerson{department=intern, level=1, tags=temp, name=X, salary=28000, score=6.5, status=active, age=22, active=true}

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

1. [1] TestOrder{customer_id=P001, amount=2, date=2024-01-15, total=1999.98, region=north, product_id=PROD001, priority=normal, discount=50, status=pending}
2. [1] TestOrder{date=2024-01-20, status=confirmed, product_id=PROD002, amount=1, priority=low, total=25.5, discount=0, customer_id=P002, region=south}
3. [1] TestOrder{total=225, date=2024-02-01, region=north, product_id=PROD003, amount=3, discount=15, customer_id=P001, status=shipped, priority=high}
4. [1] TestOrder{total=299.99, date=2024-02-05, amount=1, status=delivered, priority=normal, region=east, product_id=PROD004, customer_id=P004, discount=0}
5. [1] TestOrder{priority=high, customer_id=P002, discount=100, region=south, total=999.99, status=confirmed, product_id=PROD001, amount=1, date=2024-02-10}
6. [1] TestOrder{amount=2, total=999.98, date=2024-02-15, discount=0, status=cancelled, customer_id=P005, priority=low, region=west, product_id=PROD005}
7. [1] TestOrder{product_id=PROD006, amount=4, date=2024-03-01, customer_id=P007, total=600, priority=urgent, region=north, status=shipped, discount=50}
8. [1] TestOrder{total=255, customer_id=P010, amount=10, date=2024-03-05, priority=normal, region=south, product_id=PROD002, status=pending, discount=0}
9. [1] TestOrder{product_id=PROD007, amount=1, discount=10, customer_id=P001, total=89.99, date=2024-03-10, status=completed, priority=low, region=north}
10. [1] TestOrder{date=2024-03-15, discount=0, amount=1, total=75000, priority=urgent, customer_id=P006, region=east, product_id=PROD001, status=refunded}

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

1. [1] TestPerson{active=true, department=sales, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, age=25}
2. [1] TestPerson{tags=senior, name=Bob, active=true, score=9.2, department=engineering, salary=75000, level=5, status=active, age=35}
3. [1] TestPerson{department=hr, score=6, status=inactive, level=1, name=Charlie, salary=0, active=false, tags=intern, age=16}
4. [1] TestPerson{department=marketing, level=7, age=45, salary=85000, name=Diana, tags=manager, active=true, score=7.8, status=active}
5. [1] TestPerson{name=Eve, active=false, age=30, score=8, department=sales, level=3, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{department=qa, age=0, active=true, tags=test, status=active, salary=-5000, level=1, name=Frank, score=0}
7. [1] TestPerson{active=true, salary=95000, score=10, level=9, name=Grace, tags=executive, status=active, department=management, age=65}
8. [1] TestPerson{department=support, name=Henry, active=false, tags=junior, level=1, age=18, status=inactive, salary=25000, score=5.5}
9. [1] TestPerson{age=40, score=8.7, status=active, salary=68000, active=true, name=Ivy, department=engineering, tags=senior, level=6}
10. [1] TestPerson{score=6.5, status=active, age=22, active=true, department=intern, level=1, tags=temp, name=X, salary=28000}

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

1. [1] TestOrder{date=2024-01-15, total=1999.98, region=north, product_id=PROD001, priority=normal, discount=50, status=pending, customer_id=P001, amount=2}
2. [1] TestOrder{discount=0, customer_id=P002, region=south, date=2024-01-20, status=confirmed, product_id=PROD002, amount=1, priority=low, total=25.5}
3. [1] TestOrder{priority=high, total=225, date=2024-02-01, region=north, product_id=PROD003, amount=3, discount=15, customer_id=P001, status=shipped}
4. [1] TestOrder{status=delivered, priority=normal, region=east, product_id=PROD004, customer_id=P004, discount=0, total=299.99, date=2024-02-05, amount=1}
5. [1] TestOrder{status=confirmed, product_id=PROD001, amount=1, date=2024-02-10, priority=high, customer_id=P002, discount=100, region=south, total=999.99}
6. [1] TestOrder{amount=2, total=999.98, date=2024-02-15, discount=0, status=cancelled, customer_id=P005, priority=low, region=west, product_id=PROD005}
7. [1] TestOrder{customer_id=P007, total=600, priority=urgent, region=north, status=shipped, discount=50, product_id=PROD006, amount=4, date=2024-03-01}
8. [1] TestOrder{product_id=PROD002, status=pending, discount=0, total=255, customer_id=P010, amount=10, date=2024-03-05, priority=normal, region=south}
9. [1] TestOrder{product_id=PROD007, amount=1, discount=10, customer_id=P001, total=89.99, date=2024-03-10, status=completed, priority=low, region=north}
10. [1] TestOrder{date=2024-03-15, discount=0, amount=1, total=75000, priority=urgent, customer_id=P006, region=east, product_id=PROD001, status=refunded}

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

1. [1] TestPerson{name=Alice, age=25, active=true, department=sales, salary=45000, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{age=35, tags=senior, name=Bob, active=true, score=9.2, department=engineering, salary=75000, level=5, status=active}
3. [1] TestPerson{score=6, status=inactive, level=1, name=Charlie, salary=0, active=false, tags=intern, age=16, department=hr}
4. [1] TestPerson{score=7.8, status=active, department=marketing, level=7, age=45, salary=85000, name=Diana, tags=manager, active=true}
5. [1] TestPerson{level=3, salary=55000, tags=employee, status=inactive, name=Eve, active=false, age=30, score=8, department=sales}
6. [1] TestPerson{department=qa, age=0, active=true, tags=test, status=active, salary=-5000, level=1, name=Frank, score=0}
7. [1] TestPerson{tags=executive, status=active, department=management, age=65, active=true, salary=95000, score=10, level=9, name=Grace}
8. [1] TestPerson{level=1, age=18, status=inactive, salary=25000, score=5.5, department=support, name=Henry, active=false, tags=junior}
9. [1] TestPerson{level=6, age=40, score=8.7, status=active, salary=68000, active=true, name=Ivy, department=engineering, tags=senior}
10. [1] TestPerson{salary=28000, score=6.5, status=active, age=22, active=true, department=intern, level=1, tags=temp, name=X}

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

1. [1] TestProduct{brand=TechCorp, name=Laptop, rating=4.5, keywords=computer, stock=50, supplier=TechSupply, price=999.99, available=true, category=electronics}
2. [1] TestProduct{keywords=peripheral, brand=TechCorp, stock=200, name=Mouse, rating=4, supplier=TechSupply, category=accessories, price=25.5, available=true}
3. [1] TestProduct{available=false, rating=3.5, keywords=typing, brand=KeyTech, supplier=KeySupply, name=Keyboard, category=accessories, price=75, stock=0}
4. [1] TestProduct{keywords=display, available=true, brand=ScreenPro, stock=30, name=Monitor, price=299.99, category=electronics, rating=4.8, supplier=ScreenSupply}
5. [1] TestProduct{rating=2, price=8.5, stock=0, name=OldKeyboard, keywords=obsolete, category=accessories, brand=OldTech, supplier=OldSupply, available=false}
6. [1] TestProduct{stock=75, price=150, available=true, supplier=AudioSupply, rating=4.6, name=Headphones, keywords=sound, category=audio, brand=AudioMax}
7. [1] TestProduct{price=89.99, available=true, rating=3.8, name=Webcam, category=electronics, brand=CamTech, stock=25, keywords=video, supplier=CamSupply}

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

1. [1] TestPerson{score=8.5, tags=junior, status=active, level=2, name=Alice, age=25, active=true, department=sales, salary=45000}
2. [1] TestPerson{level=5, status=active, age=35, tags=senior, name=Bob, active=true, score=9.2, department=engineering, salary=75000}
3. [1] TestPerson{salary=0, active=false, tags=intern, age=16, department=hr, score=6, status=inactive, level=1, name=Charlie}
4. [1] TestPerson{salary=85000, name=Diana, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7, age=45}
5. [1] TestPerson{name=Eve, active=false, age=30, score=8, department=sales, level=3, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{department=qa, age=0, active=true, tags=test, status=active, salary=-5000, level=1, name=Frank, score=0}
7. [1] TestPerson{active=true, salary=95000, score=10, level=9, name=Grace, tags=executive, status=active, department=management, age=65}
8. [1] TestPerson{score=5.5, department=support, name=Henry, active=false, tags=junior, level=1, age=18, status=inactive, salary=25000}
9. [1] TestPerson{tags=senior, level=6, age=40, score=8.7, status=active, salary=68000, active=true, name=Ivy, department=engineering}
10. [1] TestPerson{department=intern, level=1, tags=temp, name=X, salary=28000, score=6.5, status=active, age=22, active=true}

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

1. [1] TestOrder{region=north, product_id=PROD001, priority=normal, discount=50, status=pending, customer_id=P001, amount=2, date=2024-01-15, total=1999.98}
2. [1] TestOrder{product_id=PROD002, amount=1, priority=low, total=25.5, discount=0, customer_id=P002, region=south, date=2024-01-20, status=confirmed}
3. [1] TestOrder{total=225, date=2024-02-01, region=north, product_id=PROD003, amount=3, discount=15, customer_id=P001, status=shipped, priority=high}
4. [1] TestOrder{customer_id=P004, discount=0, total=299.99, date=2024-02-05, amount=1, status=delivered, priority=normal, region=east, product_id=PROD004}
5. [1] TestOrder{priority=high, customer_id=P002, discount=100, region=south, total=999.99, status=confirmed, product_id=PROD001, amount=1, date=2024-02-10}
6. [1] TestOrder{status=cancelled, customer_id=P005, priority=low, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, discount=0}
7. [1] TestOrder{amount=4, date=2024-03-01, customer_id=P007, total=600, priority=urgent, region=north, status=shipped, discount=50, product_id=PROD006}
8. [1] TestOrder{region=south, product_id=PROD002, status=pending, discount=0, total=255, customer_id=P010, amount=10, date=2024-03-05, priority=normal}
9. [1] TestOrder{date=2024-03-10, status=completed, priority=low, region=north, product_id=PROD007, amount=1, discount=10, customer_id=P001, total=89.99}
10. [1] TestOrder{date=2024-03-15, discount=0, amount=1, total=75000, priority=urgent, customer_id=P006, region=east, product_id=PROD001, status=refunded}

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

1. [1] TestPerson{department=sales, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, age=25, active=true}
2. [1] TestPerson{age=35, tags=senior, name=Bob, active=true, score=9.2, department=engineering, salary=75000, level=5, status=active}
3. [1] TestPerson{age=16, department=hr, score=6, status=inactive, level=1, name=Charlie, salary=0, active=false, tags=intern}
4. [1] TestPerson{status=active, department=marketing, level=7, age=45, salary=85000, name=Diana, tags=manager, active=true, score=7.8}
5. [1] TestPerson{level=3, salary=55000, tags=employee, status=inactive, name=Eve, active=false, age=30, score=8, department=sales}
6. [1] TestPerson{score=0, department=qa, age=0, active=true, tags=test, status=active, salary=-5000, level=1, name=Frank}
7. [1] TestPerson{status=active, department=management, age=65, active=true, salary=95000, score=10, level=9, name=Grace, tags=executive}
8. [1] TestPerson{salary=25000, score=5.5, department=support, name=Henry, active=false, tags=junior, level=1, age=18, status=inactive}
9. [1] TestPerson{age=40, score=8.7, status=active, salary=68000, active=true, name=Ivy, department=engineering, tags=senior, level=6}
10. [1] TestPerson{score=6.5, status=active, age=22, active=true, department=intern, level=1, tags=temp, name=X, salary=28000}
11. [1] TestOrder{region=north, product_id=PROD001, priority=normal, discount=50, status=pending, customer_id=P001, amount=2, date=2024-01-15, total=1999.98}
12. [1] TestOrder{customer_id=P002, region=south, date=2024-01-20, status=confirmed, product_id=PROD002, amount=1, priority=low, total=25.5, discount=0}
13. [1] TestOrder{priority=high, total=225, date=2024-02-01, region=north, product_id=PROD003, amount=3, discount=15, customer_id=P001, status=shipped}
14. [1] TestOrder{date=2024-02-05, amount=1, status=delivered, priority=normal, region=east, product_id=PROD004, customer_id=P004, discount=0, total=299.99}
15. [1] TestOrder{discount=100, region=south, total=999.99, status=confirmed, product_id=PROD001, amount=1, date=2024-02-10, priority=high, customer_id=P002}
16. [1] TestOrder{total=999.98, date=2024-02-15, discount=0, status=cancelled, customer_id=P005, priority=low, region=west, product_id=PROD005, amount=2}
17. [1] TestOrder{customer_id=P007, total=600, priority=urgent, region=north, status=shipped, discount=50, product_id=PROD006, amount=4, date=2024-03-01}
18. [1] TestOrder{total=255, customer_id=P010, amount=10, date=2024-03-05, priority=normal, region=south, product_id=PROD002, status=pending, discount=0}
19. [1] TestOrder{region=north, product_id=PROD007, amount=1, discount=10, customer_id=P001, total=89.99, date=2024-03-10, status=completed, priority=low}
20. [1] TestOrder{status=refunded, date=2024-03-15, discount=0, amount=1, total=75000, priority=urgent, customer_id=P006, region=east, product_id=PROD001}

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

1. [1] TestPerson{name=Alice, age=25, active=true, department=sales, salary=45000, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{score=9.2, department=engineering, salary=75000, level=5, status=active, age=35, tags=senior, name=Bob, active=true}
3. [1] TestPerson{level=1, name=Charlie, salary=0, active=false, tags=intern, age=16, department=hr, score=6, status=inactive}
4. [1] TestPerson{status=active, department=marketing, level=7, age=45, salary=85000, name=Diana, tags=manager, active=true, score=7.8}
5. [1] TestPerson{level=3, salary=55000, tags=employee, status=inactive, name=Eve, active=false, age=30, score=8, department=sales}
6. [1] TestPerson{level=1, name=Frank, score=0, department=qa, age=0, active=true, tags=test, status=active, salary=-5000}
7. [1] TestPerson{tags=executive, status=active, department=management, age=65, active=true, salary=95000, score=10, level=9, name=Grace}
8. [1] TestPerson{department=support, name=Henry, active=false, tags=junior, level=1, age=18, status=inactive, salary=25000, score=5.5}
9. [1] TestPerson{salary=68000, active=true, name=Ivy, department=engineering, tags=senior, level=6, age=40, score=8.7, status=active}
10. [1] TestPerson{department=intern, level=1, tags=temp, name=X, salary=28000, score=6.5, status=active, age=22, active=true}

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

1. [1] TestPerson{salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, age=25, active=true, department=sales}
2. [1] TestPerson{salary=75000, level=5, status=active, age=35, tags=senior, name=Bob, active=true, score=9.2, department=engineering}
3. [1] TestPerson{department=hr, score=6, status=inactive, level=1, name=Charlie, salary=0, active=false, tags=intern, age=16}
4. [1] TestPerson{name=Diana, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7, age=45, salary=85000}
5. [1] TestPerson{name=Eve, active=false, age=30, score=8, department=sales, level=3, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{salary=-5000, level=1, name=Frank, score=0, department=qa, age=0, active=true, tags=test, status=active}
7. [1] TestPerson{status=active, department=management, age=65, active=true, salary=95000, score=10, level=9, name=Grace, tags=executive}
8. [1] TestPerson{salary=25000, score=5.5, department=support, name=Henry, active=false, tags=junior, level=1, age=18, status=inactive}
9. [1] TestPerson{tags=senior, level=6, age=40, score=8.7, status=active, salary=68000, active=true, name=Ivy, department=engineering}
10. [1] TestPerson{level=1, tags=temp, name=X, salary=28000, score=6.5, status=active, age=22, active=true, department=intern}
11. [1] TestOrder{region=north, product_id=PROD001, priority=normal, discount=50, status=pending, customer_id=P001, amount=2, date=2024-01-15, total=1999.98}
12. [1] TestOrder{priority=low, total=25.5, discount=0, customer_id=P002, region=south, date=2024-01-20, status=confirmed, product_id=PROD002, amount=1}
13. [1] TestOrder{total=225, date=2024-02-01, region=north, product_id=PROD003, amount=3, discount=15, customer_id=P001, status=shipped, priority=high}
14. [1] TestOrder{discount=0, total=299.99, date=2024-02-05, amount=1, status=delivered, priority=normal, region=east, product_id=PROD004, customer_id=P004}
15. [1] TestOrder{status=confirmed, product_id=PROD001, amount=1, date=2024-02-10, priority=high, customer_id=P002, discount=100, region=south, total=999.99}
16. [1] TestOrder{date=2024-02-15, discount=0, status=cancelled, customer_id=P005, priority=low, region=west, product_id=PROD005, amount=2, total=999.98}
17. [1] TestOrder{priority=urgent, region=north, status=shipped, discount=50, product_id=PROD006, amount=4, date=2024-03-01, customer_id=P007, total=600}
18. [1] TestOrder{total=255, customer_id=P010, amount=10, date=2024-03-05, priority=normal, region=south, product_id=PROD002, status=pending, discount=0}
19. [1] TestOrder{status=completed, priority=low, region=north, product_id=PROD007, amount=1, discount=10, customer_id=P001, total=89.99, date=2024-03-10}
20. [1] TestOrder{product_id=PROD001, status=refunded, date=2024-03-15, discount=0, amount=1, total=75000, priority=urgent, customer_id=P006, region=east}

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

1. [1] TestPerson{score=8.5, tags=junior, status=active, level=2, name=Alice, age=25, active=true, department=sales, salary=45000}
2. [1] TestPerson{age=35, tags=senior, name=Bob, active=true, score=9.2, department=engineering, salary=75000, level=5, status=active}
3. [1] TestPerson{department=hr, score=6, status=inactive, level=1, name=Charlie, salary=0, active=false, tags=intern, age=16}
4. [1] TestPerson{department=marketing, level=7, age=45, salary=85000, name=Diana, tags=manager, active=true, score=7.8, status=active}
5. [1] TestPerson{level=3, salary=55000, tags=employee, status=inactive, name=Eve, active=false, age=30, score=8, department=sales}
6. [1] TestPerson{department=qa, age=0, active=true, tags=test, status=active, salary=-5000, level=1, name=Frank, score=0}
7. [1] TestPerson{tags=executive, status=active, department=management, age=65, active=true, salary=95000, score=10, level=9, name=Grace}
8. [1] TestPerson{salary=25000, score=5.5, department=support, name=Henry, active=false, tags=junior, level=1, age=18, status=inactive}
9. [1] TestPerson{tags=senior, level=6, age=40, score=8.7, status=active, salary=68000, active=true, name=Ivy, department=engineering}
10. [1] TestPerson{score=6.5, status=active, age=22, active=true, department=intern, level=1, tags=temp, name=X, salary=28000}

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

1. [1] TestOrder{customer_id=P001, amount=2, date=2024-01-15, total=1999.98, region=north, product_id=PROD001, priority=normal, discount=50, status=pending}
2. [1] TestOrder{date=2024-01-20, status=confirmed, product_id=PROD002, amount=1, priority=low, total=25.5, discount=0, customer_id=P002, region=south}
3. [1] TestOrder{amount=3, discount=15, customer_id=P001, status=shipped, priority=high, total=225, date=2024-02-01, region=north, product_id=PROD003}
4. [1] TestOrder{date=2024-02-05, amount=1, status=delivered, priority=normal, region=east, product_id=PROD004, customer_id=P004, discount=0, total=299.99}
5. [1] TestOrder{priority=high, customer_id=P002, discount=100, region=south, total=999.99, status=confirmed, product_id=PROD001, amount=1, date=2024-02-10}
6. [1] TestOrder{amount=2, total=999.98, date=2024-02-15, discount=0, status=cancelled, customer_id=P005, priority=low, region=west, product_id=PROD005}
7. [1] TestOrder{customer_id=P007, total=600, priority=urgent, region=north, status=shipped, discount=50, product_id=PROD006, amount=4, date=2024-03-01}
8. [1] TestOrder{date=2024-03-05, priority=normal, region=south, product_id=PROD002, status=pending, discount=0, total=255, customer_id=P010, amount=10}
9. [1] TestOrder{product_id=PROD007, amount=1, discount=10, customer_id=P001, total=89.99, date=2024-03-10, status=completed, priority=low, region=north}
10. [1] TestOrder{amount=1, total=75000, priority=urgent, customer_id=P006, region=east, product_id=PROD001, status=refunded, date=2024-03-15, discount=0}

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

1. [1] TestPerson{score=8.5, tags=junior, status=active, level=2, name=Alice, age=25, active=true, department=sales, salary=45000}
2. [1] TestPerson{department=engineering, salary=75000, level=5, status=active, age=35, tags=senior, name=Bob, active=true, score=9.2}
3. [1] TestPerson{active=false, tags=intern, age=16, department=hr, score=6, status=inactive, level=1, name=Charlie, salary=0}
4. [1] TestPerson{status=active, department=marketing, level=7, age=45, salary=85000, name=Diana, tags=manager, active=true, score=7.8}
5. [1] TestPerson{name=Eve, active=false, age=30, score=8, department=sales, level=3, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{tags=test, status=active, salary=-5000, level=1, name=Frank, score=0, department=qa, age=0, active=true}
7. [1] TestPerson{salary=95000, score=10, level=9, name=Grace, tags=executive, status=active, department=management, age=65, active=true}
8. [1] TestPerson{tags=junior, level=1, age=18, status=inactive, salary=25000, score=5.5, department=support, name=Henry, active=false}
9. [1] TestPerson{salary=68000, active=true, name=Ivy, department=engineering, tags=senior, level=6, age=40, score=8.7, status=active}
10. [1] TestPerson{score=6.5, status=active, age=22, active=true, department=intern, level=1, tags=temp, name=X, salary=28000}

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
