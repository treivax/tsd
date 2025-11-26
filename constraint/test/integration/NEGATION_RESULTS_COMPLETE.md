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

1. [1] TestPerson{department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active}
2. [1] TestPerson{age=35, salary=75000, name=Bob, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5}
3. [1] TestPerson{age=16, status=inactive, active=false, tags=intern, name=Charlie, salary=0, level=1, score=6, department=hr}
4. [1] TestPerson{tags=manager, level=7, salary=85000, active=true, score=7.8, department=marketing, age=45, status=active, name=Diana}
5. [1] TestPerson{score=8, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, tags=employee, active=false}
6. [1] TestPerson{name=Frank, salary=-5000, age=0, active=true, score=0, tags=test, status=active, department=qa, level=1}
7. [1] TestPerson{department=management, salary=95000, name=Grace, score=10, active=true, tags=executive, status=active, level=9, age=65}
8. [1] TestPerson{department=support, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false, score=5.5, status=inactive}
9. [1] TestPerson{status=active, name=Ivy, age=40, department=engineering, salary=68000, level=6, active=true, score=8.7, tags=senior}
10. [1] TestPerson{salary=28000, status=active, level=1, name=X, score=6.5, tags=temp, active=true, age=22, department=intern}

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

1. [1] TestOrder{discount=50, region=north, amount=2, total=1999.98, priority=normal, customer_id=P001, date=2024-01-15, status=pending, product_id=PROD001}
2. [1] TestOrder{amount=1, total=25.5, date=2024-01-20, discount=0, status=confirmed, region=south, customer_id=P002, product_id=PROD002, priority=low}
3. [1] TestOrder{date=2024-02-01, status=shipped, priority=high, product_id=PROD003, region=north, customer_id=P001, discount=15, amount=3, total=225}
4. [1] TestOrder{customer_id=P004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, amount=1, discount=0}
5. [1] TestOrder{date=2024-02-10, product_id=PROD001, status=confirmed, discount=100, region=south, priority=high, customer_id=P002, amount=1, total=999.99}
6. [1] TestOrder{amount=2, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, total=999.98, priority=low}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, date=2024-03-01, status=shipped, region=north, discount=50, total=600, priority=urgent, amount=4}
8. [1] TestOrder{customer_id=P010, region=south, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, status=pending}
9. [1] TestOrder{discount=10, customer_id=P001, amount=1, total=89.99, product_id=PROD007, status=completed, region=north, date=2024-03-10, priority=low}
10. [1] TestOrder{priority=urgent, product_id=PROD001, amount=1, customer_id=P006, date=2024-03-15, status=refunded, region=east, discount=0, total=75000}

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

1. [1] TestPerson{department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active}
2. [1] TestPerson{score=9.2, tags=senior, status=active, department=engineering, level=5, age=35, salary=75000, name=Bob, active=true}
3. [1] TestPerson{department=hr, age=16, status=inactive, active=false, tags=intern, name=Charlie, salary=0, level=1, score=6}
4. [1] TestPerson{age=45, status=active, name=Diana, tags=manager, level=7, salary=85000, active=true, score=7.8, department=marketing}
5. [1] TestPerson{score=8, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, tags=employee, active=false}
6. [1] TestPerson{name=Frank, salary=-5000, age=0, active=true, score=0, tags=test, status=active, department=qa, level=1}
7. [1] TestPerson{department=management, salary=95000, name=Grace, score=10, active=true, tags=executive, status=active, level=9, age=65}
8. [1] TestPerson{tags=junior, level=1, name=Henry, age=18, active=false, score=5.5, status=inactive, department=support, salary=25000}
9. [1] TestPerson{score=8.7, tags=senior, status=active, name=Ivy, age=40, department=engineering, salary=68000, level=6, active=true}
10. [1] TestPerson{score=6.5, tags=temp, active=true, age=22, department=intern, salary=28000, status=active, level=1, name=X}

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

