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

1. [1] TestPerson{id=P001, score=8.5, status=active, department=sales, age=25, active=true, tags=junior, level=2, name=Alice, salary=45000}
2. [1] TestPerson{id=P002, name=Bob, age=35, status=active, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, salary=0, active=false, score=6, department=hr, tags=intern, status=inactive, level=1, name=Charlie, age=16}
4. [1] TestPerson{id=P004, age=45, active=true, tags=manager, status=active, level=7, name=Diana, salary=85000, score=7.8, department=marketing}
5. [1] TestPerson{id=P005, age=30, salary=55000, active=false, tags=employee, status=inactive, department=sales, level=3, name=Eve, score=8}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, age=0, salary=-5000, tags=test, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, level=1, active=false, score=5.5, department=support, name=Henry, age=18, salary=25000, tags=junior, status=inactive}
9. [1] TestPerson{id=P009, age=40, tags=senior, status=active, department=engineering, name=Ivy, salary=68000, active=true, score=8.7, level=6}
10. [1] TestPerson{id=P010, level=1, active=true, tags=temp, status=active, name=X, age=22, salary=28000, score=6.5, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, department=marketing, age=45, active=true, tags=manager, status=active, level=7, name=Diana, salary=85000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, level=1, active=false, score=5.5, department=support}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, level=1, active=true, tags=temp, status=active, name=X, age=22, salary=28000, score=6.5, department=intern}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, age=25, active=true, tags=junior, level=2, name=Alice, salary=45000, score=8.5, status=active, department=sales}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, name=Bob, age=35, status=active, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, department=hr, tags=intern, status=inactive, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, active=false, tags=employee, status=inactive, department=sales, level=3, name=Eve, score=8, age=30, salary=55000}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, age=40, tags=senior, status=active, department=engineering, name=Ivy, salary=68000, active=true, score=8.7, level=6}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, region=north, product_id=PROD001, amount=2, priority=normal, discount=50}
2. [1] TestOrder{id=O002, discount=0, customer_id=P002, total=25.5, status=confirmed, priority=low, region=south, product_id=PROD002, amount=1, date=2024-01-20}
3. [1] TestOrder{id=O003, discount=15, region=north, customer_id=P001, product_id=PROD003, date=2024-02-01, amount=3, total=225, status=shipped, priority=high}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, status=confirmed, region=south}
6. [1] TestOrder{id=O006, priority=low, discount=0, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, region=west, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, region=north, customer_id=P007, total=600, date=2024-03-01, status=shipped, priority=urgent, product_id=PROD006, amount=4, discount=50}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, amount=10, status=pending, priority=normal, discount=0, region=south}
9. [1] TestOrder{id=O009, region=north, product_id=PROD007, total=89.99, date=2024-03-10, discount=10, customer_id=P001, amount=1, status=completed, priority=low}
10. [1] TestOrder{id=O010, amount=1, total=75000, status=refunded, discount=0, region=east, customer_id=P006, product_id=PROD001, date=2024-03-15, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O008, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, amount=10, status=pending, priority=normal}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O001, region=north, product_id=PROD001, amount=2, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, status=pending}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O002, region=south, product_id=PROD002, amount=1, date=2024-01-20, discount=0, customer_id=P002, total=25.5, status=confirmed, priority=low}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003, date=2024-02-01}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, status=confirmed, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, total=600, date=2024-03-01, status=shipped, priority=urgent, product_id=PROD006, amount=4, discount=50, region=north, customer_id=P007}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, status=completed, priority=low, region=north, product_id=PROD007, total=89.99, date=2024-03-10, discount=10}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O010, discount=0, region=east, customer_id=P006, product_id=PROD001, date=2024-03-15, priority=urgent, amount=1, total=75000, status=refunded}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, status=active, department=sales, age=25, active=true, tags=junior, level=2}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, status=active, salary=75000, active=true}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, department=hr, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, score=7.8, department=marketing, age=45, active=true, tags=manager, status=active, level=7, name=Diana, salary=85000}
5. [1] TestPerson{id=P005, name=Eve, score=8, age=30, salary=55000, active=false, tags=employee, status=inactive, department=sales, level=3}
6. [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, department=qa, level=1, name=Frank, active=true, score=0}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, level=1, active=false, score=5.5, department=support}
9. [1] TestPerson{id=P009, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, department=engineering, name=Ivy}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, department=intern, level=1, active=true, tags=temp, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, department=engineering}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, tags=junior, level=2, name=Alice, salary=45000, score=8.5, status=active, department=sales, age=25, active=true}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P002, age=35, status=active, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, name=Bob}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, department=marketing, age=45, active=true, tags=manager, status=active, level=7, name=Diana, salary=85000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, active=false, tags=employee, status=inactive, department=sales, level=3, name=Eve, score=8}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, status=pending, region=north, product_id=PROD001, amount=2, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15}
2. [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, priority=low, region=south, product_id=PROD002, amount=1, date=2024-01-20, discount=0}
3. [1] TestOrder{id=O003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003, date=2024-02-01}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}
5. [1] TestOrder{id=O005, total=999.99, date=2024-02-10, priority=high, discount=100, status=confirmed, region=south, customer_id=P002, product_id=PROD001, amount=1}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, region=west, amount=2, date=2024-02-15, priority=low, discount=0}
7. [1] TestOrder{id=O007, status=shipped, priority=urgent, product_id=PROD006, amount=4, discount=50, region=north, customer_id=P007, total=600, date=2024-03-01}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, amount=10, status=pending, priority=normal, discount=0, region=south}
9. [1] TestOrder{id=O009, total=89.99, date=2024-03-10, discount=10, customer_id=P001, amount=1, status=completed, priority=low, region=north, product_id=PROD007}
10. [1] TestOrder{id=O010, date=2024-03-15, priority=urgent, amount=1, total=75000, status=refunded, discount=0, region=east, customer_id=P006, product_id=PROD001}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, status=completed, priority=low, region=north, product_id=PROD007, total=89.99, date=2024-03-10, discount=10}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O006, status=cancelled, region=west, amount=2, date=2024-02-15, priority=low, discount=0, customer_id=P005, product_id=PROD005, total=999.98}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, status=pending, region=north, product_id=PROD001, amount=2, priority=normal, discount=50, customer_id=P001}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, priority=low, region=south, product_id=PROD002, amount=1, date=2024-01-20, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003, date=2024-02-01}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, total=999.99, date=2024-02-10, priority=high, discount=100, status=confirmed, region=south, customer_id=P002, product_id=PROD001, amount=1}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, priority=urgent, product_id=PROD006, amount=4, discount=50, region=north}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O008, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, amount=10, status=pending}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, category=electronics, price=999.99, available=true, keywords=computer, stock=50, name=Laptop, rating=4.5, brand=TechCorp, supplier=TechSupply}
2. [1] TestProduct{id=PROD002, keywords=peripheral, brand=TechCorp, price=25.5, rating=4, stock=200, supplier=TechSupply, name=Mouse, category=accessories, available=true}
3. [1] TestProduct{id=PROD003, stock=0, supplier=KeySupply, name=Keyboard, available=false, rating=3.5, keywords=typing, category=accessories, price=75, brand=KeyTech}
4. [1] TestProduct{id=PROD004, brand=ScreenPro, supplier=ScreenSupply, category=electronics, available=true, rating=4.8, stock=30, name=Monitor, price=299.99, keywords=display}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, price=8.5, available=false, rating=2, keywords=obsolete, stock=0, category=accessories, brand=OldTech, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, stock=75, name=Headphones, price=150, keywords=sound, supplier=AudioSupply, category=audio, available=true, rating=4.6, brand=AudioMax}
7. [1] TestProduct{id=PROD007, price=89.99, rating=3.8, keywords=video, brand=CamTech, name=Webcam, category=electronics, available=true, stock=25, supplier=CamSupply}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD003, price=75, brand=KeyTech, stock=0, supplier=KeySupply, name=Keyboard, available=false, rating=3.5, keywords=typing, category=accessories}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD004, rating=4.8, stock=30, name=Monitor, price=299.99, keywords=display, brand=ScreenPro, supplier=ScreenSupply, category=electronics, available=true}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD006, category=audio, available=true, rating=4.6, brand=AudioMax, stock=75, name=Headphones, price=150, keywords=sound, supplier=AudioSupply}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD007, price=89.99, rating=3.8, keywords=video, brand=CamTech, name=Webcam, category=electronics, available=true, stock=25, supplier=CamSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD001, price=999.99, available=true, keywords=computer, stock=50, name=Laptop, rating=4.5, brand=TechCorp, supplier=TechSupply, category=electronics}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, category=accessories, available=true, keywords=peripheral, brand=TechCorp, price=25.5, rating=4, stock=200, supplier=TechSupply}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, status=active, department=sales, age=25, active=true, tags=junior, level=2, name=Alice, salary=45000, score=8.5}
2. [1] TestPerson{id=P002, age=35, status=active, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, name=Bob}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, department=hr, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, salary=85000, score=7.8, department=marketing, age=45, active=true, tags=manager, status=active, level=7, name=Diana}
5. [1] TestPerson{id=P005, age=30, salary=55000, active=false, tags=employee, status=inactive, department=sales, level=3, name=Eve, score=8}
6. [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, department=qa, level=1, name=Frank, active=true, score=0}
7. [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, level=1, active=false, score=5.5, department=support}
9. [1] TestPerson{id=P009, level=6, age=40, tags=senior, status=active, department=engineering, name=Ivy, salary=68000, active=true, score=8.7}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, department=intern, level=1, active=true, tags=temp, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, department=hr, tags=intern, status=inactive, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, department=qa, level=1, name=Frank, active=true, score=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, region=north, product_id=PROD001, amount=2, priority=normal}
2. [1] TestOrder{id=O002, discount=0, customer_id=P002, total=25.5, status=confirmed, priority=low, region=south, product_id=PROD002, amount=1, date=2024-01-20}
3. [1] TestOrder{id=O003, region=north, customer_id=P001, product_id=PROD003, date=2024-02-01, amount=3, total=225, status=shipped, priority=high, discount=15}
4. [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0, customer_id=P004, amount=1}
5. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, status=confirmed, region=south, customer_id=P002, product_id=PROD001}
6. [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, region=west, amount=2, date=2024-02-15, priority=low, discount=0, customer_id=P005}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, discount=50, region=north, customer_id=P007, total=600, date=2024-03-01, status=shipped, priority=urgent}
8. [1] TestOrder{id=O008, amount=10, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05}
9. [1] TestOrder{id=O009, product_id=PROD007, total=89.99, date=2024-03-10, discount=10, customer_id=P001, amount=1, status=completed, priority=low, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, date=2024-03-15, priority=urgent, amount=1, total=75000, status=refunded, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, customer_id=P002, total=25.5, status=confirmed, priority=low, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, status=confirmed, region=south}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O006, priority=low, discount=0, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, region=west, amount=2, date=2024-02-15}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, priority=urgent, product_id=PROD006, amount=4, discount=50, region=north, customer_id=P007, total=600, date=2024-03-01, status=shipped}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O001, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, region=north, product_id=PROD001, amount=2}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003, date=2024-02-01}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O008, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, amount=10, status=pending, priority=normal}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, discount=10, customer_id=P001, amount=1, status=completed, priority=low, region=north, product_id=PROD007}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O010, discount=0, region=east, customer_id=P006, product_id=PROD001, date=2024-03-15, priority=urgent, amount=1, total=75000, status=refunded}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, status=active, department=sales, age=25, active=true, tags=junior, level=2, name=Alice, salary=45000}
2. [1] TestPerson{id=P002, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, status=active, salary=75000, active=true}
3. [1] TestPerson{id=P003, department=hr, tags=intern, status=inactive, level=1, name=Charlie, age=16, salary=0, active=false, score=6}
4. [1] TestPerson{id=P004, tags=manager, status=active, level=7, name=Diana, salary=85000, score=7.8, department=marketing, age=45, active=true}
5. [1] TestPerson{id=P005, age=30, salary=55000, active=false, tags=employee, status=inactive, department=sales, level=3, name=Eve, score=8}
6. [1] TestPerson{id=P006, level=1, name=Frank, active=true, score=0, age=0, salary=-5000, tags=test, status=active, department=qa}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, level=1, active=false, score=5.5, department=support}
9. [1] TestPerson{id=P009, department=engineering, name=Ivy, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active}
10. [1] TestPerson{id=P010, score=6.5, department=intern, level=1, active=true, tags=temp, status=active, name=X, age=22, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, age=16, salary=0, active=false, score=6, department=hr, tags=intern, status=inactive, level=1, name=Charlie}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, active=false, tags=employee, status=inactive, department=sales, level=3, name=Eve, score=8}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P008, level=1, active=false, score=5.5, department=support, name=Henry, age=18, salary=25000, tags=junior, status=inactive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, region=north, product_id=PROD001, amount=2, priority=normal, discount=50}
2. [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, priority=low, region=south, product_id=PROD002, amount=1, date=2024-01-20, discount=0}
3. [1] TestOrder{id=O003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003, date=2024-02-01}
4. [1] TestOrder{id=O004, priority=normal, region=east, product_id=PROD004, discount=0, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered}
5. [1] TestOrder{id=O005, region=south, customer_id=P002, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, status=confirmed}
6. [1] TestOrder{id=O006, status=cancelled, region=west, amount=2, date=2024-02-15, priority=low, discount=0, customer_id=P005, product_id=PROD005, total=999.98}
7. [1] TestOrder{id=O007, region=north, customer_id=P007, total=600, date=2024-03-01, status=shipped, priority=urgent, product_id=PROD006, amount=4, discount=50}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, amount=10, status=pending, priority=normal, discount=0, region=south}
9. [1] TestOrder{id=O009, status=completed, priority=low, region=north, product_id=PROD007, total=89.99, date=2024-03-10, discount=10, customer_id=P001, amount=1}
10. [1] TestOrder{id=O010, amount=1, total=75000, status=refunded, discount=0, region=east, customer_id=P006, product_id=PROD001, date=2024-03-15, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, amount=3, total=225, status=shipped, priority=high, discount=15, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O005, status=confirmed, region=south, customer_id=P002, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O006, amount=2, date=2024-02-15, priority=low, discount=0, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, region=west}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O007, region=north, customer_id=P007, total=600, date=2024-03-01, status=shipped, priority=urgent, product_id=PROD006, amount=4, discount=50}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O009, status=completed, priority=low, region=north, product_id=PROD007, total=89.99, date=2024-03-10, discount=10, customer_id=P001, amount=1}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O010, amount=1, total=75000, status=refunded, discount=0, region=east, customer_id=P006, product_id=PROD001, date=2024-03-15, priority=urgent}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, status=confirmed, priority=low, region=south, product_id=PROD002, amount=1, date=2024-01-20, discount=0, customer_id=P002, total=25.5}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, status=active, department=sales, age=25, active=true, tags=junior, level=2}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, status=active}
3. [1] TestPerson{id=P003, active=false, score=6, department=hr, tags=intern, status=inactive, level=1, name=Charlie, age=16, salary=0}
4. [1] TestPerson{id=P004, active=true, tags=manager, status=active, level=7, name=Diana, salary=85000, score=7.8, department=marketing, age=45}
5. [1] TestPerson{id=P005, age=30, salary=55000, active=false, tags=employee, status=inactive, department=sales, level=3, name=Eve, score=8}
6. [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, department=qa, level=1, name=Frank, active=true, score=0}
7. [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, level=1, active=false, score=5.5, department=support}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, age=22, salary=28000, score=6.5, department=intern, level=1, active=true, tags=temp, status=active, name=X}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, age=18, salary=25000, tags=junior, status=inactive, level=1, active=false, score=5.5, department=support, name=Henry}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, level=1, name=Frank, active=true, score=0, age=0, salary=-5000, tags=test, status=active, department=qa}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, department=engineering}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, status=active, department=sales, age=25, active=true, tags=junior, level=2, name=Alice, salary=45000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P003, score=6, department=hr, tags=intern, status=inactive, level=1, name=Charlie, age=16, salary=0, active=false}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, department=marketing, age=45, active=true, tags=manager, status=active, level=7}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P005, status=inactive, department=sales, level=3, name=Eve, score=8, age=30, salary=55000, active=false, tags=employee}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, supplier=TechSupply, category=electronics, price=999.99, available=true, keywords=computer, stock=50, name=Laptop, rating=4.5, brand=TechCorp}
2. [1] TestProduct{id=PROD002, supplier=TechSupply, name=Mouse, category=accessories, available=true, keywords=peripheral, brand=TechCorp, price=25.5, rating=4, stock=200}
3. [1] TestProduct{id=PROD003, category=accessories, price=75, brand=KeyTech, stock=0, supplier=KeySupply, name=Keyboard, available=false, rating=3.5, keywords=typing}
4. [1] TestProduct{id=PROD004, name=Monitor, price=299.99, keywords=display, brand=ScreenPro, supplier=ScreenSupply, category=electronics, available=true, rating=4.8, stock=30}
5. [1] TestProduct{id=PROD005, supplier=OldSupply, name=OldKeyboard, price=8.5, available=false, rating=2, keywords=obsolete, stock=0, category=accessories, brand=OldTech}
6. [1] TestProduct{id=PROD006, supplier=AudioSupply, category=audio, available=true, rating=4.6, brand=AudioMax, stock=75, name=Headphones, price=150, keywords=sound}
7. [1] TestProduct{id=PROD007, keywords=video, brand=CamTech, name=Webcam, category=electronics, available=true, stock=25, supplier=CamSupply, price=89.99, rating=3.8}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{id=PROD001, category=electronics, price=999.99, available=true, keywords=computer, stock=50, name=Laptop, rating=4.5, brand=TechCorp, supplier=TechSupply}

