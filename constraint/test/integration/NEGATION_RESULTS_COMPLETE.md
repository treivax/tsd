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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, name=Bob, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, age=16, active=false, tags=intern, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, active=false, status=inactive, department=sales, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, status=active, level=1, score=0, tags=test, department=qa}
7. [1] TestPerson{id=P007, name=Grace, tags=executive, level=9, age=65, salary=95000, active=true, score=10, status=active, department=management}
8. [1] TestPerson{id=P008, level=1, active=false, score=5.5, status=inactive, name=Henry, age=18, salary=25000, tags=junior, department=support}
9. [1] TestPerson{id=P009, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, salary=68000, active=true}
10. [1] TestPerson{id=P010, level=1, salary=28000, active=true, name=X, age=22, score=6.5, tags=temp, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, tags=intern, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, active=false, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, score=8}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, active=true, score=10, status=active, department=management, name=Grace, tags=executive, level=9, age=65}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, age=22, score=6.5, tags=temp, status=active, department=intern, level=1, salary=28000, active=true, name=X}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, name=Bob, status=active, department=engineering, level=5}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, department=support, level=1, active=false, score=5.5, status=inactive}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, salary=68000, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, customer_id=P001, product_id=PROD001, amount=2}
2. [1] TestOrder{id=O002, status=confirmed, discount=0, product_id=PROD002, amount=1, total=25.5, date=2024-01-20, priority=low, region=south, customer_id=P002}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, total=225, status=shipped, region=north, product_id=PROD003, date=2024-02-01, priority=high, discount=15}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, date=2024-02-05, status=delivered, region=east, total=299.99, priority=normal, discount=0}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, discount=100, region=south, total=999.99, date=2024-02-10, priority=high}
6. [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, customer_id=P005, product_id=PROD005, total=999.98, region=west, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, amount=4, date=2024-03-01, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600, status=shipped}
8. [1] TestOrder{id=O008, product_id=PROD002, total=255, status=pending, region=south, customer_id=P010, amount=10, date=2024-03-05, priority=normal, discount=0}
9. [1] TestOrder{id=O009, customer_id=P001, status=completed, priority=low, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, discount=10}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, total=75000, status=refunded, discount=0, region=east, customer_id=P006, date=2024-03-15, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, status=active, department=engineering, level=5, age=35, salary=75000, active=true, score=9.2, tags=senior}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7, name=Diana, salary=85000, active=true}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, active=false, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, score=8, tags=employee}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, tags=junior, department=support, level=1, active=false, score=5.5, status=inactive, name=Henry, age=18}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, score=6.5, tags=temp, status=active, department=intern, level=1, salary=28000, active=true}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, active=false, tags=intern, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6, age=16}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P006, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true, status=active, level=1}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, level=9, age=65, salary=95000, active=true, score=10, status=active, department=management}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, level=6, age=40, salary=68000, active=true, name=Ivy, score=8.7, tags=senior, status=active}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active}
2. [1] TestPerson{id=P002, status=active, department=engineering, level=5, age=35, salary=75000, active=true, score=9.2, tags=senior, name=Bob}
3. [1] TestPerson{id=P003, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6, age=16, active=false, tags=intern}
4. [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7, name=Diana}
5. [1] TestPerson{id=P005, tags=employee, active=false, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, score=8}
6. [1] TestPerson{id=P006, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true, status=active, level=1}
7. [1] TestPerson{id=P007, salary=95000, active=true, score=10, status=active, department=management, name=Grace, tags=executive, level=9, age=65}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, department=support, level=1, active=false, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, salary=68000, active=true, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40}
10. [1] TestPerson{id=P010, age=22, score=6.5, tags=temp, status=active, department=intern, level=1, salary=28000, active=true, name=X}

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

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, status=confirmed, discount=0, product_id=PROD002, amount=1, total=25.5, date=2024-01-20, priority=low, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, total=225, status=shipped, region=north, product_id=PROD003, date=2024-02-01, priority=high, discount=15}
4. [1] TestOrder{id=O004, date=2024-02-05, status=delivered, region=east, total=299.99, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1}
5. [1] TestOrder{id=O005, total=999.99, date=2024-02-10, priority=high, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, discount=100, region=south}
6. [1] TestOrder{id=O006, discount=0, customer_id=P005, product_id=PROD005, total=999.98, region=west, amount=2, date=2024-02-15, status=cancelled, priority=low}
7. [1] TestOrder{id=O007, total=600, status=shipped, amount=4, date=2024-03-01, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006}
8. [1] TestOrder{id=O008, region=south, customer_id=P010, amount=10, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, total=255, status=pending}
9. [1] TestOrder{id=O009, customer_id=P001, status=completed, priority=low, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, discount=10}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, total=75000, status=refunded, discount=0, region=east, customer_id=P006, date=2024-03-15, priority=urgent}

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

