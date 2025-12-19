# RAPPORT DÉTAILLÉ DES TESTS E2E - XUPLES & XUPLE-SPACES

## Informations générales

- **Date d'exécution**: 2025-12-18 11:59:45
- **Version Go**: go1.24.4
- **Système**: Linux x86_64
- **Projet**: tsd - Tuple Space Distribution

## 1. Programme TSD de test

Le test E2E utilise un programme TSD complet avec des capteurs de température et d'humidité.

### 1.1 Types définis

Le programme définit **3 types** :

```tsd
type Sensor(sensorId: string, location: string, temperature: number, humidity: number)
type Alert(level: string, message: string, sensorId: string)
type Command(action: string, target: string, priority: number)
```

**Description des types** :
- **Sensor** : Représente un capteur avec son ID, sa localisation, sa température et son humidité
- **Alert** : Représente une alerte avec un niveau de sévérité, un message et l'ID du capteur concerné
- **Command** : Représente une commande d'action avec une cible et une priorité

### 1.2 Xuple-spaces déclarés

Le programme déclare **3 xuple-spaces** avec des politiques différentes :

#### a) critical_alerts
```tsd
xuple-space critical_alerts {
  selection: lifo
  consumption: per-agent
  retention: duration(10m)
}
```
- **Sélection LIFO** : Last In First Out - le dernier xuple inséré est récupéré en premier
- **Consommation per-agent** : Chaque agent peut consommer le même xuple
- **Rétention 10 minutes** : Les xuples expirent après 10 minutes

**Utilisation** : Pour les alertes critiques qui doivent être traitées par plusieurs agents de monitoring

#### b) normal_alerts
```tsd
xuple-space normal_alerts {
  selection: random
  consumption: once
  retention: duration(30m)
}
```
- **Sélection aléatoire** : Les xuples sont récupérés dans un ordre aléatoire
- **Consommation once** : Chaque xuple ne peut être consommé qu'une seule fois
- **Rétention 30 minutes** : Les xuples expirent après 30 minutes

**Utilisation** : Pour les alertes normales qui doivent être traitées par un seul agent

#### c) command_queue
```tsd
xuple-space command_queue {
  selection: fifo
  consumption: once
  retention: duration(1h)
}
```
- **Sélection FIFO** : First In First Out - le premier xuple inséré est récupéré en premier
- **Consommation once** : Chaque xuple ne peut être consommé qu'une seule fois
- **Rétention 1 heure** : Les xuples expirent après 1 heure

**Utilisation** : Pour une file de commandes traitées dans l'ordre d'arrivée

### 1.3 Règles définies

Le programme définit **3 règles** de détection d'anomalies :

#### Règle 1 : critical_temperature
```tsd
rule critical_temperature: {s: Sensor} / s.temperature > 40
  ==> notifyCritical(s.sensorId, s.temperature)
```
- **Condition** : Température supérieure à 40°C
- **Action** : Notification critique avec l'ID du capteur et la température
- **Déclenchement** : Alerte de niveau CRITICAL

#### Règle 2 : high_temperature
```tsd
rule high_temperature: {s: Sensor} / s.temperature > 30 AND s.temperature <= 40
  ==> notifyHigh(s.sensorId, s.temperature)
```
- **Condition** : Température entre 30°C et 40°C (exclusif)
- **Action** : Notification de température élevée
- **Déclenchement** : Alerte de niveau WARNING

#### Règle 3 : high_humidity
```tsd
rule high_humidity: {s: Sensor} / s.humidity > 80
  ==> ventilate(s.location)
```
- **Condition** : Humidité supérieure à 80%
- **Action** : Activation de la ventilation pour la localisation
- **Déclenchement** : Commande de ventilation

### 1.4 Faits insérés

Le test insère **5 faits** de type Sensor :

| # | Sensor ID | Location | Température | Humidité | Règles déclenchées |
|---|-----------|----------|-------------|----------|-------------------|
| 1 | S001 | RoomA | 22.0°C | 45.0% | *(aucune)* |
| 2 | S002 | RoomB | 35.0°C | 50.0% | **high_temperature** |
| 3 | S003 | RoomC | 45.0°C | 60.0% | **critical_temperature** |
| 4 | S004 | RoomD | 25.0°C | 85.0% | **high_humidity** |
| 5 | S005 | ServerRoom | 42.0°C | 85.0% | **critical_temperature**, **high_humidity** |

**Détail des faits** :

#### Fait 1 - Capteur S001 (RoomA)
```tsd
Sensor(sensorId: "S001", location: "RoomA", temperature: 22.0, humidity: 45.0)
```
✅ Valeurs normales - Aucune règle déclenchée

#### Fait 2 - Capteur S002 (RoomB)
```tsd
Sensor(sensorId: "S002", location: "RoomB", temperature: 35.0, humidity: 50.0)
```
⚠️ Température élevée (35°C)
- **Règle déclenchée** : `high_temperature`
- **Action exécutée** : `notifyHigh("S002", 35)`

#### Fait 3 - Capteur S003 (RoomC)
```tsd
Sensor(sensorId: "S003", location: "RoomC", temperature: 45.0, humidity: 60.0)
```
🔴 Température critique (45°C)
- **Règle déclenchée** : `critical_temperature`
- **Action exécutée** : `notifyCritical("S003", 45)`

