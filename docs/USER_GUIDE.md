# TSD User Guide

Complete guide to the TSD rule engine and its features.

## Table of Contents

1. [Introduction](#introduction)
2. [Language Syntax](#language-syntax)
3. [Type System](#type-system)
4. [Pattern Matching](#pattern-matching)
5. [Conditions](#conditions)
6. [Actions](#actions)
7. [Type Casting](#type-casting)
8. [String Operations](#string-operations)
9. [Arithmetic Operations](#arithmetic-operations)
10. [Advanced Features](#advanced-features)
11. [Configuration](#configuration)
12. [Authentication](#authentication)
13. [Server Mode](#server-mode)
14. [Best Practices](#best-practices)

## Introduction

TSD (Type System Development) is a rule engine based on the RETE algorithm that allows you to define business rules using a declarative syntax. It features:

- **Pattern Matching:** Match complex fact combinations
- **Type Safety:** Strong typing with validation
- **RETE Algorithm:** Efficient rule evaluation
- **Network Optimization:** Share common conditions across rules
- **Type Casting:** Explicit type conversions
- **String Operations:** Pattern matching and containment checks
- **Authentication:** Secure API access with keys and JWT
- **Multiple Modes:** Compiler, server, client, auth management

## Language Syntax

### Program Structure

A TSD program consists of four main elements:

```tsd
// 1. Type definitions
type Person(name: string, age: number)

// 2. Action declarations
action greet(name: string)

// 3. Rule definitions
rule adult : {p: Person} / p.age >= 18 ==> greet(p.name)

// 4. Fact assertions
Person(name: "Alice", age: 25)
```

### Comments

```tsd
// Single-line comment

/* Multi-line
   comment */

/* Nested /* comments */ are NOT supported */
```

### Case Sensitivity

TSD is **case-insensitive** for keywords but **case-sensitive** for identifiers:

```tsd
// These are equivalent (keywords)
TYPE Person(name: string)
type Person(name: string)
Type Person(name: string)

// These are DIFFERENT (identifiers)
Person(name: "Alice")
person(name: "Bob")      // Different type!
PERSON(name: "Charlie")  // Different type!
```

### Identifiers

Valid identifier rules:

```tsd
// Valid identifiers
myVariable
_underscore
camelCase
PascalCase
snake_case
with123Numbers
αβγ              // Unicode support
名前              // UTF-8 support

// Invalid identifiers
123start         // Cannot start with digit
my-var           // No hyphens
my.var           // No dots (except field access)
```

## Type System

### Primitive Types

TSD supports three primitive types:

```tsd
type Example(
    text: string,     // String values
    count: number,    // Numeric values (float64)
    active: bool      // Boolean values
)
```

### Type Definitions

Define custom types with typed fields:

```tsd
// Simple type
type Product(name: string, price: number)

// Complex type
type Order(
    id: string,
    customerId: string,
    items: string,        // Could be JSON string
    total: number,
    paid: bool,
    createdAt: string
)

// Multiple types
type Customer(id: string, name: string, vip: bool)
type Invoice(orderId: string, amount: number, status: string)
```

### Type Validation

Types are validated at parse time:

```tsd
// ✓ Valid
type User(name: string, age: number)

// ✗ Invalid - unknown type
type User(name: integer, age: float)

// ✗ Invalid - missing type
type User(name, age: number)

// ✗ Invalid - duplicate field
type User(name: string, name: string)
```

## Pattern Matching

### Single Fact Patterns

Match facts of a specific type:

```tsd
type User(name: string, active: bool)
action process(name: string)

// Match any User fact, bind to variable 'u'
rule processUser : {u: User} / u.active == true ==> process(u.name)

User(name: "Alice", active: true)
```

### Multi-Fact Patterns

Match combinations of facts:

```tsd
type Customer(id: string, name: string, premium: bool)
type Order(id: string, customerId: string, amount: number)
action processOrder(customerName: string, orderId: string, amount: number)

// Match Customer AND Order where IDs match
rule premiumOrder : {c: Customer, o: Order} / 
    c.id == o.customerId AND 
    c.premium == true AND 
    o.amount > 100 
    ==> processOrder(c.name, o.id, o.amount)

Customer(id: "C001", name: "Alice", premium: true)
Order(id: "O001", customerId: "C001", amount: 250.00)
```

### Variable Binding

Variables in patterns bind to matched facts:

```tsd
// 'p' is bound to matched Product fact
rule expensive : {p: Product} / p.price > 1000 ==> alert(p.name)

// Multiple bindings
rule related : {p: Product, c: Category} / 
    p.categoryId == c.id 
    ==> link(p.name, c.name)
```

## Conditions

### Comparison Operators

```tsd
x == y          // Equal
x != y          // Not equal
x < y           // Less than
x > y           // Greater than
x <= y          // Less than or equal
x >= y          // Greater than or equal
```

Examples:

```tsd
type Product(name: string, price: number, stock: number)
action reorder(name: string)
action discount(name: string)

// Equality
rule outOfStock : {p: Product} / p.stock == 0 ==> reorder(p.name)

// Inequality
rule inStock : {p: Product} / p.stock != 0 ==> process(p.name)

// Range checks
rule expensive : {p: Product} / p.price > 1000 ==> premium(p.name)
rule cheap : {p: Product} / p.price < 50 ==> discount(p.name)
rule midRange : {p: Product} / p.price >= 50 AND p.price <= 1000 ==> standard(p.name)
```

### Logical Operators

```tsd
x AND y         // Logical AND (both must be true)
x OR y          // Logical OR (at least one must be true)
NOT(x)          // Logical NOT (negation)
```

Examples:

```tsd
type User(name: string, age: number, verified: bool, premium: bool)
action sendOffer(name: string)
action verify(name: string)

// AND - all conditions must be true
rule eligibleUser : {u: User} / 
    u.age >= 18 AND u.verified == true AND u.premium == true 
    ==> sendOffer(u.name)

// OR - at least one condition must be true
rule needsAttention : {u: User} / 
    u.age < 18 OR NOT(u.verified) 
    ==> review(u.name)

// Complex combinations
rule targetUser : {u: User} / 
    (u.age >= 18 AND u.age <= 65) AND 
    (u.verified == true OR u.premium == true) 
    ==> target(u.name)

// Negation
rule unverified : {u: User} / 
    NOT(u.verified) 
    ==> verify(u.name)
```

### Operator Precedence

From highest to lowest:

1. `NOT`
2. Comparison (`==`, `!=`, `<`, `>`, `<=`, `>=`)
3. `AND`
4. `OR`

Use parentheses for clarity:

```tsd
// Without parentheses (relies on precedence)
rule r1 : {x: Type} / x.a == 1 OR x.b == 2 AND x.c == 3 ==> action1()
// Evaluated as: x.a == 1 OR (x.b == 2 AND x.c == 3)

// With parentheses (explicit grouping)
rule r2 : {x: Type} / (x.a == 1 OR x.b == 2) AND x.c == 3 ==> action2()
```

## Actions

### Action Declarations

Declare actions before using them in rules:

```tsd
// Simple action
action log(message: string)

// Multiple parameters
action sendEmail(to: string, subject: string, body: string)

// Mixed types
action createInvoice(id: string, amount: number, paid: bool)
```

### Action Execution

Actions are triggered when rules fire:

```tsd
type Order(id: string, total: number)
action processPayment(orderId: string, amount: number)
action sendConfirmation(orderId: string)

rule newOrder : {o: Order} / o.total > 0 ==> 
    processPayment(o.id, o.total)

// Multiple actions (executed in order)
rule confirmedOrder : {o: Order} / o.total > 0 ==> {
    processPayment(o.id, o.total)
    sendConfirmation(o.id)
}

Order(id: "ORD001", total: 250.00)
```

### Built-in Actions

TSD provides a built-in `print` action for debugging:

```tsd
type Debug(value: string)

rule debug : {d: Debug} / d.value != "" ==> print(d.value)

Debug(value: "Debug message")
```

## Type Casting

### Casting Syntax

Convert values between types using explicit casts:

```tsd
(number)expression     // Cast to number
(string)expression     // Cast to string
(bool)expression       // Cast to boolean
```

### Number Casting

Convert to number (float64):

```tsd
type Data(strValue: string, boolValue: bool, numValue: number)
action process(value: number)

rule toNumber : {d: Data} / (number)d.strValue > 100 ==> process((number)d.strValue)

// String to Number
(number)"123"        → 123.0
(number)"12.5"       → 12.5
(number)"-45"        → -45.0
(number)"  123  "    → 123.0    // Whitespace trimmed

// Bool to Number
(number)true         → 1.0
(number)false        → 0.0

// Number to Number (identity)
(number)42.5         → 42.5

// Errors
(number)"abc"        → Error: cannot cast 'abc' to number
(number)""           → Error: cannot cast empty string to number

Data(strValue: "150", boolValue: true, numValue: 42.0)
```

### String Casting

Convert to string:

```tsd
type Product(name: string, price: number, inStock: bool)
action notify(message: string)

// Number to String
rule priceAlert : {p: Product} / p.price > 100 ==> 
    notify("High price: $" + (string)p.price)

// Examples
(string)123.0        → "123"       // Integer format
(string)12.5         → "12.5"      // Decimal format
(string)(-45.0)      → "-45"

// Bool to String
(string)true         → "true"
(string)false        → "false"

// String to String (identity)
(string)"hello"      → "hello"

Product(name: "Laptop", price: 999.99, inStock: true)
```

### Boolean Casting

Convert to boolean:

```tsd
type Config(enabled: string, count: number)
action activate()

rule checkEnabled : {c: Config} / (bool)c.enabled == true ==> activate()

// String to Bool
(bool)"true"         → true        // Case insensitive
(bool)"TRUE"         → true
(bool)"1"            → true
(bool)"false"        → false
(bool)"FALSE"        → false
(bool)"0"            → false
(bool)""             → false       // Empty string is false
(bool)"anything"     → false       // Other strings are false

// Number to Bool
(bool)0.0            → false       // Zero is false
(bool)1.0            → true        // Non-zero is true
(bool)(-5.0)         → true
(bool)999.0          → true

// Bool to Bool (identity)
(bool)true           → true

Config(enabled: "true", count: 5)
```

### Casting in Expressions

Use casts in complex expressions:

```tsd
type Order(id: string, priceStr: string, quantityStr: string, urgent: string)
action process(id: string, total: number, isUrgent: bool)

rule complexCast : {o: Order} / 
    (number)o.quantityStr > 5 AND 
    (number)o.priceStr < 100 AND
    (bool)o.urgent == true
    ==> process(o.id, (number)o.priceStr * (number)o.quantityStr, (bool)o.urgent)

Order(id: "O001", priceStr: "50.00", quantityStr: "10", urgent: "true")
```

## String Operations

### String Concatenation

Concatenate strings using the `+` operator:

```tsd
type User(firstName: string, lastName: string, age: number)
action greet(message: string)

// String + String
rule greeting : {u: User} / u.age >= 18 ==> 
    greet("Hello, " + u.firstName + " " + u.lastName)

// String + Casted Number
rule ageInfo : {u: User} / u.age > 0 ==> 
    greet("Age: " + (string)u.age)

// ✓ Valid: Both operands are strings
"Hello" + "World"                    → "Hello World"
"Count: " + (string)42               → "Count: 42"

// ✗ Invalid: Mixed types without cast
"Count: " + 42                       → Error: use explicit cast
42 + "items"                         → Error: use explicit cast

User(firstName: "Alice", lastName: "Smith", age: 30)
```

### CONTAINS Operator

Check if a string contains a substring:

```tsd
type Email(address: string, subject: string, body: string)
action flagAsSpam(address: string)
action important(address: string)

// Case-sensitive substring check
rule spamFilter : {e: Email} / 
    e.subject CONTAINS "URGENT" OR e.body CONTAINS "Click here"
    ==> flagAsSpam(e.address)

rule importantEmail : {e: Email} /
    e.subject CONTAINS "IMPORTANT"
    ==> important(e.address)

Email(address: "user@example.com", subject: "URGENT: Act now!", body: "Click here")
```

### LIKE Operator

SQL-style pattern matching:

```tsd
type File(name: string, extension: string, path: string)
action processImage(name: string)
action processConfig(name: string)

// % matches zero or more characters
// _ matches exactly one character

rule imageFiles : {f: File} / 
    f.name LIKE "%.png" OR f.name LIKE "%.jpg" OR f.name LIKE "%.gif"
    ==> processImage(f.name)

rule configFiles : {f: File} /
    f.name LIKE "config_%.txt"
    ==> processConfig(f.name)

rule singleChar : {f: File} /
    f.name LIKE "file_?.dat"
    ==> process(f.name)

File(name: "photo.png", extension: "png", path: "/images/photo.png")
File(name: "config_prod.txt", extension: "txt", path: "/etc/config_prod.txt")
```

### MATCHES Operator

Regular expression pattern matching:

```tsd
type Code(id: string, value: string)
type Email(address: string)
action validCode(id: string)
action validEmail(address: string)

// Full regex support (Go regex syntax)
rule codeFormat : {c: Code} /
    c.value MATCHES "^CODE[0-9]{3,6}$"
    ==> validCode(c.id)

rule emailFormat : {e: Email} /
    e.address MATCHES "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
    ==> validEmail(e.address)

Code(id: "C001", value: "CODE12345")
Email(address: "user@example.com")
```

### IN Operator

Check membership in a collection:

```tsd
type User(name: string, role: string, status: string)
action grantAccess(name: string)
action review(name: string)

// Check if value is in array
rule adminAccess : {u: User} /
    u.role IN ["admin", "superuser", "root"]
    ==> grantAccess(u.name)

rule suspiciousStatus : {u: User} /
    u.status IN ["suspended", "banned", "locked"]
    ==> review(u.name)

User(name: "Alice", role: "admin", status: "active")
```

## Arithmetic Operations

### Basic Operations

```tsd
x + y          // Addition
x - y          // Subtraction
x * y          // Multiplication
x / y          // Division
x % y          // Modulo (remainder)
```

### In Conditions

Use arithmetic in rule conditions:

```tsd
type Product(name: string, price: number, quantity: number, discount: number)
action alert(name: string, total: number)

// Calculate total in condition
rule highValue : {p: Product} /
    p.price * p.quantity > 1000
    ==> alert(p.name, p.price * p.quantity)

// With discount
rule afterDiscount : {p: Product} /
    (p.price * p.quantity) - p.discount > 500
    ==> process(p.name)

// Percentage calculation
rule discountRate : {p: Product} /
    (p.discount / p.price) * 100 > 20
    ==> highDiscount(p.name)

Product(name: "Laptop", price: 1000.00, quantity: 2, discount: 200.00)
```

### In Actions

Compute values passed to actions:

```tsd
type Order(id: string, subtotal: number, taxRate: number, shipping: number)
action createInvoice(orderId: string, total: number, tax: number)

rule calculateTotal : {o: Order} / o.subtotal > 0 ==>
    createInvoice(
        o.id,
        o.subtotal + (o.subtotal * o.taxRate) + o.shipping,
        o.subtotal * o.taxRate
    )

Order(id: "O001", subtotal: 100.00, taxRate: 0.20, shipping: 10.00)
```

### Operator Precedence

Standard mathematical precedence:

1. `()` - Parentheses (highest)
2. `*`, `/`, `%` - Multiplication, Division, Modulo
3. `+`, `-` - Addition, Subtraction (lowest)

```tsd
// Examples
2 + 3 * 4           → 14         // 3*4 first, then +2
(2 + 3) * 4         → 20         // Parentheses first
10 - 4 / 2          → 8          // 4/2 first, then 10-2
(10 - 4) / 2        → 3          // Parentheses first
10 % 3 + 2          → 5          // 10%3=1, then 1+2
```

### Division by Zero

Division and modulo by zero result in errors:

```tsd
type Data(value: number, divisor: number)
action result(value: number)

// Safe division - check divisor first
rule safeDivide : {d: Data} /
    d.divisor != 0 AND (d.value / d.divisor) > 10
    ==> result(d.value / d.divisor)

// ✓ Valid: divisor is checked
Data(value: 100.0, divisor: 5.0)

// ✗ Runtime error if divisor is zero
Data(value: 100.0, divisor: 0.0)
```

## Advanced Features

### Strong Mode

Enable strict type checking and validation:

```tsd
// Strong mode requires explicit type declarations
// and validates all field accesses at compile time

type User(name: string, age: number)

// ✓ Valid in strong mode
rule validRule : {u: User} / u.age >= 18 ==> process(u.name)

// ✗ Invalid in strong mode - unknown field
rule invalidRule : {u: User} / u.email != "" ==> process(u.email)
```

Enable strong mode:

```bash
tsd --strong-mode program.tsd
```

### Multiple Rules

Rules are evaluated independently:

```tsd
type Order(id: string, total: number, urgent: bool)
action processStandard(id: string)
action processUrgent(id: string)
action applyDiscount(id: string)

// Multiple rules can match the same fact
rule standard : {o: Order} / o.urgent == false ==> processStandard(o.id)
rule urgent : {o: Order} / o.urgent == true ==> processUrgent(o.id)
rule discount : {o: Order} / o.total > 1000 ==> applyDiscount(o.id)

// This order will trigger both 'urgent' and 'discount' rules
Order(id: "O001", total: 1500.00, urgent: true)
```

### Rule Conflicts

When multiple rules match, all fire:

```tsd
type Product(name: string, category: string, price: number)
action processElectronics(name: string)
action processExpensive(name: string)

rule electronics : {p: Product} / p.category == "Electronics" ==> processElectronics(p.name)
rule expensive : {p: Product} / p.price > 1000 ==> processExpensive(p.name)

// Both rules fire for this product
Product(name: "Laptop", category: "Electronics", price: 1500.00)
```

### Optimization

TSD automatically optimizes rule networks:

- **Alpha Node Sharing:** Common conditions shared across rules
- **Beta Network Optimization:** Efficient join computation
- **Index-based Lookups:** Fast fact retrieval

These optimizations are automatic and transparent.

## Configuration

### Command Line Options

```bash
# Basic execution
tsd program.tsd

# Verbose output
tsd -v program.tsd
tsd --verbose program.tsd

# Strong mode
tsd --strong-mode program.tsd

# Help
tsd --help
tsd -h

# Version
tsd --version
tsd -v
```

### Environment Variables

```bash
# Log level (debug, info, warn, error)
export TSD_LOG_LEVEL=debug

# Server port
export TSD_PORT=8080

# API key for authentication
export TSD_API_KEY=your-key-here

# TLS configuration
export TSD_TLS_CERT=/path/to/cert.pem
export TSD_TLS_KEY=/path/to/key.pem
```

### Configuration File

Create `config.yaml`:

```yaml
# Server configuration
server:
  port: 8080
  host: 0.0.0.0
  read_timeout: 30s
  write_timeout: 30s

# Logging
logging:
  level: info           # debug, info, warn, error
  output: stdout        # stdout, stderr, file
  file: /var/log/tsd.log
  format: json          # json, text

# Authentication
authentication:
  enabled: true
  key_file: /etc/tsd/api-key.txt
  jwt_secret_file: /etc/tsd/jwt-secret.txt
  token_expiry: 24h

# TLS/SSL
tls:
  enabled: false
  cert_file: /etc/tsd/cert.pem
  key_file: /etc/tsd/key.pem
  client_auth: none     # none, request, require

# Performance
performance:
  max_concurrent_rules: 1000
  fact_cache_size: 10000
  enable_optimizations: true
```

Load configuration:

```bash
tsd server --config config.yaml
```

## Authentication

### API Key Authentication

Generate and use API keys:

```bash
# Generate API key
tsd auth generate-key --output api-key.txt

# View generated key
cat api-key.txt

# Use with server
tsd server --auth-key-file api-key.txt

# Use with client
tsd client --url http://localhost:8080 --api-key $(cat api-key.txt) program.tsd
```

### JWT Authentication

Generate and validate JWT tokens:

```bash
# Generate JWT secret
tsd auth generate-jwt-secret --output jwt-secret.txt

# Generate token for user
tsd auth generate-jwt --user admin --secret-file jwt-secret.txt --output token.txt

# Validate token
tsd auth validate-token --token $(cat token.txt) --secret-file jwt-secret.txt

# Use with server
tsd server --jwt-secret-file jwt-secret.txt

# Use with client (token in header)
curl -H "Authorization: Bearer $(cat token.txt)" \
     -X POST http://localhost:8080/execute \
     -d @program.tsd
```

### API Key Management

```bash
# Generate with specific prefix
tsd auth generate-key --prefix tsd_ --output key.txt

# Generate with specific length (default 32)
tsd auth generate-key --length 64 --output key.txt

# Validate key format
tsd auth validate-key --key $(cat key.txt)

# Rotate keys (generate new, revoke old)
tsd auth rotate-key --old-key-file old-key.txt --output new-key.txt
```

## Server Mode

### Starting the Server

```bash
# Basic server
tsd server

# Custom port
tsd server --port 8080

# With authentication
tsd server --auth-key-file api-key.txt

# With TLS
tsd server --tls-cert cert.pem --tls-key key.pem --port 8443

# With all options
tsd server \
  --port 8080 \
  --auth-key-file api-key.txt \
  --tls-cert cert.pem \
  --tls-key key.pem \
  --config config.yaml
```

### API Endpoints

#### POST /execute

Execute a TSD program:

```bash
# Without authentication
curl -X POST http://localhost:8080/execute \
     -H "Content-Type: text/plain" \
     -d @program.tsd

# With API key
curl -X POST http://localhost:8080/execute \
     -H "Content-Type: text/plain" \
     -H "X-API-Key: YOUR_API_KEY" \
     -d @program.tsd

# With JWT
curl -X POST http://localhost:8080/execute \
     -H "Content-Type: text/plain" \
     -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     -d @program.tsd
```

#### GET /health

Health check endpoint:

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "ok",
  "version": "1.0.0",
  "uptime": "2h30m15s"
}
```

#### GET /metrics

Prometheus metrics (if enabled):

```bash
curl http://localhost:8080/metrics
```

### Using Client Mode

```bash
# Basic usage
tsd client --url http://localhost:8080 program.tsd

# With API key
tsd client --url http://localhost:8080 --api-key YOUR_KEY program.tsd

# With JWT token
tsd client --url http://localhost:8080 --token YOUR_TOKEN program.tsd

# Multiple files
tsd client --url http://localhost:8080 types.tsd rules.tsd facts.tsd
```

## Best Practices

### 1. Type Design

**✓ Do:**
```tsd
// Use descriptive names
type Customer(id: string, name: string, email: string, vip: bool)

// Group related fields
type Address(street: string, city: string, zipCode: string, country: string)
```

**✗ Don't:**
```tsd
// Avoid generic names
type Data(a: string, b: number, c: bool)

// Avoid mixing unrelated fields
type Mixed(userName: string, orderTotal: number, isActive: bool)
```

### 2. Rule Organization

**✓ Do:**
```tsd
// Use descriptive rule names
rule processVIPOrder : {c: Customer, o: Order} / 
    c.id == o.customerId AND c.vip == true 
    ==> priorityProcessing(o.id)

// Group related rules with prefixes
rule validation_email : {u: User} / u.email MATCHES "^.+@.+\..+$" ==> validateEmail(u.id)
rule validation_age : {u: User} / u.age >= 18 ==> validateAge(u.id)
```

**✗ Don't:**
```tsd
// Avoid cryptic names
rule r1 : {x: Type} / x.a == 1 ==> do(x.b)

// Avoid overly complex rules
rule complex : {a: A, b: B, c: C, d: D} / 
    a.x == b.y AND b.z == c.w AND c.v == d.u AND 
    a.p > 10 OR (b.q < 20 AND c.r != 30) AND NOT(d.s)
    ==> action(a, b, c, d)
```

### 3. Condition Writing

**✓ Do:**
```tsd
// Use clear, readable conditions
rule adult : {u: User} / u.age >= 18 ==> process(u.name)

// Use parentheses for clarity
rule complex : {x: X} / (x.a > 10 AND x.b < 20) OR x.c == 30 ==> action(x)
```

**✗ Don't:**
```tsd
// Avoid relying on precedence
rule unclear : {x: X} / x.a > 10 AND x.b < 20 OR x.c == 30 ==> action(x)
```

### 4. Type Casting

**✓ Do:**
```tsd
// Cast explicitly when mixing types
rule priceDisplay : {p: Product} / p.price > 0 ==>
    display("Price: $" + (string)p.price)

// Cast in conditions
rule stringNumber : {d: Data} / (number)d.value > 100 ==> process(d.id)
```

**✗ Don't:**
```tsd
// Don't rely on implicit conversion (will error)
rule noCast : {p: Product} / p.price > 0 ==>
    display("Price: $" + p.price)  // ERROR: need (string)p.price
```

### 5. Action Arguments

**✓ Do:**
```tsd
// Pass needed data to actions
rule orderPlaced : {o: Order} / o.total > 0 ==>
    processOrder(o.id, o.customerId, o.total)

// Use calculations
rule withTax : {o: Order} / o.subtotal > 0 ==>
    createInvoice(o.id, o.subtotal * 1.20)
```

**✗ Don't:**
```tsd
// Don't pass unnecessary data
rule tooMuch : {o: Order} / o.total > 0 ==>
    process(o.id, o.customerId, o.total, o.status, o.date, o.notes)
```

### 6. Performance

**✓ Do:**
```tsd
// Put selective conditions first
rule specific : {u: User} / 
    u.role == "admin" AND      // Selective
    u.age >= 18                // Less selective
    ==> process(u.id)

// Use indexed fields (id, primary keys)
rule join : {c: Customer, o: Order} /
    c.id == o.customerId       // Index-friendly
    ==> link(c.name, o.id)
```

**✗ Don't:**
```tsd
// Don't start with expensive operations
rule slow : {u: User} /
    u.email MATCHES "^.+@.+\..+$" AND    // Expensive regex
    u.role == "admin"                     // Should be first
    ==> process(u.id)
```

### 7. Testing

**✓ Do:**
```tsd
// Test with minimal examples
type User(name: string, age: number)
action greet(name: string)

rule test : {u: User} / u.age >= 18 ==> greet(u.name)

User(name: "Test", age: 18)    // Boundary case
User(name: "Test", age: 17)    // Should not match
User(name: "Test", age: 19)    // Should match
```

### 8. Documentation

**✓ Do:**
```tsd
// Document complex rules
// This rule processes high-value orders from VIP customers
// and applies a 10% discount automatically
rule vipDiscount : {c: Customer, o: Order} /
    c.id == o.customerId AND 
    c.vip == true AND 
    o.total > 1000
    ==> applyDiscount(o.id, o.total * 0.10)
```

## Troubleshooting

### Common Errors

**Parse Error: Unexpected Token**
```
Error: unexpected token 'xyz' at line 10
```
- Check syntax, missing semicolons, typos
- Verify keywords are spelled correctly
- Check for unmatched parentheses or braces

**Type Error: Unknown Type**
```
Error: unknown type 'Prodcut' at line 15
```
- Check type name spelling
- Ensure type is defined before use
- Remember: identifiers are case-sensitive

**Type Error: Field Not Found**
```
Error: field 'priice' not found in type Product
```
- Check field name spelling
- Verify field exists in type definition
- Remember: field names are case-sensitive

**Runtime Error: Division by Zero**
```
Error: division by zero in rule calculation
```
- Add checks before division: `x != 0 AND y / x > 10`
- Validate input data

**Cast Error: Invalid Conversion**
```
Error: cannot cast 'abc' to number
```
- Validate string format before casting
- Handle invalid data appropriately

### Debug Tips

1. **Enable Debug Logging:**
   ```bash
   TSD_LOG_LEVEL=debug tsd program.tsd
   ```

2. **Test Rules Individually:**
   ```tsd
   // Comment out other rules
   // rule other1 : ...
   // rule other2 : ...
   
   rule testThis : {x: X} / x.a == 1 ==> print("Testing")
   ```

3. **Use Print for Debugging:**
   ```tsd
   rule debug : {x: X} / x.a > 0 ==> print("x.a = " + (string)x.a)
   ```

4. **Simplify Conditions:**
   ```tsd
   // Complex
   rule complex : {x: X} / x.a > 10 AND x.b < 20 AND x.c == 30 ==> action(x)
   
   // Test parts separately
   rule testA : {x: X} / x.a > 10 ==> print("A passed")
   rule testB : {x: X} / x.b < 20 ==> print("B passed")
   rule testC : {x: X} / x.c == 30 ==> print("C passed")
   ```

## Additional Resources

- [Quick Start Guide](QUICK_START.md) - Get started in 5 minutes
- [Tutorial](TUTORIAL.md) - Step-by-step learning
- [Grammar Guide](GRAMMAR_GUIDE.md) - Complete syntax reference
- [API Reference](API_REFERENCE.md) - HTTP API documentation
- [Architecture](ARCHITECTURE.md) - Technical design
- [Examples](../examples/) - Real-world examples

## Glossary

- **Fact:** An instance of a type with specific values
- **Pattern:** A template for matching facts in rules
- **Condition:** A boolean expression that filters matches
- **Action:** An operation executed when a rule fires
- **Rule:** A pattern-condition-action triple
- **Alpha Network:** Filters facts by type and conditions
- **Beta Network:** Joins multiple facts
- **RETE:** The algorithm used for efficient rule evaluation
- **Binding:** Association of a variable with a matched fact
- **Cast:** Explicit type conversion