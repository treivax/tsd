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

1. [1] TestPerson{active=true, age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}
2. [1] TestPerson{salary=75000, department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35}
3. [1] TestPerson{age=16, level=1, department=hr, score=6, tags=intern, status=inactive, name=Charlie, salary=0, active=false}
4. [1] TestPerson{name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager, level=7}
5. [1] TestPerson{active=false, tags=employee, score=8, level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales}
6. [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
7. [1] TestPerson{active=true, tags=executive, age=65, level=9, name=Grace, salary=95000, status=active, department=management, score=10}
8. [1] TestPerson{age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support}
9. [1] TestPerson{department=engineering, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy}
10. [1] TestPerson{name=X, salary=28000, active=true, level=1, score=6.5, status=active, department=intern, age=22, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{level=7, name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager}

2. **Token 2**:
   - Fait 1: [1] TestPerson{age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8, level=3}

3. **Token 3**:
   - Fait 1: [1] TestPerson{name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive, age=65, level=9}

4. **Token 4**:
   - Fait 1: [1] TestPerson{age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support}

5. **Token 5**:
   - Fait 1: [1] TestPerson{active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering, salary=68000}

6. **Token 6**:
   - Fait 1: [1] TestPerson{active=true, age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}

7. **Token 7**:
   - Fait 1: [1] TestPerson{age=16, level=1, department=hr, score=6, tags=intern, status=inactive, name=Charlie, salary=0, active=false}

8. **Token 8**:
   - Fait 1: [1] TestPerson{name=X, salary=28000, active=true, level=1, score=6.5, status=active, department=intern, age=22, tags=temp}

9. **Token 9**:
   - Fait 1: [1] TestPerson{tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering, name=Bob, active=true}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001}
2. [1] TestOrder{status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20}
3. [1] TestOrder{amount=3, customer_id=P001, date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003}
4. [1] TestOrder{discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99}
5. [1] TestOrder{customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1}
6. [1] TestOrder{product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005, priority=low}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped}
8. [1] TestOrder{total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10, product_id=PROD002}
9. [1] TestOrder{product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10}
10. [1] TestOrder{date=2024-03-15, amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001}

2. **Token 2**:
   - Fait 1: [1] TestOrder{customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5}

3. **Token 3**:
   - Fait 1: [1] TestOrder{total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001, date=2024-02-01, discount=15}

4. **Token 4**:
   - Fait 1: [1] TestOrder{discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99}

5. **Token 5**:
   - Fait 1: [1] TestOrder{customer_id=P010, priority=normal, region=south, amount=10, product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0}

6. **Token 6**:
   - Fait 1: [1] TestOrder{priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10, product_id=PROD007, status=completed}

7. **Token 7**:
   - Fait 1: [1] TestOrder{amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15}

8. **Token 8**:
   - Fait 1: [1] TestOrder{amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high}

9. **Token 9**:
   - Fait 1: [1] TestOrder{product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{status=active, department=sales, active=true, age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice}
2. [1] TestPerson{name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering}
3. [1] TestPerson{level=1, department=hr, score=6, tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16}
4. [1] TestPerson{status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true}
5. [1] TestPerson{age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8, level=3}
6. [1] TestPerson{score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test, salary=-5000, active=true}
7. [1] TestPerson{name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive, age=65, level=9}
8. [1] TestPerson{active=false, status=inactive, department=support, age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1}
9. [1] TestPerson{salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering}
10. [1] TestPerson{name=X, salary=28000, active=true, level=1, score=6.5, status=active, department=intern, age=22, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45}

2. **Token 2**:
   - Fait 1: [1] TestPerson{score=8, level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee}

3. **Token 3**:
   - Fait 1: [1] TestPerson{age=65, level=9, name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive}

4. **Token 4**:
   - Fait 1: [1] TestPerson{age=40, tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6}

5. **Token 5**:
   - Fait 1: [1] TestPerson{tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000, level=2, score=8.5}

6. **Token 6**:
   - Fait 1: [1] TestPerson{department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/10 (60.0%)
- **Efficacité**: 🟠 Moyenne

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001, discount=50}
2. [1] TestOrder{region=south, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1}
3. [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}
4. [1] TestOrder{total=299.99, discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal}
5. [1] TestOrder{region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001}
6. [1] TestOrder{date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98}
7. [1] TestOrder{customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped}
8. [1] TestOrder{total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10, product_id=PROD002}
9. [1] TestOrder{date=2024-03-10, product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99}
10. [1] TestOrder{discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15, amount=1, status=refunded}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001, date=2024-02-01, discount=15, total=225, status=shipped}

2. **Token 2**:
   - Fait 1: [1] TestOrder{total=299.99, discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal}

3. **Token 3**:
   - Fait 1: [1] TestOrder{region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001}

4. **Token 4**:
   - Fait 1: [1] TestOrder{region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001, discount=50, amount=2, date=2024-01-15}

5. **Token 5**:
   - Fait 1: [1] TestOrder{date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south}

6. **Token 6**:
   - Fait 1: [1] TestOrder{total=999.98, date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled}

7. **Token 7**:
   - Fait 1: [1] TestOrder{region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600}

8. **Token 8**:
   - Fait 1: [1] TestOrder{amount=10, product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south}

9. **Token 9**:
   - Fait 1: [1] TestOrder{product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{rating=4.5, price=999.99, keywords=computer, stock=50, category=electronics, available=true, brand=TechCorp, supplier=TechSupply, name=Laptop}
2. [1] TestProduct{category=accessories, available=true, rating=4, stock=200, price=25.5, name=Mouse, supplier=TechSupply, keywords=peripheral, brand=TechCorp}
3. [1] TestProduct{name=Keyboard, stock=0, supplier=KeySupply, keywords=typing, brand=KeyTech, available=false, rating=3.5, category=accessories, price=75}
4. [1] TestProduct{name=Monitor, rating=4.8, keywords=display, category=electronics, available=true, brand=ScreenPro, stock=30, price=299.99, supplier=ScreenSupply}
5. [1] TestProduct{available=false, rating=2, brand=OldTech, stock=0, supplier=OldSupply, category=accessories, price=8.5, name=OldKeyboard, keywords=obsolete}
6. [1] TestProduct{category=audio, price=150, brand=AudioMax, stock=75, available=true, supplier=AudioSupply, keywords=sound, rating=4.6, name=Headphones}
7. [1] TestProduct{name=Webcam, category=electronics, brand=CamTech, stock=25, supplier=CamSupply, price=89.99, rating=3.8, available=true, keywords=video}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{name=Webcam, category=electronics, brand=CamTech, stock=25, supplier=CamSupply, price=89.99, rating=3.8, available=true, keywords=video}

2. **Token 2**:
   - Fait 1: [1] TestProduct{name=Laptop, rating=4.5, price=999.99, keywords=computer, stock=50, category=electronics, available=true, brand=TechCorp, supplier=TechSupply}

3. **Token 3**:
   - Fait 1: [1] TestProduct{brand=TechCorp, category=accessories, available=true, rating=4, stock=200, price=25.5, name=Mouse, supplier=TechSupply, keywords=peripheral}

4. **Token 4**:
   - Fait 1: [1] TestProduct{brand=KeyTech, available=false, rating=3.5, category=accessories, price=75, name=Keyboard, stock=0, supplier=KeySupply, keywords=typing}

5. **Token 5**:
   - Fait 1: [1] TestProduct{category=electronics, available=true, brand=ScreenPro, stock=30, price=299.99, supplier=ScreenSupply, name=Monitor, rating=4.8, keywords=display}

6. **Token 6**:
   - Fait 1: [1] TestProduct{keywords=sound, rating=4.6, name=Headphones, category=audio, price=150, brand=AudioMax, stock=75, available=true, supplier=AudioSupply}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000}
