# 🔐 Exemple TLS/HTTPS pour TSD

Ce répertoire contient un guide de démarrage rapide pour configurer et utiliser TLS/HTTPS avec TSD.

## 📋 Vue d'ensemble

TSD utilise **HTTPS par défaut** pour sécuriser toutes les communications entre le client et le serveur. Ce guide vous montre comment :

1. Générer des certificats TLS pour le développement
2. Démarrer un serveur HTTPS sécurisé
3. Connecter un client avec vérification TLS
4. Configurer pour la production

## 🚀 Démarrage Rapide (5 minutes)

### Étape 1 : Générer les certificats

```bash
# Générer des certificats auto-signés pour développement
tsd auth generate-cert

# Vérifier que les fichiers sont créés
ls -lh certs/
# server.crt (certificat serveur)
# server.key (clé privée serveur)
# ca.crt (certificat CA pour les clients)
```

### Étape 2 : Démarrer le serveur HTTPS

```bash
# Le serveur démarre automatiquement en HTTPS s'il trouve les certificats
tsd server

# Ou spécifier explicitement les certificats
tsd server --tls-cert ./certs/server.crt --tls-key ./certs/server.key

# Sortie attendue :
# [TSD-SERVER] 🚀 Démarrage du serveur TSD sur https://0.0.0.0:8080
# [TSD-SERVER] 🔒 TLS: activé
# [TSD-SERVER]    Certificat: ./certs/server.crt
# [TSD-SERVER]    Clé: ./certs/server.key
```

### Étape 3 : Utiliser le client HTTPS

Dans un autre terminal :

```bash
# Créer un fichier de test
cat > test.tsd << 'EOF'
type Person : <id: string, name: string, age: int>

fact p1 : Person <id: "1", name: "Alice", age: 30>
fact p2 : Person <id: "2", name: "Bob", age: 25>

rule check_adult : {p: Person} / p.age >= 18 ==> adult(p.id)
EOF

# Option 1 : Mode insecure (développement avec certificats auto-signés)
tsd client test.tsd -insecure

# Option 2 : Avec vérification du CA
tsd client test.tsd -tls-ca ./certs/ca.crt

# Option 3 : Variable d'environnement
export TSD_TLS_CA=./certs/ca.crt
tsd client test.tsd
```

## 🔧 Options de Configuration

### Génération de Certificats

```bash
# Personnaliser les hôtes autorisés
tsd auth generate-cert -hosts "localhost,127.0.0.1,192.168.1.100"

# Personnaliser la durée de validité (en jours)
tsd auth generate-cert -valid-days 730

# Répertoire de sortie personnalisé
tsd auth generate-cert -output-dir ./my-certs

# Spécifier l'organisation
tsd auth generate-cert -org "My Company"

# Tout en une fois
tsd auth generate-cert \
  -hosts "localhost,127.0.0.1,myserver.local" \
  -valid-days 365 \
  -org "My Company" \
  -output-dir ./certs
```

### Configuration du Serveur

```bash
# Certificats par défaut (./certs/server.{crt,key})
tsd server

# Certificats personnalisés
tsd server \
  --tls-cert /path/to/cert.crt \
  --tls-key /path/to/key.key

# Variables d'environnement
export TSD_TLS_CERT=/path/to/cert.crt
export TSD_TLS_KEY=/path/to/key.key
tsd server

# Mode HTTP non sécurisé (⚠️ développement uniquement)
tsd server --insecure
```

### Configuration du Client

```bash
# HTTPS par défaut avec mode insecure (certificats auto-signés)
tsd client test.tsd -insecure

# Avec vérification du CA
tsd client test.tsd -tls-ca ./certs/ca.crt

# Variables d'environnement
export TSD_TLS_CA=./certs/ca.crt
export TSD_CLIENT_INSECURE=false
tsd client test.tsd

# Serveur distant
tsd client test.tsd \
  -server https://tsd.example.com:8080 \
  -tls-ca /path/to/ca.crt

# Mode HTTP non sécurisé
tsd client test.tsd -server http://localhost:8080
```

## 🏭 Configuration Production

### Avec Let's Encrypt

```bash
# 1. Obtenir un certificat Let's Encrypt
sudo certbot certonly --standalone -d tsd.example.com

# 2. Démarrer le serveur avec le certificat
tsd server \
  --tls-cert /etc/letsencrypt/live/tsd.example.com/fullchain.pem \
  --tls-key /etc/letsencrypt/live/tsd.example.com/privkey.pem \
  --port 443

# 3. Le client n'a pas besoin de --insecure (certificat valide)
tsd client test.tsd -server https://tsd.example.com
```

### Avec un Certificat d'Entreprise

```bash
# 1. Obtenir votre certificat auprès de votre CA interne
# (fichiers : company-cert.crt, company-key.key, company-ca.crt)

# 2. Démarrer le serveur
tsd server \
  --tls-cert /etc/tsd/company-cert.crt \
  --tls-key /etc/tsd/company-key.key

# 3. Configurer le client avec le CA d'entreprise
tsd client test.tsd \
  -server https://tsd.internal.company.com:8080 \
  -tls-ca /etc/ssl/certs/company-ca.crt
```

