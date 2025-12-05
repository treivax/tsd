#!/usr/bin/env python3
"""
Exemple d'utilisation du serveur TSD avec authentification par clé API (Auth Key)

Ce script montre comment :
1. Se connecter au serveur TSD avec une clé API
2. Vérifier la santé du serveur
3. Exécuter des programmes TSD
4. Gérer les erreurs d'authentification

Prérequis:
    pip install requests

Usage:
    # Définir la clé API
    export TSD_AUTH_TOKEN="votre-cle-api-ici"

    # Exécuter le script
    python3 client_auth_key.py

    # Ou passer la clé directement
    python3 client_auth_key.py --token "votre-cle-api"
"""

import argparse
import json
import os
import sys
from typing import Any, Dict, Optional

import requests


class TSDAuthKeyClient:
    """Client Python pour le serveur TSD avec authentification par clé API"""

    def __init__(
        self,
        server_url: str = "http://localhost:8080",
        auth_token: Optional[str] = None,
    ):
        """
        Initialise le client TSD

        Args:
            server_url: URL du serveur TSD
            auth_token: Token d'authentification (clé API)
                       Si None, utilise la variable d'environnement TSD_AUTH_TOKEN
        """
        self.server_url = server_url.rstrip("/")
        self.auth_token = auth_token or os.getenv("TSD_AUTH_TOKEN")

        if not self.auth_token:
            raise ValueError(
                "Token d'authentification requis.\n"
                "Définissez-le via:\n"
                "  - Le paramètre auth_token\n"
                "  - La variable d'environnement TSD_AUTH_TOKEN\n"
                "  - L'option --token en ligne de commande"
            )

        self.session = requests.Session()
        self.session.headers.update(self._get_headers())

    def _get_headers(self) -> Dict[str, str]:
        """Retourne les headers HTTP avec authentification"""
        return {
            "Content-Type": "application/json",
            "Authorization": f"Bearer {self.auth_token}",
        }

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
                    "Authentification échouée - Clé API invalide.\n"
                    "Vérifiez que votre clé API est correcte et que le serveur "
                    "est configuré avec la même clé."
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
                    "Authentification échouée - Clé API invalide ou expirée"
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


def example_basic_usage():
    """Exemple d'utilisation basique"""
    print("📝 Exemple 1: Utilisation basique")
    print("-" * 50)

    # Créer le client
    client = TSDAuthKeyClient(server_url="http://localhost:8080")

    # Test de connexion
    print("🔍 Test de connexion...")
    health = client.health_check()
    print(f"✅ Serveur OK - Version: {health['version']}")
    print(f"⏱️  Uptime: {health['uptime_seconds']}s")
    print()

    # Exécuter un programme simple
    print("🚀 Exécution d'un programme TSD...")
    tsd_code = """
type Person : <
  id: string,
  name: string,
  age: int
>

Person("p1", "Alice", 30)
Person("p2", "Bob", 25)
Person("p3", "Charlie", 35)
"""

    result = client.execute(tsd_code)
    print_result(result)


def example_file_execution():
    """Exemple d'exécution d'un fichier"""
    print("\n\n📝 Exemple 2: Exécution d'un fichier")
    print("-" * 50)

    # Créer un fichier TSD temporaire
    import tempfile

    with tempfile.NamedTemporaryFile(mode="w", suffix=".tsd", delete=False) as f:
        f.write("""
type Product : <
  id: string,
  name: string,
  price: float
>

Product("p1", "Laptop", 999.99)
Product("p2", "Mouse", 29.99)
""")
        temp_file = f.name

    try:
        client = TSDAuthKeyClient()
        result = client.execute_file(temp_file, verbose=True)
        print_result(result)
    finally:
        os.unlink(temp_file)


def example_error_handling():
    """Exemple de gestion d'erreurs"""
    print("\n\n📝 Exemple 3: Gestion d'erreurs")
    print("-" * 50)

    client = TSDAuthKeyClient()

    # Code TSD invalide
    print("🔍 Test avec du code invalide...")
    invalid_code = "this is not valid TSD code"

    result = client.execute(invalid_code)
    print_result(result)


def example_multiple_requests():
    """Exemple de requêtes multiples"""
    print("\n\n📝 Exemple 4: Requêtes multiples")
    print("-" * 50)

    client = TSDAuthKeyClient()

    programs = [
        ("Programme 1", 'type User : <id: string>\nUser("u1")'),
        ("Programme 2", 'type Order : <id: string, total: float>\nOrder("o1", 99.99)'),
        (
            "Programme 3",
            'type Event : <id: string, name: string>\nEvent("e1", "Login")',
        ),
    ]

    for name, code in programs:
        print(f"\n🚀 Exécution: {name}")
        result = client.execute(code, source_name=name)

        if result["success"]:
            print(f"   ✅ {result['results']['facts_count']} faits créés")
        else:
            print(f"   ❌ {result['error']}")


def main():
    """Point d'entrée principal"""
    parser = argparse.ArgumentParser(
        description="Exemple d'utilisation du client TSD avec Auth Key"
    )
    parser.add_argument(
        "--server",
        default="http://localhost:8080",
        help="URL du serveur TSD (défaut: http://localhost:8080)",
    )
    parser.add_argument(
        "--token", help="Token d'authentification (défaut: variable TSD_AUTH_TOKEN)"
    )
    parser.add_argument(
        "--example",
        type=int,
        choices=[1, 2, 3, 4],
        help="Numéro de l'exemple à exécuter (1-4, tous par défaut)",
    )

    args = parser.parse_args()

    # Mettre à jour les variables globales pour les exemples
    if args.token:
        os.environ["TSD_AUTH_TOKEN"] = args.token

    # Vérifier que le token est défini
    if not os.getenv("TSD_AUTH_TOKEN"):
        print("❌ Erreur: Token d'authentification requis")
        print()
        print("Définissez-le via:")
        print("  export TSD_AUTH_TOKEN='votre-cle-api'")
        print("  ou")
        print("  python3 client_auth_key.py --token 'votre-cle-api'")
        print()
        print("Pour générer une clé API:")
        print("  tsd-auth generate-key")
        sys.exit(1)

    print("=" * 50)
    print("🔐 Client TSD avec Auth Key")
    print("=" * 50)
    print(f"Serveur: {args.server}")
    print(f"Token: {os.getenv('TSD_AUTH_TOKEN')[:20]}...")
    print()

    try:
        if args.example is None:
            # Exécuter tous les exemples
            example_basic_usage()
            example_file_execution()
            example_error_handling()
            example_multiple_requests()
        elif args.example == 1:
            example_basic_usage()
        elif args.example == 2:
            example_file_execution()
        elif args.example == 3:
            example_error_handling()
        elif args.example == 4:
            example_multiple_requests()

        print("\n\n" + "=" * 50)
        print("✅ Tous les exemples terminés avec succès!")
        print("=" * 50)

    except ValueError as e:
        print(f"\n❌ Erreur de configuration: {e}")
        sys.exit(1)
    except requests.exceptions.ConnectionError:
        print(f"\n❌ Erreur: Impossible de se connecter au serveur {args.server}")
        print("Vérifiez que le serveur TSD est démarré:")
        print("  tsd-server -auth key")
        sys.exit(1)
    except Exception as e:
        print(f"\n❌ Erreur: {e}")
        import traceback

        traceback.print_exc()
        sys.exit(1)


if __name__ == "__main__":
    main()
