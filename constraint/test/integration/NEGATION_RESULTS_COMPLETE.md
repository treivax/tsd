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

1. [1] TestPerson{id=P001, status=active, level=2, name=Alice, active=true, department=sales, age=25, salary=45000, score=8.5, tags=junior}
2. [1] TestPerson{id=P002, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior, status=active}
3. [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}
5. [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, age=0, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1, name=Frank}
7. [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, tags=junior, status=inactive, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000}
9. [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering, name=Ivy}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior, status=active}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15}
4. [1] TestOrder{id=O004, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0, total=299.99}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}
6. [1] TestOrder{id=O006, discount=0, region=west, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4}
8. [1] TestOrder{id=O008, customer_id=P010, total=255, discount=0, region=south, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10}
10. [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east, customer_id=P006}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high, discount=100}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O008, discount=0, region=south, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O010, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, discount=0, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, status=active, level=2, name=Alice, active=true, department=sales, age=25, salary=45000, score=8.5}
2. [1] TestPerson{id=P002, status=active, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, status=inactive, level=1, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0}
4. [1] TestPerson{id=P004, status=active, department=marketing, age=45, score=7.8, level=7, name=Diana, salary=85000, active=true, tags=manager}
5. [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive, name=Henry, age=18}
9. [1] TestPerson{id=P009, department=engineering, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active}
10. [1] TestPerson{id=P010, status=active, level=1, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, active=true, department=sales, age=25, salary=45000, score=8.5, tags=junior, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5, age=35}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9, age=65}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}
3. [1] TestOrder{id=O003, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15, customer_id=P001, amount=3, status=shipped}
4. [1] TestOrder{id=O004, amount=1, discount=0, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004}
5. [1] TestOrder{id=O005, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high}
6. [1] TestOrder{id=O006, discount=0, region=west, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west, product_id=PROD005, total=999.98, status=cancelled, priority=low}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O001, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending, priority=normal}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, discount=0, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, region=south, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, status=confirmed, priority=low, discount=0, total=25.5, region=south, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, brand=TechCorp, stock=50, supplier=TechSupply, name=Laptop, rating=4.5, category=electronics, price=999.99, available=true, keywords=computer}
2. [1] TestProduct{id=PROD002, supplier=TechSupply, name=Mouse, category=accessories, available=true, rating=4, keywords=peripheral, price=25.5, brand=TechCorp, stock=200}
3. [1] TestProduct{id=PROD003, name=Keyboard, price=75, rating=3.5, brand=KeyTech, stock=0, supplier=KeySupply, category=accessories, available=false, keywords=typing}
4. [1] TestProduct{id=PROD004, name=Monitor, rating=4.8, brand=ScreenPro, stock=30, supplier=ScreenSupply, category=electronics, price=299.99, available=true, keywords=display}
5. [1] TestProduct{id=PROD005, supplier=OldSupply, name=OldKeyboard, category=accessories, price=8.5, keywords=obsolete, brand=OldTech, stock=0, available=false, rating=2}
6. [1] TestProduct{id=PROD006, available=true, keywords=sound, brand=AudioMax, stock=75, name=Headphones, category=audio, price=150, rating=4.6, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, available=true, rating=3.8, keywords=video, brand=CamTech, stock=25, supplier=CamSupply, name=Webcam, category=electronics, price=89.99}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, category=electronics, price=999.99, available=true, keywords=computer, brand=TechCorp, stock=50, supplier=TechSupply, name=Laptop, rating=4.5}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, stock=200, supplier=TechSupply, name=Mouse, category=accessories, available=true, rating=4, keywords=peripheral, price=25.5, brand=TechCorp}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, name=Keyboard, price=75, rating=3.5, brand=KeyTech, stock=0, supplier=KeySupply, category=accessories, available=false, keywords=typing}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, category=electronics, price=299.99, available=true, keywords=display, name=Monitor, rating=4.8, brand=ScreenPro, stock=30, supplier=ScreenSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, name=Headphones, category=audio, price=150, rating=4.6, supplier=AudioSupply, available=true, keywords=sound, brand=AudioMax, stock=75}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, keywords=video, brand=CamTech, stock=25, supplier=CamSupply, name=Webcam, category=electronics, price=89.99, available=true, rating=3.8}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales}
2. [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, status=inactive, level=1, age=16, active=false, score=6, tags=intern, department=hr}
4. [1] TestPerson{id=P004, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8}
5. [1] TestPerson{id=P005, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales, age=30, active=false, score=8}
6. [1] TestPerson{id=P006, tags=test, status=active, department=qa, level=1, name=Frank, age=0, score=0, salary=-5000, active=true}
7. [1] TestPerson{id=P007, score=10, status=active, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true}
8. [1] TestPerson{id=P008, level=1, salary=25000, tags=junior, status=inactive, name=Henry, age=18, active=false, score=5.5, department=support}
9. [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering, name=Ivy}
10. [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1, age=16, active=false, score=6}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, active=true, tags=test, status=active, department=qa, level=1, name=Frank, age=0, score=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}
3. [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, region=north}
4. [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0}
5. [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99}
6. [1] TestOrder{id=O006, region=west, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south}
9. [1] TestOrder{id=O009, region=north, amount=1, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10}
10. [1] TestOrder{id=O010, discount=0, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O009, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, status=confirmed, priority=low, discount=0, total=25.5, region=south, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O006, region=west, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true}
2. [1] TestPerson{id=P002, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior, status=active}
3. [1] TestPerson{id=P003, level=1, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}
5. [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, level=1, salary=25000, tags=junior, status=inactive, name=Henry, age=18, active=false, score=5.5, department=support}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}
10. [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales, age=30}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15}
4. [1] TestOrder{id=O004, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0, total=299.99}
5. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high, discount=100}
6. [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west, product_id=PROD005}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south}
9. [1] TestOrder{id=O009, amount=1, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O009, region=north, amount=1, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south, customer_id=P002}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0, total=299.99}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true}
2. [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, status=inactive, level=1, age=16, active=false, score=6, tags=intern, department=hr}
4. [1] TestPerson{id=P004, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7, name=Diana, salary=85000, active=true}
5. [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, salary=-5000, active=true, tags=test, status=active, department=qa, level=1, name=Frank, age=0, score=0}
7. [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, active=true, status=active, department=engineering, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, department=support, level=1, salary=25000, tags=junior, status=inactive, name=Henry, age=18, active=false, score=5.5}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales, age=25, salary=45000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior, status=active}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales, age=30, active=false, score=8}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P006, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1, name=Frank, age=0}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, score=10, status=active, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, category=electronics, price=999.99, available=true, keywords=computer, brand=TechCorp, stock=50, supplier=TechSupply, name=Laptop, rating=4.5}
2. [1] TestProduct{id=PROD002, price=25.5, brand=TechCorp, stock=200, supplier=TechSupply, name=Mouse, category=accessories, available=true, rating=4, keywords=peripheral}
3. [1] TestProduct{id=PROD003, category=accessories, available=false, keywords=typing, name=Keyboard, price=75, rating=3.5, brand=KeyTech, stock=0, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, supplier=ScreenSupply, category=electronics, price=299.99, available=true, keywords=display, name=Monitor, rating=4.8, brand=ScreenPro, stock=30}
5. [1] TestProduct{id=PROD005, brand=OldTech, stock=0, available=false, rating=2, supplier=OldSupply, name=OldKeyboard, category=accessories, price=8.5, keywords=obsolete}
6. [1] TestProduct{id=PROD006, name=Headphones, category=audio, price=150, rating=4.6, supplier=AudioSupply, available=true, keywords=sound, brand=AudioMax, stock=75}
7. [1] TestProduct{id=PROD007, available=true, rating=3.8, keywords=video, brand=CamTech, stock=25, supplier=CamSupply, name=Webcam, category=electronics, price=89.99}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, rating=4.5, category=electronics, price=999.99, available=true, keywords=computer, brand=TechCorp, stock=50, supplier=TechSupply}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, price=25.5, brand=TechCorp, stock=200, supplier=TechSupply, name=Mouse, category=accessories, available=true, rating=4, keywords=peripheral}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, rating=3.5, brand=KeyTech, stock=0, supplier=KeySupply, category=accessories, available=false, keywords=typing, name=Keyboard, price=75}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, rating=4.8, brand=ScreenPro, stock=30, supplier=ScreenSupply, category=electronics, price=299.99, available=true, keywords=display}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, available=true, keywords=sound, brand=AudioMax, stock=75, name=Headphones, category=audio, price=150, rating=4.6, supplier=AudioSupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, available=true, rating=3.8, keywords=video, brand=CamTech, stock=25, supplier=CamSupply, name=Webcam, category=electronics, price=89.99}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales}
2. [1] TestPerson{id=P002, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob}
3. [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1}
4. [1] TestPerson{id=P004, status=active, department=marketing, age=45, score=7.8, level=7, name=Diana, salary=85000, active=true, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, status=inactive, department=sales, age=30, active=false, score=8, level=3}
6. [1] TestPerson{id=P006, tags=test, status=active, department=qa, level=1, name=Frank, age=0, score=0, salary=-5000, active=true}
7. [1] TestPerson{id=P007, level=9, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management}
8. [1] TestPerson{id=P008, salary=25000, tags=junior, status=inactive, name=Henry, age=18, active=false, score=5.5, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}
10. [1] TestPerson{id=P010, department=intern, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, active=true, department=sales, age=25, salary=45000, score=8.5, tags=junior, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, salary=0, status=inactive, level=1, age=16, active=false, score=6, tags=intern, department=hr}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, department=sales, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active, name=Grace}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, active=true, status=active, department=engineering, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, level=5, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}
3. [1] TestOrder{id=O003, total=225, date=2024-02-01, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, discount=0, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east}
5. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high, discount=100}
6. [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west}
7. [1] TestOrder{id=O007, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50}
8. [1] TestOrder{id=O008, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south, product_id=PROD002, amount=10, date=2024-03-05}
9. [1] TestOrder{id=O009, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10, customer_id=P001}
10. [1] TestOrder{id=O010, discount=0, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales}
2. [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5}
3. [1] TestPerson{id=P003, level=1, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive}
4. [1] TestPerson{id=P004, department=marketing, age=45, score=7.8, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, active=true, tags=test, status=active, department=qa, level=1, name=Frank, age=0, score=0, salary=-5000}
7. [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, status=inactive, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}
10. [1] TestPerson{id=P010, status=active, level=1, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000}
11. [1] TestOrder{id=O001, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending}
12. [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south, customer_id=P002, product_id=PROD002}
13. [1] TestOrder{id=O003, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15}
14. [1] TestOrder{id=O004, amount=1, discount=0, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004}
15. [1] TestOrder{id=O005, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002}
16. [1] TestOrder{id=O006, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west, product_id=PROD005, total=999.98, status=cancelled, priority=low}
17. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}
18. [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south}
19. [1] TestOrder{id=O009, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99}
20. [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, level=6, active=true, status=active, department=engineering, name=Ivy, age=40, salary=68000, score=8.7}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, region=north}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P007, status=active, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, status=active, level=2, name=Alice, active=true, department=sales, age=25, salary=45000, score=8.5}
   - Fait 2: [1] TestOrder{id=O005, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, discount=0, region=south, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, active=true, status=active, department=engineering, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, tags=junior, status=inactive, name=Henry, age=18, active=false, score=5.5, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1, age=16, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O005, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P006, active=true, tags=test, status=active, department=qa, level=1, name=Frank, age=0, score=0, salary=-5000}
   - Fait 2: [1] TestOrder{id=O007, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P001, active=true, department=sales, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice}
   - Fait 2: [1] TestOrder{id=O009, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales, age=30}
   - Fait 2: [1] TestOrder{id=O004, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0, total=299.99, date=2024-02-05, status=delivered, priority=normal}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O006, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west, product_id=PROD005, total=999.98, status=cancelled}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales, age=25}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5, age=35}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P010, level=1, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000, status=active}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales, age=30}
   - Fait 2: [1] TestOrder{id=O008, total=255, discount=0, region=south, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, salary=0, status=inactive, level=1, age=16, active=false, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales, age=25}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, region=north}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P006, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1, name=Frank, age=0}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, status=inactive, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P005, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales, age=30, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O002, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south, customer_id=P002, product_id=PROD002, amount=1}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O002, priority=low, discount=0, total=25.5, region=south, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, region=north}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, active=true, department=sales, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O008, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south, product_id=PROD002, amount=10, date=2024-03-05}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O003, total=225, date=2024-02-01, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P006, active=true, tags=test, status=active, department=qa, level=1, name=Frank, age=0, score=0, salary=-5000}
   - Fait 2: [1] TestOrder{id=O005, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering, name=Ivy}
   - Fait 2: [1] TestOrder{id=O006, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west, product_id=PROD005, total=999.98, status=cancelled}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P003, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1, age=16, active=false}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, salary=0, status=inactive, level=1, age=16, active=false, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west, product_id=PROD005, total=999.98}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P009, active=true, status=active, department=engineering, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, status=inactive, department=sales, age=30, active=false, score=8, level=3}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10, customer_id=P001}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P002, level=5, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P004, status=active, department=marketing, age=45, score=7.8, level=7, name=Diana, salary=85000, active=true, tags=manager}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, active=true, department=sales, age=25, salary=45000, score=8.5, tags=junior, status=active}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1, active=true}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high, discount=100}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, level=1, name=Frank, age=0, score=0, salary=-5000, active=true, tags=test, status=active}
   - Fait 2: [1] TestOrder{id=O001, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending, priority=normal}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O004, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0, total=299.99, date=2024-02-05, status=delivered}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, tags=junior, status=inactive, name=Henry, age=18, active=false, score=5.5, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5}
   - Fait 2: [1] TestOrder{id=O009, region=north, amount=1, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P004, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7, name=Diana, salary=85000}
   - Fait 2: [1] TestOrder{id=O010, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, age=45, score=7.8, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active}
   - Fait 2: [1] TestOrder{id=O001, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, active=true, tags=test, status=active, department=qa, level=1, name=Frank, age=0, score=0}
   - Fait 2: [1] TestOrder{id=O003, total=225, date=2024-02-01, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P002, status=active, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior}
   - Fait 2: [1] TestOrder{id=O004, discount=0, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, active=true, department=sales, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west, product_id=PROD005, total=999.98, status=cancelled, priority=low}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, level=1, name=Frank, age=0, score=0, salary=-5000, active=true, tags=test, status=active}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9, age=65}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P005, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales, age=30, active=false}
   - Fait 2: [1] TestOrder{id=O002, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south, customer_id=P002, product_id=PROD002, amount=1}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales, age=25, salary=45000}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, active=true, tags=test, status=active, department=qa, level=1, name=Frank, age=0, score=0}
   - Fait 2: [1] TestOrder{id=O009, amount=1, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P008, department=support, level=1, salary=25000, tags=junior, status=inactive, name=Henry, age=18, active=false, score=5.5}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, region=north}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O003, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, discount=0, region=south, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P007, level=9, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales, age=25}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, status=inactive, level=1, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, discount=0, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, status=inactive, department=sales, age=30, active=false, score=8, level=3, name=Eve, salary=55000}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south, product_id=PROD002, amount=10}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P009, level=6, active=true, status=active, department=engineering, name=Ivy, age=40, salary=68000, score=8.7, tags=senior}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, salary=0, status=inactive, level=1, age=16, active=false, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2}
   - Fait 2: [1] TestOrder{id=O003, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15, customer_id=P001, amount=3, status=shipped}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O004, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0, total=299.99, date=2024-02-05}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, tags=junior, status=inactive, name=Henry, age=18, active=false, score=5.5, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O007, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, level=1, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east, customer_id=P006}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, level=6, active=true, status=active, department=engineering, name=Ivy, age=40, salary=68000, score=8.7}
   - Fait 2: [1] TestOrder{id=O001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, tags=junior, status=inactive, name=Henry, age=18, active=false, score=5.5, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P006, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1, name=Frank, age=0}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, discount=0, region=south, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O009, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, status=active, level=2, name=Alice, active=true, department=sales, age=25, salary=45000, score=8.5}