## 🔒 Bonnes Pratiques de Sécurité

### Gestion des Clés Privées

1. **Permissions restrictives** :
   ```bash
   chmod 600 certs/server.key
   chown tsd:tsd certs/server.key
   ```

2. **Ne JAMAIS committer dans Git** :
   ```bash
   # Déjà dans .gitignore
   certs/
   *.key
   *.crt
   *.pem
   ```

3. **Utiliser des secrets managers en production** :
   ```bash
   # Exemple avec HashiCorp Vault
   vault kv get -field=tls_cert secret/tsd/certs > /tmp/cert.crt
   vault kv get -field=tls_key secret/tsd/certs > /tmp/key.key
   tsd server --tls-cert /tmp/cert.crt --tls-key /tmp/key.key
   ```

### Rotation des Certificats

```bash
# 1. Générer de nouveaux certificats
tsd auth generate-cert -output-dir ./certs-new

# 2. Tester avec les nouveaux certificats
tsd server --tls-cert ./certs-new/server.crt --tls-key ./certs-new/server.key

# 3. Si OK, remplacer les anciens
mv certs certs-old
mv certs-new certs

# 4. Redémarrer le serveur
```

### Sécurité Renforcée

```bash
# 1. Utiliser TLS 1.3 uniquement (si votre Go le supporte)
# (configuré automatiquement dans le code, MinVersion = TLS 1.2)

# 2. Combiner avec authentification
tsd server \
  --tls-cert ./certs/server.crt \
  --tls-key ./certs/server.key \
  --auth jwt \
  --jwt-secret "$(openssl rand -base64 32)"

# 3. Utiliser des certificats clients (mTLS) - fonctionnalité future
```

## 🐛 Dépannage

### Erreur : "Certificat TLS non trouvé"

```bash
❌ Certificat TLS non trouvé: ./certs/server.crt

💡 Solutions:
   1. Générer des certificats: tsd auth generate-cert
   2. Spécifier un certificat: --tls-cert /path/to/cert.crt
   3. Démarrer en mode non sécurisé: --insecure (déconseillé en production)
```

**Solution** :
```bash
tsd auth generate-cert
```

### Erreur : "x509: certificate signed by unknown authority"

```bash
❌ Erreur communication serveur: Get "https://localhost:8080/api/v1/execute": 
    x509: certificate signed by unknown authority
```

**Solutions** :
```bash
# Option 1 : Utiliser le mode insecure (développement)
tsd client test.tsd -insecure

# Option 2 : Spécifier le CA
tsd client test.tsd -tls-ca ./certs/ca.crt

# Option 3 : Variable d'environnement
export TSD_TLS_CA=./certs/ca.crt
tsd client test.tsd
```

### Erreur : "certificate is valid for localhost, not 192.168.1.100"

Le certificat ne contient pas l'IP/hostname utilisé.

**Solution** :
```bash
# Régénérer avec les bons hôtes
tsd auth generate-cert -hosts "localhost,127.0.0.1,192.168.1.100"
```

### Vérifier un Certificat

```bash
# Afficher les détails du certificat
openssl x509 -in certs/server.crt -text -noout | less

# Vérifier la validité
openssl x509 -in certs/server.crt -noout -checkend 0

# Vérifier les hôtes autorisés
openssl x509 -in certs/server.crt -noout -text | grep -A1 "Subject Alternative Name"
```

### Tester la connexion TLS

```bash
# Tester avec OpenSSL
openssl s_client -connect localhost:8080 -servername localhost

# Tester avec curl
curl -v --cacert ./certs/ca.crt https://localhost:8080/health

# Tester sans vérification (dev)
curl -k https://localhost:8080/health
```

## 📚 Ressources Complémentaires

- [Documentation TLS/HTTPS](../../docs/AUTHENTICATION_TUTORIAL.md#23-configuration-tlshttps-requis)
- [Guide d'authentification complet](../../docs/AUTHENTICATION_TUTORIAL.md)
- [README principal](../../README.md#-tlshttps-nouveau)

## ⚠️ Avertissements Importants

1. **Développement** :
   - Les certificats auto-signés sont OK pour le développement local
   - Utilisez `-insecure` ou `-tls-ca` avec le client
   - Ne jamais committer les certificats dans Git

2. **Production** :
   - TOUJOURS utiliser des certificats signés par une CA reconnue
   - JAMAIS utiliser `--insecure` ou `-insecure`
   - Mettre en place une rotation automatique des certificats
   - Utiliser un gestionnaire de secrets pour les clés privées

3. **Sécurité** :
   - Les clés privées doivent avoir les permissions 600
   - Stocker les secrets dans un vault sécurisé
   - Monitorer l'expiration des certificats
   - Combiner TLS avec l'authentification (JWT ou Auth Key)

## 🎯 Prochaines Étapes

1. [Configurer l'authentification](../../docs/AUTHENTICATION_TUTORIAL.md)
2. [Déployer en production](../../docs/DEPLOYMENT.md) (à venir)
3. [Configurer le monitoring](../../docs/MONITORING.md) (à venir)