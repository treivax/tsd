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

1. [1] TestPerson{active=true, tags=junior, status=active, score=8.5, level=2, age=25, department=sales, name=Alice, salary=45000}
2. [1] TestPerson{active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2, name=Bob, salary=75000}
3. [1] TestPerson{score=6, active=false, tags=intern, status=inactive, level=1, name=Charlie, age=16, salary=0, department=hr}
4. [1] TestPerson{score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager}
5. [1] TestPerson{age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales}
6. [1] TestPerson{active=true, tags=test, status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0}
7. [1] TestPerson{age=65, department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000}
8. [1] TestPerson{status=inactive, active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry, tags=junior}
9. [1] TestPerson{score=8.7, level=6, tags=senior, department=engineering, active=true, age=40, status=active, name=Ivy, salary=68000}
10. [1] TestPerson{department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{active=true, tags=junior, status=active, score=8.5, level=2, age=25, department=sales, name=Alice, salary=45000}

2. **Token 2**:
   - Fait 1: [1] TestPerson{active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2, name=Bob, salary=75000}

3. **Token 3**:
   - Fait 1: [1] TestPerson{score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager}

4. **Token 4**:
   - Fait 1: [1] TestPerson{age=18, name=Henry, tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support, level=1}

5. **Token 5**:
   - Fait 1: [1] TestPerson{level=6, tags=senior, department=engineering, active=true, age=40, status=active, name=Ivy, salary=68000, score=8.7}

6. **Token 6**:
   - Fait 1: [1] TestPerson{age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1, name=Charlie}

7. **Token 7**:
   - Fait 1: [1] TestPerson{salary=55000, department=sales, age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3}

8. **Token 8**:
   - Fait 1: [1] TestPerson{name=Grace, active=true, tags=executive, salary=95000, age=65, department=management, score=10, status=active, level=9}

9. **Token 9**:
   - Fait 1: [1] TestPerson{salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending, region=north, discount=50}
2. [1] TestOrder{discount=0, total=25.5, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south}
3. [1] TestOrder{status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225}
4. [1] TestOrder{customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99}
6. [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west}
7. [1] TestOrder{product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north, customer_id=P007}
8. [1] TestOrder{amount=10, status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal}
9. [1] TestOrder{date=2024-03-10, region=north, total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007}
10. [1] TestOrder{amount=1, total=75000, discount=0, region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10}

2. **Token 2**:
   - Fait 1: [1] TestOrder{product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south, customer_id=P010}

3. **Token 3**:
   - Fait 1: [1] TestOrder{date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0, region=east}

4. **Token 4**:
   - Fait 1: [1] TestOrder{date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending, region=north, discount=50}

5. **Token 5**:
   - Fait 1: [1] TestOrder{status=confirmed, region=south, discount=0, total=25.5, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1}

6. **Token 6**:
   - Fait 1: [1] TestOrder{status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225}

7. **Token 7**:
   - Fait 1: [1] TestOrder{customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}

8. **Token 8**:
   - Fait 1: [1] TestOrder{amount=4, total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50}

9. **Token 9**:
   - Fait 1: [1] TestOrder{region=north, total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}
2. [1] TestPerson{status=active, score=9.2, name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35}
3. [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}
4. [1] TestPerson{status=active, department=marketing, name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true}
5. [1] TestPerson{age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales}
6. [1] TestPerson{salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active, name=Frank, department=qa}
7. [1] TestPerson{name=Grace, active=true, tags=executive, salary=95000, age=65, department=management, score=10, status=active, level=9}
8. [1] TestPerson{level=1, age=18, name=Henry, tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support}
9. [1] TestPerson{department=engineering, active=true, age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior}
10. [1] TestPerson{score=6.5, status=active, level=1, salary=28000, department=intern, name=X, age=22, tags=temp, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true}

2. **Token 2**:
   - Fait 1: [1] TestPerson{tags=junior, status=active, score=8.5, level=2, age=25, department=sales, name=Alice, salary=45000, active=true}

3. **Token 3**:
   - Fait 1: [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}

4. **Token 4**:
   - Fait 1: [1] TestPerson{age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager, score=7.8, level=7}

5. **Token 5**:
   - Fait 1: [1] TestPerson{age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales}

6. **Token 6**:
   - Fait 1: [1] TestPerson{active=true, tags=executive, salary=95000, age=65, department=management, score=10, status=active, level=9, name=Grace}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending, region=north, discount=50}
2. [1] TestOrder{region=south, discount=0, total=25.5, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed}
3. [1] TestOrder{discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225, status=shipped}
4. [1] TestOrder{discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal}
5. [1] TestOrder{amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south}
6. [1] TestOrder{region=west, total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0}
7. [1] TestOrder{status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, priority=urgent}
8. [1] TestOrder{status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10}
9. [1] TestOrder{customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99, status=completed, priority=low, amount=1, discount=10}
10. [1] TestOrder{discount=0, region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}

2. **Token 2**:
   - Fait 1: [1] TestOrder{status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225}

3. **Token 3**:
   - Fait 1: [1] TestOrder{priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered}

4. **Token 4**:
   - Fait 1: [1] TestOrder{priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed}