2. **Token 2**:
   - Fait 1: [1] TestProduct{id=PROD002, name=Mouse, category=accessories, available=true, keywords=peripheral, brand=TechCorp, price=25.5, rating=4, stock=200, supplier=TechSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{id=PROD003, stock=0, supplier=KeySupply, name=Keyboard, available=false, rating=3.5, keywords=typing, category=accessories, price=75, brand=KeyTech}

4. **Token 4**:
   - Fait 1: [1] TestProduct{id=PROD004, category=electronics, available=true, rating=4.8, stock=30, name=Monitor, price=299.99, keywords=display, brand=ScreenPro, supplier=ScreenSupply}

5. **Token 5**:
   - Fait 1: [1] TestProduct{id=PROD006, available=true, rating=4.6, brand=AudioMax, stock=75, name=Headphones, price=150, keywords=sound, supplier=AudioSupply, category=audio}

6. **Token 6**:
   - Fait 1: [1] TestProduct{id=PROD007, brand=CamTech, name=Webcam, category=electronics, available=true, stock=25, supplier=CamSupply, price=89.99, rating=3.8, keywords=video}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, name=Alice, salary=45000, score=8.5, status=active, department=sales, age=25, active=true, tags=junior}
2. [1] TestPerson{id=P002, age=35, status=active, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, name=Bob}
3. [1] TestPerson{id=P003, tags=intern, status=inactive, level=1, name=Charlie, age=16, salary=0, active=false, score=6, department=hr}
4. [1] TestPerson{id=P004, age=45, active=true, tags=manager, status=active, level=7, name=Diana, salary=85000, score=7.8, department=marketing}
5. [1] TestPerson{id=P005, active=false, tags=employee, status=inactive, department=sales, level=3, name=Eve, score=8, age=30, salary=55000}
6. [1] TestPerson{id=P006, active=true, score=0, age=0, salary=-5000, tags=test, status=active, department=qa, level=1, name=Frank}
7. [1] TestPerson{id=P007, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, level=1, active=false, score=5.5, department=support}
9. [1] TestPerson{id=P009, age=40, tags=senior, status=active, department=engineering, name=Ivy, salary=68000, active=true, score=8.7, level=6}
10. [1] TestPerson{id=P010, level=1, active=true, tags=temp, status=active, name=X, age=22, salary=28000, score=6.5, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, salary=45000, score=8.5, status=active, department=sales, age=25, active=true, tags=junior, level=2, name=Alice}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, status=active, salary=75000, active=true}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P005, age=30, salary=55000, active=false, tags=employee, status=inactive, department=sales, level=3, name=Eve, score=8}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P006, active=true, score=0, age=0, salary=-5000, tags=test, status=active, department=qa, level=1, name=Frank}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, level=1, active=false, score=5.5, department=support}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, department=hr, tags=intern, status=inactive, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, department=marketing, age=45, active=true, tags=manager, status=active, level=7}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, department=engineering}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, amount=2, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, region=north, product_id=PROD001}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, customer_id=P002, total=25.5, status=confirmed, priority=low, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, amount=3, total=225, status=shipped, priority=high, discount=15, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, status=confirmed, region=south}
6. [1] TestOrder{id=O006, amount=2, date=2024-02-15, priority=low, discount=0, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, region=west}
7. [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, priority=urgent, product_id=PROD006, amount=4, discount=50, region=north}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, amount=10, status=pending, priority=normal, discount=0, region=south}
9. [1] TestOrder{id=O009, product_id=PROD007, total=89.99, date=2024-03-10, discount=10, customer_id=P001, amount=1, status=completed, priority=low, region=north}
10. [1] TestOrder{id=O010, amount=1, total=75000, status=refunded, discount=0, region=east, customer_id=P006, product_id=PROD001, date=2024-03-15, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O008, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, amount=10}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, status=completed, priority=low, region=north, product_id=PROD007, total=89.99, date=2024-03-10, discount=10}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O001, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, region=north, product_id=PROD001, amount=2}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, customer_id=P002, total=25.5, status=confirmed, priority=low, region=south}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, amount=3, total=225, status=shipped, priority=high, discount=15, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0, customer_id=P004, amount=1}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, status=confirmed, region=south}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, discount=50, region=north, customer_id=P007, total=600, date=2024-03-01, status=shipped, priority=urgent}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, age=25, active=true, tags=junior, level=2, name=Alice, salary=45000, score=8.5, status=active, department=sales}
2. [1] TestPerson{id=P002, name=Bob, age=35, status=active, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, department=hr, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, department=marketing, age=45, active=true, tags=manager, status=active, level=7}
5. [1] TestPerson{id=P005, tags=employee, status=inactive, department=sales, level=3, name=Eve, score=8, age=30, salary=55000, active=false}
6. [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, department=qa, level=1, name=Frank, active=true, score=0}
7. [1] TestPerson{id=P007, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, level=1, active=false, score=5.5, department=support}
9. [1] TestPerson{id=P009, level=6, age=40, tags=senior, status=active, department=engineering, name=Ivy, salary=68000, active=true, score=8.7}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, department=intern, level=1, active=true, tags=temp, status=active}
11. [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, region=north, product_id=PROD001, amount=2, priority=normal, discount=50}
12. [1] TestOrder{id=O002, total=25.5, status=confirmed, priority=low, region=south, product_id=PROD002, amount=1, date=2024-01-20, discount=0, customer_id=P002}
13. [1] TestOrder{id=O003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003, date=2024-02-01}
14. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, status=confirmed, region=south}
16. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, region=west, amount=2, date=2024-02-15, priority=low, discount=0}
17. [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, priority=urgent, product_id=PROD006, amount=4, discount=50, region=north}
18. [1] TestOrder{id=O008, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, amount=10, status=pending, priority=normal, discount=0}
19. [1] TestOrder{id=O009, customer_id=P001, amount=1, status=completed, priority=low, region=north, product_id=PROD007, total=89.99, date=2024-03-10, discount=10}
20. [1] TestOrder{id=O010, amount=1, total=75000, status=refunded, discount=0, region=east, customer_id=P006, product_id=PROD001, date=2024-03-15, priority=urgent}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, department=marketing, age=45, active=true, tags=manager, status=active, level=7}
   - Fait 2: [1] TestOrder{id=O004, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0, customer_id=P004, amount=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P010, active=true, tags=temp, status=active, name=X, age=22, salary=28000, score=6.5, department=intern, level=1}
   - Fait 2: [1] TestOrder{id=O008, amount=10, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, level=1, name=Frank, active=true, score=0, age=0, salary=-5000, tags=test, status=active}
   - Fait 2: [1] TestOrder{id=O010, total=75000, status=refunded, discount=0, region=east, customer_id=P006, product_id=PROD001, date=2024-03-15, priority=urgent, amount=1}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, status=active, department=sales, age=25, active=true, tags=junior, level=2, name=Alice, salary=45000, score=8.5}
   - Fait 2: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, region=north, product_id=PROD001, amount=2, priority=normal, discount=50}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, status=active, department=sales, age=25, active=true, tags=junior, level=2}
   - Fait 2: [1] TestOrder{id=O003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003, date=2024-02-01}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, status=active}
   - Fait 2: [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, status=confirmed, region=south}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, tags=employee, status=inactive, department=sales, level=3, name=Eve, score=8, age=30, salary=55000, active=false}
   - Fait 2: [1] TestOrder{id=O006, status=cancelled, region=west, amount=2, date=2024-02-15, priority=low, discount=0, customer_id=P005, product_id=PROD005, total=999.98}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P007, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive}
   - Fait 2: [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, priority=urgent, product_id=PROD006, amount=4, discount=50, region=north}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, status=active, department=sales, age=25, active=true, tags=junior, level=2}
   - Fait 2: [1] TestOrder{id=O009, product_id=PROD007, total=89.99, date=2024-03-10, discount=10, customer_id=P001, amount=1, status=completed, priority=low, region=north}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, status=active}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, customer_id=P002, total=25.5, status=confirmed, priority=low, region=south}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/20 (50.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, name=Alice, salary=45000, score=8.5, status=active, department=sales, age=25, active=true, tags=junior}
