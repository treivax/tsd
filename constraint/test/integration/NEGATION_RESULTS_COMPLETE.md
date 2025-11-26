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

1. [1] TestPerson{status=active, name=Alice, score=8.5, level=2, age=25, active=true, department=sales, tags=junior, salary=45000}
2. [1] TestPerson{active=true, age=35, score=9.2, status=active, level=5, name=Bob, salary=75000, department=engineering, tags=senior}
3. [1] TestPerson{active=false, level=1, age=16, department=hr, tags=intern, name=Charlie, score=6, status=inactive, salary=0}
4. [1] TestPerson{department=marketing, age=45, salary=85000, score=7.8, level=7, name=Diana, tags=manager, active=true, status=active}
5. [1] TestPerson{status=inactive, department=sales, tags=employee, level=3, name=Eve, salary=55000, score=8, age=30, active=false}
6. [1] TestPerson{active=true, level=1, age=0, salary=-5000, tags=test, status=active, score=0, department=qa, name=Frank}
7. [1] TestPerson{tags=executive, status=active, department=management, level=9, age=65, active=true, name=Grace, score=10, salary=95000}
8. [1] TestPerson{name=Henry, age=18, department=support, tags=junior, salary=25000, active=false, score=5.5, level=1, status=inactive}
9. [1] TestPerson{active=true, tags=senior, department=engineering, level=6, salary=68000, status=active, age=40, name=Ivy, score=8.7}
10. [1] TestPerson{age=22, active=true, department=intern, tags=temp, salary=28000, score=6.5, status=active, level=1, name=X}

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

1. [1] TestOrder{product_id=PROD001, status=pending, discount=50, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, priority=normal, region=north}
2. [1] TestOrder{status=confirmed, priority=low, region=south, amount=1, total=25.5, date=2024-01-20, customer_id=P002, product_id=PROD002, discount=0}
3. [1] TestOrder{status=shipped, region=north, amount=3, total=225, date=2024-02-01, priority=high, discount=15, product_id=PROD003, customer_id=P001}
4. [1] TestOrder{region=east, customer_id=P004, amount=1, status=delivered, discount=0, total=299.99, date=2024-02-05, priority=normal, product_id=PROD004}
5. [1] TestOrder{status=confirmed, product_id=PROD001, priority=high, discount=100, region=south, customer_id=P002, total=999.99, date=2024-02-10, amount=1}
6. [1] TestOrder{region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, status=cancelled, priority=low, amount=2, discount=0}
7. [1] TestOrder{amount=4, customer_id=P007, total=600, priority=urgent, discount=50, region=north, date=2024-03-01, product_id=PROD006, status=shipped}
8. [1] TestOrder{customer_id=P010, amount=10, total=255, date=2024-03-05, discount=0, status=pending, product_id=PROD002, priority=normal, region=south}
9. [1] TestOrder{amount=1, region=north, product_id=PROD007, priority=low, discount=10, total=89.99, date=2024-03-10, status=completed, customer_id=P001}
10. [1] TestOrder{status=refunded, discount=0, amount=1, date=2024-03-15, region=east, priority=urgent, customer_id=P006, product_id=PROD001, total=75000}

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

1. [1] TestPerson{status=active, name=Alice, score=8.5, level=2, age=25, active=true, department=sales, tags=junior, salary=45000}
2. [1] TestPerson{active=true, age=35, score=9.2, status=active, level=5, name=Bob, salary=75000, department=engineering, tags=senior}
3. [1] TestPerson{tags=intern, name=Charlie, score=6, status=inactive, salary=0, active=false, level=1, age=16, department=hr}
4. [1] TestPerson{score=7.8, level=7, name=Diana, tags=manager, active=true, status=active, department=marketing, age=45, salary=85000}
5. [1] TestPerson{score=8, age=30, active=false, status=inactive, department=sales, tags=employee, level=3, name=Eve, salary=55000}
6. [1] TestPerson{name=Frank, active=true, level=1, age=0, salary=-5000, tags=test, status=active, score=0, department=qa}
7. [1] TestPerson{salary=95000, tags=executive, status=active, department=management, level=9, age=65, active=true, name=Grace, score=10}
8. [1] TestPerson{salary=25000, active=false, score=5.5, level=1, status=inactive, name=Henry, age=18, department=support, tags=junior}
9. [1] TestPerson{active=true, tags=senior, department=engineering, level=6, salary=68000, status=active, age=40, name=Ivy, score=8.7}
10. [1] TestPerson{status=active, level=1, name=X, age=22, active=true, department=intern, tags=temp, salary=28000, score=6.5}

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

