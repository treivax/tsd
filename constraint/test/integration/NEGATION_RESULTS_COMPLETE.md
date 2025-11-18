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

1. [1] TestPerson{id=P001, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active, department=sales, name=Alice, active=true}
2. [1] TestPerson{id=P002, active=true, status=active, department=engineering, level=5, age=35, score=9.2, tags=senior, name=Bob, salary=75000}
3. [1] TestPerson{id=P003, status=inactive, department=hr, name=Charlie, salary=0, tags=intern, level=1, age=16, active=false, score=6}
4. [1] TestPerson{id=P004, age=45, active=true, score=7.8, tags=manager, status=active, level=7, salary=85000, department=marketing, name=Diana}
5. [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee, level=3}
6. [1] TestPerson{id=P006, status=active, department=qa, name=Frank, age=0, salary=-5000, score=0, tags=test, level=1, active=true}
7. [1] TestPerson{id=P007, tags=executive, department=management, name=Grace, salary=95000, status=active, level=9, age=65, active=true, score=10}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, level=1, salary=25000, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, active=true, score=8.7, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, status=active, level=6}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, tags=temp, department=intern, active=true, score=6.5, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, status=active, department=engineering, level=5, age=35, score=9.2, tags=senior, name=Bob}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, status=inactive, department=hr, name=Charlie, salary=0, tags=intern, level=1}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, score=7.8, tags=manager, status=active, level=7, salary=85000, department=marketing, name=Diana}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, tags=senior, status=active, level=6, active=true, score=8.7, department=engineering, name=Ivy}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, department=management, name=Grace, salary=95000, status=active, level=9, age=65, active=true}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, tags=junior, status=inactive, department=support, name=Henry, age=18, active=false, score=5.5, level=1}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, tags=temp, department=intern, active=true, score=6.5, status=active, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, priority=normal, discount=50, product_id=PROD001, amount=2, date=2024-01-15, region=north, customer_id=P001, total=1999.98, status=pending}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, discount=0, amount=1, date=2024-01-20, status=confirmed, priority=low, region=south}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, customer_id=P001, status=shipped, discount=15, region=north}
4. [1] TestOrder{id=O004, priority=normal, region=east, amount=1, total=299.99, date=2024-02-05, discount=0, customer_id=P004, product_id=PROD004, status=delivered}
5. [1] TestOrder{id=O005, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, status=confirmed, product_id=PROD001, amount=1, total=999.99}
6. [1] TestOrder{id=O006, amount=2, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, date=2024-02-15, product_id=PROD005}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}
8. [1] TestOrder{id=O008, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, total=255}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, amount=1, region=north}
10. [1] TestOrder{id=O010, region=east, product_id=PROD001, amount=1, total=75000, status=refunded, customer_id=P006, date=2024-03-15, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, total=25.5, discount=0, amount=1, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, total=255, priority=normal, discount=0, region=south, customer_id=P010}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, amount=1, region=north}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, amount=1, total=75000, status=refunded, customer_id=P006, date=2024-03-15, priority=urgent, discount=0, region=east}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O001, date=2024-01-15, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, discount=50, product_id=PROD001, amount=2}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, priority=high, customer_id=P001, status=shipped, discount=15, region=north, product_id=PROD003}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, region=east, amount=1, total=299.99, date=2024-02-05, discount=0, customer_id=P004, product_id=PROD004, status=delivered, priority=normal}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, status=confirmed, product_id=PROD001, amount=1}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O007, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, status=active, department=sales, name=Alice, active=true, score=8.5, tags=junior, level=2, age=25}
2. [1] TestPerson{id=P002, age=35, score=9.2, tags=senior, name=Bob, salary=75000, active=true, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, tags=intern, level=1, age=16, active=false, score=6, status=inactive, department=hr, name=Charlie, salary=0}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, level=7, salary=85000, department=marketing}
5. [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee, level=3}
6. [1] TestPerson{id=P006, score=0, tags=test, level=1, active=true, status=active, department=qa, name=Frank, age=0, salary=-5000}
7. [1] TestPerson{id=P007, status=active, level=9, age=65, active=true, score=10, tags=executive, department=management, name=Grace, salary=95000}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, level=1, salary=25000, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, score=8.7, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, status=active, level=6, active=true}
10. [1] TestPerson{id=P010, status=active, level=1, name=X, age=22, salary=28000, tags=temp, department=intern, active=true, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, active=true, status=active, department=engineering, level=5, age=35, score=9.2, tags=senior, name=Bob, salary=75000}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, status=active, level=7, salary=85000, department=marketing, name=Diana, age=45, active=true, score=7.8}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, status=active, level=9, age=65, active=true, score=10, tags=executive, department=management}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, status=active, level=6, active=true, score=8.7, department=engineering}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, status=pending, priority=normal, discount=50, product_id=PROD001, amount=2, date=2024-01-15, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, discount=0, amount=1, date=2024-01-20, status=confirmed, priority=low, region=south}
3. [1] TestOrder{id=O003, region=north, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, customer_id=P001, status=shipped, discount=15}
4. [1] TestOrder{id=O004, amount=1, total=299.99, date=2024-02-05, discount=0, customer_id=P004, product_id=PROD004, status=delivered, priority=normal, region=east}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, status=confirmed}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, date=2024-02-15}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007}
8. [1] TestOrder{id=O008, status=pending, total=255, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05}
9. [1] TestOrder{id=O009, amount=1, region=north, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, total=75000, status=refunded, customer_id=P006, date=2024-03-15, priority=urgent, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, date=2024-02-15}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, total=255, priority=normal, discount=0, region=south}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O003, region=north, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, customer_id=P001, status=shipped, discount=15}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, discount=10, amount=1, region=north, customer_id=P001, product_id=PROD007, total=89.99}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, discount=50}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, region=south, customer_id=P002, product_id=PROD002, total=25.5, discount=0, amount=1, date=2024-01-20, status=confirmed, priority=low}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, priority=normal, region=east, amount=1, total=299.99, date=2024-02-05, discount=0, customer_id=P004, product_id=PROD004, status=delivered}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O005, status=confirmed, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, available=true, rating=4.5, keywords=computer, stock=50, brand=TechCorp, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, brand=TechCorp, stock=200, name=Mouse, available=true, supplier=TechSupply, category=accessories, price=25.5, rating=4, keywords=peripheral}
3. [1] TestProduct{id=PROD003, category=accessories, supplier=KeySupply, price=75, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0, name=Keyboard}
4. [1] TestProduct{id=PROD004, brand=ScreenPro, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics, available=true, rating=4.8, price=299.99, keywords=display}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, price=8.5, available=false, rating=2, keywords=obsolete, stock=0, supplier=OldSupply, category=accessories, brand=OldTech}
6. [1] TestProduct{id=PROD006, name=Headphones, available=true, rating=4.6, keywords=sound, stock=75, supplier=AudioSupply, category=audio, price=150, brand=AudioMax}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, price=89.99, available=true, stock=25, supplier=CamSupply, rating=3.8, keywords=video, brand=CamTech}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD003, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0, name=Keyboard, category=accessories, supplier=KeySupply, price=75}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD004, category=electronics, available=true, rating=4.8, price=299.99, keywords=display, brand=ScreenPro, stock=30, supplier=ScreenSupply, name=Monitor}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD006, category=audio, price=150, brand=AudioMax, name=Headphones, available=true, rating=4.6, keywords=sound, stock=75, supplier=AudioSupply}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD007, keywords=video, brand=CamTech, name=Webcam, category=electronics, price=89.99, available=true, stock=25, supplier=CamSupply, rating=3.8}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, available=true, rating=4.5, keywords=computer, stock=50, brand=TechCorp, supplier=TechSupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD002, stock=200, name=Mouse, available=true, supplier=TechSupply, category=accessories, price=25.5, rating=4, keywords=peripheral, brand=TechCorp}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, status=active, department=sales, name=Alice, active=true, score=8.5, tags=junior, level=2}
2. [1] TestPerson{id=P002, age=35, score=9.2, tags=senior, name=Bob, salary=75000, active=true, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, age=16, active=false, score=6, status=inactive, department=hr, name=Charlie, salary=0, tags=intern, level=1}
4. [1] TestPerson{id=P004, level=7, salary=85000, department=marketing, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active}
5. [1] TestPerson{id=P005, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee, level=3, age=30}
6. [1] TestPerson{id=P006, score=0, tags=test, level=1, active=true, status=active, department=qa, name=Frank, age=0, salary=-5000}
7. [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, department=management, name=Grace, salary=95000, status=active, level=9}
8. [1] TestPerson{id=P008, status=inactive, department=support, name=Henry, age=18, active=false, score=5.5, level=1, salary=25000, tags=junior}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, status=active, level=6, active=true, score=8.7, department=engineering}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, tags=temp, department=intern, active=true, score=6.5, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, level=1, age=16, active=false, score=6, status=inactive, department=hr, name=Charlie}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, level=1, active=true, status=active, department=qa}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, total=1999.98, status=pending, priority=normal, discount=50, product_id=PROD001, amount=2, date=2024-01-15, region=north, customer_id=P001}
2. [1] TestOrder{id=O002, total=25.5, discount=0, amount=1, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002}
3. [1] TestOrder{id=O003, region=north, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, customer_id=P001, status=shipped, discount=15}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, status=delivered, priority=normal, region=east, amount=1, total=299.99, date=2024-02-05, discount=0}
5. [1] TestOrder{id=O005, total=999.99, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, status=confirmed, product_id=PROD001, amount=1}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, date=2024-02-15}
7. [1] TestOrder{id=O007, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600}
8. [1] TestOrder{id=O008, total=255, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, amount=1, region=north}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, total=75000, status=refunded, customer_id=P006, date=2024-03-15, priority=urgent, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O008, region=south, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, total=255, priority=normal, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002, total=25.5, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, status=shipped, discount=15, region=north, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, customer_id=P001}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, region=north, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, amount=1}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, date=2024-03-15, priority=urgent, discount=0, region=east, product_id=PROD001, amount=1, total=75000, status=refunded}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, discount=50, product_id=PROD001, amount=2, date=2024-01-15}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, status=delivered, priority=normal, region=east, amount=1, total=299.99, date=2024-02-05, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, status=confirmed}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, date=2024-02-15}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active, department=sales, name=Alice, active=true}
2. [1] TestPerson{id=P002, age=35, score=9.2, tags=senior, name=Bob, salary=75000, active=true, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, age=16, active=false, score=6, status=inactive, department=hr, name=Charlie, salary=0, tags=intern, level=1}
4. [1] TestPerson{id=P004, status=active, level=7, salary=85000, department=marketing, name=Diana, age=45, active=true, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}
6. [1] TestPerson{id=P006, salary=-5000, score=0, tags=test, level=1, active=true, status=active, department=qa, name=Frank, age=0}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, status=active, level=9, age=65, active=true, score=10, tags=executive, department=management}
8. [1] TestPerson{id=P008, salary=25000, tags=junior, status=inactive, department=support, name=Henry, age=18, active=false, score=5.5, level=1}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, status=active, level=6, active=true, score=8.7, department=engineering}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, tags=temp, department=intern, active=true, score=6.5, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, level=1, salary=25000, tags=junior, status=inactive, department=support}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, status=inactive, department=hr, name=Charlie, salary=0, tags=intern, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, status=pending, priority=normal, discount=50, product_id=PROD001, amount=2, date=2024-01-15, region=north}
2. [1] TestOrder{id=O002, product_id=PROD002, total=25.5, discount=0, amount=1, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, customer_id=P001, status=shipped, discount=15, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, status=delivered, priority=normal, region=east, amount=1, total=299.99, date=2024-02-05, discount=0}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, status=confirmed}
6. [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, date=2024-02-15, product_id=PROD005, amount=2, total=999.98}
7. [1] TestOrder{id=O007, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600}
8. [1] TestOrder{id=O008, total=255, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, amount=1, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, date=2024-03-15, priority=urgent, discount=0, region=east, product_id=PROD001, amount=1, total=75000, status=refunded}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, amount=1, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O010, total=75000, status=refunded, customer_id=P006, date=2024-03-15, priority=urgent, discount=0, region=east, product_id=PROD001, amount=1}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O002, discount=0, amount=1, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002, total=25.5}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O003, region=north, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, customer_id=P001, status=shipped, discount=15}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O004, date=2024-02-05, discount=0, customer_id=P004, product_id=PROD004, status=delivered, priority=normal, region=east, amount=1, total=299.99}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O005, discount=100, region=south, customer_id=P002, status=confirmed, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, date=2024-02-15}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, age=25, salary=45000, status=active, department=sales, name=Alice, active=true, score=8.5, tags=junior}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, name=Bob, salary=75000, active=true, status=active, department=engineering, level=5, age=35}
3. [1] TestPerson{id=P003, status=inactive, department=hr, name=Charlie, salary=0, tags=intern, level=1, age=16, active=false, score=6}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, level=7, salary=85000, department=marketing}
5. [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, level=1, active=true, status=active, department=qa}
7. [1] TestPerson{id=P007, status=active, level=9, age=65, active=true, score=10, tags=executive, department=management, name=Grace, salary=95000}
8. [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, name=Henry, age=18, active=false, score=5.5, level=1, salary=25000}
9. [1] TestPerson{id=P009, status=active, level=6, active=true, score=8.7, department=engineering, name=Ivy, age=40, salary=68000, tags=senior}
10. [1] TestPerson{id=P010, salary=28000, tags=temp, department=intern, active=true, score=6.5, status=active, level=1, name=X, age=22}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, level=1, active=true, status=active, department=qa}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, department=support, name=Henry, age=18, active=false, score=5.5, level=1, salary=25000, tags=junior}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, status=active, department=engineering, level=5, age=35, score=9.2, tags=senior}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, level=7, salary=85000, department=marketing}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, department=management, name=Grace, salary=95000, status=active, level=9}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, status=active, level=6, active=true, score=8.7, department=engineering}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, name=Alice, active=true, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, salary=0, tags=intern, level=1, age=16, active=false, score=6, status=inactive, department=hr}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, available=true, rating=4.5, keywords=computer, stock=50, brand=TechCorp, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, name=Mouse, available=true, supplier=TechSupply, category=accessories, price=25.5, rating=4, keywords=peripheral, brand=TechCorp, stock=200}
3. [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, supplier=KeySupply, price=75, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0}
4. [1] TestProduct{id=PROD004, price=299.99, keywords=display, brand=ScreenPro, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics, available=true, rating=4.8}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, price=8.5, available=false, rating=2, keywords=obsolete, stock=0, supplier=OldSupply, category=accessories, brand=OldTech}
6. [1] TestProduct{id=PROD006, available=true, rating=4.6, keywords=sound, stock=75, supplier=AudioSupply, category=audio, price=150, brand=AudioMax, name=Headphones}
7. [1] TestProduct{id=PROD007, stock=25, supplier=CamSupply, rating=3.8, keywords=video, brand=CamTech, name=Webcam, category=electronics, price=89.99, available=true}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, available=true, rating=4.5, keywords=computer, stock=50, brand=TechCorp, supplier=TechSupply}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, available=true, supplier=TechSupply, category=accessories, price=25.5, rating=4, keywords=peripheral, brand=TechCorp, stock=200}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, brand=KeyTech, stock=0, name=Keyboard, category=accessories, supplier=KeySupply, price=75, available=false, rating=3.5, keywords=typing}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, price=299.99, keywords=display, brand=ScreenPro, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics, available=true, rating=4.8}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, category=audio, price=150, brand=AudioMax, name=Headphones, available=true, rating=4.6, keywords=sound, stock=75, supplier=AudioSupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, category=electronics, price=89.99, available=true, stock=25, supplier=CamSupply, rating=3.8, keywords=video, brand=CamTech, name=Webcam}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active, department=sales}
2. [1] TestPerson{id=P002, age=35, score=9.2, tags=senior, name=Bob, salary=75000, active=true, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, age=16, active=false, score=6, status=inactive, department=hr, name=Charlie, salary=0, tags=intern, level=1}
4. [1] TestPerson{id=P004, active=true, score=7.8, tags=manager, status=active, level=7, salary=85000, department=marketing, name=Diana, age=45}
5. [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, level=1, active=true, status=active, department=qa}
7. [1] TestPerson{id=P007, level=9, age=65, active=true, score=10, tags=executive, department=management, name=Grace, salary=95000, status=active}
8. [1] TestPerson{id=P008, department=support, name=Henry, age=18, active=false, score=5.5, level=1, salary=25000, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, status=active, level=6, active=true, score=8.7, department=engineering, name=Ivy, age=40, salary=68000, tags=senior}
10. [1] TestPerson{id=P010, level=1, name=X, age=22, salary=28000, tags=temp, department=intern, active=true, score=6.5, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, status=inactive, department=hr, name=Charlie, salary=0, tags=intern, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, status=active, level=9, age=65, active=true, score=10, tags=executive, department=management}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, active=true, score=8.7, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, status=active, level=6}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, status=active, department=sales, name=Alice, active=true, score=8.5, tags=junior, level=2, age=25, salary=45000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, age=35, score=9.2, tags=senior, name=Bob, salary=75000, active=true, status=active, department=engineering, level=5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, level=7, salary=85000}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, score=0, tags=test, level=1, active=true, status=active, department=qa, name=Frank, age=0}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, level=1, salary=25000, tags=junior, status=inactive, department=support, name=Henry, age=18, active=false}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, amount=2, date=2024-01-15, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, discount=50, product_id=PROD001}
2. [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002, total=25.5, discount=0}
3. [1] TestOrder{id=O003, date=2024-02-01, priority=high, customer_id=P001, status=shipped, discount=15, region=north, product_id=PROD003, amount=3, total=225}
4. [1] TestOrder{id=O004, total=299.99, date=2024-02-05, discount=0, customer_id=P004, product_id=PROD004, status=delivered, priority=normal, region=east, amount=1}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, status=confirmed}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, date=2024-02-15}
7. [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006}
8. [1] TestOrder{id=O008, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, total=255, priority=normal}
9. [1] TestOrder{id=O009, region=north, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, amount=1}
10. [1] TestOrder{id=O010, region=east, product_id=PROD001, amount=1, total=75000, status=refunded, customer_id=P006, date=2024-03-15, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, total=299.99, date=2024-02-05, discount=0, customer_id=P004, product_id=PROD004, status=delivered, priority=normal, region=east}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, status=confirmed}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, date=2024-03-05, status=pending, total=255, priority=normal, discount=0, region=south, customer_id=P010}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, amount=1, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, status=pending, priority=normal, discount=50, product_id=PROD001, amount=2, date=2024-01-15, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, priority=low, region=south, customer_id=P002, product_id=PROD002, total=25.5, discount=0, amount=1, date=2024-01-20, status=confirmed}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, priority=high, customer_id=P001, status=shipped, discount=15, region=north, product_id=PROD003, amount=3, total=225, date=2024-02-01}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, status=active, department=sales, name=Alice, active=true, score=8.5, tags=junior, level=2}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, status=active, department=engineering, level=5, age=35, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, tags=intern, level=1, age=16, active=false, score=6, status=inactive, department=hr, name=Charlie, salary=0}
4. [1] TestPerson{id=P004, tags=manager, status=active, level=7, salary=85000, department=marketing, name=Diana, age=45, active=true, score=7.8}
5. [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}
6. [1] TestPerson{id=P006, active=true, status=active, department=qa, name=Frank, age=0, salary=-5000, score=0, tags=test, level=1}
7. [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, department=management, name=Grace, salary=95000, status=active, level=9}
8. [1] TestPerson{id=P008, level=1, salary=25000, tags=junior, status=inactive, department=support, name=Henry, age=18, active=false, score=5.5}
9. [1] TestPerson{id=P009, status=active, level=6, active=true, score=8.7, department=engineering, name=Ivy, age=40, salary=68000, tags=senior}
10. [1] TestPerson{id=P010, status=active, level=1, name=X, age=22, salary=28000, tags=temp, department=intern, active=true, score=6.5}
11. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, status=pending, priority=normal, discount=50, product_id=PROD001, amount=2, date=2024-01-15, region=north}
12. [1] TestOrder{id=O002, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002, total=25.5, discount=0, amount=1}
13. [1] TestOrder{id=O003, date=2024-02-01, priority=high, customer_id=P001, status=shipped, discount=15, region=north, product_id=PROD003, amount=3, total=225}
14. [1] TestOrder{id=O004, amount=1, total=299.99, date=2024-02-05, discount=0, customer_id=P004, product_id=PROD004, status=delivered, priority=normal, region=east}
15. [1] TestOrder{id=O005, customer_id=P002, status=confirmed, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, region=south}
16. [1] TestOrder{id=O006, date=2024-02-15, product_id=PROD005, amount=2, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005}
17. [1] TestOrder{id=O007, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north, customer_id=P007, product_id=PROD006, amount=4, total=600}
18. [1] TestOrder{id=O008, total=255, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending}
19. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, amount=1, region=north}
20. [1] TestOrder{id=O010, date=2024-03-15, priority=urgent, discount=0, region=east, product_id=PROD001, amount=1, total=75000, status=refunded, customer_id=P006}

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

