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

1. [1] TestPerson{id=P001, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2, name=Alice, age=25, active=true}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5}
3. [1] TestPerson{id=P003, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr}
4. [1] TestPerson{id=P004, level=7, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales}
6. [1] TestPerson{id=P006, status=active, level=1, salary=-5000, department=qa, name=Frank, age=0, active=true, score=0, tags=test}
7. [1] TestPerson{id=P007, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000, age=65, active=true}
8. [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive, level=3}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, department=engineering, level=6, name=Ivy, age=40, salary=68000, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, discount=50, region=north, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15}
2. [1] TestOrder{id=O002, total=25.5, status=confirmed, priority=low, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15}
4. [1] TestOrder{id=O004, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1}
5. [1] TestOrder{id=O005, status=confirmed, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north}
8. [1] TestOrder{id=O008, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002, amount=10}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O008, date=2024-03-05, region=south, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O010, total=75000, region=east, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south, customer_id=P002}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active}
2. [1] TestPerson{id=P002, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5, name=Bob, salary=75000}
3. [1] TestPerson{id=P003, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0, active=false}
4. [1] TestPerson{id=P004, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana}
5. [1] TestPerson{id=P005, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve}
6. [1] TestPerson{id=P006, status=active, level=1, salary=-5000, department=qa, name=Frank, age=0, active=true, score=0, tags=test}
7. [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}
8. [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, department=engineering, level=6, name=Ivy, age=40, salary=68000, active=true}
10. [1] TestPerson{id=P010, tags=temp, status=active, level=1, age=22, salary=28000, department=intern, name=X, active=true, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, level=5, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive, level=3}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, level=9, name=Grace, salary=95000, age=65, active=true, score=10, tags=executive, status=active, department=management}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6, name=Ivy, age=40, salary=68000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, discount=50, region=north, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}
3. [1] TestOrder{id=O003, status=shipped, priority=high, region=north, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225}
4. [1] TestOrder{id=O004, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal}
5. [1] TestOrder{id=O005, status=confirmed, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north}
8. [1] TestOrder{id=O008, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002, amount=10}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O006, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2, total=999.98}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O008, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002, amount=10, total=255, status=pending}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, product_id=PROD001, status=pending, priority=normal}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, region=south, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, category=electronics, price=999.99, rating=4.5, name=Laptop, available=true, keywords=computer, brand=TechCorp, stock=50, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, keywords=peripheral, brand=TechCorp, price=25.5, available=true, rating=4, stock=200, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply, category=accessories, price=75, available=false, rating=3.5, name=Keyboard}
4. [1] TestProduct{id=PROD004, category=electronics, available=true, brand=ScreenPro, stock=30, name=Monitor, price=299.99, rating=4.8, keywords=display, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, brand=OldTech, name=OldKeyboard, rating=2, keywords=obsolete, stock=0, supplier=OldSupply, category=accessories, price=8.5, available=false}
6. [1] TestProduct{id=PROD006, available=true, stock=75, supplier=AudioSupply, name=Headphones, category=audio, rating=4.6, keywords=sound, brand=AudioMax, price=150}
7. [1] TestProduct{id=PROD007, brand=CamTech, stock=25, name=Webcam, available=true, supplier=CamSupply, category=electronics, price=89.99, rating=3.8, keywords=video}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD006, price=150, available=true, stock=75, supplier=AudioSupply, name=Headphones, category=audio, rating=4.6, keywords=sound, brand=AudioMax}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD007, category=electronics, price=89.99, rating=3.8, keywords=video, brand=CamTech, stock=25, name=Webcam, available=true, supplier=CamSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD001, rating=4.5, name=Laptop, available=true, keywords=computer, brand=TechCorp, stock=50, supplier=TechSupply, category=electronics, price=999.99}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, category=accessories, keywords=peripheral, brand=TechCorp, price=25.5, available=true, rating=4, stock=200, supplier=TechSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD003, rating=3.5, name=Keyboard, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply, category=accessories, price=75, available=false}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, price=299.99, rating=4.8, keywords=display, supplier=ScreenSupply, category=electronics, available=true, brand=ScreenPro, stock=30}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales}
2. [1] TestPerson{id=P002, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35}
3. [1] TestPerson{id=P003, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern}
4. [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}
5. [1] TestPerson{id=P005, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false, score=8}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
7. [1] TestPerson{id=P007, level=9, name=Grace, salary=95000, age=65, active=true, score=10, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, department=qa, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north}
2. [1] TestOrder{id=O002, status=confirmed, priority=low, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15}
4. [1] TestOrder{id=O004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south}
6. [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north}
8. [1] TestOrder{id=O008, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002, amount=10, total=255, status=pending}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north}
10. [1] TestOrder{id=O010, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east, product_id=PROD001}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, discount=0, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, status=confirmed, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, product_id=PROD001, status=pending, priority=normal}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2}
2. [1] TestPerson{id=P002, age=35, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2, department=engineering}
3. [1] TestPerson{id=P003, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr}
4. [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}
5. [1] TestPerson{id=P005, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
7. [1] TestPerson{id=P007, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000, age=65, active=true, score=10}
8. [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, department=engineering, level=6, name=Ivy, age=40, salary=68000, active=true}
10. [1] TestPerson{id=P010, age=22, salary=28000, department=intern, name=X, active=true, score=6.5, tags=temp, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive, level=3}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, discount=50, region=north, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2}
2. [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, priority=low, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15}
4. [1] TestOrder{id=O004, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low}
7. [1] TestOrder{id=O007, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O009, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O003, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, status=confirmed, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2, name=Alice, age=25}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5, name=Bob}
3. [1] TestPerson{id=P003, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr}
4. [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}
5. [1] TestPerson{id=P005, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve}
6. [1] TestPerson{id=P006, department=qa, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000}
7. [1] TestPerson{id=P007, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000, age=65}
8. [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6, name=Ivy, age=40}
10. [1] TestPerson{id=P010, status=active, level=1, age=22, salary=28000, department=intern, name=X, active=true, score=6.5, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, age=35, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0, active=false, score=6}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, keywords=computer, brand=TechCorp, stock=50, supplier=TechSupply, category=electronics, price=999.99, rating=4.5, name=Laptop, available=true}
2. [1] TestProduct{id=PROD002, rating=4, stock=200, supplier=TechSupply, name=Mouse, category=accessories, keywords=peripheral, brand=TechCorp, price=25.5, available=true}
3. [1] TestProduct{id=PROD003, category=accessories, price=75, available=false, rating=3.5, name=Keyboard, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, name=Monitor, price=299.99, rating=4.8, keywords=display, supplier=ScreenSupply, category=electronics, available=true, brand=ScreenPro, stock=30}
5. [1] TestProduct{id=PROD005, keywords=obsolete, stock=0, supplier=OldSupply, category=accessories, price=8.5, available=false, brand=OldTech, name=OldKeyboard, rating=2}
6. [1] TestProduct{id=PROD006, brand=AudioMax, price=150, available=true, stock=75, supplier=AudioSupply, name=Headphones, category=audio, rating=4.6, keywords=sound}
7. [1] TestProduct{id=PROD007, brand=CamTech, stock=25, name=Webcam, available=true, supplier=CamSupply, category=electronics, price=89.99, rating=3.8, keywords=video}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, category=electronics, price=999.99, rating=4.5, name=Laptop, available=true, keywords=computer, brand=TechCorp, stock=50, supplier=TechSupply}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, category=accessories, keywords=peripheral, brand=TechCorp, price=25.5, available=true, rating=4, stock=200, supplier=TechSupply, name=Mouse}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, category=accessories, price=75, available=false, rating=3.5, name=Keyboard, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, price=299.99, rating=4.8, keywords=display, supplier=ScreenSupply, category=electronics, available=true, brand=ScreenPro, stock=30}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, name=Headphones, category=audio, rating=4.6, keywords=sound, brand=AudioMax, price=150, available=true, stock=75, supplier=AudioSupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, category=electronics, price=89.99, rating=3.8, keywords=video, brand=CamTech, stock=25, name=Webcam, available=true, supplier=CamSupply}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5}
3. [1] TestPerson{id=P003, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr}
4. [1] TestPerson{id=P004, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7}
5. [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive, level=3}
6. [1] TestPerson{id=P006, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa, name=Frank}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9}
8. [1] TestPerson{id=P008, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern, name=X, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive, level=3}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6, name=Ivy, age=40, salary=68000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, status=active, salary=45000, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, product_id=PROD001, status=pending, priority=normal}
2. [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low, product_id=PROD002, amount=1, date=2024-01-20}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15}
4. [1] TestOrder{id=O004, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05}
5. [1] TestOrder{id=O005, amount=1, status=confirmed, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low}
7. [1] TestOrder{id=O007, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600}
8. [1] TestOrder{id=O008, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002, amount=10, total=255}
9. [1] TestOrder{id=O009, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, discount=50, region=north, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, priority=low, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, region=north, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O008, date=2024-03-05, region=south, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales}
2. [1] TestPerson{id=P002, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35}
3. [1] TestPerson{id=P003, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0, active=false}
4. [1] TestPerson{id=P004, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true}
5. [1] TestPerson{id=P005, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales}
6. [1] TestPerson{id=P006, score=0, tags=test, status=active, level=1, salary=-5000, department=qa, name=Frank, age=0, active=true}
7. [1] TestPerson{id=P007, status=active, department=management, level=9, name=Grace, salary=95000, age=65, active=true, score=10, tags=executive}
8. [1] TestPerson{id=P008, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1}
9. [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6, name=Ivy, age=40}
10. [1] TestPerson{id=P010, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern, name=X, active=true}
11. [1] TestOrder{id=O001, date=2024-01-15, discount=50, region=north, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98}
12. [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}
13. [1] TestOrder{id=O003, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north}
14. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east}
15. [1] TestOrder{id=O005, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high}
16. [1] TestOrder{id=O006, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2, total=999.98, status=cancelled}
17. [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north}
18. [1] TestOrder{id=O008, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002}
19. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10}
20. [1] TestOrder{id=O010, priority=urgent, discount=0, customer_id=P006, total=75000, region=east, product_id=PROD001, amount=1, date=2024-03-15, status=refunded}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P010, level=1, age=22, salary=28000, department=intern, name=X, active=true, score=6.5, tags=temp, status=active}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}
   - Fait 2: [1] TestOrder{id=O007, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, level=1, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}
   - Fait 2: [1] TestOrder{id=O004, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6, name=Ivy}
   - Fait 2: [1] TestOrder{id=O006, region=west, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}
   - Fait 2: [1] TestOrder{id=O010, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east, product_id=PROD001}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east, product_id=PROD001, amount=1, date=2024-03-15}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P009, level=6, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east, customer_id=P004}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}
   - Fait 2: [1] TestOrder{id=O006, priority=low, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6, name=Ivy, age=40}
   - Fait 2: [1] TestOrder{id=O010, region=east, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0}
   - Fait 2: [1] TestOrder{id=O001, date=2024-01-15, discount=50, region=north, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O003, priority=high, region=north, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P002, age=35, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2, department=engineering}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee}
   - Fait 2: [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low, product_id=PROD002, amount=1, date=2024-01-20}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O005, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, status=active, level=1, age=22, salary=28000, department=intern, name=X, active=true, score=6.5}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2, total=999.98}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O009, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P003, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie}
   - Fait 2: [1] TestOrder{id=O010, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east, product_id=PROD001}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, level=6, name=Ivy, age=40, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, department=intern, name=X, active=true, score=6.5, tags=temp, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O003, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O008, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002, amount=10, total=255}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}
   - Fait 2: [1] TestOrder{id=O001, region=north, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, priority=low, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, department=intern, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P010, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern}
   - Fait 2: [1] TestOrder{id=O004, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O005, amount=1, status=confirmed, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}
   - Fait 2: [1] TestOrder{id=O009, region=north, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P010, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern}
   - Fait 2: [1] TestOrder{id=O010, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east, product_id=PROD001}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P008, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O005, amount=1, status=confirmed, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P001, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2, name=Alice}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2, total=999.98}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000}
   - Fait 2: [1] TestOrder{id=O003, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, status=active, level=1, age=22, salary=28000, department=intern, name=X, active=true, score=6.5}
   - Fait 2: [1] TestOrder{id=O008, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002, amount=10, total=255, status=pending, priority=normal}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, product_id=PROD001, status=pending, priority=normal}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P009, level=6, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low, product_id=PROD002, amount=1, date=2024-01-20}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, age=35, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O009, region=north, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, discount=50, region=north, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P003, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6, name=Ivy}
   - Fait 2: [1] TestOrder{id=O008, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002, amount=10}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{id=O004, discount=0, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P008, level=1, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P002, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P005, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve}
   - Fait 2: [1] TestOrder{id=O007, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O009, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales, salary=55000}
   - Fait 2: [1] TestOrder{id=O010, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east, product_id=PROD001}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P005, level=3, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15, customer_id=P001}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35}
   - Fait 2: [1] TestOrder{id=O003, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P008, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}
   - Fait 2: [1] TestOrder{id=O006, amount=2, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P006, level=1, salary=-5000, department=qa, name=Frank, age=0, active=true, score=0, tags=test, status=active}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east, product_id=PROD001, amount=1}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east, product_id=PROD001, amount=1, date=2024-03-15}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, total=75000, region=east, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, department=intern, name=X, active=true, score=6.5, tags=temp, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low, product_id=PROD002, amount=1, date=2024-01-20}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{id=O003, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true}
   - Fait 2: [1] TestOrder{id=O003, priority=high, region=north, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P005, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000, age=65, active=true}
   - Fait 2: [1] TestOrder{id=O008, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, age=35, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, total=75000, region=east, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2}
