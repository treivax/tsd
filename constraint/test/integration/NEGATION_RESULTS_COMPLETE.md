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

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, department=engineering, level=5, name=Bob, age=35, tags=senior, status=active}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, score=6, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, tags=manager, status=active, age=45, department=marketing, level=7, name=Diana, salary=85000, active=true, score=7.8}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales, level=3}
6. [1] TestPerson{id=P006, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}
8. [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, status=inactive, name=Henry, age=18, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, tags=senior, status=active, name=Ivy, age=40, active=true, score=8.7, department=engineering, level=6, salary=68000}
10. [1] TestPerson{id=P010, tags=temp, level=1, name=X, status=active, department=intern, age=22, salary=28000, active=true, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, tags=senior, status=active, salary=75000, active=true, score=9.2, department=engineering, level=5}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, tags=junior, department=support, level=1, salary=25000, active=false, score=5.5, status=inactive}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, active=true, score=8.7, department=engineering, level=6, salary=68000, tags=senior, status=active}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P010, name=X, status=active, department=intern, age=22, salary=28000, active=true, score=6.5, tags=temp, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, active=false, tags=intern, score=6, status=inactive, department=hr, level=1, name=Charlie, age=16, salary=0}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, level=7, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales, level=3, name=Eve}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, discount=50, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north}
2. [1] TestOrder{id=O002, region=south, amount=1, total=25.5, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, product_id=PROD002, discount=0}
3. [1] TestOrder{id=O003, total=225, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003, amount=3}
4. [1] TestOrder{id=O004, product_id=PROD004, amount=1, region=east, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, customer_id=P004}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, discount=100, region=south, date=2024-02-10, priority=high}
6. [1] TestOrder{id=O006, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, customer_id=P005, amount=2, status=cancelled}
7. [1] TestOrder{id=O007, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600, priority=urgent, amount=4, date=2024-03-01}
8. [1] TestOrder{id=O008, date=2024-03-05, status=pending, discount=0, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10}
9. [1] TestOrder{id=O009, total=89.99, priority=low, region=north, product_id=PROD007, date=2024-03-10, status=completed, discount=10, customer_id=P001, amount=1}
10. [1] TestOrder{id=O010, customer_id=P006, priority=urgent, region=east, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, region=south, amount=1, total=25.5, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, product_id=PROD002, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, discount=15, region=north, customer_id=P001, product_id=PROD003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, region=east, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, date=2024-02-10, priority=high, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, discount=100, region=south}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, total=600, priority=urgent, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, discount=0, total=255, priority=normal, region=south}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, date=2024-03-10, status=completed, discount=10, customer_id=P001, amount=1, total=89.99, priority=low, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, status=refunded, discount=0, customer_id=P006, priority=urgent, region=east, product_id=PROD001}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north, product_id=PROD001, total=1999.98, discount=50}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25}
2. [1] TestPerson{id=P002, score=9.2, department=engineering, level=5, name=Bob, age=35, tags=senior, status=active, salary=75000, active=true}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, score=6, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, age=45, department=marketing, level=7, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales, level=3}
6. [1] TestPerson{id=P006, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, score=10, tags=executive, department=management, level=9, active=true, status=active, name=Grace, age=65, salary=95000}
8. [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, status=inactive, name=Henry, age=18, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, salary=68000, tags=senior, status=active, name=Ivy, age=40, active=true, score=8.7, department=engineering, level=6}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, tags=temp, level=1, name=X, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales, level=3}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, active=true, status=active, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, salary=68000, tags=senior, status=active, name=Ivy, age=40, active=true, score=8.7, department=engineering, level=6}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active, active=true}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, department=engineering, level=5, name=Bob, age=35, tags=senior, status=active, salary=75000, active=true, score=9.2}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north, product_id=PROD001, total=1999.98, discount=50}
2. [1] TestOrder{id=O002, amount=1, total=25.5, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, product_id=PROD002, discount=0, region=south}
3. [1] TestOrder{id=O003, status=shipped, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003, amount=3, total=225, date=2024-02-01}
4. [1] TestOrder{id=O004, status=delivered, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1, region=east, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, discount=100, region=south, date=2024-02-10, priority=high}
6. [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, priority=low}
7. [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600, priority=urgent}
8. [1] TestOrder{id=O008, date=2024-03-05, status=pending, discount=0, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10}
9. [1] TestOrder{id=O009, product_id=PROD007, date=2024-03-10, status=completed, discount=10, customer_id=P001, amount=1, total=89.99, priority=low, region=north}
10. [1] TestOrder{id=O010, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, discount=0, customer_id=P006, priority=urgent, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north, product_id=PROD001, total=1999.98, discount=50}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O003, total=225, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003, amount=3}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1, region=east}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, discount=100, region=south, date=2024-02-10, priority=high}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, customer_id=P005, amount=2, status=cancelled}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O002, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, product_id=PROD002, discount=0, region=south, amount=1, total=25.5}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600, priority=urgent}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, discount=0, total=255, priority=normal, region=south}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, priority=low, region=north, product_id=PROD007, date=2024-03-10, status=completed, discount=10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, category=electronics, price=999.99, keywords=computer, brand=TechCorp, supplier=TechSupply, available=true, rating=4.5, stock=50, name=Laptop}
2. [1] TestProduct{id=PROD002, keywords=peripheral, stock=200, category=accessories, brand=TechCorp, supplier=TechSupply, name=Mouse, price=25.5, available=true, rating=4}
3. [1] TestProduct{id=PROD003, supplier=KeySupply, name=Keyboard, category=accessories, price=75, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0}
4. [1] TestProduct{id=PROD004, name=Monitor, category=electronics, price=299.99, rating=4.8, keywords=display, supplier=ScreenSupply, available=true, brand=ScreenPro, stock=30}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, price=8.5, available=false, stock=0, supplier=OldSupply, rating=2, keywords=obsolete, brand=OldTech}
6. [1] TestProduct{id=PROD006, category=audio, price=150, keywords=sound, name=Headphones, available=true, rating=4.6, brand=AudioMax, stock=75, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, available=true, rating=3.8, keywords=video, brand=CamTech, stock=25, supplier=CamSupply, name=Webcam, category=electronics, price=89.99}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, stock=50, name=Laptop, category=electronics, price=999.99, keywords=computer, brand=TechCorp, supplier=TechSupply, available=true, rating=4.5}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, price=25.5, available=true, rating=4, keywords=peripheral, stock=200, category=accessories, brand=TechCorp, supplier=TechSupply, name=Mouse}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, brand=KeyTech, stock=0, supplier=KeySupply, name=Keyboard, category=accessories, price=75, available=false, rating=3.5, keywords=typing}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, name=Monitor, category=electronics, price=299.99, rating=4.8, keywords=display, supplier=ScreenSupply, available=true, brand=ScreenPro, stock=30}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, category=audio, price=150, keywords=sound, name=Headphones, available=true, rating=4.6, brand=AudioMax, stock=75, supplier=AudioSupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, stock=25, supplier=CamSupply, name=Webcam, category=electronics, price=89.99, available=true, rating=3.8, keywords=video, brand=CamTech}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active}
2. [1] TestPerson{id=P002, active=true, score=9.2, department=engineering, level=5, name=Bob, age=35, tags=senior, status=active, salary=75000}
3. [1] TestPerson{id=P003, score=6, status=inactive, department=hr, level=1, name=Charlie, age=16, salary=0, active=false, tags=intern}
4. [1] TestPerson{id=P004, department=marketing, level=7, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales, level=3}
6. [1] TestPerson{id=P006, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}
8. [1] TestPerson{id=P008, name=Henry, age=18, tags=junior, department=support, level=1, salary=25000, active=false, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, score=8.7, department=engineering, level=6, salary=68000, tags=senior, status=active, name=Ivy, age=40, active=true}
10. [1] TestPerson{id=P010, score=6.5, tags=temp, level=1, name=X, status=active, department=intern, age=22, salary=28000, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, score=6, status=inactive, department=hr, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, discount=50, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north}
2. [1] TestOrder{id=O002, status=confirmed, priority=low, customer_id=P002, product_id=PROD002, discount=0, region=south, amount=1, total=25.5, date=2024-01-20}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, discount=15, region=north}
4. [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1, region=east}
5. [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, status=confirmed, discount=100, region=south, date=2024-02-10, priority=high, customer_id=P002}
6. [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, priority=low}
7. [1] TestOrder{id=O007, priority=urgent, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, discount=0, total=255, priority=normal, region=south}
9. [1] TestOrder{id=O009, discount=10, customer_id=P001, amount=1, total=89.99, priority=low, region=north, product_id=PROD007, date=2024-03-10, status=completed}
10. [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, discount=0, customer_id=P006, priority=urgent, region=east, product_id=PROD001, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1, region=east}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, total=999.99, status=confirmed, discount=100, region=south, date=2024-02-10, priority=high, customer_id=P002, product_id=PROD001}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O006, date=2024-02-15, priority=low, customer_id=P005, amount=2, status=cancelled, discount=0, region=west, product_id=PROD005, total=999.98}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600, priority=urgent}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O008, status=pending, discount=0, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, region=north, product_id=PROD007, date=2024-03-10, status=completed, discount=10, customer_id=P001, amount=1, total=89.99, priority=low}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north, product_id=PROD001, total=1999.98, discount=50}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O002, total=25.5, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, product_id=PROD002, discount=0, region=south, amount=1}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, amount=1, total=75000, date=2024-03-15, status=refunded, discount=0, customer_id=P006, priority=urgent, region=east, product_id=PROD001}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2}
2. [1] TestPerson{id=P002, name=Bob, age=35, tags=senior, status=active, salary=75000, active=true, score=9.2, department=engineering, level=5}
3. [1] TestPerson{id=P003, score=6, status=inactive, department=hr, level=1, name=Charlie, age=16, salary=0, active=false, tags=intern}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7}
5. [1] TestPerson{id=P005, tags=employee, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, active=false, score=8}
6. [1] TestPerson{id=P006, level=1, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}
8. [1] TestPerson{id=P008, name=Henry, age=18, tags=junior, department=support, level=1, salary=25000, active=false, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, level=6, salary=68000, tags=senior, status=active, name=Ivy, age=40, active=true, score=8.7, department=engineering}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, tags=temp, level=1, name=X, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, active=false, tags=intern, score=6, status=inactive, department=hr, level=1, name=Charlie, age=16, salary=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales, level=3, name=Eve, age=30}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, tags=junior, department=support, level=1, salary=25000, active=false, score=5.5, status=inactive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, discount=50, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north}
2. [1] TestOrder{id=O002, region=south, amount=1, total=25.5, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, product_id=PROD002, discount=0}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, discount=15, region=north}
4. [1] TestOrder{id=O004, region=east, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, discount=100, region=south, date=2024-02-10, priority=high}
6. [1] TestOrder{id=O006, priority=low, customer_id=P005, amount=2, status=cancelled, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15}
7. [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600, priority=urgent}
8. [1] TestOrder{id=O008, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, discount=0, total=255}
9. [1] TestOrder{id=O009, region=north, product_id=PROD007, date=2024-03-10, status=completed, discount=10, customer_id=P001, amount=1, total=89.99, priority=low}
10. [1] TestOrder{id=O010, customer_id=P006, priority=urgent, region=east, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, region=east, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, total=999.99, status=confirmed, discount=100, region=south, date=2024-02-10, priority=high, customer_id=P002, product_id=PROD001}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O006, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, customer_id=P005, amount=2, status=cancelled}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O007, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600, priority=urgent, amount=4, date=2024-03-01}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O009, product_id=PROD007, date=2024-03-10, status=completed, discount=10, customer_id=P001, amount=1, total=89.99, priority=low, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O010, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, discount=0, customer_id=P006, priority=urgent, region=east}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, discount=0, region=south, amount=1, total=25.5, date=2024-01-20, status=confirmed, priority=low}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O003, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003, amount=3, total=225, date=2024-02-01, status=shipped}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice}
2. [1] TestPerson{id=P002, status=active, salary=75000, active=true, score=9.2, department=engineering, level=5, name=Bob, age=35, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, score=6, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7, name=Diana}
5. [1] TestPerson{id=P005, level=3, name=Eve, age=30, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales}
6. [1] TestPerson{id=P006, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, active=true, status=active, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, status=inactive, name=Henry, age=18, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, score=8.7, department=engineering, level=6, salary=68000, tags=senior, status=active, name=Ivy, age=40, active=true}
10. [1] TestPerson{id=P010, status=active, department=intern, age=22, salary=28000, active=true, score=6.5, tags=temp, level=1, name=X}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, tags=senior, status=active, salary=75000, active=true, score=9.2, department=engineering, level=5}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, active=false, tags=intern, score=6, status=inactive, department=hr, level=1, name=Charlie}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, age=45, department=marketing, level=7, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, level=3, name=Eve, age=30, salary=55000, active=false, score=8, tags=employee, status=inactive}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, tags=junior, department=support, level=1, salary=25000, active=false, score=5.5, status=inactive}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, status=active, name=Ivy, age=40, active=true, score=8.7, department=engineering, level=6, salary=68000, tags=senior}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, available=true, rating=4.5, stock=50, name=Laptop, category=electronics, price=999.99, keywords=computer, brand=TechCorp, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, name=Mouse, price=25.5, available=true, rating=4, keywords=peripheral, stock=200, category=accessories, brand=TechCorp, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, price=75, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, rating=4.8, keywords=display, supplier=ScreenSupply, available=true, brand=ScreenPro, stock=30, name=Monitor, category=electronics, price=299.99}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, price=8.5, available=false, stock=0, supplier=OldSupply, rating=2, keywords=obsolete, brand=OldTech}
6. [1] TestProduct{id=PROD006, name=Headphones, available=true, rating=4.6, brand=AudioMax, stock=75, supplier=AudioSupply, category=audio, price=150, keywords=sound}
7. [1] TestProduct{id=PROD007, brand=CamTech, stock=25, supplier=CamSupply, name=Webcam, category=electronics, price=89.99, available=true, rating=3.8, keywords=video}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, name=Laptop, category=electronics, price=999.99, keywords=computer, brand=TechCorp, supplier=TechSupply, available=true, rating=4.5, stock=50}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, available=true, rating=4, keywords=peripheral, stock=200, category=accessories, brand=TechCorp, supplier=TechSupply, name=Mouse, price=25.5}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, category=accessories, price=75, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply, name=Keyboard}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, keywords=display, supplier=ScreenSupply, available=true, brand=ScreenPro, stock=30, name=Monitor, category=electronics, price=299.99, rating=4.8}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, price=150, keywords=sound, name=Headphones, available=true, rating=4.6, brand=AudioMax, stock=75, supplier=AudioSupply, category=audio}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, supplier=CamSupply, name=Webcam, category=electronics, price=89.99, available=true, rating=3.8, keywords=video, brand=CamTech, stock=25}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, level=2, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior}
2. [1] TestPerson{id=P002, name=Bob, age=35, tags=senior, status=active, salary=75000, active=true, score=9.2, department=engineering, level=5}
3. [1] TestPerson{id=P003, status=inactive, department=hr, level=1, name=Charlie, age=16, salary=0, active=false, tags=intern, score=6}
4. [1] TestPerson{id=P004, status=active, age=45, department=marketing, level=7, name=Diana, salary=85000, active=true, score=7.8, tags=manager}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales, level=3}
6. [1] TestPerson{id=P006, status=active, department=qa, level=1, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}
8. [1] TestPerson{id=P008, name=Henry, age=18, tags=junior, department=support, level=1, salary=25000, active=false, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, department=engineering, level=6, salary=68000, tags=senior, status=active, name=Ivy, age=40, active=true, score=8.7}
10. [1] TestPerson{id=P010, score=6.5, tags=temp, level=1, name=X, status=active, department=intern, age=22, salary=28000, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, status=active, salary=75000, active=true, score=9.2, department=engineering, level=5, name=Bob, age=35}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, level=1, name=Charlie, age=16, salary=0, active=false, tags=intern, score=6, status=inactive, department=hr}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales, level=3}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, status=active, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, active=true, status=active, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, tags=junior, department=support, level=1, salary=25000, active=false, score=5.5, status=inactive}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, age=40, active=true, score=8.7, department=engineering, level=6, salary=68000, tags=senior, status=active, name=Ivy}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north, product_id=PROD001, total=1999.98, discount=50}
2. [1] TestOrder{id=O002, amount=1, total=25.5, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, product_id=PROD002, discount=0, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, discount=15, region=north}
4. [1] TestOrder{id=O004, status=delivered, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1, region=east, total=299.99, date=2024-02-05}
5. [1] TestOrder{id=O005, amount=1, total=999.99, status=confirmed, discount=100, region=south, date=2024-02-10, priority=high, customer_id=P002, product_id=PROD001}
6. [1] TestOrder{id=O006, region=west, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, customer_id=P005, amount=2, status=cancelled, discount=0}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, priority=urgent, amount=4, date=2024-03-01, status=shipped, discount=50, region=north}
8. [1] TestOrder{id=O008, discount=0, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending}
9. [1] TestOrder{id=O009, product_id=PROD007, date=2024-03-10, status=completed, discount=10, customer_id=P001, amount=1, total=89.99, priority=low, region=north}
10. [1] TestOrder{id=O010, region=east, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, discount=0, customer_id=P006, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O001, product_id=PROD001, total=1999.98, discount=50, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O002, discount=0, region=south, amount=1, total=25.5, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, product_id=PROD002}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, discount=15, region=north}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, region=east, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O005, product_id=PROD001, amount=1, total=999.99, status=confirmed, discount=100, region=south, date=2024-02-10, priority=high, customer_id=P002}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, priority=urgent, amount=4, date=2024-03-01, status=shipped, discount=50, region=north}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O008, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, total=89.99, priority=low, region=north, product_id=PROD007, date=2024-03-10, status=completed, discount=10, customer_id=P001, amount=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active}
2. [1] TestPerson{id=P002, name=Bob, age=35, tags=senior, status=active, salary=75000, active=true, score=9.2, department=engineering, level=5}
3. [1] TestPerson{id=P003, active=false, tags=intern, score=6, status=inactive, department=hr, level=1, name=Charlie, age=16, salary=0}
4. [1] TestPerson{id=P004, age=45, department=marketing, level=7, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales, level=3}
6. [1] TestPerson{id=P006, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}
8. [1] TestPerson{id=P008, salary=25000, active=false, score=5.5, status=inactive, name=Henry, age=18, tags=junior, department=support, level=1}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, score=8.7, department=engineering, level=6, salary=68000, tags=senior, status=active}
10. [1] TestPerson{id=P010, name=X, status=active, department=intern, age=22, salary=28000, active=true, score=6.5, tags=temp, level=1}
11. [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north, product_id=PROD001, total=1999.98, discount=50}
12. [1] TestOrder{id=O002, amount=1, total=25.5, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, product_id=PROD002, discount=0, region=south}
13. [1] TestOrder{id=O003, product_id=PROD003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, customer_id=P001}
14. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, region=east, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, discount=100, region=south, date=2024-02-10, priority=high}
16. [1] TestOrder{id=O006, amount=2, status=cancelled, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, customer_id=P005}
17. [1] TestOrder{id=O007, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600, priority=urgent, amount=4, date=2024-03-01, status=shipped}
18. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, discount=0, total=255, priority=normal, region=south}
19. [1] TestOrder{id=O009, region=north, product_id=PROD007, date=2024-03-10, status=completed, discount=10, customer_id=P001, amount=1, total=89.99, priority=low}
20. [1] TestOrder{id=O010, total=75000, date=2024-03-15, status=refunded, discount=0, customer_id=P006, priority=urgent, region=east, product_id=PROD001, amount=1}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, amount=2, status=cancelled, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, priority=low}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, priority=urgent, amount=4, date=2024-03-01, status=shipped, discount=50, region=north}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, status=active, department=intern, age=22, salary=28000, active=true, score=6.5, tags=temp, level=1, name=X}
   - Fait 2: [1] TestOrder{id=O008, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, discount=0, total=255}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active}
   - Fait 2: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, priority=low, region=north, product_id=PROD007, date=2024-03-10, status=completed, discount=10}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O010, customer_id=P006, priority=urgent, region=east, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P001, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice}
   - Fait 2: [1] TestOrder{id=O001, region=north, product_id=PROD001, total=1999.98, discount=50, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, discount=15, region=north}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7}
   - Fait 2: [1] TestOrder{id=O004, discount=0, customer_id=P004, product_id=PROD004, amount=1, region=east, total=299.99, date=2024-02-05, status=delivered, priority=normal}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, department=engineering, level=5, name=Bob, age=35, tags=senior, status=active, salary=75000}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, discount=100, region=south, date=2024-02-10, priority=high}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, department=engineering, level=5, name=Bob, age=35, tags=senior, status=active}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, discount=0, region=south, amount=1, total=25.5, date=2024-01-20, status=confirmed, priority=low, customer_id=P002}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/20 (50.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice}
