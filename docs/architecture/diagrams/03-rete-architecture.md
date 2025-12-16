# 🧠 Architecture Moteur RETE

**Date** : 2025-12-16  
**Version** : 1.0.0  
**Statut** : Documentation officielle

---

## Vue d'Ensemble du Moteur RETE

Le moteur RETE (Rete Algorithm) est au cœur de TSD. Il permet l'évaluation efficace de règles complexes.

```mermaid
graph TB
    subgraph "Input"
        PROG[📋 Programme TSD<br/>Types + Rules + Facts]
    end

    subgraph "RETE Engine Core"
        direction TB
        
        subgraph "Alpha Network"
            A1[🔍 Type Tests]
            A2[🔍 Simple Conditions]
            A3[🔍 Field Checks]
            AM[💾 Alpha Memory]
            
            A1 --> A2
            A2 --> A3
            A3 --> AM
        end

        subgraph "Beta Network"
            B1[🔗 Join Nodes]
            B2[🧮 Aggregate Nodes]
            B3[❌ NOT Nodes]
            BM[💾 Beta Memory]
            
            B1 --> B2
            B2 --> B3
            B3 --> BM
        end

        subgraph "Terminal Layer"
            T1[🎯 Terminal Node 1]
            T2[🎯 Terminal Node 2]
            TN[🎯 Terminal Node N]
            
            T1 -.-> T2
            T2 -.-> TN
        end
    end

    subgraph "Output"
        ACT[✅ Activations<br/>Règles déclenchées]
    end

    PROG -->|Facts| A1
    AM -->|Tokens| B1
    BM --> T1
    BM --> T2
    BM --> TN
    
    T1 --> ACT
    T2 --> ACT
    TN --> ACT

    style PROG fill:#E8F4F8,stroke:#333
    style A1 fill:#FFE5CC,stroke:#333
    style A2 fill:#FFE5CC,stroke:#333
    style A3 fill:#FFE5CC,stroke:#333
    style AM fill:#FFF4CC,stroke:#333,stroke-width:2px
    style B1 fill:#CCE5FF,stroke:#333
    style B2 fill:#CCE5FF,stroke:#333
    style B3 fill:#CCE5FF,stroke:#333
    style BM fill:#CCF0FF,stroke:#333,stroke-width:2px
    style T1 fill:#D5F4E6,stroke:#333
    style T2 fill:#D5F4E6,stroke:#333
    style TN fill:#D5F4E6,stroke:#333
    style ACT fill:#C3F0CA,stroke:#333,stroke-width:3px
```

---

## Types de Nœuds RETE

```mermaid
classDiagram
    class Node {
        <<interface>>
        +ProcessToken(token)
        +GetMemory()
    }

    class AlphaNode {
        -typeTest TypeTest
        -condition Condition
        -memory AlphaMemory
        +EvaluateCondition(fact)
    }

    class BetaNode {
        <<interface>>
        +ProcessLeftToken(token)
        +ProcessRightToken(token)
    }

    class JoinNode {
        -leftMemory LeftMemory
        -rightMemory RightMemory
        -joinCondition Condition
        +PerformJoin()
    }

    class AggregateNode {
        -aggregateType Type
        -sourcePattern Pattern
        -aggregateCondition Condition
        +ComputeAggregate()
    }

    class NotNode {
        -negatedPattern Pattern
        -negatedCondition Condition
        +CheckAbsence()
    }

    class TerminalNode {
        -ruleID string
        -activations []Activation
        +RecordActivation(tokens)
    }

    Node <|-- AlphaNode
    Node <|-- BetaNode
    BetaNode <|-- JoinNode
    BetaNode <|-- AggregateNode
    BetaNode <|-- NotNode
    Node <|-- TerminalNode

    AlphaNode --> JoinNode : "Tokens"
    JoinNode --> AggregateNode : "Joined Tokens"
    AggregateNode --> NotNode : "Aggregated Tokens"
    NotNode --> TerminalNode : "Final Tokens"
```

---

## Réseau Alpha - Filtrage des Faits

```mermaid
graph TB
    subgraph "Input Facts"
        F1[Person{id:1, age:25}]
        F2[Person{id:2, age:17}]
        F3[Order{id:1, amount:150}]
        F4[Order{id:2, amount:50}]
    end

    subgraph "Type Filter Layer"
        TP[Type=Person]
        TO[Type=Order]
    end

    subgraph "Condition Filter Layer"
        C1[age >= 18]
        C2[age < 18]
        C3[amount > 100]
        C4[amount <= 100]
    end

    subgraph "Alpha Memory"
        M1[Person age>=18<br/>✓ id:1]
        M2[Person age<18<br/>✓ id:2]
        M3[Order amount>100<br/>✓ id:1]
        M4[Order amount<=100<br/>✓ id:2]
    end

    F1 --> TP
    F2 --> TP
    F3 --> TO
    F4 --> TO

    TP --> C1
    TP --> C2
    TO --> C3
    TO --> C4

    C1 -->|Match| M1
    C2 -->|Match| M2
    C3 -->|Match| M3
    C4 -->|Match| M4

    style F1 fill:#E8F4F8
    style F2 fill:#E8F4F8
    style F3 fill:#E8F4F8
    style F4 fill:#E8F4F8
    style M1 fill:#D5F4E6,stroke:#333,stroke-width:2px
    style M2 fill:#D5F4E6,stroke:#333,stroke-width:2px
    style M3 fill:#D5F4E6,stroke:#333,stroke-width:2px
    style M4 fill:#D5F4E6,stroke:#333,stroke-width:2px
```