1. [1] TestOrder{priority=normal, customer_id=P001, date=2024-01-15, status=pending, product_id=PROD001, discount=50, region=north, amount=2, total=1999.98}
2. [1] TestOrder{discount=0, status=confirmed, region=south, customer_id=P002, product_id=PROD002, priority=low, amount=1, total=25.5, date=2024-01-20}
3. [1] TestOrder{discount=15, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, product_id=PROD003, region=north, customer_id=P001}
4. [1] TestOrder{discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, amount=1}
5. [1] TestOrder{priority=high, customer_id=P002, amount=1, total=999.99, date=2024-02-10, product_id=PROD001, status=confirmed, discount=100, region=south}
6. [1] TestOrder{region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, total=999.98, priority=low, amount=2}
7. [1] TestOrder{priority=urgent, amount=4, customer_id=P007, product_id=PROD006, date=2024-03-01, status=shipped, region=north, discount=50, total=600}
8. [1] TestOrder{total=255, date=2024-03-05, priority=normal, discount=0, status=pending, customer_id=P010, region=south, product_id=PROD002, amount=10}
9. [1] TestOrder{status=completed, region=north, date=2024-03-10, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, product_id=PROD007}
10. [1] TestOrder{priority=urgent, product_id=PROD001, amount=1, customer_id=P006, date=2024-03-15, status=refunded, region=east, discount=0, total=75000}

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

1. [1] TestProduct{available=true, category=electronics, brand=TechCorp, stock=50, name=Laptop, price=999.99, rating=4.5, keywords=computer, supplier=TechSupply}
2. [1] TestProduct{category=accessories, rating=4, brand=TechCorp, supplier=TechSupply, name=Mouse, price=25.5, available=true, stock=200, keywords=peripheral}
3. [1] TestProduct{name=Keyboard, price=75, rating=3.5, supplier=KeySupply, keywords=typing, brand=KeyTech, category=accessories, available=false, stock=0}
4. [1] TestProduct{price=299.99, keywords=display, supplier=ScreenSupply, name=Monitor, available=true, rating=4.8, brand=ScreenPro, stock=30, category=electronics}
5. [1] TestProduct{supplier=OldSupply, price=8.5, rating=2, available=false, name=OldKeyboard, keywords=obsolete, stock=0, category=accessories, brand=OldTech}
6. [1] TestProduct{category=audio, price=150, rating=4.6, brand=AudioMax, supplier=AudioSupply, name=Headphones, available=true, keywords=sound, stock=75}
7. [1] TestProduct{available=true, brand=CamTech, supplier=CamSupply, rating=3.8, category=electronics, stock=25, keywords=video, name=Webcam, price=89.99}

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

1. [1] TestPerson{salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, level=2, name=Alice, age=25}
2. [1] TestPerson{score=9.2, tags=senior, status=active, department=engineering, level=5, age=35, salary=75000, name=Bob, active=true}
3. [1] TestPerson{salary=0, level=1, score=6, department=hr, age=16, status=inactive, active=false, tags=intern, name=Charlie}
4. [1] TestPerson{salary=85000, active=true, score=7.8, department=marketing, age=45, status=active, name=Diana, tags=manager, level=7}
5. [1] TestPerson{score=8, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, tags=employee, active=false}
6. [1] TestPerson{name=Frank, salary=-5000, age=0, active=true, score=0, tags=test, status=active, department=qa, level=1}
7. [1] TestPerson{active=true, tags=executive, status=active, level=9, age=65, department=management, salary=95000, name=Grace, score=10}
8. [1] TestPerson{salary=25000, tags=junior, level=1, name=Henry, age=18, active=false, score=5.5, status=inactive, department=support}
9. [1] TestPerson{status=active, name=Ivy, age=40, department=engineering, salary=68000, level=6, active=true, score=8.7, tags=senior}
10. [1] TestPerson{salary=28000, status=active, level=1, name=X, score=6.5, tags=temp, active=true, age=22, department=intern}

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

