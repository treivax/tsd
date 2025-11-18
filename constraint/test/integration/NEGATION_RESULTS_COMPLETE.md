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

1. [1] TestPerson{id=P001, name=Alice, age=25, active=true, status=active, level=2, salary=45000, score=8.5, tags=junior, department=sales}
2. [1] TestPerson{id=P002, status=active, level=5, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2}
3. [1] TestPerson{id=P003, status=inactive, department=hr, salary=0, active=false, level=1, name=Charlie, age=16, score=6, tags=intern}
4. [1] TestPerson{id=P004, salary=85000, tags=manager, department=marketing, level=7, age=45, active=true, score=7.8, status=active, name=Diana}
5. [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, name=Eve, active=false, score=8, department=sales, level=3, age=30}
6. [1] TestPerson{id=P006, department=qa, level=1, salary=-5000, score=0, status=active, name=Frank, age=0, active=true, tags=test}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, tags=executive, status=active, department=management, level=9, salary=95000, active=true}
8. [1] TestPerson{id=P008, age=18, score=5.5, tags=junior, department=support, name=Henry, salary=25000, active=false, status=inactive, level=1}
9. [1] TestPerson{id=P009, tags=senior, salary=68000, active=true, score=8.7, status=active, department=engineering, level=6, name=Ivy, age=40}
10. [1] TestPerson{id=P010, department=intern, salary=28000, level=1, name=X, age=22, active=true, score=6.5, tags=temp, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, level=2, salary=45000, score=8.5, tags=junior, department=sales, name=Alice, age=25, active=true, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, active=false, level=1, name=Charlie, age=16, score=6, tags=intern, status=inactive, department=hr, salary=0}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, tags=manager, department=marketing, level=7, age=45, active=true, score=7.8, status=active}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, score=10, tags=executive, status=active, department=management, level=9, salary=95000, active=true}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, salary=68000, active=true, score=8.7, status=active, department=engineering, level=6, name=Ivy, age=40}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, name=Eve, active=false, score=8, department=sales, level=3, age=30, salary=55000, tags=employee}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, department=support, name=Henry, salary=25000, active=false, status=inactive, level=1, age=18, score=5.5, tags=junior}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, amount=2, date=2024-01-15, priority=normal, discount=50, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, region=south, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0}
3. [1] TestOrder{id=O003, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, region=north}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, region=east, customer_id=P004, date=2024-02-05, priority=normal}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, discount=100, region=south, total=999.99, status=confirmed}
6. [1] TestOrder{id=O006, total=999.98, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, priority=low, region=west, customer_id=P005}
7. [1] TestOrder{id=O007, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, amount=4, total=600, date=2024-03-01, product_id=PROD006}
8. [1] TestOrder{id=O008, amount=10, total=255, priority=normal, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, status=completed, amount=1, total=89.99, date=2024-03-10, priority=low, discount=10, region=north}
10. [1] TestOrder{id=O010, product_id=PROD001, total=75000, date=2024-03-15, priority=urgent, discount=0, amount=1, status=refunded, region=east, customer_id=P006}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, region=north, amount=2, date=2024-01-15, priority=normal, discount=50}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, discount=0, customer_id=P002, product_id=PROD002, total=25.5, region=south, amount=1, date=2024-01-20, status=confirmed, priority=low}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, date=2024-03-01, product_id=PROD006, status=shipped, priority=urgent, discount=50, region=north}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O008, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal, customer_id=P010}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, region=north, customer_id=P001, product_id=PROD007, status=completed, amount=1, total=89.99, date=2024-03-10, priority=low, discount=10}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O010, region=east, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, priority=urgent, discount=0, amount=1, status=refunded}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O003, region=north, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01, status=shipped}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, region=east, customer_id=P004, date=2024-02-05, priority=normal}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O005, date=2024-02-10, priority=high, discount=100, region=south, total=999.99, status=confirmed, customer_id=P002, product_id=PROD001, amount=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, active=true, status=active, level=2, salary=45000, score=8.5, tags=junior, department=sales}
2. [1] TestPerson{id=P002, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5}
3. [1] TestPerson{id=P003, department=hr, salary=0, active=false, level=1, name=Charlie, age=16, score=6, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, tags=manager, department=marketing, level=7, age=45, active=true, score=7.8, status=active}
5. [1] TestPerson{id=P005, level=3, age=30, salary=55000, tags=employee, status=inactive, name=Eve, active=false, score=8, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, department=qa, level=1, salary=-5000, score=0, status=active}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, tags=executive, status=active, department=management, level=9, salary=95000, active=true}
8. [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, salary=25000, active=false, status=inactive, level=1, age=18, score=5.5}
9. [1] TestPerson{id=P009, score=8.7, status=active, department=engineering, level=6, name=Ivy, age=40, tags=senior, salary=68000, active=true}
10. [1] TestPerson{id=P010, status=active, department=intern, salary=28000, level=1, name=X, age=22, active=true, score=6.5, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, department=sales, name=Alice, age=25, active=true, status=active, level=2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, status=active, level=5, name=Bob, active=true, tags=senior, department=engineering}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, score=7.8, status=active, name=Diana, salary=85000, tags=manager, department=marketing, level=7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, department=sales, level=3, age=30, salary=55000, tags=employee, status=inactive}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, score=10, tags=executive, status=active, department=management, level=9, salary=95000, active=true}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, status=active, department=engineering, level=6, name=Ivy, age=40, tags=senior}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, total=1999.98, status=pending, region=north, amount=2, date=2024-01-15, priority=normal, discount=50, customer_id=P001, product_id=PROD001}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, region=south, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0}
3. [1] TestOrder{id=O003, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, region=north, customer_id=P001, total=225}
4. [1] TestOrder{id=O004, date=2024-02-05, priority=normal, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, region=east, customer_id=P004}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, discount=100, region=south, total=999.99, status=confirmed}
6. [1] TestOrder{id=O006, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, priority=low, region=west, customer_id=P005, total=999.98}
7. [1] TestOrder{id=O007, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, amount=4, total=600, date=2024-03-01, product_id=PROD006}
8. [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal}
9. [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, status=completed}
10. [1] TestOrder{id=O010, status=refunded, region=east, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, priority=urgent, discount=0, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, status=pending, region=north, amount=2, date=2024-01-15, priority=normal, discount=50, customer_id=P001, product_id=PROD001, total=1999.98}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, product_id=PROD002, total=25.5, region=south}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, region=north, customer_id=P001, total=225, priority=high, discount=15}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, discount=0, region=east, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, date=2024-02-10, priority=high, discount=100, region=south, total=999.99, status=confirmed, customer_id=P002, product_id=PROD001}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, priority=low, region=west}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, date=2024-03-01, product_id=PROD006, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, amount=4, total=600}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O008, status=pending, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal, customer_id=P010, date=2024-03-05}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, status=completed}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, rating=4.5, keywords=computer, category=electronics, price=999.99, available=true, brand=TechCorp, stock=50, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, stock=200, supplier=TechSupply, price=25.5, keywords=peripheral, name=Mouse, category=accessories, available=true, rating=4, brand=TechCorp}
3. [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, price=75, available=false, brand=KeyTech, stock=0, supplier=KeySupply, rating=3.5, keywords=typing}
4. [1] TestProduct{id=PROD004, keywords=display, brand=ScreenPro, price=299.99, available=true, rating=4.8, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics}
5. [1] TestProduct{id=PROD005, brand=OldTech, stock=0, available=false, rating=2, supplier=OldSupply, name=OldKeyboard, category=accessories, price=8.5, keywords=obsolete}
6. [1] TestProduct{id=PROD006, rating=4.6, keywords=sound, supplier=AudioSupply, name=Headphones, category=audio, price=150, available=true, brand=AudioMax, stock=75}
7. [1] TestProduct{id=PROD007, category=electronics, price=89.99, available=true, rating=3.8, keywords=video, stock=25, supplier=CamSupply, name=Webcam, brand=CamTech}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, rating=4.5, keywords=computer, category=electronics, price=999.99, available=true, brand=TechCorp, stock=50, supplier=TechSupply}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, price=25.5, keywords=peripheral, name=Mouse, category=accessories, available=true, rating=4, brand=TechCorp, stock=200, supplier=TechSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, brand=KeyTech, stock=0, supplier=KeySupply, rating=3.5, keywords=typing, name=Keyboard, category=accessories, price=75, available=false}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, category=electronics, keywords=display, brand=ScreenPro, price=299.99, available=true, rating=4.8, stock=30, supplier=ScreenSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, supplier=AudioSupply, name=Headphones, category=audio, price=150, available=true, brand=AudioMax, stock=75, rating=4.6, keywords=sound}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, category=electronics, price=89.99, available=true, rating=3.8, keywords=video, stock=25, supplier=CamSupply, name=Webcam, brand=CamTech}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, active=true, status=active, level=2, salary=45000, score=8.5, tags=junior, department=sales, name=Alice}
2. [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, status=active, level=5, name=Bob, active=true, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, salary=0, active=false, level=1, name=Charlie, age=16, score=6, tags=intern, status=inactive, department=hr}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, tags=manager, department=marketing, level=7, age=45, active=true, score=7.8, status=active}
5. [1] TestPerson{id=P005, name=Eve, active=false, score=8, department=sales, level=3, age=30, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, department=qa, level=1, salary=-5000, score=0, status=active}
7. [1] TestPerson{id=P007, tags=executive, status=active, department=management, level=9, salary=95000, active=true, name=Grace, age=65, score=10}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, status=inactive, level=1, age=18, score=5.5, tags=junior, department=support}
9. [1] TestPerson{id=P009, name=Ivy, age=40, tags=senior, salary=68000, active=true, score=8.7, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1, name=X, age=22}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, active=false, level=1, name=Charlie, age=16, score=6, tags=intern, status=inactive, department=hr, salary=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, status=active, name=Frank, age=0, active=true, tags=test, department=qa, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, amount=2, date=2024-01-15, priority=normal, discount=50, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, region=north}
2. [1] TestOrder{id=O002, priority=low, discount=0, customer_id=P002, product_id=PROD002, total=25.5, region=south, amount=1, date=2024-01-20, status=confirmed}
3. [1] TestOrder{id=O003, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, region=north}
4. [1] TestOrder{id=O004, status=delivered, discount=0, region=east, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, amount=1, total=299.99}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, discount=100, region=south, total=999.99, status=confirmed}
6. [1] TestOrder{id=O006, amount=2, priority=low, region=west, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005}
7. [1] TestOrder{id=O007, product_id=PROD006, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, amount=4, total=600, date=2024-03-01}
8. [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal}
9. [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, status=completed}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, priority=urgent, discount=0, amount=1, status=refunded, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, status=completed, amount=1, total=89.99, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, priority=urgent, discount=0, amount=1, status=refunded, region=east}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, region=north, customer_id=P001, total=225, priority=high, discount=15}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O004, discount=0, region=east, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, amount=1, total=299.99, status=delivered}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, amount=2, priority=low, region=west, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O001, amount=2, date=2024-01-15, priority=normal, discount=50, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, region=south, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, product_id=PROD002, total=25.5}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O005, priority=high, discount=100, region=south, total=999.99, status=confirmed, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, date=2024-03-01, product_id=PROD006, status=shipped, priority=urgent, discount=50, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, active=true, status=active, level=2, salary=45000, score=8.5, tags=junior, department=sales}
2. [1] TestPerson{id=P002, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5}
3. [1] TestPerson{id=P003, age=16, score=6, tags=intern, status=inactive, department=hr, salary=0, active=false, level=1, name=Charlie}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, tags=manager, department=marketing, level=7, age=45, active=true, score=7.8, status=active}
5. [1] TestPerson{id=P005, level=3, age=30, salary=55000, tags=employee, status=inactive, name=Eve, active=false, score=8, department=sales}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, name=Frank, age=0, active=true, tags=test, department=qa, level=1}
7. [1] TestPerson{id=P007, status=active, department=management, level=9, salary=95000, active=true, name=Grace, age=65, score=10, tags=executive}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, status=inactive, level=1, age=18, score=5.5, tags=junior, department=support}
9. [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, status=active, department=engineering, level=6, name=Ivy, age=40, tags=senior}
10. [1] TestPerson{id=P010, salary=28000, level=1, name=X, age=22, active=true, score=6.5, tags=temp, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, status=inactive, department=hr, salary=0, active=false, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, department=sales, level=3, age=30, salary=55000, tags=employee, status=inactive}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, active=false, status=inactive, level=1, age=18, score=5.5, tags=junior, department=support, name=Henry, salary=25000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, status=pending, region=north, amount=2, date=2024-01-15, priority=normal, discount=50, customer_id=P001, product_id=PROD001, total=1999.98}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, region=south, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0}
3. [1] TestOrder{id=O003, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, region=north}
4. [1] TestOrder{id=O004, discount=0, region=east, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, amount=1, total=299.99, status=delivered}
5. [1] TestOrder{id=O005, total=999.99, status=confirmed, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, discount=100, region=south}
6. [1] TestOrder{id=O006, total=999.98, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, priority=low, region=west, customer_id=P005}
7. [1] TestOrder{id=O007, priority=urgent, discount=50, region=north, customer_id=P007, amount=4, total=600, date=2024-03-01, product_id=PROD006, status=shipped}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, priority=normal, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, status=completed, amount=1, total=89.99, date=2024-03-10, priority=low, discount=10, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, priority=urgent, discount=0, amount=1, status=refunded, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, status=delivered, discount=0, region=east, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, amount=1, total=299.99}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, discount=100, region=south, total=999.99, status=confirmed}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, amount=2, priority=low, region=west, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, amount=4, total=600, date=2024-03-01}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, status=completed, amount=1, total=89.99}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, priority=urgent, discount=0, amount=1, status=refunded, region=east}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, product_id=PROD002, total=25.5, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, date=2024-02-01, status=shipped, region=north, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003, amount=3}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, tags=junior, department=sales, name=Alice, age=25, active=true, status=active, level=2, salary=45000}
2. [1] TestPerson{id=P002, status=active, level=5, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2}
3. [1] TestPerson{id=P003, age=16, score=6, tags=intern, status=inactive, department=hr, salary=0, active=false, level=1, name=Charlie}
4. [1] TestPerson{id=P004, department=marketing, level=7, age=45, active=true, score=7.8, status=active, name=Diana, salary=85000, tags=manager}
5. [1] TestPerson{id=P005, age=30, salary=55000, tags=employee, status=inactive, name=Eve, active=false, score=8, department=sales, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, department=qa, level=1, salary=-5000, score=0, status=active}
7. [1] TestPerson{id=P007, tags=executive, status=active, department=management, level=9, salary=95000, active=true, name=Grace, age=65, score=10}
8. [1] TestPerson{id=P008, active=false, status=inactive, level=1, age=18, score=5.5, tags=junior, department=support, name=Henry, salary=25000}
9. [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, status=active, department=engineering, level=6, name=Ivy, age=40, tags=senior}
10. [1] TestPerson{id=P010, name=X, age=22, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, score=9.2, status=active, level=5, name=Bob, active=true, tags=senior, department=engineering, age=35}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, status=inactive, department=hr, salary=0, active=false, level=1}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, name=Eve, active=false, score=8, department=sales, level=3, age=30, salary=55000, tags=employee}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, status=active, department=management, level=9, salary=95000, active=true, name=Grace, age=65, score=10, tags=executive}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, tags=manager, department=marketing, level=7, age=45, active=true, score=7.8, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, tags=test, department=qa, level=1, salary=-5000, score=0, status=active, name=Frank}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P008, age=18, score=5.5, tags=junior, department=support, name=Henry, salary=25000, active=false, status=inactive, level=1}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, status=active, department=engineering, level=6, name=Ivy, age=40, tags=senior}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P001, active=true, status=active, level=2, salary=45000, score=8.5, tags=junior, department=sales, name=Alice, age=25}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, rating=4.5, keywords=computer, category=electronics, price=999.99, available=true, brand=TechCorp, stock=50, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, brand=TechCorp, stock=200, supplier=TechSupply, price=25.5, keywords=peripheral, name=Mouse, category=accessories, available=true, rating=4}
3. [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, price=75, available=false, brand=KeyTech, stock=0, supplier=KeySupply, rating=3.5, keywords=typing}
4. [1] TestProduct{id=PROD004, category=electronics, keywords=display, brand=ScreenPro, price=299.99, available=true, rating=4.8, stock=30, supplier=ScreenSupply, name=Monitor}
5. [1] TestProduct{id=PROD005, brand=OldTech, stock=0, available=false, rating=2, supplier=OldSupply, name=OldKeyboard, category=accessories, price=8.5, keywords=obsolete}
6. [1] TestProduct{id=PROD006, stock=75, rating=4.6, keywords=sound, supplier=AudioSupply, name=Headphones, category=audio, price=150, available=true, brand=AudioMax}
7. [1] TestProduct{id=PROD007, category=electronics, price=89.99, available=true, rating=3.8, keywords=video, stock=25, supplier=CamSupply, name=Webcam, brand=CamTech}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, category=electronics, price=999.99, available=true, brand=TechCorp, stock=50, supplier=TechSupply, name=Laptop, rating=4.5, keywords=computer}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, category=accessories, available=true, rating=4, brand=TechCorp, stock=200, supplier=TechSupply, price=25.5, keywords=peripheral}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, rating=3.5, keywords=typing, name=Keyboard, category=accessories, price=75, available=false, brand=KeyTech, stock=0, supplier=KeySupply}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, available=true, rating=4.8, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics, keywords=display, brand=ScreenPro, price=299.99}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, name=Headphones, category=audio, price=150, available=true, brand=AudioMax, stock=75, rating=4.6, keywords=sound, supplier=AudioSupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, stock=25, supplier=CamSupply, name=Webcam, brand=CamTech, category=electronics, price=89.99, available=true, rating=3.8, keywords=video}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, salary=45000, score=8.5, tags=junior, department=sales, name=Alice, age=25, active=true, status=active}
2. [1] TestPerson{id=P002, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5, name=Bob}
3. [1] TestPerson{id=P003, age=16, score=6, tags=intern, status=inactive, department=hr, salary=0, active=false, level=1, name=Charlie}
4. [1] TestPerson{id=P004, department=marketing, level=7, age=45, active=true, score=7.8, status=active, name=Diana, salary=85000, tags=manager}
5. [1] TestPerson{id=P005, tags=employee, status=inactive, name=Eve, active=false, score=8, department=sales, level=3, age=30, salary=55000}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, name=Frank, age=0, active=true, tags=test, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, tags=executive, status=active, department=management, level=9, salary=95000, active=true}
8. [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, salary=25000, active=false, status=inactive, level=1, age=18, score=5.5}
9. [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, status=active, department=engineering, level=6, name=Ivy, age=40, tags=senior}
10. [1] TestPerson{id=P010, department=intern, salary=28000, level=1, name=X, age=22, active=true, score=6.5, tags=temp, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, status=active, name=Diana, salary=85000, tags=manager, department=marketing, level=7, age=45, active=true}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, tags=employee, status=inactive, name=Eve, active=false, score=8, department=sales, level=3}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, department=qa, level=1, salary=-5000, score=0, status=active, name=Frank, age=0, active=true}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, score=10, tags=executive, status=active, department=management, level=9, salary=95000, active=true}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, level=6, name=Ivy, age=40, tags=senior, salary=68000, active=true, score=8.7, status=active, department=engineering}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, department=sales, name=Alice, age=25, active=true, status=active, level=2, salary=45000}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5, name=Bob}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P003, score=6, tags=intern, status=inactive, department=hr, salary=0, active=false, level=1, name=Charlie, age=16}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, salary=25000, active=false, status=inactive, level=1, age=18, score=5.5}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, region=north, amount=2, date=2024-01-15, priority=normal, discount=50, customer_id=P001, product_id=PROD001, total=1999.98, status=pending}
2. [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, product_id=PROD002, total=25.5, region=south}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, region=north, customer_id=P001, total=225, priority=high, discount=15}
4. [1] TestOrder{id=O004, total=299.99, status=delivered, discount=0, region=east, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, amount=1}
5. [1] TestOrder{id=O005, status=confirmed, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, discount=100, region=south, total=999.99}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, priority=low, region=west, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, discount=0}
7. [1] TestOrder{id=O007, product_id=PROD006, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, amount=4, total=600, date=2024-03-01}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, priority=normal, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south}
9. [1] TestOrder{id=O009, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, status=completed, amount=1, total=89.99, date=2024-03-10}
10. [1] TestOrder{id=O010, total=75000, date=2024-03-15, priority=urgent, discount=0, amount=1, status=refunded, region=east, customer_id=P006, product_id=PROD001}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, priority=normal, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, region=east, customer_id=P004, date=2024-02-05}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, discount=100, region=south, total=999.99, status=confirmed}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, amount=4, total=600, date=2024-03-01}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O008, total=255, priority=normal, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002, amount=10}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, status=completed}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, discount=50, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, region=north, amount=2, date=2024-01-15, priority=normal}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, region=south, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, product_id=PROD002, total=25.5}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, department=sales, name=Alice, age=25, active=true, status=active, level=2}
2. [1] TestPerson{id=P002, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5}
3. [1] TestPerson{id=P003, salary=0, active=false, level=1, name=Charlie, age=16, score=6, tags=intern, status=inactive, department=hr}
4. [1] TestPerson{id=P004, age=45, active=true, score=7.8, status=active, name=Diana, salary=85000, tags=manager, department=marketing, level=7}
5. [1] TestPerson{id=P005, name=Eve, active=false, score=8, department=sales, level=3, age=30, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, department=qa, level=1, salary=-5000, score=0, status=active}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, tags=executive, status=active, department=management, level=9, salary=95000, active=true}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, status=inactive, level=1, age=18, score=5.5, tags=junior, department=support}
9. [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, status=active, department=engineering, level=6, name=Ivy, age=40, tags=senior}
10. [1] TestPerson{id=P010, name=X, age=22, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1}
11. [1] TestOrder{id=O001, amount=2, date=2024-01-15, priority=normal, discount=50, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, region=north}
12. [1] TestOrder{id=O002, total=25.5, region=south, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, product_id=PROD002}
13. [1] TestOrder{id=O003, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, region=north}
14. [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, region=east}
15. [1] TestOrder{id=O005, total=999.99, status=confirmed, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, discount=100, region=south}
16. [1] TestOrder{id=O006, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, priority=low, region=west, customer_id=P005, total=999.98}
17. [1] TestOrder{id=O007, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, amount=4, total=600, date=2024-03-01, product_id=PROD006}
18. [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal}
19. [1] TestOrder{id=O009, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, status=completed, amount=1, total=89.99}
20. [1] TestOrder{id=O010, amount=1, status=refunded, region=east, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, priority=urgent, discount=0}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, status=active, level=5, name=Bob, active=true, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O002, discount=0, customer_id=P002, product_id=PROD002, total=25.5, region=south, amount=1, date=2024-01-20, status=confirmed, priority=low}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, department=sales, name=Alice, age=25, active=true, status=active, level=2, salary=45000}
   - Fait 2: [1] TestOrder{id=O003, region=north, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01, status=shipped}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, score=7.8, status=active, name=Diana, salary=85000, tags=manager, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O004, discount=0, region=east, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, amount=1, total=299.99, status=delivered}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, department=management, level=9, salary=95000, active=true, name=Grace, age=65, score=10, tags=executive, status=active}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, date=2024-03-01, product_id=PROD006, status=shipped, priority=urgent, discount=50, region=north}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, active=true, status=active, level=2, salary=45000, score=8.5, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O001, amount=2, date=2024-01-15, priority=normal, discount=50, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, region=north}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O005, priority=high, discount=100, region=south, total=999.99, status=confirmed, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, department=sales, level=3, age=30, salary=55000, tags=employee, status=inactive, name=Eve}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, priority=low, region=west}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, active=true, status=active, level=2, salary=45000, score=8.5, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, status=completed, amount=1, total=89.99, date=2024-03-10, priority=low, discount=10, region=north}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, department=qa, level=1, salary=-5000, score=0, status=active, name=Frank, age=0, active=true}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, priority=urgent, discount=0, amount=1, status=refunded, region=east, customer_id=P006, product_id=PROD001, total=75000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/20 (50.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, level=2, salary=45000, score=8.5, tags=junior, department=sales, name=Alice, age=25, active=true}
