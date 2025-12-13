// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

//go:build performance

package performance

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/tests/shared/testutil"
)

func BenchmarkTSDExecution_Simple(b *testing.B) {
	fixture := filepath.Join(testutil.GetTestDataPath(), "fixtures/alpha/alpha_abs_positive.tsd")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipeline := rete.NewConstraintPipeline()
		storage := rete.NewMemoryStorage()
		_, _, _ = pipeline.IngestFile(fixture, nil, storage)
	}
}

func BenchmarkTSDExecution_Complex(b *testing.B) {
	fixture := filepath.Join(testutil.GetTestDataPath(), "fixtures/integration/alpha_exhaustive_coverage.tsd")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipeline := rete.NewConstraintPipeline()
		storage := rete.NewMemoryStorage()
		_, _, _ = pipeline.IngestFile(fixture, nil, storage)
	}
}

func BenchmarkParallel(b *testing.B) {
	fixture := filepath.Join(testutil.GetTestDataPath(), "fixtures/alpha/alpha_abs_positive.tsd")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pipeline := rete.NewConstraintPipeline()
			storage := rete.NewMemoryStorage()
			_, _, _ = pipeline.IngestFile(fixture, nil, storage)
		}
	})
}

func BenchmarkTSDExecution_AlphaFixtures(b *testing.B) {
	fixtures := []string{
		"fixtures/alpha/alpha_abs_positive.tsd",
		"fixtures/alpha/alpha_abs_negative.tsd",
		"fixtures/alpha/alpha_addition.tsd",
		"fixtures/alpha/alpha_subtraction.tsd",
		"fixtures/alpha/alpha_multiplication.tsd",
	}

	for _, fixture := range fixtures {
		fixturePath := filepath.Join(testutil.GetTestDataPath(), fixture)
		b.Run(filepath.Base(fixture), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				pipeline := rete.NewConstraintPipeline()
				storage := rete.NewMemoryStorage()
				_, _, _ = pipeline.IngestFile(fixturePath, nil, storage)
			}
		})
	}
}

func BenchmarkTSDExecution_BetaFixtures(b *testing.B) {
	fixtures := []string{
		"fixtures/beta/beta_basic_join.tsd",
		"fixtures/beta/beta_multiple_joins.tsd",
		"fixtures/beta/beta_complex_constraints.tsd",
	}

	for _, fixture := range fixtures {
		fixturePath := filepath.Join(testutil.GetTestDataPath(), fixture)
		b.Run(filepath.Base(fixture), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				pipeline := rete.NewConstraintPipeline()
				storage := rete.NewMemoryStorage()
				_, _, _ = pipeline.IngestFile(fixturePath, nil, storage)
			}
		})
	}
}

func BenchmarkPipelineCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rete.NewConstraintPipeline()
	}
}

func BenchmarkStorageCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rete.NewMemoryStorage()
	}
}

func BenchmarkFactProcessing_10Facts(b *testing.B) {
	rule := generateBenchmarkRule(10)
	tempFile := testutil.CreateTempTSDFile(&testing.T{}, rule)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipeline := rete.NewConstraintPipeline()
		storage := rete.NewMemoryStorage()
		_, _, _ = pipeline.IngestFile(tempFile, nil, storage)
	}
}

func BenchmarkFactProcessing_100Facts(b *testing.B) {
	rule := generateBenchmarkRule(100)
	tempFile := testutil.CreateTempTSDFile(&testing.T{}, rule)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipeline := rete.NewConstraintPipeline()
		storage := rete.NewMemoryStorage()
		_, _, _ = pipeline.IngestFile(tempFile, nil, storage)
	}
}

func BenchmarkFactProcessing_1000Facts(b *testing.B) {
	rule := generateBenchmarkRule(1000)
	tempFile := testutil.CreateTempTSDFile(&testing.T{}, rule)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipeline := rete.NewConstraintPipeline()
		storage := rete.NewMemoryStorage()
		_, _, _ = pipeline.IngestFile(tempFile, nil, storage)
	}
}

func BenchmarkConstraintEvaluation_Simple(b *testing.B) {
	rule := `type Item(value: number)

rule r1 : {i: Item} / i.value > 0 ==> print("positive")

Item(value:1)
Item(value:2)
Item(value:3)
`

	tempFile := testutil.CreateTempTSDFile(&testing.T{}, rule)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipeline := rete.NewConstraintPipeline()
		storage := rete.NewMemoryStorage()
		_, _, _ = pipeline.IngestFile(tempFile, nil, storage)
	}
}

