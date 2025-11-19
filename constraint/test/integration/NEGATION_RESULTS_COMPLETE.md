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

1. [1] TestPerson{id=P001, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering, age=35, salary=75000}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive}
4. [1] TestPerson{id=P004, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing, name=Diana}
5. [1] TestPerson{id=P005, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3, name=Eve, salary=55000}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active}
10. [1] TestPerson{id=P010, department=intern, level=1, age=22, tags=temp, status=active, name=X, salary=28000, active=true, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales, name=Alice, age=25}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing, name=Diana, salary=85000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active, name=X, salary=28000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, level=1, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, status=active, department=management, level=9, salary=95000, score=10, tags=executive, name=Grace, age=65, active=true}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, tags=senior, status=active, name=Ivy, salary=68000, score=8.7, department=engineering, level=6}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15}
2. [1] TestOrder{id=O002, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, customer_id=P002, date=2024-01-20}
3. [1] TestOrder{id=O003, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high, product_id=PROD003, amount=3, total=225}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100}
6. [1] TestOrder{id=O006, amount=2, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west}
7. [1] TestOrder{id=O007, discount=50, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped}
8. [1] TestOrder{id=O008, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10}
9. [1] TestOrder{id=O009, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}
10. [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, customer_id=P002, date=2024-01-20, status=confirmed, priority=low}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive}
4. [1] TestPerson{id=P004, active=true, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45}
5. [1] TestPerson{id=P005, status=inactive, level=3, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8}
6. [1] TestPerson{id=P006, tags=test, status=active, level=1, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000}
7. [1] TestPerson{id=P007, salary=95000, score=10, tags=executive, name=Grace, age=65, active=true, status=active, department=management, level=9}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1}
9. [1] TestPerson{id=P009, tags=senior, status=active, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, level=3, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, active=true, tags=senior, status=active, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98}
2. [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}
3. [1] TestOrder{id=O003, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high, product_id=PROD003, amount=3}
4. [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal, discount=0}
5. [1] TestOrder{id=O005, priority=high, region=south, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10}
6. [1] TestOrder{id=O006, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}
8. [1] TestOrder{id=O008, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10}
9. [1] TestOrder{id=O009, status=completed, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10}
10. [1] TestOrder{id=O010, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O009, status=completed, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low, customer_id=P005}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O001, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, rating=4.5, brand=TechCorp, stock=50, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, available=true, keywords=computer}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, price=25.5, supplier=TechSupply, available=true, rating=4, keywords=peripheral, brand=TechCorp, stock=200}
3. [1] TestProduct{id=PROD003, category=accessories, available=false, rating=3.5, keywords=typing, brand=KeyTech, supplier=KeySupply, name=Keyboard, price=75, stock=0}
4. [1] TestProduct{id=PROD004, keywords=display, rating=4.8, brand=ScreenPro, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics, price=299.99, available=true}
5. [1] TestProduct{id=PROD005, keywords=obsolete, stock=0, name=OldKeyboard, category=accessories, price=8.5, rating=2, brand=OldTech, supplier=OldSupply, available=false}
6. [1] TestProduct{id=PROD006, category=audio, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply, price=150, available=true, rating=4.6, name=Headphones}
7. [1] TestProduct{id=PROD007, available=true, rating=3.8, keywords=video, stock=25, brand=CamTech, supplier=CamSupply, name=Webcam, category=electronics, price=89.99}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, available=true, keywords=computer, rating=4.5, brand=TechCorp, stock=50, supplier=TechSupply}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, available=true, rating=4, keywords=peripheral, brand=TechCorp, stock=200, name=Mouse, category=accessories, price=25.5, supplier=TechSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, brand=KeyTech, supplier=KeySupply, name=Keyboard, price=75, stock=0, category=accessories, available=false, rating=3.5, keywords=typing}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, category=electronics, price=299.99, available=true, keywords=display, rating=4.8, brand=ScreenPro, stock=30, supplier=ScreenSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, rating=4.6, name=Headphones, category=audio, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply, price=150, available=true}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, brand=CamTech, supplier=CamSupply, name=Webcam, category=electronics, price=89.99, available=true, rating=3.8, keywords=video, stock=25}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, status=active, level=2, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5}
2. [1] TestPerson{id=P002, level=5, name=Bob, department=engineering, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive}
4. [1] TestPerson{id=P004, age=45, active=true, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7}
5. [1] TestPerson{id=P005, status=inactive, level=3, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8}
6. [1] TestPerson{id=P006, active=true, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank}
7. [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}
8. [1] TestPerson{id=P008, level=1, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support}
9. [1] TestPerson{id=P009, tags=senior, status=active, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true}
10. [1] TestPerson{id=P010, age=22, tags=temp, status=active, name=X, salary=28000, active=true, score=6.5, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true, score=0, department=qa}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, level=1, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north, customer_id=P001}
2. [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}
3. [1] TestOrder{id=O003, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high, product_id=PROD003, amount=3}
4. [1] TestOrder{id=O004, region=east, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100}
6. [1] TestOrder{id=O006, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}
10. [1] TestOrder{id=O010, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering, age=35}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1}
7. [1] TestPerson{id=P007, score=10, tags=executive, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000}
8. [1] TestPerson{id=P008, status=inactive, age=18, department=support, level=1, name=Henry, salary=25000, active=false, score=5.5, tags=junior}
9. [1] TestPerson{id=P009, active=true, tags=senior, status=active, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40}
10. [1] TestPerson{id=P010, age=22, tags=temp, status=active, name=X, salary=28000, active=true, score=6.5, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, department=support, level=1, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive, name=Charlie}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, status=inactive, level=3, name=Eve, salary=55000, tags=employee, department=sales, age=30}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north}
2. [1] TestOrder{id=O002, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, customer_id=P002, date=2024-01-20, status=confirmed}
3. [1] TestOrder{id=O003, customer_id=P001, status=shipped, priority=high, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal, discount=0}
5. [1] TestOrder{id=O005, priority=high, region=south, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10}
6. [1] TestOrder{id=O006, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}
9. [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007}
10. [1] TestOrder{id=O010, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O005, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, customer_id=P001, priority=low, discount=10}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O010, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O003, priority=high, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active}
2. [1] TestPerson{id=P002, department=engineering, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob}
3. [1] TestPerson{id=P003, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, active=true, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45}
5. [1] TestPerson{id=P005, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3, name=Eve, salary=55000}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}
8. [1] TestPerson{id=P008, department=support, level=1, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, tags=senior, status=active, name=Ivy, salary=68000, score=8.7, department=engineering, level=6}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, status=inactive, level=3, name=Eve, salary=55000, tags=employee, department=sales}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, status=active, level=1, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, name=Henry}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering, age=35, salary=75000}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing, name=Diana, salary=85000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, keywords=computer, rating=4.5, brand=TechCorp, stock=50, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, available=true}
2. [1] TestProduct{id=PROD002, available=true, rating=4, keywords=peripheral, brand=TechCorp, stock=200, name=Mouse, category=accessories, price=25.5, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, category=accessories, available=false, rating=3.5, keywords=typing, brand=KeyTech, supplier=KeySupply, name=Keyboard, price=75, stock=0}
4. [1] TestProduct{id=PROD004, rating=4.8, brand=ScreenPro, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics, price=299.99, available=true, keywords=display}
5. [1] TestProduct{id=PROD005, price=8.5, rating=2, brand=OldTech, supplier=OldSupply, available=false, keywords=obsolete, stock=0, name=OldKeyboard, category=accessories}
6. [1] TestProduct{id=PROD006, name=Headphones, category=audio, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply, price=150, available=true, rating=4.6}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, price=89.99, available=true, rating=3.8, keywords=video, stock=25, brand=CamTech, supplier=CamSupply}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD003, category=accessories, available=false, rating=3.5, keywords=typing, brand=KeyTech, supplier=KeySupply, name=Keyboard, price=75, stock=0}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, category=electronics, price=299.99, available=true, keywords=display, rating=4.8, brand=ScreenPro, stock=30, supplier=ScreenSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD006, stock=75, supplier=AudioSupply, price=150, available=true, rating=4.6, name=Headphones, category=audio, keywords=sound, brand=AudioMax}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD007, brand=CamTech, supplier=CamSupply, name=Webcam, category=electronics, price=89.99, available=true, rating=3.8, keywords=video, stock=25}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD001, supplier=TechSupply, name=Laptop, category=electronics, price=999.99, available=true, keywords=computer, rating=4.5, brand=TechCorp, stock=50}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD002, stock=200, name=Mouse, category=accessories, price=25.5, supplier=TechSupply, available=true, rating=4, keywords=peripheral, brand=TechCorp}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, tags=junior, status=active, level=2, active=true, department=sales, name=Alice, age=25, salary=45000}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering}
3. [1] TestPerson{id=P003, level=1, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr}
4. [1] TestPerson{id=P004, status=active, level=7, age=45, active=true, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3, name=Eve}
6. [1] TestPerson{id=P006, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true}
7. [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active}
10. [1] TestPerson{id=P010, status=active, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, level=2, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, active=true, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3, name=Eve}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true, score=0, department=qa}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive, name=Grace}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending}
2. [1] TestOrder{id=O002, region=south, product_id=PROD002, amount=1, total=25.5, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0}
3. [1] TestOrder{id=O003, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high, product_id=PROD003, amount=3}
4. [1] TestOrder{id=O004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004}
5. [1] TestOrder{id=O005, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south}
6. [1] TestOrder{id=O006, amount=2, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}
9. [1] TestOrder{id=O009, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, customer_id=P001, priority=low}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high, product_id=PROD003, amount=3, total=225}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal, discount=0}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, customer_id=P001, priority=low, discount=10, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, status=active, level=5, name=Bob, department=engineering, age=35, salary=75000, active=true, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive}
4. [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing, name=Diana, salary=85000}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3}
6. [1] TestPerson{id=P006, level=1, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test, status=active}
7. [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}
11. [1] TestOrder{id=O001, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending}
12. [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south}
13. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high}
14. [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal, discount=0}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100}
16. [1] TestOrder{id=O006, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled}
17. [1] TestOrder{id=O007, status=shipped, discount=50, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01}
18. [1] TestOrder{id=O008, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255}
19. [1] TestOrder{id=O009, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, customer_id=P001, priority=low}
20. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, level=7, age=45, active=true, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager, status=active}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, level=2, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100, customer_id=P002}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O007, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, status=active, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, customer_id=P001, priority=low, discount=10, region=north}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3, name=Eve}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true, score=0, department=qa}
   - Fait 2: [1] TestOrder{id=O001, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, status=active, level=2, active=true, department=sales, name=Alice, age=25, salary=45000}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales}
   - Fait 2: [1] TestOrder{id=O006, region=west, amount=2, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true, score=0, department=qa}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south, customer_id=P010}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive, name=Charlie, age=16}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal, discount=0}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, level=1, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false}
   - Fait 2: [1] TestOrder{id=O009, status=completed, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P006, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P001, status=active, level=2, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, status=inactive, level=3, name=Eve, salary=55000, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100, customer_id=P002}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true, score=0}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P001, status=active, level=2, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent, customer_id=P007}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P008, department=support, level=1, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P010, level=1, age=22, tags=temp, status=active, name=X, salary=28000, active=true, score=6.5, department=intern}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P007, level=9, salary=95000, score=10, tags=executive, name=Grace, age=65, active=true, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, level=6, age=40, active=true, tags=senior, status=active, name=Ivy, salary=68000, score=8.7}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, status=inactive, level=3, name=Eve, salary=55000, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high, product_id=PROD003}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal, discount=0}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100, customer_id=P002}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering, age=35, salary=75000}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive, name=Grace}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, level=6, age=40, active=true, tags=senior, status=active, name=Ivy, salary=68000, score=8.7}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, status=shipped, priority=high, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, department=intern, level=1, age=22, tags=temp, status=active, name=X, salary=28000, active=true}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000}
   - Fait 2: [1] TestOrder{id=O004, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering}
   - Fait 2: [1] TestOrder{id=O004, region=east, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7}
   - Fait 2: [1] TestOrder{id=O006, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering, age=35}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active, name=X, salary=28000}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P007, department=management, level=9, salary=95000, score=10, tags=executive, name=Grace, age=65, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal, discount=0}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100, customer_id=P002}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, department=engineering, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent, customer_id=P007}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales, name=Alice}
   - Fait 2: [1] TestOrder{id=O002, amount=1, total=25.5, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7}
   - Fait 2: [1] TestOrder{id=O005, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true, score=0, department=qa}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, customer_id=P001, priority=low, discount=10, region=north}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing, name=Diana}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active, name=X, salary=28000}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P007, level=9, salary=95000, score=10, tags=executive, name=Grace, age=65, active=true, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, level=1, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P001, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, name=Henry, salary=25000}
   - Fait 2: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales, name=Alice, age=25}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true, score=0, department=qa}
   - Fait 2: [1] TestOrder{id=O002, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, customer_id=P002, date=2024-01-20, status=confirmed, priority=low}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, level=1, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false, tags=intern}
   - Fait 2: [1] TestOrder{id=O004, region=east, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P001, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering, age=35, salary=75000}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P003, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true, score=0, department=qa}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing, name=Diana, salary=85000}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, age=30, active=false, score=8, status=inactive, level=3, name=Eve, salary=55000, tags=employee}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active, name=Ivy, salary=68000}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true}
   - Fait 2: [1] TestOrder{id=O006, amount=2, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, score=10, tags=executive, name=Grace, age=65, active=true, status=active, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O006, amount=2, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P002, level=5, name=Bob, department=engineering, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O003, status=shipped, priority=high, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P007, department=management, level=9, salary=95000, score=10, tags=executive, name=Grace, age=65, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, amount=2, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, amount=2, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3, name=Eve}
   - Fait 2: [1] TestOrder{id=O008, discount=0, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, name=Henry, salary=25000, active=false}
   - Fait 2: [1] TestOrder{id=O003, discount=15, region=north, customer_id=P001, status=shipped, priority=high, product_id=PROD003, amount=3, total=225, date=2024-02-01}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, status=inactive, level=3, name=Eve, salary=55000, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, name=Henry, salary=25000, active=false}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P004, active=true, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, name=Henry, salary=25000, active=false}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, customer_id=P001, priority=low, discount=10, region=north}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering, age=35, salary=75000}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive}
