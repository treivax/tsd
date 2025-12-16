// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"reflect"
	"testing"
)

func TestTypeDefinitionClone(t *testing.T) {
	t.Log("🧪 TEST TypeDefinition.Clone()")
	t.Log("==============================")

	t.Run("clone avec champs multiples", func(t *testing.T) {
		original := TypeDefinition{
			Type: "typeDefinition",
			Name: "User",
			Fields: []Field{
				{Name: "id", Type: "string"},
				{Name: "email", Type: "string"},
				{Name: "age", Type: "number"},
			},
		}

		clone := original.Clone()

		// Vérifier que les valeurs sont identiques
		if clone.Type != original.Type {
			t.Errorf("❌ Type: attendu %q, reçu %q", original.Type, clone.Type)
		}
		if clone.Name != original.Name {
			t.Errorf("❌ Name: attendu %q, reçu %q", original.Name, clone.Name)
		}
		if len(clone.Fields) != len(original.Fields) {
			t.Errorf("❌ Fields length: attendu %d, reçu %d", len(original.Fields), len(clone.Fields))
		}

		// Vérifier que les champs sont identiques
		for i, field := range original.Fields {
			if clone.Fields[i] != field {
				t.Errorf("❌ Field[%d]: attendu %+v, reçu %+v", i, field, clone.Fields[i])
			}
		}

		// Vérifier l'indépendance (copie profonde)
		clone.Name = "ModifiedUser"
		if original.Name == clone.Name {
			t.Error("❌ La modification du clone a affecté l'original (Name)")
		}

		clone.Fields[0].Name = "userId"
		if original.Fields[0].Name == clone.Fields[0].Name {
			t.Error("❌ La modification du clone a affecté l'original (Fields)")
		}

		t.Log("✅ Clone avec champs multiples réussi")
	})

	t.Run("clone avec champs vides", func(t *testing.T) {
		original := TypeDefinition{
			Type:   "typeDefinition",
			Name:   "EmptyType",
			Fields: []Field{},
		}

		clone := original.Clone()

		if clone.Type != original.Type {
			t.Errorf("❌ Type: attendu %q, reçu %q", original.Type, clone.Type)
		}
		if clone.Name != original.Name {
			t.Errorf("❌ Name: attendu %q, reçu %q", original.Name, clone.Name)
		}
		if len(clone.Fields) != 0 {
			t.Errorf("❌ Fields devrait être vide, longueur: %d", len(clone.Fields))
		}

		t.Log("✅ Clone avec champs vides réussi")
	})

	t.Run("clone sans initialisation de Fields", func(t *testing.T) {
		original := TypeDefinition{
			Type:   "typeDefinition",
			Name:   "NoFieldsType",
			Fields: nil,
		}

		clone := original.Clone()

		if clone.Type != original.Type {
			t.Errorf("❌ Type: attendu %q, reçu %q", original.Type, clone.Type)
		}
		if clone.Name != original.Name {
			t.Errorf("❌ Name: attendu %q, reçu %q", original.Name, clone.Name)
		}
		if clone.Fields == nil {
			t.Error("❌ Fields ne devrait pas être nil après clone")
		}
		if len(clone.Fields) != 0 {
			t.Errorf("❌ Fields devrait être vide, longueur: %d", len(clone.Fields))
		}

		t.Log("✅ Clone sans Fields initialisés réussi")
	})

	t.Run("indépendance du slice Fields", func(t *testing.T) {
		original := TypeDefinition{
			Type: "typeDefinition",
			Name: "TestType",
			Fields: []Field{
				{Name: "field1", Type: "string"},
			},
		}

		clone := original.Clone()

		// Ajouter un champ au clone
		clone.Fields = append(clone.Fields, Field{Name: "field2", Type: "number"})

		// Vérifier que l'original n'est pas affecté
		if len(original.Fields) != 1 {
			t.Errorf("❌ L'ajout au clone a affecté l'original: longueur = %d", len(original.Fields))
		}

		t.Log("✅ Indépendance du slice Fields vérifiée")
	})
}

