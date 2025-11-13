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

1. [1] TestPerson{id=P001, age=25, salary=45000, active=true, tags=junior, name=Alice, score=8.5, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, level=5, name=Bob, active=true, score=9.2, tags=senior, age=35, salary=75000, status=active, department=engineering}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, department=hr, level=1, active=false, score=6, status=inactive}
4. [1] TestPerson{id=P004, department=marketing, name=Diana, salary=85000, active=true, score=7.8, level=7, age=45, tags=manager, status=active}
5. [1] TestPerson{id=P005, tags=employee, status=inactive, salary=55000, active=false, department=sales, level=3, name=Eve, age=30, score=8}
6. [1] TestPerson{id=P006, name=Frank, score=0, status=active, department=qa, level=1, age=0, salary=-5000, active=true, tags=test}
7. [1] TestPerson{id=P007, active=true, status=active, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, active=false, tags=junior, level=1, score=5.5, status=inactive, department=support, name=Henry, age=18, salary=25000}
9. [1] TestPerson{id=P009, level=6, age=40, salary=68000, active=true, tags=senior, status=active, name=Ivy, score=8.7, department=engineering}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, name=X, active=true, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, name=X, active=true, status=active, department=intern, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, score=8.5, status=active, department=sales, level=2, age=25, salary=45000, active=true, tags=junior}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P002, level=5, name=Bob, active=true, score=9.2, tags=senior, age=35, salary=75000, status=active, department=engineering}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, age=45, tags=manager, status=active, department=marketing, name=Diana, salary=85000, active=true, score=7.8, level=7}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, score=8, tags=employee, status=inactive, salary=55000, active=false, department=sales, level=3}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, status=inactive, name=Charlie, age=16, salary=0, tags=intern, department=hr, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, active=true, status=active, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, active=false, tags=junior, level=1, score=5.5, status=inactive, department=support, name=Henry}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, active=true, tags=senior, status=active, name=Ivy, score=8.7, department=engineering, level=6}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98}
2. [1] TestOrder{id=O002, product_id=PROD002, total=25.5, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed, priority=low, customer_id=P002}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, status=shipped, region=north, product_id=PROD003, total=225, priority=high, discount=15}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east, customer_id=P004}
5. [1] TestOrder{id=O005, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, discount=0, amount=2, date=2024-02-15, status=cancelled, priority=low, region=west}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, total=600, discount=50, region=north}
8. [1] TestOrder{id=O008, discount=0, region=south, total=255, status=pending, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, priority=normal}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, region=north}
10. [1] TestOrder{id=O010, product_id=PROD001, total=75000, date=2024-03-15, status=refunded, customer_id=P006, amount=1, priority=urgent, discount=0, region=east}

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

1. [1] TestPerson{id=P001, name=Alice, score=8.5, status=active, department=sales, level=2, age=25, salary=45000, active=true, tags=junior}
2. [1] TestPerson{id=P002, name=Bob, active=true, score=9.2, tags=senior, age=35, salary=75000, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, level=1, active=false, score=6, status=inactive, name=Charlie, age=16, salary=0, tags=intern, department=hr}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, level=7, age=45, tags=manager, status=active, department=marketing}
5. [1] TestPerson{id=P005, level=3, name=Eve, age=30, score=8, tags=employee, status=inactive, salary=55000, active=false, department=sales}
6. [1] TestPerson{id=P006, age=0, salary=-5000, active=true, tags=test, name=Frank, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}
8. [1] TestPerson{id=P008, score=5.5, status=inactive, department=support, name=Henry, age=18, salary=25000, active=false, tags=junior, level=1}
9. [1] TestPerson{id=P009, score=8.7, department=engineering, level=6, age=40, salary=68000, active=true, tags=senior, status=active, name=Ivy}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, name=X, active=true, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, score=8.5, status=active, department=sales, level=2, age=25, salary=45000, active=true, tags=junior}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, active=true, score=9.2, tags=senior, age=35, salary=75000, status=active, department=engineering, level=5}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, age=45, tags=manager, status=active, department=marketing, name=Diana, salary=85000, active=true, score=7.8, level=7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, score=8, tags=employee, status=inactive, salary=55000, active=false, department=sales, level=3}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, score=8.7, department=engineering, level=6, age=40, salary=68000, active=true, tags=senior, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north}
2. [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, product_id=PROD002, total=25.5, discount=0, region=south}
3. [1] TestOrder{id=O003, amount=3, date=2024-02-01, status=shipped, region=north, product_id=PROD003, total=225, priority=high, discount=15, customer_id=P001}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100}
6. [1] TestOrder{id=O006, amount=2, date=2024-02-15, status=cancelled, priority=low, region=west, customer_id=P005, product_id=PROD005, total=999.98, discount=0}
7. [1] TestOrder{id=O007, discount=50, region=north, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, total=600}
8. [1] TestOrder{id=O008, total=255, status=pending, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, priority=normal, discount=0, region=south}
9. [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, region=north, customer_id=P001, amount=1, total=89.99}
10. [1] TestOrder{id=O010, status=refunded, customer_id=P006, amount=1, priority=urgent, discount=0, region=east, product_id=PROD001, total=75000, date=2024-03-15}

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