4. [1] TestPerson{id=P004, active=true, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45}
5. [1] TestPerson{id=P005, level=3, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive}
6. [1] TestPerson{id=P006, status=active, level=1, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test}
7. [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}
8. [1] TestPerson{id=P008, status=inactive, age=18, department=support, level=1, name=Henry, salary=25000, active=false, score=5.5, tags=junior}
9. [1] TestPerson{id=P009, department=engineering, level=6, age=40, active=true, tags=senior, status=active, name=Ivy, salary=68000, score=8.7}
10. [1] TestPerson{id=P010, status=active, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, level=6, age=40, active=true, tags=senior, status=active, name=Ivy, salary=68000, score=8.7, department=engineering}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, status=active, level=1, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering}
3. [1] TestPerson{id=P003, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, tags=manager, status=active, level=7, age=45, active=true, department=marketing, name=Diana, salary=85000, score=7.8}
5. [1] TestPerson{id=P005, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3, name=Eve}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}
8. [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, name=Henry}
9. [1] TestPerson{id=P009, age=40, active=true, tags=senior, status=active, name=Ivy, salary=68000, score=8.7, department=engineering, level=6}
10. [1] TestPerson{id=P010, age=22, tags=temp, status=active, name=X, salary=28000, active=true, score=6.5, department=intern, level=1}
11. [1] TestOrder{id=O001, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98}
12. [1] TestOrder{id=O002, amount=1, total=25.5, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002}
13. [1] TestOrder{id=O003, customer_id=P001, status=shipped, priority=high, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north}
14. [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east}
15. [1] TestOrder{id=O005, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001}
16. [1] TestOrder{id=O006, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2}
17. [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006}
18. [1] TestOrder{id=O008, discount=0, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending}
19. [1] TestOrder{id=O009, date=2024-03-10, status=completed, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99}
20. [1] TestOrder{id=O010, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering, age=35, salary=75000}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, status=active, level=5, name=Bob, department=engineering, age=35, salary=75000, active=true, score=9.2, tags=senior}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high, product_id=PROD003}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7}
   - Fait 2: [1] TestOrder{id=O010, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O002, region=south, product_id=PROD002, amount=1, total=25.5, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive, name=Grace}
   - Fait 2: [1] TestOrder{id=O001, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north, customer_id=P001}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7}
   - Fait 2: [1] TestOrder{id=O006, amount=2, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O007, total=600, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south, customer_id=P010}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P007, status=active, department=management, level=9, salary=95000, score=10, tags=executive, name=Grace, age=65, active=true}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true, score=0, department=qa}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales}
   - Fait 2: [1] TestOrder{id=O001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north, customer_id=P001}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, level=1, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false}
   - Fait 2: [1] TestOrder{id=O001, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, tags=senior, status=active, name=Ivy, salary=68000, score=8.7, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, status=inactive, level=3, name=Eve, salary=55000, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, status=inactive, level=3, name=Eve, salary=55000, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales}
   - Fait 2: [1] TestOrder{id=O008, status=pending, discount=0, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P001, status=active, level=2, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales, name=Alice}
   - Fait 2: [1] TestOrder{id=O002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, customer_id=P002}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P003, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, customer_id=P002}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, tags=senior, status=active, name=Ivy, salary=68000, score=8.7, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O003, discount=15, region=north, customer_id=P001, status=shipped, priority=high, product_id=PROD003, amount=3, total=225, date=2024-02-01}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P005, score=8, status=inactive, level=3, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P003, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, customer_id=P001}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, level=6, age=40, active=true, tags=senior, status=active, name=Ivy, salary=68000, score=8.7}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing, name=Diana, salary=85000}
   - Fait 2: [1] TestOrder{id=O001, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}
   - Fait 2: [1] TestOrder{id=O006, amount=2, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true, score=0, department=qa}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true, score=0, department=qa}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, department=marketing, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high, product_id=PROD003}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true, score=0, department=qa}
   - Fait 2: [1] TestOrder{id=O004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal, discount=0}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering, age=35, salary=75000, active=true}
   - Fait 2: [1] TestOrder{id=O010, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, status=inactive, level=3, name=Eve, salary=55000, tags=employee, department=sales, age=30}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active, name=X, salary=28000}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, status=shipped, priority=high, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P008, department=support, level=1, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18}
   - Fait 2: [1] TestOrder{id=O003, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high, product_id=PROD003, amount=3}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, level=1, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P001, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal, discount=0}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P010, age=22, tags=temp, status=active, name=X, salary=28000, active=true, score=6.5, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, level=3, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100, customer_id=P002}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P003, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales}
   - Fait 2: [1] TestOrder{id=O006, amount=2, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O008, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O001, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive, name=Charlie}
   - Fait 2: [1] TestOrder{id=O005, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, status=active, level=2, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, amount=2, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, name=Henry, salary=25000, active=false}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active, name=X, salary=28000}
   - Fait 2: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent, customer_id=P007, product_id=PROD006}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active, name=Ivy}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, status=shipped, priority=high, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, status=active, level=2, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O004, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, customer_id=P001}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1, score=6, status=inactive}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, score=10, tags=executive, name=Grace, age=65, active=true, status=active, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, status=shipped, priority=high, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal, discount=0}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, status=inactive, level=3, name=Eve, salary=55000, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true, score=0, department=qa}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P007, level=9, salary=95000, score=10, tags=executive, name=Grace, age=65, active=true, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O010, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P007, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive, name=Grace, age=65}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P002, level=5, name=Bob, department=engineering, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, department=engineering, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, status=inactive, level=3, name=Eve, salary=55000, tags=employee, department=sales, age=30}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active, name=Ivy, salary=68000}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal, discount=0}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O010, region=east, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P008, department=support, level=1, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P010, age=22, tags=temp, status=active, name=X, salary=28000, active=true, score=6.5, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P009, level=6, age=40, active=true, tags=senior, status=active, name=Ivy, salary=68000, score=8.7, department=engineering}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P009, status=active, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O008, discount=0, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P008, level=1, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, status=completed, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, level=5, name=Bob, department=engineering, age=35, salary=75000, active=true, score=9.2}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north, customer_id=P001, amount=2, date=2024-01-15, status=pending}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2}
