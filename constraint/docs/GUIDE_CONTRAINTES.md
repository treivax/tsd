# Guide de définition des types et règles de contraintes

Ce document présente la syntaxe et les conventions pour écrire des fichiers de définition de types et de règles de contraintes selon la grammaire PEG fournie dans `constraint.peg`.

## Table des matières

1. [Concepts fondamentaux](#concepts-fondamentaux)
2. [Définition des types](#définition-des-types)
3. [Écriture des règles de contraintes](#écriture-des-règles-de-contraintes)
4. [Actions et appels de fonctions](#actions-et-appels-de-fonctions)
5. [Exemples complets](#exemples-complets)
6. [Bonnes pratiques](#bonnes-pratiques)
7. [Référence syntaxique](#référence-syntaxique)

---

## Concepts fondamentaux

### Terminologie

- **Type** : Structure de données définie avec des champs typés
- **Fait** : Instance/objet concret d'un type défini
- **Variable typée** : Variable liée à un type spécifique dans une règle
- **Contrainte** : Condition logique appliquée aux variables
- **Action** : Fonction déclenchée quand une règle est satisfaite

### Structure générale d'un fichier

Un fichier de définition contient :
1. **Définitions de types** (optionnel)
2. **Règles de contraintes** avec leurs actions (optionnel)

```
type TypeName : < field1: type, field2: type, ... >

{ var1: Type1, var2: Type2 } / constraints ==> action(args)
```

---

## Définition des types

### Syntaxe de base

```
type NomDuType : < champ1: type_atomique, champ2: type_atomique, ... >
```

### Types atomiques supportés

- `string` : Chaîne de caractères
- `number` : Nombre (entier ou décimal)
- `bool` : Booléen (true/false)

### Exemples de définitions

```
// Type simple
type Personne : < nom: string, age: number, adulte: bool >

// Type pour un animal
type Animal : < espece: string, domestique: bool, age: number >

// Type pour un produit
type Produit : < prix: number, stock: number, categorie: string >
```

### Règles de nommage

- Les noms de types commencent par une majuscule (convention)
- Les noms de champs commencent par une minuscule
- Utilisez des caractères alphanumériques et underscore `_`
- Pas d'espaces dans les noms

---

## Écriture des règles de contraintes

### Structure d'une règle

```
{ variables_typées } / contraintes [==> action]
```

- **Variables typées** : Variables liées à des types
- **Contraintes** : Conditions logiques à satisfaire
- **Action** : Fonction à exécuter (optionnelle)

### Variables typées

Définissez les variables avec leur type dans des accolades :

```
{ client: Personne, pet: Animal }
{ p1: Personne, p2: Personne }
{ prod: Produit }
```

### Contraintes

#### Accès aux champs

Utilisez la notation pointée pour accéder aux champs :

```
client.age        // Accès au champ 'age' de la variable 'client'
pet.domestique    // Accès au champ 'domestique' de la variable 'pet'
```

#### Opérateurs de comparaison

| Opérateur | Description | Exemple |
|-----------|-------------|---------|
| `==` | Égalité | `client.age == 25` |
| `!=` | Différent | `client.nom != "inconnu"` |
| `<` | Inférieur | `client.age < 18` |
| `>` | Supérieur | `client.age > 65` |
| `<=` | Inférieur ou égal | `produit.prix <= 100` |
| `>=` | Supérieur ou égal | `client.age >= 18` |
| `=` | Égalité (alternative) | `client.adulte = true` |

#### Opérateurs logiques

| Opérateur | Description | Exemple |
|-----------|-------------|---------|
| `AND` ou `&&` | Et logique | `client.adulte = true AND pet.domestique = true` |
| `OR` ou `\|\|` | Ou logique | `client.age < 18 OR client.age > 65` |

#### Parenthèses

Utilisez des parenthèses pour grouper les conditions :

```
(client.age > 18 AND client.adulte = true) OR client.nom = "admin"
```

#### Expressions arithmétiques

Supports des opérations arithmétiques de base :

```
client.age + 5 > 30
produit.prix * 1.2 <= budget
(stock1 + stock2) / 2 > seuil_minimum
```

### Littéraux

#### Nombres
```
42        // Entier
3.14      // Décimal
```

#### Chaînes de caractères
```
"Hello"     // Guillemets doubles
'World'     // Guillemets simples
```

#### Booléens
```
true
false
```

---

## Actions et appels de fonctions

### Syntaxe des actions

```
==> nomFonction(argument1, argument2, ...)
```

### Arguments

Les arguments sont des noms de variables définies dans la partie "variables typées" :

```
{ client: Personne, pet: Animal } / ... ==> adoption(client)
{ acheteur: Personne, vehicule: Voiture } / ... ==> vente(acheteur, vehicule)
```

---

## Exemples complets

### Exemple 1 : Règle simple avec action

```
type Personne : < nom: string, age: number, adulte: bool >

{ client: Personne } / client.age >= 18 AND client.adulte = true ==> validerClient(client)
```

### Exemple 2 : Règles multiples

```
type Personne : < nom: string, age: number, adulte: bool >
type Animal : < espece: string, domestique: bool, age: number >
type Voiture : < marque: string, annee: number, electrique: bool >

// Règle d'adoption
{ client: Personne, pet: Animal } / client.adulte = true AND pet.domestique = true ==> adoption(client)

// Règle de vente de véhicule
{ acheteur: Personne, vehicule: Voiture } / acheteur.age >= 18 AND vehicule.electrique = true ==> vente(acheteur, vehicule)
```

### Exemple 3 : Contraintes complexes

```
type Produit : < prix: number, stock: number, categorie: string >
type Client : < nom: string, age: number, budget: number, vip: bool >

{ prod: Produit, client: Client } / 
    prod.prix <= client.budget AND 
    prod.stock > 0 AND 
    (client.vip = true OR prod.categorie != "premium") AND
    client.age >= 18 
==> processCommande(client, prod)
```

### Exemple 4 : Comparaison entre variables

```
type Personne : < nom: string, age: number, salaire: number >

// Trouver les personnes avec un salaire supérieur à une autre
{ p1: Personne, p2: Personne } / 
    p1.salaire > p2.salaire AND 
    p1.age < p2.age 
==> analyseEcartSalarial(p1, p2)
```

---

## Bonnes pratiques

### Nommage

1. **Types** : Utilisez PascalCase (`Personne`, `AnimalDomestique`)
2. **Champs** : Utilisez camelCase (`nom`, `estAdulte`)
3. **Variables** : Utilisez des noms descriptifs (`client`, `produit`)

### Structure

1. **Définissez tous les types en début de fichier**
2. **Groupez les règles par domaine métier**
3. **Ajoutez des commentaires pour expliquer les règles complexes**

### Performance

1. **Ordonnez les contraintes** : mettez les plus sélectives en premier
2. **Évitez les contraintes redondantes**
3. **Utilisez des parenthèses** pour clarifier la précédence

### Exemple bien structuré

```
// ===== DÉFINITIONS DE TYPES =====

type Client : < nom: string, age: number, vip: bool >
type Produit : < nom: string, prix: number, stock: number >
type Commande : < montant: number, urgent: bool >

// ===== RÈGLES MÉTIER =====

// Règle de validation client VIP
{ client: Client } / 
    client.age >= 18 AND client.vip = true 
==> activerServiceVIP(client)

// Règle de gestion stock
{ prod: Produit, cmd: Commande } / 
    prod.stock > 0 AND 
    cmd.montant >= prod.prix 
==> traiterCommande(cmd, prod)
```

---

## Référence syntaxique

### Grammaire complète (résumé)

```
// Définition de type
type <nom> : < <champ>: <type>, ... >

// Règle
{ <var>: <type>, ... } / <contraintes> [==> <action>(<args>)]

// Types atomiques
string | number | bool

// Opérateurs de comparaison
== | != | < | > | <= | >= | =

// Opérateurs logiques
AND | OR | && | || | & | |

// Arithmétique
+ | - | * | /

// Accès aux champs
<variable>.<champ>

// Littéraux
"string" | 'string' | number | true | false
```

### Commentaires

```
// Commentaire sur une ligne
/* Commentaire 
   sur plusieurs lignes */
```

---

## Débogage et erreurs courantes

### Erreurs de syntaxe fréquentes

1. **Oubli de deux-points** dans la définition de type
   ```
   // ❌ Incorrect
   type Personne < nom: string >
   
   // ✅ Correct  
   type Personne : < nom: string >
   ```

2. **Mauvais séparateur** dans les listes
   ```
   // ❌ Incorrect
   { client: Personne; prod: Produit }
   
   // ✅ Correct
   { client: Personne, prod: Produit }
   ```

3. **Type non défini**
   ```
   // ❌ Référence à un type non défini
   { user: Utilisateur } / user.age > 18
   
   // ✅ Défini d'abord le type
   type Utilisateur : < age: number >
   { user: Utilisateur } / user.age > 18
   ```

### Conseils de débogage

1. **Vérifiez la correspondance des types** entre définition et utilisation
2. **Testez avec des règles simples** avant d'ajouter de la complexité
3. **Utilisez des parenthèses** pour éviter les ambiguïtés
4. **Vérifiez l'orthographe** des noms de champs et variables

---

*Cette documentation couvre la syntaxe complète supportée par la grammaire PEG définie dans `constraint.peg`. Pour des cas d'usage spécifiques ou des questions, référez-vous aux fichiers de test fournis avec le projet.*