1. [1] TestOrder{amount=2, total=1999.98, priority=normal, customer_id=P001, date=2024-01-15, status=pending, product_id=PROD001, discount=50, region=north}
2. [1] TestOrder{customer_id=P002, product_id=PROD002, priority=low, amount=1, total=25.5, date=2024-01-20, discount=0, status=confirmed, region=south}
3. [1] TestOrder{amount=3, total=225, date=2024-02-01, status=shipped, priority=high, product_id=PROD003, region=north, customer_id=P001, discount=15}
4. [1] TestOrder{status=delivered, priority=normal, region=east, product_id=PROD004, amount=1, discount=0, customer_id=P004, total=299.99, date=2024-02-05}
5. [1] TestOrder{priority=high, customer_id=P002, amount=1, total=999.99, date=2024-02-10, product_id=PROD001, status=confirmed, discount=100, region=south}
6. [1] TestOrder{discount=0, product_id=PROD005, total=999.98, priority=low, amount=2, region=west, customer_id=P005, date=2024-02-15, status=cancelled}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, date=2024-03-01, status=shipped, region=north, discount=50, total=600, priority=urgent, amount=4}
8. [1] TestOrder{amount=10, total=255, date=2024-03-05, priority=normal, discount=0, status=pending, customer_id=P010, region=south, product_id=PROD002}
9. [1] TestOrder{discount=10, customer_id=P001, amount=1, total=89.99, product_id=PROD007, status=completed, region=north, date=2024-03-10, priority=low}
10. [1] TestOrder{region=east, discount=0, total=75000, priority=urgent, product_id=PROD001, amount=1, customer_id=P006, date=2024-03-15, status=refunded}

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

1. [1] TestPerson{active=true, score=8.5, tags=junior, status=active, department=sales, level=2, name=Alice, age=25, salary=45000}
2. [1] TestPerson{active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, age=35, salary=75000, name=Bob}
3. [1] TestPerson{age=16, status=inactive, active=false, tags=intern, name=Charlie, salary=0, level=1, score=6, department=hr}
4. [1] TestPerson{salary=85000, active=true, score=7.8, department=marketing, age=45, status=active, name=Diana, tags=manager, level=7}
5. [1] TestPerson{score=8, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, tags=employee, active=false}
6. [1] TestPerson{score=0, tags=test, status=active, department=qa, level=1, name=Frank, salary=-5000, age=0, active=true}
7. [1] TestPerson{active=true, tags=executive, status=active, level=9, age=65, department=management, salary=95000, name=Grace, score=10}
8. [1] TestPerson{department=support, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false, score=5.5, status=inactive}
9. [1] TestPerson{active=true, score=8.7, tags=senior, status=active, name=Ivy, age=40, department=engineering, salary=68000, level=6}
10. [1] TestPerson{age=22, department=intern, salary=28000, status=active, level=1, name=X, score=6.5, tags=temp, active=true}

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

1. [1] TestOrder{status=pending, product_id=PROD001, discount=50, region=north, amount=2, total=1999.98, priority=normal, customer_id=P001, date=2024-01-15}
2. [1] TestOrder{amount=1, total=25.5, date=2024-01-20, discount=0, status=confirmed, region=south, customer_id=P002, product_id=PROD002, priority=low}
3. [1] TestOrder{amount=3, total=225, date=2024-02-01, status=shipped, priority=high, product_id=PROD003, region=north, customer_id=P001, discount=15}
4. [1] TestOrder{amount=1, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004}
5. [1] TestOrder{priority=high, customer_id=P002, amount=1, total=999.99, date=2024-02-10, product_id=PROD001, status=confirmed, discount=100, region=south}
6. [1] TestOrder{total=999.98, priority=low, amount=2, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, date=2024-03-01, status=shipped, region=north, discount=50, total=600, priority=urgent, amount=4}
8. [1] TestOrder{total=255, date=2024-03-05, priority=normal, discount=0, status=pending, customer_id=P010, region=south, product_id=PROD002, amount=10}
9. [1] TestOrder{customer_id=P001, amount=1, total=89.99, product_id=PROD007, status=completed, region=north, date=2024-03-10, priority=low, discount=10}
10. [1] TestOrder{amount=1, customer_id=P006, date=2024-03-15, status=refunded, region=east, discount=0, total=75000, priority=urgent, product_id=PROD001}

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