func TestActionClone(t *testing.T) {
	t.Log("🧪 TEST Action.Clone()")
	t.Log("======================")

	t.Run("clone avec Jobs multiples", func(t *testing.T) {
		original := &Action{
			Type: "action",
			Jobs: []JobCall{
				{Type: "jobCall", Name: "notify", Args: []interface{}{"message1"}},
				{Type: "jobCall", Name: "log", Args: []interface{}{"info", "test"}},
			},
			Job: nil,
		}

		clone := original.Clone()

		// Vérifier que les valeurs sont identiques
		if clone.Type != original.Type {
			t.Errorf("❌ Type: attendu %q, reçu %q", original.Type, clone.Type)
		}
		if len(clone.Jobs) != len(original.Jobs) {
			t.Errorf("❌ Jobs length: attendu %d, reçu %d", len(original.Jobs), len(clone.Jobs))
		}

		// Vérifier que les jobs sont identiques
		for i, job := range original.Jobs {
			if !reflect.DeepEqual(clone.Jobs[i], job) {
				t.Errorf("❌ Job[%d]: attendu %+v, reçu %+v", i, job, clone.Jobs[i])
			}
		}

		// Vérifier l'indépendance
		clone.Jobs[0].Name = "modified"
		if original.Jobs[0].Name == clone.Jobs[0].Name {
			t.Error("❌ La modification du clone a affecté l'original (Jobs)")
		}

		t.Log("✅ Clone avec Jobs multiples réussi")
	})

	t.Run("clone avec Job unique (backward compatibility)", func(t *testing.T) {
		singleJob := &JobCall{
			Type: "jobCall",
			Name: "notify",
			Args: []interface{}{"test message"},
		}

		original := &Action{
			Type: "action",
			Job:  singleJob,
			Jobs: []JobCall{},
		}

		clone := original.Clone()

		// Vérifier que Job est cloné
		if clone.Job == nil {
			t.Fatal("❌ Job ne devrait pas être nil")
		}
		if clone.Job == original.Job {
			t.Error("❌ Job devrait être une copie, pas la même référence")
		}
		if !reflect.DeepEqual(*clone.Job, *original.Job) {
			t.Errorf("❌ Job: attendu %+v, reçu %+v", *original.Job, *clone.Job)
		}

		// Vérifier l'indépendance
		clone.Job.Name = "modified"
		if original.Job.Name == clone.Job.Name {
			t.Error("❌ La modification du clone a affecté l'original (Job)")
		}

		t.Log("✅ Clone avec Job unique réussi")
	})

	t.Run("clone avec Job et Jobs combinés", func(t *testing.T) {
		singleJob := &JobCall{
			Type: "jobCall",
			Name: "notify",
			Args: []interface{}{"single"},
		}

		original := &Action{
			Type: "action",
			Job:  singleJob,
			Jobs: []JobCall{
				{Type: "jobCall", Name: "log", Args: []interface{}{"multi1"}},
				{Type: "jobCall", Name: "alert", Args: []interface{}{"multi2"}},
			},
		}

		clone := original.Clone()

		// Vérifier Job
		if clone.Job == nil {
			t.Error("❌ Job ne devrait pas être nil")
		} else if clone.Job == original.Job {
			t.Error("❌ Job devrait être une nouvelle instance")
		}

		// Vérifier Jobs
		if len(clone.Jobs) != len(original.Jobs) {
			t.Errorf("❌ Jobs length: attendu %d, reçu %d", len(original.Jobs), len(clone.Jobs))
		}

		t.Log("✅ Clone avec Job et Jobs combinés réussi")
	})

	t.Run("clone avec Jobs vides", func(t *testing.T) {
		original := &Action{
			Type: "action",
			Jobs: []JobCall{},
			Job:  nil,
		}

		clone := original.Clone()

		if clone.Type != original.Type {
			t.Errorf("❌ Type: attendu %q, reçu %q", original.Type, clone.Type)
		}
		if len(clone.Jobs) != 0 {
			t.Errorf("❌ Jobs devrait être vide, longueur: %d", len(clone.Jobs))
		}
		if clone.Job != nil {
			t.Error("❌ Job devrait être nil")
		}

		t.Log("✅ Clone avec Jobs vides réussi")
	})

	t.Run("clone sans initialisation de Jobs", func(t *testing.T) {
		original := &Action{
			Type: "action",
			Jobs: nil,
			Job:  nil,
		}

		clone := original.Clone()

		if clone.Jobs == nil {
			t.Error("❌ Jobs ne devrait pas être nil après clone")
		}
		if len(clone.Jobs) != 0 {
			t.Errorf("❌ Jobs devrait être vide, longueur: %d", len(clone.Jobs))
		}

		t.Log("✅ Clone sans Jobs initialisés réussi")
	})

	t.Run("indépendance du slice Jobs", func(t *testing.T) {
		original := &Action{
			Type: "action",
			Jobs: []JobCall{
				{Type: "jobCall", Name: "job1", Args: []interface{}{"arg1"}},
			},
			Job: nil,
		}

		clone := original.Clone()

		// Ajouter un job au clone
		clone.Jobs = append(clone.Jobs, JobCall{Type: "jobCall", Name: "job2", Args: []interface{}{"arg2"}})

		// Vérifier que l'original n'est pas affecté
		if len(original.Jobs) != 1 {
			t.Errorf("❌ L'ajout au clone a affecté l'original: longueur = %d", len(original.Jobs))
		}

		t.Log("✅ Indépendance du slice Jobs vérifiée")
	})

	t.Run("clone avec Args complexes", func(t *testing.T) {
		original := &Action{
			Type: "action",
			Jobs: []JobCall{
				{
					Type: "jobCall",
					Name: "complexJob",
					Args: []interface{}{
						"string",
						123,
						true,
						map[string]interface{}{"key": "value"},
						[]interface{}{1, 2, 3},
					},
				},
			},
		}

		clone := original.Clone()

		// Vérifier que les args sont copiés
		if !reflect.DeepEqual(clone.Jobs[0].Args, original.Jobs[0].Args) {
			t.Error("❌ Args ne sont pas identiques après clone")
		}

		t.Log("✅ Clone avec Args complexes réussi")
	})
}

