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

1. [1] TestPerson{id=P001, score=8.5, age=25, salary=45000, tags=junior, status=active, department=sales, level=2, name=Alice, active=true}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true, status=active}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, score=6, department=hr, level=1, active=false, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, age=45, active=true, status=active, salary=85000, score=7.8, tags=manager, department=marketing, level=7, name=Diana}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, level=3, age=30, active=false, score=8, status=inactive, department=sales}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, tags=test, age=0, salary=-5000, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, salary=95000, score=10, tags=executive, status=active, name=Grace, active=true, department=management, level=9, age=65}
8. [1] TestPerson{id=P008, department=support, level=1, name=Henry, age=18, score=5.5, salary=25000, active=false, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, age=40, active=true, department=engineering, level=6}
10. [1] TestPerson{id=P010, age=22, score=6.5, tags=temp, status=active, level=1, name=X, salary=28000, active=true, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, level=5, name=Bob, age=35, salary=75000, active=true, status=active, score=9.2, tags=senior, department=engineering}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, active=false, tags=intern, status=inactive, name=Charlie, age=16, salary=0, score=6, department=hr, level=1}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, score=7.8, tags=manager, department=marketing, level=7, name=Diana, age=45, active=true, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, tags=employee, level=3, age=30, active=false, score=8, status=inactive, department=sales, name=Eve}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, active=false, tags=junior, status=inactive, department=support, level=1, name=Henry, age=18, score=5.5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, department=intern, age=22, score=6.5, tags=temp, status=active, level=1}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, level=2, name=Alice, active=true, score=8.5, age=25, salary=45000, tags=junior, status=active}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, score=10, tags=executive, status=active, name=Grace, active=true, department=management, level=9}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, region=north, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, priority=normal, discount=50, amount=2, date=2024-01-15}
2. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, discount=0, region=south, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, priority=low}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, discount=15, amount=3, date=2024-02-01, status=shipped, priority=high, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, product_id=PROD004, discount=0, region=east}
5. [1] TestOrder{id=O005, amount=1, total=999.99, status=confirmed, priority=high, discount=100, date=2024-02-10, region=south, customer_id=P002, product_id=PROD001}
6. [1] TestOrder{id=O006, product_id=PROD005, status=cancelled, region=west, customer_id=P005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0}
7. [1] TestOrder{id=O007, date=2024-03-01, status=shipped, discount=50, product_id=PROD006, amount=4, total=600, priority=urgent, region=north, customer_id=P007}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255}
9. [1] TestOrder{id=O009, discount=10, customer_id=P001, amount=1, total=89.99, region=north, product_id=PROD007, date=2024-03-10, status=completed, priority=low}
10. [1] TestOrder{id=O010, date=2024-03-15, customer_id=P006, amount=1, status=refunded, priority=urgent, discount=0, region=east, product_id=PROD001, total=75000}

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

1. [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, age=25, salary=45000, tags=junior, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true, status=active}
3. [1] TestPerson{id=P003, tags=intern, status=inactive, name=Charlie, age=16, salary=0, score=6, department=hr, level=1, active=false}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, status=active, salary=85000, score=7.8, tags=manager, department=marketing, level=7}
5. [1] TestPerson{id=P005, active=false, score=8, status=inactive, department=sales, name=Eve, salary=55000, tags=employee, level=3, age=30}
6. [1] TestPerson{id=P006, age=0, salary=-5000, status=active, department=qa, level=1, name=Frank, active=true, score=0, tags=test}
7. [1] TestPerson{id=P007, name=Grace, active=true, department=management, level=9, age=65, salary=95000, score=10, tags=executive, status=active}
8. [1] TestPerson{id=P008, age=18, score=5.5, salary=25000, active=false, tags=junior, status=inactive, department=support, level=1, name=Henry}
9. [1] TestPerson{id=P009, age=40, active=true, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior, status=active}
10. [1] TestPerson{id=P010, active=true, department=intern, age=22, score=6.5, tags=temp, status=active, level=1, name=X, salary=28000}

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

