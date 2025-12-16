// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package authcmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

// formatCertOutputJSON formate la sortie en JSON
func formatCertOutputJSON(result *certGenerationResult, config *certConfig, stdout io.Writer) {
	output := map[string]interface{}{
		"success":      true,
		"cert_path":    result.certPath,
		"key_path":     result.keyPath,
		"ca_path":      result.caPath,
		"hosts":        config.hosts,
		"valid_days":   config.validDays,
		"not_before":   result.notBefore.Format(time.RFC3339),
		"not_after":    result.notAfter.Format(time.RFC3339),
		"organization": config.org,
	}
	data, _ := json.MarshalIndent(output, "", "  ")
	fmt.Fprintln(stdout, string(data))
}

// formatCertOutputText formate la sortie en texte lisible
func formatCertOutputText(result *certGenerationResult, config *certConfig, stdout io.Writer) {
	fmt.Fprintln(stdout, "🔐 Certificats TLS générés avec succès!")
	fmt.Fprintln(stdout, "=====================================")
	fmt.Fprintf(stdout, "\n📁 Répertoire: %s\n\n", config.outputDir)
	fmt.Fprintf(stdout, "📄 Fichiers générés:\n")
	fmt.Fprintf(stdout, "   - %s (certificat serveur)\n", result.certPath)
	fmt.Fprintf(stdout, "   - %s (clé privée serveur)\n", result.keyPath)
	fmt.Fprintf(stdout, "   - %s (copie du certificat pour trust store client)\n\n", result.caPath)
	fmt.Fprintf(stdout, "🏷️  Hôtes autorisés: %s\n", strings.Join(config.hosts, ", "))
	fmt.Fprintf(stdout, "📅 Valide de %s à %s\n", result.notBefore.Format("2006-01-02"), result.notAfter.Format("2006-01-02"))
	fmt.Fprintf(stdout, "🏢 Organisation: %s\n\n", config.org)
	fmt.Fprintln(stdout, "⚠️  IMPORTANT:")
	fmt.Fprintf(stdout, "   - La clé privée (%s) doit rester SECRÈTE\n", result.keyPath)
	fmt.Fprintln(stdout, "   - Ne JAMAIS committer les certificats dans Git")
	fmt.Fprintln(stdout, "   - Ces certificats sont auto-signés (pour développement)")
	fmt.Fprintln(stdout, "")
	fmt.Fprintln(stdout, "📝 Utilisation:")
	fmt.Fprintf(stdout, "   Serveur: tsd server --tls-cert %s --tls-key %s\n", result.certPath, result.keyPath)
	fmt.Fprintf(stdout, "   Client:  tsd client --server https://localhost:8080 --tls-ca %s\n", result.caPath)
	fmt.Fprintln(stdout, "")
	fmt.Fprintln(stdout, "💡 Pour production, utilisez des certificats signés par une CA reconnue (Let's Encrypt, etc.)")
}