1. [1] TestPerson{status=active, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior}
2. [1] TestPerson{name=Bob, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, age=35, salary=75000}
3. [1] TestPerson{status=inactive, active=false, tags=intern, name=Charlie, salary=0, level=1, score=6, department=hr, age=16}
4. [1] TestPerson{salary=85000, active=true, score=7.8, department=marketing, age=45, status=active, name=Diana, tags=manager, level=7}
5. [1] TestPerson{score=8, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, tags=employee, active=false}
6. [1] TestPerson{tags=test, status=active, department=qa, level=1, name=Frank, salary=-5000, age=0, active=true, score=0}
7. [1] TestPerson{active=true, tags=executive, status=active, level=9, age=65, department=management, salary=95000, name=Grace, score=10}
8. [1] TestPerson{department=support, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false, score=5.5, status=inactive}
9. [1] TestPerson{score=8.7, tags=senior, status=active, name=Ivy, age=40, department=engineering, salary=68000, level=6, active=true}
10. [1] TestPerson{level=1, name=X, score=6.5, tags=temp, active=true, age=22, department=intern, salary=28000, status=active}

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

1. [1] TestProduct{rating=4.5, keywords=computer, supplier=TechSupply, available=true, category=electronics, brand=TechCorp, stock=50, name=Laptop, price=999.99}
2. [1] TestProduct{stock=200, keywords=peripheral, category=accessories, rating=4, brand=TechCorp, supplier=TechSupply, name=Mouse, price=25.5, available=true}
3. [1] TestProduct{keywords=typing, brand=KeyTech, category=accessories, available=false, stock=0, name=Keyboard, price=75, rating=3.5, supplier=KeySupply}
4. [1] TestProduct{rating=4.8, brand=ScreenPro, stock=30, category=electronics, price=299.99, keywords=display, supplier=ScreenSupply, name=Monitor, available=true}
5. [1] TestProduct{supplier=OldSupply, price=8.5, rating=2, available=false, name=OldKeyboard, keywords=obsolete, stock=0, category=accessories, brand=OldTech}
6. [1] TestProduct{keywords=sound, stock=75, category=audio, price=150, rating=4.6, brand=AudioMax, supplier=AudioSupply, name=Headphones, available=true}
7. [1] TestProduct{name=Webcam, price=89.99, available=true, brand=CamTech, supplier=CamSupply, rating=3.8, category=electronics, stock=25, keywords=video}

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

1. [1] TestPerson{department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active}
2. [1] TestPerson{tags=senior, status=active, department=engineering, level=5, age=35, salary=75000, name=Bob, active=true, score=9.2}
3. [1] TestPerson{age=16, status=inactive, active=false, tags=intern, name=Charlie, salary=0, level=1, score=6, department=hr}
4. [1] TestPerson{age=45, status=active, name=Diana, tags=manager, level=7, salary=85000, active=true, score=7.8, department=marketing}
5. [1] TestPerson{score=8, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, tags=employee, active=false}
6. [1] TestPerson{name=Frank, salary=-5000, age=0, active=true, score=0, tags=test, status=active, department=qa, level=1}
7. [1] TestPerson{salary=95000, name=Grace, score=10, active=true, tags=executive, status=active, level=9, age=65, department=management}
8. [1] TestPerson{salary=25000, tags=junior, level=1, name=Henry, age=18, active=false, score=5.5, status=inactive, department=support}
9. [1] TestPerson{status=active, name=Ivy, age=40, department=engineering, salary=68000, level=6, active=true, score=8.7, tags=senior}
10. [1] TestPerson{age=22, department=intern, salary=28000, status=active, level=1, name=X, score=6.5, tags=temp, active=true}

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

