#!/bin/bash
# Script pour archiver les fichiers TODO obsolètes
# Suite à la correction du bug de partage JoinNode, plusieurs fichiers TODO sont maintenant obsolètes

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
ARCHIVE_DIR="$PROJECT_ROOT/REPORTS/ARCHIVE"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

echo "🗂️  Script d'Archivage des TODOs Obsolètes"
echo "=========================================="
echo ""
echo "Répertoire projet: $PROJECT_ROOT"
echo "Répertoire archive: $ARCHIVE_DIR"
echo ""

# Créer le répertoire d'archive s'il n'existe pas
if [ ! -d "$ARCHIVE_DIR" ]; then
    echo "📁 Création du répertoire d'archive..."
    mkdir -p "$ARCHIVE_DIR"
fi

# Liste des fichiers TODO obsolètes (bug bindings résolu)
OBSOLETE_FILES=(
    "TODO_BINDINGS_CASCADE.md"
    "TODO_CASCADE_BINDINGS_FIX.md"
    "TODO_DEBUG_E2E_BINDINGS.md"
    "TODO_FIX_BINDINGS_3_VARIABLES.md"
)

echo "📋 Fichiers à archiver:"
for file in "${OBSOLETE_FILES[@]}"; do
    echo "   - $file"
done
echo ""

# Compteurs
moved=0
not_found=0
failed=0

# Archiver chaque fichier
for file in "${OBSOLETE_FILES[@]}"; do
    filepath="$PROJECT_ROOT/$file"

    if [ -f "$filepath" ]; then
        # Créer un nom avec timestamp pour éviter les collisions
        archive_name="${file%.md}_archived_${TIMESTAMP}.md"
        archive_path="$ARCHIVE_DIR/$archive_name"

        echo "📦 Archivage de $file..."

        # Ajouter un header au fichier archivé
        {
            echo "# ARCHIVED - $(date '+%Y-%m-%d %H:%M:%S')"
            echo ""
            echo "**Raison:** Bug de partage JoinNode résolu - Ce TODO n'est plus pertinent"
            echo "**Référence:** Voir VALIDATION_FINALE_POST_FIX.md"
            echo ""
            echo "---"
            echo ""
            cat "$filepath"
        } > "$archive_path"

        if [ $? -eq 0 ]; then
            # Supprimer l'original
            rm "$filepath"
            echo "   ✅ Archivé vers: $archive_name"
            ((moved++))
        else
            echo "   ❌ Erreur lors de l'archivage"
            ((failed++))
        fi
    else
        echo "   ⚠️  Fichier non trouvé: $file"
        ((not_found++))
    fi
done

echo ""
echo "=========================================="
echo "📊 Résumé:"
echo "   ✅ Archivés:    $moved"
echo "   ⚠️  Non trouvés: $not_found"
echo "   ❌ Échecs:      $failed"
echo ""

# Créer un fichier index dans l'archive
INDEX_FILE="$ARCHIVE_DIR/INDEX.md"
echo "📝 Création de l'index d'archive..."

{
    echo "# Index des Fichiers Archivés"
    echo ""
    echo "**Dernière mise à jour:** $(date '+%Y-%m-%d %H:%M:%S')"
    echo ""
    echo "## Fichiers TODOs Bindings (Bug Résolu)"
    echo ""
    echo "Ces fichiers concernaient le bug de partage JoinNode qui a été résolu."
    echo "Voir \`VALIDATION_FINALE_POST_FIX.md\` pour les détails."
    echo ""
    echo "| Fichier Original | Archive | Date |"
    echo "|------------------|---------|------|"

    for file in "${OBSOLETE_FILES[@]}"; do
        archive_name="${file%.md}_archived_${TIMESTAMP}.md"
        if [ -f "$ARCHIVE_DIR/$archive_name" ]; then
            echo "| $file | $archive_name | $(date '+%Y-%m-%d') |"
        fi
    done

    echo ""
    echo "## Contexte"
    echo ""
    echo "Suite à la correction du bug de partage des JoinNodes (prefix sharing / beta sharing),"
    echo "ces fichiers TODO ne sont plus nécessaires car:"
    echo ""
    echo "1. Le bug a été identifié et corrigé"
    echo "2. Tests de régression ajoutés"
    echo "3. Tous les tests passent (100%)"
    echo "4. Documentation mise à jour"
    echo ""
    echo "## Fichiers Modifiés Lors de la Correction"
    echo ""
    echo "- \`rete/beta_chain_optimizer.go\`"
    echo "- \`rete/beta_sharing_interface.go\`"
    echo "- \`rete/beta_sharing.go\`"
    echo "- \`rete/beta_sharing_hash.go\`"
    echo "- \`rete/beta_chain_builder_orchestration.go\`"
    echo "- Tests de régression ajoutés dans \`rete/\`"
    echo ""
    echo "## Références"
    echo ""
    echo "- \`VALIDATION_FINALE_POST_FIX.md\` - Rapport de validation complète"
    echo "- \`CHANGELOG.md\` - Historique des changements"
    echo "- Thread Zed: \"Immutable Bindings JoinNode Sharing Bug\""
    echo ""
} > "$INDEX_FILE"

