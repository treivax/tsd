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

1. [1] TestPerson{id=P001, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active}
2. [1] TestPerson{id=P002, level=5, name=Bob, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering}
3. [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, age=45, active=true, status=active, department=marketing, level=7, name=Diana, salary=85000, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, age=30, status=inactive, level=3, salary=55000, active=false, score=8, tags=employee, department=sales}
6. [1] TestPerson{id=P006, tags=test, status=active, department=qa, active=true, level=1, name=Frank, age=0, salary=-5000, score=0}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}
8. [1] TestPerson{id=P008, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000}
9. [1] TestPerson{id=P009, salary=68000, score=8.7, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior, status=active}
10. [1] TestPerson{id=P010, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1, age=22, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management, name=Grace}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false, score=5.5}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, level=3, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve, age=30}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P010, name=X, active=true, level=1, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, department=hr, level=1, name=Charlie, age=16, active=false, score=6, status=inactive}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7}

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
2. [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south}
3. [1] TestOrder{id=O003, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, customer_id=P001}
4. [1] TestOrder{id=O004, total=299.99, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004}
5. [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high}
6. [1] TestOrder{id=O006, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}
8. [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0}
9. [1] TestOrder{id=O009, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed}
10. [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O008, discount=0, total=255, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O009, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, region=south, product_id=PROD001, amount=1, total=999.99, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, score=8.5, status=active, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35}
3. [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3, salary=55000}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}
9. [1] TestPerson{id=P009, salary=68000, score=8.7, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior, status=active}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, age=30, status=inactive, level=3, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6, name=Ivy}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, status=active, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true, score=9.2}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}
2. [1] TestOrder{id=O002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002}
3. [1] TestOrder{id=O003, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, customer_id=P001}
4. [1] TestOrder{id=O004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004}
5. [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled, region=west}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}
8. [1] TestOrder{id=O008, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1}
10. [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, total=999.99, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O008, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O006, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, stock=50, name=Laptop, category=electronics, price=999.99, available=true, rating=4.5, keywords=computer, brand=TechCorp, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, category=accessories, available=true, supplier=TechSupply, price=25.5, rating=4, keywords=peripheral, brand=TechCorp, stock=200, name=Mouse}
3. [1] TestProduct{id=PROD003, available=false, name=Keyboard, category=accessories, rating=3.5, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply, price=75}
4. [1] TestProduct{id=PROD004, category=electronics, brand=ScreenPro, price=299.99, available=true, rating=4.8, keywords=display, stock=30, supplier=ScreenSupply, name=Monitor}
5. [1] TestProduct{id=PROD005, keywords=obsolete, stock=0, name=OldKeyboard, category=accessories, rating=2, brand=OldTech, supplier=OldSupply, price=8.5, available=false}
6. [1] TestProduct{id=PROD006, available=true, rating=4.6, brand=AudioMax, name=Headphones, category=audio, price=150, keywords=sound, stock=75, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, rating=3.8, keywords=video, brand=CamTech, stock=25, supplier=CamSupply, price=89.99, available=true}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD002, available=true, supplier=TechSupply, price=25.5, rating=4, keywords=peripheral, brand=TechCorp, stock=200, name=Mouse, category=accessories}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD003, price=75, available=false, name=Keyboard, category=accessories, rating=3.5, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD004, rating=4.8, keywords=display, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics, brand=ScreenPro, price=299.99, available=true}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD006, keywords=sound, stock=75, supplier=AudioSupply, available=true, rating=4.6, brand=AudioMax, name=Headphones, category=audio, price=150}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD007, name=Webcam, category=electronics, rating=3.8, keywords=video, brand=CamTech, stock=25, supplier=CamSupply, price=89.99, available=true}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD001, stock=50, name=Laptop, category=electronics, price=999.99, available=true, rating=4.5, keywords=computer, brand=TechCorp, supplier=TechSupply}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35}
3. [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, department=marketing, level=7, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active}
5. [1] TestPerson{id=P005, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3, salary=55000, active=false}
6. [1] TestPerson{id=P006, status=active, department=qa, active=true, level=1, name=Frank, age=0, salary=-5000, score=0, tags=test}
7. [1] TestPerson{id=P007, score=10, status=active, level=9, tags=executive, department=management, name=Grace, age=65, salary=95000, active=true}
8. [1] TestPerson{id=P008, level=1, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, salary=68000, score=8.7, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior, status=active}
10. [1] TestPerson{id=P010, active=true, level=1, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, department=hr, level=1, name=Charlie, age=16, active=false, score=6, status=inactive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99}
6. [1] TestOrder{id=O006, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0}
7. [1] TestOrder{id=O007, priority=urgent, discount=50, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north}
10. [1] TestOrder{id=O010, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0, amount=1, total=75000, date=2024-03-15}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, total=999.99, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled, region=west}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35}
3. [1] TestPerson{id=P003, level=1, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr}
4. [1] TestPerson{id=P004, active=true, status=active, department=marketing, level=7, name=Diana, salary=85000, score=7.8, tags=manager, age=45}
5. [1] TestPerson{id=P005, name=Eve, age=30, status=inactive, level=3, salary=55000, active=false, score=8, tags=employee, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}
8. [1] TestPerson{id=P008, department=support, level=1, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, status=active, salary=68000, score=8.7, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, level=1, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}

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
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south}
6. [1] TestOrder{id=O006, priority=low, discount=0, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15}
7. [1] TestOrder{id=O007, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006}
8. [1] TestOrder{id=O008, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05}
9. [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed, priority=low, discount=10}
10. [1] TestOrder{id=O010, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0, amount=1, total=75000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O010, status=refunded, discount=0, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O004, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O006, priority=low, discount=0, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north, customer_id=P007}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35}
3. [1] TestPerson{id=P003, salary=0, tags=intern, department=hr, level=1, name=Charlie, age=16, active=false, score=6, status=inactive}
4. [1] TestPerson{id=P004, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7, name=Diana}
5. [1] TestPerson{id=P005, name=Eve, age=30, status=inactive, level=3, salary=55000, active=false, score=8, tags=employee, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}
9. [1] TestPerson{id=P009, status=active, salary=68000, score=8.7, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior}
10. [1] TestPerson{id=P010, tags=temp, status=active, department=intern, name=X, active=true, level=1, age=22, salary=28000, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3, salary=55000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, score=0, tags=test, status=active, department=qa, active=true, level=1, name=Frank, age=0, salary=-5000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, department=management, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, name=Henry}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, age=25, score=8.5, status=active, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, level=6, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, price=999.99, available=true, rating=4.5, keywords=computer, brand=TechCorp, supplier=TechSupply, stock=50, name=Laptop, category=electronics}
2. [1] TestProduct{id=PROD002, category=accessories, available=true, supplier=TechSupply, price=25.5, rating=4, keywords=peripheral, brand=TechCorp, stock=200, name=Mouse}
3. [1] TestProduct{id=PROD003, category=accessories, rating=3.5, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply, price=75, available=false, name=Keyboard}
4. [1] TestProduct{id=PROD004, price=299.99, available=true, rating=4.8, keywords=display, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics, brand=ScreenPro}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, rating=2, brand=OldTech, supplier=OldSupply, price=8.5, available=false, keywords=obsolete, stock=0}
6. [1] TestProduct{id=PROD006, available=true, rating=4.6, brand=AudioMax, name=Headphones, category=audio, price=150, keywords=sound, stock=75, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, brand=CamTech, stock=25, supplier=CamSupply, price=89.99, available=true, name=Webcam, category=electronics, rating=3.8, keywords=video}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, category=accessories, available=true, supplier=TechSupply, price=25.5, rating=4, keywords=peripheral, brand=TechCorp, stock=200}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, rating=3.5, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply, price=75, available=false}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD004, brand=ScreenPro, price=299.99, available=true, rating=4.8, keywords=display, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD006, name=Headphones, category=audio, price=150, keywords=sound, stock=75, supplier=AudioSupply, available=true, rating=4.6, brand=AudioMax}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD007, name=Webcam, category=electronics, rating=3.8, keywords=video, brand=CamTech, stock=25, supplier=CamSupply, price=89.99, available=true}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, available=true, rating=4.5, keywords=computer, brand=TechCorp, supplier=TechSupply, stock=50}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35}
3. [1] TestPerson{id=P003, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1, name=Charlie, age=16}
4. [1] TestPerson{id=P004, age=45, active=true, status=active, department=marketing, level=7, name=Diana, salary=85000, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3, salary=55000}
6. [1] TestPerson{id=P006, status=active, department=qa, active=true, level=1, name=Frank, age=0, salary=-5000, score=0, tags=test}
7. [1] TestPerson{id=P007, status=active, level=9, tags=executive, department=management, name=Grace, age=65, salary=95000, active=true, score=10}
8. [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, name=Henry, age=18}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6}
10. [1] TestPerson{id=P010, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1, age=22, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35, salary=75000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, department=hr, level=1, name=Charlie, age=16, active=false, score=6, status=inactive}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3, salary=55000, active=false}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}
4. [1] TestOrder{id=O004, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{id=O005, total=999.99, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1}
6. [1] TestOrder{id=O006, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98}
7. [1] TestOrder{id=O007, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north}
10. [1] TestOrder{id=O010, discount=0, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed, priority=low, discount=10}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O002, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O003, region=north, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, region=south, product_id=PROD001, amount=1, total=999.99, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2, name=Alice, salary=45000}
2. [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true, score=9.2}
3. [1] TestPerson{id=P003, department=hr, level=1, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern}
4. [1] TestPerson{id=P004, tags=manager, age=45, active=true, status=active, department=marketing, level=7, name=Diana, salary=85000, score=7.8}
5. [1] TestPerson{id=P005, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}
7. [1] TestPerson{id=P007, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management, name=Grace, age=65}
8. [1] TestPerson{id=P008, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, name=Henry}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6}
10. [1] TestPerson{id=P010, tags=temp, status=active, department=intern, name=X, active=true, level=1, age=22, salary=28000, score=6.5}
11. [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}
12. [1] TestOrder{id=O002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002}
13. [1] TestOrder{id=O003, discount=15, region=north, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003, status=shipped, priority=high}
14. [1] TestOrder{id=O004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004}
15. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south}
16. [1] TestOrder{id=O006, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0}
17. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}
18. [1] TestOrder{id=O008, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255}
19. [1] TestOrder{id=O009, region=north, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10}
20. [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing}
   - Fait 2: [1] TestOrder{id=O010, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0, amount=1, total=75000, date=2024-03-15}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3, salary=55000, active=false}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, name=X, active=true, level=1, age=22, salary=28000, score=6.5, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed, priority=low, discount=10}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, discount=0, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, age=30, status=inactive, level=3, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled, region=west}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, score=8.7, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O007, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal, region=south}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P006, score=0, tags=test, status=active, department=qa, active=true, level=1, name=Frank, age=0, salary=-5000}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal, region=south}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P010, active=true, level=1, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1, name=Frank}
   - Fait 2: [1] TestOrder{id=O001, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, salary=0, tags=intern, department=hr, level=1, name=Charlie, age=16, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, status=active, department=qa, active=true, level=1, name=Frank, age=0, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true, score=9.2}
   - Fait 2: [1] TestOrder{id=O006, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7, name=Diana, salary=85000}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O002, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35, salary=75000}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, status=refunded, discount=0, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P008, level=1, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, status=refunded, discount=0, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, name=X, active=true, level=1, age=22, salary=28000, score=6.5, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003, status=shipped}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P004, status=active, department=marketing, level=7, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true}
   - Fait 2: [1] TestOrder{id=O003, region=north, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P001, age=25, score=8.5, status=active, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false, score=5.5}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O005, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high, customer_id=P002, date=2024-02-10}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, age=45, active=true, status=active, department=marketing, level=7, name=Diana, salary=85000, score=7.8}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1, name=Charlie}
   - Fait 2: [1] TestOrder{id=O002, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O002, region=south, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P009, level=6, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P008, level=1, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled, region=west}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, status=active, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior, age=25}
   - Fait 2: [1] TestOrder{id=O006, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O009, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed, priority=low}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management, name=Grace}
   - Fait 2: [1] TestOrder{id=O006, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}
   - Fait 2: [1] TestOrder{id=O004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, status=active, department=marketing, level=7, name=Diana, salary=85000, score=7.8, tags=manager}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, status=refunded, discount=0, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1, age=22, salary=28000}
   - Fait 2: [1] TestOrder{id=O006, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing}
   - Fait 2: [1] TestOrder{id=O006, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1, name=Frank, age=0}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, region=north, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, status=refunded, discount=0, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, active=true, level=1, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, status=refunded, discount=0, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active}
   - Fait 2: [1] TestOrder{id=O003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003, status=shipped}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003, status=shipped}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, status=active, department=marketing, level=7, name=Diana, salary=85000, score=7.8, tags=manager}
   - Fait 2: [1] TestOrder{id=O004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P001, age=25, score=8.5, status=active, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, department=management, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3, salary=55000}
   - Fait 2: [1] TestOrder{id=O002, region=south, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, department=management, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O006, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P005, level=3, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O008, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P008, department=support, level=1, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal, region=south}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, status=inactive, level=3, salary=55000, active=false, score=8, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P005, level=3, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O005, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O007, discount=50, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal, region=south}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O006, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled, region=west}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P007, department=management, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, level=1, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1, name=Frank, age=0}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6, name=Ivy}
   - Fait 2: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P010, active=true, level=1, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, status=active, department=marketing, level=7, name=Diana, salary=85000, score=7.8, tags=manager}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, score=8.5, status=active, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior}
