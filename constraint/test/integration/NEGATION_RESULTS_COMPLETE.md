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

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, age=25, level=2}
2. [1] TestPerson{id=P002, score=9.2, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, name=Bob}
3. [1] TestPerson{id=P003, level=1, name=Charlie, salary=0, active=false, tags=intern, department=hr, age=16, score=6, status=inactive}
4. [1] TestPerson{id=P004, salary=85000, tags=manager, status=active, department=marketing, age=45, active=true, score=7.8, level=7, name=Diana}
5. [1] TestPerson{id=P005, age=30, salary=55000, active=false, score=8, status=inactive, department=sales, name=Eve, tags=employee, level=3}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, score=0, status=active, department=qa, age=0, active=true, tags=test, level=1}
7. [1] TestPerson{id=P007, active=true, score=10, tags=executive, status=active, department=management, name=Grace, age=65, salary=95000, level=9}
8. [1] TestPerson{id=P008, age=18, score=5.5, tags=junior, status=inactive, level=1, name=Henry, salary=25000, active=false, department=support}
9. [1] TestPerson{id=P009, active=true, score=8.7, tags=senior, status=active, name=Ivy, salary=68000, department=engineering, level=6, age=40}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, level=1, name=X, score=6.5, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, tags=intern, department=hr, age=16, score=6, status=inactive, level=1}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P007, salary=95000, level=9, active=true, score=10, tags=executive, status=active, department=management, name=Grace, age=65}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, age=25, level=2}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, name=Bob, score=9.2}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P004, age=45, active=true, score=7.8, level=7, name=Diana, salary=85000, tags=manager, status=active, department=marketing}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P005, name=Eve, tags=employee, level=3, age=30, salary=55000, active=false, score=8, status=inactive, department=sales}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, department=support, age=18, score=5.5, tags=junior, status=inactive, level=1}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, active=true, score=8.7, tags=senior, status=active, name=Ivy, salary=68000, department=engineering, level=6, age=40}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, level=1, name=X, score=6.5, department=intern, age=22, salary=28000, active=true, tags=temp, status=active}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, product_id=PROD001, amount=2, total=1999.98, region=north, date=2024-01-15, status=pending, priority=normal, discount=50, customer_id=P001}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, amount=1, total=25.5, discount=0, region=south}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, discount=15, product_id=PROD003, total=225, status=shipped, priority=high, region=north}
4. [1] TestOrder{id=O004, total=299.99, status=delivered, priority=normal, amount=1, date=2024-02-05, discount=0, region=east, customer_id=P004, product_id=PROD004}
5. [1] TestOrder{id=O005, total=999.99, status=confirmed, amount=1, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001}
6. [1] TestOrder{id=O006, status=cancelled, priority=low, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005, product_id=PROD005, total=999.98}
7. [1] TestOrder{id=O007, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, total=600, discount=50, region=north, date=2024-03-01}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, region=south, amount=10, total=255, priority=normal, discount=0}
9. [1] TestOrder{id=O009, product_id=PROD007, total=89.99, date=2024-03-10, discount=10, region=north, amount=1, status=completed, priority=low, customer_id=P001}
10. [1] TestOrder{id=O010, date=2024-03-15, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, amount=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestOrder{id=O008, amount=10, total=255, priority=normal, discount=0, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, region=south}

2. **Token 2**:
   - Fait 1: [1] TestOrder{id=O010, amount=1, date=2024-03-15, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006}

3. **Token 3**:
   - Fait 1: [1] TestOrder{id=O001, status=pending, priority=normal, discount=50, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, region=north, date=2024-01-15}

4. **Token 4**:
   - Fait 1: [1] TestOrder{id=O002, total=25.5, discount=0, region=south, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, amount=1}

5. **Token 5**:
   - Fait 1: [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, discount=15, product_id=PROD003, total=225, status=shipped, priority=high, region=north}

6. **Token 6**:
   - Fait 1: [1] TestOrder{id=O004, priority=normal, amount=1, date=2024-02-05, discount=0, region=east, customer_id=P004, product_id=PROD004, total=299.99, status=delivered}

7. **Token 7**:
   - Fait 1: [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, discount=50, region=north, date=2024-03-01, status=shipped, priority=urgent}

8. **Token 8**:
   - Fait 1: [1] TestOrder{id=O009, total=89.99, date=2024-03-10, discount=10, region=north, amount=1, status=completed, priority=low, customer_id=P001, product_id=PROD007}

9. **Token 9**:
   - Fait 1: [1] TestOrder{id=O005, amount=1, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, status=confirmed}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, age=25, level=2}
