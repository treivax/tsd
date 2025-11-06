# Module RETE - Moteur d'inférence avec persistance etcd

Le module RETE implémente un réseau d'inférence basé sur l'algorithme RETE qui construit automatiquement un réseau de nœuds à partir d'un AST de règles métier et permet l'exécution efficace d'actions basées sur des faits.

## 🏗️ Architecture

```
AST (constraint) → Réseau RETE → Actions déclenchées
                      ↓
                   etcd (persistance)
```

### Types de nœuds

1. **RootNode** : Point d'entrée pour tous les faits
2. **TypeNode** : Filtre les faits par type et valide leur structure  
3. **AlphaNode** : Teste les conditions sur les faits individuels
4. **TerminalNode** : Déclenche les actions quand les conditions sont remplies

### Persistance

Chaque nœud sauvegarde automatiquement son état (Working Memory) dans etcd :
- Faits correspondants aux conditions du nœud
- Tokens de propagation 
- Timestamps de dernière modification

## 🚀 Utilisation

### Exemple basique

```go
package main

import (
    "github.com/treivax/tsd/rete"
)

func main() {
    // 1. Créer le storage
    storage := rete.NewMemoryStorage()
    
    // 2. Créer le réseau
    network := rete.NewReteNetwork(storage)
    
    // 3. Charger les règles depuis un AST
    err := network.LoadFromAST(program)
    if err != nil {
        panic(err)
    }
    
    // 4. Soumettre des faits
    fact := &rete.Fact{
        ID:   "person1",
        Type: "Person",
        Fields: map[string]interface{}{
            "age": 25,
            "name": "Alice",
        },
    }
    
    err = network.SubmitFact(fact)
    if err != nil {
        panic(err)
    }
    
    // Les actions sont automatiquement déclenchées !
}
```

### Avec etcd

```go

```

## 📊 Fonctionnalités

### ✅ Implémenté

- [x] Construction automatique du réseau depuis AST
- [x] Propagation efficace des faits 
- [x] Filtrage par type avec validation
- [x] Déclenchement d'actions conditionnelles
- [x] Persistance etcd de l'état complet
- [x] Storage en mémoire pour les tests
- [x] Logging détaillé du flux d'exécution
- [x] API complète de gestion du réseau

### 🔄 Améliorations futures possibles

- [ ] Évaluation complète des expressions de condition
- [ ] Nœuds Beta pour les jointures multi-faits
- [ ] Optimisations de performance (indexing)
- [ ] Interface web de monitoring
- [ ] Métriques et observabilité

## 🏃 Exécution

### Démo interactive

```bash
# Compiler et exécuter la démo
go build -o rete-demo ./rete/cmd/
./rete-demo

# Sortie attendue :
# 🔥 DÉMONSTRATION DU RÉSEAU RETE
# ===============================================
# 
# 📋 ÉTAPE 1: Création du programme RETE
# ✅ Programme créé avec 1 type(s) et 1 expression(s)
# 
# [... construction du réseau ...]
# 
# 🎯 ACTION DÉCLENCHÉE: action
#    Arguments: [client]
#    Faits correspondants:
#      - { "id": "personne_1", "type": "Personne", ... }
```

### Tests

```bash
# Exécuter les tests (à venir)
go test ./rete/
```

## 🛠️ API

### Interfaces principales

```go
// Network principal
type ReteNetwork struct {
    LoadFromAST(program *Program) error
    SubmitFact(fact *Fact) error
    GetNetworkState() (map[string]*WorkingMemory, error)
}

// Storage pour la persistance
type Storage interface {
    SaveMemory(nodeID string, memory *WorkingMemory) error
    LoadMemory(nodeID string) (*WorkingMemory, error) 
    DeleteMemory(nodeID string) error
    ListNodes() ([]string, error)
}

// Nœud du réseau
type Node interface {
    ActivateLeft(token *Token) error
    ActivateRight(fact *Fact) error
}
```

## 📈 Performance

Le système est conçu pour :
- **Scalabilité** : Ajout dynamique de règles et faits
- **Persistance** : État complet sauvé en temps réel dans etcd  
- **Concurrence** : Safe pour l'utilisation multi-thread
- **Efficacité** : Propagation optimisée selon l'algorithme RETE

## 🔗 Intégration

Ce module s'intègre parfaitement avec :
- **Module constraint** : Parse les règles métier
- **etcd** : Stockage distribué de l'état
- **Systèmes distribués** : Multiple instances avec état partagé

---

*Le module RETE fournit une base solide pour des systèmes experts, moteurs de règles métier, et systèmes d'inférence nécessitant une persistance robuste.*