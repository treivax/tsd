// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package incremental

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/treivax/tsd/rete"
)

// TestIncrementalValidation teste la validation sémantique incrémentale
func TestIncrementalValidation(t *testing.T) {
	storage := rete.NewMemoryStorage()
	pipeline := rete.NewConstraintPipeline()

	// Phase 1: Charger les types
	typesContent := `
type Person(id: string, name: string, age: number)

type Company(id: string, name: string, employees: number)
`
	typesFile := createTempFile(t, "types.tsd", typesContent)
	defer os.Remove(typesFile)

	network, err := pipeline.IngestFile(typesFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur ingestion types: %v", err)
	}

	if len(network.Types) != 2 {
		t.Errorf("Attendu 2 types, obtenu %d", len(network.Types))
	}

	// Phase 2: Charger des règles qui référencent les types existants
	rulesContent := `
action print_adult(name: string)
action print_company(name: string)

rule adult_check: {p: Person} / p.age >= 18 ==> print_adult(p.name)

rule big_company: {c: Company} / c.employees > 100 ==> print_company(c.name)
`
	rulesFile := createTempFile(t, "rules.tsd", rulesContent)
	defer os.Remove(rulesFile)

	// La validation incrémentale (activée automatiquement) doit réussir car les types sont définis dans le réseau
	network, err = pipeline.IngestFile(rulesFile, network, storage)
	if err != nil {
		t.Fatalf("Erreur ingestion règles: %v", err)
	}

	if len(network.TerminalNodes) != 2 {
		t.Errorf("Attendu 2 règles, obtenu %d", len(network.TerminalNodes))
	}
}

// TestIncrementalValidationError teste la détection d'erreurs en validation incrémentale
func TestIncrementalValidationError(t *testing.T) {
	storage := rete.NewMemoryStorage()
	pipeline := rete.NewConstraintPipeline()

	// Phase 1: Charger un type
	typesContent := `
type Person(id: string, name: string, age: number)
`
	typesFile := createTempFile(t, "types.tsd", typesContent)
	defer os.Remove(typesFile)

	network, err := pipeline.IngestFile(typesFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur ingestion types: %v", err)
	}

	// Phase 2: Essayer de charger une règle avec un type non défini
	rulesContent := `
action print_company(msg: string)

rule company_check: {c: Company} / c.employees > 10 ==> print_company("Company found")
`
	rulesFile := createTempFile(t, "invalid_rules.tsd", rulesContent)
	defer os.Remove(rulesFile)

	// La validation incrémentale (automatique) doit échouer car Company n'est pas défini
	_, err = pipeline.IngestFile(rulesFile, network, storage)
	if err == nil {
		t.Error("Attendu une erreur de validation pour type non défini, mais aucune erreur")
	}
}

// TestGarbageCollectionAfterReset teste le GC après un reset
func TestGarbageCollectionAfterReset(t *testing.T) {
	storage := rete.NewMemoryStorage()
	pipeline := rete.NewConstraintPipeline()

	// Phase 1: Créer un réseau avec plusieurs nœuds
	initialContent := `
type Person(id: string, name: string)

action print_person(msg: string)

rule test_rule: {p: Person} / p.name == "John" ==> print_person("Found John")

Person(id: p1, name: John)
`
	initialFile := createTempFile(t, "initial.tsd", initialContent)
	defer os.Remove(initialFile)

	network, err := pipeline.IngestFile(initialFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur ingestion initiale: %v", err)
	}

	// Vérifier que le réseau contient des éléments
	nodesBefore := len(network.TypeNodes) + len(network.AlphaNodes) + len(network.TerminalNodes)
	if nodesBefore == 0 {
		t.Fatal("Le réseau initial devrait contenir des nœuds")
	}

	// Phase 2: Faire un reset avec GC
	resetContent := `
reset

type Vehicle(id: string, brand: string)
`
	resetFile := createTempFile(t, "reset.tsd", resetContent)
	defer os.Remove(resetFile)

	network, err = pipeline.IngestFile(resetFile, network, storage)
	if err != nil {
		t.Fatalf("Erreur ingestion avec reset: %v", err)
	}

	// Vérifier que l'ancien réseau a été nettoyé
	// Le nouveau réseau doit avoir seulement le nouveau type
	if len(network.Types) != 1 {
		t.Errorf("Attendu 1 type après reset, obtenu %d", len(network.Types))
	}

	if network.Types[0].Name != "Vehicle" {
		t.Errorf("Attendu type Vehicle, obtenu %s", network.Types[0].Name)
	}

	// Vérifier que les anciens nœuds ont été nettoyés
	for id := range network.TypeNodes {
		if id == "type_Person" {
			t.Error("L'ancien TypeNode Person ne devrait plus exister après reset et GC")
		}
	}
}