1. [1] TestProduct{id=PROD001, name=Laptop, price=999.99, available=true, rating=4.5, keywords=computer, stock=50, supplier=TechSupply, category=electronics, brand=TechCorp}
2. [1] TestProduct{id=PROD002, available=true, keywords=peripheral, brand=TechCorp, name=Mouse, category=accessories, price=25.5, rating=4, stock=200, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, price=75, rating=3.5, keywords=typing, brand=KeyTech, stock=0, name=Keyboard, category=accessories, available=false, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, name=Monitor, category=electronics, price=299.99, rating=4.8, brand=ScreenPro, available=true, keywords=display, stock=30, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, available=false, rating=2, brand=OldTech, stock=0, supplier=OldSupply, price=8.5, keywords=obsolete}
6. [1] TestProduct{id=PROD006, name=Headphones, available=true, rating=4.6, keywords=sound, brand=AudioMax, supplier=AudioSupply, category=audio, price=150, stock=75}
7. [1] TestProduct{id=PROD007, keywords=video, brand=CamTech, stock=25, supplier=CamSupply, name=Webcam, category=electronics, price=89.99, available=true, rating=3.8}

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

1. [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, name=Bob, status=active, department=engineering, level=5, age=35, salary=75000, active=true}
3. [1] TestPerson{id=P003, age=16, active=false, tags=intern, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6}
4. [1] TestPerson{id=P004, department=marketing, level=7, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, active=false, status=inactive, department=sales, level=3}
6. [1] TestPerson{id=P006, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true, status=active, level=1}
7. [1] TestPerson{id=P007, score=10, status=active, department=management, name=Grace, tags=executive, level=9, age=65, salary=95000, active=true}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, department=support, level=1, active=false, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, level=6, age=40, salary=68000, active=true, name=Ivy, score=8.7}
10. [1] TestPerson{id=P010, active=true, name=X, age=22, score=6.5, tags=temp, status=active, department=intern, level=1, salary=28000}

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

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, status=confirmed, discount=0, product_id=PROD002, amount=1, total=25.5, date=2024-01-20, priority=low, region=south}
3. [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, priority=high, discount=15, customer_id=P001, amount=3, total=225, status=shipped, region=north}
4. [1] TestOrder{id=O004, status=delivered, region=east, total=299.99, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1, date=2024-02-05}
5. [1] TestOrder{id=O005, status=confirmed, discount=100, region=south, total=999.99, date=2024-02-10, priority=high, customer_id=P002, product_id=PROD001, amount=1}
6. [1] TestOrder{id=O006, amount=2, date=2024-02-15, status=cancelled, priority=low, discount=0, customer_id=P005, product_id=PROD005, total=999.98, region=west}
7. [1] TestOrder{id=O007, status=shipped, amount=4, date=2024-03-01, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600}
8. [1] TestOrder{id=O008, product_id=PROD002, total=255, status=pending, region=south, customer_id=P010, amount=10, date=2024-03-05, priority=normal, discount=0}
9. [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, discount=10, customer_id=P001, status=completed, priority=low, region=north}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, total=75000, status=refunded, discount=0, region=east, customer_id=P006, date=2024-03-15, priority=urgent}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, name=Bob, status=active, department=engineering, level=5, age=35}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, score=6, age=16, active=false, tags=intern, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, status=active, age=45, department=marketing, level=7, name=Diana, salary=85000, active=true, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, active=false, status=inactive, department=sales, level=3}
6. [1] TestPerson{id=P006, active=true, status=active, level=1, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000}
7. [1] TestPerson{id=P007, active=true, score=10, status=active, department=management, name=Grace, tags=executive, level=9, age=65, salary=95000}
8. [1] TestPerson{id=P008, level=1, active=false, score=5.5, status=inactive, name=Henry, age=18, salary=25000, tags=junior, department=support}
9. [1] TestPerson{id=P009, age=40, salary=68000, active=true, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, name=X, age=22, score=6.5, tags=temp, status=active, department=intern, level=1, salary=28000, active=true}

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

