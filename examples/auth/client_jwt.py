#!/usr/bin/env python3
"""
Exemple d'utilisation du serveur TSD avec authentification JWT

Ce script montre comment :
1. Générer un JWT en Python
2. Se connecter au serveur TSD avec un JWT
3. Gérer l'expiration des tokens
4. Rafraîchir les tokens

Prérequis:
    pip install requests PyJWT

Usage:
    # Option 1: Utiliser un JWT existant
    export TSD_AUTH_TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    python3 client_jwt.py

    # Option 2: Générer un JWT en Python
    export TSD_JWT_SECRET="votre-secret-jwt"
    python3 client_jwt.py --generate --username alice

    # Option 3: Passer le JWT directement
    python3 client_jwt.py --token "eyJhbGciOi..."
"""

import argparse
import json
import os
import sys
from datetime import datetime, timedelta
from typing import Any, Dict, List, Optional

import requests

try:
    import jwt as pyjwt

    JWT_AVAILABLE = True
except ImportError:
    JWT_AVAILABLE = False


class TSDJWTClient:
    """Client Python pour le serveur TSD avec authentification JWT"""

    def __init__(
        self,
        server_url: str = "http://localhost:8080",
        jwt_token: Optional[str] = None,
    ):
        """
        Initialise le client TSD avec JWT

        Args:
            server_url: URL du serveur TSD
            jwt_token: Token JWT
                      Si None, utilise la variable d'environnement TSD_AUTH_TOKEN
        """
        self.server_url = server_url.rstrip("/")
        self.jwt_token = jwt_token or os.getenv("TSD_AUTH_TOKEN")

        if not self.jwt_token:
            raise ValueError(
                "JWT token requis.\n"
                "Définissez-le via:\n"
                "  - Le paramètre jwt_token\n"
                "  - La variable d'environnement TSD_AUTH_TOKEN\n"
                "  - L'option --token en ligne de commande\n"
                "  - Générez-le avec: tsd-auth generate-jwt"
            )

        self.session = requests.Session()
        self.session.headers.update(self._get_headers())

    def _get_headers(self) -> Dict[str, str]:
        """Retourne les headers HTTP avec JWT"""
        return {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {self.jwt_token}",
        }

    def update_token(self, new_token: str) -> None:
        """
        Met à jour le token JWT

        Args:
            new_token: Nouveau token JWT
        """
        self.jwt_token = new_token
        self.session.headers.update(self._get_headers())

    def health_check(self) -> Dict[str, Any]:
        """
        Vérifie la santé du serveur

        Returns:
            Dict contenant le statut, la version et l'uptime

        Raises:
            requests.HTTPError: Si la requête échoue
        """
        try:
            response = self.session.get(f"{self.server_url}/health")
            response.raise_for_status()
            return response.json()
        except requests.HTTPError as e:
            if e.response.status_code == 401:
                raise Exception(
                    "Authentification échouée - JWT invalide ou expiré.\n"
                    "Générez un nouveau token avec: tsd-auth generate-jwt"
                ) from e
            raise

    def get_version(self) -> Dict[str, Any]:
        """
        Récupère la version du serveur

        Returns:
            Dict contenant la version et la version Go
        """
        response = self.session.get(f"{self.server_url}/api/v1/version")
        response.raise_for_status()
        return response.json()

    def execute(
        self, source: str, source_name: str = "<python>", verbose: bool = False
    ) -> Dict[str, Any]:
        """
        Exécute un programme TSD

        Args:
            source: Code source TSD
            source_name: Nom du fichier source (pour les messages d'erreur)
            verbose: Mode verbeux

        Returns:
            Dict contenant les résultats de l'exécution

        Raises:
            requests.HTTPError: Si la requête échoue
        """
        payload = {"source": source, "source_name": source_name, "verbose": verbose}

        try:
            response = self.session.post(
                f"{self.server_url}/api/v1/execute", json=payload, timeout=30
            )
            response.raise_for_status()
            return response.json()
        except requests.HTTPError as e:
            if e.response.status_code == 401:
                raise Exception(
                    "Authentification échouée - JWT invalide ou expiré.\n"
                    "Votre token a peut-être expiré. Générez-en un nouveau."
                ) from e
            raise

    def execute_file(self, file_path: str, verbose: bool = False) -> Dict[str, Any]:
        """
        Exécute un fichier TSD

        Args:
            file_path: Chemin vers le fichier TSD
            verbose: Mode verbeux

        Returns:
            Dict contenant les résultats de l'exécution
        """
        with open(file_path, "r") as f:
            source = f.read()

        return self.execute(source, source_name=file_path, verbose=verbose)