2. [1] TestPerson{id=P002, name=Bob, age=35, tags=senior, status=active, salary=75000, active=true, score=9.2, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, score=6, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7, name=Diana, salary=85000}
5. [1] TestPerson{id=P005, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, active=false, score=8, tags=employee}
6. [1] TestPerson{id=P006, department=qa, level=1, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active}
7. [1] TestPerson{id=P007, active=true, status=active, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, tags=junior, department=support, level=1, salary=25000, active=false, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, age=40, active=true, score=8.7, department=engineering, level=6, salary=68000, tags=senior, status=active}
10. [1] TestPerson{id=P010, department=intern, age=22, salary=28000, active=true, score=6.5, tags=temp, level=1, name=X, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, department=engineering, level=5, name=Bob, age=35, tags=senior, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P009, status=active, name=Ivy, age=40, active=true, score=8.7, department=engineering, level=6, salary=68000, tags=senior}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, tags=temp, level=1, name=X, status=active, department=intern, age=22, salary=28000, active=true}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2}
2. [1] TestPerson{id=P002, status=active, salary=75000, active=true, score=9.2, department=engineering, level=5, name=Bob, age=35, tags=senior}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, score=6, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7}
5. [1] TestPerson{id=P005, score=8, tags=employee, status=inactive, department=sales, level=3, name=Eve, age=30, salary=55000, active=false}
6. [1] TestPerson{id=P006, status=active, department=qa, level=1, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0}
7. [1] TestPerson{id=P007, tags=executive, department=management, level=9, active=true, status=active, name=Grace, age=65, salary=95000, score=10}
8. [1] TestPerson{id=P008, name=Henry, age=18, tags=junior, department=support, level=1, salary=25000, active=false, score=5.5, status=inactive}
9. [1] TestPerson{id=P009, department=engineering, level=6, salary=68000, tags=senior, status=active, name=Ivy, age=40, active=true, score=8.7}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, tags=temp, level=1, name=X, status=active, department=intern}
11. [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north, product_id=PROD001, total=1999.98, discount=50}
12. [1] TestOrder{id=O002, amount=1, total=25.5, date=2024-01-20, status=confirmed, priority=low, customer_id=P002, product_id=PROD002, discount=0, region=south}
13. [1] TestOrder{id=O003, date=2024-02-01, status=shipped, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003, amount=3, total=225}
14. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, amount=1, region=east, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0}
15. [1] TestOrder{id=O005, total=999.99, status=confirmed, discount=100, region=south, date=2024-02-10, priority=high, customer_id=P002, product_id=PROD001, amount=1}
16. [1] TestOrder{id=O006, date=2024-02-15, priority=low, customer_id=P005, amount=2, status=cancelled, discount=0, region=west, product_id=PROD005, total=999.98}
17. [1] TestOrder{id=O007, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600, priority=urgent, amount=4, date=2024-03-01, status=shipped}
18. [1] TestOrder{id=O008, status=pending, discount=0, total=255, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05}
19. [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, priority=low, region=north, product_id=PROD007, date=2024-03-10, status=completed, discount=10}
20. [1] TestOrder{id=O010, date=2024-03-15, status=refunded, discount=0, customer_id=P006, priority=urgent, region=east, product_id=PROD001, amount=1, total=75000}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales, level=3}
   - Fait 2: [1] TestOrder{id=O006, region=west, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, customer_id=P005, amount=2, status=cancelled, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, status=active, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true}
   - Fait 2: [1] TestOrder{id=O007, product_id=PROD006, total=600, priority=urgent, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, name=X, status=active, department=intern, age=22, salary=28000, active=true, score=6.5, tags=temp, level=1}
   - Fait 2: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, discount=0, total=255, priority=normal, region=south}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, status=active, active=true, score=8.5, tags=junior, department=sales, level=2, name=Alice, age=25, salary=45000}
   - Fait 2: [1] TestOrder{id=O001, priority=normal, region=north, product_id=PROD001, total=1999.98, discount=50, customer_id=P001, amount=2, date=2024-01-15, status=pending}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, discount=15, region=north}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P004, department=marketing, level=7, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active, age=45}
   - Fait 2: [1] TestOrder{id=O004, status=delivered, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1, region=east, total=299.99, date=2024-02-05}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, tags=senior, status=active, salary=75000, active=true, score=9.2, department=engineering, level=5}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, discount=100, region=south, date=2024-02-10, priority=high}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2}
   - Fait 2: [1] TestOrder{id=O009, discount=10, customer_id=P001, amount=1, total=89.99, priority=low, region=north, product_id=PROD007, date=2024-03-10, status=completed}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P006, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{id=O010, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, discount=0, customer_id=P006, priority=urgent, region=east}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P002, status=active, salary=75000, active=true, score=9.2, department=engineering, level=5, name=Bob, age=35, tags=senior}
   - Fait 2: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, discount=0, region=south, amount=1, total=25.5, date=2024-01-20, status=confirmed, priority=low}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/20 (50.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales}
