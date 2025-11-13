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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales}
2. [1] TestPerson{id=P002, age=35, department=engineering, level=5, name=Bob, salary=75000, active=true, score=9.2, tags=senior, status=active}
3. [1] TestPerson{id=P003, level=1, name=Charlie, active=false, tags=intern, status=inactive, age=16, salary=0, score=6, department=hr}
4. [1] TestPerson{id=P004, level=7, name=Diana, salary=85000, status=active, age=45, active=true, score=7.8, tags=manager, department=marketing}
5. [1] TestPerson{id=P005, score=8, status=inactive, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales, level=3}
6. [1] TestPerson{id=P006, tags=test, level=1, salary=-5000, status=active, department=qa, name=Frank, age=0, active=true, score=0}
7. [1] TestPerson{id=P007, name=Grace, status=active, age=65, salary=95000, active=true, score=10, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, age=18, salary=25000, active=false, score=5.5, level=1, name=Henry}
9. [1] TestPerson{id=P009, score=8.7, department=engineering, name=Ivy, age=40, active=true, tags=senior, status=active, level=6, salary=68000}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true, department=intern, name=X}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, level=5, name=Bob, salary=75000, active=true, score=9.2, tags=senior, status=active, age=35}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, score=7.8, tags=manager, department=marketing, level=7, name=Diana, salary=85000, status=active}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true, department=intern}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales, name=Alice, age=25}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, score=6, department=hr, level=1, name=Charlie, active=false, tags=intern, status=inactive}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, department=sales, level=3, score=8, status=inactive, name=Eve, age=30, salary=55000, active=false}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, tags=executive, department=management, level=9, name=Grace, status=active}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, level=1, name=Henry, tags=junior, status=inactive, department=support, age=18}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, score=8.7, department=engineering, name=Ivy, age=40, active=true, tags=senior, status=active, level=6}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, status=pending, discount=50, customer_id=P001, product_id=PROD001, amount=2, priority=normal, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, priority=low, discount=0, region=south, total=25.5, date=2024-01-20, status=confirmed}
3. [1] TestOrder{id=O003, discount=15, region=north, total=225, status=shipped, priority=high, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01}
4. [1] TestOrder{id=O004, customer_id=P004, total=299.99, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, amount=1, date=2024-02-05}
5. [1] TestOrder{id=O005, customer_id=P002, total=999.99, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, date=2024-02-10, priority=high}
6. [1] TestOrder{id=O006, customer_id=P005, total=999.98, date=2024-02-15, discount=0, region=west, product_id=PROD005, amount=2, status=cancelled, priority=low}
7. [1] TestOrder{id=O007, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, region=north, customer_id=P007, amount=4, total=600, priority=urgent}
8. [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, customer_id=P010, product_id=PROD002, region=south}
9. [1] TestOrder{id=O009, amount=1, date=2024-03-10, discount=10, region=north, total=89.99, status=completed, priority=low, customer_id=P001, product_id=PROD007}
10. [1] TestOrder{id=O010, priority=urgent, region=east, total=75000, discount=0, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded}

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

1. [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales, name=Alice, age=25}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, status=active, age=35, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, active=false, tags=intern, status=inactive, age=16, salary=0, score=6, department=hr, level=1}
4. [1] TestPerson{id=P004, age=45, active=true, score=7.8, tags=manager, department=marketing, level=7, name=Diana, salary=85000, status=active}
5. [1] TestPerson{id=P005, active=false, tags=employee, department=sales, level=3, score=8, status=inactive, name=Eve, age=30, salary=55000}
6. [1] TestPerson{id=P006, level=1, salary=-5000, status=active, department=qa, name=Frank, age=0, active=true, score=0, tags=test}
7. [1] TestPerson{id=P007, level=9, name=Grace, status=active, age=65, salary=95000, active=true, score=10, tags=executive, department=management}
8. [1] TestPerson{id=P008, age=18, salary=25000, active=false, score=5.5, level=1, name=Henry, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true, department=intern}

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