1. [1] TestOrder{id=O001, status=pending, priority=normal, discount=50, region=north, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15}
2. [1] TestOrder{id=O002, customer_id=P002, status=confirmed, discount=0, product_id=PROD002, amount=1, total=25.5, date=2024-01-20, priority=low, region=south}
3. [1] TestOrder{id=O003, amount=3, total=225, status=shipped, region=north, product_id=PROD003, date=2024-02-01, priority=high, discount=15, customer_id=P001}
4. [1] TestOrder{id=O004, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1, date=2024-02-05, status=delivered, region=east, total=299.99}
5. [1] TestOrder{id=O005, total=999.99, date=2024-02-10, priority=high, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, discount=100, region=south}
6. [1] TestOrder{id=O006, amount=2, date=2024-02-15, status=cancelled, priority=low, discount=0, customer_id=P005, product_id=PROD005, total=999.98, region=west}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, status=shipped, amount=4, date=2024-03-01, priority=urgent, discount=50, region=north}
8. [1] TestOrder{id=O008, amount=10, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, total=255, status=pending, region=south, customer_id=P010}
9. [1] TestOrder{id=O009, customer_id=P001, status=completed, priority=low, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, discount=10}
10. [1] TestOrder{id=O010, customer_id=P006, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, total=75000, status=refunded, discount=0, region=east}

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

1. [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active}
2. [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, name=Bob, status=active, department=engineering, level=5, age=35, salary=75000}
3. [1] TestPerson{id=P003, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6, age=16, active=false, tags=intern}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7}
5. [1] TestPerson{id=P005, score=8, tags=employee, active=false, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000}
6. [1] TestPerson{id=P006, active=true, status=active, level=1, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, department=management, name=Grace, tags=executive, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, department=support, level=1, active=false, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, salary=68000, active=true}
10. [1] TestPerson{id=P010, name=X, age=22, score=6.5, tags=temp, status=active, department=intern, level=1, salary=28000, active=true}

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

1. [1] TestProduct{id=PROD001, category=electronics, brand=TechCorp, name=Laptop, price=999.99, available=true, rating=4.5, keywords=computer, stock=50, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, price=25.5, rating=4, stock=200, supplier=TechSupply, available=true, keywords=peripheral, brand=TechCorp}
3. [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, available=false, supplier=KeySupply, price=75, rating=3.5, keywords=typing, brand=KeyTech, stock=0}
4. [1] TestProduct{id=PROD004, name=Monitor, category=electronics, price=299.99, rating=4.8, brand=ScreenPro, available=true, keywords=display, stock=30, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, price=8.5, keywords=obsolete, name=OldKeyboard, category=accessories, available=false, rating=2, brand=OldTech, stock=0, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, name=Headphones, available=true, rating=4.6, keywords=sound, brand=AudioMax, supplier=AudioSupply, category=audio, price=150, stock=75}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, price=89.99, available=true, rating=3.8, keywords=video, brand=CamTech, stock=25, supplier=CamSupply}

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

1. [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, name=Bob, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6, age=16, active=false}
4. [1] TestPerson{id=P004, level=7, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, active=false, status=inactive, department=sales, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, status=active, level=1, score=0, tags=test, department=qa}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, department=management, name=Grace, tags=executive, level=9}
8. [1] TestPerson{id=P008, active=false, score=5.5, status=inactive, name=Henry, age=18, salary=25000, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, salary=68000, active=true, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40}
10. [1] TestPerson{id=P010, salary=28000, active=true, name=X, age=22, score=6.5, tags=temp, status=active, department=intern, level=1}

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

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}
2. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, priority=low, region=south, customer_id=P002, status=confirmed, discount=0, product_id=PROD002, amount=1}
3. [1] TestOrder{id=O003, priority=high, discount=15, customer_id=P001, amount=3, total=225, status=shipped, region=north, product_id=PROD003, date=2024-02-01}
4. [1] TestOrder{id=O004, total=299.99, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1, date=2024-02-05, status=delivered, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, discount=100, region=south, total=999.99, date=2024-02-10, priority=high}
6. [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, customer_id=P005, product_id=PROD005, total=999.98, region=west, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, status=shipped, amount=4, date=2024-03-01, priority=urgent, discount=50, region=north}
8. [1] TestOrder{id=O008, product_id=PROD002, total=255, status=pending, region=south, customer_id=P010, amount=10, date=2024-03-05, priority=normal, discount=0}
9. [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, discount=10, customer_id=P001, status=completed, priority=low, region=north}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, total=75000, status=refunded, discount=0, region=east, customer_id=P006, date=2024-03-15, priority=urgent}

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

