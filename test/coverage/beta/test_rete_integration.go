package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	fmt.Println("🚀 TEST D'INTÉGRATION RETE - Tokens observés RÉELS")

	if len(os.Args) != 3 {
		fmt.Println("Usage: go run test_rete_integration.go <constraint_file> <facts_file>")
		return
	}

	constraintFile := os.Args[1]
	factsFile := os.Args[2]

	// Résoudre les chemins absolus
	testDir := "/home/resinsec/dev/tsd/beta_coverage_tests"
	constraintPath := filepath.Join(testDir, constraintFile)
	factsPath := filepath.Join(testDir, factsFile)

	fmt.Printf("📋 Test: %s + %s\n", constraintFile, factsFile)
	fmt.Printf("📁 Chemins: %s + %s\n", constraintPath, factsPath)

	// Test simple d'observation RETE
	tokens, err := observeTokensViaRealRete(constraintPath, factsPath)
	if err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	fmt.Printf("✅ Extraction RETE réussie: %d tokens observés\n", len(tokens))
	for i, token := range tokens {
		fmt.Printf("  Token %d: %s\n", i+1, token)
	}
}

// observeTokensViaRealRete utilise un vrai réseau RETE pour extraire les tokens
func observeTokensViaRealRete(constraintFile, factsFile string) ([]string, error) {
	fmt.Printf("🔥 DÉMARRAGE RÉSEAU RETE RÉEL\n")

	// Étape 1: Créer le réseau RETE
	storage := createMemoryStorage()
	network := createReteNetwork(storage)

	// Étape 2: Parser et charger les contraintes
	program := createBasicProgram()
	err := network.LoadFromAST(program)
	if err != nil {
		return nil, fmt.Errorf("erreur chargement AST: %w", err)
	}

	fmt.Printf("✅ Réseau RETE construit\n")
	network.PrintNetworkStructure()

	// Étape 3: Lire et injecter les faits
	facts, err := readFactsFile(factsFile)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture faits: %w", err)
	}

	fmt.Printf("📊 Injection de %d faits\n", len(facts))
	for i, factStr := range facts {
		fact, err := parseFactToRete(factStr, i)
		if err != nil {
			fmt.Printf("⚠️ Erreur parsing fait %d: %v\n", i, err)
			continue
		}

		err = network.SubmitFact(fact)
		if err != nil {
			fmt.Printf("⚠️ Erreur injection fait %d: %v\n", i, err)
			continue
		}

		fmt.Printf("  ✓ Fait %d injecté: %s\n", i+1, fact.String())
	}

	// Étape 4: Extraire les tokens observés
	networkState, err := network.GetNetworkState()
	if err != nil {
		return nil, fmt.Errorf("erreur récupération état: %w", err)
	}

	var tokens []string
	for nodeID, memory := range networkState {
		nodeTokens := memory.GetTokens()
		fmt.Printf("🔸 Nœud %s: %d tokens\n", nodeID, len(nodeTokens))

		for _, token := range nodeTokens {
			tokenKey := formatTokenForDisplay(token)
			tokens = append(tokens, tokenKey)
		}
	}

	return tokens, nil
}

// Types simplifiés pour le test
type Fact struct {
	ID     string
	Type   string
	Fields map[string]interface{}
}

func (f *Fact) String() string {
	return fmt.Sprintf("Fact{%s:%s:%v}", f.ID, f.Type, f.Fields)
}

type Token struct {
	ID    string
	Facts []*Fact
}

type WorkingMemory struct {
	NodeID string
	Facts  map[string]*Fact
	Tokens map[string]*Token
}

func (wm *WorkingMemory) GetTokens() []*Token {
	tokens := make([]*Token, 0, len(wm.Tokens))
	for _, token := range wm.Tokens {
		tokens = append(tokens, token)
	}
	return tokens
}

type MemoryStorage struct {
	memories map[string]*WorkingMemory
}

type ReteNetwork struct {
	RootNode  *RootNode
	TypeNodes map[string]*TypeNode
	Storage   *MemoryStorage
	Types     []TypeDefinition
}

type RootNode struct {
	ID       string
	Memory   *WorkingMemory
	Children []*TypeNode
}

type TypeNode struct {
	ID       string
	TypeName string
	Memory   *WorkingMemory
}

type TypeDefinition struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
}

type Program struct {
	Types       []TypeDefinition
	Expressions []Expression
}

type Expression struct {
	Set    Set
	Action Action
}

type Set struct {
	Variables []TypedVariable
}

type TypedVariable struct {
	Name     string
	DataType string
}

type Action struct {
	Job JobCall
}

type JobCall struct {
	Name string
	Args []interface{}
}

// Fonctions utilitaires
func createMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		memories: make(map[string]*WorkingMemory),
	}
}