1. [1] TestOrder{region=north, product_id=PROD001, status=pending, discount=50, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, priority=normal}
2. [1] TestOrder{customer_id=P002, product_id=PROD002, discount=0, status=confirmed, priority=low, region=south, amount=1, total=25.5, date=2024-01-20}
3. [1] TestOrder{discount=15, product_id=PROD003, customer_id=P001, status=shipped, region=north, amount=3, total=225, date=2024-02-01, priority=high}
4. [1] TestOrder{date=2024-02-05, priority=normal, product_id=PROD004, region=east, customer_id=P004, amount=1, status=delivered, discount=0, total=299.99}
5. [1] TestOrder{discount=100, region=south, customer_id=P002, total=999.99, date=2024-02-10, amount=1, status=confirmed, product_id=PROD001, priority=high}
6. [1] TestOrder{customer_id=P005, status=cancelled, priority=low, amount=2, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}
7. [1] TestOrder{total=600, priority=urgent, discount=50, region=north, date=2024-03-01, product_id=PROD006, status=shipped, amount=4, customer_id=P007}
8. [1] TestOrder{discount=0, status=pending, product_id=PROD002, priority=normal, region=south, customer_id=P010, amount=10, total=255, date=2024-03-05}
9. [1] TestOrder{status=completed, customer_id=P001, amount=1, region=north, product_id=PROD007, priority=low, discount=10, total=89.99, date=2024-03-10}
10. [1] TestOrder{status=refunded, discount=0, amount=1, date=2024-03-15, region=east, priority=urgent, customer_id=P006, product_id=PROD001, total=75000}

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

1. [1] TestProduct{stock=50, rating=4.5, name=Laptop, category=electronics, price=999.99, available=true, keywords=computer, supplier=TechSupply, brand=TechCorp}
2. [1] TestProduct{category=accessories, rating=4, supplier=TechSupply, name=Mouse, keywords=peripheral, stock=200, price=25.5, available=true, brand=TechCorp}
3. [1] TestProduct{available=false, rating=3.5, keywords=typing, stock=0, supplier=KeySupply, category=accessories, brand=KeyTech, name=Keyboard, price=75}
4. [1] TestProduct{name=Monitor, available=true, brand=ScreenPro, price=299.99, keywords=display, rating=4.8, stock=30, supplier=ScreenSupply, category=electronics}
5. [1] TestProduct{price=8.5, keywords=obsolete, stock=0, name=OldKeyboard, rating=2, category=accessories, available=false, brand=OldTech, supplier=OldSupply}
6. [1] TestProduct{stock=75, supplier=AudioSupply, name=Headphones, category=audio, price=150, rating=4.6, available=true, keywords=sound, brand=AudioMax}
7. [1] TestProduct{rating=3.8, brand=CamTech, stock=25, supplier=CamSupply, name=Webcam, available=true, keywords=video, category=electronics, price=89.99}

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

