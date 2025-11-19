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

1. [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2, name=Bob}
3. [1] TestPerson{id=P003, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000}
5. [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa}
7. [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management, name=Grace}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}
9. [1] TestPerson{id=P009, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true}
10. [1] TestPerson{id=P010, status=active, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr, level=1}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee, level=3, age=30}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, score=5.5, level=1, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true, status=active, level=6}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active, name=X, age=22, active=true}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8, status=active, department=marketing}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15}
2. [1] TestOrder{id=O002, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed}
3. [1] TestOrder{id=O003, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001, amount=3}
4. [1] TestOrder{id=O004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal, customer_id=P004}
5. [1] TestOrder{id=O005, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100}
6. [1] TestOrder{id=O006, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2, customer_id=P005, product_id=PROD005}
7. [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010}
9. [1] TestOrder{id=O009, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}
10. [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20, priority=low}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, total=299.99, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2}
3. [1] TestPerson{id=P003, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000}
5. [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa}
7. [1] TestPerson{id=P007, salary=95000, status=active, department=management, name=Grace, age=65, active=true, score=10, tags=executive, level=9}
8. [1] TestPerson{id=P008, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1, name=Henry}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering}
10. [1] TestPerson{id=P010, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee, level=3, age=30}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, level=9, salary=95000, status=active, department=management, name=Grace, age=65, active=true, score=10, tags=executive}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south}
6. [1] TestOrder{id=O006, amount=2, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west}
7. [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50}
8. [1] TestOrder{id=O008, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending, priority=normal}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north}
10. [1] TestOrder{id=O010, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, region=west, amount=2, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O001, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4, total=600, date=2024-03-01}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, available=true, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, rating=4.5, keywords=computer, brand=TechCorp, stock=50}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, rating=4, keywords=peripheral, brand=TechCorp, price=25.5, available=true, stock=200, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0, name=Keyboard, category=accessories, price=75, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, name=Monitor, available=true, keywords=display, brand=ScreenPro, stock=30, category=electronics, price=299.99, rating=4.8, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, price=8.5, available=false, brand=OldTech, stock=0, rating=2, keywords=obsolete, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, supplier=AudioSupply, name=Headphones, category=audio, keywords=sound, price=150, available=true, rating=4.6, brand=AudioMax, stock=75}
7. [1] TestProduct{id=PROD007, rating=3.8, brand=CamTech, supplier=CamSupply, name=Webcam, category=electronics, available=true, keywords=video, stock=25, price=89.99}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, price=75, supplier=KeySupply, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, available=true, keywords=display, brand=ScreenPro, stock=30, category=electronics, price=299.99, rating=4.8, supplier=ScreenSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD006, brand=AudioMax, stock=75, supplier=AudioSupply, name=Headphones, category=audio, keywords=sound, price=150, available=true, rating=4.6}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD007, brand=CamTech, supplier=CamSupply, name=Webcam, category=electronics, available=true, keywords=video, stock=25, price=89.99, rating=3.8}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD001, category=electronics, price=999.99, rating=4.5, keywords=computer, brand=TechCorp, stock=50, available=true, supplier=TechSupply, name=Laptop}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD002, keywords=peripheral, brand=TechCorp, price=25.5, available=true, stock=200, supplier=TechSupply, name=Mouse, category=accessories, rating=4}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales}
2. [1] TestPerson{id=P002, score=9.2, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive, name=Charlie}
4. [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}
6. [1] TestPerson{id=P006, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1}
7. [1] TestPerson{id=P007, salary=95000, status=active, department=management, name=Grace, age=65, active=true, score=10, tags=executive, level=9}
8. [1] TestPerson{id=P008, department=support, salary=25000, score=5.5, level=1, name=Henry, age=18, active=false, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering, name=Ivy}
10. [1] TestPerson{id=P010, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, discount=50, region=north, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal}
2. [1] TestOrder{id=O002, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1}
3. [1] TestOrder{id=O003, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001, amount=3, date=2024-02-01}
4. [1] TestOrder{id=O004, product_id=PROD004, total=299.99, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east}
5. [1] TestOrder{id=O005, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed}
6. [1] TestOrder{id=O006, amount=2, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west}
7. [1] TestOrder{id=O007, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O006, region=west, amount=2, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20, priority=low}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, priority=urgent, discount=50, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000}
2. [1] TestPerson{id=P002, status=active, department=engineering, level=5, score=9.2, name=Bob, age=35, salary=75000, active=true, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager, active=true}
5. [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}
6. [1] TestPerson{id=P006, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa, name=Frank, age=0}
7. [1] TestPerson{id=P007, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management, name=Grace, age=65}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}
9. [1] TestPerson{id=P009, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true, status=active, level=6, age=40}
10. [1] TestPerson{id=P010, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr, level=1, age=16, salary=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, department=support, salary=25000, score=5.5, level=1, name=Henry, age=18, active=false, tags=junior, status=inactive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98}
2. [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20, priority=low}
3. [1] TestOrder{id=O003, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001, amount=3}
4. [1] TestOrder{id=O004, total=299.99, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004}
5. [1] TestOrder{id=O005, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2}
7. [1] TestOrder{id=O007, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4, total=600}
8. [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}
9. [1] TestOrder{id=O009, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}
10. [1] TestOrder{id=O010, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O006, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2, customer_id=P005, product_id=PROD005}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales}
2. [1] TestPerson{id=P002, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2, name=Bob, age=35}
3. [1] TestPerson{id=P003, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000}
5. [1] TestPerson{id=P005, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa}
7. [1] TestPerson{id=P007, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}
9. [1] TestPerson{id=P009, level=6, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true, status=active}
10. [1] TestPerson{id=P010, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, level=9, salary=95000, status=active, department=management, name=Grace, age=65, active=true, score=10, tags=executive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering, name=Ivy}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2, name=Bob}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, rating=4.5, keywords=computer, brand=TechCorp, stock=50, available=true, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, rating=4, keywords=peripheral, brand=TechCorp, price=25.5, available=true, stock=200, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, stock=0, name=Keyboard, category=accessories, price=75, supplier=KeySupply, available=false, rating=3.5, keywords=typing, brand=KeyTech}
4. [1] TestProduct{id=PROD004, category=electronics, price=299.99, rating=4.8, supplier=ScreenSupply, name=Monitor, available=true, keywords=display, brand=ScreenPro, stock=30}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, price=8.5, available=false, brand=OldTech, stock=0, rating=2, keywords=obsolete, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, price=150, available=true, rating=4.6, brand=AudioMax, stock=75, supplier=AudioSupply, name=Headphones, category=audio, keywords=sound}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, available=true, keywords=video, stock=25, price=89.99, rating=3.8, brand=CamTech, supplier=CamSupply}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD004, brand=ScreenPro, stock=30, category=electronics, price=299.99, rating=4.8, supplier=ScreenSupply, name=Monitor, available=true, keywords=display}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD006, keywords=sound, price=150, available=true, rating=4.6, brand=AudioMax, stock=75, supplier=AudioSupply, name=Headphones, category=audio}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD007, category=electronics, available=true, keywords=video, stock=25, price=89.99, rating=3.8, brand=CamTech, supplier=CamSupply, name=Webcam}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, rating=4.5, keywords=computer, brand=TechCorp, stock=50, available=true, supplier=TechSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD002, stock=200, supplier=TechSupply, name=Mouse, category=accessories, rating=4, keywords=peripheral, brand=TechCorp, price=25.5, available=true}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, price=75, supplier=KeySupply, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales}
2. [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, level=5, score=9.2, name=Bob, age=35, salary=75000, active=true}
3. [1] TestPerson{id=P003, age=16, salary=0, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr, level=1}
4. [1] TestPerson{id=P004, level=7, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8, status=active, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}
6. [1] TestPerson{id=P006, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa, name=Frank}
7. [1] TestPerson{id=P007, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering}
10. [1] TestPerson{id=P010, level=1, salary=28000, score=6.5, status=active, name=X, age=22, active=true, tags=temp, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee, level=3}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, department=support, salary=25000, score=5.5, level=1, name=Henry, age=18, active=false, tags=junior}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, level=6, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, level=5, score=9.2, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8, status=active, department=marketing}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa, name=Frank, age=0}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, discount=50, region=north, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal}
2. [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20, priority=low}
3. [1] TestOrder{id=O003, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal}
5. [1] TestOrder{id=O005, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high}
6. [1] TestOrder{id=O006, amount=2, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west}
7. [1] TestOrder{id=O007, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north}
8. [1] TestOrder{id=O008, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05}
9. [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O009, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, total=299.99, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50, customer_id=P007}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales}
2. [1] TestPerson{id=P002, status=active, department=engineering, level=5, score=9.2, name=Bob, age=35, salary=75000, active=true, tags=senior}
3. [1] TestPerson{id=P003, age=16, salary=0, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr, level=1}
4. [1] TestPerson{id=P004, level=7, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8, status=active, department=marketing}
5. [1] TestPerson{id=P005, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee, level=3, age=30}
6. [1] TestPerson{id=P006, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa, name=Frank}
7. [1] TestPerson{id=P007, salary=95000, status=active, department=management, name=Grace, age=65, active=true, score=10, tags=executive, level=9}
8. [1] TestPerson{id=P008, salary=25000, score=5.5, level=1, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000}
10. [1] TestPerson{id=P010, department=intern, level=1, salary=28000, score=6.5, status=active, name=X, age=22, active=true, tags=temp}
11. [1] TestOrder{id=O001, region=north, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50}
12. [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20, priority=low}
13. [1] TestOrder{id=O003, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001}
14. [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal}
15. [1] TestOrder{id=O005, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100}
16. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2}
17. [1] TestOrder{id=O007, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4}
18. [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}
19. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north}
20. [1] TestOrder{id=O010, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, score=5.5, level=1, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support}
   - Fait 2: [1] TestOrder{id=O003, total=225, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee, level=3}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, status=active, name=X, age=22, active=true, tags=temp, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O003, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001, amount=3}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa}
   - Fait 2: [1] TestOrder{id=O004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal, customer_id=P004}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, status=active, name=X, age=22, active=true, tags=temp, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2, customer_id=P005}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50, customer_id=P007}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P003, level=1, age=16, salary=0, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr}
   - Fait 2: [1] TestOrder{id=O003, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P002, level=5, score=9.2, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O003, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001, amount=3, date=2024-02-01}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O005, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, level=1, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa}
   - Fait 2: [1] TestOrder{id=O001, region=north, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa}
   - Fait 2: [1] TestOrder{id=O002, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{id=O003, region=north, product_id=PROD003, total=225, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive}
   - Fait 2: [1] TestOrder{id=O003, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P001, status=active, level=2, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50, customer_id=P007}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true, status=active, level=6, age=40}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr, level=1, age=16}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P002, level=5, score=9.2, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, status=active, name=X, age=22, active=true, tags=temp, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2, name=Bob, age=35}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal, customer_id=P004, amount=1}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, total=299.99, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales, age=25}
   - Fait 2: [1] TestOrder{id=O005, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P009, level=6, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P008, department=support, salary=25000, score=5.5, level=1, name=Henry, age=18, active=false, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20, priority=low}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal, customer_id=P004, amount=1}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, level=5, score=9.2, name=Bob, age=35, salary=75000, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O006, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2, customer_id=P005, product_id=PROD005, total=999.98}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south, product_id=PROD002}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P005, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee}
   - Fait 2: [1] TestOrder{id=O005, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P010, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active, name=X, age=22}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P009, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true, status=active, level=6}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000}
   - Fait 2: [1] TestOrder{id=O008, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending, priority=normal}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, level=9, salary=95000, status=active, department=management, name=Grace, age=65, active=true}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8, status=active}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee, level=3, age=30}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, total=225, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P008, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, product_id=PROD001}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P008, level=1, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, status=active, name=X, age=22, active=true, tags=temp, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O004, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}
   - Fait 2: [1] TestOrder{id=O005, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa, name=Frank}
   - Fait 2: [1] TestOrder{id=O006, priority=low, discount=0, region=west, amount=2, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering, name=Ivy}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O008, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O010, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}
   - Fait 2: [1] TestOrder{id=O004, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management, name=Grace}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, status=inactive, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P007, department=management, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active}
   - Fait 2: [1] TestOrder{id=O004, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, department=support, salary=25000, score=5.5, level=1, name=Henry, age=18, active=false, tags=junior}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, discount=50, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P007, level=9, salary=95000, status=active, department=management, name=Grace, age=65, active=true, score=10, tags=executive}
   - Fait 2: [1] TestOrder{id=O005, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P004, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, status=active, name=X, age=22, active=true, tags=temp, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, discount=50, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2, name=Bob, age=35}
   - Fait 2: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P003, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr, level=1, age=16, salary=0}
   - Fait 2: [1] TestOrder{id=O009, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P003, level=1, age=16, salary=0, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr}
   - Fait 2: [1] TestOrder{id=O006, priority=low, discount=0, region=west, amount=2, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, level=9, salary=95000, status=active, department=management, name=Grace, age=65, active=true}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, level=9, salary=95000, status=active, department=management, name=Grace, age=65, active=true}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20, priority=low}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales}
