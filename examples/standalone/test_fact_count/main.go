package main

import (
	"fmt"
	"os"

	"github.com/treivax/tsd/rete"
)

func main() {
	rule := `type Person(name: string, age: number)

action print(arg1: string)

rule r1 : {p: Person} / p.age > 18 ==> print("adult")

Person(name:"Alice", age:25)
Person(name:"Bob", age:15)
`

	// Create temp file
	tmpfile, err := os.CreateTemp("", "test*.tsd")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(rule)); err != nil {
		panic(err)
	}
	if err := tmpfile.Close(); err != nil {
		panic(err)
	}

	// Execute
	pipeline := rete.NewConstraintPipeline()
	storage := rete.NewMemoryStorage()

	network, _, err := pipeline.IngestFile(tmpfile.Name(), nil, storage)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	facts := storage.GetAllFacts()
	fmt.Printf("Facts in storage: %d\n", len(facts))
	for i, f := range facts {
		fmt.Printf("  Fact %d: %s (Type: %s, Fields: %+v)\n", i, f.ID, f.Type, f.Fields)
	}

	fmt.Printf("TypeNodes: %d\n", len(network.TypeNodes))
	fmt.Printf("TerminalNodes: %d\n", len(network.TerminalNodes))
}