1. [1] TestPerson{status=active, name=Alice, score=8.5, level=2, age=25, active=true, department=sales, tags=junior, salary=45000}
2. [1] TestPerson{tags=senior, active=true, age=35, score=9.2, status=active, level=5, name=Bob, salary=75000, department=engineering}
3. [1] TestPerson{name=Charlie, score=6, status=inactive, salary=0, active=false, level=1, age=16, department=hr, tags=intern}
4. [1] TestPerson{active=true, status=active, department=marketing, age=45, salary=85000, score=7.8, level=7, name=Diana, tags=manager}
5. [1] TestPerson{name=Eve, salary=55000, score=8, age=30, active=false, status=inactive, department=sales, tags=employee, level=3}
6. [1] TestPerson{department=qa, name=Frank, active=true, level=1, age=0, salary=-5000, tags=test, status=active, score=0}
7. [1] TestPerson{salary=95000, tags=executive, status=active, department=management, level=9, age=65, active=true, name=Grace, score=10}
8. [1] TestPerson{department=support, tags=junior, salary=25000, active=false, score=5.5, level=1, status=inactive, name=Henry, age=18}
9. [1] TestPerson{active=true, tags=senior, department=engineering, level=6, salary=68000, status=active, age=40, name=Ivy, score=8.7}
10. [1] TestPerson{active=true, department=intern, tags=temp, salary=28000, score=6.5, status=active, level=1, name=X, age=22}

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

1. [1] TestOrder{status=pending, discount=50, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, priority=normal, region=north, product_id=PROD001}
2. [1] TestOrder{status=confirmed, priority=low, region=south, amount=1, total=25.5, date=2024-01-20, customer_id=P002, product_id=PROD002, discount=0}
3. [1] TestOrder{product_id=PROD003, customer_id=P001, status=shipped, region=north, amount=3, total=225, date=2024-02-01, priority=high, discount=15}
4. [1] TestOrder{region=east, customer_id=P004, amount=1, status=delivered, discount=0, total=299.99, date=2024-02-05, priority=normal, product_id=PROD004}
5. [1] TestOrder{region=south, customer_id=P002, total=999.99, date=2024-02-10, amount=1, status=confirmed, product_id=PROD001, priority=high, discount=100}
6. [1] TestOrder{date=2024-02-15, customer_id=P005, status=cancelled, priority=low, amount=2, discount=0, region=west, product_id=PROD005, total=999.98}
7. [1] TestOrder{status=shipped, amount=4, customer_id=P007, total=600, priority=urgent, discount=50, region=north, date=2024-03-01, product_id=PROD006}
8. [1] TestOrder{amount=10, total=255, date=2024-03-05, discount=0, status=pending, product_id=PROD002, priority=normal, region=south, customer_id=P010}
9. [1] TestOrder{customer_id=P001, amount=1, region=north, product_id=PROD007, priority=low, discount=10, total=89.99, date=2024-03-10, status=completed}
10. [1] TestOrder{product_id=PROD001, total=75000, status=refunded, discount=0, amount=1, date=2024-03-15, region=east, priority=urgent, customer_id=P006}

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

1. [1] TestPerson{status=active, name=Alice, score=8.5, level=2, age=25, active=true, department=sales, tags=junior, salary=45000}
2. [1] TestPerson{salary=75000, department=engineering, tags=senior, active=true, age=35, score=9.2, status=active, level=5, name=Bob}
3. [1] TestPerson{tags=intern, name=Charlie, score=6, status=inactive, salary=0, active=false, level=1, age=16, department=hr}
4. [1] TestPerson{department=marketing, age=45, salary=85000, score=7.8, level=7, name=Diana, tags=manager, active=true, status=active}
5. [1] TestPerson{status=inactive, department=sales, tags=employee, level=3, name=Eve, salary=55000, score=8, age=30, active=false}
6. [1] TestPerson{score=0, department=qa, name=Frank, active=true, level=1, age=0, salary=-5000, tags=test, status=active}
7. [1] TestPerson{active=true, name=Grace, score=10, salary=95000, tags=executive, status=active, department=management, level=9, age=65}
8. [1] TestPerson{name=Henry, age=18, department=support, tags=junior, salary=25000, active=false, score=5.5, level=1, status=inactive}
9. [1] TestPerson{active=true, tags=senior, department=engineering, level=6, salary=68000, status=active, age=40, name=Ivy, score=8.7}
10. [1] TestPerson{score=6.5, status=active, level=1, name=X, age=22, active=true, department=intern, tags=temp, salary=28000}

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

