// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package rete

import (
	"fmt"
	"strings"
)

// BindingChain représente une chaîne immuable de bindings variable → fact.
//
// La structure utilise le pattern "Cons list" (liste chaînée fonctionnelle)
// pour permettre le partage structurel entre différents tokens, tout en
// garantissant l'immutabilité complète.
//
// Propriétés garanties (invariants):
//   - Une fois créée, une BindingChain ne change JAMAIS
//   - Add() retourne une NOUVELLE chaîne, ne modifie pas l'existante
//   - La racine (chaîne vide) est représentée par Parent == nil
//   - Pas de cycles : Parent pointe toujours vers une chaîne plus courte
//   - Thread-safe grâce à l'immutabilité
//
// Structure:
//   - Variable: nom de la variable (ex: "u", "order", "task")
//   - Fact: pointeur vers le fait lié à cette variable
//   - Parent: chaîne parente (nil si racine/vide)
//
// Exemple d'utilisation:
//
//	// Créer une chaîne vide
//	chain := NewBindingChain()
//
//	// Ajouter des bindings (chaque Add retourne une nouvelle chaîne)
//	chain1 := chain.Add("u", userFact)        // chain1: [u]
//	chain2 := chain1.Add("order", orderFact)  // chain2: [u, order]
//	chain3 := chain2.Add("task", taskFact)    // chain3: [u, order, task]
//
//	// L'ancienne chaîne reste inchangée
//	fmt.Println(chain1.Len())  // 1 (toujours [u])
//	fmt.Println(chain2.Len())  // 2 (toujours [u, order])
//	fmt.Println(chain3.Len())  // 3 ([u, order, task])
//
// Partage structurel:
//
//	chain1: [u] → nil
//	chain2: [order] → chain1 → nil
//	chain3: [task] → chain2 → chain1 → nil
//
// Les chaînes chain1 et chain2 sont partagées (réutilisées) dans chain3,
// économisant de la mémoire et permettant une construction efficace.
type BindingChain struct {
	Variable string        // Nom de la variable (ex: "u", "order", "task")
	Fact     *Fact         // Fait lié à cette variable
	Parent   *BindingChain // Chaîne parente (nil si racine/vide)
}

// NewBindingChain crée une chaîne de bindings vide.
//
// La chaîne vide est représentée par un pointeur nil.
// Cette représentation permet une optimisation mémoire et simplifie les algorithmes.
//
// Retourne:
//   - *BindingChain: pointeur nil représentant une chaîne vide
//
// Exemple:
//
//	empty := NewBindingChain()
//	fmt.Println(empty == nil)     // true
//	fmt.Println(empty.Len())      // 0 (grâce à la gestion du nil)
func NewBindingChain() *BindingChain {
	return nil // La chaîne vide est représentée par nil
}

// NewBindingChainWith crée une chaîne de bindings avec un binding initial.
//
// Cette fonction est un raccourci pour créer une chaîne non-vide en une seule étape.
//
// Paramètres:
//   - variable: nom de la variable
//   - fact: pointeur vers le fait à lier
//
// Retourne:
//   - *BindingChain: nouvelle chaîne contenant un seul binding
//
// Exemple:
//
//	chain := NewBindingChainWith("u", userFact)
//	fmt.Println(chain.Len())  // 1
//	fmt.Println(chain.Get("u") == userFact)  // true
func NewBindingChainWith(variable string, fact *Fact) *BindingChain {
	return &BindingChain{
		Variable: variable,
		Fact:     fact,
		Parent:   nil,
	}
}

// Add ajoute un binding à la chaîne et retourne une NOUVELLE chaîne.
//
// ⚠️ IMPORTANT: Cette méthode NE MODIFIE PAS la chaîne existante.
// Elle crée et retourne une nouvelle chaîne qui pointe vers l'ancienne comme parent.
//
// Si la variable existe déjà dans la chaîne, la nouvelle valeur masque l'ancienne
// (shadowing), mais l'ancienne reste accessible via la chaîne parent.
//
// Complexité: O(1) - Création d'un seul nœud
//
// Paramètres:
//   - variable: nom de la variable à ajouter
//   - fact: pointeur vers le fait à lier
//
// Retourne:
//   - *BindingChain: nouvelle chaîne contenant le binding ajouté
//
// Exemple:
//
//	chain1 := NewBindingChain()
//	chain2 := chain1.Add("u", userFact)
//
//	// chain1 est toujours vide (nil)
//	fmt.Println(chain1 == nil)  // true
//
//	// chain2 contient le binding
//	fmt.Println(chain2.Len())  // 1
//	fmt.Println(chain2.Get("u") == userFact)  // true
//
// Shadowing (écrasement avec préservation):
//
//	chain1 := NewBindingChainWith("u", fact1)
//	chain2 := chain1.Add("u", fact2)  // Nouvelle valeur pour "u"
//
//	fmt.Println(chain1.Get("u") == fact1)  // true (inchangé)
//	fmt.Println(chain2.Get("u") == fact2)  // true (nouvelle valeur)
func (bc *BindingChain) Add(variable string, fact *Fact) *BindingChain {
	return &BindingChain{
		Variable: variable,
		Fact:     fact,
		Parent:   bc, // L'ancienne chaîne devient le parent
	}
}