2. [1] TestPerson{id=P002, score=9.2, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, name=Bob}
3. [1] TestPerson{id=P003, level=1, name=Charlie, salary=0, active=false, tags=intern, department=hr, age=16, score=6, status=inactive}
4. [1] TestPerson{id=P004, age=45, active=true, score=7.8, level=7, name=Diana, salary=85000, tags=manager, status=active, department=marketing}
5. [1] TestPerson{id=P005, name=Eve, tags=employee, level=3, age=30, salary=55000, active=false, score=8, status=inactive, department=sales}
6. [1] TestPerson{id=P006, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active, department=qa}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, level=9, active=true, score=10, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, department=support, age=18, score=5.5, tags=junior, status=inactive, level=1}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, department=engineering, level=6, age=40, active=true, score=8.7, tags=senior, status=active}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, level=1, name=X, score=6.5, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 3: not_high_total

**Condition**: `NOT (o.total > 50000)`
**Types concernés**: [TestOrder]
**Terminal**: rule_3_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, region=north, date=2024-01-15, status=pending, priority=normal, discount=50}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, amount=1, total=25.5, discount=0, region=south}
3. [1] TestOrder{id=O003, product_id=PROD003, total=225, status=shipped, priority=high, region=north, customer_id=P001, amount=3, date=2024-02-01, discount=15}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal, amount=1, date=2024-02-05, discount=0, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, status=confirmed, amount=1, date=2024-02-10, priority=high, discount=100, region=south}
6. [1] TestOrder{id=O006, amount=2, date=2024-02-15, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0, region=west}
7. [1] TestOrder{id=O007, amount=4, total=600, discount=50, region=north, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, region=south, amount=10, total=255, priority=normal, discount=0}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, discount=10, region=north, amount=1, status=completed, priority=low}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, date=2024-03-15, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0, region=east}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 4: not_cheap_product

**Condition**: `NOT (prod.price <= 10)`
**Types concernés**: [TestProduct]
**Terminal**: rule_4_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, supplier=TechSupply, available=true, rating=4.5, keywords=computer, name=Laptop, category=electronics, price=999.99, brand=TechCorp, stock=50}
2. [1] TestProduct{id=PROD002, name=Mouse, price=25.5, rating=4, keywords=peripheral, supplier=TechSupply, category=accessories, available=true, brand=TechCorp, stock=200}
3. [1] TestProduct{id=PROD003, rating=3.5, stock=0, category=accessories, available=false, keywords=typing, brand=KeyTech, supplier=KeySupply, name=Keyboard, price=75}
4. [1] TestProduct{id=PROD004, category=electronics, available=true, stock=30, supplier=ScreenSupply, name=Monitor, price=299.99, rating=4.8, keywords=display, brand=ScreenPro}
5. [1] TestProduct{id=PROD005, name=OldKeyboard, category=accessories, price=8.5, rating=2, stock=0, available=false, keywords=obsolete, brand=OldTech, supplier=OldSupply}
6. [1] TestProduct{id=PROD006, price=150, available=true, rating=4.6, supplier=AudioSupply, name=Headphones, category=audio, keywords=sound, brand=AudioMax, stock=75}
7. [1] TestProduct{id=PROD007, name=Webcam, category=electronics, available=true, stock=25, supplier=CamSupply, price=89.99, rating=3.8, keywords=video, brand=CamTech}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/7 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 5: not_age_times_thousand_less_salary