1. [1] TestOrder{product_id=PROD001, status=pending, discount=50, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, priority=normal, region=north}
2. [1] TestOrder{date=2024-01-20, customer_id=P002, product_id=PROD002, discount=0, status=confirmed, priority=low, region=south, amount=1, total=25.5}
3. [1] TestOrder{amount=3, total=225, date=2024-02-01, priority=high, discount=15, product_id=PROD003, customer_id=P001, status=shipped, region=north}
4. [1] TestOrder{status=delivered, discount=0, total=299.99, date=2024-02-05, priority=normal, product_id=PROD004, region=east, customer_id=P004, amount=1}
5. [1] TestOrder{status=confirmed, product_id=PROD001, priority=high, discount=100, region=south, customer_id=P002, total=999.99, date=2024-02-10, amount=1}
6. [1] TestOrder{discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, status=cancelled, priority=low, amount=2}
7. [1] TestOrder{total=600, priority=urgent, discount=50, region=north, date=2024-03-01, product_id=PROD006, status=shipped, amount=4, customer_id=P007}
8. [1] TestOrder{customer_id=P010, amount=10, total=255, date=2024-03-05, discount=0, status=pending, product_id=PROD002, priority=normal, region=south}
9. [1] TestOrder{amount=1, region=north, product_id=PROD007, priority=low, discount=10, total=89.99, date=2024-03-10, status=completed, customer_id=P001}
10. [1] TestOrder{status=refunded, discount=0, amount=1, date=2024-03-15, region=east, priority=urgent, customer_id=P006, product_id=PROD001, total=75000}

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

1. [1] TestPerson{department=sales, tags=junior, salary=45000, status=active, name=Alice, score=8.5, level=2, age=25, active=true}
2. [1] TestPerson{name=Bob, salary=75000, department=engineering, tags=senior, active=true, age=35, score=9.2, status=active, level=5}
3. [1] TestPerson{name=Charlie, score=6, status=inactive, salary=0, active=false, level=1, age=16, department=hr, tags=intern}
4. [1] TestPerson{name=Diana, tags=manager, active=true, status=active, department=marketing, age=45, salary=85000, score=7.8, level=7}
5. [1] TestPerson{status=inactive, department=sales, tags=employee, level=3, name=Eve, salary=55000, score=8, age=30, active=false}
6. [1] TestPerson{salary=-5000, tags=test, status=active, score=0, department=qa, name=Frank, active=true, level=1, age=0}
7. [1] TestPerson{age=65, active=true, name=Grace, score=10, salary=95000, tags=executive, status=active, department=management, level=9}
8. [1] TestPerson{salary=25000, active=false, score=5.5, level=1, status=inactive, name=Henry, age=18, department=support, tags=junior}
9. [1] TestPerson{name=Ivy, score=8.7, active=true, tags=senior, department=engineering, level=6, salary=68000, status=active, age=40}
10. [1] TestPerson{department=intern, tags=temp, salary=28000, score=6.5, status=active, level=1, name=X, age=22, active=true}

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

1. [1] TestProduct{supplier=TechSupply, brand=TechCorp, stock=50, rating=4.5, name=Laptop, category=electronics, price=999.99, available=true, keywords=computer}
2. [1] TestProduct{keywords=peripheral, stock=200, price=25.5, available=true, brand=TechCorp, category=accessories, rating=4, supplier=TechSupply, name=Mouse}
3. [1] TestProduct{available=false, rating=3.5, keywords=typing, stock=0, supplier=KeySupply, category=accessories, brand=KeyTech, name=Keyboard, price=75}
4. [1] TestProduct{name=Monitor, available=true, brand=ScreenPro, price=299.99, keywords=display, rating=4.8, stock=30, supplier=ScreenSupply, category=electronics}
5. [1] TestProduct{category=accessories, available=false, brand=OldTech, supplier=OldSupply, price=8.5, keywords=obsolete, stock=0, name=OldKeyboard, rating=2}
6. [1] TestProduct{price=150, rating=4.6, available=true, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply, name=Headphones, category=audio}
7. [1] TestProduct{price=89.99, rating=3.8, brand=CamTech, stock=25, supplier=CamSupply, name=Webcam, available=true, keywords=video, category=electronics}

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

