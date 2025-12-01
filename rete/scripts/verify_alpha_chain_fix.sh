#!/bin/bash

# Alpha Chain Integration Fix - Verification Script
# This script verifies that the 4 previously failing Alpha chain tests now pass

set -e

echo "╔════════════════════════════════════════════════════════════════╗"
echo "║   Alpha Chain Integration Fix - Verification Script            ║"
echo "╚════════════════════════════════════════════════════════════════╝"
echo ""

cd "$(dirname "$0")/.."

echo "📋 Testing the 4 previously failing Alpha chain tests..."
echo ""

# Test 1
echo "1️⃣  Testing: TestAlphaChain_TwoRules_SameConditions_DifferentOrder"
if go test -run "^TestAlphaChain_TwoRules_SameConditions_DifferentOrder$" . > /dev/null 2>&1; then
    echo "   ✅ PASS"
else
    echo "   ❌ FAIL"
    exit 1
fi

# Test 2
echo "2️⃣  Testing: TestAlphaChain_FactPropagation_ThroughChain"
if go test -run "^TestAlphaChain_FactPropagation_ThroughChain$" . > /dev/null 2>&1; then
    echo "   ✅ PASS"
else
    echo "   ❌ FAIL"
    exit 1
fi

# Test 3
echo "3️⃣  Testing: TestAlphaChain_NetworkStats_Accurate"
if go test -run "^TestAlphaChain_NetworkStats_Accurate$" . > /dev/null 2>&1; then
    echo "   ✅ PASS"
else
    echo "   ❌ FAIL"
    exit 1
fi

# Test 4
echo "4️⃣  Testing: TestAlphaChain_MixedConditions_ComplexSharing"
if go test -run "^TestAlphaChain_MixedConditions_ComplexSharing$" . > /dev/null 2>&1; then
    echo "   ✅ PASS"
else
    echo "   ❌ FAIL"
    exit 1
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🎯 Running all Alpha chain tests..."
if go test -run "TestAlphaChain" . > /dev/null 2>&1; then
    TEST_COUNT=$(go test -run "TestAlphaChain" . 2>&1 | grep -c "PASS: TestAlphaChain" || echo "0")
    echo "   ✅ All Alpha chain tests passing ($TEST_COUNT tests)"
else
    echo "   ❌ Some Alpha chain tests failed"
    exit 1
fi

echo ""
echo "🧪 Running AlphaRuleBuilder unit tests..."
if go test -run "TestAlphaRuleBuilder" . > /dev/null 2>&1; then
    TEST_COUNT=$(go test -run "TestAlphaRuleBuilder" . 2>&1 | grep -c "PASS: TestAlphaRuleBuilder" || echo "0")
    echo "   ✅ All AlphaRuleBuilder tests passing ($TEST_COUNT tests)"
else
    echo "   ❌ Some AlphaRuleBuilder tests failed"
    exit 1
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "✨ Demonstrating the fix with verbose output..."
echo ""
echo "Running: TestAlphaChain_TwoRules_SameConditions_DifferentOrder"
echo "This test creates two rules with the same conditions in different order"
echo "and verifies they share the same AlphaNodes."
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

go test -v -run "^TestAlphaChain_TwoRules_SameConditions_DifferentOrder$" . 2>&1 | \
    grep -E "(Multi-condition AND detected|Décomposition en chaîne|AlphaChainBuilder|Chaîne construite|partagé)" | \
    head -20

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "🎉 SUCCESS! All target tests are now passing!"
echo ""
echo "📊 Summary:"
echo "   - 4 previously failing tests: ✅ FIXED"
echo "   - AlphaChainBuilder: ✅ Integrated with constraint pipeline"
echo "   - Alpha node sharing: ✅ Working correctly for chains"
echo "   - Backward compatibility: ✅ Simple alpha rules still work"
echo ""
echo "For more details, see: rete/ALPHA_CHAIN_INTEGRATION_FIX.md"
echo ""
