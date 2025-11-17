#!/bin/bash

# Script pour générer tous les tests de couverture Alpha complémentaires
# Couvre TOUS les opérateurs et fonctions définis dans la grammaire

ALPHA_DIR="alpha_coverage_tests_extended"

# Créer le répertoire pour les tests étendus
mkdir -p "$ALPHA_DIR"

echo "🔧 GÉNÉRATION TESTS ALPHA ÉTENDUS"
echo "================================="

# =============================================================================
# 1. OPÉRATEUR = (égalité alternative)
# =============================================================================

cat > "$ALPHA_DIR/alpha_equal_sign_positive.constraint" << 'EOF'
// Test opérateur = (égalité alternative)
type Customer : <id: string, tier: string, points: number>

{c: Customer} / c.tier = "gold" ==> gold_customer_found(c.id, c.points)
EOF

cat > "$ALPHA_DIR/alpha_equal_sign_positive.facts" << 'EOF'
Customer{id: "C001", tier: "gold", points: 5000}
Customer{id: "C002", tier: "silver", points: 2000}
Customer{id: "C003", tier: "gold", points: 7500}
EOF

cat > "$ALPHA_DIR/alpha_equal_sign_negative.constraint" << 'EOF'
// Test opérateur = (égalité alternative) - négation
type Customer : <id: string, tier: string, points: number>

{c: Customer} / NOT(c.tier = "gold") ==> non_gold_customer_found(c.id, c.tier)
EOF

cat > "$ALPHA_DIR/alpha_equal_sign_negative.facts" << 'EOF'
Customer{id: "C001", tier: "gold", points: 5000}
Customer{id: "C002", tier: "silver", points: 2000}
Customer{id: "C003", tier: "bronze", points: 1000}
EOF

# =============================================================================
# 2. OPÉRATEUR IN (appartenance à un ensemble)
# =============================================================================

cat > "$ALPHA_DIR/alpha_in_positive.constraint" << 'EOF'
// Test opérateur IN (appartenance)
type Status : <id: string, state: string, priority: number>

{s: Status} / s.state IN ["active", "pending", "review"] ==> valid_state_found(s.id, s.state)
EOF

cat > "$ALPHA_DIR/alpha_in_positive.facts" << 'EOF'
Status{id: "S001", state: "active", priority: 1}
Status{id: "S002", state: "inactive", priority: 3}
Status{id: "S003", state: "pending", priority: 2}
EOF

cat > "$ALPHA_DIR/alpha_in_negative.constraint" << 'EOF'
// Test opérateur IN (appartenance) - négation
type Status : <id: string, state: string, priority: number>

{s: Status} / NOT(s.state IN ["active", "pending"]) ==> invalid_state_found(s.id, s.state)
EOF

cat > "$ALPHA_DIR/alpha_in_negative.facts" << 'EOF'
Status{id: "S001", state: "active", priority: 1}
Status{id: "S002", state: "inactive", priority: 3}
Status{id: "S003", state: "archived", priority: 5}
EOF

# =============================================================================
# 3. OPÉRATEUR LIKE (correspondance de motif)
# =============================================================================

cat > "$ALPHA_DIR/alpha_like_positive.constraint" << 'EOF'
// Test opérateur LIKE (motif)
type Email : <id: string, address: string, verified: bool>

{e: Email} / e.address LIKE "%@company.com" ==> company_email_found(e.id, e.address)
EOF

cat > "$ALPHA_DIR/alpha_like_positive.facts" << 'EOF'
Email{id: "E001", address: "john@company.com", verified: true}
Email{id: "E002", address: "jane@external.org", verified: false}
Email{id: "E003", address: "admin@company.com", verified: true}
EOF

cat > "$ALPHA_DIR/alpha_like_negative.constraint" << 'EOF'
// Test opérateur LIKE (motif) - négation
type Email : <id: string, address: string, verified: bool>

{e: Email} / NOT(e.address LIKE "%@company.com") ==> external_email_found(e.id, e.address)
EOF