2. [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, status=inactive, level=1, age=16, active=false, score=6, tags=intern, department=hr}
4. [1] TestPerson{id=P004, department=marketing, age=45, score=7.8, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, salary=-5000, active=true, tags=test, status=active, department=qa, level=1, name=Frank, age=0, score=0}
7. [1] TestPerson{id=P007, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, tags=senior, level=6, active=true, status=active, department=engineering, name=Ivy, age=40, salary=68000, score=8.7}
10. [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, status=active, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, age=45, score=7.8, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active, name=Grace}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, level=2, name=Alice, active=true, department=sales, age=25, salary=45000, score=8.5, tags=junior}
2. [1] TestPerson{id=P002, status=active, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, status=inactive, level=1, age=16, active=false, score=6, tags=intern, department=hr}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}
5. [1] TestPerson{id=P005, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales, age=30, active=false}
6. [1] TestPerson{id=P006, salary=-5000, active=true, tags=test, status=active, department=qa, level=1, name=Frank, age=0, score=0}
7. [1] TestPerson{id=P007, level=9, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management}
8. [1] TestPerson{id=P008, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive, name=Henry}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}
10. [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1}
11. [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001}
12. [1] TestOrder{id=O002, total=25.5, region=south, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0}
13. [1] TestOrder{id=O003, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15}
14. [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0}
15. [1] TestOrder{id=O005, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1}
16. [1] TestOrder{id=O006, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west, product_id=PROD005, total=999.98, status=cancelled, priority=low}
17. [1] TestOrder{id=O007, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50}
18. [1] TestOrder{id=O008, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south, product_id=PROD002, amount=10, date=2024-03-05}
19. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10}
20. [1] TestOrder{id=O010, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive, name=Henry, age=18}
   - Fait 2: [1] TestOrder{id=O002, priority=low, discount=0, total=25.5, region=south, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp, department=intern, name=X, age=22}
   - Fait 2: [1] TestOrder{id=O003, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15, customer_id=P001, amount=3, status=shipped}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, status=active, level=2, name=Alice, active=true, department=sales, age=25, salary=45000, score=8.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high, discount=100}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5, age=35}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{id=O008, discount=0, region=south, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1, age=16}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, status=active, department=qa, level=1, name=Frank, age=0, score=0, salary=-5000, active=true}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000}
   - Fait 2: [1] TestOrder{id=O003, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales, age=25}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, discount=0, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1, active=true}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P008, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive, name=Henry}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, age=45, score=7.8, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P009, level=6, active=true, status=active, department=engineering, name=Ivy, age=40, salary=68000, score=8.7, tags=senior}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, department=sales, age=30, active=false, score=8, level=3, name=Eve}
   - Fait 2: [1] TestOrder{id=O006, date=2024-02-15, discount=0, region=west, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P005, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales, age=30, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O003, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15, customer_id=P001, amount=3}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, salary=0, status=inactive, level=1, age=16, active=false, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, discount=0, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering, name=Ivy, age=40}
   - Fait 2: [1] TestOrder{id=O009, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, age=45, score=7.8, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active}
   - Fait 2: [1] TestOrder{id=O003, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior}
   - Fait 2: [1] TestOrder{id=O004, amount=1, discount=0, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, active=true, department=sales, age=25, salary=45000, score=8.5, tags=junior, status=active}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales}
   - Fait 2: [1] TestOrder{id=O002, region=south, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O002, region=south, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive, name=Henry, age=18, active=false}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}
   - Fait 2: [1] TestOrder{id=O007, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P001, active=true, department=sales, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2}
   - Fait 2: [1] TestOrder{id=O008, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south, product_id=PROD002}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1, age=16}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south, customer_id=P002, product_id=PROD002}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, discount=0, region=south, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O009, amount=1, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering, name=Ivy}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O002, discount=0, total=25.5, region=south, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P008, level=1, salary=25000, tags=junior, status=inactive, name=Henry, age=18, active=false, score=5.5, department=support}
   - Fait 2: [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15, customer_id=P001}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P004, age=45, score=7.8, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing}
   - Fait 2: [1] TestOrder{id=O004, discount=0, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O008, discount=0, region=south, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P007, status=active, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O009, amount=1, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P006, level=1, name=Frank, age=0, score=0, salary=-5000, active=true, tags=test, status=active, department=qa}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1, total=75000}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, active=true, tags=test, status=active, department=qa, level=1, name=Frank, age=0, score=0}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south, customer_id=P002}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O003, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15, customer_id=P001, amount=3, status=shipped, priority=high}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, status=inactive, department=sales, age=30, active=false, score=8, level=3}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive, name=Henry, age=18}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high, discount=100}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, salary=0, status=inactive, level=1, age=16, active=false, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west, product_id=PROD005, total=999.98, status=cancelled, priority=low}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, status=active, level=2, name=Alice, active=true, department=sales, age=25, salary=45000, score=8.5}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15, customer_id=P001}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P007, score=10, status=active, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true}
   - Fait 2: [1] TestOrder{id=O004, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0, total=299.99, date=2024-02-05}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P004, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7, name=Diana, salary=85000}
   - Fait 2: [1] TestOrder{id=O009, region=north, amount=1, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O001, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending, priority=normal}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P001, status=active, level=2, name=Alice, active=true, department=sales, age=25, salary=45000, score=8.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O006, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west, product_id=PROD005, total=999.98, status=cancelled}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales}
   - Fait 2: [1] TestOrder{id=O007, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O007, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P001, status=active, level=2, name=Alice, active=true, department=sales, age=25, salary=45000, score=8.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O001, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O001, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending, priority=normal}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P005, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales, age=30, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P006, active=true, tags=test, status=active, department=qa, level=1, name=Frank, age=0, score=0, salary=-5000}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering, name=Ivy}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, discount=0, region=south, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O008, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south, product_id=PROD002}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P009, status=active, department=engineering, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, discount=15, customer_id=P001, amount=3, status=shipped, priority=high, region=north}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, product_id=PROD001, priority=high, discount=100}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, amount=1, discount=0}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P006, age=0, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1, name=Frank}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, level=1, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0}
   - Fait 2: [1] TestOrder{id=O009, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P007, level=9, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1, age=16}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, region=south, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, status=active, department=qa, level=1, name=Frank, age=0, score=0, salary=-5000, active=true}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, discount=0, region=south, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, product_id=PROD001}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales}