1. [1] TestOrder{discount=50, region=north, amount=2, total=1999.98, priority=normal, customer_id=P001, date=2024-01-15, status=pending, product_id=PROD001}
2. [1] TestOrder{discount=0, status=confirmed, region=south, customer_id=P002, product_id=PROD002, priority=low, amount=1, total=25.5, date=2024-01-20}
3. [1] TestOrder{date=2024-02-01, status=shipped, priority=high, product_id=PROD003, region=north, customer_id=P001, discount=15, amount=3, total=225}
4. [1] TestOrder{total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, amount=1, discount=0, customer_id=P004}
5. [1] TestOrder{priority=high, customer_id=P002, amount=1, total=999.99, date=2024-02-10, product_id=PROD001, status=confirmed, discount=100, region=south}
6. [1] TestOrder{region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, total=999.98, priority=low, amount=2}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, date=2024-03-01, status=shipped, region=north, discount=50, total=600, priority=urgent, amount=4}
8. [1] TestOrder{total=255, date=2024-03-05, priority=normal, discount=0, status=pending, customer_id=P010, region=south, product_id=PROD002, amount=10}
9. [1] TestOrder{date=2024-03-10, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, product_id=PROD007, status=completed, region=north}
10. [1] TestOrder{amount=1, customer_id=P006, date=2024-03-15, status=refunded, region=east, discount=0, total=75000, priority=urgent, product_id=PROD001}

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

1. [1] TestPerson{department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active}
2. [1] TestPerson{active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, age=35, salary=75000, name=Bob}
3. [1] TestPerson{name=Charlie, salary=0, level=1, score=6, department=hr, age=16, status=inactive, active=false, tags=intern}
4. [1] TestPerson{salary=85000, active=true, score=7.8, department=marketing, age=45, status=active, name=Diana, tags=manager, level=7}
5. [1] TestPerson{score=8, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, tags=employee, active=false}
6. [1] TestPerson{department=qa, level=1, name=Frank, salary=-5000, age=0, active=true, score=0, tags=test, status=active}
7. [1] TestPerson{department=management, salary=95000, name=Grace, score=10, active=true, tags=executive, status=active, level=9, age=65}
8. [1] TestPerson{salary=25000, tags=junior, level=1, name=Henry, age=18, active=false, score=5.5, status=inactive, department=support}
9. [1] TestPerson{level=6, active=true, score=8.7, tags=senior, status=active, name=Ivy, age=40, department=engineering, salary=68000}
10. [1] TestPerson{age=22, department=intern, salary=28000, status=active, level=1, name=X, score=6.5, tags=temp, active=true}
11. [1] TestOrder{discount=50, region=north, amount=2, total=1999.98, priority=normal, customer_id=P001, date=2024-01-15, status=pending, product_id=PROD001}
12. [1] TestOrder{total=25.5, date=2024-01-20, discount=0, status=confirmed, region=south, customer_id=P002, product_id=PROD002, priority=low, amount=1}
13. [1] TestOrder{priority=high, product_id=PROD003, region=north, customer_id=P001, discount=15, amount=3, total=225, date=2024-02-01, status=shipped}
14. [1] TestOrder{customer_id=P004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, amount=1, discount=0}
15. [1] TestOrder{date=2024-02-10, product_id=PROD001, status=confirmed, discount=100, region=south, priority=high, customer_id=P002, amount=1, total=999.99}
16. [1] TestOrder{discount=0, product_id=PROD005, total=999.98, priority=low, amount=2, region=west, customer_id=P005, date=2024-02-15, status=cancelled}
17. [1] TestOrder{amount=4, customer_id=P007, product_id=PROD006, date=2024-03-01, status=shipped, region=north, discount=50, total=600, priority=urgent}
18. [1] TestOrder{total=255, date=2024-03-05, priority=normal, discount=0, status=pending, customer_id=P010, region=south, product_id=PROD002, amount=10}
19. [1] TestOrder{discount=10, customer_id=P001, amount=1, total=89.99, product_id=PROD007, status=completed, region=north, date=2024-03-10, priority=low}
20. [1] TestOrder{amount=1, customer_id=P006, date=2024-03-15, status=refunded, region=east, discount=0, total=75000, priority=urgent, product_id=PROD001}

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