def generate_jwt(
    secret: str,
    username: str,
    roles: Optional[List[str]] = None,
    expiration_hours: int = 24,
    issuer: str = "tsd-server",
) -> str:
    """
    Génère un JWT pour TSD

    Args:
        secret: Secret JWT (doit correspondre au serveur)
        username: Nom d'utilisateur
        roles: Liste des rôles (optionnel)
        expiration_hours: Durée de validité en heures
        issuer: Émetteur du JWT

    Returns:
        Token JWT signé

    Raises:
        ImportError: Si PyJWT n'est pas installé
    """
    if not JWT_AVAILABLE:
        raise ImportError(
            "PyJWT n'est pas installé.\n"
            "Installez-le avec: pip install PyJWT\n"
            "Ou utilisez: tsd-auth generate-jwt"
        )

    now = datetime.utcnow()

    payload = {
        "username": username,
        "roles": roles or [],
        "exp": now + timedelta(hours=expiration_hours),
        "iat": now,
        "nbf": now,
        "iss": issuer,
    }

    token = pyjwt.encode(payload, secret, algorithm="HS256")
    return token


def decode_jwt(
    token: str, verify: bool = False, secret: Optional[str] = None
) -> Dict[str, Any]:
    """
    Décode un JWT (sans vérifier la signature par défaut)

    Args:
        token: Token JWT à décoder
        verify: Si True, vérifie la signature (nécessite secret)
        secret: Secret JWT (requis si verify=True)

    Returns:
        Dict contenant les claims du JWT
    """
    if not JWT_AVAILABLE:
        raise ImportError(
            "PyJWT n'est pas installé. Installez-le avec: pip install PyJWT"
        )

    if verify:
        if not secret:
            raise ValueError("Le secret est requis pour vérifier la signature")
        return pyjwt.decode(token, secret, algorithms=["HS256"])
    else:
        return pyjwt.decode(token, options={"verify_signature": False})


def print_result(result: Dict[str, Any]) -> None:
    """Affiche les résultats de manière formatée"""
    if result["success"]:
        print("\n✅ EXÉCUTION RÉUSSIE")
        print("=" * 50)
        print(f"Temps d'exécution: {result['execution_time_ms']}ms")
        print(f"Faits injectés: {result['results']['facts_count']}")
        print(f"Activations: {result['results']['activations_count']}")

        if result["results"]["activations_count"] > 0:
            print("\n🎯 Actions déclenchées:")
            for i, activation in enumerate(result["results"]["activations"], 1):
                print(f"\n{i}. {activation['action_name']}")
                if activation.get("arguments"):
                    print("   Arguments:")
                    for arg in activation["arguments"]:
                        print(
                            f"     [{arg['position']}] {arg['value']} ({arg['type']})"
                        )
    else:
        print("\n❌ ERREUR D'EXÉCUTION")
        print("=" * 50)
        print(f"Type: {result['error_type']}")
        print(f"Message: {result['error']}")
        print(f"Temps: {result['execution_time_ms']}ms")


def print_token_info(token: str) -> None:
    """Affiche les informations d'un JWT"""
    if not JWT_AVAILABLE:
        print("⚠️  PyJWT non installé - impossible de décoder le token")
        return

    try:
        claims = decode_jwt(token)
        print("\n📋 Informations du JWT:")
        print("-" * 50)
        print(f"Utilisateur: {claims.get('username', 'N/A')}")
        print(f"Rôles: {', '.join(claims.get('roles', [])) or 'Aucun'}")
        print(f"Émetteur: {claims.get('iss', 'N/A')}")

        # Dates
        if "iat" in claims:
            iat = datetime.fromtimestamp(claims["iat"])
            print(f"Émis le: {iat.strftime('%Y-%m-%d %H:%M:%S')}")

        if "exp" in claims:
            exp = datetime.fromtimestamp(claims["exp"])
            now = datetime.utcnow()
            print(f"Expire le: {exp.strftime('%Y-%m-%d %H:%M:%S')}")

            if exp > now:
                remaining = exp - now
                hours = remaining.total_seconds() / 3600
                print(f"Temps restant: {hours:.1f}h")
            else:
                print("⚠️  Token EXPIRÉ!")

    except Exception as e:
        print(f"⚠️  Erreur décodage JWT: {e}")


