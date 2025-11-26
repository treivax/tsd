# 📄 Vérifier la Conformité de Licence (Verify License Compliance)

## Contexte

Projet TSD (Type System with Dependencies) - Moteur de règles RETE avec système de contraintes en Go.

Le projet est sous licence MIT et utilise uniquement des dépendances avec licences permissives compatibles. Il est essentiel de vérifier régulièrement que :
- Tous les fichiers .go ont les en-têtes de copyright appropriés
- Aucune dépendance incompatible n'a été ajoutée
- Toutes les dépendances tierces sont documentées
- Le code ne contient pas de copie non attribuée

## Objectif

Effectuer un audit complet de conformité de licence du projet pour garantir :
- ✅ Conformité légale à 100%
- ✅ Aucun code sous licence incompatible
- ✅ Documentation complète des dépendances
- ✅ En-têtes de copyright présents partout
- ✅ Prêt pour distribution (open-source et commerciale)

## 📄 LICENCE DU PROJET

**Licence principale:** MIT  
**Fichier:** `LICENSE`  
**Licences compatibles acceptées:** MIT, BSD-3-Clause, BSD-2-Clause, Apache-2.0, ISC  
**Licences incompatibles:** GPL, AGPL, LGPL (copyleft), code propriétaire sans licence

## Instructions

### PHASE 1 : Vérification des En-têtes de Copyright

#### 1.1 Vérifier tous les fichiers .go

**Compter les fichiers avec en-tête de copyright :**
```bash
# Compter les fichiers .go avec copyright TSD
grep -r "Copyright (c) 2025 TSD Contributors" --include="*.go" | wc -l

# Compter tous les fichiers .go (hors .git et vendor)
find . -name "*.go" -type f ! -path "./.git/*" ! -path "./vendor/*" | wc -l

# Vérifier la couverture
echo "Couverture: [nombre avec copyright] / [nombre total] fichiers"
```

**Identifier les fichiers sans en-tête :**
```bash
for file in $(find . -name "*.go" -type f ! -path "./.git/*" ! -path "./vendor/*"); do
    if ! head -1 "$file" | grep -q "Copyright\|Code generated"; then
        echo "⚠️  EN-TÊTE MANQUANT: $file"
    fi
done
```

**Critère de succès:** 100% des fichiers .go ont un en-tête approprié (copyright TSD ou "Code generated")

#### 1.2 Vérifier le format des en-têtes

**Format attendu :**
```go
// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text
```

**Vérifier un échantillon :**
```bash
# Vérifier quelques fichiers clés
head -4 rete/network.go
head -4 constraint/api.go
head -4 cmd/tsd/main.go
```

**Critère de succès:** Format cohérent et correct dans tous les fichiers

---

### PHASE 2 : Vérification des Dépendances Go

#### 2.1 Lister toutes les dépendances

```bash
# Lister les dépendances directes et indirectes
go list -m all

# Sauvegarder pour analyse
go list -m all > /tmp/dependencies_list.txt
```

#### 2.2 Vérifier les licences des dépendances

**Installer go-licenses (si nécessaire) :**
```bash
go install github.com/google/go-licenses@latest
```

**Générer le rapport de licences :**
```bash
go-licenses report github.com/treivax/tsd 2>/dev/null || echo "⚠️ go-licenses non installé"
```

#### 2.3 Vérifier manuellement les dépendances principales

**Dépendances attendues (toutes compatibles MIT) :**

| Dépendance | Version | Licence | Statut |
|------------|---------|---------|---------|
| testify | v1.8.1 | MIT | ✅ Compatible |
| go-spew | v1.1.1 | ISC | ✅ Compatible |
| go-difflib | v1.0.0 | BSD-3-Clause | ✅ Compatible |
| yaml.v3 | v3.0.1 | MIT/Apache-2.0 | ✅ Compatible |

**Vérifier qu'aucune nouvelle dépendance n'a été ajoutée :**
```bash
# Comparer avec la liste attendue
cat go.mod
```

