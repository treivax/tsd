# Quick Start Guide

Get started with TSD in 5 minutes!

## Installation

```bash
# Clone and build
git clone https://github.com/treivax/tsd.git
cd tsd
make build

# Verify installation
./bin/tsd --version
```

Or using Go:

```bash
go install github.com/treivax/tsd/cmd/tsd@latest
```

## Your First TSD Program

### 1. Create a Simple Rule

Create a file named `hello.tsd`:

```tsd
// Define a type
type Person(name: string, age: number)

// Define an action
action greet(name: string)

// Define a rule
rule welcome : {p: Person} / p.age >= 18 ==> greet(p.name)

// Add a fact
Person(name: "Alice", age: 25)
Person(name: "Bob", age: 16)
```

### 2. Run Your Program

```bash
tsd hello.tsd
```

**Output:**
```
🎯 ACTION EXÉCUTÉE: greet("Alice")
```

Only Alice's greeting is executed because Bob is under 18!

## Core Concepts

### 1. Types

Define the structure of your data:

```tsd
type Product(name: string, price: number, inStock: bool)
type Order(id: string, quantity: number, total: number)
```

### 2. Facts

Create instances of your types:

```tsd
Product(name: "Laptop", price: 999.99, inStock: true)
Product(name: "Mouse", price: 29.99, inStock: false)
```

### 3. Rules

Define business logic with pattern matching:

```tsd
// Rule structure: rule name : {pattern} / condition ==> action
rule expensive : {p: Product} / p.price > 500 ==> markAsPremium(p.name)
```

**Pattern:** `{p: Product}` - Match facts of type Product, bind to variable `p`
**Condition:** `p.price > 500` - Filter for expensive products
**Action:** `markAsPremium(p.name)` - Execute action with product name

### 4. Actions

Declare actions that rules can trigger:

```tsd
action markAsPremium(name: string)
action sendEmail(to: string, subject: string, body: string)
action createInvoice(orderId: string, amount: number)
```

## Common Patterns

### Pattern 1: Multiple Conditions

Combine conditions with `AND` and `OR`:

```tsd
type User(name: string, age: number, premium: bool)
action sendOffer(name: string)

rule targetUser : {u: User} / 
    u.age >= 18 AND u.age <= 65 AND u.premium == true 
    ==> sendOffer(u.name)

User(name: "Alice", age: 30, premium: true)
```

### Pattern 2: Multiple Facts

Match multiple facts together:

```tsd
type Customer(id: string, name: string, vip: bool)
type Order(customerId: string, total: number)
action applyDiscount(customerName: string, orderId: string)

rule vipDiscount : {c: Customer, o: Order} / 
    c.id == o.customerId AND c.vip == true AND o.total > 100 
    ==> applyDiscount(c.name, o.customerId)

Customer(id: "C001", name: "Alice", vip: true)
Order(customerId: "C001", total: 250.00)
```

### Pattern 3: String Operations

Use string operators for pattern matching:

```tsd
type Email(address: string, subject: string)
action flagAsSpam(address: string)

// Check if subject contains "URGENT"
rule spamFilter : {e: Email} / 
    e.subject CONTAINS "URGENT" 
    ==> flagAsSpam(e.address)

Email(address: "spam@example.com", subject: "URGENT: Act now!")
```

### Pattern 4: Type Casting

Convert values between types:

```tsd
type Product(name: string, price: number, quantity: number)
action notify(message: string)

// Cast number to string for concatenation
rule priceAlert : {p: Product} / p.price > 100 ==> 
    notify("High price: $" + (string)p.price)

Product(name: "Laptop", price: 999.99, quantity: 5)
```

**Available casts:**
- `(number)value` - Convert to number
- `(string)value` - Convert to string  
- `(bool)value` - Convert to boolean

### Pattern 5: Arithmetic Operations

Perform calculations in conditions and actions:

```tsd
type Order(id: string, price: number, quantity: number)
action createInvoice(orderId: string, total: number)

rule calculateTotal : {o: Order} / o.quantity > 0 ==> 
    createInvoice(o.id, o.price * o.quantity)

Order(id: "ORD001", price: 50.00, quantity: 3)
```

**Supported operators:** `+`, `-`, `*`, `/`, `%`

## Advanced Features

### Negation (NOT)

Match when a condition is false:

```tsd
type User(email: string, verified: bool)
action sendVerificationEmail(email: string)

rule needsVerification : {u: User} / 
    NOT(u.verified) 
    ==> sendVerificationEmail(u.email)

User(email: "user@example.com", verified: false)
```

### String Patterns

Use `LIKE` for SQL-style patterns or `MATCHES` for regex:

