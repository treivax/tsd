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

1. [1] TestPerson{id=P001, status=active, level=2, age=25, score=8.5, tags=junior, department=sales, name=Alice, salary=45000, active=true}
2. [1] TestPerson{id=P002, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2}
3. [1] TestPerson{id=P003, score=6, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false}
4. [1] TestPerson{id=P004, salary=85000, active=true, department=marketing, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7}
5. [1] TestPerson{id=P005, level=3, salary=55000, active=false, score=8, name=Eve, age=30, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, level=1, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, level=1, age=18, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7, tags=senior, status=active}
10. [1] TestPerson{id=P010, active=true, status=active, level=1, name=X, score=6.5, tags=temp, department=intern, age=22, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, level=1, age=16, salary=0, active=false, score=6, name=Charlie, tags=intern, status=inactive, department=hr}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales, name=Alice, salary=45000}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior, department=support}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, department=intern, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, status=pending, priority=normal, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south}
3. [1] TestOrder{id=O003, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003, amount=3, date=2024-02-01}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10}
6. [1] TestOrder{id=O006, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, total=999.98, priority=low}
7. [1] TestOrder{id=O007, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4}
8. [1] TestOrder{id=O008, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending, product_id=PROD002}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}
10. [1] TestOrder{id=O010, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, region=south, product_id=PROD002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O005, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10, customer_id=P002, amount=1}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, tags=junior, department=sales, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5}
2. [1] TestPerson{id=P002, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior}
3. [1] TestPerson{id=P003, age=16, salary=0, active=false, score=6, name=Charlie, tags=intern, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
5. [1] TestPerson{id=P005, salary=55000, active=false, score=8, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3}
6. [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, department=qa, level=1, name=Frank, active=true, score=0, status=active}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior, department=support}
9. [1] TestPerson{id=P009, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5, tags=temp, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior, name=Bob}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, department=marketing, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, score=10, status=active, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, active=true, level=6, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal}
2. [1] TestOrder{id=O002, region=south, product_id=PROD002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0}
3. [1] TestOrder{id=O003, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10}
6. [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, total=999.98, priority=low, region=west}
7. [1] TestOrder{id=O007, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}
9. [1] TestOrder{id=O009, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north, customer_id=P001}
10. [1] TestOrder{id=O010, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O007, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003, amount=3}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, status=completed, priority=low, product_id=PROD007, discount=10, region=north, customer_id=P001, amount=1, total=89.99, date=2024-03-10}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal, product_id=PROD001}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east, product_id=PROD004, amount=1}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, date=2024-02-10, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O006, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, brand=TechCorp, stock=50, supplier=TechSupply, price=999.99, available=true, name=Laptop, category=electronics, rating=4.5, keywords=computer}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, price=25.5, available=true, rating=4, stock=200, keywords=peripheral, brand=TechCorp, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, category=accessories, price=75, rating=3.5, keywords=typing, brand=KeyTech, supplier=KeySupply, name=Keyboard, available=false, stock=0}
4. [1] TestProduct{id=PROD004, name=Monitor, category=electronics, price=299.99, brand=ScreenPro, available=true, rating=4.8, keywords=display, stock=30, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, rating=2, supplier=OldSupply, price=8.5, available=false, keywords=obsolete, brand=OldTech, stock=0}
6. [1] TestProduct{id=PROD006, name=Headphones, category=audio, rating=4.6, price=150, available=true, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, keywords=video, brand=CamTech, supplier=CamSupply, name=Webcam, category=electronics, available=true, rating=3.8, stock=25, price=89.99}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, category=electronics, rating=4.5, keywords=computer, brand=TechCorp, stock=50, supplier=TechSupply, price=999.99, available=true}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, category=accessories, price=25.5, available=true, rating=4, stock=200, keywords=peripheral, brand=TechCorp, supplier=TechSupply, name=Mouse}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, rating=3.5, keywords=typing, brand=KeyTech, supplier=KeySupply, name=Keyboard, available=false, stock=0, category=accessories, price=75}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, brand=ScreenPro, available=true, rating=4.8, keywords=display, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics, price=299.99}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply, name=Headphones, category=audio, rating=4.6, price=150, available=true}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, name=Webcam, category=electronics, available=true, rating=3.8, stock=25, price=89.99, keywords=video, brand=CamTech, supplier=CamSupply}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, score=8.5, tags=junior, department=sales, name=Alice, salary=45000, active=true, status=active, level=2}
2. [1] TestPerson{id=P002, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
5. [1] TestPerson{id=P005, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8, name=Eve, age=30}
6. [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, department=qa, level=1, name=Frank, active=true, score=0, status=active}
7. [1] TestPerson{id=P007, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active, age=65, salary=95000}
8. [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior, department=support}
9. [1] TestPerson{id=P009, level=6, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5, tags=temp, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal}
2. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed, customer_id=P002}
3. [1] TestOrder{id=O003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003}
4. [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, priority=normal, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0}
5. [1] TestOrder{id=O005, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10}
6. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0}
7. [1] TestOrder{id=O007, priority=urgent, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}
10. [1] TestOrder{id=O010, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, priority=normal, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, total=999.99, date=2024-02-10, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15, discount=0, region=east}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales}
2. [1] TestPerson{id=P002, active=true, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000}
3. [1] TestPerson{id=P003, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6, name=Charlie, tags=intern}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
5. [1] TestPerson{id=P005, active=false, score=8, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, status=inactive, level=1, age=18, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5}
9. [1] TestPerson{id=P009, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, department=intern, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6, name=Charlie}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal, product_id=PROD001}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south}
3. [1] TestOrder{id=O003, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003, amount=3}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east}
5. [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001}
6. [1] TestOrder{id=O006, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, total=999.98}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}
9. [1] TestOrder{id=O009, priority=low, product_id=PROD007, discount=10, region=north, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O009, discount=10, region=north, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O010, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O003, discount=15, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003, amount=3, date=2024-02-01, status=shipped}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, priority=normal, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O006, status=cancelled, discount=0, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, age=25, score=8.5, tags=junior, department=sales, name=Alice, salary=45000, active=true, status=active}
2. [1] TestPerson{id=P002, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2}
3. [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}
6. [1] TestPerson{id=P006, status=active, age=0, salary=-5000, tags=test, department=qa, level=1, name=Frank, active=true, score=0}
7. [1] TestPerson{id=P007, name=Grace, active=true, score=10, status=active, age=65, salary=95000, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1}
9. [1] TestPerson{id=P009, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, level=1, name=X, score=6.5, tags=temp, department=intern, age=22, salary=28000, active=true, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, score=8, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales, name=Alice, salary=45000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, active=true, department=marketing, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, price=999.99, available=true, name=Laptop, category=electronics, rating=4.5, keywords=computer, brand=TechCorp, stock=50, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, supplier=TechSupply, name=Mouse, category=accessories, price=25.5, available=true, rating=4, stock=200, keywords=peripheral, brand=TechCorp}
3. [1] TestProduct{id=PROD003, category=accessories, price=75, rating=3.5, keywords=typing, brand=KeyTech, supplier=KeySupply, name=Keyboard, available=false, stock=0}
4. [1] TestProduct{id=PROD004, price=299.99, brand=ScreenPro, available=true, rating=4.8, keywords=display, stock=30, supplier=ScreenSupply, name=Monitor, category=electronics}
5. [1] TestProduct{id=PROD005, stock=0, name=OldKeyboard, category=accessories, rating=2, supplier=OldSupply, price=8.5, available=false, keywords=obsolete, brand=OldTech}
6. [1] TestProduct{id=PROD006, category=audio, rating=4.6, price=150, available=true, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply, name=Headphones}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, available=true, rating=3.8, stock=25, price=89.99, keywords=video, brand=CamTech, supplier=CamSupply}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD007, brand=CamTech, supplier=CamSupply, name=Webcam, category=electronics, available=true, rating=3.8, stock=25, price=89.99, keywords=video}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, category=electronics, rating=4.5, keywords=computer, brand=TechCorp, stock=50, supplier=TechSupply, price=999.99, available=true}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD002, supplier=TechSupply, name=Mouse, category=accessories, price=25.5, available=true, rating=4, stock=200, keywords=peripheral, brand=TechCorp}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD003, supplier=KeySupply, name=Keyboard, available=false, stock=0, category=accessories, price=75, rating=3.5, keywords=typing, brand=KeyTech}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, category=electronics, price=299.99, brand=ScreenPro, available=true, rating=4.8, keywords=display, stock=30, supplier=ScreenSupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD006, supplier=AudioSupply, name=Headphones, category=audio, rating=4.6, price=150, available=true, keywords=sound, brand=AudioMax, stock=75}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales, name=Alice}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6, name=Charlie}
4. [1] TestPerson{id=P004, department=marketing, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true}
5. [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}
6. [1] TestPerson{id=P006, tags=test, department=qa, level=1, name=Frank, active=true, score=0, status=active, age=0, salary=-5000}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior}
9. [1] TestPerson{id=P009, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5, tags=temp, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, status=active, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7, tags=senior}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal, product_id=PROD001, amount=2, total=1999.98}
2. [1] TestOrder{id=O002, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low}
3. [1] TestOrder{id=O003, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15}
4. [1] TestOrder{id=O004, date=2024-02-05, priority=normal, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004}
5. [1] TestOrder{id=O005, product_id=PROD001, total=999.99, date=2024-02-10, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south}
6. [1] TestOrder{id=O006, discount=0, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled}
7. [1] TestOrder{id=O007, priority=urgent, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped}
8. [1] TestOrder{id=O008, discount=0, region=south, customer_id=P010, status=pending, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}
10. [1] TestOrder{id=O010, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, region=south, product_id=PROD001, total=999.99, date=2024-02-10, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50, customer_id=P007}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, priority=normal, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, age=25, score=8.5, tags=junior, department=sales, name=Alice, salary=45000, active=true, status=active}
2. [1] TestPerson{id=P002, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, salary=85000, active=true, department=marketing, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7}
5. [1] TestPerson{id=P005, level=3, salary=55000, active=false, score=8, name=Eve, age=30, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, tags=test, department=qa, level=1, name=Frank, active=true, score=0, status=active, age=0, salary=-5000}
7. [1] TestPerson{id=P007, status=active, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10}
8. [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior, department=support}
9. [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy}
10. [1] TestPerson{id=P010, tags=temp, department=intern, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5}
11. [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal}
12. [1] TestOrder{id=O002, region=south, product_id=PROD002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0}
13. [1] TestOrder{id=O003, priority=high, region=north, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225}
14. [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east}
15. [1] TestOrder{id=O005, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10, customer_id=P002, amount=1, status=confirmed}
16. [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0}
17. [1] TestOrder{id=O007, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent}
18. [1] TestOrder{id=O008, customer_id=P010, status=pending, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south}
19. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}
20. [1] TestOrder{id=O010, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, status=active, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7, tags=senior}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending, product_id=PROD002, amount=10}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active, age=65, salary=95000}
   - Fait 2: [1] TestOrder{id=O008, discount=0, region=south, customer_id=P010, status=pending, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, level=1, age=18, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales, name=Alice, salary=45000}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P006, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1, name=Frank}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P005, level=3, salary=55000, active=false, score=8, name=Eve, age=30, tags=employee, status=inactive, department=sales}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P010, name=X, score=6.5, tags=temp, department=intern, age=22, salary=28000, active=true, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active, age=65}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, discount=10, region=north, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P010, name=X, score=6.5, tags=temp, department=intern, age=22, salary=28000, active=true, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, total=999.98, priority=low, region=west}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P010, status=active, level=1, name=X, score=6.5, tags=temp, department=intern, age=22, salary=28000, active=true}
   - Fait 2: [1] TestOrder{id=O007, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P002, active=true, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, status=pending, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15, discount=0, region=east}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O004, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O009, priority=low, product_id=PROD007, discount=10, region=north, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior, department=support}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P009, status=active, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7, tags=senior}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, active=false, score=6, name=Charlie, tags=intern, status=inactive, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P010, active=true, status=active, level=1, name=X, score=6.5, tags=temp, department=intern, age=22, salary=28000}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing, name=Diana, age=45, score=7.8}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P007, level=9, name=Grace, active=true, score=10, status=active, age=65, salary=95000, tags=executive, department=management}
   - Fait 2: [1] TestOrder{id=O003, priority=high, region=north, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, level=1, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior, department=support}
   - Fait 2: [1] TestOrder{id=O008, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending, product_id=PROD002, amount=10}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, discount=10, region=north, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P010, level=1, name=X, score=6.5, tags=temp, department=intern, age=22, salary=28000, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O006, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, total=999.98}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O007, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P005, score=8, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15, discount=0, region=east}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O010, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15, discount=0}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O001, region=north, customer_id=P001, status=pending, priority=normal, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O001, region=north, customer_id=P001, status=pending, priority=normal, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O006, discount=0, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P007, score=10, status=active, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true}
   - Fait 2: [1] TestOrder{id=O006, discount=0, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O007, status=shipped, priority=urgent, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, active=false, score=6, name=Charlie, tags=intern, status=inactive, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O007, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, status=pending, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, department=qa, level=1, name=Frank, active=true, score=0, status=active}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P007, score=10, status=active, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed, customer_id=P002}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P010, name=X, score=6.5, tags=temp, department=intern, age=22, salary=28000, active=true, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O003, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003, amount=3, date=2024-02-01}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, level=3, salary=55000, active=false, score=8, name=Eve, age=30, tags=employee, status=inactive}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, department=qa, level=1, name=Frank, active=true, score=0, status=active, age=0, salary=-5000}
   - Fait 2: [1] TestOrder{id=O005, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10, customer_id=P002, amount=1, status=confirmed, priority=high}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P010, active=true, status=active, level=1, name=X, score=6.5, tags=temp, department=intern, age=22, salary=28000}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, discount=10, region=north, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, priority=normal, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O008, discount=0, region=south, customer_id=P010, status=pending, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6}
   - Fait 2: [1] TestOrder{id=O010, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior, department=support}
   - Fait 2: [1] TestOrder{id=O001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal, product_id=PROD001}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P001, age=25, score=8.5, tags=junior, department=sales, name=Alice, salary=45000, active=true, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P004, status=active, level=7, salary=85000, active=true, department=marketing, name=Diana, age=45, score=7.8, tags=manager}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P007, status=active, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P007, level=9, name=Grace, active=true, score=10, status=active, age=65, salary=95000, tags=executive, department=management}
   - Fait 2: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P010, tags=temp, department=intern, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5}
   - Fait 2: [1] TestOrder{id=O010, total=75000, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O002, region=south, product_id=PROD002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P004, active=true, department=marketing, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active, age=65}
   - Fait 2: [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east, product_id=PROD004}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5, tags=temp}
   - Fait 2: [1] TestOrder{id=O005, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10, customer_id=P002, amount=1, status=confirmed}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, department=qa, level=1, name=Frank, active=true, score=0, status=active}
   - Fait 2: [1] TestOrder{id=O009, status=completed, priority=low, product_id=PROD007, discount=10, region=north, customer_id=P001, amount=1, total=89.99, date=2024-03-10}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, department=sales, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, status=pending, priority=normal, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior, department=support}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, department=sales, name=Alice, salary=45000, active=true, status=active, level=2, age=25}
   - Fait 2: [1] TestOrder{id=O007, priority=urgent, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2}
   - Fait 2: [1] TestOrder{id=O007, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5, tags=temp}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, status=pending, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, department=sales, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, priority=normal, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, priority=normal, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, total=999.98, priority=low, region=west}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy}
   - Fait 2: [1] TestOrder{id=O009, status=completed, priority=low, product_id=PROD007, discount=10, region=north, customer_id=P001, amount=1, total=89.99, date=2024-03-10}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true}
   - Fait 2: [1] TestOrder{id=O005, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10, customer_id=P002, amount=1, status=confirmed}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, tags=junior, department=sales, name=Alice, salary=45000, active=true, status=active, level=2, age=25}
