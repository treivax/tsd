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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, name=Bob, tags=senior, status=active, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering}
3. [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, level=1, name=Charlie, age=16}
4. [1] TestPerson{id=P004, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active, name=Diana, age=45}
5. [1] TestPerson{id=P005, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, active=false, status=inactive}
6. [1] TestPerson{id=P006, level=1, name=Frank, age=0, active=true, status=active, department=qa, salary=-5000, score=0, tags=test}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, status=active, department=management, name=Grace, score=10, level=9}
8. [1] TestPerson{id=P008, active=false, tags=junior, score=5.5, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000}
9. [1] TestPerson{id=P009, age=40, salary=68000, tags=senior, department=engineering, active=true, score=8.7, status=active, level=6, name=Ivy}
10. [1] TestPerson{id=P010, name=X, score=6.5, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, active=false, status=inactive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false, tags=junior}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, department=engineering, active=true, score=8.7, status=active, level=6}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, level=2, name=Alice, age=25}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, name=Bob, tags=senior, status=active, level=5, age=35, salary=75000, active=true, score=9.2}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, score=10, level=9, age=65, salary=95000, active=true, tags=executive, status=active, department=management}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, name=X, score=6.5}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, name=Diana, age=45, score=7.8, department=marketing, level=7}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, status=pending, discount=50, region=north, product_id=PROD001, priority=normal}
2. [1] TestOrder{id=O002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, amount=1, total=25.5, discount=0, region=south}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, status=shipped, region=north}
4. [1] TestOrder{id=O004, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, amount=1, total=299.99, status=delivered, priority=normal, discount=0}
5. [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, discount=100, total=999.99, priority=high, region=south, customer_id=P002, product_id=PROD001, amount=1}
6. [1] TestOrder{id=O006, customer_id=P005, amount=2, total=999.98, priority=low, discount=0, region=west, product_id=PROD005, date=2024-02-15, status=cancelled}
7. [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, priority=urgent}
8. [1] TestOrder{id=O008, amount=10, total=255, status=pending, region=south, customer_id=P010, date=2024-03-05, priority=normal, discount=0, product_id=PROD002}
9. [1] TestOrder{id=O009, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, region=north, customer_id=P001, date=2024-03-10, status=completed}
10. [1] TestOrder{id=O010, product_id=PROD001, total=75000, priority=urgent, discount=0, amount=1, date=2024-03-15, status=refunded, region=east, customer_id=P006}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, total=255, status=pending, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, status=pending, discount=50, region=north}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, priority=high, region=south, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, discount=100, total=999.99}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, priority=urgent}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, region=north, customer_id=P001, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O010, total=75000, priority=urgent, discount=0, amount=1, date=2024-03-15, status=refunded, region=east, customer_id=P006, product_id=PROD001}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, discount=0, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, status=shipped, region=north, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, tags=senior, status=active, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering, name=Bob}
3. [1] TestPerson{id=P003, score=6, tags=intern, status=inactive, department=hr, level=1, name=Charlie, age=16, salary=0, active=false}
4. [1] TestPerson{id=P004, active=true, tags=manager, status=active, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000}
5. [1] TestPerson{id=P005, department=sales, level=3, name=Eve, age=30, salary=55000, active=false, status=inactive, score=8, tags=employee}
6. [1] TestPerson{id=P006, salary=-5000, score=0, tags=test, level=1, name=Frank, age=0, active=true, status=active, department=qa}
7. [1] TestPerson{id=P007, score=10, level=9, age=65, salary=95000, active=true, tags=executive, status=active, department=management, name=Grace}
8. [1] TestPerson{id=P008, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false, tags=junior, score=5.5}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, department=engineering, active=true, score=8.7, status=active, level=6}
10. [1] TestPerson{id=P010, score=6.5, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, name=X}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering, name=Bob, tags=senior, status=active}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, active=false, status=inactive}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, status=active, department=management, name=Grace, score=10, level=9, age=65, salary=95000, active=true}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, department=engineering, active=true, score=8.7, status=active, level=6}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, status=pending, discount=50, region=north, product_id=PROD001, priority=normal}
2. [1] TestOrder{id=O002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, amount=1, total=25.5, discount=0, region=south}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, status=shipped, region=north}
4. [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05}
5. [1] TestOrder{id=O005, status=confirmed, discount=100, total=999.99, priority=high, region=south, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10}
6. [1] TestOrder{id=O006, customer_id=P005, amount=2, total=999.98, priority=low, discount=0, region=west, product_id=PROD005, date=2024-02-15, status=cancelled}
7. [1] TestOrder{id=O007, product_id=PROD006, priority=urgent, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007}
8. [1] TestOrder{id=O008, priority=normal, discount=0, product_id=PROD002, amount=10, total=255, status=pending, region=south, customer_id=P010, date=2024-03-05}
9. [1] TestOrder{id=O009, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, region=north, customer_id=P001, date=2024-03-10, status=completed}
10. [1] TestOrder{id=O010, total=75000, priority=urgent, discount=0, amount=1, date=2024-03-15, status=refunded, region=east, customer_id=P006, product_id=PROD001}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, discount=100, total=999.99, priority=high, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, total=255, status=pending, region=south}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O001, discount=50, region=north, product_id=PROD001, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, status=pending}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, priority=low, customer_id=P002, amount=1, total=25.5, discount=0, region=south, product_id=PROD002, date=2024-01-20, status=confirmed}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, status=shipped, region=north, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O006, date=2024-02-15, status=cancelled, customer_id=P005, amount=2, total=999.98, priority=low, discount=0, region=west, product_id=PROD005}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, priority=urgent}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, region=north}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, date=2024-02-05, amount=1, total=299.99, status=delivered, priority=normal, discount=0, region=east}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, category=electronics, price=999.99, stock=50, supplier=TechSupply, name=Laptop, available=true, rating=4.5, keywords=computer, brand=TechCorp}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, price=25.5, rating=4, supplier=TechSupply, available=true, keywords=peripheral, brand=TechCorp, stock=200}
3. [1] TestProduct{id=PROD003, category=accessories, available=false, rating=3.5, brand=KeyTech, name=Keyboard, price=75, keywords=typing, stock=0, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, name=Monitor, category=electronics, available=true, rating=4.8, keywords=display, supplier=ScreenSupply, price=299.99, brand=ScreenPro, stock=30}
5. [1] TestProduct{id=PROD005, supplier=OldSupply, name=OldKeyboard, category=accessories, rating=2, keywords=obsolete, brand=OldTech, stock=0, price=8.5, available=false}
6. [1] TestProduct{id=PROD006, name=Headphones, category=audio, available=true, keywords=sound, brand=AudioMax, stock=75, price=150, rating=4.6, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, supplier=CamSupply, name=Webcam, available=true, rating=3.8, brand=CamTech, category=electronics, price=89.99, keywords=video, stock=25}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, category=electronics, price=999.99, stock=50, supplier=TechSupply, name=Laptop, available=true, rating=4.5, keywords=computer, brand=TechCorp}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, category=accessories, price=25.5, rating=4, supplier=TechSupply, available=true, keywords=peripheral, brand=TechCorp, stock=200}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, name=Keyboard, price=75, keywords=typing, stock=0, supplier=KeySupply, category=accessories, available=false, rating=3.5, brand=KeyTech}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, category=electronics, available=true, rating=4.8, keywords=display, supplier=ScreenSupply, price=299.99, brand=ScreenPro, stock=30}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, price=150, rating=4.6, supplier=AudioSupply, name=Headphones, category=audio, available=true, keywords=sound, brand=AudioMax, stock=75}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, available=true, rating=3.8, brand=CamTech, category=electronics, price=89.99, keywords=video, stock=25, supplier=CamSupply, name=Webcam}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, department=engineering, name=Bob, tags=senior, status=active, level=5, age=35}
3. [1] TestPerson{id=P003, department=hr, level=1, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, active=false, status=inactive}
6. [1] TestPerson{id=P006, salary=-5000, score=0, tags=test, level=1, name=Frank, age=0, active=true, status=active, department=qa}
7. [1] TestPerson{id=P007, name=Grace, score=10, level=9, age=65, salary=95000, active=true, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, active=false, tags=junior, score=5.5, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, department=engineering, active=true, score=8.7, status=active, level=6}
10. [1] TestPerson{id=P010, department=intern, level=1, name=X, score=6.5, age=22, salary=28000, active=true, tags=temp, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, level=1, name=Charlie, age=16}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, tags=test, level=1, name=Frank, age=0, active=true, status=active, department=qa}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, status=pending, discount=50, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, discount=0, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low}
3. [1] TestOrder{id=O003, date=2024-02-01, priority=high, discount=15, customer_id=P001, status=shipped, region=north, product_id=PROD003, amount=3, total=225}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, date=2024-02-05, amount=1, total=299.99, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, discount=100, total=999.99, priority=high, region=south}
6. [1] TestOrder{id=O006, status=cancelled, customer_id=P005, amount=2, total=999.98, priority=low, discount=0, region=west, product_id=PROD005, date=2024-02-15}
7. [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, priority=urgent}
8. [1] TestOrder{id=O008, discount=0, product_id=PROD002, amount=10, total=255, status=pending, region=south, customer_id=P010, date=2024-03-05, priority=normal}
9. [1] TestOrder{id=O009, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, region=north, customer_id=P001, date=2024-03-10, status=completed}
10. [1] TestOrder{id=O010, amount=1, date=2024-03-15, status=refunded, region=east, customer_id=P006, product_id=PROD001, total=75000, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, status=pending, discount=50, region=north, product_id=PROD001}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, status=confirmed, priority=low, customer_id=P002, amount=1, total=25.5, discount=0, region=south, product_id=PROD002, date=2024-01-20}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, date=2024-02-05, amount=1, total=299.99, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, date=2024-02-10, status=confirmed, discount=100, total=999.99, priority=high, region=south, customer_id=P002, product_id=PROD001}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, amount=2, total=999.98, priority=low, discount=0, region=west, product_id=PROD005, date=2024-02-15, status=cancelled}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, priority=urgent}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, status=shipped, region=north}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, status=pending, region=south, customer_id=P010, date=2024-03-05, priority=normal, discount=0}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, discount=0, amount=1, date=2024-03-15, status=refunded, region=east, customer_id=P006, product_id=PROD001, total=75000, priority=urgent}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, name=Bob, tags=senior, status=active, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering}
3. [1] TestPerson{id=P003, department=hr, level=1, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, active=false, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, status=active, department=qa, salary=-5000, score=0, tags=test, level=1}
7. [1] TestPerson{id=P007, department=management, name=Grace, score=10, level=9, age=65, salary=95000, active=true, tags=executive, status=active}
8. [1] TestPerson{id=P008, score=5.5, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false, tags=junior}
9. [1] TestPerson{id=P009, status=active, level=6, name=Ivy, age=40, salary=68000, tags=senior, department=engineering, active=true, score=8.7}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, name=X, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, department=hr, level=1, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, status=inactive, score=8, tags=employee, department=sales, level=3}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false, tags=junior}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, region=north, product_id=PROD001, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, status=pending, discount=50}
2. [1] TestOrder{id=O002, amount=1, total=25.5, discount=0, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, customer_id=P002}
3. [1] TestOrder{id=O003, status=shipped, region=north, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001}
4. [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05}
5. [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, discount=100, total=999.99, priority=high, region=south, customer_id=P002, product_id=PROD001, amount=1}
6. [1] TestOrder{id=O006, product_id=PROD005, date=2024-02-15, status=cancelled, customer_id=P005, amount=2, total=999.98, priority=low, discount=0, region=west}
7. [1] TestOrder{id=O007, total=600, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, priority=urgent, amount=4}
8. [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, total=255, status=pending, region=south}
9. [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, region=north, customer_id=P001}
10. [1] TestOrder{id=O010, priority=urgent, discount=0, amount=1, date=2024-03-15, status=refunded, region=east, customer_id=P006, product_id=PROD001, total=75000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, status=shipped, region=north, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, date=2024-02-05, amount=1, total=299.99, status=delivered, priority=normal, discount=0, region=east}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, discount=100, total=999.99, priority=high, region=south, customer_id=P002, product_id=PROD001, amount=1}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O006, region=west, product_id=PROD005, date=2024-02-15, status=cancelled, customer_id=P005, amount=2, total=999.98, priority=low, discount=0}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, priority=urgent, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, total=75000, priority=urgent, discount=0, amount=1, date=2024-03-15, status=refunded, region=east}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, total=25.5, discount=0, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, customer_id=P002}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior}
2. [1] TestPerson{id=P002, name=Bob, tags=senior, status=active, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering}
3. [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, level=1, name=Charlie, age=16, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, name=Diana, age=45, score=7.8, department=marketing, level=7}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, status=inactive, score=8, tags=employee, department=sales, level=3}
6. [1] TestPerson{id=P006, salary=-5000, score=0, tags=test, level=1, name=Frank, age=0, active=true, status=active, department=qa}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, status=active, department=management, name=Grace, score=10, level=9}
8. [1] TestPerson{id=P008, score=5.5, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false, tags=junior}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, department=engineering, active=true, score=8.7, status=active, level=6}
10. [1] TestPerson{id=P010, status=active, department=intern, level=1, name=X, score=6.5, age=22, salary=28000, active=true, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, status=active, department=sales, level=2, name=Alice, age=25, salary=45000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, level=1}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, active=false, status=inactive}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, tags=test, level=1, name=Frank, age=0, active=true, status=active, department=qa}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, department=support, level=1, name=Henry, age=18, salary=25000, active=false, tags=junior, score=5.5, status=inactive}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, active=true, score=8.7, status=active, level=6, name=Ivy, age=40, salary=68000, tags=senior}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, tags=senior, status=active, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, active=true, tags=executive, status=active, department=management, name=Grace, score=10, level=9, age=65, salary=95000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, price=999.99, stock=50, supplier=TechSupply, name=Laptop, available=true, rating=4.5, keywords=computer, brand=TechCorp, category=electronics}
2. [1] TestProduct{id=PROD002, category=accessories, price=25.5, rating=4, supplier=TechSupply, available=true, keywords=peripheral, brand=TechCorp, stock=200, name=Mouse}
3. [1] TestProduct{id=PROD003, name=Keyboard, price=75, keywords=typing, stock=0, supplier=KeySupply, category=accessories, available=false, rating=3.5, brand=KeyTech}
4. [1] TestProduct{id=PROD004, name=Monitor, category=electronics, available=true, rating=4.8, keywords=display, supplier=ScreenSupply, price=299.99, brand=ScreenPro, stock=30}
5. [1] TestProduct{id=PROD005, keywords=obsolete, brand=OldTech, stock=0, price=8.5, available=false, supplier=OldSupply, name=OldKeyboard, category=accessories, rating=2}
6. [1] TestProduct{id=PROD006, name=Headphones, category=audio, available=true, keywords=sound, brand=AudioMax, stock=75, price=150, rating=4.6, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, rating=3.8, brand=CamTech, category=electronics, price=89.99, keywords=video, stock=25, supplier=CamSupply, name=Webcam, available=true}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD007, category=electronics, price=89.99, keywords=video, stock=25, supplier=CamSupply, name=Webcam, available=true, rating=3.8, brand=CamTech}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD001, supplier=TechSupply, name=Laptop, available=true, rating=4.5, keywords=computer, brand=TechCorp, category=electronics, price=999.99, stock=50}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD002, available=true, keywords=peripheral, brand=TechCorp, stock=200, name=Mouse, category=accessories, price=25.5, rating=4, supplier=TechSupply}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD003, stock=0, supplier=KeySupply, category=accessories, available=false, rating=3.5, brand=KeyTech, name=Keyboard, price=75, keywords=typing}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD004, price=299.99, brand=ScreenPro, stock=30, name=Monitor, category=electronics, available=true, rating=4.8, keywords=display, supplier=ScreenSupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD006, keywords=sound, brand=AudioMax, stock=75, price=150, rating=4.6, supplier=AudioSupply, name=Headphones, category=audio, available=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active}
2. [1] TestPerson{id=P002, score=9.2, department=engineering, name=Bob, tags=senior, status=active, level=5, age=35, salary=75000, active=true}
3. [1] TestPerson{id=P003, level=1, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, active=false, status=inactive, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, status=active, department=qa, salary=-5000, score=0, tags=test, level=1}
7. [1] TestPerson{id=P007, level=9, age=65, salary=95000, active=true, tags=executive, status=active, department=management, name=Grace, score=10}
8. [1] TestPerson{id=P008, score=5.5, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false, tags=junior}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, department=engineering, active=true, score=8.7, status=active, level=6}
10. [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, name=X, score=6.5, age=22}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, status=active, department=management, name=Grace, score=10, level=9}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, active=false, tags=junior, score=5.5, status=inactive, department=support, level=1, name=Henry, age=18}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, tags=intern, status=inactive, department=hr, level=1, name=Charlie, age=16, salary=0}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, active=false, status=inactive}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, status=active, level=6, name=Ivy, age=40, salary=68000, tags=senior, department=engineering, active=true, score=8.7}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, status=active, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, tags=senior, status=active, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, status=active, department=qa, salary=-5000, score=0, tags=test, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, status=pending, discount=50, region=north}
2. [1] TestOrder{id=O002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, amount=1, total=25.5, discount=0, region=south}
3. [1] TestOrder{id=O003, priority=high, discount=15, customer_id=P001, status=shipped, region=north, product_id=PROD003, amount=3, total=225, date=2024-02-01}
4. [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, discount=100, total=999.99, priority=high, region=south, customer_id=P002}
6. [1] TestOrder{id=O006, status=cancelled, customer_id=P005, amount=2, total=999.98, priority=low, discount=0, region=west, product_id=PROD005, date=2024-02-15}
7. [1] TestOrder{id=O007, product_id=PROD006, priority=urgent, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007}
8. [1] TestOrder{id=O008, discount=0, product_id=PROD002, amount=10, total=255, status=pending, region=south, customer_id=P010, date=2024-03-05, priority=normal}
9. [1] TestOrder{id=O009, amount=1, total=89.99, region=north, customer_id=P001, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007}
10. [1] TestOrder{id=O010, amount=1, date=2024-03-15, status=refunded, region=east, customer_id=P006, product_id=PROD001, total=75000, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, discount=100, total=999.99, priority=high, region=south}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, priority=urgent}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, total=255, status=pending, region=south}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, discount=10, product_id=PROD007, amount=1, total=89.99, region=north, customer_id=P001, date=2024-03-10, status=completed, priority=low}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, status=pending, discount=50, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, amount=1, total=25.5, discount=0, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, status=shipped, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, department=engineering, name=Bob, tags=senior, status=active, level=5, age=35, salary=75000, active=true, score=9.2}
3. [1] TestPerson{id=P003, department=hr, level=1, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, status=inactive, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, active=false}
6. [1] TestPerson{id=P006, level=1, name=Frank, age=0, active=true, status=active, department=qa, salary=-5000, score=0, tags=test}
7. [1] TestPerson{id=P007, name=Grace, score=10, level=9, age=65, salary=95000, active=true, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, level=1, name=Henry, age=18, salary=25000, active=false, tags=junior, score=5.5, status=inactive, department=support}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, department=engineering, active=true, score=8.7, status=active, level=6}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, name=X, score=6.5}
11. [1] TestOrder{id=O001, status=pending, discount=50, region=north, product_id=PROD001, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15}
12. [1] TestOrder{id=O002, status=confirmed, priority=low, customer_id=P002, amount=1, total=25.5, discount=0, region=south, product_id=PROD002, date=2024-01-20}
13. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, status=shipped, region=north}
14. [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, discount=100, total=999.99, priority=high, region=south}
16. [1] TestOrder{id=O006, priority=low, discount=0, region=west, product_id=PROD005, date=2024-02-15, status=cancelled, customer_id=P005, amount=2, total=999.98}
17. [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, priority=urgent}
18. [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, total=255, status=pending, region=south}
19. [1] TestOrder{id=O009, customer_id=P001, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, region=north}
20. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, total=75000, priority=urgent, discount=0, amount=1, date=2024-03-15, status=refunded, region=east}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales}
   - Fait 2: [1] TestOrder{id=O001, amount=2, total=1999.98, date=2024-01-15, status=pending, discount=50, region=north, product_id=PROD001, priority=normal, customer_id=P001}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, department=engineering, name=Bob, tags=senior, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, amount=1, total=25.5, discount=0, region=south}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, status=active, department=sales, level=2, name=Alice, age=25, salary=45000}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, status=shipped, region=north}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active, name=Diana, age=45, score=7.8}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, date=2024-02-05, amount=1, total=299.99, status=delivered, priority=normal, discount=0, region=east}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, status=active, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering, name=Bob, tags=senior}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, discount=100, total=999.99, priority=high, region=south}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, status=inactive, score=8, tags=employee, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, amount=2, total=999.98, priority=low, discount=0, region=west, product_id=PROD005, date=2024-02-15, status=cancelled}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, score=10, level=9, age=65, salary=95000, active=true, tags=executive, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, priority=urgent}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, level=1, name=Frank, age=0, active=true, status=active, department=qa, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O010, discount=0, amount=1, date=2024-03-15, status=refunded, region=east, customer_id=P006, product_id=PROD001, total=75000, priority=urgent}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, name=X, score=6.5}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, status=pending, region=south, customer_id=P010, date=2024-03-05, priority=normal, discount=0, product_id=PROD002}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, level=2, name=Alice}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, region=north, customer_id=P001}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/20 (50.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, active=true, score=9.2, department=engineering, name=Bob, tags=senior, status=active, level=5, age=35, salary=75000}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, name=Diana, age=45, score=7.8, department=marketing, level=7}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, status=inactive, score=8, tags=employee, department=sales, level=3}
