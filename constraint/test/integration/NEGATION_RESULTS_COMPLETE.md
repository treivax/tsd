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

1. [1] TestPerson{id=P001, name=Alice, tags=junior, department=sales, age=25, salary=45000, active=true, score=8.5, status=active, level=2}
2. [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, department=engineering, salary=75000, status=active, level=5, name=Bob, age=35}
3. [1] TestPerson{id=P003, tags=intern, department=hr, name=Charlie, salary=0, active=false, score=6, status=inactive, level=1, age=16}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager, level=7}
5. [1] TestPerson{id=P005, level=3, score=8, tags=employee, status=inactive, name=Eve, age=30, salary=55000, active=false, department=sales}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, score=0, status=active, department=qa, age=0, active=true, tags=test, level=1}
7. [1] TestPerson{id=P007, active=true, score=10, level=9, age=65, salary=95000, tags=executive, status=active, department=management, name=Grace}
8. [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1, active=false, score=5.5, name=Henry}
9. [1] TestPerson{id=P009, name=Ivy, score=8.7, status=active, department=engineering, age=40, salary=68000, active=true, tags=senior, level=6}
10. [1] TestPerson{id=P010, status=active, level=1, age=22, department=intern, name=X, salary=28000, active=true, score=6.5, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, tags=senior, level=6, name=Ivy, score=8.7, status=active, department=engineering, age=40}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, status=active, level=5, name=Bob, age=35, active=true, score=9.2, tags=senior, department=engineering, salary=75000}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, status=active, level=1, age=22, department=intern}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice, tags=junior, department=sales}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, status=inactive, level=1, age=16, tags=intern, department=hr, name=Charlie}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager, level=7}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, score=8, tags=employee, status=inactive, name=Eve, age=30, salary=55000, active=false, department=sales, level=3}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, status=active, department=management, name=Grace, active=true, score=10, level=9, age=65, salary=95000, tags=executive}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1, active=false, score=5.5, name=Henry}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, status=pending, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, discount=50, region=north}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, discount=0, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low}
3. [1] TestOrder{id=O003, customer_id=P001, date=2024-02-01, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north}
4. [1] TestOrder{id=O004, product_id=PROD004, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, amount=1, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south}
6. [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, date=2024-02-15, region=west, customer_id=P005, product_id=PROD005, amount=2, total=999.98}
7. [1] TestOrder{id=O007, priority=urgent, region=north, product_id=PROD006, amount=4, date=2024-03-01, discount=50, customer_id=P007, total=600, status=shipped}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, status=pending, amount=10, priority=normal, discount=0, region=south}
9. [1] TestOrder{id=O009, product_id=PROD007, total=89.99, status=completed, priority=low, discount=10, customer_id=P001, amount=1, date=2024-03-10, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, status=refunded, priority=urgent, region=east, amount=1, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O007, total=600, status=shipped, priority=urgent, region=north, product_id=PROD006, amount=4, date=2024-03-01, discount=50, customer_id=P007}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O008, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, status=pending, amount=10, priority=normal}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O009, priority=low, discount=10, customer_id=P001, amount=1, date=2024-03-10, region=north, product_id=PROD007, total=89.99, status=completed}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, total=75000, date=2024-03-15, status=refunded, priority=urgent, region=east, amount=1, discount=0, customer_id=P006}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, amount=1, total=299.99, date=2024-02-05}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, total=1999.98, status=pending}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, product_id=PROD002, amount=1, total=25.5, discount=0, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, date=2024-02-01}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, tags=junior, department=sales, age=25, salary=45000, active=true, score=8.5, status=active, level=2}
2. [1] TestPerson{id=P002, salary=75000, status=active, level=5, name=Bob, age=35, active=true, score=9.2, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, age=16, tags=intern, department=hr, name=Charlie, salary=0, active=false, score=6, status=inactive, level=1}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager, level=7}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, department=sales, level=3, score=8, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, score=0, status=active, department=qa, age=0, active=true, tags=test, level=1}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, level=9, age=65, salary=95000, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, active=false, score=5.5, name=Henry, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, score=8.7, status=active, department=engineering, age=40, salary=68000, active=true, tags=senior, level=6}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, status=active, level=1, age=22, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice, tags=junior, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, department=engineering, salary=75000, status=active, level=5, name=Bob, age=35, active=true, score=9.2}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager, level=7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, active=false, department=sales, level=3, score=8, tags=employee, status=inactive, name=Eve, age=30, salary=55000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, level=9, age=65, salary=95000, tags=executive, status=active, department=management}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, status=active, department=engineering, age=40, salary=68000, active=true, tags=senior, level=6, name=Ivy}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, total=1999.98, status=pending}
2. [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, product_id=PROD002, amount=1, total=25.5, discount=0, region=south}
3. [1] TestOrder{id=O003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, date=2024-02-01, product_id=PROD003}
4. [1] TestOrder{id=O004, product_id=PROD004, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, amount=1, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0, date=2024-02-15, region=west}
7. [1] TestOrder{id=O007, customer_id=P007, total=600, status=shipped, priority=urgent, region=north, product_id=PROD006, amount=4, date=2024-03-01, discount=50}
8. [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, amount=10, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, date=2024-03-10, region=north, product_id=PROD007, total=89.99, status=completed, priority=low, discount=10}
10. [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, region=east, amount=1, discount=0, customer_id=P006, product_id=PROD001, total=75000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, total=1999.98, status=pending}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10, priority=high}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O006, date=2024-02-15, region=west, customer_id=P005, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, date=2024-03-01, discount=50, customer_id=P007, total=600, status=shipped, priority=urgent, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, total=89.99, status=completed, priority=low, discount=10, customer_id=P001, amount=1, date=2024-03-10, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, discount=0, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, date=2024-02-01, product_id=PROD003}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, product_id=PROD004, status=delivered, priority=normal, discount=0, region=east}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, date=2024-03-05, status=pending, amount=10, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, available=true, rating=4.5, keywords=computer, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, brand=TechCorp, stock=50}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, available=true, brand=TechCorp, price=25.5, rating=4, keywords=peripheral, stock=200, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, price=75, rating=3.5, brand=KeyTech, supplier=KeySupply, name=Keyboard, category=accessories, available=false, keywords=typing, stock=0}
4. [1] TestProduct{id=PROD004, name=Monitor, category=electronics, price=299.99, available=true, rating=4.8, keywords=display, brand=ScreenPro, stock=30, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, price=8.5, available=false, rating=2, supplier=OldSupply, name=OldKeyboard, category=accessories, keywords=obsolete, brand=OldTech, stock=0}
6. [1] TestProduct{id=PROD006, category=audio, available=true, rating=4.6, supplier=AudioSupply, name=Headphones, price=150, keywords=sound, brand=AudioMax, stock=75}
7. [1] TestProduct{id=PROD007, rating=3.8, keywords=video, brand=CamTech, name=Webcam, category=electronics, price=89.99, available=true, stock=25, supplier=CamSupply}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, price=999.99, brand=TechCorp, stock=50, available=true, rating=4.5, keywords=computer, supplier=TechSupply, name=Laptop, category=electronics}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, price=25.5, rating=4, keywords=peripheral, stock=200, supplier=TechSupply, name=Mouse, category=accessories, available=true, brand=TechCorp}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, available=false, keywords=typing, stock=0, price=75, rating=3.5, brand=KeyTech, supplier=KeySupply}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, available=true, rating=4.8, keywords=display, brand=ScreenPro, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics, price=299.99}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, category=audio, available=true, rating=4.6, supplier=AudioSupply, name=Headphones, price=150, keywords=sound, brand=AudioMax, stock=75}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, keywords=video, brand=CamTech, name=Webcam, category=electronics, price=89.99, available=true, stock=25, supplier=CamSupply, rating=3.8}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice, tags=junior, department=sales}
2. [1] TestPerson{id=P002, age=35, active=true, score=9.2, tags=senior, department=engineering, salary=75000, status=active, level=5, name=Bob}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, status=inactive, level=1, age=16, tags=intern, department=hr}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager, level=7}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, department=sales, level=3, score=8, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active, department=qa, age=0}
7. [1] TestPerson{id=P007, salary=95000, tags=executive, status=active, department=management, name=Grace, active=true, score=10, level=9, age=65}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1, active=false, score=5.5}
9. [1] TestPerson{id=P009, department=engineering, age=40, salary=68000, active=true, tags=senior, level=6, name=Ivy, score=8.7, status=active}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, status=active, level=1, age=22, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, status=inactive, level=1, age=16, tags=intern, department=hr, name=Charlie}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active, department=qa}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, total=1999.98, status=pending}
2. [1] TestOrder{id=O002, total=25.5, discount=0, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, product_id=PROD002, amount=1}
3. [1] TestOrder{id=O003, status=shipped, priority=high, discount=15, region=north, customer_id=P001, date=2024-02-01, product_id=PROD003, amount=3, total=225}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, product_id=PROD004, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south}
6. [1] TestOrder{id=O006, date=2024-02-15, region=west, customer_id=P005, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0}
7. [1] TestOrder{id=O007, customer_id=P007, total=600, status=shipped, priority=urgent, region=north, product_id=PROD006, amount=4, date=2024-03-01, discount=50}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, status=pending, amount=10, priority=normal, discount=0, region=south}
9. [1] TestOrder{id=O009, product_id=PROD007, total=89.99, status=completed, priority=low, discount=10, customer_id=P001, amount=1, date=2024-03-10, region=north}
10. [1] TestOrder{id=O010, priority=urgent, region=east, amount=1, discount=0, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, status=refunded}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, region=north, customer_id=P001, total=1999.98, status=pending, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, discount=50}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, date=2024-01-20, status=confirmed, priority=low, product_id=PROD002, amount=1, total=25.5, discount=0, region=south, customer_id=P002}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, date=2024-02-01}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, date=2024-03-01, discount=50, customer_id=P007, total=600, status=shipped, priority=urgent, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, status=pending}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O010, status=refunded, priority=urgent, region=east, amount=1, discount=0, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, region=east, customer_id=P004, amount=1, total=299.99, date=2024-02-05, product_id=PROD004, status=delivered, priority=normal, discount=0}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O006, discount=0, date=2024-02-15, region=west, customer_id=P005, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, date=2024-03-10, region=north, product_id=PROD007, total=89.99, status=completed, priority=low, discount=10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, department=sales, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice}
2. [1] TestPerson{id=P002, name=Bob, age=35, active=true, score=9.2, tags=senior, department=engineering, salary=75000, status=active, level=5}
3. [1] TestPerson{id=P003, age=16, tags=intern, department=hr, name=Charlie, salary=0, active=false, score=6, status=inactive, level=1}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager, level=7}
5. [1] TestPerson{id=P005, salary=55000, active=false, department=sales, level=3, score=8, tags=employee, status=inactive, name=Eve, age=30}
6. [1] TestPerson{id=P006, department=qa, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, level=9, age=65, salary=95000, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, score=5.5, name=Henry, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1, active=false}
9. [1] TestPerson{id=P009, name=Ivy, score=8.7, status=active, department=engineering, age=40, salary=68000, active=true, tags=senior, level=6}
10. [1] TestPerson{id=P010, age=22, department=intern, name=X, salary=28000, active=true, score=6.5, tags=temp, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, active=false, department=sales, level=3, score=8, tags=employee, status=inactive, name=Eve, age=30, salary=55000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, name=Henry, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1, active=false}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, status=inactive, level=1, age=16, tags=intern, department=hr}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, status=pending, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, discount=50, region=north}
2. [1] TestOrder{id=O002, status=confirmed, priority=low, product_id=PROD002, amount=1, total=25.5, discount=0, region=south, customer_id=P002, date=2024-01-20}
3. [1] TestOrder{id=O003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, date=2024-02-01, product_id=PROD003}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, product_id=PROD004, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0, date=2024-02-15, region=west}
7. [1] TestOrder{id=O007, date=2024-03-01, discount=50, customer_id=P007, total=600, status=shipped, priority=urgent, region=north, product_id=PROD006, amount=4}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, status=pending, amount=10, priority=normal, discount=0, region=south}
9. [1] TestOrder{id=O009, product_id=PROD007, total=89.99, status=completed, priority=low, discount=10, customer_id=P001, amount=1, date=2024-03-10, region=north}
10. [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, region=east, amount=1, discount=0, customer_id=P006, product_id=PROD001, total=75000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O006, date=2024-02-15, region=west, customer_id=P005, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, total=600, status=shipped, priority=urgent, region=north, product_id=PROD006, amount=4, date=2024-03-01, discount=50}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, date=2024-03-10, region=north, product_id=PROD007, total=89.99, status=completed, priority=low, discount=10}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O010, amount=1, discount=0, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, status=refunded, priority=urgent, region=east}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, discount=0, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O003, discount=15, region=north, customer_id=P001, date=2024-02-01, product_id=PROD003, amount=3, total=225, status=shipped, priority=high}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, amount=1, total=299.99, date=2024-02-05}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice, tags=junior}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, department=engineering, salary=75000, status=active, level=5, name=Bob, age=35, active=true}
3. [1] TestPerson{id=P003, level=1, age=16, tags=intern, department=hr, name=Charlie, salary=0, active=false, score=6, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager, level=7}
5. [1] TestPerson{id=P005, active=false, department=sales, level=3, score=8, tags=employee, status=inactive, name=Eve, age=30, salary=55000}
6. [1] TestPerson{id=P006, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active, department=qa}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, level=9, age=65, salary=95000, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, active=false, score=5.5, name=Henry, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, score=8.7, status=active, department=engineering, age=40, salary=68000, active=true, tags=senior, level=6}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, status=active, level=1, age=22, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, status=active, level=2, name=Alice, tags=junior, department=sales, age=25, salary=45000, active=true}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, department=engineering, salary=75000, status=active, level=5, name=Bob, age=35}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, level=1, age=16, tags=intern, department=hr, name=Charlie, salary=0, active=false, score=6}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager, level=7, name=Diana}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, department=sales, level=3, score=8, tags=employee, status=inactive}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, age=40, salary=68000, active=true, tags=senior, level=6, name=Ivy, score=8.7, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, name=Grace, active=true, score=10, level=9}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1, active=false, score=5.5}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, available=true, rating=4.5, keywords=computer, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, brand=TechCorp, stock=50}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, available=true, brand=TechCorp, price=25.5, rating=4, keywords=peripheral, stock=200, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, price=75, rating=3.5, brand=KeyTech, supplier=KeySupply, name=Keyboard, category=accessories, available=false, keywords=typing, stock=0}
4. [1] TestProduct{id=PROD004, available=true, rating=4.8, keywords=display, brand=ScreenPro, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics, price=299.99}
5. [1] TestProduct{id=PROD005, price=8.5, available=false, rating=2, supplier=OldSupply, name=OldKeyboard, category=accessories, keywords=obsolete, brand=OldTech, stock=0}
6. [1] TestProduct{id=PROD006, category=audio, available=true, rating=4.6, supplier=AudioSupply, name=Headphones, price=150, keywords=sound, brand=AudioMax, stock=75}
7. [1] TestProduct{id=PROD007, rating=3.8, keywords=video, brand=CamTech, name=Webcam, category=electronics, price=89.99, available=true, stock=25, supplier=CamSupply}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD007, available=true, stock=25, supplier=CamSupply, rating=3.8, keywords=video, brand=CamTech, name=Webcam, category=electronics, price=89.99}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD001, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, brand=TechCorp, stock=50, available=true, rating=4.5, keywords=computer}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, category=accessories, available=true, brand=TechCorp, price=25.5, rating=4, keywords=peripheral, stock=200, supplier=TechSupply}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD003, price=75, rating=3.5, brand=KeyTech, supplier=KeySupply, name=Keyboard, category=accessories, available=false, keywords=typing, stock=0}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD004, rating=4.8, keywords=display, brand=ScreenPro, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics, price=299.99, available=true}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD006, category=audio, available=true, rating=4.6, supplier=AudioSupply, name=Headphones, price=150, keywords=sound, brand=AudioMax, stock=75}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice, tags=junior, department=sales}
2. [1] TestPerson{id=P002, status=active, level=5, name=Bob, age=35, active=true, score=9.2, tags=senior, department=engineering, salary=75000}
3. [1] TestPerson{id=P003, age=16, tags=intern, department=hr, name=Charlie, salary=0, active=false, score=6, status=inactive, level=1}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager, level=7}
5. [1] TestPerson{id=P005, tags=employee, status=inactive, name=Eve, age=30, salary=55000, active=false, department=sales, level=3, score=8}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, department=qa, age=0, active=true, tags=test, level=1, name=Frank}
7. [1] TestPerson{id=P007, department=management, name=Grace, active=true, score=10, level=9, age=65, salary=95000, tags=executive, status=active}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1, active=false, score=5.5}
9. [1] TestPerson{id=P009, age=40, salary=68000, active=true, tags=senior, level=6, name=Ivy, score=8.7, status=active, department=engineering}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, status=active, level=1, age=22, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, status=inactive, level=1, age=16, tags=intern, department=hr, name=Charlie, salary=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, department=sales, level=3, score=8, tags=employee, status=inactive}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, level=1, active=false, score=5.5, name=Henry, age=18, salary=25000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, age=45, tags=manager, level=7, name=Diana, salary=85000, active=true, score=7.8, status=active}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active, department=qa}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, name=Grace, active=true, score=10, level=9}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, score=8.7, status=active, department=engineering, age=40, salary=68000, active=true, tags=senior, level=6}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice, tags=junior, department=sales}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, department=engineering, salary=75000, status=active, level=5, name=Bob, age=35, active=true, score=9.2}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, status=pending, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, discount=50, region=north}
2. [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, product_id=PROD002, amount=1, total=25.5}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, date=2024-02-01}
4. [1] TestOrder{id=O004, product_id=PROD004, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, amount=1, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0, date=2024-02-15, region=west}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, date=2024-03-01, discount=50, customer_id=P007, total=600, status=shipped, priority=urgent, region=north}
8. [1] TestOrder{id=O008, amount=10, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, status=pending}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, date=2024-03-10, region=north, product_id=PROD007, total=89.99, status=completed, priority=low, discount=10}
10. [1] TestOrder{id=O010, region=east, amount=1, discount=0, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, status=refunded, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, date=2024-02-01, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, amount=1, total=299.99, date=2024-02-05}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, date=2024-03-01, discount=50, customer_id=P007, total=600, status=shipped, priority=urgent, region=north, product_id=PROD006}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, status=pending, amount=10, priority=normal, discount=0, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, date=2024-03-10, region=north, product_id=PROD007, total=89.99, status=completed, priority=low, discount=10}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, total=1999.98, status=pending}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, product_id=PROD002, amount=1, total=25.5, discount=0, region=south}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice, tags=junior, department=sales, age=25}
2. [1] TestPerson{id=P002, name=Bob, age=35, active=true, score=9.2, tags=senior, department=engineering, salary=75000, status=active, level=5}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, status=inactive, level=1, age=16, tags=intern, department=hr}
4. [1] TestPerson{id=P004, tags=manager, level=7, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45}
5. [1] TestPerson{id=P005, active=false, department=sales, level=3, score=8, tags=employee, status=inactive, name=Eve, age=30, salary=55000}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, score=0, status=active, department=qa, age=0, active=true, tags=test, level=1}
7. [1] TestPerson{id=P007, level=9, age=65, salary=95000, tags=executive, status=active, department=management, name=Grace, active=true, score=10}
8. [1] TestPerson{id=P008, score=5.5, name=Henry, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1, active=false}
9. [1] TestPerson{id=P009, name=Ivy, score=8.7, status=active, department=engineering, age=40, salary=68000, active=true, tags=senior, level=6}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, status=active, level=1, age=22, department=intern}
11. [1] TestOrder{id=O001, total=1999.98, status=pending, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001}
12. [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, product_id=PROD002, amount=1, total=25.5}
13. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, date=2024-02-01}
14. [1] TestOrder{id=O004, product_id=PROD004, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, amount=1, total=299.99, date=2024-02-05}
15. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south}
16. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0, date=2024-02-15, region=west}
17. [1] TestOrder{id=O007, customer_id=P007, total=600, status=shipped, priority=urgent, region=north, product_id=PROD006, amount=4, date=2024-03-01, discount=50}
18. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, status=pending, amount=10, priority=normal, discount=0, region=south}
19. [1] TestOrder{id=O009, customer_id=P001, amount=1, date=2024-03-10, region=north, product_id=PROD007, total=89.99, status=completed, priority=low, discount=10}
20. [1] TestOrder{id=O010, priority=urgent, region=east, amount=1, discount=0, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, status=refunded}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, active=true, score=9.2, tags=senior, department=engineering, salary=75000, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, level=9, age=65, salary=95000, tags=executive, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, total=600, status=shipped, priority=urgent, region=north, product_id=PROD006, amount=4, date=2024-03-01, discount=50}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, status=active, level=2, name=Alice, tags=junior, department=sales, age=25, salary=45000}
   - Fait 2: [1] TestOrder{id=O009, amount=1, date=2024-03-10, region=north, product_id=PROD007, total=89.99, status=completed, priority=low, discount=10, customer_id=P001}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, active=true, score=9.2, tags=senior, department=engineering, salary=75000, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, discount=0, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice, tags=junior}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, date=2024-02-01}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, product_id=PROD004, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, amount=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, department=sales, level=3, score=8, tags=employee, status=inactive}
   - Fait 2: [1] TestOrder{id=O006, date=2024-02-15, region=west, customer_id=P005, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, status=active, level=1, age=22, department=intern}
   - Fait 2: [1] TestOrder{id=O008, amount=10, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, status=pending}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active, department=qa}
   - Fait 2: [1] TestOrder{id=O010, amount=1, discount=0, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, status=refunded, priority=urgent, region=east}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, total=1999.98, status=pending}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/20 (50.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, name=Alice, tags=junior, department=sales, age=25, salary=45000, active=true, score=8.5, status=active}