1. [1] TestOrder{id=O001, date=2024-01-15, region=north, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, priority=normal, discount=50, amount=2}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, priority=low, total=25.5, date=2024-01-20, discount=0, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, discount=15, amount=3, date=2024-02-01, status=shipped, priority=high, region=north}
4. [1] TestOrder{id=O004, date=2024-02-05, status=delivered, priority=normal, product_id=PROD004, discount=0, region=east, customer_id=P004, amount=1, total=299.99}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, priority=high, discount=100, date=2024-02-10, region=south}
6. [1] TestOrder{id=O006, customer_id=P005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, product_id=PROD005, status=cancelled, region=west}
7. [1] TestOrder{id=O007, total=600, priority=urgent, region=north, customer_id=P007, date=2024-03-01, status=shipped, discount=50, product_id=PROD006, amount=4}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, region=north, product_id=PROD007, date=2024-03-10, status=completed, priority=low, discount=10}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, status=refunded, priority=urgent, discount=0, region=east, product_id=PROD001, total=75000, date=2024-03-15}

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

1. [1] TestProduct{id=PROD001, available=true, supplier=TechSupply, price=999.99, rating=4.5, keywords=computer, brand=TechCorp, stock=50, name=Laptop, category=electronics}
2. [1] TestProduct{id=PROD002, category=accessories, price=25.5, available=true, supplier=TechSupply, name=Mouse, rating=4, keywords=peripheral, brand=TechCorp, stock=200}
3. [1] TestProduct{id=PROD003, category=accessories, rating=3.5, keywords=typing, name=Keyboard, price=75, available=false, brand=KeyTech, stock=0, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, category=electronics, rating=4.8, keywords=display, brand=ScreenPro, name=Monitor, price=299.99, available=true, stock=30, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, price=8.5, available=false, rating=2, keywords=obsolete, supplier=OldSupply, brand=OldTech, stock=0}
6. [1] TestProduct{id=PROD006, price=150, available=true, rating=4.6, keywords=sound, brand=AudioMax, stock=75, category=audio, supplier=AudioSupply, name=Headphones}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, available=true, keywords=video, price=89.99, rating=3.8, brand=CamTech, stock=25, supplier=CamSupply}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, price=299.99, available=true, stock=30, supplier=ScreenSupply, category=electronics, rating=4.8, keywords=display, brand=ScreenPro}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD006, name=Headphones, price=150, available=true, rating=4.6, keywords=sound, brand=AudioMax, stock=75, category=audio, supplier=AudioSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD007, name=Webcam, category=electronics, available=true, keywords=video, price=89.99, rating=3.8, brand=CamTech, stock=25, supplier=CamSupply}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD001, price=999.99, rating=4.5, keywords=computer, brand=TechCorp, stock=50, name=Laptop, category=electronics, available=true, supplier=TechSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD002, category=accessories, price=25.5, available=true, supplier=TechSupply, name=Mouse, rating=4, keywords=peripheral, brand=TechCorp, stock=200}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD003, category=accessories, rating=3.5, keywords=typing, name=Keyboard, price=75, available=false, brand=KeyTech, stock=0, supplier=KeySupply}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, age=25, salary=45000, tags=junior, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, level=5, name=Bob, age=35, salary=75000, active=true, status=active, score=9.2, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, department=hr, level=1, active=false, tags=intern, status=inactive, name=Charlie, age=16, salary=0, score=6}
4. [1] TestPerson{id=P004, active=true, status=active, salary=85000, score=7.8, tags=manager, department=marketing, level=7, name=Diana, age=45}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, level=3, age=30, active=false, score=8, status=inactive, department=sales}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, tags=test, age=0, salary=-5000, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, active=true, department=management, level=9, age=65, salary=95000, score=10, tags=executive, status=active, name=Grace}
8. [1] TestPerson{id=P008, age=18, score=5.5, salary=25000, active=false, tags=junior, status=inactive, department=support, level=1, name=Henry}
9. [1] TestPerson{id=P009, age=40, active=true, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior, status=active}
10. [1] TestPerson{id=P010, age=22, score=6.5, tags=temp, status=active, level=1, name=X, salary=28000, active=true, department=intern}

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

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, priority=normal, discount=50, amount=2, date=2024-01-15, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, priority=low, total=25.5, date=2024-01-20, discount=0, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, discount=15, amount=3, date=2024-02-01, status=shipped, priority=high, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, product_id=PROD004, discount=0, region=east}
5. [1] TestOrder{id=O005, discount=100, date=2024-02-10, region=south, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, priority=high}
6. [1] TestOrder{id=O006, customer_id=P005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, product_id=PROD005, status=cancelled, region=west}
7. [1] TestOrder{id=O007, customer_id=P007, date=2024-03-01, status=shipped, discount=50, product_id=PROD006, amount=4, total=600, priority=urgent, region=north}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255}
9. [1] TestOrder{id=O009, status=completed, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, region=north, product_id=PROD007, date=2024-03-10}
10. [1] TestOrder{id=O010, product_id=PROD001, total=75000, date=2024-03-15, customer_id=P006, amount=1, status=refunded, priority=urgent, discount=0, region=east}

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