2. [1] TestPerson{id=P002, level=5, name=Bob, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering}
3. [1] TestPerson{id=P003, salary=0, tags=intern, department=hr, level=1, name=Charlie, age=16, active=false, score=6, status=inactive}
4. [1] TestPerson{id=P004, level=7, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing}
5. [1] TestPerson{id=P005, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3, salary=55000, active=false}
6. [1] TestPerson{id=P006, tags=test, status=active, department=qa, active=true, level=1, name=Frank, age=0, salary=-5000, score=0}
7. [1] TestPerson{id=P007, score=10, status=active, level=9, tags=executive, department=management, name=Grace, age=65, salary=95000, active=true}
8. [1] TestPerson{id=P008, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2, name=Alice}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, age=45, active=true, status=active, department=marketing, level=7, name=Diana, salary=85000, score=7.8}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, level=1, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, score=10, status=active, level=9, tags=executive, department=management, name=Grace, age=65, salary=95000, active=true}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob}
3. [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7, name=Diana}
5. [1] TestPerson{id=P005, age=30, status=inactive, level=3, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}
7. [1] TestPerson{id=P007, tags=executive, department=management, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9}
8. [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false}
9. [1] TestPerson{id=P009, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6, name=Ivy, age=40, active=true}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}
11. [1] TestOrder{id=O001, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98}
12. [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0}
13. [1] TestOrder{id=O003, region=north, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15}
14. [1] TestOrder{id=O004, region=east, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0}
15. [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high}
16. [1] TestOrder{id=O006, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0}
17. [1] TestOrder{id=O007, discount=50, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}
18. [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0}
19. [1] TestOrder{id=O009, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001}
20. [1] TestOrder{id=O010, status=refunded, discount=0, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, score=0, tags=test, status=active, department=qa, active=true, level=1, name=Frank, age=0, salary=-5000}
   - Fait 2: [1] TestOrder{id=O001, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P006, status=active, department=qa, active=true, level=1, name=Frank, age=0, salary=-5000, score=0, tags=test}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3, salary=55000}
   - Fait 2: [1] TestOrder{id=O005, region=south, product_id=PROD001, amount=1, total=999.99, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, department=hr, level=1, name=Charlie, age=16, active=false, score=6, status=inactive}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled, region=west}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O006, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, department=management, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high, customer_id=P002}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, level=1, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern}
   - Fait 2: [1] TestOrder{id=O008, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal, region=south}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed, priority=low, discount=10}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true, score=9.2, tags=senior}
   - Fait 2: [1] TestOrder{id=O001, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true, score=9.2}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P007, status=active, level=9, tags=executive, department=management, name=Grace, age=65, salary=95000, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal, region=south}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, name=Eve, age=30, status=inactive, level=3, salary=55000, active=false, score=8, tags=employee}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, status=refunded, discount=0, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, department=hr, level=1, name=Charlie, age=16, active=false, score=6, status=inactive}
   - Fait 2: [1] TestOrder{id=O002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P010, name=X, active=true, level=1, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P007, department=management, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O006, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35}
   - Fait 2: [1] TestOrder{id=O007, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, amount=4}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P001, age=25, score=8.5, status=active, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P010, name=X, active=true, level=1, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern}
   - Fait 2: [1] TestOrder{id=O002, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O003, discount=15, region=north, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003, status=shipped, priority=high}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high, customer_id=P002}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled, region=west}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, status=active, department=qa, active=true, level=1, name=Frank, age=0, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0, amount=1}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P001, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2, name=Alice, salary=45000}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P009, status=active, salary=68000, score=8.7, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O009, region=north, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P005, age=30, status=inactive, level=3, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, status=inactive, level=3, salary=55000, active=false, score=8, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O004, region=east, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O005, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, status=active, department=marketing, level=7, name=Diana, salary=85000, score=7.8, tags=manager}
   - Fait 2: [1] TestOrder{id=O008, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0, amount=1, total=75000}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35}
   - Fait 2: [1] TestOrder{id=O002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P003, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1, name=Charlie, age=16, active=false}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2, name=Alice}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed, priority=low, discount=10}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35}
   - Fait 2: [1] TestOrder{id=O003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, level=1, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern}
   - Fait 2: [1] TestOrder{id=O005, region=south, product_id=PROD001, amount=1, total=999.99, priority=high, customer_id=P002, date=2024-02-10, status=confirmed, discount=100}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P005, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3, salary=55000, active=false}
   - Fait 2: [1] TestOrder{id=O006, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O007, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, amount=4}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P007, score=10, status=active, level=9, tags=executive, department=management, name=Grace, age=65, salary=95000, active=true}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P006, score=0, tags=test, status=active, department=qa, active=true, level=1, name=Frank, age=0, salary=-5000}
   - Fait 2: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, status=active, department=marketing, level=7, name=Diana, salary=85000, score=7.8, tags=manager}
   - Fait 2: [1] TestOrder{id=O001, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, active=true, level=1, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active}
   - Fait 2: [1] TestOrder{id=O002, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003, status=shipped}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, department=management, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9}
   - Fait 2: [1] TestOrder{id=O006, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O008, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0, amount=1}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6, name=Ivy}
   - Fait 2: [1] TestOrder{id=O001, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P001, age=25, score=8.5, status=active, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior}
   - Fait 2: [1] TestOrder{id=O002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7}
   - Fait 2: [1] TestOrder{id=O004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true, score=9.2, tags=senior}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0, amount=1, total=75000}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, status=refunded, discount=0, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3, salary=55000, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O005, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high, customer_id=P002, date=2024-02-10, status=confirmed}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north, customer_id=P007}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O008, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001, status=completed}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3, salary=55000}
   - Fait 2: [1] TestOrder{id=O003, discount=15, region=north, amount=3, total=225, date=2024-02-01, customer_id=P001, product_id=PROD003, status=shipped, priority=high}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, score=8.7, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0, amount=2}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7, name=Diana}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, discount=0, total=255, date=2024-03-05, status=pending, priority=normal, region=south}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, status=active, level=9, tags=executive, department=management, name=Grace, age=65, salary=95000}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1, age=22, salary=28000}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O006, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, age=45, active=true, status=active, department=marketing, level=7, name=Diana, salary=85000, score=7.8}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P003, level=1, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, status=active, level=9, tags=executive, department=management, name=Grace, age=65, salary=95000}
   - Fait 2: [1] TestOrder{id=O002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false, score=5.5}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, department=hr, level=1, name=Charlie, age=16, active=false, score=6, status=inactive}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1, name=Frank}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled, region=west}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, status=inactive, level=3, salary=55000, active=false, score=8, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, score=8.7, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2, name=Alice}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6, name=Ivy}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, age=25, score=8.5, status=active, department=sales, level=2, name=Alice, salary=45000, active=true}
   - Fait 2: [1] TestOrder{id=O004, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004, amount=1}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1}
   - Fait 2: [1] TestOrder{id=O006, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, status=active, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior, age=25}
