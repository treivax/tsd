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

1. [1] TestPerson{id=P001, score=8.5, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, tags=junior, status=active}
2. [1] TestPerson{id=P002, tags=senior, status=active, name=Bob, age=35, salary=75000, department=engineering, level=5, active=true, score=9.2}
3. [1] TestPerson{id=P003, age=16, tags=intern, status=inactive, level=1, salary=0, active=false, score=6, department=hr, name=Charlie}
4. [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}
5. [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, name=Eve, age=30, active=false, score=8, department=sales, level=3}
6. [1] TestPerson{id=P006, age=0, salary=-5000, active=true, score=0, tags=test, status=active, department=qa, level=1, name=Frank}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, level=9, name=Grace, active=true, score=10, department=management}
8. [1] TestPerson{id=P008, name=Henry, active=false, department=support, level=1, age=18, salary=25000, score=5.5, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, name=Ivy, tags=senior, level=6}
10. [1] TestPerson{id=P010, tags=temp, age=22, active=true, score=6.5, status=active, department=intern, level=1, name=X, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, tags=intern, status=inactive, level=1, salary=0, active=false, score=6, department=hr}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, department=sales, level=3, salary=55000, tags=employee, status=inactive}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, tags=temp, age=22, active=true, score=6.5, status=active, department=intern, level=1, name=X}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, department=management, age=65, salary=95000, tags=executive, status=active, level=9, name=Grace, active=true, score=10}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, department=support, level=1, age=18, salary=25000, score=5.5, tags=junior, status=inactive}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, tags=senior, level=6, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, active=true, tags=junior, status=active, score=8.5, department=sales, level=2, name=Alice, age=25, salary=45000}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, department=engineering, level=5, active=true, score=9.2, tags=senior, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, priority=normal, discount=50, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, status=pending, region=north, customer_id=P001}
2. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, region=south, customer_id=P002, product_id=PROD002, status=confirmed, priority=low, discount=0, amount=1}
3. [1] TestOrder{id=O003, amount=3, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, total=225, status=shipped, product_id=PROD003}
4. [1] TestOrder{id=O004, status=delivered, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, status=confirmed, priority=high, discount=100, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, region=south}
6. [1] TestOrder{id=O006, date=2024-02-15, discount=0, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, priority=low, region=west, amount=2}
7. [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north}
8. [1] TestOrder{id=O008, discount=0, region=south, customer_id=P010, amount=10, status=pending, product_id=PROD002, total=255, date=2024-03-05, priority=normal}
9. [1] TestOrder{id=O009, region=north, customer_id=P001, product_id=PROD007, amount=1, priority=low, total=89.99, date=2024-03-10, status=completed, discount=10}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, date=2024-03-15, status=refunded, priority=urgent, region=east, product_id=PROD001, total=75000, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O009, date=2024-03-10, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, priority=low, total=89.99}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O010, discount=0, customer_id=P006, amount=1, date=2024-03-15, status=refunded, priority=urgent, region=east, product_id=PROD001, total=75000}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, status=pending, region=north, customer_id=P001, priority=normal, discount=50}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, total=225, status=shipped, product_id=PROD003, amount=3}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, region=south, customer_id=P002, status=confirmed, priority=high, discount=100}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O008, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, amount=10, status=pending, product_id=PROD002, total=255}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, status=confirmed, priority=low, discount=0, amount=1, total=25.5, date=2024-01-20, region=south}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, region=east, status=delivered, priority=normal, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, active=true, tags=junior, status=active, score=8.5, department=sales, level=2, name=Alice, age=25}
2. [1] TestPerson{id=P002, status=active, name=Bob, age=35, salary=75000, department=engineering, level=5, active=true, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, score=6, department=hr, name=Charlie, age=16, tags=intern, status=inactive, level=1, salary=0, active=false}
4. [1] TestPerson{id=P004, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana}
5. [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, department=sales, level=3, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, salary=-5000, active=true, score=0, tags=test, status=active, department=qa, level=1, name=Frank, age=0}
7. [1] TestPerson{id=P007, department=management, age=65, salary=95000, tags=executive, status=active, level=9, name=Grace, active=true, score=10}
8. [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, tags=junior, status=inactive, name=Henry, active=false, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, tags=senior, level=6, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering}
10. [1] TestPerson{id=P010, level=1, name=X, salary=28000, tags=temp, age=22, active=true, score=6.5, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing, age=45, salary=85000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, department=sales, level=3, salary=55000, tags=employee, status=inactive}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, department=management, age=65, salary=95000, tags=executive, status=active, level=9}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, status=active, department=engineering, name=Ivy, tags=senior, level=6, age=40}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, tags=junior, status=active, score=8.5, department=sales, level=2}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, department=engineering, level=5, active=true, score=9.2, tags=senior, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, status=pending, region=north, customer_id=P001, priority=normal, discount=50, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15}
2. [1] TestOrder{id=O002, amount=1, total=25.5, date=2024-01-20, region=south, customer_id=P002, product_id=PROD002, status=confirmed, priority=low, discount=0}
3. [1] TestOrder{id=O003, discount=15, region=north, customer_id=P001, total=225, status=shipped, product_id=PROD003, amount=3, date=2024-02-01, priority=high}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, region=east, status=delivered, priority=normal, discount=0}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, region=south, customer_id=P002, status=confirmed, priority=high, discount=100}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, priority=low, region=west, amount=2, date=2024-02-15, discount=0}
7. [1] TestOrder{id=O007, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent}
8. [1] TestOrder{id=O008, product_id=PROD002, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, amount=10, status=pending}
9. [1] TestOrder{id=O009, product_id=PROD007, amount=1, priority=low, total=89.99, date=2024-03-10, status=completed, discount=10, region=north, customer_id=P001}
10. [1] TestOrder{id=O010, product_id=PROD001, total=75000, discount=0, customer_id=P006, amount=1, date=2024-03-15, status=refunded, priority=urgent, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, status=confirmed, priority=high, discount=100, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, priority=low, region=west, amount=2, date=2024-02-15, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, status=pending, region=north, customer_id=P001, priority=normal, discount=50}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, total=25.5, date=2024-01-20, region=south, customer_id=P002, product_id=PROD002, status=confirmed, priority=low, discount=0}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, total=225, status=shipped, product_id=PROD003}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, region=east, status=delivered, priority=normal, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O008, priority=normal, discount=0, region=south, customer_id=P010, amount=10, status=pending, product_id=PROD002, total=255, date=2024-03-05}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, priority=low, total=89.99, date=2024-03-10, status=completed}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, price=999.99, supplier=TechSupply, category=electronics, available=true, rating=4.5, keywords=computer, brand=TechCorp, stock=50}
2. [1] TestProduct{id=PROD002, price=25.5, rating=4, brand=TechCorp, stock=200, supplier=TechSupply, name=Mouse, category=accessories, available=true, keywords=peripheral}
3. [1] TestProduct{id=PROD003, name=Keyboard, price=75, keywords=typing, brand=KeyTech, supplier=KeySupply, category=accessories, available=false, rating=3.5, stock=0}
4. [1] TestProduct{id=PROD004, category=electronics, rating=4.8, keywords=display, brand=ScreenPro, supplier=ScreenSupply, name=Monitor, price=299.99, available=true, stock=30}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, available=false, rating=2, keywords=obsolete, brand=OldTech, stock=0, price=8.5, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, rating=4.6, keywords=sound, supplier=AudioSupply, name=Headphones, brand=AudioMax, stock=75, category=audio, price=150, available=true}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, price=89.99, available=true, stock=25, supplier=CamSupply, rating=3.8, keywords=video, brand=CamTech}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, keywords=computer, brand=TechCorp, stock=50, name=Laptop, price=999.99, supplier=TechSupply, category=electronics, available=true, rating=4.5}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, category=accessories, available=true, keywords=peripheral, price=25.5, rating=4, brand=TechCorp, stock=200, supplier=TechSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, brand=KeyTech, supplier=KeySupply, category=accessories, available=false, rating=3.5, stock=0, name=Keyboard, price=75, keywords=typing}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, rating=4.8, keywords=display, brand=ScreenPro, supplier=ScreenSupply, name=Monitor, price=299.99, available=true, stock=30, category=electronics}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, rating=4.6, keywords=sound, supplier=AudioSupply, name=Headphones, brand=AudioMax, stock=75, category=audio, price=150, available=true}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, name=Webcam, category=electronics, price=89.99, available=true, stock=25, supplier=CamSupply, rating=3.8, keywords=video, brand=CamTech}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, tags=junior, status=active}
2. [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, name=Bob, age=35, salary=75000, department=engineering, level=5}
3. [1] TestPerson{id=P003, salary=0, active=false, score=6, department=hr, name=Charlie, age=16, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, department=sales, level=3, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, salary=-5000, active=true, score=0, tags=test, status=active, department=qa, level=1, name=Frank, age=0}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, department=management, age=65, salary=95000, tags=executive, status=active, level=9}
8. [1] TestPerson{id=P008, name=Henry, active=false, department=support, level=1, age=18, salary=25000, score=5.5, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, tags=senior, level=6, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering}
10. [1] TestPerson{id=P010, name=X, salary=28000, tags=temp, age=22, active=true, score=6.5, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, status=inactive, level=1, salary=0, active=false, score=6, department=hr, name=Charlie, age=16}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, active=true, score=0, tags=test, status=active, department=qa, level=1, name=Frank}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, priority=normal, discount=50, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, status=pending, region=north}
2. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, region=south, customer_id=P002, product_id=PROD002, status=confirmed, priority=low, discount=0, amount=1}
3. [1] TestOrder{id=O003, customer_id=P001, total=225, status=shipped, product_id=PROD003, amount=3, date=2024-02-01, priority=high, discount=15, region=north}
4. [1] TestOrder{id=O004, discount=0, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, region=east, status=delivered, priority=normal}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, region=south, customer_id=P002, status=confirmed, priority=high, discount=100}
6. [1] TestOrder{id=O006, discount=0, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, priority=low, region=west, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north}
8. [1] TestOrder{id=O008, product_id=PROD002, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, amount=10, status=pending}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, priority=low}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, date=2024-03-15, status=refunded, priority=urgent, region=east, product_id=PROD001, total=75000, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, status=pending, region=north, customer_id=P001, priority=normal, discount=50}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, status=confirmed, priority=low, discount=0, amount=1, total=25.5, date=2024-01-20, region=south}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, region=east, status=delivered, priority=normal, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, status=confirmed, priority=high, discount=100, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, region=south}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, region=west, amount=2, date=2024-02-15, discount=0, customer_id=P005}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, priority=low}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, total=225, status=shipped}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, amount=10, status=pending, product_id=PROD002, total=255, date=2024-03-05, priority=normal, discount=0, region=south}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, total=75000, discount=0, customer_id=P006, amount=1, date=2024-03-15, status=refunded, priority=urgent, region=east, product_id=PROD001}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, score=8.5, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, tags=junior}
2. [1] TestPerson{id=P002, department=engineering, level=5, active=true, score=9.2, tags=senior, status=active, name=Bob, age=35, salary=75000}
3. [1] TestPerson{id=P003, department=hr, name=Charlie, age=16, tags=intern, status=inactive, level=1, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true}
5. [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, name=Eve, age=30, active=false, score=8, department=sales, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, score=0, tags=test, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, level=9, name=Grace, active=true, score=10, department=management}
8. [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, tags=junior, status=inactive, name=Henry, active=false, department=support, level=1}
9. [1] TestPerson{id=P009, status=active, department=engineering, name=Ivy, tags=senior, level=6, age=40, salary=68000, active=true, score=8.7}
10. [1] TestPerson{id=P010, age=22, active=true, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, name=Charlie, age=16, tags=intern, status=inactive, level=1, salary=0, active=false, score=6}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, name=Eve, age=30, active=false, score=8, department=sales, level=3, salary=55000, tags=employee}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, tags=junior, status=inactive, name=Henry, active=false, department=support, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, discount=50, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, status=pending, region=north, customer_id=P001, priority=normal}
2. [1] TestOrder{id=O002, product_id=PROD002, status=confirmed, priority=low, discount=0, amount=1, total=25.5, date=2024-01-20, region=south, customer_id=P002}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, total=225, status=shipped}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, region=east, status=delivered, priority=normal, discount=0, customer_id=P004}
5. [1] TestOrder{id=O005, date=2024-02-10, region=south, customer_id=P002, status=confirmed, priority=high, discount=100, product_id=PROD001, amount=1, total=999.99}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, priority=low, region=west, amount=2, date=2024-02-15, discount=0}
7. [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north}
8. [1] TestOrder{id=O008, product_id=PROD002, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, amount=10, status=pending}
9. [1] TestOrder{id=O009, amount=1, priority=low, total=89.99, date=2024-03-10, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, date=2024-03-15, status=refunded, priority=urgent, region=east, product_id=PROD001, total=75000, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, status=confirmed, priority=low, discount=0, amount=1, total=25.5, date=2024-01-20, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, total=225, status=shipped, product_id=PROD003, amount=3, date=2024-02-01, priority=high, discount=15, region=north}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, region=east, status=delivered, priority=normal, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, status=confirmed, priority=high, discount=100, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, region=south}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, priority=low, region=west, amount=2, date=2024-02-15, discount=0, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, priority=low}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, amount=1, date=2024-03-15, status=refunded, priority=urgent, region=east, product_id=PROD001, total=75000, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, tags=junior, status=active}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, department=engineering, level=5, active=true, score=9.2, tags=senior, status=active}
3. [1] TestPerson{id=P003, salary=0, active=false, score=6, department=hr, name=Charlie, age=16, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7}
5. [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, department=sales, level=3, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, salary=-5000, active=true, score=0, tags=test, status=active, department=qa, level=1, name=Frank, age=0}
7. [1] TestPerson{id=P007, tags=executive, status=active, level=9, name=Grace, active=true, score=10, department=management, age=65, salary=95000}
8. [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, tags=junior, status=inactive, name=Henry, active=false, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, tags=senior, level=6, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering}
10. [1] TestPerson{id=P010, level=1, name=X, salary=28000, tags=temp, age=22, active=true, score=6.5, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, tags=junior, status=active, score=8.5}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, active=false, department=support, level=1, age=18, salary=25000, score=5.5, tags=junior, status=inactive, name=Henry}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, department=engineering, level=5, active=true, score=9.2, tags=senior, status=active}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P003, score=6, department=hr, name=Charlie, age=16, tags=intern, status=inactive, level=1, salary=0, active=false}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, level=3, salary=55000, tags=employee, status=inactive, name=Eve, age=30, active=false, score=8}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, active=true, score=0, tags=test, status=active, department=qa, level=1, name=Frank, age=0}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, department=management, age=65, salary=95000, tags=executive, status=active, level=9, name=Grace}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, name=Ivy, tags=senior, level=6, age=40, salary=68000, active=true, score=8.7, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, price=999.99, supplier=TechSupply, category=electronics, available=true, rating=4.5, keywords=computer, brand=TechCorp, stock=50}
2. [1] TestProduct{id=PROD002, price=25.5, rating=4, brand=TechCorp, stock=200, supplier=TechSupply, name=Mouse, category=accessories, available=true, keywords=peripheral}
3. [1] TestProduct{id=PROD003, brand=KeyTech, supplier=KeySupply, category=accessories, available=false, rating=3.5, stock=0, name=Keyboard, price=75, keywords=typing}
4. [1] TestProduct{id=PROD004, category=electronics, rating=4.8, keywords=display, brand=ScreenPro, supplier=ScreenSupply, name=Monitor, price=299.99, available=true, stock=30}
5. [1] TestProduct{id=PROD005, supplier=OldSupply, name=OldKeyboard, category=accessories, available=false, rating=2, keywords=obsolete, brand=OldTech, stock=0, price=8.5}
6. [1] TestProduct{id=PROD006, category=audio, price=150, available=true, rating=4.6, keywords=sound, supplier=AudioSupply, name=Headphones, brand=AudioMax, stock=75}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, price=89.99, available=true, stock=25, supplier=CamSupply, rating=3.8, keywords=video, brand=CamTech}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, category=electronics, available=true, rating=4.5, keywords=computer, brand=TechCorp, stock=50, name=Laptop, price=999.99, supplier=TechSupply}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, price=25.5, rating=4, brand=TechCorp, stock=200, supplier=TechSupply, name=Mouse, category=accessories, available=true, keywords=peripheral}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, stock=0, name=Keyboard, price=75, keywords=typing, brand=KeyTech, supplier=KeySupply, category=accessories, available=false, rating=3.5}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, keywords=display, brand=ScreenPro, supplier=ScreenSupply, name=Monitor, price=299.99, available=true, stock=30, category=electronics, rating=4.8}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, stock=75, category=audio, price=150, available=true, rating=4.6, keywords=sound, supplier=AudioSupply, name=Headphones, brand=AudioMax}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, brand=CamTech, name=Webcam, category=electronics, price=89.99, available=true, stock=25, supplier=CamSupply, rating=3.8, keywords=video}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, tags=junior, status=active}
2. [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, name=Bob, age=35, salary=75000, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, age=16, tags=intern, status=inactive, level=1, salary=0, active=false, score=6, department=hr}
4. [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}
5. [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, name=Eve, age=30, active=false, score=8, department=sales, level=3}
6. [1] TestPerson{id=P006, score=0, tags=test, status=active, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, level=9, name=Grace, active=true, score=10, department=management}
8. [1] TestPerson{id=P008, name=Henry, active=false, department=support, level=1, age=18, salary=25000, score=5.5, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, tags=senior, level=6, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering}
10. [1] TestPerson{id=P010, age=22, active=true, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, name=Bob, age=35, salary=75000, department=engineering, level=5}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, score=6, department=hr, name=Charlie, age=16, tags=intern, status=inactive, level=1, salary=0, active=false}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, score=0, tags=test, status=active, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, department=management, age=65, salary=95000, tags=executive, status=active, level=9}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, tags=junior, status=active, score=8.5, department=sales, level=2}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, level=3, salary=55000, tags=employee, status=inactive, name=Eve, age=30, active=false, score=8, department=sales}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, score=5.5, tags=junior, status=inactive, name=Henry, active=false, department=support, level=1, age=18}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, name=Ivy, tags=senior, level=6}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, status=pending, region=north, customer_id=P001, priority=normal, discount=50}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, status=confirmed, priority=low, discount=0, amount=1, total=25.5, date=2024-01-20, region=south}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, total=225, status=shipped}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, region=east, status=delivered, priority=normal, discount=0}
5. [1] TestOrder{id=O005, customer_id=P002, status=confirmed, priority=high, discount=100, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, region=south}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, priority=low, region=west, amount=2, date=2024-02-15, discount=0}
7. [1] TestOrder{id=O007, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600}
8. [1] TestOrder{id=O008, product_id=PROD002, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, amount=10, status=pending}
9. [1] TestOrder{id=O009, amount=1, priority=low, total=89.99, date=2024-03-10, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007}
10. [1] TestOrder{id=O010, region=east, product_id=PROD001, total=75000, discount=0, customer_id=P006, amount=1, date=2024-03-15, status=refunded, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O007, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O008, status=pending, product_id=PROD002, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, amount=10}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O009, date=2024-03-10, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, priority=low, total=89.99}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, status=pending, region=north, customer_id=P001, priority=normal, discount=50}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O002, region=south, customer_id=P002, product_id=PROD002, status=confirmed, priority=low, discount=0, amount=1, total=25.5, date=2024-01-20}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, total=225, status=shipped, product_id=PROD003}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, region=east, status=delivered, priority=normal, discount=0, customer_id=P004}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, discount=100, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, region=south, customer_id=P002, status=confirmed, priority=high}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, score=8.5, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, tags=junior}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, department=engineering, level=5, active=true, score=9.2, tags=senior, status=active}
3. [1] TestPerson{id=P003, status=inactive, level=1, salary=0, active=false, score=6, department=hr, name=Charlie, age=16, tags=intern}
4. [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, department=sales, level=3, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, salary=-5000, active=true, score=0, tags=test, status=active, department=qa, level=1, name=Frank, age=0}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, department=management, age=65, salary=95000, tags=executive, status=active, level=9}
8. [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, tags=junior, status=inactive, name=Henry, active=false, department=support, level=1}
9. [1] TestPerson{id=P009, active=true, score=8.7, status=active, department=engineering, name=Ivy, tags=senior, level=6, age=40, salary=68000}
10. [1] TestPerson{id=P010, age=22, active=true, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, tags=temp}
11. [1] TestOrder{id=O001, date=2024-01-15, status=pending, region=north, customer_id=P001, priority=normal, discount=50, product_id=PROD001, amount=2, total=1999.98}
12. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, status=confirmed, priority=low, discount=0, amount=1, total=25.5, date=2024-01-20, region=south}
13. [1] TestOrder{id=O003, status=shipped, product_id=PROD003, amount=3, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, total=225}
14. [1] TestOrder{id=O004, region=east, status=delivered, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05}
15. [1] TestOrder{id=O005, date=2024-02-10, region=south, customer_id=P002, status=confirmed, priority=high, discount=100, product_id=PROD001, amount=1, total=999.99}
16. [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, region=west, amount=2, date=2024-02-15, discount=0, customer_id=P005, product_id=PROD005}
17. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50}
18. [1] TestOrder{id=O008, region=south, customer_id=P010, amount=10, status=pending, product_id=PROD002, total=255, date=2024-03-05, priority=normal, discount=0}
19. [1] TestOrder{id=O009, priority=low, total=89.99, date=2024-03-10, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1}
20. [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, region=east, product_id=PROD001, total=75000, discount=0, customer_id=P006, amount=1}

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

