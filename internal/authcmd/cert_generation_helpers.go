// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package authcmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// certConfig contient la configuration pour la génération de certificats
type certConfig struct {
	outputDir string
	hosts     []string
	validDays int
	org       string
	format    string
}

// certGenerationResult contient le résultat de la génération de certificats
type certGenerationResult struct {
	certPath  string
	keyPath   string
	caPath    string
	notBefore time.Time
	notAfter  time.Time
}

// parseCertFlags parse les flags de ligne de commande pour la génération de certificats
func parseCertFlags(args []string, stderr io.Writer) (*certConfig, error) {
	fs := flag.NewFlagSet("generate-cert", flag.ContinueOnError)
	fs.SetOutput(stderr)

	outputDir := fs.String("output-dir", "./certs", "Répertoire de sortie pour les certificats")
	hosts := fs.String("hosts", "localhost,127.0.0.1", "Hôtes/IPs séparés par des virgules")
	validDays := fs.Int("valid-days", 365, "Durée de validité en jours")
	org := fs.String("org", "TSD Development", "Nom de l'organisation")
	format := fs.String("format", "text", "Format de sortie (text, json)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	return &certConfig{
		outputDir: *outputDir,
		hosts:     parseHostsList(*hosts),
		validDays: *validDays,
		org:       *org,
		format:    *format,
	}, nil
}

// parseHostsList parse une liste d'hôtes séparés par des virgules
func parseHostsList(hostsStr string) []string {
	hostList := strings.Split(hostsStr, ",")
	for i, h := range hostList {
		hostList[i] = strings.TrimSpace(h)
	}
	return hostList
}

// generateECDSAPrivateKey génère une clé privée ECDSA
func generateECDSAPrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// createCertificateTemplate crée un template de certificat X.509
func createCertificateTemplate(config *certConfig) (*x509.Certificate, error) {
	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(config.validDays) * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("erreur génération numéro série: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{config.org},
			CommonName:   config.hosts[0],
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	// Ajouter les hôtes au certificat
	for _, h := range config.hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	return template, nil
}

// createSelfSignedCertificate crée un certificat auto-signé
func createSelfSignedCertificate(template *x509.Certificate, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, fmt.Errorf("erreur création certificat: %w", err)
	}
	return derBytes, nil
}

// writeCertificateFiles écrit les fichiers de certificat et de clé
func writeCertificateFiles(config *certConfig, certDER []byte, privateKey *ecdsa.PrivateKey) (*certGenerationResult, error) {
	// Sauvegarder le certificat
	certPath := filepath.Join(config.outputDir, "server.crt")
	certOut, err := os.Create(certPath)
	if err != nil {
		return nil, fmt.Errorf("erreur création fichier certificat: %w", err)
	}

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		certOut.Close()
		return nil, fmt.Errorf("erreur encodage certificat: %w", err)
	}
	certOut.Close()

	// Sauvegarder la clé privée
	keyPath := filepath.Join(config.outputDir, "server.key")
	keyOut, err := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, fmt.Errorf("erreur création fichier clé: %w", err)
	}

	privBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		keyOut.Close()
		return nil, fmt.Errorf("erreur marshalling clé: %w", err)
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}); err != nil {
		keyOut.Close()
		return nil, fmt.Errorf("erreur encodage clé: %w", err)
	}
	keyOut.Close()

	// Créer aussi une copie du certificat comme CA (pour les clients)
	caPath := filepath.Join(config.outputDir, "ca.crt")
	if err := copyFile(certPath, caPath); err != nil {
		// Erreur non fatale, on continue
		fmt.Fprintf(os.Stderr, "⚠️  Avertissement: impossible de créer ca.crt: %v\n", err)
	}

	// Extraire les dates du certificat pour le résultat
	cert, _ := x509.ParseCertificate(certDER)

	return &certGenerationResult{
		certPath:  certPath,
		keyPath:   keyPath,
		caPath:    caPath,
		notBefore: cert.NotBefore,
		notAfter:  cert.NotAfter,
	}, nil
}