2. [1] TestPerson{id=P002, status=active, department=engineering, level=5, name=Bob, age=35, salary=75000, active=true, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1, name=Charlie, age=16, active=false}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}
7. [1] TestPerson{id=P007, tags=executive, department=management, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}
9. [1] TestPerson{id=P009, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6, name=Ivy, age=40}
10. [1] TestPerson{id=P010, name=X, active=true, level=1, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, department=management, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P010, name=X, active=true, level=1, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, department=support, level=1, name=Henry, age=18, salary=25000, active=false}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, level=6, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, tags=junior, age=25, score=8.5, status=active, department=sales, level=2}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true, status=active, department=marketing, level=7}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P005, level=3, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, region=north, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50}
2. [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5, status=confirmed, priority=low, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}
4. [1] TestOrder{id=O004, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east, product_id=PROD004, total=299.99}
5. [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high}
6. [1] TestOrder{id=O006, total=999.98, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled, region=west, customer_id=P005, product_id=PROD005}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north}
8. [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0}
9. [1] TestOrder{id=O009, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, region=north, customer_id=P001}
10. [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O009, date=2024-03-10, region=north, customer_id=P001, status=completed, priority=low, discount=10, product_id=PROD007, amount=1, total=89.99}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O010, priority=urgent, region=east, customer_id=P006, product_id=PROD001, status=refunded, discount=0, amount=1, total=75000, date=2024-03-15}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, status=shipped, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, total=299.99, date=2024-02-05, customer_id=P004, amount=1, status=delivered, priority=normal, discount=0, region=east}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, region=west, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, discount=0, amount=2, status=cancelled}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, customer_id=P001, total=1999.98, discount=50, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, status=confirmed, priority=low, region=south, customer_id=P002, date=2024-01-20, discount=0, product_id=PROD002, amount=1, total=25.5}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, status=confirmed, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, priority=high}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O007, status=shipped, priority=urgent, discount=50, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, score=8.5, status=active, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob, age=35}
3. [1] TestPerson{id=P003, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, status=active, department=marketing, level=7, name=Diana, salary=85000, score=7.8, tags=manager, age=45, active=true}
5. [1] TestPerson{id=P005, name=Eve, age=30, status=inactive, level=3, salary=55000, active=false, score=8, tags=employee, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, status=active, department=qa, active=true, level=1}
7. [1] TestPerson{id=P007, department=management, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6}
10. [1] TestPerson{id=P010, active=true, level=1, age=22, salary=28000, score=6.5, tags=temp, status=active, department=intern, name=X}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, tags=employee, department=sales, name=Eve, age=30, status=inactive, level=3}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, active=true, score=10, status=active, level=9, tags=executive, department=management}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, tags=senior, status=active, salary=68000, score=8.7, department=engineering, level=6, name=Ivy}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, tags=temp, status=active, department=intern, name=X, active=true, level=1, age=22, salary=28000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, age=25, score=8.5, status=active, department=sales, level=2, name=Alice, salary=45000, active=true, tags=junior}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, level=1, name=Charlie, age=16, active=false, score=6, status=inactive, salary=0}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P006, score=0, tags=test, status=active, department=qa, active=true, level=1, name=Frank, age=0, salary=-5000}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, department=support, level=1, name=Henry, age=18, salary=25000, active=false, score=5.5, tags=junior, status=inactive}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, department=engineering, level=5, name=Bob}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, status=active, department=marketing, level=7, name=Diana, salary=85000, score=7.8, tags=manager}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
