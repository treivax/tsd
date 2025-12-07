// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package constraint

import (
	"testing"
)

// TestBug_CaseInsensitiveKeywords_Fixed vérifie que les mots-clés
// (AND, OR, NOT, EXISTS, AVG, COUNT, etc.) sont acceptés dans les trois formes:
// UPPERCASE, lowercase et Capitalized (mais pas les formes mixtes arbitraires)
func TestBug_CaseInsensitiveKeywords_Fixed(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldParse bool
		description string
	}{
		// Opérateurs logiques
		{
			name:        "AND uppercase (current)",
			input:       `rule test1 : {p:Person} / p.age > 18 AND p.age < 65 ==> action1()`,
			shouldParse: true,
			description: "AND en majuscules devrait fonctionner (état actuel)",
		},
		{
			name:        "and lowercase",
			input:       `rule test2 : {p:Person} / p.age > 18 and p.age < 65 ==> action1()`,
			shouldParse: true,
			description: "and en minuscules devrait fonctionner",
		},
		{
			name:        "And mixed case",
			input:       `rule test3 : {p:Person} / p.age > 18 And p.age < 65 ==> action1()`,
			shouldParse: true,
			description: "And en casse mixte devrait fonctionner",
		},
		{
			name:        "OR uppercase",
			input:       `rule test4 : {p:Person} / p.age < 18 OR p.age > 65 ==> action1()`,
			shouldParse: true,
			description: "OR en majuscules devrait fonctionner",
		},
		{
			name:        "or lowercase",
			input:       `rule test5 : {p:Person} / p.age < 18 or p.age > 65 ==> action1()`,
			shouldParse: true,
			description: "or en minuscules devrait fonctionner",
		},
		{
			name:        "Or mixed case",
			input:       `rule test6 : {p:Person} / p.age < 18 Or p.age > 65 ==> action1()`,
			shouldParse: true,
			description: "Or en casse mixte devrait fonctionner",
		},
		// NOT constraint
		{
			name:        "NOT uppercase",
			input:       `rule test7 : {p:Person} / NOT(p.age > 100) ==> action1()`,
			shouldParse: true,
			description: "NOT en majuscules devrait fonctionner",
		},
		{
			name:        "not lowercase",
			input:       `rule test8 : {p:Person} / not(p.age > 100) ==> action1()`,
			shouldParse: true,
			description: "not en minuscules devrait fonctionner",
		},
		{
			name:        "Not mixed case",
			input:       `rule test9 : {p:Person} / Not(p.age > 100) ==> action1()`,
			shouldParse: true,
			description: "Not en casse mixte devrait fonctionner",
		},
		// EXISTS constraint
		{
			name:        "EXISTS uppercase",
			input:       `rule test10 : {p:Person} / EXISTS(c:Child / c.parent == p.id) ==> action1()`,
			shouldParse: true,
			description: "EXISTS en majuscules devrait fonctionner",
		},
		{
			name:        "exists lowercase",
			input:       `rule test11 : {p:Person} / exists(c:Child / c.parent == p.id) ==> action1()`,
			shouldParse: true,
			description: "exists en minuscules devrait fonctionner",
		},
		{
			name:        "Exists mixed case",
			input:       `rule test12 : {p:Person} / Exists(c:Child / c.parent == p.id) ==> action1()`,
			shouldParse: true,
			description: "Exists en casse mixte devrait fonctionner",
		},
		// Fonctions d'agrégation (syntaxe correcte sans conditions dans les accolades)
		{
			name:        "AVG uppercase in aggregation variable",
			input:       `rule test13 : {s:Sale, total:AVG(s.amount)} / total > 100 ==> action1()`,
			shouldParse: true,
			description: "AVG en majuscules devrait fonctionner",
		},
		{
			name:        "avg lowercase in aggregation variable",
			input:       `rule test14 : {s:Sale, total:avg(s.amount)} / total > 100 ==> action1()`,
			shouldParse: true,
			description: "avg en minuscules devrait fonctionner",
		},
		{
			name:        "Avg mixed case in aggregation variable",
			input:       `rule test15 : {s:Sale, total:Avg(s.amount)} / total > 100 ==> action1()`,
			shouldParse: true,
			description: "Avg en casse mixte devrait fonctionner",
		},
		{
			name:        "COUNT uppercase in constraint",
			input:       `rule test16 : {p:Person} / COUNT(c:Child / c.parent == p.id) > 2 ==> action1()`,
			shouldParse: true,
			description: "COUNT en majuscules devrait fonctionner",
		},
		{
			name:        "count lowercase in constraint",
			input:       `rule test17 : {p:Person} / count(c:Child / c.parent == p.id) > 2 ==> action1()`,
			shouldParse: true,
			description: "count en minuscules devrait fonctionner",
		},
		{
			name:        "Count mixed case in constraint",
			input:       `rule test18 : {p:Person} / Count(c:Child / c.parent == p.id) > 2 ==> action1()`,
			shouldParse: true,
			description: "Count en casse mixte devrait fonctionner",
		},
		{
			name:        "SUM uppercase in aggregation variable",
			input:       `rule test19 : {s:Sale, total:SUM(s.amount)} / total > 1000 ==> action1()`,
			shouldParse: true,
			description: "SUM en majuscules devrait fonctionner",
		},
		{
			name:        "sum lowercase in aggregation variable",
			input:       `rule test20 : {s:Sale, total:sum(s.amount)} / total > 1000 ==> action1()`,
			shouldParse: true,
			description: "sum en minuscules devrait fonctionner",
		},
		{
			name:        "Sum mixed case in aggregation variable",
			input:       `rule test21 : {s:Sale, total:Sum(s.amount)} / total > 1000 ==> action1()`,
			shouldParse: true,
			description: "Sum en casse mixte devrait fonctionner",
		},
		{
			name:        "MIN uppercase in aggregation variable",
			input:       `rule test22 : {s:Sale, min:MIN(s.amount)} / min < 10 ==> action1()`,
			shouldParse: true,
			description: "MIN en majuscules devrait fonctionner",
		},
		{
			name:        "min lowercase in aggregation variable",
			input:       `rule test23 : {s:Sale, min:min(s.amount)} / min < 10 ==> action1()`,
			shouldParse: true,
			description: "min en minuscules devrait fonctionner",
		},
		{
			name:        "Max mixed case in aggregation variable",
			input:       `rule test24 : {s:Sale, max:Max(s.amount)} / max > 1000 ==> action1()`,
			shouldParse: true,
			description: "Max en casse mixte devrait fonctionner",
		},
		// Opérateurs de comparaison
		{
			name:        "IN uppercase",
			input:       `rule test25 : {p:Person} / p.status IN ["active", "pending"] ==> action1()`,
			shouldParse: true,
			description: "IN en majuscules devrait fonctionner",
		},
		{
			name:        "in lowercase",
			input:       `rule test26 : {p:Person} / p.status in ["active", "pending"] ==> action1()`,
			shouldParse: true,
			description: "in en minuscules devrait fonctionner",
		},
		{
			name:        "In mixed case",
			input:       `rule test27 : {p:Person} / p.status In ["active", "pending"] ==> action1()`,
			shouldParse: true,
			description: "In en casse mixte devrait fonctionner",
		},
		{
			name:        "LIKE uppercase",
			input:       `rule test28 : {p:Person} / p.name LIKE "John%" ==> action1()`,
			shouldParse: true,
			description: "LIKE en majuscules devrait fonctionner",
		},
		{
			name:        "like lowercase",
			input:       `rule test29 : {p:Person} / p.name like "John%" ==> action1()`,
			shouldParse: true,
			description: "like en minuscules devrait fonctionner",
		},
		{
			name:        "Like mixed case",
			input:       `rule test30 : {p:Person} / p.name Like "John%" ==> action1()`,
			shouldParse: true,
			description: "Like en casse mixte devrait fonctionner",
		},
		{
			name:        "MATCHES uppercase",
			input:       `rule test31 : {p:Person} / p.email MATCHES ".*@example.com" ==> action1()`,
			shouldParse: true,
			description: "MATCHES en majuscules devrait fonctionner",
		},
		{
			name:        "matches lowercase",
			input:       `rule test32 : {p:Person} / p.email matches ".*@example.com" ==> action1()`,
			shouldParse: true,
			description: "matches en minuscules devrait fonctionner",
		},
		{
			name:        "Matches mixed case",
			input:       `rule test33 : {p:Person} / p.email Matches ".*@example.com" ==> action1()`,
			shouldParse: true,
			description: "Matches en casse mixte devrait fonctionner",
		},
		{
			name:        "CONTAINS uppercase",
			input:       `rule test34 : {p:Person} / p.tags CONTAINS "vip" ==> action1()`,
			shouldParse: true,
			description: "CONTAINS en majuscules devrait fonctionner",
		},
		{
			name:        "contains lowercase",
			input:       `rule test35 : {p:Person} / p.tags contains "vip" ==> action1()`,
			shouldParse: true,
			description: "contains en minuscules devrait fonctionner",
		},
		{
			name:        "Contains mixed case",
			input:       `rule test36 : {p:Person} / p.tags Contains "vip" ==> action1()`,
			shouldParse: true,
			description: "Contains en casse mixte devrait fonctionner",
		},
		// Fonctions de manipulation
		{
			name:        "LENGTH uppercase",
			input:       `rule test37 : {p:Person} / LENGTH(p.name) > 5 ==> action1()`,
			shouldParse: true,
			description: "LENGTH en majuscules devrait fonctionner",
		},
		{
			name:        "length lowercase",
			input:       `rule test38 : {p:Person} / length(p.name) > 5 ==> action1()`,
			shouldParse: true,
			description: "length en minuscules devrait fonctionner",
		},
		{
			name:        "Length mixed case",
			input:       `rule test39 : {p:Person} / Length(p.name) > 5 ==> action1()`,
			shouldParse: true,
			description: "Length en casse mixte devrait fonctionner",
		},
		{
			name:        "UPPER/LOWER/TRIM variations",
			input:       `rule test40 : {p:Person} / upper(p.name) == lower(p.status) AND trim(p.email) != "" ==> action1()`,
			shouldParse: true,
			description: "Combinaison de fonctions en minuscules devrait fonctionner",
		},
		{
			name:        "ABS/ROUND/FLOOR/CEIL variations",
			input:       `rule test41 : {p:Person} / abs(p.balance) > round(p.limit) ==> action1()`,
			shouldParse: true,
			description: "Fonctions mathématiques en minuscules devraient fonctionner",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("test", []byte(tt.input))
			if tt.shouldParse && err != nil {
				t.Errorf("%s: parsing devrait réussir mais a échoué: %v", tt.description, err)
			}
			if !tt.shouldParse && err == nil {
				t.Errorf("%s: parsing devrait échouer mais a réussi", tt.description)
			}
		})
	}
}

