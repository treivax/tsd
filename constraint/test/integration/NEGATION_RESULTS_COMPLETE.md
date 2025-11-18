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

1. [1] TestPerson{id=P001, salary=45000, tags=junior, level=2, active=true, score=8.5, status=active, department=sales, name=Alice, age=25}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, score=9.2, status=active, level=5, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}
4. [1] TestPerson{id=P004, score=7.8, status=active, department=marketing, name=Diana, age=45, active=true, tags=manager, level=7, salary=85000}
5. [1] TestPerson{id=P005, tags=employee, status=inactive, department=sales, name=Eve, age=30, active=false, score=8, level=3, salary=55000}
6. [1] TestPerson{id=P006, age=0, salary=-5000, score=0, tags=test, active=true, status=active, department=qa, level=1, name=Frank}
7. [1] TestPerson{id=P007, name=Grace, score=10, tags=executive, level=9, age=65, salary=95000, active=true, status=active, department=management}
8. [1] TestPerson{id=P008, level=1, age=18, salary=25000, score=5.5, department=support, name=Henry, active=false, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, score=8.7, tags=senior, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000}
10. [1] TestPerson{id=P010, salary=28000, score=6.5, tags=temp, active=true, status=active, department=intern, level=1, name=X, age=22}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, tags=junior, level=2, active=true, score=8.5, status=active, department=sales, name=Alice}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, department=engineering, name=Bob, age=35, salary=75000, active=true, score=9.2, status=active, level=5}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, active=true, tags=manager, level=7, salary=85000, score=7.8, status=active, department=marketing}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, active=false, tags=junior, status=inactive, level=1, age=18, salary=25000, score=5.5, department=support, name=Henry}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, active=true, status=active, department=intern, level=1, name=X, age=22, salary=28000, score=6.5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, level=3, salary=55000, tags=employee, status=inactive, department=sales}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, status=active, department=management, name=Grace, score=10, tags=executive, level=9}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, date=2024-01-15, status=pending, discount=50, region=north, customer_id=P001, amount=2, total=1999.98, priority=normal, product_id=PROD001}
2. [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, product_id=PROD002, amount=1, total=25.5, region=south}
3. [1] TestOrder{id=O003, priority=high, region=north, customer_id=P001, product_id=PROD003, amount=3, status=shipped, discount=15, total=225, date=2024-02-01}
4. [1] TestOrder{id=O004, amount=1, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, total=299.99, priority=normal, region=east, customer_id=P004}
5. [1] TestOrder{id=O005, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, date=2024-02-15, discount=0, region=west, total=999.98, status=cancelled, priority=low, customer_id=P005}
7. [1] TestOrder{id=O007, customer_id=P007, amount=4, date=2024-03-01, status=shipped, discount=50, product_id=PROD006, total=600, priority=urgent, region=north}
8. [1] TestOrder{id=O008, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal, customer_id=P010, date=2024-03-05, status=pending}
9. [1] TestOrder{id=O009, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, total=75000, status=refunded, priority=urgent, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, priority=normal, product_id=PROD001, date=2024-01-15, status=pending, discount=50, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, date=2024-01-20, status=confirmed, priority=low, discount=0, product_id=PROD002, amount=1, total=25.5, region=south, customer_id=P002}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, total=600, priority=urgent, region=north, customer_id=P007, amount=4, date=2024-03-01, status=shipped, discount=50}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O008, status=pending, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal, customer_id=P010, date=2024-03-05}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O003, total=225, date=2024-02-01, priority=high, region=north, customer_id=P001, product_id=PROD003, amount=3, status=shipped, discount=15}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, total=299.99, priority=normal, region=east}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O010, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, total=75000, status=refunded, priority=urgent, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, level=2, active=true, score=8.5, status=active, department=sales, name=Alice, age=25, salary=45000}
2. [1] TestPerson{id=P002, score=9.2, status=active, level=5, tags=senior, department=engineering, name=Bob, age=35, salary=75000, active=true}
3. [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, tags=manager, level=7, salary=85000, score=7.8, status=active, department=marketing}
5. [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, department=sales, name=Eve, age=30, active=false, score=8, level=3}
6. [1] TestPerson{id=P006, active=true, status=active, department=qa, level=1, name=Frank, age=0, salary=-5000, score=0, tags=test}
7. [1] TestPerson{id=P007, score=10, tags=executive, level=9, age=65, salary=95000, active=true, status=active, department=management, name=Grace}
8. [1] TestPerson{id=P008, name=Henry, active=false, tags=junior, status=inactive, level=1, age=18, salary=25000, score=5.5, department=support}
9. [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, active=true, status=active, department=intern, level=1, name=X, age=22, salary=28000, score=6.5, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, level=9, age=65, salary=95000, active=true, status=active, department=management, name=Grace}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior, age=40, active=true}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, level=2, active=true, score=8.5, status=active, department=sales}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, name=Bob, age=35, salary=75000, active=true, score=9.2, status=active, level=5, tags=senior}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, active=true, tags=manager, level=7, salary=85000, score=7.8, status=active, department=marketing}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, level=3, salary=55000, tags=employee, status=inactive, department=sales}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, total=1999.98, priority=normal, product_id=PROD001, date=2024-01-15, status=pending, discount=50, region=north, customer_id=P001, amount=2}
2. [1] TestOrder{id=O002, amount=1, total=25.5, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, product_id=PROD002}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, discount=15, total=225, date=2024-02-01, priority=high, region=north}
4. [1] TestOrder{id=O004, discount=0, product_id=PROD004, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, status=delivered}
5. [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed}
6. [1] TestOrder{id=O006, priority=low, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, discount=0, region=west, total=999.98, status=cancelled}
7. [1] TestOrder{id=O007, customer_id=P007, amount=4, date=2024-03-01, status=shipped, discount=50, product_id=PROD006, total=600, priority=urgent, region=north}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, priority=normal, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south}
9. [1] TestOrder{id=O009, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, total=75000, status=refunded, priority=urgent, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, total=225, date=2024-02-01, priority=high, region=north, customer_id=P001, product_id=PROD003, amount=3, status=shipped, discount=15}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O008, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal, customer_id=P010, date=2024-03-05, status=pending}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, product_id=PROD002, amount=1, total=25.5, region=south}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, total=299.99, priority=normal, region=east}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, discount=0, region=west, total=999.98, status=cancelled, priority=low}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, amount=4, date=2024-03-01, status=shipped, discount=50, product_id=PROD006, total=600, priority=urgent, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O001, total=1999.98, priority=normal, product_id=PROD001, date=2024-01-15, status=pending, discount=50, region=north, customer_id=P001, amount=2}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, category=electronics, available=true, supplier=TechSupply, price=999.99, rating=4.5, keywords=computer, brand=TechCorp, stock=50}
2. [1] TestProduct{id=PROD002, supplier=TechSupply, name=Mouse, category=accessories, price=25.5, available=true, rating=4, keywords=peripheral, brand=TechCorp, stock=200}
3. [1] TestProduct{id=PROD003, name=Keyboard, price=75, brand=KeyTech, category=accessories, available=false, rating=3.5, keywords=typing, stock=0, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, name=Monitor, price=299.99, keywords=display, brand=ScreenPro, stock=30, category=electronics, available=true, rating=4.8, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, category=accessories, price=8.5, available=false, rating=2, keywords=obsolete, brand=OldTech, stock=0, name=OldKeyboard, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, keywords=sound, name=Headphones, category=audio, brand=AudioMax, stock=75, supplier=AudioSupply, price=150, available=true, rating=4.6}
7. [1] TestProduct{id=PROD007, stock=25, supplier=CamSupply, category=electronics, price=89.99, available=true, name=Webcam, rating=3.8, keywords=video, brand=CamTech}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, category=electronics, available=true, supplier=TechSupply, price=999.99, rating=4.5, keywords=computer, brand=TechCorp, stock=50}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, price=25.5, available=true, rating=4, keywords=peripheral, brand=TechCorp, stock=200, supplier=TechSupply, name=Mouse, category=accessories}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, category=accessories, available=false, rating=3.5, keywords=typing, stock=0, supplier=KeySupply, name=Keyboard, price=75, brand=KeyTech}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, keywords=display, brand=ScreenPro, stock=30, category=electronics, available=true, rating=4.8, supplier=ScreenSupply, name=Monitor, price=299.99}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, keywords=sound, name=Headphones, category=audio, brand=AudioMax, stock=75, supplier=AudioSupply, price=150, available=true, rating=4.6}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, stock=25, supplier=CamSupply, category=electronics, price=89.99, available=true, name=Webcam, rating=3.8, keywords=video, brand=CamTech}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, level=2, active=true, score=8.5, status=active, department=sales}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, score=9.2, status=active, level=5, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie}
4. [1] TestPerson{id=P004, age=45, active=true, tags=manager, level=7, salary=85000, score=7.8, status=active, department=marketing, name=Diana}
5. [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, department=sales, name=Eve, age=30, active=false, score=8, level=3}
6. [1] TestPerson{id=P006, active=true, status=active, department=qa, level=1, name=Frank, age=0, salary=-5000, score=0, tags=test}
7. [1] TestPerson{id=P007, name=Grace, score=10, tags=executive, level=9, age=65, salary=95000, active=true, status=active, department=management}
8. [1] TestPerson{id=P008, tags=junior, status=inactive, level=1, age=18, salary=25000, score=5.5, department=support, name=Henry, active=false}
9. [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, active=true, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, active=true, status=active, department=qa, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, priority=normal, product_id=PROD001, date=2024-01-15, status=pending, discount=50, region=north}
2. [1] TestOrder{id=O002, amount=1, total=25.5, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, product_id=PROD002}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, discount=15, total=225, date=2024-02-01, priority=high, region=north}
4. [1] TestOrder{id=O004, amount=1, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, total=299.99, priority=normal, region=east, customer_id=P004}
5. [1] TestOrder{id=O005, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10}
6. [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, discount=0, region=west}
7. [1] TestOrder{id=O007, product_id=PROD006, total=600, priority=urgent, region=north, customer_id=P007, amount=4, date=2024-03-01, status=shipped, discount=50}
8. [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal}
9. [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99}
10. [1] TestOrder{id=O010, total=75000, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, priority=normal, product_id=PROD001, date=2024-01-15, status=pending, discount=50, region=north, customer_id=P001, amount=2, total=1999.98}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, product_id=PROD002, amount=1, total=25.5, region=south}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, discount=15, total=225, date=2024-02-01, priority=high, region=north, customer_id=P001, product_id=PROD003, amount=3, status=shipped}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O004, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, total=299.99, priority=normal, region=east, customer_id=P004, amount=1}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O005, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10, priority=high, discount=100}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O006, region=west, total=999.98, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, total=600, priority=urgent, region=north, customer_id=P007, amount=4, date=2024-03-01, status=shipped, discount=50}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, total=75000, status=refunded, priority=urgent, discount=0, region=east}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, priority=normal, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O009, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, score=8.5, status=active, department=sales, name=Alice, age=25, salary=45000, tags=junior, level=2}
2. [1] TestPerson{id=P002, status=active, level=5, tags=senior, department=engineering, name=Bob, age=35, salary=75000, active=true, score=9.2}
3. [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}
4. [1] TestPerson{id=P004, status=active, department=marketing, name=Diana, age=45, active=true, tags=manager, level=7, salary=85000, score=7.8}
5. [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, level=3, salary=55000, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, active=true, status=active, department=qa, level=1, name=Frank, age=0, salary=-5000, score=0, tags=test}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, status=active, department=management, name=Grace, score=10, tags=executive, level=9}
8. [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, name=Henry, active=false, tags=junior, status=inactive, level=1}
9. [1] TestPerson{id=P009, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior, age=40}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, active=true, status=active, department=intern, level=1, name=X}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, department=sales, name=Eve, age=30, active=false, score=8, level=3}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, tags=junior, status=inactive, level=1, age=18, salary=25000, score=5.5, department=support}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, region=north, customer_id=P001, amount=2, total=1999.98, priority=normal, product_id=PROD001, date=2024-01-15, status=pending, discount=50}
2. [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, product_id=PROD002, amount=1, total=25.5, region=south}
3. [1] TestOrder{id=O003, amount=3, status=shipped, discount=15, total=225, date=2024-02-01, priority=high, region=north, customer_id=P001, product_id=PROD003}
4. [1] TestOrder{id=O004, status=delivered, discount=0, product_id=PROD004, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05}
5. [1] TestOrder{id=O005, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10, priority=high, discount=100}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, discount=0, region=west, total=999.98, status=cancelled, priority=low}
7. [1] TestOrder{id=O007, product_id=PROD006, total=600, priority=urgent, region=north, customer_id=P007, amount=4, date=2024-03-01, status=shipped, discount=50}
8. [1] TestOrder{id=O008, amount=10, total=255, priority=normal, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002}
9. [1] TestOrder{id=O009, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10}
10. [1] TestOrder{id=O010, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, total=75000, status=refunded}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, status=confirmed, priority=low, discount=0, product_id=PROD002, amount=1, total=25.5, region=south, customer_id=P002, date=2024-01-20}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, status=shipped, discount=15, total=225, date=2024-02-01, priority=high, region=north, customer_id=P001, product_id=PROD003, amount=3}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, total=299.99}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, status=confirmed, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, date=2024-02-15, discount=0, region=west, total=999.98, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, amount=2}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, total=600, priority=urgent, region=north, customer_id=P007, amount=4, date=2024-03-01, status=shipped, discount=50}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O010, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, total=75000, status=refunded, priority=urgent, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, score=8.5, status=active, department=sales, name=Alice, age=25, salary=45000, tags=junior, level=2}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, status=active, level=5, tags=senior, department=engineering, name=Bob, age=35}
3. [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, tags=manager, level=7, salary=85000, score=7.8, status=active, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, level=3, salary=55000, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, active=true, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, score=10, tags=executive, level=9, age=65, salary=95000, active=true, status=active, department=management}
8. [1] TestPerson{id=P008, level=1, age=18, salary=25000, score=5.5, department=support, name=Henry, active=false, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior, age=40}
10. [1] TestPerson{id=P010, status=active, department=intern, level=1, name=X, age=22, salary=28000, score=6.5, tags=temp, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0, active=false}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, age=30, active=false, score=8, level=3, salary=55000, tags=employee, status=inactive, department=sales, name=Eve}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, score=10, tags=executive, level=9, age=65, salary=95000, active=true, status=active, department=management}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, department=support, name=Henry, active=false, tags=junior, status=inactive, level=1, age=18, salary=25000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, level=2, active=true, score=8.5, status=active, department=sales, name=Alice, age=25, salary=45000, tags=junior}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, active=true, tags=manager, level=7, salary=85000, score=7.8, status=active, department=marketing}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P006, active=true, status=active, department=qa, level=1, name=Frank, age=0, salary=-5000, score=0, tags=test}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, age=40, active=true, status=active, department=engineering, level=6}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, status=active, level=5, tags=senior, department=engineering, name=Bob}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, category=electronics, available=true, supplier=TechSupply, price=999.99, rating=4.5, keywords=computer, brand=TechCorp, stock=50}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, price=25.5, available=true, rating=4, keywords=peripheral, brand=TechCorp, stock=200, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, keywords=typing, stock=0, supplier=KeySupply, name=Keyboard, price=75, brand=KeyTech, category=accessories, available=false, rating=3.5}
4. [1] TestProduct{id=PROD004, available=true, rating=4.8, supplier=ScreenSupply, name=Monitor, price=299.99, keywords=display, brand=ScreenPro, stock=30, category=electronics}
5. [1] TestProduct{id=PROD005, category=accessories, price=8.5, available=false, rating=2, keywords=obsolete, brand=OldTech, stock=0, name=OldKeyboard, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, price=150, available=true, rating=4.6, keywords=sound, name=Headphones, category=audio, brand=AudioMax, stock=75, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, name=Webcam, rating=3.8, keywords=video, brand=CamTech, stock=25, supplier=CamSupply, category=electronics, price=89.99, available=true}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, category=electronics, available=true, supplier=TechSupply, price=999.99, rating=4.5, keywords=computer, brand=TechCorp, stock=50}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, category=accessories, price=25.5, available=true, rating=4, keywords=peripheral, brand=TechCorp, stock=200, supplier=TechSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, category=accessories, available=false, rating=3.5, keywords=typing, stock=0, supplier=KeySupply, name=Keyboard, price=75, brand=KeyTech}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, price=299.99, keywords=display, brand=ScreenPro, stock=30, category=electronics, available=true, rating=4.8, supplier=ScreenSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, rating=4.6, keywords=sound, name=Headphones, category=audio, brand=AudioMax, stock=75, supplier=AudioSupply, price=150, available=true}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, category=electronics, price=89.99, available=true, name=Webcam, rating=3.8, keywords=video, brand=CamTech, stock=25, supplier=CamSupply}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, status=active, department=sales, name=Alice, age=25, salary=45000, tags=junior, level=2, active=true}
2. [1] TestPerson{id=P002, score=9.2, status=active, level=5, tags=senior, department=engineering, name=Bob, age=35, salary=75000, active=true}
3. [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}
4. [1] TestPerson{id=P004, active=true, tags=manager, level=7, salary=85000, score=7.8, status=active, department=marketing, name=Diana, age=45}
5. [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, level=3, salary=55000, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, active=true, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, department=management, name=Grace, score=10, tags=executive, level=9, age=65, salary=95000, active=true, status=active}
8. [1] TestPerson{id=P008, tags=junior, status=inactive, level=1, age=18, salary=25000, score=5.5, department=support, name=Henry, active=false}
9. [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, active=true, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, active=true, tags=manager, level=7, salary=85000, score=7.8, status=active, department=marketing}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, level=3, salary=55000, tags=employee, status=inactive, department=sales}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, level=9, age=65, salary=95000, active=true, status=active, department=management, name=Grace}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, name=Henry, active=false, tags=junior, status=inactive, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, status=active, department=sales, name=Alice, age=25, salary=45000, tags=junior, level=2}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, level=1, name=Frank, age=0, salary=-5000, score=0, tags=test, active=true, status=active}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, tags=senior, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, status=active, level=5, tags=senior, department=engineering, name=Bob, age=35, salary=75000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, priority=normal, product_id=PROD001, date=2024-01-15, status=pending, discount=50, region=north}
2. [1] TestOrder{id=O002, total=25.5, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, product_id=PROD002, amount=1}
3. [1] TestOrder{id=O003, region=north, customer_id=P001, product_id=PROD003, amount=3, status=shipped, discount=15, total=225, date=2024-02-01, priority=high}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, total=299.99, priority=normal, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, discount=0, region=west, total=999.98, status=cancelled, priority=low}
7. [1] TestOrder{id=O007, region=north, customer_id=P007, amount=4, date=2024-03-01, status=shipped, discount=50, product_id=PROD006, total=600, priority=urgent}
8. [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north}
10. [1] TestOrder{id=O010, total=75000, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, total=1999.98, priority=normal, product_id=PROD001, date=2024-01-15, status=pending, discount=50, region=north, customer_id=P001, amount=2}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, status=shipped, discount=15, total=225, date=2024-02-01, priority=high, region=north, customer_id=P001, product_id=PROD003}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O004, region=east, customer_id=P004, amount=1, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, total=299.99, priority=normal}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, status=shipped, discount=50, product_id=PROD006, total=600, priority=urgent, region=north, customer_id=P007, amount=4, date=2024-03-01}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O008, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal, customer_id=P010, date=2024-03-05, status=pending}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, department=sales, name=Alice, age=25, salary=45000, tags=junior, level=2, active=true, score=8.5}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, score=9.2, status=active, level=5, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}
4. [1] TestPerson{id=P004, salary=85000, score=7.8, status=active, department=marketing, name=Diana, age=45, active=true, tags=manager, level=7}
5. [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, level=3, salary=55000, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, active=true, status=active, department=qa, level=1, name=Frank, age=0, salary=-5000, score=0, tags=test}
7. [1] TestPerson{id=P007, name=Grace, score=10, tags=executive, level=9, age=65, salary=95000, active=true, status=active, department=management}
8. [1] TestPerson{id=P008, status=inactive, level=1, age=18, salary=25000, score=5.5, department=support, name=Henry, active=false, tags=junior}
9. [1] TestPerson{id=P009, score=8.7, tags=senior, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, active=true, status=active, department=intern, level=1, name=X}
11. [1] TestOrder{id=O001, priority=normal, product_id=PROD001, date=2024-01-15, status=pending, discount=50, region=north, customer_id=P001, amount=2, total=1999.98}
12. [1] TestOrder{id=O002, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, product_id=PROD002, amount=1, total=25.5, region=south}
13. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, status=shipped, discount=15, total=225, date=2024-02-01, priority=high, region=north}
14. [1] TestOrder{id=O004, product_id=PROD004, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, status=delivered, discount=0}
15. [1] TestOrder{id=O005, status=confirmed, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99}
16. [1] TestOrder{id=O006, product_id=PROD005, amount=2, date=2024-02-15, discount=0, region=west, total=999.98, status=cancelled, priority=low, customer_id=P005}
17. [1] TestOrder{id=O007, total=600, priority=urgent, region=north, customer_id=P007, amount=4, date=2024-03-01, status=shipped, discount=50, product_id=PROD006}
18. [1] TestOrder{id=O008, region=south, product_id=PROD002, amount=10, total=255, priority=normal, customer_id=P010, date=2024-03-05, status=pending, discount=0}
19. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north}
20. [1] TestOrder{id=O010, amount=1, date=2024-03-15, total=75000, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001}

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

