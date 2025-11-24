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

1. [1] TestPerson{level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales}
2. [1] TestPerson{active=true, tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5}
3. [1] TestPerson{name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive}
4. [1] TestPerson{name=Diana, active=true, status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing}
5. [1] TestPerson{salary=55000, name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales}
6. [1] TestPerson{department=qa, age=0, score=0, tags=test, active=true, level=1, name=Frank, status=active, salary=-5000}
7. [1] TestPerson{active=true, salary=95000, department=management, age=65, score=10, tags=executive, status=active, level=9, name=Grace}
8. [1] TestPerson{status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18}
9. [1] TestPerson{level=6, salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7}
10. [1] TestPerson{department=intern, level=1, name=X, score=6.5, tags=temp, active=true, age=22, salary=28000, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{age=18, status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry}

2. **Token 2**:
   - Fait 1: [1] TestPerson{active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy}

3. **Token 3**:
   - Fait 1: [1] TestPerson{level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales}

4. **Token 4**:
   - Fait 1: [1] TestPerson{salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true, age=22}

5. **Token 5**:
   - Fait 1: [1] TestPerson{tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5, active=true}

6. **Token 6**:
   - Fait 1: [1] TestPerson{name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive}

7. **Token 7**:
   - Fait 1: [1] TestPerson{status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true}

8. **Token 8**:
   - Fait 1: [1] TestPerson{salary=55000, name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales}

9. **Token 9**:
   - Fait 1: [1] TestPerson{age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}
2. [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}
3. [1] TestOrder{product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped}
4. [1] TestOrder{product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05}
5. [1] TestOrder{status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1}
6. [1] TestOrder{customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low}
7. [1] TestOrder{date=2024-03-01, priority=urgent, amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007}
8. [1] TestOrder{priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002}
9. [1] TestOrder{customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10, date=2024-03-10, region=north}
10. [1] TestOrder{product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10}

2. **Token 2**:
   - Fait 1: [1] TestOrder{amount=10, discount=0, customer_id=P010, product_id=PROD002, priority=normal, region=south, total=255, date=2024-03-05, status=pending}

3. **Token 3**:
   - Fait 1: [1] TestOrder{date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10}

4. **Token 4**:
   - Fait 1: [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{priority=low, customer_id=P002, region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed}

6. **Token 6**:
   - Fait 1: [1] TestOrder{discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4}

7. **Token 7**:
   - Fait 1: [1] TestOrder{date=2024-03-15, discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded}

8. **Token 8**:
   - Fait 1: [1] TestOrder{customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north}

9. **Token 9**:
   - Fait 1: [1] TestOrder{date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales}
2. [1] TestPerson{salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering, age=35}
3. [1] TestPerson{status=inactive, name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0}
4. [1] TestPerson{status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true}
5. [1] TestPerson{name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000}
6. [1] TestPerson{age=0, score=0, tags=test, active=true, level=1, name=Frank, status=active, salary=-5000, department=qa}
7. [1] TestPerson{age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management}
8. [1] TestPerson{age=18, status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry}
9. [1] TestPerson{status=active, name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40}
10. [1] TestPerson{status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true, age=22, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales, level=2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering, age=35, salary=75000}

3. **Token 3**:
   - Fait 1: [1] TestPerson{active=true, status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana}

4. **Token 4**:
   - Fait 1: [1] TestPerson{name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management}

6. **Token 6**:
   - Fait 1: [1] TestPerson{age=40, status=active, name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}
2. [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}
3. [1] TestOrder{status=shipped, product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15, amount=3, total=225}
4. [1] TestOrder{customer_id=P004, priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99}
5. [1] TestOrder{priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99}
6. [1] TestOrder{customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low}
7. [1] TestOrder{total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4, discount=50, product_id=PROD006}
8. [1] TestOrder{priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002}
9. [1] TestOrder{total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1}
10. [1] TestOrder{discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{product_id=PROD002, priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010}

2. **Token 2**:
   - Fait 1: [1] TestOrder{customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north}

3. **Token 3**:
   - Fait 1: [1] TestOrder{customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed}

4. **Token 4**:
   - Fait 1: [1] TestOrder{total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low, customer_id=P005, amount=2}

5. **Token 5**:
   - Fait 1: [1] TestOrder{discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4}

6. **Token 6**:
   - Fait 1: [1] TestOrder{product_id=PROD007, priority=low, amount=1, total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed}

7. **Token 7**:
   - Fait 1: [1] TestOrder{amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, customer_id=P001, product_id=PROD001, total=1999.98}

8. **Token 8**:
   - Fait 1: [1] TestOrder{total=25.5, status=confirmed, priority=low, customer_id=P002, region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1}

9. **Token 9**:
   - Fait 1: [1] TestOrder{product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{supplier=TechSupply, category=electronics, available=true, brand=TechCorp, stock=50, name=Laptop, rating=4.5, keywords=computer, price=999.99}
2. [1] TestProduct{stock=200, supplier=TechSupply, name=Mouse, category=accessories, rating=4, brand=TechCorp, price=25.5, keywords=peripheral, available=true}
3. [1] TestProduct{supplier=KeySupply, name=Keyboard, keywords=typing, stock=0, rating=3.5, brand=KeyTech, available=false, category=accessories, price=75}
4. [1] TestProduct{name=Monitor, category=electronics, price=299.99, available=true, keywords=display, brand=ScreenPro, stock=30, rating=4.8, supplier=ScreenSupply}
5. [1] TestProduct{stock=0, supplier=OldSupply, category=accessories, price=8.5, available=false, rating=2, keywords=obsolete, name=OldKeyboard, brand=OldTech}
6. [1] TestProduct{available=true, rating=4.6, keywords=sound, supplier=AudioSupply, name=Headphones, stock=75, category=audio, price=150, brand=AudioMax}
7. [1] TestProduct{name=Webcam, available=true, keywords=video, brand=CamTech, price=89.99, rating=3.8, stock=25, supplier=CamSupply, category=electronics}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{name=Webcam, available=true, keywords=video, brand=CamTech, price=89.99, rating=3.8, stock=25, supplier=CamSupply, category=electronics}

2. **Token 2**:
   - Fait 1: [1] TestProduct{keywords=computer, price=999.99, supplier=TechSupply, category=electronics, available=true, brand=TechCorp, stock=50, name=Laptop, rating=4.5}

3. **Token 3**:
   - Fait 1: [1] TestProduct{stock=200, supplier=TechSupply, name=Mouse, category=accessories, rating=4, brand=TechCorp, price=25.5, keywords=peripheral, available=true}

4. **Token 4**:
   - Fait 1: [1] TestProduct{supplier=KeySupply, name=Keyboard, keywords=typing, stock=0, rating=3.5, brand=KeyTech, available=false, category=accessories, price=75}

5. **Token 5**:
   - Fait 1: [1] TestProduct{category=electronics, price=299.99, available=true, keywords=display, brand=ScreenPro, stock=30, rating=4.8, supplier=ScreenSupply, name=Monitor}

6. **Token 6**:
   - Fait 1: [1] TestProduct{stock=75, category=audio, price=150, brand=AudioMax, available=true, rating=4.6, keywords=sound, supplier=AudioSupply, name=Headphones}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales}
2. [1] TestPerson{name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior}
3. [1] TestPerson{age=16, tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false, score=6}
4. [1] TestPerson{status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true}
5. [1] TestPerson{name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000}
6. [1] TestPerson{age=0, score=0, tags=test, active=true, level=1, name=Frank, status=active, salary=-5000, department=qa}
7. [1] TestPerson{active=true, salary=95000, department=management, age=65, score=10, tags=executive, status=active, level=9, name=Grace}
8. [1] TestPerson{department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive}
9. [1] TestPerson{salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6}
10. [1] TestPerson{age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false, score=6, age=16}

2. **Token 2**:
   - Fait 1: [1] TestPerson{tags=test, active=true, level=1, name=Frank, status=active, salary=-5000, department=qa, age=0, score=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{priority=normal, discount=50, region=north, customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending}
2. [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}
3. [1] TestOrder{customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north}
4. [1] TestOrder{status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1}
5. [1] TestOrder{discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high}
6. [1] TestOrder{customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low}
7. [1] TestOrder{amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent}
8. [1] TestOrder{total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002, priority=normal, region=south}
9. [1] TestOrder{status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001}
10. [1] TestOrder{priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15, discount=0, region=east, product_id=PROD001}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{customer_id=P002, region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low}

2. **Token 2**:
   - Fait 1: [1] TestOrder{date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003}

3. **Token 3**:
   - Fait 1: [1] TestOrder{customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low}

4. **Token 4**:
   - Fait 1: [1] TestOrder{date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10}

5. **Token 5**:
   - Fait 1: [1] TestOrder{region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{region=north, customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50}

7. **Token 7**:
   - Fait 1: [1] TestOrder{customer_id=P004, priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99}

8. **Token 8**:
   - Fait 1: [1] TestOrder{region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100}

9. **Token 9**:
   - Fait 1: [1] TestOrder{priority=urgent, amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01}

10. **Token 10**:
   - Fait 1: [1] TestOrder{priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{active=true, name=Alice, age=25, tags=junior, department=sales, level=2, salary=45000, score=8.5, status=active}
2. [1] TestPerson{salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering, age=35}
3. [1] TestPerson{age=16, tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false, score=6}
4. [1] TestPerson{score=7.8, department=marketing, name=Diana, active=true, status=active, age=45, level=7, salary=85000, tags=manager}
5. [1] TestPerson{age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000, name=Eve}
6. [1] TestPerson{name=Frank, status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1}
7. [1] TestPerson{salary=95000, department=management, age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true}
8. [1] TestPerson{salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1}
9. [1] TestPerson{salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6}
10. [1] TestPerson{score=6.5, tags=temp, active=true, age=22, salary=28000, status=active, department=intern, level=1, name=X}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{age=16, tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false, score=6}

2. **Token 2**:
   - Fait 1: [1] TestPerson{tags=employee, department=sales, salary=55000, name=Eve, age=30, active=false, status=inactive, score=8, level=3}

3. **Token 3**:
   - Fait 1: [1] TestPerson{age=18, status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, customer_id=P001, product_id=PROD001, total=1999.98}
2. [1] TestOrder{amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002, region=south, product_id=PROD002, date=2024-01-20, discount=0}
3. [1] TestOrder{amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15}
4. [1] TestOrder{product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05}
5. [1] TestOrder{status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1}
6. [1] TestOrder{product_id=PROD005, date=2024-02-15, discount=0, priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west}
7. [1] TestOrder{amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent}
8. [1] TestOrder{amount=10, discount=0, customer_id=P010, product_id=PROD002, priority=normal, region=south, total=255, date=2024-03-05, status=pending}
9. [1] TestOrder{amount=1, total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low}
10. [1] TestOrder{customer_id=P006, total=75000, status=refunded, date=2024-03-15, discount=0, region=east, product_id=PROD001, priority=urgent, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{total=25.5, status=confirmed, priority=low, customer_id=P002, region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1}

2. **Token 2**:
   - Fait 1: [1] TestOrder{product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped}

3. **Token 3**:
   - Fait 1: [1] TestOrder{customer_id=P004, priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99}

4. **Token 4**:
   - Fait 1: [1] TestOrder{total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002}

5. **Token 5**:
   - Fait 1: [1] TestOrder{product_id=PROD005, date=2024-02-15, discount=0, priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west}

6. **Token 6**:
   - Fait 1: [1] TestOrder{amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent}

7. **Token 7**:
   - Fait 1: [1] TestOrder{status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001}

8. **Token 8**:
   - Fait 1: [1] TestOrder{discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{active=true, name=Alice, age=25, tags=junior, department=sales, level=2, salary=45000, score=8.5, status=active}
2. [1] TestPerson{salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering, age=35}
3. [1] TestPerson{name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive}
4. [1] TestPerson{status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true}
5. [1] TestPerson{name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000}
6. [1] TestPerson{age=0, score=0, tags=test, active=true, level=1, name=Frank, status=active, salary=-5000, department=qa}
7. [1] TestPerson{active=true, salary=95000, department=management, age=65, score=10, tags=executive, status=active, level=9, name=Grace}
8. [1] TestPerson{salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1}
9. [1] TestPerson{name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40, status=active}
10. [1] TestPerson{department=intern, level=1, name=X, score=6.5, tags=temp, active=true, age=22, salary=28000, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive}

3. **Token 3**:
   - Fait 1: [1] TestPerson{status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true}

4. **Token 4**:
   - Fait 1: [1] TestPerson{active=true, salary=95000, department=management, age=65, score=10, tags=executive, status=active, level=9, name=Grace}

5. **Token 5**:
   - Fait 1: [1] TestPerson{age=18, status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry}

6. **Token 6**:
   - Fait 1: [1] TestPerson{name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering, age=35, salary=75000}

8. **Token 8**:
   - Fait 1: [1] TestPerson{department=sales, salary=55000, name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee}

9. **Token 9**:
   - Fait 1: [1] TestPerson{age=0, score=0, tags=test, active=true, level=1, name=Frank, status=active, salary=-5000, department=qa}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{brand=TechCorp, stock=50, name=Laptop, rating=4.5, keywords=computer, price=999.99, supplier=TechSupply, category=electronics, available=true}
2. [1] TestProduct{stock=200, supplier=TechSupply, name=Mouse, category=accessories, rating=4, brand=TechCorp, price=25.5, keywords=peripheral, available=true}
3. [1] TestProduct{category=accessories, price=75, supplier=KeySupply, name=Keyboard, keywords=typing, stock=0, rating=3.5, brand=KeyTech, available=false}
4. [1] TestProduct{stock=30, rating=4.8, supplier=ScreenSupply, name=Monitor, category=electronics, price=299.99, available=true, keywords=display, brand=ScreenPro}
5. [1] TestProduct{category=accessories, price=8.5, available=false, rating=2, keywords=obsolete, name=OldKeyboard, brand=OldTech, stock=0, supplier=OldSupply}
6. [1] TestProduct{available=true, rating=4.6, keywords=sound, supplier=AudioSupply, name=Headphones, stock=75, category=audio, price=150, brand=AudioMax}
7. [1] TestProduct{brand=CamTech, price=89.99, rating=3.8, stock=25, supplier=CamSupply, category=electronics, name=Webcam, available=true, keywords=video}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{available=true, keywords=video, brand=CamTech, price=89.99, rating=3.8, stock=25, supplier=CamSupply, category=electronics, name=Webcam}

2. **Token 2**:
   - Fait 1: [1] TestProduct{category=electronics, available=true, brand=TechCorp, stock=50, name=Laptop, rating=4.5, keywords=computer, price=999.99, supplier=TechSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{stock=200, supplier=TechSupply, name=Mouse, category=accessories, rating=4, brand=TechCorp, price=25.5, keywords=peripheral, available=true}

4. **Token 4**:
   - Fait 1: [1] TestProduct{rating=3.5, brand=KeyTech, available=false, category=accessories, price=75, supplier=KeySupply, name=Keyboard, keywords=typing, stock=0}

5. **Token 5**:
   - Fait 1: [1] TestProduct{stock=30, rating=4.8, supplier=ScreenSupply, name=Monitor, category=electronics, price=299.99, available=true, keywords=display, brand=ScreenPro}

6. **Token 6**:
   - Fait 1: [1] TestProduct{name=Headphones, stock=75, category=audio, price=150, brand=AudioMax, available=true, rating=4.6, keywords=sound, supplier=AudioSupply}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales, level=2, salary=45000}
2. [1] TestPerson{tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5, active=true}
3. [1] TestPerson{score=6, age=16, tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false}
4. [1] TestPerson{status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true}
5. [1] TestPerson{department=sales, salary=55000, name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee}
6. [1] TestPerson{name=Frank, status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1}
7. [1] TestPerson{active=true, salary=95000, department=management, age=65, score=10, tags=executive, status=active, level=9, name=Grace}
8. [1] TestPerson{score=5.5, name=Henry, age=18, status=inactive, department=support, level=1, salary=25000, tags=junior, active=false}
9. [1] TestPerson{salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6}
10. [1] TestPerson{age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering, age=35}

3. **Token 3**:
   - Fait 1: [1] TestPerson{status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1, name=Frank}

4. **Token 4**:
   - Fait 1: [1] TestPerson{name=Grace, active=true, salary=95000, department=management, age=65, score=10, tags=executive, status=active, level=9}

5. **Token 5**:
   - Fait 1: [1] TestPerson{score=5.5, name=Henry, age=18, status=inactive, department=support, level=1, salary=25000, tags=junior, active=false}

6. **Token 6**:
   - Fait 1: [1] TestPerson{age=16, tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false, score=6}

7. **Token 7**:
   - Fait 1: [1] TestPerson{age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true, status=active}

8. **Token 8**:
   - Fait 1: [1] TestPerson{salary=55000, name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales}

9. **Token 9**:
   - Fait 1: [1] TestPerson{salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}
2. [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}
3. [1] TestOrder{customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north}
4. [1] TestOrder{customer_id=P004, priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99}
5. [1] TestOrder{discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high}
6. [1] TestOrder{priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0}
7. [1] TestOrder{discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4}
8. [1] TestOrder{priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002}
9. [1] TestOrder{total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1}
10. [1] TestOrder{region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent}

2. **Token 2**:
   - Fait 1: [1] TestOrder{priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002}

3. **Token 3**:
   - Fait 1: [1] TestOrder{priority=low, amount=1, total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007}

4. **Token 4**:
   - Fait 1: [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

5. **Token 5**:
   - Fait 1: [1] TestOrder{amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002, region=south, product_id=PROD002, date=2024-01-20, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15, amount=3}

7. **Token 7**:
   - Fait 1: [1] TestOrder{priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004}

8. **Token 8**:
   - Fait 1: [1] TestOrder{priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{department=sales, level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior}
2. [1] TestPerson{active=true, tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5}
3. [1] TestPerson{status=inactive, name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0}
4. [1] TestPerson{status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true}
5. [1] TestPerson{name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000}
6. [1] TestPerson{score=0, tags=test, active=true, level=1, name=Frank, status=active, salary=-5000, department=qa, age=0}
7. [1] TestPerson{score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management, age=65}
8. [1] TestPerson{department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive}
9. [1] TestPerson{name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40, status=active}
10. [1] TestPerson{department=intern, level=1, name=X, score=6.5, tags=temp, active=true, age=22, salary=28000, status=active}
11. [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}
12. [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}
13. [1] TestOrder{product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped}
14. [1] TestOrder{date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0}
15. [1] TestOrder{amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10}
16. [1] TestOrder{customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low}
17. [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4, discount=50, product_id=PROD006, total=600, status=shipped}
18. [1] TestOrder{discount=0, customer_id=P010, product_id=PROD002, priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10}
19. [1] TestOrder{discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99}
20. [1] TestOrder{discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{active=true, salary=95000, department=management, age=65, score=10, tags=executive, status=active, level=9, name=Grace}
   - Fait 2: [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}

2. **Token 2**:
   - Fait 1: [1] TestPerson{tags=employee, department=sales, salary=55000, name=Eve, age=30, active=false, status=inactive, score=8, level=3}
   - Fait 2: [1] TestOrder{status=shipped, product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15, amount=3, total=225}

3. **Token 3**:
   - Fait 1: [1] TestPerson{tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5, active=true}
   - Fait 2: [1] TestOrder{region=east, total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered}

4. **Token 4**:
   - Fait 1: [1] TestPerson{age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management}
   - Fait 2: [1] TestOrder{region=east, total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered}

5. **Token 5**:
   - Fait 1: [1] TestPerson{name=Frank, status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1}
   - Fait 2: [1] TestOrder{customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low}

6. **Token 6**:
   - Fait 1: [1] TestPerson{tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6, salary=68000}
   - Fait 2: [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4, discount=50, product_id=PROD006, total=600, status=shipped}

7. **Token 7**:
   - Fait 1: [1] TestPerson{active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1, salary=25000, tags=junior}
   - Fait 2: [1] TestOrder{region=north, customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50}

8. **Token 8**:
   - Fait 1: [1] TestPerson{level=1, name=X, score=6.5, tags=temp, active=true, age=22, salary=28000, status=active, department=intern}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10}

9. **Token 9**:
   - Fait 1: [1] TestPerson{salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{customer_id=P007, date=2024-03-01, priority=urgent, amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north}

10. **Token 10**:
   - Fait 1: [1] TestPerson{age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering}
   - Fait 2: [1] TestOrder{product_id=PROD002, priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010}

11. **Token 11**:
   - Fait 1: [1] TestPerson{name=Grace, active=true, salary=95000, department=management, age=65, score=10, tags=executive, status=active, level=9}
   - Fait 2: [1] TestOrder{discount=0, customer_id=P010, product_id=PROD002, priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10}

12. **Token 12**:
   - Fait 1: [1] TestPerson{age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000, name=Eve}
   - Fait 2: [1] TestOrder{date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10}

13. **Token 13**:
   - Fait 1: [1] TestPerson{salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering, age=35}
   - Fait 2: [1] TestOrder{customer_id=P006, total=75000, status=refunded, date=2024-03-15, discount=0, region=east, product_id=PROD001, priority=urgent, amount=1}

14. **Token 14**:
   - Fait 1: [1] TestPerson{active=true, salary=95000, department=management, age=65, score=10, tags=executive, status=active, level=9, name=Grace}
   - Fait 2: [1] TestOrder{region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15, discount=0}

15. **Token 15**:
   - Fait 1: [1] TestPerson{age=18, status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry}
   - Fait 2: [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}

16. **Token 16**:
   - Fait 1: [1] TestPerson{age=0, score=0, tags=test, active=true, level=1, name=Frank, status=active, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north}

17. **Token 17**:
   - Fait 1: [1] TestPerson{salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05}

18. **Token 18**:
   - Fait 1: [1] TestPerson{tags=employee, department=sales, salary=55000, name=Eve, age=30, active=false, status=inactive, score=8, level=3}
   - Fait 2: [1] TestOrder{region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled}

19. **Token 19**:
   - Fait 1: [1] TestPerson{salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{customer_id=P010, product_id=PROD002, priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0}

20. **Token 20**:
   - Fait 1: [1] TestPerson{active=true, level=1, name=Frank, status=active, salary=-5000, department=qa, age=0, score=0, tags=test}
   - Fait 2: [1] TestOrder{date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10}

21. **Token 21**:
   - Fait 1: [1] TestPerson{salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{customer_id=P006, total=75000, status=refunded, date=2024-03-15, discount=0, region=east, product_id=PROD001, priority=urgent, amount=1}

22. **Token 22**:
   - Fait 1: [1] TestPerson{age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management}
   - Fait 2: [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

23. **Token 23**:
   - Fait 1: [1] TestPerson{score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering, age=35, salary=75000}
   - Fait 2: [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

24. **Token 24**:
   - Fait 1: [1] TestPerson{salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6}
   - Fait 2: [1] TestOrder{amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002, region=south, product_id=PROD002, date=2024-01-20, discount=0}

25. **Token 25**:
   - Fait 1: [1] TestPerson{salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6}
   - Fait 2: [1] TestOrder{discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal}

26. **Token 26**:
   - Fait 1: [1] TestPerson{age=16, tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false, score=6}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10}

27. **Token 27**:
   - Fait 1: [1] TestPerson{level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10}

28. **Token 28**:
   - Fait 1: [1] TestPerson{department=engineering, age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob}
   - Fait 2: [1] TestOrder{priority=urgent, amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01}

29. **Token 29**:
   - Fait 1: [1] TestPerson{department=management, age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000}
   - Fait 2: [1] TestOrder{amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent}

30. **Token 30**:
   - Fait 1: [1] TestPerson{score=7.8, department=marketing, name=Diana, active=true, status=active, age=45, level=7, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{date=2024-01-15, status=pending, priority=normal, discount=50, region=north, customer_id=P001, product_id=PROD001, total=1999.98, amount=2}

31. **Token 31**:
   - Fait 1: [1] TestPerson{salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true, age=22}
   - Fait 2: [1] TestOrder{status=shipped, product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15, amount=3, total=225}

32. **Token 32**:
   - Fait 1: [1] TestPerson{level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{product_id=PROD005, date=2024-02-15, discount=0, priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west}

33. **Token 33**:
   - Fait 1: [1] TestPerson{score=6, age=16, tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false}
   - Fait 2: [1] TestOrder{date=2024-02-15, discount=0, priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005}

34. **Token 34**:
   - Fait 1: [1] TestPerson{score=7.8, department=marketing, name=Diana, active=true, status=active, age=45, level=7, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{customer_id=P007, date=2024-03-01, priority=urgent, amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north}

35. **Token 35**:
   - Fait 1: [1] TestPerson{name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40, status=active}
   - Fait 2: [1] TestOrder{product_id=PROD002, priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010}

36. **Token 36**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true}
   - Fait 2: [1] TestOrder{discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99}

37. **Token 37**:
   - Fait 1: [1] TestPerson{salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6}
   - Fait 2: [1] TestOrder{discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15}

38. **Token 38**:
   - Fait 1: [1] TestPerson{name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000}
   - Fait 2: [1] TestOrder{date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001}

39. **Token 39**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true}
   - Fait 2: [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

40. **Token 40**:
   - Fait 1: [1] TestPerson{age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true, status=active}
   - Fait 2: [1] TestOrder{status=shipped, product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15, amount=3, total=225}

41. **Token 41**:
   - Fait 1: [1] TestPerson{salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high}

42. **Token 42**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true}
   - Fait 2: [1] TestOrder{total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4, discount=50, product_id=PROD006}

43. **Token 43**:
   - Fait 1: [1] TestPerson{status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true}
   - Fait 2: [1] TestOrder{region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10, date=2024-03-10}

44. **Token 44**:
   - Fait 1: [1] TestPerson{active=true, name=Alice, age=25, tags=junior, department=sales, level=2, salary=45000, score=8.5, status=active}
   - Fait 2: [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}

45. **Token 45**:
   - Fait 1: [1] TestPerson{name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive}
   - Fait 2: [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}

46. **Token 46**:
   - Fait 1: [1] TestPerson{name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive}
   - Fait 2: [1] TestOrder{priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004}

47. **Token 47**:
   - Fait 1: [1] TestPerson{department=sales, level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior}
   - Fait 2: [1] TestOrder{customer_id=P004, priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99}

48. **Token 48**:
   - Fait 1: [1] TestPerson{salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6}
   - Fait 2: [1] TestOrder{discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high}

49. **Token 49**:
   - Fait 1: [1] TestPerson{active=true, name=Alice, age=25, tags=junior, department=sales, level=2, salary=45000, score=8.5, status=active}
   - Fait 2: [1] TestOrder{priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002}

50. **Token 50**:
   - Fait 1: [1] TestPerson{name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive}
   - Fait 2: [1] TestOrder{status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002, priority=normal, region=south, total=255, date=2024-03-05}

51. **Token 51**:
   - Fait 1: [1] TestPerson{active=false, score=6, age=16, tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1}
   - Fait 2: [1] TestOrder{customer_id=P006, total=75000, status=refunded, date=2024-03-15, discount=0, region=east, product_id=PROD001, priority=urgent, amount=1}

52. **Token 52**:
   - Fait 1: [1] TestPerson{age=0, score=0, tags=test, active=true, level=1, name=Frank, status=active, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{status=pending, priority=normal, discount=50, region=north, customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15}

53. **Token 53**:
   - Fait 1: [1] TestPerson{name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40, status=active}
   - Fait 2: [1] TestOrder{customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low}

54. **Token 54**:
   - Fait 1: [1] TestPerson{tags=test, active=true, level=1, name=Frank, status=active, salary=-5000, department=qa, age=0, score=0}
   - Fait 2: [1] TestOrder{total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4, discount=50, product_id=PROD006}

55. **Token 55**:
   - Fait 1: [1] TestPerson{salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales, level=2}
   - Fait 2: [1] TestOrder{discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15}

56. **Token 56**:
   - Fait 1: [1] TestPerson{level=3, tags=employee, department=sales, salary=55000, name=Eve, age=30, active=false, status=inactive, score=8}
   - Fait 2: [1] TestOrder{total=25.5, status=confirmed, priority=low, customer_id=P002, region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1}

57. **Token 57**:
   - Fait 1: [1] TestPerson{level=9, name=Grace, active=true, salary=95000, department=management, age=65, score=10, tags=executive, status=active}
   - Fait 2: [1] TestOrder{product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped}

58. **Token 58**:
   - Fait 1: [1] TestPerson{level=5, active=true, tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active}
   - Fait 2: [1] TestOrder{priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north, customer_id=P001}

59. **Token 59**:
   - Fait 1: [1] TestPerson{salary=55000, name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales}
   - Fait 2: [1] TestOrder{date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0}

60. **Token 60**:
   - Fait 1: [1] TestPerson{salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low}

61. **Token 61**:
   - Fait 1: [1] TestPerson{name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000}
   - Fait 2: [1] TestOrder{priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002}

62. **Token 62**:
   - Fait 1: [1] TestPerson{status=active, level=5, active=true, tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2}
   - Fait 2: [1] TestOrder{customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10, date=2024-03-10, region=north}

63. **Token 63**:
   - Fait 1: [1] TestPerson{age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management}
   - Fait 2: [1] TestOrder{date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10}

64. **Token 64**:
   - Fait 1: [1] TestPerson{tags=employee, department=sales, salary=55000, name=Eve, age=30, active=false, status=inactive, score=8, level=3}
   - Fait 2: [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

65. **Token 65**:
   - Fait 1: [1] TestPerson{tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6, salary=68000}
   - Fait 2: [1] TestOrder{customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north}

66. **Token 66**:
   - Fait 1: [1] TestPerson{name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000}
   - Fait 2: [1] TestOrder{total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4, discount=50, product_id=PROD006}

67. **Token 67**:
   - Fait 1: [1] TestPerson{salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6}
   - Fait 2: [1] TestOrder{priority=low, amount=1, total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007}

68. **Token 68**:
   - Fait 1: [1] TestPerson{department=sales, salary=55000, name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee}
   - Fait 2: [1] TestOrder{customer_id=P006, total=75000, status=refunded, date=2024-03-15, discount=0, region=east, product_id=PROD001, priority=urgent, amount=1}

69. **Token 69**:
   - Fait 1: [1] TestPerson{age=0, score=0, tags=test, active=true, level=1, name=Frank, status=active, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}

70. **Token 70**:
   - Fait 1: [1] TestPerson{status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18}
   - Fait 2: [1] TestOrder{amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15}

71. **Token 71**:
   - Fait 1: [1] TestPerson{name=Frank, status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1}
   - Fait 2: [1] TestOrder{discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal}

72. **Token 72**:
   - Fait 1: [1] TestPerson{status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10}

73. **Token 73**:
   - Fait 1: [1] TestPerson{age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management}
   - Fait 2: [1] TestOrder{discount=0, priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15}

74. **Token 74**:
   - Fait 1: [1] TestPerson{level=5, active=true, tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active}
   - Fait 2: [1] TestOrder{total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low, customer_id=P005, amount=2}

75. **Token 75**:
   - Fait 1: [1] TestPerson{status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1, name=Frank}
   - Fait 2: [1] TestOrder{priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002}

76. **Token 76**:
   - Fait 1: [1] TestPerson{status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18}
   - Fait 2: [1] TestOrder{amount=1, total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low}

77. **Token 77**:
   - Fait 1: [1] TestPerson{salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern}
   - Fait 2: [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

78. **Token 78**:
   - Fait 1: [1] TestPerson{salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales, level=2}
   - Fait 2: [1] TestOrder{total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, customer_id=P001, product_id=PROD001}

79. **Token 79**:
   - Fait 1: [1] TestPerson{salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering, age=35}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10}

80. **Token 80**:
   - Fait 1: [1] TestPerson{active=true, salary=95000, department=management, age=65, score=10, tags=executive, status=active, level=9, name=Grace}
   - Fait 2: [1] TestOrder{discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high}

81. **Token 81**:
   - Fait 1: [1] TestPerson{tags=manager, score=7.8, department=marketing, name=Diana, active=true, status=active, age=45, level=7, salary=85000}
   - Fait 2: [1] TestOrder{customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low}

82. **Token 82**:
   - Fait 1: [1] TestPerson{name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive}
   - Fait 2: [1] TestOrder{product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4, discount=50}

83. **Token 83**:
   - Fait 1: [1] TestPerson{active=true, name=Alice, age=25, tags=junior, department=sales, level=2, salary=45000, score=8.5, status=active}
   - Fait 2: [1] TestOrder{product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4, discount=50}

84. **Token 84**:
   - Fait 1: [1] TestPerson{name=Frank, status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1}
   - Fait 2: [1] TestOrder{amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15, discount=0, region=east, product_id=PROD001, priority=urgent}

85. **Token 85**:
   - Fait 1: [1] TestPerson{level=1, name=X, score=6.5, tags=temp, active=true, age=22, salary=28000, status=active, department=intern}
   - Fait 2: [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}

86. **Token 86**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true}
   - Fait 2: [1] TestOrder{date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0}

87. **Token 87**:
   - Fait 1: [1] TestPerson{department=intern, level=1, name=X, score=6.5, tags=temp, active=true, age=22, salary=28000, status=active}
   - Fait 2: [1] TestOrder{region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002, priority=normal}

88. **Token 88**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true}
   - Fait 2: [1] TestOrder{discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15}

89. **Token 89**:
   - Fait 1: [1] TestPerson{name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive}
   - Fait 2: [1] TestOrder{product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped}

90. **Token 90**:
   - Fait 1: [1] TestPerson{level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped}

91. **Token 91**:
   - Fait 1: [1] TestPerson{level=1, name=X, score=6.5, tags=temp, active=true, age=22, salary=28000, status=active, department=intern}
   - Fait 2: [1] TestOrder{priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0}

92. **Token 92**:
   - Fait 1: [1] TestPerson{name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive}
   - Fait 2: [1] TestOrder{date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10}

93. **Token 93**:
   - Fait 1: [1] TestPerson{score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales, level=2, salary=45000}
   - Fait 2: [1] TestOrder{date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10}

94. **Token 94**:
   - Fait 1: [1] TestPerson{department=marketing, name=Diana, active=true, status=active, age=45, level=7, salary=85000, tags=manager, score=7.8}
   - Fait 2: [1] TestOrder{total=25.5, status=confirmed, priority=low, customer_id=P002, region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1}

95. **Token 95**:
   - Fait 1: [1] TestPerson{level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true, status=active, age=45}
   - Fait 2: [1] TestOrder{status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1}

96. **Token 96**:
   - Fait 1: [1] TestPerson{salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1, name=Frank, status=active}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10}

97. **Token 97**:
   - Fait 1: [1] TestPerson{tags=manager, score=7.8, department=marketing, name=Diana, active=true, status=active, age=45, level=7, salary=85000}
   - Fait 2: [1] TestOrder{priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002}

98. **Token 98**:
   - Fait 1: [1] TestPerson{level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true, status=active, age=45}
   - Fait 2: [1] TestOrder{region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15, discount=0}

99. **Token 99**:
   - Fait 1: [1] TestPerson{salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6}
   - Fait 2: [1] TestOrder{date=2024-01-15, status=pending, priority=normal, discount=50, region=north, customer_id=P001, product_id=PROD001, total=1999.98, amount=2}

100. **Token 100**:
   - Fait 1: [1] TestPerson{active=true, tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5}
   - Fait 2: [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{name=Alice, age=25, tags=junior, department=sales, level=2, salary=45000, score=8.5, status=active, active=true}
2. [1] TestPerson{department=engineering, age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob}
3. [1] TestPerson{name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive}
4. [1] TestPerson{name=Diana, active=true, status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing}
5. [1] TestPerson{tags=employee, department=sales, salary=55000, name=Eve, age=30, active=false, status=inactive, score=8, level=3}
6. [1] TestPerson{status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1, name=Frank}
7. [1] TestPerson{age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management}
8. [1] TestPerson{name=Henry, age=18, status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5}
9. [1] TestPerson{salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6}
10. [1] TestPerson{status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true, age=22, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{score=7.8, department=marketing, name=Diana, active=true, status=active, age=45, level=7, salary=85000, tags=manager}

2. **Token 2**:
   - Fait 1: [1] TestPerson{active=true, level=1, name=Frank, status=active, salary=-5000, department=qa, age=0, score=0, tags=test}

3. **Token 3**:
   - Fait 1: [1] TestPerson{age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management}

4. **Token 4**:
   - Fait 1: [1] TestPerson{tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6, salary=68000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{active=true, age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp}

6. **Token 6**:
   - Fait 1: [1] TestPerson{active=true, name=Alice, age=25, tags=junior, department=sales, level=2, salary=45000, score=8.5, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering, age=35}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales, level=2}
2. [1] TestPerson{name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior}
3. [1] TestPerson{name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive}
4. [1] TestPerson{status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true}
5. [1] TestPerson{name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000}
6. [1] TestPerson{name=Frank, status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1}
7. [1] TestPerson{age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management}
8. [1] TestPerson{salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1}
9. [1] TestPerson{department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior}
10. [1] TestPerson{score=6.5, tags=temp, active=true, age=22, salary=28000, status=active, department=intern, level=1, name=X}
11. [1] TestOrder{amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north, customer_id=P001, product_id=PROD001, total=1999.98}
12. [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}
13. [1] TestOrder{discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high}
14. [1] TestOrder{product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05}
15. [1] TestOrder{discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high}
16. [1] TestOrder{priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0}
17. [1] TestOrder{amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent}
18. [1] TestOrder{priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002}
19. [1] TestOrder{total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1}
20. [1] TestOrder{discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{age=0, score=0, tags=test, active=true, level=1, name=Frank, status=active, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002, region=south, product_id=PROD002, date=2024-01-20, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{tags=test, active=true, level=1, name=Frank, status=active, salary=-5000, department=qa, age=0, score=0}
   - Fait 2: [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

3. **Token 3**:
   - Fait 1: [1] TestPerson{status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18}
   - Fait 2: [1] TestOrder{product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05}

4. **Token 4**:
   - Fait 1: [1] TestPerson{department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive}
   - Fait 2: [1] TestOrder{product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south}

5. **Token 5**:
   - Fait 1: [1] TestPerson{age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering}
   - Fait 2: [1] TestOrder{date=2024-02-15, discount=0, priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005}

6. **Token 6**:
   - Fait 1: [1] TestPerson{name=Frank, status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1}
   - Fait 2: [1] TestOrder{date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10}

7. **Token 7**:
   - Fait 1: [1] TestPerson{score=6, age=16, tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false}
   - Fait 2: [1] TestOrder{discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15}

8. **Token 8**:
   - Fait 1: [1] TestPerson{name=Grace, active=true, salary=95000, department=management, age=65, score=10, tags=executive, status=active, level=9}
   - Fait 2: [1] TestOrder{customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low}

9. **Token 9**:
   - Fait 1: [1] TestPerson{level=6, salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7}
   - Fait 2: [1] TestOrder{total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low, customer_id=P005, amount=2}

10. **Token 10**:
   - Fait 1: [1] TestPerson{tags=temp, active=true, age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5}
   - Fait 2: [1] TestOrder{product_id=PROD005, date=2024-02-15, discount=0, priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west}

11. **Token 11**:
   - Fait 1: [1] TestPerson{score=7.8, department=marketing, name=Diana, active=true, status=active, age=45, level=7, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low}

12. **Token 12**:
   - Fait 1: [1] TestPerson{age=18, status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry}
   - Fait 2: [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4, discount=50, product_id=PROD006, total=600, status=shipped}

13. **Token 13**:
   - Fait 1: [1] TestPerson{name=Frank, status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1}
   - Fait 2: [1] TestOrder{region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002, priority=normal}

14. **Token 14**:
   - Fait 1: [1] TestPerson{tags=junior, department=sales, level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25}
   - Fait 2: [1] TestOrder{region=north, customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01}

15. **Token 15**:
   - Fait 1: [1] TestPerson{tags=employee, department=sales, salary=55000, name=Eve, age=30, active=false, status=inactive, score=8, level=3}
   - Fait 2: [1] TestOrder{customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north}

16. **Token 16**:
   - Fait 1: [1] TestPerson{status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1, name=Frank}
   - Fait 2: [1] TestOrder{status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4, discount=50, product_id=PROD006, total=600}

17. **Token 17**:
   - Fait 1: [1] TestPerson{salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002}

18. **Token 18**:
   - Fait 1: [1] TestPerson{age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10}

19. **Token 19**:
   - Fait 1: [1] TestPerson{name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40, status=active}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10}

20. **Token 20**:
   - Fait 1: [1] TestPerson{age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true, status=active}
   - Fait 2: [1] TestOrder{product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05}

21. **Token 21**:
   - Fait 1: [1] TestPerson{score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management, age=65}
   - Fait 2: [1] TestOrder{product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05}

22. **Token 22**:
   - Fait 1: [1] TestPerson{name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40, status=active}
   - Fait 2: [1] TestOrder{priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004}

23. **Token 23**:
   - Fait 1: [1] TestPerson{score=6.5, tags=temp, active=true, age=22, salary=28000, status=active, department=intern, level=1, name=X}
   - Fait 2: [1] TestOrder{customer_id=P004, priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99}

24. **Token 24**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true}
   - Fait 2: [1] TestOrder{total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002}

25. **Token 25**:
   - Fait 1: [1] TestPerson{score=7.8, department=marketing, name=Diana, active=true, status=active, age=45, level=7, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed}

26. **Token 26**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true}
   - Fait 2: [1] TestOrder{discount=50, region=north, customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal}

27. **Token 27**:
   - Fait 1: [1] TestPerson{tags=manager, score=7.8, department=marketing, name=Diana, active=true, status=active, age=45, level=7, salary=85000}
   - Fait 2: [1] TestOrder{discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002, region=south, product_id=PROD002, date=2024-01-20}

28. **Token 28**:
   - Fait 1: [1] TestPerson{salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6}
   - Fait 2: [1] TestOrder{total=25.5, status=confirmed, priority=low, customer_id=P002, region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1}

29. **Token 29**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true}
   - Fait 2: [1] TestOrder{customer_id=P002, region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low}

30. **Token 30**:
   - Fait 1: [1] TestPerson{active=true, salary=95000, department=management, age=65, score=10, tags=executive, status=active, level=9, name=Grace}
   - Fait 2: [1] TestOrder{total=25.5, status=confirmed, priority=low, customer_id=P002, region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1}

31. **Token 31**:
   - Fait 1: [1] TestPerson{tags=manager, score=7.8, department=marketing, name=Diana, active=true, status=active, age=45, level=7, salary=85000}
   - Fait 2: [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

32. **Token 32**:
   - Fait 1: [1] TestPerson{tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management, age=65, score=10}
   - Fait 2: [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

33. **Token 33**:
   - Fait 1: [1] TestPerson{name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40, status=active}
   - Fait 2: [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

34. **Token 34**:
   - Fait 1: [1] TestPerson{name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40, status=active}
   - Fait 2: [1] TestOrder{total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15, amount=3}

35. **Token 35**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true}
   - Fait 2: [1] TestOrder{customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north}

36. **Token 36**:
   - Fait 1: [1] TestPerson{score=7.8, department=marketing, name=Diana, active=true, status=active, age=45, level=7, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north}

37. **Token 37**:
   - Fait 1: [1] TestPerson{name=Grace, active=true, salary=95000, department=management, age=65, score=10, tags=executive, status=active, level=9}
   - Fait 2: [1] TestOrder{amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15}

38. **Token 38**:
   - Fait 1: [1] TestPerson{tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5, active=true}
   - Fait 2: [1] TestOrder{customer_id=P010, product_id=PROD002, priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0}

39. **Token 39**:
   - Fait 1: [1] TestPerson{tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false, score=6, age=16}
   - Fait 2: [1] TestOrder{product_id=PROD002, priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010}

40. **Token 40**:
   - Fait 1: [1] TestPerson{name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40, status=active}
   - Fait 2: [1] TestOrder{priority=low, amount=1, total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007}

41. **Token 41**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true}
   - Fait 2: [1] TestOrder{priority=low, amount=1, total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007}

42. **Token 42**:
   - Fait 1: [1] TestPerson{name=Alice, age=25, tags=junior, department=sales, level=2, salary=45000, score=8.5, status=active, active=true}
   - Fait 2: [1] TestOrder{customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0, priority=low}

43. **Token 43**:
   - Fait 1: [1] TestPerson{tags=employee, department=sales, salary=55000, name=Eve, age=30, active=false, status=inactive, score=8, level=3}
   - Fait 2: [1] TestOrder{product_id=PROD005, date=2024-02-15, discount=0, priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west}

44. **Token 44**:
   - Fait 1: [1] TestPerson{status=inactive, name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0}
   - Fait 2: [1] TestOrder{customer_id=P007, date=2024-03-01, priority=urgent, amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north}

45. **Token 45**:
   - Fait 1: [1] TestPerson{level=5, active=true, tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active}
   - Fait 2: [1] TestOrder{amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent}

46. **Token 46**:
   - Fait 1: [1] TestPerson{level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15}

47. **Token 47**:
   - Fait 1: [1] TestPerson{status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000, name=Eve, age=30, active=false}
   - Fait 2: [1] TestOrder{region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15, discount=0}

48. **Token 48**:
   - Fait 1: [1] TestPerson{active=true, name=Alice, age=25, tags=junior, department=sales, level=2, salary=45000, score=8.5, status=active}
   - Fait 2: [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

49. **Token 49**:
   - Fait 1: [1] TestPerson{name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000}
   - Fait 2: [1] TestOrder{priority=normal, discount=50, region=north, customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending}

50. **Token 50**:
   - Fait 1: [1] TestPerson{active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000, name=Eve, age=30}
   - Fait 2: [1] TestOrder{total=25.5, status=confirmed, priority=low, customer_id=P002, region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1}

51. **Token 51**:
   - Fait 1: [1] TestPerson{active=true, name=Alice, age=25, tags=junior, department=sales, level=2, salary=45000, score=8.5, status=active}
   - Fait 2: [1] TestOrder{product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002, region=south}

52. **Token 52**:
   - Fait 1: [1] TestPerson{tags=junior, department=sales, level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25}
   - Fait 2: [1] TestOrder{date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10}

53. **Token 53**:
   - Fait 1: [1] TestPerson{name=Diana, active=true, status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing}
   - Fait 2: [1] TestOrder{date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10}

54. **Token 54**:
   - Fait 1: [1] TestPerson{name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000}
   - Fait 2: [1] TestOrder{priority=low, amount=1, total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007}

55. **Token 55**:
   - Fait 1: [1] TestPerson{age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management}
   - Fait 2: [1] TestOrder{priority=low, amount=1, total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007}

56. **Token 56**:
   - Fait 1: [1] TestPerson{name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000}
   - Fait 2: [1] TestOrder{priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004}

57. **Token 57**:
   - Fait 1: [1] TestPerson{salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales, level=2}
   - Fait 2: [1] TestOrder{amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05, product_id=PROD004}

58. **Token 58**:
   - Fait 1: [1] TestPerson{level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10}

59. **Token 59**:
   - Fait 1: [1] TestPerson{name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001, date=2024-02-10}

60. **Token 60**:
   - Fait 1: [1] TestPerson{score=7.8, department=marketing, name=Diana, active=true, status=active, age=45, level=7, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{status=refunded, date=2024-03-15, discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000}

61. **Token 61**:
   - Fait 1: [1] TestPerson{age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management}
   - Fait 2: [1] TestOrder{discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15}

62. **Token 62**:
   - Fait 1: [1] TestPerson{salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6}
   - Fait 2: [1] TestOrder{customer_id=P006, total=75000, status=refunded, date=2024-03-15, discount=0, region=east, product_id=PROD001, priority=urgent, amount=1}

63. **Token 63**:
   - Fait 1: [1] TestPerson{tags=temp, active=true, age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5}
   - Fait 2: [1] TestOrder{discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15}

64. **Token 64**:
   - Fait 1: [1] TestPerson{name=Frank, status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1}
   - Fait 2: [1] TestOrder{region=north, customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01}

65. **Token 65**:
   - Fait 1: [1] TestPerson{age=16, tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false, score=6}
   - Fait 2: [1] TestOrder{product_id=PROD005, date=2024-02-15, discount=0, priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west}

66. **Token 66**:
   - Fait 1: [1] TestPerson{tags=junior, department=sales, level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25}
   - Fait 2: [1] TestOrder{discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4}

67. **Token 67**:
   - Fait 1: [1] TestPerson{active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000, name=Eve, age=30}
   - Fait 2: [1] TestOrder{amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent}

68. **Token 68**:
   - Fait 1: [1] TestPerson{active=true, tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5}
   - Fait 2: [1] TestOrder{total=75000, status=refunded, date=2024-03-15, discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006}

69. **Token 69**:
   - Fait 1: [1] TestPerson{active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1, salary=25000, tags=junior}
   - Fait 2: [1] TestOrder{date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003}

70. **Token 70**:
   - Fait 1: [1] TestPerson{tags=junior, department=sales, level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25}
   - Fait 2: [1] TestOrder{region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002, priority=normal}

71. **Token 71**:
   - Fait 1: [1] TestPerson{age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000, name=Eve}
   - Fait 2: [1] TestOrder{priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002}

72. **Token 72**:
   - Fait 1: [1] TestPerson{name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive}
   - Fait 2: [1] TestOrder{customer_id=P004, priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99}

73. **Token 73**:
   - Fait 1: [1] TestPerson{name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive}
   - Fait 2: [1] TestOrder{discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high}

74. **Token 74**:
   - Fait 1: [1] TestPerson{salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{status=refunded, date=2024-03-15, discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000}

75. **Token 75**:
   - Fait 1: [1] TestPerson{active=true, tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5}
   - Fait 2: [1] TestOrder{customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north}

76. **Token 76**:
   - Fait 1: [1] TestPerson{tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false, score=6, age=16}
   - Fait 2: [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}

77. **Token 77**:
   - Fait 1: [1] TestPerson{name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive}
   - Fait 2: [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

78. **Token 78**:
   - Fait 1: [1] TestPerson{salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false, score=6, age=16, tags=intern}
   - Fait 2: [1] TestOrder{priority=low, amount=1, total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007}

79. **Token 79**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true}
   - Fait 2: [1] TestOrder{customer_id=P006, total=75000, status=refunded, date=2024-03-15, discount=0, region=east, product_id=PROD001, priority=urgent, amount=1}

80. **Token 80**:
   - Fait 1: [1] TestPerson{name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior}
   - Fait 2: [1] TestOrder{date=2024-01-15, status=pending, priority=normal, discount=50, region=north, customer_id=P001, product_id=PROD001, total=1999.98, amount=2}

81. **Token 81**:
   - Fait 1: [1] TestPerson{salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering, age=35}
   - Fait 2: [1] TestOrder{discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002, region=south, product_id=PROD002, date=2024-01-20}

82. **Token 82**:
   - Fait 1: [1] TestPerson{age=16, tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false, score=6}
   - Fait 2: [1] TestOrder{discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high}

83. **Token 83**:
   - Fait 1: [1] TestPerson{name=Frank, status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1}
   - Fait 2: [1] TestOrder{discount=0, priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15}

84. **Token 84**:
   - Fait 1: [1] TestPerson{age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management}
   - Fait 2: [1] TestOrder{region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002, priority=normal}

85. **Token 85**:
   - Fait 1: [1] TestPerson{name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40, status=active}
   - Fait 2: [1] TestOrder{discount=0, customer_id=P010, product_id=PROD002, priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10}

86. **Token 86**:
   - Fait 1: [1] TestPerson{active=true, age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp}
   - Fait 2: [1] TestOrder{total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002, priority=normal, region=south}

87. **Token 87**:
   - Fait 1: [1] TestPerson{status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true}
   - Fait 2: [1] TestOrder{amount=10, discount=0, customer_id=P010, product_id=PROD002, priority=normal, region=south, total=255, date=2024-03-05, status=pending}

88. **Token 88**:
   - Fait 1: [1] TestPerson{tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5, active=true}
   - Fait 2: [1] TestOrder{product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05}

89. **Token 89**:
   - Fait 1: [1] TestPerson{age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering}
   - Fait 2: [1] TestOrder{date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south, product_id=PROD001}

90. **Token 90**:
   - Fait 1: [1] TestPerson{salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1}
   - Fait 2: [1] TestOrder{priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west, product_id=PROD005, date=2024-02-15, discount=0}

91. **Token 91**:
   - Fait 1: [1] TestPerson{status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana, active=true}
   - Fait 2: [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4, discount=50, product_id=PROD006, total=600, status=shipped}

92. **Token 92**:
   - Fait 1: [1] TestPerson{tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management, age=65, score=10}
   - Fait 2: [1] TestOrder{customer_id=P007, date=2024-03-01, priority=urgent, amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north}

93. **Token 93**:
   - Fait 1: [1] TestPerson{salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6}
   - Fait 2: [1] TestOrder{discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent, amount=4}

94. **Token 94**:
   - Fait 1: [1] TestPerson{name=X, score=6.5, tags=temp, active=true, age=22, salary=28000, status=active, department=intern, level=1}
   - Fait 2: [1] TestOrder{amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent}

95. **Token 95**:
   - Fait 1: [1] TestPerson{active=true, tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2, status=active, level=5}
   - Fait 2: [1] TestOrder{priority=low, amount=1, total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007}

96. **Token 96**:
   - Fait 1: [1] TestPerson{status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18}
   - Fait 2: [1] TestOrder{priority=normal, discount=50, region=north, customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending}

97. **Token 97**:
   - Fait 1: [1] TestPerson{age=18, status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry}
   - Fait 2: [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}

98. **Token 98**:
   - Fait 1: [1] TestPerson{department=qa, age=0, score=0, tags=test, active=true, level=1, name=Frank, status=active, salary=-5000}
   - Fait 2: [1] TestOrder{total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east}

99. **Token 99**:
   - Fait 1: [1] TestPerson{age=0, score=0, tags=test, active=true, level=1, name=Frank, status=active, salary=-5000, department=qa}
   - Fait 2: [1] TestOrder{discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high}

100. **Token 100**:
   - Fait 1: [1] TestPerson{status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18}
   - Fait 2: [1] TestOrder{discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales}
2. [1] TestPerson{status=active, level=5, active=true, tags=senior, name=Bob, department=engineering, age=35, salary=75000, score=9.2}
3. [1] TestPerson{age=16, tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false, score=6}
4. [1] TestPerson{active=true, status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana}
5. [1] TestPerson{name=Eve, age=30, active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000}
6. [1] TestPerson{active=true, level=1, name=Frank, status=active, salary=-5000, department=qa, age=0, score=0, tags=test}
7. [1] TestPerson{department=management, age=65, score=10, tags=executive, status=active, level=9, name=Grace, active=true, salary=95000}
8. [1] TestPerson{active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1, salary=25000, tags=junior}
9. [1] TestPerson{salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy, active=true, score=8.7, level=6}
10. [1] TestPerson{level=1, name=X, score=6.5, tags=temp, active=true, age=22, salary=28000, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering, age=35}

2. **Token 2**:
   - Fait 1: [1] TestPerson{tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1, active=false, score=6, age=16}

3. **Token 3**:
   - Fait 1: [1] TestPerson{score=7.8, department=marketing, name=Diana, active=true, status=active, age=45, level=7, salary=85000, tags=manager}

4. **Token 4**:
   - Fait 1: [1] TestPerson{name=Grace, active=true, salary=95000, department=management, age=65, score=10, tags=executive, status=active, level=9}

5. **Token 5**:
   - Fait 1: [1] TestPerson{age=18, status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5, name=Henry}

6. **Token 6**:
   - Fait 1: [1] TestPerson{salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales, level=2}

7. **Token 7**:
   - Fait 1: [1] TestPerson{status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000, name=Eve, age=30, active=false}

8. **Token 8**:
   - Fait 1: [1] TestPerson{name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40, status=active}

9. **Token 9**:
   - Fait 1: [1] TestPerson{tags=temp, active=true, age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}
2. [1] TestOrder{total=25.5, status=confirmed, priority=low, customer_id=P002, region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1}
3. [1] TestOrder{product_id=PROD003, date=2024-02-01, region=north, customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped}
4. [1] TestOrder{discount=0, date=2024-02-05, product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal}
5. [1] TestOrder{product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high, discount=100, region=south}
6. [1] TestOrder{product_id=PROD005, date=2024-02-15, discount=0, priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west}
7. [1] TestOrder{amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01, priority=urgent}
8. [1] TestOrder{discount=0, customer_id=P010, product_id=PROD002, priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10}
9. [1] TestOrder{date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007, priority=low, amount=1, total=89.99, discount=10}
10. [1] TestOrder{discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{region=south, product_id=PROD002, date=2024-01-20, discount=0, amount=1, total=25.5, status=confirmed, priority=low, customer_id=P002}

2. **Token 2**:
   - Fait 1: [1] TestOrder{product_id=PROD004, amount=1, status=delivered, region=east, total=299.99, customer_id=P004, priority=normal, discount=0, date=2024-02-05}

3. **Token 3**:
   - Fait 1: [1] TestOrder{product_id=PROD005, date=2024-02-15, discount=0, priority=low, customer_id=P005, amount=2, total=999.98, status=cancelled, region=west}

4. **Token 4**:
   - Fait 1: [1] TestOrder{priority=urgent, amount=4, discount=50, product_id=PROD006, total=600, status=shipped, region=north, customer_id=P007, date=2024-03-01}

5. **Token 5**:
   - Fait 1: [1] TestOrder{priority=low, amount=1, total=89.99, discount=10, date=2024-03-10, region=north, customer_id=P001, status=completed, product_id=PROD007}

6. **Token 6**:
   - Fait 1: [1] TestOrder{discount=0, region=east, product_id=PROD001, priority=urgent, amount=1, customer_id=P006, total=75000, status=refunded, date=2024-03-15}

7. **Token 7**:
   - Fait 1: [1] TestOrder{customer_id=P001, product_id=PROD001, total=1999.98, amount=2, date=2024-01-15, status=pending, priority=normal, discount=50, region=north}

8. **Token 8**:
   - Fait 1: [1] TestOrder{customer_id=P001, priority=high, discount=15, amount=3, total=225, status=shipped, product_id=PROD003, date=2024-02-01, region=north}

9. **Token 9**:
   - Fait 1: [1] TestOrder{discount=100, region=south, product_id=PROD001, date=2024-02-10, amount=1, status=confirmed, customer_id=P002, total=999.99, priority=high}

10. **Token 10**:
   - Fait 1: [1] TestOrder{priority=normal, region=south, total=255, date=2024-03-05, status=pending, amount=10, discount=0, customer_id=P010, product_id=PROD002}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{salary=45000, score=8.5, status=active, active=true, name=Alice, age=25, tags=junior, department=sales, level=2}
2. [1] TestPerson{age=35, salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering}
3. [1] TestPerson{active=false, score=6, age=16, tags=intern, salary=0, status=inactive, name=Charlie, department=hr, level=1}
4. [1] TestPerson{name=Diana, active=true, status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing}
5. [1] TestPerson{active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000, name=Eve, age=30}
6. [1] TestPerson{tags=test, active=true, level=1, name=Frank, status=active, salary=-5000, department=qa, age=0, score=0}
7. [1] TestPerson{tags=executive, status=active, level=9, name=Grace, active=true, salary=95000, department=management, age=65, score=10}
8. [1] TestPerson{name=Henry, age=18, status=inactive, department=support, level=1, salary=25000, tags=junior, active=false, score=5.5}
9. [1] TestPerson{active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40, status=active, name=Ivy}
10. [1] TestPerson{salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true, age=22}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{tags=junior, department=sales, level=2, salary=45000, score=8.5, status=active, active=true, name=Alice, age=25}

2. **Token 2**:
   - Fait 1: [1] TestPerson{salary=75000, score=9.2, status=active, level=5, active=true, tags=senior, name=Bob, department=engineering, age=35}

3. **Token 3**:
   - Fait 1: [1] TestPerson{name=Ivy, active=true, score=8.7, level=6, salary=68000, tags=senior, department=engineering, age=40, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, status=active, department=intern, level=1, name=X, score=6.5, tags=temp, active=true}

5. **Token 5**:
   - Fait 1: [1] TestPerson{department=hr, level=1, active=false, score=6, age=16, tags=intern, salary=0, status=inactive, name=Charlie}

6. **Token 6**:
   - Fait 1: [1] TestPerson{active=true, status=active, age=45, level=7, salary=85000, tags=manager, score=7.8, department=marketing, name=Diana}

7. **Token 7**:
   - Fait 1: [1] TestPerson{active=false, status=inactive, score=8, level=3, tags=employee, department=sales, salary=55000, name=Eve, age=30}

8. **Token 8**:
   - Fait 1: [1] TestPerson{name=Frank, status=active, salary=-5000, department=qa, age=0, score=0, tags=test, active=true, level=1}

9. **Token 9**:
   - Fait 1: [1] TestPerson{status=active, level=9, name=Grace, active=true, salary=95000, department=management, age=65, score=10, tags=executive}

10. **Token 10**:
   - Fait 1: [1] TestPerson{salary=25000, tags=junior, active=false, score=5.5, name=Henry, age=18, status=inactive, department=support, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
