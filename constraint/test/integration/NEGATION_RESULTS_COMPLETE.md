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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior, department=sales}
2. [1] TestPerson{id=P002, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5, name=Bob, age=35}
3. [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, level=1, salary=0, tags=intern, status=inactive, department=hr}
4. [1] TestPerson{id=P004, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active, name=Diana, salary=85000}
5. [1] TestPerson{id=P005, tags=employee, department=sales, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10, name=Grace}
8. [1] TestPerson{id=P008, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5}
9. [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active, name=Diana, salary=85000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, tags=employee, department=sales, age=30, score=8, status=inactive, level=3}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior, department=sales, name=Alice, age=25}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5, name=Bob, age=35}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10, name=Grace}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}
2. [1] TestOrder{id=O002, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}
3. [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225, discount=15, customer_id=P001}
4. [1] TestOrder{id=O004, total=299.99, date=2024-02-05, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10}
6. [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2}
7. [1] TestOrder{id=O007, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped}
8. [1] TestOrder{id=O008, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0}
9. [1] TestOrder{id=O009, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1, total=89.99, status=completed}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O005, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O008, discount=0, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, discount=10, region=north, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1, total=89.99, status=completed, priority=low}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O010, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0, region=east, customer_id=P006, total=75000}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225, discount=15}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, level=2, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5}
2. [1] TestPerson{id=P002, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5, name=Bob, age=35}
3. [1] TestPerson{id=P003, department=hr, name=Charlie, age=16, active=false, score=6, level=1, salary=0, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, department=marketing, level=7, age=45, active=true, status=active, name=Diana, salary=85000, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales}
6. [1] TestPerson{id=P006, status=active, department=qa, level=1, age=0, tags=test, name=Frank, salary=-5000, active=true, score=0}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10}
8. [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1}
9. [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}
10. [1] TestPerson{id=P010, department=intern, name=X, score=6.5, level=1, age=22, salary=28000, active=true, tags=temp, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, status=active, level=2, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, age=45, active=true, status=active, name=Diana, salary=85000, score=7.8, tags=manager}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales, age=30, score=8}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, status=active, department=management, level=9, active=true, score=10, name=Grace, age=65, salary=95000, tags=executive}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, active=true, score=8.7, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}
2. [1] TestOrder{id=O002, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1}
3. [1] TestOrder{id=O003, amount=3, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, customer_id=P002, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10}
6. [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}
7. [1] TestOrder{id=O007, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50}
8. [1] TestOrder{id=O008, product_id=PROD002, discount=0, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south}
9. [1] TestOrder{id=O009, discount=10, region=north, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1, total=89.99, status=completed, priority=low}
10. [1] TestOrder{id=O010, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0, region=east, customer_id=P006}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, priority=high, region=north, amount=3, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05, customer_id=P004}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O005, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O006, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10, customer_id=P001}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, rating=4.5, keywords=computer, brand=TechCorp, name=Laptop, category=electronics, price=999.99, stock=50, supplier=TechSupply, available=true}
2. [1] TestProduct{id=PROD002, supplier=TechSupply, category=accessories, stock=200, name=Mouse, price=25.5, available=true, rating=4, keywords=peripheral, brand=TechCorp}
3. [1] TestProduct{id=PROD003, available=false, rating=3.5, keywords=typing, brand=KeyTech, supplier=KeySupply, stock=0, name=Keyboard, category=accessories, price=75}
4. [1] TestProduct{id=PROD004, stock=30, name=Monitor, category=electronics, price=299.99, available=true, keywords=display, brand=ScreenPro, supplier=ScreenSupply, rating=4.8}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, price=8.5, rating=2, stock=0, supplier=OldSupply, category=accessories, available=false, keywords=obsolete, brand=OldTech}
6. [1] TestProduct{id=PROD006, name=Headphones, category=audio, price=150, available=true, rating=4.6, keywords=sound, stock=75, brand=AudioMax, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, price=89.99, rating=3.8, keywords=video, brand=CamTech, stock=25, supplier=CamSupply, name=Webcam, available=true, category=electronics}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD007, brand=CamTech, stock=25, supplier=CamSupply, name=Webcam, available=true, category=electronics, price=89.99, rating=3.8, keywords=video}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD001, available=true, rating=4.5, keywords=computer, brand=TechCorp, name=Laptop, category=electronics, price=999.99, stock=50, supplier=TechSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD002, supplier=TechSupply, category=accessories, stock=200, name=Mouse, price=25.5, available=true, rating=4, keywords=peripheral, brand=TechCorp}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, price=75, available=false, rating=3.5, keywords=typing, brand=KeyTech, supplier=KeySupply, stock=0}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, category=electronics, price=299.99, available=true, keywords=display, brand=ScreenPro, supplier=ScreenSupply, rating=4.8, stock=30}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD006, stock=75, brand=AudioMax, supplier=AudioSupply, name=Headphones, category=audio, price=150, available=true, rating=4.6, keywords=sound}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior, department=sales}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5}
3. [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active}
5. [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}
7. [1] TestPerson{id=P007, score=10, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true}
8. [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support}
9. [1] TestPerson{id=P009, department=engineering, level=6, salary=68000, active=true, score=8.7, name=Ivy, age=40, tags=senior, status=active}
10. [1] TestPerson{id=P010, name=X, score=6.5, level=1, age=22, salary=28000, active=true, tags=temp, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, level=1, salary=0, tags=intern, status=inactive, department=hr, name=Charlie}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, score=0, status=active, department=qa, level=1, age=0, tags=test, name=Frank, salary=-5000, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98}
2. [1] TestOrder{id=O002, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1}
3. [1] TestOrder{id=O003, status=shipped, priority=high, region=north, amount=3, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south}
6. [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}
7. [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north}
8. [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0, customer_id=P010}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10}
10. [1] TestOrder{id=O010, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O008, priority=normal, region=south, product_id=PROD002, discount=0, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O010, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O003, status=shipped, priority=high, region=north, amount=3, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5, status=active}
2. [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000}
3. [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, level=1, salary=0, tags=intern, status=inactive, department=hr}
4. [1] TestPerson{id=P004, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active, name=Diana}
5. [1] TestPerson{id=P005, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales, age=30, score=8, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10}
8. [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support}
9. [1] TestPerson{id=P009, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7, name=Ivy, age=40, tags=senior}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1, salary=0, tags=intern}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225, discount=15}
4. [1] TestOrder{id=O004, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05, customer_id=P004, product_id=PROD004, amount=1, status=delivered}
5. [1] TestOrder{id=O005, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed}
6. [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent}
8. [1] TestOrder{id=O008, region=south, product_id=PROD002, discount=0, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, amount=1}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1, total=89.99, status=completed}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, amount=1, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5, status=active}
2. [1] TestPerson{id=P002, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering}
3. [1] TestPerson{id=P003, department=hr, name=Charlie, age=16, active=false, score=6, level=1, salary=0, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active, name=Diana, salary=85000}
5. [1] TestPerson{id=P005, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales, age=30}
6. [1] TestPerson{id=P006, status=active, department=qa, level=1, age=0, tags=test, name=Frank, salary=-5000, active=true, score=0}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10}
8. [1] TestPerson{id=P008, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5}
9. [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}
10. [1] TestPerson{id=P010, name=X, score=6.5, level=1, age=22, salary=28000, active=true, tags=temp, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, level=1, salary=0, tags=intern, status=inactive, department=hr}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, status=active, department=qa, level=1, age=0, tags=test, name=Frank, salary=-5000, active=true, score=0}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, active=true, status=active, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior, department=sales, name=Alice}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, stock=50, supplier=TechSupply, available=true, rating=4.5, keywords=computer, brand=TechCorp}
2. [1] TestProduct{id=PROD002, name=Mouse, price=25.5, available=true, rating=4, keywords=peripheral, brand=TechCorp, supplier=TechSupply, category=accessories, stock=200}
3. [1] TestProduct{id=PROD003, brand=KeyTech, supplier=KeySupply, stock=0, name=Keyboard, category=accessories, price=75, available=false, rating=3.5, keywords=typing}
4. [1] TestProduct{id=PROD004, rating=4.8, stock=30, name=Monitor, category=electronics, price=299.99, available=true, keywords=display, brand=ScreenPro, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, category=accessories, available=false, keywords=obsolete, brand=OldTech, name=OldKeyboard, price=8.5, rating=2, stock=0, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, available=true, rating=4.6, keywords=sound, stock=75, brand=AudioMax, supplier=AudioSupply, name=Headphones, category=audio, price=150}
7. [1] TestProduct{id=PROD007, category=electronics, price=89.99, rating=3.8, keywords=video, brand=CamTech, stock=25, supplier=CamSupply, name=Webcam, available=true}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, price=999.99, stock=50, supplier=TechSupply, available=true, rating=4.5, keywords=computer, brand=TechCorp, name=Laptop, category=electronics}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, price=25.5, available=true, rating=4, keywords=peripheral, brand=TechCorp, supplier=TechSupply, category=accessories, stock=200}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, category=accessories, price=75, available=false, rating=3.5, keywords=typing, brand=KeyTech, supplier=KeySupply, stock=0, name=Keyboard}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, available=true, keywords=display, brand=ScreenPro, supplier=ScreenSupply, rating=4.8, stock=30, name=Monitor, category=electronics, price=299.99}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, keywords=sound, stock=75, brand=AudioMax, supplier=AudioSupply, name=Headphones, category=audio, price=150, available=true, rating=4.6}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, category=electronics, price=89.99, rating=3.8, keywords=video, brand=CamTech, stock=25, supplier=CamSupply, name=Webcam, available=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2}
2. [1] TestPerson{id=P002, age=35, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5, name=Bob}
3. [1] TestPerson{id=P003, department=hr, name=Charlie, age=16, active=false, score=6, level=1, salary=0, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, age=45, active=true, status=active, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7}
5. [1] TestPerson{id=P005, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales, age=30, score=8}
6. [1] TestPerson{id=P006, age=0, tags=test, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, level=9, active=true, score=10, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false}
9. [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7, name=Ivy, age=40}
10. [1] TestPerson{id=P010, status=active, department=intern, name=X, score=6.5, level=1, age=22, salary=28000, active=true, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, active=true, status=active, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, age=0, tags=test, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior, department=sales}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001}
2. [1] TestOrder{id=O002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225, discount=15}
4. [1] TestOrder{id=O004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05, customer_id=P004, product_id=PROD004}
5. [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, amount=1}
6. [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}
7. [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north}
8. [1] TestOrder{id=O008, region=south, product_id=PROD002, discount=0, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10}
10. [1] TestOrder{id=O010, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O007, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O001, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O003, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior, department=sales}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5}
3. [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
4. [1] TestPerson{id=P004, active=true, status=active, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45}
5. [1] TestPerson{id=P005, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales, age=30, score=8, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}
7. [1] TestPerson{id=P007, department=management, level=9, active=true, score=10, name=Grace, age=65, salary=95000, tags=executive, status=active}
8. [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1}
9. [1] TestPerson{id=P009, score=8.7, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true}
10. [1] TestPerson{id=P010, status=active, department=intern, name=X, score=6.5, level=1, age=22, salary=28000, active=true, tags=temp}
11. [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}
12. [1] TestOrder{id=O002, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed}
13. [1] TestOrder{id=O003, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3}
14. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05}
15. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south}
16. [1] TestOrder{id=O006, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low}
17. [1] TestOrder{id=O007, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50}
18. [1] TestOrder{id=O008, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0}
19. [1] TestOrder{id=O009, date=2024-03-10, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007}
20. [1] TestOrder{id=O010, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225, discount=15}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, status=active, level=2, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P003, level=1, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O006, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, department=management, level=9, active=true, score=10, name=Grace, age=65, salary=95000, tags=executive, status=active}
   - Fait 2: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, level=6, salary=68000, active=true, score=8.7, name=Ivy, age=40, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0, customer_id=P010}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10, name=Grace}
   - Fait 2: [1] TestOrder{id=O003, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support, age=18, salary=25000}
   - Fait 2: [1] TestOrder{id=O007, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O008, discount=0, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, tags=employee, department=sales, age=30, score=8, status=inactive, level=3, name=Eve}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P005, active=false, tags=employee, department=sales, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, amount=1, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, level=1, salary=0, tags=intern, status=inactive, department=hr}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P004, status=active, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true}
   - Fait 2: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, department=marketing, level=7, age=45, active=true, status=active, name=Diana, salary=85000, score=7.8}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225, discount=15}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P001, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, name=X, score=6.5, level=1, age=22, salary=28000, active=true, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P010, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1, age=22, salary=28000}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0, region=east}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, level=1, salary=0, tags=intern, status=inactive, department=hr}
   - Fait 2: [1] TestOrder{id=O002, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, discount=0, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior, department=sales, name=Alice, age=25}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P001, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O010, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering, active=true}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, department=sales, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P010, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1, age=22, salary=28000}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, status=active, level=2, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225, discount=15, customer_id=P001}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7, name=Ivy, age=40}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, tags=employee, department=sales, age=30, score=8, status=inactive, level=3, name=Eve}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1, total=89.99, status=completed}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P003, level=1, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, amount=1, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active, name=Diana}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{id=O009, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10, customer_id=P001}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5}
   - Fait 2: [1] TestOrder{id=O010, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5, name=Bob, age=35}
   - Fait 2: [1] TestOrder{id=O001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, name=X, score=6.5, level=1, age=22, salary=28000, active=true, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O003, priority=high, region=north, amount=3, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, department=marketing, level=7, age=45, active=true, status=active, name=Diana, salary=85000, score=7.8}
   - Fait 2: [1] TestOrder{id=O009, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1, total=89.99}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, status=active, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O010, amount=1, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7, name=Ivy, age=40}
   - Fait 2: [1] TestOrder{id=O010, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test, name=Frank}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1}
   - Fait 2: [1] TestOrder{id=O007, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north, amount=4, date=2024-03-01}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active, name=Diana, salary=85000}
   - Fait 2: [1] TestOrder{id=O002, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, tags=employee, department=sales, age=30, score=8, status=inactive, level=3, name=Eve}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test, name=Frank}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0, customer_id=P010, amount=10}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, amount=1, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, status=active, level=2, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10, name=Grace, age=65}
   - Fait 2: [1] TestOrder{id=O005, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}
   - Fait 2: [1] TestOrder{id=O006, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, department=sales, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false}
   - Fait 2: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O003, priority=high, region=north, amount=3, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225, discount=15}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, status=active, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O004, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05, customer_id=P004, product_id=PROD004, amount=1}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P006, status=active, department=qa, level=1, age=0, tags=test, name=Frank, salary=-5000, active=true, score=0}
   - Fait 2: [1] TestOrder{id=O005, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P007, score=10, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O007, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0, customer_id=P010, amount=10, total=255}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P003, level=1, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O003, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5, name=Bob, age=35}
   - Fait 2: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P007, department=management, level=9, active=true, score=10, name=Grace, age=65, salary=95000, tags=executive, status=active}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, level=1, age=0, tags=test, name=Frank, salary=-5000, active=true, score=0, status=active}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, amount=1, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, tags=employee, department=sales, age=30, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O001, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, status=active, department=management, level=9, active=true, score=10, name=Grace, age=65, salary=95000}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05, customer_id=P004}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P001, status=active, level=2, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5}
   - Fait 2: [1] TestOrder{id=O007, discount=50, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test, name=Frank}
   - Fait 2: [1] TestOrder{id=O002, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P001, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225, discount=15}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support, age=18, salary=25000}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O008, discount=0, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2, active=true}