2. [1] TestPerson{id=P002, status=active, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35}
3. [1] TestPerson{id=P003, salary=0, active=false, score=6, department=hr, tags=intern, status=inactive, level=1, name=Charlie, age=16}
4. [1] TestPerson{id=P004, tags=manager, status=active, level=7, name=Diana, salary=85000, score=7.8, department=marketing, age=45, active=true}
5. [1] TestPerson{id=P005, age=30, salary=55000, active=false, tags=employee, status=inactive, department=sales, level=3, name=Eve, score=8}
6. [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, department=qa, level=1, name=Frank, active=true, score=0}
7. [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}
8. [1] TestPerson{id=P008, active=false, score=5.5, department=support, name=Henry, age=18, salary=25000, tags=junior, status=inactive, level=1}
9. [1] TestPerson{id=P009, age=40, tags=senior, status=active, department=engineering, name=Ivy, salary=68000, active=true, score=8.7, level=6}
10. [1] TestPerson{id=P010, tags=temp, status=active, name=X, age=22, salary=28000, score=6.5, department=intern, level=1, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P006, department=qa, level=1, name=Frank, active=true, score=0, age=0, salary=-5000, tags=test, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, department=engineering}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P010, status=active, name=X, age=22, salary=28000, score=6.5, department=intern, level=1, active=true, tags=temp}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, age=25, active=true, tags=junior, level=2, name=Alice, salary=45000, score=8.5, status=active, department=sales}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P002, active=true, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, status=active, salary=75000}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, tags=manager, status=active, level=7, name=Diana, salary=85000, score=7.8, department=marketing}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, level=2, name=Alice, salary=45000, score=8.5, status=active, department=sales, age=25, active=true, tags=junior}
2. [1] TestPerson{id=P002, name=Bob, age=35, status=active, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5}
3. [1] TestPerson{id=P003, level=1, name=Charlie, age=16, salary=0, active=false, score=6, department=hr, tags=intern, status=inactive}
4. [1] TestPerson{id=P004, age=45, active=true, tags=manager, status=active, level=7, name=Diana, salary=85000, score=7.8, department=marketing}
5. [1] TestPerson{id=P005, salary=55000, active=false, tags=employee, status=inactive, department=sales, level=3, name=Eve, score=8, age=30}
6. [1] TestPerson{id=P006, name=Frank, active=true, score=0, age=0, salary=-5000, tags=test, status=active, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9}
8. [1] TestPerson{id=P008, score=5.5, department=support, name=Henry, age=18, salary=25000, tags=junior, status=inactive, level=1, active=false}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, department=intern, level=1, active=true, tags=temp, status=active}
11. [1] TestOrder{id=O001, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, region=north, product_id=PROD001, amount=2}
12. [1] TestOrder{id=O002, priority=low, region=south, product_id=PROD002, amount=1, date=2024-01-20, discount=0, customer_id=P002, total=25.5, status=confirmed}
13. [1] TestOrder{id=O003, amount=3, total=225, status=shipped, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003, date=2024-02-01}
14. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}
15. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, status=confirmed, region=south}
16. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, region=west, amount=2, date=2024-02-15, priority=low, discount=0}
17. [1] TestOrder{id=O007, customer_id=P007, total=600, date=2024-03-01, status=shipped, priority=urgent, product_id=PROD006, amount=4, discount=50, region=north}
18. [1] TestOrder{id=O008, amount=10, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05}
19. [1] TestOrder{id=O009, customer_id=P001, amount=1, status=completed, priority=low, region=north, product_id=PROD007, total=89.99, date=2024-03-10, discount=10}
20. [1] TestOrder{id=O010, priority=urgent, amount=1, total=75000, status=refunded, discount=0, region=east, customer_id=P006, product_id=PROD001, date=2024-03-15}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P010, score=6.5, department=intern, level=1, active=true, tags=temp, status=active, name=X, age=22, salary=28000}
   - Fait 2: [1] TestOrder{id=O008, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, amount=10, status=pending, priority=normal}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P001, age=25, active=true, tags=junior, level=2, name=Alice, salary=45000, score=8.5, status=active, department=sales}
   - Fait 2: [1] TestOrder{id=O009, amount=1, status=completed, priority=low, region=north, product_id=PROD007, total=89.99, date=2024-03-10, discount=10, customer_id=P001}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, tags=manager, status=active, level=7, name=Diana, salary=85000, score=7.8, department=marketing, age=45, active=true}
   - Fait 2: [1] TestOrder{id=O004, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0, customer_id=P004, amount=1, total=299.99}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, score=8, age=30, salary=55000, active=false, tags=employee, status=inactive, department=sales, level=3, name=Eve}
   - Fait 2: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, region=west, amount=2, date=2024-02-15, priority=low, discount=0}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9}
   - Fait 2: [1] TestOrder{id=O007, amount=4, discount=50, region=north, customer_id=P007, total=600, date=2024-03-01, status=shipped, priority=urgent, product_id=PROD006}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, department=qa, level=1, name=Frank, active=true, score=0}
   - Fait 2: [1] TestOrder{id=O010, date=2024-03-15, priority=urgent, amount=1, total=75000, status=refunded, discount=0, region=east, customer_id=P006, product_id=PROD001}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P001, age=25, active=true, tags=junior, level=2, name=Alice, salary=45000, score=8.5, status=active, department=sales}
   - Fait 2: [1] TestOrder{id=O001, total=1999.98, date=2024-01-15, status=pending, region=north, product_id=PROD001, amount=2, priority=normal, discount=50, customer_id=P001}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, status=active}
   - Fait 2: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, customer_id=P002, total=25.5, status=confirmed, priority=low, region=south}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, status=active, department=sales, age=25, active=true, tags=junior, level=2}
   - Fait 2: [1] TestOrder{id=O003, status=shipped, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003, date=2024-02-01, amount=3, total=225}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P002, tags=senior, department=engineering, level=5, name=Bob, age=35, status=active, salary=75000, active=true, score=9.2}
   - Fait 2: [1] TestOrder{id=O005, date=2024-02-10, priority=high, discount=100, status=confirmed, region=south, customer_id=P002, product_id=PROD001, amount=1, total=999.99}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/20 (50.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, score=8.5, status=active, department=sales, age=25, active=true, tags=junior, level=2, name=Alice, salary=45000}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, status=active}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, department=hr, tags=intern, status=inactive, level=1}