1. [1] TestPerson{id=P001, active=true, score=8.5, status=active, department=sales, name=Alice, age=25, salary=45000, tags=junior, level=2}
2. [1] TestPerson{id=P002, active=true, score=9.2, status=active, level=5, tags=senior, department=engineering, name=Bob, age=35, salary=75000}
3. [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, tags=manager, level=7, salary=85000, score=7.8, status=active, department=marketing}
5. [1] TestPerson{id=P005, status=inactive, department=sales, name=Eve, age=30, active=false, score=8, level=3, salary=55000, tags=employee}
6. [1] TestPerson{id=P006, tags=test, active=true, status=active, department=qa, level=1, name=Frank, age=0, salary=-5000, score=0}
7. [1] TestPerson{id=P007, status=active, department=management, name=Grace, score=10, tags=executive, level=9, age=65, salary=95000, active=true}
8. [1] TestPerson{id=P008, name=Henry, active=false, tags=junior, status=inactive, level=1, age=18, salary=25000, score=5.5, department=support}
9. [1] TestPerson{id=P009, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior, age=40}
10. [1] TestPerson{id=P010, status=active, department=intern, level=1, name=X, age=22, salary=28000, score=6.5, tags=temp, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, score=7.8, status=active, department=marketing, name=Diana, age=45, active=true, tags=manager, level=7}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, score=0, tags=test, active=true, status=active, department=qa, level=1, name=Frank}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, active=true, status=active, department=management, name=Grace, score=10, tags=executive, level=9, age=65}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, age=40, active=true, status=active, department=engineering, level=6}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, active=true, status=active, department=intern, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, tags=junior, level=2, active=true, score=8.5, status=active, department=sales, name=Alice, age=25}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, status=active, level=5, tags=senior, department=engineering, name=Bob, age=35, salary=75000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, tags=junior, level=2, active=true, score=8.5, status=active, department=sales}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, score=9.2, status=active, level=5, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, tags=manager, level=7, salary=85000, score=7.8, status=active, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, level=3, salary=55000, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, department=qa, level=1, name=Frank, age=0, salary=-5000, score=0, tags=test, active=true, status=active}
7. [1] TestPerson{id=P007, department=management, name=Grace, score=10, tags=executive, level=9, age=65, salary=95000, active=true, status=active}
8. [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, name=Henry, active=false, tags=junior, status=inactive, level=1}
9. [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, active=true, status=active, department=intern, level=1}
11. [1] TestOrder{id=O001, product_id=PROD001, date=2024-01-15, status=pending, discount=50, region=north, customer_id=P001, amount=2, total=1999.98, priority=normal}
12. [1] TestOrder{id=O002, total=25.5, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0, product_id=PROD002, amount=1}
13. [1] TestOrder{id=O003, date=2024-02-01, priority=high, region=north, customer_id=P001, product_id=PROD003, amount=3, status=shipped, discount=15, total=225}
14. [1] TestOrder{id=O004, customer_id=P004, amount=1, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, total=299.99, priority=normal, region=east}
15. [1] TestOrder{id=O005, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed, customer_id=P002, date=2024-02-10}
16. [1] TestOrder{id=O006, amount=2, date=2024-02-15, discount=0, region=west, total=999.98, status=cancelled, priority=low, customer_id=P005, product_id=PROD005}
17. [1] TestOrder{id=O007, product_id=PROD006, total=600, priority=urgent, region=north, customer_id=P007, amount=4, date=2024-03-01, status=shipped, discount=50}
18. [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal}
19. [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1, total=89.99}
20. [1] TestOrder{id=O010, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15, total=75000, status=refunded, priority=urgent, discount=0}

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