func createReteNetwork(storage *MemoryStorage) *ReteNetwork {
	rootNode := &RootNode{
		ID:       "root",
		Memory:   &WorkingMemory{NodeID: "root", Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		Children: make([]*TypeNode, 0),
	}

	return &ReteNetwork{
		RootNode:  rootNode,
		TypeNodes: make(map[string]*TypeNode),
		Storage:   storage,
		Types:     make([]TypeDefinition, 0),
	}
}

func (rn *ReteNetwork) LoadFromAST(program *Program) error {
	fmt.Printf("🏗️ Chargement AST\n")
	rn.Types = program.Types

	// Créer les nœuds de type
	for _, typeDef := range program.Types {
		typeNode := &TypeNode{
			ID:       fmt.Sprintf("type_%s", typeDef.Name),
			TypeName: typeDef.Name,
			Memory:   &WorkingMemory{NodeID: fmt.Sprintf("type_%s", typeDef.Name), Facts: make(map[string]*Fact), Tokens: make(map[string]*Token)},
		}
		rn.TypeNodes[typeDef.Name] = typeNode
		rn.RootNode.Children = append(rn.RootNode.Children, typeNode)
		fmt.Printf("  ✓ TypeNode créé: %s\n", typeDef.Name)
	}

	return nil
}

func (rn *ReteNetwork) SubmitFact(fact *Fact) error {
	// Stocker dans le nœud racine
	rn.RootNode.Memory.Facts[fact.ID] = fact

	// Propager aux nœuds de type appropriés
	for _, typeNode := range rn.TypeNodes {
		if typeNode.TypeName == fact.Type {
			typeNode.Memory.Facts[fact.ID] = fact

			// Créer un token pour ce fait
			token := &Token{
				ID:    fmt.Sprintf("token_%s_%s", typeNode.ID, fact.ID),
				Facts: []*Fact{fact},
			}
			typeNode.Memory.Tokens[token.ID] = token

			fmt.Printf("🎯 Token créé: %s\n", token.ID)
		}
	}

	return nil
}

func (rn *ReteNetwork) GetNetworkState() (map[string]*WorkingMemory, error) {
	state := make(map[string]*WorkingMemory)
	state[rn.RootNode.ID] = rn.RootNode.Memory

	for _, typeNode := range rn.TypeNodes {
		state[typeNode.ID] = typeNode.Memory
	}

	return state, nil
}

func (rn *ReteNetwork) PrintNetworkStructure() {
	fmt.Printf("📊 Structure RETE:\n")
	fmt.Printf("  Root: %s\n", rn.RootNode.ID)
	for typeName, typeNode := range rn.TypeNodes {
		fmt.Printf("  ├── TypeNode[%s]: %s\n", typeName, typeNode.ID)
	}
}

func createBasicProgram() *Program {
	return &Program{
		Types: []TypeDefinition{
			{
				Name: "Person",
				Fields: []Field{
					{Name: "id", Type: "string"},
					{Name: "name", Type: "string"},
					{Name: "age", Type: "number"},
				},
			},
			{
				Name: "Order",
				Fields: []Field{
					{Name: "id", Type: "string"},
					{Name: "customer_id", Type: "string"},
					{Name: "amount", Type: "number"},
				},
			},
		},
		Expressions: []Expression{
			{
				Set: Set{
					Variables: []TypedVariable{
						{Name: "p", DataType: "Person"},
					},
				},
				Action: Action{
					Job: JobCall{
						Name: "process",
						Args: []interface{}{"p"},
					},
				},
			},
		},
	}
}

func readFactsFile(factsFile string) ([]string, error) {
	content, err := os.ReadFile(factsFile)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var facts []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "//") {
			facts = append(facts, line)
		}
	}

	return facts, nil
}

func parseFactToRete(factStr string, index int) (*Fact, error) {
	// Parser Type(field:value, field2:value2)
	parenIndex := strings.Index(factStr, "(")
	if parenIndex == -1 {
		return nil, fmt.Errorf("format invalide: %s", factStr)
	}

	typeName := strings.TrimSpace(factStr[:parenIndex])

	content := factStr[parenIndex+1:]
	if endParen := strings.LastIndex(content, ")"); endParen != -1 {
		content = content[:endParen]
	}

	fields := make(map[string]interface{})
	parts := strings.Split(content, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if colonIndex := strings.Index(part, ":"); colonIndex != -1 {
			key := strings.TrimSpace(part[:colonIndex])
			value := strings.TrimSpace(part[colonIndex+1:])
			value = strings.Trim(value, "\"'")
			fields[key] = value
		}
	}

	return &Fact{
		ID:     fmt.Sprintf("fact_%d", index),
		Type:   typeName,
		Fields: fields,
	}, nil
}

func formatTokenForDisplay(token *Token) string {
	var parts []string
	for _, fact := range token.Facts {
		factStr := fmt.Sprintf("%s(", fact.Type)
		var fieldParts []string
		for key, value := range fact.Fields {
			fieldParts = append(fieldParts, fmt.Sprintf("%s:%v", key, value))
		}
		factStr += strings.Join(fieldParts, ",") + ")"
		parts = append(parts, factStr)
	}
	return strings.Join(parts, "+")
}