1. [1] TestProduct{id=PROD001, brand=TechCorp, stock=50, name=Laptop, category=electronics, price=999.99, available=true, keywords=computer, supplier=TechSupply, rating=4.5}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, available=true, rating=4, brand=TechCorp, stock=200, price=25.5, keywords=peripheral, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, supplier=KeySupply, category=accessories, price=75, rating=3.5, name=Keyboard, available=false, keywords=typing, brand=KeyTech, stock=0}
4. [1] TestProduct{id=PROD004, brand=ScreenPro, price=299.99, rating=4.8, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics, available=true, keywords=display}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, price=8.5, brand=OldTech, category=accessories, available=false, rating=2, keywords=obsolete, stock=0, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, stock=75, rating=4.6, supplier=AudioSupply, name=Headphones, category=audio, price=150, available=true, keywords=sound, brand=AudioMax}
7. [1] TestProduct{id=PROD007, rating=3.8, brand=CamTech, supplier=CamSupply, price=89.99, keywords=video, stock=25, name=Webcam, category=electronics, available=true}

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

1. [1] TestPerson{id=P001, name=Alice, score=8.5, status=active, department=sales, level=2, age=25, salary=45000, active=true, tags=junior}
2. [1] TestPerson{id=P002, name=Bob, active=true, score=9.2, tags=senior, age=35, salary=75000, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, level=1, active=false, score=6, status=inactive, name=Charlie, age=16, salary=0, tags=intern, department=hr}
4. [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, level=7, age=45, tags=manager, status=active, department=marketing, name=Diana}
5. [1] TestPerson{id=P005, name=Eve, age=30, score=8, tags=employee, status=inactive, salary=55000, active=false, department=sales, level=3}
6. [1] TestPerson{id=P006, age=0, salary=-5000, active=true, tags=test, name=Frank, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}
8. [1] TestPerson{id=P008, salary=25000, active=false, tags=junior, level=1, score=5.5, status=inactive, department=support, name=Henry, age=18}
9. [1] TestPerson{id=P009, status=active, name=Ivy, score=8.7, department=engineering, level=6, age=40, salary=68000, active=true, tags=senior}
10. [1] TestPerson{id=P010, name=X, active=true, status=active, department=intern, level=1, age=22, salary=28000, score=6.5, tags=temp}

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

1. [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50}
2. [1] TestOrder{id=O002, status=confirmed, priority=low, customer_id=P002, product_id=PROD002, total=25.5, discount=0, region=south, amount=1, date=2024-01-20}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, status=shipped, region=north, product_id=PROD003, total=225, priority=high, discount=15}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{id=O005, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002}
6. [1] TestOrder{id=O006, amount=2, date=2024-02-15, status=cancelled, priority=low, region=west, customer_id=P005, product_id=PROD005, total=999.98, discount=0}
7. [1] TestOrder{id=O007, discount=50, region=north, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, total=600}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, priority=normal, discount=0, region=south, total=255, status=pending}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, region=north}
10. [1] TestOrder{id=O010, priority=urgent, discount=0, region=east, product_id=PROD001, total=75000, date=2024-03-15, status=refunded, customer_id=P006, amount=1}

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

1. [1] TestPerson{id=P001, salary=45000, active=true, tags=junior, name=Alice, score=8.5, status=active, department=sales, level=2, age=25}
2. [1] TestPerson{id=P002, name=Bob, active=true, score=9.2, tags=senior, age=35, salary=75000, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, department=hr, level=1, active=false, score=6, status=inactive}
4. [1] TestPerson{id=P004, department=marketing, name=Diana, salary=85000, active=true, score=7.8, level=7, age=45, tags=manager, status=active}
5. [1] TestPerson{id=P005, status=inactive, salary=55000, active=false, department=sales, level=3, name=Eve, age=30, score=8, tags=employee}
6. [1] TestPerson{id=P006, name=Frank, score=0, status=active, department=qa, level=1, age=0, salary=-5000, active=true, tags=test}
7. [1] TestPerson{id=P007, active=true, status=active, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, tags=junior, level=1, score=5.5, status=inactive, department=support, name=Henry, age=18, salary=25000, active=false}
9. [1] TestPerson{id=P009, tags=senior, status=active, name=Ivy, score=8.7, department=engineering, level=6, age=40, salary=68000, active=true}
10. [1] TestPerson{id=P010, level=1, age=22, salary=28000, score=6.5, tags=temp, name=X, active=true, status=active, department=intern}

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

