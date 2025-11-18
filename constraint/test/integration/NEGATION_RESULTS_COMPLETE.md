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

1. [1] TestPerson{id=P001, active=true, level=2, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales}
2. [1] TestPerson{id=P002, department=engineering, level=5, active=true, tags=senior, name=Bob, age=35, salary=75000, score=9.2, status=active}
3. [1] TestPerson{id=P003, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, salary=85000, active=true, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, score=8, department=sales, age=30, tags=employee, status=inactive, level=3}
6. [1] TestPerson{id=P006, status=active, level=1, active=true, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000}
7. [1] TestPerson{id=P007, status=active, department=management, active=true, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive}
8. [1] TestPerson{id=P008, level=1, salary=25000, name=Henry, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support}
9. [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, tags=temp, department=intern, level=1, score=6.5, status=active, name=X, age=22, salary=28000, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, status=active, department=engineering, level=5, active=true, tags=senior, name=Bob, age=35, salary=75000, score=9.2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, name=Diana, age=45, score=7.8, tags=manager, salary=85000, active=true, status=active}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, status=inactive, level=3, name=Eve, salary=55000, active=false, score=8, department=sales, age=30}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, department=management, active=true, level=9}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, age=40, active=true, status=active, department=engineering, level=6}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, active=true, tags=temp, department=intern, level=1, score=6.5, status=active, name=X, age=22, salary=28000}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P003, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1, age=16}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, salary=25000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north, total=1999.98, discount=50}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, priority=low, customer_id=P002, total=25.5, date=2024-01-20, status=confirmed, discount=0, region=south}
3. [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, status=shipped, customer_id=P001, amount=3, priority=high, discount=15, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, date=2024-02-05, discount=0, amount=1, total=299.99, status=delivered, priority=normal, region=east}
5. [1] TestOrder{id=O005, discount=100, total=999.99, region=south, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high}
6. [1] TestOrder{id=O006, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, product_id=PROD005, amount=2, status=cancelled, region=west}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50, region=north, total=600, date=2024-03-01, priority=urgent}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}
9. [1] TestOrder{id=O009, status=completed, priority=low, customer_id=P001, amount=1, total=89.99, date=2024-03-10, discount=10, region=north, product_id=PROD007}
10. [1] TestOrder{id=O010, status=refunded, priority=urgent, region=east, product_id=PROD001, total=75000, date=2024-03-15, discount=0, customer_id=P006, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O010, discount=0, customer_id=P006, amount=1, status=refunded, priority=urgent, region=east, product_id=PROD001, total=75000, date=2024-03-15}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, priority=low, customer_id=P002, total=25.5, date=2024-01-20, status=confirmed, discount=0, region=south}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, date=2024-03-10, discount=10, region=north, product_id=PROD007, status=completed, priority=low, customer_id=P001, amount=1, total=89.99}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north, total=1999.98, discount=50}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, priority=high, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped, customer_id=P001}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high, discount=100, total=999.99, region=south}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50, region=north, total=600, date=2024-03-01, priority=urgent}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}
2. [1] TestPerson{id=P002, tags=senior, name=Bob, age=35, salary=75000, score=9.2, status=active, department=engineering, level=5, active=true}
3. [1] TestPerson{id=P003, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false}
4. [1] TestPerson{id=P004, salary=85000, active=true, status=active, department=marketing, level=7, name=Diana, age=45, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, age=30, tags=employee, status=inactive, level=3, name=Eve, salary=55000, active=false, score=8, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, status=active, level=1, active=true, score=0, tags=test, department=qa}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, department=management, active=true, level=9}
8. [1] TestPerson{id=P008, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, salary=25000, name=Henry, age=18}
9. [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, tags=temp, department=intern, level=1, score=6.5, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, age=40, active=true, status=active, department=engineering, level=6}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, score=9.2, status=active, department=engineering, level=5, active=true, tags=senior}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, salary=85000, active=true, status=active, department=marketing, level=7}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, level=3, name=Eve, salary=55000, active=false, score=8, department=sales, age=30, tags=employee}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, active=true, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, department=management}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, region=north, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, priority=low, customer_id=P002, total=25.5, date=2024-01-20, status=confirmed, discount=0, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, priority=high, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped}
4. [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, discount=0}
5. [1] TestOrder{id=O005, region=south, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high, discount=100, total=999.99}
6. [1] TestOrder{id=O006, priority=low, discount=0, customer_id=P005, product_id=PROD005, amount=2, status=cancelled, region=west, total=999.98, date=2024-02-15}
7. [1] TestOrder{id=O007, region=north, total=600, date=2024-03-01, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50}
8. [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, discount=10, region=north, product_id=PROD007, status=completed, priority=low}
10. [1] TestOrder{id=O010, product_id=PROD001, total=75000, date=2024-03-15, discount=0, customer_id=P006, amount=1, status=refunded, priority=urgent, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped, customer_id=P001, amount=3, priority=high}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, date=2024-02-05, discount=0, amount=1, total=299.99, status=delivered, priority=normal, region=east, customer_id=P004}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high, discount=100, total=999.99, region=south}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O006, region=west, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, product_id=PROD005, amount=2, status=cancelled}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north, total=1999.98, discount=50}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, priority=low, customer_id=P002, total=25.5}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, status=shipped, discount=50, region=north, total=600, date=2024-03-01, priority=urgent, customer_id=P007}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, status=completed, priority=low, customer_id=P001, amount=1, total=89.99, date=2024-03-10, discount=10, region=north, product_id=PROD007}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, name=Laptop, rating=4.5, brand=TechCorp, supplier=TechSupply, category=electronics, price=999.99, available=true, keywords=computer, stock=50}
2. [1] TestProduct{id=PROD002, keywords=peripheral, stock=200, category=accessories, available=true, rating=4, brand=TechCorp, supplier=TechSupply, name=Mouse, price=25.5}
3. [1] TestProduct{id=PROD003, name=Keyboard, brand=KeyTech, stock=0, category=accessories, price=75, available=false, rating=3.5, keywords=typing, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, category=electronics, rating=4.8, keywords=display, stock=30, name=Monitor, price=299.99, available=true, brand=ScreenPro, supplier=ScreenSupply}
5. [1] TestProduct{id=PROD005, supplier=OldSupply, name=OldKeyboard, category=accessories, price=8.5, available=false, rating=2, keywords=obsolete, brand=OldTech, stock=0}
6. [1] TestProduct{id=PROD006, available=true, rating=4.6, category=audio, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply, name=Headphones, price=150}
7. [1] TestProduct{id=PROD007, price=89.99, available=true, supplier=CamSupply, name=Webcam, rating=3.8, keywords=video, brand=CamTech, stock=25, category=electronics}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD004, category=electronics, rating=4.8, keywords=display, stock=30, name=Monitor, price=299.99, available=true, brand=ScreenPro, supplier=ScreenSupply}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD006, category=audio, keywords=sound, brand=AudioMax, stock=75, supplier=AudioSupply, name=Headphones, price=150, available=true, rating=4.6}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD007, brand=CamTech, stock=25, category=electronics, price=89.99, available=true, supplier=CamSupply, name=Webcam, rating=3.8, keywords=video}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD001, category=electronics, price=999.99, available=true, keywords=computer, stock=50, name=Laptop, rating=4.5, brand=TechCorp, supplier=TechSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD002, category=accessories, available=true, rating=4, brand=TechCorp, supplier=TechSupply, name=Mouse, price=25.5, keywords=peripheral, stock=200}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD003, brand=KeyTech, stock=0, category=accessories, price=75, available=false, rating=3.5, keywords=typing, supplier=KeySupply, name=Keyboard}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, score=9.2, status=active, department=engineering, level=5, active=true, tags=senior}
3. [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, salary=85000, active=true, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, salary=55000, active=false, score=8, department=sales, age=30, tags=employee, status=inactive, level=3, name=Eve}
6. [1] TestPerson{id=P006, age=0, salary=-5000, status=active, level=1, active=true, score=0, tags=test, department=qa, name=Frank}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, department=management, active=true, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, salary=25000}
9. [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, tags=temp, department=intern, level=1, score=6.5, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P006, level=1, active=true, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, region=north, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal}
2. [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, priority=low}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, priority=high, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped}
4. [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, discount=0}
5. [1] TestOrder{id=O005, status=confirmed, priority=high, discount=100, total=999.99, region=south, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10}
6. [1] TestOrder{id=O006, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, product_id=PROD005, amount=2, status=cancelled, region=west}
7. [1] TestOrder{id=O007, region=north, total=600, date=2024-03-01, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}
9. [1] TestOrder{id=O009, region=north, product_id=PROD007, status=completed, priority=low, customer_id=P001, amount=1, total=89.99, date=2024-03-10, discount=10}
10. [1] TestOrder{id=O010, product_id=PROD001, total=75000, date=2024-03-15, discount=0, customer_id=P006, amount=1, status=refunded, priority=urgent, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, discount=0, region=south, product_id=PROD002, amount=1, priority=low, customer_id=P002, total=25.5, date=2024-01-20, status=confirmed}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O007, region=north, total=600, date=2024-03-01, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O008, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O009, discount=10, region=north, product_id=PROD007, status=completed, priority=low, customer_id=P001, amount=1, total=89.99, date=2024-03-10}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O010, priority=urgent, region=east, product_id=PROD001, total=75000, date=2024-03-15, discount=0, customer_id=P006, amount=1, status=refunded}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north, total=1999.98, discount=50, customer_id=P001}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O003, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped, customer_id=P001, amount=3, priority=high, discount=15}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, discount=0, amount=1, total=299.99, status=delivered, priority=normal}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high, discount=100, total=999.99, region=south}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O006, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, product_id=PROD005, amount=2, status=cancelled, region=west}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true}
2. [1] TestPerson{id=P002, status=active, department=engineering, level=5, active=true, tags=senior, name=Bob, age=35, salary=75000, score=9.2}
3. [1] TestPerson{id=P003, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, tags=manager, salary=85000, active=true, status=active, department=marketing, level=7, name=Diana, age=45, score=7.8}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, score=8, department=sales, age=30, tags=employee, status=inactive, level=3}
6. [1] TestPerson{id=P006, active=true, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000, status=active, level=1}
7. [1] TestPerson{id=P007, score=10, tags=executive, status=active, department=management, active=true, level=9, name=Grace, age=65, salary=95000}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, salary=25000}
9. [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, status=active, name=X, age=22, salary=28000, active=true, tags=temp, department=intern, level=1, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, tags=junior, status=inactive, department=support, level=1, salary=25000, name=Henry, age=18, active=false}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, level=3, name=Eve, salary=55000, active=false, score=8, department=sales, age=30, tags=employee}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, status=pending, priority=normal, region=north, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, amount=2, date=2024-01-15}
2. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, priority=low, customer_id=P002}
3. [1] TestOrder{id=O003, priority=high, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped, customer_id=P001, amount=3}
4. [1] TestOrder{id=O004, total=299.99, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, discount=0, amount=1}
5. [1] TestOrder{id=O005, region=south, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high, discount=100, total=999.99}
6. [1] TestOrder{id=O006, priority=low, discount=0, customer_id=P005, product_id=PROD005, amount=2, status=cancelled, region=west, total=999.98, date=2024-02-15}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50, region=north, total=600, date=2024-03-01, priority=urgent}
8. [1] TestOrder{id=O008, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010}
9. [1] TestOrder{id=O009, region=north, product_id=PROD007, status=completed, priority=low, customer_id=P001, amount=1, total=89.99, date=2024-03-10, discount=10}
10. [1] TestOrder{id=O010, product_id=PROD001, total=75000, date=2024-03-15, discount=0, customer_id=P006, amount=1, status=refunded, priority=urgent, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O010, discount=0, customer_id=P006, amount=1, status=refunded, priority=urgent, region=east, product_id=PROD001, total=75000, date=2024-03-15}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, amount=1, priority=low, customer_id=P002, total=25.5, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, status=shipped, customer_id=P001, amount=3, priority=high, discount=15, region=north}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, discount=0}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high, discount=100, total=999.99, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, amount=2, status=cancelled, region=west, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, total=600, date=2024-03-01, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, status=completed, priority=low, customer_id=P001, amount=1, total=89.99, date=2024-03-10, discount=10, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2, name=Alice, age=25}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, score=9.2, status=active, department=engineering, level=5, active=true, tags=senior}
3. [1] TestPerson{id=P003, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1}
4. [1] TestPerson{id=P004, salary=85000, active=true, status=active, department=marketing, level=7, name=Diana, age=45, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, score=8, department=sales, age=30, tags=employee, status=inactive, level=3, name=Eve, salary=55000, active=false}
6. [1] TestPerson{id=P006, active=true, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000, status=active, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, department=management, active=true, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, salary=25000}
9. [1] TestPerson{id=P009, score=8.7, tags=senior, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, tags=temp, department=intern, level=1, score=6.5, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, age=30, tags=employee, status=inactive, level=3, name=Eve, salary=55000, active=false, score=8, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, status=active, department=management, active=true, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, level=1, salary=25000, name=Henry, age=18, active=false, score=5.5}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, status=active, department=sales, active=true, level=2, name=Alice, age=25, salary=45000, score=8.5}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P003, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, salary=85000, active=true, status=active, department=marketing, level=7, name=Diana, age=45}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, status=active, level=1, active=true, score=0, tags=test, department=qa}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior, age=40, active=true}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, score=9.2, status=active, department=engineering, level=5, active=true, tags=senior}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, category=electronics, price=999.99, available=true, keywords=computer, stock=50, name=Laptop, rating=4.5, brand=TechCorp, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, name=Mouse, price=25.5, keywords=peripheral, stock=200, category=accessories, available=true, rating=4, brand=TechCorp, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, category=accessories, price=75, available=false, rating=3.5, keywords=typing, supplier=KeySupply, name=Keyboard, brand=KeyTech, stock=0}
4. [1] TestProduct{id=PROD004, rating=4.8, keywords=display, stock=30, name=Monitor, price=299.99, available=true, brand=ScreenPro, supplier=ScreenSupply, category=electronics}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, price=8.5, available=false, rating=2, keywords=obsolete, brand=OldTech, stock=0, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, brand=AudioMax, stock=75, supplier=AudioSupply, name=Headphones, price=150, available=true, rating=4.6, category=audio, keywords=sound}
7. [1] TestProduct{id=PROD007, name=Webcam, rating=3.8, keywords=video, brand=CamTech, stock=25, category=electronics, price=89.99, available=true, supplier=CamSupply}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, price=25.5, keywords=peripheral, stock=200, category=accessories, available=true, rating=4, brand=TechCorp, supplier=TechSupply}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD003, rating=3.5, keywords=typing, supplier=KeySupply, name=Keyboard, brand=KeyTech, stock=0, category=accessories, price=75, available=false}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD004, category=electronics, rating=4.8, keywords=display, stock=30, name=Monitor, price=299.99, available=true, brand=ScreenPro, supplier=ScreenSupply}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD006, supplier=AudioSupply, name=Headphones, price=150, available=true, rating=4.6, category=audio, keywords=sound, brand=AudioMax, stock=75}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD007, stock=25, category=electronics, price=89.99, available=true, supplier=CamSupply, name=Webcam, rating=3.8, keywords=video, brand=CamTech}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD001, category=electronics, price=999.99, available=true, keywords=computer, stock=50, name=Laptop, rating=4.5, brand=TechCorp, supplier=TechSupply}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}
2. [1] TestPerson{id=P002, tags=senior, name=Bob, age=35, salary=75000, score=9.2, status=active, department=engineering, level=5, active=true}
3. [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie}
4. [1] TestPerson{id=P004, salary=85000, active=true, status=active, department=marketing, level=7, name=Diana, age=45, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, score=8, department=sales, age=30, tags=employee, status=inactive, level=3}
6. [1] TestPerson{id=P006, active=true, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000, status=active, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, department=management, active=true, level=9}
8. [1] TestPerson{id=P008, department=support, level=1, salary=25000, name=Henry, age=18, active=false, score=5.5, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, level=1, score=6.5, status=active, name=X, age=22, salary=28000, active=true, tags=temp, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, age=45, score=7.8, tags=manager, salary=85000, active=true, status=active, department=marketing}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, score=8, department=sales, age=30, tags=employee, status=inactive, level=3, name=Eve, salary=55000, active=false}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, salary=25000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, status=active, level=1, active=true, score=0, tags=test, department=qa}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, department=management, active=true, level=9}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, status=active, department=sales, active=true, level=2, name=Alice, age=25, salary=45000, score=8.5, tags=junior}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, level=5, active=true, tags=senior, name=Bob, age=35, salary=75000, score=9.2, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north, total=1999.98, discount=50}
2. [1] TestOrder{id=O002, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, priority=low, customer_id=P002, total=25.5}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, priority=high, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped}
4. [1] TestOrder{id=O004, total=299.99, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, discount=0, amount=1}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high, discount=100, total=999.99, region=south}
6. [1] TestOrder{id=O006, status=cancelled, region=west, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, product_id=PROD005, amount=2}
7. [1] TestOrder{id=O007, total=600, date=2024-03-01, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50, region=north}
8. [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, discount=10, region=north, product_id=PROD007, status=completed, priority=low, customer_id=P001, amount=1}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, status=refunded, priority=urgent, region=east, product_id=PROD001, total=75000, date=2024-03-15, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O009, priority=low, customer_id=P001, amount=1, total=89.99, date=2024-03-10, discount=10, region=north, product_id=PROD007, status=completed}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north, total=1999.98, discount=50}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, priority=low, customer_id=P002, total=25.5, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, amount=3, priority=high, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, discount=0}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high, discount=100, total=999.99, region=south, customer_id=P002}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, total=600, date=2024-03-01, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50, region=north}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}
2. [1] TestPerson{id=P002, tags=senior, name=Bob, age=35, salary=75000, score=9.2, status=active, department=engineering, level=5, active=true}
3. [1] TestPerson{id=P003, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false}
4. [1] TestPerson{id=P004, salary=85000, active=true, status=active, department=marketing, level=7, name=Diana, age=45, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, score=8, department=sales, age=30, tags=employee, status=inactive, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, salary=-5000, status=active, level=1, active=true, score=0, tags=test, department=qa}
7. [1] TestPerson{id=P007, active=true, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, salary=25000}
9. [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, tags=temp, department=intern, level=1, score=6.5, status=active, name=X, age=22, salary=28000, active=true}
11. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north, total=1999.98, discount=50}
12. [1] TestOrder{id=O002, product_id=PROD002, amount=1, priority=low, customer_id=P002, total=25.5, date=2024-01-20, status=confirmed, discount=0, region=south}
13. [1] TestOrder{id=O003, amount=3, priority=high, discount=15, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped, customer_id=P001}
14. [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, discount=0}
15. [1] TestOrder{id=O005, region=south, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high, discount=100, total=999.99}
16. [1] TestOrder{id=O006, discount=0, customer_id=P005, product_id=PROD005, amount=2, status=cancelled, region=west, total=999.98, date=2024-02-15, priority=low}
17. [1] TestOrder{id=O007, date=2024-03-01, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50, region=north, total=600}
18. [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}
19. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, discount=10, region=north, product_id=PROD007, status=completed, priority=low}
20. [1] TestOrder{id=O010, product_id=PROD001, total=75000, date=2024-03-15, discount=0, customer_id=P006, amount=1, status=refunded, priority=urgent, region=east}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}
2. [1] TestPerson{id=P002, salary=75000, score=9.2, status=active, department=engineering, level=5, active=true, tags=senior, name=Bob, age=35}
3. [1] TestPerson{id=P003, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie}
4. [1] TestPerson{id=P004, salary=85000, active=true, status=active, department=marketing, level=7, name=Diana, age=45, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, salary=55000, active=false, score=8, department=sales, age=30, tags=employee, status=inactive, level=3}
6. [1] TestPerson{id=P006, active=true, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000, status=active, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, department=management, active=true, level=9}
8. [1] TestPerson{id=P008, tags=junior, status=inactive, department=support, level=1, salary=25000, name=Henry, age=18, active=false, score=5.5}
9. [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, tags=temp, department=intern, level=1, score=6.5, status=active, name=X, age=22, salary=28000, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, tags=temp, department=intern, level=1, score=6.5, status=active}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, status=active, department=sales, active=true, level=2, name=Alice, age=25, salary=45000, score=8.5, tags=junior}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, score=9.2, status=active, department=engineering, level=5, active=true, tags=senior}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, status=active, department=marketing, level=7, name=Diana, age=45, score=7.8, tags=manager}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P006, tags=test, department=qa, name=Frank, age=0, salary=-5000, status=active, level=1, active=true, score=0}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, score=10, tags=executive, status=active, department=management, active=true, level=9, name=Grace, age=65}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, score=9.2, status=active, department=engineering, level=5, active=true, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive}
4. [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, salary=85000, active=true, status=active, department=marketing, level=7}
5. [1] TestPerson{id=P005, status=inactive, level=3, name=Eve, salary=55000, active=false, score=8, department=sales, age=30, tags=employee}
6. [1] TestPerson{id=P006, active=true, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000, status=active, level=1}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, department=management, active=true, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, salary=25000}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, score=8.7, tags=senior, age=40, active=true, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, level=1, score=6.5, status=active, name=X, age=22, salary=28000, active=true, tags=temp, department=intern}
11. [1] TestOrder{id=O001, priority=normal, region=north, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, amount=2, date=2024-01-15, status=pending}
12. [1] TestOrder{id=O002, discount=0, region=south, product_id=PROD002, amount=1, priority=low, customer_id=P002, total=25.5, date=2024-01-20, status=confirmed}
13. [1] TestOrder{id=O003, product_id=PROD003, total=225, date=2024-02-01, status=shipped, customer_id=P001, amount=3, priority=high, discount=15, region=north}
14. [1] TestOrder{id=O004, discount=0, amount=1, total=299.99, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05}
15. [1] TestOrder{id=O005, priority=high, discount=100, total=999.99, region=south, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed}
16. [1] TestOrder{id=O006, status=cancelled, region=west, total=999.98, date=2024-02-15, priority=low, discount=0, customer_id=P005, product_id=PROD005, amount=2}
17. [1] TestOrder{id=O007, discount=50, region=north, total=600, date=2024-03-01, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, status=shipped}
18. [1] TestOrder{id=O008, customer_id=P010, total=255, date=2024-03-05, discount=0, product_id=PROD002, amount=10, status=pending, priority=normal, region=south}
19. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, discount=10, region=north, product_id=PROD007, status=completed, priority=low, customer_id=P001, amount=1}
20. [1] TestOrder{id=O010, priority=urgent, region=east, product_id=PROD001, total=75000, date=2024-03-15, discount=0, customer_id=P006, amount=1, status=refunded}

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