---

## Réseau Beta - Jointures

```mermaid
graph TB
    subgraph "Alpha Outputs"
        AL[Left Memory<br/>Person age>=18]
        AR[Right Memory<br/>Order amount>100]
    end

    subgraph "Join Node"
        J[🔗 Join<br/>p.id == o.customer_id]
    end

    subgraph "Join Evaluation"
        E1[Compare IDs]
        E2[Check Condition]
        E3[Create Combined Token]
    end

    subgraph "Beta Memory"
        BM[Combined Tokens<br/>Person+Order pairs]
    end

    subgraph "Downstream"
        AGG[Aggregate Node]
        TERM[Terminal Node]
    end

    AL -->|Left Token| J
    AR -->|Right Token| J
    
    J --> E1
    E1 --> E2
    E2 --> E3
    E3 --> BM
    
    BM --> AGG
    AGG --> TERM

    style AL fill:#FFE5CC
    style AR fill:#FFE5CC
    style J fill:#CCE5FF,stroke:#333,stroke-width:2px
    style BM fill:#CCF0FF,stroke:#333,stroke-width:2px
    style TERM fill:#D5F4E6,stroke:#333,stroke-width:2px
```

---

## Nœud d'Agrégation

```mermaid
graph LR
    subgraph "Source Pattern"
        SP[Order / customer_id == p.id]
    end

    subgraph "Aggregate Functions"
        SUM[SUM: o.amount]
        COUNT[COUNT: *]
        AVG[AVG: o.amount]
        MIN[MIN: o.amount]
        MAX[MAX: o.amount]
    end

    subgraph "Condition Check"
        COND[Result >= 1000]
    end

    subgraph "Output"
        OUT[Token with<br/>aggregate value]
    end

    SP --> SUM
    SP --> COUNT
    SP --> AVG
    SP --> MIN
    SP --> MAX
    
    SUM --> COND
    COUNT --> COND
    AVG --> COND
    MIN --> COND
    MAX --> COND
    
    COND -->|Pass| OUT

    style SP fill:#E8F4F8
    style SUM fill:#FFDAB9
    style COUNT fill:#FFDAB9
    style AVG fill:#FFDAB9
    style MIN fill:#FFDAB9
    style MAX fill:#FFDAB9
    style COND fill:#CCE5FF
    style OUT fill:#D5F4E6,stroke:#333,stroke-width:2px
```

---

## Nœud NOT (Négation)

```mermaid
stateDiagram-v2
    [*] --> CheckLeftToken: Left Token arrives
    
    CheckLeftToken --> SearchRightMemory: Check for matches
    SearchRightMemory --> MatchFound: Matching Right Token exists
    SearchRightMemory --> NoMatch: No matching Right Token
    
    MatchFound --> BlockToken: Block propagation
    NoMatch --> PropagateToken: Allow propagation
    
    PropagateToken --> [*]: Token continues
    BlockToken --> [*]: Token stopped
    
    note right of NoMatch
        Négation réussie:
        Aucun fait correspondant
        → Règle peut s'activer
    end note
    
    note right of MatchFound
        Négation échoue:
        Fait correspondant existe
        → Règle bloquée
    end note
```

---

## Optimisations - Alpha Sharing

Plusieurs règles peuvent partager les mêmes nœuds Alpha si elles ont les mêmes conditions.

```mermaid
graph TB
    subgraph "Sans Alpha Sharing"
        R1[Rule 1: age >= 18]
        R2[Rule 2: age >= 18]
        R3[Rule 3: age >= 18]
        
        A1[AlphaNode 1]
        A2[AlphaNode 2]
        A3[AlphaNode 3]
        
        R1 --> A1
        R2 --> A2
        R3 --> A3
    end

    subgraph "Avec Alpha Sharing"
        RR1[Rule 1]
        RR2[Rule 2]
        RR3[Rule 3]
        
        AA[Shared AlphaNode<br/>age >= 18]
        
        RR1 --> AA
        RR2 --> AA
        RR3 --> AA
    end

    style A1 fill:#FFB6C1
    style A2 fill:#FFB6C1
    style A3 fill:#FFB6C1
    style AA fill:#D5F4E6,stroke:#333,stroke-width:3px
```