1. [1] TestPerson{id=P001, salary=45000, status=active, department=sales, name=Alice, active=true, score=8.5, tags=junior, level=2, age=25}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, status=active, department=engineering, level=5, age=35, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, age=16, active=false, score=6, status=inactive, department=hr, name=Charlie, salary=0, tags=intern, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, level=7, salary=85000, department=marketing}
5. [1] TestPerson{id=P005, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve}
6. [1] TestPerson{id=P006, tags=test, level=1, active=true, status=active, department=qa, name=Frank, age=0, salary=-5000, score=0}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, status=active, level=9, age=65, active=true, score=10, tags=executive, department=management}
8. [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, name=Henry, age=18, active=false, score=5.5, level=1, salary=25000}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, status=active, level=6, active=true, score=8.7, department=engineering}
10. [1] TestPerson{id=P010, salary=28000, tags=temp, department=intern, active=true, score=6.5, status=active, level=1, name=X, age=22}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, department=management, name=Grace, salary=95000, status=active, level=9}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, status=active, level=6, active=true, score=8.7, department=engineering}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, status=active, level=1, name=X, age=22, salary=28000, tags=temp, department=intern, active=true, score=6.5}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active, department=sales, name=Alice, active=true}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, status=active, department=engineering, level=5, age=35, score=9.2, tags=senior, name=Bob}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, tags=manager, status=active, level=7, salary=85000, department=marketing, name=Diana, age=45}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, score=0, tags=test, level=1, active=true, status=active, department=qa}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, age=25, salary=45000, status=active, department=sales, name=Alice, active=true, score=8.5, tags=junior}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, status=active, department=engineering, level=5, age=35, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, salary=0, tags=intern, level=1, age=16, active=false, score=6, status=inactive, department=hr, name=Charlie}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, level=7, salary=85000, department=marketing}
5. [1] TestPerson{id=P005, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales, name=Eve}
6. [1] TestPerson{id=P006, status=active, department=qa, name=Frank, age=0, salary=-5000, score=0, tags=test, level=1, active=true}
7. [1] TestPerson{id=P007, score=10, tags=executive, department=management, name=Grace, salary=95000, status=active, level=9, age=65, active=true}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, level=1, salary=25000, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, tags=senior, status=active, level=6, active=true, score=8.7, department=engineering, name=Ivy, age=40, salary=68000}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, tags=temp, department=intern, active=true, score=6.5, status=active, level=1}
11. [1] TestOrder{id=O001, priority=normal, discount=50, product_id=PROD001, amount=2, date=2024-01-15, region=north, customer_id=P001, total=1999.98, status=pending}
12. [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002, total=25.5, discount=0}
13. [1] TestOrder{id=O003, priority=high, customer_id=P001, status=shipped, discount=15, region=north, product_id=PROD003, amount=3, total=225, date=2024-02-01}
14. [1] TestOrder{id=O004, discount=0, customer_id=P004, product_id=PROD004, status=delivered, priority=normal, region=east, amount=1, total=299.99, date=2024-02-05}
15. [1] TestOrder{id=O005, status=confirmed, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002}
16. [1] TestOrder{id=O006, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, date=2024-02-15, product_id=PROD005, amount=2}
17. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}
18. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, total=255, priority=normal, discount=0, region=south}
19. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, amount=1, region=north}
20. [1] TestOrder{id=O010, customer_id=P006, date=2024-03-15, priority=urgent, discount=0, region=east, product_id=PROD001, amount=1, total=75000, status=refunded}

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

