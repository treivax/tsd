package main

import (
	"context"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
	// Configuration de la connexion etcd
	config := clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, // Adresses des serveurs etcd
		DialTimeout: 5 * time.Second,
	}

	// Créer le client etcd
	client, err := clientv3.New(config)
	if err != nil {
		log.Fatalf("Erreur lors de la création du client etcd: %v", err)
	}
	defer client.Close()

	// Tester la connexion
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = client.Status(ctx, config.Endpoints[0])
	if err != nil {
		log.Fatalf("Impossible de se connecter à etcd: %v", err)
	}

	fmt.Println("Connexion réussie à etcd!")
	fmt.Println("Récupération des clés avec le préfixe '/a/b/c'...")

	// Lister les clés avec le préfixe '/a/b/c'
	prefix := "/a/b/c*toto"

	// Créer un contexte avec timeout pour la requête
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Effectuer la requête Get avec le préfixe
	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		log.Fatalf("Erreur lors de la récupération des clés: %v", err)
	}

	// Afficher les résultats
	fmt.Printf("\nNombre de clés trouvées avec le préfixe '%s': %d\n", prefix, resp.Count)

	if resp.Count == 0 {
		fmt.Println("Aucune clé trouvée avec ce préfixe.")
	} else {
		fmt.Println("\nClés trouvées:")
		fmt.Println("==============")

		for i, kv := range resp.Kvs {
			fmt.Printf("%d. Clé: %s\n", i+1, string(kv.Key))
			fmt.Printf("   Valeur: %s\n", string(kv.Value))
			fmt.Printf("   Version: %d\n", kv.Version)
			fmt.Printf("   Créée à: %d\n", kv.CreateRevision)
			fmt.Printf("   Modifiée à: %d\n", kv.ModRevision)
			fmt.Println("   ---")
		}
	}

	// Afficher quelques statistiques supplémentaires
	fmt.Printf("\nInformations sur la réponse:\n")
	fmt.Printf("- Révision du cluster: %d\n", resp.Header.Revision)
	fmt.Printf("- ID du membre: %d\n", resp.Header.MemberId)
	fmt.Printf("- ID du cluster: %d\n", resp.Header.ClusterId)

	// Appel de la fonction take pour l'état 'init' avec le client etcd
	fmt.Printf("\n=== Opération take('init') ===\n")
	key, value := take(client, "init")
	fmt.Printf("Résultat de take('init'):\n")
	fmt.Printf("- Clé: %s\n", key)
	fmt.Printf("- Valeur: %s\n", value)
}