// Get retourne le fait associé à une variable, ou nil si non trouvé.
//
// La recherche parcourt la chaîne depuis la tête vers la racine,
// retournant le premier binding trouvé (ce qui permet le shadowing).
//
// Complexité: O(n) où n est le nombre de bindings dans la chaîne
// Note: Pour les cas d'usage typiques (n < 10), c'est acceptable et simple
//
// Paramètres:
//   - variable: nom de la variable à rechercher
//
// Retourne:
//   - *Fact: pointeur vers le fait si trouvé, nil sinon
//
// Exemple:
//
//	chain := NewBindingChain()
//	chain = chain.Add("u", userFact)
//	chain = chain.Add("order", orderFact)
//
//	fmt.Println(chain.Get("u") == userFact)      // true
//	fmt.Println(chain.Get("order") == orderFact) // true
//	fmt.Println(chain.Get("task") == nil)        // true (non trouvé)
func (bc *BindingChain) Get(variable string) *Fact {
	// Parcourir la chaîne depuis la tête
	current := bc
	for current != nil {
		if current.Variable == variable {
			return current.Fact
		}
		current = current.Parent
	}
	return nil
}

// Has vérifie si une variable existe dans la chaîne.
//
// Équivalent à Get(variable) != nil, mais plus expressif.
//
// Complexité: O(n) où n est le nombre de bindings
//
// Paramètres:
//   - variable: nom de la variable à vérifier
//
// Retourne:
//   - bool: true si la variable existe, false sinon
//
// Exemple:
//
//	chain := NewBindingChainWith("u", userFact)
//	fmt.Println(chain.Has("u"))    // true
//	fmt.Println(chain.Has("task")) // false
func (bc *BindingChain) Has(variable string) bool {
	return bc.Get(variable) != nil
}

// Len retourne le nombre de bindings dans la chaîne.
//
// Complexité: O(n) - Nécessite un parcours complet
// Note: Pourrait être optimisé avec un cache si nécessaire
//
// Retourne:
//   - int: nombre de bindings (0 pour une chaîne vide)
//
// Exemple:
//
//	empty := NewBindingChain()
//	fmt.Println(empty.Len())  // 0
//
//	chain := empty.Add("u", userFact).Add("order", orderFact)
//	fmt.Println(chain.Len())  // 2
func (bc *BindingChain) Len() int {
	count := 0
	current := bc
	for current != nil {
		count++
		current = current.Parent
	}
	return count
}

// Variables retourne la liste des noms de variables dans l'ordre d'ajout.
//
// L'ordre retourné est du plus ancien (racine) au plus récent (tête).
// Les variables sont dédupliquées (en cas de shadowing, seule la première occurrence compte).
//
// Complexité: O(n) temps, O(n) espace
//
// Retourne:
//   - []string: slice des noms de variables
//
// Exemple:
//
//	chain := NewBindingChain()
//	chain = chain.Add("u", userFact)
//	chain = chain.Add("order", orderFact)
//	chain = chain.Add("task", taskFact)
//
//	vars := chain.Variables()
//	fmt.Println(vars)  // ["u", "order", "task"]
func (bc *BindingChain) Variables() []string {
	if bc == nil {
		return []string{}
	}

	// Collecter les variables en parcourant vers la racine
	vars := make([]string, 0, bc.Len())
	seen := make(map[string]bool)

	current := bc
	for current != nil {
		// Éviter les doublons (shadowing)
		if !seen[current.Variable] {
			vars = append(vars, current.Variable)
			seen[current.Variable] = true
		}
		current = current.Parent
	}

	// Inverser pour obtenir l'ordre d'ajout (racine → tête)
	for i, j := 0, len(vars)-1; i < j; i, j = i+1, j-1 {
		vars[i], vars[j] = vars[j], vars[i]
	}

	return vars
}

