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

1. [1] TestPerson{id=P001, name=Alice, age=25, active=true, tags=junior, salary=45000, score=8.5, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, name=Bob, tags=senior, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, tags=manager, status=active, department=marketing, level=7, name=Diana, age=45}
5. [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, status=inactive}
6. [1] TestPerson{id=P006, score=0, tags=test, status=active, name=Frank, salary=-5000, active=true, department=qa, level=1, age=0}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, status=active, department=management, level=9, age=65, active=true, score=10, tags=executive}
8. [1] TestPerson{id=P008, department=support, name=Henry, active=false, score=5.5, tags=junior, level=1, age=18, salary=25000, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, active=true, status=active, department=engineering, level=6, age=40, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, status=active, department=intern, age=22, tags=temp, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

**9 résultats obtenus**:

1. **Token 1**:
   - Fait 1: [1] TestPerson{id=P001, name=Alice, age=25, active=true, tags=junior, salary=45000, score=8.5, status=active, department=sales, level=2}

2. **Token 2**:
   - Fait 1: [1] TestPerson{id=P002, level=5, age=35, salary=75000, active=true, score=9.2, name=Bob, tags=senior, status=active, department=engineering}

3. **Token 3**:
   - Fait 1: [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, level=1}

4. **Token 4**:
   - Fait 1: [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, tags=manager, status=active, department=marketing, level=7, name=Diana, age=45}

5. **Token 5**:
   - Fait 1: [1] TestPerson{id=P005, salary=55000, status=inactive, active=false, score=8, tags=employee, department=sales, level=3, name=Eve, age=30}

6. **Token 6**:
   - Fait 1: [1] TestPerson{id=P007, name=Grace, salary=95000, status=active, department=management, level=9, age=65, active=true, score=10, tags=executive}

7. **Token 7**:
   - Fait 1: [1] TestPerson{id=P008, salary=25000, status=inactive, department=support, name=Henry, active=false, score=5.5, tags=junior, level=1, age=18}

8. **Token 8**:
   - Fait 1: [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, name=Ivy, active=true, status=active, department=engineering, level=6}

9. **Token 9**:
   - Fait 1: [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, status=active, department=intern, age=22, tags=temp, level=1}

### 📊 STATISTIQUES

- **Taux de déclenchement**: 9/10 (90.0%)
- **Efficacité**: 🟢 Très élevée

---

## 🎯 RÈGLE 1: not_cancelled_order

**Condition**: `NOT (o.status == "cancelled")`
**Types concernés**: [TestOrder]
**Terminal**: rule_1_terminal

### 📥 FAITS SOUMIS

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, amount=2, date=2024-01-15, priority=normal, region=north}
2. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, discount=0, product_id=PROD002, amount=1}
3. [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, product_id=PROD003, discount=15}
4. [1] TestOrder{id=O004, product_id=PROD004, date=2024-02-05, priority=normal, discount=0, customer_id=P004, amount=1, total=299.99, status=delivered, region=east}
5. [1] TestOrder{id=O005, customer_id=P002, product_id=PROD001, discount=100, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, region=south}
6. [1] TestOrder{id=O006, discount=0, region=west, amount=2, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, total=999.98, priority=low}
7. [1] TestOrder{id=O007, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, discount=50, customer_id=P007, date=2024-03-01, status=shipped}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, status=pending, region=south}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, date=2024-03-10, priority=low, region=north, amount=1, total=89.99, status=completed, discount=10}
10. [1] TestOrder{id=O010, date=2024-03-15, discount=0, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, region=east, product_id=PROD001}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 🎯 RÈGLE 2: not_low_salary

**Condition**: `NOT (p.salary < 30000)`
**Types concernés**: [TestPerson]
**Terminal**: rule_2_terminal

### 📥 FAITS SOUMIS