**Condition**: `NOT (p.age * 1000 < p.salary)`
**Types concernés**: [TestPerson]
**Terminal**: rule_5_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, age=25, level=2, name=Alice, salary=45000, active=true, score=8.5, tags=junior, status=active}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, name=Bob, score=9.2}
3. [1] TestPerson{id=P003, tags=intern, department=hr, age=16, score=6, status=inactive, level=1, name=Charlie, salary=0, active=false}
4. [1] TestPerson{id=P004, age=45, active=true, score=7.8, level=7, name=Diana, salary=85000, tags=manager, status=active, department=marketing}
5. [1] TestPerson{id=P005, tags=employee, level=3, age=30, salary=55000, active=false, score=8, status=inactive, department=sales, name=Eve}
6. [1] TestPerson{id=P006, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active, department=qa, age=0}
7. [1] TestPerson{id=P007, score=10, tags=executive, status=active, department=management, name=Grace, age=65, salary=95000, level=9, active=true}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, department=support, age=18, score=5.5, tags=junior, status=inactive, level=1}
9. [1] TestPerson{id=P009, age=40, active=true, score=8.7, tags=senior, status=active, name=Ivy, salary=68000, department=engineering, level=6}
10. [1] TestPerson{id=P010, department=intern, age=22, salary=28000, active=true, tags=temp, status=active, level=1, name=X, score=6.5}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 6: not_amount_plus_discount_geq_total