2. [1] TestPerson{id=P002, salary=75000, active=true, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35}
3. [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}
6. [1] TestPerson{id=P006, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1, name=Frank}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1}
9. [1] TestPerson{id=P009, level=6, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true}
10. [1] TestPerson{id=P010, tags=temp, department=intern, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, age=25, score=8.5, tags=junior, department=sales, name=Alice, salary=45000, active=true, status=active, level=2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing, name=Diana, age=45, score=7.8}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1, name=Frank, active=true}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, tags=temp, department=intern, age=22, salary=28000, active=true, status=active, level=1, name=X}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, score=8.5, tags=junior, department=sales, name=Alice, salary=45000, active=true, status=active, level=2}
2. [1] TestPerson{id=P002, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
5. [1] TestPerson{id=P005, score=8, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior, department=support}
9. [1] TestPerson{id=P009, level=6, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5, tags=temp, department=intern}
11. [1] TestOrder{id=O001, discount=50, region=north, customer_id=P001, status=pending, priority=normal, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15}
12. [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed}
13. [1] TestOrder{id=O003, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15}
14. [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, priority=normal, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0}
15. [1] TestOrder{id=O005, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10}
16. [1] TestOrder{id=O006, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, total=999.98}
17. [1] TestOrder{id=O007, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent}
18. [1] TestOrder{id=O008, status=pending, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010}
19. [1] TestOrder{id=O009, product_id=PROD007, discount=10, region=north, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}
20. [1] TestOrder{id=O010, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, product_id=PROD001}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior, name=Bob}
   - Fait 2: [1] TestOrder{id=O002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, department=qa, level=1, name=Frank, active=true, score=0, status=active}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5, tags=temp}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P008, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15, discount=0, region=east}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior, name=Bob}
   - Fait 2: [1] TestOrder{id=O003, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003, amount=3, date=2024-02-01}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, level=1, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O008, priority=normal, discount=0, region=south, customer_id=P010, status=pending, product_id=PROD002, amount=10, total=255, date=2024-03-05}