1. [1] TestPerson{id=P001, name=Alice, active=true, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active, department=sales}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, status=active, department=engineering, level=5, age=35, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, score=6, status=inactive, department=hr, name=Charlie, salary=0, tags=intern, level=1, age=16, active=false}
4. [1] TestPerson{id=P004, status=active, level=7, salary=85000, department=marketing, name=Diana, age=45, active=true, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, age=30, salary=55000, status=inactive, department=sales, name=Eve, active=false, score=8, tags=employee, level=3}
6. [1] TestPerson{id=P006, active=true, status=active, department=qa, name=Frank, age=0, salary=-5000, score=0, tags=test, level=1}
7. [1] TestPerson{id=P007, active=true, score=10, tags=executive, department=management, name=Grace, salary=95000, status=active, level=9, age=65}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, level=1, salary=25000, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, active=true, score=8.7, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, status=active, level=6}
10. [1] TestPerson{id=P010, status=active, level=1, name=X, age=22, salary=28000, tags=temp, department=intern, active=true, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, department=support, name=Henry, age=18, active=false, score=5.5, level=1, salary=25000, tags=junior}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, status=inactive, department=hr, name=Charlie, salary=0, tags=intern, level=1}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, level=7, salary=85000, department=marketing, name=Diana, age=45, active=true}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, status=active, level=9, age=65, active=true, score=10, tags=executive, department=management, name=Grace, salary=95000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, active=true, score=8.7, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, status=active, level=6}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, tags=temp, department=intern, active=true, score=6.5, status=active, level=1, name=X, age=22}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, name=Alice, active=true, score=8.5, tags=junior, level=2, age=25, salary=45000, status=active}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, status=active, department=engineering, level=5, age=35, score=9.2, tags=senior}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, status=pending, priority=normal, discount=50, product_id=PROD001, amount=2, date=2024-01-15, region=north}
2. [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002, total=25.5, discount=0}
3. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high, customer_id=P001, status=shipped, discount=15, region=north}
4. [1] TestOrder{id=O004, discount=0, customer_id=P004, product_id=PROD004, status=delivered, priority=normal, region=east, amount=1, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, status=confirmed}
6. [1] TestOrder{id=O006, amount=2, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, date=2024-02-15, product_id=PROD005}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}
8. [1] TestOrder{id=O008, region=south, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, total=255, priority=normal, discount=0}
9. [1] TestOrder{id=O009, amount=1, region=north, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10}
10. [1] TestOrder{id=O010, customer_id=P006, date=2024-03-15, priority=urgent, discount=0, region=east, product_id=PROD001, amount=1, total=75000, status=refunded}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, status=shipped, discount=15, region=north, product_id=PROD003, amount=3, total=225, date=2024-02-01, priority=high}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, date=2024-02-05, discount=0, customer_id=P004, product_id=PROD004, status=delivered, priority=normal, region=east, amount=1, total=299.99}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, total=255, priority=normal, discount=0, region=south}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, status=completed, priority=low, discount=10, amount=1, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O005, region=south, customer_id=P002, status=confirmed, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O006, amount=2, total=999.98, status=cancelled, priority=low, discount=0, region=west, customer_id=P005, date=2024-02-15, product_id=PROD005}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, priority=urgent, discount=50, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, amount=1, total=75000, status=refunded, customer_id=P006, date=2024-03-15, priority=urgent, discount=0, region=east}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, discount=50}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, product_id=PROD002, total=25.5, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, age=25, salary=45000, status=active, department=sales, name=Alice, active=true, score=8.5, tags=junior}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, status=active, department=engineering, level=5, age=35, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, age=16, active=false, score=6, status=inactive, department=hr, name=Charlie, salary=0, tags=intern, level=1}
4. [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, level=7, salary=85000, department=marketing, name=Diana, age=45, active=true}
5. [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}
6. [1] TestPerson{id=P006, status=active, department=qa, name=Frank, age=0, salary=-5000, score=0, tags=test, level=1, active=true}
7. [1] TestPerson{id=P007, level=9, age=65, active=true, score=10, tags=executive, department=management, name=Grace, salary=95000, status=active}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, level=1, salary=25000, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, name=Ivy, age=40, salary=68000, tags=senior, status=active, level=6, active=true, score=8.7, department=engineering}
10. [1] TestPerson{id=P010, score=6.5, status=active, level=1, name=X, age=22, salary=28000, tags=temp, department=intern, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, department=management, name=Grace, salary=95000, status=active, level=9, age=65, active=true, score=10, tags=executive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, department=support, name=Henry, age=18, active=false, score=5.5, level=1, salary=25000, tags=junior}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, tags=temp, department=intern, active=true, score=6.5, status=active, level=1}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, status=active, department=engineering, level=5, age=35, score=9.2, tags=senior}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, status=active, level=7, salary=85000, department=marketing, name=Diana, age=45, active=true, score=7.8}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, name=Frank, age=0, salary=-5000, score=0, tags=test, level=1, active=true, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P009, active=true, score=8.7, department=engineering, name=Ivy, age=40, salary=68000, tags=senior, status=active, level=6}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, status=active, department=sales, name=Alice, active=true, score=8.5, tags=junior, level=2}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P003, age=16, active=false, score=6, status=inactive, department=hr, name=Charlie, salary=0, tags=intern, level=1}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, active=false, score=8, tags=employee, level=3, age=30, salary=55000, status=inactive, department=sales}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 17 (89.5%)
- **Tokens générés**: 130
- **Faits traités**: 27
