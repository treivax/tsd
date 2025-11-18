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

1. [1] TestPerson{id=P001, name=Alice, department=sales, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true, score=9.2, status=active, level=5}
3. [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, level=1, name=Charlie, salary=0, status=inactive}
4. [1] TestPerson{id=P004, active=true, tags=manager, department=marketing, name=Diana, salary=85000, score=7.8, status=active, level=7, age=45}
5. [1] TestPerson{id=P005, level=3, age=30, active=false, tags=employee, name=Eve, salary=55000, score=8, status=inactive, department=sales}
6. [1] TestPerson{id=P006, tags=test, salary=-5000, status=active, department=qa, level=1, name=Frank, age=0, active=true, score=0}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, level=9, active=true, status=active, department=management}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, department=support, level=1, active=false, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, active=true, score=8.7, status=active, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, level=6}
10. [1] TestPerson{id=P010, status=active, level=1, active=true, department=intern, name=X, age=22, salary=28000, score=6.5, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, level=2, name=Alice, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, status=active, level=5, age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true, score=9.2}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, level=1, name=Charlie, salary=0, status=inactive}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, level=7, age=45, active=true, tags=manager, department=marketing, name=Diana, salary=85000, score=7.8, status=active}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, status=inactive, department=sales, level=3, age=30, active=false, tags=employee}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, active=true, status=active, department=management, name=Grace, age=65, salary=95000, score=10, tags=executive, level=9}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, level=1, active=false, score=5.5, status=inactive, name=Henry, age=18, salary=25000}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, level=6, active=true, score=8.7, status=active, department=engineering}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, priority=normal, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, region=north, amount=2, date=2024-01-15}
2. [1] TestOrder{id=O002, priority=low, discount=0, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, region=south, amount=1, total=25.5}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north}
4. [1] TestOrder{id=O004, status=delivered, priority=normal, discount=0, region=east, amount=1, date=2024-02-05, customer_id=P004, product_id=PROD004, total=299.99}
5. [1] TestOrder{id=O005, product_id=PROD001, status=confirmed, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, discount=100}
6. [1] TestOrder{id=O006, discount=0, region=west, customer_id=P005, amount=2, date=2024-02-15, product_id=PROD005, total=999.98, status=cancelled, priority=low}
7. [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, priority=urgent, product_id=PROD006, date=2024-03-01, status=shipped, discount=50, region=north}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, region=north, amount=1, status=completed, discount=10}
10. [1] TestOrder{id=O010, total=75000, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006, status=refunded, product_id=PROD001, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, amount=2, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, total=25.5, priority=low, discount=0, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, region=south}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, priority=urgent, product_id=PROD006, date=2024-03-01, status=shipped, discount=50, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, status=completed, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, region=north, amount=1}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal, discount=0, region=east, amount=1, date=2024-02-05}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, discount=100, product_id=PROD001, status=confirmed, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255, customer_id=P010}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O010, status=refunded, product_id=PROD001, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, name=Alice, department=sales, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, status=active, level=5, age=35, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, level=1, name=Charlie, salary=0, status=inactive, age=16}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, status=active, level=7, age=45, active=true, tags=manager, department=marketing}
5. [1] TestPerson{id=P005, department=sales, level=3, age=30, active=false, tags=employee, name=Eve, salary=55000, score=8, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, salary=-5000, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, level=9, active=true, status=active, department=management}
8. [1] TestPerson{id=P008, department=support, level=1, active=false, score=5.5, status=inactive, name=Henry, age=18, salary=25000, tags=junior}
9. [1] TestPerson{id=P009, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, level=6, active=true, score=8.7, status=active}
10. [1] TestPerson{id=P010, active=true, department=intern, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, level=6, active=true, score=8.7, status=active, department=engineering}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, level=2, name=Alice, department=sales}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P002, level=5, age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true, score=9.2, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, status=active, level=7, age=45, active=true, tags=manager, department=marketing}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, active=false, tags=employee, name=Eve, salary=55000, score=8, status=inactive, department=sales, level=3, age=30}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, score=10, tags=executive, level=9, active=true, status=active, department=management, name=Grace}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, region=north, amount=2, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, region=south, amount=1, total=25.5, priority=low, discount=0}
3. [1] TestOrder{id=O003, amount=3, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north, customer_id=P001, product_id=PROD003}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal, discount=0, region=east, amount=1, date=2024-02-05}
5. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, discount=100, product_id=PROD001, status=confirmed, priority=high, region=south, customer_id=P002}
6. [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, priority=urgent, product_id=PROD006, date=2024-03-01, status=shipped, discount=50, region=north}
8. [1] TestOrder{id=O008, amount=10, total=255, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south}
9. [1] TestOrder{id=O009, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, region=north, amount=1, status=completed}
10. [1] TestOrder{id=O010, customer_id=P006, status=refunded, product_id=PROD001, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, total=25.5, priority=low, discount=0, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O001, amount=2, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, region=north}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O004, region=east, amount=1, date=2024-02-05, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, status=confirmed, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, discount=100}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, amount=2, date=2024-02-15, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0, region=west}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, amount=4, total=600, priority=urgent, product_id=PROD006}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, amount=1, status=completed, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, available=true, keywords=computer, rating=4.5, brand=TechCorp, stock=50, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, name=Mouse, available=true, rating=4, keywords=peripheral, brand=TechCorp, category=accessories, price=25.5, stock=200, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0, name=Keyboard, category=accessories, price=75, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, name=Monitor, price=299.99, keywords=display, stock=30, category=electronics, available=true, rating=4.8, brand=ScreenPro, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, keywords=obsolete, stock=0, name=OldKeyboard, category=accessories, available=false, brand=OldTech, supplier=OldSupply, price=8.5, rating=2}
6. [1] TestProduct{id=PROD006, category=audio, price=150, rating=4.6, keywords=sound, stock=75, supplier=AudioSupply, name=Headphones, available=true, brand=AudioMax}
7. [1] TestProduct{id=PROD007, available=true, rating=3.8, keywords=video, stock=25, name=Webcam, category=electronics, price=89.99, brand=CamTech, supplier=CamSupply}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD006, name=Headphones, available=true, brand=AudioMax, category=audio, price=150, rating=4.6, keywords=sound, stock=75, supplier=AudioSupply}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD007, supplier=CamSupply, available=true, rating=3.8, keywords=video, stock=25, name=Webcam, category=electronics, price=89.99, brand=CamTech}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD001, rating=4.5, brand=TechCorp, stock=50, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, available=true, keywords=computer}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD002, category=accessories, price=25.5, stock=200, supplier=TechSupply, name=Mouse, available=true, rating=4, keywords=peripheral, brand=TechCorp}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, price=75, supplier=KeySupply, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD004, category=electronics, available=true, rating=4.8, brand=ScreenPro, supplier=ScreenSupply, name=Monitor, price=299.99, keywords=display, stock=30}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, department=sales, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, active=true, score=9.2, status=active, level=5, age=35, tags=senior, department=engineering, name=Bob, salary=75000}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, status=inactive, age=16, active=false, score=6, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, status=active, level=7, age=45, active=true, tags=manager, department=marketing}
5. [1] TestPerson{id=P005, score=8, status=inactive, department=sales, level=3, age=30, active=false, tags=employee, name=Eve, salary=55000}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, salary=-5000, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, status=active, department=management, name=Grace, age=65, salary=95000, score=10, tags=executive, level=9, active=true}
8. [1] TestPerson{id=P008, level=1, active=false, score=5.5, status=inactive, name=Henry, age=18, salary=25000, tags=junior, department=support}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, level=6, active=true, score=8.7, status=active, department=engineering}
10. [1] TestPerson{id=P010, active=true, department=intern, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P006, active=true, score=0, tags=test, salary=-5000, status=active, department=qa, level=1, name=Frank, age=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, level=1, name=Charlie, salary=0, status=inactive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, region=north, amount=2, date=2024-01-15, priority=normal}
2. [1] TestOrder{id=O002, priority=low, discount=0, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, region=south, amount=1, total=25.5}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north}
4. [1] TestOrder{id=O004, total=299.99, status=delivered, priority=normal, discount=0, region=east, amount=1, date=2024-02-05, customer_id=P004, product_id=PROD004}
5. [1] TestOrder{id=O005, product_id=PROD001, status=confirmed, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, discount=100}
6. [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, product_id=PROD006, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, amount=4, total=600, priority=urgent}
8. [1] TestOrder{id=O008, discount=0, region=south, amount=10, total=255, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal}
9. [1] TestOrder{id=O009, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, region=north, amount=1, status=completed, discount=10, customer_id=P001}
10. [1] TestOrder{id=O010, status=refunded, product_id=PROD001, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, amount=2, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, region=north, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high, total=225, discount=15}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal, discount=0, region=east, amount=1, date=2024-02-05}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, amount=2, date=2024-02-15}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, amount=4, total=600, priority=urgent}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255, customer_id=P010, product_id=PROD002}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, region=north, amount=1, status=completed, discount=10}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006, status=refunded, product_id=PROD001}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, total=25.5, priority=low, discount=0, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, region=south}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, status=confirmed, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, discount=100}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, status=active, level=2, name=Alice, department=sales, age=25, salary=45000, active=true, score=8.5}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, status=active, level=5, age=35, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, level=1, name=Charlie, salary=0, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, status=active, level=7, age=45, active=true, tags=manager, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, status=inactive, department=sales, level=3, age=30, active=false, tags=employee}
6. [1] TestPerson{id=P006, active=true, score=0, tags=test, salary=-5000, status=active, department=qa, level=1, name=Frank, age=0}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, level=9, active=true, status=active, department=management}
8. [1] TestPerson{id=P008, level=1, active=false, score=5.5, status=inactive, name=Henry, age=18, salary=25000, tags=junior, department=support}
9. [1] TestPerson{id=P009, score=8.7, status=active, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, level=6, active=true}
10. [1] TestPerson{id=P010, active=true, department=intern, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, tags=junior, department=support, level=1, active=false, score=5.5, status=inactive, name=Henry, age=18}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, salary=0, status=inactive, age=16, active=false, score=6, tags=intern, department=hr, level=1}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, level=3, age=30, active=false, tags=employee, name=Eve, salary=55000, score=8, status=inactive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, region=north, amount=2, date=2024-01-15, priority=normal}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, region=south, amount=1, total=25.5, priority=low, discount=0}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north}
4. [1] TestOrder{id=O004, priority=normal, discount=0, region=east, amount=1, date=2024-02-05, customer_id=P004, product_id=PROD004, total=299.99, status=delivered}
5. [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, discount=100, product_id=PROD001, status=confirmed, priority=high, region=south}
6. [1] TestOrder{id=O006, region=west, customer_id=P005, amount=2, date=2024-02-15, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0}
7. [1] TestOrder{id=O007, status=shipped, discount=50, region=north, customer_id=P007, amount=4, total=600, priority=urgent, product_id=PROD006, date=2024-03-01}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, region=north, amount=1, status=completed, discount=10}
10. [1] TestOrder{id=O010, priority=urgent, discount=0, region=east, customer_id=P006, status=refunded, product_id=PROD001, amount=1, total=75000, date=2024-03-15}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006, status=refunded, product_id=PROD001}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, total=25.5, priority=low, discount=0, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, region=south}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, status=shipped, priority=high, total=225, discount=15, region=north, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal, discount=0, region=east, amount=1, date=2024-02-05}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, discount=100, product_id=PROD001, status=confirmed, priority=high, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, amount=2, date=2024-02-15, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0, region=west}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, status=shipped, discount=50, region=north, customer_id=P007, amount=4, total=600, priority=urgent, product_id=PROD006, date=2024-03-01}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, region=north, amount=1, status=completed, discount=10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, level=2, name=Alice, department=sales, age=25, salary=45000, active=true, score=8.5, tags=junior}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, status=active, level=5, age=35, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, level=1, name=Charlie, salary=0, status=inactive}
4. [1] TestPerson{id=P004, score=7.8, status=active, level=7, age=45, active=true, tags=manager, department=marketing, name=Diana, salary=85000}
5. [1] TestPerson{id=P005, salary=55000, score=8, status=inactive, department=sales, level=3, age=30, active=false, tags=employee, name=Eve}
6. [1] TestPerson{id=P006, department=qa, level=1, name=Frank, age=0, active=true, score=0, tags=test, salary=-5000, status=active}
7. [1] TestPerson{id=P007, tags=executive, level=9, active=true, status=active, department=management, name=Grace, age=65, salary=95000, score=10}
8. [1] TestPerson{id=P008, salary=25000, tags=junior, department=support, level=1, active=false, score=5.5, status=inactive, name=Henry, age=18}
9. [1] TestPerson{id=P009, active=true, score=8.7, status=active, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, level=6}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, level=3, age=30, active=false, tags=employee, name=Eve, salary=55000, score=8, status=inactive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, active=true, status=active, department=management, name=Grace, age=65, salary=95000, score=10, tags=executive, level=9}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, status=active, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, level=6, active=true, score=8.7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, level=1, name=Charlie, salary=0, status=inactive}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, status=active, level=7, age=45, active=true, tags=manager, department=marketing}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, salary=-5000, status=active, department=qa, level=1, name=Frank, age=0, active=true, score=0}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, status=inactive, name=Henry, age=18, salary=25000, tags=junior, department=support, level=1}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, level=2, name=Alice}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, status=active, level=5, age=35, tags=senior, department=engineering}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, available=true, keywords=computer, rating=4.5, brand=TechCorp, stock=50, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, brand=TechCorp, category=accessories, price=25.5, stock=200, supplier=TechSupply, name=Mouse, available=true, rating=4, keywords=peripheral}
3. [1] TestProduct{id=PROD003, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0, name=Keyboard, category=accessories, price=75, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, category=electronics, available=true, rating=4.8, brand=ScreenPro, supplier=ScreenSupply, name=Monitor, price=299.99, keywords=display, stock=30}
5. [1] TestProduct{id=PROD005, brand=OldTech, supplier=OldSupply, price=8.5, rating=2, keywords=obsolete, stock=0, name=OldKeyboard, category=accessories, available=false}
6. [1] TestProduct{id=PROD006, category=audio, price=150, rating=4.6, keywords=sound, stock=75, supplier=AudioSupply, name=Headphones, available=true, brand=AudioMax}
7. [1] TestProduct{id=PROD007, supplier=CamSupply, available=true, rating=3.8, keywords=video, stock=25, name=Webcam, category=electronics, price=89.99, brand=CamTech}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, brand=TechCorp, stock=50, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, available=true, keywords=computer, rating=4.5}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, available=true, rating=4, keywords=peripheral, brand=TechCorp, category=accessories, price=25.5, stock=200, supplier=TechSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0, name=Keyboard, category=accessories, price=75, supplier=KeySupply}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, category=electronics, available=true, rating=4.8, brand=ScreenPro, supplier=ScreenSupply, name=Monitor, price=299.99, keywords=display, stock=30}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, stock=75, supplier=AudioSupply, name=Headphones, available=true, brand=AudioMax, category=audio, price=150, rating=4.6, keywords=sound}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, available=true, rating=3.8, keywords=video, stock=25, name=Webcam, category=electronics, price=89.99, brand=CamTech, supplier=CamSupply}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, status=active, level=2, name=Alice, department=sales, age=25, salary=45000}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, status=active, level=5, age=35, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, status=inactive, age=16, active=false, score=6, tags=intern, department=hr, level=1, name=Charlie, salary=0}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, status=active, level=7, age=45, active=true, tags=manager, department=marketing}
5. [1] TestPerson{id=P005, tags=employee, name=Eve, salary=55000, score=8, status=inactive, department=sales, level=3, age=30, active=false}
6. [1] TestPerson{id=P006, salary=-5000, status=active, department=qa, level=1, name=Frank, age=0, active=true, score=0, tags=test}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, level=9, active=true, status=active, department=management}
8. [1] TestPerson{id=P008, active=false, score=5.5, status=inactive, name=Henry, age=18, salary=25000, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, level=6, active=true, score=8.7, status=active, department=engineering}
10. [1] TestPerson{id=P010, active=true, department=intern, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, status=active, level=5, age=35, tags=senior, department=engineering}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, level=1, name=Charlie, salary=0, status=inactive}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, status=inactive, department=sales, level=3, age=30, active=false, tags=employee}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, level=1, name=Frank, age=0, active=true, score=0, tags=test, salary=-5000, status=active, department=qa}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, status=inactive, name=Henry, age=18, salary=25000, tags=junior, department=support, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, level=6, active=true, score=8.7, status=active, department=engineering, name=Ivy, age=40, salary=68000}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, department=sales, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, level=2}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, status=active, level=7, age=45, active=true, tags=manager, department=marketing}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, score=10, tags=executive, level=9, active=true, status=active, department=management, name=Grace}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, region=north, amount=2, date=2024-01-15, priority=normal}
2. [1] TestOrder{id=O002, amount=1, total=25.5, priority=low, discount=0, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal, discount=0, region=east, amount=1, date=2024-02-05}
5. [1] TestOrder{id=O005, product_id=PROD001, status=confirmed, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, discount=100}
6. [1] TestOrder{id=O006, customer_id=P005, amount=2, date=2024-02-15, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0, region=west}
7. [1] TestOrder{id=O007, amount=4, total=600, priority=urgent, product_id=PROD006, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007}
8. [1] TestOrder{id=O008, priority=normal, discount=0, region=south, amount=10, total=255, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, region=north, amount=1, status=completed, discount=10}
10. [1] TestOrder{id=O010, total=75000, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006, status=refunded, product_id=PROD001, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, amount=2, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, region=south, amount=1, total=25.5, priority=low, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O004, region=east, amount=1, date=2024-02-05, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal, discount=0}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, discount=100, product_id=PROD001, status=confirmed, priority=high, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, amount=4, total=600, priority=urgent}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, date=2024-03-10, priority=low, region=north, amount=1, status=completed, discount=10, customer_id=P001, product_id=PROD007, total=89.99}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, level=2, name=Alice, department=sales}
2. [1] TestPerson{id=P002, level=5, age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true, score=9.2, status=active}
3. [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, level=1, name=Charlie, salary=0, status=inactive}
4. [1] TestPerson{id=P004, salary=85000, score=7.8, status=active, level=7, age=45, active=true, tags=manager, department=marketing, name=Diana}
5. [1] TestPerson{id=P005, active=false, tags=employee, name=Eve, salary=55000, score=8, status=inactive, department=sales, level=3, age=30}
6. [1] TestPerson{id=P006, level=1, name=Frank, age=0, active=true, score=0, tags=test, salary=-5000, status=active, department=qa}
7. [1] TestPerson{id=P007, score=10, tags=executive, level=9, active=true, status=active, department=management, name=Grace, age=65, salary=95000}
8. [1] TestPerson{id=P008, active=false, score=5.5, status=inactive, name=Henry, age=18, salary=25000, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, status=active, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, level=6, active=true, score=8.7}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true, department=intern}
11. [1] TestOrder{id=O001, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, region=north, amount=2}
12. [1] TestOrder{id=O002, product_id=PROD002, date=2024-01-20, status=confirmed, region=south, amount=1, total=25.5, priority=low, discount=0, customer_id=P002}
13. [1] TestOrder{id=O003, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north, customer_id=P001, product_id=PROD003, amount=3}
14. [1] TestOrder{id=O004, date=2024-02-05, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal, discount=0, region=east, amount=1}
15. [1] TestOrder{id=O005, status=confirmed, priority=high, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, discount=100, product_id=PROD001}
16. [1] TestOrder{id=O006, customer_id=P005, amount=2, date=2024-02-15, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0, region=west}
17. [1] TestOrder{id=O007, customer_id=P007, amount=4, total=600, priority=urgent, product_id=PROD006, date=2024-03-01, status=shipped, discount=50, region=north}
18. [1] TestOrder{id=O008, priority=normal, discount=0, region=south, amount=10, total=255, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending}
19. [1] TestOrder{id=O009, amount=1, status=completed, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, region=north}
20. [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006, status=refunded, product_id=PROD001}

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