1. [1] TestPerson{id=P001, status=active, score=8.5, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, tags=junior}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, status=active, name=Bob, age=35, salary=75000, department=engineering, level=5, active=true}
3. [1] TestPerson{id=P003, name=Charlie, age=16, tags=intern, status=inactive, level=1, salary=0, active=false, score=6, department=hr}
4. [1] TestPerson{id=P004, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing, age=45, salary=85000}
5. [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, department=sales, level=3, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, salary=-5000, active=true, score=0, tags=test, status=active, department=qa, level=1, name=Frank, age=0}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, department=management, age=65, salary=95000, tags=executive, status=active, level=9}
8. [1] TestPerson{id=P008, tags=junior, status=inactive, name=Henry, active=false, department=support, level=1, age=18, salary=25000, score=5.5}
9. [1] TestPerson{id=P009, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, name=Ivy, tags=senior, level=6}
10. [1] TestPerson{id=P010, age=22, active=true, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, status=active, score=8.5, department=sales, level=2, name=Alice, age=25, salary=45000, active=true}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, department=engineering, level=5, active=true, score=9.2, tags=senior, status=active}
3. [1] TestPerson{id=P003, level=1, salary=0, active=false, score=6, department=hr, name=Charlie, age=16, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active}
5. [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, name=Eve, age=30, active=false, score=8, department=sales, level=3}
6. [1] TestPerson{id=P006, salary=-5000, active=true, score=0, tags=test, status=active, department=qa, level=1, name=Frank, age=0}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, department=management, age=65, salary=95000, tags=executive, status=active, level=9}
8. [1] TestPerson{id=P008, name=Henry, active=false, department=support, level=1, age=18, salary=25000, score=5.5, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, name=Ivy, tags=senior, level=6}
10. [1] TestPerson{id=P010, name=X, salary=28000, tags=temp, age=22, active=true, score=6.5, status=active, department=intern, level=1}
11. [1] TestOrder{id=O001, amount=2, total=1999.98, date=2024-01-15, status=pending, region=north, customer_id=P001, priority=normal, discount=50, product_id=PROD001}
12. [1] TestOrder{id=O002, status=confirmed, priority=low, discount=0, amount=1, total=25.5, date=2024-01-20, region=south, customer_id=P002, product_id=PROD002}
13. [1] TestOrder{id=O003, customer_id=P001, total=225, status=shipped, product_id=PROD003, amount=3, date=2024-02-01, priority=high, discount=15, region=north}
14. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, region=east, status=delivered, priority=normal, discount=0}
15. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, region=south, customer_id=P002, status=confirmed, priority=high, discount=100, product_id=PROD001}
16. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, priority=low, region=west, amount=2, date=2024-02-15, discount=0}
17. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50}
18. [1] TestOrder{id=O008, discount=0, region=south, customer_id=P010, amount=10, status=pending, product_id=PROD002, total=255, date=2024-03-05, priority=normal}
19. [1] TestOrder{id=O009, amount=1, priority=low, total=89.99, date=2024-03-10, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007}
20. [1] TestOrder{id=O010, status=refunded, priority=urgent, region=east, product_id=PROD001, total=75000, discount=0, customer_id=P006, amount=1, date=2024-03-15}

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

