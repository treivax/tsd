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

1. [1] TestPerson{id=P001, age=25, salary=45000, status=active, department=sales, name=Alice, active=true, score=8.5, tags=junior, level=2}
2. [1] TestPerson{id=P002, active=true, tags=senior, department=engineering, name=Bob, age=35, salary=75000, score=9.2, status=active, level=5}
3. [1] TestPerson{id=P003, age=16, salary=0, score=6, status=inactive, department=hr, name=Charlie, active=false, tags=intern, level=1}
4. [1] TestPerson{id=P004, active=true, score=7.8, name=Diana, age=45, salary=85000, tags=manager, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, salary=55000, score=8, tags=employee, level=3, age=30, active=false, status=inactive, department=sales, name=Eve}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, active=true, tags=test, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, level=9, salary=95000, active=true, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, salary=25000, score=5.5, tags=junior, status=inactive, department=support, name=Henry, age=18, active=false, level=1}
9. [1] TestPerson{id=P009, level=6, age=40, active=true, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, status=active, name=X, tags=temp, department=intern, level=1, age=22, salary=28000, active=true, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, active=true, tags=senior, department=engineering, name=Bob, age=35, salary=75000, score=9.2, status=active, level=5}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, status=inactive, department=sales, name=Eve, salary=55000, score=8, tags=employee, level=3}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, department=engineering, level=6}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, status=active, name=X, tags=temp, department=intern, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, score=6, status=inactive, department=hr, name=Charlie, active=false, tags=intern, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, name=Diana, age=45, salary=85000, tags=manager, status=active, department=marketing, level=7}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, department=management, name=Grace, age=65, score=10, level=9, salary=95000, active=true, tags=executive, status=active}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, score=5.5, tags=junior, status=inactive, department=support, name=Henry, age=18, active=false, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, amount=2, total=1999.98, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, status=confirmed, discount=0, region=south, date=2024-01-20, priority=low, customer_id=P002}
3. [1] TestOrder{id=O003, total=225, date=2024-02-01, status=shipped, customer_id=P001, amount=3, priority=high, discount=15, region=north, product_id=PROD003}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, priority=normal, discount=0, region=east, status=delivered}
5. [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, priority=high, region=south, total=999.99, discount=100, customer_id=P002, product_id=PROD001, amount=1}
6. [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, region=west, customer_id=P005, date=2024-02-15, discount=0, product_id=PROD005, amount=2}
7. [1] TestOrder{id=O007, status=shipped, discount=50, product_id=PROD006, amount=4, total=600, date=2024-03-01, priority=urgent, region=north, customer_id=P007}
8. [1] TestOrder{id=O008, date=2024-03-05, status=pending, priority=normal, customer_id=P010, discount=0, region=south, product_id=PROD002, amount=10, total=255}
9. [1] TestOrder{id=O009, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, amount=1, status=completed}
10. [1] TestOrder{id=O010, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006, total=75000, status=refunded, product_id=PROD001, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, amount=2, total=1999.98, status=pending}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, status=confirmed, priority=high, region=south, total=999.99, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006, total=75000, status=refunded}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, total=25.5, status=confirmed, discount=0, region=south, date=2024-01-20, priority=low}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, amount=3, priority=high, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, region=east, status=delivered, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, priority=normal, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, discount=50, product_id=PROD006, amount=4, total=600, date=2024-03-01, priority=urgent, region=north, customer_id=P007, status=shipped}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, customer_id=P010, discount=0, region=south}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, status=completed, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, amount=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active, department=sales}
2. [1] TestPerson{id=P002, active=true, tags=senior, department=engineering, name=Bob, age=35, salary=75000, score=9.2, status=active, level=5}
3. [1] TestPerson{id=P003, name=Charlie, active=false, tags=intern, level=1, age=16, salary=0, score=6, status=inactive, department=hr}
4. [1] TestPerson{id=P004, department=marketing, level=7, active=true, score=7.8, name=Diana, age=45, salary=85000, tags=manager, status=active}
5. [1] TestPerson{id=P005, level=3, age=30, active=false, status=inactive, department=sales, name=Eve, salary=55000, score=8, tags=employee}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, active=true, tags=test, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, level=9, salary=95000, active=true, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, name=Henry, age=18, active=false, level=1, salary=25000, score=5.5}
9. [1] TestPerson{id=P009, level=6, age=40, active=true, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, status=active, name=X, tags=temp, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, level=5, active=true, tags=senior, department=engineering, name=Bob, age=35, salary=75000, score=9.2, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, tags=manager, status=active, department=marketing, level=7, active=true, score=7.8, name=Diana, age=45}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, active=false, status=inactive, department=sales, name=Eve, salary=55000, score=8, tags=employee, level=3, age=30}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, active=true, tags=executive, status=active, department=management, name=Grace, age=65, score=10, level=9}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, active=true}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active, department=sales, name=Alice, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, amount=2, total=1999.98, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north}
2. [1] TestOrder{id=O002, priority=low, customer_id=P002, product_id=PROD002, amount=1, total=25.5, status=confirmed, discount=0, region=south, date=2024-01-20}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, priority=high, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped}
4. [1] TestOrder{id=O004, amount=1, total=299.99, date=2024-02-05, priority=normal, discount=0, region=east, status=delivered, customer_id=P004, product_id=PROD004}
5. [1] TestOrder{id=O005, priority=high, region=south, total=999.99, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, region=west, customer_id=P005, date=2024-02-15, discount=0}
7. [1] TestOrder{id=O007, region=north, customer_id=P007, status=shipped, discount=50, product_id=PROD006, amount=4, total=600, date=2024-03-01, priority=urgent}
8. [1] TestOrder{id=O008, date=2024-03-05, status=pending, priority=normal, customer_id=P010, discount=0, region=south, product_id=PROD002, amount=10, total=255}
9. [1] TestOrder{id=O009, region=north, customer_id=P001, amount=1, status=completed, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, discount=10}
10. [1] TestOrder{id=O010, customer_id=P006, total=75000, status=refunded, product_id=PROD001, amount=1, date=2024-03-15, priority=urgent, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, amount=2, total=1999.98}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, total=25.5, status=confirmed, discount=0, region=south}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, status=shipped, customer_id=P001, amount=3, priority=high, discount=15, region=north}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, priority=urgent, region=north, customer_id=P007, status=shipped, discount=50}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O008, status=pending, priority=normal, customer_id=P010, discount=0, region=south, product_id=PROD002, amount=10, total=255, date=2024-03-05}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, amount=1, status=completed}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, date=2024-02-05, priority=normal, discount=0, region=east, status=delivered, customer_id=P004, product_id=PROD004, amount=1, total=299.99}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high, region=south, total=999.99, discount=100}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, region=west, customer_id=P005, date=2024-02-15, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, category=electronics, price=999.99, rating=4.5, brand=TechCorp, stock=50, available=true, keywords=computer, supplier=TechSupply, name=Laptop}
2. [1] TestProduct{id=PROD002, category=accessories, price=25.5, available=true, rating=4, keywords=peripheral, brand=TechCorp, name=Mouse, stock=200, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, keywords=typing, brand=KeyTech, stock=0, category=accessories, price=75, rating=3.5, supplier=KeySupply, name=Keyboard, available=false}
4. [1] TestProduct{id=PROD004, available=true, keywords=display, stock=30, name=Monitor, category=electronics, price=299.99, rating=4.8, brand=ScreenPro, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, supplier=OldSupply, name=OldKeyboard, category=accessories, price=8.5, available=false, keywords=obsolete, stock=0, rating=2, brand=OldTech}
6. [1] TestProduct{id=PROD006, keywords=sound, stock=75, price=150, available=true, brand=AudioMax, supplier=AudioSupply, name=Headphones, category=audio, rating=4.6}
7. [1] TestProduct{id=PROD007, stock=25, category=electronics, rating=3.8, keywords=video, supplier=CamSupply, name=Webcam, price=89.99, available=true, brand=CamTech}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD007, category=electronics, rating=3.8, keywords=video, supplier=CamSupply, name=Webcam, price=89.99, available=true, brand=CamTech, stock=25}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, rating=4.5, brand=TechCorp, stock=50, available=true, keywords=computer, supplier=TechSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD002, category=accessories, price=25.5, available=true, rating=4, keywords=peripheral, brand=TechCorp, name=Mouse, stock=200, supplier=TechSupply}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD003, supplier=KeySupply, name=Keyboard, available=false, keywords=typing, brand=KeyTech, stock=0, category=accessories, price=75, rating=3.5}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, category=electronics, price=299.99, rating=4.8, brand=ScreenPro, supplier=ScreenSupply, available=true, keywords=display, stock=30}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD006, price=150, available=true, brand=AudioMax, supplier=AudioSupply, name=Headphones, category=audio, rating=4.6, keywords=sound, stock=75}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active, department=sales, name=Alice, active=true}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, name=Charlie, active=false, tags=intern, level=1, age=16, salary=0, score=6, status=inactive, department=hr}
4. [1] TestPerson{id=P004, status=active, department=marketing, level=7, active=true, score=7.8, name=Diana, age=45, salary=85000, tags=manager}
5. [1] TestPerson{id=P005, score=8, tags=employee, level=3, age=30, active=false, status=inactive, department=sales, name=Eve, salary=55000}
6. [1] TestPerson{id=P006, tags=test, level=1, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, active=true}
7. [1] TestPerson{id=P007, level=9, salary=95000, active=true, tags=executive, status=active, department=management, name=Grace, age=65, score=10}
8. [1] TestPerson{id=P008, salary=25000, score=5.5, tags=junior, status=inactive, department=support, name=Henry, age=18, active=false, level=1}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, active=true}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, status=active, name=X, tags=temp, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, level=1, age=16, salary=0, score=6, status=inactive, department=hr, name=Charlie, active=false}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, active=true, tags=test, level=1, name=Frank, age=0, salary=-5000, score=0, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, amount=2, total=1999.98, status=pending}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, total=25.5, status=confirmed, discount=0, region=south, date=2024-01-20, priority=low}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, priority=high, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, priority=normal, discount=0, region=east, status=delivered}
5. [1] TestOrder{id=O005, priority=high, region=south, total=999.99, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed}
6. [1] TestOrder{id=O006, priority=low, region=west, customer_id=P005, date=2024-02-15, discount=0, product_id=PROD005, amount=2, total=999.98, status=cancelled}
7. [1] TestOrder{id=O007, region=north, customer_id=P007, status=shipped, discount=50, product_id=PROD006, amount=4, total=600, date=2024-03-01, priority=urgent}
8. [1] TestOrder{id=O008, customer_id=P010, discount=0, region=south, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending, priority=normal}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, amount=1, status=completed, product_id=PROD007}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006, total=75000, status=refunded}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, status=completed, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, discount=10, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O010, total=75000, status=refunded, product_id=PROD001, amount=1, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O001, discount=50, region=north, amount=2, total=1999.98, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, priority=normal}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O003, status=shipped, customer_id=P001, amount=3, priority=high, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high, region=south, total=999.99, discount=100}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O006, date=2024-02-15, discount=0, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, region=west, customer_id=P005}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, total=600, date=2024-03-01, priority=urgent, region=north, customer_id=P007, status=shipped, discount=50, product_id=PROD006, amount=4}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, total=25.5, status=confirmed, discount=0, region=south, date=2024-01-20, priority=low}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O004, date=2024-02-05, priority=normal, discount=0, region=east, status=delivered, customer_id=P004, product_id=PROD004, amount=1, total=299.99}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, customer_id=P010, discount=0, region=south}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active, department=sales}
2. [1] TestPerson{id=P002, active=true, tags=senior, department=engineering, name=Bob, age=35, salary=75000, score=9.2, status=active, level=5}
3. [1] TestPerson{id=P003, name=Charlie, active=false, tags=intern, level=1, age=16, salary=0, score=6, status=inactive, department=hr}
4. [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, tags=manager, status=active, department=marketing, level=7, active=true, score=7.8}
5. [1] TestPerson{id=P005, tags=employee, level=3, age=30, active=false, status=inactive, department=sales, name=Eve, salary=55000, score=8}
6. [1] TestPerson{id=P006, age=0, salary=-5000, score=0, status=active, department=qa, active=true, tags=test, level=1, name=Frank}
7. [1] TestPerson{id=P007, salary=95000, active=true, tags=executive, status=active, department=management, name=Grace, age=65, score=10, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, level=1, salary=25000, score=5.5, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, active=true}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, status=active, name=X, tags=temp, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, level=1, age=16, salary=0, score=6, status=inactive, department=hr, name=Charlie, active=false}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, status=inactive, department=sales, name=Eve, salary=55000, score=8, tags=employee, level=3}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, name=Henry, age=18, active=false, level=1, salary=25000, score=5.5}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, amount=2, total=1999.98, status=pending}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, total=25.5, status=confirmed, discount=0, region=south, date=2024-01-20, priority=low}
3. [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, status=shipped, customer_id=P001, amount=3, priority=high, discount=15, region=north}
4. [1] TestOrder{id=O004, region=east, status=delivered, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, priority=normal, discount=0}
5. [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, priority=high, region=south, total=999.99, discount=100, customer_id=P002, product_id=PROD001, amount=1}
6. [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, region=west, customer_id=P005, date=2024-02-15, discount=0, product_id=PROD005, amount=2}
7. [1] TestOrder{id=O007, date=2024-03-01, priority=urgent, region=north, customer_id=P007, status=shipped, discount=50, product_id=PROD006, amount=4, total=600}
8. [1] TestOrder{id=O008, customer_id=P010, discount=0, region=south, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending, priority=normal}
9. [1] TestOrder{id=O009, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, amount=1, status=completed}
10. [1] TestOrder{id=O010, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006, total=75000, status=refunded, product_id=PROD001, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, priority=normal, discount=0, region=east, status=delivered}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, priority=high, region=south, total=999.99, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, discount=0, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, region=west}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, total=600, date=2024-03-01, priority=urgent, region=north, customer_id=P007, status=shipped, discount=50, product_id=PROD006, amount=4}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, status=completed, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, amount=1}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006, total=75000, status=refunded}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, total=25.5, status=confirmed, discount=0, region=south, date=2024-01-20, priority=low}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, date=2024-02-01, status=shipped, customer_id=P001, amount=3, priority=high, discount=15, region=north, product_id=PROD003, total=225}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, name=Alice, active=true, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active}
2. [1] TestPerson{id=P002, department=engineering, name=Bob, age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, active=false, tags=intern, level=1, age=16, salary=0, score=6, status=inactive, department=hr}
4. [1] TestPerson{id=P004, level=7, active=true, score=7.8, name=Diana, age=45, salary=85000, tags=manager, status=active, department=marketing}
5. [1] TestPerson{id=P005, age=30, active=false, status=inactive, department=sales, name=Eve, salary=55000, score=8, tags=employee, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, active=true, tags=test, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, level=9, salary=95000, active=true, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, level=1, salary=25000, score=5.5, tags=junior, status=inactive, department=support, name=Henry, age=18, active=false}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, active=true}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, status=active, name=X, tags=temp, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, tags=manager, status=active, department=marketing, level=7, active=true, score=7.8}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, tags=employee, level=3, age=30, active=false, status=inactive, department=sales}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, level=1, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, active=true}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, score=10, level=9, salary=95000, active=true, tags=executive, status=active, department=management}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, age=18, active=false, level=1, salary=25000, score=5.5, tags=junior, status=inactive, department=support, name=Henry}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, name=Bob, age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, score=6, status=inactive, department=hr, name=Charlie, active=false, tags=intern, level=1}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, active=true, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P001, level=2, age=25, salary=45000, status=active, department=sales, name=Alice, active=true, score=8.5, tags=junior}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, stock=50, available=true, keywords=computer, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, rating=4.5, brand=TechCorp}
2. [1] TestProduct{id=PROD002, category=accessories, price=25.5, available=true, rating=4, keywords=peripheral, brand=TechCorp, name=Mouse, stock=200, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, available=false, keywords=typing, brand=KeyTech, stock=0, category=accessories, price=75, rating=3.5, supplier=KeySupply, name=Keyboard}
4. [1] TestProduct{id=PROD004, keywords=display, stock=30, name=Monitor, category=electronics, price=299.99, rating=4.8, brand=ScreenPro, supplier=ScreenSupply, available=true}
5. [1] TestProduct{id=PROD005, keywords=obsolete, stock=0, rating=2, brand=OldTech, supplier=OldSupply, name=OldKeyboard, category=accessories, price=8.5, available=false}
6. [1] TestProduct{id=PROD006, name=Headphones, category=audio, rating=4.6, keywords=sound, stock=75, price=150, available=true, brand=AudioMax, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, name=Webcam, price=89.99, available=true, brand=CamTech, stock=25, category=electronics, rating=3.8, keywords=video, supplier=CamSupply}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD006, rating=4.6, keywords=sound, stock=75, price=150, available=true, brand=AudioMax, supplier=AudioSupply, name=Headphones, category=audio}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD007, rating=3.8, keywords=video, supplier=CamSupply, name=Webcam, price=89.99, available=true, brand=CamTech, stock=25, category=electronics}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD001, category=electronics, price=999.99, rating=4.5, brand=TechCorp, stock=50, available=true, keywords=computer, supplier=TechSupply, name=Laptop}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, stock=200, supplier=TechSupply, category=accessories, price=25.5, available=true, rating=4, keywords=peripheral, brand=TechCorp}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD003, name=Keyboard, available=false, keywords=typing, brand=KeyTech, stock=0, category=accessories, price=75, rating=3.5, supplier=KeySupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD004, keywords=display, stock=30, name=Monitor, category=electronics, price=299.99, rating=4.8, brand=ScreenPro, supplier=ScreenSupply, available=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, level=2, age=25, salary=45000, status=active, department=sales, name=Alice, active=true, score=8.5}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, tags=intern, level=1, age=16, salary=0, score=6, status=inactive, department=hr, name=Charlie, active=false}
4. [1] TestPerson{id=P004, salary=85000, tags=manager, status=active, department=marketing, level=7, active=true, score=7.8, name=Diana, age=45}
5. [1] TestPerson{id=P005, age=30, active=false, status=inactive, department=sales, name=Eve, salary=55000, score=8, tags=employee, level=3}
6. [1] TestPerson{id=P006, active=true, tags=test, level=1, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, level=9, salary=95000, active=true, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, status=inactive, department=support, name=Henry, age=18, active=false, level=1, salary=25000, score=5.5, tags=junior}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, active=true}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, status=active, name=X, tags=temp, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, level=9, salary=95000, active=true, tags=executive, status=active, department=management, name=Grace, age=65, score=10}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, level=1, salary=25000, score=5.5, tags=junior, status=inactive, department=support}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, active=true}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, department=engineering}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P003, level=1, age=16, salary=0, score=6, status=inactive, department=hr, name=Charlie, active=false, tags=intern}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, name=Diana, age=45, salary=85000, tags=manager, status=active, department=marketing, level=7, active=true}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active, department=sales}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, score=8, tags=employee, level=3, age=30, active=false, status=inactive, department=sales, name=Eve}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P006, level=1, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, active=true, tags=test}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, total=1999.98, status=pending, customer_id=P001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, amount=2}
2. [1] TestOrder{id=O002, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, total=25.5, status=confirmed, discount=0, region=south}
3. [1] TestOrder{id=O003, status=shipped, customer_id=P001, amount=3, priority=high, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, priority=normal, discount=0, region=east, status=delivered}
5. [1] TestOrder{id=O005, priority=high, region=south, total=999.99, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed}
6. [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, discount=0, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, region=west}
7. [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, priority=urgent, region=north, customer_id=P007, status=shipped, discount=50, product_id=PROD006}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, customer_id=P010, discount=0, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, status=completed, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, discount=10, region=north}
10. [1] TestOrder{id=O010, status=refunded, product_id=PROD001, amount=1, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006, total=75000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, priority=normal, discount=0, region=east, status=delivered}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high, region=south, total=999.99, discount=100}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, status=shipped, discount=50, product_id=PROD006, amount=4, total=600, date=2024-03-01, priority=urgent, region=north, customer_id=P007}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, discount=0, region=south, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending, priority=normal}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, amount=1, status=completed}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, amount=2, total=1999.98, status=pending}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, total=25.5, status=confirmed, discount=0, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, priority=high, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped, customer_id=P001, amount=3}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active, department=sales, name=Alice}
2. [1] TestPerson{id=P002, active=true, tags=senior, department=engineering, name=Bob, age=35, salary=75000, score=9.2, status=active, level=5}
3. [1] TestPerson{id=P003, name=Charlie, active=false, tags=intern, level=1, age=16, salary=0, score=6, status=inactive, department=hr}
4. [1] TestPerson{id=P004, age=45, salary=85000, tags=manager, status=active, department=marketing, level=7, active=true, score=7.8, name=Diana}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, score=8, tags=employee, level=3, age=30, active=false, status=inactive, department=sales}
6. [1] TestPerson{id=P006, score=0, status=active, department=qa, active=true, tags=test, level=1, name=Frank, age=0, salary=-5000}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, level=9, salary=95000, active=true, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, salary=25000, score=5.5, tags=junior, status=inactive, department=support, name=Henry, age=18, active=false, level=1}
9. [1] TestPerson{id=P009, department=engineering, level=6, age=40, active=true, name=Ivy, salary=68000, score=8.7, tags=senior, status=active}
10. [1] TestPerson{id=P010, name=X, tags=temp, department=intern, level=1, age=22, salary=28000, active=true, score=6.5, status=active}
11. [1] TestOrder{id=O001, date=2024-01-15, priority=normal, discount=50, region=north, amount=2, total=1999.98, status=pending, customer_id=P001, product_id=PROD001}
12. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, total=25.5, status=confirmed, discount=0, region=south, date=2024-01-20, priority=low}
13. [1] TestOrder{id=O003, amount=3, priority=high, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped, customer_id=P001}
14. [1] TestOrder{id=O004, status=delivered, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, priority=normal, discount=0, region=east}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high, region=south, total=999.99, discount=100}
16. [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, discount=0, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, region=west}
17. [1] TestOrder{id=O007, total=600, date=2024-03-01, priority=urgent, region=north, customer_id=P007, status=shipped, discount=50, product_id=PROD006, amount=4}
18. [1] TestOrder{id=O008, total=255, date=2024-03-05, status=pending, priority=normal, customer_id=P010, discount=0, region=south, product_id=PROD002, amount=10}
19. [1] TestOrder{id=O009, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, amount=1, status=completed}
20. [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006, total=75000, status=refunded}

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