func BenchmarkConstraintEvaluation_Complex(b *testing.B) {
	rule := `type Record(a: number, b: number, c: number, d: bool, e: string)

rule r1 : {r: Record} /
    r.a > 10 and
    r.b < 100 and
    r.c >= 50 and
    r.d == true
    ==> print("matched")

Record(a:20, b:80, c:60, d:true, e:"test1")
Record(a:15, b:90, c:55, d:true, e:"test2")
Record(a:25, b:75, c:65, d:true, e:"test3")
`

	tempFile := testutil.CreateTempTSDFile(&testing.T{}, rule)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipeline := rete.NewConstraintPipeline()
		storage := rete.NewMemoryStorage()
		_, _, _ = pipeline.IngestFile(tempFile, nil, storage)
	}
}

func BenchmarkJoinOperations_TwoTypes(b *testing.B) {
	rule := `type Person(id: number, name: string)
type Company(id: number, person_id: number, name: string)

rule r1 : {p: Person, c: Company} / p.id == c.person_id ==> print("match")

Person(id:1, name:"Alice")
Person(id:2, name:"Bob")
Company(id:1, person_id:1, name:"ACME")
Company(id:2, person_id:2, name:"TechCorp")
`

	tempFile := testutil.CreateTempTSDFile(&testing.T{}, rule)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipeline := rete.NewConstraintPipeline()
		storage := rete.NewMemoryStorage()
		_, _, _ = pipeline.IngestFile(tempFile, nil, storage)
	}
}

func BenchmarkJoinOperations_ThreeTypes(b *testing.B) {
	rule := `type A(id: number, value: string)
type B(id: number, a_id: number, value: string)
type C(id: number, b_id: number, value: string)

rule r1 : {a: A, b: B, c: C} /
    a.id == b.a_id and
    b.id == c.b_id
    ==> print("three_way_join")

A(id:1, value:"a1")
A(id:2, value:"a2")
B(id:10, a_id:1, value:"b1")
B(id:20, a_id:2, value:"b2")
C(id:100, b_id:10, value:"c1")
C(id:200, b_id:20, value:"c2")
`

	tempFile := testutil.CreateTempTSDFile(&testing.T{}, rule)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipeline := rete.NewConstraintPipeline()
		storage := rete.NewMemoryStorage()
		_, _, _ = pipeline.IngestFile(tempFile, nil, storage)
	}
}

func BenchmarkMultipleRules_Sequential(b *testing.B) {
	rule := `type Data(value: number)

rule r1 : {d: Data} / d.value > 0 ==> print("positive")
rule r2 : {d: Data} / d.value < 0 ==> print("negative")
rule r3 : {d: Data} / d.value == 0 ==> print("zero")

Data(value:5)
Data(value:-3)
Data(value:0)
Data(value:10)
Data(value:-7)
`

	tempFile := testutil.CreateTempTSDFile(&testing.T{}, rule)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipeline := rete.NewConstraintPipeline()
		storage := rete.NewMemoryStorage()
		_, _, _ = pipeline.IngestFile(tempFile, nil, storage)
	}
}

func BenchmarkTypeSystem_SingleField(b *testing.B) {
	rule := `type Simple(value: number)

rule r1 : {s: Simple} ==> print("matched")

Simple(value:1)
Simple(value:2)
Simple(value:3)
`

	tempFile := testutil.CreateTempTSDFile(&testing.T{}, rule)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipeline := rete.NewConstraintPipeline()
		storage := rete.NewMemoryStorage()
		_, _, _ = pipeline.IngestFile(tempFile, nil, storage)
	}
}

func BenchmarkTypeSystem_ManyFields(b *testing.B) {
	rule := `type Complex(
    f1: number, f2: string, f3: bool, f4: number, f5: string,
    f6: bool, f7: number, f8: string, f9: bool, f10: number
)

rule r1 : {c: Complex} ==> print("matched")

Complex(f1:1, f2:"a", f3:true, f4:2, f5:"b", f6:false, f7:3, f8:"c", f9:true, f10:4)
Complex(f1:5, f2:"d", f3:false, f4:6, f5:"e", f6:true, f7:7, f8:"f", f9:false, f10:8)
`

	tempFile := testutil.CreateTempTSDFile(&testing.T{}, rule)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pipeline := rete.NewConstraintPipeline()
		storage := rete.NewMemoryStorage()
		_, _, _ = pipeline.IngestFile(tempFile, nil, storage)
	}
}

func BenchmarkMemoryAllocation(b *testing.B) {
	rule := generateBenchmarkRule(50)
	tempFile := testutil.CreateTempTSDFile(&testing.T{}, rule)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		pipeline := rete.NewConstraintPipeline()
		storage := rete.NewMemoryStorage()
		_, _, _ = pipeline.IngestFile(tempFile, nil, storage)
	}
}

// generateBenchmarkRule creates a simple TSD rule with the specified number of facts
func generateBenchmarkRule(factCount int) string {
	rule := `type BenchItem(id: number, value: number)

rule r1 : {bi: BenchItem} / bi.value > 0 ==> print("positive")

`

	for i := 0; i < factCount; i++ {
		rule += fmt.Sprintf("BenchItem(id:%d, value:%d)\n", i, i+1)
	}

	return rule
}
