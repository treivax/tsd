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

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales}
2. [1] TestPerson{id=P002, active=true, tags=senior, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35}
3. [1] TestPerson{id=P003, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8, name=Diana}
5. [1] TestPerson{id=P005, age=30, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test, name=Frank, age=0}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
9. [1] TestPerson{id=P009, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive, department=management, level=9}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8, level=3}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, active=true, tags=senior, status=active, department=engineering, age=40, score=8.7, level=6, name=Ivy, salary=68000}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, age=22, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, active=true}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales, name=Alice}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, level=5, age=35, active=true, tags=senior, name=Bob, salary=75000, score=9.2}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north}
2. [1] TestOrder{id=O002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002}
3. [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0}
5. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}
6. [1] TestOrder{id=O006, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98, status=cancelled, priority=low, discount=0}
7. [1] TestOrder{id=O007, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50}
8. [1] TestOrder{id=O008, date=2024-03-05, discount=0, region=south, status=pending, priority=normal, customer_id=P010, product_id=PROD002, amount=10, total=255}
9. [1] TestOrder{id=O009, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, total=89.99}
10. [1] TestOrder{id=O010, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high, discount=15}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002, status=confirmed, priority=low, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, date=2024-02-05, priority=normal, discount=0, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2}
2. [1] TestPerson{id=P002, tags=senior, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true}
3. [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie}
4. [1] TestPerson{id=P004, department=marketing, level=7, age=45, score=7.8, name=Diana, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, age=30, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test, name=Frank, age=0}
7. [1] TestPerson{id=P007, department=management, level=9, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
9. [1] TestPerson{id=P009, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, age=22, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, active=true, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior, name=Bob, salary=75000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, age=45, score=7.8, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive, department=management, level=9}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering, age=40, score=8.7, level=6}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002, status=confirmed, priority=low, region=south}
3. [1] TestOrder{id=O003, region=north, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01}
4. [1] TestOrder{id=O004, discount=0, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal}
5. [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98, status=cancelled, priority=low, discount=0, region=west}
7. [1] TestOrder{id=O007, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}
9. [1] TestOrder{id=O009, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high, discount=15}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0, customer_id=P004, product_id=PROD004, total=299.99}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, keywords=computer, stock=50, supplier=TechSupply, name=Laptop, category=electronics, brand=TechCorp, price=999.99, available=true, rating=4.5}
2. [1] TestProduct{id=PROD002, price=25.5, available=true, stock=200, name=Mouse, category=accessories, rating=4, keywords=peripheral, brand=TechCorp, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, category=accessories, price=75, available=false, keywords=typing, stock=0, supplier=KeySupply, name=Keyboard, rating=3.5, brand=KeyTech}
4. [1] TestProduct{id=PROD004, price=299.99, available=true, rating=4.8, brand=ScreenPro, stock=30, name=Monitor, keywords=display, supplier=ScreenSupply, category=electronics}
5. [1] TestProduct{id=PROD005, brand=OldTech, supplier=OldSupply, name=OldKeyboard, category=accessories, price=8.5, available=false, keywords=obsolete, stock=0, rating=2}
6. [1] TestProduct{id=PROD006, available=true, rating=4.6, supplier=AudioSupply, name=Headphones, keywords=sound, brand=AudioMax, stock=75, category=audio, price=150}
7. [1] TestProduct{id=PROD007, price=89.99, available=true, keywords=video, stock=25, name=Webcam, category=electronics, rating=3.8, brand=CamTech, supplier=CamSupply}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD002, supplier=TechSupply, price=25.5, available=true, stock=200, name=Mouse, category=accessories, rating=4, keywords=peripheral, brand=TechCorp}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD003, category=accessories, price=75, available=false, keywords=typing, stock=0, supplier=KeySupply, name=Keyboard, rating=3.5, brand=KeyTech}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD004, category=electronics, price=299.99, available=true, rating=4.8, brand=ScreenPro, stock=30, name=Monitor, keywords=display, supplier=ScreenSupply}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD006, available=true, rating=4.6, supplier=AudioSupply, name=Headphones, keywords=sound, brand=AudioMax, stock=75, category=audio, price=150}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD007, brand=CamTech, supplier=CamSupply, price=89.99, available=true, keywords=video, stock=25, name=Webcam, category=electronics, rating=3.8}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD001, category=electronics, brand=TechCorp, price=999.99, available=true, rating=4.5, keywords=computer, stock=50, supplier=TechSupply, name=Laptop}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2}
2. [1] TestPerson{id=P002, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior, name=Bob, salary=75000}
3. [1] TestPerson{id=P003, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}
5. [1] TestPerson{id=P005, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30}
6. [1] TestPerson{id=P006, tags=test, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true}
7. [1] TestPerson{id=P007, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
9. [1] TestPerson{id=P009, department=engineering, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test, name=Frank}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north}
2. [1] TestOrder{id=O002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002}
3. [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north}
4. [1] TestOrder{id=O004, discount=0, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal}
5. [1] TestOrder{id=O005, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high}
6. [1] TestOrder{id=O006, date=2024-02-15, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent, customer_id=P007}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}
9. [1] TestOrder{id=O009, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10}
10. [1] TestOrder{id=O010, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0, customer_id=P004, product_id=PROD004}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent, customer_id=P007}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O008, status=pending, priority=normal, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15, status=refunded}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}
5. [1] TestPerson{id=P005, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30}
6. [1] TestPerson{id=P006, tags=test, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true}
7. [1] TestPerson{id=P007, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
9. [1] TestPerson{id=P009, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, department=intern, level=1, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8, level=3, name=Eve, active=false}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, age=18, level=1, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal}
2. [1] TestOrder{id=O002, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002, status=confirmed, priority=low}
3. [1] TestOrder{id=O003, region=north, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0}
5. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}
6. [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98}
7. [1] TestOrder{id=O007, status=shipped, discount=50, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, discount=0, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
3. [1] TestPerson{id=P003, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern}
4. [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8, name=Diana}
5. [1] TestPerson{id=P005, age=30, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, active=true, tags=test, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active, age=65}
8. [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1, name=Henry, salary=25000, active=false}
9. [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true}
10. [1] TestPerson{id=P010, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, active=true, tags=temp, age=22}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8, level=3}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, department=management, level=9, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1, age=16}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, age=18, level=1, name=Henry, salary=25000, active=false, score=5.5}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, status=active, department=engineering, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, price=999.99, available=true, rating=4.5, keywords=computer, stock=50, supplier=TechSupply, name=Laptop, category=electronics, brand=TechCorp}
2. [1] TestProduct{id=PROD002, rating=4, keywords=peripheral, brand=TechCorp, supplier=TechSupply, price=25.5, available=true, stock=200, name=Mouse, category=accessories}
3. [1] TestProduct{id=PROD003, supplier=KeySupply, name=Keyboard, rating=3.5, brand=KeyTech, category=accessories, price=75, available=false, keywords=typing, stock=0}
4. [1] TestProduct{id=PROD004, name=Monitor, keywords=display, supplier=ScreenSupply, category=electronics, price=299.99, available=true, rating=4.8, brand=ScreenPro, stock=30}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, price=8.5, available=false, keywords=obsolete, stock=0, rating=2, brand=OldTech, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, category=audio, price=150, available=true, rating=4.6, supplier=AudioSupply, name=Headphones, keywords=sound, brand=AudioMax, stock=75}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, rating=3.8, brand=CamTech, supplier=CamSupply, price=89.99, available=true, keywords=video, stock=25}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, category=electronics, brand=TechCorp, price=999.99, available=true, rating=4.5, keywords=computer, stock=50, supplier=TechSupply, name=Laptop}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, keywords=peripheral, brand=TechCorp, supplier=TechSupply, price=25.5, available=true, stock=200, name=Mouse, category=accessories, rating=4}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, category=accessories, price=75, available=false, keywords=typing, stock=0, supplier=KeySupply, name=Keyboard, rating=3.5, brand=KeyTech}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, rating=4.8, brand=ScreenPro, stock=30, name=Monitor, keywords=display, supplier=ScreenSupply, category=electronics, price=299.99, available=true}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, price=150, available=true, rating=4.6, supplier=AudioSupply, name=Headphones, keywords=sound, brand=AudioMax, stock=75, category=audio}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, price=89.99, available=true, keywords=video, stock=25, name=Webcam, category=electronics, rating=3.8, brand=CamTech, supplier=CamSupply}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25}
2. [1] TestPerson{id=P002, tags=senior, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true}
3. [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie}
4. [1] TestPerson{id=P004, age=45, score=7.8, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8, level=3}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test, name=Frank, age=0}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
9. [1] TestPerson{id=P009, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive, department=management, level=9}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8, level=3}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal}
2. [1] TestOrder{id=O002, date=2024-01-20, discount=0, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high, discount=15}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0}
5. [1] TestOrder{id=O005, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1}
6. [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}
9. [1] TestOrder{id=O009, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high, discount=15, customer_id=P001}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, date=2024-02-05, priority=normal, discount=0, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, discount=50, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O008, region=south, status=pending, priority=normal, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002, status=confirmed, priority=low, region=south}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
3. [1] TestPerson{id=P003, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, department=marketing, level=7, age=45, score=7.8, name=Diana, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, age=30, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, department=qa, level=1, active=true, tags=test, name=Frank, age=0, salary=-5000, score=0, status=active}
7. [1] TestPerson{id=P007, level=9, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive, department=management}
8. [1] TestPerson{id=P008, department=support, age=18, level=1, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering, age=40, score=8.7, level=6}
10. [1] TestPerson{id=P010, age=22, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, active=true, tags=temp}
11. [1] TestOrder{id=O001, product_id=PROD001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal}
12. [1] TestOrder{id=O002, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0}
13. [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north}
14. [1] TestOrder{id=O004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0, customer_id=P004, product_id=PROD004}
15. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}
16. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98, status=cancelled, priority=low, discount=0, region=west}
17. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}
18. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}
19. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north}
20. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P010, level=1, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5, status=active, department=intern}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30}
   - Fait 2: [1] TestOrder{id=O002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P006, status=active, department=qa, level=1, active=true, tags=test, name=Frank, age=0, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O005, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O009, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high, discount=15}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O004, discount=0, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test, name=Frank}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P009, status=active, department=engineering, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test, name=Frank, age=0}
   - Fait 2: [1] TestOrder{id=O002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P010, age=22, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, active=true, tags=temp}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0, customer_id=P004, product_id=PROD004}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P004, level=7, age=45, score=7.8, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing}
   - Fait 2: [1] TestOrder{id=O004, amount=1, date=2024-02-05, priority=normal, discount=0, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, department=sales, age=30, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P008, age=18, level=1, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support}
   - Fait 2: [1] TestOrder{id=O004, discount=0, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P001, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0, customer_id=P004}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P009, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering, age=40, score=8.7}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P010, age=22, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, active=true, tags=temp}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P004, level=7, age=45, score=7.8, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P005, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8, level=3, name=Eve}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal, customer_id=P010, product_id=PROD002, amount=10}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P006, status=active, department=qa, level=1, active=true, tags=test, name=Frank, age=0, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P002, active=true, tags=senior, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35}
   - Fait 2: [1] TestOrder{id=O001, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering, age=40, score=8.7, level=6}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P009, status=active, department=engineering, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, tags=senior, status=active, department=engineering, age=40, score=8.7, level=6, name=Ivy}
   - Fait 2: [1] TestOrder{id=O001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O002, discount=0, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
   - Fait 2: [1] TestOrder{id=O003, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P001, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high, discount=15}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern}
   - Fait 2: [1] TestOrder{id=O006, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98, status=cancelled, priority=low, discount=0}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P007, status=active, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal, customer_id=P010, product_id=PROD002, amount=10}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002, status=confirmed, priority=low, region=south}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, age=18, level=1, name=Henry, salary=25000, active=false, score=5.5}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales, name=Alice, salary=45000}
   - Fait 2: [1] TestOrder{id=O005, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8, level=3}
   - Fait 2: [1] TestOrder{id=O007, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P010, age=22, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, active=true, tags=temp}
   - Fait 2: [1] TestOrder{id=O008, status=pending, priority=normal, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, total=89.99}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P007, score=10, status=active, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace}
   - Fait 2: [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8, level=3}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0, customer_id=P004}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P010, active=true, tags=temp, age=22, score=6.5, status=active, department=intern, level=1, name=X, salary=28000}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8, name=Diana, salary=85000, active=true}
   - Fait 2: [1] TestOrder{id=O005, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P001, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal, customer_id=P010, product_id=PROD002, amount=10}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P010, age=22, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, active=true, tags=temp}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high, discount=15}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}
   - Fait 2: [1] TestOrder{id=O003, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, level=5, age=35, active=true, tags=senior, name=Bob, salary=75000, score=9.2, status=active}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8, level=3}
   - Fait 2: [1] TestOrder{id=O001, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8, name=Diana, salary=85000, active=true}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, level=1, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5, status=active}
   - Fait 2: [1] TestOrder{id=O002, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002, status=confirmed, priority=low}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, level=1, active=true, tags=test, name=Frank, age=0, salary=-5000, score=0, status=active}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98, status=cancelled, priority=low, discount=0, region=west}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, age=45, score=7.8, name=Diana, salary=85000, active=true, tags=manager, status=active}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P010, status=active, department=intern, level=1, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P007, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active, age=65, salary=95000}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, date=2024-01-20, discount=0, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5, status=active, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98, status=cancelled, priority=low, discount=0, region=west}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering, age=40}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, total=89.99}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering, age=40}
   - Fait 2: [1] TestOrder{id=O004, discount=0, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O006, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98, status=cancelled}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
   - Fait 2: [1] TestOrder{id=O006, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98, status=cancelled, priority=low, discount=0}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior}
   - Fait 2: [1] TestOrder{id=O002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002, status=confirmed, priority=low, region=south}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O003, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P003, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test, name=Frank}
   - Fait 2: [1] TestOrder{id=O007, discount=50, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P001, level=2, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, active=true, tags=temp, age=22}
   - Fait 2: [1] TestOrder{id=O010, amount=1, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
   - Fait 2: [1] TestOrder{id=O001, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P001, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P009, status=active, department=engineering, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive}
   - Fait 2: [1] TestOrder{id=O004, priority=normal, discount=0, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, level=1, active=true, tags=test, name=Frank, age=0, salary=-5000, score=0, status=active}
   - Fait 2: [1] TestOrder{id=O008, priority=normal, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}
   - Fait 2: [1] TestOrder{id=O010, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, age=30, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, level=5, age=35, active=true, tags=senior, name=Bob, salary=75000, score=9.2}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test, name=Frank, age=0}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high, discount=15}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8, level=3}
   - Fait 2: [1] TestOrder{id=O006, amount=2, date=2024-02-15, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent, customer_id=P007}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior}