1. [1] TestPerson{id=P001, age=25, salary=45000, status=active, department=sales, name=Alice, active=true, score=8.5, tags=junior, level=2}
2. [1] TestPerson{id=P002, score=9.2, status=active, level=5, active=true, tags=senior, department=engineering, name=Bob, age=35, salary=75000}
3. [1] TestPerson{id=P003, department=hr, name=Charlie, active=false, tags=intern, level=1, age=16, salary=0, score=6, status=inactive}
4. [1] TestPerson{id=P004, tags=manager, status=active, department=marketing, level=7, active=true, score=7.8, name=Diana, age=45, salary=85000}
5. [1] TestPerson{id=P005, age=30, active=false, status=inactive, department=sales, name=Eve, salary=55000, score=8, tags=employee, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, active=true, tags=test, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, level=9, salary=95000, active=true, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, level=1, salary=25000, score=5.5, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, age=40, active=true, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, salary=28000, active=true, score=6.5, status=active, name=X, tags=temp, department=intern, level=1, age=22}

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

1. [1] TestPerson{id=P001, age=25, salary=45000, status=active, department=sales, name=Alice, active=true, score=8.5, tags=junior, level=2}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, level=1, age=16, salary=0, score=6, status=inactive, department=hr, name=Charlie, active=false, tags=intern}
4. [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, tags=manager, status=active, department=marketing, level=7, active=true, score=7.8}
5. [1] TestPerson{id=P005, age=30, active=false, status=inactive, department=sales, name=Eve, salary=55000, score=8, tags=employee, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, active=true, tags=test, level=1}
7. [1] TestPerson{id=P007, status=active, department=management, name=Grace, age=65, score=10, level=9, salary=95000, active=true, tags=executive}
8. [1] TestPerson{id=P008, salary=25000, score=5.5, tags=junior, status=inactive, department=support, name=Henry, age=18, active=false, level=1}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, active=true}
10. [1] TestPerson{id=P010, department=intern, level=1, age=22, salary=28000, active=true, score=6.5, status=active, name=X, tags=temp}
11. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, amount=2, total=1999.98, status=pending}
12. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, amount=1, total=25.5, status=confirmed, discount=0, region=south, date=2024-01-20, priority=low}
13. [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, status=shipped, customer_id=P001, amount=3, priority=high, discount=15, region=north}
14. [1] TestOrder{id=O004, total=299.99, date=2024-02-05, priority=normal, discount=0, region=east, status=delivered, customer_id=P004, product_id=PROD004, amount=1}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high, region=south, total=999.99, discount=100}
16. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, region=west, customer_id=P005, date=2024-02-15, discount=0}
17. [1] TestOrder{id=O007, priority=urgent, region=north, customer_id=P007, status=shipped, discount=50, product_id=PROD006, amount=4, total=600, date=2024-03-01}
18. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, status=pending, priority=normal, customer_id=P010, discount=0, region=south}
19. [1] TestOrder{id=O009, customer_id=P001, amount=1, status=completed, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, discount=10, region=north}
20. [1] TestOrder{id=O010, amount=1, date=2024-03-15, priority=urgent, discount=0, region=east, customer_id=P006, total=75000, status=refunded, product_id=PROD001}

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

