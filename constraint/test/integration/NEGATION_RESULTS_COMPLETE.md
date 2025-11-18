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

1. [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
2. [1] TestPerson{id=P002, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, score=6, level=1, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false}
4. [1] TestPerson{id=P004, tags=manager, department=marketing, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true}
5. [1] TestPerson{id=P005, active=false, department=sales, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3}
6. [1] TestPerson{id=P006, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa}
7. [1] TestPerson{id=P007, department=management, level=9, age=65, salary=95000, status=active, name=Grace, active=true, score=10, tags=executive}
8. [1] TestPerson{id=P008, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support, name=Henry}
9. [1] TestPerson{id=P009, level=6, name=Ivy, active=true, status=active, age=40, salary=68000, score=8.7, tags=senior, department=engineering}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, status=active, level=7, active=true, tags=manager, department=marketing, name=Diana, age=45, salary=85000, score=7.8}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, status=inactive, level=3, active=false, department=sales, name=Eve, age=30, salary=55000, score=8}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active, name=Bob, salary=75000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5, name=X}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, date=2024-01-15}
2. [1] TestOrder{id=O002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}
4. [1] TestOrder{id=O004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east, product_id=PROD004}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10}
6. [1] TestOrder{id=O006, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}
9. [1] TestOrder{id=O009, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low}
10. [1] TestOrder{id=O010, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, region=east, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O010, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, region=south, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, region=north, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50, customer_id=P001, product_id=PROD001}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, department=sales, level=2, age=25, active=true, name=Alice, salary=45000, score=8.5, tags=junior}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active, name=Bob, salary=75000, active=true}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}
5. [1] TestPerson{id=P005, score=8, tags=employee, status=inactive, level=3, active=false, department=sales, name=Eve, age=30, salary=55000}
6. [1] TestPerson{id=P006, score=0, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}
8. [1] TestPerson{id=P008, status=inactive, level=1, tags=junior, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5}
9. [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}
10. [1] TestPerson{id=P010, department=intern, level=1, age=22, score=6.5, name=X, salary=28000, active=true, tags=temp, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, level=2, age=25, active=true, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, level=9, age=65, salary=95000, status=active, name=Grace, active=true, score=10, tags=executive, department=management}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50, customer_id=P001, product_id=PROD001}
2. [1] TestOrder{id=O002, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20}
3. [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north}
4. [1] TestOrder{id=O004, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, total=299.99}
5. [1] TestOrder{id=O005, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001}
6. [1] TestOrder{id=O006, status=cancelled, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15}
7. [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent}
8. [1] TestOrder{id=O008, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1}
10. [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001, product_id=PROD003}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O006, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, amount=1}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, region=north, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, available=true, rating=4.5, stock=50, supplier=TechSupply, category=electronics, price=999.99, keywords=computer, brand=TechCorp}
2. [1] TestProduct{id=PROD002, available=true, stock=200, category=accessories, rating=4, keywords=peripheral, brand=TechCorp, supplier=TechSupply, name=Mouse, price=25.5}
3. [1] TestProduct{id=PROD003, name=Keyboard, rating=3.5, keywords=typing, category=accessories, price=75, available=false, brand=KeyTech, stock=0, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, brand=ScreenPro, category=electronics, price=299.99, available=true, rating=4.8, keywords=display, stock=30, supplier=ScreenSupply, name=Monitor}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, price=8.5, available=false, keywords=obsolete, brand=OldTech, rating=2, stock=0, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, supplier=AudioSupply, available=true, rating=4.6, name=Headphones, category=audio, price=150, keywords=sound, brand=AudioMax, stock=75}
7. [1] TestProduct{id=PROD007, rating=3.8, stock=25, supplier=CamSupply, name=Webcam, category=electronics, keywords=video, brand=CamTech, price=89.99, available=true}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD006, category=audio, price=150, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply, available=true, rating=4.6, name=Headphones}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD007, price=89.99, available=true, rating=3.8, stock=25, supplier=CamSupply, name=Webcam, category=electronics, keywords=video, brand=CamTech}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD001, category=electronics, price=999.99, keywords=computer, brand=TechCorp, name=Laptop, available=true, rating=4.5, stock=50, supplier=TechSupply}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, price=25.5, available=true, stock=200, category=accessories, rating=4, keywords=peripheral, brand=TechCorp, supplier=TechSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD003, price=75, available=false, brand=KeyTech, stock=0, supplier=KeySupply, name=Keyboard, rating=3.5, keywords=typing, category=accessories}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD004, category=electronics, price=299.99, available=true, rating=4.8, keywords=display, stock=30, supplier=ScreenSupply, name=Monitor, brand=ScreenPro}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support}
9. [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50}
2. [1] TestOrder{id=O002, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}
4. [1] TestOrder{id=O004, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled}
7. [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent}
8. [1] TestOrder{id=O008, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending}
9. [1] TestOrder{id=O009, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, discount=50, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, region=east, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, discount=100, region=south, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
2. [1] TestPerson{id=P002, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales}
6. [1] TestPerson{id=P006, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0}
7. [1] TestPerson{id=P007, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active, name=Grace, active=true}
8. [1] TestPerson{id=P008, level=1, tags=junior, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}
10. [1] TestPerson{id=P010, level=1, age=22, score=6.5, name=X, salary=28000, active=true, tags=temp, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, department=hr, active=false, score=6, level=1, name=Charlie, age=16, salary=0, tags=intern}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales, name=Eve}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, discount=50, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10, customer_id=P002}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north}
8. [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south}
9. [1] TestOrder{id=O009, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1, total=89.99}
10. [1] TestOrder{id=O010, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, status=delivered, region=east, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
2. [1] TestPerson{id=P002, department=engineering, level=5, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}
4. [1] TestPerson{id=P004, tags=manager, department=marketing, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales}
6. [1] TestPerson{id=P006, active=true, tags=test, status=active, score=0, department=qa, level=1, name=Frank, age=0, salary=-5000}
7. [1] TestPerson{id=P007, age=65, salary=95000, status=active, name=Grace, active=true, score=10, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support}
9. [1] TestPerson{id=P009, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active, age=40, salary=68000, score=8.7}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, status=inactive, level=3, active=false, department=sales, name=Eve, age=30, salary=55000, score=8}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, status=inactive, level=1, tags=junior, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, level=1, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P006, status=active, score=0, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active, name=Grace, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, category=electronics, price=999.99, keywords=computer, brand=TechCorp, name=Laptop, available=true, rating=4.5, stock=50, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, name=Mouse, price=25.5, available=true, stock=200, category=accessories, rating=4, keywords=peripheral, brand=TechCorp, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, category=accessories, price=75, available=false, brand=KeyTech, stock=0, supplier=KeySupply, name=Keyboard, rating=3.5, keywords=typing}
4. [1] TestProduct{id=PROD004, available=true, rating=4.8, keywords=display, stock=30, supplier=ScreenSupply, name=Monitor, brand=ScreenPro, category=electronics, price=299.99}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, price=8.5, available=false, keywords=obsolete, brand=OldTech, rating=2, stock=0, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, category=audio, price=150, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply, available=true, rating=4.6, name=Headphones}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, keywords=video, brand=CamTech, price=89.99, available=true, rating=3.8, stock=25, supplier=CamSupply}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD002, category=accessories, rating=4, keywords=peripheral, brand=TechCorp, supplier=TechSupply, name=Mouse, price=25.5, available=true, stock=200}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD003, name=Keyboard, rating=3.5, keywords=typing, category=accessories, price=75, available=false, brand=KeyTech, stock=0, supplier=KeySupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD004, category=electronics, price=299.99, available=true, rating=4.8, keywords=display, stock=30, supplier=ScreenSupply, name=Monitor, brand=ScreenPro}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD006, available=true, rating=4.6, name=Headphones, category=audio, price=150, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD007, price=89.99, available=true, rating=3.8, stock=25, supplier=CamSupply, name=Webcam, category=electronics, keywords=video, brand=CamTech}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD001, price=999.99, keywords=computer, brand=TechCorp, name=Laptop, available=true, rating=4.5, stock=50, supplier=TechSupply, category=electronics}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
2. [1] TestPerson{id=P002, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, department=hr, active=false, score=6, level=1, name=Charlie, age=16, salary=0, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing, name=Diana, age=45, salary=85000}
5. [1] TestPerson{id=P005, active=false, department=sales, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3}
6. [1] TestPerson{id=P006, salary=-5000, active=true, tags=test, status=active, score=0, department=qa, level=1, name=Frank, age=0}
7. [1] TestPerson{id=P007, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active, name=Grace, active=true}
8. [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1}
9. [1] TestPerson{id=P009, level=6, name=Ivy, active=true, status=active, age=40, salary=68000, score=8.7, tags=senior, department=engineering}
10. [1] TestPerson{id=P010, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5, name=X, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, department=hr, active=false, score=6, level=1, name=Charlie, age=16, salary=0, tags=intern}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales, name=Eve}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa, level=1}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, status=active, name=Grace, active=true, score=10, tags=executive, department=management, level=9}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed}
3. [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0}
5. [1] TestOrder{id=O005, priority=high, discount=100, region=south, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, amount=1, status=confirmed}
6. [1] TestOrder{id=O006, status=cancelled, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15}
7. [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}
9. [1] TestOrder{id=O009, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north}
10. [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, region=south, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, discount=100, region=south, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O001, discount=50, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, department=sales, level=2, age=25, active=true, name=Alice, salary=45000, score=8.5, tags=junior}
2. [1] TestPerson{id=P002, level=5, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}
4. [1] TestPerson{id=P004, status=active, level=7, active=true, tags=manager, department=marketing, name=Diana, age=45, salary=85000, score=7.8}
5. [1] TestPerson{id=P005, score=8, tags=employee, status=inactive, level=3, active=false, department=sales, name=Eve, age=30, salary=55000}
6. [1] TestPerson{id=P006, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa}
7. [1] TestPerson{id=P007, age=65, salary=95000, status=active, name=Grace, active=true, score=10, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support}
9. [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}
11. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50}
12. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed}
13. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}
14. [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east}
15. [1] TestOrder{id=O005, discount=100, region=south, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high}
16. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled}
17. [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent}
18. [1] TestOrder{id=O008, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05}
19. [1] TestOrder{id=O009, status=completed, region=north, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10}
20. [1] TestOrder{id=O010, region=east, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support}
   - Fait 2: [1] TestOrder{id=O002, region=south, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true, name=Alice}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, active=true, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, status=active, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, level=9, age=65, salary=95000, status=active, name=Grace, active=true, score=10, tags=executive, department=management}
   - Fait 2: [1] TestOrder{id=O009, status=completed, region=north, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, level=7, active=true, tags=manager, department=marketing, name=Diana, age=45, salary=85000, score=7.8, status=active}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, department=marketing, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true}
   - Fait 2: [1] TestOrder{id=O001, amount=2, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P003, level=1, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P003, level=1, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, status=active, score=0, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
   - Fait 2: [1] TestOrder{id=O008, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, status=active, department=sales, level=2, age=25, active=true, name=Alice, salary=45000, score=8.5}
   - Fait 2: [1] TestOrder{id=O002, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P009, level=6, name=Ivy, active=true, status=active, age=40, salary=68000, score=8.7, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}
   - Fait 2: [1] TestOrder{id=O005, region=south, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, status=active, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P008, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true, name=Alice, salary=45000}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P005, active=false, department=sales, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O008, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, status=active, score=0, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P002, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active, name=Bob, salary=75000}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P004, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing, name=Diana}
   - Fait 2: [1] TestOrder{id=O004, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, total=299.99}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing, name=Diana, age=45}
   - Fait 2: [1] TestOrder{id=O007, region=north, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P010, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5, name=X, salary=28000}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P001, active=true, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O001, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, date=2024-01-15}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing, name=Diana, age=45}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P004, status=active, level=7, active=true, tags=manager, department=marketing, name=Diana, age=45, salary=85000, score=7.8}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
   - Fait 2: [1] TestOrder{id=O009, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active, name=Grace, active=true}
   - Fait 2: [1] TestOrder{id=O001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50, customer_id=P001, product_id=PROD001}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P009, status=active, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, region=north, amount=2, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P005, active=false, department=sales, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active, age=40}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P008, level=1, tags=junior, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support, name=Henry, age=18, salary=25000}
   - Fait 2: [1] TestOrder{id=O006, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active, age=40}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales}
   - Fait 2: [1] TestOrder{id=O003, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001, product_id=PROD003, priority=high}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P003, level=1, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, department=hr, active=false, score=6, level=1, name=Charlie, age=16, salary=0, tags=intern}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
   - Fait 2: [1] TestOrder{id=O009, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O002, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P004, active=true, tags=manager, department=marketing, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7}
   - Fait 2: [1] TestOrder{id=O003, status=shipped, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support, name=Henry, age=18, salary=25000}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P003, level=1, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P010, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5, name=X, salary=28000}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P006, score=0, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active}
   - Fait 2: [1] TestOrder{id=O005, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active, name=Bob, salary=75000, active=true}
   - Fait 2: [1] TestOrder{id=O005, discount=100, region=south, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, level=5, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P006, score=0, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1, name=Charlie, age=16}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active, age=40}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active, age=40, salary=68000}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, active=false, score=6, level=1, name=Charlie, age=16, salary=0, tags=intern, status=inactive}
   - Fait 2: [1] TestOrder{id=O010, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, department=management, level=9, age=65, salary=95000, status=active, name=Grace, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true, name=Alice}
   - Fait 2: [1] TestOrder{id=O001, discount=50, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P005, score=8, tags=employee, status=inactive, level=3, active=false, department=sales, name=Eve, age=30, salary=55000}
   - Fait 2: [1] TestOrder{id=O004, region=east, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, status=inactive, level=1, tags=junior, department=support, name=Henry, age=18, salary=25000, active=false}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O003, region=north, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001, product_id=PROD003, priority=high, discount=15}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P006, active=true, tags=test, status=active, score=0, department=qa, level=1, name=Frank, age=0, salary=-5000}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P010, status=active, department=intern, level=1, age=22, score=6.5, name=X, salary=28000, active=true, tags=temp}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P006, status=active, score=0, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P006, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active, name=Bob, salary=75000, active=true}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O009, status=completed, region=north, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P005, active=false, department=sales, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O001, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, date=2024-01-15}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