1. [1] TestPerson{id=P001, name=Alice, age=25, active=true, tags=junior, salary=45000, score=8.5, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, name=Bob, tags=senior, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, score=6, tags=intern, status=inactive, department=hr, level=1, name=Charlie, age=16, salary=0, active=false}
4. [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, status=active, department=marketing, level=7, name=Diana}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, status=inactive, active=false, score=8, tags=employee, department=sales, level=3}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, department=qa, level=1, age=0, score=0, tags=test, status=active}
7. [1] TestPerson{id=P007, department=management, level=9, age=65, active=true, score=10, tags=executive, name=Grace, salary=95000, status=active}
8. [1] TestPerson{id=P008, age=18, salary=25000, status=inactive, department=support, name=Henry, active=false, score=5.5, tags=junior, level=1}
9. [1] TestPerson{id=P009, active=true, status=active, department=engineering, level=6, age=40, salary=68000, score=8.7, tags=senior, name=Ivy}
10. [1] TestPerson{id=P010, salary=28000, active=true, score=6.5, status=active, department=intern, age=22, tags=temp, level=1, name=X}

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

1. [1] TestOrder{id=O001, total=1999.98, status=pending, discount=50, amount=2, date=2024-01-15, priority=normal, region=north, customer_id=P001, product_id=PROD001}
2. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, discount=0, product_id=PROD002, amount=1}
3. [1] TestOrder{id=O003, region=north, customer_id=P001, product_id=PROD003, discount=15, amount=3, total=225, date=2024-02-01, status=shipped, priority=high}
4. [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, region=east, product_id=PROD004, date=2024-02-05, priority=normal, discount=0, customer_id=P004}
5. [1] TestOrder{id=O005, status=confirmed, priority=high, region=south, customer_id=P002, product_id=PROD001, discount=100, amount=1, total=999.99, date=2024-02-10}
6. [1] TestOrder{id=O006, total=999.98, priority=low, discount=0, region=west, amount=2, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, discount=50, customer_id=P007, date=2024-03-01, status=shipped, priority=urgent, region=north}
8. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, status=pending, region=south}
9. [1] TestOrder{id=O009, amount=1, total=89.99, status=completed, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, priority=low, region=north}
10. [1] TestOrder{id=O010, product_id=PROD001, date=2024-03-15, discount=0, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, region=east}

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

1. [1] TestProduct{id=PROD001, supplier=TechSupply, available=true, rating=4.5, keywords=computer, brand=TechCorp, stock=50, name=Laptop, category=electronics, price=999.99}
2. [1] TestProduct{id=PROD002, name=Mouse, category=accessories, price=25.5, available=true, rating=4, keywords=peripheral, brand=TechCorp, stock=200, supplier=TechSupply}
3. [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, available=false, rating=3.5, brand=KeyTech, stock=0, price=75, keywords=typing, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, category=electronics, price=299.99, rating=4.8, brand=ScreenPro, stock=30, supplier=ScreenSupply, name=Monitor, available=true, keywords=display}
5. [1] TestProduct{id=PROD005, available=false, rating=2, brand=OldTech, stock=0, keywords=obsolete, supplier=OldSupply, name=OldKeyboard, category=accessories, price=8.5}
6. [1] TestProduct{id=PROD006, category=audio, price=150, available=true, rating=4.6, brand=AudioMax, name=Headphones, keywords=sound, stock=75, supplier=AudioSupply}
7. [1] TestProduct{id=PROD007, available=true, rating=3.8, brand=CamTech, supplier=CamSupply, name=Webcam, category=electronics, keywords=video, stock=25, price=89.99}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, active=true, tags=junior, salary=45000, score=8.5, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, tags=senior, status=active, department=engineering, level=5, age=35, salary=75000, active=true, score=9.2, name=Bob}
3. [1] TestPerson{id=P003, age=16, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, level=1, name=Charlie}
4. [1] TestPerson{id=P004, active=true, score=7.8, tags=manager, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000}
5. [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, department=qa, level=1, age=0, score=0, tags=test, status=active}
7. [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, name=Grace, salary=95000, status=active, department=management, level=9}
8. [1] TestPerson{id=P008, level=1, age=18, salary=25000, status=inactive, department=support, name=Henry, active=false, score=5.5, tags=junior}
9. [1] TestPerson{id=P009, name=Ivy, active=true, status=active, department=engineering, level=6, age=40, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, level=1, name=X, salary=28000, active=true, score=6.5, status=active, department=intern, age=22, tags=temp}

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

