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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5, name=Bob, age=35}
3. [1] TestPerson{id=P003, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern, department=hr}
4. [1] TestPerson{id=P004, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing, age=45, salary=85000, score=7.8}
5. [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales, level=3}
6. [1] TestPerson{id=P006, score=0, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000}
7. [1] TestPerson{id=P007, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7}
10. [1] TestPerson{id=P010, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, level=2, active=true, status=active, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P007, department=management, level=9, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true, status=active, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, active=true, department=marketing, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, department=sales, level=3, age=30, salary=55000, status=inactive}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98, status=pending}
2. [1] TestOrder{id=O002, amount=1, date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed}
3. [1] TestOrder{id=O003, customer_id=P001, status=shipped, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0}
5. [1] TestOrder{id=O005, discount=100, product_id=PROD001, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, status=cancelled, region=west}
7. [1] TestOrder{id=O007, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped}
8. [1] TestOrder{id=O008, region=south, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, product_id=PROD002, amount=10, priority=normal}
9. [1] TestOrder{id=O009, discount=10, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped, product_id=PROD003}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, discount=0, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O008, status=pending, discount=0, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north, amount=1}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O010, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O001, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, status=active, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2}
2. [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5}
3. [1] TestPerson{id=P003, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern, department=hr}
4. [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
5. [1] TestPerson{id=P005, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales, level=3, age=30}
6. [1] TestPerson{id=P006, status=active, department=qa, level=1, salary=-5000, score=0, name=Frank, age=0, active=true, tags=test}
7. [1] TestPerson{id=P007, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace, active=true}
8. [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1, age=18, active=false}
9. [1] TestPerson{id=P009, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active, name=Ivy}
10. [1] TestPerson{id=P010, name=X, age=22, active=true, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales, name=Alice, age=25}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, active=true, department=marketing, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, department=sales, level=3, age=30, salary=55000, status=inactive}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, department=management, level=9, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north}
2. [1] TestOrder{id=O002, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}
4. [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, amount=1, priority=normal, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, status=cancelled, region=west}
7. [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north}
8. [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10}
10. [1] TestOrder{id=O010, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O009, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped, product_id=PROD003, amount=3, total=225}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south, customer_id=P002}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98, status=pending}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0, product_id=PROD004}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, available=true, rating=4.5, keywords=computer, brand=TechCorp, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, stock=50}
2. [1] TestProduct{id=PROD002, name=Mouse, available=true, keywords=peripheral, stock=200, category=accessories, price=25.5, rating=4, brand=TechCorp, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, category=accessories, price=75, keywords=typing, supplier=KeySupply, name=Keyboard, available=false, rating=3.5, brand=KeyTech, stock=0}
4. [1] TestProduct{id=PROD004, name=Monitor, price=299.99, rating=4.8, keywords=display, brand=ScreenPro, stock=30, category=electronics, available=true, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, supplier=OldSupply, name=OldKeyboard, category=accessories, price=8.5, brand=OldTech, stock=0, available=false, rating=2, keywords=obsolete}
6. [1] TestProduct{id=PROD006, name=Headphones, price=150, available=true, rating=4.6, brand=AudioMax, category=audio, keywords=sound, stock=75, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, available=true, stock=25, supplier=CamSupply, price=89.99, rating=3.8, keywords=video, brand=CamTech}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, available=true, rating=4.5, keywords=computer, brand=TechCorp, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, stock=50}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, available=true, keywords=peripheral, stock=200, category=accessories, price=25.5, rating=4, brand=TechCorp, supplier=TechSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, category=accessories, price=75, keywords=typing, supplier=KeySupply, name=Keyboard, available=false, rating=3.5, brand=KeyTech, stock=0}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, price=299.99, rating=4.8, keywords=display, brand=ScreenPro, stock=30, category=electronics, available=true, supplier=ScreenSupply, name=Monitor}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, available=true, rating=4.6, brand=AudioMax, category=audio, keywords=sound, stock=75, supplier=AudioSupply, name=Headphones, price=150}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, keywords=video, brand=CamTech, name=Webcam, category=electronics, available=true, stock=25, supplier=CamSupply, price=89.99, rating=3.8}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, status=active, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2}
2. [1] TestPerson{id=P002, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5, name=Bob}
3. [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1}
4. [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, department=sales, level=3, age=30, salary=55000, status=inactive}
6. [1] TestPerson{id=P006, salary=-5000, score=0, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace, active=true, score=10}
8. [1] TestPerson{id=P008, active=false, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1, age=18}
9. [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6}
10. [1] TestPerson{id=P010, name=X, age=22, active=true, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, level=1, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98, status=pending}
2. [1] TestOrder{id=O002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002}
3. [1] TestOrder{id=O003, customer_id=P001, status=shipped, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north}
4. [1] TestOrder{id=O004, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, amount=1, priority=normal}
5. [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}
6. [1] TestOrder{id=O006, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2}
7. [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north}
8. [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north, amount=1}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98, status=pending}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, status=delivered, discount=0, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, product_id=PROD002}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O005, discount=100, product_id=PROD001, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales}
2. [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5}
3. [1] TestPerson{id=P003, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false}
4. [1] TestPerson{id=P004, active=true, department=marketing, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana}
5. [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, department=sales, level=3, age=30, salary=55000, status=inactive}
6. [1] TestPerson{id=P006, level=1, salary=-5000, score=0, name=Frank, age=0, active=true, tags=test, status=active, department=qa}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9}
8. [1] TestPerson{id=P008, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1}
9. [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6}
10. [1] TestPerson{id=P010, name=X, age=22, active=true, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, department=sales, level=3, age=30, salary=55000, status=inactive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98, status=pending}
2. [1] TestOrder{id=O002, date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}
4. [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, amount=1, priority=normal, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}
6. [1] TestOrder{id=O006, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0}
7. [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north}
8. [1] TestOrder{id=O008, discount=0, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent, region=east, customer_id=P006}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O006, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent, region=east}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5, name=Bob, age=35}
3. [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1}
4. [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
5. [1] TestPerson{id=P005, score=8, tags=employee, department=sales, level=3, age=30, salary=55000, status=inactive, name=Eve, active=false}
6. [1] TestPerson{id=P006, tags=test, status=active, department=qa, level=1, salary=-5000, score=0, name=Frank, age=0, active=true}
7. [1] TestPerson{id=P007, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}
10. [1] TestPerson{id=P010, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales, level=3, age=30}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, status=active, department=qa, level=1, salary=-5000, score=0, name=Frank, age=0, active=true, tags=test}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, available=true, rating=4.5, keywords=computer, brand=TechCorp, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, stock=50}
2. [1] TestProduct{id=PROD002, price=25.5, rating=4, brand=TechCorp, supplier=TechSupply, name=Mouse, available=true, keywords=peripheral, stock=200, category=accessories}
3. [1] TestProduct{id=PROD003, price=75, keywords=typing, supplier=KeySupply, name=Keyboard, available=false, rating=3.5, brand=KeyTech, stock=0, category=accessories}
4. [1] TestProduct{id=PROD004, brand=ScreenPro, stock=30, category=electronics, available=true, supplier=ScreenSupply, name=Monitor, price=299.99, rating=4.8, keywords=display}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, price=8.5, brand=OldTech, stock=0, available=false, rating=2, keywords=obsolete, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, category=audio, keywords=sound, stock=75, supplier=AudioSupply, name=Headphones, price=150, available=true, rating=4.6, brand=AudioMax}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, available=true, stock=25, supplier=CamSupply, price=89.99, rating=3.8, keywords=video, brand=CamTech}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, available=true, rating=4.5, keywords=computer, brand=TechCorp, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, stock=50}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, price=25.5, rating=4, brand=TechCorp, supplier=TechSupply, name=Mouse, available=true, keywords=peripheral, stock=200, category=accessories}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, category=accessories, price=75, keywords=typing, supplier=KeySupply, name=Keyboard, available=false, rating=3.5, brand=KeyTech, stock=0}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, price=299.99, rating=4.8, keywords=display, brand=ScreenPro, stock=30, category=electronics, available=true, supplier=ScreenSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, name=Headphones, price=150, available=true, rating=4.6, brand=AudioMax, category=audio, keywords=sound, stock=75, supplier=AudioSupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, name=Webcam, category=electronics, available=true, stock=25, supplier=CamSupply, price=89.99, rating=3.8, keywords=video, brand=CamTech}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales, name=Alice, age=25}
2. [1] TestPerson{id=P002, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true}
3. [1] TestPerson{id=P003, department=hr, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern}
4. [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
5. [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}
10. [1] TestPerson{id=P010, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, department=support, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, active=true, status=active, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, department=sales, level=3, age=30, salary=55000, status=inactive}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace, active=true, score=10}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, status=shipped, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0}
5. [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}
6. [1] TestOrder{id=O006, priority=low, discount=0, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15}
7. [1] TestOrder{id=O007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010}
9. [1] TestOrder{id=O009, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001}
10. [1] TestOrder{id=O010, discount=0, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0, product_id=PROD004}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O001, region=north, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0, region=south}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, status=active, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2}
2. [1] TestPerson{id=P002, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000}
3. [1] TestPerson{id=P003, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false}
4. [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
5. [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace, active=true, score=10}
8. [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5}
9. [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}
10. [1] TestPerson{id=P010, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true, status=active, level=1}
11. [1] TestOrder{id=O001, discount=50, region=north, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal}
12. [1] TestOrder{id=O002, date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1}
13. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}
14. [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0}
15. [1] TestOrder{id=O005, discount=100, product_id=PROD001, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed}
16. [1] TestOrder{id=O006, priority=low, discount=0, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15}
17. [1] TestOrder{id=O007, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006}
18. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0}
19. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10}
20. [1] TestOrder{id=O010, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, status=active, department=management, level=9, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0}
   - Fait 2: [1] TestOrder{id=O009, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O010, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P006, status=active, department=qa, level=1, salary=-5000, score=0, name=Frank, age=0, active=true, tags=test}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0, region=south}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}
   - Fait 2: [1] TestOrder{id=O002, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing, age=45, salary=85000}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, active=true, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O009, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales, name=Alice, age=25}
   - Fait 2: [1] TestOrder{id=O010, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, level=2, active=true, status=active, department=sales, name=Alice, age=25, salary=45000}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O006, discount=0, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing, age=45, salary=85000, score=7.8}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, status=cancelled, region=west}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O008, status=pending, discount=0, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P003, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P007, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace, active=true}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P001, active=true, status=active, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped, product_id=PROD003}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, amount=1, priority=normal, region=east}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P010, age=22, active=true, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern, name=X}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, name=X, age=22, active=true, status=active, level=1, salary=28000, score=6.5, tags=temp}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, level=3, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P008, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O006, date=2024-02-15, priority=low, discount=0, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, status=cancelled, region=west}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O008, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, product_id=PROD002}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O002, date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, department=sales, level=3, age=30, salary=55000, status=inactive}
   - Fait 2: [1] TestOrder{id=O004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0, product_id=PROD004}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O009, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1, age=18}
   - Fait 2: [1] TestOrder{id=O010, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern}
   - Fait 2: [1] TestOrder{id=O001, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P004, status=active, level=7, name=Diana, active=true, department=marketing, age=45, salary=85000, score=7.8, tags=manager}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P010, active=true, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, status=shipped, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O004, discount=0, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O005, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, status=cancelled, region=west}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales, name=Alice, age=25}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P010, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true}
   - Fait 2: [1] TestOrder{id=O001, discount=50, region=north, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1, age=18, active=false}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, active=true, department=marketing, age=45, salary=85000, score=7.8, tags=manager, status=active}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P003, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1, name=Charlie}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, amount=1, priority=normal, region=east}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O010, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, level=1, salary=-5000, score=0, name=Frank, age=0, active=true, tags=test, status=active}
   - Fait 2: [1] TestOrder{id=O001, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active, name=Ivy, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O001, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P005, score=8, tags=employee, department=sales, level=3, age=30, salary=55000, status=inactive, name=Eve, active=false}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, level=3, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee}
   - Fait 2: [1] TestOrder{id=O008, discount=0, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1, age=18, active=false}
   - Fait 2: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent, region=east}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O003, region=north, customer_id=P001, status=shipped, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P006, level=1, salary=-5000, score=0, name=Frank, age=0, active=true, tags=test, status=active, department=qa}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, amount=1, priority=normal, region=east}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P009, level=6, age=40, salary=68000, score=8.7, status=active, name=Ivy, active=true, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support}
   - Fait 2: [1] TestOrder{id=O007, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O005, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1, age=18, active=false}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support}
   - Fait 2: [1] TestOrder{id=O008, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, product_id=PROD002, amount=10}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales, name=Alice, age=25}
   - Fait 2: [1] TestOrder{id=O002, priority=low, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P006, score=0, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, discount=0, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O010, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south, customer_id=P002, amount=1, total=999.99}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped, product_id=PROD003}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P007, status=active, department=management, level=9, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, tags=temp, department=intern, name=X, age=22, active=true, status=active, level=1, salary=28000}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P001, active=true, status=active, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P003, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16}
   - Fait 2: [1] TestOrder{id=O010, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, active=true, department=marketing, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7}
   - Fait 2: [1] TestOrder{id=O010, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P008, department=support, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P010, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true}
   - Fait 2: [1] TestOrder{id=O002, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}
   - Fait 2: [1] TestOrder{id=O009, discount=10, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O002, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing, age=45, salary=85000, score=7.8}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P008, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped, product_id=PROD003}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, status=cancelled, region=west}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P006, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0, name=Frank, age=0}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}
   - Fait 2: [1] TestOrder{id=O007, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O007, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales}
