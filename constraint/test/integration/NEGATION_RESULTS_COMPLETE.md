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

1. [1] TestPerson{id=P001, status=active, level=2, name=Alice, active=true, tags=junior, department=sales, age=25, salary=45000, score=8.5}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
3. [1] TestPerson{id=P003, tags=intern, status=inactive, level=1, name=Charlie, salary=0, active=false, score=6, department=hr, age=16}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, department=marketing, level=7, age=45, active=true}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, score=8, tags=employee, status=inactive, level=3, age=30, department=sales}
6. [1] TestPerson{id=P006, age=0, salary=-5000, active=true, tags=test, name=Frank, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, tags=executive, department=management, name=Grace, status=active, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, tags=junior, department=support, level=1, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, tags=senior, name=Ivy, salary=68000, score=8.7, status=active, department=engineering, level=6, age=40, active=true}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, age=22, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, tags=executive, department=management, name=Grace, status=active, level=9}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, tags=senior, name=Ivy, salary=68000, score=8.7, status=active, department=engineering, level=6}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, age=25, salary=45000, score=8.5, status=active, level=2, name=Alice, active=true, tags=junior}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, department=hr, age=16, tags=intern, status=inactive, level=1}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, age=30, department=sales, name=Eve, salary=55000, active=false, score=8, tags=employee, status=inactive, level=3}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, tags=junior, department=support, level=1, score=5.5, status=inactive}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, age=22, status=active, department=intern, level=1, name=X, salary=28000, active=true, score=6.5, tags=temp}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P002, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, age=45, active=true, name=Diana, salary=85000, score=7.8, tags=manager, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, status=pending, priority=normal, region=north, customer_id=P001}
2. [1] TestOrder{id=O002, amount=1, date=2024-01-20, priority=low, product_id=PROD002, total=25.5, status=confirmed, discount=0, region=south, customer_id=P002}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north}
4. [1] TestOrder{id=O004, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed}
6. [1] TestOrder{id=O006, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled}
7. [1] TestOrder{id=O007, date=2024-03-01, priority=urgent, region=north, customer_id=P007, amount=4, total=600, status=shipped, discount=50, product_id=PROD006}
8. [1] TestOrder{id=O008, total=255, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10}
9. [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, region=north, customer_id=P001, product_id=PROD007, discount=10}
10. [1] TestOrder{id=O010, region=east, date=2024-03-15, status=refunded, discount=0, customer_id=P006, product_id=PROD001, amount=1, total=75000, priority=urgent}

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

1. [1] TestPerson{id=P001, salary=45000, score=8.5, status=active, level=2, name=Alice, active=true, tags=junior, department=sales, age=25}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, department=hr, age=16, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, department=marketing, level=7, age=45, active=true}
5. [1] TestPerson{id=P005, age=30, department=sales, name=Eve, salary=55000, active=false, score=8, tags=employee, status=inactive, level=3}
6. [1] TestPerson{id=P006, active=true, tags=test, name=Frank, score=0, status=active, department=qa, level=1, age=0, salary=-5000}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, tags=executive, department=management, name=Grace, status=active, level=9}
8. [1] TestPerson{id=P008, department=support, level=1, score=5.5, status=inactive, name=Henry, age=18, salary=25000, active=false, tags=junior}
9. [1] TestPerson{id=P009, active=true, tags=senior, name=Ivy, salary=68000, score=8.7, status=active, department=engineering, level=6, age=40}
10. [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, age=22, status=active, department=intern, level=1, name=X, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, department=marketing, level=7, age=45, active=true}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, score=8, tags=employee, status=inactive, level=3, age=30, department=sales}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P007, status=active, level=9, age=65, salary=95000, active=true, score=10, tags=executive, department=management, name=Grace}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, status=active, department=engineering, level=6, age=40, active=true, tags=senior, name=Ivy, salary=68000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, status=active, level=2, name=Alice, active=true, tags=junior, department=sales}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active, name=Bob}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, status=pending, priority=normal, region=north}
2. [1] TestOrder{id=O002, product_id=PROD002, total=25.5, status=confirmed, discount=0, region=south, customer_id=P002, amount=1, date=2024-01-20, priority=low}
3. [1] TestOrder{id=O003, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north, customer_id=P001, product_id=PROD003, amount=3}
4. [1] TestOrder{id=O004, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, region=south, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, date=2024-02-10, priority=high, discount=100}
6. [1] TestOrder{id=O006, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005}
7. [1] TestOrder{id=O007, product_id=PROD006, date=2024-03-01, priority=urgent, region=north, customer_id=P007, amount=4, total=600, status=shipped, discount=50}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255}
9. [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, region=north, customer_id=P001, product_id=PROD007, discount=10}
10. [1] TestOrder{id=O010, region=east, date=2024-03-15, status=refunded, discount=0, customer_id=P006, product_id=PROD001, amount=1, total=75000, priority=urgent}

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

