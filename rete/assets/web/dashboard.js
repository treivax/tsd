// RETE Monitoring Dashboard JavaScript

class RETEDashboard {
    constructor() {
        this.ws = null;
        this.wsConnected = false;
        this.charts = {};
        this.metricsHistory = [];
        this.maxHistoryPoints = 50;
        this.currentTab = 'overview';

        this.init();
    }

    init() {
        this.setupTabNavigation();
        this.setupCharts();
        this.setupWebSocket();
        this.setupEventListeners();
        this.setupModals();
        this.startDataRefresh();

        console.log('🚀 RETE Dashboard initialized');
    }

    // Tab Navigation
    setupTabNavigation() {
        const navButtons = document.querySelectorAll('.nav-btn');
        const tabContents = document.querySelectorAll('.tab-content');

        navButtons.forEach(btn => {
            btn.addEventListener('click', () => {
                const tabId = btn.dataset.tab;

                // Update navigation
                navButtons.forEach(b => b.classList.remove('active'));
                btn.classList.add('active');

                // Update content
                tabContents.forEach(content => {
                    content.classList.remove('active');
                });
                document.getElementById(tabId).classList.add('active');

                this.currentTab = tabId;
                this.onTabChange(tabId);
            });
        });
    }

    onTabChange(tabId) {
        switch(tabId) {
            case 'overview':
                this.updateOverviewTab();
                break;
            case 'performance':
                this.updatePerformanceTab();
                break;
            case 'network':
                this.updateNetworkTab();
                break;
            case 'alerts':
                this.updateAlertsTab();
                break;
            case 'system':
                this.updateSystemTab();
                break;
        }
    }