11. **Token 11**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}

12. **Token 12**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}

13. **Token 13**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, level=3, salary=55000, active=false, score=8, name=Eve, age=30, tags=employee, status=inactive}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}

14. **Token 14**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north, customer_id=P001, amount=1}

15. **Token 15**:
   - Fait 1: [1] TestPerson{id=P006, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1, name=Frank}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north}

16. **Token 16**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003}

17. **Token 17**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001}

18. **Token 18**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, department=sales, name=Alice, salary=45000, active=true, status=active, level=2, age=25}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}

19. **Token 19**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}

20. **Token 20**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, discount=10, region=north, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low}

21. **Token 21**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}

22. **Token 22**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}

23. **Token 23**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O002, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed, customer_id=P002, total=25.5}

24. **Token 24**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, department=qa, level=1, name=Frank, active=true, score=0, status=active}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed}

25. **Token 25**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east, product_id=PROD004}

26. **Token 26**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000}

27. **Token 27**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6}
   - Fait 2: [1] TestOrder{id=O001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal, product_id=PROD001}

28. **Token 28**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, status=pending, priority=normal, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north}

29. **Token 29**:
   - Fait 1: [1] TestPerson{id=P004, level=7, salary=85000, active=true, department=marketing, name=Diana, age=45, score=7.8, tags=manager, status=active}
   - Fait 2: [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal, product_id=PROD001, amount=2}