1. [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, age=25, salary=45000, tags=junior, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, status=active, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, score=6, department=hr, level=1, active=false, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, active=true, status=active, salary=85000, score=7.8, tags=manager, department=marketing, level=7, name=Diana, age=45}
5. [1] TestPerson{id=P005, status=inactive, department=sales, name=Eve, salary=55000, tags=employee, level=3, age=30, active=false, score=8}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, tags=test, age=0, salary=-5000, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, salary=95000, score=10, tags=executive, status=active, name=Grace, active=true, department=management, level=9, age=65}
8. [1] TestPerson{id=P008, name=Henry, age=18, score=5.5, salary=25000, active=false, tags=junior, status=inactive, department=support, level=1}
9. [1] TestPerson{id=P009, status=active, age=40, active=true, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, department=intern, age=22, score=6.5, tags=temp, status=active, level=1, name=X, salary=28000, active=true}

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

1. [1] TestOrder{id=O001, status=pending, priority=normal, discount=50, amount=2, date=2024-01-15, region=north, customer_id=P001, product_id=PROD001, total=1999.98}
2. [1] TestOrder{id=O002, amount=1, status=confirmed, priority=low, total=25.5, date=2024-01-20, discount=0, region=south, customer_id=P002, product_id=PROD002}
3. [1] TestOrder{id=O003, total=225, discount=15, amount=3, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, product_id=PROD003}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, product_id=PROD004, discount=0, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, priority=high, discount=100, date=2024-02-10, region=south}
6. [1] TestOrder{id=O006, customer_id=P005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, product_id=PROD005, status=cancelled, region=west}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, priority=urgent, region=north, customer_id=P007, date=2024-03-01, status=shipped, discount=50}
8. [1] TestOrder{id=O008, discount=0, region=south, amount=10, total=255, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, region=north, product_id=PROD007, date=2024-03-10, status=completed, priority=low, discount=10}
10. [1] TestOrder{id=O010, priority=urgent, discount=0, region=east, product_id=PROD001, total=75000, date=2024-03-15, customer_id=P006, amount=1, status=refunded}

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

1. [1] TestPerson{id=P001, salary=45000, tags=junior, status=active, department=sales, level=2, name=Alice, active=true, score=8.5, age=25}
2. [1] TestPerson{id=P002, tags=senior, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true, status=active, score=9.2}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, score=6, department=hr, level=1, active=false, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, status=active, salary=85000, score=7.8, tags=manager, department=marketing, level=7}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, level=3, age=30, active=false, score=8, status=inactive, department=sales}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, tags=test, age=0, salary=-5000, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, active=true, department=management, level=9, age=65, salary=95000, score=10, tags=executive, status=active}
8. [1] TestPerson{id=P008, salary=25000, active=false, tags=junior, status=inactive, department=support, level=1, name=Henry, age=18, score=5.5}
9. [1] TestPerson{id=P009, status=active, age=40, active=true, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, age=22, score=6.5, tags=temp, status=active, level=1, name=X, salary=28000, active=true, department=intern}

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

