# Tutoriel pratique : Écriture de contraintes

Ce tutoriel vous guide pas à pas dans l'apprentissage de l'écriture de fichiers de contraintes avec des exercices pratiques.

## Préparation

Avant de commencer, assurez-vous d'avoir :
- Le fichier `constraint.peg` (grammaire de référence)
- Un parser généré (via `pigeon -o parser.go constraint.peg`)
- Les fichiers de test existants comme exemples

---

## Leçon 1 : Premiers pas avec les types

### Objectif
Apprendre à définir des types simples avec des champs typés.

### Théorie
Un type se définit avec la syntaxe :
```
type NomType : < champ1: type, champ2: type >
```

### Exercice 1.1 : Définir un type Person
Créez un type `Personne` avec les champs suivants :
- `prenom` (string)
- `nom` (string) 
- `age` (number)
- `actif` (bool)

<details>
<summary>💡 Solution</summary>

```
type Personne : < prenom: string, nom: string, age: number, actif: bool >
```
</details>

### Exercice 1.2 : Types pour un e-commerce
Créez trois types pour un système e-commerce :
1. `Client` : nom, email, premium (bool)
2. `Produit` : titre, prix, disponible (bool)  
3. `Commande` : numero (number), total (number), livree (bool)

<details>
<summary>💡 Solution</summary>

```
type Client : < nom: string, email: string, premium: bool >
type Produit : < titre: string, prix: number, disponible: bool >
type Commande : < numero: number, total: number, livree: bool >
```
</details>

---

## Leçon 2 : Règles simples

### Objectif
Créer des règles avec une seule variable et des contraintes basiques.

### Théorie
Une règle basique : `{ variable: Type } / contrainte`

### Exercice 2.1 : Clients majeurs
Avec le type `Personne` de l'exercice 1.1, écrivez une règle qui sélectionne les personnes de 18 ans ou plus.

<details>
<summary>💡 Solution</summary>

```
type Personne : < prenom: string, nom: string, age: number, actif: bool >

{ p: Personne } / p.age >= 18
```
</details>

### Exercice 2.2 : Produits en stock
Avec le type `Produit`, écrivez une règle pour les produits disponibles ET à moins de 50€.

<details>
<summary>💡 Solution</summary>

```
type Produit : < titre: string, prix: number, disponible: bool >

{ prod: Produit } / prod.disponible = true AND prod.prix < 50
```
</details>

### Exercice 2.3 : Contraintes complexes
Écrivez une règle pour sélectionner les personnes qui sont :
- Âgées entre 25 et 65 ans (inclus)
- ET actives
- ET dont le prénom n'est pas "Test"

<details>
<summary>💡 Solution</summary>

```
type Personne : < prenom: string, nom: string, age: number, actif: bool >

{ p: Personne } / p.age >= 25 AND p.age <= 65 AND p.actif = true AND p.prenom != "Test"
```
</details>

---

## Leçon 3 : Règles avec plusieurs variables

### Objectif
Utiliser plusieurs variables dans une même règle et créer des relations entre elles.

### Théorie
Syntaxe : `{ var1: Type1, var2: Type2 } / contraintes_utilisant_var1_et_var2`

### Exercice 3.1 : Commandes client premium
Écrivez une règle qui associe un client premium avec une commande de plus de 100€.

<details>
<summary>💡 Solution</summary>

```
type Client : < nom: string, email: string, premium: bool >
type Commande : < numero: number, total: number, livree: bool >

{ c: Client, cmd: Commande } / c.premium = true AND cmd.total > 100
```
</details>

### Exercice 3.2 : Comparaison entre personnes
Créez une règle qui trouve les paires de personnes où :
- La première est plus âgée que la seconde
- La différence d'âge est d'au moins 10 ans
- Les deux sont actives

<details>
<summary>💡 Solution</summary>

```
type Personne : < prenom: string, nom: string, age: number, actif: bool >

{ p1: Personne, p2: Personne } / p1.age > p2.age AND (p1.age - p2.age) >= 10 AND p1.actif = true AND p2.actif = true
```
</details>

---

## Leçon 4 : Actions et appels de fonctions

### Objectif
Ajouter des actions aux règles avec la syntaxe `==> fonction(args)`.

### Théorie
Les actions se déclenchent quand une règle est satisfaite :
```
{ vars } / contraintes ==> nomFonction(var1, var2)
```

### Exercice 4.1 : Validation client
Reprenez la règle des clients majeurs et ajoutez une action `validerClient(client)`.

<details>
<summary>💡 Solution</summary>

```
type Personne : < prenom: string, nom: string, age: number, actif: bool >

{ p: Personne } / p.age >= 18 ==> validerClient(p)
```
</details>

### Exercice 4.2 : Processus de commande
Créez une règle qui :
- Associe un client premium avec un produit disponible
- Le produit coûte moins de 200€
- Déclenche l'action `traiterCommande(client, produit)`

<details>
<summary>💡 Solution</summary>

```
type Client : < nom: string, email: string, premium: bool >
type Produit : < titre: string, prix: number, disponible: bool >

{ c: Client, p: Produit } / c.premium = true AND p.disponible = true AND p.prix < 200 ==> traiterCommande(c, p)
```
</details>

---

## Leçon 5 : Cas d'usage complexes

### Objectif
Combiner tout ce qui a été appris pour créer des systèmes de règles complets.

### Exercice 5.1 : Système de recommandation
Contexte : Un site de streaming vidéo

