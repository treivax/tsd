# Guide Utilisateur - Système de Types

## Introduction

Le **système de types** de TSD v2.0 permet de définir des types personnalisés avec des champs pouvant être de types primitifs **ou de types de faits**, créant ainsi des relations explicites et type-safe entre entités.

---

## Types Primitifs

Les types primitifs disponibles :

| Type | Description | Exemple |
|------|-------------|---------|
| `string` | Chaîne de caractères | `"hello"` |
| `number` | Nombre (entier ou décimal) | `42`, `3.14` |
| `boolean` | Booléen | `true`, `false` |

---

## Types de Faits dans les Champs

**Nouveau en v2.0** : Les champs peuvent être d'un type de fait, créant des relations directes.

### Syntaxe

```tsd
type TypeName(field1: OtherFactType, field2: string, ...)
```

### Exemple

```tsd
type User(#username: string, email: string)
type Order(customer: User, #orderNum: string, total: number)
```

Le champ `customer` est de type `User` (un type de fait).

---

## Clés Primaires

Préfixe `#` pour marquer les champs formant l'identifiant unique :

```tsd
// Clé simple
type User(#username: string, email: string)

// Clé composite
type OrderLine(#orderId: string, #productId: string, quantity: number)
```

**Voir** : [Identifiants Internes](../internal-ids.md) pour plus de détails.

---

## Relations

### One-to-Many

```tsd
type Author(#authorId: string, name: string)
type Book(author: Author, #isbn: string, title: string)

tolkien = Author("A001", "J.R.R. Tolkien")
Book(tolkien, "ISBN-001", "The Hobbit")
Book(tolkien, "ISBN-002", "The Lord of the Rings")
```

### Many-to-Many

```tsd
type Student(#studentId: string, name: string)
type Course(#courseId: string, title: string)
type Enrollment(student: Student, course: Course, grade: string)

alice = Student("S001", "Alice")
math = Course("C001", "Math 101")
Enrollment(alice, math, "A")
```

### Hiérarchies

```tsd
type Company(#companyId: string, name: string)
type Department(company: Company, #deptId: string, name: string)
type Employee(dept: Department, #empId: string, name: string)

acme = Company("COMP-001", "ACME Corp")
eng = Department(acme, "DEPT-001", "Engineering")
Employee(eng, "EMP-001", "Alice")
```

---

## Validation de Types

Le parser valide les types automatiquement :

```tsd
type User(#username: string)
type Order(customer: User, #orderNum: string)

laptop = Product("LAP-001", "Laptop", 1200.00)

// ❌ ERREUR : laptop est un Product, pas un User
Order(laptop, "ORD-001")
```

**Erreur** : `type mismatch: expected User, got Product`

---

## Bonnes Pratiques

### 1. Relations Explicites

```tsd
// ✅ BON - Relation explicite via type
type Order(customer: Customer, #orderNum: string)

// ❌ ÉVITER - Relation via string
type Order(customerId: string, #orderNum: string)
```

### 2. Nommage Cohérent

```tsd
// ✅ BON
type Book(author: Author, ...)  // Nom singulier descriptif

// ❌ ÉVITER
type Book(a: Author, ...)  // Nom cryptique
```

### 3. Clés Primaires Appropriées

```tsd
// ✅ BON - Clé naturelle stable
type User(#username: string, email: string)

// ⚠️ ATTENTION - Email peut changer
type User(username: string, #email: string)
```

---

## Résumé

| Aspect | Description |
|--------|-------------|
| **Primitifs** | string, number, boolean |
| **Types de faits** | Champs peuvent être d'autres faits |
| **Clés primaires** | Préfixe `#` |
| **Relations** | One-to-Many, Many-to-Many, Hiérarchiques |
| **Validation** | Type-checking automatique au parsing |

---

## Voir Aussi

- [Affectations de Faits](fact-assignments.md)
- [Comparaisons de Faits](fact-comparisons.md)
- [Identifiants Internes](../internal-ids.md)

---

**Version** : 2.0.0  
**Dernière mise à jour** : 2025-12-19