4. [1] TestPerson{id=P004, age=45, active=true, tags=manager, status=active, level=7, name=Diana, salary=85000, score=7.8, department=marketing}
5. [1] TestPerson{id=P005, age=30, salary=55000, active=false, tags=employee, status=inactive, department=sales, level=3, name=Eve, score=8}
6. [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, department=qa, level=1, name=Frank, active=true, score=0}
7. [1] TestPerson{id=P007, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000}
8. [1] TestPerson{id=P008, score=5.5, department=support, name=Henry, age=18, salary=25000, tags=junior, status=inactive, level=1, active=false}
9. [1] TestPerson{id=P009, level=6, age=40, tags=senior, status=active, department=engineering, name=Ivy, salary=68000, active=true, score=8.7}
10. [1] TestPerson{id=P010, level=1, active=true, tags=temp, status=active, name=X, age=22, salary=28000, score=6.5, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, department=hr, tags=intern, status=inactive, level=1}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P004, score=7.8, department=marketing, age=45, active=true, tags=manager, status=active, level=7, name=Diana, salary=85000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, score=8, age=30, salary=55000, active=false, tags=employee, status=inactive, department=sales, level=3}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, score=8.5, status=active, department=sales, age=25, active=true, tags=junior, level=2}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000, active=true, score=10, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P008, department=support, name=Henry, age=18, salary=25000, tags=junior, status=inactive, level=1, active=false, score=5.5}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, age=40, tags=senior, status=active, department=engineering, name=Ivy, salary=68000, active=true, score=8.7, level=6}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, status=active, name=X, age=22, salary=28000, score=6.5, department=intern, level=1, active=true, tags=temp}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, priority=normal, discount=50, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, total=25.5, status=confirmed, priority=low, region=south, product_id=PROD002, amount=1, date=2024-01-20, discount=0}
3. [1] TestOrder{id=O003, priority=high, discount=15, region=north, customer_id=P001, product_id=PROD003, date=2024-02-01, amount=3, total=225, status=shipped}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east, product_id=PROD004, discount=0}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, status=confirmed, region=south}
6. [1] TestOrder{id=O006, amount=2, date=2024-02-15, priority=low, discount=0, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, region=west}
7. [1] TestOrder{id=O007, priority=urgent, product_id=PROD006, amount=4, discount=50, region=north, customer_id=P007, total=600, date=2024-03-01, status=shipped}
8. [1] TestOrder{id=O008, date=2024-03-05, amount=10, status=pending, priority=normal, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255}
9. [1] TestOrder{id=O009, product_id=PROD007, total=89.99, date=2024-03-10, discount=10, customer_id=P001, amount=1, status=completed, priority=low, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, product_id=PROD001, date=2024-03-15, priority=urgent, amount=1, total=75000, status=refunded, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O008, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, amount=10, status=pending, priority=normal}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O009, customer_id=P001, amount=1, status=completed, priority=low, region=north, product_id=PROD007, total=89.99, date=2024-03-10, discount=10}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O001, customer_id=P001, total=1999.98, date=2024-01-15, status=pending, region=north, product_id=PROD001, amount=2, priority=normal, discount=50}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O004, product_id=PROD004, discount=0, customer_id=P004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, region=east}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, region=west, amount=2, date=2024-02-15, priority=low, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O007, product_id=PROD006, amount=4, discount=50, region=north, customer_id=P007, total=600, date=2024-03-01, status=shipped, priority=urgent}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O010, date=2024-03-15, priority=urgent, amount=1, total=75000, status=refunded, discount=0, region=east, customer_id=P006, product_id=PROD001}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O002, product_id=PROD002, amount=1, date=2024-01-20, discount=0, customer_id=P002, total=25.5, status=confirmed, priority=low, region=south}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, product_id=PROD003, date=2024-02-01, amount=3, total=225, status=shipped, priority=high, discount=15, region=north}