1. [1] TestProduct{id=PROD001, available=true, brand=TechCorp, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, rating=4.5, keywords=computer, stock=50}
2. [1] TestProduct{id=PROD002, price=25.5, keywords=peripheral, stock=200, name=Mouse, category=accessories, available=true, rating=4, brand=TechCorp, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, available=false, rating=3.5, price=75, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, price=299.99, available=true, rating=4.8, keywords=display, stock=30, name=Monitor, category=electronics, brand=ScreenPro, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, category=accessories, price=8.5, available=false, brand=OldTech, stock=0, rating=2, keywords=obsolete, supplier=OldSupply, name=OldKeyboard}
6. [1] TestProduct{id=PROD006, category=audio, price=150, rating=4.6, keywords=sound, brand=AudioMax, name=Headphones, available=true, stock=75, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, available=true, rating=3.8, keywords=video, brand=CamTech, price=89.99, stock=25, supplier=CamSupply}

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

1. [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, status=active, level=2, name=Alice, active=true, tags=junior, department=sales}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
3. [1] TestPerson{id=P003, active=false, score=6, department=hr, age=16, tags=intern, status=inactive, level=1, name=Charlie, salary=0}
4. [1] TestPerson{id=P004, age=45, active=true, name=Diana, salary=85000, score=7.8, tags=manager, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, score=8, tags=employee, status=inactive, level=3, age=30, department=sales}
6. [1] TestPerson{id=P006, tags=test, name=Frank, score=0, status=active, department=qa, level=1, age=0, salary=-5000, active=true}
7. [1] TestPerson{id=P007, salary=95000, active=true, score=10, tags=executive, department=management, name=Grace, status=active, level=9, age=65}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, tags=junior, department=support, level=1, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, status=active, department=engineering, level=6, age=40, active=true, tags=senior}
10. [1] TestPerson{id=P010, age=22, status=active, department=intern, level=1, name=X, salary=28000, active=true, score=6.5, tags=temp}

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

1. [1] TestOrder{id=O001, date=2024-01-15, discount=50, status=pending, priority=normal, region=north, customer_id=P001, product_id=PROD001, amount=2, total=1999.98}
2. [1] TestOrder{id=O002, region=south, customer_id=P002, amount=1, date=2024-01-20, priority=low, product_id=PROD002, total=25.5, status=confirmed, discount=0}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{id=O005, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, date=2024-02-10}
6. [1] TestOrder{id=O006, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005}
7. [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, status=shipped, discount=50, product_id=PROD006, date=2024-03-01, priority=urgent, region=north}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255}
9. [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, region=north, customer_id=P001, product_id=PROD007, discount=10}
10. [1] TestOrder{id=O010, date=2024-03-15, status=refunded, discount=0, customer_id=P006, product_id=PROD001, amount=1, total=75000, priority=urgent, region=east}

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

