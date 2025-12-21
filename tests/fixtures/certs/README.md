# Certificats de Test TLS

Ce répertoire contient des certificats auto-signés pour les tests TLS uniquement.

## ⚠️ AVERTISSEMENT

**CES CERTIFICATS SONT UNIQUEMENT POUR LES TESTS !**

- ❌ Ne JAMAIS utiliser en production
- ❌ Ne JAMAIS committer dans le dépôt Git
- ✅ Générer localement avec le script fourni

## 📋 Fichiers

- `generate_certs.sh` - Script de génération des certificats
- `test-server.crt` - Certificat auto-signé (généré localement, ignoré par Git)
- `test-server.key` - Clé privée (générée localement, ignorée par Git)

## 🔧 Génération

Pour générer les certificats de test :

```bash
cd tests/fixtures/certs
./generate_certs.sh
```

**Prérequis** : OpenSSL doit être installé sur votre système.

## 🧪 Usage dans les Tests

Les tests TLS utilisent automatiquement ces certificats s'ils existent :

```go
// Les tests TLS vérifient l'existence des certificats
// et skip gracieusement s'ils ne sont pas disponibles
func TestTimeoutsWithTLS(t *testing.T) {
    if testing.Short() {
        t.Skip("⏭️  Test long, skip en mode -short")
    }
    
    certFile, keyFile, skip := createTestCertificates(t)
    if skip {
        t.Skip("⏭️  Certificats de test non disponibles")
    }
    // ... test avec TLS
}
```

## 🔐 Caractéristiques des Certificats

- **Type** : Auto-signé (self-signed)
- **Algorithme** : RSA 2048 bits
- **Hash** : SHA-256
- **Validité** : 365 jours à partir de la génération
- **CN** : localhost
- **Organisation** : TSD Test

## 🔄 Régénération

Les certificats peuvent être régénérés à tout moment en ré-exécutant le script. 
Cela peut être nécessaire si :

- Les certificats ont expiré (après 365 jours)
- Les fichiers ont été supprimés
- Vous voulez changer les paramètres

## 📝 Notes

- Les certificats sont ignorés par `.gitignore` pour éviter de committer des clés
- Chaque développeur doit générer ses propres certificats localement
- Les tests TLS sont automatiquement skippés en CI si les certificats ne sont pas disponibles
- La génération est rapide (< 1 seconde) et peut être faite à la demande

## 🛡️ Sécurité

Ces certificats n'offrent AUCUNE sécurité réelle car :

1. Ils sont auto-signés (non vérifiés par une autorité de certification)
2. La clé privée est générée localement sans protection
3. Ils sont destinés uniquement aux tests fonctionnels

Pour la production, utilisez toujours des certificats émis par une autorité de certification reconnue (Let's Encrypt, etc.).