6. [1] TestPerson{id=P006, score=0, tags=test, level=1, name=Frank, age=0, active=true, status=active, department=qa, salary=-5000}
7. [1] TestPerson{id=P007, status=active, department=management, name=Grace, score=10, level=9, age=65, salary=95000, active=true, tags=executive}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, tags=junior, score=5.5, status=inactive, department=support, level=1}
9. [1] TestPerson{id=P009, age=40, salary=68000, tags=senior, department=engineering, active=true, score=8.7, status=active, level=6, name=Ivy}
10. [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, name=X, score=6.5, age=22}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, status=active, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, department=engineering, name=Bob, tags=senior, status=active, level=5, age=35, salary=75000}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active, name=Diana, age=45, score=7.8}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, active=true, status=active, department=qa, salary=-5000, score=0, tags=test, level=1, name=Frank, age=0}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, status=active, department=management, name=Grace, score=10, level=9}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, active=true, score=8.7, status=active, level=6, name=Ivy, age=40, salary=68000, tags=senior, department=engineering}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, name=X, score=6.5}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active}
2. [1] TestPerson{id=P002, department=engineering, name=Bob, tags=senior, status=active, level=5, age=35, salary=75000, active=true, score=9.2}
3. [1] TestPerson{id=P003, department=hr, level=1, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active, name=Diana, age=45}
5. [1] TestPerson{id=P005, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, active=false, status=inactive, score=8}
6. [1] TestPerson{id=P006, salary=-5000, score=0, tags=test, level=1, name=Frank, age=0, active=true, status=active, department=qa}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, status=active, department=management, name=Grace, score=10, level=9}
8. [1] TestPerson{id=P008, score=5.5, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false, tags=junior}
9. [1] TestPerson{id=P009, department=engineering, active=true, score=8.7, status=active, level=6, name=Ivy, age=40, salary=68000, tags=senior}
10. [1] TestPerson{id=P010, department=intern, level=1, name=X, score=6.5, age=22, salary=28000, active=true, tags=temp, status=active}
11. [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, status=pending, discount=50, region=north, product_id=PROD001, priority=normal}
12. [1] TestOrder{id=O002, status=confirmed, priority=low, customer_id=P002, amount=1, total=25.5, discount=0, region=south, product_id=PROD002, date=2024-01-20}
13. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, status=shipped, region=north}
14. [1] TestOrder{id=O004, total=299.99, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, amount=1}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, discount=100, total=999.99, priority=high, region=south}
16. [1] TestOrder{id=O006, customer_id=P005, amount=2, total=999.98, priority=low, discount=0, region=west, product_id=PROD005, date=2024-02-15, status=cancelled}
17. [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, priority=urgent}
18. [1] TestOrder{id=O008, amount=10, total=255, status=pending, region=south, customer_id=P010, date=2024-03-05, priority=normal, discount=0, product_id=PROD002}
19. [1] TestOrder{id=O009, customer_id=P001, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, region=north}
20. [1] TestOrder{id=O010, amount=1, date=2024-03-15, status=refunded, region=east, customer_id=P006, product_id=PROD001, total=75000, priority=urgent, discount=0}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, status=inactive, score=8, tags=employee, department=sales, level=3, name=Eve, age=30}
   - Fait 2: [1] TestOrder{id=O006, amount=2, total=999.98, priority=low, discount=0, region=west, product_id=PROD005, date=2024-02-15, status=cancelled, customer_id=P005}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, level=9, age=65, salary=95000, active=true, tags=executive, status=active, department=management, name=Grace, score=10}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, priority=urgent}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, active=true, status=active, department=qa, salary=-5000, score=0, tags=test, level=1, name=Frank, age=0}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, region=east, customer_id=P006, product_id=PROD001, total=75000, priority=urgent, discount=0, amount=1, date=2024-03-15}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, status=active, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, status=pending, discount=50, region=north, product_id=PROD001, priority=normal}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, amount=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, department=engineering, name=Bob, tags=senior, status=active, level=5, age=35}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, discount=100, total=999.99, priority=high, region=south, customer_id=P002, product_id=PROD001, amount=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, name=X, score=6.5}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, status=pending, region=south, customer_id=P010, date=2024-03-05, priority=normal, discount=0, product_id=PROD002}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active}
   - Fait 2: [1] TestOrder{id=O009, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, region=north, customer_id=P001, date=2024-03-10}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, department=engineering, name=Bob, tags=senior, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, amount=1, total=25.5, discount=0, region=south}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O003, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, status=shipped, region=north, product_id=PROD003, amount=3}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/20 (50.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, status=active, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5}