1. [1] TestPerson{id=P001, name=Alice, active=true, tags=junior, department=sales, age=25, salary=45000, score=8.5, status=active, level=2}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
3. [1] TestPerson{id=P003, status=inactive, level=1, name=Charlie, salary=0, active=false, score=6, department=hr, age=16, tags=intern}
4. [1] TestPerson{id=P004, department=marketing, level=7, age=45, active=true, name=Diana, salary=85000, score=7.8, tags=manager, status=active}
5. [1] TestPerson{id=P005, department=sales, name=Eve, salary=55000, active=false, score=8, tags=employee, status=inactive, level=3, age=30}
6. [1] TestPerson{id=P006, level=1, age=0, salary=-5000, active=true, tags=test, name=Frank, score=0, status=active, department=qa}
7. [1] TestPerson{id=P007, status=active, level=9, age=65, salary=95000, active=true, score=10, tags=executive, department=management, name=Grace}
8. [1] TestPerson{id=P008, salary=25000, active=false, tags=junior, department=support, level=1, score=5.5, status=inactive, name=Henry, age=18}
9. [1] TestPerson{id=P009, age=40, active=true, tags=senior, name=Ivy, salary=68000, score=8.7, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, age=22, status=active, department=intern, level=1, name=X, salary=28000}

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

1. [1] TestOrder{id=O001, status=pending, priority=normal, region=north, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50}
2. [1] TestOrder{id=O002, product_id=PROD002, total=25.5, status=confirmed, discount=0, region=south, customer_id=P002, amount=1, date=2024-01-20, priority=low}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{id=O005, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed}
6. [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}
7. [1] TestOrder{id=O007, product_id=PROD006, date=2024-03-01, priority=urgent, region=north, customer_id=P007, amount=4, total=600, status=shipped, discount=50}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255}
9. [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, region=north, customer_id=P001, product_id=PROD007, discount=10}
10. [1] TestOrder{id=O010, total=75000, priority=urgent, region=east, date=2024-03-15, status=refunded, discount=0, customer_id=P006, product_id=PROD001, amount=1}

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

1. [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, status=active, level=2, name=Alice, active=true, tags=junior, department=sales}
2. [1] TestPerson{id=P002, department=engineering, level=5, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, salary=0, active=false, score=6, department=hr, age=16, tags=intern, status=inactive, level=1, name=Charlie}
4. [1] TestPerson{id=P004, age=45, active=true, name=Diana, salary=85000, score=7.8, tags=manager, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, score=8, tags=employee, status=inactive, level=3, age=30, department=sales}
6. [1] TestPerson{id=P006, age=0, salary=-5000, active=true, tags=test, name=Frank, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, tags=executive, department=management, name=Grace, status=active, level=9}
8. [1] TestPerson{id=P008, score=5.5, status=inactive, name=Henry, age=18, salary=25000, active=false, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, salary=68000, score=8.7, status=active, department=engineering, level=6, age=40, active=true, tags=senior, name=Ivy}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, age=22, status=active, department=intern, level=1}

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

1. [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, rating=4.5, keywords=computer, stock=50, available=true, brand=TechCorp, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, available=true, rating=4, brand=TechCorp, supplier=TechSupply, price=25.5, keywords=peripheral, stock=200}
3. [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, available=false, rating=3.5, price=75, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, price=299.99, available=true, rating=4.8, keywords=display, stock=30, name=Monitor, category=electronics, brand=ScreenPro, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, price=8.5, available=false, brand=OldTech, stock=0, rating=2, keywords=obsolete, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, name=Headphones, available=true, stock=75, supplier=AudioSupply, category=audio, price=150, rating=4.6, keywords=sound, brand=AudioMax}
7. [1] TestProduct{id=PROD007, rating=3.8, keywords=video, brand=CamTech, price=89.99, stock=25, supplier=CamSupply, name=Webcam, category=electronics, available=true}

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