cat > "$ALPHA_DIR/alpha_like_negative.facts" << 'EOF'
Email{id: "E001", address: "john@company.com", verified: true}
Email{id: "E002", address: "jane@external.org", verified: false}
Email{id: "E003", address: "user@other.net", verified: true}
EOF

# =============================================================================
# 4. OPÉRATEUR MATCHES (regex)
# =============================================================================

cat > "$ALPHA_DIR/alpha_matches_positive.constraint" << 'EOF'
// Test opérateur MATCHES (regex)
type Code : <id: string, value: string, active: bool>

{c: Code} / c.value MATCHES "CODE[0-9]+" ==> valid_code_found(c.id, c.value)
EOF

cat > "$ALPHA_DIR/alpha_matches_positive.facts" << 'EOF'
Code{id: "C001", value: "CODE123", active: true}
Code{id: "C002", value: "INVALID", active: false}
Code{id: "C003", value: "CODE999", active: true}
EOF

cat > "$ALPHA_DIR/alpha_matches_negative.constraint" << 'EOF'
// Test opérateur MATCHES (regex) - négation
type Code : <id: string, value: string, active: bool>

{c: Code} / NOT(c.value MATCHES "CODE[0-9]+") ==> invalid_code_found(c.id, c.value)
EOF

cat > "$ALPHA_DIR/alpha_matches_negative.facts" << 'EOF'
Code{id: "C001", value: "CODE123", active: true}
Code{id: "C002", value: "INVALID", active: false}
Code{id: "C003", value: "BADFORMAT", active: true}
EOF

# =============================================================================
# 5. OPÉRATEUR CONTAINS (contenance)
# =============================================================================

cat > "$ALPHA_DIR/alpha_contains_positive.constraint" << 'EOF'
// Test opérateur CONTAINS (contenance)
type Message : <id: string, content: string, urgent: bool>

{m: Message} / m.content CONTAINS "urgent" ==> urgent_message_found(m.id, m.content)
EOF

cat > "$ALPHA_DIR/alpha_contains_positive.facts" << 'EOF'
Message{id: "M001", content: "This is urgent please respond", urgent: true}
Message{id: "M002", content: "Regular message content", urgent: false}
Message{id: "M003", content: "Very urgent matter!", urgent: true}
EOF

cat > "$ALPHA_DIR/alpha_contains_negative.constraint" << 'EOF'
// Test opérateur CONTAINS (contenance) - négation
type Message : <id: string, content: string, urgent: bool>

{m: Message} / NOT(m.content CONTAINS "urgent") ==> normal_message_found(m.id, m.content)
EOF

cat > "$ALPHA_DIR/alpha_contains_negative.facts" << 'EOF'
Message{id: "M001", content: "This is urgent please respond", urgent: true}
Message{id: "M002", content: "Regular message content", urgent: false}
Message{id: "M003", content: "Simple notification", urgent: false}
EOF

# =============================================================================
# 6. FONCTION LENGTH()
# =============================================================================

cat > "$ALPHA_DIR/alpha_length_positive.constraint" << 'EOF'
// Test fonction LENGTH()
type Password : <id: string, value: string, secure: bool>

{p: Password} / LENGTH(p.value) >= 8 ==> secure_password_found(p.id, p.value)
EOF

cat > "$ALPHA_DIR/alpha_length_positive.facts" << 'EOF'
Password{id: "P001", value: "password123", secure: true}
Password{id: "P002", value: "123", secure: false}
Password{id: "P003", value: "verysecurepass", secure: true}
EOF

cat > "$ALPHA_DIR/alpha_length_negative.constraint" << 'EOF'
// Test fonction LENGTH() - négation
type Password : <id: string, value: string, secure: bool>

{p: Password} / NOT(LENGTH(p.value) >= 8) ==> weak_password_found(p.id, p.value)
EOF

cat > "$ALPHA_DIR/alpha_length_negative.facts" << 'EOF'
Password{id: "P001", value: "password123", secure: true}
Password{id: "P002", value: "123", secure: false}
Password{id: "P003", value: "pass", secure: false}
EOF

# =============================================================================
# 7. FONCTION ABS() (valeur absolue)
# =============================================================================

