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

1. [1] TestPerson{id=P001, status=active, department=sales, active=true, level=2, name=Alice, age=25, salary=45000, score=8.5, tags=junior}
2. [1] TestPerson{id=P002, age=35, score=9.2, name=Bob, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, level=1, age=16, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, status=active, department=marketing, level=7, salary=85000, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales, level=3, score=8, status=inactive}
6. [1] TestPerson{id=P006, status=active, department=qa, name=Frank, active=true, score=0, tags=test, level=1, age=0, salary=-5000}
7. [1] TestPerson{id=P007, name=Grace, status=active, department=management, age=65, salary=95000, active=true, score=10, tags=executive, level=9}
8. [1] TestPerson{id=P008, salary=25000, department=support, level=1, name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18}
9. [1] TestPerson{id=P009, name=Ivy, active=true, score=8.7, tags=senior, status=active, level=6, age=40, salary=68000, department=engineering}
10. [1] TestPerson{id=P010, department=intern, name=X, active=true, status=active, level=1, age=22, salary=28000, score=6.5, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, active=true, tags=senior, status=active, department=engineering, level=5, age=35, score=9.2, name=Bob, salary=75000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, status=inactive, level=1, age=16, salary=0, score=6, department=hr, name=Charlie, active=false}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, status=active, department=marketing, level=7, salary=85000, tags=manager}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, department=support, level=1, name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, salary=25000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, department=engineering, name=Ivy, active=true, score=8.7, tags=senior, status=active, level=6}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, department=intern, name=X, active=true, status=active, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, level=3, score=8, status=inactive, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, level=9, name=Grace, status=active, department=management, age=65, salary=95000, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, amount=2, status=pending, priority=normal}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, discount=0, amount=1, total=25.5, priority=low, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, discount=15, date=2024-02-01, priority=high, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east}
5. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, date=2024-02-15, status=cancelled, priority=low, discount=0, customer_id=P005, total=999.98, region=west}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50, region=north, amount=4, priority=urgent}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, discount=0, total=255, date=2024-03-05, status=pending, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, product_id=PROD007}
10. [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, customer_id=P006, total=75000, discount=0, region=east, product_id=PROD001, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, tags=junior, status=inactive, age=18, salary=25000, department=support, level=1, name=Henry}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, level=6, age=40, salary=68000, department=engineering, name=Ivy, active=true, score=8.7}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, active=true, level=2, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales, level=3, score=8, status=inactive}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P010, status=active, level=1, age=22, salary=28000, score=6.5, tags=temp, department=intern, name=X, active=true}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, age=35, score=9.2, name=Bob}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, tags=intern, status=inactive, level=1, age=16, salary=0, score=6, department=hr}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, status=active, department=marketing, level=7, salary=85000, tags=manager}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, status=active, department=qa, name=Frank, active=true, score=0, tags=test, level=1}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P007, status=active, department=management, age=65, salary=95000, active=true, score=10, tags=executive, level=9, name=Grace}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, age=35, score=9.2}
3. [1] TestPerson{id=P003, active=false, tags=intern, status=inactive, level=1, age=16, salary=0, score=6, department=hr, name=Charlie}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, status=active, department=marketing, level=7, salary=85000, tags=manager}
5. [1] TestPerson{id=P005, active=false, tags=employee, department=sales, level=3, score=8, status=inactive, name=Eve, age=30, salary=55000}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, tags=test, level=1, age=0, salary=-5000, status=active, department=qa}
7. [1] TestPerson{id=P007, level=9, name=Grace, status=active, department=management, age=65, salary=95000, active=true, score=10, tags=executive}
8. [1] TestPerson{id=P008, department=support, level=1, name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, salary=25000}
9. [1] TestPerson{id=P009, tags=senior, status=active, level=6, age=40, salary=68000, department=engineering, name=Ivy, active=true, score=8.7}
10. [1] TestPerson{id=P010, name=X, active=true, status=active, level=1, age=22, salary=28000, score=6.5, tags=temp, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, tags=executive, level=9, name=Grace, status=active, department=management, age=65, salary=95000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, score=8.7, tags=senior, status=active, level=6, age=40, salary=68000, department=engineering}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, age=35, score=9.2}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, status=active, department=marketing, level=7, salary=85000, tags=manager}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales, level=3, score=8, status=inactive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, status=pending, priority=normal, product_id=PROD001, total=1999.98, date=2024-01-15, discount=50, region=north}
2. [1] TestOrder{id=O002, discount=0, amount=1, total=25.5, priority=low, region=south, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, discount=15, date=2024-02-01, priority=high, region=north}
4. [1] TestOrder{id=O004, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100}
6. [1] TestOrder{id=O006, total=999.98, region=west, product_id=PROD005, amount=2, date=2024-02-15, status=cancelled, priority=low, discount=0, customer_id=P005}
7. [1] TestOrder{id=O007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50, region=north, amount=4, priority=urgent, customer_id=P007}
8. [1] TestOrder{id=O008, status=pending, region=south, customer_id=P010, product_id=PROD002, amount=10, priority=normal, discount=0, total=255, date=2024-03-05}
9. [1] TestOrder{id=O009, discount=10, region=north, product_id=PROD007, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}
10. [1] TestOrder{id=O010, amount=1, date=2024-03-15, status=refunded, priority=urgent, customer_id=P006, total=75000, discount=0, region=east, product_id=PROD001}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O006, region=west, product_id=PROD005, amount=2, date=2024-02-15, status=cancelled, priority=low, discount=0, customer_id=P005, total=999.98}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, discount=0, total=255, date=2024-03-05, status=pending, region=south}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, product_id=PROD007}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50, region=north, amount=4, priority=urgent}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O001, amount=2, status=pending, priority=normal, product_id=PROD001, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, discount=0, amount=1, total=25.5, priority=low, region=south}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O003, region=north, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, discount=15, date=2024-02-01, priority=high}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, category=electronics, price=999.99, available=true, rating=4.5, keywords=computer, stock=50, brand=TechCorp, supplier=TechSupply, name=Laptop}
2. [1] TestProduct{id=PROD002, category=accessories, available=true, keywords=peripheral, brand=TechCorp, stock=200, name=Mouse, price=25.5, rating=4, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, keywords=typing, brand=KeyTech, supplier=KeySupply, price=75, available=false, rating=3.5, stock=0}
4. [1] TestProduct{id=PROD004, category=electronics, price=299.99, rating=4.8, keywords=display, brand=ScreenPro, stock=30, name=Monitor, available=true, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, keywords=obsolete, supplier=OldSupply, price=8.5, available=false, rating=2, brand=OldTech, stock=0}
6. [1] TestProduct{id=PROD006, stock=75, supplier=AudioSupply, price=150, available=true, name=Headphones, category=audio, rating=4.6, keywords=sound, brand=AudioMax}
7. [1] TestProduct{id=PROD007, category=electronics, price=89.99, available=true, rating=3.8, brand=CamTech, stock=25, name=Webcam, keywords=video, supplier=CamSupply}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}
2. [1] TestPerson{id=P002, age=35, score=9.2, name=Bob, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1, age=16}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, status=active, department=marketing, level=7, salary=85000, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales, level=3, score=8, status=inactive}
6. [1] TestPerson{id=P006, age=0, salary=-5000, status=active, department=qa, name=Frank, active=true, score=0, tags=test, level=1}
7. [1] TestPerson{id=P007, status=active, department=management, age=65, salary=95000, active=true, score=10, tags=executive, level=9, name=Grace}
8. [1] TestPerson{id=P008, age=18, salary=25000, department=support, level=1, name=Henry, active=false, score=5.5, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, level=6, age=40, salary=68000, department=engineering, name=Ivy, active=true, score=8.7, tags=senior, status=active}
10. [1] TestPerson{id=P010, status=active, level=1, age=22, salary=28000, score=6.5, tags=temp, department=intern, name=X, active=true}

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