2. [1] TestPerson{id=P002, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1, name=Charlie, age=16}
4. [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}
5. [1] TestPerson{id=P005, tags=employee, status=inactive, level=3, active=false, department=sales, name=Eve, age=30, salary=55000, score=8}
6. [1] TestPerson{id=P006, salary=-5000, active=true, tags=test, status=active, score=0, department=qa, level=1, name=Frank, age=0}
7. [1] TestPerson{id=P007, salary=95000, status=active, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support}
9. [1] TestPerson{id=P009, department=engineering, level=6, name=Ivy, active=true, status=active, age=40, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, age=25, active=true, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active, name=Bob, salary=75000, active=true}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa, level=1}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, status=active, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, age=22, score=6.5, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
2. [1] TestPerson{id=P002, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, active=false, score=6, level=1, name=Charlie, age=16, salary=0}
4. [1] TestPerson{id=P004, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing, name=Diana, age=45, salary=85000}
5. [1] TestPerson{id=P005, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales, name=Eve}
6. [1] TestPerson{id=P006, score=0, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active}
7. [1] TestPerson{id=P007, tags=executive, department=management, level=9, age=65, salary=95000, status=active, name=Grace, active=true, score=10}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support}
9. [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}
11. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50}
12. [1] TestOrder{id=O002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed, customer_id=P002}
13. [1] TestOrder{id=O003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001, product_id=PROD003}
14. [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10}
16. [1] TestOrder{id=O006, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west, total=999.98}
17. [1] TestOrder{id=O007, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north}
18. [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south}
19. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10}
20. [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P010, status=active, department=intern, level=1, age=22, score=6.5, name=X, salary=28000, active=true, tags=temp}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support, name=Henry}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, status=active, department=sales, level=2, age=25, active=true, name=Alice, salary=45000, score=8.5, tags=junior}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35}
   - Fait 2: [1] TestOrder{id=O001, amount=2, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, status=active, name=Grace, active=true, score=10, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O002, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south, amount=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, status=active, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior}
   - Fait 2: [1] TestOrder{id=O008, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
   - Fait 2: [1] TestOrder{id=O008, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P002, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50, customer_id=P001, product_id=PROD001}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, active=true, status=active, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1, name=Charlie, age=16}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1, total=89.99}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support, name=Henry, age=18, salary=25000}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, level=2, age=25, active=true, name=Alice, salary=45000, score=8.5, tags=junior, status=active}
   - Fait 2: [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10, customer_id=P001}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, level=5, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P005, active=false, department=sales, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, level=6, name=Ivy, active=true, status=active, age=40, salary=68000, score=8.7, tags=senior}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, region=north, amount=2, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}
   - Fait 2: [1] TestOrder{id=O002, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1, name=Charlie}
   - Fait 2: [1] TestOrder{id=O004, status=delivered, region=east, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa, level=1, name=Frank}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P007, status=active, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000}
   - Fait 2: [1] TestOrder{id=O004, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, total=299.99}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, department=marketing, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P002, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P007, status=active, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O002, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P006, score=0, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, status=active, department=intern, level=1, age=22, score=6.5, name=X, salary=28000, active=true}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P006, score=0, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active}
   - Fait 2: [1] TestOrder{id=O006, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P008, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, active=false, score=6, level=1, name=Charlie, age=16, salary=0, tags=intern, status=inactive}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}
   - Fait 2: [1] TestOrder{id=O004, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, total=299.99}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, status=inactive, level=3, active=false, department=sales, name=Eve, age=30, salary=55000, score=8}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, department=hr, active=false, score=6, level=1, name=Charlie, age=16, salary=0, tags=intern}
   - Fait 2: [1] TestOrder{id=O007, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, status=active, name=Grace, active=true, score=10, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P007, level=9, age=65, salary=95000, status=active, name=Grace, active=true, score=10, tags=executive, department=management}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing, name=Diana, age=45}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north, amount=2, total=1999.98, discount=50}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active, age=40, salary=68000, score=8.7}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P004, level=7, active=true, tags=manager, department=marketing, name=Diana, age=45, salary=85000, score=7.8, status=active}
   - Fait 2: [1] TestOrder{id=O008, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P010, age=22, score=6.5, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales, name=Eve}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa, level=1, name=Frank}
   - Fait 2: [1] TestOrder{id=O004, region=east, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, level=1, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}
   - Fait 2: [1] TestOrder{id=O006, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O001, amount=2, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O008, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P004, active=true, tags=manager, department=marketing, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7}
   - Fait 2: [1] TestOrder{id=O010, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, active=true, tags=test, status=active, score=0, department=qa, level=1, name=Frank, age=0}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}
   - Fait 2: [1] TestOrder{id=O010, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}
   - Fait 2: [1] TestOrder{id=O006, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P004, level=7, active=true, tags=manager, department=marketing, name=Diana, age=45, salary=85000, score=7.8, status=active}
   - Fait 2: [1] TestOrder{id=O003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001, product_id=PROD003}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P001, active=true, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25}
   - Fait 2: [1] TestOrder{id=O005, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P002, level=5, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales}
   - Fait 2: [1] TestOrder{id=O001, amount=2, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales}
   - Fait 2: [1] TestOrder{id=O008, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
   - Fait 2: [1] TestOrder{id=O006, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P002, level=5, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa, level=1, name=Frank}
   - Fait 2: [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
   - Fait 2: [1] TestOrder{id=O003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001, product_id=PROD003}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales}
   - Fait 2: [1] TestOrder{id=O004, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east, product_id=PROD004, amount=1, total=299.99, priority=normal}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}
   - Fait 2: [1] TestOrder{id=O005, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10, customer_id=P002, product_id=PROD001}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support, name=Henry, age=18}
   - Fait 2: [1] TestOrder{id=O006, region=west, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P003, status=inactive, department=hr, active=false, score=6, level=1, name=Charlie, age=16, salary=0, tags=intern}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west, total=999.98, date=2024-02-15}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
   - Fait 2: [1] TestOrder{id=O006, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west, total=999.98}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P006, score=0, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, level=6, name=Ivy, active=true, status=active, age=40, salary=68000, score=8.7, tags=senior}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, status=active, name=Grace, active=true, score=10, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10, customer_id=P002}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O006, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, amount=2, priority=low}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P005, active=false, department=sales, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
   - Fait 2: [1] TestOrder{id=O010, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active, name=Bob, salary=75000}
   - Fait 2: [1] TestOrder{id=O010, region=east, product_id=PROD001, status=refunded, priority=urgent, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales, name=Eve}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support, name=Henry}
   - Fait 2: [1] TestOrder{id=O004, status=delivered, region=east, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, status=active, department=sales, level=2, age=25, active=true, name=Alice, salary=45000, score=8.5}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false}
   - Fait 2: [1] TestOrder{id=O006, priority=low, discount=0, region=west, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, amount=2}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