**Critère de succès:** Toutes les dépendances sont sous licences permissives (MIT, BSD, ISC, Apache-2.0)

---

### PHASE 3 : Vérification de la Documentation

#### 3.1 Vérifier les fichiers de licence

**Fichiers obligatoires :**
```bash
# Vérifier la présence
ls -lh LICENSE THIRD_PARTY_LICENSES.md NOTICE

# Vérifier le contenu
echo "=== LICENSE ==="
head -3 LICENSE

echo "=== THIRD_PARTY_LICENSES.md ==="
grep -A 2 "## Table of Contents" THIRD_PARTY_LICENSES.md || echo "⚠️ Table des matières manquante"

echo "=== NOTICE ==="
head -5 NOTICE
```

**Critère de succès:** 
- ✅ LICENSE existe et contient le texte MIT complet
- ✅ THIRD_PARTY_LICENSES.md existe et liste toutes les dépendances
- ✅ NOTICE existe avec les attributions consolidées

#### 3.2 Vérifier la cohérence de la documentation

**Vérifier que README.md mentionne la licence :**
```bash
grep -A 5 "## 📄 License" README.md || echo "⚠️ Section License manquante dans README"
grep "THIRD_PARTY_LICENSES.md" README.md || echo "⚠️ Lien vers THIRD_PARTY_LICENSES.md manquant"
```

**Vérifier que THIRD_PARTY_LICENSES.md liste toutes les dépendances :**
```bash
# Vérifier présence de chaque dépendance principale
grep -i "testify" THIRD_PARTY_LICENSES.md || echo "⚠️ testify manquant"
grep -i "pigeon" THIRD_PARTY_LICENSES.md || echo "⚠️ pigeon manquant"
grep -i "go-spew" THIRD_PARTY_LICENSES.md || echo "⚠️ go-spew manquant"
```

**Critère de succès:** Documentation complète et cohérente

---

### PHASE 4 : Recherche de Code Non Attribué

#### 4.1 Rechercher des références à du code externe

**Rechercher des commentaires indicateurs :**
```bash
# Rechercher des références à du code copié
grep -ri "stackoverflow\|stack overflow" --include="*.go" || echo "✅ Aucune référence StackOverflow"
grep -ri "copied from\|taken from\|borrowed from" --include="*.go" || echo "✅ Aucune copie détectée"
grep -ri "source:\|adapted from\|based on:" --include="*.go" | head -10
```

**Rechercher des TODOs liés aux licences :**
```bash
grep -ri "TODO.*license\|FIXME.*license\|XXX.*license" --include="*.go" || echo "✅ Aucun TODO de licence"
```

**Critère de succès:** 
- Aucune référence à du code copié non documenté
- Toutes les sources d'inspiration sont citées correctement
- Aucun TODO de licence non résolu

#### 4.2 Vérifier le code généré

**Identifier tous les fichiers générés :**
```bash
grep -r "^// Code generated" --include="*.go"
```

**Vérifier constraint/parser.go (généré par Pigeon) :**
```bash
head -1 constraint/parser.go
# Devrait afficher: // Code generated by pigeon; DO NOT EDIT.
```

**Critère de succès:** 
- Tous les fichiers générés sont identifiés
- Pigeon (BSD-3-Clause) documenté dans THIRD_PARTY_LICENSES.md

---

### PHASE 5 : Vérification de Compatibilité Légale

#### 5.1 Vérifier l'absence de licences incompatibles

**Rechercher des mentions de licences copyleft :**
```bash
grep -ri "GPL\|AGPL\|LGPL" --include="*.go" --include="*.md" go.mod go.sum || echo "✅ Aucune licence copyleft détectée"
```

**Vérifier go.mod pour licences suspectes :**
```bash
# Chercher des packages connus pour être GPL
grep -i "gnu\|copyleft" go.mod || echo "✅ Pas de package GPL dans go.mod"
```

**Critère de succès:** Aucune dépendance GPL/AGPL/LGPL

#### 5.2 Vérifier la compatibilité avec MIT

**Tableau de compatibilité :**