2. [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering, age=35, salary=75000}
3. [1] TestPerson{id=P003, active=false, tags=intern, department=hr, level=1, score=6, status=inactive, name=Charlie, age=16, salary=0}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3}
6. [1] TestPerson{id=P006, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true, score=0, department=qa, age=0}
7. [1] TestPerson{id=P007, department=management, level=9, salary=95000, score=10, tags=executive, name=Grace, age=65, active=true, status=active}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active}
10. [1] TestPerson{id=P010, age=22, tags=temp, status=active, name=X, salary=28000, active=true, score=6.5, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, department=intern, level=1, age=22, tags=temp, status=active, name=X, salary=28000, active=true}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, department=engineering, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, level=1, name=Henry, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P003, level=1, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, score=10, tags=executive, name=Grace, age=65, active=true, status=active, department=management, level=9}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, active=true, tags=senior, status=active, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, region=south}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high}
4. [1] TestOrder{id=O004, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, priority=normal}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100}
6. [1] TestOrder{id=O006, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, total=999.98}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}
8. [1] TestOrder{id=O008, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0}
9. [1] TestOrder{id=O009, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, customer_id=P001, priority=low, discount=10}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, status=confirmed, priority=low, discount=0, region=south, product_id=PROD002, amount=1, total=25.5, customer_id=P002, date=2024-01-20}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, date=2024-02-10, priority=high, region=south, total=999.99, status=confirmed, discount=100, customer_id=P002, product_id=PROD001, amount=1}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, date=2024-02-15, discount=0, region=west, amount=2, status=cancelled, priority=low}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, total=600, priority=urgent}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, product_id=PROD001, total=1999.98, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, discount=15, region=north, customer_id=P001, status=shipped, priority=high}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, amount=1, priority=normal, discount=0, customer_id=P004, total=299.99, date=2024-02-05, status=delivered, region=east}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, priority=low, discount=10, region=north, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering}
3. [1] TestPerson{id=P003, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing}
5. [1] TestPerson{id=P005, score=8, status=inactive, level=3, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, department=qa, age=0, salary=-5000, tags=test, status=active, level=1}
7. [1] TestPerson{id=P007, status=active, department=management, level=9, salary=95000, score=10, tags=executive, name=Grace, age=65, active=true}
8. [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, name=Henry}
9. [1] TestPerson{id=P009, department=engineering, level=6, age=40, active=true, tags=senior, status=active, name=Ivy, salary=68000, score=8.7}
10. [1] TestPerson{id=P010, department=intern, level=1, age=22, tags=temp, status=active, name=X, salary=28000, active=true, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, age=0, salary=-5000, tags=test, status=active, level=1, name=Frank, active=true, score=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, name=Henry}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, active=true, score=6.5, department=intern, level=1, age=22, tags=temp, status=active, name=X, salary=28000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, status=active, level=2, active=true, department=sales, name=Alice, age=25}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, salary=55000, tags=employee, department=sales, age=30, active=false, score=8, status=inactive, level=3}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, active=true, status=active, department=management, level=9, salary=95000, score=10, tags=executive}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, department=engineering, level=6, age=40, active=true, tags=senior, status=active}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, status=active, level=5, name=Bob, department=engineering, age=35, salary=75000}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P003, score=6, status=inactive, name=Charlie, age=16, salary=0, active=false, tags=intern, department=hr, level=1}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, level=7, age=45, active=true, department=marketing, name=Diana, salary=85000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