2. [1] TestPerson{id=P002, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}
4. [1] TestPerson{id=P004, active=true, tags=manager, department=marketing, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7}
5. [1] TestPerson{id=P005, department=sales, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false}
6. [1] TestPerson{id=P006, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}
8. [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1}
9. [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}
10. [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5, name=X}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active, age=40}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5, name=X, salary=28000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3, active=false, department=sales}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true, name=Alice, salary=45000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, tags=intern, status=inactive, department=hr, active=false, score=6, level=1}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true, tags=manager, department=marketing}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, amount=2, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal, region=north}
2. [1] TestOrder{id=O002, total=25.5, priority=low, discount=0, region=south, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002}
3. [1] TestOrder{id=O003, discount=15, region=north, amount=3, total=225, date=2024-02-01, status=shipped, customer_id=P001, product_id=PROD003, priority=high}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10}
6. [1] TestOrder{id=O006, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0, region=west}
7. [1] TestOrder{id=O007, status=shipped, discount=50, region=north, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, region=north, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, region=north, amount=2, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, date=2024-01-15, status=pending, priority=normal}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, status=shipped, customer_id=P001, product_id=PROD003, priority=high, discount=15, region=north, amount=3, total=225, date=2024-02-01}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, total=299.99, priority=normal, discount=0, customer_id=P004, date=2024-02-05, status=delivered, region=east, product_id=PROD004}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, status=confirmed, priority=high, discount=100, region=south, total=999.99, date=2024-02-10}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, region=west, total=999.98, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, amount=2, priority=low, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, priority=low, discount=10, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, date=2024-03-15, discount=0, region=east, product_id=PROD001, status=refunded, priority=urgent}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, date=2024-01-20, status=confirmed, customer_id=P002, product_id=PROD002, total=25.5, priority=low, discount=0, region=south}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, priority=urgent, customer_id=P007, total=600, date=2024-03-01, status=shipped, discount=50, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}
2. [1] TestPerson{id=P002, age=35, status=active, name=Bob, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, department=hr, active=false, score=6, level=1, name=Charlie, age=16, salary=0, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, active=true, tags=manager, department=marketing, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7}
5. [1] TestPerson{id=P005, active=false, department=sales, name=Eve, age=30, salary=55000, score=8, tags=employee, status=inactive, level=3}
6. [1] TestPerson{id=P006, department=qa, level=1, name=Frank, age=0, salary=-5000, active=true, tags=test, status=active, score=0}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active}
8. [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, age=18, salary=25000, active=false, score=5.5, status=inactive, level=1}
9. [1] TestPerson{id=P009, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active, age=40, salary=68000, score=8.7}
10. [1] TestPerson{id=P010, age=22, score=6.5, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, score=10, tags=executive, department=management, level=9, age=65, salary=95000, status=active, name=Grace, active=true}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, tags=temp, status=active, department=intern, level=1, age=22, score=6.5}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, tags=junior, status=active, department=sales, level=2, age=25, active=true}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, department=marketing, name=Diana, age=45, salary=85000, score=7.8, status=active, level=7, active=true}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, level=3, active=false, department=sales, name=Eve, age=30, salary=55000, score=8, tags=employee}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, status=inactive, level=1, tags=junior, department=support, name=Henry, age=18, salary=25000}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, department=engineering, level=6, name=Ivy, active=true, status=active}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, department=engineering, level=5, age=35, status=active, name=Bob, salary=75000}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, active=false, score=6, level=1, name=Charlie, age=16, salary=0}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, active=true, tags=test, status=active, score=0, department=qa, level=1, name=Frank}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
