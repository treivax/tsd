package constraint

import (
	"testing"
)

// TestParserCallbacks_BasicSyntax tests parser callbacks for basic constraint features
func TestParserCallbacks_BasicSyntax(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		description string
	}{
		{
			name: "Simple comparison with <=",
			input: `
type Person(name: string, age: number)
action notify(msg: string)

rule r1 : {p: Person} / p.age <= 18 ==> notify("minor")
`,
			expectError: false,
			description: "Tests <= comparison operator parser callback",
		},
		{
			name: "Simple comparison with >=",
			input: `
type Person(name: string, age: number)
action notify(msg: string)

rule r2 : {p: Person} / p.age >= 65 ==> notify("senior")
`,
			expectError: false,
			description: "Tests >= comparison operator parser callback",
		},
		{
			name: "Simple comparison with !=",
			input: `
type Person(name: string, age: number)
action notify(msg: string)

rule r3 : {p: Person} / p.name != "John" ==> notify("not john")
`,
			expectError: false,
			description: "Tests != comparison operator parser callback",
		},
		{
			name: "Logical OR operator",
			input: `
type Person(name: string, age: number)
action notify(msg: string)

rule r4 : {p: Person} / p.age < 18 OR p.age > 65 ==> notify("eligible")
`,
			expectError: false,
			description: "Tests OR logical operator parser callback",
		},
		{
			name: "String with escape sequences",
			input: `
type Message(text: string)
action log(msg: string)

rule r5 : {m: Message} / m.text == "Hello\nWorld" ==> log("found")
`,
			expectError: false,
			description: "Tests escape sequence parser callback",
		},
		{
			name: "Parenthesized complex expression",
			input: `
type Person(name: string, age: number)
action process(id: string)

rule r6 : {p: Person} / (p.age > 18 AND p.age < 30) OR p.age > 65 ==> process("eligible")
`,
			expectError: false,
			description: "Tests parenthesized constraint parser callback",
		},
		{
			name: "Multiple types in rule",
			input: `
type Person(id: string, name: string, age: number)
type Order(person_id: string, amount: number)
action process(name: string, amt: number)

rule r7 : {p: Person, o: Order} / p.id == o.person_id AND o.amount > 100 ==> process(p.name, o.amount)
`,
			expectError: false,
			description: "Multiple type definitions and join",
		},
		{
			name: "Arithmetic expression with multiplication",
			input: `
type Product(price: number, quantity: number)
action alert(val: number)

rule r8 : {pr: Product} / pr.price * pr.quantity > 100 ==> alert(100)
`,
			expectError: false,
			description: "Tests arithmetic expression parser callback",
		},
		{
			name: "Arithmetic with subtraction",
			input: `
type Account(balance: number, fee: number)
action warn(msg: string)

rule r9 : {a: Account} / a.balance - a.fee < 0 ==> warn("insufficient")
`,
			expectError: false,
			description: "Tests subtraction in arithmetic expression",
		},
		{
			name: "Arithmetic with division",
			input: `
type Stats(total: number, count: number)
action report(avg: number)

rule r10 : {s: Stats} / s.total / s.count > 10 ==> report(10)
`,
			expectError: false,
			description: "Tests division in arithmetic expression",
		},
		{
			name: "Negative number literal",
			input: `
type Temperature(celsius: number)
action alert(msg: string)

rule r11 : {t: Temperature} / t.celsius < -10 ==> alert("cold")
`,
			expectError: false,
			description: "Tests negative number parser callback",
		},
		{
			name: "Float number literal",
			input: `
type Measurement(value: number)
action log(val: number)

rule r12 : {m: Measurement} / m.value > 3.14159 ==> log(3.14)
`,
			expectError: false,
			description: "Tests float number parser callback",
		},
		{
			name: "Boolean field",
			input: `
type Account(active: bool, balance: number)
action process(id: string)

rule r13 : {a: Account} / a.active == true AND a.balance > 0 ==> process("id")
`,
			expectError: false,
			description: "Boolean field constraint",
		},
		{
			name: "Multiple comparison operators",
			input: `
type Person(age: number)
action notify(msg: string)

rule r14 : {p: Person} / p.age >= 18 AND p.age <= 65 ==> notify("working age")
`,
			expectError: false,
			description: "Multiple comparison operators in one rule",
		},
		{
			name: "String concatenation",
			input: `
type Person(firstName: string, lastName: string)
action greet(name: string)

rule r15 : {p: Person} / p.firstName == "John" ==> greet(p.lastName)
`,
			expectError: false,
			description: "String field access",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse("test", []byte(tt.input))

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none for %s", tt.description)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.description, err)
			}

			if !tt.expectError && err == nil {
				if result == nil {
					t.Errorf("Expected non-nil result for %s", tt.description)
				}
			}
		})
	}
}

// TestParserCallbacks_AllComparisonOperators ensures all comparison operators are tested
func TestParserCallbacks_AllComparisonOperators(t *testing.T) {
	operators := []struct {
		op       string
		expected string
	}{
		{"==", "equal"},
		{"!=", "not equal"},
		{"<", "less than"},
		{">", "greater than"},
		{"<=", "less than or equal"},
		{">=", "greater than or equal"},
	}

	for _, op := range operators {
		t.Run(op.op, func(t *testing.T) {
			input := `
type Person(name: string, age: number)
action notify(msg: string)

rule test_op : {p: Person} / p.age ` + op.op + ` 18 ==> notify("match")
`
			result, err := Parse("test", []byte(input))
			if err != nil {
				t.Errorf("Parse error for operator %s: %v", op.op, err)
				return
			}

			if result == nil {
				t.Errorf("Expected non-nil result for operator %s", op.op)
			}
		})
	}
}

