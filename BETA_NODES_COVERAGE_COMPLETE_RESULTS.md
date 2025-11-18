# RAPPORT COMPLET DE COUVERTURE DES NŒUDS BETA
================================================

**📊 Tests exécutés:** 3
**✅ Tests réussis:** 3 (100.0%)
**🧠 Score sémantique moyen:** 80.0%
**📅 Date d'exécution:** 2025-11-18 11:06:42

## 🎯 NŒUDS BETA ANALYSÉS
| Type de Nœud | Tests | Succès | Score Sémantique |
|---------------|--------|--------|------------------|
| ExistsNode | 1 | 1 | 80.0% |
| JoinNode | 1 | 1 | 60.0% |
| NotNode | 1 | 1 | 100.0% |

## 🧪 TEST 1: exists_simple
---

### 📋 Informations générales
- **Description:** Test existence simple
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/exists_simple.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/exists_simple.facts`
- **Temps d'exécution:** 690.802µs
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 80.0%
- **Actions valides:** ❌
- **Jointures valides:** ✅
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

**⚠️ Erreurs de validation:**
- Action person_has_orders: attendu 1-1 déclenchements, observé 2

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
├── 🔍 AlphaNodes
│   ├── rule_0_alpha
│   │   ├── Condition: Condition positive
│   │   └── Variable: p
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
3. **Order[fact_Order_3]** - `Order[customer_id=P001, amount=100, id=fact_Order_3]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| person_has_orders | 2 | AlphaNode | ✅ |

#### 🎯 Activation détaillée: `person_has_orders`
- **Nombre de déclenchements:** 2
- **Type de nœud déclencheur:** AlphaNode

**📋 Tokens et couples de faits activant l'action:**

##### Token 1
**Fait activateur:** `Person[id=P001, name=Alice]`

##### Token 2
**Fait activateur:** `Person[id=P002, name=Bob]`

### 🎯 Comparaison attendu vs observé
#### Actions
| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| person_has_orders | 1-1 | 2 | P001 | P001, P002 | ❌ |

#### 📋 Détails des tuples Beta attendus

**Action `person_has_orders`:**
- **Description:** Une personne a des commandes
- **Déclenchements attendus:** 1-1
- **IDs de faits attendus:**
  1. `P001`

---

## 🧪 TEST 2: join_simple
---

### 📋 Informations générales
- **Description:** Test jointure simple entre deux faits
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_simple.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/join_simple.facts`
- **Temps d'exécution:** 656.417µs
- **Résultat:** ✅ Succès

### 🧠 Validation sémantique
- **Score global:** 60.0%
- **Actions valides:** ❌
- **Jointures valides:** ❌
- **Négations valides:** ✅
- **Existences valides:** ✅
- **Agrégations valides:** ✅

**⚠️ Erreurs de validation:**
- Action attendue manquante: join_person_order
- Jointure attendue: Person -> Order, 2 correspondances

### 📜 Règles analysées
#### Règle 1
- **Texte original:** `{p: Person, o: Order} / p.id == o.customer_id ==> join_person_order(p.id, o.id)`
- **Action:** join_person_order
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
└── 🎯 TerminalNodes (Actions)
    ├── rule_0_terminal
    │   └── Action: join_person_order
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

- **Order:** 2 faits
- **Person:** 2 faits

**📋 Détail des faits parsés:**
1. **Person[P001]** - `Person[id=P001, name=Alice, age=25]`
2. **Person[P002]** - `Person[id=P002, name=Bob, age=30]`
3. **Order[O001]** - `Order[id=O001, customer_id=P001, amount=100]`
4. **Order[O002]** - `Order[amount=200, id=O002, customer_id=P002]`

### ⚡ Résultats des actions
*Aucune action déclenchée*

### 🔗 Analyse des jointures (JoinNodes)
| Nœud | Paires de Variables | Correspondances | Type | Validation |
|------|---------------------|-----------------|------|------------|
| join_0 | p <-> o | 0 | inner | ❌ |

### 🎯 Comparaison attendu vs observé
#### Actions
| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| join_person_order | 2-2 | 0 | P001, P002, O001, O002 |  | ❌ |

#### 📋 Détails des tuples Beta attendus

**Action `join_person_order`:**
- **Description:** Deux personnes ont chacune une commande
- **Déclenchements attendus:** 2-2
- **IDs de faits attendus:**
  1. `P001`
  2. `P002`
  3. `O001`
  4. `O002`

---

## 🧪 TEST 3: not_simple
---

### 📋 Informations générales
- **Description:** Test négation simple
- **Fichier contraintes:** `/home/resinsec/dev/tsd/beta_coverage_tests/not_simple.constraint`
- **Fichier faits:** `/home/resinsec/dev/tsd/beta_coverage_tests/not_simple.facts`
- **Temps d'exécution:** 884.243µs
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
1. **Person[P001]** - `Person[id=P001, name=Alice, active=true]`
2. **Person[P002]** - `Person[name=Bob, active=false, id=P002]`

### ⚡ Résultats des actions
| Action | Déclenchements | Type de Nœud | Correspondance Sémantique |
|--------|----------------|-------------|---------------------------|
| active_person | 1 | AlphaNode | ✅ |

#### 🎯 Activation détaillée: `active_person`
- **Nombre de déclenchements:** 1
- **Type de nœud déclencheur:** AlphaNode

**📋 Tokens et couples de faits activant l'action:**

##### Token 1
**Fait activateur:** `Person[id=P001, name=Alice, active=true]`

### 🚫 Analyse des négations (NotNodes)
| Nœud | Condition Niée | Faits Filtrés | Type | Validation |
|------|----------------|---------------|------|------------|
| alpha_%!d(string=rule_0_alpha) | map[condition:map[left:map[field:active object:p type:fieldAccess] operator:== right:map[type:boolean value:false] type:comparison] negated:true type:negation] | 0 | simple | ❌ |

### 🎯 Comparaison attendu vs observé
#### Actions
| Action | Attendu | Observé | Faits Attendus | Faits Observés | Statut |
|--------|---------|---------|----------------|----------------|--------|
| active_person | 1-1 | 1 | P001 | P001 | ✅ |

#### 📋 Détails des tuples Beta attendus

**Action `active_person`:**
- **Description:** Une seule personne active
- **Déclenchements attendus:** 1-1
- **IDs de faits attendus:**
  1. `P001`

---

## 💡 RECOMMANDATIONS
### Amélioration de la couverture Beta
### Prochaines étapes
1. **Ajouter plus de tests complexes** avec jointures multiples
2. **Tester les négations imbriquées** et conditions complexes
3. **Valider les performances** des nœuds Beta avec de gros volumes
4. **Enrichir la validation sémantique** avec plus de critères