1. [1] TestOrder{id=O001, amount=2, date=2024-01-15, priority=normal, region=north, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50}
2. [1] TestOrder{id=O002, product_id=PROD002, amount=1, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, discount=0}
3. [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, product_id=PROD003, discount=15}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, status=delivered, region=east, product_id=PROD004, date=2024-02-05, priority=normal, discount=0}
5. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, region=south, customer_id=P002, product_id=PROD001, discount=100}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, priority=low, discount=0, region=west, amount=2, date=2024-02-15, status=cancelled}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, discount=50, customer_id=P007, date=2024-03-01, status=shipped, priority=urgent, region=north}
8. [1] TestOrder{id=O008, discount=0, status=pending, region=south, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal}
9. [1] TestOrder{id=O009, status=completed, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, priority=low, region=north, amount=1, total=89.99}
10. [1] TestOrder{id=O010, amount=1, total=75000, status=refunded, priority=urgent, region=east, product_id=PROD001, date=2024-03-15, discount=0, customer_id=P006}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, active=true, tags=junior, salary=45000, score=8.5, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, name=Bob, tags=senior, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, tags=manager, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true, score=7.8}
5. [1] TestPerson{id=P005, department=sales, level=3, name=Eve, age=30, salary=55000, status=inactive, active=false, score=8, tags=employee}
6. [1] TestPerson{id=P006, age=0, score=0, tags=test, status=active, name=Frank, salary=-5000, active=true, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, name=Grace, salary=95000, status=active, department=management, level=9}
8. [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, tags=junior, level=1, age=18, salary=25000, status=inactive, department=support}
9. [1] TestPerson{id=P009, name=Ivy, active=true, status=active, department=engineering, level=6, age=40, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, department=intern, age=22, tags=temp, level=1, name=X, salary=28000, active=true, score=6.5, status=active}

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

1. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, amount=2, date=2024-01-15, priority=normal, region=north}
2. [1] TestOrder{id=O002, status=confirmed, priority=low, region=south, customer_id=P002, discount=0, product_id=PROD002, amount=1, total=25.5, date=2024-01-20}
3. [1] TestOrder{id=O003, priority=high, region=north, customer_id=P001, product_id=PROD003, discount=15, amount=3, total=225, date=2024-02-01, status=shipped}
4. [1] TestOrder{id=O004, priority=normal, discount=0, customer_id=P004, amount=1, total=299.99, status=delivered, region=east, product_id=PROD004, date=2024-02-05}
5. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, region=south, customer_id=P002, product_id=PROD001, discount=100}
6. [1] TestOrder{id=O006, amount=2, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, total=999.98, priority=low, discount=0, region=west}
7. [1] TestOrder{id=O007, product_id=PROD006, amount=4, total=600, discount=50, customer_id=P007, date=2024-03-01, status=shipped, priority=urgent, region=north}
8. [1] TestOrder{id=O008, date=2024-03-05, priority=normal, discount=0, status=pending, region=south, customer_id=P010, product_id=PROD002, amount=10, total=255}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, date=2024-03-10, priority=low, region=north, amount=1, total=89.99, status=completed, discount=10}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, region=east, product_id=PROD001, date=2024-03-15, discount=0}

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

