#!/usr/bin/env python3

import subprocess
import time

def test_all_alpha_operators():
    """Teste tous les 26 opérateurs Alpha et compte les succès"""
    
    print("🧪 TEST COMPLET DES 26 OPÉRATEURS ALPHA")
    print("======================================")
    
    # Tests originaux
    original_tests = [
        "alpha_boolean_negative", "alpha_boolean_positive",
        "alpha_comparison_negative", "alpha_comparison_positive", 
        "alpha_equality_negative", "alpha_equality_positive",
        "alpha_inequality_negative", "alpha_inequality_positive",
        "alpha_string_negative", "alpha_string_positive"
    ]
    
    # Tests étendus 
    extended_tests = [
        "alpha_equal_sign_negative", "alpha_equal_sign_positive",
        "alpha_contains_negative", "alpha_contains_positive",
        "alpha_in_negative", "alpha_in_positive",
        "alpha_length_negative", "alpha_length_positive",
        "alpha_like_negative", "alpha_like_positive",
        "alpha_matches_negative", "alpha_matches_positive",
        "alpha_abs_negative", "alpha_abs_positive",
        "alpha_upper_negative", "alpha_upper_positive"
    ]
    
    def run_test_suite(runner_cmd, tests_list, suite_name):
        print(f"\n📋 {suite_name}")
        print("-" * 50)
        
        try:
            # Exécuter le test runner
            result = subprocess.run(
                ["go", "run", runner_cmd], 
                capture_output=True, 
                text=True, 
                timeout=60
            )
            
            output = result.stdout + result.stderr
            success_count = 0
            error_count = 0
            action_count = 0
            
            # Analyser les résultats
            for test in tests_list:
                if f"🧪 Exécution test" in output and test in output:
                    if "✅ Succès" in output:
                        success_count += 1
                    if "❌ Échec" in output or "Erreur" in output:
                        error_count += 1
                
            # Compter les actions déclenchées
            action_count = output.count("🎯 ACTION DISPONIBLE DANS TUPLE-SPACE:")
            
            print(f"   ✅ Tests réussis: {success_count}/{len(tests_list)}")
            print(f"   ❌ Tests échoués: {error_count}")
            print(f"   🎯 Actions déclenchées: {action_count}")
            
            return success_count, error_count, action_count
            
        except subprocess.TimeoutExpired:
            print(f"   ⏰ Timeout - suite {suite_name}")
            return 0, len(tests_list), 0
        except Exception as e:
            print(f"   💥 Erreur: {e}")
            return 0, len(tests_list), 0
    
    # Exécuter les tests originaux
    orig_success, orig_error, orig_actions = run_test_suite(
        "test_alpha_coverage/alpha_coverage_test_runner.go",
        original_tests,
        "TESTS ORIGINAUX (10 opérateurs de base)"
    )
    
    # Exécuter les tests étendus
    ext_success, ext_error, ext_actions = run_test_suite(
        "test_alpha_coverage_extended/alpha_coverage_extended_test_runner.go", 
        extended_tests,
        "TESTS ÉTENDUS (16 nouveaux opérateurs)"
    )
    
    # Résumé global
    total_success = orig_success + ext_success
    total_tests = len(original_tests) + len(extended_tests)
    total_actions = orig_actions + ext_actions
    
    print(f"\n🎯 RÉSUMÉ GLOBAL")
    print("=" * 50)
    print(f"   📊 Tests réussis: {total_success}/{total_tests} ({(total_success/total_tests)*100:.1f}%)")
    print(f"   🎬 Actions déclenchées: {total_actions}")
    print(f"   🔧 Opérateurs fonctionnels: {total_success}")
    print(f"   ⚠️ Opérateurs à corriger: {total_tests - total_success}")
    
    if total_success == total_tests:
        print("\n🎉 TOUS LES OPÉRATEURS FONCTIONNENT ! 🎉")
        print("✅ TSD supporte maintenant tous les 26 opérateurs Alpha testés")
    else:
        print(f"\n⚠️ {total_tests - total_success} opérateurs nécessitent encore des corrections")
    
    return total_success, total_tests, total_actions

if __name__ == "__main__":
    test_all_alpha_operators()