// TestParserCallbacks_AllLogicalOperators ensures all logical operators are tested
func TestParserCallbacks_AllLogicalOperators(t *testing.T) {
	operators := []struct {
		op       string
		expected string
	}{
		{"AND", "and"},
		{"OR", "or"},
	}

	for _, op := range operators {
		t.Run(op.op, func(t *testing.T) {
			input := `
type Person(name: string, age: number)
action notify(msg: string)

rule test_logical : {p: Person} / p.age > 18 ` + op.op + ` p.age < 65 ==> notify("match")
`
			result, err := Parse("test", []byte(input))
			if err != nil {
				t.Errorf("Parse error for logical operator %s: %v", op.op, err)
				return
			}

			if result == nil {
				t.Errorf("Expected non-nil result for logical operator %s", op.op)
			}
		})
	}
}

// TestParserCallbacks_EscapeSequences tests all escape sequences in strings
func TestParserCallbacks_EscapeSequences(t *testing.T) {
	escapes := []struct {
		sequence string
		desc     string
	}{
		{`\n`, "newline"},
		{`\t`, "tab"},
		{`\r`, "carriage return"},
		{`\"`, "quote"},
		{`\\`, "backslash"},
	}

	for _, esc := range escapes {
		t.Run(esc.desc, func(t *testing.T) {
			input := `
type Message(text: string)
action log(msg: string)

rule test_escape : {m: Message} / m.text == "test` + esc.sequence + `value" ==> log("match")
`
			result, err := Parse("test", []byte(input))
			if err != nil {
				t.Errorf("Parse error for escape %s: %v", esc.desc, err)
				return
			}

			if result == nil {
				t.Errorf("Expected non-nil result for escape %s", esc.desc)
			}
		})
	}
}

// TestParserCallbacks_EdgeCases tests edge cases and complex scenarios
func TestParserCallbacks_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		description string
	}{
		{
			name: "Multiple types in program",
			input: `
type Person(name: string, age: number)
type Order(id: string, amount: number)
type Product(name: string, price: number)

action process(msg: string)

rule multi : {p: Person} / p.age > 18 ==> process("ok")
`,
			expectError: false,
			description: "Multiple type definitions",
		},
		{
			name: "Multiple rules",
			input: `
type Person(name: string, age: number)
action notify(msg: string)

rule r1 : {p: Person} / p.age < 18 ==> notify("minor")
rule r2 : {p: Person} / p.age >= 18 AND p.age < 65 ==> notify("adult")
rule r3 : {p: Person} / p.age >= 65 ==> notify("senior")
`,
			expectError: false,
			description: "Multiple rules in one program",
		},
		{
			name: "Complex nested expression",
			input: `
type Person(name: string, age: number, active: bool)
action process(id: string)

rule complex : {p: Person} / (p.age > 18 AND p.age < 65) OR (p.active == true AND p.age >= 65) ==> process("eligible")
`,
			expectError: false,
			description: "Complex nested boolean expression",
		},
		{
			name: "Action with multiple parameters",
			input: `
type Person(name: string, age: number)
action notify(recipient: string, message: string, priority: number)

rule multi_param : {p: Person} / p.age > 18 ==> notify(p.name, "hello", 1)
`,
			expectError: false,
			description: "Action call with multiple parameters",
		},
		{
			name: "Optional action parameter",
			input: `
type Person(name: string)
action greet(name: string, greeting: string = "Hello")

rule optional : {p: Person} / p.name == "John" ==> greet(p.name)
`,
			expectError: false,
			description: "Action with optional parameter",
		},
		{
			name: "Numeric field comparisons",
			input: `
type Product(price: number, stock: number)
action reorder(id: string)

rule reorder_check : {p: Product} / p.stock < 10 AND p.price > 0 ==> reorder("id")
`,
			expectError: false,
			description: "Multiple numeric field comparisons",
		},
		{
			name: "Empty program with just types",
			input: `
type Person(name: string, age: number)
type Order(id: string)
`,
			expectError: false,
			description: "Program with only type definitions",
		},
		{
			name: "Facts in program",
			input: `
type Person(name: string, age: number)

Person(name: "Alice", age: 30)
Person(name: "Bob", age: 25)
`,
			expectError: false,
			description: "Program with facts",
		},
		{
			name: "Mixed arithmetic operators",
			input: `
type Calculation(a: number, b: number, c: number, d: number)
action result(val: number)

rule calc : {x: Calculation} / x.a + x.b * x.c - x.d > 100 ==> result(100)
`,
			expectError: false,
			description: "Multiple arithmetic operators in one expression",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse("test", []byte(tt.input))

			if tt.expectError && err == nil {
				t.Errorf("Expected error for %s but got none", tt.description)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.description, err)
			}

			if !tt.expectError && err == nil {
				if result == nil {
					t.Errorf("Expected non-nil result for %s", tt.description)
				}
			}
		})
	}
}

// TestParserCallbacks_NumberFormats tests various number formats
func TestParserCallbacks_NumberFormats(t *testing.T) {
	numbers := []struct {
		value string
		desc  string
	}{
		{"42", "integer"},
		{"0", "zero"},
		{"-10", "negative integer"},
		{"3.14", "float"},
		{"-99.99", "negative float"},
		{"0.5", "fraction"},
		{"1000000", "large number"},
	}

	for _, num := range numbers {
		t.Run(num.desc, func(t *testing.T) {
			input := `
type Value(amount: number)
action process(val: number)

rule test_num : {v: Value} / v.amount > ` + num.value + ` ==> process(1)
`
			result, err := Parse("test", []byte(input))
			if err != nil {
				t.Errorf("Parse error for number %s: %v", num.desc, err)
				return
			}

			if result == nil {
				t.Errorf("Parse returned nil result for number %s", num.desc)
			}
		})
	}
}