1. [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50}
2. [1] TestOrder{id=O002, total=25.5, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, product_id=PROD002}
3. [1] TestOrder{id=O003, total=225, priority=high, discount=15, customer_id=P001, amount=3, date=2024-02-01, status=shipped, region=north, product_id=PROD003}
4. [1] TestOrder{id=O004, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, discount=0, amount=2, date=2024-02-15, status=cancelled, priority=low, region=west}
7. [1] TestOrder{id=O007, discount=50, region=north, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, total=600}
8. [1] TestOrder{id=O008, status=pending, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, priority=normal, discount=0, region=south, total=255}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, priority=urgent, discount=0, region=east, product_id=PROD001, total=75000, date=2024-03-15, status=refunded}

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

1. [1] TestPerson{id=P001, salary=45000, active=true, tags=junior, name=Alice, score=8.5, status=active, department=sales, level=2, age=25}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, age=35, salary=75000, status=active, department=engineering, level=5, name=Bob, active=true}
3. [1] TestPerson{id=P003, level=1, active=false, score=6, status=inactive, name=Charlie, age=16, salary=0, tags=intern, department=hr}
4. [1] TestPerson{id=P004, department=marketing, name=Diana, salary=85000, active=true, score=7.8, level=7, age=45, tags=manager, status=active}
5. [1] TestPerson{id=P005, salary=55000, active=false, department=sales, level=3, name=Eve, age=30, score=8, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, level=1, age=0, salary=-5000, active=true, tags=test, name=Frank, score=0, status=active, department=qa}
7. [1] TestPerson{id=P007, active=true, status=active, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, salary=25000, active=false, tags=junior, level=1, score=5.5, status=inactive, department=support, name=Henry, age=18}
9. [1] TestPerson{id=P009, score=8.7, department=engineering, level=6, age=40, salary=68000, active=true, tags=senior, status=active, name=Ivy}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, name=X, active=true, status=active, department=intern, level=1}

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

1. [1] TestProduct{id=PROD001, supplier=TechSupply, rating=4.5, brand=TechCorp, stock=50, name=Laptop, category=electronics, price=999.99, available=true, keywords=computer}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, available=true, rating=4, brand=TechCorp, stock=200, price=25.5, keywords=peripheral, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, available=false, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply, category=accessories, price=75, rating=3.5, name=Keyboard}
4. [1] TestProduct{id=PROD004, name=Monitor, category=electronics, available=true, keywords=display, brand=ScreenPro, price=299.99, rating=4.8, stock=30, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, price=8.5, brand=OldTech, category=accessories, available=false, rating=2, keywords=obsolete, stock=0, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, brand=AudioMax, stock=75, rating=4.6, supplier=AudioSupply, name=Headphones, category=audio, price=150, available=true, keywords=sound}
7. [1] TestProduct{id=PROD007, brand=CamTech, supplier=CamSupply, price=89.99, keywords=video, stock=25, name=Webcam, category=electronics, available=true, rating=3.8}

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

1. [1] TestPerson{id=P001, salary=45000, active=true, tags=junior, name=Alice, score=8.5, status=active, department=sales, level=2, age=25}
2. [1] TestPerson{id=P002, name=Bob, active=true, score=9.2, tags=senior, age=35, salary=75000, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, department=hr, level=1, active=false, score=6, status=inactive}
4. [1] TestPerson{id=P004, level=7, age=45, tags=manager, status=active, department=marketing, name=Diana, salary=85000, active=true, score=7.8}
5. [1] TestPerson{id=P005, tags=employee, status=inactive, salary=55000, active=false, department=sales, level=3, name=Eve, age=30, score=8}
6. [1] TestPerson{id=P006, age=0, salary=-5000, active=true, tags=test, name=Frank, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active, name=Grace}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, tags=junior, level=1, score=5.5, status=inactive, department=support}
9. [1] TestPerson{id=P009, name=Ivy, score=8.7, department=engineering, level=6, age=40, salary=68000, active=true, tags=senior, status=active}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, name=X, active=true, status=active, department=intern, level=1}

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