2. [1] TestPerson{id=P002, department=engineering, level=5, age=35, active=true, tags=senior, name=Bob, salary=75000, score=9.2, status=active}
3. [1] TestPerson{id=P003, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}
5. [1] TestPerson{id=P005, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test}
7. [1] TestPerson{id=P007, tags=executive, department=management, level=9, name=Grace, score=10, status=active, age=65, salary=95000, active=true}
8. [1] TestPerson{id=P008, status=inactive, department=support, age=18, level=1, name=Henry, salary=25000, active=false, score=5.5, tags=junior}
9. [1] TestPerson{id=P009, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering, age=40, score=8.7}
10. [1] TestPerson{id=P010, active=true, tags=temp, age=22, score=6.5, status=active, department=intern, level=1, name=X, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, active=true, tags=senior, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8, name=Diana}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, status=active, department=engineering, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, age=22, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, active=true, tags=temp}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2}
2. [1] TestPerson{id=P002, active=true, tags=senior, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35}
3. [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}
5. [1] TestPerson{id=P005, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8, level=3}
6. [1] TestPerson{id=P006, active=true, tags=test, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
9. [1] TestPerson{id=P009, status=active, department=engineering, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior}
10. [1] TestPerson{id=P010, status=active, department=intern, level=1, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5}
11. [1] TestOrder{id=O001, priority=normal, product_id=PROD001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending}
12. [1] TestOrder{id=O002, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0}
13. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high, discount=15}
14. [1] TestOrder{id=O004, amount=1, date=2024-02-05, priority=normal, discount=0, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east}
15. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}
16. [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15}
17. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}
18. [1] TestOrder{id=O008, status=pending, priority=normal, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south}
19. [1] TestOrder{id=O009, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1}
20. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering, age=40, score=8.7, level=6}
   - Fait 2: [1] TestOrder{id=O003, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5, status=active, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test}
   - Fait 2: [1] TestOrder{id=O006, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98, status=cancelled, priority=low, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002, status=confirmed, priority=low, region=south}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8, level=3, name=Eve, active=false}
   - Fait 2: [1] TestOrder{id=O004, amount=1, date=2024-02-05, priority=normal, discount=0, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P006, active=true, tags=test, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P003, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, age=30, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering, age=40, score=8.7, level=6}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P007, department=management, level=9, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O003, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P008, level=1, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P007, score=10, status=active, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace}
   - Fait 2: [1] TestOrder{id=O005, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P007, department=management, level=9, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1, name=Henry, salary=25000}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales, name=Alice}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, total=89.99}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P009, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O007, discount=50, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O007, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P003, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O009, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P004, level=7, age=45, score=7.8, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P010, status=active, department=intern, level=1, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P010, age=22, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, active=true, tags=temp}
   - Fait 2: [1] TestOrder{id=O003, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8, level=3}
   - Fait 2: [1] TestOrder{id=O002, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002, status=confirmed}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0, customer_id=P004}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
   - Fait 2: [1] TestOrder{id=O005, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1, name=Henry, salary=25000}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P005, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8, level=3, name=Eve}
   - Fait 2: [1] TestOrder{id=O001, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
   - Fait 2: [1] TestOrder{id=O002, discount=0, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high, discount=15}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}
   - Fait 2: [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test, name=Frank}
   - Fait 2: [1] TestOrder{id=O004, discount=0, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O005, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P005, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8}
   - Fait 2: [1] TestOrder{id=O006, date=2024-02-15, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5, status=active, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test, name=Frank, age=0}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, level=5, age=35, active=true, tags=senior, name=Bob, salary=75000, score=9.2}
   - Fait 2: [1] TestOrder{id=O009, discount=10, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007, total=89.99, priority=low}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001, product_id=PROD007}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P009, active=true, tags=senior, status=active, department=engineering, age=40, score=8.7, level=6, name=Ivy, salary=68000}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P001, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P009, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering, age=40, score=8.7}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering, age=40, score=8.7, level=6}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98, status=cancelled, priority=low}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, score=10, status=active, age=65, salary=95000, active=true, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P006, active=true, tags=test, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie}
   - Fait 2: [1] TestOrder{id=O002, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002, status=confirmed, priority=low}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, level=1, active=true, tags=test, name=Frank, age=0, salary=-5000, score=0, status=active}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north, amount=3, status=shipped, priority=high, discount=15}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P002, age=35, active=true, tags=senior, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O004, amount=1, date=2024-02-05, priority=normal, discount=0, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P004, age=45, score=7.8, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, priority=normal, discount=0, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal, customer_id=P010, product_id=PROD002, amount=10}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1, name=Henry, salary=25000}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active, age=65}
   - Fait 2: [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, active=true, tags=temp, age=22}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5, status=active, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O008, region=south, status=pending, priority=normal, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P001, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering, age=40, score=8.7, level=6}
   - Fait 2: [1] TestOrder{id=O002, discount=0, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P005, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true}
   - Fait 2: [1] TestOrder{id=O010, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P001, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P006, active=true, tags=test, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000, date=2024-03-15}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O003, region=north, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, active=true, tags=temp, age=22}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}
   - Fait 2: [1] TestOrder{id=O002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, age=22, score=6.5, status=active, department=intern, level=1, name=X}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, age=22, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, active=true}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98, status=cancelled, priority=low, discount=0, region=west}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, age=22, score=6.5, status=active, department=intern, level=1, name=X}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P005, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000}
   - Fait 2: [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P002, level=5, age=35, active=true, tags=senior, name=Bob, salary=75000, score=9.2, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}
   - Fait 2: [1] TestOrder{id=O005, discount=100, customer_id=P002, product_id=PROD001, region=south, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P002, age=35, active=true, tags=senior, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98, status=cancelled, priority=low, discount=0, region=west}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, total=999.98, status=cancelled, priority=low, discount=0, region=west}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0, customer_id=P004}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P009, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1, name=Henry, salary=25000, active=false}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, tags=senior, status=active, department=engineering, age=40, score=8.7, level=6, name=Ivy}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal, customer_id=P010, product_id=PROD002}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P001, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P005, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east, total=75000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
