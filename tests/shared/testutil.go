// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package shared

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/treivax/tsd/api"
	"github.com/treivax/tsd/rete"
	"github.com/treivax/tsd/xuples"
)

// CreatePipelineFromTSD crée un pipeline à partir d'un contenu TSD.
// Retourne le pipeline et le résultat de l'ingestion.
func CreatePipelineFromTSD(t *testing.T, tsdContent string) (*api.Pipeline, *api.Result) {
	tmpfile, err := os.CreateTemp("", "test_*.tsd")
	require.NoError(t, err, "échec création fichier temporaire")

	t.Cleanup(func() {
		os.Remove(tmpfile.Name())
	})

	_, err = tmpfile.WriteString(tsdContent)
	require.NoError(t, err, "échec écriture fichier")
	tmpfile.Close()

	pipeline := api.NewPipeline()
	require.NotNil(t, pipeline, "pipeline ne doit pas être nil")

	result, err := pipeline.IngestFile(tmpfile.Name())
	require.NoError(t, err, "échec ingestion fichier")
	require.NotNil(t, result, "result ne doit pas être nil")

	return pipeline, result
}

// CreatePipelineFromFile crée un pipeline à partir d'un fichier TSD existant.
func CreatePipelineFromFile(t *testing.T, filepath string) (*api.Pipeline, *api.Result) {
	pipeline := api.NewPipeline()
	require.NotNil(t, pipeline, "pipeline ne doit pas être nil")

	result, err := pipeline.IngestFile(filepath)
	require.NoError(t, err, "échec ingestion fichier %s", filepath)
	require.NotNil(t, result, "result ne doit pas être nil")

	return pipeline, result
}

// AssertXupleFields vérifie qu'un xuple a les champs attendus.
func AssertXupleFields(t *testing.T, xuple *xuples.Xuple, expectedType string, expectedFields map[string]interface{}) {
	require.NotNil(t, xuple, "xuple ne doit pas être nil")
	require.NotNil(t, xuple.Fact, "xuple.Fact ne doit pas être nil")
	require.Equal(t, expectedType, xuple.Fact.Type, "type du xuple incorrect")

	for field, expectedValue := range expectedFields {
		actualValue := xuple.Fact.Fields[field]
		require.Equal(t, expectedValue, actualValue, "champ '%s' incorrect", field)
	}
}

// AssertXupleSpaceExists vérifie qu'un xuple-space existe.
func AssertXupleSpaceExists(t *testing.T, result *api.Result, spaceName string) {
	spaces := result.XupleSpaceNames()
	require.Contains(t, spaces, spaceName, "xuple-space '%s' devrait exister", spaceName)
}

// AssertXupleSpaceNotExists vérifie qu'un xuple-space n'existe pas.
func AssertXupleSpaceNotExists(t *testing.T, result *api.Result, spaceName string) {
	spaces := result.XupleSpaceNames()
	require.NotContains(t, spaces, spaceName, "xuple-space '%s' ne devrait pas exister", spaceName)
}

// AssertXupleCount vérifie le nombre de xuples dans un xuple-space.
func AssertXupleCount(t *testing.T, result *api.Result, spaceName string, expectedCount int) {
	count, err := result.XupleCount(spaceName)
	require.NoError(t, err, "échec récupération count pour '%s'", spaceName)
	require.Equal(t, expectedCount, count, "nombre de xuples incorrect dans '%s'", spaceName)
}

// SubmitFact crée et soumet un fait au réseau.
// Pour les types avec clé primaire, l'ID sera généré automatiquement.
func SubmitFact(t *testing.T, result *api.Result, typeName string, data map[string]interface{}) {
	network := result.Network()
	require.NotNil(t, network, "network ne doit pas être nil")

	// Créer le fait avec un ID vide, qui sera généré par le réseau
	fact := &rete.Fact{
		ID:     "", // L'ID sera généré automatiquement selon le type
		Type:   typeName,
		Fields: data,
	}

	err := network.SubmitFact(fact)
	require.NoError(t, err, "échec soumission du fait de type '%s'", typeName)
}

// GetXuples récupère tous les xuples d'un xuple-space et vérifie qu'il n'y a pas d'erreur.
func GetXuples(t *testing.T, result *api.Result, spaceName string) []*xuples.Xuple {
	xuples, err := result.GetXuples(spaceName)
	require.NoError(t, err, "échec récupération xuples de '%s'", spaceName)
	return xuples
}