1. [1] TestOrder{id=O001, discount=50, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal}
2. [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, product_id=PROD002, total=25.5, discount=0, region=south}
3. [1] TestOrder{id=O003, date=2024-02-01, status=shipped, region=north, product_id=PROD003, total=225, priority=high, discount=15, customer_id=P001, amount=3}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}
6. [1] TestOrder{id=O006, amount=2, date=2024-02-15, status=cancelled, priority=low, region=west, customer_id=P005, product_id=PROD005, total=999.98, discount=0}
7. [1] TestOrder{id=O007, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, total=600, discount=50, region=north, product_id=PROD006, amount=4}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, priority=normal, discount=0, region=south, total=255, status=pending}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, region=north, customer_id=P001, amount=1}
10. [1] TestOrder{id=O010, region=east, product_id=PROD001, total=75000, date=2024-03-15, status=refunded, customer_id=P006, amount=1, priority=urgent, discount=0}

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

1. [1] TestPerson{id=P001, score=8.5, status=active, department=sales, level=2, age=25, salary=45000, active=true, tags=junior, name=Alice}
2. [1] TestPerson{id=P002, age=35, salary=75000, status=active, department=engineering, level=5, name=Bob, active=true, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, department=hr, level=1, active=false, score=6, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, level=7, age=45, tags=manager, status=active, department=marketing}
5. [1] TestPerson{id=P005, salary=55000, active=false, department=sales, level=3, name=Eve, age=30, score=8, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, department=qa, level=1, age=0, salary=-5000, active=true, tags=test, name=Frank, score=0, status=active}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, tags=junior, level=1, score=5.5, status=inactive, department=support}
9. [1] TestPerson{id=P009, status=active, name=Ivy, score=8.7, department=engineering, level=6, age=40, salary=68000, active=true, tags=senior}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, name=X, active=true, status=active, department=intern, level=1}
11. [1] TestOrder{id=O001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001}
12. [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, product_id=PROD002, total=25.5, discount=0, region=south}
13. [1] TestOrder{id=O003, priority=high, discount=15, customer_id=P001, amount=3, date=2024-02-01, status=shipped, region=north, product_id=PROD003, total=225}
14. [1] TestOrder{id=O004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004}
15. [1] TestOrder{id=O005, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1}
16. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, discount=0, amount=2, date=2024-02-15, status=cancelled, priority=low, region=west}
17. [1] TestOrder{id=O007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, total=600, discount=50, region=north}
18. [1] TestOrder{id=O008, discount=0, region=south, total=255, status=pending, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, priority=normal}
19. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, region=north}
20. [1] TestOrder{id=O010, region=east, product_id=PROD001, total=75000, date=2024-03-15, status=refunded, customer_id=P006, amount=1, priority=urgent, discount=0}

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

1. [1] TestPerson{id=P001, score=8.5, status=active, department=sales, level=2, age=25, salary=45000, active=true, tags=junior, name=Alice}
2. [1] TestPerson{id=P002, salary=75000, status=active, department=engineering, level=5, name=Bob, active=true, score=9.2, tags=senior, age=35}
3. [1] TestPerson{id=P003, level=1, active=false, score=6, status=inactive, name=Charlie, age=16, salary=0, tags=intern, department=hr}
4. [1] TestPerson{id=P004, tags=manager, status=active, department=marketing, name=Diana, salary=85000, active=true, score=7.8, level=7, age=45}
5. [1] TestPerson{id=P005, level=3, name=Eve, age=30, score=8, tags=employee, status=inactive, salary=55000, active=false, department=sales}
6. [1] TestPerson{id=P006, status=active, department=qa, level=1, age=0, salary=-5000, active=true, tags=test, name=Frank, score=0}
7. [1] TestPerson{id=P007, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active, name=Grace, age=65}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, tags=junior, level=1, score=5.5, status=inactive, department=support}
9. [1] TestPerson{id=P009, tags=senior, status=active, name=Ivy, score=8.7, department=engineering, level=6, age=40, salary=68000, active=true}
10. [1] TestPerson{id=P010, active=true, status=active, department=intern, level=1, age=22, salary=28000, score=6.5, tags=temp, name=X}

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

