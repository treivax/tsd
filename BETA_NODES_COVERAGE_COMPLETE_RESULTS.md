# RAPPORT COMPLET DE COUVERTURE DES NŒUDS BETA
================================================

**📊 Tests exécutés:** 3
**✅ Tests réussis:** 3 (100.0%)
**🧠 Score sémantique moyen:** 100.0%
**📅 Date d'exécution:** 2025-11-18 14:56:07

## 🎯 NŒUDS BETA ANALYSÉS
| Type de Nœud | Tests | Succès | Score Sémantique |
|---------------|--------|--------|------------------|
| ExistsNode | 1 | 1 | 100.0% |
| JoinNode | 1 | 1 | 100.0% |
| NotNode | 1 | 1 | 100.0% |

## 🧪 TEST 1: exists_simple
---

### 📋 Informations générales
- **Description:** Test existence simple
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/exists_simple.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/exists_simple.facts`
- **Temps d'exécution:** 531.654µs
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 100.0%
- **Actions valides:** ✅
- **Jointures valides:** ✅
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

### 📜 Règles analysées
#### Règle 1
- **Texte original:** `{p: Person} / EXISTS (o: Order / o.customer_id == p.id) ==> person_has_orders(p.id)`
- **Action:** person_has_orders
- **Type de nœud:** ExistsNode
- **Type sémantique:** existence
- **Complexité:** simple
- **Variables:**
  - p (Person): primary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Person (type_Person)
│   ├── Order (type_Order)
│
├── 🔗 BetaNodes (Jointures)
│   ├── rule_0_exists
│   │   └── Type: *rete.ExistsNode
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: person_has_orders
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Person[id=P001, name=Alice]
Person[id=P002, name=Bob]
Order[customer_id=P001, amount=100]

```

**Total faits:** 3

- **Person:** 2 faits
- **Order:** 1 faits

**📋 Détail des faits parsés:**
1. **Person[P001]** - `Person[id=P001, name=Alice]`
2. **Person[P002]** - `Person[id=P002, name=Bob]`
3. **Order[fact_Order_3]** - `Order[id=fact_Order_3, customer_id=P001, amount=100]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| person_has_orders | 1 | AlphaNode | ✅ |

#### 🎯 Activation détaillée: `person_has_orders`
- **Nombre de déclenchements:** 1
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`p`**: Person[P001] - `Person[id=P001, name=Alice]`

### 🎯 Comparaison attendu vs observé
#### Actions
| Action | Attendu | Observé | Statut |
|--------|---------|---------|--------|
| person_has_orders | 1-1 | 1 | ✅ |

#### 📋 TOKENS COMBINÉS ATTENDUS vs OBTENUS

**🎯 Action `person_has_orders`:**
- **Description:** Une personne (Alice) a une commande existante
- **Variables de la règle:** p

**📍 TOKENS COMBINÉS ATTENDUS:**
- **Nombre de tokens attendus:** 1-1
- **Token attendu 1:**
  * `p`: Person[P001] - `Person[id=P001, name=Alice]`

**📊 TOKENS COMBINÉS OBTENUS:**
- **Nombre de tokens obtenus:** 1
- **Token obtenu 1:**
  * `p`: Person[P001] - `Person[id=P001, name=Alice]`

**🎯 RÉSULTAT:** ✅ SUCCÈS
- ✅ Nombre de tokens correct

---

## 🧪 TEST 2: join_simple
---

### 📋 Informations générales
- **Description:** Test jointure simple entre deux faits
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_simple.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_simple.facts`
- **Temps d'exécution:** 466.332µs
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 100.0%
- **Actions valides:** ✅
- **Jointures valides:** ✅
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

### 📜 Règles analysées
#### Règle 1
- **Texte original:** `{p: Person, o: Order} / p.id == o.customer_id ==> customer_order_match(p.id, o.id)`
- **Action:** customer_order_match
- **Type de nœud:** JoinNode
- **Type sémantique:** comparison
- **Complexité:** complex
- **Variables:**
  - p (Person): primary
  - o (Order): secondary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Person (type_Person)
│   ├── Order (type_Order)
│
├── 🔗 BetaNodes (Jointures)
│   ├── rule_0_join
│   │   ├── Variables: p ⋈ o
│   │   ├── Conditions: 1
│   │   └── Type: JoinNode
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: customer_order_match
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Person[id=P001, name=Alice, age=25]
Person[id=P002, name=Bob, age=30]
Order[id=O001, customer_id=P001, amount=100]
Order[id=O002, customer_id=P002, amount=200]

```

**Total faits:** 4

- **Person:** 2 faits
- **Order:** 2 faits

**📋 Détail des faits parsés:**
1. **Person[P001]** - `Person[id=P001, name=Alice, age=25]`
2. **Person[P002]** - `Person[id=P002, name=Bob, age=30]`
3. **Order[O001]** - `Order[id=O001, customer_id=P001, amount=100]`
4. **Order[O002]** - `Order[id=O002, customer_id=P002, amount=200]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| customer_order_match | 2 | AlphaNode | ✅ |

#### 🎯 Activation détaillée: `customer_order_match`
- **Nombre de déclenchements:** 2
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`p`**: Person[P001] - `Person[id=P001, name=Alice, age=25]`
- **`o`**: Order[O001] - `Order[id=O001, customer_id=P001, amount=100]`
- **Association:** Person[P001] ⋈ Order[O001]