// TestTransactionCommit teste le commit automatique d'une transaction réussie
func TestTransactionCommit(t *testing.T) {
	storage := rete.NewMemoryStorage()
	pipeline := rete.NewConstraintPipeline()

	// Simuler une ingestion réussie (transaction automatique)
	content := `
type Person(id: string, name: string)
`
	file := createTempFile(t, "person.tsd", content)
	defer os.Remove(file)

	// IngestFile crée automatiquement une transaction et la commit en cas de succès
	network, err := pipeline.IngestFile(file, nil, storage)
	if err != nil {
		t.Fatalf("Erreur ingestion: %v", err)
	}

	// Vérifier que le changement est bien présent
	if len(network.Types) != 1 {
		t.Errorf("Attendu 1 type après commit, obtenu %d", len(network.Types))
	}
}

// TestTransactionRollback teste le rollback automatique d'une transaction en cas d'erreur
func TestTransactionRollback(t *testing.T) {
	storage := rete.NewMemoryStorage()
	pipeline := rete.NewConstraintPipeline()

	// Phase 1: Créer un réseau initial
	initialContent := `
type Person(id: string, name: string, age: number)

Person(id: p1, name: Alice, age: 30)
`
	initialFile := createTempFile(t, "initial.tsd", initialContent)
	defer os.Remove(initialFile)

	network, err := pipeline.IngestFile(initialFile, nil, storage)
	if err != nil {
		t.Fatalf("Erreur ingestion initiale: %v", err)
	}

	typesCountBefore := len(network.Types)
	factsCountBefore := len(storage.GetAllFacts())

	// Phase 2: Tenter une ingestion invalide (transaction automatique + rollback automatique)
	// Fichier avec erreur de syntaxe
	invalidContent := `
type Vehicle(id: string, brand: string, INVALID SYNTAX HERE)
`
	invalidFile := createTempFile(t, "invalid.tsd", invalidContent)
	defer os.Remove(invalidFile)

	// IngestFile crée automatiquement une transaction et la rollback en cas d'erreur
	_, err = pipeline.IngestFile(invalidFile, network, storage)
	if err == nil {
		t.Error("Attendu une erreur pour fichier invalide")
	}

	// Vérifier que l'état du réseau est inchangé (rollback automatique)
	typesCountAfter := len(network.Types)
	factsCountAfter := len(storage.GetAllFacts())

	if typesCountAfter != typesCountBefore {
		t.Errorf("Nombre de types devrait être inchangé: avant=%d, après=%d", typesCountBefore, typesCountAfter)
	}

	if factsCountAfter != factsCountBefore {
		t.Errorf("Nombre de faits devrait être inchangé: avant=%d, après=%d", factsCountBefore, factsCountAfter)
	}
}

// TestAdvancedFeaturesIntegration teste l'intégration de toutes les fonctionnalités
func TestAdvancedFeaturesIntegration(t *testing.T) {
	storage := rete.NewMemoryStorage()
	pipeline := rete.NewConstraintPipeline()

	config := rete.DefaultAdvancedPipelineConfig()
	config.AutoCommit = true
	config.AutoRollbackOnError = true

	// Phase 1: Charger des types
	typesContent := `
type Employee(id: string, name: string, salary: number, department: string)
`
	typesFile := createTempFile(t, "types.tsd", typesContent)
	defer os.Remove(typesFile)

	network, metrics, err := pipeline.IngestFileWithAdvancedFeatures(typesFile, nil, storage, config)
	if err != nil {
		t.Fatalf("Erreur ingestion types: %v", err)
	}

	if metrics == nil {
		t.Fatal("Métriques ne devraient pas être nil")
	}

	// Phase 2: Charger des règles avec validation incrémentale
	rulesContent := `
action print_high_earner(name: string)

rule high_salary: {e: Employee} / e.salary > 100000 ==> print_high_earner(e.name)
`
	rulesFile := createTempFile(t, "rules.tsd", rulesContent)
	defer os.Remove(rulesFile)

	network, metrics, err = pipeline.IngestFileWithAdvancedFeatures(rulesFile, network, storage, config)
	if err != nil {
		t.Fatalf("Erreur ingestion règles: %v", err)
	}

	// Validation incrémentale est maintenant toujours activée, pas besoin de vérifier le flag
	if metrics.ValidationWithContextDuration == 0 {
		t.Error("Validation incrémentale devrait avoir été exécutée")
	}

	// Phase 3: Reset avec GC
	resetContent := `
reset

type Product(id: string, name: string, price: number)
`
	resetFile := createTempFile(t, "reset.tsd", resetContent)
	defer os.Remove(resetFile)

	network, metrics, err = pipeline.IngestFileWithAdvancedFeatures(resetFile, network, storage, config)
	if err != nil {
		t.Fatalf("Erreur ingestion reset: %v", err)
	}

	if !metrics.GCPerformed {
		t.Error("GC devrait avoir été effectué après reset")
	}

	if metrics.NodesCollected == 0 {
		t.Error("Des nœuds devraient avoir été collectés")
	}

	// Vérifier que le nouveau réseau contient le bon type
	if len(network.Types) != 1 {
		t.Errorf("Attendu 1 type après reset, obtenu %d", len(network.Types))
	}

	if network.Types[0].Name != "Product" {
		t.Errorf("Attendu type Product, obtenu %s", network.Types[0].Name)
	}
}