2. [1] TestPerson{id=P002, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering}
3. [1] TestPerson{id=P003, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern, department=hr}
4. [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
5. [1] TestPerson{id=P005, score=8, tags=employee, department=sales, level=3, age=30, salary=55000, status=inactive, name=Eve, active=false}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
7. [1] TestPerson{id=P007, level=9, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, active=false, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1, age=18}
9. [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6}
10. [1] TestPerson{id=P010, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing, age=45, salary=85000}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, active=true, status=active, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales}
2. [1] TestPerson{id=P002, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5, name=Bob}
3. [1] TestPerson{id=P003, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0}
4. [1] TestPerson{id=P004, name=Diana, active=true, department=marketing, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7}
5. [1] TestPerson{id=P005, department=sales, level=3, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9}
8. [1] TestPerson{id=P008, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000}
9. [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6}
10. [1] TestPerson{id=P010, active=true, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22}
11. [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north}
12. [1] TestOrder{id=O002, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0}
13. [1] TestOrder{id=O003, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped, product_id=PROD003, amount=3}
14. [1] TestOrder{id=O004, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, amount=1, priority=normal}
15. [1] TestOrder{id=O005, discount=100, product_id=PROD001, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed}
16. [1] TestOrder{id=O006, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0}
17. [1] TestOrder{id=O007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007}
18. [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south}
19. [1] TestOrder{id=O009, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007}
20. [1] TestOrder{id=O010, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P006, status=active, department=qa, level=1, salary=-5000, score=0, name=Frank, age=0, active=true, tags=test}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south, customer_id=P002, amount=1}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, department=sales, level=3, age=30, salary=55000, status=inactive}
   - Fait 2: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, age=22, active=true, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern, name=X}
   - Fait 2: [1] TestOrder{id=O006, priority=low, discount=0, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P009, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern}
   - Fait 2: [1] TestOrder{id=O010, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P003, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98, status=pending}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P006, level=1, salary=-5000, score=0, name=Frank, age=0, active=true, tags=test, status=active, department=qa}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P007, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace, active=true}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16, score=6}
   - Fait 2: [1] TestOrder{id=O006, priority=low, discount=0, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, department=intern, name=X, age=22, active=true, status=active, level=1, salary=28000, score=6.5}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent, region=east}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0, region=south}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5}
   - Fait 2: [1] TestOrder{id=O004, discount=0, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P010, age=22, active=true, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern, name=X}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, status=pending, discount=0, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P007, level=9, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P005, level=3, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, amount=1, priority=normal, region=east}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P008, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1}
   - Fait 2: [1] TestOrder{id=O004, discount=0, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south, customer_id=P002, amount=1, total=999.99}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O008, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, product_id=PROD002, amount=10}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5}
   - Fait 2: [1] TestOrder{id=O001, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent, region=east}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O010, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, department=support, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98, status=pending}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98, status=pending}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales, name=Alice, age=25}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, status=shipped, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, level=3, age=30, salary=55000, status=inactive, name=Eve}
   - Fait 2: [1] TestOrder{id=O006, priority=low, discount=0, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support}
   - Fait 2: [1] TestOrder{id=O006, priority=low, discount=0, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent, region=east, customer_id=P006}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P010, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, amount=1, priority=normal, region=east}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O010, discount=0, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O001, region=north, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing, age=45, salary=85000}
   - Fait 2: [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O003, region=north, customer_id=P001, status=shipped, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P003, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, amount=1, priority=normal, region=east}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O006, priority=low, discount=0, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P001, level=2, active=true, status=active, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1, age=18, active=false}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0, region=south}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, level=3, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee}
   - Fait 2: [1] TestOrder{id=O002, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P010, active=true, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22}
   - Fait 2: [1] TestOrder{id=O003, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped, product_id=PROD003, amount=3, total=225}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O008, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, product_id=PROD002, amount=10}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support}
   - Fait 2: [1] TestOrder{id=O009, discount=10, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0, region=south, customer_id=P002}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern}
   - Fait 2: [1] TestOrder{id=O003, region=north, customer_id=P001, status=shipped, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, status=active, department=management, level=9, name=Grace, active=true, score=10, age=65, salary=95000}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P005, level=3, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P006, status=active, department=qa, level=1, salary=-5000, score=0, name=Frank, age=0, active=true, tags=test}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0, region=south}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O005, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, status=cancelled, region=west}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, status=cancelled, region=west, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, active=true, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north, amount=1}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales}
   - Fait 2: [1] TestOrder{id=O004, status=delivered, discount=0, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, department=sales, level=3, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south, customer_id=P002, amount=1, total=999.99}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P008, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P003, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1, name=Charlie}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O010, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent, region=east, customer_id=P006}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, status=shipped, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active, name=Ivy, active=true}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O009, discount=10, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent, region=east}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, level=2, active=true, status=active, department=sales, name=Alice, age=25, salary=45000}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, active=true, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O002, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, level=3, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee}
   - Fait 2: [1] TestOrder{id=O003, priority=high, discount=15, region=north, customer_id=P001, status=shipped, product_id=PROD003, amount=3, total=225, date=2024-02-01}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P008, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped, product_id=PROD003}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P001, level=2, active=true, status=active, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O006, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, status=cancelled, region=west, product_id=PROD005}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P006, level=1, salary=-5000, score=0, name=Frank, age=0, active=true, tags=test, status=active, department=qa}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, total=600, region=north, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending, discount=0}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P010, status=active, level=1, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, discount=0, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, discount=50, region=north, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active}