1. [1] TestPerson{status=active, name=Alice, score=8.5, level=2, age=25, active=true, department=sales, tags=junior, salary=45000}
2. [1] TestPerson{age=35, score=9.2, status=active, level=5, name=Bob, salary=75000, department=engineering, tags=senior, active=true}
3. [1] TestPerson{salary=0, active=false, level=1, age=16, department=hr, tags=intern, name=Charlie, score=6, status=inactive}
4. [1] TestPerson{tags=manager, active=true, status=active, department=marketing, age=45, salary=85000, score=7.8, level=7, name=Diana}
5. [1] TestPerson{tags=employee, level=3, name=Eve, salary=55000, score=8, age=30, active=false, status=inactive, department=sales}
6. [1] TestPerson{age=0, salary=-5000, tags=test, status=active, score=0, department=qa, name=Frank, active=true, level=1}
7. [1] TestPerson{status=active, department=management, level=9, age=65, active=true, name=Grace, score=10, salary=95000, tags=executive}
8. [1] TestPerson{score=5.5, level=1, status=inactive, name=Henry, age=18, department=support, tags=junior, salary=25000, active=false}
9. [1] TestPerson{active=true, tags=senior, department=engineering, level=6, salary=68000, status=active, age=40, name=Ivy, score=8.7}
10. [1] TestPerson{level=1, name=X, age=22, active=true, department=intern, tags=temp, salary=28000, score=6.5, status=active}

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

1. [1] TestOrder{priority=normal, region=north, product_id=PROD001, status=pending, discount=50, customer_id=P001, amount=2, total=1999.98, date=2024-01-15}
2. [1] TestOrder{customer_id=P002, product_id=PROD002, discount=0, status=confirmed, priority=low, region=south, amount=1, total=25.5, date=2024-01-20}
3. [1] TestOrder{priority=high, discount=15, product_id=PROD003, customer_id=P001, status=shipped, region=north, amount=3, total=225, date=2024-02-01}
4. [1] TestOrder{total=299.99, date=2024-02-05, priority=normal, product_id=PROD004, region=east, customer_id=P004, amount=1, status=delivered, discount=0}
5. [1] TestOrder{date=2024-02-10, amount=1, status=confirmed, product_id=PROD001, priority=high, discount=100, region=south, customer_id=P002, total=999.99}
6. [1] TestOrder{date=2024-02-15, customer_id=P005, status=cancelled, priority=low, amount=2, discount=0, region=west, product_id=PROD005, total=999.98}
7. [1] TestOrder{amount=4, customer_id=P007, total=600, priority=urgent, discount=50, region=north, date=2024-03-01, product_id=PROD006, status=shipped}
8. [1] TestOrder{customer_id=P010, amount=10, total=255, date=2024-03-05, discount=0, status=pending, product_id=PROD002, priority=normal, region=south}
9. [1] TestOrder{amount=1, region=north, product_id=PROD007, priority=low, discount=10, total=89.99, date=2024-03-10, status=completed, customer_id=P001}
10. [1] TestOrder{priority=urgent, customer_id=P006, product_id=PROD001, total=75000, status=refunded, discount=0, amount=1, date=2024-03-15, region=east}

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