1. [1] TestPerson{id=P001, department=sales, active=true, level=2, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active}
2. [1] TestPerson{id=P002, age=35, salary=75000, score=9.2, status=active, department=engineering, level=5, active=true, tags=senior, name=Bob}
3. [1] TestPerson{id=P003, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6, tags=intern, department=hr}
4. [1] TestPerson{id=P004, age=45, score=7.8, tags=manager, salary=85000, active=true, status=active, department=marketing, level=7, name=Diana}
5. [1] TestPerson{id=P005, salary=55000, active=false, score=8, department=sales, age=30, tags=employee, status=inactive, level=3, name=Eve}
6. [1] TestPerson{id=P006, status=active, level=1, active=true, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, department=management, active=true, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, salary=25000}
9. [1] TestPerson{id=P009, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior, age=40, active=true}
10. [1] TestPerson{id=P010, tags=temp, department=intern, level=1, score=6.5, status=active, name=X, age=22, salary=28000, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, tags=senior, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, score=9.2, status=active, department=engineering, level=5, active=true, tags=senior}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, age=45, score=7.8, tags=manager, salary=85000, active=true, status=active, department=marketing, level=7}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, salary=25000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P010, name=X, age=22, salary=28000, active=true, tags=temp, department=intern, level=1, score=6.5, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, tags=junior, status=active, department=sales, active=true, level=2, name=Alice, age=25}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, age=30, tags=employee, status=inactive, level=3, name=Eve, salary=55000, active=false, score=8, department=sales}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, department=management, active=true, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, date=2024-01-15, status=pending, priority=normal, region=north, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, amount=2}
2. [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, priority=low}
3. [1] TestOrder{id=O003, region=north, product_id=PROD003, total=225, date=2024-02-01, status=shipped, customer_id=P001, amount=3, priority=high, discount=15}
4. [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, priority=normal, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, discount=0}
5. [1] TestOrder{id=O005, discount=100, total=999.99, region=south, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed, priority=high}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, status=cancelled, region=west, total=999.98, date=2024-02-15, priority=low, discount=0}
7. [1] TestOrder{id=O007, total=600, date=2024-03-01, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, status=shipped, discount=50, region=north}
8. [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}
9. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, date=2024-03-10, discount=10, region=north, product_id=PROD007, status=completed, priority=low}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, status=refunded, priority=urgent, region=east, product_id=PROD001, total=75000, date=2024-03-15, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O009, priority=low, customer_id=P001, amount=1, total=89.99, date=2024-03-10, discount=10, region=north, product_id=PROD007, status=completed}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, total=25.5, date=2024-01-20, status=confirmed, discount=0, region=south, product_id=PROD002, amount=1, priority=low}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, total=225, date=2024-02-01, status=shipped, customer_id=P001, amount=3, priority=high, discount=15, region=north, product_id=PROD003}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, amount=1, status=refunded, priority=urgent, region=east, product_id=PROD001, total=75000, date=2024-03-15, discount=0}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O001, total=1999.98, discount=50, customer_id=P001, product_id=PROD001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, region=east, customer_id=P004, product_id=PROD004, date=2024-02-05, discount=0, amount=1, total=299.99, status=delivered, priority=normal}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, priority=high, discount=100, total=999.99, region=south, customer_id=P002, product_id=PROD001, amount=1, date=2024-02-10, status=confirmed}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, amount=2, status=cancelled, region=west, total=999.98, date=2024-02-15, priority=low, discount=0}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O007, discount=50, region=north, total=600, date=2024-03-01, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, status=shipped}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O008, product_id=PROD002, amount=10, status=pending, priority=normal, region=south, customer_id=P010, total=255, date=2024-03-05, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, level=2, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales}
2. [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, score=9.2, status=active, department=engineering, level=5, active=true, tags=senior}
3. [1] TestPerson{id=P003, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive, name=Charlie, active=false, score=6}
4. [1] TestPerson{id=P004, salary=85000, active=true, status=active, department=marketing, level=7, name=Diana, age=45, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, age=30, tags=employee, status=inactive, level=3, name=Eve, salary=55000, active=false, score=8, department=sales}
6. [1] TestPerson{id=P006, active=true, score=0, tags=test, department=qa, name=Frank, age=0, salary=-5000, status=active, level=1}
7. [1] TestPerson{id=P007, department=management, active=true, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active}
8. [1] TestPerson{id=P008, name=Henry, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1, salary=25000}
9. [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, status=active, name=X, age=22, salary=28000, active=true, tags=temp, department=intern, level=1, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, department=intern, level=1, score=6.5, status=active, name=X, age=22}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, tags=manager, salary=85000, active=true, status=active, department=marketing, level=7, name=Diana, age=45}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, status=active, level=1, active=true, score=0, tags=test, department=qa, name=Frank, age=0}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, status=active, department=engineering, level=6, name=Ivy, salary=68000, score=8.7, tags=senior}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, active=true, level=2, name=Alice, age=25, salary=45000, score=8.5, tags=junior, status=active, department=sales}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, salary=75000, score=9.2, status=active, department=engineering, level=5, active=true, tags=senior}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, active=false, score=6, tags=intern, department=hr, level=1, age=16, salary=0, status=inactive}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, age=30, tags=employee, status=inactive, level=3, name=Eve, salary=55000, active=false, score=8, department=sales}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, active=true, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, department=management}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, name=Henry, age=18, active=false, score=5.5, tags=junior, status=inactive, department=support, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 17 (89.5%)
- **Tokens générés**: 130
- **Faits traités**: 27