    // Charts Setup
    setupCharts() {
        // Vérifier que Chart.js est disponible
        if (typeof Chart === 'undefined') {
            console.error('❌ Chart.js not loaded');
            return;
        }

        // Configuration par défaut pour tous les graphiques
        Chart.defaults.responsive = true;
        Chart.defaults.maintainAspectRatio = false;
        Chart.defaults.animation = {
            duration: 750,
            easing: 'easeInOutQuart'
        };

        // Throughput Chart
        const throughputCtx = document.getElementById('throughputChart');
        if (throughputCtx) {
            this.charts.throughput = new Chart(throughputCtx, {
                type: 'line',
                data: {
                    labels: this.generateTimeLabels(10), // Générer des labels initiaux
                    datasets: [
                        {
                            label: 'Facts/sec',
                            data: new Array(10).fill(0), // Données initiales à zéro
                            borderColor: '#2563eb',
                            backgroundColor: 'rgba(37, 99, 235, 0.1)',
                            tension: 0.4,
                            fill: false
                        },
                        {
                            label: 'Tokens/sec',
                            data: new Array(10).fill(0),
                            borderColor: '#10b981',
                            backgroundColor: 'rgba(16, 185, 129, 0.1)',
                            tension: 0.4,
                            fill: false
                        },
                        {
                            label: 'Rules/sec',
                            data: new Array(10).fill(0),
                            borderColor: '#f59e0b',
                            backgroundColor: 'rgba(245, 158, 11, 0.1)',
                            tension: 0.4,
                            fill: false
                        }
                    ]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    interaction: {
                        intersect: false,
                        mode: 'index'
                    },
                    scales: {
                        y: {
                            beginAtZero: true,
                            title: {
                                display: true,
                                text: 'Operations/sec'
                            },
                            grid: {
                                color: '#e2e8f0'
                            }
                        },
                        x: {
                            title: {
                                display: true,
                                text: 'Time'
                            },
                            grid: {
                                color: '#e2e8f0'
                            }
                        }
                    },
                    plugins: {
                        legend: {
                            display: true,
                            position: 'top'
                        },
                        tooltip: {
                            backgroundColor: 'rgba(0, 0, 0, 0.8)',
                            titleColor: 'white',
                            bodyColor: 'white'
                        }
                    }
                }
            });
        }

        // Latency Chart
        const latencyCtx = document.getElementById('latencyChart');
        if (latencyCtx) {
            this.charts.latency = new Chart(latencyCtx, {
                type: 'bar',
                data: {
                    labels: ['P50', 'P75', 'P90', 'P95', 'P99'],
                    datasets: [{
                        label: 'Latency (ms)',
                        data: [0, 0, 0, 0, 0],
                        backgroundColor: [
                            '#10b981',
                            '#3b82f6',
                            '#f59e0b',
                            '#ef4444',
                            '#7c3aed'
                        ],
                        borderRadius: 4
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        y: {
                            beginAtZero: true,
                            title: {
                                display: true,
                                text: 'Latency (ms)'
                            },
                            grid: {
                                color: '#e2e8f0'
                            }
                        },
                        x: {
                            grid: {
                                display: false
                            }
                        }
                    },
                    plugins: {
                        legend: {
                            display: false
                        },
                        tooltip: {
                            backgroundColor: 'rgba(0, 0, 0, 0.8)',
                            titleColor: 'white',
                            bodyColor: 'white'
                        }
                    }
                }
            });
        }

        // Performance Chart
        const performanceCtx = document.getElementById('performanceChart');
        if (performanceCtx) {
            this.charts.performance = new Chart(performanceCtx, {
                type: 'radar',
                data: {
                    labels: ['IndexedStorage', 'HashJoin', 'EvalCache', 'TokenProp'],
                    datasets: [{
                        label: 'Performance Score',
                        data: [50, 50, 50, 50], // Valeurs initiales neutres
                        borderColor: '#2563eb',
                        backgroundColor: 'rgba(37, 99, 235, 0.2)',
                        pointBackgroundColor: '#2563eb',
                        pointBorderColor: '#fff',
                        pointBorderWidth: 2
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        r: {
                            beginAtZero: true,
                            max: 100,
                            grid: {
                                color: '#e2e8f0'
                            },
                            pointLabels: {
                                font: {
                                    size: 12
                                }
                            }
                        }
                    },
                    plugins: {
                        legend: {
                            display: false
                        }
                    }
                }
            });
        }

        // Cache Chart
        const cacheCtx = document.getElementById('cacheChart');
        if (cacheCtx) {
            this.charts.cache = new Chart(cacheCtx, {
                type: 'doughnut',
                data: {
                    labels: ['Hits', 'Misses'],
                    datasets: [{
                        data: [80, 20],
                        backgroundColor: ['#10b981', '#ef4444'],
                        borderWidth: 2,
                        borderColor: '#ffffff'
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    cutout: '60%',
                    plugins: {
                        legend: {
                            position: 'bottom',
                            labels: {
                                padding: 20
                            }
                        }
                    }
                }
            });
        }

        // Memory Chart
        const memoryCtx = document.getElementById('memoryChart');
        if (memoryCtx) {
            this.charts.memory = new Chart(memoryCtx, {
                type: 'line',
                data: {
                    labels: this.generateTimeLabels(10),
                    datasets: [{
                        label: 'Memory Usage (MB)',
                        data: new Array(10).fill(0),
                        borderColor: '#3b82f6',
                        backgroundColor: 'rgba(59, 130, 246, 0.1)',
                        tension: 0.4,
                        fill: true
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        y: {
                            beginAtZero: true,
                            title: {
                                display: true,
                                text: 'Memory (MB)'
                            },
                            grid: {
                                color: '#e2e8f0'
                            }
                        },
                        x: {
                            grid: {
                                color: '#e2e8f0'
                            }
                        }
                    },
                    plugins: {
                        legend: {
                            display: false
                        }
                    }
                }
            });
        }

        // Goroutines Chart
        const goroutinesCtx = document.getElementById('goroutinesChart');
        if (goroutinesCtx) {
            this.charts.goroutines = new Chart(goroutinesCtx, {
                type: 'line',
                data: {
                    labels: this.generateTimeLabels(10),
                    datasets: [{
                        label: 'Goroutines',
                        data: new Array(10).fill(0),
                        borderColor: '#10b981',
                        backgroundColor: 'rgba(16, 185, 129, 0.1)',
                        tension: 0.4,
                        fill: true
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        y: {
                            beginAtZero: true,
                            title: {
                                display: true,
                                text: 'Count'
                            },
                            grid: {
                                color: '#e2e8f0'
                            }
                        },
                        x: {
                            grid: {
                                color: '#e2e8f0'
                            }
                        }
                    },
                    plugins: {
                        legend: {
                            display: false
                        }
                    }
                }
            });
        }

        // Alert History Chart
        const alertHistoryCtx = document.getElementById('alertHistoryChart');
        if (alertHistoryCtx) {
            this.charts.alertHistory = new Chart(alertHistoryCtx, {
                type: 'line',
                data: {
                    labels: this.generateTimeLabels(10),
                    datasets: [
                        {
                            label: 'Critical',
                            data: new Array(10).fill(0),
                            borderColor: '#ef4444',
                            backgroundColor: 'rgba(239, 68, 68, 0.1)',
                            tension: 0.4,
                            fill: false
                        },
                        {
                            label: 'High',
                            data: new Array(10).fill(0),
                            borderColor: '#f59e0b',
                            backgroundColor: 'rgba(245, 158, 11, 0.1)',
                            tension: 0.4,
                            fill: false
                        },
                        {
                            label: 'Medium',
                            data: new Array(10).fill(0),
                            borderColor: '#3b82f6',
                            backgroundColor: 'rgba(59, 130, 246, 0.1)',
                            tension: 0.4,
                            fill: false
                        }
                    ]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    scales: {
                        y: {
                            beginAtZero: true,
                            title: {
                                display: true,
                                text: 'Alert Count'
                            },
                            grid: {
                                color: '#e2e8f0'
                            }
                        },
                        x: {
                            grid: {
                                color: '#e2e8f0'
                            }
                        }
                    },
                    plugins: {
                        legend: {
                            position: 'top'
                        }
                    }
                }
            });
        }

        console.log('📊 Charts initialized successfully');
    }

    // WebSocket Connection
    setupWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws/metrics`;

        this.connectWebSocket(wsUrl);
    }

    connectWebSocket(url) {
        try {
            this.ws = new WebSocket(url);

            this.ws.onopen = () => {
                console.log('🔌 WebSocket connected');
                this.wsConnected = true;
                this.updateConnectionStatus('connected', 'Connected');
            };

            this.ws.onmessage = (event) => {
                try {
                    const data = JSON.parse(event.data);
                    this.handleWebSocketMessage(data);
                } catch (error) {
                    console.error('❌ Error parsing WebSocket message:', error);
                }
            };

            this.ws.onclose = () => {
                console.log('🔌 WebSocket disconnected');
                this.wsConnected = false;
                this.updateConnectionStatus('disconnected', 'Disconnected');

                // Reconnect after 5 seconds
                setTimeout(() => {
                    this.connectWebSocket(url);
                }, 5000);
            };

            this.ws.onerror = (error) => {
                console.error('❌ WebSocket error:', error);
                this.updateConnectionStatus('error', 'Connection Error');
            };

        } catch (error) {
            console.error('❌ Failed to create WebSocket:', error);
            this.updateConnectionStatus('error', 'Connection Failed');
        }
    }

    handleWebSocketMessage(data) {
        switch (data.type) {
            case 'initial_data':
                this.handleInitialData(data.data);
                break;
            case 'metrics_update':
                this.handleMetricsUpdate(data.data);
                break;
            case 'alert':
                this.handleAlert(data.data);
                break;
            default:
                console.log('📊 Unknown message type:', data.type);
        }
    }

    handleInitialData(data) {
        console.log('📊 Initial data received');
        this.updateMetrics(data);

        if (data.history && data.history.length > 0) {
            this.metricsHistory = data.history;
            this.updateCharts();
        }
    }

    handleMetricsUpdate(data) {
        this.updateMetrics(data);
        this.addToHistory(data);
        this.updateCharts();
        this.updateTimestamp();
    }

    handleAlert(alert) {
        console.log('🚨 Alert received:', alert);
        this.showToast('Alert', alert.message, alert.severity);
        this.updateAlertsDisplay();
    }

    updateConnectionStatus(status, text) {
        const statusElement = document.getElementById('connectionStatus');
        const statusDot = statusElement.querySelector('.status-dot');
        const statusText = statusElement.querySelector('.status-text');

        statusDot.className = `status-dot ${status}`;
        statusText.textContent = text;
    }

    updateTimestamp() {
        const timestampElement = document.getElementById('lastUpdate');
        if (timestampElement) {
            timestampElement.textContent = `Last update: ${new Date().toLocaleTimeString()}`;
        }
    }

    // Event Listeners
    setupEventListeners() {
        // Network controls
        const refreshNetwork = document.getElementById('refreshNetwork');
        if (refreshNetwork) {
            refreshNetwork.addEventListener('click', () => this.refreshNetworkData());
        }

        const exportNetwork = document.getElementById('exportNetwork');
        if (exportNetwork) {
            exportNetwork.addEventListener('click', () => this.exportNetworkData());
        }

        // Alert controls
        const createAlert = document.getElementById('createAlert');
        if (createAlert) {
            createAlert.addEventListener('click', () => this.showCreateAlertModal());
        }

        const refreshAlerts = document.getElementById('refreshAlerts');
        if (refreshAlerts) {
            refreshAlerts.addEventListener('click', () => this.refreshAlertsData());
        }
    }

    // Modal Setup
    setupModals() {
        const alertModal = document.getElementById('alertRuleModal');
        const closeModal = alertModal.querySelector('.modal-close');
        const cancelBtn = document.getElementById('cancelRule');
        const alertForm = document.getElementById('alertRuleForm');

        closeModal.addEventListener('click', () => {
            alertModal.classList.remove('show');
        });

        cancelBtn.addEventListener('click', () => {
            alertModal.classList.remove('show');
        });

        alertForm.addEventListener('submit', (e) => {
            e.preventDefault();
            this.createAlertRule();
        });

        // Close modal on outside click
        alertModal.addEventListener('click', (e) => {
            if (e.target === alertModal) {
                alertModal.classList.remove('show');
            }
        });
    }

    // Data Updates
    updateMetrics(data) {
        if (data.system) {
            this.updateSystemMetrics(data.system);
        }

        if (data.rete) {
            this.updateReteMetrics(data.rete);
        }

        if (data.performance) {
            this.updatePerformanceMetrics(data.performance);
        }
    }

    updateSystemMetrics(system) {
        // Memory usage
        const memoryUsed = system.memory_usage_bytes || 0;
        const memoryTotal = system.memory_system_bytes || 1;
        const memoryPercent = Math.round((memoryUsed / memoryTotal) * 100);

        this.updateElement('memoryUsed', `width: ${memoryPercent}%`, 'style');
        this.updateElement('memoryText', `${this.formatBytes(memoryUsed)} / ${this.formatBytes(memoryTotal)}`);
        this.updateElement('memoryPercent', `${memoryPercent}%`);

        // Other system metrics
        this.updateElement('goroutineCount', system.goroutine_count || 0);
        this.updateElement('gcCount', system.gc_count || 0);
        this.updateElement('uptime', this.formatDuration(system.uptime_seconds || 0));
    }

    updateReteMetrics(rete) {
        // KPIs
        this.updateElement('factsPerSec', this.formatNumber(rete.facts_per_second || 0));
        this.updateElement('tokensPerSec', this.formatNumber(rete.tokens_per_second || 0));
        this.updateElement('rulesPerSec', this.formatNumber(rete.rules_per_second || 0));
        this.updateElement('avgLatency', `${this.formatNumber(rete.average_latency || 0)}ms`);

        // Network stats
        this.updateElement('totalNodes', rete.total_nodes || 0);
        this.updateElement('activeNodes', rete.active_nodes || 0);
        this.updateElement('totalFacts', this.formatNumber(rete.total_facts || 0));
        this.updateElement('errorRate', `${this.formatNumber(rete.error_rate || 0)}%`);
    }

    updatePerformanceMetrics(performance) {
        // IndexedStorage stats
        const indexStats = performance.indexed_storage_stats || {};
        this.updateElement('indexCacheHit', `${this.formatNumber(indexStats.cache_hit_ratio || 0)}%`);
        this.updateElement('totalIndexes', indexStats.total_indexes || 0);
        this.updateElement('lookupSpeed', `${this.formatNumber(indexStats.avg_lookup_time || 0)}ms`);

        // HashJoin stats
        const joinStats = performance.hash_join_stats || {};
        this.updateElement('joinCacheHits', joinStats.cache_hits || 0);
        this.updateElement('joinCacheMisses', joinStats.cache_misses || 0);
        this.updateElement('avgJoinTime', `${this.formatNumber(joinStats.avg_join_time || 0)}ms`);

        // Evaluation cache stats
        const evalStats = performance.evaluation_cache_stats || {};
        this.updateElement('evalCacheSize', evalStats.current_size || 0);
        this.updateElement('evalHitRatio', `${this.formatNumber(evalStats.hit_ratio || 0)}%`);
        this.updateElement('evalEvictions', evalStats.evictions || 0);

        // Token propagation stats
        const tokenStats = performance.token_propagation_stats || {};
        this.updateElement('queueSize', tokenStats.queue_size || 0);
        this.updateElement('parallelEfficiency', `${this.formatNumber(tokenStats.parallel_efficiency || 0)}%`);
        this.updateElement('workerUtilization', `${this.formatNumber(tokenStats.avg_worker_utilization || 0)}%`);
    }

    addToHistory(data) {
        this.metricsHistory.push({
            timestamp: new Date(),
            system: data.system,
            rete: data.rete,
            performance: data.performance
        });

        if (this.metricsHistory.length > this.maxHistoryPoints) {
            this.metricsHistory.shift();
        }
    }

    updateCharts() {
        this.updateThroughputChart();
        this.updateMemoryChart();
        this.updateGoroutinesChart();
        this.updatePerformanceRadarChart();
        this.updateCacheChart();
    }

    updateThroughputChart() {
        if (!this.charts.throughput || this.metricsHistory.length === 0) {
            return;
        }

        const labels = this.metricsHistory.map(h => h.timestamp.toLocaleTimeString());
        const factsData = this.metricsHistory.map(h => h.rete?.facts_per_second || 0);
        const tokensData = this.metricsHistory.map(h => h.rete?.tokens_per_second || 0);
        const rulesData = this.metricsHistory.map(h => h.rete?.rules_per_second || 0);

        this.charts.throughput.data.labels = labels;
        this.charts.throughput.data.datasets[0].data = factsData;
        this.charts.throughput.data.datasets[1].data = tokensData;
        this.charts.throughput.data.datasets[2].data = rulesData;

        // Utiliser 'none' pour éviter les animations qui peuvent causer des problèmes de taille
        this.charts.throughput.update('none');
    }

    updateMemoryChart() {
        if (!this.charts.memory || this.metricsHistory.length === 0) {
            return;
        }

        const labels = this.metricsHistory.map(h => h.timestamp.toLocaleTimeString());
        const memoryData = this.metricsHistory.map(h =>
            (h.system?.memory_usage_bytes || 0) / (1024 * 1024) // Convert to MB
        );

        this.charts.memory.data.labels = labels;
        this.charts.memory.data.datasets[0].data = memoryData;
        this.charts.memory.update('none');
    }

    updateGoroutinesChart() {
        if (!this.charts.goroutines || this.metricsHistory.length === 0) {
            return;
        }

        const labels = this.metricsHistory.map(h => h.timestamp.toLocaleTimeString());
        const goroutinesData = this.metricsHistory.map(h => h.system?.goroutine_count || 0);

        this.charts.goroutines.data.labels = labels;
        this.charts.goroutines.data.datasets[0].data = goroutinesData;
        this.charts.goroutines.update('none');
    }

    // Fonction utilitaire pour générer des labels de temps
    generateTimeLabels(count) {
        const labels = [];
        const now = new Date();
        for (let i = count - 1; i >= 0; i--) {
            const time = new Date(now.getTime() - (i * 3000)); // 3 secondes d'intervalle
            labels.push(time.toLocaleTimeString());
        }
        return labels;
    }

    // Fonction pour redimensionner tous les graphiques
    resizeCharts() {
        Object.values(this.charts).forEach(chart => {
            if (chart && typeof chart.resize === 'function') {
                try {
                    chart.resize();
                } catch (error) {
                    console.warn('⚠️ Error resizing chart:', error);
                }
            }
        });
    }

    updatePerformanceRadarChart() {
        if (!this.charts.performance) return;

        // Calculate performance scores (0-100) based on current metrics
        const latest = this.metricsHistory[this.metricsHistory.length - 1];
        if (!latest) return;

        const perf = latest.performance || {};
        const scores = [
            this.calculatePerformanceScore('indexedStorage', perf.indexed_storage_stats),
            this.calculatePerformanceScore('hashJoin', perf.hash_join_stats),
            this.calculatePerformanceScore('evalCache', perf.evaluation_cache_stats),
            this.calculatePerformanceScore('tokenProp', perf.token_propagation_stats)
        ];

        this.charts.performance.data.datasets[0].data = scores;
        this.charts.performance.update('none');
    }

    updateCacheChart() {
        if (!this.charts.cache) return;

        const latest = this.metricsHistory[this.metricsHistory.length - 1];
        if (!latest) return;

        const evalStats = latest.performance?.evaluation_cache_stats || {};
        const hitRatio = evalStats.hit_ratio || 80;
        const missRatio = 100 - hitRatio;

        this.charts.cache.data.datasets[0].data = [hitRatio, missRatio];
        this.charts.cache.update('none');
    }

    calculatePerformanceScore(component, stats) {
        if (!stats) return 0;

        // Simple scoring algorithm - can be enhanced
        switch(component) {
            case 'indexedStorage':
                return Math.min(100, (stats.cache_hit_ratio || 0));
            case 'hashJoin':
                const hitRate = stats.cache_hits / (stats.cache_hits + stats.cache_misses + 1) * 100;
                return Math.min(100, hitRate);
            case 'evalCache':
                return Math.min(100, (stats.hit_ratio || 0));
            case 'tokenProp':
                return Math.min(100, (stats.parallel_efficiency || 0));
            default:
                return 0;
        }
    }

    // Tab Updates
    updateOverviewTab() {
        // Overview tab is updated automatically with metrics
    }

    updatePerformanceTab() {
        // Performance tab is updated automatically with metrics
    }

    updateNetworkTab() {
        this.refreshNetworkData();
    }

    updateAlertsTab() {
        this.refreshAlertsData();
    }

    updateSystemTab() {
        // System tab is updated automatically with metrics
    }

    // Data Refresh
    startDataRefresh() {
        // Initial load
        this.refreshAllData();

        // Periodic refresh for non-WebSocket data
        setInterval(() => {
            if (this.currentTab === 'network') {
                this.refreshNetworkData();
            } else if (this.currentTab === 'alerts') {
                this.refreshAlertsData();
            }
        }, 30000); // Refresh every 30 seconds
    }

    async refreshAllData() {
        try {
            const [metrics, network, alerts] = await Promise.all([
                this.fetchAPI('/api/metrics'),
                this.fetchAPI('/api/network/status'),
                this.fetchAPI('/api/alerts')
            ]);

            if (metrics) {
                this.updateMetrics({
                    system: metrics.system_metrics,
                    rete: metrics.rete_metrics,
                    performance: metrics.performance_metrics
                });
            }

        } catch (error) {
            console.error('❌ Error refreshing data:', error);
            this.showToast('Error', 'Failed to refresh data', 'error');
        }
    }

    async refreshNetworkData() {
        try {
            const [status, nodes] = await Promise.all([
                this.fetchAPI('/api/network/status'),
                this.fetchAPI('/api/network/nodes')
            ]);

            this.updateNetworkDisplay(status, nodes);

        } catch (error) {
            console.error('❌ Error refreshing network data:', error);
        }
    }

    async refreshAlertsData() {
        try {
            const alerts = await this.fetchAPI('/api/alerts');
            this.updateAlertsDisplay(alerts);

        } catch (error) {
            console.error('❌ Error refreshing alerts data:', error);
        }
    }

    updateNetworkDisplay(status, nodes) {
        // Update network graph placeholder with actual data
        const networkGraph = document.getElementById('networkGraph');
        if (status && nodes) {
            networkGraph.innerHTML = `
                <div class="network-info">
                    <h4>Network Status: ${status.status}</h4>
                    <p>Nodes: ${nodes.length}</p>
                    <p>Uptime: ${status.uptime}</p>
                </div>
            `;
        }
    }

    updateAlertsDisplay(alertsData) {
        if (!alertsData) return;

        // Update active alerts
        const activeAlertsContainer = document.getElementById('activeAlerts');
        if (alertsData.active_alerts && alertsData.active_alerts.length > 0) {
            activeAlertsContainer.innerHTML = alertsData.active_alerts
                .map(alert => this.createAlertHTML(alert))
                .join('');
        } else {
            activeAlertsContainer.innerHTML = '<div class="alert-placeholder">No active alerts</div>';
        }

        // Update alert rules
        const rulesContainer = document.getElementById('alertRules');
        if (alertsData.rules_count > 0) {
            rulesContainer.innerHTML = `<div class="rule-placeholder">${alertsData.rules_count} rules configured</div>`;
        } else {
            rulesContainer.innerHTML = '<div class="rule-placeholder">No alert rules configured</div>';
        }
    }

    createAlertHTML(alert) {
        return `
            <div class="alert-item">
                <div class="alert-content">
                    <h4>${alert.rule_name}</h4>
                    <p>${alert.message}</p>
                    <small>${new Date(alert.timestamp).toLocaleString()}</small>
                </div>
                <span class="alert-severity ${alert.severity}">${alert.severity}</span>
            </div>
        `;
    }

    // Alert Management
    showCreateAlertModal() {
        const modal = document.getElementById('alertRuleModal');
        modal.classList.add('show');
    }

    async createAlertRule() {
        const form = document.getElementById('alertRuleForm');
        const formData = new FormData(form);

        const rule = {
            name: formData.get('name'),
            description: formData.get('description'),
            condition: formData.get('condition'),
            threshold: parseFloat(formData.get('threshold')),
            severity: formData.get('severity'),
            is_enabled: true
        };

        try {
            await this.fetchAPI('/api/alerts/rules', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(rule)
            });

            this.showToast('Success', 'Alert rule created successfully', 'success');
            document.getElementById('alertRuleModal').classList.remove('show');
            form.reset();
            this.refreshAlertsData();

        } catch (error) {
            console.error('❌ Error creating alert rule:', error);
            this.showToast('Error', 'Failed to create alert rule', 'error');
        }
    }

    exportNetworkData() {
        // Export network topology as JSON
        this.fetchAPI('/api/network/nodes')
            .then(nodes => {
                const dataStr = JSON.stringify(nodes, null, 2);
                const dataBlob = new Blob([dataStr], { type: 'application/json' });

                const link = document.createElement('a');
                link.href = URL.createObjectURL(dataBlob);
                link.download = `rete-network-${new Date().toISOString().slice(0, 10)}.json`;
                link.click();
            })
            .catch(error => {
                console.error('❌ Error exporting network data:', error);
                this.showToast('Error', 'Failed to export network data', 'error');
            });
    }

    // Utility Functions
    async fetchAPI(endpoint, options = {}) {
        try {
            const response = await fetch(endpoint, options);
            if (!response.ok) {
                throw new Error(`HTTP ${response.status}: ${response.statusText}`);
            }
            return await response.json();
        } catch (error) {
            console.error(`❌ API Error (${endpoint}):`, error);
            throw error;
        }
    }

    updateElement(id, value, attribute = 'textContent') {
        const element = document.getElementById(id);
        if (element) {
            if (attribute === 'style') {
                element.style.cssText = value;
            } else {
                element[attribute] = value;
            }
        }
    }

    formatNumber(num) {
        if (num >= 1000000) {
            return (num / 1000000).toFixed(1) + 'M';
        } else if (num >= 1000) {
            return (num / 1000).toFixed(1) + 'K';
        }
        return num.toString();
    }

    formatBytes(bytes) {
        const sizes = ['B', 'KB', 'MB', 'GB'];
        if (bytes === 0) return '0 B';
        const i = Math.floor(Math.log(bytes) / Math.log(1024));
        return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i];
    }

    formatDuration(seconds) {
        const h = Math.floor(seconds / 3600);
        const m = Math.floor((seconds % 3600) / 60);
        const s = Math.floor(seconds % 60);

        if (h > 0) return `${h}h ${m}m`;
        if (m > 0) return `${m}m ${s}s`;
        return `${s}s`;
    }

    showToast(title, message, type = 'info') {
        const toastContainer = document.getElementById('toastContainer');
        const toast = document.createElement('div');
        toast.className = `toast ${type}`;

        toast.innerHTML = `
            <div class="toast-header">
                <span class="toast-title">${title}</span>
                <span class="toast-close">&times;</span>
            </div>
            <div class="toast-message">${message}</div>
        `;

        toastContainer.appendChild(toast);

        // Show animation
        setTimeout(() => toast.classList.add('show'), 100);

        // Close button
        toast.querySelector('.toast-close').addEventListener('click', () => {
            toast.classList.remove('show');
            setTimeout(() => toast.remove(), 300);
        });

        // Auto-remove after 5 seconds
        setTimeout(() => {
            if (toast.parentNode) {
                toast.classList.remove('show');
                setTimeout(() => toast.remove(), 300);
            }
        }, 5000);
    }
}

// Initialize dashboard when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.reteDashboard = new RETEDashboard();
});

// Handle window resize for chart responsiveness
window.addEventListener('resize', () => {
    if (window.reteDashboard && window.reteDashboard.resizeCharts) {
        // Debounce le redimensionnement pour éviter les appels excessifs
        clearTimeout(window.reteDashboard.resizeTimeout);
        window.reteDashboard.resizeTimeout = setTimeout(() => {
            window.reteDashboard.resizeCharts();
        }, 150);
    }
});