1. [1] TestPerson{tags=junior, status=active, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5}
2. [1] TestPerson{tags=senior, status=active, department=engineering, level=5, age=35, salary=75000, name=Bob, active=true, score=9.2}
3. [1] TestPerson{status=inactive, active=false, tags=intern, name=Charlie, salary=0, level=1, score=6, department=hr, age=16}
4. [1] TestPerson{score=7.8, department=marketing, age=45, status=active, name=Diana, tags=manager, level=7, salary=85000, active=true}
5. [1] TestPerson{tags=employee, active=false, score=8, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000}
6. [1] TestPerson{tags=test, status=active, department=qa, level=1, name=Frank, salary=-5000, age=0, active=true, score=0}
7. [1] TestPerson{active=true, tags=executive, status=active, level=9, age=65, department=management, salary=95000, name=Grace, score=10}
8. [1] TestPerson{score=5.5, status=inactive, department=support, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false}
9. [1] TestPerson{active=true, score=8.7, tags=senior, status=active, name=Ivy, age=40, department=engineering, salary=68000, level=6}
10. [1] TestPerson{salary=28000, status=active, level=1, name=X, score=6.5, tags=temp, active=true, age=22, department=intern}

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

1. [1] TestPerson{age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, level=2, name=Alice}
2. [1] TestPerson{tags=senior, status=active, department=engineering, level=5, age=35, salary=75000, name=Bob, active=true, score=9.2}
3. [1] TestPerson{department=hr, age=16, status=inactive, active=false, tags=intern, name=Charlie, salary=0, level=1, score=6}
4. [1] TestPerson{status=active, name=Diana, tags=manager, level=7, salary=85000, active=true, score=7.8, department=marketing, age=45}
5. [1] TestPerson{tags=employee, active=false, score=8, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000}
6. [1] TestPerson{age=0, active=true, score=0, tags=test, status=active, department=qa, level=1, name=Frank, salary=-5000}
7. [1] TestPerson{level=9, age=65, department=management, salary=95000, name=Grace, score=10, active=true, tags=executive, status=active}
8. [1] TestPerson{name=Henry, age=18, active=false, score=5.5, status=inactive, department=support, salary=25000, tags=junior, level=1}
9. [1] TestPerson{salary=68000, level=6, active=true, score=8.7, tags=senior, status=active, name=Ivy, age=40, department=engineering}
10. [1] TestPerson{age=22, department=intern, salary=28000, status=active, level=1, name=X, score=6.5, tags=temp, active=true}
11. [1] TestOrder{amount=2, total=1999.98, priority=normal, customer_id=P001, date=2024-01-15, status=pending, product_id=PROD001, discount=50, region=north}
12. [1] TestOrder{amount=1, total=25.5, date=2024-01-20, discount=0, status=confirmed, region=south, customer_id=P002, product_id=PROD002, priority=low}
13. [1] TestOrder{discount=15, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, product_id=PROD003, region=north, customer_id=P001}
14. [1] TestOrder{customer_id=P004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, amount=1, discount=0}
15. [1] TestOrder{priority=high, customer_id=P002, amount=1, total=999.99, date=2024-02-10, product_id=PROD001, status=confirmed, discount=100, region=south}
16. [1] TestOrder{total=999.98, priority=low, amount=2, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005}
17. [1] TestOrder{product_id=PROD006, date=2024-03-01, status=shipped, region=north, discount=50, total=600, priority=urgent, amount=4, customer_id=P007}
18. [1] TestOrder{discount=0, status=pending, customer_id=P010, region=south, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal}
19. [1] TestOrder{amount=1, total=89.99, product_id=PROD007, status=completed, region=north, date=2024-03-10, priority=low, discount=10, customer_id=P001}
20. [1] TestOrder{priority=urgent, product_id=PROD001, amount=1, customer_id=P006, date=2024-03-15, status=refunded, region=east, discount=0, total=75000}

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