1. [1] TestProduct{id=PROD001, price=999.99, rating=4.5, keywords=computer, brand=TechCorp, stock=50, name=Laptop, category=electronics, available=true, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, category=accessories, price=25.5, available=true, supplier=TechSupply, name=Mouse, rating=4, keywords=peripheral, brand=TechCorp, stock=200}
3. [1] TestProduct{id=PROD003, name=Keyboard, price=75, available=false, brand=KeyTech, stock=0, supplier=KeySupply, category=accessories, rating=3.5, keywords=typing}
4. [1] TestProduct{id=PROD004, name=Monitor, price=299.99, available=true, stock=30, supplier=ScreenSupply, category=electronics, rating=4.8, keywords=display, brand=ScreenPro}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, price=8.5, available=false, rating=2, keywords=obsolete, supplier=OldSupply, brand=OldTech, stock=0}
6. [1] TestProduct{id=PROD006, name=Headphones, price=150, available=true, rating=4.6, keywords=sound, brand=AudioMax, stock=75, category=audio, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, rating=3.8, brand=CamTech, stock=25, supplier=CamSupply, name=Webcam, category=electronics, available=true, keywords=video, price=89.99}

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

1. [1] TestPerson{id=P001, salary=45000, tags=junior, status=active, department=sales, level=2, name=Alice, active=true, score=8.5, age=25}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true, status=active}
3. [1] TestPerson{id=P003, active=false, tags=intern, status=inactive, name=Charlie, age=16, salary=0, score=6, department=hr, level=1}
4. [1] TestPerson{id=P004, salary=85000, score=7.8, tags=manager, department=marketing, level=7, name=Diana, age=45, active=true, status=active}
5. [1] TestPerson{id=P005, age=30, active=false, score=8, status=inactive, department=sales, name=Eve, salary=55000, tags=employee, level=3}
6. [1] TestPerson{id=P006, score=0, tags=test, age=0, salary=-5000, status=active, department=qa, level=1, name=Frank, active=true}
7. [1] TestPerson{id=P007, name=Grace, active=true, department=management, level=9, age=65, salary=95000, score=10, tags=executive, status=active}
8. [1] TestPerson{id=P008, salary=25000, active=false, tags=junior, status=inactive, department=support, level=1, name=Henry, age=18, score=5.5}
9. [1] TestPerson{id=P009, age=40, active=true, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior, status=active}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, department=intern, age=22, score=6.5, tags=temp, status=active, level=1}

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

1. [1] TestOrder{id=O001, region=north, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, priority=normal, discount=50, amount=2, date=2024-01-15}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, priority=low, total=25.5, date=2024-01-20, discount=0, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, discount=15, amount=3, date=2024-02-01, status=shipped, priority=high, region=north}
4. [1] TestOrder{id=O004, discount=0, region=east, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, product_id=PROD004}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, priority=high, discount=100, date=2024-02-10, region=south}
6. [1] TestOrder{id=O006, product_id=PROD005, status=cancelled, region=west, customer_id=P005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, priority=urgent, region=north, customer_id=P007, date=2024-03-01, status=shipped, discount=50}
8. [1] TestOrder{id=O008, amount=10, total=255, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south}
9. [1] TestOrder{id=O009, product_id=PROD007, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, region=north}
10. [1] TestOrder{id=O010, priority=urgent, discount=0, region=east, product_id=PROD001, total=75000, date=2024-03-15, customer_id=P006, amount=1, status=refunded}

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