| Licence Tierce | Compatible MIT | Utilisation TSD |
|----------------|----------------|-----------------|
| MIT | ✅ Oui | testify, yaml.v3 |
| BSD-2/3-Clause | ✅ Oui | Pigeon, go-difflib |
| ISC | ✅ Oui | go-spew |
| Apache-2.0 | ✅ Oui | yaml.v3 (dual) |
| GPL/AGPL/LGPL | ❌ Non | Aucune |
| Propriétaire | ❌ Non | Aucune |

**Critère de succès:** Toutes les licences tierces sont compatibles MIT

---

### PHASE 6 : Génération du Rapport de Conformité

#### 6.1 Créer le rapport

**Générer un rapport complet :**
```bash
cat > LICENSE_COMPLIANCE_REPORT_$(date +%Y%m%d).md << 'EOFR'
# Rapport de Conformité de Licence TSD

**Date:** $(date +%Y-%m-%d)
**Auditeur:** [Nom]
**Statut:** [À compléter]

## 1. En-têtes de Copyright

- Fichiers .go totaux: [X]
- Fichiers avec copyright TSD: [X]
- Fichiers avec "Code generated": [X]
- Couverture: [X]%

**Statut:** [✅ CONFORME / ⚠️ PARTIEL / ❌ NON CONFORME]

## 2. Dépendances

### Dépendances Directes
- testify v1.8.1 (MIT) ✅
- [Lister autres]

### Dépendances Indirectes
- go-spew v1.1.1 (ISC) ✅
- go-difflib v1.0.0 (BSD-3-Clause) ✅
- yaml.v3 v3.0.1 (MIT/Apache-2.0) ✅

**Statut:** [✅ TOUTES COMPATIBLES / ⚠️ À VÉRIFIER / ❌ INCOMPATIBLES]

## 3. Documentation

- LICENSE: [✅ / ❌]
- THIRD_PARTY_LICENSES.md: [✅ / ❌]
- NOTICE: [✅ / ❌]
- README.md section License: [✅ / ❌]

**Statut:** [✅ COMPLÈTE / ⚠️ PARTIELLE / ❌ MANQUANTE]

## 4. Code Non Attribué

- Références externes trouvées: [X]
- TODOs de licence: [X]
- Code copié non documenté: [X]

**Statut:** [✅ AUCUN / ⚠️ À DOCUMENTER / ❌ PROBLÈME]

## 5. Compatibilité Légale

- Licences incompatibles détectées: [AUCUNE / Lister]
- Risque GPL/AGPL: [AUCUN / Lister]

**Statut:** [✅ 100% COMPATIBLE / ❌ INCOMPATIBLE]

## 6. Conclusion Globale

**Statut de Conformité:** [✅ 100% CONFORME / ⚠️ ACTIONS REQUISES / ❌ NON CONFORME]

**Risque Légal:** [AUCUN / FAIBLE / MOYEN / ÉLEVÉ]

**Prêt pour Distribution:** [✅ OUI / ❌ NON]

### Actions Recommandées
- [ ] [Action 1]
- [ ] [Action 2]

EOFR
```

#### 6.2 Format de rapport attendu

**Le rapport doit inclure :**

1. **Résumé Exécutif**
   - Statut global de conformité
   - Risques identifiés
   - Actions requises

2. **Détails par Section**
   - En-têtes de copyright (avec statistiques)
   - Dépendances (tableau complet)
   - Documentation (checklist)
   - Code externe (liste des références)
   - Compatibilité légale (analyse)

3. **Recommandations**
   - Actions correctives prioritaires
   - Améliorations suggérées
   - Prochaine date d'audit

---

## 🎯 Critères de Succès Globaux

### ✅ CONFORMITÉ TOTALE (100%)

- [x] Tous les fichiers .go ont un en-tête approprié
- [x] Toutes les dépendances sont sous licences permissives
- [x] LICENSE, THIRD_PARTY_LICENSES.md et NOTICE présents
- [x] Documentation complète et à jour
- [x] Aucun code non attribué
- [x] Aucune licence incompatible (GPL/AGPL)
- [x] Prêt pour distribution open-source et commerciale