1. [1] TestOrder{id=O001, discount=50, customer_id=P001, product_id=PROD001, amount=2, priority=normal, region=north, total=1999.98, date=2024-01-15, status=pending}
2. [1] TestOrder{id=O002, status=confirmed, customer_id=P002, product_id=PROD002, amount=1, priority=low, discount=0, region=south, total=25.5, date=2024-01-20}
3. [1] TestOrder{id=O003, date=2024-02-01, discount=15, region=north, total=225, status=shipped, priority=high, customer_id=P001, product_id=PROD003, amount=3}
4. [1] TestOrder{id=O004, region=east, product_id=PROD004, amount=1, date=2024-02-05, customer_id=P004, total=299.99, status=delivered, priority=normal, discount=0}
5. [1] TestOrder{id=O005, customer_id=P002, total=999.99, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, date=2024-02-10, priority=high}
6. [1] TestOrder{id=O006, customer_id=P005, total=999.98, date=2024-02-15, discount=0, region=west, product_id=PROD005, amount=2, status=cancelled, priority=low}
7. [1] TestOrder{id=O007, amount=4, total=600, priority=urgent, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, region=north, customer_id=P007}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, region=south, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0}
9. [1] TestOrder{id=O009, region=north, total=89.99, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, discount=10}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, region=east, total=75000, discount=0, customer_id=P006}

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

1. [1] TestProduct{id=PROD001, rating=4.5, brand=TechCorp, stock=50, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, available=true, keywords=computer}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, price=25.5, available=true, rating=4, brand=TechCorp, stock=200, keywords=peripheral, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, stock=0, supplier=KeySupply, price=75, available=false, rating=3.5, keywords=typing, name=Keyboard, category=accessories, brand=KeyTech}
4. [1] TestProduct{id=PROD004, rating=4.8, keywords=display, brand=ScreenPro, stock=30, name=Monitor, category=electronics, price=299.99, available=true, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, brand=OldTech, stock=0, supplier=OldSupply, name=OldKeyboard, category=accessories, price=8.5, available=false, rating=2, keywords=obsolete}
6. [1] TestProduct{id=PROD006, name=Headphones, category=audio, price=150, supplier=AudioSupply, available=true, rating=4.6, keywords=sound, brand=AudioMax, stock=75}
7. [1] TestProduct{id=PROD007, category=electronics, price=89.99, stock=25, name=Webcam, available=true, rating=3.8, keywords=video, brand=CamTech, supplier=CamSupply}

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

1. [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales, name=Alice}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, status=active, age=35, department=engineering, level=5}
3. [1] TestPerson{id=P003, level=1, name=Charlie, active=false, tags=intern, status=inactive, age=16, salary=0, score=6, department=hr}
4. [1] TestPerson{id=P004, salary=85000, status=active, age=45, active=true, score=7.8, tags=manager, department=marketing, level=7, name=Diana}
5. [1] TestPerson{id=P005, age=30, salary=55000, active=false, tags=employee, department=sales, level=3, score=8, status=inactive, name=Eve}
6. [1] TestPerson{id=P006, level=1, salary=-5000, status=active, department=qa, name=Frank, age=0, active=true, score=0, tags=test}
7. [1] TestPerson{id=P007, level=9, name=Grace, status=active, age=65, salary=95000, active=true, score=10, tags=executive, department=management}
8. [1] TestPerson{id=P008, department=support, age=18, salary=25000, active=false, score=5.5, level=1, name=Henry, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, score=8.7, department=engineering, name=Ivy, age=40, active=true, tags=senior, status=active, level=6, salary=68000}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true, department=intern}

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

