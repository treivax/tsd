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

1. [1] TestPerson{id=P001, level=2, name=Alice, age=25, score=8.5, tags=junior, status=active, department=sales, salary=45000, active=true}
2. [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, department=engineering, level=5, name=Bob, active=true, status=active}
3. [1] TestPerson{id=P003, age=16, tags=intern, status=inactive, department=hr, level=1, name=Charlie, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, score=7.8, salary=85000, tags=manager, status=active, department=marketing, level=7, name=Diana, age=45, active=true}
5. [1] TestPerson{id=P005, score=8, salary=55000, tags=employee, status=inactive, department=sales, level=3, name=Eve, age=30, active=false}
6. [1] TestPerson{id=P006, score=0, tags=test, name=Frank, age=0, salary=-5000, active=true, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, score=10, status=active, department=management, age=65, tags=executive, level=9, name=Grace, salary=95000, active=true}
8. [1] TestPerson{id=P008, status=inactive, age=18, active=false, score=5.5, tags=junior, department=support, level=1, name=Henry, salary=25000}
9. [1] TestPerson{id=P009, age=40, active=true, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering, name=Ivy}
10. [1] TestPerson{id=P010, department=intern, salary=28000, active=true, tags=temp, status=active, level=1, name=X, age=22, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, salary=45000, active=true, level=2, name=Alice, age=25, score=8.5, tags=junior, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, active=false, score=8, salary=55000, tags=employee, status=inactive, department=sales, level=3, name=Eve, age=30}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P007, active=true, score=10, status=active, department=management, age=65, tags=executive, level=9, name=Grace, salary=95000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, age=18, active=false, score=5.5, tags=junior, department=support, level=1, name=Henry, salary=25000, status=inactive}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, active=true, status=active, age=35, salary=75000, score=9.2, tags=senior, department=engineering, level=5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, age=16, tags=intern, status=inactive, department=hr, level=1, name=Charlie}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, salary=85000, tags=manager, status=active, department=marketing, level=7}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, name=Ivy, age=40, active=true, tags=senior, status=active, level=6, salary=68000, score=8.7}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, status=active, level=1, name=X, age=22, score=6.5, department=intern}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, amount=2, status=pending, region=north, customer_id=P001, product_id=PROD001, total=1999.98, date=2024-01-15, priority=normal, discount=50}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, status=confirmed, priority=low, discount=0, region=south, customer_id=P002, total=25.5, date=2024-01-20}
3. [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, amount=3, discount=15}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, date=2024-02-05, status=delivered, customer_id=P004, total=299.99, priority=normal, discount=0, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, priority=high, discount=100, amount=1, total=999.99, region=south}
6. [1] TestOrder{id=O006, customer_id=P005, amount=2, date=2024-02-15, priority=low, discount=0, product_id=PROD005, total=999.98, status=cancelled, region=west}
7. [1] TestOrder{id=O007, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4}
8. [1] TestOrder{id=O008, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south, customer_id=P010}
9. [1] TestOrder{id=O009, product_id=PROD007, amount=1, date=2024-03-10, status=completed, priority=low, region=north, customer_id=P001, total=89.99, discount=10}
10. [1] TestOrder{id=O010, date=2024-03-15, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O010, priority=urgent, discount=0, region=east, customer_id=P006, amount=1, date=2024-03-15, product_id=PROD001, total=75000, status=refunded}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O001, amount=2, status=pending, region=north, customer_id=P001, product_id=PROD001, total=1999.98, date=2024-01-15, priority=normal, discount=50}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, total=299.99, priority=normal, discount=0, region=east, product_id=PROD004, amount=1, date=2024-02-05, status=delivered}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, region=south, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, priority=high, discount=100, amount=1, total=999.99}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, product_id=PROD002, amount=1, status=confirmed, priority=low, discount=0, region=south}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, amount=3, discount=15, product_id=PROD003, total=225, date=2024-02-01, status=shipped, priority=high, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, status=completed, priority=low, region=north, customer_id=P001, total=89.99, discount=10, product_id=PROD007, amount=1, date=2024-03-10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, active=true, level=2, name=Alice, age=25, score=8.5, tags=junior, status=active, department=sales}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, department=engineering, level=5, name=Bob, active=true, status=active, age=35, salary=75000}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, age=16, tags=intern, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, salary=85000, tags=manager, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, level=3, name=Eve, age=30, active=false, score=8, salary=55000, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, status=active, department=qa, level=1, score=0, tags=test, name=Frank, age=0, salary=-5000, active=true}
7. [1] TestPerson{id=P007, active=true, score=10, status=active, department=management, age=65, tags=executive, level=9, name=Grace, salary=95000}
8. [1] TestPerson{id=P008, age=18, active=false, score=5.5, tags=junior, department=support, level=1, name=Henry, salary=25000, status=inactive}
9. [1] TestPerson{id=P009, age=40, active=true, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering, name=Ivy}
10. [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, status=active, level=1, name=X, age=22, score=6.5, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, department=sales, salary=45000, active=true, level=2, name=Alice, age=25, score=8.5, tags=junior, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, department=engineering, level=5, name=Bob, active=true, status=active}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, tags=manager, status=active, department=marketing, level=7, name=Diana, age=45, active=true, score=7.8}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, salary=55000, tags=employee, status=inactive, department=sales, level=3}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, age=65, tags=executive, level=9, name=Grace, salary=95000, active=true, score=10, status=active, department=management}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering, name=Ivy}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, amount=2, status=pending, region=north, customer_id=P001, product_id=PROD001, total=1999.98, date=2024-01-15, priority=normal, discount=50}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, status=confirmed, priority=low, discount=0, region=south, customer_id=P002, total=25.5, date=2024-01-20}
3. [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, amount=3, discount=15}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, date=2024-02-05, status=delivered, customer_id=P004, total=299.99, priority=normal, discount=0, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, priority=high, discount=100, amount=1, total=999.99, region=south}
6. [1] TestOrder{id=O006, region=west, customer_id=P005, amount=2, date=2024-02-15, priority=low, discount=0, product_id=PROD005, total=999.98, status=cancelled}
7. [1] TestOrder{id=O007, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, total=89.99, discount=10, product_id=PROD007, amount=1, date=2024-03-10, status=completed, priority=low, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, date=2024-03-15, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, amount=2, status=pending, region=north, customer_id=P001, product_id=PROD001, total=1999.98, date=2024-01-15, priority=normal, discount=50}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, total=25.5, date=2024-01-20, product_id=PROD002, amount=1, status=confirmed, priority=low}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, amount=3, discount=15}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, status=confirmed, priority=high, discount=100, amount=1, total=999.99, region=south, customer_id=P002, product_id=PROD001, date=2024-02-10}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, date=2024-02-05, status=delivered, customer_id=P004, total=299.99, priority=normal, discount=0, region=east, product_id=PROD004}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, amount=2, date=2024-02-15, priority=low, discount=0, product_id=PROD005, total=999.98, status=cancelled, region=west}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O008, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, amount=1, date=2024-03-10, status=completed, priority=low, region=north, customer_id=P001, total=89.99, discount=10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, rating=4.5, keywords=computer, name=Laptop, category=electronics, available=true, brand=TechCorp, stock=50, supplier=TechSupply, price=999.99}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, price=25.5, brand=TechCorp, stock=200, available=true, rating=4, keywords=peripheral, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, category=accessories, price=75, available=false, rating=3.5, keywords=typing, supplier=KeySupply, name=Keyboard, brand=KeyTech, stock=0}
4. [1] TestProduct{id=PROD004, available=true, keywords=display, name=Monitor, price=299.99, rating=4.8, brand=ScreenPro, stock=30, supplier=ScreenSupply, category=electronics}
5. [1] TestProduct{id=PROD005, price=8.5, rating=2, keywords=obsolete, brand=OldTech, stock=0, name=OldKeyboard, category=accessories, available=false, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, supplier=AudioSupply, category=audio, keywords=sound, name=Headphones, price=150, available=true, rating=4.6, brand=AudioMax, stock=75}
7. [1] TestProduct{id=PROD007, rating=3.8, supplier=CamSupply, price=89.99, keywords=video, brand=CamTech, stock=25, name=Webcam, category=electronics, available=true}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, category=accessories, price=25.5, brand=TechCorp, stock=200, available=true, rating=4, keywords=peripheral, supplier=TechSupply}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD003, supplier=KeySupply, name=Keyboard, brand=KeyTech, stock=0, category=accessories, price=75, available=false, rating=3.5, keywords=typing}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, price=299.99, rating=4.8, brand=ScreenPro, stock=30, supplier=ScreenSupply, category=electronics, available=true, keywords=display}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD006, price=150, available=true, rating=4.6, brand=AudioMax, stock=75, supplier=AudioSupply, category=audio, keywords=sound, name=Headphones}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD007, price=89.99, keywords=video, brand=CamTech, stock=25, name=Webcam, category=electronics, available=true, rating=3.8, supplier=CamSupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD001, available=true, brand=TechCorp, stock=50, supplier=TechSupply, price=999.99, rating=4.5, keywords=computer, name=Laptop, category=electronics}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, score=8.5, tags=junior, status=active, department=sales, salary=45000, active=true, level=2}
2. [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, department=engineering, level=5, name=Bob, active=true, status=active}
3. [1] TestPerson{id=P003, age=16, tags=intern, status=inactive, department=hr, level=1, name=Charlie, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, salary=85000, tags=manager, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, department=sales, level=3, name=Eve, age=30, active=false, score=8}
6. [1] TestPerson{id=P006, tags=test, name=Frank, age=0, salary=-5000, active=true, status=active, department=qa, level=1, score=0}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, score=10, status=active, department=management, age=65, tags=executive, level=9}
8. [1] TestPerson{id=P008, age=18, active=false, score=5.5, tags=junior, department=support, level=1, name=Henry, salary=25000, status=inactive}
9. [1] TestPerson{id=P009, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering, name=Ivy, age=40, active=true}
10. [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, status=active, level=1, name=X, age=22, score=6.5, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, age=16, tags=intern, status=inactive, department=hr, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, active=true, status=active, department=qa, level=1, score=0, tags=test, name=Frank, age=0, salary=-5000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, region=north, customer_id=P001, product_id=PROD001, total=1999.98, date=2024-01-15, priority=normal, discount=50, amount=2, status=pending}
2. [1] TestOrder{id=O002, date=2024-01-20, product_id=PROD002, amount=1, status=confirmed, priority=low, discount=0, region=south, customer_id=P002, total=25.5}
3. [1] TestOrder{id=O003, region=north, customer_id=P001, amount=3, discount=15, product_id=PROD003, total=225, date=2024-02-01, status=shipped, priority=high}
4. [1] TestOrder{id=O004, total=299.99, priority=normal, discount=0, region=east, product_id=PROD004, amount=1, date=2024-02-05, status=delivered, customer_id=P004}
5. [1] TestOrder{id=O005, discount=100, amount=1, total=999.99, region=south, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, priority=high}
6. [1] TestOrder{id=O006, date=2024-02-15, priority=low, discount=0, product_id=PROD005, total=999.98, status=cancelled, region=west, customer_id=P005, amount=2}
7. [1] TestOrder{id=O007, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01}
8. [1] TestOrder{id=O008, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal}
9. [1] TestOrder{id=O009, product_id=PROD007, amount=1, date=2024-03-10, status=completed, priority=low, region=north, customer_id=P001, total=89.99, discount=10}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, date=2024-03-15, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, total=299.99, priority=normal, discount=0, region=east, product_id=PROD004, amount=1, date=2024-02-05, status=delivered}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O008, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O009, priority=low, region=north, customer_id=P001, total=89.99, discount=10, product_id=PROD007, amount=1, date=2024-03-10, status=completed}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, status=confirmed, priority=low, discount=0, region=south, customer_id=P002, total=25.5, date=2024-01-20}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, discount=15, product_id=PROD003, total=225, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, amount=3}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O005, total=999.99, region=south, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, priority=high, discount=100, amount=1}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, region=west, customer_id=P005, amount=2, date=2024-02-15, priority=low, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O010, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, amount=1, date=2024-03-15, product_id=PROD001, total=75000}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, date=2024-01-15, priority=normal, discount=50, amount=2, status=pending, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, department=sales, salary=45000, active=true, level=2, name=Alice, age=25, score=8.5, tags=junior}
2. [1] TestPerson{id=P002, name=Bob, active=true, status=active, age=35, salary=75000, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, age=16, tags=intern, status=inactive, department=hr, level=1, name=Charlie, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, status=active, department=marketing, level=7, name=Diana, age=45, active=true, score=7.8, salary=85000, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, salary=55000, tags=employee, status=inactive, department=sales, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, status=active, department=qa, level=1, score=0, tags=test}
7. [1] TestPerson{id=P007, department=management, age=65, tags=executive, level=9, name=Grace, salary=95000, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, tags=junior, department=support, level=1, name=Henry, salary=25000, status=inactive, age=18, active=false, score=5.5}
9. [1] TestPerson{id=P009, salary=68000, score=8.7, department=engineering, name=Ivy, age=40, active=true, tags=senior, status=active, level=6}
10. [1] TestPerson{id=P010, level=1, name=X, age=22, score=6.5, department=intern, salary=28000, active=true, tags=temp, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, age=18, active=false, score=5.5, tags=junior, department=support, level=1, name=Henry, salary=25000, status=inactive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, active=false, score=6, age=16, tags=intern, status=inactive, department=hr, level=1, name=Charlie}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, salary=55000, tags=employee, status=inactive, department=sales, level=3}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, amount=2, status=pending, region=north, customer_id=P001, product_id=PROD001, total=1999.98, date=2024-01-15, priority=normal, discount=50}
2. [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, product_id=PROD002, amount=1, status=confirmed, priority=low, discount=0, region=south}
3. [1] TestOrder{id=O003, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, amount=3, discount=15, product_id=PROD003, total=225}
4. [1] TestOrder{id=O004, priority=normal, discount=0, region=east, product_id=PROD004, amount=1, date=2024-02-05, status=delivered, customer_id=P004, total=299.99}
5. [1] TestOrder{id=O005, discount=100, amount=1, total=999.99, region=south, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, priority=high}
6. [1] TestOrder{id=O006, region=west, customer_id=P005, amount=2, date=2024-02-15, priority=low, discount=0, product_id=PROD005, total=999.98, status=cancelled}
7. [1] TestOrder{id=O007, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent}
8. [1] TestOrder{id=O008, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0}
9. [1] TestOrder{id=O009, status=completed, priority=low, region=north, customer_id=P001, total=89.99, discount=10, product_id=PROD007, amount=1, date=2024-03-10}
10. [1] TestOrder{id=O010, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, amount=1, date=2024-03-15}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O005, priority=high, discount=100, amount=1, total=999.99, region=south, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, region=west, customer_id=P005, amount=2, date=2024-02-15, priority=low, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O007, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, discount=10, product_id=PROD007, amount=1, date=2024-03-10, status=completed, priority=low, region=north, customer_id=P001, total=89.99}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, amount=1, date=2024-03-15, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0, region=east}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O002, total=25.5, date=2024-01-20, product_id=PROD002, amount=1, status=confirmed, priority=low, discount=0, region=south, customer_id=P002}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, amount=3, discount=15}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, amount=1, date=2024-02-05, status=delivered, customer_id=P004, total=299.99, priority=normal, discount=0, region=east}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, active=true, level=2, name=Alice, age=25, score=8.5, tags=junior, status=active, department=sales}
2. [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, department=engineering, level=5, name=Bob, active=true, status=active}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, age=16, tags=intern, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, salary=85000, tags=manager, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, score=8, salary=55000, tags=employee, status=inactive, department=sales, level=3, name=Eve, age=30, active=false}
6. [1] TestPerson{id=P006, score=0, tags=test, name=Frank, age=0, salary=-5000, active=true, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, level=9, name=Grace, salary=95000, active=true, score=10, status=active, department=management, age=65, tags=executive}
8. [1] TestPerson{id=P008, age=18, active=false, score=5.5, tags=junior, department=support, level=1, name=Henry, salary=25000, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering}
10. [1] TestPerson{id=P010, tags=temp, status=active, level=1, name=X, age=22, score=6.5, department=intern, salary=28000, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, salary=55000, tags=employee, status=inactive, department=sales, level=3}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, active=true, score=10, status=active, department=management, age=65, tags=executive, level=9}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, status=inactive, age=18, active=false, score=5.5, tags=junior, department=support, level=1}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, score=8.5, tags=junior, status=active, department=sales, salary=45000, active=true, level=2}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, department=engineering, level=5, name=Bob, active=true, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, salary=85000, tags=manager, status=active, department=marketing, level=7}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, name=Frank, age=0, salary=-5000, active=true, status=active, department=qa, level=1, score=0}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P003, department=hr, level=1, name=Charlie, salary=0, active=false, score=6, age=16, tags=intern, status=inactive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, supplier=TechSupply, price=999.99, rating=4.5, keywords=computer, name=Laptop, category=electronics, available=true, brand=TechCorp, stock=50}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, price=25.5, brand=TechCorp, stock=200, available=true, rating=4, keywords=peripheral, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, brand=KeyTech, stock=0, category=accessories, price=75, available=false, rating=3.5, keywords=typing, supplier=KeySupply, name=Keyboard}
4. [1] TestProduct{id=PROD004, supplier=ScreenSupply, category=electronics, available=true, keywords=display, name=Monitor, price=299.99, rating=4.8, brand=ScreenPro, stock=30}
5. [1] TestProduct{id=PROD005, brand=OldTech, stock=0, name=OldKeyboard, category=accessories, available=false, supplier=OldSupply, price=8.5, rating=2, keywords=obsolete}
6. [1] TestProduct{id=PROD006, price=150, available=true, rating=4.6, brand=AudioMax, stock=75, supplier=AudioSupply, category=audio, keywords=sound, name=Headphones}
7. [1] TestProduct{id=PROD007, brand=CamTech, stock=25, name=Webcam, category=electronics, available=true, rating=3.8, supplier=CamSupply, price=89.99, keywords=video}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, price=299.99, rating=4.8, brand=ScreenPro, stock=30, supplier=ScreenSupply, category=electronics, available=true, keywords=display}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD006, stock=75, supplier=AudioSupply, category=audio, keywords=sound, name=Headphones, price=150, available=true, rating=4.6, brand=AudioMax}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD007, supplier=CamSupply, price=89.99, keywords=video, brand=CamTech, stock=25, name=Webcam, category=electronics, available=true, rating=3.8}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD001, stock=50, supplier=TechSupply, price=999.99, rating=4.5, keywords=computer, name=Laptop, category=electronics, available=true, brand=TechCorp}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD002, price=25.5, brand=TechCorp, stock=200, available=true, rating=4, keywords=peripheral, supplier=TechSupply, name=Mouse, category=accessories}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD003, brand=KeyTech, stock=0, category=accessories, price=75, available=false, rating=3.5, keywords=typing, supplier=KeySupply, name=Keyboard}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, score=8.5, tags=junior, status=active, department=sales, salary=45000, active=true, level=2}
2. [1] TestPerson{id=P002, level=5, name=Bob, active=true, status=active, age=35, salary=75000, score=9.2, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, score=6, age=16, tags=intern, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, salary=85000, tags=manager, status=active, department=marketing, level=7, name=Diana, age=45, active=true, score=7.8}
5. [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, department=sales, level=3, name=Eve, age=30, active=false, score=8}
6. [1] TestPerson{id=P006, tags=test, name=Frank, age=0, salary=-5000, active=true, status=active, department=qa, level=1, score=0}
7. [1] TestPerson{id=P007, age=65, tags=executive, level=9, name=Grace, salary=95000, active=true, score=10, status=active, department=management}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, status=inactive, age=18, active=false, score=5.5, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering}
10. [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, status=active, level=1, name=X, age=22, score=6.5, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P006, score=0, tags=test, name=Frank, age=0, salary=-5000, active=true, status=active, department=qa, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, tags=executive, level=9, name=Grace, salary=95000, active=true, score=10, status=active, department=management, age=65}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, level=6, salary=68000, score=8.7, department=engineering, name=Ivy, age=40, active=true, tags=senior, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, score=8.5, tags=junior, status=active, department=sales, salary=45000, active=true, level=2}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, active=true, status=active, age=35, salary=75000, score=9.2, tags=senior, department=engineering, level=5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, age=16, tags=intern, status=inactive, department=hr, level=1, name=Charlie, salary=0, active=false, score=6}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, active=true, score=7.8, salary=85000, tags=manager, status=active, department=marketing, level=7, name=Diana, age=45}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, department=support, level=1, name=Henry, salary=25000, status=inactive, age=18, active=false, score=5.5, tags=junior}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, salary=55000, tags=employee, status=inactive, department=sales, level=3}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, date=2024-01-15, priority=normal, discount=50, amount=2, status=pending, region=north, customer_id=P001}
2. [1] TestOrder{id=O002, amount=1, status=confirmed, priority=low, discount=0, region=south, customer_id=P002, total=25.5, date=2024-01-20, product_id=PROD002}
3. [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, amount=3, discount=15}
4. [1] TestOrder{id=O004, customer_id=P004, total=299.99, priority=normal, discount=0, region=east, product_id=PROD004, amount=1, date=2024-02-05, status=delivered}
5. [1] TestOrder{id=O005, status=confirmed, priority=high, discount=100, amount=1, total=999.99, region=south, customer_id=P002, product_id=PROD001, date=2024-02-10}
6. [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, region=west, customer_id=P005, amount=2, date=2024-02-15, priority=low, discount=0}
7. [1] TestOrder{id=O007, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50}
8. [1] TestOrder{id=O008, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south, customer_id=P010}
9. [1] TestOrder{id=O009, product_id=PROD007, amount=1, date=2024-03-10, status=completed, priority=low, region=north, customer_id=P001, total=89.99, discount=10}
10. [1] TestOrder{id=O010, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, amount=1, date=2024-03-15}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, priority=low, discount=0, region=south, customer_id=P002, total=25.5, date=2024-01-20, product_id=PROD002, amount=1, status=confirmed}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, amount=3, discount=15, product_id=PROD003, total=225}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, discount=0, region=east, product_id=PROD004, amount=1, date=2024-02-05, status=delivered, customer_id=P004, total=299.99, priority=normal}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, priority=high, discount=100, amount=1, total=999.99, region=south}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, amount=1, date=2024-03-10, status=completed, priority=low, region=north, customer_id=P001, total=89.99, discount=10}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, date=2024-01-15, priority=normal, discount=50, amount=2, status=pending, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, score=8.5, tags=junior, status=active, department=sales, salary=45000, active=true, level=2}
2. [1] TestPerson{id=P002, level=5, name=Bob, active=true, status=active, age=35, salary=75000, score=9.2, tags=senior, department=engineering}
3. [1] TestPerson{id=P003, status=inactive, department=hr, level=1, name=Charlie, salary=0, active=false, score=6, age=16, tags=intern}
4. [1] TestPerson{id=P004, salary=85000, tags=manager, status=active, department=marketing, level=7, name=Diana, age=45, active=true, score=7.8}
5. [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, department=sales, level=3, name=Eve, age=30, active=false, score=8}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, status=active, department=qa, level=1, score=0, tags=test}
7. [1] TestPerson{id=P007, salary=95000, active=true, score=10, status=active, department=management, age=65, tags=executive, level=9, name=Grace}
8. [1] TestPerson{id=P008, level=1, name=Henry, salary=25000, status=inactive, age=18, active=false, score=5.5, tags=junior, department=support}
9. [1] TestPerson{id=P009, salary=68000, score=8.7, department=engineering, name=Ivy, age=40, active=true, tags=senior, status=active, level=6}
10. [1] TestPerson{id=P010, department=intern, salary=28000, active=true, tags=temp, status=active, level=1, name=X, age=22, score=6.5}
11. [1] TestOrder{id=O001, discount=50, amount=2, status=pending, region=north, customer_id=P001, product_id=PROD001, total=1999.98, date=2024-01-15, priority=normal}
12. [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, product_id=PROD002, amount=1, status=confirmed, priority=low, discount=0, region=south}
13. [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, amount=3, discount=15}
14. [1] TestOrder{id=O004, product_id=PROD004, amount=1, date=2024-02-05, status=delivered, customer_id=P004, total=299.99, priority=normal, discount=0, region=east}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, priority=high, discount=100, amount=1, total=999.99, region=south}
16. [1] TestOrder{id=O006, status=cancelled, region=west, customer_id=P005, amount=2, date=2024-02-15, priority=low, discount=0, product_id=PROD005, total=999.98}
17. [1] TestOrder{id=O007, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, date=2024-03-01, status=shipped, discount=50}
18. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}
19. [1] TestOrder{id=O009, customer_id=P001, total=89.99, discount=10, product_id=PROD007, amount=1, date=2024-03-10, status=completed, priority=low, region=north}
20. [1] TestOrder{id=O010, customer_id=P006, amount=1, date=2024-03-15, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0, region=east}

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