2. [1] TestPerson{tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering, name=Bob, active=true}
3. [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}
4. [1] TestPerson{level=7, name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager}
5. [1] TestPerson{level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8}
6. [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
7. [1] TestPerson{score=10, active=true, tags=executive, age=65, level=9, name=Grace, salary=95000, status=active, department=management}
8. [1] TestPerson{department=support, age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive}
9. [1] TestPerson{salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering}
10. [1] TestPerson{name=X, salary=28000, active=true, level=1, score=6.5, status=active, department=intern, age=22, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**2 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}

2. **Token 2**:
   - Fait 1: [1] TestPerson{tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 2/10 (20.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001}
2. [1] TestOrder{total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20, status=confirmed, priority=low, discount=0}
3. [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}
4. [1] TestOrder{amount=1, priority=normal, total=299.99, discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004}
5. [1] TestOrder{amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high}
6. [1] TestOrder{priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005}
7. [1] TestOrder{region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600}
8. [1] TestOrder{amount=10, product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south}
9. [1] TestOrder{product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10}
10. [1] TestOrder{date=2024-03-15, amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east, product_id=PROD004, date=2024-02-05}

3. **Token 3**:
   - Fait 1: [1] TestOrder{customer_id=P010, priority=normal, region=south, amount=10, product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0}

4. **Token 4**:
   - Fait 1: [1] TestOrder{product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10}

5. **Token 5**:
   - Fait 1: [1] TestOrder{status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001, date=2024-02-01, discount=15, total=225}

6. **Token 6**:
   - Fait 1: [1] TestOrder{region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001}

7. **Token 7**:
   - Fait 1: [1] TestOrder{region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0}

8. **Token 8**:
   - Fait 1: [1] TestOrder{region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600}

9. **Token 9**:
   - Fait 1: [1] TestOrder{customer_id=P006, date=2024-03-15, amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001}

10. **Token 10**:
   - Fait 1: [1] TestOrder{discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{active=true, age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}
2. [1] TestPerson{tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering, name=Bob, active=true}
3. [1] TestPerson{age=16, level=1, department=hr, score=6, tags=intern, status=inactive, name=Charlie, salary=0, active=false}
4. [1] TestPerson{level=7, name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager}
5. [1] TestPerson{age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8, level=3}
6. [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
7. [1] TestPerson{age=65, level=9, name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive}
8. [1] TestPerson{name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support, age=18, score=5.5}
9. [1] TestPerson{tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6, age=40}
10. [1] TestPerson{name=X, salary=28000, active=true, level=1, score=6.5, status=active, department=intern, age=22, tags=temp}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**3 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{salary=0, active=false, age=16, level=1, department=hr, score=6, tags=intern, status=inactive, name=Charlie}

2. **Token 2**:
   - Fait 1: [1] TestPerson{status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8, level=3, age=30, salary=55000}

3. **Token 3**:
   - Fait 1: [1] TestPerson{score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support, age=18}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 3/10 (30.0%)
- **Efficacité**: 🔴 Faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001}
2. [1] TestOrder{status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20}
3. [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}
4. [1] TestOrder{product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east}
5. [1] TestOrder{amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high}
6. [1] TestOrder{priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005}
7. [1] TestOrder{region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600}
8. [1] TestOrder{region=south, amount=10, product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal}
9. [1] TestOrder{product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10}
10. [1] TestOrder{date=2024-03-15, amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99, discount=0}

2. **Token 2**:
   - Fait 1: [1] TestOrder{amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high}

3. **Token 3**:
   - Fait 1: [1] TestOrder{status=cancelled, total=999.98, date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0, region=west}

4. **Token 4**:
   - Fait 1: [1] TestOrder{amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent}

5. **Token 5**:
   - Fait 1: [1] TestOrder{total=89.99, date=2024-03-10, product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1}

6. **Token 6**:
   - Fait 1: [1] TestOrder{discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15, amount=1, status=refunded}

7. **Token 7**:
   - Fait 1: [1] TestOrder{amount=1, region=south, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002}

8. **Token 8**:
   - Fait 1: [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{active=true, age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}
2. [1] TestPerson{department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000}
3. [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}
4. [1] TestPerson{level=7, name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager}
5. [1] TestPerson{active=false, tags=employee, score=8, level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales}
6. [1] TestPerson{active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test, salary=-5000}
7. [1] TestPerson{age=65, level=9, name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive}
8. [1] TestPerson{status=inactive, department=support, age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false}
9. [1] TestPerson{department=engineering, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy}
10. [1] TestPerson{score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}

2. **Token 2**:
   - Fait 1: [1] TestPerson{name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager, level=7}

3. **Token 3**:
   - Fait 1: [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}

4. **Token 4**:
   - Fait 1: [1] TestPerson{status=active, department=management, score=10, active=true, tags=executive, age=65, level=9, name=Grace, salary=95000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{name=Alice, status=active, department=sales, active=true, age=25, salary=45000, level=2, score=8.5, tags=junior}

6. **Token 6**:
   - Fait 1: [1] TestPerson{tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering, name=Bob, active=true}

7. **Token 7**:
   - Fait 1: [1] TestPerson{level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8}

8. **Token 8**:
   - Fait 1: [1] TestPerson{level=1, active=false, status=inactive, department=support, age=18, score=5.5, name=Henry, salary=25000, tags=junior}

9. **Token 9**:
   - Fait 1: [1] TestPerson{tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6, age=40}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{stock=50, category=electronics, available=true, brand=TechCorp, supplier=TechSupply, name=Laptop, rating=4.5, price=999.99, keywords=computer}
2. [1] TestProduct{category=accessories, available=true, rating=4, stock=200, price=25.5, name=Mouse, supplier=TechSupply, keywords=peripheral, brand=TechCorp}
3. [1] TestProduct{price=75, name=Keyboard, stock=0, supplier=KeySupply, keywords=typing, brand=KeyTech, available=false, rating=3.5, category=accessories}
4. [1] TestProduct{available=true, brand=ScreenPro, stock=30, price=299.99, supplier=ScreenSupply, name=Monitor, rating=4.8, keywords=display, category=electronics}
5. [1] TestProduct{category=accessories, price=8.5, name=OldKeyboard, keywords=obsolete, available=false, rating=2, brand=OldTech, stock=0, supplier=OldSupply}
6. [1] TestProduct{rating=4.6, name=Headphones, category=audio, price=150, brand=AudioMax, stock=75, available=true, supplier=AudioSupply, keywords=sound}
7. [1] TestProduct{stock=25, supplier=CamSupply, price=89.99, rating=3.8, available=true, keywords=video, name=Webcam, category=electronics, brand=CamTech}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

**6 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestProduct{available=true, keywords=video, name=Webcam, category=electronics, brand=CamTech, stock=25, supplier=CamSupply, price=89.99, rating=3.8}

2. **Token 2**:
   - Fait 1: [1] TestProduct{category=electronics, available=true, brand=TechCorp, supplier=TechSupply, name=Laptop, rating=4.5, price=999.99, keywords=computer, stock=50}

3. **Token 3**:
   - Fait 1: [1] TestProduct{keywords=peripheral, brand=TechCorp, category=accessories, available=true, rating=4, stock=200, price=25.5, name=Mouse, supplier=TechSupply}

4. **Token 4**:
   - Fait 1: [1] TestProduct{rating=3.5, category=accessories, price=75, name=Keyboard, stock=0, supplier=KeySupply, keywords=typing, brand=KeyTech, available=false}

5. **Token 5**:
   - Fait 1: [1] TestProduct{category=electronics, available=true, brand=ScreenPro, stock=30, price=299.99, supplier=ScreenSupply, name=Monitor, rating=4.8, keywords=display}

6. **Token 6**:
   - Fait 1: [1] TestProduct{category=audio, price=150, brand=AudioMax, stock=75, available=true, supplier=AudioSupply, keywords=sound, rating=4.6, name=Headphones}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 6/7 (85.7%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{active=true, age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}
2. [1] TestPerson{status=active, level=5, score=9.2, age=35, salary=75000, department=engineering, name=Bob, active=true, tags=senior}
3. [1] TestPerson{score=6, tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr}
4. [1] TestPerson{status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true}
5. [1] TestPerson{level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8}
6. [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
7. [1] TestPerson{active=true, tags=executive, age=65, level=9, name=Grace, salary=95000, status=active, department=management, score=10}
8. [1] TestPerson{active=false, status=inactive, department=support, age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1}
9. [1] TestPerson{tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6, age=40}
10. [1] TestPerson{score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{tags=employee, score=8, level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false}

2. **Token 2**:
   - Fait 1: [1] TestPerson{active=true, tags=executive, age=65, level=9, name=Grace, salary=95000, status=active, department=management, score=10}

3. **Token 3**:
   - Fait 1: [1] TestPerson{score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true}

4. **Token 4**:
   - Fait 1: [1] TestPerson{level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000}

5. **Token 5**:
   - Fait 1: [1] TestPerson{age=35, salary=75000, department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2}

6. **Token 6**:
   - Fait 1: [1] TestPerson{name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager, level=7}

7. **Token 7**:
   - Fait 1: [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}

8. **Token 8**:
   - Fait 1: [1] TestPerson{age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support}

9. **Token 9**:
   - Fait 1: [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{priority=normal, product_id=PROD001, customer_id=P001, discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending}
2. [1] TestOrder{status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20}
3. [1] TestOrder{priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001, date=2024-02-01, discount=15, total=225, status=shipped}
4. [1] TestOrder{priority=normal, total=299.99, discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1}
5. [1] TestOrder{date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99}
6. [1] TestOrder{total=999.98, date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled}
7. [1] TestOrder{status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01}
8. [1] TestOrder{product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10}
9. [1] TestOrder{status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10, product_id=PROD007}
10. [1] TestOrder{customer_id=P006, date=2024-03-15, amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**8 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east}

2. **Token 2**:
   - Fait 1: [1] TestOrder{product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100}

3. **Token 3**:
   - Fait 1: [1] TestOrder{status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01}

4. **Token 4**:
   - Fait 1: [1] TestOrder{status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10, product_id=PROD002, total=255, date=2024-03-05}

5. **Token 5**:
   - Fait 1: [1] TestOrder{region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10, product_id=PROD007, status=completed, priority=low}

6. **Token 6**:
   - Fait 1: [1] TestOrder{status=pending, priority=normal, product_id=PROD001, customer_id=P001, discount=50, amount=2, date=2024-01-15, region=north, total=1999.98}

7. **Token 7**:
   - Fait 1: [1] TestOrder{discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20, status=confirmed, priority=low}

8. **Token 8**:
   - Fait 1: [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 8/10 (80.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000, level=2}
2. [1] TestPerson{name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering}
3. [1] TestPerson{score=6, tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr}
4. [1] TestPerson{score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45, salary=85000}
5. [1] TestPerson{age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8, level=3}
6. [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
7. [1] TestPerson{salary=95000, status=active, department=management, score=10, active=true, tags=executive, age=65, level=9, name=Grace}
8. [1] TestPerson{age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support}
9. [1] TestPerson{tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6, age=40}
10. [1] TestPerson{name=X, salary=28000, active=true, level=1, score=6.5, status=active, department=intern, age=22, tags=temp}
11. [1] TestOrder{date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001, discount=50, amount=2}
12. [1] TestOrder{customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5}
13. [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}
14. [1] TestOrder{product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east}
15. [1] TestOrder{status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002}
16. [1] TestOrder{priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005}
17. [1] TestOrder{customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped}
18. [1] TestOrder{date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10, product_id=PROD002, total=255}
19. [1] TestOrder{customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10, product_id=PROD007, status=completed, priority=low, region=north}
20. [1] TestOrder{discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15, amount=1, status=refunded}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{level=1, active=false, status=inactive, department=support, age=18, score=5.5, name=Henry, salary=25000, tags=junior}
   - Fait 2: [1] TestOrder{status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20}

2. **Token 2**:
   - Fait 1: [1] TestPerson{score=10, active=true, tags=executive, age=65, level=9, name=Grace, salary=95000, status=active, department=management}
   - Fait 2: [1] TestOrder{date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south}

3. **Token 3**:
   - Fait 1: [1] TestPerson{age=22, tags=temp, name=X, salary=28000, active=true, level=1, score=6.5, status=active, department=intern}
   - Fait 2: [1] TestOrder{status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20}

4. **Token 4**:
   - Fait 1: [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}
   - Fait 2: [1] TestOrder{total=299.99, discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal}

5. **Token 5**:
   - Fait 1: [1] TestPerson{name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{status=cancelled, total=999.98, date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0, region=west}

6. **Token 6**:
   - Fait 1: [1] TestPerson{status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6, tags=intern}
   - Fait 2: [1] TestOrder{customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped}

7. **Token 7**:
   - Fait 1: [1] TestPerson{age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8, level=3}
   - Fait 2: [1] TestOrder{status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20}

8. **Token 8**:
   - Fait 1: [1] TestPerson{active=true, age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}
   - Fait 2: [1] TestOrder{customer_id=P001, discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001}

9. **Token 9**:
   - Fait 1: [1] TestPerson{name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager, level=7}
   - Fait 2: [1] TestOrder{product_id=PROD001, customer_id=P001, discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal}

10. **Token 10**:
   - Fait 1: [1] TestPerson{active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering, salary=68000}
   - Fait 2: [1] TestOrder{total=299.99, discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal}

11. **Token 11**:
   - Fait 1: [1] TestPerson{department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000}
   - Fait 2: [1] TestOrder{discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed}

12. **Token 12**:
   - Fait 1: [1] TestPerson{department=engineering, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy}
   - Fait 2: [1] TestOrder{status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01}

13. **Token 13**:
   - Fait 1: [1] TestPerson{level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000}
   - Fait 2: [1] TestOrder{date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10, product_id=PROD002, total=255}

14. **Token 14**:
   - Fait 1: [1] TestPerson{level=7, name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager}
   - Fait 2: [1] TestOrder{amount=10, product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south}

15. **Token 15**:
   - Fait 1: [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}
   - Fait 2: [1] TestOrder{discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001}

16. **Token 16**:
   - Fait 1: [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}
   - Fait 2: [1] TestOrder{product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10}

17. **Token 17**:
   - Fait 1: [1] TestPerson{name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering}
   - Fait 2: [1] TestOrder{amount=1, total=89.99, date=2024-03-10, product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10}

18. **Token 18**:
   - Fait 1: [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
   - Fait 2: [1] TestOrder{customer_id=P006, date=2024-03-15, amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001}

19. **Token 19**:
   - Fait 1: [1] TestPerson{salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering}
   - Fait 2: [1] TestOrder{discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001}

20. **Token 20**:
   - Fait 1: [1] TestPerson{age=35, salary=75000, department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2}
   - Fait 2: [1] TestOrder{amount=3, customer_id=P001, date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003}

21. **Token 21**:
   - Fait 1: [1] TestPerson{level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000}
   - Fait 2: [1] TestOrder{amount=1, priority=normal, total=299.99, discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004}

22. **Token 22**:
   - Fait 1: [1] TestPerson{tags=manager, level=7, name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing}
   - Fait 2: [1] TestOrder{product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east}

23. **Token 23**:
   - Fait 1: [1] TestPerson{level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000}
   - Fait 2: [1] TestOrder{status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01}

24. **Token 24**:
   - Fait 1: [1] TestPerson{active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana}
   - Fait 2: [1] TestOrder{region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600}

25. **Token 25**:
   - Fait 1: [1] TestPerson{status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior}
   - Fait 2: [1] TestOrder{region=south, amount=10, product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal}

26. **Token 26**:
   - Fait 1: [1] TestPerson{level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000}
   - Fait 2: [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}

27. **Token 27**:
   - Fait 1: [1] TestPerson{name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager, level=7}
   - Fait 2: [1] TestOrder{customer_id=P001, date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3}

28. **Token 28**:
   - Fait 1: [1] TestPerson{tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering, name=Bob, active=true}
   - Fait 2: [1] TestOrder{total=299.99, discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal}

29. **Token 29**:
   - Fait 1: [1] TestPerson{name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active}
   - Fait 2: [1] TestOrder{amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high}

30. **Token 30**:
   - Fait 1: [1] TestPerson{department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000}
   - Fait 2: [1] TestOrder{customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped}

31. **Token 31**:
   - Fait 1: [1] TestPerson{status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior}
   - Fait 2: [1] TestOrder{product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10}

32. **Token 32**:
   - Fait 1: [1] TestPerson{status=active, department=management, score=10, active=true, tags=executive, age=65, level=9, name=Grace, salary=95000}
   - Fait 2: [1] TestOrder{customer_id=P006, date=2024-03-15, amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001}

33. **Token 33**:
   - Fait 1: [1] TestPerson{level=1, score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true}
   - Fait 2: [1] TestOrder{customer_id=P006, date=2024-03-15, amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001}

34. **Token 34**:
   - Fait 1: [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}
   - Fait 2: [1] TestOrder{product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100}

35. **Token 35**:
   - Fait 1: [1] TestPerson{active=false, age=16, level=1, department=hr, score=6, tags=intern, status=inactive, name=Charlie, salary=0}
   - Fait 2: [1] TestOrder{region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10, product_id=PROD007, status=completed, priority=low}

36. **Token 36**:
   - Fait 1: [1] TestPerson{tags=employee, score=8, level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false}
   - Fait 2: [1] TestOrder{discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15, amount=1, status=refunded}

37. **Token 37**:
   - Fait 1: [1] TestPerson{age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support}
   - Fait 2: [1] TestOrder{customer_id=P006, date=2024-03-15, amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001}

38. **Token 38**:
   - Fait 1: [1] TestPerson{age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support}
   - Fait 2: [1] TestOrder{total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south}

39. **Token 39**:
   - Fait 1: [1] TestPerson{score=10, active=true, tags=executive, age=65, level=9, name=Grace, salary=95000, status=active, department=management}
   - Fait 2: [1] TestOrder{product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100}

40. **Token 40**:
   - Fait 1: [1] TestPerson{level=1, score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true}
   - Fait 2: [1] TestOrder{amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high}

41. **Token 41**:
   - Fait 1: [1] TestPerson{level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000}
   - Fait 2: [1] TestOrder{priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005}

42. **Token 42**:
   - Fait 1: [1] TestPerson{score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45, salary=85000}
   - Fait 2: [1] TestOrder{customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15}

43. **Token 43**:
   - Fait 1: [1] TestPerson{score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1}
   - Fait 2: [1] TestOrder{discount=10, amount=1, total=89.99, date=2024-03-10, product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001}

44. **Token 44**:
   - Fait 1: [1] TestPerson{level=1, active=false, status=inactive, department=support, age=18, score=5.5, name=Henry, salary=25000, tags=junior}
   - Fait 2: [1] TestOrder{customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10, product_id=PROD007, status=completed, priority=low, region=north}

45. **Token 45**:
   - Fait 1: [1] TestPerson{name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive, age=65, level=9}
   - Fait 2: [1] TestOrder{date=2024-03-10, product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99}

46. **Token 46**:
   - Fait 1: [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
   - Fait 2: [1] TestOrder{discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001}

47. **Token 47**:
   - Fait 1: [1] TestPerson{tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering, name=Bob, active=true}
   - Fait 2: [1] TestOrder{priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20, status=confirmed}

48. **Token 48**:
   - Fait 1: [1] TestPerson{age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8, level=3}
   - Fait 2: [1] TestOrder{discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed}

49. **Token 49**:
   - Fait 1: [1] TestPerson{name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10}

50. **Token 50**:
   - Fait 1: [1] TestPerson{active=false, tags=employee, score=8, level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10}

51. **Token 51**:
   - Fait 1: [1] TestPerson{salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering}
   - Fait 2: [1] TestOrder{region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15, amount=1, status=refunded, discount=0}

52. **Token 52**:
   - Fait 1: [1] TestPerson{name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6, tags=intern, status=inactive}
   - Fait 2: [1] TestOrder{customer_id=P006, date=2024-03-15, amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001}

53. **Token 53**:
   - Fait 1: [1] TestPerson{name=Eve, department=sales, active=false, tags=employee, score=8, level=3, age=30, salary=55000, status=inactive}
   - Fait 2: [1] TestOrder{status=pending, priority=normal, product_id=PROD001, customer_id=P001, discount=50, amount=2, date=2024-01-15, region=north, total=1999.98}

54. **Token 54**:
   - Fait 1: [1] TestPerson{salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45}
   - Fait 2: [1] TestOrder{status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20}

55. **Token 55**:
   - Fait 1: [1] TestPerson{name=Alice, status=active, department=sales, active=true, age=25, salary=45000, level=2, score=8.5, tags=junior}
   - Fait 2: [1] TestOrder{customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5}

56. **Token 56**:
   - Fait 1: [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
   - Fait 2: [1] TestOrder{total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south}

57. **Token 57**:
   - Fait 1: [1] TestPerson{salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8, level=3, age=30}
   - Fait 2: [1] TestOrder{priority=normal, region=south, amount=10, product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010}

58. **Token 58**:
   - Fait 1: [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
   - Fait 2: [1] TestOrder{discount=10, amount=1, total=89.99, date=2024-03-10, product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001}

59. **Token 59**:
   - Fait 1: [1] TestPerson{salary=25000, tags=junior, level=1, active=false, status=inactive, department=support, age=18, score=5.5, name=Henry}
   - Fait 2: [1] TestOrder{region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001, discount=50, amount=2, date=2024-01-15}

60. **Token 60**:
   - Fait 1: [1] TestPerson{tags=executive, age=65, level=9, name=Grace, salary=95000, status=active, department=management, score=10, active=true}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001, discount=50, amount=2, date=2024-01-15, region=north}

61. **Token 61**:
   - Fait 1: [1] TestPerson{score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1}
   - Fait 2: [1] TestOrder{status=pending, priority=normal, product_id=PROD001, customer_id=P001, discount=50, amount=2, date=2024-01-15, region=north, total=1999.98}

62. **Token 62**:
   - Fait 1: [1] TestPerson{name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering}
   - Fait 2: [1] TestOrder{priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005}

63. **Token 63**:
   - Fait 1: [1] TestPerson{name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive, age=65, level=9}
   - Fait 2: [1] TestOrder{customer_id=P010, priority=normal, region=south, amount=10, product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0}

64. **Token 64**:
   - Fait 1: [1] TestPerson{name=X, salary=28000, active=true, level=1, score=6.5, status=active, department=intern, age=22, tags=temp}
   - Fait 2: [1] TestOrder{amount=10, product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south}

65. **Token 65**:
   - Fait 1: [1] TestPerson{tags=junior, level=1, active=false, status=inactive, department=support, age=18, score=5.5, name=Henry, salary=25000}
   - Fait 2: [1] TestOrder{priority=normal, region=south, amount=10, product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010}

66. **Token 66**:
   - Fait 1: [1] TestPerson{salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering}
   - Fait 2: [1] TestOrder{date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south}

67. **Token 67**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa}
   - Fait 2: [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}

68. **Token 68**:
   - Fait 1: [1] TestPerson{name=Eve, department=sales, active=false, tags=employee, score=8, level=3, age=30, salary=55000, status=inactive}
   - Fait 2: [1] TestOrder{priority=normal, total=299.99, discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1}

69. **Token 69**:
   - Fait 1: [1] TestPerson{active=false, tags=employee, score=8, level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped}

70. **Token 70**:
   - Fait 1: [1] TestPerson{level=1, department=hr, score=6, tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16}
   - Fait 2: [1] TestOrder{discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20, status=confirmed, priority=low}

71. **Token 71**:
   - Fait 1: [1] TestPerson{name=X, salary=28000, active=true, level=1, score=6.5, status=active, department=intern, age=22, tags=temp}
   - Fait 2: [1] TestOrder{product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east}

72. **Token 72**:
   - Fait 1: [1] TestPerson{active=false, status=inactive, department=support, age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1}
   - Fait 2: [1] TestOrder{discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99}

73. **Token 73**:
   - Fait 1: [1] TestPerson{name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive, age=65, level=9}
   - Fait 2: [1] TestOrder{customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered}

74. **Token 74**:
   - Fait 1: [1] TestPerson{level=9, name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive, age=65}
   - Fait 2: [1] TestOrder{region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600}

75. **Token 75**:
   - Fait 1: [1] TestPerson{score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1}
   - Fait 2: [1] TestOrder{total=600, region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50}

76. **Token 76**:
   - Fait 1: [1] TestPerson{score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support, age=18}
   - Fait 2: [1] TestOrder{discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006}

77. **Token 77**:
   - Fait 1: [1] TestPerson{department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000}
   - Fait 2: [1] TestOrder{total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15, amount=1, status=refunded, discount=0, region=east}

78. **Token 78**:
   - Fait 1: [1] TestPerson{name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive, age=65, level=9}
   - Fait 2: [1] TestOrder{discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001, date=2024-02-01}

79. **Token 79**:
   - Fait 1: [1] TestPerson{active=false, status=inactive, department=support, age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1}
   - Fait 2: [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}

80. **Token 80**:
   - Fait 1: [1] TestPerson{score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1}
   - Fait 2: [1] TestOrder{amount=3, customer_id=P001, date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003}

81. **Token 81**:
   - Fait 1: [1] TestPerson{salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering}
   - Fait 2: [1] TestOrder{priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005}

82. **Token 82**:
   - Fait 1: [1] TestPerson{score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000, level=2}
   - Fait 2: [1] TestOrder{status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15, amount=1}

83. **Token 83**:
   - Fait 1: [1] TestPerson{salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45}
   - Fait 2: [1] TestOrder{date=2024-03-15, amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006}

84. **Token 84**:
   - Fait 1: [1] TestPerson{salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8, level=3, age=30}
   - Fait 2: [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}

85. **Token 85**:
   - Fait 1: [1] TestPerson{active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test, salary=-5000}
   - Fait 2: [1] TestOrder{product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east}

86. **Token 86**:
   - Fait 1: [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}
   - Fait 2: [1] TestOrder{priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005}

87. **Token 87**:
   - Fait 1: [1] TestPerson{name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped}

88. **Token 88**:
   - Fait 1: [1] TestPerson{name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering}
   - Fait 2: [1] TestOrder{product_id=PROD001, customer_id=P001, discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal}

89. **Token 89**:
   - Fait 1: [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
   - Fait 2: [1] TestOrder{customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5}

90. **Token 90**:
   - Fait 1: [1] TestPerson{tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6, age=40}
   - Fait 2: [1] TestOrder{region=north, product_id=PROD003, amount=3, customer_id=P001, date=2024-02-01, discount=15, total=225, status=shipped, priority=high}

91. **Token 91**:
   - Fait 1: [1] TestPerson{level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000}
   - Fait 2: [1] TestOrder{date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99}

92. **Token 92**:
   - Fait 1: [1] TestPerson{level=7, name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager}
   - Fait 2: [1] TestOrder{customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1}

93. **Token 93**:
   - Fait 1: [1] TestPerson{active=false, status=inactive, department=support, age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1}
   - Fait 2: [1] TestOrder{status=cancelled, total=999.98, date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0, region=west}

94. **Token 94**:
   - Fait 1: [1] TestPerson{name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive, age=65, level=9}
   - Fait 2: [1] TestOrder{total=999.98, date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled}

95. **Token 95**:
   - Fait 1: [1] TestPerson{active=true, level=1, score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000}
   - Fait 2: [1] TestOrder{total=999.98, date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled}

96. **Token 96**:
   - Fait 1: [1] TestPerson{salary=0, active=false, age=16, level=1, department=hr, score=6, tags=intern, status=inactive, name=Charlie}
   - Fait 2: [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}

97. **Token 97**:
   - Fait 1: [1] TestPerson{tags=employee, score=8, level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false}
   - Fait 2: [1] TestOrder{total=999.98, date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled}

98. **Token 98**:
   - Fait 1: [1] TestPerson{department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000}
   - Fait 2: [1] TestOrder{product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10}

99. **Token 99**:
   - Fait 1: [1] TestPerson{level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000}
   - Fait 2: [1] TestOrder{status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10, product_id=PROD007}

100. **Token 100**:
   - Fait 1: [1] TestPerson{salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45}
   - Fait 2: [1] TestOrder{product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 14: double_not_active

**Condition**: `NOT (NOT (p.active == true))`
**Types concernés**: [TestPerson]
**Terminal**: rule_14_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000}
2. [1] TestPerson{score=9.2, age=35, salary=75000, department=engineering, name=Bob, active=true, tags=senior, status=active, level=5}
3. [1] TestPerson{salary=0, active=false, age=16, level=1, department=hr, score=6, tags=intern, status=inactive, name=Charlie}
4. [1] TestPerson{salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45}
5. [1] TestPerson{active=false, tags=employee, score=8, level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales}
6. [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
7. [1] TestPerson{level=9, name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive, age=65}
8. [1] TestPerson{age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support}
9. [1] TestPerson{score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true}
10. [1] TestPerson{department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1, score=6.5, status=active}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**7 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{active=true, age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}

2. **Token 2**:
   - Fait 1: [1] TestPerson{tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering, name=Bob, active=true}

3. **Token 3**:
   - Fait 1: [1] TestPerson{score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45, salary=85000}

4. **Token 4**:
   - Fait 1: [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}

5. **Token 5**:
   - Fait 1: [1] TestPerson{active=true, tags=executive, age=65, level=9, name=Grace, salary=95000, status=active, department=management, score=10}

6. **Token 6**:
   - Fait 1: [1] TestPerson{tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6, age=40}

7. **Token 7**:
   - Fait 1: [1] TestPerson{department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1, score=6.5, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 7/10 (70.0%)
- **Efficacité**: 🟡 Élevée

---

## 🎯 RÈGLE 15: not_minor_poor_large_urgent_order

**Condition**: `p.id == o.customer_id AND NOT ((p.age < 18 OR p.salary < 25000) AND (o.total > 1000 OR o.priority == "urgent"))`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_15_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{active=true, age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}
2. [1] TestPerson{name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering}
3. [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}
4. [1] TestPerson{name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager, level=7}
5. [1] TestPerson{active=false, tags=employee, score=8, level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales}
6. [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
7. [1] TestPerson{active=true, tags=executive, age=65, level=9, name=Grace, salary=95000, status=active, department=management, score=10}
8. [1] TestPerson{age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support}
9. [1] TestPerson{age=40, tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6}
10. [1] TestPerson{score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1}
11. [1] TestOrder{discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001}
12. [1] TestOrder{discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20, status=confirmed, priority=low}
13. [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}
14. [1] TestOrder{product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east}
15. [1] TestOrder{customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1}
16. [1] TestOrder{product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005, priority=low}
17. [1] TestOrder{discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006}
18. [1] TestOrder{product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10}
19. [1] TestOrder{product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10}
20. [1] TestOrder{status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15, amount=1}

**Total**: 20 faits soumis

### 📤 RÉSULTATS TERMINAL

**100 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{active=true, age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}
   - Fait 2: [1] TestOrder{customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered}

2. **Token 2**:
   - Fait 1: [1] TestPerson{active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering, salary=68000}
   - Fait 2: [1] TestOrder{discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99}

3. **Token 3**:
   - Fait 1: [1] TestPerson{status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8, level=3, age=30, salary=55000}
   - Fait 2: [1] TestOrder{amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high}

4. **Token 4**:
   - Fait 1: [1] TestPerson{tags=junior, level=1, active=false, status=inactive, department=support, age=18, score=5.5, name=Henry, salary=25000}
   - Fait 2: [1] TestOrder{amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001, discount=50}

5. **Token 5**:
   - Fait 1: [1] TestPerson{name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6, tags=intern, status=inactive}
   - Fait 2: [1] TestOrder{amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001, discount=50}

6. **Token 6**:
   - Fait 1: [1] TestPerson{level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8}
   - Fait 2: [1] TestOrder{status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001, date=2024-02-01, discount=15, total=225}

7. **Token 7**:
   - Fait 1: [1] TestPerson{name=Eve, department=sales, active=false, tags=employee, score=8, level=3, age=30, salary=55000, status=inactive}
   - Fait 2: [1] TestOrder{product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east}

8. **Token 8**:
   - Fait 1: [1] TestPerson{tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000, level=2, score=8.5}
   - Fait 2: [1] TestOrder{discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed}

9. **Token 9**:
   - Fait 1: [1] TestPerson{active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering, salary=68000}
   - Fait 2: [1] TestOrder{product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100}

10. **Token 10**:
   - Fait 1: [1] TestPerson{active=true, tags=executive, age=65, level=9, name=Grace, salary=95000, status=active, department=management, score=10}
   - Fait 2: [1] TestOrder{discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20, status=confirmed, priority=low}

11. **Token 11**:
   - Fait 1: [1] TestPerson{level=1, department=hr, score=6, tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16}
   - Fait 2: [1] TestOrder{total=999.98, date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled}

12. **Token 12**:
   - Fait 1: [1] TestPerson{age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support}
   - Fait 2: [1] TestOrder{priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005}

13. **Token 13**:
   - Fait 1: [1] TestPerson{active=true, age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}
   - Fait 2: [1] TestOrder{customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped}

14. **Token 14**:
   - Fait 1: [1] TestPerson{status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior}
   - Fait 2: [1] TestOrder{status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01}

15. **Token 15**:
   - Fait 1: [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}
   - Fait 2: [1] TestOrder{product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10}

16. **Token 16**:
   - Fait 1: [1] TestPerson{age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support}
   - Fait 2: [1] TestOrder{product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10}

17. **Token 17**:
   - Fait 1: [1] TestPerson{age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8, level=3}
   - Fait 2: [1] TestOrder{customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10, product_id=PROD007, status=completed, priority=low, region=north}

18. **Token 18**:
   - Fait 1: [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}
   - Fait 2: [1] TestOrder{status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20}

19. **Token 19**:
   - Fait 1: [1] TestPerson{score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support, age=18}
   - Fait 2: [1] TestOrder{amount=1, region=south, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002}

20. **Token 20**:
   - Fait 1: [1] TestPerson{score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1}
   - Fait 2: [1] TestOrder{product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100}

21. **Token 21**:
   - Fait 1: [1] TestPerson{salary=95000, status=active, department=management, score=10, active=true, tags=executive, age=65, level=9, name=Grace}
   - Fait 2: [1] TestOrder{priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005}

22. **Token 22**:
   - Fait 1: [1] TestPerson{salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45}
   - Fait 2: [1] TestOrder{region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600}

23. **Token 23**:
   - Fait 1: [1] TestPerson{age=65, level=9, name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive}
   - Fait 2: [1] TestOrder{product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10}

24. **Token 24**:
   - Fait 1: [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
   - Fait 2: [1] TestOrder{product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10}

25. **Token 25**:
   - Fait 1: [1] TestPerson{name=Eve, department=sales, active=false, tags=employee, score=8, level=3, age=30, salary=55000, status=inactive}
   - Fait 2: [1] TestOrder{status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15, amount=1}

26. **Token 26**:
   - Fait 1: [1] TestPerson{score=10, active=true, tags=executive, age=65, level=9, name=Grace, salary=95000, status=active, department=management}
   - Fait 2: [1] TestOrder{discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001}

27. **Token 27**:
   - Fait 1: [1] TestPerson{name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}

28. **Token 28**:
   - Fait 1: [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
   - Fait 2: [1] TestOrder{customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered}

29. **Token 29**:
   - Fait 1: [1] TestPerson{name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager, level=7}
   - Fait 2: [1] TestOrder{amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high}

30. **Token 30**:
   - Fait 1: [1] TestPerson{salary=28000, active=true, level=1, score=6.5, status=active, department=intern, age=22, tags=temp, name=X}
   - Fait 2: [1] TestOrder{date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4}

31. **Token 31**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa}
   - Fait 2: [1] TestOrder{region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15, amount=1, status=refunded, discount=0}

32. **Token 32**:
   - Fait 1: [1] TestPerson{age=35, salary=75000, department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2}
   - Fait 2: [1] TestOrder{status=pending, priority=normal, product_id=PROD001, customer_id=P001, discount=50, amount=2, date=2024-01-15, region=north, total=1999.98}

33. **Token 33**:
   - Fait 1: [1] TestPerson{status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true}
   - Fait 2: [1] TestOrder{customer_id=P001, date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3}

34. **Token 34**:
   - Fait 1: [1] TestPerson{name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager, level=7}
   - Fait 2: [1] TestOrder{product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east}

35. **Token 35**:
   - Fait 1: [1] TestPerson{name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001}

36. **Token 36**:
   - Fait 1: [1] TestPerson{score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1}
   - Fait 2: [1] TestOrder{product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10}

37. **Token 37**:
   - Fait 1: [1] TestPerson{department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1, score=6.5, status=active}
   - Fait 2: [1] TestOrder{discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15, amount=1, status=refunded}

38. **Token 38**:
   - Fait 1: [1] TestPerson{score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1}
   - Fait 2: [1] TestOrder{total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001, date=2024-02-01, discount=15}

39. **Token 39**:
   - Fait 1: [1] TestPerson{level=1, score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true}
   - Fait 2: [1] TestOrder{total=299.99, discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal}

40. **Token 40**:
   - Fait 1: [1] TestPerson{tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering, name=Bob, active=true}
   - Fait 2: [1] TestOrder{priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005}

41. **Token 41**:
   - Fait 1: [1] TestPerson{active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test, salary=-5000}
   - Fait 2: [1] TestOrder{total=600, region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50}

42. **Token 42**:
   - Fait 1: [1] TestPerson{tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering, name=Bob, active=true}
   - Fait 2: [1] TestOrder{region=south, amount=10, product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal}

43. **Token 43**:
   - Fait 1: [1] TestPerson{salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45}
   - Fait 2: [1] TestOrder{region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10, product_id=PROD007, status=completed, priority=low}

44. **Token 44**:
   - Fait 1: [1] TestPerson{salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45}
   - Fait 2: [1] TestOrder{customer_id=P006, date=2024-03-15, amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001}

45. **Token 45**:
   - Fait 1: [1] TestPerson{active=true, age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}
   - Fait 2: [1] TestOrder{total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001, discount=50, amount=2, date=2024-01-15, region=north}

46. **Token 46**:
   - Fait 1: [1] TestPerson{age=40, tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6}
   - Fait 2: [1] TestOrder{discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001}

47. **Token 47**:
   - Fait 1: [1] TestPerson{name=X, salary=28000, active=true, level=1, score=6.5, status=active, department=intern, age=22, tags=temp}
   - Fait 2: [1] TestOrder{date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south}

48. **Token 48**:
   - Fait 1: [1] TestPerson{status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6, tags=intern}
   - Fait 2: [1] TestOrder{product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100}

49. **Token 49**:
   - Fait 1: [1] TestPerson{age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support}
   - Fait 2: [1] TestOrder{status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002}

50. **Token 50**:
   - Fait 1: [1] TestPerson{age=45, salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active}
   - Fait 2: [1] TestOrder{priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20, status=confirmed}

51. **Token 51**:
   - Fait 1: [1] TestPerson{score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000, level=2}
   - Fait 2: [1] TestOrder{total=999.98, date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled}

52. **Token 52**:
   - Fait 1: [1] TestPerson{salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering}
   - Fait 2: [1] TestOrder{date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98}

53. **Token 53**:
   - Fait 1: [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}
   - Fait 2: [1] TestOrder{region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600}

54. **Token 54**:
   - Fait 1: [1] TestPerson{age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support}
   - Fait 2: [1] TestOrder{region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600}

55. **Token 55**:
   - Fait 1: [1] TestPerson{level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000}
   - Fait 2: [1] TestOrder{region=south, amount=10, product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal}

56. **Token 56**:
   - Fait 1: [1] TestPerson{tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6, age=40}
   - Fait 2: [1] TestOrder{total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10, product_id=PROD002}

57. **Token 57**:
   - Fait 1: [1] TestPerson{name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south}

58. **Token 58**:
   - Fait 1: [1] TestPerson{active=false, tags=employee, score=8, level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{total=999.98, date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled}

59. **Token 59**:
   - Fait 1: [1] TestPerson{name=Eve, department=sales, active=false, tags=employee, score=8, level=3, age=30, salary=55000, status=inactive}
   - Fait 2: [1] TestOrder{product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10}

60. **Token 60**:
   - Fait 1: [1] TestPerson{salary=25000, tags=junior, level=1, active=false, status=inactive, department=support, age=18, score=5.5, name=Henry}
   - Fait 2: [1] TestOrder{status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10, product_id=PROD007}

61. **Token 61**:
   - Fait 1: [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}
   - Fait 2: [1] TestOrder{amount=1, total=89.99, date=2024-03-10, product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10}

62. **Token 62**:
   - Fait 1: [1] TestPerson{level=1, department=hr, score=6, tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16}
   - Fait 2: [1] TestOrder{amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15}

63. **Token 63**:
   - Fait 1: [1] TestPerson{age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support}
   - Fait 2: [1] TestOrder{customer_id=P006, date=2024-03-15, amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001}

64. **Token 64**:
   - Fait 1: [1] TestPerson{status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8, level=3, age=30, salary=55000}
   - Fait 2: [1] TestOrder{region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001, discount=50, amount=2, date=2024-01-15}

65. **Token 65**:
   - Fait 1: [1] TestPerson{level=1, department=hr, score=6, tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16}
   - Fait 2: [1] TestOrder{product_id=PROD003, amount=3, customer_id=P001, date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north}

66. **Token 66**:
   - Fait 1: [1] TestPerson{score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support, age=18}
   - Fait 2: [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}

67. **Token 67**:
   - Fait 1: [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}
   - Fait 2: [1] TestOrder{date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east, product_id=PROD004}

68. **Token 68**:
   - Fait 1: [1] TestPerson{age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support}
   - Fait 2: [1] TestOrder{discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99}

69. **Token 69**:
   - Fait 1: [1] TestPerson{score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test, salary=-5000, active=true}
   - Fait 2: [1] TestOrder{product_id=PROD001, customer_id=P001, discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal}

70. **Token 70**:
   - Fait 1: [1] TestPerson{active=true, tags=executive, age=65, level=9, name=Grace, salary=95000, status=active, department=management, score=10}
   - Fait 2: [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}

71. **Token 71**:
   - Fait 1: [1] TestPerson{status=active, department=management, score=10, active=true, tags=executive, age=65, level=9, name=Grace, salary=95000}
   - Fait 2: [1] TestOrder{date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east, product_id=PROD004}

72. **Token 72**:
   - Fait 1: [1] TestPerson{salary=75000, department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35}
   - Fait 2: [1] TestOrder{region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001}

73. **Token 73**:
   - Fait 1: [1] TestPerson{active=false, tags=employee, score=8, level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales}
   - Fait 2: [1] TestOrder{product_id=PROD002, amount=1, region=south, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002}

74. **Token 74**:
   - Fait 1: [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
   - Fait 2: [1] TestOrder{amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005}

75. **Token 75**:
   - Fait 1: [1] TestPerson{tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering, name=Bob, active=true}
   - Fait 2: [1] TestOrder{date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4}

76. **Token 76**:
   - Fait 1: [1] TestPerson{name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1}
   - Fait 2: [1] TestOrder{product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10}

77. **Token 77**:
   - Fait 1: [1] TestPerson{tags=executive, age=65, level=9, name=Grace, salary=95000, status=active, department=management, score=10, active=true}
   - Fait 2: [1] TestOrder{product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10}

78. **Token 78**:
   - Fait 1: [1] TestPerson{active=true, tags=executive, age=65, level=9, name=Grace, salary=95000, status=active, department=management, score=10}
   - Fait 2: [1] TestOrder{discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15, amount=1, status=refunded}

79. **Token 79**:
   - Fait 1: [1] TestPerson{level=1, score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true}
   - Fait 2: [1] TestOrder{region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001, discount=50, amount=2, date=2024-01-15}

80. **Token 80**:
   - Fait 1: [1] TestPerson{salary=45000, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25}
   - Fait 2: [1] TestOrder{date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south}

81. **Token 81**:
   - Fait 1: [1] TestPerson{level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7}
   - Fait 2: [1] TestOrder{region=south, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1}

82. **Token 82**:
   - Fait 1: [1] TestPerson{salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45}
   - Fait 2: [1] TestOrder{priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005}

83. **Token 83**:
   - Fait 1: [1] TestPerson{name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive, age=65, level=9}
   - Fait 2: [1] TestOrder{region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600}

84. **Token 84**:
   - Fait 1: [1] TestPerson{department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45, salary=85000, score=7.8}
   - Fait 2: [1] TestOrder{priority=normal, region=south, amount=10, product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010}

85. **Token 85**:
   - Fait 1: [1] TestPerson{tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering, name=Bob, active=true}
   - Fait 2: [1] TestOrder{product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10}

86. **Token 86**:
   - Fait 1: [1] TestPerson{salary=75000, department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35}
   - Fait 2: [1] TestOrder{date=2024-03-15, amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006}

87. **Token 87**:
   - Fait 1: [1] TestPerson{salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45}
   - Fait 2: [1] TestOrder{discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001}

88. **Token 88**:
   - Fait 1: [1] TestPerson{active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000, department=engineering, name=Bob}
   - Fait 2: [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}

89. **Token 89**:
   - Fait 1: [1] TestPerson{age=35, salary=75000, department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2}
   - Fait 2: [1] TestOrder{amount=1, priority=normal, total=299.99, discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004}

90. **Token 90**:
   - Fait 1: [1] TestPerson{department=management, score=10, active=true, tags=executive, age=65, level=9, name=Grace, salary=95000, status=active}
   - Fait 2: [1] TestOrder{product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100}

91. **Token 91**:
   - Fait 1: [1] TestPerson{tags=temp, name=X, salary=28000, active=true, level=1, score=6.5, status=active, department=intern, age=22}
   - Fait 2: [1] TestOrder{priority=low, product_id=PROD005, amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005}

92. **Token 92**:
   - Fait 1: [1] TestPerson{score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1}
   - Fait 2: [1] TestOrder{product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10}

93. **Token 93**:
   - Fait 1: [1] TestPerson{department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000}
   - Fait 2: [1] TestOrder{priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1, region=south, date=2024-01-20, status=confirmed}

94. **Token 94**:
   - Fait 1: [1] TestPerson{level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8}
   - Fait 2: [1] TestOrder{customer_id=P007, product_id=PROD006, discount=50, total=600, region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped}

95. **Token 95**:
   - Fait 1: [1] TestPerson{level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000}
   - Fait 2: [1] TestOrder{product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10}

96. **Token 96**:
   - Fait 1: [1] TestPerson{tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6, age=40}
   - Fait 2: [1] TestOrder{status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10, product_id=PROD007}

97. **Token 97**:
   - Fait 1: [1] TestPerson{active=true, age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}
   - Fait 2: [1] TestOrder{customer_id=P006, date=2024-03-15, amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001}

98. **Token 98**:
   - Fait 1: [1] TestPerson{status=active, name=Ivy, department=engineering, salary=68000, active=true, score=8.7, level=6, age=40, tags=senior}
   - Fait 2: [1] TestOrder{discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15, amount=1, status=refunded}

99. **Token 99**:
   - Fait 1: [1] TestPerson{level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000}
   - Fait 2: [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}

100. **Token 100**:
   - Fait 1: [1] TestPerson{salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering}
   - Fait 2: [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 100/20 (500.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 16: valid_non_zero_person

**Condition**: `p.age != 0 AND p.salary > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_16_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{active=true, age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales}
2. [1] TestPerson{score=9.2, age=35, salary=75000, department=engineering, name=Bob, active=true, tags=senior, status=active, level=5}
3. [1] TestPerson{age=16, level=1, department=hr, score=6, tags=intern, status=inactive, name=Charlie, salary=0, active=false}
4. [1] TestPerson{salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45}
5. [1] TestPerson{department=sales, active=false, tags=employee, score=8, level=3, age=30, salary=55000, status=inactive, name=Eve}
6. [1] TestPerson{age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank}
7. [1] TestPerson{age=65, level=9, name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive}
8. [1] TestPerson{active=false, status=inactive, department=support, age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1}
9. [1] TestPerson{score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true}
10. [1] TestPerson{score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true}

2. **Token 2**:
   - Fait 1: [1] TestPerson{department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000}

3. **Token 3**:
   - Fait 1: [1] TestPerson{name=Grace, salary=95000, status=active, department=management, score=10, active=true, tags=executive, age=65, level=9}

4. **Token 4**:
   - Fait 1: [1] TestPerson{status=inactive, department=support, age=18, score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false}

5. **Token 5**:
   - Fait 1: [1] TestPerson{salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering}

6. **Token 6**:
   - Fait 1: [1] TestPerson{score=6.5, status=active, department=intern, age=22, tags=temp, name=X, salary=28000, active=true, level=1}

7. **Token 7**:
   - Fait 1: [1] TestPerson{level=1, department=hr, score=6, tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16}

8. **Token 8**:
   - Fait 1: [1] TestPerson{name=Diana, active=true, status=active, age=45, salary=85000, score=7.8, department=marketing, tags=manager, level=7}

9. **Token 9**:
   - Fait 1: [1] TestPerson{tags=employee, score=8, level=3, age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 17: valid_positive_order

**Condition**: `o.amount > 0 AND o.total > 0`
**Types concernés**: [TestOrder]
**Terminal**: rule_17_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{discount=50, amount=2, date=2024-01-15, region=north, total=1999.98, status=pending, priority=normal, product_id=PROD001, customer_id=P001}
2. [1] TestOrder{amount=1, region=south, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002}
3. [1] TestOrder{date=2024-02-01, discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001}
4. [1] TestOrder{product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99, discount=0, region=east}
5. [1] TestOrder{date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100, product_id=PROD001, region=south, total=999.99}
6. [1] TestOrder{amount=2, discount=0, region=west, status=cancelled, total=999.98, date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005}
7. [1] TestOrder{region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50, total=600}
8. [1] TestOrder{total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal, region=south, amount=10, product_id=PROD002}
9. [1] TestOrder{product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10, amount=1, total=89.99, date=2024-03-10}
10. [1] TestOrder{amount=1, status=refunded, discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{total=600, region=north, priority=urgent, amount=4, date=2024-03-01, status=shipped, customer_id=P007, product_id=PROD006, discount=50}

2. **Token 2**:
   - Fait 1: [1] TestOrder{discount=0, region=east, total=75000, priority=urgent, product_id=PROD001, customer_id=P006, date=2024-03-15, amount=1, status=refunded}

3. **Token 3**:
   - Fait 1: [1] TestOrder{discount=15, total=225, status=shipped, priority=high, region=north, product_id=PROD003, amount=3, customer_id=P001, date=2024-02-01}

4. **Token 4**:
   - Fait 1: [1] TestOrder{status=cancelled, total=999.98, date=2024-02-15, customer_id=P005, priority=low, product_id=PROD005, amount=2, discount=0, region=west}

5. **Token 5**:
   - Fait 1: [1] TestOrder{region=south, amount=10, product_id=PROD002, total=255, date=2024-03-05, status=pending, discount=0, customer_id=P010, priority=normal}

6. **Token 6**:
   - Fait 1: [1] TestOrder{amount=1, total=89.99, date=2024-03-10, product_id=PROD007, status=completed, priority=low, region=north, customer_id=P001, discount=10}

7. **Token 7**:
   - Fait 1: [1] TestOrder{status=pending, priority=normal, product_id=PROD001, customer_id=P001, discount=50, amount=2, date=2024-01-15, region=north, total=1999.98}

8. **Token 8**:
   - Fait 1: [1] TestOrder{region=south, date=2024-01-20, status=confirmed, priority=low, discount=0, total=25.5, customer_id=P002, product_id=PROD002, amount=1}

9. **Token 9**:
   - Fait 1: [1] TestOrder{discount=0, region=east, product_id=PROD004, date=2024-02-05, status=delivered, customer_id=P004, amount=1, priority=normal, total=299.99}

10. **Token 10**:
   - Fait 1: [1] TestOrder{product_id=PROD001, region=south, total=999.99, date=2024-02-10, priority=high, amount=1, customer_id=P002, status=confirmed, discount=100}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 18: valid_person_name

**Condition**: `LENGTH(p.name) > 0`
**Types concernés**: [TestPerson]
**Terminal**: rule_18_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{level=2, score=8.5, tags=junior, name=Alice, status=active, department=sales, active=true, age=25, salary=45000}
2. [1] TestPerson{department=engineering, name=Bob, active=true, tags=senior, status=active, level=5, score=9.2, age=35, salary=75000}
3. [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}
4. [1] TestPerson{salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45}
5. [1] TestPerson{age=30, salary=55000, status=inactive, name=Eve, department=sales, active=false, tags=employee, score=8, level=3}
6. [1] TestPerson{salary=-5000, active=true, score=0, status=active, department=qa, level=1, name=Frank, age=0, tags=test}
7. [1] TestPerson{status=active, department=management, score=10, active=true, tags=executive, age=65, level=9, name=Grace, salary=95000}
8. [1] TestPerson{score=5.5, name=Henry, salary=25000, tags=junior, level=1, active=false, status=inactive, department=support, age=18}
9. [1] TestPerson{score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering, salary=68000, active=true}
10. [1] TestPerson{age=22, tags=temp, name=X, salary=28000, active=true, level=1, score=6.5, status=active, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**10 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{status=active, department=sales, active=true, age=25, salary=45000, level=2, score=8.5, tags=junior, name=Alice}

2. **Token 2**:
   - Fait 1: [1] TestPerson{tags=intern, status=inactive, name=Charlie, salary=0, active=false, age=16, level=1, department=hr, score=6}

3. **Token 3**:
   - Fait 1: [1] TestPerson{department=sales, active=false, tags=employee, score=8, level=3, age=30, salary=55000, status=inactive, name=Eve}

4. **Token 4**:
   - Fait 1: [1] TestPerson{score=10, active=true, tags=executive, age=65, level=9, name=Grace, salary=95000, status=active, department=management}

5. **Token 5**:
   - Fait 1: [1] TestPerson{salary=68000, active=true, score=8.7, level=6, age=40, tags=senior, status=active, name=Ivy, department=engineering}

6. **Token 6**:
   - Fait 1: [1] TestPerson{score=9.2, age=35, salary=75000, department=engineering, name=Bob, active=true, tags=senior, status=active, level=5}

7. **Token 7**:
   - Fait 1: [1] TestPerson{salary=85000, score=7.8, department=marketing, tags=manager, level=7, name=Diana, active=true, status=active, age=45}

8. **Token 8**:
   - Fait 1: [1] TestPerson{level=1, name=Frank, age=0, tags=test, salary=-5000, active=true, score=0, status=active, department=qa}

9. **Token 9**:
   - Fait 1: [1] TestPerson{level=1, active=false, status=inactive, department=support, age=18, score=5.5, name=Henry, salary=25000, tags=junior}

10. **Token 10**:
   - Fait 1: [1] TestPerson{name=X, salary=28000, active=true, level=1, score=6.5, status=active, department=intern, age=22, tags=temp}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 10/10 (100.0%)
- **Efficacité**: 🟢 Très élevée

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 19 (100.0%)
- **Tokens générés**: 330
- **Faits traités**: 27