1. [1] TestPerson{id=P001, age=25, salary=45000, status=active, department=sales, name=Alice, active=true, score=8.5, tags=junior, level=2}
2. [1] TestPerson{id=P002, level=5, active=true, tags=senior, department=engineering, name=Bob, age=35, salary=75000, score=9.2, status=active}
3. [1] TestPerson{id=P003, tags=intern, level=1, age=16, salary=0, score=6, status=inactive, department=hr, name=Charlie, active=false}
4. [1] TestPerson{id=P004, status=active, department=marketing, level=7, active=true, score=7.8, name=Diana, age=45, salary=85000, tags=manager}
5. [1] TestPerson{id=P005, tags=employee, level=3, age=30, active=false, status=inactive, department=sales, name=Eve, salary=55000, score=8}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa, active=true, tags=test, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, score=10, level=9, salary=95000, active=true, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, level=1, salary=25000, score=5.5, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, status=active, department=engineering, level=6, age=40, active=true}
10. [1] TestPerson{id=P010, status=active, name=X, tags=temp, department=intern, level=1, age=22, salary=28000, active=true, score=6.5}

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

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, priority=normal, discount=50, region=north, amount=2, total=1999.98, status=pending}
2. [1] TestOrder{id=O002, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, total=25.5, status=confirmed, discount=0, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, priority=high, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, priority=normal, discount=0, region=east, status=delivered}
5. [1] TestOrder{id=O005, status=confirmed, priority=high, region=south, total=999.99, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10}
6. [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, discount=0, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, region=west}
7. [1] TestOrder{id=O007, discount=50, product_id=PROD006, amount=4, total=600, date=2024-03-01, priority=urgent, region=north, customer_id=P007, status=shipped}
8. [1] TestOrder{id=O008, status=pending, priority=normal, customer_id=P010, discount=0, region=south, product_id=PROD002, amount=10, total=255, date=2024-03-05}
9. [1] TestOrder{id=O009, status=completed, product_id=PROD007, total=89.99, date=2024-03-10, priority=low, discount=10, region=north, customer_id=P001, amount=1}
10. [1] TestOrder{id=O010, priority=urgent, discount=0, region=east, customer_id=P006, total=75000, status=refunded, product_id=PROD001, amount=1, date=2024-03-15}

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

1. [1] TestPerson{id=P001, salary=45000, status=active, department=sales, name=Alice, active=true, score=8.5, tags=junior, level=2, age=25}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, name=Charlie, active=false, tags=intern, level=1, age=16, salary=0, score=6, status=inactive, department=hr}
4. [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, tags=manager, status=active, department=marketing, level=7, active=true, score=7.8}
5. [1] TestPerson{id=P005, salary=55000, score=8, tags=employee, level=3, age=30, active=false, status=inactive, department=sales, name=Eve}
6. [1] TestPerson{id=P006, active=true, tags=test, level=1, name=Frank, age=0, salary=-5000, score=0, status=active, department=qa}
7. [1] TestPerson{id=P007, salary=95000, active=true, tags=executive, status=active, department=management, name=Grace, age=65, score=10, level=9}
8. [1] TestPerson{id=P008, active=false, level=1, salary=25000, score=5.5, tags=junior, status=inactive, department=support, name=Henry, age=18}
9. [1] TestPerson{id=P009, status=active, department=engineering, level=6, age=40, active=true, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, name=X, tags=temp, department=intern, level=1, age=22, salary=28000, active=true, score=6.5, status=active}

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