2. [1] TestPerson{id=P002, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering}
3. [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
4. [1] TestPerson{id=P004, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active, name=Diana}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, tags=employee, department=sales, age=30, score=8, status=inactive, level=3}
6. [1] TestPerson{id=P006, department=qa, level=1, age=0, tags=test, name=Frank, salary=-5000, active=true, score=0, status=active}
7. [1] TestPerson{id=P007, level=9, active=true, score=10, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, active=false, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry}
9. [1] TestPerson{id=P009, level=6, salary=68000, active=true, score=8.7, name=Ivy, age=40, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1, age=22, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7, name=Ivy, age=40}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, status=active, level=2, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, level=2, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5}
3. [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1, salary=0}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active}
5. [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10}
8. [1] TestPerson{id=P008, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support, age=18, salary=25000}
9. [1] TestPerson{id=P009, score=8.7, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true}
10. [1] TestPerson{id=P010, name=X, score=6.5, level=1, age=22, salary=28000, active=true, tags=temp, status=active, department=intern}
11. [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}
12. [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south}
13. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225, discount=15}
14. [1] TestOrder{id=O004, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05, customer_id=P004, product_id=PROD004, amount=1, status=delivered}
15. [1] TestOrder{id=O005, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100}
16. [1] TestOrder{id=O006, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98}
17. [1] TestOrder{id=O007, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50}
18. [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0, customer_id=P010, amount=10}
19. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10}
20. [1] TestOrder{id=O010, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0, region=east}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, status=active, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north, amount=4}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O003, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, tags=employee, department=sales, age=30, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, discount=0, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{id=O008, region=south, product_id=PROD002, discount=0, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P005, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales, age=30}
   - Fait 2: [1] TestOrder{id=O002, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, level=1, salary=0, tags=intern, status=inactive, department=hr, name=Charlie}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225, discount=15}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support, age=18}
   - Fait 2: [1] TestOrder{id=O007, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P006, level=1, age=0, tags=test, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1, total=89.99, status=completed}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10, name=Grace}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05, customer_id=P004}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O006, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O003, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P001, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, status=active, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O005, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{id=O009, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1, total=89.99}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0, region=east}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1, salary=0, tags=intern}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, tags=employee, department=sales, age=30, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O005, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O008, priority=normal, region=south, product_id=PROD002, discount=0, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P001, status=active, level=2, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O005, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5}
   - Fait 2: [1] TestOrder{id=O009, discount=10, region=north, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1, total=89.99, status=completed, priority=low}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0, region=east}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, department=sales, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, amount=1, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P009, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7, name=Ivy}
   - Fait 2: [1] TestOrder{id=O006, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P010, level=1, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5}
   - Fait 2: [1] TestOrder{id=O006, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, status=active, department=management, level=9, active=true, score=10, name=Grace, age=65, salary=95000}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P006, age=0, tags=test, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0, region=east}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5, name=Bob, age=35}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P006, status=active, department=qa, level=1, age=0, tags=test, name=Frank, salary=-5000, active=true, score=0}
   - Fait 2: [1] TestOrder{id=O002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{id=O002, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P010, status=active, department=intern, name=X, score=6.5, level=1, age=22, salary=28000, active=true, tags=temp}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, priority=low}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P006, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test, name=Frank, salary=-5000}
   - Fait 2: [1] TestOrder{id=O003, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{id=O003, priority=high, region=north, amount=3, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225, discount=15}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P007, department=management, level=9, active=true, score=10, name=Grace, age=65, salary=95000, tags=executive, status=active}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225, discount=15}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5, level=1, age=22}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north, amount=4}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P006, level=1, age=0, tags=test, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{id=O007, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north, amount=4, date=2024-03-01}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P005, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales, age=30}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P010, status=active, department=intern, name=X, score=6.5, level=1, age=22, salary=28000, active=true, tags=temp}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O004, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05, customer_id=P004, product_id=PROD004, amount=1}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, department=sales, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east, total=299.99}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, status=active, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O008, priority=normal, region=south, product_id=PROD002, discount=0, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P003, level=1, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O008, status=pending, priority=normal, region=south, product_id=PROD002, discount=0, customer_id=P010, amount=10, total=255, date=2024-03-05}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering, active=true}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P001, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O006, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, age=45, active=true, status=active, name=Diana, salary=85000, score=7.8, tags=manager}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P006, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test, name=Frank, salary=-5000}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0, region=east}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0, region=east}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, amount=1, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P001, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O002, region=south, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225, discount=15, customer_id=P001, product_id=PROD003}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P005, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales, age=30, score=8, status=inactive}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O009, region=north, product_id=PROD007, date=2024-03-10, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P010, level=1, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0, region=east}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225, discount=15}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P010, name=X, score=6.5, level=1, age=22, salary=28000, active=true, tags=temp, status=active, department=intern}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, date=2024-02-10}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior, department=sales, name=Alice}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P001, status=active, level=2, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5}
   - Fait 2: [1] TestOrder{id=O001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000}
   - Fait 2: [1] TestOrder{id=O006, region=west, product_id=PROD005, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10, name=Grace}
   - Fait 2: [1] TestOrder{id=O010, amount=1, status=refunded, discount=0, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, tags=employee, department=sales, age=30, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{id=O004, status=delivered, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05, customer_id=P004, product_id=PROD004, amount=1}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, level=1, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X}
   - Fait 2: [1] TestOrder{id=O004, region=east, total=299.99, date=2024-02-05, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}
   - Fait 2: [1] TestOrder{id=O004, region=east, total=299.99, date=2024-02-05, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2}
