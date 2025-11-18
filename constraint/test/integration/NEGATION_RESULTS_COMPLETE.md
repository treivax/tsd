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

1. [1] TestPerson{id=P001, active=true, score=8.5, status=active, department=sales, level=2, age=25, tags=junior, name=Alice, salary=45000}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, level=5, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering}
3. [1] TestPerson{id=P003, age=16, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, active=true, tags=manager, status=active, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000}
5. [1] TestPerson{id=P005, score=8, department=sales, level=3, active=false, tags=employee, status=inactive, name=Eve, age=30, salary=55000}
6. [1] TestPerson{id=P006, salary=-5000, active=true, score=0, tags=test, name=Frank, age=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, status=active, name=Grace, salary=95000, active=true, tags=executive, department=management, level=9, age=65, score=10}
8. [1] TestPerson{id=P008, score=5.5, tags=junior, department=support, level=1, age=18, status=inactive, name=Henry, salary=25000, active=false}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, score=8.7, tags=senior, department=engineering, level=6, salary=68000, status=active}
10. [1] TestPerson{id=P010, name=X, age=22, score=6.5, tags=temp, level=1, salary=28000, active=true, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1, age=16}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, tags=executive, department=management, level=9, age=65, score=10, status=active}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, department=support, level=1, age=18, status=inactive, name=Henry, salary=25000, active=false, score=5.5, tags=junior}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, status=active, name=Ivy, age=40, active=true, score=8.7, tags=senior, department=engineering, level=6, salary=68000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, score=8.5, status=active, department=sales, level=2, age=25, tags=junior, name=Alice}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, salary=75000, level=5}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, name=Diana, age=45, score=7.8, department=marketing, level=7}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, score=8, department=sales, level=3, active=false, tags=employee, status=inactive, name=Eve}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, status=active, department=intern, name=X, age=22, score=6.5, tags=temp, level=1, salary=28000, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, date=2024-01-15, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, priority=low, amount=1, discount=0, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, region=north, amount=3, total=225, date=2024-02-01, priority=high, discount=15}
4. [1] TestOrder{id=O004, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0, customer_id=P004, amount=1, total=299.99}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, priority=high, discount=100}
6. [1] TestOrder{id=O006, discount=0, amount=2, date=2024-02-15, priority=low, region=west, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled}
7. [1] TestOrder{id=O007, amount=4, priority=urgent, region=north, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50}
8. [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, region=south, product_id=PROD002, amount=10, total=255, priority=normal, discount=0}
9. [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, priority=low, discount=10}
10. [1] TestOrder{id=O010, date=2024-03-15, status=refunded, discount=0, customer_id=P006, product_id=PROD001, amount=1, priority=urgent, region=east, total=75000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, priority=low, amount=1, discount=0, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, status=shipped, region=north, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, discount=100, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, priority=high}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, priority=low, discount=10, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, customer_id=P001}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O010, region=east, total=75000, date=2024-03-15, status=refunded, discount=0, customer_id=P006, product_id=PROD001, amount=1, priority=urgent}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, discount=50, product_id=PROD001, total=1999.98, date=2024-01-15, region=north, customer_id=P001, amount=2, status=pending, priority=normal}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, discount=50, amount=4, priority=urgent, region=north, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, discount=0, customer_id=P010, date=2024-03-05, status=pending, region=south, product_id=PROD002, amount=10, total=255, priority=normal}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, age=25, tags=junior, name=Alice, salary=45000, active=true, score=8.5, status=active, department=sales}
2. [1] TestPerson{id=P002, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, salary=75000, level=5}
3. [1] TestPerson{id=P003, name=Charlie, active=false, tags=intern, status=inactive, level=1, age=16, salary=0, score=6, department=hr}
4. [1] TestPerson{id=P004, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active, name=Diana, age=45}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, department=sales, level=3, active=false, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, salary=-5000, active=true, score=0, tags=test, name=Frank, age=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, active=true, tags=executive, department=management, level=9, age=65, score=10, status=active, name=Grace, salary=95000}
8. [1] TestPerson{id=P008, tags=junior, department=support, level=1, age=18, status=inactive, name=Henry, salary=25000, active=false, score=5.5}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, score=8.7, tags=senior, department=engineering, level=6, salary=68000, status=active}
10. [1] TestPerson{id=P010, name=X, age=22, score=6.5, tags=temp, level=1, salary=28000, active=true, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, score=8.5, status=active, department=sales, level=2, age=25, tags=junior, name=Alice}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, salary=75000, level=5, age=35}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, name=Diana, age=45, score=7.8, department=marketing, level=7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, department=sales, level=3, active=false, tags=employee, status=inactive}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, score=10, status=active, name=Grace, salary=95000, active=true, tags=executive, department=management, level=9, age=65}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, active=true, score=8.7, tags=senior, department=engineering, level=6, salary=68000, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, discount=50, product_id=PROD001, total=1999.98, date=2024-01-15, region=north, customer_id=P001, amount=2, status=pending, priority=normal}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, priority=low, amount=1, discount=0, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, region=north, amount=3, total=225, date=2024-02-01, priority=high, discount=15}
4. [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0, customer_id=P004, amount=1}
5. [1] TestOrder{id=O005, customer_id=P002, priority=high, discount=100, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south}
6. [1] TestOrder{id=O006, amount=2, date=2024-02-15, priority=low, region=west, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, discount=0}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50, amount=4, priority=urgent, region=north}
8. [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, region=south, product_id=PROD002, amount=10, total=255, priority=normal, discount=0}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, priority=low, discount=10, amount=1}
10. [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, discount=0, customer_id=P006, product_id=PROD001, amount=1, priority=urgent, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, priority=low, amount=1, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, product_id=PROD003, status=shipped, region=north, amount=3}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, priority=high, discount=100}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, priority=normal, discount=0, customer_id=P010, date=2024-03-05, status=pending, region=south}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O001, date=2024-01-15, region=north, customer_id=P001, amount=2, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, discount=0, amount=2, date=2024-02-15, priority=low, region=west}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, priority=urgent, region=north, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, region=north, customer_id=P001, product_id=PROD007, priority=low, discount=10, amount=1, total=89.99, date=2024-03-10, status=completed}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, category=electronics, rating=4.5, brand=TechCorp, stock=50, price=999.99, available=true, keywords=computer, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, keywords=peripheral, brand=TechCorp, supplier=TechSupply, name=Mouse, category=accessories, price=25.5, available=true, rating=4, stock=200}
3. [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, price=75, available=false, keywords=typing, stock=0, rating=3.5, brand=KeyTech, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, supplier=ScreenSupply, name=Monitor, category=electronics, price=299.99, available=true, keywords=display, brand=ScreenPro, stock=30, rating=4.8}
5. [1] TestProduct{id=PROD005, brand=OldTech, stock=0, supplier=OldSupply, name=OldKeyboard, rating=2, category=accessories, price=8.5, available=false, keywords=obsolete}
6. [1] TestProduct{id=PROD006, brand=AudioMax, supplier=AudioSupply, category=audio, price=150, available=true, rating=4.6, keywords=sound, stock=75, name=Headphones}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, price=89.99, stock=25, supplier=CamSupply, available=true, rating=3.8, keywords=video, brand=CamTech}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, category=electronics, rating=4.5, brand=TechCorp, stock=50, price=999.99, available=true, keywords=computer, supplier=TechSupply}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, keywords=peripheral, brand=TechCorp, supplier=TechSupply, name=Mouse, category=accessories, price=25.5, available=true, rating=4, stock=200}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, rating=3.5, brand=KeyTech, supplier=KeySupply, name=Keyboard, category=accessories, price=75, available=false, keywords=typing, stock=0}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, category=electronics, price=299.99, available=true, keywords=display, brand=ScreenPro, stock=30, rating=4.8, supplier=ScreenSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, category=audio, price=150, available=true, rating=4.6, keywords=sound, stock=75, name=Headphones, brand=AudioMax, supplier=AudioSupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, keywords=video, brand=CamTech, name=Webcam, category=electronics, price=89.99, stock=25, supplier=CamSupply, available=true, rating=3.8}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, department=sales, level=2, age=25, tags=junior}
2. [1] TestPerson{id=P002, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, salary=75000, level=5}
3. [1] TestPerson{id=P003, name=Charlie, active=false, tags=intern, status=inactive, level=1, age=16, salary=0, score=6, department=hr}
4. [1] TestPerson{id=P004, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active, name=Diana}
5. [1] TestPerson{id=P005, score=8, department=sales, level=3, active=false, tags=employee, status=inactive, name=Eve, age=30, salary=55000}
6. [1] TestPerson{id=P006, department=qa, level=1, salary=-5000, active=true, score=0, tags=test, name=Frank, age=0, status=active}
7. [1] TestPerson{id=P007, age=65, score=10, status=active, name=Grace, salary=95000, active=true, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, age=18, status=inactive, name=Henry, salary=25000, active=false, score=5.5, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, score=8.7, tags=senior, department=engineering, level=6, salary=68000, status=active}
10. [1] TestPerson{id=P010, name=X, age=22, score=6.5, tags=temp, level=1, salary=28000, active=true, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, level=1, age=16, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, age=0, status=active, department=qa, level=1, salary=-5000, active=true, score=0, tags=test, name=Frank}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, date=2024-01-15, region=north, customer_id=P001, amount=2, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98}
2. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, status=confirmed, priority=low, amount=1, discount=0, region=south, customer_id=P002, product_id=PROD002}
3. [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, product_id=PROD003, status=shipped, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, priority=high, discount=100}
6. [1] TestOrder{id=O006, status=cancelled, discount=0, amount=2, date=2024-02-15, priority=low, region=west, customer_id=P005, product_id=PROD005, total=999.98}
7. [1] TestOrder{id=O007, amount=4, priority=urgent, region=north, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50}
8. [1] TestOrder{id=O008, priority=normal, discount=0, customer_id=P010, date=2024-03-05, status=pending, region=south, product_id=PROD002, amount=10, total=255}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, priority=low, discount=10, amount=1}
10. [1] TestOrder{id=O010, amount=1, priority=urgent, region=east, total=75000, date=2024-03-15, status=refunded, discount=0, customer_id=P006, product_id=PROD001}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, product_id=PROD003, status=shipped, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, discount=0, amount=2, date=2024-02-15, priority=low, region=west}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, region=south, product_id=PROD002, amount=10, total=255, priority=normal, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, priority=low, discount=10}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, date=2024-01-15, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O002, region=south, customer_id=P002, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, priority=low, amount=1, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0, customer_id=P004, amount=1}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, priority=high, discount=100}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50, amount=4, priority=urgent, region=north}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, status=refunded, discount=0, customer_id=P006, product_id=PROD001, amount=1, priority=urgent, region=east, total=75000, date=2024-03-15}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, name=Alice, salary=45000, active=true, score=8.5, status=active, department=sales, level=2, age=25}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, status=active, department=engineering, name=Bob, salary=75000, level=5, age=35, active=true}
3. [1] TestPerson{id=P003, name=Charlie, active=false, tags=intern, status=inactive, level=1, age=16, salary=0, score=6, department=hr}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, department=sales, level=3, active=false, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, age=0, status=active, department=qa, level=1, salary=-5000, active=true, score=0, tags=test, name=Frank}
7. [1] TestPerson{id=P007, age=65, score=10, status=active, name=Grace, salary=95000, active=true, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, department=support, level=1, age=18, status=inactive}
9. [1] TestPerson{id=P009, level=6, salary=68000, status=active, name=Ivy, age=40, active=true, score=8.7, tags=senior, department=engineering}
10. [1] TestPerson{id=P010, status=active, department=intern, name=X, age=22, score=6.5, tags=temp, level=1, salary=28000, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, department=sales, level=3, active=false, tags=employee, status=inactive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, level=1, age=18, status=inactive, name=Henry, salary=25000, active=false, score=5.5, tags=junior, department=support}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, level=1, age=16, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, amount=2, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, date=2024-01-15, region=north, customer_id=P001}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, priority=low, amount=1, discount=0, region=south}
3. [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, product_id=PROD003, status=shipped, region=north}
4. [1] TestOrder{id=O004, priority=normal, region=east, product_id=PROD004, discount=0, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered}
5. [1] TestOrder{id=O005, discount=100, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, priority=high}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, discount=0, amount=2, date=2024-02-15, priority=low, region=west}
7. [1] TestOrder{id=O007, region=north, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50, amount=4, priority=urgent}
8. [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, region=south, product_id=PROD002, amount=10, total=255, priority=normal, discount=0}
9. [1] TestOrder{id=O009, status=completed, region=north, customer_id=P001, product_id=PROD007, priority=low, discount=10, amount=1, total=89.99, date=2024-03-10}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, priority=urgent, region=east, total=75000, date=2024-03-15, status=refunded, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O009, status=completed, region=north, customer_id=P001, product_id=PROD007, priority=low, discount=10, amount=1, total=89.99, date=2024-03-10}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, amount=1, priority=urgent, region=east, total=75000, date=2024-03-15, status=refunded, discount=0, customer_id=P006}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, priority=low, amount=1, discount=0, region=south}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O003, priority=high, discount=15, customer_id=P001, product_id=PROD003, status=shipped, region=north, amount=3, total=225, date=2024-02-01}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O004, priority=normal, region=east, product_id=PROD004, discount=0, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, priority=high, discount=100}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, discount=0, amount=2, date=2024-02-15, priority=low, region=west}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, region=north, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50, amount=4, priority=urgent}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, tags=junior, name=Alice, salary=45000, active=true, score=8.5, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, salary=75000, level=5}
3. [1] TestPerson{id=P003, age=16, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, department=sales, level=3, active=false, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, salary=-5000, active=true, score=0, tags=test, name=Frank, age=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, tags=executive, department=management, level=9, age=65, score=10, status=active}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, department=support, level=1, age=18, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, score=8.7, tags=senior, department=engineering, level=6, salary=68000, status=active}
10. [1] TestPerson{id=P010, name=X, age=22, score=6.5, tags=temp, level=1, salary=28000, active=true, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, status=active, department=sales, level=2, age=25, tags=junior, name=Alice, salary=45000, active=true, score=8.5}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, level=5, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, score=8, department=sales, level=3, active=false, tags=employee, status=inactive, name=Eve, age=30, salary=55000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, active=true, score=0, tags=test, name=Frank, age=0, status=active, department=qa, level=1}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, tags=junior, department=support, level=1, age=18, status=inactive, name=Henry, salary=25000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, level=7, salary=85000, active=true, tags=manager, status=active, name=Diana, age=45, score=7.8, department=marketing}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, tags=executive, department=management, level=9, age=65, score=10, status=active}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, tags=senior, department=engineering, level=6, salary=68000, status=active, name=Ivy, age=40, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, price=999.99, available=true, keywords=computer, supplier=TechSupply, name=Laptop, category=electronics, rating=4.5, brand=TechCorp, stock=50}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, price=25.5, available=true, rating=4, stock=200, keywords=peripheral, brand=TechCorp, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, keywords=typing, stock=0, rating=3.5, brand=KeyTech, supplier=KeySupply, name=Keyboard, category=accessories, price=75, available=false}
4. [1] TestProduct{id=PROD004, rating=4.8, supplier=ScreenSupply, name=Monitor, category=electronics, price=299.99, available=true, keywords=display, brand=ScreenPro, stock=30}
5. [1] TestProduct{id=PROD005, category=accessories, price=8.5, available=false, keywords=obsolete, brand=OldTech, stock=0, supplier=OldSupply, name=OldKeyboard, rating=2}
6. [1] TestProduct{id=PROD006, name=Headphones, brand=AudioMax, supplier=AudioSupply, category=audio, price=150, available=true, rating=4.6, keywords=sound, stock=75}
7. [1] TestProduct{id=PROD007, brand=CamTech, name=Webcam, category=electronics, price=89.99, stock=25, supplier=CamSupply, available=true, rating=3.8, keywords=video}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD004, rating=4.8, supplier=ScreenSupply, name=Monitor, category=electronics, price=299.99, available=true, keywords=display, brand=ScreenPro, stock=30}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD006, supplier=AudioSupply, category=audio, price=150, available=true, rating=4.6, keywords=sound, stock=75, name=Headphones, brand=AudioMax}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD007, rating=3.8, keywords=video, brand=CamTech, name=Webcam, category=electronics, price=89.99, stock=25, supplier=CamSupply, available=true}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD001, rating=4.5, brand=TechCorp, stock=50, price=999.99, available=true, keywords=computer, supplier=TechSupply, name=Laptop, category=electronics}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD002, category=accessories, price=25.5, available=true, rating=4, stock=200, keywords=peripheral, brand=TechCorp, supplier=TechSupply, name=Mouse}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, price=75, available=false, keywords=typing, stock=0, rating=3.5, brand=KeyTech, supplier=KeySupply}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, department=sales, level=2, age=25, tags=junior}
2. [1] TestPerson{id=P002, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, salary=75000, level=5}
3. [1] TestPerson{id=P003, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1, age=16, salary=0, score=6}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, department=sales, level=3, active=false, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, age=0, status=active, department=qa, level=1, salary=-5000, active=true, score=0, tags=test}
7. [1] TestPerson{id=P007, department=management, level=9, age=65, score=10, status=active, name=Grace, salary=95000, active=true, tags=executive}
8. [1] TestPerson{id=P008, level=1, age=18, status=inactive, name=Henry, salary=25000, active=false, score=5.5, tags=junior, department=support}
9. [1] TestPerson{id=P009, department=engineering, level=6, salary=68000, status=active, name=Ivy, age=40, active=true, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, name=X, age=22, score=6.5, tags=temp, level=1, salary=28000, active=true, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, tags=senior, department=engineering, level=6, salary=68000, status=active, name=Ivy, age=40, active=true}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, department=sales, level=3, active=false, tags=employee, status=inactive}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, active=true, score=0, tags=test, name=Frank, age=0, status=active, department=qa, level=1, salary=-5000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, department=management, level=9, age=65, score=10, status=active, name=Grace, salary=95000, active=true}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, department=support, level=1, age=18, status=inactive, name=Henry, salary=25000, active=false, score=5.5, tags=junior}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, status=active, department=sales, level=2, age=25, tags=junior, name=Alice, salary=45000}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, level=5, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, region=north, customer_id=P001, amount=2, status=pending, priority=normal, discount=50, product_id=PROD001}
2. [1] TestOrder{id=O002, region=south, customer_id=P002, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, priority=low, amount=1, discount=0}
3. [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, product_id=PROD003, status=shipped, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, priority=high, discount=100}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, discount=0, amount=2, date=2024-02-15, priority=low, region=west}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50, amount=4, priority=urgent, region=north}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, priority=normal, discount=0, customer_id=P010, date=2024-03-05, status=pending, region=south}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, priority=low, discount=10, amount=1}
10. [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, discount=0, customer_id=P006, product_id=PROD001, amount=1, priority=urgent, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, priority=high, discount=100}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, priority=urgent, region=north, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, region=south, product_id=PROD002, amount=10, total=255, priority=normal, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, priority=low, discount=10, amount=1, total=89.99, date=2024-03-10, status=completed, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O001, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, date=2024-01-15, region=north, customer_id=P001, amount=2}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, priority=low, amount=1}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, product_id=PROD003, status=shipped, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, status=active, department=sales, level=2, age=25, tags=junior, name=Alice, salary=45000, active=true}
2. [1] TestPerson{id=P002, level=5, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, salary=75000}
3. [1] TestPerson{id=P003, age=16, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, active=false, tags=employee, status=inactive, name=Eve, age=30, salary=55000, score=8, department=sales, level=3}
6. [1] TestPerson{id=P006, status=active, department=qa, level=1, salary=-5000, active=true, score=0, tags=test, name=Frank, age=0}
7. [1] TestPerson{id=P007, department=management, level=9, age=65, score=10, status=active, name=Grace, salary=95000, active=true, tags=executive}
8. [1] TestPerson{id=P008, age=18, status=inactive, name=Henry, salary=25000, active=false, score=5.5, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, department=engineering, level=6, salary=68000, status=active, name=Ivy, age=40, active=true, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, salary=28000, active=true, status=active, department=intern, name=X, age=22, score=6.5, tags=temp, level=1}
11. [1] TestOrder{id=O001, customer_id=P001, amount=2, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, date=2024-01-15, region=north}
12. [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, priority=low, amount=1}
13. [1] TestOrder{id=O003, region=north, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, product_id=PROD003, status=shipped}
14. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}
15. [1] TestOrder{id=O005, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, priority=high, discount=100, product_id=PROD001, amount=1}
16. [1] TestOrder{id=O006, priority=low, region=west, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, discount=0, amount=2, date=2024-02-15}
17. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50, amount=4, priority=urgent, region=north}
18. [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, region=south, product_id=PROD002, amount=10, total=255, priority=normal, discount=0}
19. [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, priority=low, discount=10}
20. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, priority=urgent, region=east, total=75000, date=2024-03-15, status=refunded, discount=0}

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

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, department=sales, level=2, age=25, tags=junior}
2. [1] TestPerson{id=P002, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, salary=75000, level=5}
3. [1] TestPerson{id=P003, tags=intern, status=inactive, level=1, age=16, salary=0, score=6, department=hr, name=Charlie, active=false}
4. [1] TestPerson{id=P004, active=true, tags=manager, status=active, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000}
5. [1] TestPerson{id=P005, active=false, tags=employee, status=inactive, name=Eve, age=30, salary=55000, score=8, department=sales, level=3}
6. [1] TestPerson{id=P006, level=1, salary=-5000, active=true, score=0, tags=test, name=Frank, age=0, status=active, department=qa}
7. [1] TestPerson{id=P007, status=active, name=Grace, salary=95000, active=true, tags=executive, department=management, level=9, age=65, score=10}
8. [1] TestPerson{id=P008, tags=junior, department=support, level=1, age=18, status=inactive, name=Henry, salary=25000, active=false, score=5.5}
9. [1] TestPerson{id=P009, tags=senior, department=engineering, level=6, salary=68000, status=active, name=Ivy, age=40, active=true, score=8.7}
10. [1] TestPerson{id=P010, name=X, age=22, score=6.5, tags=temp, level=1, salary=28000, active=true, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, name=X, age=22, score=6.5, tags=temp, level=1, salary=28000, active=true, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, status=active, department=sales, level=2, age=25, tags=junior, name=Alice, salary=45000, active=true}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, level=5, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, name=Diana, age=45, score=7.8, department=marketing, level=7}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, status=active, department=qa, level=1, salary=-5000, active=true, score=0, tags=test}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, tags=executive, department=management, level=9, age=65, score=10, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, tags=senior, department=engineering, level=6, salary=68000, status=active, name=Ivy, age=40, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, department=sales, level=2, age=25, tags=junior}
2. [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, salary=75000, level=5, age=35}
3. [1] TestPerson{id=P003, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1, age=16, salary=0, score=6}
4. [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, name=Diana, age=45, score=7.8, department=marketing, level=7}
5. [1] TestPerson{id=P005, level=3, active=false, tags=employee, status=inactive, name=Eve, age=30, salary=55000, score=8, department=sales}
6. [1] TestPerson{id=P006, score=0, tags=test, name=Frank, age=0, status=active, department=qa, level=1, salary=-5000, active=true}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, tags=executive, department=management, level=9, age=65, score=10, status=active}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, department=support, level=1, age=18, status=inactive}
9. [1] TestPerson{id=P009, salary=68000, status=active, name=Ivy, age=40, active=true, score=8.7, tags=senior, department=engineering, level=6}
10. [1] TestPerson{id=P010, name=X, age=22, score=6.5, tags=temp, level=1, salary=28000, active=true, status=active, department=intern}
11. [1] TestOrder{id=O001, customer_id=P001, amount=2, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, date=2024-01-15, region=north}
12. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, priority=low, amount=1, discount=0, region=south}
13. [1] TestOrder{id=O003, product_id=PROD003, status=shipped, region=north, amount=3, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001}
14. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}
15. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, priority=high, discount=100}
16. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, discount=0, amount=2, date=2024-02-15, priority=low, region=west}
17. [1] TestOrder{id=O007, date=2024-03-01, status=shipped, discount=50, amount=4, priority=urgent, region=north, customer_id=P007, product_id=PROD006, total=600}
18. [1] TestOrder{id=O008, status=pending, region=south, product_id=PROD002, amount=10, total=255, priority=normal, discount=0, customer_id=P010, date=2024-03-05}
19. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, priority=low, discount=10, amount=1, total=89.99, date=2024-03-10, status=completed, region=north}
20. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, priority=urgent, region=east, total=75000, date=2024-03-15, status=refunded, discount=0}

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

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, department=sales, level=2, age=25, tags=junior}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, level=5, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering}
3. [1] TestPerson{id=P003, active=false, tags=intern, status=inactive, level=1, age=16, salary=0, score=6, department=hr, name=Charlie}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, active=false, tags=employee, status=inactive, name=Eve, age=30, salary=55000, score=8, department=sales, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, status=active, department=qa, level=1, salary=-5000, active=true, score=0, tags=test}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, tags=executive, department=management, level=9, age=65, score=10, status=active}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, department=support, level=1, age=18, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, score=8.7, tags=senior, department=engineering, level=6, salary=68000, status=active}
10. [1] TestPerson{id=P010, name=X, age=22, score=6.5, tags=temp, level=1, salary=28000, active=true, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, name=Bob, salary=75000, level=5, age=35, active=true, score=9.2, tags=senior}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, active=false, tags=intern, status=inactive, level=1, age=16, salary=0, score=6, department=hr, name=Charlie}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, score=10, status=active, name=Grace, salary=95000, active=true, tags=executive, department=management, level=9, age=65}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, department=support, level=1, age=18, status=inactive, name=Henry, salary=25000, active=false, score=5.5, tags=junior}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, active=true, score=8.7, tags=senior, department=engineering, level=6, salary=68000, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, department=sales, level=2, age=25, tags=junior}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, active=false, tags=employee, status=inactive, name=Eve, age=30, salary=55000, score=8, department=sales, level=3}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, age=22, score=6.5, tags=temp, level=1, salary=28000, active=true, status=active, department=intern, name=X}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, date=2024-01-15, region=north, customer_id=P001, amount=2, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98}
2. [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, priority=low, amount=1}
3. [1] TestOrder{id=O003, total=225, date=2024-02-01, priority=high, discount=15, customer_id=P001, product_id=PROD003, status=shipped, region=north, amount=3}
4. [1] TestOrder{id=O004, discount=0, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004}
5. [1] TestOrder{id=O005, region=south, customer_id=P002, priority=high, discount=100, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, status=confirmed}
6. [1] TestOrder{id=O006, region=west, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, discount=0, amount=2, date=2024-02-15, priority=low}
7. [1] TestOrder{id=O007, product_id=PROD006, total=600, date=2024-03-01, status=shipped, discount=50, amount=4, priority=urgent, region=north, customer_id=P007}
8. [1] TestOrder{id=O008, total=255, priority=normal, discount=0, customer_id=P010, date=2024-03-05, status=pending, region=south, product_id=PROD002, amount=10}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, priority=low, discount=10, amount=1, total=89.99, date=2024-03-10, status=completed, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, priority=urgent, region=east, total=75000, date=2024-03-15, status=refunded, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O008, status=pending, region=south, product_id=PROD002, amount=10, total=255, priority=normal, discount=0, customer_id=P010, date=2024-03-05}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O009, discount=10, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, priority=low}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O010, date=2024-03-15, status=refunded, discount=0, customer_id=P006, product_id=PROD001, amount=1, priority=urgent, region=east, total=75000}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, date=2024-01-20, status=confirmed, priority=low}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, status=confirmed, region=south, customer_id=P002, priority=high, discount=100}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O006, date=2024-02-15, priority=low, region=west, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, discount=0, amount=2}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, discount=50, amount=4, priority=urgent, region=north, customer_id=P007, product_id=PROD006, total=600, date=2024-03-01, status=shipped}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, date=2024-01-15, region=north}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, region=north, amount=3, total=225, date=2024-02-01, priority=high, discount=15}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, level=2, age=25, tags=junior, name=Alice, salary=45000, active=true, score=8.5, status=active}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, status=active, department=engineering, name=Bob, salary=75000, level=5, age=35, active=true}
3. [1] TestPerson{id=P003, age=16, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, tags=employee, status=inactive, name=Eve, age=30, salary=55000, score=8, department=sales, level=3, active=false}
6. [1] TestPerson{id=P006, age=0, status=active, department=qa, level=1, salary=-5000, active=true, score=0, tags=test, name=Frank}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, tags=executive, department=management, level=9, age=65, score=10, status=active}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, department=support, level=1, age=18, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, score=8.7, tags=senior, department=engineering, level=6, salary=68000, status=active}
10. [1] TestPerson{id=P010, score=6.5, tags=temp, level=1, salary=28000, active=true, status=active, department=intern, name=X, age=22}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, age=18, status=inactive, name=Henry, salary=25000, active=false, score=5.5, tags=junior, department=support, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, status=active, department=sales, level=2, age=25, tags=junior, name=Alice, salary=45000}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P002, age=35, active=true, score=9.2, tags=senior, status=active, department=engineering, name=Bob, salary=75000, level=5}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, score=6, department=hr, name=Charlie, active=false, tags=intern, status=inactive, level=1, age=16}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, department=marketing, level=7, salary=85000, active=true, tags=manager, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, department=sales, level=3, active=false, tags=employee, status=inactive}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P009, active=true, score=8.7, tags=senior, department=engineering, level=6, salary=68000, status=active, name=Ivy, age=40}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P010, age=22, score=6.5, tags=temp, level=1, salary=28000, active=true, status=active, department=intern, name=X}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, level=1, salary=-5000, active=true, score=0, tags=test, name=Frank, age=0, status=active}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P007, score=10, status=active, name=Grace, salary=95000, active=true, tags=executive, department=management, level=9, age=65}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 17 (89.5%)
- **Tokens générés**: 130
- **Faits traités**: 27