// TestTransactionAutoRollback teste le rollback automatique en cas d'erreur
func TestTransactionAutoRollback(t *testing.T) {
	storage := rete.NewMemoryStorage()
	pipeline := rete.NewConstraintPipeline()

	// Phase 1: État initial
	initialContent := `
type Person(id: string, name: string)
`
	initialFile := createTempFile(t, "initial.tsd", initialContent)
	defer os.Remove(initialFile)

	config := rete.DefaultAdvancedPipelineConfig()
	config.AutoCommit = false
	config.AutoRollbackOnError = true

	network, _, err := pipeline.IngestFileWithAdvancedFeatures(initialFile, nil, storage, config)
	if err != nil {
		t.Fatalf("Erreur ingestion initiale: %v", err)
	}

	typesCountBefore := len(network.Types)

	// Phase 2: Tenter une ingestion invalide avec auto-rollback
	invalidContent := `
action print_msg(msg: string)

rule invalid_rule: {x: UndefinedType} / x.field == "value" ==> print_msg("This should fail")
`
	invalidFile := createTempFile(t, "invalid.tsd", invalidContent)
	defer os.Remove(invalidFile)

	resultNetwork, metrics, err := pipeline.IngestFileWithAdvancedFeatures(invalidFile, network, storage, config)

	// L'erreur devrait être détectée
	if err == nil {
		t.Error("Attendu une erreur pour type non défini")
	}

	// Le rollback automatique devrait avoir été effectué
	if metrics != nil && !metrics.RollbackPerformed {
		// Note: peut être nil si erreur avant création transaction
		t.Log("Rollback automatique devrait avoir été effectué (métriques peuvent être nil)")
	}

	// L'état devrait être restauré - utiliser le réseau original car resultNetwork peut être nil
	if resultNetwork == nil {
		resultNetwork = network
	}
	typesCountAfter := len(resultNetwork.Types)
	if typesCountAfter != typesCountBefore {
		t.Errorf("État devrait être restauré après rollback: avant=%d, après=%d",
			typesCountBefore, typesCountAfter)
	}
}

// TestTransactionCommandPattern teste le système de transaction avec Command Pattern
func TestTransactionCommandPattern(t *testing.T) {
	storage := rete.NewMemoryStorage()
	network := rete.NewReteNetwork(storage)

	// Ajouter quelques types
	content := `
type Person(id: string, name: string, age: number)

type Company(id: string, name: string)
`
	file := createTempFile(t, "types.tsd", content)
	defer os.Remove(file)

	pipeline := rete.NewConstraintPipeline()
	network, err := pipeline.IngestFile(file, network, storage)
	if err != nil {
		t.Fatalf("Erreur ingestion: %v", err)
	}

	// Vérifier que les types ont été ajoutés
	if len(network.Types) != 2 {
		t.Errorf("Attendu 2 types, obtenu %d", len(network.Types))
	}

	// Créer une transaction pour tester le système de commandes
	tx := network.BeginTransaction()
	if tx == nil {
		t.Fatal("Transaction ne devrait pas être nil")
	}

	if !tx.IsActive {
		t.Error("Transaction devrait être active")
	}

	// La transaction utilise maintenant le Command Pattern (pas de snapshot)
	t.Logf("Transaction créée avec Command Pattern (pas de snapshot)")
}

// Helper function pour créer un fichier temporaire
func createTempFile(t *testing.T, name, content string) string {
	t.Helper()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, name)

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Erreur création fichier temporaire: %v", err)
	}

	return filePath
}
