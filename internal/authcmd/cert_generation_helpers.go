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

// Constantes pour la génération de certificats
const (
	// DefaultCertOutputDir répertoire par défaut pour les certificats
	DefaultCertOutputDir = "./certs"

	// DefaultCertHosts hôtes par défaut pour les certificats
	DefaultCertHosts = "localhost,127.0.0.1"

	// DefaultCertValidDays durée de validité par défaut en jours
	DefaultCertValidDays = 365

	// DefaultCertOrganization organisation par défaut
	DefaultCertOrganization = "TSD Development"

	// DefaultOutputFormat format de sortie par défaut
	DefaultOutputFormat = "text"

	// SerialNumberBitLength longueur en bits du numéro de série
	SerialNumberBitLength = 128

	// CertFileName nom du fichier certificat
	CertFileName = "server.crt"

	// KeyFileName nom du fichier clé privée
	KeyFileName = "server.key"

	// CAFileName nom du fichier CA (copie du certificat auto-signé)
	CAFileName = "ca.crt"

	// KeyFilePermissions permissions du fichier de clé privée
	KeyFilePermissions = 0600

	// CertFilePermissions permissions du fichier certificat
	CertFilePermissions = 0644
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

	outputDir := fs.String("output-dir", DefaultCertOutputDir, "Répertoire de sortie pour les certificats")
	hosts := fs.String("hosts", DefaultCertHosts, "Hôtes/IPs séparés par des virgules")
	validDays := fs.Int("valid-days", DefaultCertValidDays, "Durée de validité en jours")
	org := fs.String("org", DefaultCertOrganization, "Nom de l'organisation")
	format := fs.String("format", DefaultOutputFormat, "Format de sortie (text, json)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	hostList := parseHostsList(*hosts)
	if len(hostList) == 0 || (len(hostList) == 1 && hostList[0] == "") {
		return nil, fmt.Errorf("au moins un hôte doit être spécifié")
	}

	return &certConfig{
		outputDir: *outputDir,
		hosts:     hostList,
		validDays: *validDays,
		org:       *org,
		format:    *format,
	}, nil
}

// parseHostsList parse une liste d'hôtes séparés par des virgules
func parseHostsList(hostsStr string) []string {
	if hostsStr == "" {
		return []string{}
	}

	hostList := strings.Split(hostsStr, ",")
	result := make([]string, 0, len(hostList))

	for _, h := range hostList {
		trimmed := strings.TrimSpace(h)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// generateECDSAPrivateKey génère une clé privée ECDSA
func generateECDSAPrivateKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
}

// createCertificateTemplate crée un template de certificat X.509
func createCertificateTemplate(config *certConfig) (*x509.Certificate, error) {
	if len(config.hosts) == 0 {
		return nil, fmt.Errorf("au moins un hôte doit être spécifié")
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(config.validDays) * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), SerialNumberBitLength)
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
		IsCA:                  false, // Certificat serveur/client, pas CA (sécurité: CWE-295)
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
	certPath := filepath.Join(config.outputDir, CertFileName)
	if err := writeCertificatePEM(certPath, certDER); err != nil {
		return nil, err
	}

	keyPath := filepath.Join(config.outputDir, KeyFileName)
	if err := writePrivateKeyPEM(keyPath, privateKey); err != nil {
		return nil, err
	}

	caPath := filepath.Join(config.outputDir, CAFileName)
	if err := createCACopy(certPath, caPath); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Avertissement: impossible de créer ca.crt: %v\n", err)
	}

	cert, _ := x509.ParseCertificate(certDER)

	return &certGenerationResult{
		certPath:  certPath,
		keyPath:   keyPath,
		caPath:    caPath,
		notBefore: cert.NotBefore,
		notAfter:  cert.NotAfter,
	}, nil
}

// writeCertificatePEM écrit un certificat au format PEM
func writeCertificatePEM(path string, certDER []byte) error {
	certOut, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("erreur création fichier certificat: %w", err)
	}
	defer certOut.Close()

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER}); err != nil {
		return fmt.Errorf("erreur encodage certificat: %w", err)
	}

	return nil
}

// writePrivateKeyPEM écrit une clé privée ECDSA au format PEM
func writePrivateKeyPEM(path string, privateKey *ecdsa.PrivateKey) error {
	keyOut, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, KeyFilePermissions)
	if err != nil {
		return fmt.Errorf("erreur création fichier clé: %w", err)
	}
	defer keyOut.Close()

	privBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("erreur marshalling clé: %w", err)
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes}); err != nil {
		return fmt.Errorf("erreur encodage clé: %w", err)
	}

	return nil
}

// createCACopy crée une copie du certificat auto-signé pour la confiance des clients
// Note: Ce n'est PAS un certificat CA - c'est simplement le certificat serveur auto-signé
// que les clients doivent ajouter à leur trust store pour accepter la connexion TLS
func createCACopy(certPath, caPath string) error {
	return copyFile(certPath, caPath)
}