1. [1] TestPerson{id=P001, score=8.5, status=active, level=2, name=Alice, active=true, tags=junior, department=sales, age=25, salary=45000}
2. [1] TestPerson{id=P002, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, active=false, score=6, department=hr, age=16, tags=intern, status=inactive, level=1, name=Charlie, salary=0}
4. [1] TestPerson{id=P004, tags=manager, status=active, department=marketing, level=7, age=45, active=true, name=Diana, salary=85000, score=7.8}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, score=8, tags=employee, status=inactive, level=3, age=30, department=sales}
6. [1] TestPerson{id=P006, name=Frank, score=0, status=active, department=qa, level=1, age=0, salary=-5000, active=true, tags=test}
7. [1] TestPerson{id=P007, score=10, tags=executive, department=management, name=Grace, status=active, level=9, age=65, salary=95000, active=true}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, tags=junior, department=support, level=1, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, status=active, department=engineering, level=6, age=40, active=true, tags=senior}
10. [1] TestPerson{id=P010, level=1, name=X, salary=28000, active=true, score=6.5, tags=temp, age=22, status=active, department=intern}

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

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, status=pending, priority=normal, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, amount=1, date=2024-01-20, priority=low, product_id=PROD002, total=25.5, status=confirmed, discount=0, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{id=O005, status=confirmed, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, amount=1, total=999.99}
6. [1] TestOrder{id=O006, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled}
7. [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, status=shipped, discount=50, product_id=PROD006, date=2024-03-01, priority=urgent, region=north}
8. [1] TestOrder{id=O008, status=pending, priority=normal, discount=0, region=south, amount=10, total=255, customer_id=P010, product_id=PROD002, date=2024-03-05}
9. [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, region=north, customer_id=P001, product_id=PROD007, discount=10}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, priority=urgent, region=east, date=2024-03-15, status=refunded, discount=0}

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

1. [1] TestPerson{id=P001, name=Alice, active=true, tags=junior, department=sales, age=25, salary=45000, score=8.5, status=active, level=2}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, department=hr, age=16, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, level=7, age=45, active=true, name=Diana, salary=85000, score=7.8, tags=manager, status=active, department=marketing}
5. [1] TestPerson{id=P005, score=8, tags=employee, status=inactive, level=3, age=30, department=sales, name=Eve, salary=55000, active=false}
6. [1] TestPerson{id=P006, level=1, age=0, salary=-5000, active=true, tags=test, name=Frank, score=0, status=active, department=qa}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, tags=executive, department=management, name=Grace, status=active, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, tags=junior, department=support, level=1, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, active=true, tags=senior, name=Ivy, salary=68000, score=8.7, status=active, department=engineering, level=6, age=40}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, age=22, status=active, department=intern, level=1}
11. [1] TestOrder{id=O001, discount=50, status=pending, priority=normal, region=north, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15}
12. [1] TestOrder{id=O002, product_id=PROD002, total=25.5, status=confirmed, discount=0, region=south, customer_id=P002, amount=1, date=2024-01-20, priority=low}
13. [1] TestOrder{id=O003, total=225, discount=15, region=north, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high}
14. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}
15. [1] TestOrder{id=O005, discount=100, region=south, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, date=2024-02-10, priority=high}
16. [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}
17. [1] TestOrder{id=O007, product_id=PROD006, date=2024-03-01, priority=urgent, region=north, customer_id=P007, amount=4, total=600, status=shipped, discount=50}
18. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255}
19. [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, region=north, customer_id=P001, product_id=PROD007, discount=10}
20. [1] TestOrder{id=O010, status=refunded, discount=0, customer_id=P006, product_id=PROD001, amount=1, total=75000, priority=urgent, region=east, date=2024-03-15}

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

1. [1] TestPerson{id=P001, tags=junior, department=sales, age=25, salary=45000, score=8.5, status=active, level=2, name=Alice, active=true}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, department=hr, age=16, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, age=45, active=true, name=Diana, salary=85000, score=7.8, tags=manager, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, score=8, tags=employee, status=inactive, level=3, age=30, department=sales}
6. [1] TestPerson{id=P006, age=0, salary=-5000, active=true, tags=test, name=Frank, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, status=active, level=9, age=65, salary=95000, active=true, score=10, tags=executive, department=management, name=Grace}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, tags=junior, department=support, level=1, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, status=active, department=engineering, level=6, age=40, active=true, tags=senior, name=Ivy, salary=68000, score=8.7}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, age=22, status=active, department=intern, level=1}

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