1. [1] TestPerson{id=P001, salary=45000, active=true, score=8.5, tags=junior, status=active, level=2, name=Alice, department=sales, age=25}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, status=active, level=5, age=35, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, level=1, name=Charlie, salary=0, status=inactive, age=16, active=false, score=6, tags=intern, department=hr}
4. [1] TestPerson{id=P004, status=active, level=7, age=45, active=true, tags=manager, department=marketing, name=Diana, salary=85000, score=7.8}
5. [1] TestPerson{id=P005, active=false, tags=employee, name=Eve, salary=55000, score=8, status=inactive, department=sales, level=3, age=30}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, salary=-5000, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, salary=95000, score=10, tags=executive, level=9, active=true, status=active, department=management, name=Grace, age=65}
8. [1] TestPerson{id=P008, active=false, score=5.5, status=inactive, name=Henry, age=18, salary=25000, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, active=true, score=8.7, status=active, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, level=6}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, department=sales, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, level=2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, status=active, level=5, age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true, score=9.2}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, status=active, level=7, age=45, active=true, tags=manager, department=marketing}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, salary=-5000, status=active, department=qa, level=1, name=Frank, age=0, active=true, score=0}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, level=9, active=true, status=active, department=management}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, level=6, active=true, score=8.7, status=active, department=engineering}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, active=true, department=intern, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, tags=junior, status=active, level=2, name=Alice, department=sales, age=25, salary=45000, active=true}
2. [1] TestPerson{id=P002, age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true, score=9.2, status=active, level=5}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, status=inactive, age=16, active=false, score=6, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, department=marketing, name=Diana, salary=85000, score=7.8, status=active, level=7, age=45, active=true, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, status=inactive, department=sales, level=3, age=30, active=false, tags=employee}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, salary=-5000, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, level=9, active=true, status=active, department=management, name=Grace, age=65, salary=95000, score=10, tags=executive}
8. [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, department=support, level=1, active=false, score=5.5, status=inactive, name=Henry}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, level=6, active=true, score=8.7, status=active, department=engineering}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true, department=intern}
11. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, region=north, amount=2, date=2024-01-15, priority=normal}
12. [1] TestOrder{id=O002, status=confirmed, region=south, amount=1, total=25.5, priority=low, discount=0, customer_id=P002, product_id=PROD002, date=2024-01-20}
13. [1] TestOrder{id=O003, total=225, discount=15, region=north, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high}
14. [1] TestOrder{id=O004, discount=0, region=east, amount=1, date=2024-02-05, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal}
15. [1] TestOrder{id=O005, region=south, customer_id=P002, amount=1, total=999.99, date=2024-02-10, discount=100, product_id=PROD001, status=confirmed, priority=high}
16. [1] TestOrder{id=O006, customer_id=P005, amount=2, date=2024-02-15, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0, region=west}
17. [1] TestOrder{id=O007, amount=4, total=600, priority=urgent, product_id=PROD006, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007}
18. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south, amount=10, total=255}
19. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, region=north, amount=1, status=completed, discount=10}
20. [1] TestOrder{id=O010, product_id=PROD001, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006, status=refunded}

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