2. [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5}
3. [1] TestPerson{id=P003, department=hr, name=Charlie, salary=0, status=inactive, level=1, age=16, active=false, score=6, tags=intern}
4. [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7, name=Diana}
5. [1] TestPerson{id=P005, status=inactive, department=sales, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee}
6. [1] TestPerson{id=P006, name=Frank, age=0, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering, name=Ivy, age=40, salary=68000}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, level=1, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0, status=inactive}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, status=inactive, department=sales, age=30, active=false, score=8, level=3}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior, status=active}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, age=45, score=7.8, level=7}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, tags=junior, status=inactive, name=Henry, age=18, active=false, score=5.5, department=support, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001, amount=2, status=pending, priority=normal, discount=50}
2. [1] TestOrder{id=O002, status=confirmed, priority=low, discount=0, total=25.5, region=south, customer_id=P002, product_id=PROD002, amount=1, date=2024-01-20}
3. [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15, customer_id=P001}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, discount=0, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}
6. [1] TestOrder{id=O006, amount=2, date=2024-02-15, discount=0, region=west, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south}
9. [1] TestOrder{id=O009, amount=1, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north}
10. [1] TestOrder{id=O010, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, amount=2, status=pending, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, region=north, product_id=PROD001}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, region=south, customer_id=P002, product_id=PROD002}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, amount=3, status=shipped, priority=high, region=north, product_id=PROD003, total=225, date=2024-02-01, discount=15}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, priority=high, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, customer_id=P005, amount=2, date=2024-02-15, discount=0, region=west}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, status=pending, priority=normal, customer_id=P010, total=255, discount=0, region=south, product_id=PROD002, amount=10, date=2024-03-05}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, region=east, customer_id=P006, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, discount=0, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, region=north, amount=1, status=completed, priority=low, discount=10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales, age=25, salary=45000}
2. [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob, active=true, department=engineering, level=5}
3. [1] TestPerson{id=P003, status=inactive, level=1, age=16, active=false, score=6, tags=intern, department=hr, name=Charlie, salary=0}
4. [1] TestPerson{id=P004, department=marketing, age=45, score=7.8, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, name=Eve, salary=55000, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9, age=65}
8. [1] TestPerson{id=P008, active=false, score=5.5, department=support, level=1, salary=25000, tags=junior, status=inactive, name=Henry, age=18}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}
10. [1] TestPerson{id=P010, salary=28000, status=active, level=1, active=true, score=6.5, tags=temp, department=intern, name=X, age=22}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, active=true, department=engineering, level=5, age=35, salary=75000, score=9.2, tags=senior, status=active, name=Bob}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, name=Charlie, salary=0, status=inactive, level=1, age=16, active=false, score=6}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, status=inactive, department=sales, age=30, active=false, score=8, level=3}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, score=8.7, tags=senior, level=6, active=true, status=active, department=engineering}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, name=Alice, active=true, department=sales}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, age=45, score=7.8, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P006, score=0, salary=-5000, active=true, tags=test, status=active, department=qa, level=1, name=Frank, age=0}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, status=active, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, status=inactive, name=Henry, age=18, active=false, score=5.5, department=support, level=1, salary=25000}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, department=intern, name=X, age=22, salary=28000, status=active, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