Créez les types et règles pour :
1. `Utilisateur` : nom, age, abonne (bool), genre_prefere (string)
2. `Film` : titre, duree (number), genre (string), note (number) 
3. `Visionnage` : termine (bool), note_utilisateur (number)

Règles à implémenter :
- Recommander des films du genre préféré de l'utilisateur avec une note >= 7
- Action : `recommander(utilisateur, film)`

<details>
<summary>💡 Solution</summary>

```
type Utilisateur : < nom: string, age: number, abonne: bool, genre_prefere: string >
type Film : < titre: string, duree: number, genre: string, note: number >

{ user: Utilisateur, film: Film } / 
    user.abonne = true AND 
    film.genre = user.genre_prefere AND 
    film.note >= 7 
==> recommander(user, film)
```
</details>

### Exercice 5.2 : Système de gestion RH
Contexte : Gestion des employés et des projets

Types nécessaires :
1. `Employe` : nom, experience (number), disponible (bool), competence (string)
2. `Projet` : nom, duree (number), competence_requise (string), urgent (bool)

Règles :
1. Affecter des employés disponibles aux projets urgents correspondant à leur compétence
2. Les employés doivent avoir au moins 2 ans d'expérience pour les projets urgents
3. Action : `affecterProjet(employe, projet)`

<details>
<summary>💡 Solution</summary>

```
type Employe : < nom: string, experience: number, disponible: bool, competence: string >
type Projet : < nom: string, duree: number, competence_requise: string, urgent: bool >

{ emp: Employe, proj: Projet } / 
    emp.disponible = true AND 
    proj.urgent = true AND 
    emp.competence = proj.competence_requise AND 
    emp.experience >= 2 
==> affecterProjet(emp, proj)
```
</details>

---

## Leçon 6 : Règles multiples et organisation

### Objectif
Structurer un fichier avec plusieurs règles et maintenir la lisibilité.

### Exercice 6.1 : E-commerce complet
Créez un système complet pour un e-commerce avec :

Types :
- `Client` : nom, age, vip (bool), budget (number)
- `Produit` : nom, prix (number), stock (number), categorie (string)
- `Reduction` : pourcentage (number), categorie_cible (string), actif (bool)

Règles :
1. Clients VIP peuvent acheter des produits "premium" en stock
2. Appliquer des réductions actives aux produits de la bonne catégorie si le client a le budget
3. Valider les commandes de clients majeurs avec budget suffisant

Actions : `vendreVIP(client, produit)`, `appliquerReduction(client, produit, reduction)`, `validerCommande(client, produit)`

<details>
<summary>💡 Solution</summary>

```
// ===== TYPES =====
type Client : < nom: string, age: number, vip: bool, budget: number >
type Produit : < nom: string, prix: number, stock: number, categorie: string >
type Reduction : < pourcentage: number, categorie_cible: string, actif: bool >

// ===== RÈGLES MÉTIER =====

// Vente VIP pour produits premium
{ c: Client, p: Produit } / 
    c.vip = true AND 
    p.categorie = "premium" AND 
    p.stock > 0 
==> vendreVIP(c, p)

// Application de réductions
{ c: Client, p: Produit, r: Reduction } / 
    r.actif = true AND 
    r.categorie_cible = p.categorie AND 
    c.budget >= (p.prix * (100 - r.pourcentage) / 100) 
==> appliquerReduction(c, p, r)

// Validation commandes standard
{ c: Client, p: Produit } / 
    c.age >= 18 AND 
    c.budget >= p.prix AND 
    p.stock > 0 
==> validerCommande(c, p)
```
</details>

---

## Exercices d'évaluation

### Défi 1 : Système de transport
Créez un système pour une compagnie de transport avec bus, chauffeurs et trajets.

**Contraintes métier :**
- Seuls les chauffeurs avec un permis valide peuvent conduire
- Les bus doivent être en état de marche
- Les trajets longue distance (> 500km) nécessitent deux chauffeurs
- Les chauffeurs seniors (> 10 ans d'expérience) peuvent faire des trajets de nuit

### Défi 2 : Plateforme de cours en ligne
Système de gestion d'une plateforme éducative avec étudiants, cours et certifications.

**Contraintes métier :**
- Les étudiants peuvent s'inscrire aux cours de leur niveau ou inférieur
- Les cours avancés nécessitent d'avoir complété les prérequis
- Seuls les étudiants ayant validé >= 80% peuvent obtenir la certification
- Les cours premium sont réservés aux abonnés payants

---

## Conseils pour la suite

### Débogage
1. **Testez étape par étape** : Commencez par des règles simples
2. **Vérifiez la syntaxe** : Attention aux deux-points, virgules, parenthèses
3. **Validez les types** : Assurez-vous que tous les types référencés sont définis
4. **Utilisez des noms explicites** : `client` plutôt que `c` pour la lisibilité

### Optimisation
1. **Ordonnez les contraintes** : Mettez les plus sélectives en premier
2. **Groupez les règles** par domaine métier
3. **Commentez** les règles complexes
4. **Évitez la duplication** de contraintes

### Évolution
1. **Commencez simple** puis enrichissez progressivement
2. **Testez chaque modification** avec des données réelles
3. **Documentez** les règles métier importantes
4. **Versionnez** vos fichiers de contraintes

---

*Ce tutoriel couvre les bases nécessaires pour maîtriser l'écriture de contraintes. Pratiquez régulièrement et n'hésitez pas à consulter les fichiers de test pour des exemples supplémentaires.*