1. [1] TestPerson{id=P001, status=active, department=sales, salary=45000, active=true, level=2, name=Alice, age=25, score=8.5, tags=junior}
2. [1] TestPerson{id=P002, department=engineering, level=5, name=Bob, active=true, status=active, age=35, salary=75000, score=9.2, tags=senior}
3. [1] TestPerson{id=P003, salary=0, active=false, score=6, age=16, tags=intern, status=inactive, department=hr, level=1, name=Charlie}
4. [1] TestPerson{id=P004, level=7, name=Diana, age=45, active=true, score=7.8, salary=85000, tags=manager, status=active, department=marketing}
5. [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, department=sales, level=3, name=Eve, age=30, active=false, score=8}
6. [1] TestPerson{id=P006, salary=-5000, active=true, status=active, department=qa, level=1, score=0, tags=test, name=Frank, age=0}
7. [1] TestPerson{id=P007, department=management, age=65, tags=executive, level=9, name=Grace, salary=95000, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, score=5.5, tags=junior, department=support, level=1, name=Henry, salary=25000, status=inactive, age=18, active=false}
9. [1] TestPerson{id=P009, age=40, active=true, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering, name=Ivy}
10. [1] TestPerson{id=P010, name=X, age=22, score=6.5, department=intern, salary=28000, active=true, tags=temp, status=active, level=1}

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