1. [1] TestOrder{id=O001, region=north, customer_id=P001, amount=2, status=pending, priority=normal, product_id=PROD001, total=1999.98, date=2024-01-15, discount=50}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, discount=0, amount=1, total=25.5, priority=low, region=south}
3. [1] TestOrder{id=O003, discount=15, date=2024-02-01, priority=high, region=north, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped}
4. [1] TestOrder{id=O004, discount=0, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal}
5. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}
6. [1] TestOrder{id=O006, discount=0, customer_id=P005, total=999.98, region=west, product_id=PROD005, amount=2, date=2024-02-15, status=cancelled, priority=low}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50, region=north, amount=4, priority=urgent}
8. [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, region=south, customer_id=P010, product_id=PROD002, amount=10, priority=normal, discount=0}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, product_id=PROD007}
10. [1] TestOrder{id=O010, amount=1, date=2024-03-15, status=refunded, priority=urgent, customer_id=P006, total=75000, discount=0, region=east, product_id=PROD001}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, age=35, score=9.2}
3. [1] TestPerson{id=P003, age=16, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, status=active, department=marketing, level=7, salary=85000, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales, level=3, score=8, status=inactive}
6. [1] TestPerson{id=P006, department=qa, name=Frank, active=true, score=0, tags=test, level=1, age=0, salary=-5000, status=active}
7. [1] TestPerson{id=P007, name=Grace, status=active, department=management, age=65, salary=95000, active=true, score=10, tags=executive, level=9}
8. [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, salary=25000, department=support, level=1}
9. [1] TestPerson{id=P009, status=active, level=6, age=40, salary=68000, department=engineering, name=Ivy, active=true, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, department=intern, name=X, active=true, status=active, level=1}

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