2. [1] TestPerson{id=P002, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2, name=Bob, age=35}
3. [1] TestPerson{id=P003, level=1, age=16, salary=0, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr}
4. [1] TestPerson{id=P004, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000}
5. [1] TestPerson{id=P005, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false}
6. [1] TestPerson{id=P006, active=true, score=0, status=active, level=1, tags=test, department=qa, name=Frank, age=0, salary=-5000}
7. [1] TestPerson{id=P007, score=10, tags=executive, level=9, salary=95000, status=active, department=management, name=Grace, age=65, active=true}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}
9. [1] TestPerson{id=P009, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true, status=active, level=6}
10. [1] TestPerson{id=P010, salary=28000, score=6.5, status=active, name=X, age=22, active=true, tags=temp, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, status=active, department=management, name=Grace, age=65, active=true, score=10, tags=executive, level=9}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, level=1, salary=28000, score=6.5, status=active, name=X, age=22, active=true, tags=temp, department=intern}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2, name=Bob}
3. [1] TestPerson{id=P003, salary=0, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr, level=1, age=16}
4. [1] TestPerson{id=P004, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000}
5. [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee, level=3}
6. [1] TestPerson{id=P006, score=0, status=active, level=1, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true}
7. [1] TestPerson{id=P007, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management, name=Grace, age=65}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}
9. [1] TestPerson{id=P009, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true, status=active, level=6, age=40}
10. [1] TestPerson{id=P010, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active, name=X, age=22}
11. [1] TestOrder{id=O001, region=north, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50}
12. [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20, priority=low}
13. [1] TestOrder{id=O003, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high}
14. [1] TestOrder{id=O004, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05, discount=0}
15. [1] TestOrder{id=O005, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1}
16. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2}
17. [1] TestOrder{id=O007, region=north, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped}
18. [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}
19. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1}
20. [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, score=5.5, level=1, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O006, priority=low, discount=0, region=west, amount=2, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true, status=active, level=6}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4, total=600}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O009, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr, level=1, age=16, salary=0}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000}
   - Fait 2: [1] TestOrder{id=O001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, product_id=PROD001}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa, name=Frank, age=0}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, region=west, amount=2, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O001, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee, level=3, age=30}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, status=active, name=X, age=22, active=true, tags=temp, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2, name=Bob}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O002, region=south, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20, priority=low}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, total=225, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager, active=true}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2, customer_id=P005}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, amount=2, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P008, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1, name=Henry, age=18}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P005, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false}
   - Fait 2: [1] TestOrder{id=O003, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active, name=X, age=22, active=true}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4, total=600}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, level=1, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000}
   - Fait 2: [1] TestOrder{id=O008, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active, name=X, age=22, active=true}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr, level=1, age=16}
   - Fait 2: [1] TestOrder{id=O001, discount=50, region=north, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P006, status=active, level=1, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true, score=0}
   - Fait 2: [1] TestOrder{id=O004, priority=normal, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P004, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P002, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2, name=Bob, age=35, salary=75000}
   - Fait 2: [1] TestOrder{id=O005, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, discount=50, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P003, level=1, age=16, salary=0, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P009, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true, status=active, level=6}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, level=5, score=9.2, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management, name=Grace}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active}
   - Fait 2: [1] TestOrder{id=O002, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive}
   - Fait 2: [1] TestOrder{id=O003, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001, amount=3}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, level=9, salary=95000, status=active, department=management, name=Grace, age=65, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O004, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P002, level=5, score=9.2, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa, name=Frank}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, status=active, name=X, age=22, active=true, tags=temp, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, name=Ivy, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior}
   - Fait 2: [1] TestOrder{id=O001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, product_id=PROD001}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O007, region=north, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, level=5, score=9.2, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O009, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{id=O009, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P002, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2, name=Bob, age=35, salary=75000}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P006, score=0, status=active, level=1, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true}
   - Fait 2: [1] TestOrder{id=O010, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P003, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive, name=Charlie}
   - Fait 2: [1] TestOrder{id=O002, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, level=1, salary=28000, score=6.5, status=active, name=X, age=22, active=true, tags=temp}
   - Fait 2: [1] TestOrder{id=O003, total=225, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P001, status=active, level=2, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, total=299.99, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive}
   - Fait 2: [1] TestOrder{id=O006, region=west, amount=2, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1, name=Henry, age=18, active=false}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P005, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}
   - Fait 2: [1] TestOrder{id=O005, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O007, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P005, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O005, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active, name=X, age=22, active=true}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P009, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, score=5.5, level=1, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support}
   - Fait 2: [1] TestOrder{id=O004, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, discount=50, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa}
   - Fait 2: [1] TestOrder{id=O001, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, level=5, score=9.2, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, total=225, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, status=active, department=management, name=Grace, age=65, active=true, score=10, tags=executive, level=9}
   - Fait 2: [1] TestOrder{id=O003, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001, amount=3}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45}
   - Fait 2: [1] TestOrder{id=O003, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225, customer_id=P001, amount=3, date=2024-02-01}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee, level=3, age=30}
   - Fait 2: [1] TestOrder{id=O006, region=west, amount=2, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O008, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive, name=Charlie, score=6}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2}
   - Fait 2: [1] TestOrder{id=O001, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98, customer_id=P001, amount=2}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, status=inactive, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P007, level=9, salary=95000, status=active, department=management, name=Grace, age=65, active=true, score=10, tags=executive}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10, total=255}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8, status=active, department=marketing}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10, total=255}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2}