1. [1] TestPerson{id=P001, salary=45000, active=true, level=2, name=Alice, age=25, score=8.5, tags=junior, status=active, department=sales}
2. [1] TestPerson{id=P002, name=Bob, active=true, status=active, age=35, salary=75000, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, tags=intern, status=inactive, department=hr, level=1, name=Charlie, salary=0, active=false, score=6, age=16}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, salary=85000, tags=manager, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, name=Eve, age=30, active=false, score=8, salary=55000, tags=employee, status=inactive, department=sales, level=3}
6. [1] TestPerson{id=P006, age=0, salary=-5000, active=true, status=active, department=qa, level=1, score=0, tags=test, name=Frank}
7. [1] TestPerson{id=P007, age=65, tags=executive, level=9, name=Grace, salary=95000, active=true, score=10, status=active, department=management}
8. [1] TestPerson{id=P008, age=18, active=false, score=5.5, tags=junior, department=support, level=1, name=Henry, salary=25000, status=inactive}
9. [1] TestPerson{id=P009, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering, name=Ivy, age=40, active=true}
10. [1] TestPerson{id=P010, tags=temp, status=active, level=1, name=X, age=22, score=6.5, department=intern, salary=28000, active=true}
11. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, date=2024-01-15, priority=normal, discount=50, amount=2, status=pending, region=north}
12. [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, product_id=PROD002, amount=1, status=confirmed, priority=low, discount=0, region=south}
13. [1] TestOrder{id=O003, amount=3, discount=15, product_id=PROD003, total=225, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001}
14. [1] TestOrder{id=O004, priority=normal, discount=0, region=east, product_id=PROD004, amount=1, date=2024-02-05, status=delivered, customer_id=P004, total=299.99}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, priority=high, discount=100, amount=1, total=999.99, region=south}
16. [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, region=west, customer_id=P005, amount=2, date=2024-02-15, priority=low, discount=0}
17. [1] TestOrder{id=O007, amount=4, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006}
18. [1] TestOrder{id=O008, status=pending, discount=0, amount=10, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, date=2024-03-05}
19. [1] TestOrder{id=O009, product_id=PROD007, amount=1, date=2024-03-10, status=completed, priority=low, region=north, customer_id=P001, total=89.99, discount=10}
20. [1] TestOrder{id=O010, date=2024-03-15, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, amount=1}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, score=8.5, tags=junior, status=active, department=sales, salary=45000, active=true, level=2}
2. [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, tags=senior, department=engineering, level=5, name=Bob, active=true, status=active}
3. [1] TestPerson{id=P003, age=16, tags=intern, status=inactive, department=hr, level=1, name=Charlie, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, salary=85000, tags=manager, status=active, department=marketing, level=7, name=Diana, age=45, active=true, score=7.8}
5. [1] TestPerson{id=P005, status=inactive, department=sales, level=3, name=Eve, age=30, active=false, score=8, salary=55000, tags=employee}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, status=active, department=qa, level=1, score=0, tags=test}
7. [1] TestPerson{id=P007, age=65, tags=executive, level=9, name=Grace, salary=95000, active=true, score=10, status=active, department=management}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, status=inactive, age=18, active=false, score=5.5, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, tags=senior, status=active, level=6, salary=68000, score=8.7, department=engineering}
10. [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, status=active, level=1, name=X, age=22, score=6.5, department=intern}

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

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, date=2024-01-15, priority=normal, discount=50, amount=2, status=pending, region=north}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, status=confirmed, priority=low, discount=0, region=south, customer_id=P002, total=25.5, date=2024-01-20}
3. [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, amount=3, discount=15}
4. [1] TestOrder{id=O004, customer_id=P004, total=299.99, priority=normal, discount=0, region=east, product_id=PROD004, amount=1, date=2024-02-05, status=delivered}
5. [1] TestOrder{id=O005, amount=1, total=999.99, region=south, customer_id=P002, product_id=PROD001, date=2024-02-10, status=confirmed, priority=high, discount=100}
6. [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, region=west, customer_id=P005, amount=2, date=2024-02-15, priority=low, discount=0}
7. [1] TestOrder{id=O007, total=600, date=2024-03-01, status=shipped, discount=50, customer_id=P007, priority=urgent, region=north, product_id=PROD006, amount=4}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, discount=0, amount=10, total=255, priority=normal, region=south}
9. [1] TestOrder{id=O009, product_id=PROD007, amount=1, date=2024-03-10, status=completed, priority=low, region=north, customer_id=P001, total=89.99, discount=10}
10. [1] TestOrder{id=O010, amount=1, date=2024-03-15, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, score=8.5, tags=junior, status=active, department=sales, salary=45000, active=true, level=2}
2. [1] TestPerson{id=P002, tags=senior, department=engineering, level=5, name=Bob, active=true, status=active, age=35, salary=75000, score=9.2}
3. [1] TestPerson{id=P003, age=16, tags=intern, status=inactive, department=hr, level=1, name=Charlie, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, name=Diana, age=45, active=true, score=7.8, salary=85000, tags=manager, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, salary=55000, tags=employee, status=inactive, department=sales, level=3, name=Eve, age=30, active=false, score=8}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, active=true, status=active, department=qa, level=1, score=0, tags=test}
7. [1] TestPerson{id=P007, age=65, tags=executive, level=9, name=Grace, salary=95000, active=true, score=10, status=active, department=management}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, status=inactive, age=18, active=false, score=5.5, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, salary=68000, score=8.7, department=engineering, name=Ivy, age=40, active=true, tags=senior, status=active, level=6}
10. [1] TestPerson{id=P010, active=true, tags=temp, status=active, level=1, name=X, age=22, score=6.5, department=intern, salary=28000}

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