1. [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, amount=2, status=pending, priority=normal}
2. [1] TestOrder{id=O002, total=25.5, priority=low, region=south, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, discount=0, amount=1}
3. [1] TestOrder{id=O003, date=2024-02-01, priority=high, region=north, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, discount=15}
4. [1] TestOrder{id=O004, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}
5. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}
6. [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, customer_id=P005, total=999.98, region=west, product_id=PROD005, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50, region=north, amount=4, priority=urgent}
8. [1] TestOrder{id=O008, region=south, customer_id=P010, product_id=PROD002, amount=10, priority=normal, discount=0, total=255, date=2024-03-05, status=pending}
9. [1] TestOrder{id=O009, status=completed, priority=low, discount=10, region=north, product_id=PROD007, customer_id=P001, amount=1, total=89.99, date=2024-03-10}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, customer_id=P006, total=75000, discount=0, region=east}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, age=35, score=9.2}
3. [1] TestPerson{id=P003, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1, age=16, salary=0, score=6}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, status=active, department=marketing, level=7, salary=85000, tags=manager}
5. [1] TestPerson{id=P005, level=3, score=8, status=inactive, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, tags=test, level=1, age=0, salary=-5000, status=active, department=qa}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, tags=executive, level=9, name=Grace, status=active, department=management}
8. [1] TestPerson{id=P008, active=false, score=5.5, tags=junior, status=inactive, age=18, salary=25000, department=support, level=1, name=Henry}
9. [1] TestPerson{id=P009, level=6, age=40, salary=68000, department=engineering, name=Ivy, active=true, score=8.7, tags=senior, status=active}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, department=intern, name=X, active=true, status=active, level=1}

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

1. [1] TestProduct{id=PROD001, stock=50, brand=TechCorp, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, available=true, rating=4.5, keywords=computer}
2. [1] TestProduct{id=PROD002, available=true, keywords=peripheral, brand=TechCorp, stock=200, name=Mouse, price=25.5, rating=4, supplier=TechSupply, category=accessories}
3. [1] TestProduct{id=PROD003, price=75, available=false, rating=3.5, stock=0, name=Keyboard, category=accessories, keywords=typing, brand=KeyTech, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, category=electronics, price=299.99, rating=4.8, keywords=display, brand=ScreenPro, stock=30, name=Monitor, available=true, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, keywords=obsolete, supplier=OldSupply, price=8.5, available=false, rating=2, brand=OldTech, stock=0}
6. [1] TestProduct{id=PROD006, name=Headphones, category=audio, rating=4.6, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply, price=150, available=true}
7. [1] TestProduct{id=PROD007, category=electronics, price=89.99, available=true, rating=3.8, brand=CamTech, stock=25, name=Webcam, keywords=video, supplier=CamSupply}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}
2. [1] TestPerson{id=P002, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, age=35, score=9.2, name=Bob}
3. [1] TestPerson{id=P003, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1, age=16}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, status=active, department=marketing, level=7, salary=85000, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales, level=3, score=8, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, tags=test, level=1, age=0, salary=-5000, status=active, department=qa}
7. [1] TestPerson{id=P007, tags=executive, level=9, name=Grace, status=active, department=management, age=65, salary=95000, active=true, score=10}
8. [1] TestPerson{id=P008, salary=25000, department=support, level=1, name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18}
9. [1] TestPerson{id=P009, name=Ivy, active=true, score=8.7, tags=senior, status=active, level=6, age=40, salary=68000, department=engineering}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, department=intern, name=X, active=true, status=active, level=1}

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

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, status=pending, priority=normal, product_id=PROD001, total=1999.98, date=2024-01-15, discount=50, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, discount=0, amount=1, total=25.5, priority=low, region=south}
3. [1] TestOrder{id=O003, region=north, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, discount=15, date=2024-02-01, priority=high}
4. [1] TestOrder{id=O004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004}
5. [1] TestOrder{id=O005, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, date=2024-02-15, status=cancelled, priority=low, discount=0, customer_id=P005, total=999.98, region=west}
7. [1] TestOrder{id=O007, date=2024-03-01, status=shipped, discount=50, region=north, amount=4, priority=urgent, customer_id=P007, product_id=PROD006, total=600}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, discount=0, total=255, date=2024-03-05, status=pending, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, product_id=PROD007}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, customer_id=P006, total=75000, discount=0, region=east}

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