// TestBug_CaseInsensitiveKeywords_AccumulateConstraints teste les fonctions d'accumulation dans AccumulateConstraint
func TestBug_CaseInsensitiveKeywords_AccumulateConstraints(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldParse bool
		description string
	}{
		// SUM dans AccumulateConstraint
		{
			name:        "SUM uppercase in AccumulateConstraint",
			input:       `rule test1 : {o:Order} / SUM(x:Order / x.id == o.id ; x.amount) > 1000 ==> alert()`,
			shouldParse: true,
			description: "SUM en majuscules dans AccumulateConstraint devrait fonctionner",
		},
		{
			name:        "sum lowercase in AccumulateConstraint",
			input:       `rule test2 : {o:Order} / sum(x:Order / x.id == o.id ; x.amount) > 1000 ==> alert()`,
			shouldParse: true,
			description: "sum en minuscules dans AccumulateConstraint devrait fonctionner",
		},
		{
			name:        "Sum capitalized in AccumulateConstraint",
			input:       `rule test3 : {o:Order} / Sum(x:Order / x.id == o.id ; x.amount) > 1000 ==> alert()`,
			shouldParse: true,
			description: "Sum capitalisé dans AccumulateConstraint devrait fonctionner",
		},
		// AVG dans AccumulateConstraint
		{
			name:        "AVG uppercase in AccumulateConstraint",
			input:       `rule test4 : {m:Metric} / AVG(x:Metric / x.sensor == m.sensor ; x.value) > 50 ==> notify()`,
			shouldParse: true,
			description: "AVG en majuscules dans AccumulateConstraint devrait fonctionner",
		},
		{
			name:        "avg lowercase in AccumulateConstraint",
			input:       `rule test5 : {m:Metric} / avg(x:Metric / x.sensor == m.sensor ; x.value) > 50 ==> notify()`,
			shouldParse: true,
			description: "avg en minuscules dans AccumulateConstraint devrait fonctionner",
		},
		{
			name:        "Avg capitalized in AccumulateConstraint",
			input:       `rule test6 : {m:Metric} / Avg(x:Metric / x.sensor == m.sensor ; x.value) > 50 ==> notify()`,
			shouldParse: true,
			description: "Avg capitalisé dans AccumulateConstraint devrait fonctionner",
		},
		// MIN dans AccumulateConstraint
		{
			name:        "MIN uppercase in AccumulateConstraint",
			input:       `rule test7 : {t:Temperature} / MIN(x:Temperature / x.location == t.location ; x.temp) < 0 ==> alert()`,
			shouldParse: true,
			description: "MIN en majuscules dans AccumulateConstraint devrait fonctionner",
		},
		{
			name:        "min lowercase in AccumulateConstraint",
			input:       `rule test8 : {t:Temperature} / min(x:Temperature / x.location == t.location ; x.temp) < 0 ==> alert()`,
			shouldParse: true,
			description: "min en minuscules dans AccumulateConstraint devrait fonctionner",
		},
		{
			name:        "Min capitalized in AccumulateConstraint",
			input:       `rule test9 : {t:Temperature} / Min(x:Temperature / x.location == t.location ; x.temp) < 0 ==> alert()`,
			shouldParse: true,
			description: "Min capitalisé dans AccumulateConstraint devrait fonctionner",
		},
		// MAX dans AccumulateConstraint
		{
			name:        "MAX uppercase in AccumulateConstraint",
			input:       `rule test10 : {p:Pressure} / MAX(x:Pressure / x.device == p.device ; x.value) > 100 ==> warn()`,
			shouldParse: true,
			description: "MAX en majuscules dans AccumulateConstraint devrait fonctionner",
		},
		{
			name:        "max lowercase in AccumulateConstraint",
			input:       `rule test11 : {p:Pressure} / max(x:Pressure / x.device == p.device ; x.value) > 100 ==> warn()`,
			shouldParse: true,
			description: "max en minuscules dans AccumulateConstraint devrait fonctionner",
		},
		{
			name:        "Max capitalized in AccumulateConstraint",
			input:       `rule test12 : {p:Pressure} / Max(x:Pressure / x.device == p.device ; x.value) > 100 ==> warn()`,
			shouldParse: true,
			description: "Max capitalisé dans AccumulateConstraint devrait fonctionner",
		},
		// COUNT dans AccumulateConstraint (sans champ)
		{
			name:        "COUNT uppercase in AccumulateConstraint",
			input:       `rule test13 : {e:Event} / COUNT(x:Event / x.type == e.type) > 10 ==> escalate()`,
			shouldParse: true,
			description: "COUNT en majuscules dans AccumulateConstraint devrait fonctionner",
		},
		{
			name:        "count lowercase in AccumulateConstraint",
			input:       `rule test14 : {e:Event} / count(x:Event / x.type == e.type) > 10 ==> escalate()`,
			shouldParse: true,
			description: "count en minuscules dans AccumulateConstraint devrait fonctionner",
		},
		{
			name:        "Count capitalized in AccumulateConstraint",
			input:       `rule test15 : {e:Event} / Count(x:Event / x.type == e.type) > 10 ==> escalate()`,
			shouldParse: true,
			description: "Count capitalisé dans AccumulateConstraint devrait fonctionner",
		},
		// Combinaison avec AND/OR dans la condition
		{
			name:        "Mixed case SUM with and in condition",
			input:       `rule test16 : {o:Order} / sum(x:Order / x.id == o.id and x.valid == true ; x.amount) > 1000 ==> alert()`,
			shouldParse: true,
			description: "sum minuscule avec and minuscule dans la condition devrait fonctionner",
		},
		{
			name:        "Mixed case AVG with Or in condition",
			input:       `rule test17 : {m:Metric} / Avg(x:Metric / x.sensor == m.sensor Or x.backup == true ; x.value) > 50 ==> notify()`,
			shouldParse: true,
			description: "Avg capitalisé avec Or capitalisé dans la condition devrait fonctionner",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("test", []byte(tt.input))
			if tt.shouldParse && err != nil {
				t.Errorf("%s: parsing devrait réussir mais a échoué: %v", tt.description, err)
			}
			if !tt.shouldParse && err == nil {
				t.Errorf("%s: parsing devrait échouer mais a réussi", tt.description)
			}
		})
	}
}

