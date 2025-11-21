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

1. [1] TestPerson{id=P001, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25, salary=45000}
2. [1] TestPerson{id=P002, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active}
3. [1] TestPerson{id=P003, score=6, level=1, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern}
4. [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true, tags=manager}
5. [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}
6. [1] TestPerson{id=P006, age=0, active=true, tags=test, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank}
7. [1] TestPerson{id=P007, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65}
8. [1] TestPerson{id=P008, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry, tags=junior}
9. [1] TestPerson{id=P009, level=6, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true}
10. [1] TestPerson{id=P010, name=X, score=6.5, level=1, salary=28000, active=true, tags=temp, department=intern, age=22, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, active=true, level=6, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25, salary=45000, active=true}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16, salary=0, active=false}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P008, age=18, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, age=22, status=active, name=X, score=6.5, level=1, salary=28000, active=true, tags=temp, department=intern}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001}
2. [1] TestOrder{id=O002, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low}
3. [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high}
4. [1] TestOrder{id=O004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east, customer_id=P004}
5. [1] TestOrder{id=O005, amount=1, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south}
6. [1] TestOrder{id=O006, priority=low, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2}
7. [1] TestOrder{id=O007, customer_id=P007, total=600, priority=urgent, region=north, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped}
8. [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north}
10. [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, priority=low, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O008, status=pending, amount=10, total=255, discount=0, product_id=PROD002, date=2024-03-05, priority=normal, region=south, customer_id=P010}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, total=600, priority=urgent, region=north, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded, priority=urgent}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, status=active, department=sales, age=25, salary=45000, active=true, level=2, score=8.5, tags=junior}
2. [1] TestPerson{id=P002, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern, score=6, level=1}
4. [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45}
5. [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa}
7. [1] TestPerson{id=P007, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65, name=Grace}
8. [1] TestPerson{id=P008, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive}
9. [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering}
10. [1] TestPerson{id=P010, level=1, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, status=active, department=sales, age=25, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, status=active, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, name=Eve, department=sales, age=30, score=8, status=inactive, level=3, active=false, salary=55000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north}
2. [1] TestOrder{id=O002, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0}
3. [1] TestOrder{id=O003, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225}
4. [1] TestOrder{id=O004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east, customer_id=P004}
5. [1] TestOrder{id=O005, amount=1, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west}
7. [1] TestOrder{id=O007, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north, amount=4}
8. [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}
9. [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1}
10. [1] TestOrder{id=O010, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98, priority=normal}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O006, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, rating=4.5, keywords=computer, brand=TechCorp, category=electronics, available=true, name=Laptop, stock=50, supplier=TechSupply, price=999.99}
2. [1] TestProduct{id=PROD002, available=true, supplier=TechSupply, name=Mouse, category=accessories, price=25.5, rating=4, keywords=peripheral, brand=TechCorp, stock=200}
3. [1] TestProduct{id=PROD003, name=Keyboard, available=false, price=75, supplier=KeySupply, category=accessories, keywords=typing, stock=0, rating=3.5, brand=KeyTech}
4. [1] TestProduct{id=PROD004, rating=4.8, price=299.99, available=true, stock=30, supplier=ScreenSupply, keywords=display, brand=ScreenPro, name=Monitor, category=electronics}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, price=8.5, rating=2, available=false, keywords=obsolete, supplier=OldSupply, brand=OldTech, category=accessories, stock=0}
6. [1] TestProduct{id=PROD006, available=true, supplier=AudioSupply, name=Headphones, keywords=sound, brand=AudioMax, category=audio, price=150, rating=4.6, stock=75}
7. [1] TestProduct{id=PROD007, brand=CamTech, supplier=CamSupply, category=electronics, price=89.99, available=true, rating=3.8, stock=25, name=Webcam, keywords=video}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD003, name=Keyboard, available=false, price=75, supplier=KeySupply, category=accessories, keywords=typing, stock=0, rating=3.5, brand=KeyTech}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD004, brand=ScreenPro, name=Monitor, category=electronics, rating=4.8, price=299.99, available=true, stock=30, supplier=ScreenSupply, keywords=display}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD006, stock=75, available=true, supplier=AudioSupply, name=Headphones, keywords=sound, brand=AudioMax, category=audio, price=150, rating=4.6}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD007, keywords=video, brand=CamTech, supplier=CamSupply, category=electronics, price=89.99, available=true, rating=3.8, stock=25, name=Webcam}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD001, price=999.99, rating=4.5, keywords=computer, brand=TechCorp, category=electronics, available=true, name=Laptop, stock=50, supplier=TechSupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, category=accessories, price=25.5, rating=4, keywords=peripheral, brand=TechCorp, stock=200, available=true, supplier=TechSupply}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25}
2. [1] TestPerson{id=P002, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering, score=9.2, level=5}
3. [1] TestPerson{id=P003, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16}
4. [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true, tags=manager}
5. [1] TestPerson{id=P005, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales, age=30, score=8}
6. [1] TestPerson{id=P006, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa, salary=-5000}
7. [1] TestPerson{id=P007, department=management, active=true, score=10, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000}
8. [1] TestPerson{id=P008, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry}
9. [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering}
10. [1] TestPerson{id=P010, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1, salary=28000, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa, salary=-5000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98}
2. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002}
3. [1] TestOrder{id=O003, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225}
4. [1] TestOrder{id=O004, date=2024-02-05, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1}
5. [1] TestOrder{id=O005, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002}
6. [1] TestOrder{id=O006, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low}
7. [1] TestOrder{id=O007, customer_id=P007, total=600, priority=urgent, region=north, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped}
8. [1] TestOrder{id=O008, total=255, discount=0, product_id=PROD002, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10}
9. [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1}
10. [1] TestOrder{id=O010, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002, total=25.5}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, discount=0, product_id=PROD002, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, discount=0, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25}
2. [1] TestPerson{id=P002, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active}
3. [1] TestPerson{id=P003, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern, score=6, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}
6. [1] TestPerson{id=P006, age=0, active=true, tags=test, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank}
7. [1] TestPerson{id=P007, department=management, active=true, score=10, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000}
8. [1] TestPerson{id=P008, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000}
9. [1] TestPerson{id=P009, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy}
10. [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, score=6, level=1, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales, age=30}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, age=18, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001}
2. [1] TestOrder{id=O002, status=confirmed, priority=low, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1}
3. [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high}
4. [1] TestOrder{id=O004, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered, discount=0}
5. [1] TestOrder{id=O005, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed}
6. [1] TestOrder{id=O006, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005}
7. [1] TestOrder{id=O007, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north}
8. [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}
9. [1] TestOrder{id=O009, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007}
10. [1] TestOrder{id=O010, discount=0, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O010, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}
2. [1] TestPerson{id=P002, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active}
3. [1] TestPerson{id=P003, age=16, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr}
4. [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true, tags=manager}
5. [1] TestPerson{id=P005, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales, age=30, score=8, status=inactive}
6. [1] TestPerson{id=P006, active=true, tags=test, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank, age=0}
7. [1] TestPerson{id=P007, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65}
8. [1] TestPerson{id=P008, age=18, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5}
9. [1] TestPerson{id=P009, level=6, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true}
10. [1] TestPerson{id=P010, age=22, status=active, name=X, score=6.5, level=1, salary=28000, active=true, tags=temp, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, age=18, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7, name=Diana}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, department=management, active=true, score=10, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, price=999.99, rating=4.5, keywords=computer, brand=TechCorp, category=electronics, available=true, name=Laptop, stock=50, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, available=true, supplier=TechSupply, name=Mouse, category=accessories, price=25.5, rating=4, keywords=peripheral, brand=TechCorp, stock=200}
3. [1] TestProduct{id=PROD003, keywords=typing, stock=0, rating=3.5, brand=KeyTech, name=Keyboard, available=false, price=75, supplier=KeySupply, category=accessories}
4. [1] TestProduct{id=PROD004, rating=4.8, price=299.99, available=true, stock=30, supplier=ScreenSupply, keywords=display, brand=ScreenPro, name=Monitor, category=electronics}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, price=8.5, rating=2, available=false, keywords=obsolete, supplier=OldSupply, brand=OldTech, category=accessories, stock=0}
6. [1] TestProduct{id=PROD006, category=audio, price=150, rating=4.6, stock=75, available=true, supplier=AudioSupply, name=Headphones, keywords=sound, brand=AudioMax}
7. [1] TestProduct{id=PROD007, category=electronics, price=89.99, available=true, rating=3.8, stock=25, name=Webcam, keywords=video, brand=CamTech, supplier=CamSupply}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD006, stock=75, available=true, supplier=AudioSupply, name=Headphones, keywords=sound, brand=AudioMax, category=audio, price=150, rating=4.6}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD007, category=electronics, price=89.99, available=true, rating=3.8, stock=25, name=Webcam, keywords=video, brand=CamTech, supplier=CamSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD001, stock=50, supplier=TechSupply, price=999.99, rating=4.5, keywords=computer, brand=TechCorp, category=electronics, available=true, name=Laptop}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, category=accessories, price=25.5, rating=4, keywords=peripheral, brand=TechCorp, stock=200, available=true, supplier=TechSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD003, category=accessories, keywords=typing, stock=0, rating=3.5, brand=KeyTech, name=Keyboard, available=false, price=75, supplier=KeySupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD004, category=electronics, rating=4.8, price=299.99, available=true, stock=30, supplier=ScreenSupply, keywords=display, brand=ScreenPro, name=Monitor}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25}
2. [1] TestPerson{id=P002, salary=75000, status=active, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35}
3. [1] TestPerson{id=P003, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16}
4. [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa}
7. [1] TestPerson{id=P007, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10}
8. [1] TestPerson{id=P008, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1}
9. [1] TestPerson{id=P009, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6, salary=68000}
10. [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, name=Alice, status=active, department=sales, age=25, salary=45000, active=true, level=2, score=8.5}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7, name=Diana}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern, score=6, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6, salary=68000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north, amount=2}
2. [1] TestOrder{id=O002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002}
3. [1] TestOrder{id=O003, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north}
4. [1] TestOrder{id=O004, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered}
5. [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high}
6. [1] TestOrder{id=O006, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15}
7. [1] TestOrder{id=O007, status=shipped, customer_id=P007, total=600, priority=urgent, region=north, amount=4, discount=50, product_id=PROD006, date=2024-03-01}
8. [1] TestOrder{id=O008, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002, date=2024-03-05, priority=normal, region=south}
9. [1] TestOrder{id=O009, product_id=PROD007, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10}
10. [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, total=255, discount=0, product_id=PROD002, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O001, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3, status=shipped}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north, amount=4}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25}
2. [1] TestPerson{id=P002, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, score=6, level=1, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern}
4. [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, salary=55000, tags=employee, name=Eve, department=sales, age=30, score=8, status=inactive, level=3, active=false}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa}
7. [1] TestPerson{id=P007, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65, name=Grace, tags=executive}
8. [1] TestPerson{id=P008, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false}
9. [1] TestPerson{id=P009, tags=senior, status=active, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7}
10. [1] TestPerson{id=P010, score=6.5, level=1, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X}
11. [1] TestOrder{id=O001, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15}
12. [1] TestOrder{id=O002, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south}
13. [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high}
14. [1] TestOrder{id=O004, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered, discount=0, total=299.99}
15. [1] TestOrder{id=O005, amount=1, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south}
16. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west}
17. [1] TestOrder{id=O007, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north}
18. [1] TestOrder{id=O008, amount=10, total=255, discount=0, product_id=PROD002, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending}
19. [1] TestOrder{id=O009, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007}
20. [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25, salary=45000, active=true, level=2}
   - Fait 2: [1] TestOrder{id=O004, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry}
   - Fait 2: [1] TestOrder{id=O005, amount=1, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1, salary=28000}
   - Fait 2: [1] TestOrder{id=O005, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales, age=30, score=8, status=inactive}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, level=5, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering, score=9.2}
   - Fait 2: [1] TestOrder{id=O007, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, discount=0, product_id=PROD002, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16, salary=0, active=false}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P008, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5, age=18}
   - Fait 2: [1] TestOrder{id=O003, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, age=22, status=active, name=X, score=6.5, level=1, salary=28000, active=true, tags=temp}
   - Fait 2: [1] TestOrder{id=O003, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P007, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O010, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering, score=9.2, level=5}
   - Fait 2: [1] TestOrder{id=O001, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P006, active=true, tags=test, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank, age=0}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north, amount=4, discount=50}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16}
   - Fait 2: [1] TestOrder{id=O008, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002, date=2024-03-05, priority=normal}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P006, level=1, name=Frank, age=0, active=true, tags=test, department=qa, salary=-5000, score=0, status=active}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P004, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, department=management, active=true, score=10, age=65, name=Grace, tags=executive, status=active, level=9}
   - Fait 2: [1] TestOrder{id=O004, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25, salary=45000, active=true, level=2}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1, total=89.99, status=completed}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P003, score=6, level=1, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded, priority=urgent}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O004, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered, discount=0, total=299.99}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25}
   - Fait 2: [1] TestOrder{id=O006, amount=2, priority=low, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}
   - Fait 2: [1] TestOrder{id=O008, status=pending, amount=10, total=255, discount=0, product_id=PROD002, date=2024-03-05, priority=normal, region=south, customer_id=P010}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P006, active=true, tags=test, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank, age=0}
   - Fait 2: [1] TestOrder{id=O001, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98, priority=normal}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P002, level=5, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering, score=9.2}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P005, active=false, salary=55000, tags=employee, name=Eve, department=sales, age=30, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north, amount=4, discount=50, product_id=PROD006}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry}
   - Fait 2: [1] TestOrder{id=O007, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P010, age=22, status=active, name=X, score=6.5, level=1, salary=28000, active=true, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north, amount=4, discount=50, product_id=PROD006}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16}
   - Fait 2: [1] TestOrder{id=O009, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy}
   - Fait 2: [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P003, score=6, level=1, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered, discount=0}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P007, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, tags=employee, name=Eve, department=sales, age=30, score=8, status=inactive, level=3, active=false}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false}
   - Fait 2: [1] TestOrder{id=O001, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25}
   - Fait 2: [1] TestOrder{id=O002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P009, active=true, level=6, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P006, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O007, total=600, priority=urgent, region=north, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P007, department=management, active=true, score=10, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000}
   - Fait 2: [1] TestOrder{id=O008, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002, date=2024-03-05, priority=normal}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active}
   - Fait 2: [1] TestOrder{id=O002, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002, total=25.5, date=2024-01-20}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P007, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa}
   - Fait 2: [1] TestOrder{id=O004, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O009, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6}
   - Fait 2: [1] TestOrder{id=O010, discount=0, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true, tags=manager}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65, name=Grace}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P009, status=active, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7, tags=senior}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P003, level=1, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern, score=6}
   - Fait 2: [1] TestOrder{id=O007, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6, salary=68000}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P005, active=false, salary=55000, tags=employee, name=Eve, department=sales, age=30, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O009, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P008, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5, age=18}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P010, name=X, score=6.5, level=1, salary=28000, active=true, tags=temp, department=intern, age=22, status=active}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{id=O004, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P010, name=X, score=6.5, level=1, salary=28000, active=true, tags=temp, department=intern, age=22, status=active}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high, total=999.99}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering, score=9.2, level=5}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded, priority=urgent}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16}
   - Fait 2: [1] TestOrder{id=O001, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, status=active, department=sales, age=25, salary=45000, active=true, level=2, score=8.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O003, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior}
   - Fait 2: [1] TestOrder{id=O006, amount=2, priority=low, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active}
   - Fait 2: [1] TestOrder{id=O007, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, status=active, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true}
   - Fait 2: [1] TestOrder{id=O008, total=255, discount=0, product_id=PROD002, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P003, score=6, level=1, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P007, department=management, active=true, score=10, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000}
   - Fait 2: [1] TestOrder{id=O007, region=north, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P002, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering, score=9.2, level=5, name=Bob}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, tags=employee, name=Eve, department=sales, age=30, score=8, status=inactive, level=3, active=false}
   - Fait 2: [1] TestOrder{id=O002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25, salary=45000, active=true, level=2}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false}
   - Fait 2: [1] TestOrder{id=O002, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002, total=25.5}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P010, level=1, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P003, score=6, level=1, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern}
   - Fait 2: [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered, discount=0}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, tags=test, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, department=management, active=true, score=10, age=65, name=Grace, tags=executive, status=active, level=9}
   - Fait 2: [1] TestOrder{id=O001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7, name=Diana}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, name=Eve, department=sales, age=30, score=8, status=inactive, level=3, active=false, salary=55000}
   - Fait 2: [1] TestOrder{id=O006, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P008, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5, age=18}
   - Fait 2: [1] TestOrder{id=O006, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P010, age=22, status=active, name=X, score=6.5, level=1, salary=28000, active=true, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P008, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5, age=18}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, discount=0, product_id=PROD002, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, tags=employee, name=Eve, department=sales, age=30, score=8, status=inactive, level=3, active=false}
   - Fait 2: [1] TestOrder{id=O008, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002, date=2024-03-05}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true}
   - Fait 2: [1] TestOrder{id=O003, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, tags=test, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank}
   - Fait 2: [1] TestOrder{id=O002, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, age=25, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, total=600, priority=urgent, region=north, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7}
   - Fait 2: [1] TestOrder{id=O009, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P008, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5, age=18}
   - Fait 2: [1] TestOrder{id=O010, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P010, name=X, score=6.5, level=1, salary=28000, active=true, tags=temp, department=intern, age=22, status=active}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P009, status=active, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7, tags=senior}
   - Fait 2: [1] TestOrder{id=O001, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}