**Avantages :**
- ✅ Réduction mémoire : 1 nœud au lieu de 3
- ✅ Réduction CPU : 1 évaluation au lieu de 3
- ✅ Cache plus efficace

---

## Optimisations - Beta Sharing

```mermaid
graph TB
    subgraph "Sans Beta Sharing"
        R1[Rule 1]
        R2[Rule 2]
        
        J1[JoinNode 1<br/>p.id == o.customer_id]
        J2[JoinNode 2<br/>p.id == o.customer_id]
        
        R1 --> J1
        R2 --> J2
    end

    subgraph "Avec Beta Sharing"
        RR1[Rule 1]
        RR2[Rule 2]
        
        JJ[Shared JoinNode<br/>p.id == o.customer_id]
        
        RR1 --> JJ
        RR2 --> JJ
    end

    style J1 fill:#FFB6C1
    style J2 fill:#FFB6C1
    style JJ fill:#D5F4E6,stroke:#333,stroke-width:3px
```

**Impact :**
- 🎯 60-80% réduction nœuds Beta
- 🚀 Jusqu'à 10x amélioration performance

---

## Architecture Complète d'une Règle

Exemple de règle : `rule vip : {p: Person} / SUM(o: Order / o.customer_id == p.id ; o.amount) >= 1000 ==> vip(p.id)`

```mermaid
graph TB
    subgraph "Facts Input"
        P[Person{id:1}]
        O1[Order{customer_id:1, amount:600}]
        O2[Order{customer_id:1, amount:500}]
    end

    subgraph "Alpha Network"
        AP[TypeTest: Person]
        AO[TypeTest: Order]
        AOM[AlphaMemory: Order]
    end

    subgraph "Beta Network"
        J[JoinNode<br/>o.customer_id == p.id]
        AGG[AggregateNode<br/>SUM o.amount]
        CHK[ConditionCheck<br/>result >= 1000]
    end

    subgraph "Terminal"
        T[TerminalNode<br/>rule: vip]
    end

    subgraph "Action"
        ACT[Action: vip]
        NEW[NewFact: vip]
    end

    P --> AP
    O1 --> AO
    O2 --> AO
    
    AO --> AOM
    
    AP -->|Left| J
    AOM -->|Right| J
    
    J --> AGG
    AGG -->|1100| CHK
    CHK -->|Pass| T
    
    T --> ACT
    ACT --> NEW

    style P fill:#E8F4F8
    style O1 fill:#E8F4F8
    style O2 fill:#E8F4F8
    style AGG fill:#FFE5CC,stroke:#333,stroke-width:2px
    style T fill:#D5F4E6,stroke:#333,stroke-width:2px
    style NEW fill:#C3F0CA,stroke:#333,stroke-width:3px
```

---

## Métriques Performance

```mermaid
graph LR
    subgraph "Input"
        F[1000 Facts]
    end

    subgraph "Performances"
        T1[⏱️ <1ms par règle]
        T2[🚀 10-50K facts/sec]
        T3[💾 -80% mémoire<br/>avec sharing]
        T4[🎯 O log n<br/>complexité]
    end

    F --> T1
    F --> T2
    F --> T3
    F --> T4

    style F fill:#E8F4F8
    style T1 fill:#D5F4E6
    style T2 fill:#D5F4E6
    style T3 fill:#D5F4E6
    style T4 fill:#D5F4E6
```

---

## Structures de Données Clés

```mermaid
classDiagram
    class Engine {
        -alphaNetwork AlphaNetwork
        -betaNetwork BetaNetwork
        -workingMemory WorkingMemory
        -ruleIndex map[string]*TerminalNode
        +CompileProgram(ast)
        +AssertFact(fact)
        +ExecuteActions()
    }

    class AlphaNetwork {
        -rootNode *AlphaNode
        -typeIndex map[string]*AlphaNode
        -chainCache map[string]*AlphaChain
        +GetOrCreateChain(conditions)
        +EvaluateFact(fact)
    }

    class BetaNetwork {
        -joinNodes []*JoinNode
        -aggregateNodes []*AggregateNode
        -notNodes []*NotNode
        -sharingIndex map[string]BetaNode
        +GetSharedNode(signature)
        +ProcessToken(token)
    }

    class WorkingMemory {
        -facts map[string]*Fact
        -factIndex map[string][]*Fact
        +StoreFact(fact)
        +GetFactsByType(typeName)
    }

    class Token {
        -facts []*Fact
        -parent *Token
        -timestamp int64
        +GetValue(path)
        +Clone()
    }

    Engine --> AlphaNetwork
    Engine --> BetaNetwork
    Engine --> WorkingMemory
    BetaNetwork --> Token
    AlphaNetwork --> Token
```

---

## Références

- [Architecture Globale](01-global-architecture.md)
- [Flux de Données](02-data-flow.md)
- [Optimisations RETE](../architecture.md#optimisations)

---

**Maintenu par** : TSD Contributors  
**Dernière mise à jour** : 2025-12-16