1. [1] TestPerson{id=P001, age=25, active=true, tags=junior, salary=45000, score=8.5, status=active, department=sales, level=2, name=Alice}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, name=Bob, tags=senior, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, score=6, tags=intern, status=inactive, department=hr, level=1, name=Charlie, age=16, salary=0, active=false}
4. [1] TestPerson{id=P004, age=45, salary=85000, active=true, score=7.8, tags=manager, status=active, department=marketing, level=7, name=Diana}
5. [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, department=qa, level=1, age=0, score=0, tags=test, status=active}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, status=active, department=management, level=9, age=65, active=true, score=10, tags=executive}
8. [1] TestPerson{id=P008, status=inactive, department=support, name=Henry, active=false, score=5.5, tags=junior, level=1, age=18, salary=25000}
9. [1] TestPerson{id=P009, name=Ivy, active=true, status=active, department=engineering, level=6, age=40, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, department=intern, age=22, tags=temp, level=1, name=X, salary=28000, active=true, score=6.5, status=active}

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

1. [1] TestProduct{id=PROD001, brand=TechCorp, stock=50, name=Laptop, category=electronics, price=999.99, supplier=TechSupply, available=true, rating=4.5, keywords=computer}
2. [1] TestProduct{id=PROD002, brand=TechCorp, stock=200, supplier=TechSupply, name=Mouse, category=accessories, price=25.5, available=true, rating=4, keywords=peripheral}
3. [1] TestProduct{id=PROD003, name=Keyboard, category=accessories, available=false, rating=3.5, brand=KeyTech, stock=0, price=75, keywords=typing, supplier=KeySupply}
4. [1] TestProduct{id=PROD004, price=299.99, rating=4.8, brand=ScreenPro, stock=30, supplier=ScreenSupply, name=Monitor, available=true, keywords=display, category=electronics}
5. [1] TestProduct{id=PROD005, stock=0, keywords=obsolete, supplier=OldSupply, name=OldKeyboard, category=accessories, price=8.5, available=false, rating=2, brand=OldTech}
6. [1] TestProduct{id=PROD006, rating=4.6, brand=AudioMax, name=Headphones, keywords=sound, stock=75, supplier=AudioSupply, category=audio, price=150, available=true}
7. [1] TestProduct{id=PROD007, stock=25, price=89.99, available=true, rating=3.8, brand=CamTech, supplier=CamSupply, name=Webcam, category=electronics, keywords=video}

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

1. [1] TestPerson{id=P001, salary=45000, score=8.5, status=active, department=sales, level=2, name=Alice, age=25, active=true, tags=junior}
2. [1] TestPerson{id=P002, name=Bob, tags=senior, status=active, department=engineering, level=5, age=35, salary=75000, active=true, score=9.2}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true, score=7.8, tags=manager, status=active}
5. [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, status=inactive}
6. [1] TestPerson{id=P006, status=active, name=Frank, salary=-5000, active=true, department=qa, level=1, age=0, score=0, tags=test}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, status=active, department=management, level=9, age=65, active=true, score=10, tags=executive}
8. [1] TestPerson{id=P008, salary=25000, status=inactive, department=support, name=Henry, active=false, score=5.5, tags=junior, level=1, age=18}
9. [1] TestPerson{id=P009, name=Ivy, active=true, status=active, department=engineering, level=6, age=40, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, score=6.5, status=active, department=intern, age=22, tags=temp, level=1, name=X, salary=28000, active=true}

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

1. [1] TestOrder{id=O001, amount=2, date=2024-01-15, priority=normal, region=north, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50}
2. [1] TestOrder{id=O002, status=confirmed, priority=low, region=south, customer_id=P002, discount=0, product_id=PROD002, amount=1, total=25.5, date=2024-01-20}
3. [1] TestOrder{id=O003, status=shipped, priority=high, region=north, customer_id=P001, product_id=PROD003, discount=15, amount=3, total=225, date=2024-02-01}
4. [1] TestOrder{id=O004, product_id=PROD004, date=2024-02-05, priority=normal, discount=0, customer_id=P004, amount=1, total=299.99, status=delivered, region=east}
5. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, region=south, customer_id=P002, product_id=PROD001, discount=100}
6. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, priority=low, discount=0, region=west, amount=2, date=2024-02-15, status=cancelled}
7. [1] TestOrder{id=O007, status=shipped, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, discount=50, customer_id=P007, date=2024-03-01}
8. [1] TestOrder{id=O008, status=pending, region=south, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0}
9. [1] TestOrder{id=O009, amount=1, total=89.99, status=completed, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, priority=low, region=north}
10. [1] TestOrder{id=O010, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, region=east, product_id=PROD001, date=2024-03-15, discount=0}

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