cat > "$ALPHA_DIR/alpha_abs_positive.constraint" << 'EOF'
// Test fonction ABS() (valeur absolue)
type Balance : <id: string, amount: number, type: string>

{b: Balance} / ABS(b.amount) > 100 ==> significant_balance_found(b.id, b.amount)
EOF

cat > "$ALPHA_DIR/alpha_abs_positive.facts" << 'EOF'
Balance{id: "B001", amount: 150.0, type: "credit"}
Balance{id: "B002", amount: -200.0, type: "debit"}
Balance{id: "B003", amount: 50.0, type: "credit"}
EOF

cat > "$ALPHA_DIR/alpha_abs_negative.constraint" << 'EOF'
// Test fonction ABS() (valeur absolue) - négation
type Balance : <id: string, amount: number, type: string>

{b: Balance} / NOT(ABS(b.amount) > 100) ==> small_balance_found(b.id, b.amount)
EOF

cat > "$ALPHA_DIR/alpha_abs_negative.facts" << 'EOF'
Balance{id: "B001", amount: 150.0, type: "credit"}
Balance{id: "B002", amount: -25.0, type: "debit"}
Balance{id: "B003", amount: 75.0, type: "credit"}
EOF

# =============================================================================
# 8. FONCTION UPPER() (majuscules)
# =============================================================================

cat > "$ALPHA_DIR/alpha_upper_positive.constraint" << 'EOF'
// Test fonction UPPER() (majuscules)
type Department : <id: string, name: string, active: bool>

{d: Department} / UPPER(d.name) == "FINANCE" ==> finance_dept_found(d.id, d.name)
EOF

cat > "$ALPHA_DIR/alpha_upper_positive.facts" << 'EOF'
Department{id: "D001", name: "finance", active: true}
Department{id: "D002", name: "IT", active: true}
Department{id: "D003", name: "Finance", active: true}
EOF

cat > "$ALPHA_DIR/alpha_upper_negative.constraint" << 'EOF'
// Test fonction UPPER() (majuscules) - négation
type Department : <id: string, name: string, active: bool>

{d: Department} / NOT(UPPER(d.name) == "FINANCE") ==> non_finance_dept_found(d.id, d.name)
EOF

cat > "$ALPHA_DIR/alpha_upper_negative.facts" << 'EOF'
Department{id: "D001", name: "finance", active: true}
Department{id: "D002", name: "IT", active: true}
Department{id: "D003", name: "HR", active: true}
EOF

echo "📊 RÉSUMÉ GÉNÉRATION"
echo "=================="
echo "✅ 16 nouveaux tests Alpha générés (8 opérateurs × 2 polarités)"
echo "📂 Répertoire: $ALPHA_DIR"
echo ""
echo "🧪 NOUVEAUX OPÉRATEURS/FONCTIONS TESTÉS:"
echo "  • Opérateur = (égalité alternative)"
echo "  • Opérateur IN (appartenance)"
echo "  • Opérateur LIKE (motif)"
echo "  • Opérateur MATCHES (regex)"
echo "  • Opérateur CONTAINS (contenance)"
echo "  • Fonction LENGTH() (longueur)"
echo "  • Fonction ABS() (valeur absolue)"
echo "  • Fonction UPPER() (majuscules)"
echo ""
echo "🎯 COUVERTURE TOTALE: $(ls "$ALPHA_DIR"/*.constraint | wc -l) tests"

# Compter tous les tests (originaux + étendus)
TOTAL_ORIGINAL=$(ls alpha_coverage_tests/*.constraint 2>/dev/null | wc -l)
TOTAL_EXTENDED=$(ls "$ALPHA_DIR"/*.constraint 2>/dev/null | wc -l)
TOTAL_ALL=$((TOTAL_ORIGINAL + TOTAL_EXTENDED))

echo "📋 COMPARAISON:"
echo "  • Tests originaux: $TOTAL_ORIGINAL"
echo "  • Tests étendus: $TOTAL_EXTENDED"
echo "  • TOTAL: $TOTAL_ALL tests de couverture Alpha"