##### Token combiné 2
- **`p`**: Person[P002] - `Person[id=P002, name=Bob, age=30]`
- **`o`**: Order[O002] - `Order[amount=200, id=O002, customer_id=P002]`
- **Association:** Person[P002] ⋈ Order[O002]

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | p <-> o | 2 | inner | ✅ |

### 🎯 Comparaison attendu vs observé
#### Actions
| Action | Attendu | Observé | Statut |
|--------|---------|---------|--------|
| customer_order_match | 2-2 | 2 | ✅ |

#### 📋 TOKENS COMBINÉS ATTENDUS vs OBTENUS

**🎯 Action `customer_order_match`:**
- **Description:** Deux customers avec leurs commandes matchent
- **Variables de la règle:** p, o

**📍 TOKENS COMBINÉS ATTENDUS:**
- **Nombre de tokens attendus:** 2-2
- **Token attendu 1:**
  * `p`: Person[P001] - `Person[id=P001, name=Alice, age=25]`
  * `o`: Order[O001] - `Order[id=O001, customer_id=P001, amount=100]`
- **Token attendu 2:**
  * `p`: Person[P002] - `Person[age=30, id=P002, name=Bob]`
  * `o`: Order[O002] - `Order[id=O002, customer_id=P002, amount=200]`

**📊 TOKENS COMBINÉS OBTENUS:**
- **Nombre de tokens obtenus:** 2
- **Token obtenu 1:**
  * `p`: Person[P001] - `Person[id=P001, name=Alice, age=25]`
  * `o`: Order[O001] - `Order[customer_id=P001, amount=100, id=O001]`
- **Token obtenu 2:**
  * `p`: Person[P002] - `Person[id=P002, name=Bob, age=30]`
  * `o`: Order[O002] - `Order[id=O002, customer_id=P002, amount=200]`

**🎯 RÉSULTAT:** ✅ SUCCÈS
- ✅ Nombre de tokens correct

---

## 🧪 TEST 3: not_simple
---

### 📋 Informations générales
- **Description:** Test négation simple
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/not_simple.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/not_simple.facts`
- **Temps d'exécution:** 275.485µs
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 100.0%
- **Actions valides:** ✅
- **Jointures valides:** ✅
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

### 📜 Règles analysées
#### Règle 1
- **Texte original:** `{p: Person} / NOT (p.active == false) ==> active_person(p.id)`
- **Action:** active_person
- **Type de nœud:** NotNode
- **Type sémantique:** negation
- **Complexité:** simple
- **Variables:**
  - p (Person): primary

### 🕸️ Structure du réseau RETE

```
RÉSEAU RETE BETA - STRUCTURE HIÉRARCHIQUE
==========================================

🌳 RootNode
│
├── 📁 TypeNodes
│   ├── Person (type_Person)
│
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: NOT(...) [Négation]
│   │   └── Variable: p
│
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: active_person
```

### 📄 Faits traités
**📄 Contenu fichier facts:**
```
Person[id=P001, name=Alice, active=true]
Person[id=P002, name=Bob, active=false]

```

**Total faits:** 2

- **Person:** 2 faits

**📋 Détail des faits parsés:**
1. **Person[P001]** - `Person[name=Alice, active=true, id=P001]`
2. **Person[P002]** - `Person[id=P002, name=Bob, active=false]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| active_person | 1 | AlphaNode | ✅ |

#### 🎯 Activation détaillée: `active_person`
- **Nombre de déclenchements:** 1
- **Type de nœud déclencheur:** AlphaNode

**📋 TOKENS COMBINÉS activant l'action:**

##### Token combiné 1
- **`p`**: Person[P001] - `Person[id=P001, name=Alice, active=true]`

### 🚫 Analyse des négations (NotNodes)
| Nœud | Condition Niée | Faits Filtrés | Type | Validation |
|------|----------------|---------------|------|------------|
| alpha_rule_0_alpha | map[condition:map[left:map[field:active object:p type:fieldAccess] operator:== right:map[type:boolean value:false] type:comparison] negated:true type:negation] | 0 | simple | ❌ |

### 🎯 Comparaison attendu vs observé
#### Actions
| Action | Attendu | Observé | Statut |
|--------|---------|---------|--------|
| active_person | 1-1 | 1 | ✅ |

#### 📋 TOKENS COMBINÉS ATTENDUS vs OBTENUS

**🎯 Action `active_person`:**
- **Description:** Une personne active (Alice) passe le filtre NOT
- **Variables de la règle:** p

**📍 TOKENS COMBINÉS ATTENDUS:**
- **Nombre de tokens attendus:** 1-1
- **Token attendu 1:**
  * `p`: Person[P001] - `Person[active=true, id=P001, name=Alice]`

**📊 TOKENS COMBINÉS OBTENUS:**
- **Nombre de tokens obtenus:** 1
- **Token obtenu 1:**
  * `p`: Person[P001] - `Person[id=P001, name=Alice, active=true]`

**🎯 RÉSULTAT:** ✅ SUCCÈS
- ✅ Nombre de tokens correct

---

## 💡 RECOMMANDATIONS
### Amélioration de la couverture Beta
### Prochaines étapes
1. **Ajouter plus de tests complexes** avec jointures multiples
2. **Tester les négations imbriquées** et conditions complexes
3. **Valider les performances** des nœuds Beta avec de gros volumes
4. **Enrichir la validation sémantique** avec plus de critères