2. [1] TestPerson{id=P002, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5}
3. [1] TestPerson{id=P003, salary=0, active=false, level=1, name=Charlie, age=16, score=6, tags=intern, status=inactive, department=hr}
4. [1] TestPerson{id=P004, status=active, name=Diana, salary=85000, tags=manager, department=marketing, level=7, age=45, active=true, score=7.8}
5. [1] TestPerson{id=P005, tags=employee, status=inactive, name=Eve, active=false, score=8, department=sales, level=3, age=30, salary=55000}
6. [1] TestPerson{id=P006, score=0, status=active, name=Frank, age=0, active=true, tags=test, department=qa, level=1, salary=-5000}
7. [1] TestPerson{id=P007, level=9, salary=95000, active=true, name=Grace, age=65, score=10, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, active=false, status=inactive, level=1, age=18, score=5.5, tags=junior, department=support, name=Henry, salary=25000}
9. [1] TestPerson{id=P009, department=engineering, level=6, name=Ivy, age=40, tags=senior, salary=68000, active=true, score=8.7, status=active}
10. [1] TestPerson{id=P010, department=intern, salary=28000, level=1, name=X, age=22, active=true, score=6.5, tags=temp, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, department=sales, name=Alice, age=25, active=true, status=active, level=2, salary=45000, score=8.5}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, age=45, active=true, score=7.8, status=active, name=Diana, salary=85000, tags=manager}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, score=0, status=active, name=Frank, age=0, active=true, tags=test, department=qa, level=1, salary=-5000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, status=active, department=management, level=9, salary=95000, active=true, name=Grace, age=65}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, status=active, department=engineering, level=6, name=Ivy, age=40, tags=senior}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, level=1, name=X, age=22, active=true, score=6.5, tags=temp, status=active, department=intern}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, active=true, status=active, level=2, salary=45000, score=8.5, tags=junior, department=sales}
2. [1] TestPerson{id=P002, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5}
3. [1] TestPerson{id=P003, salary=0, active=false, level=1, name=Charlie, age=16, score=6, tags=intern, status=inactive, department=hr}
4. [1] TestPerson{id=P004, level=7, age=45, active=true, score=7.8, status=active, name=Diana, salary=85000, tags=manager, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, active=false, score=8, department=sales, level=3, age=30, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, name=Frank, age=0, active=true, tags=test, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, tags=executive, status=active, department=management, level=9, salary=95000, active=true}
8. [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, salary=25000, active=false, status=inactive, level=1, age=18, score=5.5}
9. [1] TestPerson{id=P009, level=6, name=Ivy, age=40, tags=senior, salary=68000, active=true, score=8.7, status=active, department=engineering}
10. [1] TestPerson{id=P010, name=X, age=22, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1}
11. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, region=north, amount=2, date=2024-01-15, priority=normal, discount=50}
12. [1] TestOrder{id=O002, priority=low, discount=0, customer_id=P002, product_id=PROD002, total=25.5, region=south, amount=1, date=2024-01-20, status=confirmed}
13. [1] TestOrder{id=O003, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, region=north}
14. [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, region=east}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, discount=100, region=south, total=999.99, status=confirmed}
16. [1] TestOrder{id=O006, total=999.98, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, priority=low, region=west, customer_id=P005}
17. [1] TestOrder{id=O007, region=north, customer_id=P007, amount=4, total=600, date=2024-03-01, product_id=PROD006, status=shipped, priority=urgent, discount=50}
18. [1] TestOrder{id=O008, total=255, priority=normal, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002, amount=10}
19. [1] TestOrder{id=O009, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, status=completed, amount=1, total=89.99}
20. [1] TestOrder{id=O010, discount=0, amount=1, status=refunded, region=east, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, priority=urgent}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, department=sales, level=3, age=30, salary=55000, tags=employee, status=inactive}
   - Fait 2: [1] TestOrder{id=O006, priority=low, region=west, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, active=true, status=active, level=2, salary=45000, score=8.5, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, status=completed, amount=1}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, active=true, status=active, level=2, salary=45000, score=8.5, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O003, amount=3, date=2024-02-01, status=shipped, region=north, customer_id=P001, total=225, priority=high, discount=15, product_id=PROD003}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, status=active, level=5, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000}
   - Fait 2: [1] TestOrder{id=O005, status=confirmed, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, discount=100, region=south, total=999.99}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, active=true, name=Grace, age=65, score=10, tags=executive, status=active, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O007, region=north, customer_id=P007, amount=4, total=600, date=2024-03-01, product_id=PROD006, status=shipped, priority=urgent, discount=50}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, level=1, name=X, age=22, active=true, score=6.5, tags=temp, status=active, department=intern}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, priority=normal, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, department=qa, level=1, salary=-5000, score=0, status=active}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, region=east, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, priority=urgent, discount=0, amount=1}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, name=Alice, age=25, active=true, status=active, level=2, salary=45000, score=8.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, region=north, amount=2, date=2024-01-15, priority=normal, discount=50}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, level=5, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2, status=active}
   - Fait 2: [1] TestOrder{id=O002, discount=0, customer_id=P002, product_id=PROD002, total=25.5, region=south, amount=1, date=2024-01-20, status=confirmed, priority=low}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, tags=manager, department=marketing, level=7, age=45, active=true, score=7.8, status=active}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, region=east}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/20 (50.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, tags=junior, department=sales, name=Alice, age=25, active=true, status=active, level=2, salary=45000}