10. **Token 10**:
   - Fait 1: [1] TestOrder{id=O005, region=south, customer_id=P002, product_id=PROD001, amount=1, total=999.99, date=2024-02-10, priority=high, discount=100, status=confirmed}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, score=8.5, status=active, department=sales, age=25, active=true, tags=junior, level=2, name=Alice}
2. [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, status=active}
3. [1] TestPerson{id=P003, status=inactive, level=1, name=Charlie, age=16, salary=0, active=false, score=6, department=hr, tags=intern}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, department=marketing, age=45, active=true, tags=manager, status=active, level=7}
5. [1] TestPerson{id=P005, status=inactive, department=sales, level=3, name=Eve, score=8, age=30, salary=55000, active=false, tags=employee}
6. [1] TestPerson{id=P006, salary=-5000, tags=test, status=active, department=qa, level=1, name=Frank, active=true, score=0, age=0}
7. [1] TestPerson{id=P007, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9, age=65, salary=95000}
8. [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, level=1, active=false, score=5.5, department=support}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, department=engineering}
10. [1] TestPerson{id=P010, name=X, age=22, salary=28000, score=6.5, department=intern, level=1, active=true, tags=temp, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P004, name=Diana, salary=85000, score=7.8, department=marketing, age=45, active=true, tags=manager, status=active, level=7}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, age=18, salary=25000, tags=junior, status=inactive, level=1, active=false, score=5.5, department=support}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P010, salary=28000, score=6.5, department=intern, level=1, active=true, tags=temp, status=active, name=X, age=22}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P001, score=8.5, status=active, department=sales, age=25, active=true, tags=junior, level=2, name=Alice, salary=45000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P002, salary=75000, active=true, score=9.2, tags=senior, department=engineering, level=5, name=Bob, age=35, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P003, level=1, name=Charlie, age=16, salary=0, active=false, score=6, department=hr, tags=intern, status=inactive}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P005, department=sales, level=3, name=Eve, score=8, age=30, salary=55000, active=false, tags=employee, status=inactive}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P006, age=0, salary=-5000, tags=test, status=active, department=qa, level=1, name=Frank, active=true, score=0}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P007, age=65, salary=95000, active=true, score=10, status=active, name=Grace, tags=executive, department=management, level=9}

10. **Token 10**:
   - Fait 1: [1] TestPerson{id=P009, score=8.7, level=6, age=40, tags=senior, status=active, department=engineering, name=Ivy, salary=68000, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 150
- **Faits traités**: 27