echo "   ✅ Index créé: INDEX.md"
echo ""

# Créer un fichier TODO_ACTIFS.md pour les TODOs restants
ACTIVE_TODO="$PROJECT_ROOT/TODO_ACTIFS.md"
echo "📝 Création de TODO_ACTIFS.md..."

{
    echo "# TODOs Actifs - Liste Consolidée"
    echo ""
    echo "**Date:** $(date '+%Y-%m-%d')"
    echo "**Status:** Tous non-critiques (améliorations futures)"
    echo ""
    echo "## 📋 TODOs Non-Critiques (7 items)"
    echo ""
    echo "### 1. Migration ParseInput (Compatibilité)"
    echo ""
    echo "**Fichier:** \`constraint/cmd/main.go:248\`"
    echo "**Priorité:** Basse"
    echo "**Impact:** Aucun (wrapper de compatibilité)"
    echo ""
    echo '```go'
    echo "// TODO: Les tests doivent être migrés pour utiliser ParseInput au lieu de ParseFile"
    echo '```'
    echo ""
    echo "### 2. Validation Types Personnalisés"
    echo ""
    echo "**Fichier:** \`constraint/constraint_facts.go:71\`"
    echo "**Priorité:** Moyenne"
    echo "**Impact:** Extensibilité future"
    echo ""
    echo '```go'
    echo "// TODO: Valider que le type personnalisé existe dans le programme"
    echo '```'
    echo ""
    echo "### 3. Génération Table de Règles"
    echo ""
    echo "**Fichier:** \`constraint/parser.go:6034\`"
    echo "**Priorité:** Basse"
    echo "**Impact:** Performance marginale"
    echo ""
    echo '```go'
    echo "// TODO : not super critical but this could be generated"
    echo '```'
    echo ""
    echo "### 4. Support Opérateur Modulo (%)"
    echo ""
    echo "**Fichier:** \`rete/arithmetic_alpha_extraction_test.go:317\`"
    echo "**Priorité:** Moyenne"
    echo "**Impact:** Feature manquante (test commenté)"
    echo ""
    echo '```go'
    echo "// TODO: Enable when parser supports % operator"
    echo '```'
    echo ""
    echo "### 5. Comparaison Profonde Conditions"
    echo ""
    echo "**Fichier:** \`rete/beta_sharing_interface.go:444\`"
    echo "**Priorité:** Moyenne"
    echo "**Impact:** Partage plus agressif (optimisation)"
    echo ""
    echo '```go'
    echo "// TODO: Deep comparison of normalized conditions"
    echo '```'
    echo ""
    echo "### 6. Métriques Détaillées JoinNode"
    echo ""
    echo "**Fichier:** \`rete/beta_sharing_stats.go:135-136\`"
    echo "**Priorité:** Basse"
    echo "**Impact:** Observabilité (non-bloquant)"
    echo ""
    echo '```go'
    echo "CreatedAt:        time.Time{}, // TODO: Track creation time"
    echo "ActivationCount:  0,            // TODO: Track activation count"
    echo '```'
    echo ""
    echo "### 7. AlphaConditionEvaluator Arithmétique"
    echo ""
    echo "**Fichier:** \`rete/condition_splitter.go:86\`"
    echo "**Priorité:** Moyenne"
    echo "**Impact:** Performance (conditions arithmétiques en alpha)"
    echo ""
    echo '```go'
    echo "// TODO: Enhance AlphaConditionEvaluator to handle arithmetic"
    echo '```'
    echo ""
    echo "---"
    echo ""
    echo "## 📊 Statistiques"
    echo ""
    echo "- **Total TODOs:** 7"
    echo "- **Priorité Haute:** 0"
    echo "- **Priorité Moyenne:** 4"
    echo "- **Priorité Basse:** 3"
    echo ""
    echo "## 🎯 Recommandations"
    echo ""
    echo "1. **Court terme:** Aucune action requise (tous non-bloquants)"
    echo "2. **Moyen terme:** Prioriser les TODOs moyenne priorité selon roadmap"
    echo "3. **Long terme:** Implémenter les optimisations et features manquantes"
    echo ""
    echo "## 📝 Notes"
    echo ""
    echo "- Tous les TODOs critiques ont été résolus"
    echo "- Le système est production ready"
    echo "- Ces TODOs sont des améliorations futures optionnelles"
    echo ""
    echo "---"
    echo ""
    echo "*Dernière révision: $(date '+%Y-%m-%d %H:%M:%S')*"
} > "$ACTIVE_TODO"

echo "   ✅ TODO_ACTIFS.md créé"
echo ""

echo "=========================================="
echo "✅ Archivage terminé avec succès!"
echo ""
echo "📁 Fichiers archivés dans: $ARCHIVE_DIR"
echo "📋 Index créé: $ARCHIVE_DIR/INDEX.md"
echo "📝 TODOs actifs: TODO_ACTIFS.md"
echo ""
echo "🎯 Prochaines étapes:"
echo "   1. Réviser TODO_ACTIFS.md"
echo "   2. Vérifier VALIDATION_FINALE_POST_FIX.md"
echo "   3. Lancer 'make test' pour validation finale"
echo ""