1. [1] TestPerson{id=P001, status=active, department=sales, level=2, name=Alice, age=25, active=true, tags=junior, salary=45000, score=8.5}
2. [1] TestPerson{id=P002, name=Bob, tags=senior, status=active, department=engineering, level=5, age=35, salary=75000, active=true, score=9.2}
3. [1] TestPerson{id=P003, level=1, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr}
4. [1] TestPerson{id=P004, level=7, name=Diana, age=45, salary=85000, active=true, score=7.8, tags=manager, status=active, department=marketing}
5. [1] TestPerson{id=P005, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, status=inactive, active=false}
6. [1] TestPerson{id=P006, status=active, name=Frank, salary=-5000, active=true, department=qa, level=1, age=0, score=0, tags=test}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, status=active, department=management, level=9, age=65, active=true, score=10, tags=executive}
8. [1] TestPerson{id=P008, level=1, age=18, salary=25000, status=inactive, department=support, name=Henry, active=false, score=5.5, tags=junior}
9. [1] TestPerson{id=P009, age=40, salary=68000, score=8.7, tags=senior, name=Ivy, active=true, status=active, department=engineering, level=6}
10. [1] TestPerson{id=P010, level=1, name=X, salary=28000, active=true, score=6.5, status=active, department=intern, age=22, tags=temp}
11. [1] TestOrder{id=O001, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, amount=2, date=2024-01-15, priority=normal, region=north}
12. [1] TestOrder{id=O002, total=25.5, date=2024-01-20, status=confirmed, priority=low, region=south, customer_id=P002, discount=0, product_id=PROD002, amount=1}
13. [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, product_id=PROD003, discount=15}
14. [1] TestOrder{id=O004, amount=1, total=299.99, status=delivered, region=east, product_id=PROD004, date=2024-02-05, priority=normal, discount=0, customer_id=P004}
15. [1] TestOrder{id=O005, date=2024-02-10, status=confirmed, priority=high, region=south, customer_id=P002, product_id=PROD001, discount=100, amount=1, total=999.99}
16. [1] TestOrder{id=O006, customer_id=P005, product_id=PROD005, total=999.98, priority=low, discount=0, region=west, amount=2, date=2024-02-15, status=cancelled}
17. [1] TestOrder{id=O007, customer_id=P007, date=2024-03-01, status=shipped, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, discount=50}
18. [1] TestOrder{id=O008, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, status=pending, region=south}
19. [1] TestOrder{id=O009, total=89.99, status=completed, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, priority=low, region=north, amount=1}
20. [1] TestOrder{id=O010, amount=1, total=75000, status=refunded, priority=urgent, region=east, product_id=PROD001, date=2024-03-15, discount=0, customer_id=P006}

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

1. [1] TestPerson{id=P001, age=25, active=true, tags=junior, salary=45000, score=8.5, status=active, department=sales, level=2, name=Alice}
2. [1] TestPerson{id=P002, score=9.2, name=Bob, tags=senior, status=active, department=engineering, level=5, age=35, salary=75000, active=true}
3. [1] TestPerson{id=P003, status=inactive, department=hr, level=1, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern}
4. [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, tags=manager, status=active, department=marketing, level=7, name=Diana, age=45}
5. [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, status=inactive}
6. [1] TestPerson{id=P006, age=0, score=0, tags=test, status=active, name=Frank, salary=-5000, active=true, department=qa, level=1}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, status=active, department=management, level=9, age=65, active=true, score=10, tags=executive}
8. [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, tags=junior, level=1, age=18, salary=25000, status=inactive, department=support}
9. [1] TestPerson{id=P009, name=Ivy, active=true, status=active, department=engineering, level=6, age=40, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, status=active, department=intern, age=22, tags=temp, level=1, name=X, salary=28000, active=true, score=6.5}

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