1. [1] TestPerson{id=P001, level=2, name=Alice, active=true, score=8.5, age=25, salary=45000, tags=junior, status=active, department=sales}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true, status=active}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, score=6, department=hr, level=1, active=false, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, salary=85000, score=7.8, tags=manager, department=marketing, level=7, name=Diana, age=45, active=true, status=active}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, level=3, age=30, active=false, score=8, status=inactive, department=sales}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, tags=test, age=0, salary=-5000, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, level=9, age=65, salary=95000, score=10, tags=executive, status=active, name=Grace, active=true, department=management}
8. [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, level=1, name=Henry, age=18, score=5.5, salary=25000, active=false}
9. [1] TestPerson{id=P009, age=40, active=true, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior, status=active}
10. [1] TestPerson{id=P010, salary=28000, active=true, department=intern, age=22, score=6.5, tags=temp, status=active, level=1, name=X}
11. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, priority=normal, discount=50, amount=2, date=2024-01-15, region=north}
12. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, discount=0, region=south, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, priority=low}
13. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, discount=15, amount=3, date=2024-02-01, status=shipped, priority=high, region=north}
14. [1] TestOrder{id=O004, date=2024-02-05, status=delivered, priority=normal, product_id=PROD004, discount=0, region=east, customer_id=P004, amount=1, total=299.99}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, priority=high, discount=100, date=2024-02-10, region=south}
16. [1] TestOrder{id=O006, customer_id=P005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, product_id=PROD005, status=cancelled, region=west}
17. [1] TestOrder{id=O007, status=shipped, discount=50, product_id=PROD006, amount=4, total=600, priority=urgent, region=north, customer_id=P007, date=2024-03-01}
18. [1] TestOrder{id=O008, status=pending, priority=normal, discount=0, region=south, amount=10, total=255, customer_id=P010, product_id=PROD002, date=2024-03-05}
19. [1] TestOrder{id=O009, product_id=PROD007, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, region=north}
20. [1] TestOrder{id=O010, total=75000, date=2024-03-15, customer_id=P006, amount=1, status=refunded, priority=urgent, discount=0, region=east, product_id=PROD001}

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

1. [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, age=25, salary=45000, tags=junior, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true, status=active}
3. [1] TestPerson{id=P003, active=false, tags=intern, status=inactive, name=Charlie, age=16, salary=0, score=6, department=hr, level=1}
4. [1] TestPerson{id=P004, salary=85000, score=7.8, tags=manager, department=marketing, level=7, name=Diana, age=45, active=true, status=active}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, level=3, age=30, active=false, score=8, status=inactive, department=sales}
6. [1] TestPerson{id=P006, age=0, salary=-5000, status=active, department=qa, level=1, name=Frank, active=true, score=0, tags=test}
7. [1] TestPerson{id=P007, active=true, department=management, level=9, age=65, salary=95000, score=10, tags=executive, status=active, name=Grace}
8. [1] TestPerson{id=P008, name=Henry, age=18, score=5.5, salary=25000, active=false, tags=junior, status=inactive, department=support, level=1}
9. [1] TestPerson{id=P009, age=40, active=true, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior, status=active}
10. [1] TestPerson{id=P010, age=22, score=6.5, tags=temp, status=active, level=1, name=X, salary=28000, active=true, department=intern}

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

1. [1] TestPerson{id=P001, department=sales, level=2, name=Alice, active=true, score=8.5, age=25, salary=45000, tags=junior, status=active}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, status=active, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, score=6, department=hr, level=1, active=false, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, score=7.8, tags=manager, department=marketing, level=7, name=Diana, age=45, active=true, status=active, salary=85000}
5. [1] TestPerson{id=P005, active=false, score=8, status=inactive, department=sales, name=Eve, salary=55000, tags=employee, level=3, age=30}
6. [1] TestPerson{id=P006, status=active, department=qa, level=1, name=Frank, active=true, score=0, tags=test, age=0, salary=-5000}
7. [1] TestPerson{id=P007, name=Grace, active=true, department=management, level=9, age=65, salary=95000, score=10, tags=executive, status=active}
8. [1] TestPerson{id=P008, salary=25000, active=false, tags=junior, status=inactive, department=support, level=1, name=Henry, age=18, score=5.5}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, age=40, active=true, department=engineering, level=6}
10. [1] TestPerson{id=P010, age=22, score=6.5, tags=temp, status=active, level=1, name=X, salary=28000, active=true, department=intern}
11. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, priority=normal, discount=50, amount=2, date=2024-01-15, region=north}
12. [1] TestOrder{id=O002, priority=low, total=25.5, date=2024-01-20, discount=0, region=south, customer_id=P002, product_id=PROD002, amount=1, status=confirmed}
13. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, discount=15, amount=3, date=2024-02-01, status=shipped, priority=high, region=north}
14. [1] TestOrder{id=O004, region=east, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, product_id=PROD004, discount=0}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, priority=high, discount=100, date=2024-02-10, region=south}
16. [1] TestOrder{id=O006, customer_id=P005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, product_id=PROD005, status=cancelled, region=west}
17. [1] TestOrder{id=O007, discount=50, product_id=PROD006, amount=4, total=600, priority=urgent, region=north, customer_id=P007, date=2024-03-01, status=shipped}
18. [1] TestOrder{id=O008, priority=normal, discount=0, region=south, amount=10, total=255, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending}
19. [1] TestOrder{id=O009, total=89.99, region=north, product_id=PROD007, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, amount=1}
20. [1] TestOrder{id=O010, discount=0, region=east, product_id=PROD001, total=75000, date=2024-03-15, customer_id=P006, amount=1, status=refunded, priority=urgent}

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

