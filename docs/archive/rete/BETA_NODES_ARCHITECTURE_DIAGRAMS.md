# Diagrammes d'Architecture des BetaNodes

**Date**: 2025-01-27  
**Version**: 1.0  
**Complément à**: BETA_NODES_ANALYSIS.md

---

## Table des Matières

1. [Architecture Actuelle - Sans Partage](#architecture-actuelle---sans-partage)
2. [Architecture Proposée - Avec Partage](#architecture-proposée---avec-partage)
3. [Flux de Données](#flux-de-données)
4. [Cascades de Jointures](#cascades-de-jointures)
5. [Gestion du Cycle de Vie](#gestion-du-cycle-de-vie)
6. [Comparaison Alpha vs Beta](#comparaison-alpha-vs-beta)

---

## Architecture Actuelle - Sans Partage

### Scénario: 3 Règles avec Même Jointure

```
Règles:
  rule validate : {u: User, o: Order} / o.user_id == u.id ==> validate(o)
  rule notify   : {u: User, o: Order} / o.user_id == u.id ==> notify(u, o)
  rule invoice  : {u: User, o: Order} / o.user_id == u.id ==> invoice(u, o)

┌─────────────────────────────────────────────────────────────────────────┐
│                          RETE Network (ACTUEL)                          │
└─────────────────────────────────────────────────────────────────────────┘

                           ┌──────────────┐
                           │   RootNode   │
                           └──────┬───────┘
                                  │
                    ┌─────────────┴─────────────┐
                    │                           │
              ┌─────▼──────┐            ┌──────▼─────┐
              │ TypeNode   │            │ TypeNode   │
              │   (User)   │            │  (Order)   │
              └─────┬──────┘            └──────┬─────┘
                    │                           │
        ┌───────────┼───────────┐   ┌───────────┼───────────┐
        │           │           │   │           │           │
    ┌───▼───┐   ┌───▼───┐   ┌───▼───┐ ┌───▼───┐ ┌───▼───┐ ┌───▼───┐
    │ Alpha │   │ Alpha │   │ Alpha │ │ Alpha │ │ Alpha │ │ Alpha │
    │ Pass1 │   │ Pass2 │   │ Pass3 │ │ Pass1 │ │ Pass2 │ │ Pass3 │
    │ [LEFT]│   │ [LEFT]│   │ [LEFT]│ │[RIGHT]│ │[RIGHT]│ │[RIGHT]│
    └───┬───┘   └───┬───┘   └───┬───┘ └───┬───┘ └───┬───┘ └───┬───┘
        │           │           │         │         │         │
        └─────┐     │     ┏━━━━━┷━━━━┓    │     ┏━━━┷━━━━┓    │
              │     │     ┃          ┃    │     ┃        ┃    │
        ┌─────▼─────▼─────▼──────┐   ┃    │     ┃  ┌─────▼────▼────┐
        │   JoinNode #1          │◄──┛    │     ┃  │  JoinNode #3  │
        │ o.user_id == u.id      │        │     ┃  │o.user_id==u.id│
        │ LeftMemory: [...]      │        │     ┃  │LeftMemory:[..]│
        │ RightMemory: [...]     │◄───────┘     ┃  │RightMem: [..] │
        │ ResultMemory: [...]    │              ┃  │ResultMem:[..] │
        └─────────┬──────────────┘              ┃  └───────┬───────┘
                  │                             ┃          │
            ┌─────▼──────┐                      ┃    ┌─────▼──────┐
            │ Terminal   │                      ┃    │ Terminal   │
            │ (validate) │                      ┃    │ (invoice)  │
            └────────────┘                      ┃    └────────────┘
                                                ┃
                                    ┌───────────┘
                                    │
                              ┌─────▼─────────────┐
                              │   JoinNode #2     │
                              │ o.user_id == u.id │
                              │ LeftMemory: [...] │
                              │ RightMemory: [...]│
                              │ ResultMemory:[...]│
                              └─────────┬─────────┘
                                        │
                                  ┌─────▼──────┐
                                  │ Terminal   │
                                  │  (notify)  │
                                  └────────────┘

❌ PROBLÈMES:
  • 3 JoinNodes DUPLIQUÉS avec condition IDENTIQUE
  • 3 × LeftMemory stockant les mêmes tokens User
  • 3 × RightMemory stockant les mêmes tokens Order
  • 3 × évaluations de o.user_id == u.id
  • Mémoire: 3x plus que nécessaire
  • Calculs: 3x plus que nécessaire
```

---

## Architecture Proposée - Avec Partage

### Même Scénario avec Partage Activé

```
Règles: (identiques)
  rule validate : {u: User, o: Order} / o.user_id == u.id ==> validate(o)
  rule notify   : {u: User, o: Order} / o.user_id == u.id ==> notify(u, o)
  rule invoice  : {u: User, o: Order} / o.user_id == u.id ==> invoice(u, o)

┌─────────────────────────────────────────────────────────────────────────┐
│                    RETE Network (AVEC PARTAGE)                          │
│                                                                         │
│  ┌────────────────────────────────────────────────────────────────┐    │
│  │ BetaSharingRegistry                                            │    │
│  │  sharedJoinNodes["join_a3f8b92e"] → JoinNode (RefCount=3)    │    │
│  │  lruHashCache: 1000 entries                                   │    │
│  │  metrics: created=1, reused=2, ratio=66%                      │    │
│  └────────────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────────────┘

                           ┌──────────────┐
                           │   RootNode   │
                           └──────┬───────┘
                                  │
                    ┌─────────────┴─────────────┐
                    │                           │
              ┌─────▼──────┐            ┌──────▼─────┐
              │ TypeNode   │            │ TypeNode   │
              │   (User)   │            │  (Order)   │
              └─────┬──────┘            └──────┬─────┘
                    │                           │
                    │                           │
              ┌─────▼──────┐            ┌──────▼─────┐
              │   Alpha    │            │   Alpha    │
              │   Pass_u   │            │   Pass_o   │
              │   [LEFT]   │            │  [RIGHT]   │
              └─────┬──────┘            └──────┬─────┘
                    │                           │
                    └─────────┐     ┌───────────┘
                              │     │
                        ┏━━━━━▼━━━━━▼━━━━━┓
                        ┃                  ┃
                  ┌─────▼──────────────────▼────┐
                  │   SHARED JoinNode           │
                  │   ID: join_a3f8b92e         │
                  │   Condition: o.user_id==u.id│
                  │   RefCount: 3 ◄──────────┐  │
                  │                          │  │
                  │   LeftMemory: [...]      │  │
                  │   RightMemory: [...]     │  │
                  │   ResultMemory: [...]    │  │
                  │                          │  │
                  └─────┬──────┬──────┬──────┘  │
                        │      │      │         │
              ┌─────────┘      │      └────────┐│
              │                │               ││
        ┌─────▼──────┐  ┌──────▼──────┐ ┌─────▼▼─────┐
        │ Terminal   │  │  Terminal   │ │ Terminal   │
        │ (validate) │  │  (notify)   │ │ (invoice)  │
        └────────────┘  └─────────────┘ └────────────┘

┌────────────────────────────────────────────────────────────┐
│ LifecycleManager                                           │
│  NodeLifecycle["join_a3f8b92e"]                           │
│    RefCount: 3                                            │
│    Rules: ["validate", "notify", "invoice"]              │
│    CreatedBy: ["validate", "notify", "invoice"]          │
└────────────────────────────────────────────────────────────┘

✅ AVANTAGES:
  • 1 seul JoinNode partagé (vs 3)
  • 1 × mémoire (LeftMemory, RightMemory, ResultMemory)
  • 1 × évaluation de condition par fait
  • Résultats propagés aux 3 TerminalNodes
  • Réduction mémoire: 66%
  • Réduction calculs: 66%
  • RefCount automatique via LifecycleManager
```

---

## Flux de Données

### Propagation d'un Fait dans un JoinNode Partagé

```
┌───────────────────────────────────────────────────────────────────────┐
│  Étape 1: Soumission d'un nouveau fait User                          │
└───────────────────────────────────────────────────────────────────────┘

  User{id:"U1", name:"Alice"} ──┐
                                 │
                           ┌─────▼──────┐
                           │ RootNode   │
                           └─────┬──────┘
                                 │
                           ┌─────▼──────┐
                           │ TypeNode   │
                           │  (User)    │
                           └─────┬──────┘
                                 │
                           ┌─────▼──────┐
                           │ AlphaPass  │
                           │   [LEFT]   │
                           └─────┬──────┘
                                 │
                    ┌────────────▼────────────┐
                    │ SharedJoinNode          │
                    │ ActivateLeft(token_u1)  │
                    │                         │
                    │ LeftMemory.Add(token_u1)│
                    │   {"u": User_U1}        │
                    │                         │
                    │ RightMemory: [...]      │◄─── Mémoire persistante
                    │   (vide pour l'instant) │     entre soumissions
                    └─────────────────────────┘

  → Aucun token joint créé (RightMemory vide)
  → Aucune propagation aux Terminals


┌───────────────────────────────────────────────────────────────────────┐
│  Étape 2: Soumission d'un fait Order CORRESPONDANT                   │
└───────────────────────────────────────────────────────────────────────┘

  Order{id:"O1", user_id:"U1", amount:100} ──┐
                                              │
                           ┌──────────────────▼─┐
                           │ TypeNode (Order)   │
                           └──────────┬─────────┘
                                      │
                           ┌──────────▼─────────┐
                           │ AlphaPass [RIGHT]  │
                           └──────────┬─────────┘
                                      │
                    ┌─────────────────▼──────────────┐
                    │ SharedJoinNode                 │
                    │ ActivateRight(fact_o1)         │
                    │                                │
                    │ 1️⃣ RightMemory.Add(token_o1)    │
                    │    {"o": Order_O1}             │
                    │                                │
                    │ 2️⃣ Parcourir LeftMemory:        │
                    │    token_u1 = {"u": User_U1}   │
                    │                                │
                    │ 3️⃣ Évaluer jointure:            │
                    │    u.id ("U1") == o.user_id ("U1") ✅
                    │                                │
                    │ 4️⃣ Créer token joint:           │
                    │    token_joined = {            │
                    │      "u": User_U1,             │
                    │      "o": Order_O1             │
                    │    }                           │
                    │                                │
                    │ 5️⃣ ResultMemory.Add(token_joined) │
                    │                                │
                    │ 6️⃣ Propager aux TOUS enfants:   │
                    └───┬────────┬───────────┬───────┘
                        │        │           │
              ┌─────────▼──┐  ┌──▼──────┐ ┌─▼────────┐
              │ Terminal   │  │Terminal │ │Terminal  │
              │ validate() │  │notify() │ │invoice() │
              │   FIRED!   │  │ FIRED!  │ │  FIRED!  │
              └────────────┘  └─────────┘ └──────────┘

  ✅ UNE évaluation → TROIS règles activées


┌───────────────────────────────────────────────────────────────────────┐
│  Étape 3: Soumission d'un fait Order NON CORRESPONDANT               │
└───────────────────────────────────────────────────────────────────────┘

  Order{id:"O2", user_id:"U999", amount:50} ──┐
                                               │
                           ┌───────────────────▼┐
                           │ TypeNode (Order)   │
                           └───────────┬────────┘
                                       │
                    ┌──────────────────▼─────────────┐
                    │ SharedJoinNode                 │
                    │ ActivateRight(fact_o2)         │
                    │                                │
                    │ 1️⃣ RightMemory.Add(token_o2)    │
                    │    {"o": Order_O2}             │
                    │                                │
                    │ 2️⃣ Parcourir LeftMemory:        │
                    │    token_u1 = {"u": User_U1}   │
                    │                                │
                    │ 3️⃣ Évaluer jointure:            │
                    │    u.id ("U1") == o.user_id ("U999") ❌
                    │                                │
                    │ 4️⃣ Jointure échoue → Filtrage   │
                    │    Aucun token créé            │
                    └────────────────────────────────┘

  ❌ Aucune propagation (filtrage correct)
```

---

## Cascades de Jointures

### Architecture en Cascade (3 Variables)

```
Règle: {u: User, o: Order, p: Product} / 
       o.user_id == u.id AND o.product_id == p.id

┌─────────────────────────────────────────────────────────────────────────┐
│                      CASCADE SANS PARTAGE                               │
└─────────────────────────────────────────────────────────────────────────┘

TypeNode(User) ──► AlphaPass_u ──┐
                                 │
                            ┌────▼────────────┐
                            │ JoinNode1       │
TypeNode(Order) ──► Alpha ──┤ (u ⋈ o)         │
                    Pass_o  │ u.id==o.user_id │
                            └────┬────────────┘
                                 │
                                 │ token_uo = {"u": User, "o": Order}
                                 │
                            ┌────▼────────────────┐
                            │ JoinNode2           │
TypeNode(Product)──►Alpha──┤ ((u,o) ⋈ p)         │
                   Pass_p  │ o.product_id==p.id  │
                            └────┬────────────────┘
                                 │
                                 │ token_uop = {"u", "o", "p"}
                                 │
                            ┌────▼────────┐
                            │ Terminal    │
                            └─────────────┘


┌─────────────────────────────────────────────────────────────────────────┐
│              CASCADE AVEC PARTAGE PARTIEL (2 Règles)                    │
│                                                                         │
│  Rule A: {u: User, o: Order, p: Product} /                            │
│          o.user_id == u.id AND o.product_id == p.id                   │
│                                                                         │
│  Rule B: {u: User, o: Order, s: Shipment} /                           │
│          o.user_id == u.id AND s.order_id == o.id                     │
└─────────────────────────────────────────────────────────────────────────┘

TypeNode(User) ──► AlphaPass_u [LEFT] ──┐
                                        │
                                   ┏━━━━▼━━━━━━━━━━━━━━┓
                                   ┃ SHARED JoinNode1  ┃
TypeNode(Order) ──► AlphaPass_o ───┃ (u ⋈ o)          ┃
                    [RIGHT]        ┃ RefCount=2 ◄─────┃──┐
                                   ┃ u.id==o.user_id  ┃  │
                                   ┗━━━━┳━━━━━━━━━━━━━━┛  │
                                        │                 │
                        ┌───────────────┴────────┐        │
                        │                        │        │
                 token_uo                  token_uo       │
                 {"u","o"}                 {"u","o"}      │
                        │                        │        │
              ┌─────────▼─────────┐    ┌─────────▼────────┐
              │ JoinNode2_A       │    │ JoinNode2_B      │
              │ ((u,o) ⋈ p)       │    │ ((u,o) ⋈ s)      │
              │ o.prod_id==p.id   │    │ s.order_id==o.id │
              │ RefCount=1        │    │ RefCount=1       │
              └─────────┬─────────┘    └─────────┬────────┘
                        │                        │
                        ▲                        ▲
                        │                        │
TypeNode(Product)──►AlphaPass_p          TypeNode(Shipment)──►AlphaPass_s
                    [RIGHT]                                    [RIGHT]
                        │                        │
                 ┌──────▼──────┐        ┌───────▼──────┐
                 │ Terminal_A  │        │ Terminal_B   │
                 └─────────────┘        └──────────────┘

✅ OPTIMISATION:
  • JoinNode1 (u ⋈ o) PARTAGÉ entre Rule A et Rule B
  • JoinNode2_A et JoinNode2_B UNIQUES (conditions différentes)
  • Réduction: 1 JoinNode économisé sur 4 (25%)
  • Évaluation u.id==o.user_id faite UNE fois pour les 2 règles
```

### Cascade 4 Variables (Maximum Sharing)

```
Rule A: {u, o, p, i} / u⋈o AND o⋈p AND p⋈i
Rule B: {u, o, p, s} / u⋈o AND o⋈p AND p⋈s
Rule C: {u, o, q, r} / u⋈o AND o⋈q AND q⋈r

                    ┏━━━━━━━━━━━━━━━━━━┓
                    ┃ Shared: u ⋈ o    ┃ RefCount=3 (A,B,C)
                    ┗━━━━━━━┳━━━━━━━━━━┛
                            │
                ┌───────────┴───────────┐
                │                       │
         ┏━━━━━━▼━━━━━━━┓        ┌──────▼──────┐
         ┃ Shared:(u,o)⋈p┃        │ (u,o) ⋈ q   │
         ┃ RefCount=2    ┃        │ RefCount=1  │
         ┗━━━━━━━┳━━━━━━━┛        └──────┬──────┘
                 │                       │
        ┌────────┴────────┐              │
        │                 │              │
   ┌────▼─────┐      ┌────▼─────┐  ┌────▼─────┐
   │(u,o,p)⋈i │      │(u,o,p)⋈s │  │(u,o,q)⋈r │
   │Rule A    │      │Rule B    │  │Rule C    │
   └──────────┘      └──────────┘  └──────────┘

Partage: 2 JoinNodes partagés sur 9 total (22%)
Mais: Les 2 partagés sont les plus coûteux (début cascade)
```

---

## Gestion du Cycle de Vie

### Ajout Progressif de Règles

```
┌─────────────────────────────────────────────────────────────────────────┐
│  T0: Réseau Vide                                                        │
└─────────────────────────────────────────────────────────────────────────┘

BetaSharingRegistry: {}
LifecycleManager: {}


┌─────────────────────────────────────────────────────────────────────────┐
│  T1: Ajout Rule1                                                        │
│      rule r1 : {u: User, o: Order} / o.user_id == u.id ==> action1()  │
└─────────────────────────────────────────────────────────────────────────┘

1️⃣ Calculer hash: hash("o.user_id == u.id", ["u"], ["o"]) = "join_a3f8"
2️⃣ BetaSharingRegistry.GetOrCreate("join_a3f8")
   → Pas trouvé → Créer JoinNode("join_a3f8")
3️⃣ LifecycleManager.RegisterNode("join_a3f8", "join")
4️⃣ LifecycleManager.AddRuleToNode("join_a3f8", "r1")

État:
  BetaSharingRegistry: {
    "join_a3f8": JoinNode(condition, RefCount via Lifecycle)
  }
  LifecycleManager: {
    "join_a3f8": NodeLifecycle {
      RefCount: 1,
      Rules: ["r1"]
    }
  }

Network:
  TypeNode(User) ──► AlphaPass ──┐
                                 │
                            ┌────▼─────┐
  TypeNode(Order)──►AlphaPass──┤join_a3f8│──► Terminal(r1)
                            └──────────┘


┌─────────────────────────────────────────────────────────────────────────┐
│  T2: Ajout Rule2 (MÊME jointure)                                       │
│      rule r2 : {u: User, o: Order} / o.user_id == u.id ==> action2()  │
└─────────────────────────────────────────────────────────────────────────┘

1️⃣ Calculer hash: hash("o.user_id == u.id", ["u"], ["o"]) = "join_a3f8"
2️⃣ BetaSharingRegistry.GetOrCreate("join_a3f8")
   → TROUVÉ! ♻️  Réutiliser JoinNode existant
3️⃣ LifecycleManager.AddRuleToNode("join_a3f8", "r2")
   → Incrémenter RefCount

État:
  BetaSharingRegistry: {
    "join_a3f8": JoinNode (même instance)
  }
  LifecycleManager: {
    "join_a3f8": NodeLifecycle {
      RefCount: 2 ⬆️,
      Rules: ["r1", "r2"]
    }
  }
  Metrics:
    TotalCreated: 1
    TotalReused: 1
    SharingRatio: 50%

Network:
  TypeNode(User) ──► AlphaPass ──┐
                                 │
                            ┌────▼─────┐
  TypeNode(Order)──►AlphaPass──┤join_a3f8├──┬──► Terminal(r1)
                            └──────────┘  └──► Terminal(r2)


┌─────────────────────────────────────────────────────────────────────────┐
│  T3: Ajout Rule3 (jointure DIFFÉRENTE)                                 │
│      rule r3 : {e: Emp, d: Dept} / e.dept_id == d.id ==> action3()    │
└─────────────────────────────────────────────────────────────────────────┘

1️⃣ Calculer hash: hash("e.dept_id == d.id", ["e"], ["d"]) = "join_f7c2"
2️⃣ BetaSharingRegistry.GetOrCreate("join_f7c2")
   → Pas trouvé → Créer nouveau JoinNode("join_f7c2")
3️⃣ LifecycleManager.RegisterNode("join_f7c2", "join")
4️⃣ LifecycleManager.AddRuleToNode("join_f7c2", "r3")

État:
  BetaSharingRegistry: {
    "join_a3f8": JoinNode (RefCount=2),
    "join_f7c2": JoinNode (RefCount=1) ⬅️ NOUVEAU
  }
  LifecycleManager: {
    "join_a3f8": NodeLifecycle { RefCount: 2, Rules: ["r1", "r2"] },
    "join_f7c2": NodeLifecycle { RefCount: 1, Rules: ["r3"] }
  }
  Metrics:
    TotalCreated: 2
    TotalReused: 1
    SharingRatio: 33%
```

### Suppression Progressive de Règles

```
┌─────────────────────────────────────────────────────────────────────────┐
│  T4: Suppression Rule2                                                  │
└─────────────────────────────────────────────────────────────────────────┘

1️⃣ LifecycleManager.RemoveRuleFromNode("join_a3f8", "r2")
   → Décrémenter RefCount: 2 → 1
   → shouldDelete = false (RefCount > 0)

2️⃣ Supprimer Terminal(r2) uniquement
3️⃣ CONSERVER JoinNode("join_a3f8") ← Encore utilisé par r1

État:
  BetaSharingRegistry: {
    "join_a3f8": JoinNode (RefCount=1) ⬇️,
    "join_f7c2": JoinNode (RefCount=1)
  }
  LifecycleManager: {
    "join_a3f8": NodeLifecycle { 
      RefCount: 1 ⬇️, 
      Rules: ["r1"]  ⬅️ r2 retiré
    },
    "join_f7c2": NodeLifecycle { RefCount: 1, Rules: ["r3"] }
  }

Network:
  TypeNode(User) ──► AlphaPass ──┐
                                 │
                            ┌────▼─────┐
  TypeNode(Order)──►AlphaPass──┤join_a3f8├──► Terminal(r1) ✅
                            └──────────┘


┌─────────────────────────────────────────────────────────────────────────┐
│  T5: Suppression Rule1 (dernière règle utilisant join_a3f8)            │
└─────────────────────────────────────────────────────────────────────────┘

1️⃣ LifecycleManager.RemoveRuleFromNode("join_a3f8", "r1")
   → Décrémenter RefCount: 1 → 0
   → shouldDelete = true ✅

2️⃣ BetaSharingRegistry.RemoveJoinNode("join_a3f8")
   → Supprimer du registre

3️⃣ Network.BetaNodes: delete("join_a3f8")

4️⃣ LifecycleManager.RemoveNode("join_a3f8")

État:
  BetaSharingRegistry: {
    "join_f7c2": JoinNode (RefCount=1)
  }
  LifecycleManager: {
    "join_f7c2": NodeLifecycle { RefCount: 1, Rules: ["r3"] }
  }

Network:
  TypeNode(User) ──► (AlphaPass supprimé)
                            
  TypeNode(Order)──► (AlphaPass supprimé)
                            
  ⚠️ join_a3f8 COMPLÈTEMENT SUPPRIMÉ ⚠️
```

---

## Comparaison Alpha vs Beta

### AlphaNodes: Partage Simple (1 Variable)

```
Rule 1: {p: Person} / p.age > 18 ==> action1()
Rule 2: {p: Person} / p.age > 18 ==> action2()
Rule 3: {p: Person} / p.age > 18 AND p.city == "Paris" ==> action3()

┌────────────────────────────────────────────────────────────────────┐
│                   Alpha Node Sharing                               │
└────────────────────────────────────────────────────────────────────┘

                      ┌──────────────┐
                      │ TypeNode(P)  │
                      └──────┬───────┘
                             │
                   ┏━━━━━━━━━▼━━━━━━━━━┓
                   ┃ SHARED AlphaNode  ┃
                   ┃ p.age > 18        ┃
                   ┃ RefCount=3        ┃
                   ┗━━━━━━━━━┳━━━━━━━━━┛
                             │
                 ┌───────────┼───────────┐
                 │           │           │
          ┌──────▼───┐       │      ┌────▼─────┐
          │Terminal1 │       │      │AlphaNode │
          └──────────┘       │      │city="Paris"│
                             │      └────┬─────┘
                      ┌──────▼───┐       │
                      │Terminal2 │  ┌────▼─────┐
                      └──────────┘  │Terminal3 │
                                    └──────────┘

Hash Basé Sur:
  • Condition (p.age > 18)
  • Variable name ("p")

Partage: 3 règles partagent 1 AlphaNode
```

### BetaNodes: Partage Complexe (2+ Variables)

```
Rule 1: {u: User, o: Order} / o.user_id == u.id ==> action1()
Rule 2: {u: User, o: Order} / o.user_id == u.id ==> action2()

┌────────────────────────────────────────────────────────────────────┐
│                   Beta Node Sharing                                │
└────────────────────────────────────────────────────────────────────┘

TypeNode(User) ──► AlphaPass [LEFT] ──┐
                                      │
                                ┏━━━━━▼━━━━━━━━━━━━┓
                                ┃ SHARED JoinNode  ┃
TypeNode(Order)──►AlphaPass ────┃ o.user_id==u.id ┃
                   [RIGHT]      ┃ RefCount=2       ┃
                                ┃                  ┃
                                ┃ 2 ENTRÉES:       ┃
                                ┃  • LEFT: User    ┃
                                ┃  • RIGHT: Order  ┃
                                ┗━━━━━━┳━━━━━━━━━━━┛
                                       │
                                ┌──────┴──────┐
                                │             │
                         ┌──────▼───┐  ┌──────▼───┐
                         │Terminal1 │  │Terminal2 │
                         └──────────┘  └──────────┘

Hash Basé Sur:
  • Condition (o.user_id == u.id)
  • LeftVariables (["u"])
  • RightVariables (["o"])
  • VariableTypes ({"u":"User", "o":"Order"})

Complexité Additionnelle:
  ✓ 2 points d'entrée (left/right)
  ✓ 3 mémoires (Left, Right, Result)
  ✓ Gestion des tokens (pas juste faits)
  ✓ Coordination bidirectionnelle
```

### Comparaison Détaillée

```
┌──────────────────┬─────────────────────┬─────────────────────────┐
│   Aspect         │    AlphaNodes       │      BetaNodes          │
├──────────────────┼─────────────────────┼─────────────────────────┤
│ Entrées          │ 1 (TypeNode)        │ 2 (Left + Right)        │
│ Variables        │ 1 (ex: "p")         │ 2+ (ex: "u", "o")       │
│ Propagation      │ Faits               │ Tokens (multi-faits)    │
│ Mémoire          │ Simple (Facts map)  │ Triple (L/R/Result)     │
│ Hash             │ condition + var     │ condition + vars + types│
│ Complexité Hash  │ Faible              │ Moyenne                 │
│ Partage Potentiel│ Très Élevé (70-85%) │ Élevé (30-70%)          │
│ Bénéfice/Nœud    │ Moyen               │ Élevé (plus coûteux)    │
└──────────────────┴─────────────────────┴─────────────────────────┘
```

---

## Performance et Métriques

### Dashboard Visuel des Métriques

```
┌─────────────────────────────────────────────────────────────────────┐
│             BETA NODE SHARING METRICS DASHBOARD                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│  Total JoinNodes Created:  47                                      │
│  Total JoinNodes Reused:   103                                     │
│  Sharing Ratio:            68.7% ████████████████░░░░               │
│                                                                     │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  Sharing Distribution                                        │  │
│  │                                                              │  │
│  │  RefCount=1: ████████  (20 nodes) - 42.6%                   │  │
│  │  RefCount=2: ███████████ (15 nodes) - 31.9%                 │  │
│  │  RefCount=3: ████ (7 nodes) - 14.9%                         │  │
│  │  RefCount=4+: ██ (5 nodes) - 10.6%                          │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                     │
│  Hash Cache Stats:                                                 │
│  ├─ Size: 8,523 / 10,000 (85.2% full)                             │
│  ├─ Hits: 42,891                                                   │
│  ├─ Misses: 7,109                                                  │
│  ├─ Hit Rate: 85.8% ████████████████░░░                            │
│  └─ Evictions: 1,234                                               │
│                                                                     │
│  Memory Usage:                                                     │
│  ├─ Without Sharing (estimated): 235 MB                           │
│  ├─ With Sharing (actual):       147 MB                           │
│  └─ Reduction: 37.4% ████████████████████████                      │
│                                                                     │
│  Build Time:                                                       │
│  ├─ Average per JoinNode: 38.2 µs                                 │
│  ├─ Total for 150 rules: 5.73 ms                                  │
│  └─ Improvement vs No-Cache: 41.3% faster                         │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

### Évolution Temporelle

```
┌─────────────────────────────────────────────────────────────────────┐
│  Sharing Ratio Over Time (as rules are added)                      │
└─────────────────────────────────────────────────────────────────────┘

Sharing 
Ratio
  100% │                                    ╭────────────
       │                              ╭────╯
   80% │                        ╭────╯
       │                  ╭────╯
   60% │            ╭────╯
       │      ╭────╯
   40% │ ╭───╯
       │╯
   20% │
       │
    0% └─────────────────────────────────────────────────────────►
       0   10   20   30   40   50   60   70   80   90  100  Rules

Observation: Sharing ratio augmente avec le nombre de règles
             (plus de règles = plus de patterns communs)
```

---

**Fin des Diagrammes d'Architecture**

*Ce document complète l'analyse technique BETA_NODES_ANALYSIS.md avec des représentations visuelles.*