1. [1] TestPerson{id=P001, level=2, name=Alice, age=25, active=true, tags=junior, salary=45000, score=8.5, status=active, department=sales}
2. [1] TestPerson{id=P002, age=35, salary=75000, active=true, score=9.2, name=Bob, tags=senior, status=active, department=engineering, level=5}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, score=7.8, tags=manager, status=active, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true}
5. [1] TestPerson{id=P005, name=Eve, age=30, salary=55000, status=inactive, active=false, score=8, tags=employee, department=sales, level=3}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, department=qa, level=1, age=0, score=0, tags=test, status=active}
7. [1] TestPerson{id=P007, tags=executive, name=Grace, salary=95000, status=active, department=management, level=9, age=65, active=true, score=10}
8. [1] TestPerson{id=P008, department=support, name=Henry, active=false, score=5.5, tags=junior, level=1, age=18, salary=25000, status=inactive}
9. [1] TestPerson{id=P009, name=Ivy, active=true, status=active, department=engineering, level=6, age=40, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, status=active, department=intern, age=22, tags=temp, level=1}
11. [1] TestOrder{id=O001, region=north, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, amount=2, date=2024-01-15, priority=normal}
12. [1] TestOrder{id=O002, status=confirmed, priority=low, region=south, customer_id=P002, discount=0, product_id=PROD002, amount=1, total=25.5, date=2024-01-20}
13. [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, product_id=PROD003, discount=15}
14. [1] TestOrder{id=O004, product_id=PROD004, date=2024-02-05, priority=normal, discount=0, customer_id=P004, amount=1, total=299.99, status=delivered, region=east}
15. [1] TestOrder{id=O005, status=confirmed, priority=high, region=south, customer_id=P002, product_id=PROD001, discount=100, amount=1, total=999.99, date=2024-02-10}
16. [1] TestOrder{id=O006, amount=2, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, total=999.98, priority=low, discount=0, region=west}
17. [1] TestOrder{id=O007, region=north, product_id=PROD006, amount=4, total=600, discount=50, customer_id=P007, date=2024-03-01, status=shipped, priority=urgent}
18. [1] TestOrder{id=O008, priority=normal, discount=0, status=pending, region=south, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05}
19. [1] TestOrder{id=O009, status=completed, discount=10, customer_id=P001, product_id=PROD007, date=2024-03-10, priority=low, region=north, amount=1, total=89.99}
20. [1] TestOrder{id=O010, date=2024-03-15, discount=0, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, region=east, product_id=PROD001}

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

1. [1] TestPerson{id=P001, level=2, name=Alice, age=25, active=true, tags=junior, salary=45000, score=8.5, status=active, department=sales}
2. [1] TestPerson{id=P002, status=active, department=engineering, level=5, age=35, salary=75000, active=true, score=9.2, name=Bob, tags=senior}
3. [1] TestPerson{id=P003, age=16, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, level=1, name=Charlie}
4. [1] TestPerson{id=P004, department=marketing, level=7, name=Diana, age=45, salary=85000, active=true, score=7.8, tags=manager, status=active}
5. [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, status=inactive}
6. [1] TestPerson{id=P006, age=0, score=0, tags=test, status=active, name=Frank, salary=-5000, active=true, department=qa, level=1}
7. [1] TestPerson{id=P007, age=65, active=true, score=10, tags=executive, name=Grace, salary=95000, status=active, department=management, level=9}
8. [1] TestPerson{id=P008, name=Henry, active=false, score=5.5, tags=junior, level=1, age=18, salary=25000, status=inactive, department=support}
9. [1] TestPerson{id=P009, name=Ivy, active=true, status=active, department=engineering, level=6, age=40, salary=68000, score=8.7, tags=senior}
10. [1] TestPerson{id=P010, salary=28000, active=true, score=6.5, status=active, department=intern, age=22, tags=temp, level=1, name=X}

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