2. [1] TestPerson{id=P002, score=9.2, status=active, level=5, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000}
3. [1] TestPerson{id=P003, status=inactive, department=hr, salary=0, active=false, level=1, name=Charlie, age=16, score=6, tags=intern}
4. [1] TestPerson{id=P004, status=active, name=Diana, salary=85000, tags=manager, department=marketing, level=7, age=45, active=true, score=7.8}
5. [1] TestPerson{id=P005, tags=employee, status=inactive, name=Eve, active=false, score=8, department=sales, level=3, age=30, salary=55000}
6. [1] TestPerson{id=P006, department=qa, level=1, salary=-5000, score=0, status=active, name=Frank, age=0, active=true, tags=test}
7. [1] TestPerson{id=P007, score=10, tags=executive, status=active, department=management, level=9, salary=95000, active=true, name=Grace, age=65}
8. [1] TestPerson{id=P008, active=false, status=inactive, level=1, age=18, score=5.5, tags=junior, department=support, name=Henry, salary=25000}
9. [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, status=active, department=engineering, level=6, name=Ivy, age=40, tags=senior}
10. [1] TestPerson{id=P010, name=X, age=22, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, active=true, status=active, level=2, salary=45000, score=8.5, tags=junior, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, status=inactive, department=hr, salary=0, active=false, level=1}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, score=7.8, status=active, name=Diana, salary=85000, tags=manager, department=marketing, level=7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, department=sales, level=3, age=30, salary=55000, tags=employee, status=inactive}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, age=18, score=5.5, tags=junior, department=support, name=Henry, salary=25000, active=false, status=inactive, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, status=active, department=engineering, level=6, name=Ivy, age=40, tags=senior}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P002, status=active, level=5, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, status=active, department=management, level=9, salary=95000, active=true, name=Grace, age=65, score=10, tags=executive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, date=2024-01-15, priority=normal, discount=50, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, region=north, amount=2}
2. [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0, customer_id=P002, product_id=PROD002, total=25.5, region=south}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, region=north, customer_id=P001, total=225, priority=high, discount=15}
4. [1] TestOrder{id=O004, region=east, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, discount=100, region=south, total=999.99, status=confirmed}
6. [1] TestOrder{id=O006, region=west, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, priority=low}
7. [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, product_id=PROD006, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007}
8. [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, status=completed, amount=1, total=89.99, date=2024-03-10, priority=low, discount=10, region=north}
10. [1] TestOrder{id=O010, region=east, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, priority=urgent, discount=0, amount=1, status=refunded}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, region=north, amount=2, date=2024-01-15, priority=normal, discount=50}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, region=south, amount=1, date=2024-01-20, status=confirmed, priority=low, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, priority=normal, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, region=east}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, status=completed, amount=1, total=89.99, date=2024-03-10, priority=low, discount=10, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, total=225, priority=high, discount=15, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, region=north, customer_id=P001}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, date=2024-02-10, priority=high, discount=100, region=south, total=999.99, status=confirmed, customer_id=P002, product_id=PROD001}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, priority=low, region=west}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, amount=4, total=600, date=2024-03-01}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, total=255, priority=normal, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002, amount=10}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, total=75000, date=2024-03-15, priority=urgent, discount=0, amount=1, status=refunded, region=east}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, active=true, status=active, level=2, salary=45000, score=8.5, tags=junior, department=sales}
2. [1] TestPerson{id=P002, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5}
3. [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, status=inactive, department=hr, salary=0, active=false, level=1}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, tags=manager, department=marketing, level=7, age=45, active=true, score=7.8, status=active}
5. [1] TestPerson{id=P005, name=Eve, active=false, score=8, department=sales, level=3, age=30, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, name=Frank, age=0, active=true, tags=test, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, tags=executive, status=active, department=management, level=9, salary=95000, active=true}
8. [1] TestPerson{id=P008, score=5.5, tags=junior, department=support, name=Henry, salary=25000, active=false, status=inactive, level=1, age=18}
9. [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, status=active, department=engineering, level=6, name=Ivy, age=40, tags=senior}
10. [1] TestPerson{id=P010, name=X, age=22, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, status=active, department=management, level=9, salary=95000, active=true, name=Grace, age=65, score=10, tags=executive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, salary=25000, active=false, status=inactive, level=1, age=18, score=5.5}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, status=active, department=intern, salary=28000, level=1, name=X, age=22}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, department=sales, name=Alice, age=25, active=true, status=active, level=2, salary=45000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, status=active, level=5, name=Bob, active=true, tags=senior, department=engineering, age=35, salary=75000, score=9.2}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, level=1, name=Charlie, age=16, score=6, tags=intern, status=inactive, department=hr}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, name=Eve, active=false, score=8, department=sales, level=3, age=30}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, department=qa, level=1, salary=-5000, score=0, status=active, name=Frank, age=0, active=true}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, salary=68000, active=true, score=8.7, status=active, department=engineering, level=6, name=Ivy, age=40}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P004, status=active, name=Diana, salary=85000, tags=manager, department=marketing, level=7, age=45, active=true, score=7.8}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 150
- **Faits traités**: 27
