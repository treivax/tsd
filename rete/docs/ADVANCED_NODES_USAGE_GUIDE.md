# Guide d'Utilisation des Nœuds RETE Avancés

## 🎯 Vue d'ensemble

Les nœuds RETE avancés permettent de créer des règles métier sophistiquées avec négation, quantification existentielle et agrégation. Ce guide présente leur utilisation pratique.

## 📋 Nœuds Disponibles

### 1. NotNode (Négation)
**Usage** : Détecter l'absence de faits satisfaisant une condition

```go
// Créer un nœud NOT
notNode := nodes.NewNotNode("not_recent_login", logger)
notNode.SetNegationCondition("type == 'login' AND timestamp > recent")

// Usage : Détecter les comptes sans connexion récente
```

### 2. ExistsNode (Quantification Existentielle)
**Usage** : Vérifier l'existence d'au moins un fait satisfaisant une condition

```go
// Créer un nœud EXISTS  
existsNode := nodes.NewExistsNode("exists_suspicious", logger)
variable := domain.TypedVariable{
    Name:     "suspicious_activity", 
    DataType: "SecurityEvent"
}
existsNode.SetExistenceCondition(variable, "risk_level == 'high'")

// Usage : Détecter la présence d'activités suspectes
```

### 3. AccumulateNode (Agrégation)
**Usage** : Calculer des agrégations sur des collections de faits

```go
// Créer un nœud d'accumulation
accumulator := domain.AccumulateFunction{
    FunctionType: "SUM",
    Field:        "amount",
}
accNode := nodes.NewAccumulateNode("daily_sum", accumulator, logger)

// Usage : Calculer la somme des transactions quotidiennes
```

## 🏦 Cas d'Usage : Détection de Fraude Bancaire

### Scénario Complet

```go
func SetupFraudDetection(logger domain.Logger) (*FraudDetectionSystem, error) {
    // 1. Nœud NOT : Pas de transaction légitime récente
    notNode := nodes.NewNotNode("no_recent_legitimate", logger)
    notNode.SetNegationCondition("type == 'legitimate' AND age_hours < 24")
    
    // 2. Nœud EXISTS : Transactions suspectes présentes
    existsNode := nodes.NewExistsNode("has_suspicious", logger)
    suspiciousVar := domain.TypedVariable{
        Name:     "suspicious_tx",
        DataType: "Transaction",
    }
    existsNode.SetExistenceCondition(suspiciousVar, "amount > 10000 AND location != 'home'")
    
    // 3. Nœud ACCUMULATE : Somme des montants
    sumAccumulator := domain.AccumulateFunction{
        FunctionType: "SUM",
        Field:        "amount",
    }
    sumNode := nodes.NewAccumulateNode("total_amount", sumAccumulator, logger)
    
    // 4. Nœud ACCUMULATE : Nombre de transactions
    countAccumulator := domain.AccumulateFunction{
        FunctionType: "COUNT",
        Field:        "",
    }
    countNode := nodes.NewAccumulateNode("tx_count", countAccumulator, logger)
    
    return &FraudDetectionSystem{
        NotNode:   notNode,
        ExistsNode: existsNode,
        SumNode:   sumNode,
        CountNode: countNode,
    }, nil
}
```

### Analyse de Fraude

```go
func (fds *FraudDetectionSystem) AnalyzeAccount(accountToken *domain.Token, transactions []*domain.Fact) *FraudReport {
    report := &FraudReport{
        AccountID: accountToken.Facts[0].Fields["id"].(string),
        Score:     0,
        Reasons:   []string{},
    }
    
    // 1. Vérifier l'absence de transactions légitimes
    legitimateFound := false
    for _, tx := range transactions {
        if tx.Fields["type"] == "legitimate" {
            ageHours := calculateAgeInHours(tx.Timestamp)
            if ageHours < 24 {
                legitimateFound = true
                break
            }
        }
    }
    
    if !legitimateFound {
        report.Score += 30
        report.Reasons = append(report.Reasons, "Pas de transaction légitime récente")
    }
    
    // 2. Vérifier la présence de transactions suspectes
    if fds.ExistsNode.CheckExistence(accountToken) {
        report.Score += 50
        report.Reasons = append(report.Reasons, "Transactions suspectes détectées")
    }
    
    // 3. Vérifier la somme totale
    totalSum, _ := fds.SumNode.ComputeAggregate(accountToken, transactions)
    if sum := totalSum.(float64); sum > 50000 {
        report.Score += 20
        report.Reasons = append(report.Reasons, fmt.Sprintf("Montant élevé: %.2f", sum))
    }
    
    // Déterminer le niveau de risque
    if report.Score >= 70 {
        report.RiskLevel = "HIGH"
    } else if report.Score >= 40 {
        report.RiskLevel = "MEDIUM"
    } else {
        report.RiskLevel = "LOW"
    }
    
    return report
}
```

## 📊 Fonctions d'Agrégation