2. [1] TestPerson{id=P002, department=engineering, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior, status=active}
3. [1] TestPerson{id=P003, level=1, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7, age=45, active=true, status=active}
5. [1] TestPerson{id=P005, active=false, tags=employee, department=sales, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa, level=1, age=0, tags=test}
7. [1] TestPerson{id=P007, score=10, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true}
8. [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1}
9. [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7, name=Ivy, age=40}
10. [1] TestPerson{id=P010, status=active, department=intern, name=X, score=6.5, level=1, age=22, salary=28000, active=true, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P010, name=X, score=6.5, level=1, age=22, salary=28000, active=true, tags=temp, status=active, department=intern}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, status=active, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, active=true, tags=junior, department=sales, name=Alice, age=25, salary=45000, score=8.5, status=active, level=2}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, level=1, salary=0, tags=intern, status=inactive, department=hr}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, score=10, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false, score=5.5, department=support}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6, salary=68000, active=true, score=8.7}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending}
2. [1] TestOrder{id=O002, region=south, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0}
3. [1] TestOrder{id=O003, status=shipped, priority=high, region=north, amount=3, total=225, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01}
4. [1] TestOrder{id=O004, priority=normal, discount=0, region=east, total=299.99, date=2024-02-05, customer_id=P004, product_id=PROD004, amount=1, status=delivered}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south}
6. [1] TestOrder{id=O006, total=999.98, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005}
7. [1] TestOrder{id=O007, priority=urgent, customer_id=P007, product_id=PROD006, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped}
8. [1] TestOrder{id=O008, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007, date=2024-03-10}
10. [1] TestOrder{id=O010, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, discount=15, customer_id=P001, product_id=PROD003, date=2024-02-01, status=shipped, priority=high, region=north, amount=3, total=225}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, region=east}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, customer_id=P002, status=confirmed, priority=high, discount=100, region=south}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O006, date=2024-02-15, customer_id=P005, amount=2, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, total=999.98}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, total=600, discount=50, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O010, region=east, customer_id=P006, total=75000, date=2024-03-15, priority=urgent, product_id=PROD001, amount=1, status=refunded, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, region=south, product_id=PROD002, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, date=2024-03-10, customer_id=P001, amount=1, total=89.99, status=completed, priority=low, discount=10, region=north, product_id=PROD007}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, priority=low, customer_id=P002, date=2024-01-20, status=confirmed, discount=0, region=south}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior, department=sales, name=Alice, age=25}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, tags=senior, status=active, department=engineering, active=true, score=9.2, level=5}
3. [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1}
4. [1] TestPerson{id=P004, age=45, active=true, status=active, name=Diana, salary=85000, score=7.8, tags=manager, department=marketing, level=7}
5. [1] TestPerson{id=P005, score=8, status=inactive, level=3, name=Eve, salary=55000, active=false, tags=employee, department=sales, age=30}
6. [1] TestPerson{id=P006, score=0, status=active, department=qa, level=1, age=0, tags=test, name=Frank, salary=-5000, active=true}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9, active=true, score=10}
8. [1] TestPerson{id=P008, active=false, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry}
9. [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, name=Ivy, age=40, tags=senior, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, name=X, score=6.5, level=1, age=22, salary=28000, active=true, tags=temp, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P006, level=1, age=0, tags=test, name=Frank, salary=-5000, active=true, score=0, status=active, department=qa}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, department=support, age=18, salary=25000, tags=junior, status=inactive, level=1, name=Henry, active=false}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, level=1, age=22, salary=28000, active=true, tags=temp, status=active, department=intern, name=X, score=6.5}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, status=active, level=2, active=true, tags=junior, department=sales, name=Alice}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, name=Charlie, age=16, active=false, score=6, level=1, salary=0}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, age=45, active=true, status=active, name=Diana, salary=85000, score=7.8, tags=manager}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, name=Grace, age=65, salary=95000, tags=executive, status=active, department=management, level=9}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, level=6, salary=68000, active=true, score=8.7, name=Ivy, age=40, tags=senior, status=active}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, active=true, score=9.2, level=5, name=Bob, age=35, salary=75000, tags=senior, status=active}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P005, active=false, tags=employee, department=sales, age=30, score=8, status=inactive, level=3, name=Eve, salary=55000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