2. [1] TestPerson{id=P002, status=active, department=engineering, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1}
4. [1] TestPerson{id=P004, department=marketing, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true}
5. [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales, level=3}
6. [1] TestPerson{id=P006, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0, name=Frank}
7. [1] TestPerson{id=P007, tags=executive, status=active, department=management, level=9, name=Grace, active=true, score=10, age=65, salary=95000}
8. [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5}
9. [1] TestPerson{id=P009, status=active, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7}
10. [1] TestPerson{id=P010, score=6.5, tags=temp, department=intern, name=X, age=22, active=true, status=active, level=1, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, department=support, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, name=X, age=22, active=true, status=active, level=1, salary=28000, score=6.5, tags=temp}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, active=true, status=active, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, active=true, department=marketing, age=45, salary=85000, score=7.8, tags=manager, status=active, level=7}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, level=3, age=30, salary=55000, status=inactive, name=Eve, active=false, score=8, tags=employee, department=sales}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, level=9, name=Grace, active=true, score=10, age=65, salary=95000, tags=executive, status=active, department=management}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true, level=5}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low, discount=0, region=south}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0}
5. [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001, priority=high, region=south}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, status=cancelled, region=west}
7. [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north}
8. [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, region=north, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0, total=75000, priority=urgent, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, product_id=PROD002, total=25.5, status=confirmed, amount=1, date=2024-01-20, priority=low}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, discount=15, region=north, customer_id=P001, status=shipped}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, priority=urgent, discount=50, customer_id=P007, product_id=PROD006, total=600, region=north, amount=4}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, status=pending, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O004, region=east, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, amount=1, priority=normal}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O005, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, status=confirmed, discount=100, product_id=PROD001}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, status=cancelled, region=west}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O008, discount=0, customer_id=P010, product_id=PROD002, amount=10, priority=normal, region=south, total=255, date=2024-03-05, status=pending}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, customer_id=P001, product_id=PROD007, region=north}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, total=75000, priority=urgent, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales}
2. [1] TestPerson{id=P002, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering, salary=75000, active=true}
3. [1] TestPerson{id=P003, name=Charlie, age=16, score=6, tags=intern, department=hr, salary=0, active=false, status=inactive, level=1}
4. [1] TestPerson{id=P004, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing, age=45, salary=85000, score=7.8}
5. [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, level=3, age=30, salary=55000, status=inactive, name=Eve}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace, active=true, score=10}
8. [1] TestPerson{id=P008, active=false, score=5.5, tags=junior, status=inactive, department=support, name=Henry, salary=25000, level=1, age=18}
9. [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}
10. [1] TestPerson{id=P010, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, status=active, department=management, level=9, name=Grace, active=true, score=10}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, tags=senior, department=engineering, level=6, age=40, salary=68000, score=8.7, status=active}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, tags=temp, department=intern, name=X, age=22, active=true, status=active, level=1}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, level=5, name=Bob, age=35, score=9.2, tags=senior, status=active, department=engineering}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, status=inactive, level=1, name=Charlie, age=16, score=6, tags=intern, department=hr}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, level=3, age=30, salary=55000, status=inactive, name=Eve}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, level=1, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, level=2, active=true, status=active, department=sales}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, level=7, name=Diana, active=true, department=marketing, age=45, salary=85000}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P006, age=0, active=true, tags=test, status=active, department=qa, level=1, salary=-5000, score=0, name=Frank}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