// TestBug_CaseInsensitiveKeywords_MixedCombinations teste des combinaisons complexes
func TestBug_CaseInsensitiveKeywords_MixedCombinations(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Mixed case in complex rule",
			input: `rule complex : {p:Person} / p.age > 18 and (p.status in ["active", "pending"] or length(p.name) > 5) ==> action1()`,
		},
		{
			name:  "Mixed aggregation and logical ops",
			input: `rule agg : {s:Sale, total:sum(s.amount)} / total > 1000 and total < 10000 ==> notify()`,
		},
		{
			name:  "Mixed NOT and EXISTS",
			input: `rule exists_not : {p:Person} / not(exists(c:Child / c.parent == p.id)) ==> alert()`,
		},
		{
			name:  "All lowercase complex",
			input: `rule all_lower : {p:Person} / p.age > 18 and (not(p.blocked == true) or exists(v:Voucher / v.user == p.id)) ==> process()`,
		},
		{
			name:  "Mixed case everything",
			input: `rule MixedCase : {s:Sale, avg:Avg(s.amount)} / avg > 100 And avg < 500 Or Count(c:Customer / c.type == "vip") > 10 ==> action()`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("test", []byte(tt.input))
			if err != nil {
				t.Errorf("Parsing devrait réussir pour %s, mais a échoué: %v", tt.name, err)
			}
		})
	}
}