1. [1] TestPerson{department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active}
2. [1] TestPerson{score=9.2, tags=senior, status=active, department=engineering, level=5, age=35, salary=75000, name=Bob, active=true}
3. [1] TestPerson{level=1, score=6, department=hr, age=16, status=inactive, active=false, tags=intern, name=Charlie, salary=0}
4. [1] TestPerson{salary=85000, active=true, score=7.8, department=marketing, age=45, status=active, name=Diana, tags=manager, level=7}
5. [1] TestPerson{score=8, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, tags=employee, active=false}
6. [1] TestPerson{name=Frank, salary=-5000, age=0, active=true, score=0, tags=test, status=active, department=qa, level=1}
7. [1] TestPerson{department=management, salary=95000, name=Grace, score=10, active=true, tags=executive, status=active, level=9, age=65}
8. [1] TestPerson{score=5.5, status=inactive, department=support, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false}
9. [1] TestPerson{score=8.7, tags=senior, status=active, name=Ivy, age=40, department=engineering, salary=68000, level=6, active=true}
10. [1] TestPerson{age=22, department=intern, salary=28000, status=active, level=1, name=X, score=6.5, tags=temp, active=true}

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

1. [1] TestOrder{discount=50, region=north, amount=2, total=1999.98, priority=normal, customer_id=P001, date=2024-01-15, status=pending, product_id=PROD001}
2. [1] TestOrder{amount=1, total=25.5, date=2024-01-20, discount=0, status=confirmed, region=south, customer_id=P002, product_id=PROD002, priority=low}
3. [1] TestOrder{amount=3, total=225, date=2024-02-01, status=shipped, priority=high, product_id=PROD003, region=north, customer_id=P001, discount=15}
4. [1] TestOrder{product_id=PROD004, amount=1, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east}
5. [1] TestOrder{customer_id=P002, amount=1, total=999.99, date=2024-02-10, product_id=PROD001, status=confirmed, discount=100, region=south, priority=high}
6. [1] TestOrder{amount=2, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, total=999.98, priority=low}
7. [1] TestOrder{total=600, priority=urgent, amount=4, customer_id=P007, product_id=PROD006, date=2024-03-01, status=shipped, region=north, discount=50}
8. [1] TestOrder{total=255, date=2024-03-05, priority=normal, discount=0, status=pending, customer_id=P010, region=south, product_id=PROD002, amount=10}
9. [1] TestOrder{priority=low, discount=10, customer_id=P001, amount=1, total=89.99, product_id=PROD007, status=completed, region=north, date=2024-03-10}
10. [1] TestOrder{priority=urgent, product_id=PROD001, amount=1, customer_id=P006, date=2024-03-15, status=refunded, region=east, discount=0, total=75000}

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

1. [1] TestPerson{department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active}
2. [1] TestPerson{name=Bob, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, age=35, salary=75000}
3. [1] TestPerson{status=inactive, active=false, tags=intern, name=Charlie, salary=0, level=1, score=6, department=hr, age=16}
4. [1] TestPerson{active=true, score=7.8, department=marketing, age=45, status=active, name=Diana, tags=manager, level=7, salary=85000}
5. [1] TestPerson{score=8, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, tags=employee, active=false}
6. [1] TestPerson{department=qa, level=1, name=Frank, salary=-5000, age=0, active=true, score=0, tags=test, status=active}
7. [1] TestPerson{tags=executive, status=active, level=9, age=65, department=management, salary=95000, name=Grace, score=10, active=true}
8. [1] TestPerson{status=inactive, department=support, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false, score=5.5}
9. [1] TestPerson{score=8.7, tags=senior, status=active, name=Ivy, age=40, department=engineering, salary=68000, level=6, active=true}
10. [1] TestPerson{active=true, age=22, department=intern, salary=28000, status=active, level=1, name=X, score=6.5, tags=temp}

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