3. [1] TestPerson{id=P003, age=16, salary=0, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr, level=1}
4. [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager}
5. [1] TestPerson{id=P005, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8}
6. [1] TestPerson{id=P006, tags=test, department=qa, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1}
7. [1] TestPerson{id=P007, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management, name=Grace, age=65}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}
9. [1] TestPerson{id=P009, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true, status=active, level=6}
10. [1] TestPerson{id=P010, score=6.5, status=active, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, score=6, department=hr, level=1, age=16, salary=0, active=false, tags=intern, status=inactive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, tags=manager, active=true, score=7.8, status=active, department=marketing, level=7}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active, department=management}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, status=active, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales, age=25}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P002, level=5, score=9.2, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, level=6, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, discount=50, region=north, product_id=PROD001, total=1999.98, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal}
2. [1] TestOrder{id=O002, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20, priority=low, customer_id=P002, total=25.5}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225}
4. [1] TestOrder{id=O004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004, total=299.99, status=delivered, priority=normal, customer_id=P004}
5. [1] TestOrder{id=O005, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south, product_id=PROD001, amount=1}
6. [1] TestOrder{id=O006, priority=low, discount=0, region=west, amount=2, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled}
7. [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50}
8. [1] TestOrder{id=O008, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10, total=255}
9. [1] TestOrder{id=O009, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}
10. [1] TestOrder{id=O010, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, product_id=PROD001, total=1999.98}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, priority=low, customer_id=P002, total=25.5, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, date=2024-01-20}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, product_id=PROD003, total=225}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, amount=2}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, total=600, date=2024-03-01, status=shipped, region=north, product_id=PROD006, priority=urgent, discount=50, customer_id=P007, amount=4}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, total=299.99, status=delivered, priority=normal, customer_id=P004, amount=1, date=2024-02-05, discount=0, region=east, product_id=PROD004}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, status=confirmed, discount=100, customer_id=P002, total=999.99, date=2024-02-10, priority=high, region=south}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, department=sales, age=25, salary=45000, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2}
3. [1] TestPerson{id=P003, level=1, age=16, salary=0, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr}
4. [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager}
5. [1] TestPerson{id=P005, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa}
7. [1] TestPerson{id=P007, department=management, name=Grace, age=65, active=true, score=10, tags=executive, level=9, salary=95000, status=active}
8. [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, salary=25000, score=5.5, level=1, name=Henry, age=18, active=false}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, status=active, level=6, age=40, score=8.7, tags=senior, department=engineering}
10. [1] TestPerson{id=P010, name=X, age=22, active=true, tags=temp, department=intern, level=1, salary=28000, score=6.5, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, status=active, level=2, name=Alice, active=true, score=8.5, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, score=9.2}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, active=true, score=0, status=active, level=1, tags=test, department=qa, name=Frank}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, score=5.5, level=1, name=Henry, age=18, active=false, tags=junior, status=inactive, department=support}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, age=40, score=8.7, tags=senior, department=engineering, name=Ivy, salary=68000, active=true, status=active, level=6}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, active=false, tags=intern, status=inactive, name=Charlie, score=6, department=hr, level=1, age=16, salary=0}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, tags=manager}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, level=9, salary=95000, status=active, department=management, name=Grace, age=65, active=true}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, status=active, name=X, age=22, active=true, tags=temp, department=intern, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