1. [1] TestPerson{id=P001, tags=junior, level=2, active=true, score=8.5, status=active, department=sales, name=Alice, age=25, salary=45000}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, score=9.2, status=active, level=5, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, tags=manager, level=7, salary=85000, score=7.8, status=active, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, level=3, salary=55000, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, status=active, department=qa, level=1, name=Frank, age=0, salary=-5000, score=0, tags=test, active=true}
7. [1] TestPerson{id=P007, name=Grace, score=10, tags=executive, level=9, age=65, salary=95000, active=true, status=active, department=management}
8. [1] TestPerson{id=P008, status=inactive, level=1, age=18, salary=25000, score=5.5, department=support, name=Henry, active=false, tags=junior}
9. [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, tags=temp, active=true, status=active, department=intern, level=1, name=X, age=22, salary=28000, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, status=active, department=sales, name=Alice, age=25, salary=45000, tags=junior, level=2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, active=true, score=9.2, status=active, level=5, tags=senior, department=engineering}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, status=active, department=marketing, name=Diana, age=45, active=true, tags=manager, level=7, salary=85000, score=7.8}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, level=9, age=65, salary=95000, active=true, status=active, department=management, name=Grace}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, name=Henry, active=false, tags=junior, status=inactive, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, tags=temp, active=true, status=active, department=intern, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, department=hr, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, department=sales, name=Eve, age=30, active=false, score=8, level=3}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior, age=40, active=true, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, priority=normal, product_id=PROD001, date=2024-01-15, status=pending, discount=50, region=north}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low, discount=0}
3. [1] TestOrder{id=O003, region=north, customer_id=P001, product_id=PROD003, amount=3, status=shipped, discount=15, total=225, date=2024-02-01, priority=high}
4. [1] TestOrder{id=O004, status=delivered, discount=0, product_id=PROD004, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05}
5. [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed}
6. [1] TestOrder{id=O006, priority=low, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, discount=0, region=west, total=999.98, status=cancelled}
7. [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, discount=50, product_id=PROD006, total=600, priority=urgent, region=north, customer_id=P007}
8. [1] TestOrder{id=O008, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002, amount=10, total=255, priority=normal}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north, customer_id=P001, product_id=PROD007, amount=1}
10. [1] TestOrder{id=O010, total=75000, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, total=1999.98, priority=normal, product_id=PROD001, date=2024-01-15, status=pending, discount=50, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, discount=0, product_id=PROD002, amount=1, total=25.5, region=south, customer_id=P002, date=2024-01-20, status=confirmed, priority=low}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, priority=normal, region=east, customer_id=P004, amount=1, date=2024-02-05, status=delivered, discount=0, product_id=PROD004, total=299.99}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, date=2024-02-10, priority=high, discount=100, region=south, product_id=PROD001, amount=1, total=999.99, status=confirmed}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, total=600, priority=urgent, region=north, customer_id=P007, amount=4, date=2024-03-01, status=shipped, discount=50}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, status=shipped, discount=15, total=225, date=2024-02-01, priority=high, region=north, customer_id=P001}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O006, status=cancelled, priority=low, customer_id=P005, product_id=PROD005, amount=2, date=2024-02-15, discount=0, region=west, total=999.98}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, priority=normal, customer_id=P010, date=2024-03-05, status=pending, discount=0, region=south, product_id=PROD002, amount=10, total=255}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, total=75000, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, product_id=PROD001, amount=1, date=2024-03-15}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, tags=junior, level=2, active=true, score=8.5, status=active, department=sales, name=Alice, age=25}
2. [1] TestPerson{id=P002, department=engineering, name=Bob, age=35, salary=75000, active=true, score=9.2, status=active, level=5, tags=senior}
3. [1] TestPerson{id=P003, department=hr, name=Charlie, age=16, level=1, salary=0, active=false, score=6, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, salary=85000, score=7.8, status=active, department=marketing, name=Diana, age=45, active=true, tags=manager, level=7}
5. [1] TestPerson{id=P005, department=sales, name=Eve, age=30, active=false, score=8, level=3, salary=55000, tags=employee, status=inactive}
6. [1] TestPerson{id=P006, status=active, department=qa, level=1, name=Frank, age=0, salary=-5000, score=0, tags=test, active=true}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, status=active, department=management, name=Grace, score=10, tags=executive, level=9}
8. [1] TestPerson{id=P008, age=18, salary=25000, score=5.5, department=support, name=Henry, active=false, tags=junior, status=inactive, level=1}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, age=40, active=true, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, tags=temp, active=true, status=active, department=intern, level=1, name=X}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, department=engineering, name=Bob, age=35, salary=75000, active=true, score=9.2, status=active, level=5}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, name=Charlie, age=16, level=1}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, level=3, salary=55000, tags=employee, status=inactive, department=sales}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, active=true, status=active, department=qa, level=1, name=Frank, age=0, salary=-5000, score=0}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, status=active, department=management, name=Grace, score=10, tags=executive, level=9}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P008, active=false, tags=junior, status=inactive, level=1, age=18, salary=25000, score=5.5, department=support, name=Henry}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, active=true, status=active, department=intern, level=1, name=X, age=22, salary=28000, score=6.5, tags=temp}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, status=active, department=sales, name=Alice, age=25, salary=45000, tags=junior, level=2}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, active=true, tags=manager, level=7, salary=85000, score=7.8, status=active, department=marketing}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P009, level=6, name=Ivy, salary=68000, score=8.7, tags=senior, age=40, active=true, status=active, department=engineering}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 17 (89.5%)
- **Tokens générés**: 130
- **Faits traités**: 27