2. [1] TestPerson{id=P002, age=35, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2, department=engineering}
3. [1] TestPerson{id=P003, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr}
4. [1] TestPerson{id=P004, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7}
5. [1] TestPerson{id=P005, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
7. [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}
8. [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, status=active, level=1, salary=-5000, department=qa, name=Frank, age=0, active=true, score=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P010, level=1, age=22, salary=28000, department=intern, name=X, active=true, score=6.5, tags=temp, status=active}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, status=active, salary=45000, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5, name=Bob, salary=75000}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2, name=Alice, age=25}
2. [1] TestPerson{id=P002, level=5, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active}
3. [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active}
5. [1] TestPerson{id=P005, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
7. [1] TestPerson{id=P007, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000, age=65, active=true}
8. [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern}
11. [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, product_id=PROD001, status=pending, priority=normal}
12. [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}
13. [1] TestOrder{id=O003, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north}
14. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east}
15. [1] TestOrder{id=O005, amount=1, status=confirmed, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10}
16. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low}
17. [1] TestOrder{id=O007, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50}
18. [1] TestOrder{id=O008, date=2024-03-05, region=south, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010}
19. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10}
20. [1] TestOrder{id=O010, customer_id=P006, total=75000, region=east, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, total=75000, region=east, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, product_id=PROD001, status=pending, priority=normal}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O004, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active}
   - Fait 2: [1] TestOrder{id=O005, amount=1, status=confirmed, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P006, score=0, tags=test, status=active, level=1, salary=-5000, department=qa, name=Frank, age=0, active=true}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{id=O007, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, product_id=PROD001, status=pending, priority=normal}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry, active=false}
   - Fait 2: [1] TestOrder{id=O004, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, level=6, name=Ivy, age=40, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2, total=999.98}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P006, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa, name=Frank, age=0}
   - Fait 2: [1] TestOrder{id=O001, date=2024-01-15, discount=50, region=north, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15, customer_id=P001}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{id=O004, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P002, age=35, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2, department=engineering}
   - Fait 2: [1] TestOrder{id=O007, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6, name=Ivy}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0}
   - Fait 2: [1] TestOrder{id=O001, discount=50, region=north, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low, product_id=PROD002, amount=1, date=2024-01-20}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P005, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false}
   - Fait 2: [1] TestOrder{id=O003, region=north, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P003, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0, active=false}
   - Fait 2: [1] TestOrder{id=O004, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P008, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern, name=X, active=true}
   - Fait 2: [1] TestOrder{id=O002, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low, product_id=PROD002, amount=1}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, age=35, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2}
   - Fait 2: [1] TestOrder{id=O003, status=shipped, priority=high, region=north, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P007, level=9, name=Grace, salary=95000, age=65, active=true, score=10, tags=executive, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing, age=45}
   - Fait 2: [1] TestOrder{id=O006, priority=low, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern}
   - Fait 2: [1] TestOrder{id=O008, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002, amount=10, total=255, status=pending, priority=normal}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P002, age=35, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2, department=engineering}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P008, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O003, status=shipped, priority=high, region=north, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P005, level=3, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive}
   - Fait 2: [1] TestOrder{id=O007, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000}
   - Fait 2: [1] TestOrder{id=O008, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002, amount=10}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{id=O009, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P010, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east, product_id=PROD001, amount=1}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr}
   - Fait 2: [1] TestOrder{id=O005, amount=1, status=confirmed, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P009, status=active, department=engineering, level=6, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P007, status=active, department=management, level=9, name=Grace, salary=95000, age=65, active=true, score=10, tags=executive}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P002, age=35, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2, department=engineering}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, department=intern, name=X, active=true, score=6.5, tags=temp, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O003, priority=high, region=north, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P006, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa, name=Frank, age=0}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000}
   - Fait 2: [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P010, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1}
   - Fait 2: [1] TestOrder{id=O010, discount=0, customer_id=P006, total=75000, region=east, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, level=6, name=Ivy, age=40, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6, name=Ivy}
   - Fait 2: [1] TestOrder{id=O009, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P004, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000, age=65}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low, product_id=PROD002}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P001, status=active, salary=45000, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O003, status=shipped, priority=high, region=north, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P007, level=9, name=Grace, salary=95000, age=65, active=true, score=10, tags=executive, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O004, discount=0, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P010, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern}
   - Fait 2: [1] TestOrder{id=O006, amount=2, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P010, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east, product_id=PROD001, amount=1}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P001, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P010, status=active, level=1, age=22, salary=28000, department=intern, name=X, active=true, score=6.5, tags=temp}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, product_id=PROD001, status=pending, priority=normal}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern, name=X, active=true}
   - Fait 2: [1] TestOrder{id=O004, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2, total=999.98, status=cancelled}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P007, level=9, name=Grace, salary=95000, age=65, active=true, score=10, tags=executive, status=active, department=management}
   - Fait 2: [1] TestOrder{id=O007, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, product_id=PROD001, status=pending, priority=normal}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr}
   - Fait 2: [1] TestOrder{id=O002, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low, product_id=PROD002, amount=1, date=2024-01-20, discount=0}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, department=engineering, level=6, name=Ivy, age=40, salary=68000, active=true}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P008, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P010, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern}
   - Fait 2: [1] TestOrder{id=O008, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002, amount=10, total=255, status=pending}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{id=O010, discount=0, customer_id=P006, total=75000, region=east, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2, name=Alice, age=25, active=true}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, priority=low, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east, customer_id=P004}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true}
   - Fait 2: [1] TestOrder{id=O005, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}
   - Fait 2: [1] TestOrder{id=O009, region=north, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000, region=east}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P004, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana}
   - Fait 2: [1] TestOrder{id=O001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, product_id=PROD001}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}
   - Fait 2: [1] TestOrder{id=O002, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low, product_id=PROD002, amount=1, date=2024-01-20, discount=0}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active}
   - Fait 2: [1] TestOrder{id=O004, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P010, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P008, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support, level=1, name=Henry}
   - Fait 2: [1] TestOrder{id=O008, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002, amount=10, total=255, status=pending, priority=normal}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, status=active, level=1, salary=-5000, department=qa, name=Frank, age=0, active=true, score=0}
   - Fait 2: [1] TestOrder{id=O009, priority=low, region=north, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O001, region=north, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}
   - Fait 2: [1] TestOrder{id=O003, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P005, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales, salary=55000, active=false}
   - Fait 2: [1] TestOrder{id=O004, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5}
   - Fait 2: [1] TestOrder{id=O008, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south, product_id=PROD002, amount=10}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2}