5. **Token 5**:
   - Fait 1: [1] TestOrder{customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{product_id=PROD007, date=2024-03-10, region=north, total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001}

7. **Token 7**:
   - Fait 1: [1] TestOrder{date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west, total=999.98, status=cancelled, product_id=PROD005, amount=2}

8. **Token 8**:
   - Fait 1: [1] TestOrder{amount=4, total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50}

9. **Token 9**:
   - Fait 1: [1] TestOrder{total=1999.98, status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{keywords=computer, supplier=TechSupply, name=Laptop, available=true, brand=TechCorp, price=999.99, rating=4.5, stock=50, category=electronics}
2. [1] TestProduct{name=Mouse, rating=4, stock=200, price=25.5, available=true, category=accessories, keywords=peripheral, brand=TechCorp, supplier=TechSupply}
3. [1] TestProduct{available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply, category=accessories, name=Keyboard, price=75}
4. [1] TestProduct{name=Monitor, brand=ScreenPro, supplier=ScreenSupply, rating=4.8, stock=30, category=electronics, price=299.99, available=true, keywords=display}
5. [1] TestProduct{category=accessories, keywords=obsolete, stock=0, brand=OldTech, price=8.5, available=false, name=OldKeyboard, rating=2, supplier=OldSupply}
6. [1] TestProduct{price=150, supplier=AudioSupply, category=audio, available=true, keywords=sound, brand=AudioMax, stock=75, name=Headphones, rating=4.6}
7. [1] TestProduct{category=electronics, price=89.99, brand=CamTech, stock=25, supplier=CamSupply, rating=3.8, keywords=video, name=Webcam, available=true}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{price=150, supplier=AudioSupply, category=audio, available=true, keywords=sound, brand=AudioMax, stock=75, name=Headphones, rating=4.6}

2. **Token 2**:
   - Fait 1: [1] TestProduct{brand=CamTech, stock=25, supplier=CamSupply, rating=3.8, keywords=video, name=Webcam, available=true, category=electronics, price=89.99}

3. **Token 3**:
   - Fait 1: [1] TestProduct{keywords=computer, supplier=TechSupply, name=Laptop, available=true, brand=TechCorp, price=999.99, rating=4.5, stock=50, category=electronics}

4. **Token 4**:
   - Fait 1: [1] TestProduct{supplier=TechSupply, name=Mouse, rating=4, stock=200, price=25.5, available=true, category=accessories, keywords=peripheral, brand=TechCorp}

5. **Token 5**:
   - Fait 1: [1] TestProduct{price=75, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply, category=accessories, name=Keyboard}

6. **Token 6**:
   - Fait 1: [1] TestProduct{name=Monitor, brand=ScreenPro, supplier=ScreenSupply, rating=4.8, stock=30, category=electronics, price=299.99, available=true, keywords=display}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}
2. [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
3. [1] TestPerson{tags=intern, status=inactive, level=1, name=Charlie, age=16, salary=0, department=hr, score=6, active=false}
4. [1] TestPerson{score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager}
5. [1] TestPerson{age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales}
6. [1] TestPerson{age=0, active=true, tags=test, status=active, name=Frank, department=qa, salary=-5000, score=0, level=1}
7. [1] TestPerson{age=65, department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000}
8. [1] TestPerson{active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry, tags=junior, status=inactive}
9. [1] TestPerson{department=engineering, active=true, age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior}
10. [1] TestPerson{level=1, salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{status=inactive, level=1, name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern}

2. **Token 2**:
   - Fait 1: [1] TestPerson{active=true, tags=test, status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending, region=north, discount=50}
2. [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}
3. [1] TestOrder{status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225}
4. [1] TestOrder{amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004}
5. [1] TestOrder{amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south}
6. [1] TestOrder{date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west, total=999.98, status=cancelled, product_id=PROD005, amount=2}
7. [1] TestOrder{amount=4, total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50}
8. [1] TestOrder{total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002}
9. [1] TestOrder{region=north, total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10}
10. [1] TestOrder{status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0, region=east, date=2024-03-15, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{total=25.5, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestOrder{customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}

3. **Token 3**:
   - Fait 1: [1] TestOrder{total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002}

4. **Token 4**:
   - Fait 1: [1] TestOrder{status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0, region=east, date=2024-03-15, priority=urgent}

5. **Token 5**:
   - Fait 1: [1] TestOrder{product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225, status=shipped, discount=15, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99}

7. **Token 7**:
   - Fait 1: [1] TestOrder{region=west, total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestOrder{customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north}

9. **Token 9**:
   - Fait 1: [1] TestOrder{date=2024-03-10, region=north, total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007}

10. **Token 10**:
   - Fait 1: [1] TestOrder{total=1999.98, status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5, level=2}
2. [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
3. [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}
4. [1] TestPerson{score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager}
5. [1] TestPerson{age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales}
6. [1] TestPerson{status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0, active=true, tags=test}
7. [1] TestPerson{active=true, tags=executive, salary=95000, age=65, department=management, score=10, status=active, level=9, name=Grace}
8. [1] TestPerson{salary=25000, department=support, level=1, age=18, name=Henry, tags=junior, status=inactive, active=false, score=5.5}
9. [1] TestPerson{salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true, age=40, status=active, name=Ivy}
10. [1] TestPerson{salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1, name=Charlie, age=16}

2. **Token 2**:
   - Fait 1: [1] TestPerson{level=3, salary=55000, department=sales, age=30, score=8, tags=employee, status=inactive, active=false, name=Eve}

3. **Token 3**:
   - Fait 1: [1] TestPerson{level=1, age=18, name=Henry, tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending, region=north}
2. [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}
3. [1] TestOrder{customer_id=P001, amount=3, total=225, status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high}
4. [1] TestOrder{customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001}
6. [1] TestOrder{date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west, total=999.98, status=cancelled, product_id=PROD005, amount=2}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north}
8. [1] TestOrder{customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south}
9. [1] TestOrder{status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99}
10. [1] TestOrder{total=75000, discount=0, region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{discount=0, region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000}

2. **Token 2**:
   - Fait 1: [1] TestOrder{region=south, discount=0, total=25.5, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed}

3. **Token 3**:
   - Fait 1: [1] TestOrder{total=225, status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3}

4. **Token 4**:
   - Fait 1: [1] TestOrder{discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal}

5. **Token 5**:
   - Fait 1: [1] TestOrder{discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99}

6. **Token 6**:
   - Fait 1: [1] TestOrder{date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west, total=999.98, status=cancelled, product_id=PROD005, amount=2}

7. **Token 7**:
   - Fait 1: [1] TestOrder{priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006}

8. **Token 8**:
   - Fait 1: [1] TestOrder{status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{active=true, tags=junior, status=active, score=8.5, level=2, age=25, department=sales, name=Alice, salary=45000}
2. [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
3. [1] TestPerson{age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1, name=Charlie}
4. [1] TestPerson{score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager}
5. [1] TestPerson{department=sales, age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000}
6. [1] TestPerson{salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active, name=Frank, department=qa}
7. [1] TestPerson{age=65, department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000}
8. [1] TestPerson{name=Henry, tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support, level=1, age=18}
9. [1] TestPerson{active=true, age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering}
10. [1] TestPerson{salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true, status=active, department=marketing}

2. **Token 2**:
   - Fait 1: [1] TestPerson{score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales, age=30}

3. **Token 3**:
   - Fait 1: [1] TestPerson{age=65, department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{tags=senior, department=engineering, active=true, age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6}

5. **Token 5**:
   - Fait 1: [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{score=9.2, name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active}

7. **Token 7**:
   - Fait 1: [1] TestPerson{age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1, name=Charlie}

8. **Token 8**:
   - Fait 1: [1] TestPerson{tags=test, status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0, active=true}

9. **Token 9**:
   - Fait 1: [1] TestPerson{active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry, tags=junior, status=inactive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{keywords=computer, supplier=TechSupply, name=Laptop, available=true, brand=TechCorp, price=999.99, rating=4.5, stock=50, category=electronics}
2. [1] TestProduct{name=Mouse, rating=4, stock=200, price=25.5, available=true, category=accessories, keywords=peripheral, brand=TechCorp, supplier=TechSupply}
3. [1] TestProduct{stock=0, supplier=KeySupply, category=accessories, name=Keyboard, price=75, available=false, rating=3.5, keywords=typing, brand=KeyTech}
4. [1] TestProduct{category=electronics, price=299.99, available=true, keywords=display, name=Monitor, brand=ScreenPro, supplier=ScreenSupply, rating=4.8, stock=30}
5. [1] TestProduct{name=OldKeyboard, rating=2, supplier=OldSupply, category=accessories, keywords=obsolete, stock=0, brand=OldTech, price=8.5, available=false}
6. [1] TestProduct{name=Headphones, rating=4.6, price=150, supplier=AudioSupply, category=audio, available=true, keywords=sound, brand=AudioMax, stock=75}
7. [1] TestProduct{rating=3.8, keywords=video, name=Webcam, available=true, category=electronics, price=89.99, brand=CamTech, stock=25, supplier=CamSupply}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{available=true, keywords=sound, brand=AudioMax, stock=75, name=Headphones, rating=4.6, price=150, supplier=AudioSupply, category=audio}

2. **Token 2**:
   - Fait 1: [1] TestProduct{category=electronics, price=89.99, brand=CamTech, stock=25, supplier=CamSupply, rating=3.8, keywords=video, name=Webcam, available=true}

3. **Token 3**:
   - Fait 1: [1] TestProduct{available=true, brand=TechCorp, price=999.99, rating=4.5, stock=50, category=electronics, keywords=computer, supplier=TechSupply, name=Laptop}

4. **Token 4**:
   - Fait 1: [1] TestProduct{rating=4, stock=200, price=25.5, available=true, category=accessories, keywords=peripheral, brand=TechCorp, supplier=TechSupply, name=Mouse}

5. **Token 5**:
   - Fait 1: [1] TestProduct{price=75, available=false, rating=3.5, keywords=typing, brand=KeyTech, stock=0, supplier=KeySupply, category=accessories, name=Keyboard}

6. **Token 6**:
   - Fait 1: [1] TestProduct{price=299.99, available=true, keywords=display, name=Monitor, brand=ScreenPro, supplier=ScreenSupply, rating=4.8, stock=30, category=electronics}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5, level=2, age=25, department=sales}
2. [1] TestPerson{age=35, status=active, score=9.2, name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5}
3. [1] TestPerson{active=false, tags=intern, status=inactive, level=1, name=Charlie, age=16, salary=0, department=hr, score=6}
4. [1] TestPerson{status=active, department=marketing, name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true}
5. [1] TestPerson{age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales}
6. [1] TestPerson{active=true, tags=test, status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0}
7. [1] TestPerson{age=65, department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000}
8. [1] TestPerson{salary=25000, department=support, level=1, age=18, name=Henry, tags=junior, status=inactive, active=false, score=5.5}
9. [1] TestPerson{age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true}
10. [1] TestPerson{score=6.5, status=active, level=1, salary=28000, department=intern, name=X, age=22, tags=temp, active=true}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}

2. **Token 2**:
   - Fait 1: [1] TestPerson{score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager}

3. **Token 3**:
   - Fait 1: [1] TestPerson{age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true}

4. **Token 4**:
   - Fait 1: [1] TestPerson{department=engineering, level=5, age=35, status=active, score=9.2, name=Bob, salary=75000, active=true, tags=senior}

5. **Token 5**:
   - Fait 1: [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}

6. **Token 6**:
   - Fait 1: [1] TestPerson{age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales}

7. **Token 7**:
   - Fait 1: [1] TestPerson{tags=test, status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0, active=true}

8. **Token 8**:
   - Fait 1: [1] TestPerson{department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000, age=65}

9. **Token 9**:
   - Fait 1: [1] TestPerson{active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry, tags=junior, status=inactive}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{product_id=PROD001, amount=2, total=1999.98, status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001}
2. [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}
3. [1] TestOrder{product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225, status=shipped, discount=15, region=north}
4. [1] TestOrder{status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05}
5. [1] TestOrder{product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100}
6. [1] TestOrder{region=west, total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0}
7. [1] TestOrder{status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, priority=urgent}
8. [1] TestOrder{total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002}
9. [1] TestOrder{amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99, status=completed, priority=low}
10. [1] TestOrder{discount=0, region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north}

3. **Token 3**:
   - Fait 1: [1] TestOrder{discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10, status=pending}

4. **Token 4**:
   - Fait 1: [1] TestOrder{product_id=PROD007, date=2024-03-10, region=north, total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001}

5. **Token 5**:
   - Fait 1: [1] TestOrder{discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}

7. **Token 7**:
   - Fait 1: [1] TestOrder{discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225, status=shipped}

8. **Token 8**:
   - Fait 1: [1] TestOrder{customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{salary=45000, active=true, tags=junior, status=active, score=8.5, level=2, age=25, department=sales, name=Alice}
2. [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
3. [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}
4. [1] TestPerson{score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager}
5. [1] TestPerson{age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales}
6. [1] TestPerson{name=Frank, department=qa, salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active}
7. [1] TestPerson{age=65, department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000}
8. [1] TestPerson{status=inactive, active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry, tags=junior}
9. [1] TestPerson{age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true}
10. [1] TestPerson{age=22, tags=temp, active=true, score=6.5, status=active, level=1, salary=28000, department=intern, name=X}
11. [1] TestOrder{total=1999.98, status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2}
12. [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}
13. [1] TestOrder{status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225}
14. [1] TestOrder{discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal}
15. [1] TestOrder{amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south}
16. [1] TestOrder{discount=0, region=west, total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low}
17. [1] TestOrder{discount=50, amount=4, total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, priority=urgent, status=shipped}
18. [1] TestOrder{customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south}
19. [1] TestOrder{status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99}
20. [1] TestOrder{discount=0, region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{active=false, name=Eve, level=3, salary=55000, department=sales, age=30, score=8, tags=employee, status=inactive}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south}

2. **Token 2**:
   - Fait 1: [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
   - Fait 2: [1] TestOrder{discount=50, amount=4, total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, priority=urgent, status=shipped}

3. **Token 3**:
   - Fait 1: [1] TestPerson{salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1, name=Charlie, age=16}
   - Fait 2: [1] TestOrder{discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10, status=pending}

4. **Token 4**:
   - Fait 1: [1] TestPerson{age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1, name=Charlie}
   - Fait 2: [1] TestOrder{amount=3, total=225, status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001}

5. **Token 5**:
   - Fait 1: [1] TestPerson{salary=55000, department=sales, age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3}
   - Fait 2: [1] TestOrder{amount=3, total=225, status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001}

6. **Token 6**:
   - Fait 1: [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}
   - Fait 2: [1] TestOrder{status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05}

7. **Token 7**:
   - Fait 1: [1] TestPerson{tags=test, status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0, active=true}
   - Fait 2: [1] TestOrder{discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99}

8. **Token 8**:
   - Fait 1: [1] TestPerson{active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100}

9. **Token 9**:
   - Fait 1: [1] TestPerson{department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1, salary=28000}
   - Fait 2: [1] TestOrder{discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99}

10. **Token 10**:
   - Fait 1: [1] TestPerson{status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0, active=true, tags=test}
   - Fait 2: [1] TestOrder{date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255}

11. **Token 11**:
   - Fait 1: [1] TestPerson{tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry}
   - Fait 2: [1] TestOrder{customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south}

12. **Token 12**:
   - Fait 1: [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}
   - Fait 2: [1] TestOrder{region=south, discount=0, total=25.5, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed}

13. **Token 13**:
   - Fait 1: [1] TestPerson{tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales, age=30, score=8}
   - Fait 2: [1] TestOrder{discount=0, total=25.5, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south}

14. **Token 14**:
   - Fait 1: [1] TestPerson{age=65, department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000}
   - Fait 2: [1] TestOrder{discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99}

15. **Token 15**:
   - Fait 1: [1] TestPerson{level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager, score=7.8}
   - Fait 2: [1] TestOrder{priority=low, discount=0, region=west, total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005}

16. **Token 16**:
   - Fait 1: [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}
   - Fait 2: [1] TestOrder{amount=4, total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50}

17. **Token 17**:
   - Fait 1: [1] TestPerson{salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1}
   - Fait 2: [1] TestOrder{customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south}

18. **Token 18**:
   - Fait 1: [1] TestPerson{name=Grace, active=true, tags=executive, salary=95000, age=65, department=management, score=10, status=active, level=9}
   - Fait 2: [1] TestOrder{customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south}

19. **Token 19**:
   - Fait 1: [1] TestPerson{salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active, name=Frank, department=qa}
   - Fait 2: [1] TestOrder{region=north, total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10}

20. **Token 20**:
   - Fait 1: [1] TestPerson{score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{amount=2, total=1999.98, status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001}

21. **Token 21**:
   - Fait 1: [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}
   - Fait 2: [1] TestOrder{status=confirmed, region=south, discount=0, total=25.5, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1}

22. **Token 22**:
   - Fait 1: [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
   - Fait 2: [1] TestOrder{customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high}

23. **Token 23**:
   - Fait 1: [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}
   - Fait 2: [1] TestOrder{amount=4, total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50}

24. **Token 24**:
   - Fait 1: [1] TestPerson{name=Eve, level=3, salary=55000, department=sales, age=30, score=8, tags=employee, status=inactive, active=false}
   - Fait 2: [1] TestOrder{customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north}

25. **Token 25**:
   - Fait 1: [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
   - Fait 2: [1] TestOrder{customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south}

26. **Token 26**:
   - Fait 1: [1] TestPerson{tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales, age=30, score=8}
   - Fait 2: [1] TestOrder{status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10}

27. **Token 27**:
   - Fait 1: [1] TestPerson{level=1, age=18, name=Henry, tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support}
   - Fait 2: [1] TestOrder{status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99}

28. **Token 28**:
   - Fait 1: [1] TestPerson{name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5, level=2, age=25, department=sales}
   - Fait 2: [1] TestOrder{status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225}

29. **Token 29**:
   - Fait 1: [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}
   - Fait 2: [1] TestOrder{customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}

30. **Token 30**:
   - Fait 1: [1] TestPerson{tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales, age=30, score=8}
   - Fait 2: [1] TestOrder{customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}

31. **Token 31**:
   - Fait 1: [1] TestPerson{active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2, name=Bob, salary=75000}
   - Fait 2: [1] TestOrder{total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north}

32. **Token 32**:
   - Fait 1: [1] TestPerson{score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000, age=65, department=management}
   - Fait 2: [1] TestOrder{discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99, status=completed, priority=low, amount=1}

33. **Token 33**:
   - Fait 1: [1] TestPerson{score=6, active=false, tags=intern, status=inactive, level=1, name=Charlie, age=16, salary=0, department=hr}
   - Fait 2: [1] TestOrder{total=75000, discount=0, region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1}

34. **Token 34**:
   - Fait 1: [1] TestPerson{salary=55000, department=sales, age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3}
   - Fait 2: [1] TestOrder{total=75000, discount=0, region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1}

35. **Token 35**:
   - Fait 1: [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
   - Fait 2: [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}

36. **Token 36**:
   - Fait 1: [1] TestPerson{tags=test, status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0, active=true}
   - Fait 2: [1] TestOrder{region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0}

37. **Token 37**:
   - Fait 1: [1] TestPerson{tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry}
   - Fait 2: [1] TestOrder{customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}

38. **Token 38**:
   - Fait 1: [1] TestPerson{active=true, score=6.5, status=active, level=1, salary=28000, department=intern, name=X, age=22, tags=temp}
   - Fait 2: [1] TestOrder{status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05}

39. **Token 39**:
   - Fait 1: [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south}

40. **Token 40**:
   - Fait 1: [1] TestPerson{score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000, age=65, department=management}
   - Fait 2: [1] TestOrder{total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4}

41. **Token 41**:
   - Fait 1: [1] TestPerson{active=true, tags=junior, status=active, score=8.5, level=2, age=25, department=sales, name=Alice, salary=45000}
   - Fait 2: [1] TestOrder{date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255}

42. **Token 42**:
   - Fait 1: [1] TestPerson{active=true, score=6.5, status=active, level=1, salary=28000, department=intern, name=X, age=22, tags=temp}
   - Fait 2: [1] TestOrder{region=north, total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10}

43. **Token 43**:
   - Fait 1: [1] TestPerson{active=true, age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2}

44. **Token 44**:
   - Fait 1: [1] TestPerson{active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2, name=Bob, salary=75000}
   - Fait 2: [1] TestOrder{status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225}

45. **Token 45**:
   - Fait 1: [1] TestPerson{tags=executive, salary=95000, age=65, department=management, score=10, status=active, level=9, name=Grace, active=true}
   - Fait 2: [1] TestOrder{customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}

46. **Token 46**:
   - Fait 1: [1] TestPerson{active=true, tags=test, status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0}
   - Fait 2: [1] TestOrder{customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north}

47. **Token 47**:
   - Fait 1: [1] TestPerson{level=1, age=18, name=Henry, tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support}
   - Fait 2: [1] TestOrder{total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4}

48. **Token 48**:
   - Fait 1: [1] TestPerson{salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1}
   - Fait 2: [1] TestOrder{amount=4, total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50}

49. **Token 49**:
   - Fait 1: [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}
   - Fait 2: [1] TestOrder{date=2024-03-10, region=north, total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007}

50. **Token 50**:
   - Fait 1: [1] TestPerson{age=65, department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000}
   - Fait 2: [1] TestOrder{total=75000, discount=0, region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1}

51. **Token 51**:
   - Fait 1: [1] TestPerson{name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true, status=active, department=marketing}
   - Fait 2: [1] TestOrder{discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal}

52. **Token 52**:
   - Fait 1: [1] TestPerson{age=0, active=true, tags=test, status=active, name=Frank, department=qa, salary=-5000, score=0, level=1}
   - Fait 2: [1] TestOrder{region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0}

53. **Token 53**:
   - Fait 1: [1] TestPerson{score=5.5, salary=25000, department=support, level=1, age=18, name=Henry, tags=junior, status=inactive, active=false}
   - Fait 2: [1] TestOrder{region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0}

54. **Token 54**:
   - Fait 1: [1] TestPerson{active=true, score=6.5, status=active, level=1, salary=28000, department=intern, name=X, age=22, tags=temp}
   - Fait 2: [1] TestOrder{discount=0, region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000}

55. **Token 55**:
   - Fait 1: [1] TestPerson{active=true, tags=junior, status=active, score=8.5, level=2, age=25, department=sales, name=Alice, salary=45000}
   - Fait 2: [1] TestOrder{status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west, total=999.98}

56. **Token 56**:
   - Fait 1: [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}
   - Fait 2: [1] TestOrder{status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98}

57. **Token 57**:
   - Fait 1: [1] TestPerson{salary=55000, department=sales, age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2}

58. **Token 58**:
   - Fait 1: [1] TestPerson{age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager, score=7.8, level=7}
   - Fait 2: [1] TestOrder{customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north}

59. **Token 59**:
   - Fait 1: [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}
   - Fait 2: [1] TestOrder{status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0, region=east, date=2024-03-15, priority=urgent}

60. **Token 60**:
   - Fait 1: [1] TestPerson{tags=manager, score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000}
   - Fait 2: [1] TestOrder{amount=1, total=75000, discount=0, region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001}

61. **Token 61**:
   - Fait 1: [1] TestPerson{salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active, name=Frank, department=qa}
   - Fait 2: [1] TestOrder{date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending, region=north, discount=50}

62. **Token 62**:
   - Fait 1: [1] TestPerson{age=18, name=Henry, tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support, level=1}
   - Fait 2: [1] TestOrder{amount=2, total=1999.98, status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001}

63. **Token 63**:
   - Fait 1: [1] TestPerson{salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2}

64. **Token 64**:
   - Fait 1: [1] TestPerson{active=true, age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north}

65. **Token 65**:
   - Fait 1: [1] TestPerson{age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true}
   - Fait 2: [1] TestOrder{status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0, region=east, date=2024-03-15, priority=urgent}

66. **Token 66**:
   - Fait 1: [1] TestPerson{department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000, age=65}
   - Fait 2: [1] TestOrder{date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending, region=north, discount=50}

67. **Token 67**:
   - Fait 1: [1] TestPerson{level=6, tags=senior, department=engineering, active=true, age=40, status=active, name=Ivy, salary=68000, score=8.7}
   - Fait 2: [1] TestOrder{discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal}

68. **Token 68**:
   - Fait 1: [1] TestPerson{level=5, age=35, status=active, score=9.2, name=Bob, salary=75000, active=true, tags=senior, department=engineering}
   - Fait 2: [1] TestOrder{status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west, total=999.98}

69. **Token 69**:
   - Fait 1: [1] TestPerson{status=active, score=9.2, name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35}
   - Fait 2: [1] TestOrder{status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0, region=east, date=2024-03-15, priority=urgent}

70. **Token 70**:
   - Fait 1: [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2}

71. **Token 71**:
   - Fait 1: [1] TestPerson{level=6, tags=senior, department=engineering, active=true, age=40, status=active, name=Ivy, salary=68000, score=8.7}
   - Fait 2: [1] TestOrder{status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225}

72. **Token 72**:
   - Fait 1: [1] TestPerson{status=active, department=marketing, name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true}
   - Fait 2: [1] TestOrder{discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99}

73. **Token 73**:
   - Fait 1: [1] TestPerson{age=65, department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000}
   - Fait 2: [1] TestOrder{priority=low, discount=0, region=west, total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005}

74. **Token 74**:
   - Fait 1: [1] TestPerson{name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true, status=active, department=marketing}
   - Fait 2: [1] TestOrder{status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10}

75. **Token 75**:
   - Fait 1: [1] TestPerson{name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true, age=40, status=active}
   - Fait 2: [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}

76. **Token 76**:
   - Fait 1: [1] TestPerson{salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active, name=Frank, department=qa}
   - Fait 2: [1] TestOrder{status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west, total=999.98}

77. **Token 77**:
   - Fait 1: [1] TestPerson{department=support, level=1, age=18, name=Henry, tags=junior, status=inactive, active=false, score=5.5, salary=25000}
   - Fait 2: [1] TestOrder{customer_id=P005, priority=low, discount=0, region=west, total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15}

78. **Token 78**:
   - Fait 1: [1] TestPerson{age=22, tags=temp, active=true, score=6.5, status=active, level=1, salary=28000, department=intern, name=X}
   - Fait 2: [1] TestOrder{status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west, total=999.98}

79. **Token 79**:
   - Fait 1: [1] TestPerson{status=active, department=marketing, name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true}
   - Fait 2: [1] TestOrder{status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99}

80. **Token 80**:
   - Fait 1: [1] TestPerson{status=active, score=8.5, level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2}

81. **Token 81**:
   - Fait 1: [1] TestPerson{name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true, status=active, department=marketing}
   - Fait 2: [1] TestOrder{status=confirmed, region=south, discount=0, total=25.5, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1}

82. **Token 82**:
   - Fait 1: [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}
   - Fait 2: [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west}

83. **Token 83**:
   - Fait 1: [1] TestPerson{tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales, age=30, score=8}
   - Fait 2: [1] TestOrder{product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west, total=999.98, status=cancelled}

84. **Token 84**:
   - Fait 1: [1] TestPerson{age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true}
   - Fait 2: [1] TestOrder{status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99}

85. **Token 85**:
   - Fait 1: [1] TestPerson{name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true, status=active, department=marketing}
   - Fait 2: [1] TestOrder{product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225, status=shipped, discount=15, region=north}

86. **Token 86**:
   - Fait 1: [1] TestPerson{age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true}
   - Fait 2: [1] TestOrder{priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed}

87. **Token 87**:
   - Fait 1: [1] TestPerson{salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true, age=40, status=active, name=Ivy}
   - Fait 2: [1] TestOrder{date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255}

88. **Token 88**:
   - Fait 1: [1] TestPerson{department=qa, salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active, name=Frank}
   - Fait 2: [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}

89. **Token 89**:
   - Fait 1: [1] TestPerson{active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}

90. **Token 90**:
   - Fait 1: [1] TestPerson{status=active, level=1, salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5}
   - Fait 2: [1] TestOrder{status=confirmed, region=south, discount=0, total=25.5, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1}

91. **Token 91**:
   - Fait 1: [1] TestPerson{age=65, department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000}
   - Fait 2: [1] TestOrder{discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225, status=shipped}

92. **Token 92**:
   - Fait 1: [1] TestPerson{status=active, score=9.2, name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35}
   - Fait 2: [1] TestOrder{customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}

93. **Token 93**:
   - Fait 1: [1] TestPerson{score=8.7, level=6, tags=senior, department=engineering, active=true, age=40, status=active, name=Ivy, salary=68000}
   - Fait 2: [1] TestOrder{region=west, total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0}

94. **Token 94**:
   - Fait 1: [1] TestPerson{age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales}
   - Fait 2: [1] TestOrder{status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99}

95. **Token 95**:
   - Fait 1: [1] TestPerson{active=false, tags=intern, status=inactive, level=1, name=Charlie, age=16, salary=0, department=hr, score=6}
   - Fait 2: [1] TestOrder{date=2024-03-10, region=north, total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007}

96. **Token 96**:
   - Fait 1: [1] TestPerson{tags=executive, salary=95000, age=65, department=management, score=10, status=active, level=9, name=Grace, active=true}
   - Fait 2: [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}

97. **Token 97**:
   - Fait 1: [1] TestPerson{tags=test, status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0, active=true}
   - Fait 2: [1] TestOrder{status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225}

98. **Token 98**:
   - Fait 1: [1] TestPerson{salary=25000, department=support, level=1, age=18, name=Henry, tags=junior, status=inactive, active=false, score=5.5}
   - Fait 2: [1] TestOrder{discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225, status=shipped}

99. **Token 99**:
   - Fait 1: [1] TestPerson{active=true, score=6.5, status=active, level=1, salary=28000, department=intern, name=X, age=22, tags=temp}
   - Fait 2: [1] TestOrder{product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225, status=shipped, discount=15, region=north}

100. **Token 100**:
   - Fait 1: [1] TestPerson{tags=intern, status=inactive, level=1, name=Charlie, age=16, salary=0, department=hr, score=6, active=false}
   - Fait 2: [1] TestOrder{discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}
2. [1] TestPerson{status=active, score=9.2, name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35}
3. [1] TestPerson{salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1, name=Charlie, age=16}
4. [1] TestPerson{age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager, score=7.8, level=7}
5. [1] TestPerson{active=false, name=Eve, level=3, salary=55000, department=sales, age=30, score=8, tags=employee, status=inactive}
6. [1] TestPerson{tags=test, status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0, active=true}
7. [1] TestPerson{department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000, age=65}
8. [1] TestPerson{active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry, tags=junior, status=inactive}
9. [1] TestPerson{age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true}
10. [1] TestPerson{active=true, score=6.5, status=active, level=1, salary=28000, department=intern, name=X, age=22, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{level=6, tags=senior, department=engineering, active=true, age=40, status=active, name=Ivy, salary=68000, score=8.7}

2. **Token 2**:
   - Fait 1: [1] TestPerson{salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1}

3. **Token 3**:
   - Fait 1: [1] TestPerson{tags=junior, status=active, score=8.5, level=2, age=25, department=sales, name=Alice, salary=45000, active=true}

4. **Token 4**:
   - Fait 1: [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}

5. **Token 5**:
   - Fait 1: [1] TestPerson{department=marketing, name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true, status=active}

6. **Token 6**:
   - Fait 1: [1] TestPerson{status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0, active=true, tags=test}

7. **Token 7**:
   - Fait 1: [1] TestPerson{status=active, level=9, name=Grace, active=true, tags=executive, salary=95000, age=65, department=management, score=10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}
2. [1] TestPerson{salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2, name=Bob}
3. [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}
4. [1] TestPerson{department=marketing, name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true, status=active}
5. [1] TestPerson{age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales}
6. [1] TestPerson{department=qa, salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active, name=Frank}
7. [1] TestPerson{score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000, age=65, department=management}
8. [1] TestPerson{active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry, tags=junior, status=inactive}
9. [1] TestPerson{age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true}
10. [1] TestPerson{level=1, salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active}
11. [1] TestOrder{date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending, region=north, discount=50}
12. [1] TestOrder{discount=0, total=25.5, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south}
13. [1] TestOrder{total=225, status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3}
14. [1] TestOrder{discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal}
15. [1] TestOrder{discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99}
16. [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west}
17. [1] TestOrder{status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, priority=urgent}
18. [1] TestOrder{status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10}
19. [1] TestOrder{total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north}
20. [1] TestOrder{discount=0, region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{level=6, tags=senior, department=engineering, active=true, age=40, status=active, name=Ivy, salary=68000, score=8.7}
   - Fait 2: [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west}

2. **Token 2**:
   - Fait 1: [1] TestPerson{name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1, salary=28000, department=intern}
   - Fait 2: [1] TestOrder{customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south}

3. **Token 3**:
   - Fait 1: [1] TestPerson{age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales}
   - Fait 2: [1] TestOrder{region=north, total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10}

4. **Token 4**:
   - Fait 1: [1] TestPerson{name=Henry, tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support, level=1, age=18}
   - Fait 2: [1] TestOrder{product_id=PROD001, amount=2, total=1999.98, status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001}

5. **Token 5**:
   - Fait 1: [1] TestPerson{level=1, name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive}
   - Fait 2: [1] TestOrder{region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending}

6. **Token 6**:
   - Fait 1: [1] TestPerson{name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true, age=40, status=active}
   - Fait 2: [1] TestOrder{discount=0, total=25.5, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south}

7. **Token 7**:
   - Fait 1: [1] TestPerson{age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true}
   - Fait 2: [1] TestOrder{region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestPerson{tags=senior, department=engineering, level=5, age=35, status=active, score=9.2, name=Bob, salary=75000, active=true}
   - Fait 2: [1] TestOrder{discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99}

9. **Token 9**:
   - Fait 1: [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
   - Fait 2: [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west}

10. **Token 10**:
   - Fait 1: [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}
   - Fait 2: [1] TestOrder{customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north}

11. **Token 11**:
   - Fait 1: [1] TestPerson{status=inactive, active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry, tags=junior}
   - Fait 2: [1] TestOrder{customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north}

12. **Token 12**:
   - Fait 1: [1] TestPerson{name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true, age=40, status=active}
   - Fait 2: [1] TestOrder{date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending, region=north, discount=50}

13. **Token 13**:
   - Fait 1: [1] TestPerson{tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry}
   - Fait 2: [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}

14. **Token 14**:
   - Fait 1: [1] TestPerson{active=false, tags=intern, status=inactive, level=1, name=Charlie, age=16, salary=0, department=hr, score=6}
   - Fait 2: [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}

15. **Token 15**:
   - Fait 1: [1] TestPerson{name=Henry, tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support, level=1, age=18}
   - Fait 2: [1] TestOrder{region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0}

16. **Token 16**:
   - Fait 1: [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}
   - Fait 2: [1] TestOrder{discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal}

17. **Token 17**:
   - Fait 1: [1] TestPerson{score=8.7, level=6, tags=senior, department=engineering, active=true, age=40, status=active, name=Ivy, salary=68000}
   - Fait 2: [1] TestOrder{customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north}

18. **Token 18**:
   - Fait 1: [1] TestPerson{level=9, name=Grace, active=true, tags=executive, salary=95000, age=65, department=management, score=10, status=active}
   - Fait 2: [1] TestOrder{amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99, status=completed, priority=low}

19. **Token 19**:
   - Fait 1: [1] TestPerson{tags=test, status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0, active=true}
   - Fait 2: [1] TestOrder{status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0, region=east, date=2024-03-15, priority=urgent}

20. **Token 20**:
   - Fait 1: [1] TestPerson{age=35, status=active, score=9.2, name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5}
   - Fait 2: [1] TestOrder{date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending, region=north, discount=50}

21. **Token 21**:
   - Fait 1: [1] TestPerson{level=1, age=18, name=Henry, tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support}
   - Fait 2: [1] TestOrder{date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001}

22. **Token 22**:
   - Fait 1: [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}
   - Fait 2: [1] TestOrder{discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99}

23. **Token 23**:
   - Fait 1: [1] TestPerson{active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{customer_id=P005, priority=low, discount=0, region=west, total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15}

24. **Token 24**:
   - Fait 1: [1] TestPerson{level=1, name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive}
   - Fait 2: [1] TestOrder{region=west, total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0}

25. **Token 25**:
   - Fait 1: [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
   - Fait 2: [1] TestOrder{product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north, customer_id=P007}

26. **Token 26**:
   - Fait 1: [1] TestPerson{salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active, name=Frank, department=qa}
   - Fait 2: [1] TestOrder{date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255}

27. **Token 27**:
   - Fait 1: [1] TestPerson{salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active, name=Frank, department=qa}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south}

28. **Token 28**:
   - Fait 1: [1] TestPerson{status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0, active=true, tags=test}
   - Fait 2: [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west}

29. **Token 29**:
   - Fait 1: [1] TestPerson{active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry, tags=junior, status=inactive}
   - Fait 2: [1] TestOrder{amount=10, status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal}

30. **Token 30**:
   - Fait 1: [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}
   - Fait 2: [1] TestOrder{amount=10, status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal}

31. **Token 31**:
   - Fait 1: [1] TestPerson{salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1}
   - Fait 2: [1] TestOrder{status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0, region=east, date=2024-03-15, priority=urgent}

32. **Token 32**:
   - Fait 1: [1] TestPerson{active=true, tags=executive, salary=95000, age=65, department=management, score=10, status=active, level=9, name=Grace}
   - Fait 2: [1] TestOrder{customer_id=P001, amount=3, total=225, status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high}

33. **Token 33**:
   - Fait 1: [1] TestPerson{name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1, salary=28000, department=intern}
   - Fait 2: [1] TestOrder{priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending, region=north, discount=50, date=2024-01-15}

34. **Token 34**:
   - Fait 1: [1] TestPerson{status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0, active=true, tags=test}
   - Fait 2: [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}

35. **Token 35**:
   - Fait 1: [1] TestPerson{score=0, level=1, age=0, active=true, tags=test, status=active, name=Frank, department=qa, salary=-5000}
   - Fait 2: [1] TestOrder{date=2024-02-05, status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99}

36. **Token 36**:
   - Fait 1: [1] TestPerson{salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1}
   - Fait 2: [1] TestOrder{priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006}

37. **Token 37**:
   - Fait 1: [1] TestPerson{level=1, age=18, name=Henry, tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support}
   - Fait 2: [1] TestOrder{status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0, region=east, date=2024-03-15, priority=urgent}

38. **Token 38**:
   - Fait 1: [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}
   - Fait 2: [1] TestOrder{status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0, region=east, date=2024-03-15, priority=urgent}

39. **Token 39**:
   - Fait 1: [1] TestPerson{salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active, name=Frank, department=qa}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2}

40. **Token 40**:
   - Fait 1: [1] TestPerson{score=6.5, status=active, level=1, salary=28000, department=intern, name=X, age=22, tags=temp, active=true}
   - Fait 2: [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}

41. **Token 41**:
   - Fait 1: [1] TestPerson{status=active, level=1, salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5}
   - Fait 2: [1] TestOrder{status=delivered, priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05}

42. **Token 42**:
   - Fait 1: [1] TestPerson{salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active, name=Frank, department=qa}
   - Fait 2: [1] TestOrder{date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600}

43. **Token 43**:
   - Fait 1: [1] TestPerson{salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2, name=Bob}
   - Fait 2: [1] TestOrder{total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002}

44. **Token 44**:
   - Fait 1: [1] TestPerson{status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true, age=40}
   - Fait 2: [1] TestOrder{status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0, region=east, date=2024-03-15, priority=urgent}

45. **Token 45**:
   - Fait 1: [1] TestPerson{name=Eve, level=3, salary=55000, department=sales, age=30, score=8, tags=employee, status=inactive, active=false}
   - Fait 2: [1] TestOrder{status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225}

46. **Token 46**:
   - Fait 1: [1] TestPerson{active=true, score=6.5, status=active, level=1, salary=28000, department=intern, name=X, age=22, tags=temp}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south}

47. **Token 47**:
   - Fait 1: [1] TestPerson{salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1}
   - Fait 2: [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west}

48. **Token 48**:
   - Fait 1: [1] TestPerson{tags=senior, department=engineering, active=true, age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6}
   - Fait 2: [1] TestOrder{customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south}

49. **Token 49**:
   - Fait 1: [1] TestPerson{salary=85000, tags=manager, score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana}
   - Fait 2: [1] TestOrder{status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99}

50. **Token 50**:
   - Fait 1: [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}
   - Fait 2: [1] TestOrder{status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99}

51. **Token 51**:
   - Fait 1: [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
   - Fait 2: [1] TestOrder{discount=0, region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000}

52. **Token 52**:
   - Fait 1: [1] TestPerson{age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true}
   - Fait 2: [1] TestOrder{status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225}

53. **Token 53**:
   - Fait 1: [1] TestPerson{tags=manager, score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000}
   - Fait 2: [1] TestOrder{customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high}

54. **Token 54**:
   - Fait 1: [1] TestPerson{tags=junior, status=active, score=8.5, level=2, age=25, department=sales, name=Alice, salary=45000, active=true}
   - Fait 2: [1] TestOrder{priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed}

55. **Token 55**:
   - Fait 1: [1] TestPerson{score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west}

56. **Token 56**:
   - Fait 1: [1] TestPerson{department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5, level=2, age=25}
   - Fait 2: [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west}

57. **Token 57**:
   - Fait 1: [1] TestPerson{status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales, age=30, score=8, tags=employee}
   - Fait 2: [1] TestOrder{customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south}

58. **Token 58**:
   - Fait 1: [1] TestPerson{tags=temp, active=true, score=6.5, status=active, level=1, salary=28000, department=intern, name=X, age=22}
   - Fait 2: [1] TestOrder{amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99, status=completed, priority=low}

59. **Token 59**:
   - Fait 1: [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
   - Fait 2: [1] TestOrder{status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225}

60. **Token 60**:
   - Fait 1: [1] TestPerson{salary=85000, tags=manager, score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana}
   - Fait 2: [1] TestOrder{total=25.5, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0}

61. **Token 61**:
   - Fait 1: [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}
   - Fait 2: [1] TestOrder{customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5, date=2024-01-20, priority=low}

62. **Token 62**:
   - Fait 1: [1] TestPerson{score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}

63. **Token 63**:
   - Fait 1: [1] TestPerson{tags=junior, status=active, score=8.5, level=2, age=25, department=sales, name=Alice, salary=45000, active=true}
   - Fait 2: [1] TestOrder{product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east, customer_id=P004}

64. **Token 64**:
   - Fait 1: [1] TestPerson{age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales}
   - Fait 2: [1] TestOrder{region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0}

65. **Token 65**:
   - Fait 1: [1] TestPerson{name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true, status=active, department=marketing}
   - Fait 2: [1] TestOrder{product_id=PROD001, amount=2, total=1999.98, status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001}

66. **Token 66**:
   - Fait 1: [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}
   - Fait 2: [1] TestOrder{region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending}

67. **Token 67**:
   - Fait 1: [1] TestPerson{age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5, level=2}
   - Fait 2: [1] TestOrder{region=north, customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01}

68. **Token 68**:
   - Fait 1: [1] TestPerson{name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true, status=active, department=marketing}
   - Fait 2: [1] TestOrder{customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north}

69. **Token 69**:
   - Fait 1: [1] TestPerson{name=Grace, active=true, tags=executive, salary=95000, age=65, department=management, score=10, status=active, level=9}
   - Fait 2: [1] TestOrder{status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10}

70. **Token 70**:
   - Fait 1: [1] TestPerson{level=1, age=18, name=Henry, tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support}
   - Fait 2: [1] TestOrder{customer_id=P001, amount=3, total=225, status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high}

71. **Token 71**:
   - Fait 1: [1] TestPerson{active=false, tags=intern, status=inactive, level=1, name=Charlie, age=16, salary=0, department=hr, score=6}
   - Fait 2: [1] TestOrder{status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225}

72. **Token 72**:
   - Fait 1: [1] TestPerson{name=Frank, department=qa, salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active}
   - Fait 2: [1] TestOrder{date=2024-03-10, region=north, total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007}

73. **Token 73**:
   - Fait 1: [1] TestPerson{active=true, tags=executive, salary=95000, age=65, department=management, score=10, status=active, level=9, name=Grace}
   - Fait 2: [1] TestOrder{region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0}

74. **Token 74**:
   - Fait 1: [1] TestPerson{department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000, age=65}
   - Fait 2: [1] TestOrder{priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5, date=2024-01-20}

75. **Token 75**:
   - Fait 1: [1] TestPerson{salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active, name=Frank, department=qa}
   - Fait 2: [1] TestOrder{priority=high, customer_id=P001, amount=3, total=225, status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01}

76. **Token 76**:
   - Fait 1: [1] TestPerson{age=65, department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000}
   - Fait 2: [1] TestOrder{priority=normal, discount=0, region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered}

77. **Token 77**:
   - Fait 1: [1] TestPerson{tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support, level=1, age=18, name=Henry}
   - Fait 2: [1] TestOrder{status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99}

78. **Token 78**:
   - Fait 1: [1] TestPerson{status=inactive, level=1, name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern}
   - Fait 2: [1] TestOrder{status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99}

79. **Token 79**:
   - Fait 1: [1] TestPerson{age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, region=north, discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2}

80. **Token 80**:
   - Fait 1: [1] TestPerson{age=65, department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south}

81. **Token 81**:
   - Fait 1: [1] TestPerson{name=Grace, active=true, tags=executive, salary=95000, age=65, department=management, score=10, status=active, level=9}
   - Fait 2: [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west}

82. **Token 82**:
   - Fait 1: [1] TestPerson{department=sales, age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000}
   - Fait 2: [1] TestOrder{amount=4, total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50}

83. **Token 83**:
   - Fait 1: [1] TestPerson{level=3, salary=55000, department=sales, age=30, score=8, tags=employee, status=inactive, active=false, name=Eve}
   - Fait 2: [1] TestOrder{discount=0, total=25.5, date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south}

84. **Token 84**:
   - Fait 1: [1] TestPerson{name=Eve, level=3, salary=55000, department=sales, age=30, score=8, tags=employee, status=inactive, active=false}
   - Fait 2: [1] TestOrder{customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}

85. **Token 85**:
   - Fait 1: [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
   - Fait 2: [1] TestOrder{amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99, status=completed, priority=low}

86. **Token 86**:
   - Fait 1: [1] TestPerson{score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager}
   - Fait 2: [1] TestOrder{status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0, region=east, date=2024-03-15, priority=urgent}

87. **Token 87**:
   - Fait 1: [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}
   - Fait 2: [1] TestOrder{discount=0, region=east, date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000}

88. **Token 88**:
   - Fait 1: [1] TestPerson{name=Grace, active=true, tags=executive, salary=95000, age=65, department=management, score=10, status=active, level=9}
   - Fait 2: [1] TestOrder{date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending, region=north, discount=50}

89. **Token 89**:
   - Fait 1: [1] TestPerson{active=true, score=6.5, status=active, level=1, salary=28000, department=intern, name=X, age=22, tags=temp}
   - Fait 2: [1] TestOrder{status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225}

90. **Token 90**:
   - Fait 1: [1] TestPerson{salary=55000, department=sales, age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south}

91. **Token 91**:
   - Fait 1: [1] TestPerson{tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales, age=30, score=8}
   - Fait 2: [1] TestOrder{customer_id=P005, priority=low, discount=0, region=west, total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15}

92. **Token 92**:
   - Fait 1: [1] TestPerson{name=Grace, active=true, tags=executive, salary=95000, age=65, department=management, score=10, status=active, level=9}
   - Fait 2: [1] TestOrder{priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north, customer_id=P007, product_id=PROD006}

93. **Token 93**:
   - Fait 1: [1] TestPerson{active=true, tags=junior, status=active, score=8.5, level=2, age=25, department=sales, name=Alice, salary=45000}
   - Fait 2: [1] TestOrder{total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002}

94. **Token 94**:
   - Fait 1: [1] TestPerson{name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true, status=active, department=marketing}
   - Fait 2: [1] TestOrder{amount=10, status=pending, discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal}

95. **Token 95**:
   - Fait 1: [1] TestPerson{department=engineering, active=true, age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior}
   - Fait 2: [1] TestOrder{status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north, total=89.99}

96. **Token 96**:
   - Fait 1: [1] TestPerson{score=9.2, name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active}
   - Fait 2: [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}

97. **Token 97**:
   - Fait 1: [1] TestPerson{status=active, department=marketing, name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true}
   - Fait 2: [1] TestOrder{product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225, status=shipped, discount=15, region=north}

98. **Token 98**:
   - Fait 1: [1] TestPerson{salary=45000, active=true, tags=junior, status=active, score=8.5, level=2, age=25, department=sales, name=Alice}
   - Fait 2: [1] TestOrder{customer_id=P001, amount=3, total=225, status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high}

99. **Token 99**:
   - Fait 1: [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
   - Fait 2: [1] TestOrder{region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0}

100. **Token 100**:
   - Fait 1: [1] TestPerson{age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true}
   - Fait 2: [1] TestOrder{amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}
2. [1] TestPerson{active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2, name=Bob, salary=75000}
3. [1] TestPerson{salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1, name=Charlie, age=16}
4. [1] TestPerson{name=Diana, salary=85000, tags=manager, score=7.8, level=7, age=45, active=true, status=active, department=marketing}
5. [1] TestPerson{age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales}
6. [1] TestPerson{status=active, name=Frank, department=qa, salary=-5000, score=0, level=1, age=0, active=true, tags=test}
7. [1] TestPerson{age=65, department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000}
8. [1] TestPerson{score=5.5, salary=25000, department=support, level=1, age=18, name=Henry, tags=junior, status=inactive, active=false}
9. [1] TestPerson{level=6, tags=senior, department=engineering, active=true, age=40, status=active, name=Ivy, salary=68000, score=8.7}
10. [1] TestPerson{status=active, level=1, salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{active=false, name=Eve, level=3, salary=55000, department=sales, age=30, score=8, tags=employee, status=inactive}

2. **Token 2**:
   - Fait 1: [1] TestPerson{age=65, department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000}

3. **Token 3**:
   - Fait 1: [1] TestPerson{level=1, age=18, name=Henry, tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support}

4. **Token 4**:
   - Fait 1: [1] TestPerson{name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true, age=40, status=active}

5. **Token 5**:
   - Fait 1: [1] TestPerson{active=true, score=6.5, status=active, level=1, salary=28000, department=intern, name=X, age=22, tags=temp}

6. **Token 6**:
   - Fait 1: [1] TestPerson{department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5, level=2, age=25}

7. **Token 7**:
   - Fait 1: [1] TestPerson{score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager}

8. **Token 8**:
   - Fait 1: [1] TestPerson{active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2, name=Bob, salary=75000}

9. **Token 9**:
   - Fait 1: [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending, region=north, discount=50}
2. [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}
3. [1] TestOrder{amount=3, total=225, status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001}
4. [1] TestOrder{customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0, region=east}
5. [1] TestOrder{discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1, status=confirmed, priority=high, customer_id=P002, total=999.99}
6. [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north}
8. [1] TestOrder{discount=0, region=south, customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10, status=pending}
9. [1] TestOrder{total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, region=north}
10. [1] TestOrder{date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{discount=50, date=2024-01-15, priority=normal, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, status=pending, region=north}

2. **Token 2**:
   - Fait 1: [1] TestOrder{status=shipped, discount=15, region=north, product_id=PROD003, date=2024-02-01, priority=high, customer_id=P001, amount=3, total=225}

3. **Token 3**:
   - Fait 1: [1] TestOrder{region=east, customer_id=P004, product_id=PROD004, amount=1, total=299.99, date=2024-02-05, status=delivered, priority=normal, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{product_id=PROD006, priority=urgent, status=shipped, discount=50, amount=4, total=600, date=2024-03-01, region=north, customer_id=P007}

5. **Token 5**:
   - Fait 1: [1] TestOrder{customer_id=P010, product_id=PROD002, total=255, date=2024-03-05, priority=normal, amount=10, status=pending, discount=0, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{date=2024-03-15, priority=urgent, status=refunded, customer_id=P006, product_id=PROD001, amount=1, total=75000, discount=0, region=east}

7. **Token 7**:
   - Fait 1: [1] TestOrder{date=2024-01-20, priority=low, customer_id=P002, product_id=PROD002, amount=1, status=confirmed, region=south, discount=0, total=25.5}

8. **Token 8**:
   - Fait 1: [1] TestOrder{status=confirmed, priority=high, customer_id=P002, total=999.99, discount=100, product_id=PROD001, date=2024-02-10, region=south, amount=1}

9. **Token 9**:
   - Fait 1: [1] TestOrder{total=999.98, status=cancelled, product_id=PROD005, amount=2, date=2024-02-15, customer_id=P005, priority=low, discount=0, region=west}

10. **Token 10**:
   - Fait 1: [1] TestOrder{date=2024-03-10, region=north, total=89.99, status=completed, priority=low, amount=1, discount=10, customer_id=P001, product_id=PROD007}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5, level=2, age=25, department=sales}
2. [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}
3. [1] TestPerson{level=1, name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive}
4. [1] TestPerson{score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000, tags=manager}
5. [1] TestPerson{age=30, score=8, tags=employee, status=inactive, active=false, name=Eve, level=3, salary=55000, department=sales}
6. [1] TestPerson{salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active, name=Frank, department=qa}
7. [1] TestPerson{name=Grace, active=true, tags=executive, salary=95000, age=65, department=management, score=10, status=active, level=9}
8. [1] TestPerson{score=5.5, salary=25000, department=support, level=1, age=18, name=Henry, tags=junior, status=inactive, active=false}
9. [1] TestPerson{level=6, tags=senior, department=engineering, active=true, age=40, status=active, name=Ivy, salary=68000, score=8.7}
10. [1] TestPerson{salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{salary=-5000, score=0, level=1, age=0, active=true, tags=test, status=active, name=Frank, department=qa}

2. **Token 2**:
   - Fait 1: [1] TestPerson{department=management, score=10, status=active, level=9, name=Grace, active=true, tags=executive, salary=95000, age=65}

3. **Token 3**:
   - Fait 1: [1] TestPerson{age=18, name=Henry, tags=junior, status=inactive, active=false, score=5.5, salary=25000, department=support, level=1}

4. **Token 4**:
   - Fait 1: [1] TestPerson{salary=28000, department=intern, name=X, age=22, tags=temp, active=true, score=6.5, status=active, level=1}

5. **Token 5**:
   - Fait 1: [1] TestPerson{level=2, age=25, department=sales, name=Alice, salary=45000, active=true, tags=junior, status=active, score=8.5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{name=Charlie, age=16, salary=0, department=hr, score=6, active=false, tags=intern, status=inactive, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{age=40, status=active, name=Ivy, salary=68000, score=8.7, level=6, tags=senior, department=engineering, active=true}

8. **Token 8**:
   - Fait 1: [1] TestPerson{name=Bob, salary=75000, active=true, tags=senior, department=engineering, level=5, age=35, status=active, score=9.2}

9. **Token 9**:
   - Fait 1: [1] TestPerson{tags=manager, score=7.8, level=7, age=45, active=true, status=active, department=marketing, name=Diana, salary=85000}

10. **Token 10**:
   - Fait 1: [1] TestPerson{name=Eve, level=3, salary=55000, department=sales, age=30, score=8, tags=employee, status=inactive, active=false}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
