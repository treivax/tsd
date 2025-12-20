// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/treivax/tsd/tsdio"
)

// testCompleteRoundtrip teste le scénario complet HTTPS + JWT + Exécution TSD
func testCompleteRoundtrip(ctx *testContext) {
	ctx.t.Log("\n🔄 TEST: Roundtrip complet HTTPS + JWT")
	ctx.t.Log("=======================================")

	// Créer requête d'exécution
	requestBody := tsdio.NewExecuteRequest(SimpleTSDProgram)
	requestBody.Verbose = true

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		ctx.t.Fatalf("❌ Erreur marshalling requête: %v", err)
	}

	// Créer client HTTP avec TLS
	client := ctx.createHTTPClient()

	// Créer requête HTTP
	executeURL := ctx.serverURL + "/api/v1/execute"
	req, err := http.NewRequest("POST", executeURL, bytes.NewReader(jsonData))
	if err != nil {
		ctx.t.Fatalf("❌ Erreur création requête: %v", err)
	}

	// Headers requis
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ctx.jwtToken)

	ctx.t.Log("📤 Envoi requête au serveur...")
	ctx.t.Logf("   URL: %s", executeURL)
	ctx.t.Logf("   JWT: %s...", ctx.jwtToken[:32])

	// Exécuter requête
	resp, err := client.Do(req)
	if err != nil {
		ctx.t.Fatalf("❌ Erreur requête HTTP: %v", err)
	}
	defer resp.Body.Close()

	// Vérifier status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		ctx.t.Fatalf("❌ Status code = %d, attendu 200\nBody: %s", resp.StatusCode, string(body))
	}
	ctx.t.Logf("✅ Status code: %d", resp.StatusCode)

	// Vérifier Content-Type
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		ctx.t.Errorf("❌ Content-Type = %s, attendu application/json", contentType)
	}
	ctx.t.Logf("✅ Content-Type: %s", contentType)

	// Parser réponse
	var response tsdio.ExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		ctx.t.Fatalf("❌ Erreur parsing réponse: %v", err)
	}

	// Valider succès
	if !response.Success {
		ctx.t.Fatalf("❌ Exécution échouée: %s (type: %s)", response.Error, response.ErrorType)
	}
	ctx.t.Log("✅ Exécution réussie")

	// Valider résultats
	if response.Results == nil {
		ctx.t.Fatal("❌ Résultats manquants")
	}

	ctx.t.Logf("📊 Résultats:")
	ctx.t.Logf("   - Facts: %d", response.Results.FactsCount)
	ctx.t.Logf("   - Activations: %d", response.Results.ActivationsCount)
	ctx.t.Logf("   - Temps: %dms", response.ExecutionTimeMs)

	// Vérifications détaillées
	if response.Results.FactsCount != ExpectedFactsCount {
		ctx.t.Errorf("❌ Facts = %d, attendu %d", response.Results.FactsCount, ExpectedFactsCount)
	}

	if response.Results.ActivationsCount != ExpectedActivationsCount {
		ctx.t.Errorf("❌ Activations = %d, attendu %d", response.Results.ActivationsCount, ExpectedActivationsCount)
	}

	if len(response.Results.Activations) != ExpectedActivationsCount {
		ctx.t.Fatalf("❌ Activations détaillées = %d, attendu %d",
			len(response.Results.Activations), ExpectedActivationsCount)
	}

	// Vérifier action déclenchée
	activation := response.Results.Activations[0]
	expectedActionName := "Xuple"
	if activation.ActionName != expectedActionName {
		ctx.t.Errorf("❌ Action = '%s', attendu '%s'", activation.ActionName, expectedActionName)
	} else {
		ctx.t.Logf("✅ Activation correcte: %s avec %d arguments", activation.ActionName, len(activation.Arguments))

		// Vérifier que c'est bien le xuple-space "adults"
		if len(activation.Arguments) > 0 {
			if xuplespace, ok := activation.Arguments[0].Value.(string); ok && xuplespace == ExpectedXupleSpaceName {
				ctx.t.Logf("✅ Xuple créé dans l'espace '%s'", xuplespace)
			}
		}
	}

	ctx.t.Log("✅ TEST ROUNDTRIP COMPLET RÉUSSI")
}