1. [1] TestPerson{id=P001, salary=45000, score=8.5, status=active, level=2, name=Alice, active=true, tags=junior, department=sales, age=25}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, department=hr, age=16, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, salary=85000, score=7.8, tags=manager, status=active, department=marketing, level=7, age=45, active=true, name=Diana}
5. [1] TestPerson{id=P005, age=30, department=sales, name=Eve, salary=55000, active=false, score=8, tags=employee, status=inactive, level=3}
6. [1] TestPerson{id=P006, age=0, salary=-5000, active=true, tags=test, name=Frank, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, tags=executive, department=management, name=Grace, status=active, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, tags=junior, department=support, level=1, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, level=6, age=40, active=true, tags=senior, name=Ivy, salary=68000, score=8.7, status=active, department=engineering}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, age=22, status=active, department=intern, level=1}
11. [1] TestOrder{id=O001, status=pending, priority=normal, region=north, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50}
12. [1] TestOrder{id=O002, region=south, customer_id=P002, amount=1, date=2024-01-20, priority=low, product_id=PROD002, total=25.5, status=confirmed, discount=0}
13. [1] TestOrder{id=O003, priority=high, total=225, discount=15, region=north, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped}
14. [1] TestOrder{id=O004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004}
15. [1] TestOrder{id=O005, total=999.99, status=confirmed, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, amount=1}
16. [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}
17. [1] TestOrder{id=O007, total=600, status=shipped, discount=50, product_id=PROD006, date=2024-03-01, priority=urgent, region=north, customer_id=P007, amount=4}
18. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255}
19. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, priority=low, region=north, customer_id=P001, product_id=PROD007, discount=10, amount=1}
20. [1] TestOrder{id=O010, status=refunded, discount=0, customer_id=P006, product_id=PROD001, amount=1, total=75000, priority=urgent, region=east, date=2024-03-15}

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

1. [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, status=active, level=2, name=Alice, active=true, tags=junior, department=sales}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, department=hr, age=16, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, age=45, active=true, name=Diana, salary=85000, score=7.8, tags=manager, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, score=8, tags=employee, status=inactive, level=3, age=30, department=sales}
6. [1] TestPerson{id=P006, name=Frank, score=0, status=active, department=qa, level=1, age=0, salary=-5000, active=true, tags=test}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, tags=executive, department=management, name=Grace, status=active, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, tags=junior, department=support, level=1, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, status=active, department=engineering, level=6, age=40, active=true, tags=senior, name=Ivy, salary=68000, score=8.7}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, age=22, status=active, department=intern, level=1}

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

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, status=pending, priority=normal, region=north}
2. [1] TestOrder{id=O002, region=south, customer_id=P002, amount=1, date=2024-01-20, priority=low, product_id=PROD002, total=25.5, status=confirmed, discount=0}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north, customer_id=P001}
4. [1] TestOrder{id=O004, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, discount=100, region=south, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, date=2024-02-10, priority=high}
6. [1] TestOrder{id=O006, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled}
7. [1] TestOrder{id=O007, priority=urgent, region=north, customer_id=P007, amount=4, total=600, status=shipped, discount=50, product_id=PROD006, date=2024-03-01}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255}
9. [1] TestOrder{id=O009, status=completed, priority=low, region=north, customer_id=P001, product_id=PROD007, discount=10, amount=1, total=89.99, date=2024-03-10}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, priority=urgent, region=east, date=2024-03-15, status=refunded, discount=0}

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

1. [1] TestPerson{id=P001, active=true, tags=junior, department=sales, age=25, salary=45000, score=8.5, status=active, level=2, name=Alice}
2. [1] TestPerson{id=P002, level=5, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, department=hr, age=16, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, age=45, active=true, name=Diana, salary=85000, score=7.8, tags=manager, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, score=8, tags=employee, status=inactive, level=3, age=30, department=sales}
6. [1] TestPerson{id=P006, age=0, salary=-5000, active=true, tags=test, name=Frank, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, score=10, tags=executive, department=management, name=Grace, status=active, level=9, age=65, salary=95000, active=true}
8. [1] TestPerson{id=P008, score=5.5, status=inactive, name=Henry, age=18, salary=25000, active=false, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, age=40, active=true, tags=senior, name=Ivy, salary=68000, score=8.7, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, level=1, name=X, salary=28000, active=true, score=6.5, tags=temp, age=22, status=active, department=intern}

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
