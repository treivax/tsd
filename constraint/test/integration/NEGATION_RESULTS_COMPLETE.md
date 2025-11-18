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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, department=sales, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, salary=75000, active=true, status=active, name=Bob, age=35, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, score=6, tags=intern, age=16, active=false, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, age=45, score=7.8, department=marketing, level=7, name=Diana}
5. [1] TestPerson{id=P005, status=inactive, age=30, active=false, department=sales, level=3, name=Eve, salary=55000, score=8, tags=employee}
6. [1] TestPerson{id=P006, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, status=active, department=management, salary=95000, active=true, tags=executive, level=9}
8. [1] TestPerson{id=P008, department=support, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, tags=senior, level=6}
10. [1] TestPerson{id=P010, salary=28000, active=true, score=6.5, tags=temp, status=active, name=X, age=22, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, status=active, name=Bob, age=35, score=9.2, tags=senior, department=engineering, level=5}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6, tags=intern}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, age=45, score=7.8, department=marketing, level=7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, tags=employee, status=inactive, age=30, active=false, department=sales, level=3}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, tags=senior, level=6, name=Ivy}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, age=22, department=intern, level=1, salary=28000, active=true, score=6.5, tags=temp, status=active, name=X}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, status=active, level=2, name=Alice, age=25, salary=45000, active=true, department=sales}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, score=10, status=active, department=management, salary=95000, active=true, tags=executive, level=9}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, age=18, active=false, score=5.5, status=inactive, department=support, salary=25000, tags=junior, level=1, name=Henry}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, region=north, customer_id=P001, total=1999.98, date=2024-01-15, discount=50}
2. [1] TestOrder{id=O002, discount=0, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002, amount=1}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, discount=15, region=north, date=2024-02-01, status=shipped, priority=high}
4. [1] TestOrder{id=O004, status=delivered, priority=normal, product_id=PROD004, amount=1, total=299.99, discount=0, region=east, customer_id=P004, date=2024-02-05}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, region=south, date=2024-02-10, status=confirmed, priority=high, discount=100}
6. [1] TestOrder{id=O006, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, amount=2}
7. [1] TestOrder{id=O007, priority=urgent, discount=50, region=north, date=2024-03-01, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, discount=0}
9. [1] TestOrder{id=O009, region=north, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, priority=low, total=89.99, status=completed, discount=10}
10. [1] TestOrder{id=O010, product_id=PROD001, total=75000, status=refunded, region=east, amount=1, date=2024-03-15, priority=urgent, discount=0, customer_id=P006}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, total=75000, status=refunded, region=east, amount=1, date=2024-03-15, priority=urgent, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, discount=0, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, discount=15, region=north, date=2024-02-01, status=shipped, priority=high}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, region=south, date=2024-02-10, status=confirmed, priority=high, discount=100}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, discount=50, region=north, date=2024-03-01, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, total=255, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, discount=0, product_id=PROD002, amount=10}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, priority=low, total=89.99}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O001, priority=normal, region=north, customer_id=P001, total=1999.98, date=2024-01-15, discount=50, product_id=PROD001, amount=2, status=pending}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O004, region=east, customer_id=P004, date=2024-02-05, status=delivered, priority=normal, product_id=PROD004, amount=1, total=299.99, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, department=sales, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, department=engineering, level=5, salary=75000, active=true, status=active, name=Bob, age=35, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, department=hr, level=1, name=Charlie, salary=0, score=6, tags=intern, age=16, active=false, status=inactive}
4. [1] TestPerson{id=P004, age=45, score=7.8, department=marketing, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, tags=employee, status=inactive, age=30, active=false, department=sales, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, salary=95000, active=true, tags=executive, level=9, name=Grace, age=65, score=10, status=active, department=management}
8. [1] TestPerson{id=P008, score=5.5, status=inactive, department=support, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false}
9. [1] TestPerson{id=P009, status=active, department=engineering, tags=senior, level=6, name=Ivy, age=40, salary=68000, active=true, score=8.7}
10. [1] TestPerson{id=P010, level=1, salary=28000, active=true, score=6.5, tags=temp, status=active, name=X, age=22, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, age=45, score=7.8, department=marketing, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, age=30, active=false, department=sales, level=3, name=Eve, salary=55000, score=8, tags=employee}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, score=10, status=active, department=management, salary=95000, active=true, tags=executive, level=9}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, tags=senior, level=6}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, department=sales, score=8.5, tags=junior, status=active, level=2}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, status=active, name=Bob, age=35, score=9.2, tags=senior, department=engineering, level=5, salary=75000, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, region=north, customer_id=P001, total=1999.98, date=2024-01-15, discount=50}
2. [1] TestOrder{id=O002, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002, amount=1, discount=0, total=25.5, date=2024-01-20}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, discount=15, region=north, date=2024-02-01, status=shipped, priority=high}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, discount=0, region=east, customer_id=P004, date=2024-02-05, status=delivered, priority=normal}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, region=south, date=2024-02-10, status=confirmed, priority=high, discount=100}
6. [1] TestOrder{id=O006, priority=low, discount=0, region=west, product_id=PROD005, amount=2, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled}
7. [1] TestOrder{id=O007, priority=urgent, discount=50, region=north, date=2024-03-01, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped}
8. [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, discount=0, product_id=PROD002}
9. [1] TestOrder{id=O009, total=89.99, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, priority=low}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, total=75000, status=refunded, region=east, amount=1, date=2024-03-15, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, discount=15, region=north, date=2024-02-01, status=shipped, priority=high}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O006, priority=low, discount=0, region=west, product_id=PROD005, amount=2, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, discount=0, product_id=PROD002}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, priority=low, total=89.99, status=completed, discount=10, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O001, region=north, customer_id=P001, total=1999.98, date=2024-01-15, discount=50, product_id=PROD001, amount=2, status=pending, priority=normal}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O002, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002, amount=1, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, discount=0, region=east, customer_id=P004, date=2024-02-05, status=delivered, priority=normal}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, total=999.99, region=south, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, amount=1}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, region=north, date=2024-03-01}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, price=999.99, available=true, rating=4.5, keywords=computer, stock=50, supplier=TechSupply, name=Laptop, brand=TechCorp, category=electronics}
2. [1] TestProduct{id=PROD002, category=accessories, price=25.5, brand=TechCorp, stock=200, name=Mouse, available=true, rating=4, keywords=peripheral, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, category=accessories, price=75, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0, name=Keyboard, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, price=299.99, available=true, rating=4.8, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics, keywords=display, brand=ScreenPro}
5. [1] TestProduct{id=PROD005, category=accessories, price=8.5, available=false, keywords=obsolete, brand=OldTech, stock=0, name=OldKeyboard, rating=2, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, stock=75, supplier=AudioSupply, name=Headphones, category=audio, price=150, available=true, rating=4.6, keywords=sound, brand=AudioMax}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, price=89.99, available=true, rating=3.8, keywords=video, brand=CamTech, supplier=CamSupply, stock=25}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, stock=50, supplier=TechSupply, name=Laptop, brand=TechCorp, category=electronics, price=999.99, available=true, rating=4.5, keywords=computer}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, stock=200, name=Mouse, available=true, rating=4, keywords=peripheral, supplier=TechSupply, category=accessories, price=25.5, brand=TechCorp}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, category=accessories, price=75, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0, name=Keyboard, supplier=KeySupply}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, category=electronics, keywords=display, brand=ScreenPro, price=299.99, available=true, rating=4.8, stock=30, supplier=ScreenSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, name=Headphones, category=audio, price=150, available=true, rating=4.6, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, name=Webcam, category=electronics, price=89.99, available=true, rating=3.8, keywords=video, brand=CamTech, supplier=CamSupply, stock=25}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, department=sales, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, department=engineering, level=5, salary=75000, active=true, status=active}
3. [1] TestPerson{id=P003, level=1, name=Charlie, salary=0, score=6, tags=intern, age=16, active=false, status=inactive, department=hr}
4. [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, age=45, score=7.8, department=marketing, level=7, name=Diana}
5. [1] TestPerson{id=P005, age=30, active=false, department=sales, level=3, name=Eve, salary=55000, score=8, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
7. [1] TestPerson{id=P007, level=9, name=Grace, age=65, score=10, status=active, department=management, salary=95000, active=true, tags=executive}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, status=inactive, department=support, salary=25000, tags=junior, level=1}
9. [1] TestPerson{id=P009, tags=senior, level=6, name=Ivy, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering}
10. [1] TestPerson{id=P010, salary=28000, active=true, score=6.5, tags=temp, status=active, name=X, age=22, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6, tags=intern}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, discount=50, product_id=PROD001, amount=2, status=pending, priority=normal, region=north}
2. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002, amount=1, discount=0}
3. [1] TestOrder{id=O003, region=north, date=2024-02-01, status=shipped, priority=high, customer_id=P001, product_id=PROD003, amount=3, total=225, discount=15}
4. [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, status=delivered, priority=normal, product_id=PROD004, amount=1, total=299.99, discount=0, region=east}
5. [1] TestOrder{id=O005, region=south, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, amount=1, total=999.99}
6. [1] TestOrder{id=O006, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, amount=2}
7. [1] TestOrder{id=O007, status=shipped, priority=urgent, discount=50, region=north, date=2024-03-01, customer_id=P007, product_id=PROD006, amount=4, total=600}
8. [1] TestOrder{id=O008, customer_id=P010, status=pending, discount=0, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, priority=low, total=89.99, status=completed, discount=10, region=north}
10. [1] TestOrder{id=O010, amount=1, date=2024-03-15, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, total=75000, status=refunded, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, discount=0, region=east, customer_id=P004, date=2024-02-05, status=delivered, priority=normal}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, region=south, date=2024-02-10, status=confirmed, priority=high, discount=100}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O006, region=west, product_id=PROD005, amount=2, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, region=north, date=2024-03-01, customer_id=P007}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, status=pending, discount=0, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, priority=normal, region=north, customer_id=P001, total=1999.98, date=2024-01-15, discount=50, product_id=PROD001, amount=2, status=pending}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002, amount=1, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, discount=15, region=north, date=2024-02-01, status=shipped, priority=high}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, date=2024-03-10, priority=low, total=89.99, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, total=75000, status=refunded, region=east, amount=1, date=2024-03-15}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, department=sales, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, department=engineering, level=5, salary=75000, active=true, status=active}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, score=6, tags=intern, age=16, active=false, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, age=45, score=7.8, department=marketing, level=7, name=Diana}
5. [1] TestPerson{id=P005, tags=employee, status=inactive, age=30, active=false, department=sales, level=3, name=Eve, salary=55000, score=8}
6. [1] TestPerson{id=P006, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, status=active, department=management, salary=95000, active=true, tags=executive, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, status=inactive, department=support, salary=25000, tags=junior, level=1}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, tags=senior, level=6}
10. [1] TestPerson{id=P010, department=intern, level=1, salary=28000, active=true, score=6.5, tags=temp, status=active, name=X, age=22}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6, tags=intern}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, tags=employee, status=inactive, age=30, active=false, department=sales, level=3}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, status=inactive, department=support, salary=25000, tags=junior, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, priority=normal, region=north, customer_id=P001, total=1999.98, date=2024-01-15, discount=50, product_id=PROD001, amount=2, status=pending}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, discount=0, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south}
3. [1] TestOrder{id=O003, amount=3, total=225, discount=15, region=north, date=2024-02-01, status=shipped, priority=high, customer_id=P001, product_id=PROD003}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, discount=0, region=east, customer_id=P004, date=2024-02-05, status=delivered, priority=normal}
5. [1] TestOrder{id=O005, region=south, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, amount=1, total=999.99}
6. [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, amount=2, customer_id=P005, total=999.98, date=2024-02-15}
7. [1] TestOrder{id=O007, priority=urgent, discount=50, region=north, date=2024-03-01, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, discount=0}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, priority=low, total=89.99, status=completed, discount=10, region=north}
10. [1] TestOrder{id=O010, product_id=PROD001, total=75000, status=refunded, region=east, amount=1, date=2024-03-15, priority=urgent, discount=0, customer_id=P006}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, priority=normal, product_id=PROD004, amount=1, total=299.99, discount=0, region=east, customer_id=P004, date=2024-02-05, status=delivered}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, region=south, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, amount=1, total=999.99}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O006, amount=2, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, discount=50, region=north, date=2024-03-01, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, priority=low, total=89.99, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O010, date=2024-03-15, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, total=75000, status=refunded, region=east, amount=1}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, discount=0, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, discount=15, region=north, date=2024-02-01, status=shipped, priority=high}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, department=sales, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, salary=75000, active=true, status=active, name=Bob, age=35, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, salary=0, score=6, tags=intern, age=16, active=false, status=inactive, department=hr, level=1, name=Charlie}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, age=45, score=7.8, department=marketing, level=7}
5. [1] TestPerson{id=P005, age=30, active=false, department=sales, level=3, name=Eve, salary=55000, score=8, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, status=active, department=management, salary=95000, active=true, tags=executive, level=9}
8. [1] TestPerson{id=P008, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false, score=5.5, status=inactive, department=support}
9. [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, status=active, department=engineering, tags=senior, level=6, name=Ivy, age=40}
10. [1] TestPerson{id=P010, salary=28000, active=true, score=6.5, tags=temp, status=active, name=X, age=22, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6, tags=intern}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, age=30, active=false, department=sales, level=3, name=Eve, salary=55000, score=8, tags=employee}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, score=8.5, tags=junior, status=active, level=2, name=Alice, age=25, salary=45000, active=true}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, department=engineering, level=5, salary=75000, active=true, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, age=45, score=7.8, department=marketing, level=7}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, score=10, status=active, department=management, salary=95000, active=true, tags=executive, level=9}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false, score=5.5, status=inactive, department=support}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, level=6, name=Ivy, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, tags=senior}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, rating=4.5, keywords=computer, stock=50, supplier=TechSupply, name=Laptop, brand=TechCorp, category=electronics, price=999.99, available=true}
2. [1] TestProduct{id=PROD002, name=Mouse, available=true, rating=4, keywords=peripheral, supplier=TechSupply, category=accessories, price=25.5, brand=TechCorp, stock=200}
3. [1] TestProduct{id=PROD003, name=Keyboard, supplier=KeySupply, category=accessories, price=75, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0}
4. [1] TestProduct{id=PROD004, name=Monitor, category=electronics, keywords=display, brand=ScreenPro, price=299.99, available=true, rating=4.8, stock=30, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, category=accessories, price=8.5, available=false, keywords=obsolete, brand=OldTech, stock=0, name=OldKeyboard, rating=2, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, name=Headphones, category=audio, price=150, available=true, rating=4.6, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, brand=CamTech, supplier=CamSupply, stock=25, name=Webcam, category=electronics, price=89.99, available=true, rating=3.8, keywords=video}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD007, rating=3.8, keywords=video, brand=CamTech, supplier=CamSupply, stock=25, name=Webcam, category=electronics, price=89.99, available=true}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, brand=TechCorp, category=electronics, price=999.99, available=true, rating=4.5, keywords=computer, stock=50, supplier=TechSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD002, supplier=TechSupply, category=accessories, price=25.5, brand=TechCorp, stock=200, name=Mouse, available=true, rating=4, keywords=peripheral}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD003, keywords=typing, brand=KeyTech, stock=0, name=Keyboard, supplier=KeySupply, category=accessories, price=75, available=false, rating=3.5}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, category=electronics, keywords=display, brand=ScreenPro, price=299.99, available=true, rating=4.8, stock=30, supplier=ScreenSupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD006, brand=AudioMax, stock=75, supplier=AudioSupply, name=Headphones, category=audio, price=150, available=true, rating=4.6, keywords=sound}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, department=sales, score=8.5, tags=junior, status=active, level=2, name=Alice, age=25, salary=45000}
2. [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, department=engineering, level=5, salary=75000, active=true, status=active}
3. [1] TestPerson{id=P003, age=16, active=false, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6, tags=intern}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, age=45, score=7.8, department=marketing, level=7}
5. [1] TestPerson{id=P005, score=8, tags=employee, status=inactive, age=30, active=false, department=sales, level=3, name=Eve, salary=55000}
6. [1] TestPerson{id=P006, salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
7. [1] TestPerson{id=P007, department=management, salary=95000, active=true, tags=executive, level=9, name=Grace, age=65, score=10, status=active}
8. [1] TestPerson{id=P008, status=inactive, department=support, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false, score=5.5}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, tags=senior, level=6}
10. [1] TestPerson{id=P010, score=6.5, tags=temp, status=active, name=X, age=22, department=intern, level=1, salary=28000, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, status=active, level=2, name=Alice, age=25, salary=45000, active=true, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, level=5, salary=75000, active=true, status=active, name=Bob, age=35, score=9.2, tags=senior}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, age=45, score=7.8, department=marketing, level=7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, department=sales, level=3, name=Eve, salary=55000, score=8, tags=employee, status=inactive}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, status=active, department=management, salary=95000, active=true, tags=executive, level=9, name=Grace, age=65, score=10}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false, score=5.5, status=inactive, department=support}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, salary=0, score=6, tags=intern, age=16, active=false, status=inactive, department=hr, level=1}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, tags=senior, level=6}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, status=pending, priority=normal, region=north, customer_id=P001, total=1999.98, date=2024-01-15, discount=50}
2. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002, amount=1, discount=0}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, discount=15, region=north, date=2024-02-01, status=shipped, priority=high, customer_id=P001}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, discount=0, region=east, customer_id=P004, date=2024-02-05, status=delivered, priority=normal}
5. [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, amount=1, total=999.99, region=south}
6. [1] TestOrder{id=O006, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, amount=2}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, region=north, date=2024-03-01}
8. [1] TestOrder{id=O008, status=pending, discount=0, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, region=south, customer_id=P010}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, priority=low, total=89.99, status=completed, discount=10, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, total=75000, status=refunded, region=east, amount=1, date=2024-03-15, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, discount=0, region=east, customer_id=P004, date=2024-02-05, status=delivered, priority=normal}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, region=south, date=2024-02-10, status=confirmed, priority=high, discount=100}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, region=north, date=2024-03-01}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O008, status=pending, discount=0, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, region=south, customer_id=P010}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, amount=1, date=2024-03-10, priority=low, total=89.99, status=completed, discount=10, region=north, customer_id=P001}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, discount=50, product_id=PROD001, amount=2, status=pending, priority=normal, region=north, customer_id=P001, total=1999.98, date=2024-01-15}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, discount=0, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, discount=15, region=north, date=2024-02-01, status=shipped, priority=high, customer_id=P001, product_id=PROD003, amount=3, total=225}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, score=8.5, tags=junior, status=active, level=2, name=Alice, age=25, salary=45000, active=true}
2. [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, department=engineering, level=5, salary=75000, active=true, status=active}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, score=6, tags=intern, age=16, active=false, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, active=true, tags=manager, status=active, age=45, score=7.8, department=marketing, level=7, name=Diana, salary=85000}
5. [1] TestPerson{id=P005, status=inactive, age=30, active=false, department=sales, level=3, name=Eve, salary=55000, score=8, tags=employee}
6. [1] TestPerson{id=P006, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, department=management, salary=95000, active=true, tags=executive, level=9, name=Grace, age=65, score=10, status=active}
8. [1] TestPerson{id=P008, active=false, score=5.5, status=inactive, department=support, salary=25000, tags=junior, level=1, name=Henry, age=18}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, tags=senior, level=6}
10. [1] TestPerson{id=P010, age=22, department=intern, level=1, salary=28000, active=true, score=6.5, tags=temp, status=active, name=X}
11. [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, discount=50, product_id=PROD001, amount=2, status=pending, priority=normal, region=north, customer_id=P001}
12. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002, amount=1, discount=0}
13. [1] TestOrder{id=O003, status=shipped, priority=high, customer_id=P001, product_id=PROD003, amount=3, total=225, discount=15, region=north, date=2024-02-01}
14. [1] TestOrder{id=O004, priority=normal, product_id=PROD004, amount=1, total=299.99, discount=0, region=east, customer_id=P004, date=2024-02-05, status=delivered}
15. [1] TestOrder{id=O005, status=confirmed, priority=high, discount=100, customer_id=P002, product_id=PROD001, amount=1, total=999.99, region=south, date=2024-02-10}
16. [1] TestOrder{id=O006, amount=2, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005}
17. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, region=north, date=2024-03-01}
18. [1] TestOrder{id=O008, total=255, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, discount=0, product_id=PROD002, amount=10}
19. [1] TestOrder{id=O009, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, priority=low, total=89.99, status=completed}
20. [1] TestOrder{id=O010, region=east, amount=1, date=2024-03-15, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, total=75000, status=refunded}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, department=sales, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, department=engineering, level=5, salary=75000, active=true, status=active, name=Bob, age=35, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, score=6, tags=intern, age=16, active=false, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, age=45, score=7.8, department=marketing, level=7}
5. [1] TestPerson{id=P005, level=3, name=Eve, salary=55000, score=8, tags=employee, status=inactive, age=30, active=false, department=sales}
6. [1] TestPerson{id=P006, salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
7. [1] TestPerson{id=P007, active=true, tags=executive, level=9, name=Grace, age=65, score=10, status=active, department=management, salary=95000}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, status=inactive, department=support, salary=25000, tags=junior, level=1}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, tags=senior, level=6}
10. [1] TestPerson{id=P010, salary=28000, active=true, score=6.5, tags=temp, status=active, name=X, age=22, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, department=sales, score=8.5, tags=junior, status=active, level=2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, department=engineering, level=5, salary=75000, active=true, status=active}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, tags=manager, status=active, age=45, score=7.8, department=marketing, level=7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, score=10, status=active, department=management, salary=95000, active=true, tags=executive, level=9, name=Grace, age=65}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, tags=senior, level=6}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, tags=temp, status=active, name=X, age=22, department=intern, level=1, salary=28000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, tags=junior, status=active, level=2, name=Alice, age=25, salary=45000, active=true, department=sales}
2. [1] TestPerson{id=P002, salary=75000, active=true, status=active, name=Bob, age=35, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, tags=intern, age=16, active=false, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6}
4. [1] TestPerson{id=P004, age=45, score=7.8, department=marketing, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, tags=employee, status=inactive, age=30, active=false, department=sales, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, status=active, department=management, salary=95000, active=true, tags=executive, level=9}
8. [1] TestPerson{id=P008, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false, score=5.5, status=inactive, department=support}
9. [1] TestPerson{id=P009, active=true, score=8.7, status=active, department=engineering, tags=senior, level=6, name=Ivy, age=40, salary=68000}
10. [1] TestPerson{id=P010, status=active, name=X, age=22, department=intern, level=1, salary=28000, active=true, score=6.5, tags=temp}
11. [1] TestOrder{id=O001, status=pending, priority=normal, region=north, customer_id=P001, total=1999.98, date=2024-01-15, discount=50, product_id=PROD001, amount=2}
12. [1] TestOrder{id=O002, amount=1, discount=0, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002}
13. [1] TestOrder{id=O003, date=2024-02-01, status=shipped, priority=high, customer_id=P001, product_id=PROD003, amount=3, total=225, discount=15, region=north}
14. [1] TestOrder{id=O004, total=299.99, discount=0, region=east, customer_id=P004, date=2024-02-05, status=delivered, priority=normal, product_id=PROD004, amount=1}
15. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, region=south, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002}
16. [1] TestOrder{id=O006, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, amount=2}
17. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, region=north, date=2024-03-01}
18. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, discount=0}
19. [1] TestOrder{id=O009, total=89.99, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, priority=low}
20. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, total=75000, status=refunded, region=east, amount=1, date=2024-03-15, priority=urgent, discount=0}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, department=sales, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, department=engineering, level=5, salary=75000, active=true, status=active}
3. [1] TestPerson{id=P003, age=16, active=false, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6, tags=intern}
4. [1] TestPerson{id=P004, age=45, score=7.8, department=marketing, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, score=8, tags=employee, status=inactive, age=30, active=false, department=sales, level=3, name=Eve, salary=55000}
6. [1] TestPerson{id=P006, salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, status=active, department=management, salary=95000, active=true, tags=executive, level=9}
8. [1] TestPerson{id=P008, score=5.5, status=inactive, department=support, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false}
9. [1] TestPerson{id=P009, active=true, score=8.7, status=active, department=engineering, tags=senior, level=6, name=Ivy, age=40, salary=68000}
10. [1] TestPerson{id=P010, salary=28000, active=true, score=6.5, tags=temp, status=active, name=X, age=22, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, department=sales, level=3, name=Eve, salary=55000, score=8, tags=employee, status=inactive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, status=active, level=2, name=Alice, age=25, salary=45000, active=true, department=sales, score=8.5}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, score=9.2, tags=senior, department=engineering, level=5, salary=75000, active=true, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6, tags=intern}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, tags=manager, status=active, age=45, score=7.8, department=marketing, level=7, name=Diana}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, score=10, status=active, department=management, salary=95000, active=true, tags=executive, level=9}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P008, department=support, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false, score=5.5, status=inactive}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, tags=senior, level=6}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, level=1, salary=28000, active=true, score=6.5, tags=temp, status=active, name=X, age=22, department=intern}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, amount=2, status=pending, priority=normal, region=north, customer_id=P001, total=1999.98, date=2024-01-15, discount=50, product_id=PROD001}
2. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002, amount=1, discount=0}
3. [1] TestOrder{id=O003, date=2024-02-01, status=shipped, priority=high, customer_id=P001, product_id=PROD003, amount=3, total=225, discount=15, region=north}
4. [1] TestOrder{id=O004, total=299.99, discount=0, region=east, customer_id=P004, date=2024-02-05, status=delivered, priority=normal, product_id=PROD004, amount=1}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, region=south, date=2024-02-10, status=confirmed, priority=high, discount=100}
6. [1] TestOrder{id=O006, priority=low, discount=0, region=west, product_id=PROD005, amount=2, customer_id=P005, total=999.98, date=2024-02-15, status=cancelled}
7. [1] TestOrder{id=O007, discount=50, region=north, date=2024-03-01, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}
8. [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, priority=normal, region=south, customer_id=P010, status=pending, discount=0, product_id=PROD002}
9. [1] TestOrder{id=O009, total=89.99, status=completed, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, priority=low}
10. [1] TestOrder{id=O010, product_id=PROD001, total=75000, status=refunded, region=east, amount=1, date=2024-03-15, priority=urgent, discount=0, customer_id=P006}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, discount=50, product_id=PROD001, amount=2, status=pending, priority=normal, region=north, customer_id=P001, total=1999.98, date=2024-01-15}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, priority=high, customer_id=P001, product_id=PROD003, amount=3, total=225, discount=15, region=north, date=2024-02-01, status=shipped}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, total=299.99, discount=0, region=east, customer_id=P004, date=2024-02-05, status=delivered, priority=normal, product_id=PROD004}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, total=600, status=shipped, priority=urgent, discount=50, region=north, date=2024-03-01, customer_id=P007, product_id=PROD006}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, status=pending, discount=0, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, date=2024-03-10, priority=low, total=89.99, status=completed, discount=10, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O010, priority=urgent, discount=0, customer_id=P006, product_id=PROD001, total=75000, status=refunded, region=east, amount=1, date=2024-03-15}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, priority=low, region=south, customer_id=P002, product_id=PROD002, amount=1, discount=0, total=25.5, date=2024-01-20, status=confirmed}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, region=south, date=2024-02-10, status=confirmed, priority=high, discount=100, customer_id=P002}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, region=west, product_id=PROD005, amount=2, customer_id=P005, total=999.98, date=2024-02-15}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, department=sales, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, status=active, name=Bob, age=35, score=9.2, tags=senior, department=engineering, level=5, salary=75000, active=true}
3. [1] TestPerson{id=P003, age=16, active=false, status=inactive, department=hr, level=1, name=Charlie, salary=0, score=6, tags=intern}
4. [1] TestPerson{id=P004, age=45, score=7.8, department=marketing, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active}
5. [1] TestPerson{id=P005, active=false, department=sales, level=3, name=Eve, salary=55000, score=8, tags=employee, status=inactive, age=30}
6. [1] TestPerson{id=P006, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, status=active, department=management, salary=95000, active=true, tags=executive, level=9}
8. [1] TestPerson{id=P008, salary=25000, tags=junior, level=1, name=Henry, age=18, active=false, score=5.5, status=inactive, department=support}
9. [1] TestPerson{id=P009, department=engineering, tags=senior, level=6, name=Ivy, age=40, salary=68000, active=true, score=8.7, status=active}
10. [1] TestPerson{id=P010, age=22, department=intern, level=1, salary=28000, active=true, score=6.5, tags=temp, status=active, name=X}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, active=true, department=sales, score=8.5, tags=junior, status=active, level=2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, score=6, tags=intern, age=16, active=false, status=inactive, department=hr, level=1, name=Charlie, salary=0}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, department=marketing, level=7, name=Diana, salary=85000, active=true, tags=manager, status=active, age=45}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, status=active, department=engineering, tags=senior, level=6}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, status=active, name=Bob, age=35, score=9.2, tags=senior, department=engineering, level=5, salary=75000, active=true}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, tags=employee, status=inactive, age=30, active=false, department=sales, level=3}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, level=9, name=Grace, age=65, score=10, status=active, department=management, salary=95000, active=true}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, status=inactive, department=support, salary=25000, tags=junior, level=1}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, status=active, name=X, age=22, department=intern, level=1, salary=28000, active=true, score=6.5}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 17 (89.5%)
- **Tokens générés**: 130
- **Faits traités**: 27
