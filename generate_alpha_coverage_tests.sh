#!/bin/bash

# 🎯 TESTS DE COUVERTURE OPTIMALE NOEUDS ALPHA
# ============================================
# Ce script génère des tests complets pour valider le traitement sémantique
# des conditions positives et négatives par les nœuds Alpha

echo "🔬 GÉNÉRATION DES TESTS DE COUVERTURE ALPHA NODES"
echo "================================================="

# Configuration
WORKSPACE_DIR="/home/resinsec/dev/tsd"
TEST_OUTPUT_DIR="$WORKSPACE_DIR/alpha_coverage_tests"
RESULTS_FILE="$WORKSPACE_DIR/ALPHA_NODES_COVERAGE_COMPLETE_RESULTS.md"

# Créer le répertoire de tests
mkdir -p "$TEST_OUTPUT_DIR"

echo "📁 Répertoire de tests: $TEST_OUTPUT_DIR"
echo "📄 Fichier résultats: $RESULTS_FILE"
echo ""

# Fonction pour créer un test
create_alpha_test() {
    local test_name="$1"
    local description="$2" 
    local constraint_content="$3"
    local facts_content="$4"
    
    local constraint_file="$TEST_OUTPUT_DIR/${test_name}.constraint"
    local facts_file="$TEST_OUTPUT_DIR/${test_name}.facts"
    
    echo "// $description" > "$constraint_file"
    echo "$constraint_content" >> "$constraint_file"
    
    echo "$facts_content" > "$facts_file"
    
    echo "✅ Test créé: $test_name"
}

echo "🧪 CRÉATION DES TESTS DE COUVERTURE"
echo "==================================="

# Test 1: Égalité positive vs négative
create_alpha_test "alpha_equality_positive" \
"Test égalité positive simple" \
"type Person : <id: string, age: number, status: string>

{p: Person} / p.age == 25 ==> age_is_twenty_five(p.id, p.age)" \
"Person[id=P001, age=25, status=active]
Person[id=P002, age=30, status=active]
Person[id=P003, age=25, status=inactive]"

create_alpha_test "alpha_equality_negative" \
"Test égalité négative simple" \
"type Person : <id: string, age: number, status: string>

{p: Person} / NOT(p.age == 25) ==> age_is_not_twenty_five(p.id, p.age)" \
"Person[id=P001, age=25, status=active]
Person[id=P002, age=30, status=active]  
Person[id=P003, age=25, status=inactive]"

# Test 2: Comparaisons numériques positive vs négative  
create_alpha_test "alpha_comparison_positive" \
"Test comparaison numérique positive" \
"type Product : <id: string, price: number, category: string>

{prod: Product} / prod.price > 100 ==> expensive_product(prod.id, prod.price)" \
"Product[id=PROD001, price=150, category=electronics]
Product[id=PROD002, price=50, category=books]
Product[id=PROD003, price=200, category=electronics]"

create_alpha_test "alpha_comparison_negative" \
"Test comparaison numérique négative" \
"type Product : <id: string, price: number, category: string>

{prod: Product} / NOT(prod.price > 100) ==> affordable_product(prod.id, prod.price)" \
"Product[id=PROD001, price=150, category=electronics]
Product[id=PROD002, price=50, category=books]
Product[id=PROD003, price=200, category=electronics]"

# Test 3: Conditions string positive vs négative
create_alpha_test "alpha_string_positive" \
"Test condition string positive" \
"type User : <id: string, name: string, role: string>

{u: User} / u.role == \"admin\" ==> admin_user_found(u.id, u.name)" \
"User[id=U001, name=Alice, role=admin]
User[id=U002, name=Bob, role=user]
User[id=U003, name=Charlie, role=admin]"

create_alpha_test "alpha_string_negative" \
"Test condition string négative" \
"type User : <id: string, name: string, role: string>

{u: User} / NOT(u.role == \"admin\") ==> non_admin_user_found(u.id, u.name)" \
"User[id=U001, name=Alice, role=admin]
User[id=U002, name=Bob, role=user]
User[id=U003, name=Charlie, role=admin]"

# Test 4: Conditions booléennes positive vs négative
create_alpha_test "alpha_boolean_positive" \
"Test condition booléenne positive" \
"type Account : <id: string, balance: number, active: bool>

{a: Account} / a.active == true ==> active_account_found(a.id, a.balance)" \
"Account[id=ACC001, balance=1000, active=true]
Account[id=ACC002, balance=500, active=false] 
Account[id=ACC003, balance=2000, active=true]"

create_alpha_test "alpha_boolean_negative" \
"Test condition booléenne négative" \
"type Account : <id: string, balance: number, active: bool>

{a: Account} / NOT(a.active == true) ==> inactive_account_found(a.id, a.balance)" \
"Account[id=ACC001, balance=1000, active=true]
Account[id=ACC002, balance=500, active=false]
Account[id=ACC003, balance=2000, active=true]"

# Test 5: Inégalité positive vs négative
create_alpha_test "alpha_inequality_positive" \
"Test inégalité positive" \
"type Order : <id: string, total: number, status: string>

{o: Order} / o.status != \"cancelled\" ==> valid_order_found(o.id, o.total)" \
"Order[id=ORD001, total=100, status=pending]
Order[id=ORD002, total=200, status=cancelled]
Order[id=ORD003, total=300, status=completed]"

create_alpha_test "alpha_inequality_negative" \
"Test inégalité négative" \
"type Order : <id: string, total: number, status: string>

{o: Order} / NOT(o.status != \"cancelled\") ==> cancelled_order_found(o.id, o.total)" \
"Order[id=ORD001, total=100, status=pending]
Order[id=ORD002, total=200, status=cancelled]
Order[id=ORD003, total=300, status=completed]"

echo ""
echo "📋 TESTS CRÉÉS:"
echo "==============="
ls -la "$TEST_OUTPUT_DIR"/*.constraint | wc -l | xargs echo "- Fichiers .constraint:"
ls -la "$TEST_OUTPUT_DIR"/*.facts | wc -l | xargs echo "- Fichiers .facts:"

echo ""
echo "🚀 EXÉCUTION DES TESTS ET GÉNÉRATION DU RAPPORT"
echo "==============================================="

# Lancer le script Go de test complet
cd "$WORKSPACE_DIR"
go run alpha_coverage_test_runner.go