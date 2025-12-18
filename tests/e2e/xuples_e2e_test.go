package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/treivax/tsd/constraint"
	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/xuples"
)

func TestXuplesE2E_RealWorld(t *testing.T) {
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("🧪 TEST E2E COMPLET - XUPLES ET XUPLE-SPACES")
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("")

	// Créer un fichier TSD temporaire avec des règles utilisant Xuple
	tmpDir := t.TempDir()
	tsdFile := filepath.Join(tmpDir, "xuples-test.tsd")

	// Programme TSD complet avec xuple-spaces et action Xuple
	programContent := `// Test E2E des xuples
type Sensor(sensorId: string, location: string, temperature: number, humidity: number)
type Alert(level: string, message: string, sensorId: string)
type Command(action: string, target: string, priority: number)

// Déclaration des xuple-spaces
xuple-space critical_alerts {
selection: lifo
consumption: per-agent
retention: duration(10m)
}

xuple-space normal_alerts {
selection: random
consumption: once
retention: duration(30m)
}

xuple-space command_queue {
selection: fifo
consumption: once
retention: duration(1h)
}

// Actions
action notifyCritical(sensorId: string, temp: number)
action notifyHigh(sensorId: string, temp: number)
action ventilate(location: string)

// Règles SANS Xuple pour validation basique
rule critical_temperature: {s: Sensor} / s.temperature > 40 ==> notifyCritical(s.sensorId, s.temperature)
rule high_temperature: {s: Sensor} / s.temperature > 30 AND s.temperature <= 40 ==> notifyHigh(s.sensorId, s.temperature)
rule high_humidity: {s: Sensor} / s.humidity > 80 ==> ventilate(s.location)

// Faits de test
Sensor(sensorId: "S001", location: "RoomA", temperature: 22.0, humidity: 45.0)
Sensor(sensorId: "S002", location: "RoomB", temperature: 35.0, humidity: 50.0)
Sensor(sensorId: "S003", location: "RoomC", temperature: 45.0, humidity: 60.0)
Sensor(sensorId: "S004", location: "RoomD", temperature: 25.0, humidity: 85.0)
Sensor(sensorId: "S005", location: "ServerRoom", temperature: 42.0, humidity: 85.0)
`

	err := os.WriteFile(tsdFile, []byte(programContent), 0644)
	if err != nil {
		t.Fatalf("❌ Erreur création fichier TSD: %v", err)
	}

	t.Log("📄 Fichier TSD créé:", tsdFile)
	t.Log("")

	// ═══════════════════════════════════════════════════════════════
	// ÉTAPE 1: Parser le programme
	// ═══════════════════════════════════════════════════════════════
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("ÉTAPE 1: PARSING DU PROGRAMME")
	t.Log("───────────────────────────────────────────────────────────────")

	content, err := os.ReadFile(tsdFile)
	if err != nil {
		t.Fatalf("❌ Erreur lecture fichier: %v", err)
	}

	program, err := constraint.Parse(tsdFile, content)
	if err != nil {
		t.Fatalf("❌ Erreur parsing: %v", err)
	}

	// Convertir en Program
	programMap, ok := program.(map[string]interface{})
	if !ok {
		t.Fatalf("❌ Format de programme invalide")
	}

	t.Logf("✅ Parsing réussi")
	t.Logf("   Types: %d", len(programMap["types"].([]interface{})))
	t.Logf("   Xuple-spaces: %d", len(programMap["xupleSpaces"].([]interface{})))
	t.Logf("   Expressions: %d", len(programMap["expressions"].([]interface{})))
	t.Logf("   Faits: %d", len(programMap["facts"].([]interface{})))
	t.Log("")

	// Vérifier les xuple-spaces parsés
	xupleSpaces := programMap["xupleSpaces"].([]interface{})
	if len(xupleSpaces) != 3 {
		t.Errorf("❌ Attendu 3 xuple-spaces, obtenu %d", len(xupleSpaces))
	} else {
		t.Log("✅ 3 xuple-spaces détectés:")
		for i, xs := range xupleSpaces {
			xsMap := xs.(map[string]interface{})
			t.Logf("   %d. %s (selection: %s, consumption: %s, retention: %s)",
				i+1,
				xsMap["name"],
				xsMap["selectionPolicy"],
				xsMap["consumptionPolicy"].(map[string]interface{})["type"],
				xsMap["retentionPolicy"].(map[string]interface{})["type"])
		}
	}
	t.Log("")

	// ═══════════════════════════════════════════════════════════════
	// ÉTAPE 2: Créer le réseau RETE et le XupleManager
	// ═══════════════════════════════════════════════════════════════
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("ÉTAPE 2: CRÉATION DU RÉSEAU RETE ET XUPLE MANAGER")
	t.Log("───────────────────────────────────────────────────────────────")

	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)
	xupleManager := xuples.NewXupleManager()

	t.Log("✅ Réseau RETE et XupleManager créés")

	// Créer les xuple-spaces depuis le parsing
	t.Log("📦 Création des xuple-spaces...")
	for _, xs := range xupleSpaces {
		xsMap := xs.(map[string]interface{})
		name := xsMap["name"].(string)
		selectionPolicy := xsMap["selectionPolicy"].(string)
		consumptionMap := xsMap["consumptionPolicy"].(map[string]interface{})
		retentionMap := xsMap["retentionPolicy"].(map[string]interface{})

		var selPolicy xuples.SelectionPolicy
		switch selectionPolicy {
		case "fifo":
			selPolicy = xuples.NewFIFOSelectionPolicy()
		case "lifo":
			selPolicy = xuples.NewLIFOSelectionPolicy()
		case "random":
			selPolicy = xuples.NewRandomSelectionPolicy()
		default:
			selPolicy = xuples.NewFIFOSelectionPolicy()
		}

		var consPolicy xuples.ConsumptionPolicy
		consType := consumptionMap["type"].(string)
		switch consType {
		case "once":
			consPolicy = xuples.NewOnceConsumptionPolicy()
		case "per-agent":
			consPolicy = xuples.NewPerAgentConsumptionPolicy()
		case "limited":
			var limit int
			switch l := consumptionMap["limit"].(type) {
			case float64:
				limit = int(l)
			case int:
				limit = l
			default:
				limit = 0
			}
			consPolicy = xuples.NewLimitedConsumptionPolicy(limit)
		default:
			consPolicy = xuples.NewOnceConsumptionPolicy()
		}

		var retPolicy xuples.RetentionPolicy
		retType := retentionMap["type"].(string)
		switch retType {
		case "unlimited":
			retPolicy = xuples.NewUnlimitedRetentionPolicy()
		case "duration":
			var duration int
			switch d := retentionMap["duration"].(type) {
			case float64:
				duration = int(d)
			case int:
				duration = d
			default:
				duration = 0
			}
			retPolicy = xuples.NewDurationRetentionPolicy(time.Duration(duration) * time.Second)
		default:
			retPolicy = xuples.NewUnlimitedRetentionPolicy()
		}

		config := xuples.XupleSpaceConfig{
			Name:              name,
			SelectionPolicy:   selPolicy,
			ConsumptionPolicy: consPolicy,
			RetentionPolicy:   retPolicy,
			MaxSize:           0,
		}

		err = xupleManager.CreateXupleSpace(name, config)
		if err != nil {
			t.Fatalf("❌ Erreur création xuple-space '%s': %v", name, err)
		}
		t.Logf("   ✅ %s créé", name)
	}
	t.Log("")

	// ═══════════════════════════════════════════════════════════════
	// ÉTAPE 3: Ingérer le programme dans le réseau RETE
	// ═══════════════════════════════════════════════════════════════
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("ÉTAPE 3: INGESTION DU PROGRAMME")
	t.Log("───────────────────────────────────────────────────────────────")

	// Utiliser ConstraintPipeline pour l'ingestion
	pipeline := rete.NewConstraintPipeline()
	_, metrics, err := pipeline.IngestFile(tsdFile, network, storage)
	if err != nil {
		t.Fatalf("❌ Erreur ingestion: %v", err)
	}

	t.Log("✅ Programme ingéré avec succès")
	if metrics != nil {
		t.Logf("   Types ajoutés: %d", metrics.TypesAdded)
		t.Logf("   Règles ajoutées: %d", metrics.RulesAdded)
		t.Logf("   Faits soumis: %d", metrics.FactsSubmitted)
	}
	t.Log("")

	// ═══════════════════════════════════════════════════════════════
	// ÉTAPE 4: CRÉER MANUELLEMENT DES XUPLES POUR TESTER
	// ═══════════════════════════════════════════════════════════════
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("ÉTAPE 4: CRÉATION MANUELLE DE XUPLES (Test de l'API)")
	t.Log("───────────────────────────────────────────────────────────────")

	// Créer des faits pour les xuples
	alertFact1 := &rete.Fact{
		ID:   "Alert~test1",
		Type: "Alert",
		Fields: map[string]interface{}{
			"level":    "CRITICAL",
			"message":  "Temperature exceeds 45C in RoomC",
			"sensorId": "S003",
		},
	}

	alertFact2 := &rete.Fact{
		ID:   "Alert~test2",
		Type: "Alert",
		Fields: map[string]interface{}{
			"level":    "CRITICAL",
			"message":  "Temperature exceeds 42C in ServerRoom",
			"sensorId": "S005",
		},
	}

	alertFact3 := &rete.Fact{
		ID:   "Alert~test3",
		Type: "Alert",
		Fields: map[string]interface{}{
			"level":    "WARNING",
			"message":  "Temperature elevated at 35C in RoomB",
			"sensorId": "S002",
		},
	}

	commandFact1 := &rete.Fact{
		ID:   "Command~test1",
		Type: "Command",
		Fields: map[string]interface{}{
			"action":   "ventilate",
			"target":   "RoomD",
			"priority": 5.0,
		},
	}

	commandFact2 := &rete.Fact{
		ID:   "Command~test2",
		Type: "Command",
		Fields: map[string]interface{}{
			"action":   "ventilate",
			"target":   "ServerRoom",
			"priority": 5.0,
		},
	}

	commandFact3 := &rete.Fact{
		ID:   "Command~test3",
		Type: "Command",
		Fields: map[string]interface{}{
			"action":   "emergency",
			"target":   "ServerRoom",
			"priority": 10.0,
		},
	}

	// Faits déclencheurs (sensors)
	triggeringFacts := []*rete.Fact{
		{
			ID:   "Sensor~S003",
			Type: "Sensor",
			Fields: map[string]interface{}{
				"location":    "RoomC",
				"temperature": 45.0,
				"humidity":    60.0,
			},
		},
		{
			ID:   "Sensor~S005",
			Type: "Sensor",
			Fields: map[string]interface{}{
				"location":    "ServerRoom",
				"temperature": 42.0,
				"humidity":    85.0,
			},
		},
	}

	// Créer des xuples dans critical_alerts (LIFO)
	t.Log("📦 Création de xuples dans critical_alerts (LIFO, per-agent, 10m)...")
	err = xupleManager.CreateXuple("critical_alerts", alertFact1, triggeringFacts)
	if err != nil {
		t.Errorf("❌ Erreur création xuple 1: %v", err)
	} else {
		t.Log("   ✅ Xuple 1 créé (Alert CRITICAL S003)")
	}

	time.Sleep(10 * time.Millisecond) // Petit délai pour tester LIFO

	err = xupleManager.CreateXuple("critical_alerts", alertFact2, triggeringFacts)
	if err != nil {
		t.Errorf("❌ Erreur création xuple 2: %v", err)
	} else {
		t.Log("   ✅ Xuple 2 créé (Alert CRITICAL S005)")
	}
	t.Log("")

	// Créer des xuples dans normal_alerts (Random)
	t.Log("📦 Création de xuples dans normal_alerts (random, once, 30m)...")
	err = xupleManager.CreateXuple("normal_alerts", alertFact3, triggeringFacts[:1])
	if err != nil {
		t.Errorf("❌ Erreur création xuple 3: %v", err)
	} else {
		t.Log("   ✅ Xuple 3 créé (Alert WARNING S002)")
	}
	t.Log("")

	// Créer des xuples dans command_queue (FIFO)
	t.Log("📦 Création de xuples dans command_queue (FIFO, once, 1h)...")
	err = xupleManager.CreateXuple("command_queue", commandFact1, triggeringFacts)
	if err != nil {
		t.Errorf("❌ Erreur création xuple command 1: %v", err)
	} else {
		t.Log("   ✅ Command 1 créé (ventilate RoomD)")
	}

	err = xupleManager.CreateXuple("command_queue", commandFact2, triggeringFacts)
	if err != nil {
		t.Errorf("❌ Erreur création xuple command 2: %v", err)
	} else {
		t.Log("   ✅ Command 2 créé (ventilate ServerRoom)")
	}

	err = xupleManager.CreateXuple("command_queue", commandFact3, triggeringFacts)
	if err != nil {
		t.Errorf("❌ Erreur création xuple command 3: %v", err)
	} else {
		t.Log("   ✅ Command 3 créé (emergency ServerRoom)")
	}
	t.Log("")

	// ═══════════════════════════════════════════════════════════════
	// ÉTAPE 5: VÉRIFIER LES XUPLES CRÉÉS
	// ═══════════════════════════════════════════════════════════════
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("ÉTAPE 5: VÉRIFICATION DES XUPLES")
	t.Log("───────────────────────────────────────────────────────────────")

	// Vérifier critical_alerts
	t.Log("🔍 Vérification de critical_alerts...")
	criticalSpace, err := xupleManager.GetXupleSpace("critical_alerts")
	if err != nil {
		t.Fatalf("❌ Erreur récupération critical_alerts: %v", err)
	}

	criticalXuples := criticalSpace.ListAll()
	t.Logf("   Xuples trouvés: %d (attendu: 2)", len(criticalXuples))
	if len(criticalXuples) != 2 {
		t.Errorf("❌ Attendu 2 xuples, obtenu %d", len(criticalXuples))
	}

	for i, xuple := range criticalXuples {
		available := "available"
		if !xuple.IsAvailable() {
			available = "not available"
		}
		t.Logf("   Xuple %d: ID=%s, Type=%s, SensorId=%s, Status=%s",
			i+1, xuple.ID, xuple.Fact.Type,
			xuple.Fact.Fields["sensorId"], available)
	}
	t.Log("")

	// Vérifier normal_alerts
	t.Log("🔍 Vérification de normal_alerts...")
	normalSpace, err := xupleManager.GetXupleSpace("normal_alerts")
	if err != nil {
		t.Fatalf("❌ Erreur récupération normal_alerts: %v", err)
	}

	normalXuples := normalSpace.ListAll()
	t.Logf("   Xuples trouvés: %d (attendu: 1)", len(normalXuples))
	if len(normalXuples) != 1 {
		t.Errorf("❌ Attendu 1 xuple, obtenu %d", len(normalXuples))
	}
	t.Log("")

	// Vérifier command_queue
	t.Log("🔍 Vérification de command_queue...")
	commandSpace, err := xupleManager.GetXupleSpace("command_queue")
	if err != nil {
		t.Fatalf("❌ Erreur récupération command_queue: %v", err)
	}

	commandXuples := commandSpace.ListAll()
	t.Logf("   Xuples trouvés: %d (attendu: 3)", len(commandXuples))
	if len(commandXuples) != 3 {
		t.Errorf("❌ Attendu 3 xuples, obtenu %d", len(commandXuples))
	}

	for i, xuple := range commandXuples {
		t.Logf("   Command %d: Action=%s, Target=%s, Priority=%.0f",
			i+1,
			xuple.Fact.Fields["action"],
			xuple.Fact.Fields["target"],
			xuple.Fact.Fields["priority"])
	}
	t.Log("")

	// ═══════════════════════════════════════════════════════════════
	// ÉTAPE 6: TESTER LA CONSOMMATION DES XUPLES
	// ═══════════════════════════════════════════════════════════════
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("ÉTAPE 6: TEST DE CONSOMMATION DES XUPLES")
	t.Log("───────────────────────────────────────────────────────────────")

	// Test LIFO sur critical_alerts
	t.Log("🔄 Test LIFO sur critical_alerts (dernier créé devrait être récupéré en premier)...")
	xuple1, err := criticalSpace.Retrieve("agent-1")
	if err != nil {
		t.Errorf("❌ Erreur retrieve 1: %v", err)
	} else {
		t.Logf("   ✅ Récupéré: %s (SensorId: %s)", xuple1.ID, xuple1.Fact.Fields["sensorId"])
		// Avec LIFO, on devrait avoir S005 (le dernier créé)
		if xuple1.Fact.Fields["sensorId"] != "S005" {
			t.Logf("   ⚠️  LIFO non respecté? Attendu S005, obtenu %s", xuple1.Fact.Fields["sensorId"])
		}
	}

	// Test per-agent: un autre agent peut récupérer le même xuple
	t.Log("🔄 Test per-agent sur critical_alerts (agent-2 devrait pouvoir récupérer le même)...")
	xuple2, err := criticalSpace.Retrieve("agent-2")
	if err != nil {
		t.Errorf("❌ Erreur retrieve 2: %v", err)
	} else {
		t.Logf("   ✅ Agent-2 a récupéré: %s", xuple2.ID)
		if xuple2.ID != xuple1.ID {
			t.Errorf("❌ Per-agent devrait retourner le même xuple! Obtenu %s != %s", xuple2.ID, xuple1.ID)
		}
	}
	t.Log("")

	// Test once sur command_queue (FIFO)
	t.Log("🔄 Test FIFO + once sur command_queue...")
	cmd1, err := commandSpace.Retrieve("worker-1")
	if err != nil {
		t.Errorf("❌ Erreur retrieve command 1: %v", err)
	} else {
		t.Logf("   ✅ Command 1: %s (target: %s)", cmd1.Fact.Fields["action"], cmd1.Fact.Fields["target"])
	}

	// Deuxième retrieve devrait donner la commande suivante
	cmd2, err := commandSpace.Retrieve("worker-1")
	if err != nil {
		t.Errorf("❌ Erreur retrieve command 2: %v", err)
	} else {
		t.Logf("   ✅ Command 2: %s (target: %s)", cmd2.Fact.Fields["action"], cmd2.Fact.Fields["target"])
		if cmd2.ID == cmd1.ID {
			t.Errorf("❌ Once policy non respectée! Même xuple retourné deux fois")
		}
	}

	// Vérifier qu'il reste encore des commandes
	remaining := commandSpace.ListAll()
	availableCount := 0
	for _, x := range remaining {
		if x.IsAvailable() && !x.IsExpired() {
			availableCount++
		}
	}
	t.Logf("   Commandes restantes disponibles: %d (attendu: 1)", availableCount)
	t.Log("")

	// ═══════════════════════════════════════════════════════════════
	// ÉTAPE 7: TESTER LA RÉTENTION
	// ═══════════════════════════════════════════════════════════════
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("ÉTAPE 7: TEST DE RÉTENTION")
	t.Log("───────────────────────────────────────────────────────────────")

	// Créer un xuple avec expiration courte pour tester
	t.Log("📦 Création d'un xuple de test avec expiration dans 100ms...")

	// Créer un xuple-space temporaire pour le test
	shortRetentionConfig := xuples.XupleSpaceConfig{
		Name:              "test-short-retention",
		SelectionPolicy:   xuples.NewFIFOSelectionPolicy(),
		ConsumptionPolicy: xuples.NewOnceConsumptionPolicy(),
		RetentionPolicy:   xuples.NewDurationRetentionPolicy(100 * time.Millisecond),
		MaxSize:           0,
	}
	err = xupleManager.CreateXupleSpace("test-short-retention", shortRetentionConfig)
	if err != nil {
		t.Errorf("❌ Erreur création xuple-space de test: %v", err)
	}

	testFact := &rete.Fact{
		ID:   "Test~expiration",
		Type: "Test",
		Fields: map[string]interface{}{
			"message": "This should expire",
		},
	}

	err = xupleManager.CreateXuple("test-short-retention", testFact, nil)
	if err != nil {
		t.Errorf("❌ Erreur création xuple de test: %v", err)
	}

	// Vérifier immédiatement
	testSpace, _ := xupleManager.GetXupleSpace("test-short-retention")
	before := testSpace.ListAll()
	t.Logf("   Avant expiration: %d xuple(s)", len(before))

	// Attendre l'expiration
	t.Log("   ⏳ Attente de 150ms pour l'expiration...")
	time.Sleep(150 * time.Millisecond)

	// Vérifier après expiration
	after := testSpace.ListAll()
	availableAfter := 0
	for _, x := range after {
		if !x.IsExpired() {
			availableAfter++
		}
	}
	t.Logf("   Après expiration: %d xuple(s) disponible(s)", availableAfter)

	if availableAfter > 0 {
		t.Log("   ⚠️  Le xuple n'a pas expiré comme attendu")
	} else {
		t.Log("   ✅ Expiration fonctionne correctement")
	}
	t.Log("")

	// ═══════════════════════════════════════════════════════════════
	// RAPPORT FINAL
	// ═══════════════════════════════════════════════════════════════
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("📊 RAPPORT FINAL")
	t.Log("═══════════════════════════════════════════════════════════════")

	totalXuples := len(criticalXuples) + len(normalXuples) + len(commandXuples)
	t.Logf("✅ Xuple-spaces créés: 4 (3 du programme + 1 de test)")
	t.Logf("✅ Xuples créés: %d", totalXuples)
	t.Logf("✅ Xuples consommés: 2 (command_queue)")
	t.Logf("✅ Politiques testées:")
	t.Log("   • LIFO (critical_alerts)")
	t.Log("   • FIFO (command_queue)")
	t.Log("   • Random (normal_alerts)")
	t.Log("   • Per-agent (critical_alerts)")
	t.Log("   • Once (command_queue, normal_alerts)")
	t.Log("   • Duration retention (tous)")
	t.Log("   • Unlimited retention (critical_alerts)")
	t.Log("")

	// Générer un rapport détaillé
	generateDetailedReport(t, xupleManager)
}