1. [1] TestOrder{id=O001, priority=normal, region=north, customer_id=P001, product_id=PROD001, total=1999.98, status=pending, discount=50, amount=2, date=2024-01-15}
2. [1] TestOrder{id=O002, priority=low, region=south, customer_id=P002, discount=0, product_id=PROD002, amount=1, total=25.5, date=2024-01-20, status=confirmed}
3. [1] TestOrder{id=O003, amount=3, total=225, date=2024-02-01, status=shipped, priority=high, region=north, customer_id=P001, product_id=PROD003, discount=15}
4. [1] TestOrder{id=O004, customer_id=P004, amount=1, total=299.99, status=delivered, region=east, product_id=PROD004, date=2024-02-05, priority=normal, discount=0}
5. [1] TestOrder{id=O005, amount=1, total=999.99, date=2024-02-10, status=confirmed, priority=high, region=south, customer_id=P002, product_id=PROD001, discount=100}
6. [1] TestOrder{id=O006, discount=0, region=west, amount=2, date=2024-02-15, status=cancelled, customer_id=P005, product_id=PROD005, total=999.98, priority=low}
7. [1] TestOrder{id=O007, customer_id=P007, date=2024-03-01, status=shipped, priority=urgent, region=north, product_id=PROD006, amount=4, total=600, discount=50}
8. [1] TestOrder{id=O008, region=south, customer_id=P010, product_id=PROD002, amount=10, total=255, date=2024-03-05, priority=normal, discount=0, status=pending}
9. [1] TestOrder{id=O009, customer_id=P001, product_id=PROD007, date=2024-03-10, priority=low, region=north, amount=1, total=89.99, status=completed, discount=10}
10. [1] TestOrder{id=O010, discount=0, customer_id=P006, amount=1, total=75000, status=refunded, priority=urgent, region=east, product_id=PROD001, date=2024-03-15}

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

1. [1] TestPerson{id=P001, name=Alice, age=25, active=true, tags=junior, salary=45000, score=8.5, status=active, department=sales, level=2}
2. [1] TestPerson{id=P002, name=Bob, tags=senior, status=active, department=engineering, level=5, age=35, salary=75000, active=true, score=9.2}
3. [1] TestPerson{id=P003, name=Charlie, age=16, salary=0, active=false, score=6, tags=intern, status=inactive, department=hr, level=1}
4. [1] TestPerson{id=P004, salary=85000, active=true, score=7.8, tags=manager, status=active, department=marketing, level=7, name=Diana, age=45}
5. [1] TestPerson{id=P005, active=false, score=8, tags=employee, department=sales, level=3, name=Eve, age=30, salary=55000, status=inactive}
6. [1] TestPerson{id=P006, name=Frank, salary=-5000, active=true, department=qa, level=1, age=0, score=0, tags=test, status=active}
7. [1] TestPerson{id=P007, name=Grace, salary=95000, status=active, department=management, level=9, age=65, active=true, score=10, tags=executive}
8. [1] TestPerson{id=P008, tags=junior, level=1, age=18, salary=25000, status=inactive, department=support, name=Henry, active=false, score=5.5}
9. [1] TestPerson{id=P009, active=true, status=active, department=engineering, level=6, age=40, salary=68000, score=8.7, tags=senior, name=Ivy}
10. [1] TestPerson{id=P010, name=X, salary=28000, active=true, score=6.5, status=active, department=intern, age=22, tags=temp, level=1}

**Total**: 10 faits soumis

### 📤 RÉSULTATS TERMINAL

Aucun résultat (règle non déclenchée)

### 📊 STATISTIQUES

- **Taux de déclenchement**: 0/10 (0.0%)
- **Efficacité**: ⚫ Très faible

---

## 📊 RÉSUMÉ GLOBAL

- **Terminaux totaux**: 19
- **Terminaux actifs**: 1 (5.3%)
- **Tokens générés**: 9
- **Faits traités**: 27