// testSimpleHTTPRoundtrip teste le scénario simple HTTP sans auth
func testSimpleHTTPRoundtrip(ctx *testContext) {
	ctx.t.Log("\n🔄 TEST: Roundtrip simple HTTP")
	ctx.t.Log("===============================")

	// Créer requête d'exécution
	requestBody := tsdio.NewExecuteRequest(SimpleTSDProgram)

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		ctx.t.Fatalf("❌ Erreur marshalling requête: %v", err)
	}

	// Créer client HTTP
	client := &http.Client{Timeout: RequestTimeout}

	// Créer requête HTTP
	executeURL := ctx.serverURL + "/api/v1/execute"
	req, err := http.NewRequest("POST", executeURL, bytes.NewReader(jsonData))
	if err != nil {
		ctx.t.Fatalf("❌ Erreur création requête: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	ctx.t.Log("📤 Envoi requête au serveur...")

	// Exécuter requête
	resp, err := client.Do(req)
	if err != nil {
		ctx.t.Fatalf("❌ Erreur requête HTTP: %v", err)
	}
	defer resp.Body.Close()

	// Vérifier status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		ctx.t.Fatalf("❌ Status code = %d, attendu 200\nBody: %s", resp.StatusCode, string(body))
	}

	// Parser réponse
	var response tsdio.ExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		ctx.t.Fatalf("❌ Erreur parsing réponse: %v", err)
	}

	// Valider succès
	if !response.Success {
		ctx.t.Fatalf("❌ Exécution échouée: %s", response.Error)
	}

	// Valider résultats basiques
	if response.Results.FactsCount != ExpectedFactsCount {
		ctx.t.Errorf("❌ Facts = %d, attendu %d", response.Results.FactsCount, ExpectedFactsCount)
	}

	ctx.t.Log("✅ TEST HTTP SIMPLE RÉUSSI")
}

// testInvalidAuthentication teste le rejet d'authentifications invalides
func testInvalidAuthentication(ctx *testContext) {
	ctx.t.Log("\n🔒 TEST: Authentification invalide")
	ctx.t.Log("===================================")

	requestBody := tsdio.NewExecuteRequest(SimpleTSDProgram)
	jsonData, _ := json.Marshal(requestBody)

	client := ctx.createHTTPClient()
	executeURL := ctx.serverURL + "/api/v1/execute"

	tests := []struct {
		name         string
		token        string
		expectedCode int
	}{
		{"Token_Invalide", "invalid-token-xyz", http.StatusUnauthorized},
		{"Token_Malformé", "not.a.jwt.token", http.StatusUnauthorized},
		{"Token_Vide", "", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		tt := tt // Capture loop variable
		ctx.t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", executeURL, bytes.NewReader(jsonData))
			req.Header.Set("Content-Type", "application/json")

			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("❌ Erreur requête: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedCode {
				t.Errorf("❌ Status = %d, attendu %d", resp.StatusCode, tt.expectedCode)
			} else {
				t.Logf("✅ Requête rejetée correctement: %d", resp.StatusCode)
			}
		})
	}

	ctx.t.Log("✅ TEST AUTHENTIFICATION INVALIDE RÉUSSI")
}

// testUnauthorizedRequest teste les requêtes sans token
func testUnauthorizedRequest(ctx *testContext) {
	ctx.t.Log("\n🚫 TEST: Requête sans token")
	ctx.t.Log("============================")

	requestBody := tsdio.NewExecuteRequest(SimpleTSDProgram)
	jsonData, _ := json.Marshal(requestBody)

	client := ctx.createHTTPClient()
	executeURL := ctx.serverURL + "/api/v1/execute"

	// Requête sans header Authorization
	req, _ := http.NewRequest("POST", executeURL, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		ctx.t.Fatalf("❌ Erreur requête: %v", err)
	}
	defer resp.Body.Close()

	// Devrait retourner 401 Unauthorized
	if resp.StatusCode != http.StatusUnauthorized {
		ctx.t.Errorf("❌ Status = %d, attendu 401", resp.StatusCode)
	} else {
		ctx.t.Logf("✅ Requête non autorisée rejetée: %d", resp.StatusCode)
	}

	ctx.t.Log("✅ TEST REQUÊTE NON AUTORISÉE RÉUSSI")
}

// testInvalidProgram teste la gestion d'un programme TSD invalide
func testInvalidProgram(ctx *testContext) {
	ctx.t.Log("\n❌ TEST: Programme TSD invalide")
	ctx.t.Log("================================")

	requestBody := tsdio.NewExecuteRequest(InvalidTSDProgram)
	jsonData, _ := json.Marshal(requestBody)

	client := ctx.createHTTPClient()
	executeURL := ctx.serverURL + "/api/v1/execute"

	req, _ := http.NewRequest("POST", executeURL, bytes.NewReader(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ctx.jwtToken)

	resp, err := client.Do(req)
	if err != nil {
		ctx.t.Fatalf("❌ Erreur requête: %v", err)
	}
	defer resp.Body.Close()

	// Parser réponse
	var errResp tsdio.ExecuteResponse
	if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
		ctx.t.Fatalf("❌ Erreur parsing réponse: %v", err)
	}

	// Vérifier que l'exécution a échoué
	if errResp.Success {
		ctx.t.Error("❌ L'exécution devrait échouer pour un programme invalide")
	} else {
		ctx.t.Logf("✅ Programme invalide rejeté correctement")
	}

	// Vérifier qu'il y a un message d'erreur
	if errResp.Error == "" {
		ctx.t.Error("❌ Message d'erreur manquant")
	} else {
		ctx.t.Logf("✅ Message d'erreur: %s", errResp.Error)
	}

	ctx.t.Log("✅ TEST PROGRAMME INVALIDE RÉUSSI")
}