1. [1] TestPerson{id=P001, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, name=Bob, status=active, department=engineering, level=5, age=35, salary=75000, active=true}
3. [1] TestPerson{id=P003, department=hr, level=1, name=Charlie, salary=0, score=6, age=16, active=false, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7}
5. [1] TestPerson{id=P005, level=3, name=Eve, age=30, salary=55000, score=8, tags=employee, active=false, status=inactive, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, status=active, level=1, score=0, tags=test, department=qa}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, department=management, name=Grace, tags=executive, level=9}
8. [1] TestPerson{id=P008, salary=25000, tags=junior, department=support, level=1, active=false, score=5.5, status=inactive, name=Henry, age=18}
9. [1] TestPerson{id=P009, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, salary=68000, active=true}
10. [1] TestPerson{id=P010, active=true, name=X, age=22, score=6.5, tags=temp, status=active, department=intern, level=1, salary=28000}
11. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}
12. [1] TestOrder{id=O002, customer_id=P002, status=confirmed, discount=0, product_id=PROD002, amount=1, total=25.5, date=2024-01-20, priority=low, region=south}
13. [1] TestOrder{id=O003, customer_id=P001, amount=3, total=225, status=shipped, region=north, product_id=PROD003, date=2024-02-01, priority=high, discount=15}
14. [1] TestOrder{id=O004, status=delivered, region=east, total=299.99, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1, date=2024-02-05}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, discount=100, region=south, total=999.99, date=2024-02-10, priority=high}
16. [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, customer_id=P005, product_id=PROD005, total=999.98, region=west, amount=2, date=2024-02-15}
17. [1] TestOrder{id=O007, total=600, status=shipped, amount=4, date=2024-03-01, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006}
18. [1] TestOrder{id=O008, product_id=PROD002, total=255, status=pending, region=south, customer_id=P010, amount=10, date=2024-03-05, priority=normal, discount=0}
19. [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, discount=10, customer_id=P001, status=completed, priority=low, region=north}
20. [1] TestOrder{id=O010, priority=urgent, product_id=PROD001, amount=1, total=75000, status=refunded, discount=0, region=east, customer_id=P006, date=2024-03-15}

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

1. [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, name=Bob, status=active, department=engineering, level=5, age=35}
3. [1] TestPerson{id=P003, age=16, active=false, tags=intern, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, active=false, status=inactive, department=sales, level=3}
6. [1] TestPerson{id=P006, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true, status=active, level=1}
7. [1] TestPerson{id=P007, score=10, status=active, department=management, name=Grace, tags=executive, level=9, age=65, salary=95000, active=true}
8. [1] TestPerson{id=P008, status=inactive, name=Henry, age=18, salary=25000, tags=junior, department=support, level=1, active=false, score=5.5}
9. [1] TestPerson{id=P009, active=true, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, salary=68000}
10. [1] TestPerson{id=P010, tags=temp, status=active, department=intern, level=1, salary=28000, active=true, name=X, age=22, score=6.5}

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