// TestBug_CaseInsensitiveKeywords_InvalidCases teste que les formes de casse invalides sont rejetées
func TestBug_CaseInsensitiveKeywords_InvalidCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		description string
	}{
		{
			name:        "aNd invalid case",
			input:       `rule test : {p:Person} / p.age > 18 aNd p.age < 65 ==> action1()`,
			description: "aNd avec casse invalide devrait échouer",
		},
		{
			name:        "AnD invalid case",
			input:       `rule test : {p:Person} / p.age > 18 AnD p.age < 65 ==> action1()`,
			description: "AnD avec casse invalide devrait échouer",
		},
		{
			name:        "oR invalid case",
			input:       `rule test : {p:Person} / p.age < 18 oR p.age > 65 ==> action1()`,
			description: "oR avec casse invalide devrait échouer",
		},
		{
			name:        "nOt invalid case",
			input:       `rule test : {p:Person} / nOt(p.age > 100) ==> action1()`,
			description: "nOt avec casse invalide devrait échouer",
		},
		{
			name:        "eXiStS invalid case",
			input:       `rule test : {p:Person} / eXiStS(c:Child / c.parent == p.id) ==> action1()`,
			description: "eXiStS avec casse invalide devrait échouer",
		},
		{
			name:        "aVg invalid case",
			input:       `rule test : {s:Sale, total:aVg(s.amount)} / total > 100 ==> action1()`,
			description: "aVg avec casse invalide devrait échouer",
		},
		{
			name:        "CoUnT invalid case",
			input:       `rule test : {p:Person} / CoUnT(c:Child / c.parent == p.id) > 2 ==> action1()`,
			description: "CoUnT avec casse invalide devrait échouer",
		},
		{
			name:        "iN invalid case",
			input:       `rule test : {p:Person} / p.status iN ["active", "pending"] ==> action1()`,
			description: "iN avec casse invalide devrait échouer",
		},
		{
			name:        "LiKe invalid case",
			input:       `rule test : {p:Person} / p.name LiKe "John%" ==> action1()`,
			description: "LiKe avec casse invalide devrait échouer",
		},
		{
			name:        "lEnGtH invalid case",
			input:       `rule test : {p:Person} / lEnGtH(p.name) > 5 ==> action1()`,
			description: "lEnGtH avec casse invalide devrait échouer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse("test", []byte(tt.input))
			if err == nil {
				t.Errorf("%s: le parsing devrait échouer mais a réussi", tt.description)
			}
		})
	}
}
