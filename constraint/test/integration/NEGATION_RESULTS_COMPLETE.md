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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}
3. [1] TestPerson{id=P003, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0}
4. [1] TestPerson{id=P004, tags=manager, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, level=3, age=30, active=false, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1, age=0, active=true, score=0}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, active=false, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, status=active, department=engineering, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior}
10. [1] TestPerson{id=P010, level=1, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true, score=10, level=9}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, active=false, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending, product_id=PROD001, amount=2}
2. [1] TestOrder{id=O002, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east}
5. [1] TestOrder{id=O005, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10}
6. [1] TestOrder{id=O006, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled, product_id=PROD005, amount=2}
7. [1] TestOrder{id=O007, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600}
8. [1] TestOrder{id=O008, status=pending, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10}
9. [1] TestOrder{id=O009, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O008, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south, customer_id=P010}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending, product_id=PROD001, amount=2}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O002, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, score=8.5, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}
3. [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1}
4. [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}
5. [1] TestPerson{id=P005, score=8, level=3, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000}
6. [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1}
7. [1] TestPerson{id=P007, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000}
8. [1] TestPerson{id=P008, active=false, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering, name=Ivy, age=40, salary=68000}
10. [1] TestPerson{id=P010, tags=temp, status=active, department=intern, level=1, name=X, age=22, salary=28000, active=true, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3, age=30, active=false}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, status=pending, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north}
2. [1] TestOrder{id=O002, priority=low, discount=0, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed}
3. [1] TestOrder{id=O003, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225}
4. [1] TestOrder{id=O004, discount=0, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal}
5. [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled}
7. [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north}
8. [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low}
10. [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east, customer_id=P006, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending, product_id=PROD001, amount=2}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, discount=10, region=north, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, status=pending, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, rating=4.5, keywords=computer, stock=50, available=true, brand=TechCorp, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, keywords=peripheral, brand=TechCorp, supplier=TechSupply, category=accessories, price=25.5, stock=200, name=Mouse, available=true, rating=4}
3. [1] TestProduct{id=PROD003, category=accessories, price=75, keywords=typing, supplier=KeySupply, available=false, rating=3.5, brand=KeyTech, stock=0, name=Keyboard}
4. [1] TestProduct{id=PROD004, price=299.99, rating=4.8, stock=30, supplier=ScreenSupply, available=true, keywords=display, brand=ScreenPro, name=Monitor, category=electronics}
5. [1] TestProduct{id=PROD005, category=accessories, price=8.5, rating=2, keywords=obsolete, brand=OldTech, supplier=OldSupply, name=OldKeyboard, available=false, stock=0}
6. [1] TestProduct{id=PROD006, name=Headphones, category=audio, available=true, keywords=sound, brand=AudioMax, price=150, rating=4.6, stock=75, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, brand=CamTech, supplier=CamSupply, name=Webcam, category=electronics, available=true, rating=3.8, keywords=video, stock=25, price=89.99}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD007, keywords=video, stock=25, price=89.99, brand=CamTech, supplier=CamSupply, name=Webcam, category=electronics, available=true, rating=3.8}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, rating=4.5, keywords=computer, stock=50, available=true, brand=TechCorp, supplier=TechSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD002, stock=200, name=Mouse, available=true, rating=4, keywords=peripheral, brand=TechCorp, supplier=TechSupply, category=accessories, price=25.5}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD003, category=accessories, price=75, keywords=typing, supplier=KeySupply, available=false, rating=3.5, brand=KeyTech, stock=0, name=Keyboard}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD004, supplier=ScreenSupply, available=true, keywords=display, brand=ScreenPro, name=Monitor, category=electronics, price=299.99, rating=4.8, stock=30}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD006, category=audio, available=true, keywords=sound, brand=AudioMax, price=150, rating=4.6, stock=75, supplier=AudioSupply, name=Headphones}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, department=sales, active=true, score=8.5, level=2, name=Alice, age=25, salary=45000, tags=junior}
2. [1] TestPerson{id=P002, level=5, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35}
3. [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1, name=Charlie}
4. [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}
5. [1] TestPerson{id=P005, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, level=1, age=0, active=true, score=0, tags=test, status=active, department=qa}
7. [1] TestPerson{id=P007, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true, score=10, level=9, age=65}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, priority=normal, discount=50, region=north, customer_id=P001, status=pending, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15}
2. [1] TestOrder{id=O002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1, total=25.5, region=south}
3. [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003}
4. [1] TestOrder{id=O004, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped}
8. [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south}
9. [1] TestOrder{id=O009, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1, total=25.5, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, total=225, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O006, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled, product_id=PROD005}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, discount=0, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O008, status=pending, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5}
2. [1] TestPerson{id=P002, level=5, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35}
3. [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1}
4. [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager, salary=85000, active=true}
5. [1] TestPerson{id=P005, level=3, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8}
6. [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1}
7. [1] TestPerson{id=P007, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000}
8. [1] TestPerson{id=P008, active=false, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, status=active, department=engineering, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior}
10. [1] TestPerson{id=P010, level=1, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending}
2. [1] TestOrder{id=O002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1, total=25.5, region=south}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15, customer_id=P001}
4. [1] TestOrder{id=O004, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100}
6. [1] TestOrder{id=O006, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled, product_id=PROD005, amount=2}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0}
9. [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O005, region=south, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north, date=2024-03-10, status=completed}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, active=true, score=8.5, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active}
2. [1] TestPerson{id=P002, name=Bob, age=35, level=5, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering}
3. [1] TestPerson{id=P003, age=16, salary=0, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr}
4. [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager, salary=85000, active=true}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, level=3, age=30, active=false, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1, age=0, active=true}
7. [1] TestPerson{id=P007, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true, score=10, level=9}
8. [1] TestPerson{id=P008, active=false, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering, name=Ivy}
10. [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, status=active, department=intern, level=1, name=X, age=22, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, status=active, department=qa, name=Frank, salary=-5000, level=1, age=0, active=true, score=0, tags=test}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, active=true, score=8.5, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, name=Bob, age=35, level=5, salary=75000, active=true, score=9.2, tags=senior}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1, name=Charlie}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true, score=10, level=9, age=65}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, rating=4.5, keywords=computer, stock=50, available=true, brand=TechCorp, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, supplier=TechSupply, category=accessories, price=25.5, stock=200, name=Mouse, available=true, rating=4, keywords=peripheral, brand=TechCorp}
3. [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, price=75, keywords=typing, supplier=KeySupply, available=false, rating=3.5, brand=KeyTech, stock=0}
4. [1] TestProduct{id=PROD004, name=Monitor, category=electronics, price=299.99, rating=4.8, stock=30, supplier=ScreenSupply, available=true, keywords=display, brand=ScreenPro}
5. [1] TestProduct{id=PROD005, category=accessories, price=8.5, rating=2, keywords=obsolete, brand=OldTech, supplier=OldSupply, name=OldKeyboard, available=false, stock=0}
6. [1] TestProduct{id=PROD006, category=audio, available=true, keywords=sound, brand=AudioMax, price=150, rating=4.6, stock=75, supplier=AudioSupply, name=Headphones}
7. [1] TestProduct{id=PROD007, keywords=video, stock=25, price=89.99, brand=CamTech, supplier=CamSupply, name=Webcam, category=electronics, available=true, rating=3.8}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD003, available=false, rating=3.5, brand=KeyTech, stock=0, name=Keyboard, category=accessories, price=75, keywords=typing, supplier=KeySupply}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, category=electronics, price=299.99, rating=4.8, stock=30, supplier=ScreenSupply, available=true, keywords=display, brand=ScreenPro}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD006, price=150, rating=4.6, stock=75, supplier=AudioSupply, name=Headphones, category=audio, available=true, keywords=sound, brand=AudioMax}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD007, name=Webcam, category=electronics, available=true, rating=3.8, keywords=video, stock=25, price=89.99, brand=CamTech, supplier=CamSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD001, stock=50, available=true, brand=TechCorp, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, rating=4.5, keywords=computer}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD002, keywords=peripheral, brand=TechCorp, supplier=TechSupply, category=accessories, price=25.5, stock=200, name=Mouse, available=true, rating=4}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2, name=Alice}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}
3. [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1}
4. [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}
5. [1] TestPerson{id=P005, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, level=1, age=0, active=true, score=0, tags=test, status=active, department=qa}
7. [1] TestPerson{id=P007, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true, score=10, level=9}
8. [1] TestPerson{id=P008, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, score=10, level=9, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5, salary=75000, active=true, score=9.2}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager, salary=85000, active=true}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending}
2. [1] TestOrder{id=O002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1, total=25.5, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}
4. [1] TestOrder{id=O004, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}
5. [1] TestOrder{id=O005, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south}
6. [1] TestOrder{id=O006, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled, product_id=PROD005, amount=2}
7. [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north}
8. [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south}
9. [1] TestOrder{id=O009, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, discount=10, region=north, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, department=sales, active=true, score=8.5, level=2, name=Alice, age=25, salary=45000, tags=junior}
2. [1] TestPerson{id=P002, department=engineering, name=Bob, age=35, level=5, salary=75000, active=true, score=9.2, tags=senior, status=active}
3. [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1}
4. [1] TestPerson{id=P004, age=45, tags=manager, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana}
5. [1] TestPerson{id=P005, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3, age=30, active=false}
6. [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support}
9. [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering, name=Ivy, age=40}
10. [1] TestPerson{id=P010, department=intern, level=1, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active}
11. [1] TestOrder{id=O001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending, product_id=PROD001, amount=2, total=1999.98}
12. [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0}
13. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}
14. [1] TestOrder{id=O004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004}
15. [1] TestOrder{id=O005, product_id=PROD001, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100, customer_id=P002}
16. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled}
17. [1] TestOrder{id=O007, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006}
18. [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0}
19. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low}
20. [1] TestOrder{id=O010, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, age=45, tags=manager, salary=85000, active=true, score=7.8, status=active, department=marketing}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6}
   - Fait 2: [1] TestOrder{id=O006, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled, product_id=PROD005}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, status=active, department=sales, active=true, score=8.5, level=2, name=Alice, age=25, salary=45000, tags=junior}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, status=active, department=intern, level=1, name=X, age=22, salary=28000, active=true, score=6.5}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, status=inactive, active=false, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000}
   - Fait 2: [1] TestOrder{id=O004, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5, salary=75000, active=true, score=9.2}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east, customer_id=P006, amount=1}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, name=Bob, age=35, level=5, salary=75000, active=true, score=9.2, tags=senior}
   - Fait 2: [1] TestOrder{id=O001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending, product_id=PROD001, amount=2, total=1999.98}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management, name=Grace}
   - Fait 2: [1] TestOrder{id=O009, discount=10, region=north, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P005, score=8, level=3, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east, customer_id=P006}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O010, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1}
   - Fait 2: [1] TestOrder{id=O001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending, product_id=PROD001}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3, age=30, active=false}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6}
   - Fait 2: [1] TestOrder{id=O003, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager, salary=85000, active=true}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15, customer_id=P001}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P010, level=1, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3}
   - Fait 2: [1] TestOrder{id=O010, priority=urgent, product_id=PROD001, discount=0, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P007, department=management, name=Grace, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active}
   - Fait 2: [1] TestOrder{id=O010, priority=urgent, product_id=PROD001, discount=0, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering, name=Ivy, age=40}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O002, discount=0, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, score=8, level=3, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1, total=25.5, region=south}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, name=Bob, age=35, level=5, salary=75000, active=true, score=9.2, tags=senior}
   - Fait 2: [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, active=true, score=8.5, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, level=1, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled, product_id=PROD005, amount=2}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, level=1, age=0, active=true, score=0, tags=test, status=active, department=qa}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, level=1, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank}
   - Fait 2: [1] TestOrder{id=O008, discount=0, product_id=PROD002, amount=10, status=pending, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1, total=25.5, region=south}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, level=3, age=30, active=false, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, customer_id=P005, status=cancelled, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management, name=Grace}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, customer_id=P005, status=cancelled, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, level=5, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5, salary=75000, active=true, score=9.2}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, status=pending, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, status=pending, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P007, score=10, level=9, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, score=8, level=3, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve}
   - Fait 2: [1] TestOrder{id=O005, region=south, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P008, active=false, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive}
   - Fait 2: [1] TestOrder{id=O007, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false}
   - Fait 2: [1] TestOrder{id=O008, status=pending, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, status=cancelled, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O003, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2, name=Alice, age=25}
   - Fait 2: [1] TestOrder{id=O006, priority=low, discount=0, region=west, customer_id=P005, status=cancelled, product_id=PROD005, amount=2, total=999.98, date=2024-02-15}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1, age=0, active=true, score=0}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1, name=Charlie}
   - Fait 2: [1] TestOrder{id=O003, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, active=false, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P002, level=5, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35}
   - Fait 2: [1] TestOrder{id=O006, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled, product_id=PROD005}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O007, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P007, level=9, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, level=3, age=30, active=false, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O008, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south, customer_id=P010, total=255, date=2024-03-05}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}
   - Fait 2: [1] TestOrder{id=O002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P005, score=8, level=3, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000}
   - Fait 2: [1] TestOrder{id=O003, priority=high, region=north, total=225, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, status=shipped}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1, age=0, active=true, score=0}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0}
   - Fait 2: [1] TestOrder{id=O007, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O008, discount=0, product_id=PROD002, amount=10, status=pending, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P004, age=45, tags=manager, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, discount=0, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O002, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1, total=25.5}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P006, status=active, department=qa, name=Frank, salary=-5000, level=1, age=0, active=true, score=0, tags=test}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, status=active, department=sales, active=true, score=8.5, level=2, name=Alice, age=25, salary=45000}
   - Fait 2: [1] TestOrder{id=O008, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2}
   - Fait 2: [1] TestOrder{id=O003, region=north, total=225, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1, name=X}
   - Fait 2: [1] TestOrder{id=O003, region=north, total=225, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5, salary=75000, active=true, score=9.2}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P008, active=false, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, status=cancelled, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low, customer_id=P001}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O001, region=north, customer_id=P001, status=pending, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P008, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P007, department=management, name=Grace, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active}
   - Fait 2: [1] TestOrder{id=O001, region=north, customer_id=P001, status=pending, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P005, score=8, level=3, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, level=1, age=0, active=true, score=0, tags=test, status=active, department=qa}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager, salary=85000}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south, amount=1}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O005, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P002, level=5, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35}
   - Fait 2: [1] TestOrder{id=O009, discount=10, region=north, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}
