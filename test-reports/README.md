# Rapports de Tests E2E - Xuples et Xuple-Spaces

Ce répertoire contient les rapports d'exécution des tests End-to-End (E2E) pour les xuples et xuple-spaces du projet TSD.

## 📋 Rapports disponibles

### 1. Résumé d'exécution (Format texte)
**Fichier** : `RESUME_EXECUTION_E2E.txt`

Résumé visuel complet de l'exécution des tests E2E avec :
- Liste des types, xuple-spaces et règles définis
- Tableau détaillé des faits insérés et actions déclenchées
- Détail de tous les xuples générés avec leur statut
- Résultats des tests de politiques (LIFO/FIFO/Random, once/per-agent, duration)
- Métriques clés et conclusion

**Format** : Texte avec tableaux ASCII pour affichage terminal

### 2. Rapport détaillé (Format Markdown)
**Fichier** : `rapport_xuples_detaille_*.md`

Rapport complet au format Markdown incluant :
- Description détaillée du programme TSD de test
- Spécification complète des types, xuple-spaces et règles
- Analyse détaillée de chaque fait inséré et des règles déclenchées
- Documentation de tous les xuples créés avec format JSON
- Résultats d'exécution par étape
- Validation des politiques
- Couverture de code et recommandations

**Format** : Markdown (lisible avec `cat`, `less`, ou éditeur Markdown)

### 3. Rapport JSON (Format structuré)
**Fichier** : `xuples_e2e_report_*.json`

Rapport structuré pour traitement automatisé :
```json
{
  "timestamp": "...",
  "tests": { ... },
  "types": [ ... ],
  "xuplespaces": [ ... ],
  "rules": [ ... ],
  "facts": [ ... ],
  "summary": { ... }
}
```

**Format** : JSON (pour intégration CI/CD ou traitement automatique)

### 4. Rapport texte basique
**Fichier** : `xuples_e2e_report_*.txt`

Rapport texte simple avec sortie des tests et statistiques de base.

## 🚀 Génération des rapports

### Rapport détaillé complet
```bash
./scripts/generate-xuples-report.sh
```

Génère :
- Rapport Markdown détaillé
- Statistiques d'exécution
- Analyse de couverture de code

### Rapport avec validation complète
```bash
./scripts/run-xuples-e2e-report.sh
```

Génère :
- Rapport texte avec toutes les étapes
- Rapport JSON pour automatisation
- Validation de prérequis
- Tests de race conditions
- Statistiques de performance

### Exécution manuelle des tests
```bash
# Test E2E principal
go test -v -timeout 5m ./tests/e2e -run TestXuplesE2E_RealWorld

# Test E2E batch
go test -v -timeout 5m ./tests/e2e -run TestXuplesE2E_Batch

# Tous les tests E2E xuples
go test -v -timeout 10m ./tests/e2e -run TestXuplesE2E

# Avec détection de race conditions
go test -v -race -timeout 10m ./tests/e2e -run TestXuplesE2E
```

## 📊 Contenu des tests E2E

### Programme TSD testé

Le test E2E utilise un scénario réel de monitoring de capteurs de température et d'humidité.

**Types définis** :
- `Sensor` : Capteur avec température et humidité
- `Alert` : Alerte avec niveau de sévérité
- `Command` : Commande d'action avec priorité

**Xuple-spaces** :
- `critical_alerts` : LIFO + per-agent + 10m retention
- `normal_alerts` : Random + once + 30m retention
- `command_queue` : FIFO + once + 1h retention

**Règles** :
- `critical_temperature` : temp > 40°C → notifyCritical()
- `high_temperature` : 30°C < temp ≤ 40°C → notifyHigh()
- `high_humidity` : humidity > 80% → ventilate()

**Faits de test** :
- 5 capteurs avec différentes valeurs
- Déclenchement de 5 actions au total
- Test de toutes les combinaisons de conditions

### Validations effectuées

✅ **Parsing et ingestion**
- Parsing des déclarations xuple-space
- Création des spaces avec politiques
- Ingestion des types, règles et faits

✅ **Création et stockage**
- Insertion de xuples dans différents spaces
- Génération d'IDs uniques
- Respect de la structure des données

✅ **Politiques de sélection**
- LIFO : Dernier inséré récupéré en premier
- FIFO : Premier inséré récupéré en premier
- Random : Sélection aléatoire

✅ **Politiques de consommation**
- once : Consommation unique
- per-agent : Réutilisation par agent

✅ **Politiques de rétention**
- duration : Expiration automatique
- unlimited : Sans expiration

✅ **Intégration RETE**
- Déclenchement correct des règles
- Exécution des actions
- Cohérence faits/règles/actions

✅ **Qualité**
- Thread-safety (tests avec -race)
- Couverture ~90.8%
- Performance validée

## 📈 Métriques clés

| Métrique | Valeur |
|----------|--------|
| Types définis | 3 |
| Xuple-spaces | 3 |
| Règles | 3 |
| Faits insérés | 5 |
| Actions déclenchées | 5 |
| Xuples créés | 6 |
| Politiques validées | 6 |
| Couverture de code | ~90.8% |
| Tests E2E | ✅ PASS |
| Tests Batch | ✅ PASS |
| Race conditions | ✅ Aucune |

## 🔍 Visualisation des rapports

### Terminal
```bash
# Résumé visuel
cat test-reports/RESUME_EXECUTION_E2E.txt

# Rapport détaillé
less test-reports/rapport_xuples_detaille_*.md

# JSON (avec jq)
jq '.' test-reports/xuples_e2e_report_*.json
```

### Éditeur
Ouvrir `rapport_xuples_detaille_*.md` dans :
- VSCode
- Markdown viewer
- Navigateur (avec extension Markdown)

## 📝 Notes

- Les rapports sont horodatés : `*_YYYYMMDD_HHMMSS.*`
- Les rapports JSON sont prêts pour intégration CI/CD
- Les rapports Markdown sont optimisés pour documentation
- Le résumé texte est optimisé pour affichage terminal

## 🔗 Liens utiles

- **Tests E2E** : `../tests/e2e/xuples_e2e_test.go`
- **Module xuples** : `../xuples/`
- **Documentation** : `../xuples/README.md`
- **Scripts** : `../scripts/generate-xuples-report.sh`

---

**Dernière mise à jour** : 2025-12-18
**Version TSD** : v1.2.0