1. [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, level=2, name=Alice, department=sales}
2. [1] TestPerson{id=P002, age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true, score=9.2, status=active, level=5}
3. [1] TestPerson{id=P003, age=16, active=false, score=6, tags=intern, department=hr, level=1, name=Charlie, salary=0, status=inactive}
4. [1] TestPerson{id=P004, age=45, active=true, tags=manager, department=marketing, name=Diana, salary=85000, score=7.8, status=active, level=7}
5. [1] TestPerson{id=P005, tags=employee, name=Eve, salary=55000, score=8, status=inactive, department=sales, level=3, age=30, active=false}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, salary=-5000, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, score=10, tags=executive, level=9, active=true, status=active, department=management, name=Grace, age=65, salary=95000}
8. [1] TestPerson{id=P008, department=support, level=1, active=false, score=5.5, status=inactive, name=Henry, age=18, salary=25000, tags=junior}
9. [1] TestPerson{id=P009, salary=68000, tags=senior, level=6, active=true, score=8.7, status=active, department=engineering, name=Ivy, age=40}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, active=true, score=8.7, status=active, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, level=6}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, name=Eve, salary=55000, score=8, status=inactive, department=sales, level=3, age=30, active=false}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true, department=intern}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, department=sales, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, level=2}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, level=5, age=35, tags=senior, department=engineering, name=Bob, salary=75000, active=true, score=9.2, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, age=16, active=false, score=6, tags=intern, department=hr, level=1, name=Charlie, salary=0}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, status=active, level=7, age=45, active=true, tags=manager, department=marketing}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, level=9, active=true, status=active, department=management}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, department=support, level=1, active=false, score=5.5, status=inactive, name=Henry}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, region=north, amount=2, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50}
2. [1] TestOrder{id=O002, date=2024-01-20, status=confirmed, region=south, amount=1, total=25.5, priority=low, discount=0, customer_id=P002, product_id=PROD002}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high, total=225, discount=15, region=north}
4. [1] TestOrder{id=O004, discount=0, region=east, amount=1, date=2024-02-05, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal}
5. [1] TestOrder{id=O005, customer_id=P002, amount=1, total=999.99, date=2024-02-10, discount=100, product_id=PROD001, status=confirmed, priority=high, region=south}
6. [1] TestOrder{id=O006, region=west, customer_id=P005, amount=2, date=2024-02-15, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0}
7. [1] TestOrder{id=O007, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, amount=4, total=600, priority=urgent, product_id=PROD006}
8. [1] TestOrder{id=O008, amount=10, total=255, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south}
9. [1] TestOrder{id=O009, amount=1, status=completed, discount=10, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, region=north}
10. [1] TestOrder{id=O010, discount=0, region=east, customer_id=P006, status=refunded, product_id=PROD001, amount=1, total=75000, date=2024-03-15, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, amount=4, total=600, priority=urgent}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, total=255, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, priority=normal, discount=0, region=south}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, region=north, amount=1, status=completed, discount=10}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, amount=1, total=75000, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006, status=refunded}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, total=25.5, priority=low, discount=0, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O003, total=225, discount=15, region=north, customer_id=P001, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, priority=high}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, discount=100, product_id=PROD001, status=confirmed, priority=high, region=south, customer_id=P002}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O001, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, region=north, amount=2}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal, discount=0, region=east, amount=1, date=2024-02-05}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, amount=2, date=2024-02-15, product_id=PROD005}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, active=true, score=8.5, tags=junior, status=active, level=2, name=Alice, department=sales, age=25}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, status=active, level=5, age=35, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, salary=0, status=inactive, age=16, active=false, score=6, tags=intern, department=hr, level=1, name=Charlie}
4. [1] TestPerson{id=P004, age=45, active=true, tags=manager, department=marketing, name=Diana, salary=85000, score=7.8, status=active, level=7}
5. [1] TestPerson{id=P005, salary=55000, score=8, status=inactive, department=sales, level=3, age=30, active=false, tags=employee, name=Eve}
6. [1] TestPerson{id=P006, tags=test, salary=-5000, status=active, department=qa, level=1, name=Frank, age=0, active=true, score=0}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, level=9, active=true, status=active, department=management}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, department=support, level=1, active=false, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, active=true, score=8.7, status=active, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, level=6}
10. [1] TestPerson{id=P010, active=true, department=intern, name=X, age=22, salary=28000, score=6.5, tags=temp, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, active=true, score=8.5, tags=junior, status=active, level=2, name=Alice, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, department=engineering, name=Bob, salary=75000, active=true, score=9.2, status=active, level=5, age=35}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, status=active, level=7, age=45, active=true, tags=manager, department=marketing}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, status=inactive, department=sales, level=3, age=30, active=false, tags=employee}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P006, active=true, score=0, tags=test, salary=-5000, status=active, department=qa, level=1, name=Frank, age=0}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, department=management, name=Grace, age=65, salary=95000, score=10, tags=executive, level=9, active=true, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, department=support, level=1, active=false, score=5.5, status=inactive}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, tags=temp, status=active, level=1, active=true, department=intern, name=X, age=22}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, salary=0, status=inactive, age=16, active=false, score=6, tags=intern, department=hr, level=1}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, status=active, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, level=6, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 17 (89.5%)
- **Tokens générés**: 130
- **Faits traités**: 27