2. [1] TestPerson{id=P002, name=Bob, age=35, active=true, score=9.2, tags=senior, department=engineering, salary=75000, status=active, level=5}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, status=inactive, level=1, age=16, tags=intern, department=hr}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager, level=7}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, department=sales, level=3, score=8, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, department=qa, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, level=9, age=65, salary=95000, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1, active=false, score=5.5}
9. [1] TestPerson{id=P009, age=40, salary=68000, active=true, tags=senior, level=6, name=Ivy, score=8.7, status=active, department=engineering}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, status=active, level=1, age=22, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, status=active, department=management, name=Grace, active=true, score=10, level=9, age=65, salary=95000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, status=active, department=engineering, age=40, salary=68000, active=true, tags=senior, level=6, name=Ivy}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, status=active, level=1, age=22, department=intern}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice, tags=junior, department=sales}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, active=true, score=9.2, tags=senior, department=engineering, salary=75000, status=active, level=5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active, department=qa}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice, tags=junior, department=sales}
2. [1] TestPerson{id=P002, salary=75000, status=active, level=5, name=Bob, age=35, active=true, score=9.2, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, status=inactive, level=1, age=16, tags=intern, department=hr, name=Charlie, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, age=45, tags=manager, level=7, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing}
5. [1] TestPerson{id=P005, score=8, tags=employee, status=inactive, name=Eve, age=30, salary=55000, active=false, department=sales, level=3}
6. [1] TestPerson{id=P006, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active, department=qa}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, name=Grace, active=true, score=10, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1, active=false, score=5.5}
9. [1] TestPerson{id=P009, age=40, salary=68000, active=true, tags=senior, level=6, name=Ivy, score=8.7, status=active, department=engineering}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, tags=temp, status=active, level=1, age=22, department=intern}
11. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, status=pending, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, discount=50, region=north}
12. [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, product_id=PROD002, amount=1, total=25.5, discount=0, region=south}
13. [1] TestOrder{id=O003, customer_id=P001, date=2024-02-01, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north}
14. [1] TestOrder{id=O004, region=east, customer_id=P004, amount=1, total=299.99, date=2024-02-05, product_id=PROD004, status=delivered, priority=normal, discount=0}
15. [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed}
16. [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, date=2024-02-15, region=west, customer_id=P005, product_id=PROD005, amount=2, total=999.98}
17. [1] TestOrder{id=O007, customer_id=P007, total=600, status=shipped, priority=urgent, region=north, product_id=PROD006, amount=4, date=2024-03-01, discount=50}
18. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, status=pending, amount=10, priority=normal, discount=0, region=south}
19. [1] TestOrder{id=O009, discount=10, customer_id=P001, amount=1, date=2024-03-10, region=north, product_id=PROD007, total=89.99, status=completed, priority=low}
20. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, status=refunded, priority=urgent, region=east, amount=1, discount=0}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P006, score=0, status=active, department=qa, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, total=75000, date=2024-03-15, status=refunded, priority=urgent, region=east, amount=1, discount=0, customer_id=P006}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, active=true, score=9.2, tags=senior, department=engineering, salary=75000, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O002, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, product_id=PROD002, amount=1, total=25.5, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager, level=7, name=Diana}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, product_id=PROD004, status=delivered, priority=normal, discount=0, region=east}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, status=active, level=5, name=Bob, age=35, active=true, score=9.2, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, name=Eve, age=30, salary=55000, active=false, department=sales, level=3, score=8, tags=employee}
   - Fait 2: [1] TestOrder{id=O006, region=west, customer_id=P005, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0, date=2024-02-15}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, status=active, level=1, age=22, department=intern, name=X, salary=28000}
   - Fait 2: [1] TestOrder{id=O008, amount=10, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, status=pending}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, tags=junior, department=sales, age=25, salary=45000, active=true, score=8.5, status=active}
   - Fait 2: [1] TestOrder{id=O009, amount=1, date=2024-03-10, region=north, product_id=PROD007, total=89.99, status=completed, priority=low, discount=10, customer_id=P001}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, total=1999.98, status=pending}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, department=sales, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice}
   - Fait 2: [1] TestOrder{id=O003, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, date=2024-02-01, product_id=PROD003, amount=3}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, status=active, department=management, name=Grace, active=true, score=10, level=9, age=65, salary=95000}
   - Fait 2: [1] TestOrder{id=O007, amount=4, date=2024-03-01, discount=50, customer_id=P007, total=600, status=shipped, priority=urgent, region=north, product_id=PROD006}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/20 (50.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, department=sales, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice}