2. [1] TestPerson{id=P002, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active}
3. [1] TestPerson{id=P003, score=6, level=1, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern}
4. [1] TestPerson{id=P004, tags=manager, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true}
5. [1] TestPerson{id=P005, tags=employee, name=Eve, department=sales, age=30, score=8, status=inactive, level=3, active=false, salary=55000}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa}
7. [1] TestPerson{id=P007, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65}
8. [1] TestPerson{id=P008, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7, tags=senior, status=active}
10. [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, tags=test, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25}
2. [1] TestPerson{id=P002, level=5, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering, score=9.2}
3. [1] TestPerson{id=P003, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16}
4. [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa}
7. [1] TestPerson{id=P007, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65}
8. [1] TestPerson{id=P008, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5, age=18}
9. [1] TestPerson{id=P009, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy}
10. [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1}
11. [1] TestOrder{id=O001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001}
12. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002}
13. [1] TestOrder{id=O003, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225}
14. [1] TestOrder{id=O004, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered, discount=0}
15. [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high}
16. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west}
17. [1] TestOrder{id=O007, customer_id=P007, total=600, priority=urgent, region=north, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped}
18. [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}
19. [1] TestOrder{id=O009, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007}
20. [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65, name=Grace}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25, salary=45000}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P003, score=6, level=1, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern}
   - Fait 2: [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, tags=employee, name=Eve, department=sales, age=30, score=8, status=inactive, level=3, active=false}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry}
   - Fait 2: [1] TestOrder{id=O002, priority=low, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, age=22, status=active, name=X, score=6.5, level=1, salary=28000, active=true, tags=temp}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P007, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O007, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P002, level=5, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering, score=9.2}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002, date=2024-03-05, priority=normal, region=south}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P003, score=6, level=1, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, tags=test, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending, customer_id=P001}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, age=25, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true, tags=manager}
   - Fait 2: [1] TestOrder{id=O004, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P007, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65, name=Grace, tags=executive}
   - Fait 2: [1] TestOrder{id=O010, discount=0, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P003, score=6, level=1, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, status=active, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true}
   - Fait 2: [1] TestOrder{id=O002, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O003, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, name=Eve, department=sales, age=30, score=8, status=inactive, level=3, active=false, salary=55000}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true, tags=manager}
   - Fait 2: [1] TestOrder{id=O007, status=shipped, customer_id=P007, total=600, priority=urgent, region=north, amount=4, discount=50, product_id=PROD006, date=2024-03-01}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P008, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5, age=18}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6}
   - Fait 2: [1] TestOrder{id=O009, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, department=sales, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P001, status=active, department=sales, age=25, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice}
   - Fait 2: [1] TestOrder{id=O003, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1}
   - Fait 2: [1] TestOrder{id=O007, total=600, priority=urgent, region=north, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, department=qa, salary=-5000, score=0, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P001, status=active, department=sales, age=25, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice}
   - Fait 2: [1] TestOrder{id=O009, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P004, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded, priority=urgent}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P010, level=1, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior}
   - Fait 2: [1] TestOrder{id=O003, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{id=O002, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P009, status=active, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7, tags=senior}
   - Fait 2: [1] TestOrder{id=O002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1, salary=28000, active=true}
   - Fait 2: [1] TestOrder{id=O004, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true, tags=manager}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P007, department=management, active=true, score=10, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, tags=test, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25}
   - Fait 2: [1] TestOrder{id=O005, amount=1, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16}
   - Fait 2: [1] TestOrder{id=O007, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north, amount=4}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P007, score=10, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P003, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16, salary=0}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16, salary=0, active=false}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P002, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering, score=9.2, level=5, name=Bob}
   - Fait 2: [1] TestOrder{id=O005, region=south, amount=1, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, age=25, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active}
   - Fait 2: [1] TestOrder{id=O006, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true}
   - Fait 2: [1] TestOrder{id=O007, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, level=1, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002, date=2024-03-05, priority=normal, region=south}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test}
   - Fait 2: [1] TestOrder{id=O010, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P008, active=false, status=inactive, department=support, name=Henry, tags=junior, level=1, salary=25000, score=5.5, age=18}
   - Fait 2: [1] TestOrder{id=O007, total=600, priority=urgent, region=north, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, tags=employee, name=Eve, department=sales, age=30, score=8, status=inactive, level=3, active=false}
   - Fait 2: [1] TestOrder{id=O005, amount=1, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P010, name=X, score=6.5, level=1, salary=28000, active=true, tags=temp, department=intern, age=22, status=active}
   - Fait 2: [1] TestOrder{id=O001, amount=2, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P009, status=active, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7, tags=senior}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high, total=999.99}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P004, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa}
   - Fait 2: [1] TestOrder{id=O005, amount=1, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active}
   - Fait 2: [1] TestOrder{id=O007, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management}
   - Fait 2: [1] TestOrder{id=O008, discount=0, product_id=PROD002, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1, salary=28000, active=true}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O004, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16, salary=0, active=false}
   - Fait 2: [1] TestOrder{id=O005, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, tags=test, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north, amount=4, discount=50}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, status=active, department=sales, age=25, salary=45000, active=true, level=2, score=8.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10, age=65}
   - Fait 2: [1] TestOrder{id=O001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending, customer_id=P001, product_id=PROD001}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O003, region=north, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P008, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry, tags=junior}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north, amount=4, discount=50, product_id=PROD006}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, total=600, priority=urgent, region=north, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true, tags=manager}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10, product_id=PROD007, region=north}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6, salary=68000}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P009, level=6, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P005, active=false, salary=55000, tags=employee, name=Eve, department=sales, age=30, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, product_id=PROD004, status=delivered, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P008, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry, tags=junior}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, name=Alice, status=active, department=sales, age=25, salary=45000, active=true, level=2, score=8.5}
   - Fait 2: [1] TestOrder{id=O001, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north, amount=2}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa}
   - Fait 2: [1] TestOrder{id=O002, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low, discount=0, product_id=PROD002, total=25.5}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7, name=Diana}
   - Fait 2: [1] TestOrder{id=O005, region=south, amount=1, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, status=active, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16}
   - Fait 2: [1] TestOrder{id=O002, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed, priority=low}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry}
   - Fait 2: [1] TestOrder{id=O003, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true, tags=manager}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0, region=west}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, tags=test, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25, salary=45000, active=true, level=2}