30. **Token 30**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal}

31. **Token 31**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, tags=temp, department=intern, age=22, salary=28000, active=true, status=active, level=1, name=X}
   - Fait 2: [1] TestOrder{id=O003, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003, amount=3}

32. **Token 32**:
   - Fait 1: [1] TestPerson{id=P007, status=active, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east}

33. **Token 33**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales, name=Alice}
   - Fait 2: [1] TestOrder{id=O007, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north}

34. **Token 34**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent}

35. **Token 35**:
   - Fait 1: [1] TestPerson{id=P001, status=active, level=2, age=25, score=8.5, tags=junior, department=sales, name=Alice, salary=45000, active=true}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, total=999.98, priority=low, region=west}

36. **Token 36**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O002, region=south, product_id=PROD002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0}

37. **Token 37**:
   - Fait 1: [1] TestPerson{id=P007, score=10, status=active, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true}
   - Fait 2: [1] TestOrder{id=O005, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10, customer_id=P002}

38. **Token 38**:
   - Fait 1: [1] TestPerson{id=P007, status=active, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10}
   - Fait 2: [1] TestOrder{id=O009, region=north, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10}

39. **Token 39**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O010, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000}

40. **Token 40**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed}

41. **Token 41**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed}

42. **Token 42**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002}

