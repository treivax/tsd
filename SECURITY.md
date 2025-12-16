# 🔒 Security Policy

> **Version actuelle** : TSD 1.0.0  
> **Dernière mise à jour** : 16 décembre 2024

---

## 📋 Table des Matières

1. [Versions Supportées](#-versions-supportées)
2. [Reporting d'une Vulnérabilité](#-reporting-dune-vulnérabilité)
3. [Process de Gestion](#-process-de-gestion)
4. [Divulgation Responsable](#-divulgation-responsable)
5. [Best Practices de Déploiement](#️-best-practices-de-déploiement)
6. [Hall of Fame](#-hall-of-fame)
7. [Ressources](#-ressources)
8. [Contact](#-contact)

---

## 📊 Versions Supportées

Nous prenons la sécurité de TSD très au sérieux. Voici les versions actuellement supportées avec des mises à jour de sécurité :

| Version | Supportée          | Fin de Support     | Notes |
| ------- | ------------------ | ------------------ | ----- |
| 1.0.x   | ✅ Oui             | En cours           | Version stable actuelle |
| 0.x     | ⚠️ Maintenance limitée | 30 juin 2025   | Migration recommandée vers 1.0.x |
| < 0.x   | ❌ Non             | Non supporté       | Aucun correctif de sécurité |

### Politique de Support

**Version actuelle (1.x)** :
- Correctifs de sécurité pour toutes les vulnérabilités
- Mises à jour régulières et patches mineurs
- Support complet et documentation

**Version précédente (0.x)** :
- Correctifs uniquement pour vulnérabilités critiques et hautes
- Support limité jusqu'au 30 juin 2025
- **Migration fortement recommandée vers 1.0.x**

**Versions obsolètes (< 0.x)** :
- Aucun support
- Aucun correctif de sécurité
- **Mise à jour obligatoire**

### Cycle de Vie des Versions

À partir de la version 1.0, nous suivons :
- **Version majeure actuelle** : Support complet
- **Version majeure N-1** : Support de sécurité pendant 6 mois après release N+1
- **Versions plus anciennes** : Fin de support

---

## 🚨 Reporting d'une Vulnérabilité

### ⚠️ Important : NE PAS créer d'issue publique

Si vous découvrez une vulnérabilité de sécurité dans TSD, merci de **NE PAS** créer d'issue publique sur GitHub. Cela pourrait mettre en danger les utilisateurs du projet.

### 📧 Méthodes de Contact Privé

**Option 1 - GitHub Security Advisory** (recommandé) :
1. Accédez à la page Security du projet GitHub
2. Cliquez sur "Report a vulnerability"
3. Remplissez le formulaire privé avec les détails

**Option 2 - Email sécurisé** :
- Créez une issue privée via GitHub Security Advisory
- Ou contactez directement les mainteneurs du projet via GitHub

**Option 3 - Communication directe** :
- Contactez directement les mainteneurs du projet
- Demandez un canal de communication privé

### 📝 Informations à Fournir

Pour nous aider à traiter rapidement et efficacement votre rapport, veuillez inclure :

**Template de Rapport** :

\`\`\`markdown
Sujet : [SECURITY] [Gravité: Critique/Haute/Moyenne/Basse] Description courte

## 📌 Résumé
Brève description de la vulnérabilité (1-2 lignes)

## 🎯 Impact
- **Gravité** : Critique / Haute / Moyenne / Basse (selon CVSS)
- **CVSS Score** : X.X (si calculé)
- **Type** : Injection / XSS / CSRF / DoS / Autre
- **Vecteur d'attaque** : Local / Adjacent / Network / Physical
- **Privilèges requis** : None / Low / High
- **Interaction utilisateur** : None / Required

**Description de l'impact** :
[Que peut faire un attaquant ? Quelles données sont compromises ?]

## 🔍 Description Technique
[Description détaillée de la vulnérabilité]

**Composants affectés** :
- Module : [rete / constraint / auth / autre]
- Fichier(s) : [chemin/vers/fichier.go]
- Fonction(s) : [NomFonction()]
- Ligne(s) : [numéros de lignes]

## 🔄 Étapes de Reproduction

**Prérequis** :
- TSD version : X.X.X
- OS : Linux / macOS / Windows
- Go version : 1.21+
- Configuration particulière : [si applicable]

**Étapes** :
1. [Étape détaillée 1]
2. [Étape détaillée 2]
3. [Étape détaillée 3]
...

**Résultat obtenu** :
[Ce qui se passe actuellement]

**Résultat attendu** :
[Ce qui devrait se passer]

## 💻 Preuve de Concept (PoC)

\`\`\`go
// Code de démonstration
\`\`\`

\`\`\`bash
# Commandes à exécuter
\`\`\`

**Logs / Captures d'écran** :
[Si applicable, joindre logs ou captures]

## 🔧 Versions Affectées
- TSD v1.0.0 : ✅ Vulnérable
- TSD v0.9.x : ✅ Vulnérable
- TSD v0.8.x : ❓ Non testé
- TSD < 0.8 : ❓ Non testé

## 💡 Suggestions de Correctif
[Optionnel : vos idées pour corriger la vulnérabilité]

**Patch proposé** :
\`\`\`go
// Si vous avez un correctif
\`\`\`

## 📚 Références
- [CVE similaires]
- [Documentation technique]
- [Articles de recherche]

## 👤 Informations du Reporter
- **Nom/Pseudo** : [votre nom ou pseudo]
- **Organisation** : [optionnel]
- **Email** : [pour communication]
- **Préférence de crédit** : Public / Anonyme / Pseudo uniquement
\`\`\`

### ⏱️ Délai de Réponse

Nous nous engageons à :

- **Accusé de réception** : Sous 48 heures (2 jours ouvrés)
- **Évaluation initiale** : Sous 7 jours
- **Plan de correction** : 
  - Critique : 7 jours
  - Haute : 14 jours
  - Moyenne : 30 jours
  - Basse : 60 jours

**Note** : Ces délais sont des objectifs. Nous communiquerons régulièrement sur l'avancement.

---

## 🔄 Process de Gestion

### Timeline de Traitement

\`\`\`
J+0 → Réception
 ↓
J+2 → Accusé de réception + ID de tracking
 ↓
J+7 → Évaluation (reproduction, gravité CVSS, versions affectées)
 ↓
J+14-60 → Développement correctif (selon gravité)
 ↓    → Tests de non-régression
 ↓    → Review interne sécurité
 ↓
J+30-90 → Coordination avec reporter
 ↓    → Validation du correctif
 ↓    → Planification release
 ↓
J+45-120 → Publication
 ↓     → Release version corrigée
 ↓     → Advisory de sécurité publique
 ↓     → Notification utilisateurs
 ↓
      → Clôture
\`\`\`

### Étapes Détaillées

#### 1. Réception (J+0)
- Réception du rapport
- Attribution d'un identifiant unique (ex: TSD-SEC-2024-001)
- Assignation à un responsable sécurité

#### 2. Accusé de Réception (J+2)
- Confirmation de réception au reporter
- Fourniture de l'ID de tracking
- Demande de clarifications si nécessaire

#### 3. Évaluation (J+7)
- **Reproduction** : Validation que la vulnérabilité existe
- **Gravité** : Calcul du score CVSS v3.1
- **Périmètre** : Identification des versions affectées
- **Impact** : Évaluation de l'impact réel
- **Classification** : Attribution du niveau de priorité

#### 4. Développement du Correctif (J+14 à J+60)

Selon la gravité :

| Gravité | CVSS Score | Délai Correctif | Priorité |
|---------|------------|-----------------|----------|
| 🔴 **Critique** | 9.0-10.0 | 7 jours | P0 - Immédiate |
| 🟠 **Haute** | 7.0-8.9 | 14 jours | P1 - Urgente |
| 🟡 **Moyenne** | 4.0-6.9 | 30 jours | P2 - Normale |
| 🟢 **Basse** | 0.1-3.9 | 60 jours | P3 - Planifiée |

**Processus de développement** :
- Création d'une branche privée
- Développement du correctif
- Tests unitaires et d'intégration
- Tests de non-régression complets
- Review de code sécurité (double validation)
- Validation par au moins 2 mainteneurs

#### 5. Coordination (J+30 à J+90)
- **Communication avec le reporter** :
  - Partage du correctif pour validation
  - Accord sur le timing de divulgation
  - Discussion sur les crédits
- **Préparation de l'advisory** :
  - Rédaction de l'advisory de sécurité
  - Identification CVE si applicable
  - Préparation des release notes

#### 6. Publication (J+45 à J+120)
- **Release** :
  - Publication de la version corrigée
  - Tag git avec mention sécurité
  - Build et distribution des binaires
- **Advisory** :
  - Publication GitHub Security Advisory
  - Entrée dans CHANGELOG.md
  - Notification sur les canaux du projet
- **Communication** :
  - Annonce sur GitHub Discussions/Issues (pinned)
  - Email aux utilisateurs connus (si liste existe)
  - Mise à jour de la documentation

### Communication Continue

Pendant tout le processus :
- Mises à jour régulières au reporter (au moins hebdomadaires)
- Transparence sur les délais et obstacles
- Coordination étroite pour la divulgation

---

## 🤝 Divulgation Responsable

### Principes de Confidentialité

Nous demandons aux chercheurs en sécurité de respecter les principes suivants :

#### ✅ À Faire

- **Garder la confidentialité** : Ne pas divulguer publiquement la vulnérabilité avant notre accord
- **Délai raisonnable** : Nous laisser un délai suffisant pour corriger (selon gravité)
- **Coordonner** : Travailler avec nous pour la divulgation publique
- **Communication privée** : Utiliser uniquement les canaux sécurisés listés
- **Bienveillance** : Agir de bonne foi pour protéger les utilisateurs

#### ❌ À Ne Pas Faire

- **Exploitation** : Ne pas exploiter la vulnérabilité à des fins malveillantes ou personnelles
- **Divulgation prématurée** : Ne pas divulguer publiquement avant notre accord
- **Accès non autorisé** : Ne pas accéder aux données utilisateurs ou systèmes de production
- **Déni de service** : Ne pas lancer d'attaques DoS ou perturber le service
- **Escalade** : Ne pas effectuer d'actions au-delà du nécessaire pour démontrer la vulnérabilité

### Politique de Non-Poursuite

**Nous nous engageons à** :

- ✅ Ne pas poursuivre légalement les chercheurs qui suivent cette politique
- ✅ Travailler de bonne foi avec la communauté sécurité
- ✅ Respecter les délais convenus pour la divulgation
- ✅ Créditer publiquement les découvertes (sauf demande d'anonymat)
- ✅ Traiter tous les rapports avec sérieux et respect

**Activités autorisées dans le cadre de la recherche** :
- Tests sur votre propre installation locale de TSD
- Analyse statique du code source
- Review de code pour identifier des vulnérabilités potentielles
- Proof of Concept démontrant la vulnérabilité (sans exploitation réelle)

**Activités interdites** :
- Tests sur des instances de production sans autorisation
- Accès à des données utilisateurs réelles
- Attaques par déni de service
- Social engineering des utilisateurs ou mainteneurs

### Attribution et Crédits

Nous créditons publiquement les chercheurs qui :

- Reportent de manière responsable via les canaux appropriés
- Respectent notre processus de divulgation coordonnée
- Le souhaitent (nous respectons l'anonymat si demandé)

**Les crédits sont inclus dans** :

1. **GitHub Security Advisory** : Nom/pseudo du reporter
2. **CHANGELOG.md** : Section dédiée aux contributions sécurité
3. **Release Notes** : Mention dans les notes de version
4. **Hall of Fame** : Section ci-dessous de ce document
5. **Commit message** : Attribution dans le message de commit du correctif

**Options de crédit** :
- Nom complet + organisation
- Pseudo uniquement
- Lien vers profil GitHub/Twitter
- Anonyme (aucune mention publique)

### Coordination de la Divulgation

Nous travaillons avec le reporter pour :

1. **Valider le correctif** : Le reporter peut tester la version corrigée avant publication
2. **Timing** : Coordonner la date et heure de divulgation publique
3. **Contenu** : Rédiger conjointement l'advisory si souhaité
4. **Communication** : Décider du niveau de détail à divulguer
5. **CVE** : Coordonner l'attribution CVE si applicable

**Délais de divulgation par défaut** :

| Gravité | Délai depuis rapport | Délai depuis correctif |
|---------|---------------------|----------------------|
| Critique | 30 jours | 7 jours |
| Haute | 45 jours | 14 jours |
| Moyenne | 60 jours | 21 jours |
| Basse | 90 jours | 30 jours |

**Note** : Ces délais peuvent être ajustés en accord avec le reporter.

---

## 🛡️ Best Practices de Déploiement

### Configuration Sécurisée

#### 1. Authentification et Autorisation

**JWT et API Keys** :

\`\`\`bash
# Générer une clé API forte
tsd auth generate-api-key

# Générer un JWT avec durée limitée (60 minutes)
tsd auth generate-jwt --api-key=YOUR_API_KEY --duration=60

# Utiliser des rôles pour les permissions
tsd auth generate-jwt --api-key=KEY --roles=admin,editor
\`\`\`

**Rotation des secrets** :
- Clés API : rotation tous les 90 jours minimum
- JWT secret : rotation tous les 180 jours
- Certificats TLS : renouvellement avant expiration

#### 2. TLS/HTTPS en Production

**Toujours utiliser HTTPS** :

\`\`\`bash
# Avec certificats auto-signés (développement uniquement)
tsd auth generate-cert --host localhost --days 365

# En production : utilisez Let's Encrypt ou certificats valides
tsd server \\
  --tls-cert=/path/to/fullchain.pem \\
  --tls-key=/path/to/privkey.pem \\
  --port=8443
\`\`\`

**Recommandations TLS** :
- ✅ TLS 1.2+ uniquement (désactiver TLS 1.0 et 1.1)
- ✅ Cipher suites fortes uniquement
- ✅ HSTS (HTTP Strict Transport Security)
- ✅ Certificats valides (pas auto-signés en production)

#### 3. Gestion des Secrets

**Jamais en clair** :

\`\`\`bash
# ❌ MAUVAIS - secrets en ligne de commande
tsd server --jwt-secret=my-secret-123

# ✅ BON - secrets via variables d'environnement
export TSD_JWT_SECRET=\$(cat /secure/path/jwt-secret.txt)
tsd server

# ✅ MEILLEUR - secrets via gestionnaire de secrets
# (Vault, AWS Secrets Manager, etc.)
\`\`\`

**Fichiers de configuration** :
- Ne jamais committer les secrets dans Git
- Utiliser \`.env.example\` avec valeurs factices
- Ajouter \`.env\` dans \`.gitignore\`
- Permissions restrictives : \`chmod 600 config/secrets.yaml\`

#### 4. Réseau et Firewall

**Exposition minimale** :

\`\`\`bash
# Écouter uniquement sur localhost si pas besoin d'accès externe
tsd server --host 127.0.0.1 --port 8080

# Avec reverse proxy (recommandé)
# TSD écoute sur localhost:8080
# Nginx/Caddy expose HTTPS sur 443
\`\`\`

**Firewall** :
\`\`\`bash
# Autoriser uniquement port HTTPS
ufw allow 443/tcp
ufw enable

# Ou avec iptables
iptables -A INPUT -p tcp --dport 443 -j ACCEPT
iptables -A INPUT -p tcp --dport 80 -j ACCEPT  # Redirection HTTP→HTTPS
\`\`\`

**Rate Limiting** :
- Limiter les requêtes par IP
- Utiliser un reverse proxy (Nginx, Caddy) avec rate limiting
- Protéger les endpoints sensibles (auth, compilation)

#### 5. Logging et Monitoring

**Activer les logs sécurité** :

\`\`\`bash
# Logs détaillés pour audit
tsd server --log-level=info --log-format=json --log-file=/var/log/tsd/security.log
\`\`\`

**Surveiller** :
- Tentatives d'authentification échouées
- Utilisation suspecte des API keys
- Erreurs de validation JWT
- Accès à des ressources non autorisées
- Patterns de requêtes anormaux

**Alertes** :
- Configurer des alertes sur les échecs d'authentification répétés
- Notifier sur les erreurs critiques
- Monitorer l'utilisation des ressources (CPU, mémoire, disque)

#### 6. Mises à Jour

**Maintenir TSD à jour** :

\`\`\`bash
# Vérifier la version actuelle
tsd --version

# Surveiller les releases
# GitHub : https://github.com/OWNER/tsd/releases
# RSS : https://github.com/OWNER/tsd/releases.atom
\`\`\`

**Processus de mise à jour** :
1. Lire le CHANGELOG et release notes
2. Vérifier les breaking changes
3. Tester en environnement de staging
4. Backup des données et configuration
5. Mise à jour en production
6. Vérifier les logs post-déploiement

### Audits de Sécurité

#### Outils Automatisés

**Scanner les vulnérabilités Go** :

\`\`\`bash
# govulncheck - détecte vulnérabilités dans dépendances
go install golang.org/x/vuln/cmd/govulncheck@latest
govulncheck ./...

# Intégré dans le Makefile
make security-vulncheck
\`\`\`

**Analyse statique de sécurité** :

\`\`\`bash
# gosec - analyse le code pour problèmes de sécurité
go install github.com/securego/gosec/v2/cmd/gosec@latest
gosec -exclude-dir=tests ./...

# Intégré dans le Makefile
make security-gosec

# Scan complet
make security-scan  # Exécute gosec + govulncheck
\`\`\`

**Scanner les dépendances** :

\`\`\`bash
# Nancy - vérifie vulnérabilités dans go.mod
go list -json -m all | nancy sleuth

# Dependabot - activé sur GitHub pour alertes automatiques
\`\`\`

#### Audits Manuels

**Review de code sécurité** :
- [ ] Validation de toutes les entrées utilisateur
- [ ] Gestion des erreurs robuste (pas de panic)
- [ ] Pas d'injection SQL/commande/code
- [ ] Gestion des cas nil/vides
- [ ] Pas de race conditions
- [ ] Pas de fuites mémoire
- [ ] Secrets jamais en clair dans le code

**Checklist sécurité** :
\`\`\`bash
# Exécuter la validation complète
make validate

# Inclut :
# - go fmt / goimports
# - go vet
# - staticcheck
# - errcheck
# - gosec
# - govulncheck
# - Tests complets (> 80% couverture)
\`\`\`

### Déploiement en Production

#### Environnement Sécurisé

**Utilisateur dédié** :

\`\`\`bash
# Créer utilisateur système sans shell
sudo useradd -r -s /bin/false tsd

# Permissions restrictives
sudo chown -R tsd:tsd /opt/tsd
sudo chmod 750 /opt/tsd
\`\`\`

**Systemd service** :

\`\`\`ini
[Unit]
Description=TSD Server
After=network.target

[Service]
Type=simple
User=tsd
Group=tsd
WorkingDirectory=/opt/tsd
ExecStart=/opt/tsd/bin/tsd server --config=/etc/tsd/config.yaml
Restart=on-failure
RestartSec=10

# Sécurité
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/lib/tsd /var/log/tsd

[Install]
WantedBy=multi-user.target
\`\`\`

**Ressources limitées** :

\`\`\`ini
# Dans [Service]
LimitNOFILE=65536
LimitNPROC=512
MemoryLimit=2G
CPUQuota=200%
\`\`\`

#### Reverse Proxy (Nginx)

\`\`\`nginx
upstream tsd_backend {
    server 127.0.0.1:8080;
}

server {
    listen 443 ssl http2;
    server_name tsd.example.com;

    # TLS
    ssl_certificate /etc/letsencrypt/live/tsd.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/tsd.example.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Headers de sécurité
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Rate limiting
    limit_req_zone \$binary_remote_addr zone=tsd_limit:10m rate=10r/s;
    limit_req zone=tsd_limit burst=20 nodelay;

    location / {
        proxy_pass http://tsd_backend;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }
}

# Redirection HTTP → HTTPS
server {
    listen 80;
    server_name tsd.example.com;
    return 301 https://\$server_name\$request_uri;
}
\`\`\`

### Backup et Disaster Recovery

**Sauvegardes régulières** :
- Configuration : quotidien
- Données : selon criticité (quotidien/hebdomadaire)
- Secrets : stockage sécurisé hors-ligne

**Plan de récupération** :
- Documenter la procédure de restauration
- Tester régulièrement les backups
- Conserver plusieurs versions

---

## 🏆 Hall of Fame

Merci aux chercheurs en sécurité qui ont contribué à améliorer la sécurité de TSD :

| Date | Chercheur | Vulnérabilité | Gravité | CVE |
|------|-----------|---------------|---------|-----|
| *À venir* | *À venir* | *À venir* | - | - |

<!--
Format pour ajout futur :
| YYYY-MM-DD | Nom/Pseudo (lien GitHub/site) | Description courte | Critique/Haute/Moyenne/Basse | CVE-YYYY-XXXXX |

Exemple :
| 2024-12-16 | [@security_researcher](https://github.com/user) | JWT signature bypass | Critique | CVE-2024-12345 |
-->

**Nous remercions tous les chercheurs en sécurité qui contribuent à protéger TSD et ses utilisateurs.**

---

## �� Ressources

### Documentation Interne

- **[Configuration TLS](docs/security/HTTP_SECURITY_HEADERS.md)** : En-têtes de sécurité HTTP
- **[Scanner de Vulnérabilités](docs/security/VULNERABILITY_SCANNING.md)** : govulncheck et gosec
- **[Guide Auth](auth/README.md)** : Authentification JWT et API keys
- **[CONTRIBUTING.md](CONTRIBUTING.md)** : Standards de code sécurisé

### Standards et Guidelines

**Sécurité Web** :
- [OWASP Top 10](https://owasp.org/www-project-top-ten/) - Top vulnérabilités web
- [OWASP Go Security Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Go_Security_Cheat_Sheet.html)
- [CWE Top 25](https://cwe.mitre.org/top25/) - Faiblesses logicielles communes

**Scoring et CVE** :
- [CVSS v3.1 Calculator](https://www.first.org/cvss/calculator/3.1) - Calculateur CVSS
- [CVE Program](https://cve.mitre.org/) - Programme CVE
- [NVD](https://nvd.nist.gov/) - National Vulnerability Database

**Go Security** :
- [Go Security Policy](https://go.dev/security/policy) - Politique de sécurité Go
- [Go Vulnerability Database](https://pkg.go.dev/vuln/) - Base vulnérabilités Go
- [Secure Go Guidelines](https://github.com/OWASP/Go-SCP) - OWASP Secure Coding Practices

### Outils de Sécurité

**Scanners** :
- [govulncheck](https://pkg.go.dev/golang.org/x/vuln/cmd/govulncheck) - Scanner vulnérabilités Go
- [gosec](https://github.com/securego/gosec) - Analyseur sécurité statique Go
- [Nancy](https://github.com/sonatype-nexus-community/nancy) - Scanner dépendances

**Analyse Statique** :
- [staticcheck](https://staticcheck.io/) - Linter Go avancé
- [errcheck](https://github.com/kisielk/errcheck) - Vérification gestion erreurs
- [golangci-lint](https://golangci-lint.run/) - Meta-linter avec règles sécurité

**TLS/Certificats** :
- [Let's Encrypt](https://letsencrypt.org/) - Certificats TLS gratuits
- [SSL Labs Test](https://www.ssllabs.com/ssltest/) - Test configuration TLS
- [testssl.sh](https://testssl.sh/) - Test TLS en ligne de commande

### Divulgation Responsable

**Guides** :
- [ISO 29147](https://www.iso.org/standard/72311.html) - Vulnerability disclosure
- [GitHub Security Advisories](https://docs.github.com/en/code-security/security-advisories)
- [CERT Coordination Center](https://www.sei.cmu.edu/about/divisions/cert/)

---

## 📧 Contact

### GitHub Security Advisory

**Pour reporter une vulnérabilité** :

1. Accédez à la page Security du projet GitHub
2. Cliquez sur "Report a vulnerability" 
3. Remplissez le formulaire privé avec les détails fournis dans la section "Informations à Fournir"

### Autres Canaux

**GitHub Discussions** : Pour questions générales de sécurité (non sensibles)

**Issue Tracker** : **UNIQUEMENT pour bugs non-sécuritaires**

⚠️ **Important** : Ne jamais reporter de vulnérabilités via issues publiques ou discussions !

---

## ⚖️ Politique de Divulgation

Cette politique de sécurité suit les principes de **divulgation coordonnée responsable** :

### Nos Engagements

**Nous nous engageons à** :

1. ✅ **Traiter tous les rapports sérieusement** - Chaque rapport est évalué et traité
2. ✅ **Répondre rapidement** - Accusé de réception sous 48h, évaluation sous 7 jours
3. ✅ **Communiquer régulièrement** - Mises à jour hebdomadaires minimum
4. ✅ **Travailler de bonne foi** - Collaboration transparente avec les reporters
5. ✅ **Créditer publiquement** - Attribution appropriée (sauf demande d'anonymat)
6. ✅ **Ne pas poursuivre** - Aucune action légale contre chercheurs suivant cette politique
7. ✅ **Divulguer de manière responsable** - Publication coordonnée avec le reporter

### Vos Responsabilités

**Nous demandons aux reporters de** :

1. ✅ **Garder la confidentialité** - Jusqu'à la publication du correctif
2. ✅ **Donner un délai raisonnable** - Permettre la correction avant divulgation
3. ✅ **Coordonner la publication** - Travailler avec nous sur le timing
4. ✅ **Agir de bonne foi** - Objectif de protection, pas d'exploitation
5. ✅ **Utiliser les canaux appropriés** - Rapport privé via GitHub Security Advisory
6. ✅ **Fournir des détails** - Information suffisante pour reproduire
7. ✅ **Respecter les limites** - Pas d'accès non autorisé ou DoS

### Principes de Divulgation

**Divulgation coordonnée** :
- Nous travaillons avec le reporter pour choisir la date de divulgation publique
- Nous respectons les délais standards de l'industrie (30-90 jours selon gravité)
- Nous publions l'advisory et les crédits simultanément avec la release corrigée

**Transparence** :
- Publication d'advisories détaillés après correction
- Documentation des vulnérabilités dans CHANGELOG.md
- Notification proactive des utilisateurs affectés

**Équité** :
- Tous les reporters sont traités équitablement
- Les crédits sont donnés selon les préférences du reporter
- Nous ne discriminons pas selon l'origine du rapport

---

## 📝 Changelog de cette Politique

| Date | Version | Changements |
|------|---------|-------------|
| 2024-12-16 | 1.0 | Version initiale de la politique de sécurité |

---

## 🙏 Remerciements

Nous remercions :

- **La communauté sécurité** pour leur vigilance et contributions
- **Les projets open source** dont nous nous inspirons (Go, Kubernetes, Node.js)
- **Nos utilisateurs** pour leur confiance et leurs retours

---

**Merci de contribuer à la sécurité de TSD ! 🛡️**

**TSD Security Team**  
*Protéger nos utilisateurs est notre priorité.*