def example_basic_usage():
    """Exemple d'utilisation basique avec JWT"""
    print("📝 Exemple 1: Utilisation basique avec JWT")
    print("-" * 50)

    # Créer le client
    client = TSDJWTClient(server_url="http://localhost:8080")

    # Afficher les infos du token
    print_token_info(client.jwt_token)

    # Test de connexion
    print("\n🔍 Test de connexion...")
    health = client.health_check()
    print(f"✅ Serveur OK - Version: {health['version']}")
    print(f"⏱️  Uptime: {health['uptime_seconds']}s")

    # Exécuter un programme simple
    print("\n🚀 Exécution d'un programme TSD...")
    tsd_code = """
type User : <
  id: string,
  username: string,
  role: string
>

User("u1", "alice", "admin")
User("u2", "bob", "developer")
User("u3", "charlie", "user")
"""

    result = client.execute(tsd_code)
    print_result(result)


def example_with_generation():
    """Exemple avec génération de JWT en Python"""
    print("\n\n📝 Exemple 2: Génération de JWT en Python")
    print("-" * 50)

    if not JWT_AVAILABLE:
        print("⚠️  PyJWT non installé - cet exemple nécessite PyJWT")
        print("   Installez-le avec: pip install PyJWT")
        return

    secret = os.getenv("TSD_JWT_SECRET")
    if not secret:
        print("⚠️  Variable TSD_JWT_SECRET non définie")
        print("   Définissez-la avec: export TSD_JWT_SECRET='votre-secret'")
        return

    # Générer un JWT
    print("🔐 Génération d'un JWT...")
    token = generate_jwt(
        secret=secret,
        username="python_user",
        roles=["developer", "api"],
        expiration_hours=1,
    )

    print(f"Token généré: {token[:50]}...")
    print_token_info(token)

    # Utiliser le token
    print("\n🚀 Test avec le token généré...")
    client = TSDJWTClient(jwt_token=token)

    health = client.health_check()
    print(f"✅ Connexion OK - Version: {health['version']}")


def example_token_expiration():
    """Exemple de gestion de l'expiration"""
    print("\n\n📝 Exemple 3: Gestion de l'expiration")
    print("-" * 50)

    if not JWT_AVAILABLE:
        print("⚠️  PyJWT non installé")
        return

    secret = os.getenv("TSD_JWT_SECRET")
    if not secret:
        print("⚠️  Variable TSD_JWT_SECRET non définie")
        return

    # Générer un token avec expiration courte (1 seconde pour la démo)
    print("🔐 Génération d'un token expirant dans 1 seconde...")
    token = generate_jwt(
        secret=secret,
        username="short_lived",
        expiration_hours=1 / 3600,  # 1 seconde
    )

    client = TSDJWTClient(jwt_token=token)
    print_token_info(token)

    # Test immédiat (devrait fonctionner)
    print("\n🔍 Test immédiat...")
    try:
        health = client.health_check()
        print(f"✅ OK - {health['status']}")
    except Exception as e:
        print(f"❌ Erreur: {e}")

    # Attendre l'expiration
    print("\n⏳ Attente de 2 secondes...")
    import time

    time.sleep(2)

    # Test après expiration (devrait échouer)
    print("\n🔍 Test après expiration...")
    try:
        health = client.health_check()
        print(f"✅ OK - {health['status']}")
    except Exception as e:
        print(f"❌ Erreur attendue: Token expiré")

    # Régénérer un token
    print("\n🔄 Régénération d'un nouveau token...")
    new_token = generate_jwt(secret=secret, username="short_lived", expiration_hours=1)
    client.update_token(new_token)

    try:
        health = client.health_check()
        print(f"✅ OK avec nouveau token - {health['status']}")
    except Exception as e:
        print(f"❌ Erreur: {e}")


def example_multiple_users():
    """Exemple avec plusieurs utilisateurs"""
    print("\n\n📝 Exemple 4: Multi-utilisateurs")
    print("-" * 50)

    if not JWT_AVAILABLE:
        print("⚠️  PyJWT non installé")
        return

    secret = os.getenv("TSD_JWT_SECRET")
    if not secret:
        print("⚠️  Variable TSD_JWT_SECRET non définie")
        return

    users = [
        ("alice", ["admin", "developer"]),
        ("bob", ["developer"]),
        ("charlie", ["readonly"]),
    ]

    for username, roles in users:
        print(f"\n👤 Utilisateur: {username} (rôles: {', '.join(roles)})")

        # Générer un token pour cet utilisateur
        token = generate_jwt(secret=secret, username=username, roles=roles)
        client = TSDJWTClient(jwt_token=token)

        # Test
        try:
            result = client.execute(f'type User : <id: string>\nUser("{username}")')
            if result["success"]:
                print(f"   ✅ Exécution réussie")
            else:
                print(f"   ❌ {result['error']}")
        except Exception as e:
            print(f"   ❌ Erreur: {e}")