43. **Token 43**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed, customer_id=P002}

44. **Token 44**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, department=sales, name=Alice, salary=45000, active=true, status=active, level=2, age=25}
   - Fait 2: [1] TestOrder{id=O004, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered}

45. **Token 45**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50}

46. **Token 46**:
   - Fait 1: [1] TestPerson{id=P010, level=1, name=X, score=6.5, tags=temp, department=intern, age=22, salary=28000, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending}

47. **Token 47**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15}

48. **Token 48**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, active=false, score=6, name=Charlie, tags=intern, status=inactive, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O003, priority=high, region=north, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225}

49. **Token 49**:
   - Fait 1: [1] TestPerson{id=P004, level=7, salary=85000, active=true, department=marketing, name=Diana, age=45, score=7.8, tags=manager, status=active}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north}

50. **Token 50**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O003, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north}

51. **Token 51**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, department=sales, name=Alice, salary=45000, active=true, status=active, level=2, age=25}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99}

52. **Token 52**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, discount=0, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15}

53. **Token 53**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, status=pending, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south}

54. **Token 54**:
   - Fait 1: [1] TestPerson{id=P002, active=true, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0}

55. **Token 55**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, active=true, score=10, status=active, age=65, salary=95000, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O008, region=south, customer_id=P010, status=pending, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0}