1. [1] TestPerson{id=P001, status=active, department=sales, level=2, age=25, salary=45000, active=true, tags=junior, name=Alice, score=8.5}
2. [1] TestPerson{id=P002, tags=senior, age=35, salary=75000, status=active, department=engineering, level=5, name=Bob, active=true, score=9.2}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, department=hr, level=1, active=false, score=6, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, level=7, age=45, tags=manager, status=active, department=marketing}
5. [1] TestPerson{id=P005, salary=55000, active=false, department=sales, level=3, name=Eve, age=30, score=8, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, level=1, age=0, salary=-5000, active=true, tags=test, name=Frank, score=0, status=active, department=qa}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}
8. [1] TestPerson{id=P008, score=5.5, status=inactive, department=support, name=Henry, age=18, salary=25000, active=false, tags=junior, level=1}
9. [1] TestPerson{id=P009, age=40, salary=68000, active=true, tags=senior, status=active, name=Ivy, score=8.7, department=engineering, level=6}
10. [1] TestPerson{id=P010, name=X, active=true, status=active, department=intern, level=1, age=22, salary=28000, score=6.5, tags=temp}
11. [1] TestOrder{id=O001, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15}
12. [1] TestOrder{id=O002, total=25.5, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, product_id=PROD002}
13. [1] TestOrder{id=O003, priority=high, discount=15, customer_id=P001, amount=3, date=2024-02-01, status=shipped, region=north, product_id=PROD003, total=225}
14. [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1}
15. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}
16. [1] TestOrder{id=O006, amount=2, date=2024-02-15, status=cancelled, priority=low, region=west, customer_id=P005, product_id=PROD005, total=999.98, discount=0}
17. [1] TestOrder{id=O007, customer_id=P007, total=600, discount=50, region=north, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent}
18. [1] TestOrder{id=O008, amount=10, date=2024-03-05, priority=normal, discount=0, region=south, total=255, status=pending, customer_id=P010, product_id=PROD002}
19. [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, region=north, customer_id=P001}
20. [1] TestOrder{id=O010, discount=0, region=east, product_id=PROD001, total=75000, date=2024-03-15, status=refunded, customer_id=P006, amount=1, priority=urgent}

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

1. [1] TestPerson{id=P001, score=8.5, status=active, department=sales, level=2, age=25, salary=45000, active=true, tags=junior, name=Alice}
2. [1] TestPerson{id=P002, age=35, salary=75000, status=active, department=engineering, level=5, name=Bob, active=true, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, level=1, active=false, score=6, status=inactive, name=Charlie, age=16, salary=0, tags=intern, department=hr}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, level=7, age=45, tags=manager, status=active, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, age=30, score=8, tags=employee, status=inactive, salary=55000, active=false, department=sales, level=3}
6. [1] TestPerson{id=P006, tags=test, name=Frank, score=0, status=active, department=qa, level=1, age=0, salary=-5000, active=true}
7. [1] TestPerson{id=P007, score=10, tags=executive, department=management, level=9, active=true, status=active, name=Grace, age=65, salary=95000}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, tags=junior, level=1, score=5.5, status=inactive, department=support}
9. [1] TestPerson{id=P009, status=active, name=Ivy, score=8.7, department=engineering, level=6, age=40, salary=68000, active=true, tags=senior}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, name=X, active=true, status=active, department=intern, level=1}

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

1. [1] TestOrder{id=O001, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed, priority=low}
3. [1] TestOrder{id=O003, discount=15, customer_id=P001, amount=3, date=2024-02-01, status=shipped, region=north, product_id=PROD003, total=225, priority=high}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{id=O005, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high}
6. [1] TestOrder{id=O006, status=cancelled, priority=low, region=west, customer_id=P005, product_id=PROD005, total=999.98, discount=0, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, total=600, discount=50, region=north}
8. [1] TestOrder{id=O008, total=255, status=pending, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, priority=normal, discount=0, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, priority=urgent, discount=0, region=east, product_id=PROD001, total=75000, date=2024-03-15, status=refunded}

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

1. [1] TestPerson{id=P001, age=25, salary=45000, active=true, tags=junior, name=Alice, score=8.5, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, name=Bob, active=true, score=9.2, tags=senior, age=35, salary=75000, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, tags=intern, department=hr, level=1, active=false, score=6, status=inactive, name=Charlie, age=16, salary=0}
4. [1] TestPerson{id=P004, age=45, tags=manager, status=active, department=marketing, name=Diana, salary=85000, active=true, score=7.8, level=7}
5. [1] TestPerson{id=P005, level=3, name=Eve, age=30, score=8, tags=employee, status=inactive, salary=55000, active=false, department=sales}
6. [1] TestPerson{id=P006, name=Frank, score=0, status=active, department=qa, level=1, age=0, salary=-5000, active=true, tags=test}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}
8. [1] TestPerson{id=P008, status=inactive, department=support, name=Henry, age=18, salary=25000, active=false, tags=junior, level=1, score=5.5}
9. [1] TestPerson{id=P009, department=engineering, level=6, age=40, salary=68000, active=true, tags=senior, status=active, name=Ivy, score=8.7}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, name=X, active=true, status=active, department=intern, level=1}

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