2. [1] TestPerson{id=P002, name=Bob, tags=senior, status=active, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering}
3. [1] TestPerson{id=P003, score=6, tags=intern, status=inactive, department=hr, level=1, name=Charlie, age=16, salary=0, active=false}
4. [1] TestPerson{id=P004, tags=manager, status=active, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true}
5. [1] TestPerson{id=P005, age=30, salary=55000, active=false, status=inactive, score=8, tags=employee, department=sales, level=3, name=Eve}
6. [1] TestPerson{id=P006, score=0, tags=test, level=1, name=Frank, age=0, active=true, status=active, department=qa, salary=-5000}
7. [1] TestPerson{id=P007, name=Grace, score=10, level=9, age=65, salary=95000, active=true, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, score=5.5, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false, tags=junior}
9. [1] TestPerson{id=P009, active=true, score=8.7, status=active, level=6, name=Ivy, age=40, salary=68000, tags=senior, department=engineering}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, name=X, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, status=active, department=sales, level=2, name=Alice, age=25, salary=45000, active=true}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, tags=senior, status=active, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, level=1}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, score=10, level=9, age=65, salary=95000, active=true, tags=executive, status=active, department=management}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false, tags=junior}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, name=X, score=6.5}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, name=Diana, age=45, score=7.8, department=marketing, level=7}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, status=inactive, score=8, tags=employee, department=sales, level=3}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, active=true, score=8.7, status=active, level=6, name=Ivy, age=40, salary=68000, tags=senior, department=engineering}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, status=pending, discount=50, region=north, product_id=PROD001, priority=normal}
2. [1] TestOrder{id=O002, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, amount=1, total=25.5, discount=0, region=south, product_id=PROD002}
3. [1] TestOrder{id=O003, discount=15, customer_id=P001, status=shipped, region=north, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high}
4. [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05}
5. [1] TestOrder{id=O005, region=south, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, discount=100, total=999.99, priority=high}
6. [1] TestOrder{id=O006, customer_id=P005, amount=2, total=999.98, priority=low, discount=0, region=west, product_id=PROD005, date=2024-02-15, status=cancelled}
7. [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, priority=urgent}
8. [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, total=255, status=pending, region=south}
9. [1] TestOrder{id=O009, discount=10, product_id=PROD007, amount=1, total=89.99, region=north, customer_id=P001, date=2024-03-10, status=completed, priority=low}
10. [1] TestOrder{id=O010, amount=1, date=2024-03-15, status=refunded, region=east, customer_id=P006, product_id=PROD001, total=75000, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, priority=low, customer_id=P002, amount=1, total=25.5, discount=0, region=south, product_id=PROD002, date=2024-01-20, status=confirmed}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, status=shipped, region=north}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, discount=100, total=999.99, priority=high, region=south}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O006, priority=low, discount=0, region=west, product_id=PROD005, date=2024-02-15, status=cancelled, customer_id=P005, amount=2, total=999.98}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O008, total=255, status=pending, region=south, customer_id=P010, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, priority=urgent, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, total=89.99, region=north, customer_id=P001, date=2024-03-10, status=completed, priority=low, discount=10, product_id=PROD007, amount=1}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O010, region=east, customer_id=P006, product_id=PROD001, total=75000, priority=urgent, discount=0, amount=1, date=2024-03-15, status=refunded}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, status=pending, discount=50, region=north, product_id=PROD001, priority=normal}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales}
2. [1] TestPerson{id=P002, tags=senior, status=active, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering, name=Bob}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, status=inactive, score=8, tags=employee, department=sales, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, status=active, department=qa, salary=-5000, score=0, tags=test, level=1}
7. [1] TestPerson{id=P007, level=9, age=65, salary=95000, active=true, tags=executive, status=active, department=management, name=Grace, score=10}
8. [1] TestPerson{id=P008, score=5.5, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false, tags=junior}
9. [1] TestPerson{id=P009, salary=68000, tags=senior, department=engineering, active=true, score=8.7, status=active, level=6, name=Ivy, age=40}
10. [1] TestPerson{id=P010, name=X, score=6.5, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, active=false, tags=junior, score=5.5, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P010, name=X, score=6.5, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, level=1}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, level=2}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, level=1, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, name=Diana, age=45, score=7.8, department=marketing, level=7}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, active=true, score=8.7, status=active, level=6, name=Ivy, age=40, salary=68000, tags=senior, department=engineering}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, status=active, level=5, age=35, salary=75000, active=true, score=9.2, department=engineering, name=Bob, tags=senior}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, status=inactive, score=8, tags=employee, department=sales, level=3}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, status=active, department=qa, salary=-5000, score=0, tags=test, level=1, name=Frank}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, score=10, level=9, age=65, salary=95000, active=true, tags=executive, status=active, department=management}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 150
- **Faits traités**: 27