1. [1] TestPerson{id=P001, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5}
2. [1] TestPerson{id=P002, department=engineering, level=5, age=35, salary=75000, active=true, score=9.2, tags=senior, name=Bob, status=active}
3. [1] TestPerson{id=P003, department=hr, level=1, name=Charlie, salary=0, score=6, age=16, active=false, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7}
5. [1] TestPerson{id=P005, level=3, name=Eve, age=30, salary=55000, score=8, tags=employee, active=false, status=inactive, department=sales}
6. [1] TestPerson{id=P006, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true, status=active, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, department=management, name=Grace, tags=executive, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, department=support, level=1, active=false, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, salary=68000, active=true, name=Ivy}
10. [1] TestPerson{id=P010, score=6.5, tags=temp, status=active, department=intern, level=1, salary=28000, active=true, name=X, age=22}
11. [1] TestOrder{id=O001, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, customer_id=P001, product_id=PROD001, amount=2, total=1999.98}
12. [1] TestOrder{id=O002, customer_id=P002, status=confirmed, discount=0, product_id=PROD002, amount=1, total=25.5, date=2024-01-20, priority=low, region=south}
13. [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, priority=high, discount=15, customer_id=P001, amount=3, total=225, status=shipped, region=north}
14. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, date=2024-02-05, status=delivered, region=east, total=299.99, priority=normal, discount=0}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, discount=100, region=south, total=999.99, date=2024-02-10, priority=high}
16. [1] TestOrder{id=O006, discount=0, customer_id=P005, product_id=PROD005, total=999.98, region=west, amount=2, date=2024-02-15, status=cancelled, priority=low}
17. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, status=shipped, amount=4, date=2024-03-01, priority=urgent, discount=50, region=north}
18. [1] TestOrder{id=O008, customer_id=P010, amount=10, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, total=255, status=pending, region=south}
19. [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, discount=10, customer_id=P001, status=completed, priority=low, region=north}
20. [1] TestOrder{id=O010, total=75000, status=refunded, discount=0, region=east, customer_id=P006, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1}

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

1. [1] TestPerson{id=P001, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, name=Bob, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, age=16, active=false, tags=intern, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7}
5. [1] TestPerson{id=P005, active=false, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, score=8, tags=employee}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, status=active, level=1, score=0, tags=test, department=qa}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, department=management, name=Grace, tags=executive, level=9}
8. [1] TestPerson{id=P008, department=support, level=1, active=false, score=5.5, status=inactive, name=Henry, age=18, salary=25000, tags=junior}
9. [1] TestPerson{id=P009, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, salary=68000, active=true}
10. [1] TestPerson{id=P010, tags=temp, status=active, department=intern, level=1, salary=28000, active=true, name=X, age=22, score=6.5}

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

1. [1] TestOrder{id=O001, discount=50, region=north, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, status=pending, priority=normal}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, date=2024-01-20, priority=low, region=south, customer_id=P002, status=confirmed, discount=0}
3. [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, priority=high, discount=15, customer_id=P001, amount=3, total=225, status=shipped, region=north}
4. [1] TestOrder{id=O004, date=2024-02-05, status=delivered, region=east, total=299.99, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1}
5. [1] TestOrder{id=O005, total=999.99, date=2024-02-10, priority=high, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, discount=100, region=south}
6. [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, customer_id=P005, product_id=PROD005, total=999.98, region=west, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, status=shipped, amount=4, date=2024-03-01, priority=urgent, discount=50, region=north}
8. [1] TestOrder{id=O008, region=south, customer_id=P010, amount=10, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, total=255, status=pending}
9. [1] TestOrder{id=O009, priority=low, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, discount=10, customer_id=P001, status=completed}
10. [1] TestOrder{id=O010, region=east, customer_id=P006, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, total=75000, status=refunded, discount=0}

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

1. [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active}
2. [1] TestPerson{id=P002, name=Bob, status=active, department=engineering, level=5, age=35, salary=75000, active=true, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, age=16, active=false, tags=intern, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7}
5. [1] TestPerson{id=P005, active=false, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, score=8, tags=employee}
6. [1] TestPerson{id=P006, status=active, level=1, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true}
7. [1] TestPerson{id=P007, name=Grace, tags=executive, level=9, age=65, salary=95000, active=true, score=10, status=active, department=management}
8. [1] TestPerson{id=P008, department=support, level=1, active=false, score=5.5, status=inactive, name=Henry, age=18, salary=25000, tags=junior}
9. [1] TestPerson{id=P009, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, salary=68000, active=true}
10. [1] TestPerson{id=P010, name=X, age=22, score=6.5, tags=temp, status=active, department=intern, level=1, salary=28000, active=true}

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
- **Tokens générés**: 19
- **Faits traités**: 27