56. **Token 56**:
   - Fait 1: [1] TestPerson{id=P001, age=25, score=8.5, tags=junior, department=sales, name=Alice, salary=45000, active=true, status=active, level=2}
   - Fait 2: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north, customer_id=P001, amount=1}

57. **Token 57**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, level=1, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test}
   - Fait 2: [1] TestOrder{id=O001, status=pending, priority=normal, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001}

58. **Token 58**:
   - Fait 1: [1] TestPerson{id=P008, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, status=pending, priority=normal, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north}

59. **Token 59**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O007, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent}

60. **Token 60**:
   - Fait 1: [1] TestPerson{id=P001, status=active, level=2, age=25, score=8.5, tags=junior, department=sales, name=Alice, salary=45000, active=true}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15, discount=0, region=east}

61. **Token 61**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal, product_id=PROD001}

62. **Token 62**:
   - Fait 1: [1] TestPerson{id=P006, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1, name=Frank}
   - Fait 2: [1] TestOrder{id=O007, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600}

63. **Token 63**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50}

64. **Token 64**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}

65. **Token 65**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7}
   - Fait 2: [1] TestOrder{id=O010, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15}

66. **Token 66**:
   - Fait 1: [1] TestPerson{id=P003, score=6, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15, discount=0, region=east}