3. [1] TestPerson{id=P003, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0}
4. [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}
5. [1] TestPerson{id=P005, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3, age=30}
6. [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1}
7. [1] TestPerson{id=P007, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true, score=10, level=9}
8. [1] TestPerson{id=P008, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true, score=10, level=9}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1, name=X, age=22}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, status=active, department=sales, active=true, score=8.5, level=2, name=Alice, age=25, salary=45000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, name=Bob, age=35, level=5, salary=75000, active=true, score=9.2, tags=senior, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}
3. [1] TestPerson{id=P003, age=16, salary=0, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr}
4. [1] TestPerson{id=P004, name=Diana, age=45, tags=manager, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, department=sales, name=Eve, salary=55000, score=8, level=3, age=30, active=false, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, level=1, age=0, active=true, score=0, tags=test, status=active, department=qa}
7. [1] TestPerson{id=P007, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true, score=10, level=9}
8. [1] TestPerson{id=P008, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support}
9. [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6}
10. [1] TestPerson{id=P010, status=active, department=intern, level=1, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp}
11. [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending, product_id=PROD001, amount=2}
12. [1] TestOrder{id=O002, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20}
13. [1] TestOrder{id=O003, total=225, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north}
14. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100}
16. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled}
17. [1] TestOrder{id=O007, region=north, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50}
18. [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south}
19. [1] TestOrder{id=O009, discount=10, region=north, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99}
20. [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east, customer_id=P006, amount=1}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager, salary=85000, active=true, score=7.8}
   - Fait 2: [1] TestOrder{id=O002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, tags=temp, status=active, department=intern, level=1, name=X, age=22, salary=28000, active=true}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15, customer_id=P001}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, active=false, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, active=false, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1, name=X}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east, customer_id=P006, amount=1, total=75000}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O001, discount=50, region=north, customer_id=P001, status=pending, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, age=45, tags=manager, salary=85000, active=true, score=7.8, status=active, department=marketing}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, status=active, department=intern, level=1, name=X, age=22, salary=28000}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, level=3, age=30, active=false, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, name=Frank, salary=-5000, level=1, age=0, active=true, score=0, tags=test, status=active}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, score=8, level=3, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve}
   - Fait 2: [1] TestOrder{id=O010, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P008, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support}
   - Fait 2: [1] TestOrder{id=O007, region=north, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, age=16, salary=0, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern}
   - Fait 2: [1] TestOrder{id=O002, discount=0, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P010, status=active, department=intern, level=1, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, discount=50, region=north, customer_id=P001, status=pending, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2}
   - Fait 2: [1] TestOrder{id=O007, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P007, status=active, department=management, name=Grace, salary=95000, active=true, score=10, level=9, age=65, tags=executive}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P009, status=active, department=engineering, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O009, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north, date=2024-03-10}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O003, discount=15, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5, salary=75000, active=true}
   - Fait 2: [1] TestOrder{id=O009, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north, date=2024-03-10}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2, name=Alice, age=25}
   - Fait 2: [1] TestOrder{id=O010, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management, name=Grace}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, discount=0, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, status=pending, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P003, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1, name=Charlie, active=false}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager, salary=85000}
   - Fait 2: [1] TestOrder{id=O005, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2}
   - Fait 2: [1] TestOrder{id=O008, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2}
   - Fait 2: [1] TestOrder{id=O009, region=north, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support, level=1, name=Henry, age=18}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P007, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true, score=10, level=9}
   - Fait 2: [1] TestOrder{id=O009, discount=10, region=north, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, discount=0, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}
   - Fait 2: [1] TestOrder{id=O001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending, product_id=PROD001}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, level=1, age=0, active=true, score=0, tags=test, status=active, department=qa}
   - Fait 2: [1] TestOrder{id=O004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1, age=0, active=true, score=0}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1, total=25.5, region=south}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, status=active, department=intern, level=1, name=X, age=22, salary=28000}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P010, status=active, department=intern, level=1, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, age=45, tags=manager, salary=85000, active=true, score=7.8, status=active, department=marketing}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east, customer_id=P006}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending, product_id=PROD001, amount=2}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, level=3, age=30, active=false, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, level=1, age=0, active=true, score=0, tags=test, status=active, department=qa}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P005, level=3, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, status=pending, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P006, status=active, department=qa, name=Frank, salary=-5000, level=1, age=0, active=true, score=0, tags=test}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P002, level=5, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3}
   - Fait 2: [1] TestOrder{id=O005, region=south, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales}
   - Fait 2: [1] TestOrder{id=O006, priority=low, discount=0, region=west, customer_id=P005, status=cancelled, product_id=PROD005, amount=2, total=999.98, date=2024-02-15}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P007, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true, score=10, level=9}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, status=pending}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, active=true, score=8.5, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O010, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P007, score=10, level=9, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O006, priority=low, discount=0, region=west, customer_id=P005, status=cancelled, product_id=PROD005, amount=2, total=999.98, date=2024-02-15}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south, customer_id=P010, total=255}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0}
   - Fait 2: [1] TestOrder{id=O009, discount=10, region=north, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1, total=89.99}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}
   - Fait 2: [1] TestOrder{id=O009, amount=1, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2}
   - Fait 2: [1] TestOrder{id=O003, discount=15, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P010, level=1, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, level=1, age=0, active=true, score=0, tags=test, status=active, department=qa}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3, age=30, active=false, tags=employee}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, level=3, age=30, active=false, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales, active=true, score=8.5, level=2}
