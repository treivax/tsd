#!/bin/bash

echo "🚀 LANCEMENT DES TESTS DE COUVERTURE DÉTAILLÉS"
echo "=============================================="
echo

echo "📊 Lancement des tests Alpha..."
cd /home/resinsec/dev/tsd/test/coverage/alpha
timeout 60s go run alpha_detailed_runner.go > alpha_run.log 2>&1 &
ALPHA_PID=$!

echo "📊 Lancement des tests Beta..."
cd /home/resinsec/dev/tsd/test/coverage/beta
timeout 60s go run beta_detailed_runner.go > beta_run.log 2>&1 &
BETA_PID=$!

echo "⏳ Attente des tests Alpha..."
wait $ALPHA_PID
ALPHA_STATUS=$?

echo "⏳ Attente des tests Beta..."
wait $BETA_PID
BETA_STATUS=$?

echo
echo "📋 RÉSULTATS:"
echo "============"

if [ $ALPHA_STATUS -eq 0 ]; then
    echo "✅ Tests Alpha: SUCCÈS"
else
    echo "❌ Tests Alpha: ÉCHEC (code $ALPHA_STATUS)"
fi

if [ $BETA_STATUS -eq 0 ]; then
    echo "✅ Tests Beta: SUCCÈS"
else
    echo "❌ Tests Beta: ÉCHEC (code $BETA_STATUS)"
fi

echo
echo "📄 Rapports générés:"
echo "- Alpha: /home/resinsec/dev/tsd/ALPHA_NODES_DETAILED_RESULTS.md"
echo "- Beta: /home/resinsec/dev/tsd/BETA_NODES_DETAILED_RESULTS.md"

echo
echo "📊 Affichage des logs Alpha (dernières 10 lignes):"
echo "================================================="
tail -10 /home/resinsec/dev/tsd/test/coverage/alpha/alpha_run.log

echo
echo "📊 Affichage des logs Beta (dernières 10 lignes):"
echo "================================================="
tail -10 /home/resinsec/dev/tsd/test/coverage/beta/beta_run.log