1. [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2, name=Alice, age=25}
2. [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, level=5, age=35, score=9.2, name=Bob, salary=75000, active=true}
3. [1] TestPerson{id=P003, age=16, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, status=active, department=marketing, level=7, salary=85000, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales, level=3, score=8, status=inactive}
6. [1] TestPerson{id=P006, level=1, age=0, salary=-5000, status=active, department=qa, name=Frank, active=true, score=0, tags=test}
7. [1] TestPerson{id=P007, status=active, department=management, age=65, salary=95000, active=true, score=10, tags=executive, level=9, name=Grace}
8. [1] TestPerson{id=P008, salary=25000, department=support, level=1, name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18}
9. [1] TestPerson{id=P009, name=Ivy, active=true, score=8.7, tags=senior, status=active, level=6, age=40, salary=68000, department=engineering}
10. [1] TestPerson{id=P010, name=X, active=true, status=active, level=1, age=22, salary=28000, score=6.5, tags=temp, department=intern}
11. [1] TestOrder{id=O001, date=2024-01-15, discount=50, region=north, customer_id=P001, amount=2, status=pending, priority=normal, product_id=PROD001, total=1999.98}
12. [1] TestOrder{id=O002, status=confirmed, discount=0, amount=1, total=25.5, priority=low, region=south, customer_id=P002, product_id=PROD002, date=2024-01-20}
13. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, discount=15, date=2024-02-01, priority=high, region=north}
14. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east}
15. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}
16. [1] TestOrder{id=O006, date=2024-02-15, status=cancelled, priority=low, discount=0, customer_id=P005, total=999.98, region=west, product_id=PROD005, amount=2}
17. [1] TestOrder{id=O007, priority=urgent, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50, region=north, amount=4}
18. [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, region=south, customer_id=P010, product_id=PROD002, amount=10, priority=normal, discount=0}
19. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, product_id=PROD007}
20. [1] TestOrder{id=O010, discount=0, region=east, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, customer_id=P006, total=75000}

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

1. [1] TestPerson{id=P001, score=8.5, tags=junior, status=active, department=sales, active=true, level=2, name=Alice, age=25, salary=45000}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, age=35, score=9.2}
3. [1] TestPerson{id=P003, age=16, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, tags=manager, name=Diana, age=45, active=true, score=7.8, status=active, department=marketing, level=7, salary=85000}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales, level=3, score=8, status=inactive}
6. [1] TestPerson{id=P006, department=qa, name=Frank, active=true, score=0, tags=test, level=1, age=0, salary=-5000, status=active}
7. [1] TestPerson{id=P007, active=true, score=10, tags=executive, level=9, name=Grace, status=active, department=management, age=65, salary=95000}
8. [1] TestPerson{id=P008, active=false, score=5.5, tags=junior, status=inactive, age=18, salary=25000, department=support, level=1, name=Henry}
9. [1] TestPerson{id=P009, level=6, age=40, salary=68000, department=engineering, name=Ivy, active=true, score=8.7, tags=senior, status=active}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, department=intern, name=X, active=true, status=active, level=1}

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