1. [1] TestOrder{id=O001, amount=2, priority=normal, region=north, total=1999.98, date=2024-01-15, status=pending, discount=50, customer_id=P001, product_id=PROD001}
2. [1] TestOrder{id=O002, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, amount=1, priority=low, discount=0, region=south, total=25.5}
3. [1] TestOrder{id=O003, date=2024-02-01, discount=15, region=north, total=225, status=shipped, priority=high, customer_id=P001, product_id=PROD003, amount=3}
4. [1] TestOrder{id=O004, customer_id=P004, total=299.99, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, amount=1, date=2024-02-05}
5. [1] TestOrder{id=O005, customer_id=P002, total=999.99, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, date=2024-02-10, priority=high}
6. [1] TestOrder{id=O006, customer_id=P005, total=999.98, date=2024-02-15, discount=0, region=west, product_id=PROD005, amount=2, status=cancelled, priority=low}
7. [1] TestOrder{id=O007, product_id=PROD006, date=2024-03-01, status=shipped, region=north, customer_id=P007, amount=4, total=600, priority=urgent, discount=50}
8. [1] TestOrder{id=O008, discount=0, customer_id=P010, product_id=PROD002, region=south, amount=10, total=255, date=2024-03-05, status=pending, priority=normal}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, discount=10, region=north, total=89.99, status=completed, priority=low}
10. [1] TestOrder{id=O010, amount=1, date=2024-03-15, status=refunded, priority=urgent, region=east, total=75000, discount=0, customer_id=P006, product_id=PROD001}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, status=active, age=35, department=engineering, level=5, name=Bob, salary=75000, active=true}
3. [1] TestPerson{id=P003, age=16, salary=0, score=6, department=hr, level=1, name=Charlie, active=false, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, age=45, active=true, score=7.8, tags=manager, department=marketing, level=7, name=Diana, salary=85000, status=active}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales, level=3, score=8, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, level=1, salary=-5000, status=active, department=qa}
7. [1] TestPerson{id=P007, department=management, level=9, name=Grace, status=active, age=65, salary=95000, active=true, score=10, tags=executive}
8. [1] TestPerson{id=P008, name=Henry, tags=junior, status=inactive, department=support, age=18, salary=25000, active=false, score=5.5, level=1}
9. [1] TestPerson{id=P009, status=active, level=6, salary=68000, score=8.7, department=engineering, name=Ivy, age=40, active=true, tags=senior}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true, department=intern}

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

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, priority=normal, region=north, total=1999.98, date=2024-01-15, status=pending, discount=50}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, priority=low, discount=0, region=south, total=25.5, date=2024-01-20, status=confirmed}
3. [1] TestOrder{id=O003, priority=high, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, discount=15, region=north, total=225, status=shipped}
4. [1] TestOrder{id=O004, customer_id=P004, total=299.99, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, amount=1, date=2024-02-05}
5. [1] TestOrder{id=O005, customer_id=P002, total=999.99, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, date=2024-02-10, priority=high}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, status=cancelled, priority=low, customer_id=P005, total=999.98, date=2024-02-15, discount=0, region=west}
7. [1] TestOrder{id=O007, amount=4, total=600, priority=urgent, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, region=north, customer_id=P007}
8. [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, customer_id=P010, product_id=PROD002, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, discount=10, region=north, total=89.99, status=completed, priority=low}
10. [1] TestOrder{id=O010, total=75000, discount=0, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, region=east}

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

1. [1] TestPerson{id=P001, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true}
2. [1] TestPerson{id=P002, age=35, department=engineering, level=5, name=Bob, salary=75000, active=true, score=9.2, tags=senior, status=active}
3. [1] TestPerson{id=P003, level=1, name=Charlie, active=false, tags=intern, status=inactive, age=16, salary=0, score=6, department=hr}
4. [1] TestPerson{id=P004, status=active, age=45, active=true, score=7.8, tags=manager, department=marketing, level=7, name=Diana, salary=85000}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales, level=3, score=8, status=inactive}
6. [1] TestPerson{id=P006, active=true, score=0, tags=test, level=1, salary=-5000, status=active, department=qa, name=Frank, age=0}
7. [1] TestPerson{id=P007, department=management, level=9, name=Grace, status=active, age=65, salary=95000, active=true, score=10, tags=executive}
8. [1] TestPerson{id=P008, score=5.5, level=1, name=Henry, tags=junior, status=inactive, department=support, age=18, salary=25000, active=false}
9. [1] TestPerson{id=P009, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering, name=Ivy, age=40, active=true}
10. [1] TestPerson{id=P010, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true, department=intern, name=X, age=22}

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