### ⚠️ CONFORMITÉ PARTIELLE

- Actions mineures requises
- Documentation à compléter
- Quelques en-têtes manquants

### ❌ NON CONFORMITÉ

- Dépendances incompatibles présentes
- Documentation manquante
- Nombreux fichiers sans en-têtes
- Code copié non attribué

---

## 📋 Checklist Rapide

Utiliser cette checklist pour un audit rapide :

```
VÉRIFICATION RAPIDE DE CONFORMITÉ TSD

□ LICENSE existe et contient MIT complet
□ THIRD_PARTY_LICENSES.md existe et liste toutes dépendances
□ NOTICE existe avec attributions
□ README.md a section License avec liens
□ 100% fichiers .go ont en-tête copyright ou "Code generated"
□ go.mod ne contient que dépendances permissives
□ Aucune mention GPL/AGPL dans le projet
□ constraint/parser.go a en-tête "Code generated by pigeon"
□ Pigeon documenté dans THIRD_PARTY_LICENSES.md
□ Aucune référence StackOverflow non documentée
□ Aucun TODO de licence non résolu

RÉSULTAT: ____ / 11 critères validés

Statut: 
  11/11 = ✅ CONFORME
  9-10/11 = ⚠️ PARTIEL
  <9/11 = ❌ NON CONFORME
```

---

## 🚨 Actions Correctives Courantes

### Si en-têtes manquants :

```bash
# Utiliser le script existant
bash scripts/add_copyright_headers.sh
```

### Si dépendance incompatible détectée :

1. **Identifier des alternatives** sous licence permissive
2. **Remplacer la dépendance** incompatible
3. **Mettre à jour go.mod** et THIRD_PARTY_LICENSES.md
4. **Tester** que tout fonctionne encore

### Si documentation incomplète :

1. **Mettre à jour THIRD_PARTY_LICENSES.md** avec nouvelles dépendances
2. **Ajouter NOTICE** si manquant
3. **Mettre à jour README.md** section License

---

## 📊 Fréquence d'Audit Recommandée

- **Audit complet:** Tous les 3 mois ou avant chaque release majeure
- **Vérification rapide:** Avant chaque merge de dépendance
- **Automatisation:** Intégrer dans CI/CD (optionnel)

---

## 🔗 Ressources

### Fichiers du Projet
- `LICENSE` - Licence MIT du projet
- `THIRD_PARTY_LICENSES.md` - Licences des dépendances
- `NOTICE` - Attributions consolidées
- `CODE_TIERS_IDENTIFIE.md` - Identification du code tiers
- `COPYRIGHT_HEADERS_COMPLETE.md` - Guide des en-têtes

### Outils Externes
- go-licenses: https://github.com/google/go-licenses
- SPDX License List: https://spdx.org/licenses/
- Choose a License: https://choosealicense.com/

### Références Légales
- MIT License: https://opensource.org/licenses/MIT
- BSD Licenses: https://opensource.org/licenses/BSD-3-Clause
- Apache 2.0: https://www.apache.org/licenses/LICENSE-2.0

---

## Exemple d'Utilisation

**Commande simple :**
```
Lance une vérification complète de conformité de licence
```

**Commande détaillée :**
```
Utilise le prompt "verify-license-compliance" pour auditer 
toutes les licences du projet et générer un rapport complet
```

**Vérification rapide :**
```
Fais juste une vérification rapide des en-têtes et dépendances
```

---

## 📝 Notes Importantes

- ⚠️ Ce prompt doit être exécuté **avant chaque release publique**
- ⚠️ Toute nouvelle dépendance doit déclencher un audit
- ⚠️ Les licences des dépendances peuvent changer entre versions
- ✅ La conformité est une responsabilité continue, pas ponctuelle
- ✅ En cas de doute sur une licence, **NE PAS UTILISER** la dépendance

---

**Prompt créé le:** 2025-01-XX  
**Version:** 1.0  
**Prochaine révision:** Lors de l'ajout de nouvelles dépendances majeures