// RetrieveXuple récupère un xuple d'un xuple-space.
func RetrieveXuple(t *testing.T, result *api.Result, spaceName string, agentID string) *xuples.Xuple {
	xuple, err := result.Retrieve(spaceName, agentID)
	require.NoError(t, err, "échec retrieve de '%s' par agent '%s'", spaceName, agentID)
	return xuple
}

// RetrieveAndAssert récupère un xuple et vérifie ses propriétés.
func RetrieveAndAssert(t *testing.T, result *api.Result, spaceName string, agentID string, expectedType string, expectedFields map[string]interface{}) {
	xuple := RetrieveXuple(t, result, spaceName, agentID)
	if xuple != nil {
		AssertXupleFields(t, xuple, expectedType, expectedFields)
	}
}

// AssertMetrics vérifie que les métriques sont cohérentes.
func AssertMetrics(t *testing.T, result *api.Result, minTypes, minRules, minFacts int) {
	metrics := result.Metrics()
	require.NotNil(t, metrics, "metrics ne doivent pas être nil")

	if minTypes > 0 {
		require.GreaterOrEqual(t, metrics.TypeCount, minTypes, "nombre de types insuffisant")
	}
	if minRules > 0 {
		require.GreaterOrEqual(t, metrics.RuleCount, minRules, "nombre de règles insuffisant")
	}
	if minFacts > 0 {
		require.GreaterOrEqual(t, metrics.FactCount, minFacts, "nombre de faits insuffisant")
	}
}

// LogTestSection affiche un séparateur de section dans les logs de test.
func LogTestSection(t *testing.T, title string) {
	t.Log("")
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log(title)
	t.Log("═══════════════════════════════════════════════════════════════")
	t.Log("")
}

// LogTestSubsection affiche un sous-titre dans les logs de test.
func LogTestSubsection(t *testing.T, title string) {
	t.Log("")
	t.Log("───────────────────────────────────────────────────────────────")
	t.Log(title)
	t.Log("───────────────────────────────────────────────────────────────")
}

// FilterXuples filtre les xuples selon un prédicat.
func FilterXuples(xuplesSlice []*xuples.Xuple, predicate func(*xuples.Xuple) bool) []*xuples.Xuple {
	var result []*xuples.Xuple
	for _, x := range xuplesSlice {
		if predicate(x) {
			result = append(result, x)
		}
	}
	return result
}

// GetXupleField récupère la valeur d'un champ d'un xuple de manière sûre.
func GetXupleField(xuple *xuples.Xuple, fieldName string) (interface{}, error) {
	if xuple == nil {
		return nil, fmt.Errorf("xuple est nil")
	}
	if xuple.Fact == nil {
		return nil, fmt.Errorf("xuple.Fact est nil")
	}
	value, ok := xuple.Fact.Fields[fieldName]
	if !ok {
		return nil, fmt.Errorf("champ '%s' non trouvé dans xuple", fieldName)
	}
	return value, nil
}

// GetXupleFieldString récupère la valeur string d'un champ.
func GetXupleFieldString(t *testing.T, xuple *xuples.Xuple, fieldName string) string {
	value, err := GetXupleField(xuple, fieldName)
	require.NoError(t, err, "échec récupération champ '%s'", fieldName)
	str, ok := value.(string)
	require.True(t, ok, "champ '%s' n'est pas un string: %T", fieldName, value)
	return str
}

// GetXupleFieldFloat récupère la valeur float64 d'un champ.
func GetXupleFieldFloat(t *testing.T, xuple *xuples.Xuple, fieldName string) float64 {
	value, err := GetXupleField(xuple, fieldName)
	require.NoError(t, err, "échec récupération champ '%s'", fieldName)

	// Gérer int et float64
	switch v := value.(type) {
	case float64:
		return v
	case int:
		return float64(v)
	default:
		t.Fatalf("champ '%s' n'est pas un nombre: %T", fieldName, value)
		return 0
	}
}

// GetXupleFieldInt récupère la valeur int d'un champ.
func GetXupleFieldInt(t *testing.T, xuple *xuples.Xuple, fieldName string) int {
	value, err := GetXupleField(xuple, fieldName)
	require.NoError(t, err, "échec récupération champ '%s'", fieldName)

	// Gérer int et float64
	switch v := value.(type) {
	case int:
		return v
	case float64:
		return int(v)
	default:
		t.Fatalf("champ '%s' n'est pas un nombre entier: %T", fieldName, value)
		return 0
	}
}