2. [1] TestPerson{id=P002, name=Bob, age=35, active=true, score=9.2, tags=senior, department=engineering, salary=75000, status=active, level=5}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, status=inactive, level=1, age=16, tags=intern, department=hr}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager, level=7}
5. [1] TestPerson{id=P005, status=inactive, name=Eve, age=30, salary=55000, active=false, department=sales, level=3, score=8, tags=employee}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, score=0, status=active, department=qa, age=0, active=true, tags=test, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, name=Grace, active=true, score=10, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1, active=false, score=5.5}
9. [1] TestPerson{id=P009, name=Ivy, score=8.7, status=active, department=engineering, age=40, salary=68000, active=true, tags=senior, level=6}
10. [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, status=active, level=1, age=22, department=intern, name=X, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, tags=junior, department=sales, age=25, salary=45000, active=true, score=8.5, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, active=true, score=9.2, tags=senior, department=engineering, salary=75000, status=active, level=5}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, status=inactive, level=1, age=16, tags=intern, department=hr}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, name=Eve, age=30, salary=55000, active=false, department=sales, level=3, score=8, tags=employee}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, age=22, department=intern, name=X, salary=28000, active=true, score=6.5, tags=temp, status=active, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, name=Grace, active=true, score=10, level=9}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1, active=false, score=5.5}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, active=true, tags=senior, level=6, name=Ivy, score=8.7, status=active, department=engineering, age=40, salary=68000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, total=1999.98, status=pending}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, discount=0, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, date=2024-02-01}
4. [1] TestOrder{id=O004, product_id=PROD004, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, amount=1, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south}
6. [1] TestOrder{id=O006, region=west, customer_id=P005, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0, date=2024-02-15}
7. [1] TestOrder{id=O007, customer_id=P007, total=600, status=shipped, priority=urgent, region=north, product_id=PROD006, amount=4, date=2024-03-01, discount=50}
8. [1] TestOrder{id=O008, product_id=PROD002, total=255, date=2024-03-05, status=pending, amount=10, priority=normal, discount=0, region=south, customer_id=P010}
9. [1] TestOrder{id=O009, product_id=PROD007, total=89.99, status=completed, priority=low, discount=10, customer_id=P001, amount=1, date=2024-03-10, region=north}
10. [1] TestOrder{id=O010, priority=urgent, region=east, amount=1, discount=0, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, status=refunded}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, product_id=PROD002, amount=1, total=25.5, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, product_id=PROD004, status=delivered, priority=normal, discount=0, region=east}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, status=shipped, priority=urgent, region=north, product_id=PROD006, amount=4, date=2024-03-01, discount=50, customer_id=P007, total=600}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, total=1999.98, status=pending, product_id=PROD001, amount=2}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, region=north, customer_id=P001, date=2024-02-01, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, discount=15}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O006, priority=low, discount=0, date=2024-02-15, region=west, customer_id=P005, product_id=PROD005, amount=2, total=999.98, status=cancelled}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, status=pending, amount=10, priority=normal, discount=0, region=south}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, region=north, product_id=PROD007, total=89.99, status=completed, priority=low, discount=10, customer_id=P001, amount=1, date=2024-03-10}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, status=refunded, priority=urgent, region=east, amount=1, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice, tags=junior, department=sales}
2. [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, department=engineering, salary=75000, status=active, level=5, name=Bob, age=35}
3. [1] TestPerson{id=P003, level=1, age=16, tags=intern, department=hr, name=Charlie, salary=0, active=false, score=6, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager, level=7}
5. [1] TestPerson{id=P005, level=3, score=8, tags=employee, status=inactive, name=Eve, age=30, salary=55000, active=false, department=sales}
6. [1] TestPerson{id=P006, score=0, status=active, department=qa, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000}
7. [1] TestPerson{id=P007, active=true, score=10, level=9, age=65, salary=95000, tags=executive, status=active, department=management, name=Grace}
8. [1] TestPerson{id=P008, active=false, score=5.5, name=Henry, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1}
9. [1] TestPerson{id=P009, level=6, name=Ivy, score=8.7, status=active, department=engineering, age=40, salary=68000, active=true, tags=senior}
10. [1] TestPerson{id=P010, score=6.5, tags=temp, status=active, level=1, age=22, department=intern, name=X, salary=28000, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, status=active, level=2, name=Alice, tags=junior, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, name=Charlie, salary=0, active=false, score=6, status=inactive, level=1, age=16}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, department=sales, level=3, score=8, tags=employee, status=inactive}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, status=active, department=management, name=Grace, active=true, score=10, level=9, age=65, salary=95000, tags=executive}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, name=Henry, age=18, salary=25000, tags=junior, status=inactive, department=support, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, active=true, score=9.2, tags=senior, department=engineering, salary=75000, status=active, level=5}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, age=45, tags=manager, level=7, name=Diana}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active, department=qa}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, active=true, tags=senior, level=6, name=Ivy, score=8.7, status=active, department=engineering}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, score=6.5, tags=temp, status=active, level=1, age=22, department=intern, name=X}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 150
- **Faits traités**: 27
