# Plan de Refactorisation Structure Packages Go

## 🎯 Objectif
Réorganiser le projet selon les conventions Go standards pour améliorer la maintenabilité et la lisibilité.

## 📋 État Actuel
```
/
├── cmd/main.go (CLI principal)
├── constraint/
│   ├── cmd/main.go (CLI parsing debug)
│   ├── pkg/
│   ├── internal/
│   └── *.go (fonctions publiques)
├── rete/
│   ├── cmd/main.go (CLI benchmark)
│   ├── cmd/monitoring/main.go
│   ├── pkg/
│   ├── internal/
│   └── *.go (fonctions publiques)
└── test/ (tests intégration)
```

## 🎯 Structure Cible (Conventions Go)
```
/
├── cmd/
│   ├── tsd/main.go (CLI principal unifié)
│   ├── constraint-debug/main.go (debug parsing)
│   ├── rete-benchmark/main.go (benchmark)
│   └── rete-monitor/main.go (monitoring)
├── pkg/
│   ├── constraint/
│   │   ├── parser.go
│   │   ├── types.go
│   │   ├── api.go
│   │   └── validator/
│   └── rete/
│       ├── network.go
│       ├── nodes/
│       └── storage/
├── internal/
│   ├── config/
│   └── utils/
└── test/ (tests intégration)
```

## 📝 Actions à Effectuer

### 1. Réorganisation des commandes (cmd/)
- [x] Analyser les commandes existantes
- [ ] Créer cmd/tsd/ (CLI principal)
- [ ] Créer cmd/constraint-debug/ (outil debug)
- [ ] Créer cmd/rete-benchmark/ (tests performance)
- [ ] Créer cmd/rete-monitor/ (interface monitoring)
- [ ] Supprimer anciennes commandes dispersées

### 2. Migration vers pkg/ (code public)
- [ ] Déplacer constraint/*.go → pkg/constraint/
- [ ] Déplacer rete/*.go → pkg/rete/
- [ ] Préserver les APIs publiques
- [ ] Adapter les imports

### 3. Consolidation internal/
- [ ] Fusionner constraint/internal/ et rete/internal/
- [ ] Créer internal/config/ unifié
- [ ] Créer internal/utils/ pour code partagé

### 4. Mise à jour des imports
- [ ] Adapter tous les imports (*.go)
- [ ] Mettre à jour tests/
- [ ] Vérifier go.mod

### 5. Tests et validation
- [ ] Exécuter tous les tests
- [ ] Vérifier build complet
- [ ] Valider APIs publiques inchangées

## 🔄 Ordre d'Exécution
1. Créer nouvelle structure (cmd/, pkg/)
2. Copier et adapter les fichiers
3. Mettre à jour imports
4. Tests et validation
5. Suppression ancienne structure

## ⚠️ Points d'Attention
- Ne pas casser l'API publique existante
- Maintenir tous les tests fonctionnels
- Préserver la compatibilité backward
- Documenter les changements
