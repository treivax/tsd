package constraint

import (
	"testing"
)

func TestCommentsWithSlashes(t *testing.T) {
	input := `
// Commentaire ligne
/* Commentaire
   bloc */
type Person(name: string)
Person(name: "Alice")
`
	_, err := ParseConstraint("test.tsd", []byte(input))
	if err != nil {
		t.Fatalf("Les commentaires // et /* */ devraient fonctionner: %v", err)
	}
	t.Log("✅ Commentaires // et /* */ fonctionnent")
}

func TestHashAsCommentShouldFail(t *testing.T) {
	input := `
# Ce commentaire ne devrait plus fonctionner
type Person(name: string)
`
	_, err := ParseConstraint("test.tsd", []byte(input))
	if err == nil {
		t.Fatal("❌ Le commentaire # devrait être rejeté")
	}
	t.Logf("✅ Le commentaire # est bien rejeté: %v", err)
}

func TestHashAsPrimaryKeyStillWorks(t *testing.T) {
	input := `type Person(#name: string, age: number)`
	result, err := ParseConstraint("test.tsd", []byte(input))
	if err != nil {
		t.Fatalf("# comme marqueur de clé primaire devrait fonctionner: %v", err)
	}
	t.Log("✅ # comme marqueur de clé primaire fonctionne")
	t.Logf("Résultat: %+v", result)
}