67. **Token 67**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent}

68. **Token 68**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, date=2024-03-15, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent}

69. **Token 69**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, department=qa, level=1, name=Frank, active=true, score=0, status=active}
   - Fait 2: [1] TestOrder{id=O006, discount=0, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled}

70. **Token 70**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0}

71. **Token 71**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, active=false, score=6, name=Charlie, tags=intern, status=inactive, department=hr, level=1}
   - Fait 2: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001}

72. **Token 72**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10}

73. **Token 73**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10}

74. **Token 74**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003}

75. **Token 75**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south, product_id=PROD001, total=999.99, date=2024-02-10}

76. **Token 76**:
   - Fait 1: [1] TestPerson{id=P010, name=X, score=6.5, tags=temp, department=intern, age=22, salary=28000, active=true, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006}

77. **Token 77**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}

78. **Token 78**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior, department=support}
   - Fait 2: [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}

79. **Token 79**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, status=active, level=1, name=X, score=6.5, tags=temp, department=intern, age=22}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, discount=0, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15}

80. **Token 80**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}
   - Fait 2: [1] TestOrder{id=O001, discount=50, region=north, customer_id=P001, status=pending, priority=normal, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15}

81. **Token 81**:
   - Fait 1: [1] TestPerson{id=P001, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales, name=Alice, salary=45000}
   - Fait 2: [1] TestOrder{id=O002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002}

82. **Token 82**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, department=sales, level=3, salary=55000, active=false, score=8, name=Eve, age=30, tags=employee}
   - Fait 2: [1] TestOrder{id=O004, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east}

83. **Token 83**:
   - Fait 1: [1] TestPerson{id=P009, level=6, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, priority=normal, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004}

84. **Token 84**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, priority=normal, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004}

85. **Token 85**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O004, customer_id=P004, date=2024-02-05, priority=normal, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0}

86. **Token 86**:
   - Fait 1: [1] TestPerson{id=P009, status=active, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7, tags=senior}
   - Fait 2: [1] TestOrder{id=O009, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north, customer_id=P001, amount=1, total=89.99}

87. **Token 87**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, status=active, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace}
   - Fait 2: [1] TestOrder{id=O003, region=north, product_id=PROD003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high}

88. **Token 88**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5, tags=temp, department=intern}
   - Fait 2: [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east, product_id=PROD004}

89. **Token 89**:
   - Fait 1: [1] TestPerson{id=P009, level=6, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true}
   - Fait 2: [1] TestOrder{id=O006, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, total=999.98}

90. **Token 90**:
   - Fait 1: [1] TestPerson{id=P003, active=false, score=6, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, total=999.98, priority=low, region=west}

91. **Token 91**:
   - Fait 1: [1] TestPerson{id=P004, level=7, salary=85000, active=true, department=marketing, name=Diana, age=45, score=7.8, tags=manager, status=active}
   - Fait 2: [1] TestOrder{id=O006, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0}

92. **Token 92**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005, amount=2, total=999.98, priority=low, region=west}

93. **Token 93**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{id=O001, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal, product_id=PROD001, amount=2, total=1999.98}

94. **Token 94**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active, age=65, salary=95000}
   - Fait 2: [1] TestOrder{id=O002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south, product_id=PROD002, amount=1, status=confirmed, customer_id=P002}

95. **Token 95**:
   - Fait 1: [1] TestPerson{id=P010, name=X, score=6.5, tags=temp, department=intern, age=22, salary=28000, active=true, status=active, level=1}
   - Fait 2: [1] TestOrder{id=O005, product_id=PROD001, total=999.99, date=2024-02-10, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south}

96. **Token 96**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7, tags=senior, status=active, department=engineering}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50}

97. **Token 97**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50}

98. **Token 98**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50, customer_id=P007}

99. **Token 99**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50}

100. **Token 100**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
   - Fait 2: [1] TestOrder{id=O009, region=north, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, age=16, salary=0, active=false, score=6, name=Charlie, tags=intern, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}