def main():
    """Point d'entrée principal"""
    parser = argparse.ArgumentParser(
        description="Exemple d'utilisation du client TSD avec JWT"
    )
    parser.add_argument(
        "--server",
        default="http://localhost:8080",
        help="URL du serveur TSD (défaut: http://localhost:8080)",
    )
    parser.add_argument("--token", help="Token JWT (défaut: variable TSD_AUTH_TOKEN)")
    parser.add_argument(
        "--generate",
        action="store_true",
        help="Générer un JWT en Python (nécessite TSD_JWT_SECRET)",
    )
    parser.add_argument(
        "--username", default="python_user", help="Nom d'utilisateur pour le JWT généré"
    )
    parser.add_argument(
        "--roles", help="Rôles séparés par des virgules (ex: admin,user)"
    )
    parser.add_argument(
        "--expiration",
        type=int,
        default=24,
        help="Durée de validité en heures (défaut: 24)",
    )
    parser.add_argument(
        "--example",
        type=int,
        choices=[1, 2, 3, 4],
        help="Numéro de l'exemple à exécuter (1-4, tous par défaut)",
    )
    parser.add_argument(
        "--decode", help="Décoder un JWT (affiche les claims sans vérifier)"
    )

    args = parser.parse_args()

    # Décoder un JWT
    if args.decode:
        if not JWT_AVAILABLE:
            print("❌ PyJWT non installé. Installez-le avec: pip install PyJWT")
            sys.exit(1)
        print_token_info(args.decode)
        return

    # Générer un JWT
    if args.generate:
        if not JWT_AVAILABLE:
            print("❌ PyJWT non installé. Installez-le avec: pip install PyJWT")
            sys.exit(1)

        secret = os.getenv("TSD_JWT_SECRET")
        if not secret:
            print("❌ Variable TSD_JWT_SECRET non définie")
            print("\nDéfinissez-la avec:")
            print("  export TSD_JWT_SECRET='votre-secret-32-chars-minimum'")
            print("\nOu utilisez:")
            print("  tsd-auth generate-jwt -secret 'votre-secret' -username alice")
            sys.exit(1)

        roles = args.roles.split(",") if args.roles else []
        token = generate_jwt(
            secret=secret,
            username=args.username,
            roles=roles,
            expiration_hours=args.expiration,
        )

        print("🎫 JWT généré:")
        print(token)
        print()
        print_token_info(token)
        print("\nUtilisez-le avec:")
        print(f'  export TSD_AUTH_TOKEN="{token}"')
        return

    # Mettre à jour les variables pour les exemples
    if args.token:
        os.environ["TSD_AUTH_TOKEN"] = args.token

    # Vérifier que le token est défini
    if not os.getenv("TSD_AUTH_TOKEN"):
        print("❌ Erreur: JWT token requis")
        print()
        print("Définissez-le via:")
        print("  export TSD_AUTH_TOKEN='eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...'")
        print("  ou")
        print("  python3 client_jwt.py --token 'votre-jwt'")
        print()
        print("Pour générer un JWT:")
        print("  # Avec tsd-auth (recommandé)")
        print("  tsd-auth generate-jwt -secret 'votre-secret' -username alice")
        print()
        print("  # En Python")
        print("  export TSD_JWT_SECRET='votre-secret'")
        print("  python3 client_jwt.py --generate --username alice")
        sys.exit(1)

    print("=" * 50)
    print("🔐 Client TSD avec JWT")
    print("=" * 50)
    print(f"Serveur: {args.server}")
    print(f"Token: {os.getenv('TSD_AUTH_TOKEN')[:30]}...")
    print()

    try:
        if args.example is None:
            # Exécuter tous les exemples
            example_basic_usage()
            example_with_generation()
            example_token_expiration()
            example_multiple_users()
        elif args.example == 1:
            example_basic_usage()
        elif args.example == 2:
            example_with_generation()
        elif args.example == 3:
            example_token_expiration()
        elif args.example == 4:
            example_multiple_users()

        print("\n\n" + "=" * 50)
        print("✅ Tous les exemples terminés!")
        print("=" * 50)

    except ValueError as e:
        print(f"\n❌ Erreur de configuration: {e}")
        sys.exit(1)
    except requests.exceptions.ConnectionError:
        print(f"\n❌ Erreur: Impossible de se connecter au serveur {args.server}")
        print("Vérifiez que le serveur TSD est démarré:")
        print("  tsd-server -auth jwt")
        sys.exit(1)
    except Exception as e:
        print(f"\n❌ Erreur: {e}")
        import traceback

        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()
