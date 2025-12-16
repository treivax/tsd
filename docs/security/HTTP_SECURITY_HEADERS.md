# Headers de Sécurité HTTP

## 📋 Vue d'ensemble

Le serveur TSD implémente 7 headers de sécurité HTTP critiques pour protéger l'API contre les principales vulnérabilités web identifiées par l'OWASP.

## 🔒 Headers Implémentés

### 1. Strict-Transport-Security (HSTS)

**Valeur** : `max-age=31536000; includeSubDomains`

**Protection** : Force HTTPS et empêche les downgrade attacks

**Détails** :
- `max-age=31536000` : 1 an (365 jours)
- `includeSubDomains` : Applique HSTS à tous les sous-domaines
- Le navigateur refusera toute connexion HTTP pendant 1 an après la première visite en HTTPS

**Référence** : [RFC 6797](https://tools.ietf.org/html/rfc6797)

---

### 2. X-Content-Type-Options

**Valeur** : `nosniff`

**Protection** : Empêche le MIME sniffing

**Détails** :
- Force le navigateur à respecter strictement le Content-Type déclaré
- Empêche l'exécution de JavaScript déguisé en image, JSON, etc.
- Protection critique contre l'injection de contenu malveillant

**Référence** : [Mozilla MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Content-Type-Options)

---

### 3. X-Frame-Options

**Valeur** : `DENY`

**Protection** : Bloque le clickjacking

**Détails** :
- Empêche totalement l'affichage de l'API dans une iframe
- Protection contre les attaques par superposition d'interface (UI redressing)
- Mode `DENY` plus strict que `SAMEORIGIN`

**Référence** : [RFC 7034](https://tools.ietf.org/html/rfc7034)

---

### 4. Content-Security-Policy (CSP)

**Valeur** : `default-src 'none'; frame-ancestors 'none'`

**Protection** : Politique de sécurité stricte pour API

**Détails** :
- `default-src 'none'` : Bloque toutes les ressources par défaut
- `frame-ancestors 'none'` : Empêche l'inclusion dans des frames (remplace X-Frame-Options)
- Adapté à une API JSON pure (pas de HTML/JS à servir)
- Version moderne et plus puissante que X-Frame-Options

**Note** : Si le serveur sert du contenu HTML/JS à l'avenir, cette policy devra être adaptée.

**Référence** : [W3C CSP Level 3](https://www.w3.org/TR/CSP3/)

---

### 5. X-XSS-Protection

**Valeur** : `1; mode=block`

**Protection** : Protection XSS pour navigateurs legacy

**Détails** :
- `1` : Active la protection XSS
- `mode=block` : Bloque complètement la page en cas de détection XSS
- Header legacy maintenu pour compatibilité avec anciens navigateurs
- Modern browsers utilisent CSP à la place

**Référence** : [Mozilla MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-XSS-Protection)

---

### 6. Referrer-Policy

**Valeur** : `no-referrer`

**Protection** : Contrôle des informations de navigation

**Détails** :
- Aucune information de referrer n'est envoyée aux ressources externes
- Protège la confidentialité des URLs contenant potentiellement des tokens/IDs
- Empêche la fuite d'informations sensibles dans les requêtes sortantes

**Référence** : [W3C Referrer Policy](https://www.w3.org/TR/referrer-policy/)

---

### 7. Server

**Valeur** : `TSD`

**Protection** : Masquage de la version du serveur

**Détails** :
- Affiche uniquement "TSD" au lieu de "Go http/1.1" ou version détaillée
- Réduit la surface d'attaque en cachant les détails d'implémentation
- Empêche le fingerprinting automatisé du serveur

**Référence** : [OWASP - Information Leakage](https://owasp.org/www-community/vulnerabilities/Information_exposure_through_server_headers)

---

## 🧪 Tests

### Tests Unitaires

Le fichier `internal/servercmd/servercmd_test.go` contient des tests complets :

1. **TestSecurityHeaders** : Vérifie chaque header individuellement
2. **TestSecurityHeadersOnAllEndpoints** : Vérifie que les headers sont présents sur tous les endpoints

```bash
# Exécuter les tests
go test ./internal/servercmd/... -v -run TestSecurityHeaders

# Avec couverture
go test ./internal/servercmd/... -v -cover
```

### Test Manuel

```bash
# Démarrer le serveur (avec certificats TLS configurés)
tsd server

# Vérifier les headers avec curl
curl -I https://localhost:8443/health

# Exemple de réponse attendue :
# HTTP/2 200
# Strict-Transport-Security: max-age=31536000; includeSubDomains
# X-Content-Type-Options: nosniff
# X-Frame-Options: DENY
# Content-Security-Policy: default-src 'none'; frame-ancestors 'none'
# X-XSS-Protection: 1; mode=block
# Referrer-Policy: no-referrer
# Server: TSD
# Content-Type: application/json
```

## 📊 Validation Externe

### Mozilla Observatory

Pour obtenir un score A+ sur Mozilla Observatory :

```bash
# Scanner votre serveur TSD déployé
https://observatory.mozilla.org/analyze/votre-domaine.com
```

**Score attendu** : A+ avec les 7 headers implémentés

### Security Headers

```bash
# Vérifier via securityheaders.com
https://securityheaders.com/?q=votre-domaine.com
```

**Score attendu** : A+ (100/100)

## 🔧 Configuration

### Mode Développement

En mode développement avec `--insecure`, HSTS reste actif mais n'a pas d'effet car le serveur utilise HTTP.

**Recommandation** : Toujours tester avec TLS activé avant la production.

### Personnalisation

Les valeurs des headers sont définies comme constantes dans `internal/servercmd/servercmd.go` :

```go
const (
    HeaderStrictTransportSecurity = "Strict-Transport-Security"
    ValueHSTS = "max-age=31536000; includeSubDomains"
    // ... autres headers
)
```

**Note** : Modifier ces valeurs nécessite une justification de sécurité et validation des tests.

## ⚠️ Points d'Attention

### HSTS et Développement Local

Le header HSTS peut compliquer le développement local :

- Une fois HSTS activé, le navigateur refusera HTTP pendant 1 an
- Solution : Utiliser un domaine différent pour dev/prod
- Solution alternative : Nettoyer les données HSTS du navigateur

**Chrome** : `chrome://net-internals/#hsts` → Delete domain  
**Firefox** : Effacer l'historique récent → Cookies et données de sites

### CSP et Contenu Futur

La CSP actuelle (`default-src 'none'`) est stricte et adaptée à une API JSON.

**Si vous devez servir du HTML/JS** :
1. Modifier `ValueCSP` dans `servercmd.go`
2. Exemple : `default-src 'self'; script-src 'self'; style-src 'self'`
3. Mettre à jour les tests correspondants
4. Documenter le changement dans CHANGELOG.md

### Compatibilité

| Header | IE | Edge | Chrome | Firefox | Safari |
|--------|----|----|--------|---------|--------|
| HSTS | 11+ | ✅ | ✅ | ✅ | ✅ |
| X-Content-Type-Options | 8+ | ✅ | ✅ | ✅ | ✅ |
| X-Frame-Options | ✅ | ✅ | ✅ | ✅ | ✅ |
| CSP | 10+ | ✅ | ✅ | ✅ | ✅ |
| X-XSS-Protection | 8+ | ✅ | ✅ | ✅ | ✅ |
| Referrer-Policy | ❌ | ✅ | ✅ | ✅ | ✅ |

**Compatibilité globale** : 95%+ des navigateurs modernes

## 📚 Références

### Standards et Spécifications

- [OWASP Secure Headers Project](https://owasp.org/www-project-secure-headers/)
- [RFC 6797 - HTTP Strict Transport Security](https://tools.ietf.org/html/rfc6797)
- [RFC 7034 - X-Frame-Options](https://tools.ietf.org/html/rfc7034)
- [W3C CSP Level 3](https://www.w3.org/TR/CSP3/)
- [W3C Referrer Policy](https://www.w3.org/TR/referrer-policy/)

### Outils de Validation

- [Mozilla Observatory](https://observatory.mozilla.org/)
- [Security Headers](https://securityheaders.com/)
- [SSL Labs Server Test](https://www.ssllabs.com/ssltest/)

### Documentation Projet

- [SECURITY.md](../../SECURITY.md) - Politique de sécurité du projet
- [CHANGELOG.md](../../CHANGELOG.md) - Historique des modifications
- [servercmd.go](../../internal/servercmd/servercmd.go) - Implémentation
- [servercmd_test.go](../../internal/servercmd/servercmd_test.go) - Tests

---

**Dernière mise à jour** : 2025-12-15  
**Statut** : ✅ Implémenté et testé  
**Couverture de tests** : 100% des headers
