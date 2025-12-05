# Configuration TLS/HTTPS pour TSD

**Version**: 1.0.0  
**Dernière mise à jour**: Janvier 2025

---

## Vue d'ensemble

TSD utilise **HTTPS/TLS par défaut** pour toutes les communications client-serveur. Cette documentation détaille la configuration, les options avancées et les meilleures pratiques pour sécuriser votre déploiement TSD.

## Table des Matières

1. [Architecture TLS](#architecture-tls)
2. [Génération de Certificats](#génération-de-certificats)
3. [Configuration Serveur](#configuration-serveur)
4. [Configuration Client](#configuration-client)
5. [Variables d'Environnement](#variables-denvironnement)
6. [Scénarios de Déploiement](#scénarios-de-déploiement)
7. [Sécurité Avancée](#sécurité-avancée)
8. [Résolution de Problèmes](#résolution-de-problèmes)
9. [Référence API](#référence-api)

---

## Architecture TLS

### Flux de Communication

```
┌─────────────┐                    ┌─────────────┐
│             │   HTTPS/TLS 1.2+   │             │
│  TSD Client ├───────────────────>│ TSD Server  │
│             │   Certificat vérifié│             │
└─────────────┘                    └─────────────┘
      │                                    │
      │ CA Certificate                     │ Server Certificate
      │ (ca.crt)                          │ (server.crt)
      │                                    │ Private Key
      │                                    │ (server.key)
      v                                    v
   Verification                      Authentication
```

### Composants TLS

1. **Certificat Serveur** (`server.crt`)
   - Identifie le serveur
   - Contient la clé publique
   - Signé par une CA (auto-signé en dev)

2. **Clé Privée Serveur** (`server.key`)
   - Doit rester secrète
   - Utilisée pour déchiffrer les communications
   - Permissions 600 obligatoires

3. **Certificat CA** (`ca.crt`)
   - Utilisé par le client pour vérifier le serveur
   - Certificat racine de confiance
   - Peut être auto-signé (dev) ou d'une CA reconnue (prod)

### Versions et Cipher Suites

**Versions TLS Supportées**:
- TLS 1.2 (minimum)
- TLS 1.3 (recommandé si disponible)

**Cipher Suites** (ordre de préférence):
1. `TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384`
2. `TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256`
3. `TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384`
4. `TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256`

**Caractéristiques**:
- ✅ Perfect Forward Secrecy (ECDHE)
- ✅ Authenticated Encryption (GCM)
- ✅ Pas de ciphers obsolètes (RC4, DES, MD5)
- ✅ Préférence serveur activée

---

## Génération de Certificats

### Commande de Base

```bash
tsd auth generate-cert
```

**Génère**:
- `./certs/server.crt` - Certificat serveur
- `./certs/server.key` - Clé privée (permissions 600)
- `./certs/ca.crt` - Certificat CA pour clients

### Options Avancées

```bash
tsd auth generate-cert [options]

Options:
  -output-dir string
        Répertoire de sortie (défaut: "./certs")
  
  -hosts string
        Liste d'hôtes/IPs séparés par virgules
        (défaut: "localhost,127.0.0.1")
  
  -valid-days int
        Durée de validité en jours (défaut: 365)
  
  -org string
        Nom de l'organisation (défaut: "TSD Development")
  
  -format string
        Format de sortie: text ou json (défaut: "text")
```

### Exemples

#### Certificat pour Développement Local

```bash
tsd auth generate-cert
```

#### Certificat avec Plusieurs Hôtes

```bash
tsd auth generate-cert \
  -hosts "localhost,127.0.0.1,192.168.1.100,myserver.local"
```

#### Certificat Longue Durée

```bash
tsd auth generate-cert \
  -valid-days 730 \
  -org "My Company Production"
```

#### Sortie JSON (Automation)

```bash
tsd auth generate-cert -format json > cert-info.json
```

**Sortie**:
```json
{
  "success": true,
  "cert_path": "./certs/server.crt",
  "key_path": "./certs/server.key",
  "ca_path": "./certs/ca.crt",
  "hosts": ["localhost", "127.0.0.1"],
  "valid_days": 365,
  "not_before": "2025-01-15T10:00:00Z",
  "not_after": "2026-01-15T10:00:00Z",
  "organization": "TSD Development"
}
```

### Vérification des Certificats

```bash
# Afficher les détails
openssl x509 -in certs/server.crt -text -noout

# Vérifier la validité
openssl x509 -in certs/server.crt -noout -checkend 0

# Vérifier les hôtes autorisés (SAN)
openssl x509 -in certs/server.crt -noout -text | \
  grep -A1 "Subject Alternative Name"

# Vérifier la correspondance clé/certificat
openssl x509 -noout -modulus -in certs/server.crt | openssl md5
openssl rsa -noout -modulus -in certs/server.key | openssl md5
```

---

## Configuration Serveur

### Démarrage de Base

```bash
# Détection automatique des certificats dans ./certs/
tsd server
```

Le serveur cherche automatiquement:
1. `./certs/server.crt`
2. `./certs/server.key`

Si trouvés → démarrage en HTTPS  
Si absents → erreur avec suggestions

### Certificats Personnalisés

```bash
tsd server \
  --tls-cert /path/to/certificate.crt \
  --tls-key /path/to/private-key.key
```

### Options Complètes

```bash
tsd server [options]

Options TLS:
  --tls-cert string
        Chemin vers le certificat TLS
        (défaut: "./certs/server.crt")
  
  --tls-key string
        Chemin vers la clé privée TLS
        (défaut: "./certs/server.key")
  
  --insecure
        Désactiver TLS (mode HTTP non sécurisé)
        ⚠️ Développement uniquement!

Options Serveur:
  --host string      Hôte (défaut: "0.0.0.0")
  --port int         Port (défaut: 8080)
  -v                 Mode verbeux
  --auth string      Type auth: none, key, jwt
  --auth-keys        Clés API (si --auth key)
  --jwt-secret       Secret JWT (si --auth jwt)
```

### Exemples de Configuration

#### Développement Local

```bash
# Avec certificats auto-signés
tsd auth generate-cert
tsd server
```

#### Production avec Let's Encrypt

```bash
tsd server \
  --tls-cert /etc/letsencrypt/live/tsd.example.com/fullchain.pem \
  --tls-key /etc/letsencrypt/live/tsd.example.com/privkey.pem \
  --port 443 \
  --auth jwt \
  --jwt-secret "$(cat /var/secrets/jwt-secret)"
```

#### Multi-environnement

```bash
# Production
export ENVIRONMENT=production
export TSD_TLS_CERT=/etc/tsd/prod-cert.crt
export TSD_TLS_KEY=/etc/tsd/prod-key.key
tsd server --port 443

# Staging
export ENVIRONMENT=staging
export TSD_TLS_CERT=/etc/tsd/staging-cert.crt
export TSD_TLS_KEY=/etc/tsd/staging-key.key
tsd server --port 8443
```

### Mode HTTP Non Sécurisé

⚠️ **ATTENTION**: À utiliser UNIQUEMENT en développement local!

```bash
tsd server --insecure
```

**Logs**:
```
[TSD-SERVER] 🚀 Démarrage du serveur TSD sur http://0.0.0.0:8080
[TSD-SERVER] ⚠️  TLS: désactivé (mode HTTP non sécurisé)
[TSD-SERVER] ⚠️  AVERTISSEMENT: Ne pas utiliser en production!
```

---

## Configuration Client

### URL par Défaut

L'URL par défaut du client est **HTTPS**:
```bash
# Équivalent à --server https://localhost:8080
tsd client program.tsd
```

### Options TLS

```bash
tsd client [file] [options]

Options TLS:
  --server string
        URL du serveur (défaut: "https://localhost:8080")
  
  --tls-ca string
        Certificat CA pour vérifier le serveur
        (défaut: "./certs/ca.crt")
  
  --insecure
        Désactiver la vérification TLS
        ⚠️ Développement uniquement!

Options Client:
  -v                 Mode verbeux
  --token string     Token d'authentification
  --format string    Format sortie: text ou json
  --timeout duration Timeout requêtes (défaut: 30s)
```

### Scénarios d'Utilisation

#### Développement avec Certificats Auto-signés

**Option 1**: Mode insecure (rapide)
```bash
tsd client program.tsd -insecure
```

**Option 2**: Avec CA (recommandé)
```bash
tsd client program.tsd --tls-ca ./certs/ca.crt
```

**Option 3**: Variable d'environnement
```bash
export TSD_TLS_CA=./certs/ca.crt
tsd client program.tsd
```

#### Production avec Certificat Valide

```bash
# Pas besoin de --insecure ou --tls-ca
tsd client program.tsd --server https://tsd.example.com
```

#### Serveur Distant avec CA d'Entreprise

```bash
tsd client program.tsd \
  --server https://tsd.internal.company.com:8080 \
  --tls-ca /etc/ssl/certs/company-ca.crt
```

#### Health Check

```bash
# Avec vérification TLS
tsd client --health --tls-ca ./certs/ca.crt

# Mode insecure
tsd client --health --insecure
```

### Mode HTTP Non Sécurisé

Pour se connecter à un serveur HTTP:

```bash
tsd client program.tsd --server http://localhost:8080
```

---

## Variables d'Environnement

### Serveur

```bash
# Certificat TLS
export TSD_TLS_CERT=/path/to/cert.crt

# Clé privée TLS
export TSD_TLS_KEY=/path/to/key.key

# Mode insecure (true pour désactiver TLS)
export TSD_INSECURE=true
```

### Client

```bash
# Certificat CA pour vérification
export TSD_TLS_CA=/path/to/ca.crt

# Mode insecure (true pour désactiver vérification)
export TSD_CLIENT_INSECURE=true

# Token d'authentification
export TSD_AUTH_TOKEN=your-token-here
```

### Exemple Complet

```bash
#!/bin/bash
# Configuration environnement TSD

# === Serveur ===
export TSD_TLS_CERT=/etc/tsd/certs/server.crt
export TSD_TLS_KEY=/etc/tsd/certs/server.key
export TSD_JWT_SECRET=$(cat /var/secrets/tsd-jwt)

# === Client ===
export TSD_TLS_CA=/etc/tsd/certs/ca.crt
export TSD_AUTH_TOKEN=$(cat ~/.tsd/token)

# Lancer serveur
tsd server --port 8080 --auth jwt &

# Attendre démarrage
sleep 2

# Utiliser client
tsd client program.tsd -v
```

---

## Scénarios de Déploiement

### 1. Développement Local

**Objectif**: Test rapide avec certificats auto-signés

```bash
# Setup (une fois)
tsd auth generate-cert

# Serveur (terminal 1)
tsd server

# Client (terminal 2)
tsd client test.tsd -insecure
```

### 2. Staging avec Certificats Internes

**Objectif**: Environnement de test avec certificats d'entreprise

```bash
# Serveur
tsd server \
  --tls-cert /etc/tsd/staging/cert.crt \
  --tls-key /etc/tsd/staging/key.key \
  --port 8080 \
  --auth jwt \
  --jwt-secret "$STAGING_JWT_SECRET"

# Client
tsd client program.tsd \
  --server https://staging-tsd.company.com:8080 \
  --tls-ca /etc/ssl/certs/company-ca.crt \
  --token "$STAGING_TOKEN"
```

### 3. Production avec Let's Encrypt

**Objectif**: Déploiement public avec certificats valides

#### Obtention du Certificat

```bash
# Installer certbot
sudo apt-get install certbot

# Obtenir certificat
sudo certbot certonly --standalone -d tsd.example.com

# Certificats générés dans:
# /etc/letsencrypt/live/tsd.example.com/
#   - fullchain.pem (certificat)
#   - privkey.pem (clé privée)
```

#### Configuration Serveur

```bash
# Script de démarrage
#!/bin/bash
set -e

# Variables
DOMAIN="tsd.example.com"
CERT_DIR="/etc/letsencrypt/live/$DOMAIN"

# Démarrer serveur
tsd server \
  --host 0.0.0.0 \
  --port 443 \
  --tls-cert "$CERT_DIR/fullchain.pem" \
  --tls-key "$CERT_DIR/privkey.pem" \
  --auth jwt \
  --jwt-secret "$(cat /var/secrets/tsd-jwt-secret)" \
  -v
```

#### Renouvellement Automatique

```bash
# Cron job pour renouvellement (tous les jours à 2h)
0 2 * * * /usr/bin/certbot renew --quiet --post-hook "systemctl restart tsd"
```

#### Client

```bash
# Pas besoin de --insecure (certificat valide)
tsd client program.tsd --server https://tsd.example.com
```

### 4. Docker avec Volumes

**Dockerfile**:
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /build
COPY . .
RUN go build -o tsd ./cmd/tsd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/tsd /usr/local/bin/
VOLUME ["/certs"]
EXPOSE 8080
ENTRYPOINT ["tsd", "server"]
CMD ["--tls-cert", "/certs/server.crt", "--tls-key", "/certs/server.key"]
```

**Docker Compose**:
```yaml
version: '3.8'

services:
  tsd-server:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./certs:/certs:ro
      - ./secrets:/secrets:ro
    environment:
      - TSD_JWT_SECRET_FILE=/secrets/jwt-secret
    command:
      - --tls-cert=/certs/server.crt
      - --tls-key=/certs/server.key
      - --auth=jwt
      - --jwt-secret=$(cat /secrets/jwt-secret)
```

### 5. Kubernetes avec Secrets

**Secret**:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: tsd-tls
type: kubernetes.io/tls
data:
  tls.crt: <base64-encoded-cert>
  tls.key: <base64-encoded-key>
```

**Deployment**:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: tsd-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: tsd-server
  template:
    metadata:
      labels:
        app: tsd-server
    spec:
      containers:
      - name: tsd
        image: tsd:latest
        ports:
        - containerPort: 8080
          name: https
        volumeMounts:
        - name: tls-certs
          mountPath: /certs
          readOnly: true
        command:
        - tsd
        - server
        - --tls-cert=/certs/tls.crt
        - --tls-key=/certs/tls.key
        - --auth=jwt
        env:
        - name: TSD_JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: tsd-jwt
              key: secret
      volumes:
      - name: tls-certs
        secret:
          secretName: tsd-tls
```

---

## Sécurité Avancée

### Durcissement TLS

#### Configuration Serveur

Pour un niveau de sécurité maximal, le code TSD est déjà configuré avec:

```go
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,
    CipherSuites: []uint16{
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
    },
    PreferServerCipherSuites: true,
}
```

#### Tests de Sécurité

```bash
# Test avec SSL Labs (si public)
# https://www.ssllabs.com/ssltest/

# Test avec testssl.sh
git clone --depth 1 https://github.com/drwetter/testssl.sh.git
cd testssl.sh
./testssl.sh https://localhost:8080

# Test avec nmap
nmap --script ssl-enum-ciphers -p 8080 localhost
```

### Rotation des Certificats

#### Script de Rotation

```bash
#!/bin/bash
# rotate-certs.sh

set -e

CERTS_DIR="/etc/tsd/certs"
BACKUP_DIR="/etc/tsd/certs-backup-$(date +%Y%m%d)"

echo "🔄 Rotation des certificats TSD"

# 1. Backup des anciens certificats
echo "📦 Backup des certificats actuels..."
mkdir -p "$BACKUP_DIR"
cp -r "$CERTS_DIR"/* "$BACKUP_DIR/"

# 2. Générer nouveaux certificats
echo "🔑 Génération nouveaux certificats..."
tsd auth generate-cert \
  -output-dir "$CERTS_DIR-new" \
  -hosts "$(hostname),$(hostname -I | tr ' ' ',')" \
  -valid-days 365

# 3. Test avec les nouveaux certificats
echo "🧪 Test nouveaux certificats..."
tsd server \
  --tls-cert "$CERTS_DIR-new/server.crt" \
  --tls-key "$CERTS_DIR-new/server.key" \
  --port 9999 &
TEST_PID=$!
sleep 2

if curl -k https://localhost:9999/health; then
    echo "✅ Nouveaux certificats OK"
    kill $TEST_PID
else
    echo "❌ Test échoué, restauration backup"
    kill $TEST_PID
    exit 1
fi

# 4. Remplacer les certificats
echo "🔄 Remplacement certificats..."
rm -rf "$CERTS_DIR-old"
mv "$CERTS_DIR" "$CERTS_DIR-old"
mv "$CERTS_DIR-new" "$CERTS_DIR"

# 5. Redémarrer le serveur
echo "🔄 Redémarrage serveur..."
systemctl restart tsd

echo "✅ Rotation terminée avec succès"
```

#### Monitoring d'Expiration

```bash
#!/bin/bash
# check-cert-expiry.sh

CERT_FILE="/etc/tsd/certs/server.crt"
ALERT_DAYS=30

# Obtenir date d'expiration
EXPIRY_DATE=$(openssl x509 -in "$CERT_FILE" -noout -enddate | cut -d= -f2)
EXPIRY_EPOCH=$(date -d "$EXPIRY_DATE" +%s)
NOW_EPOCH=$(date +%s)
DAYS_LEFT=$(( ($EXPIRY_EPOCH - $NOW_EPOCH) / 86400 ))

echo "Certificat expire dans $DAYS_LEFT jours"

if [ $DAYS_LEFT -le $ALERT_DAYS ]; then
    echo "⚠️ ALERTE: Certificat expire bientôt!"
    # Envoyer alerte (email, Slack, etc.)
fi
```

### Gestion des Clés

#### Permissions Strictes

```bash
# Certificats (lecture publique OK)
chmod 644 /etc/tsd/certs/server.crt
chmod 644 /etc/tsd/certs/ca.crt

# Clés privées (lecture propriétaire uniquement)
chmod 600 /etc/tsd/certs/server.key
chown tsd:tsd /etc/tsd/certs/server.key

# Répertoire
chmod 750 /etc/tsd/certs
chown tsd:tsd /etc/tsd/certs
```

#### Stockage Sécurisé (HashiCorp Vault)

```bash
# Stocker dans Vault
vault kv put secret/tsd/certs \
  server_crt=@/etc/tsd/certs/server.crt \
  server_key=@/etc/tsd/certs/server.key

# Récupérer au démarrage
#!/bin/bash
export VAULT_ADDR='https://vault.example.com'
export VAULT_TOKEN='your-token'

# Télécharger certificats
vault kv get -field=server_crt secret/tsd/certs > /tmp/server.crt
vault kv get -field=server_key secret/tsd/certs > /tmp/server.key
chmod 600 /tmp/server.key

# Démarrer serveur
tsd server \
  --tls-cert /tmp/server.crt \
  --tls-key /tmp/server.key

# Nettoyer à l'arrêt
trap "rm -f /tmp/server.crt /tmp/server.key" EXIT
```

---

## Résolution de Problèmes

### Erreurs Courantes

#### 1. Certificat Non Trouvé

**Erreur**:
```
❌ Certificat TLS non trouvé: ./certs/server.crt

💡 Solutions:
   1. Générer des certificats: tsd auth generate-cert
   2. Spécifier un certificat: --tls-cert /path/to/cert.crt
   3. Démarrer en mode non sécurisé: --insecure (déconseillé en production)
```

**Solution**:
```bash
tsd auth generate-cert
tsd server
```

#### 2. Certificate Signed by Unknown Authority

**Erreur**:
```
❌ Erreur communication serveur: Get "https://localhost:8080/api/v1/execute": 
    x509: certificate signed by unknown authority
```

**Solutions**:

A. Mode insecure (dev):
```bash
tsd client test.tsd -insecure
```

B. Spécifier CA:
```bash
tsd client test.tsd --tls-ca ./certs/ca.crt
```

C. Variable d'environnement:
```bash
export TSD_TLS_CA=./certs/ca.crt
tsd client test.tsd
```

#### 3. Certificate Invalid for Hostname

**Erreur**:
```
❌ x509: certificate is valid for localhost, not 192.168.1.100
```

**Solution**: Régénérer avec le bon hostname
```bash
tsd auth generate-cert -hosts "localhost,127.0.0.1,192.168.1.100"
```

#### 4. Permission Denied sur Clé Privée

**Erreur**:
```
❌ Erreur démarrage serveur: open ./certs/server.key: permission denied
```

**Solution**:
```bash
chmod 600 ./certs/server.key
chown $USER ./certs/server.key
```

### Diagnostic

#### Vérifier Configuration TLS

```bash
# Test connexion basique
curl -k https://localhost:8080/health

# Test avec CA
curl --cacert ./certs/ca.crt https://localhost:8080/health

# Test avec OpenSSL
openssl s_client -connect localhost:8080 -servername localhost

# Vérifier version TLS et ciphers
openssl s_client -connect localhost:8080 -tls1_2
```

#### Logs Détaillés

```bash
# Serveur en mode verbeux
tsd server -v

# Client en mode verbeux
tsd client test.tsd -v -insecure
```

#### Vérifier Certificats

```bash
# Détails certificat
openssl x509 -in certs/server.crt -text -noout

# Vérifier validité
openssl x509 -in certs/server.crt -noout -checkend 0 && \
  echo "Certificat valide" || echo "Certificat expiré"

# Vérifier correspondance clé/cert
diff \
  <(openssl x509 -noout -modulus -in certs/server.crt | openssl md5) \
  <(openssl rsa -noout -modulus -in certs/server.key | openssl md5)
```

---

## Référence API

### Configuration TLS Programmatique

Si vous intégrez TSD dans votre application Go:

```go
package main

import (
    "crypto/tls"
    "crypto/x509"
    "os"
    
    "github.com/treivax/tsd/internal/servercmd"
)

func main() {
    // Configuration serveur
    config := &servercmd.Config{
        Host: "0.0.0.0",
        Port: 8080,
        TLSCertFile: "/path/to/cert.crt",
        TLSKeyFile: "/path/to/key.key",
        Insecure: false,
    }
    
    // Créer et démarrer serveur
    // (voir servercmd.go pour l'implémentation complète)
}
```

### Structure Configuration TLS

```go
type Config struct {
    // ... autres champs ...
    
    // TLS
    TLSCertFile string
    TLSKeyFile  string
    Insecure    bool
}

// Configuration TLS serveur
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,
    CipherSuites: []uint16{
        tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
        tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
    },
    PreferServerCipherSuites: true,
}

// Configuration TLS client
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,
    InsecureSkipVerify: config.Insecure, // si --insecure
    RootCAs: caCertPool, // si --tls-ca fourni
}
```

---

## Checklist de Production

### Avant le Déploiement

- [ ] Certificats obtenus d'une CA reconnue (Let's Encrypt, etc.)
- [ ] Clés privées stockées de manière sécurisée (Vault, AWS Secrets Manager)
- [ ] Permissions correctes (600 pour clés, 644 pour certificats)
- [ ] Renouvellement automatique configuré
- [ ] Monitoring expiration en place
- [ ] Tests de sécurité effectués (testssl.sh, SSL Labs)
- [ ] Pas de flag `--insecure` dans configuration
- [ ] Variables d'environnement sécurisées
- [ ] Logs configurés et monitored
- [ ] Documentation d'incident prête

### Monitoring Continue

- [ ] Métriques TLS (versions, ciphers utilisés)
- [ ] Alertes expiration certificats (30 jours avant)
- [ ] Logs d'erreurs TLS monitored
- [ ] Tests automatisés connexion HTTPS
- [ ] Audits de sécurité réguliers

---

## Ressources

### Documentation

- [Tutoriel Authentification](./AUTHENTICATION_TUTORIAL.md)
- [Guide TLS Exemples](../examples/tls/README.md)
- [README Principal](../README.md)

### Outils Externes

- [Let's Encrypt](https://letsencrypt.org/)
- [testssl.sh](https://testssl.sh/)
- [SSL Labs](https://www.ssllabs.com/ssltest/)
- [Mozilla SSL Configuration Generator](https://ssl-config.mozilla.org/)

### Standards

- [RFC 5246 - TLS 1.2](https://datatracker.ietf.org/doc/html/rfc5246)
- [RFC 8446 - TLS 1.3](https://datatracker.ietf.org/doc/html/rfc8446)
- [RFC 5280 - X.509 Certificates](https://datatracker.ietf.org/doc/html/rfc5280)

---

**Document maintenu par l'équipe TSD**  
**Dernière révision**: Janvier 2025