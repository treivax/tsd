#!/bin/bash

# Script de démonstration du système de monitoring RETE
echo "=== Démonstration du Système de Monitoring RETE ==="
echo

# Construire le projet
echo "🔨 Construction du projet..."
cd /home/resinsec/dev/tsd/rete
go build -o monitoring-demo ./cmd/monitoring

if [ $? -ne 0 ]; then
    echo "❌ Erreur lors de la construction"
    exit 1
fi

echo "✅ Construction réussie"
echo

# Démarrer le serveur de monitoring en arrière-plan
echo "🚀 Démarrage du serveur de monitoring..."
./monitoring-demo &
MONITOR_PID=$!

# Attendre que le serveur démarre
sleep 3

echo "✅ Serveur de monitoring démarré (PID: $MONITOR_PID)"
echo "📊 Interface web disponible à: http://localhost:8080"
echo

# Instructions pour l'utilisateur
echo "=== Instructions ==="
echo "1. Ouvrez votre navigateur web"
echo "2. Allez à: http://localhost:8080"
echo "3. Explorez les différents onglets du dashboard:"
echo "   - 📈 Métriques Globales"
echo "   - 🔧 Composants Optimisés"
echo "   - 🎯 Performance"
echo "   - 🚨 Alertes"
echo "4. Observez les métriques en temps réel"
echo

echo "⏱️  Le serveur tournera pendant 60 secondes..."
echo "🔄 Données simulées en cours de génération..."

# Attendre 60 secondes
sleep 60

# Arrêter le serveur
echo
echo "🛑 Arrêt du serveur de monitoring..."
kill $MONITOR_PID
wait $MONITOR_PID 2>/dev/null

echo "✅ Démonstration terminée"
echo

# Nettoyer
rm -f monitoring-demo

echo "=== Résumé des Fonctionnalités Implémentées ==="
echo "✅ Serveur HTTP avec API REST"
echo "✅ Interface web responsive avec Chart.js"
echo "✅ WebSocket pour mise à jour temps réel"
echo "✅ Collecte de métriques des composants optimisés"
echo "✅ Dashboard avec visualisations interactives"
echo "✅ Système d'alertes configurable"
echo "✅ Métriques de performance et tendances"
echo "✅ Intégration complète avec RETE"
echo
echo "📁 Fichiers créés:"
echo "   - rete/monitoring_server.go (serveur principal)"
echo "   - rete/metrics_integrator.go (collecte de métriques)"
echo "   - rete/monitored_network.go (réseau RETE monitoré)"
echo "   - rete/web/ (interface web complète)"
echo "   - rete/cmd/monitoring/ (exemple d'utilisation)"
echo
echo "🎉 Interface de monitoring en temps réel complètement implémentée !"