#### Fait 4 - Capteur S004 (RoomD)
```tsd
Sensor(sensorId: "S004", location: "RoomD", temperature: 25.0, humidity: 85.0)
```
💧 Humidité élevée (85%)
- **Règle déclenchée** : `high_humidity`
- **Action exécutée** : `ventilate("RoomD")`

#### Fait 5 - Capteur S005 (ServerRoom)
```tsd
Sensor(sensorId: "S005", location: "ServerRoom", temperature: 42.0, humidity: 85.0)
```
🔴💧 **Double anomalie** : Température critique (42°C) ET humidité élevée (85%)
- **Règles déclenchées** : `critical_temperature`, `high_humidity`
- **Actions exécutées** :
  - `notifyCritical("S005", 42)`
  - `ventilate("ServerRoom")`

**Résumé des déclenchements** :
- **critical_temperature** : 2 fois (S003, S005)
- **high_temperature** : 1 fois (S002)
- **high_humidity** : 2 fois (S004, S005)
- **Total** : 5 actions déclenchées

---

## 2. Xuples générés lors du test

Le test crée **6 xuples manuellement** via l'API pour valider le fonctionnement des xuple-spaces :

### 2.1 Xuples dans critical_alerts (2 xuples)

#### Xuple 1 - Alerte critique S003
```json
{
  "type": "Alert",
  "data": {
    "level": "CRITICAL",
    "message": "Temperature too high in RoomC",
    "sensorId": "S003"
  }
}
```
- **Space** : critical_alerts
- **Politique de sélection** : LIFO (dernier en premier)
- **Politique de consommation** : per-agent (réutilisable)
- **Statut** : ✅ Available

#### Xuple 2 - Alerte critique S005
```json
{
  "type": "Alert",
  "data": {
    "level": "CRITICAL",
    "message": "Critical conditions in ServerRoom",
    "sensorId": "S005"
  }
}
```
- **Space** : critical_alerts
- **Politique de sélection** : LIFO (dernier en premier)
- **Politique de consommation** : per-agent (réutilisable)
- **Statut** : ✅ Available

**Test de sélection LIFO** : Le xuple 2 (S005) est récupéré en premier car il a été créé en dernier.

**Test per-agent** : L'agent "agent-1" et "agent-2" peuvent tous deux récupérer le même xuple.

### 2.2 Xuples dans normal_alerts (1 xuple)

#### Xuple 3 - Alerte warning S002
```json
{
  "type": "Alert",
  "data": {
    "level": "WARNING",
    "message": "Temperature slightly high",
    "sensorId": "S002"
  }
}
```
- **Space** : normal_alerts
- **Politique de sélection** : random (aléatoire)
- **Politique de consommation** : once (une seule fois)
- **Statut** : ✅ Available

### 2.3 Xuples dans command_queue (3 xuples)

#### Xuple 4 - Commande ventilate RoomD
```json
{
  "type": "Command",
  "data": {
    "action": "ventilate",
    "target": "RoomD",
    "priority": 5
  }
}
```
- **Space** : command_queue
- **Politique de sélection** : FIFO (premier en premier)
- **Politique de consommation** : once (consommé après récupération)
- **Statut** : 🔵 Consumed by agent-1

#### Xuple 5 - Commande ventilate ServerRoom
```json
{
  "type": "Command",
  "data": {
    "action": "ventilate",
    "target": "ServerRoom",
    "priority": 5
  }
}
```
- **Space** : command_queue
- **Politique de sélection** : FIFO (premier en premier)
- **Politique de consommation** : once (consommé après récupération)
- **Statut** : 🔵 Consumed by agent-1

#### Xuple 6 - Commande emergency ServerRoom
```json
{
  "type": "Command",
  "data": {
    "action": "emergency",
    "target": "ServerRoom",
    "priority": 10
  }
}
```
- **Space** : command_queue
- **Politique de sélection** : FIFO (premier en premier)
- **Politique de consommation** : once (consommé après récupération)
- **Statut** : ✅ Available (pas encore consommé)

**Test de sélection FIFO** : Les commandes sont récupérées dans l'ordre de création (4, puis 5).

**Test once** : Une fois consommées, les commandes 4 et 5 ne peuvent plus être récupérées.

---

## 3. Résultats de l'exécution des tests

### 3.1 Exécution du test E2E principal

Lancement du test...
```bash
go test -v -timeout 5m ./tests/e2e -run TestXuplesE2E_RealWorld
```

**Résultat** : ✅ **SUCCÈS**

### 3.2 Statistiques extraites

- **Types parsés** : 3
- **Xuple-spaces parsés** : 3
- **Règles parsées** : 3
- **Faits insérés** : 5

### 3.3 Xuple-spaces créés et vérifiés

xuples_e2e_test.go:627: 📦 Xuple-space: normal_alerts

### 3.4 Détail de l'exécution par étape

#### Étape 1 : Parsing du programme
    xuples_e2e_test.go:99: ✅ Parsing réussi
    xuples_e2e_test.go:111: ✅ 3 xuple-spaces détectés:

#### Étape 2 : Création du réseau RETE et XupleManager
    xuples_e2e_test.go:135: ✅ Réseau RETE et XupleManager créés

#### Étape 3 : Ingestion du programme