1. [1] TestPerson{level=2, age=25, active=true, department=sales, tags=junior, salary=45000, status=active, name=Alice, score=8.5}
2. [1] TestPerson{active=true, age=35, score=9.2, status=active, level=5, name=Bob, salary=75000, department=engineering, tags=senior}
3. [1] TestPerson{status=inactive, salary=0, active=false, level=1, age=16, department=hr, tags=intern, name=Charlie, score=6}
4. [1] TestPerson{level=7, name=Diana, tags=manager, active=true, status=active, department=marketing, age=45, salary=85000, score=7.8}
5. [1] TestPerson{status=inactive, department=sales, tags=employee, level=3, name=Eve, salary=55000, score=8, age=30, active=false}
6. [1] TestPerson{active=true, level=1, age=0, salary=-5000, tags=test, status=active, score=0, department=qa, name=Frank}
7. [1] TestPerson{salary=95000, tags=executive, status=active, department=management, level=9, age=65, active=true, name=Grace, score=10}
8. [1] TestPerson{salary=25000, active=false, score=5.5, level=1, status=inactive, name=Henry, age=18, department=support, tags=junior}
9. [1] TestPerson{active=true, tags=senior, department=engineering, level=6, salary=68000, status=active, age=40, name=Ivy, score=8.7}
10. [1] TestPerson{level=1, name=X, age=22, active=true, department=intern, tags=temp, salary=28000, score=6.5, status=active}
11. [1] TestOrder{region=north, product_id=PROD001, status=pending, discount=50, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, priority=normal}
12. [1] TestOrder{customer_id=P002, product_id=PROD002, discount=0, status=confirmed, priority=low, region=south, amount=1, total=25.5, date=2024-01-20}
13. [1] TestOrder{priority=high, discount=15, product_id=PROD003, customer_id=P001, status=shipped, region=north, amount=3, total=225, date=2024-02-01}
14. [1] TestOrder{total=299.99, date=2024-02-05, priority=normal, product_id=PROD004, region=east, customer_id=P004, amount=1, status=delivered, discount=0}
15. [1] TestOrder{date=2024-02-10, amount=1, status=confirmed, product_id=PROD001, priority=high, discount=100, region=south, customer_id=P002, total=999.99}
16. [1] TestOrder{date=2024-02-15, customer_id=P005, status=cancelled, priority=low, amount=2, discount=0, region=west, product_id=PROD005, total=999.98}
17. [1] TestOrder{amount=4, customer_id=P007, total=600, priority=urgent, discount=50, region=north, date=2024-03-01, product_id=PROD006, status=shipped}
18. [1] TestOrder{status=pending, product_id=PROD002, priority=normal, region=south, customer_id=P010, amount=10, total=255, date=2024-03-05, discount=0}
19. [1] TestOrder{discount=10, total=89.99, date=2024-03-10, status=completed, customer_id=P001, amount=1, region=north, product_id=PROD007, priority=low}
20. [1] TestOrder{amount=1, date=2024-03-15, region=east, priority=urgent, customer_id=P006, product_id=PROD001, total=75000, status=refunded, discount=0}

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

1. [1] TestPerson{status=active, name=Alice, score=8.5, level=2, age=25, active=true, department=sales, tags=junior, salary=45000}
2. [1] TestPerson{level=5, name=Bob, salary=75000, department=engineering, tags=senior, active=true, age=35, score=9.2, status=active}
3. [1] TestPerson{tags=intern, name=Charlie, score=6, status=inactive, salary=0, active=false, level=1, age=16, department=hr}
4. [1] TestPerson{department=marketing, age=45, salary=85000, score=7.8, level=7, name=Diana, tags=manager, active=true, status=active}
5. [1] TestPerson{tags=employee, level=3, name=Eve, salary=55000, score=8, age=30, active=false, status=inactive, department=sales}
6. [1] TestPerson{active=true, level=1, age=0, salary=-5000, tags=test, status=active, score=0, department=qa, name=Frank}
7. [1] TestPerson{name=Grace, score=10, salary=95000, tags=executive, status=active, department=management, level=9, age=65, active=true}
8. [1] TestPerson{salary=25000, active=false, score=5.5, level=1, status=inactive, name=Henry, age=18, department=support, tags=junior}
9. [1] TestPerson{level=6, salary=68000, status=active, age=40, name=Ivy, score=8.7, active=true, tags=senior, department=engineering}
10. [1] TestPerson{level=1, name=X, age=22, active=true, department=intern, tags=temp, salary=28000, score=6.5, status=active}

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