1. [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, available=true, keywords=computer, rating=4.5, brand=TechCorp, stock=50, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, price=25.5, available=true, rating=4, brand=TechCorp, stock=200, keywords=peripheral, supplier=TechSupply, name=Mouse, category=accessories}
3. [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, brand=KeyTech, stock=0, supplier=KeySupply, price=75, available=false, rating=3.5, keywords=typing}
4. [1] TestProduct{id=PROD004, rating=4.8, keywords=display, brand=ScreenPro, stock=30, name=Monitor, category=electronics, price=299.99, available=true, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, keywords=obsolete, brand=OldTech, stock=0, supplier=OldSupply, name=OldKeyboard, category=accessories, price=8.5, available=false, rating=2}
6. [1] TestProduct{id=PROD006, available=true, rating=4.6, keywords=sound, brand=AudioMax, stock=75, name=Headphones, category=audio, price=150, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, category=electronics, price=89.99, stock=25, name=Webcam, available=true, rating=3.8, keywords=video, brand=CamTech, supplier=CamSupply}

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

1. [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales, name=Alice}
2. [1] TestPerson{id=P002, tags=senior, status=active, age=35, department=engineering, level=5, name=Bob, salary=75000, active=true, score=9.2}
3. [1] TestPerson{id=P003, age=16, salary=0, score=6, department=hr, level=1, name=Charlie, active=false, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, age=45, active=true, score=7.8, tags=manager, department=marketing, level=7, name=Diana, salary=85000, status=active}
5. [1] TestPerson{id=P005, active=false, tags=employee, department=sales, level=3, score=8, status=inactive, name=Eve, age=30, salary=55000}
6. [1] TestPerson{id=P006, department=qa, name=Frank, age=0, active=true, score=0, tags=test, level=1, salary=-5000, status=active}
7. [1] TestPerson{id=P007, department=management, level=9, name=Grace, status=active, age=65, salary=95000, active=true, score=10, tags=executive}
8. [1] TestPerson{id=P008, active=false, score=5.5, level=1, name=Henry, tags=junior, status=inactive, department=support, age=18, salary=25000}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering}
10. [1] TestPerson{id=P010, score=6.5, tags=temp, status=active, level=1, active=true, department=intern, name=X, age=22, salary=28000}

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

1. [1] TestOrder{id=O001, amount=2, priority=normal, region=north, total=1999.98, date=2024-01-15, status=pending, discount=50, customer_id=P001, product_id=PROD001}
2. [1] TestOrder{id=O002, status=confirmed, customer_id=P002, product_id=PROD002, amount=1, priority=low, discount=0, region=south, total=25.5, date=2024-01-20}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, discount=15, region=north, total=225, status=shipped, priority=high}
4. [1] TestOrder{id=O004, customer_id=P004, total=299.99, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, amount=1, date=2024-02-05}
5. [1] TestOrder{id=O005, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, date=2024-02-10, priority=high, customer_id=P002, total=999.99}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, status=cancelled, priority=low, customer_id=P005, total=999.98, date=2024-02-15, discount=0, region=west}
7. [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, priority=urgent, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, region=north}
8. [1] TestOrder{id=O008, region=south, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, customer_id=P010, product_id=PROD002}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, discount=10, region=north, total=89.99, status=completed, priority=low}
10. [1] TestOrder{id=O010, priority=urgent, region=east, total=75000, discount=0, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded}

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