```tsd
type File(name: string, path: string)
action processImage(name: string)

// LIKE: % = any chars, _ = single char
rule imageFiles : {f: File} / 
    f.name LIKE "%.png" OR f.name LIKE "%.jpg" 
    ==> processImage(f.name)

// MATCHES: full regex support
rule configFiles : {f: File} / 
    f.path MATCHES "^/etc/.+\\.conf$" 
    ==> processConfig(f.name)

File(name: "photo.png", path: "/images/photo.png")
```

### Collection Operations

Check membership with `IN`:

```tsd
type User(name: string, role: string)
action grantAccess(name: string)

rule adminAccess : {u: User} / 
    u.role IN ["admin", "superuser", "root"] 
    ==> grantAccess(u.name)

User(name: "Alice", role: "admin")
```

## Running Different Modes

### Compiler Mode (Default)

```bash
# Run TSD program
tsd program.tsd

# With verbose output
tsd -v program.tsd

# With debug logging
TSD_LOG_LEVEL=debug tsd program.tsd
```

### Server Mode

Start TSD as an HTTP server:

```bash
# Start server
tsd server --port 8080

# With authentication
tsd server --port 8080 --auth-key-file api-key.txt

# With TLS
tsd server --port 8443 --tls-cert cert.pem --tls-key key.pem
```

### Client Mode

Send programs to a TSD server:

```bash
# Send program to server
tsd client --url http://localhost:8080 program.tsd

# With API key
tsd client --url http://localhost:8080 --api-key YOUR_KEY program.tsd
```

### Authentication Mode

Manage API keys:

```bash
# Generate API key
tsd auth generate-key --output api-key.txt

# Generate JWT
tsd auth generate-jwt --user admin --output token.txt

# Validate token
tsd auth validate-token --token YOUR_TOKEN
```

## Project Structure

Organize larger projects:

```
my-project/
├── types/
│   ├── user.tsd       # Type definitions
│   └── order.tsd
├── rules/
│   ├── validation.tsd # Business rules
│   └── pricing.tsd
├── facts/
│   └── initial.tsd    # Initial data
└── main.tsd           # Main program
```

Run with multiple files:

```bash
tsd types/*.tsd rules/*.tsd facts/*.tsd main.tsd
```

## Common Issues

### Issue: Rule Not Firing

**Problem:** Your rule doesn't execute even though you expect it to.

**Solutions:**
1. Check condition syntax (use `==` not `=`)
2. Verify fact types match pattern
3. Add debug logging: `TSD_LOG_LEVEL=debug tsd program.tsd`
4. Check for type mismatches (numbers vs strings)

### Issue: Type Errors

**Problem:** "type mismatch" or "invalid operation" errors.

**Solutions:**
1. Use explicit casts: `(string)numberValue`
2. Check field types in type definitions
3. Verify arithmetic operations use numbers
4. For string concatenation, both operands must be strings

### Issue: Pattern Not Matching

**Problem:** Multi-fact pattern doesn't match.

**Solutions:**
1. Ensure all referenced facts exist
2. Check join conditions (variable equality)
3. Verify fact types are correct
4. Test each pattern component separately

## Examples

Explore complete examples in the `examples/` directory:

```bash
# List available examples
ls examples/

# Run an example
tsd examples/basic-rules.tsd
tsd examples/type-casting.tsd
tsd examples/string-operations.tsd
```

## Next Steps

Now that you understand the basics:

1. **Read the [Tutorial](TUTORIAL.md)** - Detailed walkthrough with explanations
2. **Read the [User Guide](USER_GUIDE.md)** - Complete language reference
3. **Explore [Grammar Guide](GRAMMAR_GUIDE.md)** - Deep dive into syntax
4. **Check [Examples](../examples/)** - Real-world use cases
5. **Review [API Reference](API_REFERENCE.md)** - Server API documentation

## Need Help?

- **Documentation:** [docs/](.)
- **Examples:** [examples/](../examples/)
- **Issues:** GitHub Issues
- **Debug:** `TSD_LOG_LEVEL=debug tsd program.tsd`

## Cheat Sheet

```tsd
// Types
type Name(field: string, count: number, active: bool)

// Facts
Name(field: "value", count: 42, active: true)

// Actions
action doSomething(arg1: string, arg2: number)

// Rules
rule name : {x: Type} / condition ==> action(x.field)

// Operators
x == y          // Equal
x != y          // Not equal
x < y           // Less than
x > y           // Greater than
x <= y          // Less or equal
x >= y          // Greater or equal
x AND y         // Logical AND
x OR y          // Logical OR
NOT(x)          // Logical NOT
x + y           // Add (numbers) or concatenate (strings)
x - y           // Subtract
x * y           // Multiply
x / y           // Divide
x % y           // Modulo
x CONTAINS y    // String contains
x IN [a, b]     // In collection
x LIKE "%.txt"  // SQL-style pattern
x MATCHES "^a"  // Regex pattern

// Casts
(number)value   // To number
(string)value   // To string
(bool)value     // To boolean

// Comments
// Single line comment
/* Multi-line
   comment */
```

Happy coding with TSD! 🚀