1. [1] TestPerson{status=active, name=Alice, score=8.5, level=2, age=25, active=true, department=sales, tags=junior, salary=45000}
2. [1] TestPerson{age=35, score=9.2, status=active, level=5, name=Bob, salary=75000, department=engineering, tags=senior, active=true}
3. [1] TestPerson{salary=0, active=false, level=1, age=16, department=hr, tags=intern, name=Charlie, score=6, status=inactive}
4. [1] TestPerson{department=marketing, age=45, salary=85000, score=7.8, level=7, name=Diana, tags=manager, active=true, status=active}
5. [1] TestPerson{status=inactive, department=sales, tags=employee, level=3, name=Eve, salary=55000, score=8, age=30, active=false}
6. [1] TestPerson{name=Frank, active=true, level=1, age=0, salary=-5000, tags=test, status=active, score=0, department=qa}
7. [1] TestPerson{department=management, level=9, age=65, active=true, name=Grace, score=10, salary=95000, tags=executive, status=active}
8. [1] TestPerson{salary=25000, active=false, score=5.5, level=1, status=inactive, name=Henry, age=18, department=support, tags=junior}
9. [1] TestPerson{department=engineering, level=6, salary=68000, status=active, age=40, name=Ivy, score=8.7, active=true, tags=senior}
10. [1] TestPerson{tags=temp, salary=28000, score=6.5, status=active, level=1, name=X, age=22, active=true, department=intern}
11. [1] TestOrder{amount=2, total=1999.98, date=2024-01-15, priority=normal, region=north, product_id=PROD001, status=pending, discount=50, customer_id=P001}
12. [1] TestOrder{date=2024-01-20, customer_id=P002, product_id=PROD002, discount=0, status=confirmed, priority=low, region=south, amount=1, total=25.5}
13. [1] TestOrder{date=2024-02-01, priority=high, discount=15, product_id=PROD003, customer_id=P001, status=shipped, region=north, amount=3, total=225}
14. [1] TestOrder{total=299.99, date=2024-02-05, priority=normal, product_id=PROD004, region=east, customer_id=P004, amount=1, status=delivered, discount=0}
15. [1] TestOrder{status=confirmed, product_id=PROD001, priority=high, discount=100, region=south, customer_id=P002, total=999.99, date=2024-02-10, amount=1}
16. [1] TestOrder{date=2024-02-15, customer_id=P005, status=cancelled, priority=low, amount=2, discount=0, region=west, product_id=PROD005, total=999.98}
17. [1] TestOrder{customer_id=P007, total=600, priority=urgent, discount=50, region=north, date=2024-03-01, product_id=PROD006, status=shipped, amount=4}
18. [1] TestOrder{customer_id=P010, amount=10, total=255, date=2024-03-05, discount=0, status=pending, product_id=PROD002, priority=normal, region=south}
19. [1] TestOrder{status=completed, customer_id=P001, amount=1, region=north, product_id=PROD007, priority=low, discount=10, total=89.99, date=2024-03-10}
20. [1] TestOrder{discount=0, amount=1, date=2024-03-15, region=east, priority=urgent, customer_id=P006, product_id=PROD001, total=75000, status=refunded}

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