1. [1] TestPerson{id=P001, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, status=active, age=35, department=engineering, level=5}
3. [1] TestPerson{id=P003, status=inactive, age=16, salary=0, score=6, department=hr, level=1, name=Charlie, active=false, tags=intern}
4. [1] TestPerson{id=P004, age=45, active=true, score=7.8, tags=manager, department=marketing, level=7, name=Diana, salary=85000, status=active}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales, level=3, score=8, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, level=1, salary=-5000, status=active, department=qa}
7. [1] TestPerson{id=P007, active=true, score=10, tags=executive, department=management, level=9, name=Grace, status=active, age=65, salary=95000}
8. [1] TestPerson{id=P008, age=18, salary=25000, active=false, score=5.5, level=1, name=Henry, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering}
10. [1] TestPerson{id=P010, score=6.5, tags=temp, status=active, level=1, active=true, department=intern, name=X, age=22, salary=28000}
11. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, priority=normal, region=north, total=1999.98, date=2024-01-15, status=pending, discount=50}
12. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, priority=low, discount=0, region=south, total=25.5, date=2024-01-20, status=confirmed}
13. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, discount=15, region=north, total=225, status=shipped, priority=high}
14. [1] TestOrder{id=O004, customer_id=P004, total=299.99, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, amount=1, date=2024-02-05}
15. [1] TestOrder{id=O005, product_id=PROD001, amount=1, date=2024-02-10, priority=high, customer_id=P002, total=999.99, status=confirmed, discount=100, region=south}
16. [1] TestOrder{id=O006, discount=0, region=west, product_id=PROD005, amount=2, status=cancelled, priority=low, customer_id=P005, total=999.98, date=2024-02-15}
17. [1] TestOrder{id=O007, product_id=PROD006, date=2024-03-01, status=shipped, region=north, customer_id=P007, amount=4, total=600, priority=urgent, discount=50}
18. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, region=south, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0}
19. [1] TestOrder{id=O009, date=2024-03-10, discount=10, region=north, total=89.99, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1}
20. [1] TestOrder{id=O010, discount=0, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, region=east, total=75000}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales}
2. [1] TestPerson{id=P002, age=35, department=engineering, level=5, name=Bob, salary=75000, active=true, score=9.2, tags=senior, status=active}
3. [1] TestPerson{id=P003, name=Charlie, active=false, tags=intern, status=inactive, age=16, salary=0, score=6, department=hr, level=1}
4. [1] TestPerson{id=P004, age=45, active=true, score=7.8, tags=manager, department=marketing, level=7, name=Diana, salary=85000, status=active}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales, level=3, score=8, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, level=1, salary=-5000, status=active, department=qa}
7. [1] TestPerson{id=P007, score=10, tags=executive, department=management, level=9, name=Grace, status=active, age=65, salary=95000, active=true}
8. [1] TestPerson{id=P008, name=Henry, tags=junior, status=inactive, department=support, age=18, salary=25000, active=false, score=5.5, level=1}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering}
10. [1] TestPerson{id=P010, active=true, department=intern, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1}

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

1. [1] TestPerson{id=P001, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, level=5, name=Bob, salary=75000, active=true, score=9.2, tags=senior, status=active, age=35, department=engineering}
3. [1] TestPerson{id=P003, age=16, salary=0, score=6, department=hr, level=1, name=Charlie, active=false, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, department=marketing, level=7, name=Diana, salary=85000, status=active, age=45, active=true, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, tags=employee, department=sales, level=3, score=8, status=inactive, name=Eve, age=30, salary=55000, active=false}
6. [1] TestPerson{id=P006, salary=-5000, status=active, department=qa, name=Frank, age=0, active=true, score=0, tags=test, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, tags=executive, department=management, level=9, name=Grace, status=active}
8. [1] TestPerson{id=P008, name=Henry, tags=junior, status=inactive, department=support, age=18, salary=25000, active=false, score=5.5, level=1}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering}
10. [1] TestPerson{id=P010, active=true, department=intern, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1}
11. [1] TestOrder{id=O001, status=pending, discount=50, customer_id=P001, product_id=PROD001, amount=2, priority=normal, region=north, total=1999.98, date=2024-01-15}
12. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, priority=low, discount=0, region=south, total=25.5, date=2024-01-20, status=confirmed}
13. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, discount=15, region=north, total=225, status=shipped, priority=high}
14. [1] TestOrder{id=O004, discount=0, region=east, product_id=PROD004, amount=1, date=2024-02-05, customer_id=P004, total=299.99, status=delivered, priority=normal}
15. [1] TestOrder{id=O005, total=999.99, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, date=2024-02-10, priority=high, customer_id=P002}
16. [1] TestOrder{id=O006, customer_id=P005, total=999.98, date=2024-02-15, discount=0, region=west, product_id=PROD005, amount=2, status=cancelled, priority=low}
17. [1] TestOrder{id=O007, product_id=PROD006, date=2024-03-01, status=shipped, region=north, customer_id=P007, amount=4, total=600, priority=urgent, discount=50}
18. [1] TestOrder{id=O008, product_id=PROD002, region=south, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, customer_id=P010}
19. [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, discount=10, region=north}
20. [1] TestOrder{id=O010, total=75000, discount=0, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, region=east}

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

