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

1. [1] TestPerson{active=true, department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000}
2. [1] TestPerson{level=5, status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000}
3. [1] TestPerson{level=1, score=6, status=inactive, age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr}
4. [1] TestPerson{status=active, department=marketing, salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager}
5. [1] TestPerson{tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales}
6. [1] TestPerson{active=true, level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0}
7. [1] TestPerson{department=management, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, active=true}
8. [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
9. [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}
10. [1] TestPerson{status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{level=1, status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern}

3. **Token 3**:
   - Fait 1: [1] TestPerson{age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive}

4. **Token 4**:
   - Fait 1: [1] TestPerson{name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, active=true, department=management, level=9}

5. **Token 5**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales}

7. **Token 7**:
   - Fait 1: [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}

8. **Token 8**:
   - Fait 1: [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing}

9. **Token 9**:
   - Fait 1: [1] TestPerson{tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal}
2. [1] TestOrder{product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0}
3. [1] TestOrder{priority=high, product_id=PROD003, date=2024-02-01, status=shipped, discount=15, region=north, amount=3, customer_id=P001, total=225}
4. [1] TestOrder{priority=normal, region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99}
5. [1] TestOrder{amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001}
6. [1] TestOrder{customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15}
7. [1] TestOrder{discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01}
8. [1] TestOrder{product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010}
9. [1] TestOrder{total=89.99, region=north, customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007}
10. [1] TestOrder{status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east, customer_id=P006, product_id=PROD001, total=75000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01}

2. **Token 2**:
   - Fait 1: [1] TestOrder{product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010}

3. **Token 3**:
   - Fait 1: [1] TestOrder{status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east, customer_id=P006, product_id=PROD001, total=75000}

4. **Token 4**:
   - Fait 1: [1] TestOrder{customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped, discount=15, region=north, amount=3}

5. **Token 5**:
   - Fait 1: [1] TestOrder{amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004}

6. **Token 6**:
   - Fait 1: [1] TestOrder{amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low}

7. **Token 7**:
   - Fait 1: [1] TestOrder{product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal}

8. **Token 8**:
   - Fait 1: [1] TestOrder{product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0}

9. **Token 9**:
   - Fait 1: [1] TestOrder{region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{level=2, salary=45000, active=true, department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5}
2. [1] TestPerson{department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5, status=active}
3. [1] TestPerson{score=6, status=inactive, age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1}
4. [1] TestPerson{department=marketing, salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active}
5. [1] TestPerson{status=inactive, level=3, name=Eve, department=sales, tags=employee, age=30, salary=55000, active=false, score=8}
6. [1] TestPerson{department=qa, age=0, active=true, level=1, name=Frank, tags=test, salary=-5000, score=0, status=active}
7. [1] TestPerson{salary=95000, score=10, tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65}
8. [1] TestPerson{score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000, name=Henry, active=false}
9. [1] TestPerson{status=active, score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000}
10. [1] TestPerson{score=6.5, tags=temp, age=22, salary=28000, department=intern, level=1, status=active, name=X, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5, status=active}

2. **Token 2**:
   - Fait 1: [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing}

3. **Token 3**:
   - Fait 1: [1] TestPerson{status=inactive, level=3, name=Eve, department=sales, tags=employee, age=30, salary=55000, active=false, score=8}

4. **Token 4**:
   - Fait 1: [1] TestPerson{department=management, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, active=true}

5. **Token 5**:
   - Fait 1: [1] TestPerson{status=active, score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001}
2. [1] TestOrder{date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed}
3. [1] TestOrder{product_id=PROD003, date=2024-02-01, status=shipped, discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high}
4. [1] TestOrder{date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004}
5. [1] TestOrder{customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100}
6. [1] TestOrder{discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005}
7. [1] TestOrder{discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01}
8. [1] TestOrder{date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0}
9. [1] TestOrder{amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low}
10. [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{priority=low, amount=1, region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5}

2. **Token 2**:
   - Fait 1: [1] TestOrder{discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped}

3. **Token 3**:
   - Fait 1: [1] TestOrder{priority=normal, region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99}

4. **Token 4**:
   - Fait 1: [1] TestOrder{priority=low, amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001}

5. **Token 5**:
   - Fait 1: [1] TestOrder{amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001}

6. **Token 6**:
   - Fait 1: [1] TestOrder{product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south}

7. **Token 7**:
   - Fait 1: [1] TestOrder{customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15}

8. **Token 8**:
   - Fait 1: [1] TestOrder{priority=urgent, region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600, status=shipped}

9. **Token 9**:
   - Fait 1: [1] TestOrder{amount=10, status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{price=999.99, keywords=computer, brand=TechCorp, supplier=TechSupply, name=Laptop, category=electronics, available=true, rating=4.5, stock=50}
2. [1] TestProduct{price=25.5, available=true, keywords=peripheral, supplier=TechSupply, stock=200, category=accessories, rating=4, name=Mouse, brand=TechCorp}
3. [1] TestProduct{keywords=typing, stock=0, price=75, category=accessories, supplier=KeySupply, brand=KeyTech, name=Keyboard, available=false, rating=3.5}
4. [1] TestProduct{stock=30, price=299.99, available=true, rating=4.8, keywords=display, brand=ScreenPro, category=electronics, supplier=ScreenSupply, name=Monitor}
5. [1] TestProduct{stock=0, name=OldKeyboard, available=false, category=accessories, rating=2, keywords=obsolete, brand=OldTech, supplier=OldSupply, price=8.5}
6. [1] TestProduct{stock=75, available=true, rating=4.6, name=Headphones, brand=AudioMax, keywords=sound, supplier=AudioSupply, category=audio, price=150}
7. [1] TestProduct{category=electronics, brand=CamTech, name=Webcam, price=89.99, supplier=CamSupply, available=true, rating=3.8, stock=25, keywords=video}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{available=false, rating=3.5, keywords=typing, stock=0, price=75, category=accessories, supplier=KeySupply, brand=KeyTech, name=Keyboard}

2. **Token 2**:
   - Fait 1: [1] TestProduct{name=Monitor, stock=30, price=299.99, available=true, rating=4.8, keywords=display, brand=ScreenPro, category=electronics, supplier=ScreenSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{name=Headphones, brand=AudioMax, keywords=sound, supplier=AudioSupply, category=audio, price=150, stock=75, available=true, rating=4.6}

4. **Token 4**:
   - Fait 1: [1] TestProduct{supplier=CamSupply, available=true, rating=3.8, stock=25, keywords=video, category=electronics, brand=CamTech, name=Webcam, price=89.99}

5. **Token 5**:
   - Fait 1: [1] TestProduct{name=Laptop, category=electronics, available=true, rating=4.5, stock=50, price=999.99, keywords=computer, brand=TechCorp, supplier=TechSupply}

6. **Token 6**:
   - Fait 1: [1] TestProduct{rating=4, name=Mouse, brand=TechCorp, price=25.5, available=true, keywords=peripheral, supplier=TechSupply, stock=200, category=accessories}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales}
2. [1] TestPerson{name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5, status=active, department=engineering}
3. [1] TestPerson{name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive, age=16, active=false}
4. [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing}
5. [1] TestPerson{tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales}
6. [1] TestPerson{salary=-5000, score=0, status=active, department=qa, age=0, active=true, level=1, name=Frank, tags=test}
7. [1] TestPerson{tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10}
8. [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
9. [1] TestPerson{tags=senior, level=6, active=true, salary=68000, status=active, score=8.7, department=engineering, name=Ivy, age=40}
10. [1] TestPerson{tags=temp, age=22, salary=28000, department=intern, level=1, status=active, name=X, active=true, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{score=6, status=inactive, age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north}
2. [1] TestOrder{date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed}
3. [1] TestOrder{total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped, discount=15, region=north, amount=3, customer_id=P001}
4. [1] TestOrder{total=299.99, priority=normal, region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0}
5. [1] TestOrder{region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10}
6. [1] TestOrder{product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled}
7. [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}
8. [1] TestOrder{total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10, status=pending, region=south}
9. [1] TestOrder{discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed}
10. [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002}

2. **Token 2**:
   - Fait 1: [1] TestOrder{discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005}

3. **Token 3**:
   - Fait 1: [1] TestOrder{status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east, customer_id=P006, product_id=PROD001, total=75000}

4. **Token 4**:
   - Fait 1: [1] TestOrder{amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001}

5. **Token 5**:
   - Fait 1: [1] TestOrder{customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal}

7. **Token 7**:
   - Fait 1: [1] TestOrder{status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600}

8. **Token 8**:
   - Fait 1: [1] TestOrder{priority=normal, discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10, status=pending, region=south, total=255}

9. **Token 9**:
   - Fait 1: [1] TestOrder{amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low}

10. **Token 10**:
   - Fait 1: [1] TestOrder{discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales, status=active, name=Alice, age=25}
2. [1] TestPerson{level=5, status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000}
3. [1] TestPerson{level=1, score=6, status=inactive, age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr}
4. [1] TestPerson{name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7}
5. [1] TestPerson{tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales}
6. [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}
7. [1] TestPerson{name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, active=true, department=management, level=9}
8. [1] TestPerson{department=support, level=1, salary=25000, name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18}
9. [1] TestPerson{name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active, score=8.7, department=engineering}
10. [1] TestPerson{department=intern, level=1, status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{score=8, status=inactive, level=3, name=Eve, department=sales, tags=employee, age=30, salary=55000, active=false}

2. **Token 2**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}

3. **Token 3**:
   - Fait 1: [1] TestPerson{level=1, score=6, status=inactive, age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15}
2. [1] TestOrder{amount=1, region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low}
3. [1] TestOrder{discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped}
4. [1] TestOrder{region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal}
5. [1] TestOrder{status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002}
6. [1] TestOrder{product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled}
7. [1] TestOrder{discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01}
8. [1] TestOrder{product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010}
9. [1] TestOrder{product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10}
10. [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled}

2. **Token 2**:
   - Fait 1: [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}

3. **Token 3**:
   - Fait 1: [1] TestOrder{customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north}

4. **Token 4**:
   - Fait 1: [1] TestOrder{amount=1, date=2024-03-15, discount=0, region=east, customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent}

5. **Token 5**:
   - Fait 1: [1] TestOrder{product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped, discount=15, region=north, amount=3}

7. **Token 7**:
   - Fait 1: [1] TestOrder{region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal}

8. **Token 8**:
   - Fait 1: [1] TestOrder{customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales}
2. [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
3. [1] TestPerson{level=1, score=6, status=inactive, age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr}
4. [1] TestPerson{status=active, department=marketing, salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager}
5. [1] TestPerson{score=8, status=inactive, level=3, name=Eve, department=sales, tags=employee, age=30, salary=55000, active=false}
6. [1] TestPerson{active=true, level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0}
7. [1] TestPerson{tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10}
8. [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
9. [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}
10. [1] TestPerson{age=22, salary=28000, department=intern, level=1, status=active, name=X, active=true, score=6.5, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{status=active, department=qa, age=0, active=true, level=1, name=Frank, tags=test, salary=-5000, score=0}

2. **Token 2**:
   - Fait 1: [1] TestPerson{status=active, active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive}

3. **Token 3**:
   - Fait 1: [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}

4. **Token 4**:
   - Fait 1: [1] TestPerson{age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales, status=active, name=Alice}

5. **Token 5**:
   - Fait 1: [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{score=6, status=inactive, age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}

8. **Token 8**:
   - Fait 1: [1] TestPerson{age=45, active=true, score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana}

9. **Token 9**:
   - Fait 1: [1] TestPerson{score=8, status=inactive, level=3, name=Eve, department=sales, tags=employee, age=30, salary=55000, active=false}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{price=999.99, keywords=computer, brand=TechCorp, supplier=TechSupply, name=Laptop, category=electronics, available=true, rating=4.5, stock=50}
2. [1] TestProduct{keywords=peripheral, supplier=TechSupply, stock=200, category=accessories, rating=4, name=Mouse, brand=TechCorp, price=25.5, available=true}
3. [1] TestProduct{name=Keyboard, available=false, rating=3.5, keywords=typing, stock=0, price=75, category=accessories, supplier=KeySupply, brand=KeyTech}
4. [1] TestProduct{available=true, rating=4.8, keywords=display, brand=ScreenPro, category=electronics, supplier=ScreenSupply, name=Monitor, stock=30, price=299.99}
5. [1] TestProduct{rating=2, keywords=obsolete, brand=OldTech, supplier=OldSupply, price=8.5, stock=0, name=OldKeyboard, available=false, category=accessories}
6. [1] TestProduct{keywords=sound, supplier=AudioSupply, category=audio, price=150, stock=75, available=true, rating=4.6, name=Headphones, brand=AudioMax}
7. [1] TestProduct{stock=25, keywords=video, category=electronics, brand=CamTech, name=Webcam, price=89.99, supplier=CamSupply, available=true, rating=3.8}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{rating=4, name=Mouse, brand=TechCorp, price=25.5, available=true, keywords=peripheral, supplier=TechSupply, stock=200, category=accessories}

2. **Token 2**:
   - Fait 1: [1] TestProduct{supplier=KeySupply, brand=KeyTech, name=Keyboard, available=false, rating=3.5, keywords=typing, stock=0, price=75, category=accessories}

3. **Token 3**:
   - Fait 1: [1] TestProduct{supplier=ScreenSupply, name=Monitor, stock=30, price=299.99, available=true, rating=4.8, keywords=display, brand=ScreenPro, category=electronics}

4. **Token 4**:
   - Fait 1: [1] TestProduct{rating=4.6, name=Headphones, brand=AudioMax, keywords=sound, supplier=AudioSupply, category=audio, price=150, stock=75, available=true}

5. **Token 5**:
   - Fait 1: [1] TestProduct{brand=CamTech, name=Webcam, price=89.99, supplier=CamSupply, available=true, rating=3.8, stock=25, keywords=video, category=electronics}

6. **Token 6**:
   - Fait 1: [1] TestProduct{category=electronics, available=true, rating=4.5, stock=50, price=999.99, keywords=computer, brand=TechCorp, supplier=TechSupply, name=Laptop}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales}
2. [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
3. [1] TestPerson{score=6, status=inactive, age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1}
4. [1] TestPerson{score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45, active=true}
5. [1] TestPerson{age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales, tags=employee}
6. [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}
7. [1] TestPerson{score=10, tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65, salary=95000}
8. [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
9. [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}
10. [1] TestPerson{status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{score=9.2, tags=senior, salary=75000, level=5, status=active, department=engineering, name=Bob, age=35, active=true}

3. **Token 3**:
   - Fait 1: [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing}

4. **Token 4**:
   - Fait 1: [1] TestPerson{department=sales, tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve}

5. **Token 5**:
   - Fait 1: [1] TestPerson{active=true, level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0}

6. **Token 6**:
   - Fait 1: [1] TestPerson{name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive, age=16, active=false}

7. **Token 7**:
   - Fait 1: [1] TestPerson{name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, active=true, department=management, level=9}

8. **Token 8**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}

9. **Token 9**:
   - Fait 1: [1] TestPerson{name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active, score=8.7, department=engineering}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15}
2. [1] TestOrder{total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002}
3. [1] TestOrder{date=2024-02-01, status=shipped, discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003}
4. [1] TestOrder{region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal}
5. [1] TestOrder{total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed, priority=high}
6. [1] TestOrder{date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2}
7. [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}
8. [1] TestOrder{date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0}
9. [1] TestOrder{date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low, amount=1}
10. [1] TestOrder{region=east, customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped}

2. **Token 2**:
   - Fait 1: [1] TestOrder{date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004}

3. **Token 3**:
   - Fait 1: [1] TestOrder{date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99}

4. **Token 4**:
   - Fait 1: [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}

5. **Token 5**:
   - Fait 1: [1] TestOrder{region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10, status=pending}

6. **Token 6**:
   - Fait 1: [1] TestOrder{amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low}

7. **Token 7**:
   - Fait 1: [1] TestOrder{date=2024-01-15, discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2}

8. **Token 8**:
   - Fait 1: [1] TestOrder{status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{level=2, salary=45000, active=true, department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5}
2. [1] TestPerson{age=35, active=true, score=9.2, tags=senior, salary=75000, level=5, status=active, department=engineering, name=Bob}
3. [1] TestPerson{salary=0, tags=intern, department=hr, level=1, score=6, status=inactive, age=16, active=false, name=Charlie}
4. [1] TestPerson{age=45, active=true, score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana}
5. [1] TestPerson{age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales, tags=employee}
6. [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}
7. [1] TestPerson{active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active}
8. [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
9. [1] TestPerson{salary=68000, status=active, score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true}
10. [1] TestPerson{salary=28000, department=intern, level=1, status=active, name=X, active=true, score=6.5, tags=temp, age=22}
11. [1] TestOrder{amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001}
12. [1] TestOrder{date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed}
13. [1] TestOrder{discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped}
14. [1] TestOrder{region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal}
15. [1] TestOrder{status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002}
16. [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low}
17. [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}
18. [1] TestOrder{status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10}
19. [1] TestOrder{product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10}
20. [1] TestOrder{amount=1, date=2024-03-15, discount=0, region=east, customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped}

2. **Token 2**:
   - Fait 1: [1] TestPerson{salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales, tags=employee, age=30}
   - Fait 2: [1] TestOrder{date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004}

3. **Token 3**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}
   - Fait 2: [1] TestOrder{discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1}

4. **Token 4**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}
   - Fait 2: [1] TestOrder{product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled}

5. **Token 5**:
   - Fait 1: [1] TestPerson{tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low, amount=1, date=2024-03-10}

6. **Token 6**:
   - Fait 1: [1] TestPerson{active=false, score=8, status=inactive, level=3, name=Eve, department=sales, tags=employee, age=30, salary=55000}
   - Fait 2: [1] TestOrder{region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50}

7. **Token 7**:
   - Fait 1: [1] TestPerson{name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5, status=active, department=engineering}
   - Fait 2: [1] TestOrder{date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed}

8. **Token 8**:
   - Fait 1: [1] TestPerson{tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01}

9. **Token 9**:
   - Fait 1: [1] TestPerson{age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001}

10. **Token 10**:
   - Fait 1: [1] TestPerson{department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true}
   - Fait 2: [1] TestOrder{amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001}

11. **Token 11**:
   - Fait 1: [1] TestPerson{score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000, name=Henry, active=false}
   - Fait 2: [1] TestOrder{date=2024-02-01, status=shipped, discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003}

12. **Token 12**:
   - Fait 1: [1] TestPerson{tags=junior, status=inactive, age=18, department=support, level=1, salary=25000, name=Henry, active=false, score=5.5}
   - Fait 2: [1] TestOrder{product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, status=delivered}

13. **Token 13**:
   - Fait 1: [1] TestPerson{status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern, level=1}
   - Fait 2: [1] TestOrder{priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed}

14. **Token 14**:
   - Fait 1: [1] TestPerson{status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern, level=1}
   - Fait 2: [1] TestOrder{product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled}

15. **Token 15**:
   - Fait 1: [1] TestPerson{age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive}
   - Fait 2: [1] TestOrder{customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, region=north}

16. **Token 16**:
   - Fait 1: [1] TestPerson{status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales}
   - Fait 2: [1] TestOrder{product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01, discount=50}

17. **Token 17**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
   - Fait 2: [1] TestOrder{date=2024-01-15, discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2}

18. **Token 18**:
   - Fait 1: [1] TestPerson{age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales, status=active, name=Alice}
   - Fait 2: [1] TestOrder{discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped}

19. **Token 19**:
   - Fait 1: [1] TestPerson{age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive}
   - Fait 2: [1] TestOrder{status=shipped, discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01}

20. **Token 20**:
   - Fait 1: [1] TestPerson{name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales, status=active}
   - Fait 2: [1] TestOrder{amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004}

21. **Token 21**:
   - Fait 1: [1] TestPerson{salary=0, tags=intern, department=hr, level=1, score=6, status=inactive, age=16, active=false, name=Charlie}
   - Fait 2: [1] TestOrder{region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal}

22. **Token 22**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
   - Fait 2: [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}

23. **Token 23**:
   - Fait 1: [1] TestPerson{score=9.2, tags=senior, salary=75000, level=5, status=active, department=engineering, name=Bob, age=35, active=true}
   - Fait 2: [1] TestOrder{amount=10, status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002}

24. **Token 24**:
   - Fait 1: [1] TestPerson{name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive, age=16, active=false}
   - Fait 2: [1] TestOrder{amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low}

25. **Token 25**:
   - Fait 1: [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing}
   - Fait 2: [1] TestOrder{customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north}

26. **Token 26**:
   - Fait 1: [1] TestPerson{status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern, level=1}
   - Fait 2: [1] TestOrder{date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed}

27. **Token 27**:
   - Fait 1: [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}
   - Fait 2: [1] TestOrder{region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped, discount=15}

28. **Token 28**:
   - Fait 1: [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}
   - Fait 2: [1] TestOrder{region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal}

29. **Token 29**:
   - Fait 1: [1] TestPerson{score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45, active=true}
   - Fait 2: [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}

30. **Token 30**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}
   - Fait 2: [1] TestOrder{total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10, status=pending, region=south}

31. **Token 31**:
   - Fait 1: [1] TestPerson{level=2, salary=45000, active=true, department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5}
   - Fait 2: [1] TestOrder{region=north, customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99}

32. **Token 32**:
   - Fait 1: [1] TestPerson{level=6, active=true, salary=68000, status=active, score=8.7, department=engineering, name=Ivy, age=40, tags=senior}
   - Fait 2: [1] TestOrder{customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north}

33. **Token 33**:
   - Fait 1: [1] TestPerson{department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active, score=8.7}
   - Fait 2: [1] TestOrder{status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001, total=1999.98}

34. **Token 34**:
   - Fait 1: [1] TestPerson{age=45, active=true, score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana}
   - Fait 2: [1] TestOrder{discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped}

35. **Token 35**:
   - Fait 1: [1] TestPerson{score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45, active=true}
   - Fait 2: [1] TestOrder{region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal}

36. **Token 36**:
   - Fait 1: [1] TestPerson{active=true, salary=68000, status=active, score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6}
   - Fait 2: [1] TestOrder{discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01}

37. **Token 37**:
   - Fait 1: [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing}
   - Fait 2: [1] TestOrder{amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low}

38. **Token 38**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
   - Fait 2: [1] TestOrder{amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low}

39. **Token 39**:
   - Fait 1: [1] TestPerson{score=0, status=active, department=qa, age=0, active=true, level=1, name=Frank, tags=test, salary=-5000}
   - Fait 2: [1] TestOrder{date=2024-03-15, discount=0, region=east, customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1}

40. **Token 40**:
   - Fait 1: [1] TestPerson{tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10}
   - Fait 2: [1] TestOrder{region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1}

41. **Token 41**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}
   - Fait 2: [1] TestOrder{discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped}

42. **Token 42**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}
   - Fait 2: [1] TestOrder{region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal}

43. **Token 43**:
   - Fait 1: [1] TestPerson{department=sales, tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve}
   - Fait 2: [1] TestOrder{status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002}

44. **Token 44**:
   - Fait 1: [1] TestPerson{tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled}

45. **Token 45**:
   - Fait 1: [1] TestPerson{tags=senior, level=6, active=true, salary=68000, status=active, score=8.7, department=engineering, name=Ivy, age=40}
   - Fait 2: [1] TestOrder{date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0}

46. **Token 46**:
   - Fait 1: [1] TestPerson{status=active, department=qa, age=0, active=true, level=1, name=Frank, tags=test, salary=-5000, score=0}
   - Fait 2: [1] TestOrder{date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low, amount=1}

47. **Token 47**:
   - Fait 1: [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
   - Fait 2: [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

48. **Token 48**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}
   - Fait 2: [1] TestOrder{region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50}

49. **Token 49**:
   - Fait 1: [1] TestPerson{active=true, level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0}
   - Fait 2: [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}

50. **Token 50**:
   - Fait 1: [1] TestPerson{score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45, active=true}
   - Fait 2: [1] TestOrder{product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010}

51. **Token 51**:
   - Fait 1: [1] TestPerson{name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active, score=8.7, department=engineering}
   - Fait 2: [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

52. **Token 52**:
   - Fait 1: [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing}
   - Fait 2: [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

53. **Token 53**:
   - Fait 1: [1] TestPerson{level=5, status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000}
   - Fait 2: [1] TestOrder{discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15}

54. **Token 54**:
   - Fait 1: [1] TestPerson{tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1}

55. **Token 55**:
   - Fait 1: [1] TestPerson{salary=95000, score=10, tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65}
   - Fait 2: [1] TestOrder{date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99}

56. **Token 56**:
   - Fait 1: [1] TestPerson{tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10}
   - Fait 2: [1] TestOrder{customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15}

57. **Token 57**:
   - Fait 1: [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
   - Fait 2: [1] TestOrder{discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01}

58. **Token 58**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
   - Fait 2: [1] TestOrder{date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0}

59. **Token 59**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
   - Fait 2: [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

60. **Token 60**:
   - Fait 1: [1] TestPerson{level=5, status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000}
   - Fait 2: [1] TestOrder{priority=high, product_id=PROD003, date=2024-02-01, status=shipped, discount=15, region=north, amount=3, customer_id=P001, total=225}

61. **Token 61**:
   - Fait 1: [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
   - Fait 2: [1] TestOrder{amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004}

62. **Token 62**:
   - Fait 1: [1] TestPerson{name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive, age=16, active=false}
   - Fait 2: [1] TestOrder{amount=10, status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002}

63. **Token 63**:
   - Fait 1: [1] TestPerson{level=2, salary=45000, active=true, department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5}
   - Fait 2: [1] TestOrder{product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010}

64. **Token 64**:
   - Fait 1: [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
   - Fait 2: [1] TestOrder{region=north, customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99}

65. **Token 65**:
   - Fait 1: [1] TestPerson{level=2, salary=45000, active=true, department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5}
   - Fait 2: [1] TestOrder{priority=urgent, amount=1, date=2024-03-15, discount=0, region=east, customer_id=P006, product_id=PROD001, total=75000, status=refunded}

66. **Token 66**:
   - Fait 1: [1] TestPerson{salary=0, tags=intern, department=hr, level=1, score=6, status=inactive, age=16, active=false, name=Charlie}
   - Fait 2: [1] TestOrder{date=2024-03-15, discount=0, region=east, customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1}

67. **Token 67**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
   - Fait 2: [1] TestOrder{customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south}

68. **Token 68**:
   - Fait 1: [1] TestPerson{level=6, active=true, salary=68000, status=active, score=8.7, department=engineering, name=Ivy, age=40, tags=senior}
   - Fait 2: [1] TestOrder{customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100}

69. **Token 69**:
   - Fait 1: [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}
   - Fait 2: [1] TestOrder{customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15}

70. **Token 70**:
   - Fait 1: [1] TestPerson{department=sales, tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve}
   - Fait 2: [1] TestOrder{status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10}

71. **Token 71**:
   - Fait 1: [1] TestPerson{tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

72. **Token 72**:
   - Fait 1: [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing}
   - Fait 2: [1] TestOrder{customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15}

73. **Token 73**:
   - Fait 1: [1] TestPerson{status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales}
   - Fait 2: [1] TestOrder{date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed}

74. **Token 74**:
   - Fait 1: [1] TestPerson{age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive}
   - Fait 2: [1] TestOrder{product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0}

75. **Token 75**:
   - Fait 1: [1] TestPerson{score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45, active=true}
   - Fait 2: [1] TestOrder{region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10}

76. **Token 76**:
   - Fait 1: [1] TestPerson{department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active, score=8.7}
   - Fait 2: [1] TestOrder{customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south}

77. **Token 77**:
   - Fait 1: [1] TestPerson{salary=28000, department=intern, level=1, status=active, name=X, active=true, score=6.5, tags=temp, age=22}
   - Fait 2: [1] TestOrder{date=2024-02-01, status=shipped, discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003}

78. **Token 78**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, department=intern, level=1, status=active, name=X, active=true, score=6.5, tags=temp}
   - Fait 2: [1] TestOrder{date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004}

79. **Token 79**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
   - Fait 2: [1] TestOrder{status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002}

80. **Token 80**:
   - Fait 1: [1] TestPerson{age=18, department=support, level=1, salary=25000, name=Henry, active=false, score=5.5, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low}

81. **Token 81**:
   - Fait 1: [1] TestPerson{level=9, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, active=true, department=management}
   - Fait 2: [1] TestOrder{total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10, status=pending, region=south}

82. **Token 82**:
   - Fait 1: [1] TestPerson{tags=temp, age=22, salary=28000, department=intern, level=1, status=active, name=X, active=true, score=6.5}
   - Fait 2: [1] TestOrder{priority=low, amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001}

83. **Token 83**:
   - Fait 1: [1] TestPerson{tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10}
   - Fait 2: [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

84. **Token 84**:
   - Fait 1: [1] TestPerson{active=true, department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000}
   - Fait 2: [1] TestOrder{product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled}

85. **Token 85**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, department=intern, level=1, status=active, name=X, active=true, score=6.5, tags=temp}
   - Fait 2: [1] TestOrder{amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001}

86. **Token 86**:
   - Fait 1: [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing}
   - Fait 2: [1] TestOrder{date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed}

87. **Token 87**:
   - Fait 1: [1] TestPerson{status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales}
   - Fait 2: [1] TestOrder{amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001}

88. **Token 88**:
   - Fait 1: [1] TestPerson{level=1, score=6, status=inactive, age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr}
   - Fait 2: [1] TestOrder{region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10}

89. **Token 89**:
   - Fait 1: [1] TestPerson{tags=intern, department=hr, level=1, score=6, status=inactive, age=16, active=false, name=Charlie, salary=0}
   - Fait 2: [1] TestOrder{priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005}

90. **Token 90**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, department=intern, level=1, status=active, name=X, active=true, score=6.5, tags=temp}
   - Fait 2: [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}

91. **Token 91**:
   - Fait 1: [1] TestPerson{salary=-5000, score=0, status=active, department=qa, age=0, active=true, level=1, name=Frank, tags=test}
   - Fait 2: [1] TestOrder{region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1}

92. **Token 92**:
   - Fait 1: [1] TestPerson{name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, active=true, department=management, level=9}
   - Fait 2: [1] TestOrder{discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped}

93. **Token 93**:
   - Fait 1: [1] TestPerson{salary=95000, score=10, tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65}
   - Fait 2: [1] TestOrder{date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004}

94. **Token 94**:
   - Fait 1: [1] TestPerson{status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern, level=1}
   - Fait 2: [1] TestOrder{discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal}

95. **Token 95**:
   - Fait 1: [1] TestPerson{tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10}
   - Fait 2: [1] TestOrder{total=89.99, region=north, customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007}

96. **Token 96**:
   - Fait 1: [1] TestPerson{name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern, level=1, status=active}
   - Fait 2: [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

97. **Token 97**:
   - Fait 1: [1] TestPerson{tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10}
   - Fait 2: [1] TestOrder{customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north}

98. **Token 98**:
   - Fait 1: [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
   - Fait 2: [1] TestOrder{region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10}

99. **Token 99**:
   - Fait 1: [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
   - Fait 2: [1] TestOrder{customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15}

100. **Token 100**:
   - Fait 1: [1] TestPerson{department=management, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, active=true}
   - Fait 2: [1] TestOrder{discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales}
2. [1] TestPerson{level=5, status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000}
3. [1] TestPerson{level=1, score=6, status=inactive, age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr}
4. [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing}
5. [1] TestPerson{tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales}
6. [1] TestPerson{department=qa, age=0, active=true, level=1, name=Frank, tags=test, salary=-5000, score=0, status=active}
7. [1] TestPerson{status=active, active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive}
8. [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
9. [1] TestPerson{department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active, score=8.7}
10. [1] TestPerson{tags=temp, age=22, salary=28000, department=intern, level=1, status=active, name=X, active=true, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{score=8.5, level=2, salary=45000, active=true, department=sales, status=active, name=Alice, age=25, tags=junior}

2. **Token 2**:
   - Fait 1: [1] TestPerson{level=5, status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000}

3. **Token 3**:
   - Fait 1: [1] TestPerson{score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45, active=true}

4. **Token 4**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}

5. **Token 5**:
   - Fait 1: [1] TestPerson{level=9, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, active=true, department=management}

6. **Token 6**:
   - Fait 1: [1] TestPerson{tags=senior, level=6, active=true, salary=68000, status=active, score=8.7, department=engineering, name=Ivy, age=40}

7. **Token 7**:
   - Fait 1: [1] TestPerson{status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{score=8.5, level=2, salary=45000, active=true, department=sales, status=active, name=Alice, age=25, tags=junior}
2. [1] TestPerson{salary=75000, level=5, status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior}
3. [1] TestPerson{age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive}
4. [1] TestPerson{level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing, salary=85000}
5. [1] TestPerson{score=8, status=inactive, level=3, name=Eve, department=sales, tags=employee, age=30, salary=55000, active=false}
6. [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}
7. [1] TestPerson{tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10}
8. [1] TestPerson{level=1, salary=25000, name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support}
9. [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}
10. [1] TestPerson{status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern, level=1}
11. [1] TestOrder{discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15}
12. [1] TestOrder{region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1}
13. [1] TestOrder{total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped, discount=15, region=north, amount=3, customer_id=P001}
14. [1] TestOrder{amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004}
15. [1] TestOrder{region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10}
16. [1] TestOrder{amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west}
17. [1] TestOrder{discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01}
18. [1] TestOrder{product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010}
19. [1] TestOrder{status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low, amount=1, date=2024-03-10}
20. [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
   - Fait 2: [1] TestOrder{region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal}

2. **Token 2**:
   - Fait 1: [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}
   - Fait 2: [1] TestOrder{date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0}

3. **Token 3**:
   - Fait 1: [1] TestPerson{active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive, age=16}
   - Fait 2: [1] TestOrder{discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal}

4. **Token 4**:
   - Fait 1: [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing}
   - Fait 2: [1] TestOrder{priority=urgent, amount=1, date=2024-03-15, discount=0, region=east, customer_id=P006, product_id=PROD001, total=75000, status=refunded}

5. **Token 5**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}
   - Fait 2: [1] TestOrder{region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, department=intern, level=1, status=active, name=X, active=true, score=6.5, tags=temp}
   - Fait 2: [1] TestOrder{amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped, discount=15, region=north}

7. **Token 7**:
   - Fait 1: [1] TestPerson{tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{priority=high, product_id=PROD003, date=2024-02-01, status=shipped, discount=15, region=north, amount=3, customer_id=P001, total=225}

8. **Token 8**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}
   - Fait 2: [1] TestOrder{status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002}

9. **Token 9**:
   - Fait 1: [1] TestPerson{name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive, age=16, active=false}
   - Fait 2: [1] TestOrder{product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled}

10. **Token 10**:
   - Fait 1: [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}
   - Fait 2: [1] TestOrder{product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled}

11. **Token 11**:
   - Fait 1: [1] TestPerson{score=0, status=active, department=qa, age=0, active=true, level=1, name=Frank, tags=test, salary=-5000}
   - Fait 2: [1] TestOrder{discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01}

12. **Token 12**:
   - Fait 1: [1] TestPerson{age=65, salary=95000, score=10, tags=executive, status=active, active=true, department=management, level=9, name=Grace}
   - Fait 2: [1] TestOrder{total=89.99, region=north, customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007}

13. **Token 13**:
   - Fait 1: [1] TestPerson{age=65, salary=95000, score=10, tags=executive, status=active, active=true, department=management, level=9, name=Grace}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001}

14. **Token 14**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001}

15. **Token 15**:
   - Fait 1: [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
   - Fait 2: [1] TestOrder{region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped, discount=15}

16. **Token 16**:
   - Fait 1: [1] TestPerson{score=0, status=active, department=qa, age=0, active=true, level=1, name=Frank, tags=test, salary=-5000}
   - Fait 2: [1] TestOrder{region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal}

17. **Token 17**:
   - Fait 1: [1] TestPerson{score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45, active=true}
   - Fait 2: [1] TestOrder{status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10}

18. **Token 18**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
   - Fait 2: [1] TestOrder{total=89.99, region=north, customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007}

19. **Token 19**:
   - Fait 1: [1] TestPerson{age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive}
   - Fait 2: [1] TestOrder{status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east, customer_id=P006, product_id=PROD001, total=75000}

20. **Token 20**:
   - Fait 1: [1] TestPerson{department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active, score=8.7}
   - Fait 2: [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

21. **Token 21**:
   - Fait 1: [1] TestPerson{name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, active=true, department=management, level=9}
   - Fait 2: [1] TestOrder{amount=1, region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low}

22. **Token 22**:
   - Fait 1: [1] TestPerson{age=18, department=support, level=1, salary=25000, name=Henry, active=false, score=5.5, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1}

23. **Token 23**:
   - Fait 1: [1] TestPerson{salary=95000, score=10, tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65}
   - Fait 2: [1] TestOrder{total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed, priority=high}

24. **Token 24**:
   - Fait 1: [1] TestPerson{score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000, name=Henry, active=false}
   - Fait 2: [1] TestOrder{region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10}

25. **Token 25**:
   - Fait 1: [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing}
   - Fait 2: [1] TestOrder{product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled}

26. **Token 26**:
   - Fait 1: [1] TestPerson{status=active, active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive}
   - Fait 2: [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}

27. **Token 27**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
   - Fait 2: [1] TestOrder{priority=urgent, region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600, status=shipped}

28. **Token 28**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}
   - Fait 2: [1] TestOrder{customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north}

29. **Token 29**:
   - Fait 1: [1] TestPerson{score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45, active=true}
   - Fait 2: [1] TestOrder{customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south}

30. **Token 30**:
   - Fait 1: [1] TestPerson{score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45, active=true}
   - Fait 2: [1] TestOrder{date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99}

31. **Token 31**:
   - Fait 1: [1] TestPerson{tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10}
   - Fait 2: [1] TestOrder{customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15}

32. **Token 32**:
   - Fait 1: [1] TestPerson{tags=junior, status=inactive, age=18, department=support, level=1, salary=25000, name=Henry, active=false, score=5.5}
   - Fait 2: [1] TestOrder{customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15}

33. **Token 33**:
   - Fait 1: [1] TestPerson{score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45, active=true}
   - Fait 2: [1] TestOrder{discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01}

34. **Token 34**:
   - Fait 1: [1] TestPerson{score=6, status=inactive, age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1}
   - Fait 2: [1] TestOrder{total=89.99, region=north, customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007}

35. **Token 35**:
   - Fait 1: [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}
   - Fait 2: [1] TestOrder{customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north}

36. **Token 36**:
   - Fait 1: [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing}
   - Fait 2: [1] TestOrder{discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15}

37. **Token 37**:
   - Fait 1: [1] TestPerson{salary=0, tags=intern, department=hr, level=1, score=6, status=inactive, age=16, active=false, name=Charlie}
   - Fait 2: [1] TestOrder{date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004}

38. **Token 38**:
   - Fait 1: [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}
   - Fait 2: [1] TestOrder{status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1}

39. **Token 39**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
   - Fait 2: [1] TestOrder{product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010}

40. **Token 40**:
   - Fait 1: [1] TestPerson{name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, active=true, department=management, level=9}
   - Fait 2: [1] TestOrder{product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010}

41. **Token 41**:
   - Fait 1: [1] TestPerson{age=0, active=true, level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa}
   - Fait 2: [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

42. **Token 42**:
   - Fait 1: [1] TestPerson{level=2, salary=45000, active=true, department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5}
   - Fait 2: [1] TestOrder{priority=high, product_id=PROD003, date=2024-02-01, status=shipped, discount=15, region=north, amount=3, customer_id=P001, total=225}

43. **Token 43**:
   - Fait 1: [1] TestPerson{department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active, score=8.7}
   - Fait 2: [1] TestOrder{date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed}

44. **Token 44**:
   - Fait 1: [1] TestPerson{age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive}
   - Fait 2: [1] TestOrder{status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002}

45. **Token 45**:
   - Fait 1: [1] TestPerson{status=inactive, age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6}
   - Fait 2: [1] TestOrder{priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed}

46. **Token 46**:
   - Fait 1: [1] TestPerson{name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active, score=8.7, department=engineering}
   - Fait 2: [1] TestOrder{customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100}

47. **Token 47**:
   - Fait 1: [1] TestPerson{salary=-5000, score=0, status=active, department=qa, age=0, active=true, level=1, name=Frank, tags=test}
   - Fait 2: [1] TestOrder{amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west}

48. **Token 48**:
   - Fait 1: [1] TestPerson{tags=intern, department=hr, level=1, score=6, status=inactive, age=16, active=false, name=Charlie, salary=0}
   - Fait 2: [1] TestOrder{amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006}

49. **Token 49**:
   - Fait 1: [1] TestPerson{name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active, score=8.7, department=engineering}
   - Fait 2: [1] TestOrder{amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006}

50. **Token 50**:
   - Fait 1: [1] TestPerson{name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive, age=16, active=false}
   - Fait 2: [1] TestOrder{customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north}

51. **Token 51**:
   - Fait 1: [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001}

52. **Token 52**:
   - Fait 1: [1] TestPerson{status=active, department=marketing, salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager}
   - Fait 2: [1] TestOrder{amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004}

53. **Token 53**:
   - Fait 1: [1] TestPerson{score=0, status=active, department=qa, age=0, active=true, level=1, name=Frank, tags=test, salary=-5000}
   - Fait 2: [1] TestOrder{product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010}

54. **Token 54**:
   - Fait 1: [1] TestPerson{tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45, active=true, score=7.8}
   - Fait 2: [1] TestOrder{product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10}

55. **Token 55**:
   - Fait 1: [1] TestPerson{department=management, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, active=true}
   - Fait 2: [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

56. **Token 56**:
   - Fait 1: [1] TestPerson{department=support, level=1, salary=25000, name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18}
   - Fait 2: [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

57. **Token 57**:
   - Fait 1: [1] TestPerson{active=true, score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45}
   - Fait 2: [1] TestOrder{discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped}

58. **Token 58**:
   - Fait 1: [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
   - Fait 2: [1] TestOrder{discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal}

59. **Token 59**:
   - Fait 1: [1] TestPerson{score=8.5, level=2, salary=45000, active=true, department=sales, status=active, name=Alice, age=25, tags=junior}
   - Fait 2: [1] TestOrder{amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low}

60. **Token 60**:
   - Fait 1: [1] TestPerson{status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales}
   - Fait 2: [1] TestOrder{product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, status=delivered}

61. **Token 61**:
   - Fait 1: [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
   - Fait 2: [1] TestOrder{product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled}

62. **Token 62**:
   - Fait 1: [1] TestPerson{name=Eve, department=sales, tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{region=east, customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0}

63. **Token 63**:
   - Fait 1: [1] TestPerson{level=1, status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern}
   - Fait 2: [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

64. **Token 64**:
   - Fait 1: [1] TestPerson{status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales}
   - Fait 2: [1] TestOrder{date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed}

65. **Token 65**:
   - Fait 1: [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}
   - Fait 2: [1] TestOrder{discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped}

66. **Token 66**:
   - Fait 1: [1] TestPerson{name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive, age=16, active=false}
   - Fait 2: [1] TestOrder{discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped}

67. **Token 67**:
   - Fait 1: [1] TestPerson{tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales, status=active, name=Alice, age=25}
   - Fait 2: [1] TestOrder{amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001}

68. **Token 68**:
   - Fait 1: [1] TestPerson{name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern, level=1, status=active}
   - Fait 2: [1] TestOrder{customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15}

69. **Token 69**:
   - Fait 1: [1] TestPerson{tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15}

70. **Token 70**:
   - Fait 1: [1] TestPerson{level=2, salary=45000, active=true, department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5}
   - Fait 2: [1] TestOrder{priority=urgent, region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600, status=shipped}

71. **Token 71**:
   - Fait 1: [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
   - Fait 2: [1] TestOrder{amount=1, date=2024-03-15, discount=0, region=east, customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent}

72. **Token 72**:
   - Fait 1: [1] TestPerson{department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001}

73. **Token 73**:
   - Fait 1: [1] TestPerson{department=intern, level=1, status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000}
   - Fait 2: [1] TestOrder{status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10}

74. **Token 74**:
   - Fait 1: [1] TestPerson{tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{amount=10, status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002}

75. **Token 75**:
   - Fait 1: [1] TestPerson{level=3, name=Eve, department=sales, tags=employee, age=30, salary=55000, active=false, score=8, status=inactive}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001}

76. **Token 76**:
   - Fait 1: [1] TestPerson{department=intern, level=1, status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000}
   - Fait 2: [1] TestOrder{discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15}

77. **Token 77**:
   - Fait 1: [1] TestPerson{status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales}
   - Fait 2: [1] TestOrder{product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010}

78. **Token 78**:
   - Fait 1: [1] TestPerson{level=5, status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000}
   - Fait 2: [1] TestOrder{amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low}

79. **Token 79**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, department=intern, level=1, status=active, name=X, active=true, score=6.5, tags=temp}
   - Fait 2: [1] TestOrder{region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1}

80. **Token 80**:
   - Fait 1: [1] TestPerson{name=Eve, department=sales, tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3}
   - Fait 2: [1] TestOrder{amount=1, region=south, customer_id=P002, status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low}

81. **Token 81**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, tags=test, salary=-5000, score=0, status=active, department=qa, age=0, active=true}
   - Fait 2: [1] TestOrder{region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003, date=2024-02-01, status=shipped, discount=15}

82. **Token 82**:
   - Fait 1: [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
   - Fait 2: [1] TestOrder{region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal}

83. **Token 83**:
   - Fait 1: [1] TestPerson{salary=55000, active=false, score=8, status=inactive, level=3, name=Eve, department=sales, tags=employee, age=30}
   - Fait 2: [1] TestOrder{status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002}

84. **Token 84**:
   - Fait 1: [1] TestPerson{level=1, status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern}
   - Fait 2: [1] TestOrder{status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002}

85. **Token 85**:
   - Fait 1: [1] TestPerson{status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales}
   - Fait 2: [1] TestOrder{product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, priority=low, total=999.98, status=cancelled}

86. **Token 86**:
   - Fait 1: [1] TestPerson{score=6.5, tags=temp, age=22, salary=28000, department=intern, level=1, status=active, name=X, active=true}
   - Fait 2: [1] TestOrder{total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4}

87. **Token 87**:
   - Fait 1: [1] TestPerson{score=9.2, tags=senior, salary=75000, level=5, status=active, department=engineering, name=Bob, age=35, active=true}
   - Fait 2: [1] TestOrder{date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed}

88. **Token 88**:
   - Fait 1: [1] TestPerson{status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern, level=1}
   - Fait 2: [1] TestOrder{discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05}

89. **Token 89**:
   - Fait 1: [1] TestPerson{status=inactive, level=3, name=Eve, department=sales, tags=employee, age=30, salary=55000, active=false, score=8}
   - Fait 2: [1] TestOrder{status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1}

90. **Token 90**:
   - Fait 1: [1] TestPerson{tags=senior, salary=75000, level=5, status=active, department=engineering, name=Bob, age=35, active=true, score=9.2}
   - Fait 2: [1] TestOrder{region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10}

91. **Token 91**:
   - Fait 1: [1] TestPerson{department=sales, tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve}
   - Fait 2: [1] TestOrder{discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01}

92. **Token 92**:
   - Fait 1: [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
   - Fait 2: [1] TestOrder{amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006}

93. **Token 93**:
   - Fait 1: [1] TestPerson{status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000, active=true, department=sales}
   - Fait 2: [1] TestOrder{total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east, customer_id=P006, product_id=PROD001}

94. **Token 94**:
   - Fait 1: [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001}

95. **Token 95**:
   - Fait 1: [1] TestPerson{tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10}
   - Fait 2: [1] TestOrder{priority=high, product_id=PROD003, date=2024-02-01, status=shipped, discount=15, region=north, amount=3, customer_id=P001, total=225}

96. **Token 96**:
   - Fait 1: [1] TestPerson{department=support, level=1, salary=25000, name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18}
   - Fait 2: [1] TestOrder{date=2024-02-01, status=shipped, discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003}

97. **Token 97**:
   - Fait 1: [1] TestPerson{score=8, status=inactive, level=3, name=Eve, department=sales, tags=employee, age=30, salary=55000, active=false}
   - Fait 2: [1] TestOrder{product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10}

98. **Token 98**:
   - Fait 1: [1] TestPerson{status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern, level=1}
   - Fait 2: [1] TestOrder{total=89.99, region=north, customer_id=P001, priority=low, amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007}

99. **Token 99**:
   - Fait 1: [1] TestPerson{department=qa, age=0, active=true, level=1, name=Frank, tags=test, salary=-5000, score=0, status=active}
   - Fait 2: [1] TestOrder{amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001}

100. **Token 100**:
   - Fait 1: [1] TestPerson{score=10, tags=executive, status=active, active=true, department=management, level=9, name=Grace, age=65, salary=95000}
   - Fait 2: [1] TestOrder{status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{salary=45000, active=true, department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5, level=2}
2. [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
3. [1] TestPerson{score=6, status=inactive, age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1}
4. [1] TestPerson{department=marketing, salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active}
5. [1] TestPerson{department=sales, tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3, name=Eve}
6. [1] TestPerson{department=qa, age=0, active=true, level=1, name=Frank, tags=test, salary=-5000, score=0, status=active}
7. [1] TestPerson{name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, active=true, department=management, level=9}
8. [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
9. [1] TestPerson{score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active}
10. [1] TestPerson{salary=28000, department=intern, level=1, status=active, name=X, active=true, score=6.5, tags=temp, age=22}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{salary=75000, level=5, status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior}

2. **Token 2**:
   - Fait 1: [1] TestPerson{active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active}

3. **Token 3**:
   - Fait 1: [1] TestPerson{name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active, score=8.7, department=engineering}

4. **Token 4**:
   - Fait 1: [1] TestPerson{name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive, age=16, active=false}

5. **Token 5**:
   - Fait 1: [1] TestPerson{active=true, score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45}

6. **Token 6**:
   - Fait 1: [1] TestPerson{name=Eve, department=sales, tags=employee, age=30, salary=55000, active=false, score=8, status=inactive, level=3}

7. **Token 7**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}

8. **Token 8**:
   - Fait 1: [1] TestPerson{status=active, name=X, active=true, score=6.5, tags=temp, age=22, salary=28000, department=intern, level=1}

9. **Token 9**:
   - Fait 1: [1] TestPerson{level=2, salary=45000, active=true, department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{customer_id=P001, total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north}
2. [1] TestOrder{date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002, status=confirmed}
3. [1] TestOrder{priority=high, product_id=PROD003, date=2024-02-01, status=shipped, discount=15, region=north, amount=3, customer_id=P001, total=225}
4. [1] TestOrder{date=2024-02-05, discount=0, total=299.99, priority=normal, region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004}
5. [1] TestOrder{discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10, region=south, product_id=PROD001, amount=1}
6. [1] TestOrder{customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15}
7. [1] TestOrder{region=north, customer_id=P007, date=2024-03-01, discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent}
8. [1] TestOrder{product_id=PROD002, amount=10, status=pending, region=south, total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010}
9. [1] TestOrder{amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low}
10. [1] TestOrder{customer_id=P006, product_id=PROD001, total=75000, status=refunded, priority=urgent, amount=1, date=2024-03-15, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{priority=normal, region=east, customer_id=P004, amount=1, status=delivered, product_id=PROD004, date=2024-02-05, discount=0, total=299.99}

2. **Token 2**:
   - Fait 1: [1] TestOrder{amount=1, date=2024-03-10, status=completed, discount=10, product_id=PROD007, total=89.99, region=north, customer_id=P001, priority=low}

3. **Token 3**:
   - Fait 1: [1] TestOrder{total=1999.98, status=pending, priority=normal, product_id=PROD001, amount=2, date=2024-01-15, discount=50, region=north, customer_id=P001}

4. **Token 4**:
   - Fait 1: [1] TestOrder{status=confirmed, date=2024-01-20, discount=0, product_id=PROD002, total=25.5, priority=low, amount=1, region=south, customer_id=P002}

5. **Token 5**:
   - Fait 1: [1] TestOrder{date=2024-02-01, status=shipped, discount=15, region=north, amount=3, customer_id=P001, total=225, priority=high, product_id=PROD003}

6. **Token 6**:
   - Fait 1: [1] TestOrder{region=south, product_id=PROD001, amount=1, discount=100, customer_id=P002, status=confirmed, priority=high, total=999.99, date=2024-02-10}

7. **Token 7**:
   - Fait 1: [1] TestOrder{customer_id=P005, priority=low, total=999.98, status=cancelled, product_id=PROD005, discount=0, region=west, amount=2, date=2024-02-15}

8. **Token 8**:
   - Fait 1: [1] TestOrder{discount=50, product_id=PROD006, amount=4, total=600, status=shipped, priority=urgent, region=north, customer_id=P007, date=2024-03-01}

9. **Token 9**:
   - Fait 1: [1] TestOrder{total=255, priority=normal, discount=0, date=2024-03-05, customer_id=P010, product_id=PROD002, amount=10, status=pending, region=south}

10. **Token 10**:
   - Fait 1: [1] TestOrder{priority=urgent, amount=1, date=2024-03-15, discount=0, region=east, customer_id=P006, product_id=PROD001, total=75000, status=refunded}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{level=2, salary=45000, active=true, department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5}
2. [1] TestPerson{status=active, department=engineering, name=Bob, age=35, active=true, score=9.2, tags=senior, salary=75000, level=5}
3. [1] TestPerson{name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive, age=16, active=false}
4. [1] TestPerson{salary=85000, level=7, name=Diana, age=45, active=true, score=7.8, tags=manager, status=active, department=marketing}
5. [1] TestPerson{score=8, status=inactive, level=3, name=Eve, department=sales, tags=employee, age=30, salary=55000, active=false}
6. [1] TestPerson{department=qa, age=0, active=true, level=1, name=Frank, tags=test, salary=-5000, score=0, status=active}
7. [1] TestPerson{active=true, department=management, level=9, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active}
8. [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}
9. [1] TestPerson{status=active, score=8.7, department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000}
10. [1] TestPerson{age=22, salary=28000, department=intern, level=1, status=active, name=X, active=true, score=6.5, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{age=22, salary=28000, department=intern, level=1, status=active, name=X, active=true, score=6.5, tags=temp}

2. **Token 2**:
   - Fait 1: [1] TestPerson{active=true, score=7.8, tags=manager, status=active, department=marketing, salary=85000, level=7, name=Diana, age=45}

3. **Token 3**:
   - Fait 1: [1] TestPerson{level=9, name=Grace, age=65, salary=95000, score=10, tags=executive, status=active, active=true, department=management}

4. **Token 4**:
   - Fait 1: [1] TestPerson{department=engineering, name=Ivy, age=40, tags=senior, level=6, active=true, salary=68000, status=active, score=8.7}

5. **Token 5**:
   - Fait 1: [1] TestPerson{active=true, department=sales, status=active, name=Alice, age=25, tags=junior, score=8.5, level=2, salary=45000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{score=9.2, tags=senior, salary=75000, level=5, status=active, department=engineering, name=Bob, age=35, active=true}

7. **Token 7**:
   - Fait 1: [1] TestPerson{age=16, active=false, name=Charlie, salary=0, tags=intern, department=hr, level=1, score=6, status=inactive}

8. **Token 8**:
   - Fait 1: [1] TestPerson{score=8, status=inactive, level=3, name=Eve, department=sales, tags=employee, age=30, salary=55000, active=false}

9. **Token 9**:
   - Fait 1: [1] TestPerson{department=qa, age=0, active=true, level=1, name=Frank, tags=test, salary=-5000, score=0, status=active}

10. **Token 10**:
   - Fait 1: [1] TestPerson{name=Henry, active=false, score=5.5, tags=junior, status=inactive, age=18, department=support, level=1, salary=25000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