2. [1] TestPerson{id=P002, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering, score=9.2, level=5}
3. [1] TestPerson{id=P003, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16}
4. [1] TestPerson{id=P004, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active}
5. [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}
6. [1] TestPerson{id=P006, active=true, tags=test, department=qa, salary=-5000, score=0, status=active, level=1, name=Frank, age=0}
7. [1] TestPerson{id=P007, active=true, score=10, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management}
8. [1] TestPerson{id=P008, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry}
9. [1] TestPerson{id=P009, status=active, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X, score=6.5, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, department=management, active=true, score=10, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, status=active, department=sales, age=25, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, department=engineering, score=8.7, tags=senior, status=active, name=Ivy, age=40, active=true, level=6}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P010, age=22, status=active, name=X, score=6.5, level=1, salary=28000, active=true, tags=temp, department=intern}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, status=active, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending}
2. [1] TestOrder{id=O002, status=confirmed, priority=low, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1}
3. [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, region=north, total=225, discount=15, customer_id=P001, amount=3, status=shipped, priority=high}
4. [1] TestOrder{id=O004, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered, discount=0, total=299.99}
5. [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1, status=confirmed, discount=100, priority=high}
6. [1] TestOrder{id=O006, total=999.98, discount=0, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low}
7. [1] TestOrder{id=O007, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north}
8. [1] TestOrder{id=O008, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002}
9. [1] TestOrder{id=O009, discount=10, date=2024-03-10, product_id=PROD007, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low}
10. [1] TestOrder{id=O010, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, discount=50, product_id=PROD006, date=2024-03-01, status=shipped, customer_id=P007, total=600, priority=urgent, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O010, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, region=east, discount=0, product_id=PROD001}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, discount=50, total=1999.98, priority=normal, region=north, amount=2, status=pending}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O004, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, product_id=PROD004, status=delivered, discount=0, total=299.99}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, region=west, customer_id=P005, product_id=PROD005, date=2024-02-15, status=cancelled, amount=2, priority=low, total=999.98, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, status=pending, amount=10, total=255, discount=0, product_id=PROD002, date=2024-03-05, priority=normal, region=south}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, region=north, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, date=2024-03-10}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, priority=low, discount=0, product_id=PROD002, total=25.5, date=2024-01-20, region=south, customer_id=P002, amount=1, status=confirmed}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O003, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, product_id=PROD003, date=2024-02-01, region=north, total=225}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O005, status=confirmed, discount=100, priority=high, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, region=south, amount=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, name=Alice, status=active, department=sales, age=25, salary=45000, active=true, level=2, score=8.5}
2. [1] TestPerson{id=P002, tags=senior, department=engineering, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active}
3. [1] TestPerson{id=P003, active=false, status=inactive, tags=intern, score=6, level=1, name=Charlie, department=hr, age=16, salary=0}
4. [1] TestPerson{id=P004, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing, level=7, name=Diana}
5. [1] TestPerson{id=P005, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales, age=30}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa}
7. [1] TestPerson{id=P007, department=management, active=true, score=10, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000}
8. [1] TestPerson{id=P008, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry}
9. [1] TestPerson{id=P009, tags=senior, status=active, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7}
10. [1] TestPerson{id=P010, age=22, status=active, name=X, score=6.5, level=1, salary=28000, active=true, tags=temp, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, age=25}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, active=false, salary=55000, tags=employee, name=Eve, department=sales}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, status=active, level=1, name=Frank, age=0, active=true, tags=test, department=qa}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, age=65, name=Grace, tags=executive, status=active, level=9, salary=95000, department=management, active=true, score=10}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, level=1, salary=25000, score=5.5, age=18, active=false, status=inactive, department=support, name=Henry}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, level=5, name=Bob, active=true, age=35, salary=75000, status=active, tags=senior, department=engineering}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, level=1, name=Charlie, department=hr, age=16, salary=0, active=false, status=inactive, tags=intern, score=6}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, age=45, salary=85000, active=true, tags=manager, score=7.8, status=active, department=marketing}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, status=active, name=Ivy, age=40, active=true, level=6, salary=68000, department=engineering, score=8.7, tags=senior}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, level=1, salary=28000, active=true, tags=temp, department=intern, age=22, status=active, name=X}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