1. [1] TestPerson{id=P001, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, status=active, age=35, department=engineering, level=5}
3. [1] TestPerson{id=P003, department=hr, level=1, name=Charlie, active=false, tags=intern, status=inactive, age=16, salary=0, score=6}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, status=active, age=45, active=true, score=7.8, tags=manager, department=marketing, level=7}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales, level=3, score=8, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, level=1, salary=-5000, status=active, department=qa}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, tags=executive, department=management, level=9, name=Grace, status=active}
8. [1] TestPerson{id=P008, active=false, score=5.5, level=1, name=Henry, tags=junior, status=inactive, department=support, age=18, salary=25000}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true, department=intern}

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

1. [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, status=pending, discount=50, customer_id=P001, product_id=PROD001, amount=2, priority=normal, region=north}
2. [1] TestOrder{id=O002, amount=1, priority=low, discount=0, region=south, total=25.5, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002}
3. [1] TestOrder{id=O003, date=2024-02-01, discount=15, region=north, total=225, status=shipped, priority=high, customer_id=P001, product_id=PROD003, amount=3}
4. [1] TestOrder{id=O004, customer_id=P004, total=299.99, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, amount=1, date=2024-02-05}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, date=2024-02-10, priority=high, customer_id=P002, total=999.99, status=confirmed, discount=100, region=south}
6. [1] TestOrder{id=O006, amount=2, status=cancelled, priority=low, customer_id=P005, total=999.98, date=2024-02-15, discount=0, region=west, product_id=PROD005}
7. [1] TestOrder{id=O007, total=600, priority=urgent, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, region=north, customer_id=P007, amount=4}
8. [1] TestOrder{id=O008, priority=normal, discount=0, customer_id=P010, product_id=PROD002, region=south, amount=10, total=255, date=2024-03-05, status=pending}
9. [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, discount=10, region=north}
10. [1] TestOrder{id=O010, discount=0, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, region=east, total=75000}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, status=active, age=35, department=engineering, level=5, name=Bob, salary=75000, active=true}
3. [1] TestPerson{id=P003, score=6, department=hr, level=1, name=Charlie, active=false, tags=intern, status=inactive, age=16, salary=0}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, status=active, age=45, active=true, score=7.8, tags=manager, department=marketing, level=7}
5. [1] TestPerson{id=P005, tags=employee, department=sales, level=3, score=8, status=inactive, name=Eve, age=30, salary=55000, active=false}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, level=1, salary=-5000, status=active, department=qa}
7. [1] TestPerson{id=P007, department=management, level=9, name=Grace, status=active, age=65, salary=95000, active=true, score=10, tags=executive}
8. [1] TestPerson{id=P008, name=Henry, tags=junior, status=inactive, department=support, age=18, salary=25000, active=false, score=5.5, level=1}
9. [1] TestPerson{id=P009, salary=68000, score=8.7, department=engineering, name=Ivy, age=40, active=true, tags=senior, status=active, level=6}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true, department=intern, name=X}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 1 (5.3%)
- **Tokens générés**: 9
- **Faits traités**: 27