6. [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, department=qa, level=1, name=Frank, active=true, score=0, status=active}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior, department=support}
9. [1] TestPerson{id=P009, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5, tags=temp, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior, department=support, name=Henry}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing, name=Diana, age=45}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, score=8, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P010, department=intern, age=22, salary=28000, active=true, status=active, level=1, name=X, score=6.5, tags=temp}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior, name=Bob}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, priority=normal, product_id=PROD001, amount=2, total=1999.98, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south}
3. [1] TestOrder{id=O003, amount=3, date=2024-02-01, status=shipped, discount=15, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003}
4. [1] TestOrder{id=O004, discount=0, customer_id=P004, date=2024-02-05, priority=normal, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered}
5. [1] TestOrder{id=O005, product_id=PROD001, total=999.99, date=2024-02-10, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south}
6. [1] TestOrder{id=O006, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15, status=cancelled, discount=0, product_id=PROD005}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped, priority=urgent, region=north, total=600, discount=50}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, region=south, customer_id=P010, status=pending}
9. [1] TestOrder{id=O009, priority=low, product_id=PROD007, discount=10, region=north, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed}
10. [1] TestOrder{id=O010, discount=0, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, discount=15, customer_id=P001, total=225, priority=high, region=north, product_id=PROD003, amount=3, date=2024-02-01, status=shipped}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, region=east, product_id=PROD004, amount=1, total=299.99, status=delivered, discount=0, customer_id=P004, date=2024-02-05, priority=normal}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, priority=urgent, region=north, total=600, discount=50, customer_id=P007, product_id=PROD006, amount=4, date=2024-03-01, status=shipped}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O008, discount=0, region=south, customer_id=P010, status=pending, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, status=completed, priority=low, product_id=PROD007, discount=10, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O010, region=east, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, product_id=PROD001, date=2024-03-15, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, status=confirmed, customer_id=P002, total=25.5, date=2024-01-20, priority=low, discount=0, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, total=999.99, date=2024-02-10, customer_id=P002, amount=1, status=confirmed, priority=high, discount=100, region=south}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O006, status=cancelled, discount=0, product_id=PROD005, amount=2, total=999.98, priority=low, region=west, customer_id=P005, date=2024-02-15}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O001, date=2024-01-15, discount=50, region=north, customer_id=P001, status=pending, priority=normal, product_id=PROD001, amount=2, total=1999.98}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales}
2. [1] TestPerson{id=P002, name=Bob, score=9.2, status=active, department=engineering, level=5, age=35, salary=75000, active=true, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, tags=intern, status=inactive, department=hr, level=1, age=16, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing, name=Diana, age=45, score=7.8}
5. [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, status=active, age=0, salary=-5000, tags=test, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, department=support, name=Henry, active=false, score=5.5, status=inactive, level=1}
9. [1] TestPerson{id=P009, name=Ivy, score=8.7, tags=senior, status=active, department=engineering, age=40, salary=68000, active=true, level=6}
10. [1] TestPerson{id=P010, active=true, status=active, level=1, name=X, score=6.5, tags=temp, department=intern, age=22, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, active=false, score=6, name=Charlie, tags=intern, status=inactive, department=hr, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, status=active, level=7, salary=85000, active=true, department=marketing}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, tags=employee, status=inactive, department=sales, level=3, salary=55000, active=false, score=8}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, status=active, level=2, age=25, score=8.5, tags=junior, department=sales}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, name=Bob, score=9.2, status=active, department=engineering, level=5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, department=qa, level=1, name=Frank, active=true, score=0, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, tags=executive, department=management, level=9, name=Grace, active=true, score=10, status=active}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, status=inactive, level=1, age=18, salary=25000, tags=junior, department=support}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, age=40, salary=68000, active=true, level=6, name=Ivy, score=8.7, tags=senior, status=active}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P010, level=1, name=X, score=6.5, tags=temp, department=intern, age=22, salary=28000, active=true, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