2. [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5, salary=75000}
3. [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1}
4. [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager, salary=85000, active=true}
5. [1] TestPerson{id=P005, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3}
6. [1] TestPerson{id=P006, score=0, tags=test, status=active, department=qa, name=Frank, salary=-5000, level=1, age=0, active=true}
7. [1] TestPerson{id=P007, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true, score=10, level=9}
8. [1] TestPerson{id=P008, active=false, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, age=16, salary=0, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, age=45, tags=manager, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, score=10, level=9, age=65, tags=executive, status=active, department=management}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, active=false, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, age=35, level=5}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, level=3, age=30, active=false, tags=employee, status=inactive, department=sales}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering, name=Ivy, age=40}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1, name=X}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, priority=normal, discount=50, region=north, customer_id=P001, status=pending, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15}
2. [1] TestOrder{id=O002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, amount=1, total=25.5, region=south}
3. [1] TestOrder{id=O003, region=north, total=225, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high}
4. [1] TestOrder{id=O004, discount=0, total=299.99, date=2024-02-05, region=east, customer_id=P004, product_id=PROD004, amount=1, status=delivered, priority=normal}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south, amount=1, total=999.99, priority=high, discount=100}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005, status=cancelled}
7. [1] TestOrder{id=O007, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600}
8. [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0, product_id=PROD002, amount=10, status=pending, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low}
10. [1] TestOrder{id=O010, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, priority=high, region=north, total=225, date=2024-02-01, discount=15}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, amount=1, status=delivered, priority=normal, discount=0, total=299.99, date=2024-02-05, region=east, customer_id=P004}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, total=999.99, priority=high, discount=100, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, region=south}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O006, status=cancelled, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, region=west, customer_id=P005}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, priority=urgent, discount=50, region=north, amount=4, total=600, date=2024-03-01, status=shipped}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, total=89.99, discount=10, region=north, date=2024-03-10, status=completed, priority=low, customer_id=P001, product_id=PROD007, amount=1}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, product_id=PROD001, discount=0, region=east}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O001, region=north, customer_id=P001, status=pending, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, priority=normal, discount=50}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, region=south, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, discount=0}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, region=south, customer_id=P010, total=255, date=2024-03-05, priority=normal, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, status=active, department=sales, active=true, score=8.5, level=2, name=Alice, age=25, salary=45000}
2. [1] TestPerson{id=P002, status=active, department=engineering, name=Bob, age=35, level=5, salary=75000, active=true, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, score=6, tags=intern, department=hr, age=16, salary=0, status=inactive, level=1, name=Charlie, active=false}
4. [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, status=active, department=marketing, level=7, name=Diana, age=45, tags=manager}
5. [1] TestPerson{id=P005, level=3, age=30, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, level=1, age=0, active=true, score=0, tags=test, status=active, department=qa}
7. [1] TestPerson{id=P007, status=active, department=management, name=Grace, salary=95000, active=true, score=10, level=9, age=65, tags=executive}
8. [1] TestPerson{id=P008, active=false, tags=junior, department=support, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, name=Ivy, age=40, salary=68000, active=true, score=8.7, level=6}
10. [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, status=active, department=intern, level=1, name=X, age=22, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, level=1, name=Henry, age=18, salary=25000, score=5.5, status=inactive, active=false, tags=junior, department=support}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, active=true, score=8.7, level=6, tags=senior, status=active, department=engineering, name=Ivy}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, active=true, score=8.5, level=2, name=Alice, age=25, salary=45000, tags=junior, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, status=inactive, level=1, name=Charlie, active=false, score=6, tags=intern, department=hr}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, active=false, tags=employee, status=inactive, department=sales, name=Eve, salary=55000, score=8, level=3, age=30}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, salary=-5000, level=1, age=0, active=true, score=0, tags=test, status=active, department=qa}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, score=6.5, tags=temp, status=active, department=intern, level=1}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, name=Bob, age=35, level=5, salary=75000, active=true, score=9.2, tags=senior, status=active}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, name=Diana, age=45, tags=manager, salary=85000, active=true, score=7.8, status=active}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P007, level=9, age=65, tags=executive, status=active, department=management, name=Grace, salary=95000, active=true, score=10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