2. [1] TestPerson{id=P002, active=true, score=9.2, department=engineering, level=5, name=Bob, age=35, tags=senior, status=active, salary=75000}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, tags=intern, score=6, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, age=45, department=marketing, level=7, name=Diana, salary=85000, active=true, score=7.8, tags=manager, status=active}
5. [1] TestPerson{id=P005, age=30, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales, level=3, name=Eve}
6. [1] TestPerson{id=P006, department=qa, level=1, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}
8. [1] TestPerson{id=P008, tags=junior, department=support, level=1, salary=25000, active=false, score=5.5, status=inactive, name=Henry, age=18}
9. [1] TestPerson{id=P009, salary=68000, tags=senior, status=active, name=Ivy, age=40, active=true, score=8.7, department=engineering, level=6}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, tags=temp, level=1, name=X, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, tags=temp, level=1, name=X, status=active, department=intern}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, department=engineering, level=5, name=Bob, age=35, tags=senior, status=active, salary=75000, active=true}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, status=active, age=45, department=marketing, level=7, name=Diana, salary=85000, active=true, score=7.8, tags=manager}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales, level=3}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P009, department=engineering, level=6, salary=68000, tags=senior, status=active, name=Ivy, age=40, active=true, score=8.7}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P003, level=1, name=Charlie, age=16, salary=0, active=false, tags=intern, score=6, status=inactive, department=hr}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P008, score=5.5, status=inactive, name=Henry, age=18, tags=junior, department=support, level=1, salary=25000, active=false}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, region=north, product_id=PROD001, total=1999.98, discount=50, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, discount=0, region=south, amount=1, total=25.5, date=2024-01-20, status=confirmed, priority=low}
3. [1] TestOrder{id=O003, discount=15, region=north, customer_id=P001, product_id=PROD003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high}
4. [1] TestOrder{id=O004, amount=1, region=east, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, customer_id=P004, product_id=PROD004}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, discount=100, region=south, date=2024-02-10, priority=high}
6. [1] TestOrder{id=O006, status=cancelled, discount=0, region=west, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, customer_id=P005, amount=2}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, total=600, priority=urgent, amount=4, date=2024-03-01, status=shipped, discount=50, region=north}
8. [1] TestOrder{id=O008, priority=normal, region=south, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, discount=0, total=255}
9. [1] TestOrder{id=O009, product_id=PROD007, date=2024-03-10, status=completed, discount=10, customer_id=P001, amount=1, total=89.99, priority=low, region=north}
10. [1] TestOrder{id=O010, priority=urgent, region=east, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, discount=0, customer_id=P006}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, date=2024-03-05, status=pending, discount=0, total=255, priority=normal, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O010, customer_id=P006, priority=urgent, region=east, product_id=PROD001, amount=1, total=75000, date=2024-03-15, status=refunded, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, amount=2, date=2024-01-15, status=pending, priority=normal, region=north, product_id=PROD001, total=1999.98, discount=50}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, discount=0, region=south, amount=1, total=25.5, date=2024-01-20, status=confirmed, priority=low}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, amount=4, date=2024-03-01, status=shipped, discount=50, region=north, customer_id=P007, product_id=PROD006, total=600, priority=urgent}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, total=89.99, priority=low, region=north, product_id=PROD007, date=2024-03-10, status=completed, discount=10}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, discount=15, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O004, region=east, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, customer_id=P004, product_id=PROD004, amount=1}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O005, date=2024-02-10, priority=high, customer_id=P002, product_id=PROD001, amount=1, total=999.99, status=confirmed, discount=100, region=south}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O006, product_id=PROD005, total=999.98, date=2024-02-15, priority=low, customer_id=P005, amount=2, status=cancelled, discount=0, region=west}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, department=engineering, level=5, name=Bob, age=35, tags=senior, status=active}
3. [1] TestPerson{id=P003, salary=0, active=false, tags=intern, score=6, status=inactive, department=hr, level=1, name=Charlie, age=16}
4. [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, age=45, department=marketing, level=7, name=Diana, salary=85000, active=true}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales, level=3}
6. [1] TestPerson{id=P006, salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
7. [1] TestPerson{id=P007, tags=executive, department=management, level=9, active=true, status=active, name=Grace, age=65, salary=95000, score=10}
8. [1] TestPerson{id=P008, level=1, salary=25000, active=false, score=5.5, status=inactive, name=Henry, age=18, tags=junior, department=support}
9. [1] TestPerson{id=P009, tags=senior, status=active, name=Ivy, age=40, active=true, score=8.7, department=engineering, level=6, salary=68000}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, tags=temp, level=1, name=X, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, salary=45000, status=active, active=true, score=8.5, tags=junior, department=sales, level=2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, score=6, status=inactive, department=hr, level=1, name=Charlie, age=16, salary=0, active=false, tags=intern}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, status=active, age=45, department=marketing, level=7, name=Diana, salary=85000, active=true, score=7.8, tags=manager}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, active=false, score=8, tags=employee, status=inactive, department=sales, level=3, name=Eve, age=30}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, department=engineering, level=5, name=Bob, age=35, tags=senior, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P006, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test, salary=-5000}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, score=10, tags=executive, department=management, level=9, active=true, status=active}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P008, age=18, tags=junior, department=support, level=1, salary=25000, active=false, score=5.5, status=inactive, name=Henry}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, age=40, active=true, score=8.7, department=engineering, level=6, salary=68000, tags=senior, status=active}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P010, age=22, salary=28000, active=true, score=6.5, tags=temp, level=1, name=X, status=active, department=intern}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 150
- **Faits traités**: 27