1. [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2, name=Alice}
2. [1] TestPerson{id=P002, age=35, score=9.2, name=Bob, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, active=false, tags=intern, status=inactive, level=1, age=16, salary=0, score=6, department=hr}
4. [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, salary=85000, tags=manager, name=Diana, age=45}
5. [1] TestPerson{id=P005, active=false, tags=employee, department=sales, level=3, score=8, status=inactive, name=Eve, age=30, salary=55000}
6. [1] TestPerson{id=P006, age=0, salary=-5000, status=active, department=qa, name=Frank, active=true, score=0, tags=test, level=1}
7. [1] TestPerson{id=P007, name=Grace, status=active, department=management, age=65, salary=95000, active=true, score=10, tags=executive, level=9}
8. [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, age=18, salary=25000, department=support, level=1, name=Henry, active=false}
9. [1] TestPerson{id=P009, age=40, salary=68000, department=engineering, name=Ivy, active=true, score=8.7, tags=senior, status=active, level=6}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, department=intern, name=X, active=true, status=active, level=1}
11. [1] TestOrder{id=O001, customer_id=P001, amount=2, status=pending, priority=normal, product_id=PROD001, total=1999.98, date=2024-01-15, discount=50, region=north}
12. [1] TestOrder{id=O002, region=south, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, discount=0, amount=1, total=25.5, priority=low}
13. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, discount=15, date=2024-02-01, priority=high, region=north}
14. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100}
16. [1] TestOrder{id=O006, product_id=PROD005, amount=2, date=2024-02-15, status=cancelled, priority=low, discount=0, customer_id=P005, total=999.98, region=west}
17. [1] TestOrder{id=O007, amount=4, priority=urgent, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50, region=north}
18. [1] TestOrder{id=O008, status=pending, region=south, customer_id=P010, product_id=PROD002, amount=10, priority=normal, discount=0, total=255, date=2024-03-05}
19. [1] TestOrder{id=O009, discount=10, region=north, product_id=PROD007, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}
20. [1] TestOrder{id=O010, customer_id=P006, total=75000, discount=0, region=east, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}
2. [1] TestPerson{id=P002, age=35, score=9.2, name=Bob, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, status=inactive, level=1, age=16, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern}
4. [1] TestPerson{id=P004, tags=manager, name=Diana, age=45, active=true, score=7.8, status=active, department=marketing, level=7, salary=85000}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, tags=employee, department=sales, level=3, score=8, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, tags=test, level=1, age=0, salary=-5000, status=active, department=qa}
7. [1] TestPerson{id=P007, department=management, age=65, salary=95000, active=true, score=10, tags=executive, level=9, name=Grace, status=active}
8. [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, salary=25000, department=support, level=1}
9. [1] TestPerson{id=P009, department=engineering, name=Ivy, active=true, score=8.7, tags=senior, status=active, level=6, age=40, salary=68000}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, department=intern, name=X, active=true, status=active, level=1}

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

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, status=pending, priority=normal, product_id=PROD001, total=1999.98, date=2024-01-15, discount=50, region=north}
2. [1] TestOrder{id=O002, amount=1, total=25.5, priority=low, region=south, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, discount=0}
3. [1] TestOrder{id=O003, discount=15, date=2024-02-01, priority=high, region=north, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east, customer_id=P004}
5. [1] TestOrder{id=O005, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, date=2024-02-15, status=cancelled, priority=low, discount=0, customer_id=P005, total=999.98, region=west}
7. [1] TestOrder{id=O007, amount=4, priority=urgent, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50, region=north}
8. [1] TestOrder{id=O008, status=pending, region=south, customer_id=P010, product_id=PROD002, amount=10, priority=normal, discount=0, total=255, date=2024-03-05}
9. [1] TestOrder{id=O009, priority=low, discount=10, region=north, product_id=PROD007, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed}
10. [1] TestOrder{id=O010, status=refunded, priority=urgent, customer_id=P006, total=75000, discount=0, region=east, product_id=PROD001, amount=1, date=2024-03-15}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, age=35, score=9.2}
3. [1] TestPerson{id=P003, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1, age=16, salary=0}
4. [1] TestPerson{id=P004, tags=manager, name=Diana, age=45, active=true, score=7.8, status=active, department=marketing, level=7, salary=85000}
5. [1] TestPerson{id=P005, active=false, tags=employee, department=sales, level=3, score=8, status=inactive, name=Eve, age=30, salary=55000}
6. [1] TestPerson{id=P006, level=1, age=0, salary=-5000, status=active, department=qa, name=Frank, active=true, score=0, tags=test}
7. [1] TestPerson{id=P007, name=Grace, status=active, department=management, age=65, salary=95000, active=true, score=10, tags=executive, level=9}
8. [1] TestPerson{id=P008, status=inactive, age=18, salary=25000, department=support, level=1, name=Henry, active=false, score=5.5, tags=junior}
9. [1] TestPerson{id=P009, tags=senior, status=active, level=6, age=40, salary=68000, department=engineering, name=Ivy, active=true, score=8.7}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, department=intern, name=X, active=true, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 4 (21.1%)
- **Tokens générés**: 34
- **Faits traités**: 27