1. [1] TestPerson{name=Alice, score=8.5, level=2, age=25, active=true, department=sales, tags=junior, salary=45000, status=active}
2. [1] TestPerson{active=true, age=35, score=9.2, status=active, level=5, name=Bob, salary=75000, department=engineering, tags=senior}
3. [1] TestPerson{score=6, status=inactive, salary=0, active=false, level=1, age=16, department=hr, tags=intern, name=Charlie}
4. [1] TestPerson{department=marketing, age=45, salary=85000, score=7.8, level=7, name=Diana, tags=manager, active=true, status=active}
5. [1] TestPerson{status=inactive, department=sales, tags=employee, level=3, name=Eve, salary=55000, score=8, age=30, active=false}
6. [1] TestPerson{active=true, level=1, age=0, salary=-5000, tags=test, status=active, score=0, department=qa, name=Frank}
7. [1] TestPerson{age=65, active=true, name=Grace, score=10, salary=95000, tags=executive, status=active, department=management, level=9}
8. [1] TestPerson{salary=25000, active=false, score=5.5, level=1, status=inactive, name=Henry, age=18, department=support, tags=junior}
9. [1] TestPerson{salary=68000, status=active, age=40, name=Ivy, score=8.7, active=true, tags=senior, department=engineering, level=6}
10. [1] TestPerson{tags=temp, salary=28000, score=6.5, status=active, level=1, name=X, age=22, active=true, department=intern}

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

1. [1] TestOrder{customer_id=P001, amount=2, total=1999.98, date=2024-01-15, priority=normal, region=north, product_id=PROD001, status=pending, discount=50}
2. [1] TestOrder{amount=1, total=25.5, date=2024-01-20, customer_id=P002, product_id=PROD002, discount=0, status=confirmed, priority=low, region=south}
3. [1] TestOrder{status=shipped, region=north, amount=3, total=225, date=2024-02-01, priority=high, discount=15, product_id=PROD003, customer_id=P001}
4. [1] TestOrder{product_id=PROD004, region=east, customer_id=P004, amount=1, status=delivered, discount=0, total=299.99, date=2024-02-05, priority=normal}
5. [1] TestOrder{product_id=PROD001, priority=high, discount=100, region=south, customer_id=P002, total=999.99, date=2024-02-10, amount=1, status=confirmed}
6. [1] TestOrder{date=2024-02-15, customer_id=P005, status=cancelled, priority=low, amount=2, discount=0, region=west, product_id=PROD005, total=999.98}
7. [1] TestOrder{product_id=PROD006, status=shipped, amount=4, customer_id=P007, total=600, priority=urgent, discount=50, region=north, date=2024-03-01}
8. [1] TestOrder{amount=10, total=255, date=2024-03-05, discount=0, status=pending, product_id=PROD002, priority=normal, region=south, customer_id=P010}
9. [1] TestOrder{amount=1, region=north, product_id=PROD007, priority=low, discount=10, total=89.99, date=2024-03-10, status=completed, customer_id=P001}
10. [1] TestOrder{priority=urgent, customer_id=P006, product_id=PROD001, total=75000, status=refunded, discount=0, amount=1, date=2024-03-15, region=east}

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

1. [1] TestPerson{status=active, name=Alice, score=8.5, level=2, age=25, active=true, department=sales, tags=junior, salary=45000}
2. [1] TestPerson{score=9.2, status=active, level=5, name=Bob, salary=75000, department=engineering, tags=senior, active=true, age=35}
3. [1] TestPerson{tags=intern, name=Charlie, score=6, status=inactive, salary=0, active=false, level=1, age=16, department=hr}
4. [1] TestPerson{department=marketing, age=45, salary=85000, score=7.8, level=7, name=Diana, tags=manager, active=true, status=active}
5. [1] TestPerson{score=8, age=30, active=false, status=inactive, department=sales, tags=employee, level=3, name=Eve, salary=55000}
6. [1] TestPerson{score=0, department=qa, name=Frank, active=true, level=1, age=0, salary=-5000, tags=test, status=active}
7. [1] TestPerson{level=9, age=65, active=true, name=Grace, score=10, salary=95000, tags=executive, status=active, department=management}
8. [1] TestPerson{age=18, department=support, tags=junior, salary=25000, active=false, score=5.5, level=1, status=inactive, name=Henry}
9. [1] TestPerson{active=true, tags=senior, department=engineering, level=6, salary=68000, status=active, age=40, name=Ivy, score=8.7}
10. [1] TestPerson{department=intern, tags=temp, salary=28000, score=6.5, status=active, level=1, name=X, age=22, active=true}

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