2. [1] TestPerson{id=P002, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35}
3. [1] TestPerson{id=P003, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16}
4. [1] TestPerson{id=P004, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7}
5. [1] TestPerson{id=P005, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales}
6. [1] TestPerson{id=P006, tags=test, status=active, level=1, salary=-5000, department=qa, name=Frank, age=0, active=true, score=0}
7. [1] TestPerson{id=P007, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000, age=65, active=true, score=10}
8. [1] TestPerson{id=P008, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22, salary=28000, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, level=1, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000, score=5.5, department=support}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, department=intern, name=X, active=true, score=6.5, tags=temp, status=active, level=1, age=22}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, age=35, tags=senior, status=active, level=5, name=Bob, salary=75000, active=true, score=9.2, department=engineering}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, tags=manager, level=7, name=Diana, status=active, department=marketing, age=45}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000, age=65, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, status=pending, priority=normal, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}
3. [1] TestOrder{id=O003, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003, amount=3, total=225, status=shipped, priority=high, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0, amount=1, total=299.99, priority=normal, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10, amount=1, status=confirmed, priority=high, discount=100, region=south}
6. [1] TestOrder{id=O006, total=999.98, status=cancelled, discount=0, region=west, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2}
7. [1] TestOrder{id=O007, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010, date=2024-03-05, region=south}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, discount=10, customer_id=P001, product_id=PROD007, amount=1, status=completed, priority=low, region=north}
10. [1] TestOrder{id=O010, region=east, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0, customer_id=P006, total=75000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, region=south, customer_id=P002, total=25.5, status=confirmed, priority=low}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, total=299.99, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, status=delivered, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, priority=low, product_id=PROD005, amount=2, total=999.98, status=cancelled, discount=0, region=west}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, amount=1, status=completed, priority=low, region=north, total=89.99, date=2024-03-10, discount=10, customer_id=P001}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, product_id=PROD001, status=pending, priority=normal}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, total=225, status=shipped, priority=high, region=north, date=2024-02-01, discount=15, customer_id=P001, product_id=PROD003}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, status=confirmed, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, date=2024-02-10}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, date=2024-03-05, region=south, product_id=PROD002, amount=10, total=255, status=pending, priority=normal, discount=0, customer_id=P010}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, total=75000, region=east, product_id=PROD001, amount=1, date=2024-03-15, status=refunded, priority=urgent, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active, salary=45000, department=sales, level=2}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5}
3. [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7}
5. [1] TestPerson{id=P005, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales}
6. [1] TestPerson{id=P006, active=true, score=0, tags=test, status=active, level=1, salary=-5000, department=qa, name=Frank, age=0}
7. [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9, name=Grace, salary=95000}
8. [1] TestPerson{id=P008, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive, age=18, salary=25000}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, age=22, salary=28000, department=intern, name=X, active=true, score=6.5, tags=temp, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, age=30, tags=employee, department=sales}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, name=Frank, age=0, active=true, score=0, tags=test, status=active, level=1, salary=-5000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, active=true, score=8.7, tags=senior, status=active, department=engineering, level=6}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, department=sales, level=2, name=Alice, age=25, active=true, score=8.5, tags=junior, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, department=engineering, age=35, tags=senior, status=active, level=5}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, status=active, department=marketing, age=45, salary=85000, active=true, score=7.8, tags=manager, level=7}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, age=65, active=true, score=10, tags=executive, status=active, department=management, level=9}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, score=5.5, department=support, level=1, name=Henry, active=false, tags=junior, status=inactive, age=18}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, status=active, level=1, age=22, salary=28000, department=intern, name=X, active=true, score=6.5}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