1. [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, status=active, department=sales, level=2, name=Alice, active=true, score=8.5}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, status=active, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, level=1, active=false, tags=intern, status=inactive, name=Charlie, age=16, salary=0, score=6, department=hr}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, status=active, salary=85000, score=7.8, tags=manager, department=marketing, level=7}
5. [1] TestPerson{id=P005, active=false, score=8, status=inactive, department=sales, name=Eve, salary=55000, tags=employee, level=3, age=30}
6. [1] TestPerson{id=P006, age=0, salary=-5000, status=active, department=qa, level=1, name=Frank, active=true, score=0, tags=test}
7. [1] TestPerson{id=P007, active=true, department=management, level=9, age=65, salary=95000, score=10, tags=executive, status=active, name=Grace}
8. [1] TestPerson{id=P008, status=inactive, department=support, level=1, name=Henry, age=18, score=5.5, salary=25000, active=false, tags=junior}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, age=40, active=true, department=engineering, level=6}
10. [1] TestPerson{id=P010, age=22, score=6.5, tags=temp, status=active, level=1, name=X, salary=28000, active=true, department=intern}

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

1. [1] TestOrder{id=O001, amount=2, date=2024-01-15, region=north, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, priority=normal, discount=50}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, priority=low, total=25.5, date=2024-01-20, discount=0, region=south}
3. [1] TestOrder{id=O003, amount=3, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, product_id=PROD003, total=225, discount=15}
4. [1] TestOrder{id=O004, status=delivered, priority=normal, product_id=PROD004, discount=0, region=east, customer_id=P004, amount=1, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, priority=high, discount=100, date=2024-02-10, region=south}
6. [1] TestOrder{id=O006, product_id=PROD005, status=cancelled, region=west, customer_id=P005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0}
7. [1] TestOrder{id=O007, customer_id=P007, date=2024-03-01, status=shipped, discount=50, product_id=PROD006, amount=4, total=600, priority=urgent, region=north}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255}
9. [1] TestOrder{id=O009, discount=10, customer_id=P001, amount=1, total=89.99, region=north, product_id=PROD007, date=2024-03-10, status=completed, priority=low}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, status=refunded, priority=urgent, discount=0, region=east, product_id=PROD001, total=75000, date=2024-03-15}

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

1. [1] TestPerson{id=P001, tags=junior, status=active, department=sales, level=2, name=Alice, active=true, score=8.5, age=25, salary=45000}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, status=active, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, level=1, active=false, tags=intern, status=inactive, name=Charlie, age=16, salary=0, score=6, department=hr}
4. [1] TestPerson{id=P004, salary=85000, score=7.8, tags=manager, department=marketing, level=7, name=Diana, age=45, active=true, status=active}
5. [1] TestPerson{id=P005, age=30, active=false, score=8, status=inactive, department=sales, name=Eve, salary=55000, tags=employee, level=3}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, tags=test, age=0, salary=-5000, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, score=10, tags=executive, status=active, name=Grace, active=true, department=management, level=9}
8. [1] TestPerson{id=P008, age=18, score=5.5, salary=25000, active=false, tags=junior, status=inactive, department=support, level=1, name=Henry}
9. [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, age=40, active=true, department=engineering, level=6, name=Ivy, salary=68000}
10. [1] TestPerson{id=P010, department=intern, age=22, score=6.5, tags=temp, status=active, level=1, name=X, salary=28000, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 2 (10.5%)
- **Tokens générés**: 15
- **Faits traités**: 27