3. [1] TestPerson{id=P003, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8, name=Diana, salary=85000}
5. [1] TestPerson{id=P005, status=inactive, department=sales, age=30, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true, tags=test}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}
9. [1] TestPerson{id=P009, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8, name=Diana, salary=85000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, department=sales, age=30, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior, name=Bob, salary=75000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P003, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P009, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5, status=active, department=intern, level=1}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, status=pending, priority=normal, product_id=PROD001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15}
2. [1] TestOrder{id=O002, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002, amount=1, total=25.5, date=2024-01-20, discount=0}
3. [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north}
4. [1] TestOrder{id=O004, discount=0, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal}
5. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}
6. [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50}
8. [1] TestOrder{id=O008, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal, customer_id=P010, product_id=PROD002, amount=10}
9. [1] TestOrder{id=O009, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north, customer_id=P001}
10. [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, amount=1, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, priority=low, discount=10, amount=1, date=2024-03-10, status=completed, region=north}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O010, amount=1, region=east, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, product_id=PROD001}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, total=25.5, date=2024-01-20, discount=0, product_id=PROD002, status=confirmed, priority=low, region=south, customer_id=P002}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, region=east, amount=1, date=2024-02-05, priority=normal, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, discount=0, region=south, status=pending, priority=normal}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O001, amount=2, discount=50, region=north, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, priority=normal, product_id=PROD001}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, status=shipped, priority=high, discount=15, customer_id=P001, product_id=PROD003, total=225, date=2024-02-01, region=north}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, region=south}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, tags=junior, department=sales, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5, age=35, active=true, tags=senior}
3. [1] TestPerson{id=P003, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}
5. [1] TestPerson{id=P005, salary=55000, score=8, level=3, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30}
6. [1] TestPerson{id=P006, active=true, tags=test, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, score=10, status=active, age=65, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace}
8. [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1, name=Henry}
9. [1] TestPerson{id=P009, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, department=intern, level=1, name=X, salary=28000, active=true, tags=temp, age=22, score=6.5, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, tags=employee, status=inactive, department=sales, age=30, salary=55000, score=8, level=3}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, level=1, active=true}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, active=true, tags=executive, department=management, level=9, name=Grace, score=10, status=active, age=65}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, department=support, age=18, level=1, name=Henry, salary=25000, active=false}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, age=40, score=8.7, level=6, name=Ivy, salary=68000, active=true, tags=senior, status=active, department=engineering}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, age=22, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, active=true}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, status=active, level=2, age=25, tags=junior, department=sales}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P002, age=35, active=true, tags=senior, name=Bob, salary=75000, score=9.2, status=active, department=engineering, level=5}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, department=marketing, level=7, age=45, score=7.8}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