func TestActionCloneNilSafety(t *testing.T) {
	t.Log("🧪 TEST Action.Clone() - Sécurité nil")
	t.Log("======================================")

	t.Run("clone d'une Action avec tous les champs nil", func(t *testing.T) {
		original := &Action{
			Type: "action",
			Job:  nil,
			Jobs: nil,
		}

		clone := original.Clone()

		if clone == nil {
			t.Fatal("❌ Clone ne devrait pas être nil")
		}
		if clone.Type != original.Type {
			t.Errorf("❌ Type: attendu %q, reçu %q", original.Type, clone.Type)
		}
		if clone.Job != nil {
			t.Error("❌ Job devrait être nil")
		}
		if clone.Jobs == nil {
			t.Error("❌ Jobs ne devrait pas être nil (slice vide attendu)")
		}
		if len(clone.Jobs) != 0 {
			t.Errorf("❌ Jobs devrait être vide, longueur: %d", len(clone.Jobs))
		}

		t.Log("✅ Clone avec champs nil est sécurisé")
	})
}

func TestTypeDefinitionCloneImmutability(t *testing.T) {
	t.Log("🧪 TEST TypeDefinition.Clone() - Immutabilité")
	t.Log("==============================================")

	original := TypeDefinition{
		Type: "typeDefinition",
		Name: "Product",
		Fields: []Field{
			{Name: "id", Type: "string"},
			{Name: "price", Type: "number"},
			{Name: "available", Type: "boolean"},
		},
	}

	clone := original.Clone()

	// Modifications du clone
	clone.Type = "modifiedType"
	clone.Name = "ModifiedProduct"
	clone.Fields[0].Name = "productId"
	clone.Fields[1].Type = "decimal"
	clone.Fields = append(clone.Fields, Field{Name: "description", Type: "string"})

	// Vérifier que l'original n'est pas modifié
	if original.Type != "typeDefinition" {
		t.Errorf("❌ Type original modifié: %q", original.Type)
	}
	if original.Name != "Product" {
		t.Errorf("❌ Name original modifié: %q", original.Name)
	}
	if original.Fields[0].Name != "id" {
		t.Errorf("❌ Field[0].Name original modifié: %q", original.Fields[0].Name)
	}
	if original.Fields[1].Type != "number" {
		t.Errorf("❌ Field[1].Type original modifié: %q", original.Fields[1].Type)
	}
	if len(original.Fields) != 3 {
		t.Errorf("❌ Longueur Fields original modifiée: %d", len(original.Fields))
	}

	t.Log("✅ L'original reste immutable après modifications du clone")
}

func TestActionCloneImmutability(t *testing.T) {
	t.Log("🧪 TEST Action.Clone() - Immutabilité")
	t.Log("=====================================")

	singleJob := &JobCall{
		Type: "jobCall",
		Name: "original",
		Args: []interface{}{"arg1", "arg2"},
	}

	original := &Action{
		Type: "action",
		Job:  singleJob,
		Jobs: []JobCall{
			{Type: "jobCall", Name: "job1", Args: []interface{}{"a"}},
			{Type: "jobCall", Name: "job2", Args: []interface{}{"b"}},
		},
	}

	clone := original.Clone()

	// Modifications du clone
	clone.Type = "modifiedAction"
	if clone.Job != nil {
		clone.Job.Name = "modified"
		clone.Job.Args = []interface{}{"new1", "new2"}
	}
	clone.Jobs[0].Name = "modifiedJob1"
	clone.Jobs = append(clone.Jobs, JobCall{Type: "jobCall", Name: "job3", Args: []interface{}{"c"}})

	// Vérifier que l'original n'est pas modifié
	if original.Type != "action" {
		t.Errorf("❌ Type original modifié: %q", original.Type)
	}
	if original.Job == nil {
		t.Fatal("❌ Job original ne devrait pas être nil")
	}
	if original.Job.Name != "original" {
		t.Errorf("❌ Job.Name original modifié: %q", original.Job.Name)
	}
	if len(original.Job.Args) != 2 {
		t.Errorf("❌ Job.Args original modifié: longueur %d", len(original.Job.Args))
	}
	if original.Jobs[0].Name != "job1" {
		t.Errorf("❌ Jobs[0].Name original modifié: %q", original.Jobs[0].Name)
	}
	if len(original.Jobs) != 2 {
		t.Errorf("❌ Longueur Jobs original modifiée: %d", len(original.Jobs))
	}

	t.Log("✅ L'original reste immutable après modifications du clone")
}
