# Archives TSD

Ce répertoire contient les fichiers et documents archivés du projet TSD qui ne sont plus utilisés activement mais conservés pour référence historique.

## Structure

```
ARCHIVES/
├── README.md          # Ce fichier
└── sessions/          # Rapports de sessions de debugging/développement
```

## Contenu Archivé

### Sessions de Développement (`sessions/`)

Rapports et documents de sessions de debugging et développement archivés :

#### Bindings & Cascade
- `BINDINGS_IMPLEMENTATION_REPORT.md` - Rapport d'implémentation des bindings
- `DEBUG_BINDINGS_FINAL_REPORT.md` - Rapport final de debugging bindings
- `RESUME_DEBUG_BINDINGS_FR.md` - Résumé debugging bindings
- `TODO_CASCADE_BINDINGS_FIX.md` - TODO fix cascade bindings
- `TODO_DEBUG_E2E_BINDINGS.md` - TODO tests E2E bindings
- `TODO_FIX_BINDINGS_3_VARIABLES.md` - TODO fix bindings 3 variables

#### Sessions & Rapports
- `SESSION_06_COMMANDS.md` - Commandes session 06
- `SESSION_12_SUMMARY.md` - Résumé session 12
- `FINAL_SESSION_12_REPORT.md` - Rapport final session 12
- `SESSION_DEBUG_BINDINGS_REPORT.md` - Rapport debugging bindings
- `SESSION_DEBUG_SUMMARY.md` - Résumé debugging
- `SESSION_DEBUG_FINAL_SUMMARY.md` - Résumé final debugging
- `SESSION_REVIEW_COMPLETE.md` - Review complète
- `SESSION_REVIEW_SUMMARY.md` - Résumé review

#### Résolutions & Fixes
- `RESOLUTION_TESTS_E2E.md` - Résolution tests E2E
- `FIXES_SUMMARY.md` - Résumé des corrections
- `REFACTORING_REPORT.md` - Rapport de refactoring
- `FILES_CHANGED.md` - Fichiers modifiés

#### Validations & Synthèses
- `VALIDATION_FINALE_POST_FIX.md` - Validation finale post-fix
- `SYNTHESE_VALIDATION_FINALE.md` - Synthèse validation finale
- `VALIDATION_SUMMARY.txt` - Résumé validation
- `MAINTENANCE_SUMMARY.txt` - Résumé maintenance
- `REVIEW_SUMMARY.txt` - Résumé review

#### Refactoring
- `CHANGELOG_PERFORM_JOIN.md` - Changelog perform join
- `refactoring_perform_join_tokens.md` - Refactoring perform join tokens

## Pourquoi Archiver ?

Ces documents ont été archivés pour les raisons suivantes :

1. **Temporalité** : Documents de sessions ponctuelles de debugging/développement
2. **Obsolescence** : Problèmes résolus, TODOs complétés
3. **Consolidation** : Informations intégrées dans la documentation principale
4. **Historique** : Conservés pour référence et traçabilité

## Documentation Active

La documentation active et à jour se trouve dans :

- `/docs/` - Documentation principale
  - `/docs/configuration/` - Configuration système
  - `/docs/api/` - API publique
  - `/docs/architecture/` - Architecture et design
  - `/docs/guides/` - Guides utilisateur
- `/README.md` - README principal
- `/CHANGELOG.md` - Changelog actif
- `/TODO_ACTIFS.md` - TODOs actifs

## Accès aux Archives

Ces documents sont conservés uniquement pour référence historique. Pour les problèmes actuels ou la documentation à jour, consultez :

1. **Documentation** : `/docs/README.md`
2. **Configuration** : `/docs/configuration/README.md`
3. **TODOs actifs** : `/TODO_ACTIFS.md`
4. **Issues** : GitHub Issues

## Politique d'Archivage

Les documents sont archivés lorsque :

- ✅ Le problème/TODO est résolu et validé
- ✅ L'information est consolidée dans la doc principale
- ✅ Le document n'est plus pertinent pour le développement actuel
- ✅ Conservation requise pour traçabilité historique

Les documents **ne sont PAS** archivés s'ils :

- ❌ Contiennent des informations encore pertinentes
- ❌ Sont référencés par la documentation active
- ❌ Décrivent des problèmes non résolus

## Script d'Archivage

Pour archiver d'autres documents obsolètes :

```bash
# Utiliser le script fourni
./scripts/archive_obsolete_todos.sh <fichier_à_archiver>

# Ou manuellement
mv <fichier> ARCHIVES/sessions/
# Mettre à jour ce README
```

---

**Date de création** : 2025-01-XX  
**Dernière mise à jour** : 2025-01-XX  
**Mainteneur** : TSD Contributors