1. [1] TestPerson{id=P001, active=true, tags=junior, status=active, score=8.5, department=sales, level=2, name=Alice, age=25, salary=45000}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, department=engineering, level=5, active=true, score=9.2, tags=senior, status=active}
3. [1] TestPerson{id=P003, name=Charlie, age=16, tags=intern, status=inactive, level=1, salary=0, active=false, score=6, department=hr}
4. [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}
5. [1] TestPerson{id=P005, department=sales, level=3, salary=55000, tags=employee, status=inactive, name=Eve, age=30, active=false, score=8}
6. [1] TestPerson{id=P006, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true, score=0, tags=test, status=active}
7. [1] TestPerson{id=P007, active=true, score=10, department=management, age=65, salary=95000, tags=executive, status=active, level=9, name=Grace}
8. [1] TestPerson{id=P008, name=Henry, active=false, department=support, level=1, age=18, salary=25000, score=5.5, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, department=engineering, name=Ivy, tags=senior, level=6, age=40, salary=68000, active=true, score=8.7, status=active}
10. [1] TestPerson{id=P010, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, tags=temp, age=22, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, region=north, customer_id=P001, priority=normal, discount=50, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, status=pending}
2. [1] TestOrder{id=O002, amount=1, total=25.5, date=2024-01-20, region=south, customer_id=P002, product_id=PROD002, status=confirmed, priority=low, discount=0}
3. [1] TestOrder{id=O003, customer_id=P001, total=225, status=shipped, product_id=PROD003, amount=3, date=2024-02-01, priority=high, discount=15, region=north}
4. [1] TestOrder{id=O004, status=delivered, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, region=east}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, region=south, customer_id=P002, status=confirmed, priority=high, discount=100}
6. [1] TestOrder{id=O006, discount=0, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, priority=low, region=west, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north}
8. [1] TestOrder{id=O008, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, amount=10, status=pending, product_id=PROD002, total=255}
9. [1] TestOrder{id=O009, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, priority=low, total=89.99, date=2024-03-10}
10. [1] TestOrder{id=O010, status=refunded, priority=urgent, region=east, product_id=PROD001, total=75000, discount=0, customer_id=P006, amount=1, date=2024-03-15}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, department=sales, level=2, name=Alice, age=25, salary=45000, active=true, tags=junior, status=active}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, department=engineering, level=5, active=true, score=9.2, tags=senior, status=active}
3. [1] TestPerson{id=P003, salary=0, active=false, score=6, department=hr, name=Charlie, age=16, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active}
5. [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, department=sales, level=3, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, tags=test, status=active, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true, score=0}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, department=management, age=65, salary=95000, tags=executive, status=active, level=9}
8. [1] TestPerson{id=P008, name=Henry, active=false, department=support, level=1, age=18, salary=25000, score=5.5, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, name=Ivy, tags=senior, level=6}
10. [1] TestPerson{id=P010, active=true, score=6.5, status=active, department=intern, level=1, name=X, salary=28000, tags=temp, age=22}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 13 (68.4%)
- **Tokens générés**: 94
- **Faits traités**: 27