**Condition**: `NOT (o.amount + o.discount >= o.total)`
**Types concernés**: [TestOrder]
**Terminal**: rule_6_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, region=north, date=2024-01-15, status=pending, priority=normal, discount=50}
2. [1] TestOrder{id=O002, amount=1, total=25.5, discount=0, region=south, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, discount=15, product_id=PROD003, total=225, status=shipped, priority=high, region=north}
4. [1] TestOrder{id=O004, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal, amount=1, date=2024-02-05, discount=0, region=east}
5. [1] TestOrder{id=O005, total=999.99, status=confirmed, amount=1, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0, region=west, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, total=600, discount=50, region=north, date=2024-03-01}
8. [1] TestOrder{id=O008, amount=10, total=255, priority=normal, discount=0, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, discount=10, region=north, amount=1, status=completed, priority=low}
10. [1] TestOrder{id=O010, region=east, customer_id=P006, amount=1, date=2024-03-15, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 7: not_active_high_earner

**Condition**: `NOT (p.active == true AND p.salary > 70000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_7_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, department=sales, age=25, level=2, name=Alice, salary=45000, active=true, score=8.5, tags=junior, status=active}
2. [1] TestPerson{id=P002, active=true, tags=senior, status=active, department=engineering, level=5, name=Bob, score=9.2, age=35, salary=75000}
3. [1] TestPerson{id=P003, active=false, tags=intern, department=hr, age=16, score=6, status=inactive, level=1, name=Charlie, salary=0}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, tags=manager, status=active, department=marketing, age=45, active=true, score=7.8, level=7}
5. [1] TestPerson{id=P005, age=30, salary=55000, active=false, score=8, status=inactive, department=sales, name=Eve, tags=employee, level=3}
6. [1] TestPerson{id=P006, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active, department=qa}
7. [1] TestPerson{id=P007, department=management, name=Grace, age=65, salary=95000, level=9, active=true, score=10, tags=executive, status=active}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, department=support, age=18, score=5.5, tags=junior, status=inactive, level=1}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, department=engineering, level=6, age=40, active=true, score=8.7, tags=senior, status=active}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, level=1, name=X, score=6.5, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 8: not_pending_or_low_priority

**Condition**: `NOT (o.status == "pending" OR o.priority == "low")`
**Types concernés**: [TestOrder]
**Terminal**: rule_8_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, date=2024-01-15, status=pending, priority=normal, discount=50, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, region=north}
2. [1] TestOrder{id=O002, status=confirmed, priority=low, amount=1, total=25.5, discount=0, region=south, customer_id=P002, product_id=PROD002, date=2024-01-20}
3. [1] TestOrder{id=O003, customer_id=P001, amount=3, date=2024-02-01, discount=15, product_id=PROD003, total=225, status=shipped, priority=high, region=north}
4. [1] TestOrder{id=O004, amount=1, date=2024-02-05, discount=0, region=east, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, total=999.99, status=confirmed, amount=1, date=2024-02-10, priority=high, discount=100, region=south}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0, region=west, amount=2, date=2024-02-15}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, discount=50, region=north, date=2024-03-01, status=shipped, priority=urgent}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, region=south, amount=10, total=255, priority=normal, discount=0}
9. [1] TestOrder{id=O009, discount=10, region=north, amount=1, status=completed, priority=low, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10}
10. [1] TestOrder{id=O010, region=east, customer_id=P006, amount=1, date=2024-03-15, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 9: not_short_name

**Condition**: `NOT (LENGTH(p.name) < 3)`
**Types concernés**: [TestPerson]
**Terminal**: rule_9_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, age=25, level=2}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, name=Bob, score=9.2}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, tags=intern, department=hr, age=16, score=6, status=inactive, level=1}
4. [1] TestPerson{id=P004, department=marketing, age=45, active=true, score=7.8, level=7, name=Diana, salary=85000, tags=manager, status=active}
5. [1] TestPerson{id=P005, age=30, salary=55000, active=false, score=8, status=inactive, department=sales, name=Eve, tags=employee, level=3}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, score=0, status=active, department=qa, age=0, active=true, tags=test, level=1}
7. [1] TestPerson{id=P007, active=true, score=10, tags=executive, status=active, department=management, name=Grace, age=65, salary=95000, level=9}
8. [1] TestPerson{id=P008, status=inactive, level=1, name=Henry, salary=25000, active=false, department=support, age=18, score=5.5, tags=junior}
9. [1] TestPerson{id=P009, active=true, score=8.7, tags=senior, status=active, name=Ivy, salary=68000, department=engineering, level=6, age=40}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, level=1, name=X, score=6.5, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 10: not_obsolete_product

**Condition**: `NOT (prod.keywords CONTAINS "obsolete")`
**Types concernés**: [TestProduct]
**Terminal**: rule_10_terminal

### 📥 FAITS SOUMIS

1. [1] TestProduct{id=PROD001, stock=50, supplier=TechSupply, available=true, rating=4.5, keywords=computer, name=Laptop, category=electronics, price=999.99, brand=TechCorp}
2. [1] TestProduct{id=PROD002, category=accessories, available=true, brand=TechCorp, stock=200, name=Mouse, price=25.5, rating=4, keywords=peripheral, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, stock=0, category=accessories, available=false, keywords=typing, brand=KeyTech, supplier=KeySupply, name=Keyboard, price=75, rating=3.5}
4. [1] TestProduct{id=PROD004, available=true, stock=30, supplier=ScreenSupply, name=Monitor, price=299.99, rating=4.8, keywords=display, brand=ScreenPro, category=electronics}
5. [1] TestProduct{id=PROD005, brand=OldTech, supplier=OldSupply, name=OldKeyboard, category=accessories, price=8.5, rating=2, stock=0, available=false, keywords=obsolete}
6. [1] TestProduct{id=PROD006, stock=75, price=150, available=true, rating=4.6, supplier=AudioSupply, name=Headphones, category=audio, keywords=sound, brand=AudioMax}
7. [1] TestProduct{id=PROD007, brand=CamTech, name=Webcam, category=electronics, available=true, stock=25, supplier=CamSupply, price=89.99, rating=3.8, keywords=video}

**Total**: 7 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/7 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 11: not_temporary_employee

**Condition**: `NOT (p.department IN ["temp", "intern"])`
**Types concernés**: [TestPerson]
**Terminal**: rule_11_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, age=25, level=2}
2. [1] TestPerson{id=P002, status=active, department=engineering, level=5, name=Bob, score=9.2, age=35, salary=75000, active=true, tags=senior}
3. [1] TestPerson{id=P003, age=16, score=6, status=inactive, level=1, name=Charlie, salary=0, active=false, tags=intern, department=hr}
4. [1] TestPerson{id=P004, salary=85000, tags=manager, status=active, department=marketing, age=45, active=true, score=7.8, level=7, name=Diana}
5. [1] TestPerson{id=P005, status=inactive, department=sales, name=Eve, tags=employee, level=3, age=30, salary=55000, active=false, score=8}
6. [1] TestPerson{id=P006, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active, department=qa}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, level=9, active=true, score=10, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, department=support, age=18, score=5.5, tags=junior, status=inactive, level=1}
9. [1] TestPerson{id=P009, score=8.7, tags=senior, status=active, name=Ivy, salary=68000, department=engineering, level=6, age=40, active=true}
10. [1] TestPerson{id=P010, active=true, tags=temp, status=active, level=1, name=X, score=6.5, department=intern, age=22, salary=28000}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 12: not_cancelled_refunded_order

**Condition**: `NOT (o.status IN ["cancelled", "refunded"])`
**Types concernés**: [TestOrder]
**Terminal**: rule_12_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, date=2024-01-15, status=pending, priority=normal, discount=50, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, region=north}
2. [1] TestOrder{id=O002, discount=0, region=south, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, amount=1, total=25.5}
3. [1] TestOrder{id=O003, product_id=PROD003, total=225, status=shipped, priority=high, region=north, customer_id=P001, amount=3, date=2024-02-01, discount=15}
4. [1] TestOrder{id=O004, region=east, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal, amount=1, date=2024-02-05, discount=0}
5. [1] TestOrder{id=O005, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, status=confirmed, amount=1}
6. [1] TestOrder{id=O006, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0, region=west, amount=2, date=2024-02-15, customer_id=P005}
7. [1] TestOrder{id=O007, customer_id=P007, product_id=PROD006, amount=4, total=600, discount=50, region=north, date=2024-03-01, status=shipped, priority=urgent}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, region=south, amount=10, total=255, priority=normal, discount=0}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, discount=10, region=north, amount=1, status=completed, priority=low}
10. [1] TestOrder{id=O010, discount=0, region=east, customer_id=P006, amount=1, date=2024-03-15, product_id=PROD001, total=75000, status=refunded, priority=urgent}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 13: not_order_exceeds_monthly_salary

**Condition**: `p.id == o.customer_id AND NOT (o.total > p.salary / 12)`
**Types concernés**: [TestPerson TestOrder]
**Terminal**: rule_13_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, age=25, level=2, name=Alice}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, name=Bob, score=9.2}
3. [1] TestPerson{id=P003, age=16, score=6, status=inactive, level=1, name=Charlie, salary=0, active=false, tags=intern, department=hr}
4. [1] TestPerson{id=P004, level=7, name=Diana, salary=85000, tags=manager, status=active, department=marketing, age=45, active=true, score=7.8}
5. [1] TestPerson{id=P005, name=Eve, tags=employee, level=3, age=30, salary=55000, active=false, score=8, status=inactive, department=sales}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, department=qa, age=0, active=true, tags=test, level=1, name=Frank}
7. [1] TestPerson{id=P007, active=true, score=10, tags=executive, status=active, department=management, name=Grace, age=65, salary=95000, level=9}
8. [1] TestPerson{id=P008, tags=junior, status=inactive, level=1, name=Henry, salary=25000, active=false, department=support, age=18, score=5.5}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, department=engineering, level=6, age=40, active=true, score=8.7, tags=senior, status=active}
10. [1] TestPerson{id=P010, department=intern, age=22, salary=28000, active=true, tags=temp, status=active, level=1, name=X, score=6.5}
11. [1] TestOrder{id=O001, date=2024-01-15, status=pending, priority=normal, discount=50, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, region=north}
12. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, amount=1, total=25.5, discount=0, region=south}
13. [1] TestOrder{id=O003, status=shipped, priority=high, region=north, customer_id=P001, amount=3, date=2024-02-01, discount=15, product_id=PROD003, total=225}
14. [1] TestOrder{id=O004, product_id=PROD004, total=299.99, status=delivered, priority=normal, amount=1, date=2024-02-05, discount=0, region=east, customer_id=P004}
15. [1] TestOrder{id=O005, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, status=confirmed, amount=1, date=2024-02-10, priority=high}
16. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0, region=west, amount=2, date=2024-02-15}
17. [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, discount=50, region=north, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007}
18. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, region=south, amount=10, total=255, priority=normal, discount=0}
19. [1] TestOrder{id=O009, priority=low, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, discount=10, region=north, amount=1, status=completed}
20. [1] TestOrder{id=O010, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, amount=1, date=2024-03-15}

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

1. [1] TestPerson{id=P001, tags=junior, status=active, department=sales, age=25, level=2, name=Alice, salary=45000, active=true, score=8.5}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, name=Bob, score=9.2}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, tags=intern, department=hr, age=16, score=6, status=inactive, level=1}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, tags=manager, status=active, department=marketing, age=45, active=true, score=7.8, level=7}
5. [1] TestPerson{id=P005, status=inactive, department=sales, name=Eve, tags=employee, level=3, age=30, salary=55000, active=false, score=8}
6. [1] TestPerson{id=P006, salary=-5000, score=0, status=active, department=qa, age=0, active=true, tags=test, level=1, name=Frank}
7. [1] TestPerson{id=P007, salary=95000, level=9, active=true, score=10, tags=executive, status=active, department=management, name=Grace, age=65}
8. [1] TestPerson{id=P008, age=18, score=5.5, tags=junior, status=inactive, level=1, name=Henry, salary=25000, active=false, department=support}
9. [1] TestPerson{id=P009, active=true, score=8.7, tags=senior, status=active, name=Ivy, salary=68000, department=engineering, level=6, age=40}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, level=1, name=X, score=6.5, department=intern}

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

1. [1] TestPerson{id=P001, tags=junior, status=active, department=sales, age=25, level=2, name=Alice, salary=45000, active=true, score=8.5}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, tags=senior, status=active, department=engineering, level=5, name=Bob, score=9.2}
3. [1] TestPerson{id=P003, name=Charlie, salary=0, active=false, tags=intern, department=hr, age=16, score=6, status=inactive, level=1}
4. [1] TestPerson{id=P004, name=Diana, salary=85000, tags=manager, status=active, department=marketing, age=45, active=true, score=7.8, level=7}
5. [1] TestPerson{id=P005, age=30, salary=55000, active=false, score=8, status=inactive, department=sales, name=Eve, tags=employee, level=3}
6. [1] TestPerson{id=P006, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active, department=qa, age=0, active=true}
7. [1] TestPerson{id=P007, active=true, score=10, tags=executive, status=active, department=management, name=Grace, age=65, salary=95000, level=9}
8. [1] TestPerson{id=P008, name=Henry, salary=25000, active=false, department=support, age=18, score=5.5, tags=junior, status=inactive, level=1}
9. [1] TestPerson{id=P009, status=active, name=Ivy, salary=68000, department=engineering, level=6, age=40, active=true, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, salary=28000, active=true, tags=temp, status=active, level=1, name=X, score=6.5, department=intern, age=22}
11. [1] TestOrder{id=O001, date=2024-01-15, status=pending, priority=normal, discount=50, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, region=north}
12. [1] TestOrder{id=O002, amount=1, total=25.5, discount=0, region=south, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low}
13. [1] TestOrder{id=O003, date=2024-02-01, discount=15, product_id=PROD003, total=225, status=shipped, priority=high, region=north, customer_id=P001, amount=3}
14. [1] TestOrder{id=O004, amount=1, date=2024-02-05, discount=0, region=east, customer_id=P004, product_id=PROD004, total=299.99, status=delivered, priority=normal}
15. [1] TestOrder{id=O005, amount=1, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001, total=999.99, status=confirmed}
16. [1] TestOrder{id=O006, amount=2, date=2024-02-15, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0, region=west}
17. [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, discount=50, region=north, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007}
18. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, region=south, amount=10, total=255, priority=normal, discount=0}
19. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, discount=10, region=north, amount=1, status=completed, priority=low}
20. [1] TestOrder{id=O010, customer_id=P006, amount=1, date=2024-03-15, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0, region=east}

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

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, age=25, level=2}
2. [1] TestPerson{id=P002, department=engineering, level=5, name=Bob, score=9.2, age=35, salary=75000, active=true, tags=senior, status=active}
3. [1] TestPerson{id=P003, tags=intern, department=hr, age=16, score=6, status=inactive, level=1, name=Charlie, salary=0, active=false}
4. [1] TestPerson{id=P004, active=true, score=7.8, level=7, name=Diana, salary=85000, tags=manager, status=active, department=marketing, age=45}
5. [1] TestPerson{id=P005, age=30, salary=55000, active=false, score=8, status=inactive, department=sales, name=Eve, tags=employee, level=3}
6. [1] TestPerson{id=P006, age=0, active=true, tags=test, level=1, name=Frank, salary=-5000, score=0, status=active, department=qa}
7. [1] TestPerson{id=P007, name=Grace, age=65, salary=95000, level=9, active=true, score=10, tags=executive, status=active, department=management}
8. [1] TestPerson{id=P008, age=18, score=5.5, tags=junior, status=inactive, level=1, name=Henry, salary=25000, active=false, department=support}
9. [1] TestPerson{id=P009, active=true, score=8.7, tags=senior, status=active, name=Ivy, salary=68000, department=engineering, level=6, age=40}
10. [1] TestPerson{id=P010, name=X, score=6.5, department=intern, age=22, salary=28000, active=true, tags=temp, status=active, level=1}

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

1. [1] TestOrder{id=O001, date=2024-01-15, status=pending, priority=normal, discount=50, customer_id=P001, product_id=PROD001, amount=2, total=1999.98, region=north}
2. [1] TestOrder{id=O002, customer_id=P002, product_id=PROD002, date=2024-01-20, status=confirmed, priority=low, amount=1, total=25.5, discount=0, region=south}
3. [1] TestOrder{id=O003, product_id=PROD003, total=225, status=shipped, priority=high, region=north, customer_id=P001, amount=3, date=2024-02-01, discount=15}
4. [1] TestOrder{id=O004, product_id=PROD004, total=299.99, status=delivered, priority=normal, amount=1, date=2024-02-05, discount=0, region=east, customer_id=P004}
5. [1] TestOrder{id=O005, total=999.99, status=confirmed, amount=1, date=2024-02-10, priority=high, discount=100, region=south, customer_id=P002, product_id=PROD001}
6. [1] TestOrder{id=O006, amount=2, date=2024-02-15, customer_id=P005, product_id=PROD005, total=999.98, status=cancelled, priority=low, discount=0, region=west}
7. [1] TestOrder{id=O007, date=2024-03-01, status=shipped, priority=urgent, customer_id=P007, product_id=PROD006, amount=4, total=600, discount=50, region=north}
8. [1] TestOrder{id=O008, total=255, priority=normal, discount=0, customer_id=P010, product_id=PROD002, date=2024-03-05, status=pending, region=south, amount=10}
9. [1] TestOrder{id=O009, amount=1, status=completed, priority=low, customer_id=P001, product_id=PROD007, total=89.99, date=2024-03-10, discount=10, region=north}
10. [1] TestOrder{id=O010, product_id=PROD001, total=75000, status=refunded, priority=urgent, discount=0, region=east, customer_id=P006, amount=1, date=2024-03-15}

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

1. [1] TestPerson{id=P001, name=Alice, salary=45000, active=true, score=8.5, tags=junior, status=active, department=sales, age=25, level=2}
2. [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, level=5, name=Bob, score=9.2, age=35, salary=75000, active=true}
3. [1] TestPerson{id=P003, score=6, status=inactive, level=1, name=Charlie, salary=0, active=false, tags=intern, department=hr, age=16}
4. [1] TestPerson{id=P004, department=marketing, age=45, active=true, score=7.8, level=7, name=Diana, salary=85000, tags=manager, status=active}
5. [1] TestPerson{id=P005, level=3, age=30, salary=55000, active=false, score=8, status=inactive, department=sales, name=Eve, tags=employee}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, score=0, status=active, department=qa, age=0, active=true, tags=test, level=1}
7. [1] TestPerson{id=P007, status=active, department=management, name=Grace, age=65, salary=95000, level=9, active=true, score=10, tags=executive}
8. [1] TestPerson{id=P008, active=false, department=support, age=18, score=5.5, tags=junior, status=inactive, level=1, name=Henry, salary=25000}
9. [1] TestPerson{id=P009, name=Ivy, salary=68000, department=engineering, level=6, age=40, active=true, score=8.7, tags=senior, status=active}
10. [1] TestPerson{id=P010, age=22, salary=28000, active=true, tags=temp, status=active, level=1, name=X, score=6.5, department=intern}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 2 (10.5%)
- **Tokens générés**: 18
- **Faits traités**: 27
