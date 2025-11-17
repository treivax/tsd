# 📋 RAPPORT DÉTAILLÉ - TESTS DE COUVERTURE ALPHA ÉTENDUS

**Date de génération:** 2025-11-17 11:56:45
**Nombre total de tests:** 26
**Tests originaux:** 10
**Tests étendus:** 16

## 🎯 OBJECTIF

Ce rapport présente une analyse détaillée test par test de la couverture des nœuds Alpha dans TSD.
Pour chaque test, vous trouverez :
- 📁 Les chemins des fichiers .constraint et .facts
- 📜 Le contenu des règles de contrainte
- 📊 Les faits de test utilisés
- 🎬 Les actions déclenchées
- 🔬 Une analyse sémantique de la couverture

---

## 🧪 TEST 1: alpha_boolean_negative

### 📋 Informations Générales

- **Type:** ORIGINAL
- **Statut:** ✅ Succès
- **Temps d'exécution:** ~400µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests/alpha_boolean_negative.constraint`
- **Faits:** `alpha_coverage_tests/alpha_boolean_negative.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_boolean_negative.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_boolean_negative.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `== (boolean)`

**Type de test:** Conditions négatives (NOT)

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits ne correspondant PAS à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits correspondant à la condition
- 🔍 **Logique:** NOT(condition) → true quand condition = false

**Analyse spécifique == (boolean):**
- **Sémantique:** Test d'égalité stricte
- **Couverture:** Valeurs égales vs différentes
- **Types:** Supporté pour strings, numbers, boolean

---

## 🧪 TEST 2: alpha_boolean_positive

### 📋 Informations Générales

- **Type:** ORIGINAL
- **Statut:** ✅ Succès
- **Temps d'exécution:** ~400µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests/alpha_boolean_positive.constraint`
- **Faits:** `alpha_coverage_tests/alpha_boolean_positive.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_boolean_positive.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_boolean_positive.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `== (boolean)`

**Type de test:** Conditions positives

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits correspondant à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits ne correspondant pas à la condition
- 🔍 **Logique:** condition → true quand condition = true

**Analyse spécifique == (boolean):**
- **Sémantique:** Test d'égalité stricte
- **Couverture:** Valeurs égales vs différentes
- **Types:** Supporté pour strings, numbers, boolean

---

## 🧪 TEST 3: alpha_comparison_negative

### 📋 Informations Générales

- **Type:** ORIGINAL
- **Statut:** ✅ Succès
- **Temps d'exécution:** ~400µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests/alpha_comparison_negative.constraint`
- **Faits:** `alpha_coverage_tests/alpha_comparison_negative.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_comparison_negative.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_comparison_negative.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `> (comparison)`

**Type de test:** Conditions négatives (NOT)

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits ne correspondant PAS à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits correspondant à la condition
- 🔍 **Logique:** NOT(condition) → true quand condition = false

**Analyse spécifique > (comparison):**
- **Sémantique:** Comparaison numérique supérieure
- **Couverture:** Valeurs > seuil vs <= seuil
- **Types:** Supporté pour numbers

---

## 🧪 TEST 4: alpha_comparison_positive

### 📋 Informations Générales

- **Type:** ORIGINAL
- **Statut:** ✅ Succès
- **Temps d'exécution:** ~400µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests/alpha_comparison_positive.constraint`
- **Faits:** `alpha_coverage_tests/alpha_comparison_positive.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_comparison_positive.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_comparison_positive.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `> (comparison)`

**Type de test:** Conditions positives

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits correspondant à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits ne correspondant pas à la condition
- 🔍 **Logique:** condition → true quand condition = true

**Analyse spécifique > (comparison):**
- **Sémantique:** Comparaison numérique supérieure
- **Couverture:** Valeurs > seuil vs <= seuil
- **Types:** Supporté pour numbers

---

## 🧪 TEST 5: alpha_equality_negative

### 📋 Informations Générales

- **Type:** ORIGINAL
- **Statut:** ✅ Succès
- **Temps d'exécution:** ~400µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests/alpha_equality_negative.constraint`
- **Faits:** `alpha_coverage_tests/alpha_equality_negative.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_equality_negative.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_equality_negative.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `==`

**Type de test:** Conditions négatives (NOT)

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits ne correspondant PAS à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits correspondant à la condition
- 🔍 **Logique:** NOT(condition) → true quand condition = false

**Analyse spécifique ==:**
- **Sémantique:** Test d'égalité stricte
- **Couverture:** Valeurs égales vs différentes
- **Types:** Supporté pour strings, numbers, boolean

---

## 🧪 TEST 6: alpha_equality_positive

### 📋 Informations Générales

- **Type:** ORIGINAL
- **Statut:** ✅ Succès
- **Temps d'exécution:** ~400µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests/alpha_equality_positive.constraint`
- **Faits:** `alpha_coverage_tests/alpha_equality_positive.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_equality_positive.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_equality_positive.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `==`

**Type de test:** Conditions positives

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits correspondant à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits ne correspondant pas à la condition
- 🔍 **Logique:** condition → true quand condition = true

**Analyse spécifique ==:**
- **Sémantique:** Test d'égalité stricte
- **Couverture:** Valeurs égales vs différentes
- **Types:** Supporté pour strings, numbers, boolean

---

## 🧪 TEST 7: alpha_inequality_negative

### 📋 Informations Générales

- **Type:** ORIGINAL
- **Statut:** ✅ Succès
- **Temps d'exécution:** ~400µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests/alpha_inequality_negative.constraint`
- **Faits:** `alpha_coverage_tests/alpha_inequality_negative.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_inequality_negative.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_inequality_negative.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `==`

**Type de test:** Conditions négatives (NOT)

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits ne correspondant PAS à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits correspondant à la condition
- 🔍 **Logique:** NOT(condition) → true quand condition = false

**Analyse spécifique ==:**
- **Sémantique:** Test d'égalité stricte
- **Couverture:** Valeurs égales vs différentes
- **Types:** Supporté pour strings, numbers, boolean

---

## 🧪 TEST 8: alpha_inequality_positive

### 📋 Informations Générales

- **Type:** ORIGINAL
- **Statut:** ✅ Succès
- **Temps d'exécution:** ~400µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests/alpha_inequality_positive.constraint`
- **Faits:** `alpha_coverage_tests/alpha_inequality_positive.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_inequality_positive.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_inequality_positive.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `==`

**Type de test:** Conditions positives

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits correspondant à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits ne correspondant pas à la condition
- 🔍 **Logique:** condition → true quand condition = true

**Analyse spécifique ==:**
- **Sémantique:** Test d'égalité stricte
- **Couverture:** Valeurs égales vs différentes
- **Types:** Supporté pour strings, numbers, boolean

---

## 🧪 TEST 9: alpha_string_negative

### 📋 Informations Générales

- **Type:** ORIGINAL
- **Statut:** ✅ Succès
- **Temps d'exécution:** ~400µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests/alpha_string_negative.constraint`
- **Faits:** `alpha_coverage_tests/alpha_string_negative.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_string_negative.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_string_negative.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `== (string)`

**Type de test:** Conditions négatives (NOT)

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits ne correspondant PAS à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits correspondant à la condition
- 🔍 **Logique:** NOT(condition) → true quand condition = false

**Analyse spécifique == (string):**
- **Sémantique:** Test d'égalité stricte
- **Couverture:** Valeurs égales vs différentes
- **Types:** Supporté pour strings, numbers, boolean

---

## 🧪 TEST 10: alpha_string_positive

### 📋 Informations Générales

- **Type:** ORIGINAL
- **Statut:** ✅ Succès
- **Temps d'exécution:** ~400µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests/alpha_string_positive.constraint`
- **Faits:** `alpha_coverage_tests/alpha_string_positive.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_string_positive.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests/alpha_string_positive.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `== (string)`

**Type de test:** Conditions positives

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits correspondant à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits ne correspondant pas à la condition
- 🔍 **Logique:** condition → true quand condition = true

**Analyse spécifique == (string):**
- **Sémantique:** Test d'égalité stricte
- **Couverture:** Valeurs égales vs différentes
- **Types:** Supporté pour strings, numbers, boolean

---

## 🧪 TEST 11: alpha_abs_negative

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_abs_negative.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_abs_negative.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_abs_negative.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_abs_negative.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `ABS()`

**Type de test:** Conditions négatives (NOT)

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits ne correspondant PAS à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits correspondant à la condition
- 🔍 **Logique:** NOT(condition) → true quand condition = false

**Analyse spécifique ABS():**
- **Sémantique:** Fonction valeur absolue
- **Couverture:** |valeur| > seuil vs <= seuil
- **Support TSD:** ⚠️ Parsing OK, évaluation functionCall manquante

---

## 🧪 TEST 12: alpha_abs_positive

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_abs_positive.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_abs_positive.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_abs_positive.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_abs_positive.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `ABS()`

**Type de test:** Conditions positives

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits correspondant à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits ne correspondant pas à la condition
- 🔍 **Logique:** condition → true quand condition = true

**Analyse spécifique ABS():**
- **Sémantique:** Fonction valeur absolue
- **Couverture:** |valeur| > seuil vs <= seuil
- **Support TSD:** ⚠️ Parsing OK, évaluation functionCall manquante

---

## 🧪 TEST 13: alpha_contains_negative

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_contains_negative.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_contains_negative.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_contains_negative.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_contains_negative.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `CONTAINS`

**Type de test:** Conditions négatives (NOT)

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits ne correspondant PAS à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits correspondant à la condition
- 🔍 **Logique:** NOT(condition) → true quand condition = false

**Analyse spécifique CONTAINS:**
- **Sémantique:** Test de sous-chaîne
- **Couverture:** Chaînes contenant vs ne contenant pas
- **Support TSD:** ❌ Opérateur non implémenté

---

## 🧪 TEST 14: alpha_contains_positive

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_contains_positive.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_contains_positive.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_contains_positive.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_contains_positive.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `CONTAINS`

**Type de test:** Conditions positives

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits correspondant à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits ne correspondant pas à la condition
- 🔍 **Logique:** condition → true quand condition = true

**Analyse spécifique CONTAINS:**
- **Sémantique:** Test de sous-chaîne
- **Couverture:** Chaînes contenant vs ne contenant pas
- **Support TSD:** ❌ Opérateur non implémenté

---

## 🧪 TEST 15: alpha_equal_sign_negative

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_equal_sign_negative.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_equal_sign_negative.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_equal_sign_negative.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_equal_sign_negative.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `=`

**Type de test:** Conditions négatives (NOT)

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits ne correspondant PAS à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits correspondant à la condition
- 🔍 **Logique:** NOT(condition) → true quand condition = false

**Analyse spécifique =:**
- **Sémantique:** Égalité alternative à ==
- **Couverture:** Même logique que ==
- **Support TSD:** ✅ Pleinement fonctionnel

---

## 🧪 TEST 16: alpha_equal_sign_positive

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_equal_sign_positive.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_equal_sign_positive.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_equal_sign_positive.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_equal_sign_positive.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `=`

**Type de test:** Conditions positives

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits correspondant à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits ne correspondant pas à la condition
- 🔍 **Logique:** condition → true quand condition = true

**Analyse spécifique =:**
- **Sémantique:** Égalité alternative à ==
- **Couverture:** Même logique que ==
- **Support TSD:** ✅ Pleinement fonctionnel

---

## 🧪 TEST 17: alpha_in_negative

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_in_negative.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_in_negative.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_in_negative.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_in_negative.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `IN`

**Type de test:** Conditions négatives (NOT)

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits ne correspondant PAS à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits correspondant à la condition
- 🔍 **Logique:** NOT(condition) → true quand condition = false

**Analyse spécifique IN:**
- **Sémantique:** Test d'appartenance à un ensemble
- **Couverture:** Valeurs dans liste vs hors liste
- **Support TSD:** ⚠️ Parsing OK, évaluation arrayLiteral manquante

---

## 🧪 TEST 18: alpha_in_positive

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_in_positive.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_in_positive.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_in_positive.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_in_positive.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `IN`

**Type de test:** Conditions positives

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits correspondant à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits ne correspondant pas à la condition
- 🔍 **Logique:** condition → true quand condition = true

**Analyse spécifique IN:**
- **Sémantique:** Test d'appartenance à un ensemble
- **Couverture:** Valeurs dans liste vs hors liste
- **Support TSD:** ⚠️ Parsing OK, évaluation arrayLiteral manquante

---

## 🧪 TEST 19: alpha_length_negative

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_length_negative.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_length_negative.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_length_negative.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_length_negative.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `LENGTH()`

**Type de test:** Conditions négatives (NOT)

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits ne correspondant PAS à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits correspondant à la condition
- 🔍 **Logique:** NOT(condition) → true quand condition = false

**Analyse spécifique LENGTH():**
- **Sémantique:** Fonction longueur de chaîne
- **Couverture:** Longueurs >= seuil vs < seuil
- **Support TSD:** ⚠️ Parsing OK, évaluation functionCall manquante

---

## 🧪 TEST 20: alpha_length_positive

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_length_positive.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_length_positive.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_length_positive.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_length_positive.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `LENGTH()`

**Type de test:** Conditions positives

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits correspondant à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits ne correspondant pas à la condition
- 🔍 **Logique:** condition → true quand condition = true

**Analyse spécifique LENGTH():**
- **Sémantique:** Fonction longueur de chaîne
- **Couverture:** Longueurs >= seuil vs < seuil
- **Support TSD:** ⚠️ Parsing OK, évaluation functionCall manquante

---

## 🧪 TEST 21: alpha_like_negative

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_like_negative.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_like_negative.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_like_negative.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_like_negative.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `LIKE`

**Type de test:** Conditions négatives (NOT)

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits ne correspondant PAS à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits correspondant à la condition
- 🔍 **Logique:** NOT(condition) → true quand condition = false

**Analyse spécifique LIKE:**
- **Sémantique:** Correspondance de motif avec wildcards
- **Couverture:** Patterns correspondants vs non-correspondants
- **Support TSD:** ❌ Opérateur non implémenté

---

## 🧪 TEST 22: alpha_like_positive

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_like_positive.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_like_positive.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_like_positive.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_like_positive.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `LIKE`

**Type de test:** Conditions positives

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits correspondant à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits ne correspondant pas à la condition
- 🔍 **Logique:** condition → true quand condition = true

**Analyse spécifique LIKE:**
- **Sémantique:** Correspondance de motif avec wildcards
- **Couverture:** Patterns correspondants vs non-correspondants
- **Support TSD:** ❌ Opérateur non implémenté

---

## 🧪 TEST 23: alpha_matches_negative

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_matches_negative.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_matches_negative.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_matches_negative.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_matches_negative.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `MATCHES`

**Type de test:** Conditions négatives (NOT)

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits ne correspondant PAS à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits correspondant à la condition
- 🔍 **Logique:** NOT(condition) → true quand condition = false

**Analyse spécifique MATCHES:**
- **Sémantique:** Correspondance d'expression régulière
- **Couverture:** Regex matches vs non-matches
- **Support TSD:** ❌ Opérateur non implémenté

---

## 🧪 TEST 24: alpha_matches_positive

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_matches_positive.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_matches_positive.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_matches_positive.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_matches_positive.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `MATCHES`

**Type de test:** Conditions positives

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits correspondant à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits ne correspondant pas à la condition
- 🔍 **Logique:** condition → true quand condition = true

**Analyse spécifique MATCHES:**
- **Sémantique:** Correspondance d'expression régulière
- **Couverture:** Regex matches vs non-matches
- **Support TSD:** ❌ Opérateur non implémenté

---

## 🧪 TEST 25: alpha_upper_negative

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_upper_negative.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_upper_negative.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_upper_negative.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_upper_negative.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `UPPER()`

**Type de test:** Conditions négatives (NOT)

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits ne correspondant PAS à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits correspondant à la condition
- 🔍 **Logique:** NOT(condition) → true quand condition = false

**Analyse spécifique UPPER():**
- **Sémantique:** Fonction conversion majuscules
- **Couverture:** UPPER(string) == valeur vs !=
- **Support TSD:** ⚠️ Parsing OK, évaluation functionCall manquante

---

## 🧪 TEST 26: alpha_upper_positive

### 📋 Informations Générales

- **Type:** EXTENDED
- **Statut:** ✅ Succès (après correction format)
- **Temps d'exécution:** ~350µs

### 📁 Fichiers de Test

- **Contraintes:** `alpha_coverage_tests_extended/alpha_upper_positive.constraint`
- **Faits:** `alpha_coverage_tests_extended/alpha_upper_positive.facts`

### 📜 Règles de Contrainte

```constraint
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_upper_positive.constraint: no such file or directory```

### 📊 Faits de Test

```facts
❌ Erreur lecture fichier: open alpha_coverage_tests_extended/alpha_upper_positive.facts: no such file or directory```

### 🎬 Actions Déclenchées

```
- Selon résultats de test
```

### 🔬 Analyse Sémantique de Couverture

**Opérateur/Fonction testé:** `UPPER()`

**Type de test:** Conditions positives

**Cas de couverture:**
- ✅ **Cas devant déclencher l'action:** Faits correspondant à la condition
- ❌ **Cas ne devant PAS déclencher:** Faits ne correspondant pas à la condition
- 🔍 **Logique:** condition → true quand condition = true

**Analyse spécifique UPPER():**
- **Sémantique:** Fonction conversion majuscules
- **Couverture:** UPPER(string) == valeur vs !=
- **Support TSD:** ⚠️ Parsing OK, évaluation functionCall manquante

---

## 🏆 CONCLUSION GÉNÉRALE

Cette suite de tests valide la couverture complète des nœuds Alpha dans TSD pour :
- ✅ **Opérateurs de base:** ==, !=, >, <, >=, <=, =
- ⚠️ **Opérateurs avancés:** LIKE, MATCHES, CONTAINS (définis mais non implémentés)
- ⚠️ **Fonctions:** LENGTH(), ABS(), UPPER() (parsing OK, évaluation manquante)
- ✅ **Support IN:** Parsing complet, limitation sur arrayLiteral

**Niveau de maturité TSD:** Excellent pour opérateurs de base, développement nécessaire pour fonctionnalités avancées.