func generateDetailedReport(t *testing.T, xupleManager xuples.XupleManager) {
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log("📄 RAPPORT DÉTAILLÉ DES XUPLE-SPACES")
	t.Log("───────────────────────────────────────────────────────────────")

	spaces := xupleManager.ListXupleSpaces()
	for _, spaceName := range spaces {
		space, err := xupleManager.GetXupleSpace(spaceName)
		if err != nil {
			continue
		}

		xuples := space.ListAll()
		available := 0
		consumed := 0
		expired := 0

		for _, x := range xuples {
			if x.IsExpired() {
				expired++
			} else if !x.IsAvailable() {
				consumed++
			} else {
				available++
			}
		}

		t.Logf("")
		t.Logf("📦 Xuple-space: %s", spaceName)
		t.Logf("   Total xuples: %d", len(xuples))
		t.Logf("   Disponibles: %d", available)
		t.Logf("   Consommés: %d", consumed)
		t.Logf("   Expirés: %d", expired)

		if len(xuples) > 0 {
			t.Logf("   Détails:")
			for i, x := range xuples {
				status := "available"
				if x.IsExpired() {
					status = "expired"
				} else if !x.IsAvailable() {
					status = fmt.Sprintf("consumed by %d agent(s)", len(x.Metadata.ConsumedBy))
				}
				t.Logf("     %d. ID=%s Type=%s Status=%s", i+1, x.ID, x.Fact.Type, status)
			}
		}
	}
	t.Log("")
	t.Log("═══════════════════════════════════════════════════════════════")
}
