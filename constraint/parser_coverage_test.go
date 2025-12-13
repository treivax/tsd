// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"strings"
	"testing"
)

// TestParserNotConstraint tests the NOT constraint parsing (currently 0% coverage)
func TestParserNotConstraint(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "simple NOT constraint",
			input: `
type Person(name: string, age: number)
rule r1 : {p: Person} / NOT (p.age < 18) ==> alert("adult")
`,
			wantErr: false,
		},
		{
			name: "NOT with complex expression",
			input: `
type User(id: string, status: string)
rule r1 : {u: User} / NOT (u.status == "active" AND u.id == "admin") ==> notify("inactive")
`,
			wantErr: false,
		},
		{
			name: "nested NOT constraints",
			input: `
type Data(value: number)
rule r1 : {d: Data} / NOT (NOT (d.value > 0)) ==> process("positive")
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseReader("test", strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestParserAccumulateConstraint tests ACCUMULATE constraint parsing (currently 0% coverage)
func TestParserAccumulateConstraint(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "ACCUMULATE with SUM",
			input: `
type Order(id: string, amount: number)
rule r1 : {o: Order} / SUM(x: Order / x.id == o.id ; x.amount) > 1000 ==> alert("high")
`,
			wantErr: false,
		},
		{
			name: "ACCUMULATE with AVG",
			input: `
type Metric(sensor: string, value: number)
rule r1 : {m: Metric} / AVG(x: Metric / x.sensor == m.sensor ; x.value) > 50 ==> notify("avg")
`,
			wantErr: false,
		},
		{
			name: "ACCUMULATE with MIN",
			input: `
type Temperature(location: string, temp: number)
rule r1 : {t: Temperature} / MIN(x: Temperature / x.location == t.location ; x.temp) < 0 ==> alert("freeze")
`,
			wantErr: false,
		},
		{
			name: "ACCUMULATE with MAX",
			input: `
type Pressure(device: string, value: number)
rule r1 : {p: Pressure} / MAX(x: Pressure / x.device == p.device ; x.value) > 100 ==> warn("max")
`,
			wantErr: false,
		},
		{
			name: "ACCUMULATE with COUNT (already covered but comprehensive)",
			input: `
type Event(type: string, severity: number)
rule r1 : {e: Event} / COUNT(x: Event / x.type == e.type AND x.severity > 5) > 10 ==> escalate("many")
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseReader("test", strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestParserParenthesizedFactor tests parenthesized arithmetic expressions (Factor2 - 0% coverage)
func TestParserParenthesizedFactor(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "simple parenthesized expression",
			input: `
type Calc(a: number, b: number, c: number)
rule r1 : {x: Calc} / (x.a + x.b) * x.c > 100 ==> compute("result")
`,
			wantErr: false,
		},
		{
			name: "nested parentheses",
			input: `
type Math(x: number, y: number, z: number)
rule r1 : {m: Math} / ((m.x + m.y) * (m.z - 1)) > 50 ==> calculate("complex")
`,
			wantErr: false,
		},
		{
			name: "parentheses with division",
			input: `
type Value(numerator: number, denominator: number, factor: number)
rule r1 : {v: Value} / (v.numerator / v.denominator) * v.factor > 10 ==> process("div")
`,
			wantErr: false,
		},
		{
			name: "parentheses changing precedence",
			input: `
type Expr(a: number, b: number, c: number)
rule r1 : {e: Expr} / e.a * (e.b + e.c) != e.a * e.b + e.a * e.c ==> flag("precedence")
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseReader("test", strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestParserAdditionalComparisonOperators tests operators at 0% coverage
func TestParserAdditionalComparisonOperators(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "CONTAINS operator",
			input: `
type Text(content: string, keyword: string)
rule r1 : {t: Text} / t.content CONTAINS t.keyword ==> match("found")
`,
			wantErr: false,
		},
		{
			name: "LIKE operator",
			input: `
type Message(text: string, pattern: string)
rule r1 : {m: Message} / m.text LIKE m.pattern ==> process("like")
`,
			wantErr: false,
		},
		{
			name: "IN operator",
			input: `
type Item(category: string)
rule r1 : {i: Item} / i.category IN ["electronics", "books", "toys"] ==> validate("in")
`,
			wantErr: false,
		},
		{
			name: "MATCHES operator (regex)",
			input: `
type Pattern(input: string, regex: string)
rule r1 : {p: Pattern} / p.input MATCHES p.regex ==> check("regex")
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseReader("test", strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestParserAdditionalFunctions tests function names at 0% coverage
func TestParserAdditionalFunctions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "TRIM function",
			input: `
type Text(value: string)
rule r1 : {t: Text} / TRIM(t.value) == "test" ==> process("trimmed")
`,
			wantErr: false,
		},
		{
			name: "SUBSTRING function",
			input: `
type String(text: string)
rule r1 : {s: String} / SUBSTRING(s.text, 0, 5) == "hello" ==> match("sub")
`,
			wantErr: false,
		},
		{
			name: "ABS function",
			input: `
type Number(value: number)
rule r1 : {n: Number} / ABS(n.value) > 10 ==> alert("absolute")
`,
			wantErr: false,
		},
		{
			name: "ROUND function",
			input: `
type Decimal(value: number)
rule r1 : {d: Decimal} / ROUND(d.value) > 10 ==> alert("rounded")
`,
			wantErr: false,
		},
		{
			name: "FLOOR function",
			input: `
type Measure(reading: number)
rule r1 : {m: Measure} / FLOOR(m.reading) < 5 ==> process("floor")
`,
			wantErr: false,
		},
		{
			name: "CEIL function",
			input: `
type Value(amount: number)
rule r1 : {v: Value} / CEIL(v.amount) > 100 ==> check("ceiling")
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseReader("test", strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestParserEscapeSequences tests escape sequences at 0% coverage
func TestParserEscapeSequences(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "backspace escape",
			input: `
type Text(value: string)
rule r1 : {t: Text} / t.value == "test\bvalue" ==> process("backspace")
`,
			wantErr: false,
		},
		{
			name: "form feed escape",
			input: `
type Content(text: string)
rule r1 : {c: Content} / c.text == "page1\fpage2" ==> handle("formfeed")
`,
			wantErr: false,
		},
		{
			name: "mixed escape sequences",
			input: `
type Message(data: string)
rule r1 : {m: Message} / m.data == "line1\nline2\ttab\rcarriage\bback\fform" ==> parse("mixed")
`,
			wantErr: false,
		},
		{
			name: "single quotes with escapes",
			input: `
type Info(text: string)
rule r1 : {i: Info} / i.text == 'test\b\fvalue' ==> check("single")
`,
			wantErr: false,
		},
		{
			name: "escaped single quote",
			input: `
type Quote(text: string)
rule r1 : {q: Quote} / q.text == 'it\'s working' ==> check("quote")
`,
			wantErr: false,
		},
		{
			name: "mixed quotes and escapes",
			input: `
type Mixed(a: string, b: string)
rule r1 : {m: Mixed} / m.a == "test\"value" AND m.b == 'it\'s ok' ==> check("both")
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseReader("test", strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestParserComplexCombinations tests combinations of features to maximize coverage
func TestParserComplexCombinations(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "NOT with ACCUMULATE",
			input: `
type Event(id: string, count: number)
rule r1 : {e: Event} / NOT (COUNT(x: Event / x.id == e.id) > 10) ==> process("few")
`,
			wantErr: false,
		},
		{
			name: "parentheses with comparison operators",
			input: `
type Data(a: number, b: number, text: string)
rule r1 : {d: Data} / (d.a + d.b) > 100 AND d.text CONTAINS "test" ==> alert("combo")
`,
			wantErr: false,
		},
		{
			name: "functions with string operators",
			input: `
type Record(name: string, prefix: string)
rule r1 : {r: Record} / UPPER(r.name) LIKE UPPER(r.prefix) ==> match("upper")
`,
			wantErr: false,
		},
		{
			name: "nested parentheses with NOT",
			input: `
type Complex(x: number, y: number, z: number)
rule r1 : {c: Complex} / NOT ((c.x + c.y) * c.z < 0) ==> validate("positive")
`,
			wantErr: false,
		},
		{
			name: "multiple patterns with NOT and accumulate",
			input: `
type Sensor(id: string, value: number)
type Alert(sensor_id: string, threshold: number)
rule r1 : {s: Sensor, a: Alert} / s.id == a.sensor_id AND NOT (s.value < a.threshold) AND COUNT(x: Sensor / x.id == s.id) > 5 ==> trigger("complex")
`,
			wantErr: false,
		},
		{
			name: "all escape sequences together",
			input: `
type Document(content: string)
rule r1 : {d: Document} / d.content == "test\n\t\r\b\f\\\"test" ==> parse("escapes")
`,
			wantErr: false,
		},
		{
			name: "arithmetic with all operators and parentheses",
			input: `
type Math(a: number, b: number, c: number, d: number)
rule r1 : {m: Math} / ((m.a + m.b) * (m.c - m.d)) / (m.a + 1) > 50 ==> compute("full")
`,
			wantErr: false,
		},
		{
			name: "string functions with operators",
			input: `
type Text(value: string)
rule r1 : {t: Text} / UPPER(TRIM(t.value)) CONTAINS "NEW" AND LENGTH(TRIM(t.value)) > 5 ==> transform("string")
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseReader("test", strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestParserEdgeCases tests edge cases and boundary conditions
func TestParserEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "empty parentheses should fail",
			input: `
type Data(value: number)
rule r1 : {d: Data} / () > 10 ==> fail("empty")
`,
			wantErr: true,
		},
		{
			name: "deeply nested parentheses",
			input: `
type Deep(a: number, b: number, c: number, d: number)
rule r1 : {x: Deep} / ((((x.a + x.b) * x.c) - x.d) + 1) > 0 ==> process("deep")
`,
			wantErr: false,
		},
		{
			name: "NOT with EXISTS",
			input: `
type User(id: string, role: string)
type Permission(user_id: string, access: string)
rule r1 : {u: User} / NOT (EXISTS (p: Permission / p.user_id == u.id AND p.access == "admin")) ==> restrict("limited")
`,
			wantErr: false,
		},
		{
			name: "accumulate in complex condition",
			input: `
type Sale(product: string, amount: number, region: string)
rule r1 : {s: Sale} / s.region == "north" AND SUM(x: Sale / x.product == s.product ; x.amount) > 10000 AND s.amount > 100 ==> reward("bonus")
`,
			wantErr: false,
		},
		{
			name: "multiple NOT operators",
			input: `
type Flag(active: bool, visible: bool)
rule r1 : {f: Flag} / NOT (f.active == false) AND NOT (f.visible == false) ==> display("both")
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseReader("test", strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestParserAllAccumulateFunctions ensures all accumulate function variants are tested
func TestParserAllAccumulateFunctions(t *testing.T) {
	input := `
type Metric(category: string, value: number, score: number)

rule count_rule : {m: Metric} / COUNT(x: Metric / x.category == m.category) > 10 ==> action1("count")
rule sum_rule : {m: Metric} / SUM(x: Metric / x.category == m.category ; x.value) > 1000 ==> action2("sum")
rule avg_rule : {m: Metric} / AVG(x: Metric / x.category == m.category ; x.value) > 50 ==> action3("avg")
rule min_rule : {m: Metric} / MIN(x: Metric / x.category == m.category ; x.score) < 10 ==> action4("min")
rule max_rule : {m: Metric} / MAX(x: Metric / x.category == m.category ; x.score) > 100 ==> action5("max")

action action1(msg: string)
action action2(msg: string)
action action3(msg: string)
action action4(msg: string)
action action5(msg: string)
`

	_, err := ParseReader("test", strings.NewReader(input))
	if err != nil {
		t.Errorf("ParseReader() failed for all accumulate functions: %v", err)
	}
}

// TestParserAllStringOperators tests all string comparison operators
func TestParserAllStringOperators(t *testing.T) {
	input := `
type StringData(text: string, pattern: string, list: string)

rule contains_rule : {s: StringData} / s.text CONTAINS s.pattern ==> check1("contains")
rule like_rule : {s: StringData} / s.text LIKE s.pattern ==> check2("like")
rule matches_rule : {s: StringData} / s.text MATCHES s.pattern ==> check3("matches")
rule in_rule : {s: StringData} / s.list IN ["a", "b", "c"] ==> check4("in")

action check1(msg: string)
action check2(msg: string)
action check3(msg: string)
action check4(msg: string)
`

	_, err := ParseReader("test", strings.NewReader(input))
	if err != nil {
		t.Errorf("ParseReader() failed for all string operators: %v", err)
	}
}