// ToMap convertit la chaîne en map pour compatibilité et debug.
//
// ⚠️ ATTENTION: Cette méthode crée une copie mutable.
// Ne pas l'utiliser pour modifier les bindings, uniquement pour lecture/debug.
//
// En cas de shadowing, seule la valeur la plus récente est conservée.
//
// Complexité: O(n) temps, O(n) espace
//
// Retourne:
//   - map[string]*Fact: map des bindings (copie)
//
// Exemple:
//
//	chain := NewBindingChain()
//	chain = chain.Add("u", userFact)
//	chain = chain.Add("order", orderFact)
//
//	m := chain.ToMap()
//	fmt.Println(m["u"] == userFact)      // true
//	fmt.Println(m["order"] == orderFact) // true
//	fmt.Println(len(m))                  // 2
func (bc *BindingChain) ToMap() map[string]*Fact {
	result := make(map[string]*Fact)
	if bc == nil {
		return result
	}

	// Parcourir depuis la racine pour avoir l'ordre correct
	// (les valeurs plus récentes écraseront les anciennes)
	vars := bc.Variables()
	for _, v := range vars {
		result[v] = bc.Get(v)
	}

	return result
}

// Merge combine deux chaînes de bindings.
//
// Crée une nouvelle chaîne contenant tous les bindings des deux chaînes.
// En cas de conflit (même variable dans les deux chaînes), la valeur de 'other' est prioritaire.
//
// Complexité: O(m) où m est le nombre de bindings dans 'other'
// (chaque binding de 'other' est ajouté à la chaîne résultante)
//
// Paramètres:
//   - other: chaîne à fusionner avec la chaîne actuelle
//
// Retourne:
//   - *BindingChain: nouvelle chaîne fusionnée
//
// Exemple:
//
//	chain1 := NewBindingChain().Add("u", userFact)
//	chain2 := NewBindingChain().Add("order", orderFact)
//
//	merged := chain1.Merge(chain2)
//	fmt.Println(merged.Has("u"))     // true
//	fmt.Println(merged.Has("order")) // true
//	fmt.Println(merged.Len())        // 2
//
// Gestion des conflits:
//
//	chain1 := NewBindingChain().Add("u", fact1)
//	chain2 := NewBindingChain().Add("u", fact2)
//
//	merged := chain1.Merge(chain2)
//	fmt.Println(merged.Get("u") == fact2)  // true (priorité à 'other')
func (bc *BindingChain) Merge(other *BindingChain) *BindingChain {
	// Commencer avec la chaîne actuelle
	result := bc

	// Ajouter tous les bindings de 'other' dans l'ordre
	// (les variables de 'other' sont ajoutées en dernier, donc prioritaires)
	if other != nil {
		vars := other.Variables()
		for _, v := range vars {
			fact := other.Get(v)
			result = result.Add(v, fact)
		}
	}

	return result
}

// String retourne une représentation textuelle pour debug.
//
// Format: "BindingChain{var1:FactID1, var2:FactID2, ...}"
//
// Exemple:
//
//	chain := NewBindingChain()
//	chain = chain.Add("u", userFact)  // userFact.ID = "U001"
//	chain = chain.Add("order", orderFact)  // orderFact.ID = "O001"
//
//	fmt.Println(chain.String())
//	// Output: "BindingChain{u:U001, order:O001}"
func (bc *BindingChain) String() string {
	if bc == nil {
		return "BindingChain{}"
	}

	vars := bc.Variables()
	if len(vars) == 0 {
		return "BindingChain{}"
	}

	parts := make([]string, 0, len(vars))
	for _, v := range vars {
		fact := bc.Get(v)
		if fact != nil {
			parts = append(parts, fmt.Sprintf("%s:%s", v, fact.ID))
		}
	}

	return fmt.Sprintf("BindingChain{%s}", strings.Join(parts, ", "))
}

// Chain retourne la liste des variables depuis la racine (pour traçage).
//
// Contrairement à Variables(), cette méthode retourne toutes les occurrences,
// y compris les variables shadowed. Utile pour le debugging et le traçage.
//
// Complexité: O(n)
//
// Retourne:
//   - []string: liste complète des variables dans l'ordre racine → tête
//
// Exemple:
//
//	chain := NewBindingChain()
//	chain = chain.Add("u", fact1)
//	chain = chain.Add("order", fact2)
//	chain = chain.Add("u", fact3)  // Shadowing
//
//	fmt.Println(chain.Variables())  // ["u", "order"] (dédupliqué)
//	fmt.Println(chain.Chain())      // ["u", "order", "u"] (complet)
func (bc *BindingChain) Chain() []string {
	if bc == nil {
		return []string{}
	}

	// Collecter toutes les variables sans déduplication
	vars := make([]string, 0, bc.Len())
	current := bc
	for current != nil {
		vars = append(vars, current.Variable)
		current = current.Parent
	}

	// Inverser pour obtenir l'ordre racine → tête
	for i, j := 0, len(vars)-1; i < j; i, j = i+1, j-1 {
		vars[i], vars[j] = vars[j], vars[i]
	}

	return vars
}