### SUM - Somme
```go
accumulator := domain.AccumulateFunction{
    FunctionType: "SUM",
    Field:        "amount",
}
// Calcule la somme de tous les montants
```

### COUNT - Comptage
```go
accumulator := domain.AccumulateFunction{
    FunctionType: "COUNT",
    Field:        "", // Pas de champ spécifique pour COUNT
}
// Compte le nombre total de faits
```

### AVG - Moyenne
```go
accumulator := domain.AccumulateFunction{
    FunctionType: "AVG", 
    Field:        "response_time",
}
// Calcule la moyenne des temps de réponse
```

### MIN - Minimum
```go
accumulator := domain.AccumulateFunction{
    FunctionType: "MIN",
    Field:        "price",
}
// Trouve le prix minimum
```

### MAX - Maximum  
```go
accumulator := domain.AccumulateFunction{
    FunctionType: "MAX",
    Field:        "severity_level",
}
// Trouve le niveau de sévérité maximum
```

## 🔄 Workflow Typique

### 1. Initialisation
```go
// Créer les nœuds
notNode := nodes.NewNotNode("id", logger)
existsNode := nodes.NewExistsNode("id", logger) 
accNode := nodes.NewAccumulateNode("id", accumulator, logger)

// Configurer les conditions
notNode.SetNegationCondition(condition)
existsNode.SetExistenceCondition(variable, condition)
```

### 2. Traitement des Faits
```go
// Ajouter des faits de droite (contexte)
for _, fact := range contextFacts {
    notNode.ProcessRightFact(fact)
    existsNode.ProcessRightFact(fact)
    accNode.ProcessRightFact(fact)
}
```

### 3. Traitement des Tokens
```go
// Traiter un token de gauche (déclencheur)
token := &domain.Token{
    ID:    "analysis_token",
    Facts: []*domain.Fact{subjectFact},
}

notNode.ProcessLeftToken(token)
existsNode.ProcessLeftToken(token)  
accNode.ProcessLeftToken(token)
```

### 4. Évaluation des Résultats
```go
// Vérifier les résultats
notResult := notNode.ProcessNegation(token, someFact)
existsResult := existsNode.CheckExistence(token)
aggResult, _ := accNode.ComputeAggregate(token, allFacts)
```

## 🎨 Patterns Avancés

### Pattern 1 : Détection d'Anomalie Temporelle
```go
// NOT : Pas d'activité normale dans les dernières heures
// EXISTS : Présence d'activité anormale
// ACCUMULATE : Pic de fréquence inhabituel
```

### Pattern 2 : Analyse de Performance  
```go
// NOT : Pas de succès récent
// EXISTS : Erreurs critiques présentes  
// ACCUMULATE : Moyenne de temps de réponse élevée
```

### Pattern 3 : Contrôle de Qualité
```go
// NOT : Pas de validation passée
// EXISTS : Défauts détectés
// ACCUMULATE : Taux d'échec au-dessus du seuil
```

## ⚡ Optimisations Performance

### 1. Indexation des Faits
```go
// Organiser les faits par type pour l'EXISTS
factsByType := make(map[string][]*domain.Fact)
for _, fact := range facts {
    factsByType[fact.Type] = append(factsByType[fact.Type], fact)
}
```

### 2. Cache des Résultats
```go
// Mettre en cache les résultats d'agrégation
type CachedAccumulator struct {
    *AccumulateNodeImpl
    cache map[string]interface{}
    mutex sync.RWMutex
}
```

### 3. Évaluation Paresseuse
```go
// Évaluer les conditions seulement si nécessaire
if quickCheck(token) {
    detailedResult := expensiveEvaluation(token, facts)
}
```

## 🚨 Gestion d'Erreurs

### Patterns de Récupération
```go
func SafeAggregation(node *AccumulateNodeImpl, token *domain.Token, facts []*domain.Fact) (interface{}, error) {
    defer func() {
        if r := recover(); r != nil {
            logger.Error("Aggregation panic recovered", fmt.Errorf("%v", r), nil)
        }
    }()
    
    return node.ComputeAggregate(token, facts)
}
```

### Validation des Données
```go
func ValidateAggregationInput(facts []*domain.Fact, field string) error {
    for _, fact := range facts {
        if _, exists := fact.Fields[field]; !exists {
            return fmt.Errorf("field %s not found in fact %s", field, fact.ID)
        }
    }
    return nil
}
```

## 🧪 Tests et Validation

### Test Pattern Recommandé
```go
func TestAdvancedNodePattern(t *testing.T) {
    // 1. Setup
    logger := &MockLogger{}
    node := createNodeUnderTest(logger)
    
    // 2. Prepare data
    facts := createTestFacts()
    token := createTestToken()
    
    // 3. Execute  
    result := executeNodeLogic(node, token, facts)
    
    // 4. Verify
    validateResult(t, result, expectedOutcome)
}
```

---

Ce guide couvre l'utilisation pratique des nœuds RETE avancés. Pour des exemples plus détaillés, consultez les tests d'intégration dans `advanced_integration_test.go`.