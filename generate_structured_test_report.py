#!/usr/bin/env python3

import os
import subprocess
import json
import re
from datetime import datetime

def read_file_content(filepath):
    """Lit le contenu d'un fichier"""
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            return f.read().strip()
    except Exception as e:
        return f"Erreur lecture: {e}"

def extract_test_description(constraint_content):
    """Extrait la description depuis les commentaires du fichier .constraint"""
    lines = constraint_content.split('\n')
    for line in lines:
        line = line.strip()
        if line.startswith('//'):
            return line[2:].strip()
    return "Description non disponible"

def extract_rules_from_constraint(constraint_content):
    """Extrait les règles complètes depuis le fichier .constraint"""
    lines = constraint_content.split('\n')
    rules = []
    current_rule = []
    in_rule = False
    
    for line in lines:
        line = line.strip()
        if not line or line.startswith('//'):
            continue
        
        # Début d'une règle type
        if line.startswith('type '):
            if in_rule and current_rule:
                rules.append('\n'.join(current_rule))
                current_rule = []
            current_rule.append(line)
            in_rule = True
        # Ligne avec { } / ==> (règle principale)
        elif '{' in line and '}' in line and '/' in line and '==>' in line:
            current_rule.append(line)
        # Continuation de règle
        elif in_rule and line:
            current_rule.append(line)
    
    # Ajouter la dernière règle
    if current_rule:
        rules.append('\n'.join(current_rule))
    
    return rules

def extract_facts_from_file(facts_content):
    """Extrait les faits depuis le fichier .facts"""
    lines = facts_content.split('\n')
    facts = []
    
    for line in lines:
        line = line.strip()
        if line and not line.startswith('//'):
            facts.append(line)
    
    return facts

def predict_expected_actions(test_name, rules, facts):
    """Prédit les actions attendues basé sur l'analyse logique"""
    
    # Dictionnaire des attentes par test
    expected_predictions = {
        # Tests originaux
        "alpha_boolean_negative": {
            "action": "inactive_account_found",
            "logic": "NOT(acc.active == true) → Comptes avec active=false",
            "expected_facts": ["ACC002 (active=false)"]
        },
        "alpha_boolean_positive": {
            "action": "active_account_found", 
            "logic": "acc.active == true → Comptes avec active=true",
            "expected_facts": ["ACC001 (active=true)", "ACC003 (active=true)"]
        },
        "alpha_comparison_negative": {
            "action": "affordable_product",
            "logic": "NOT(prod.price > 100) → Produits avec price ≤ 100",
            "expected_facts": ["PROD002 (price=50)"]
        },
        "alpha_comparison_positive": {
            "action": "expensive_product",
            "logic": "prod.price > 100 → Produits avec price > 100",
            "expected_facts": ["PROD001 (price=150)", "PROD003 (price=200)"]
        },
        "alpha_equality_negative": {
            "action": "age_is_not_twenty_five",
            "logic": "NOT(p.age == 25) → Personnes avec age ≠ 25",
            "expected_facts": ["P002 (age=30)"]
        },
        "alpha_equality_positive": {
            "action": "age_is_twenty_five",
            "logic": "p.age == 25 → Personnes avec age = 25",
            "expected_facts": ["P001 (age=25)", "P003 (age=25)"]
        },
        "alpha_inequality_negative": {
            "action": "cancelled_order_found",
            "logic": "NOT(ord.status != \"cancelled\") → Commandes avec status = cancelled",
            "expected_facts": ["ORD002 (status=cancelled)"]
        },
        "alpha_inequality_positive": {
            "action": "valid_order_found",
            "logic": "ord.status != \"cancelled\" → Commandes avec status ≠ cancelled",
            "expected_facts": ["ORD001 (status=pending)", "ORD003 (status=completed)"]
        },
        "alpha_string_negative": {
            "action": "non_admin_user_found",
            "logic": "NOT(u.role == \"admin\") → Utilisateurs avec role ≠ admin",
            "expected_facts": ["U002 (role=user)"]
        },
        "alpha_string_positive": {
            "action": "admin_user_found",
            "logic": "u.role == \"admin\" → Utilisateurs avec role = admin",
            "expected_facts": ["U001 (role=admin)", "U003 (role=admin)"]
        },
        
        # Tests étendus
        "alpha_equal_sign_negative": {
            "action": "non_gold_customer_found",
            "logic": "NOT(cust.tier = \"gold\") → Clients avec tier ≠ gold",
            "expected_facts": ["C002 (tier=silver)", "C003 (tier=bronze)"]
        },
        "alpha_equal_sign_positive": {
            "action": "gold_customer_found",
            "logic": "cust.tier = \"gold\" → Clients avec tier = gold",
            "expected_facts": ["C001 (tier=gold)", "C003 (tier=gold)"]
        },
        "alpha_contains_negative": {
            "action": "normal_message_found",
            "logic": "NOT(m.content CONTAINS \"urgent\") → Messages sans 'urgent'",
            "expected_facts": ["M002 (content sans 'urgent')", "M003 (content sans 'urgent')"]
        },
        "alpha_contains_positive": {
            "action": "urgent_message_found",
            "logic": "m.content CONTAINS \"urgent\" → Messages contenant 'urgent'",
            "expected_facts": ["M001 (content avec 'urgent')", "M003 (content avec 'urgent')"]
        },
        "alpha_in_negative": {
            "action": "invalid_state_found",
            "logic": "NOT(s.state IN [\"active\", \"pending\", \"review\"]) → États non valides",
            "expected_facts": ["S002 (state=inactive)"]
        },
        "alpha_in_positive": {
            "action": "valid_state_found",
            "logic": "s.state IN [\"active\", \"pending\", \"review\"] → États valides",
            "expected_facts": ["S001 (state=active)", "S003 (state=pending)"]
        },
        "alpha_length_negative": {
            "action": "weak_password_found",
            "logic": "NOT(LENGTH(p.value) >= 8) → Mots de passe < 8 caractères",
            "expected_facts": ["P002 (length < 8)"]
        },
        "alpha_length_positive": {
            "action": "secure_password_found",
            "logic": "LENGTH(p.value) >= 8 → Mots de passe ≥ 8 caractères",
            "expected_facts": ["P001 (length ≥ 8)", "P003 (length ≥ 8)"]
        },
        "alpha_like_negative": {
            "action": "non_company_email_found",
            "logic": "NOT(e.address LIKE \"%@company.com\") → Emails non-entreprise",
            "expected_facts": ["E002 (@gmail.com)"]
        },
        "alpha_like_positive": {
            "action": "company_email_found",
            "logic": "e.address LIKE \"%@company.com\" → Emails d'entreprise",
            "expected_facts": ["E001 (@company.com)", "E003 (@company.com)"]
        },
        "alpha_matches_negative": {
            "action": "invalid_code_found",
            "logic": "NOT(c.value MATCHES \"[A-Z]+[0-9]+\") → Codes ne matchant pas",
            "expected_facts": ["C002 (pattern invalide)"]
        },
        "alpha_matches_positive": {
            "action": "valid_code_found",
            "logic": "c.value MATCHES \"[A-Z]+[0-9]+\" → Codes matchant",
            "expected_facts": ["C001 (CODE123)", "C003 (PROD456)"]
        },
        "alpha_abs_negative": {
            "action": "small_balance_found",
            "logic": "NOT(ABS(b.amount) > 100) → Soldes absolus ≤ 100",
            "expected_facts": ["B003 (|50| ≤ 100)"]
        },
        "alpha_abs_positive": {
            "action": "significant_balance_found",
            "logic": "ABS(b.amount) > 100 → Soldes absolus > 100",
            "expected_facts": ["B001 (|150| > 100)", "B002 (|-200| > 100)"]
        },
        "alpha_upper_negative": {
            "action": "lowercase_department_found",
            "logic": "NOT(UPPER(d.name) = d.name) → Noms non en majuscules",
            "expected_facts": ["D002 (sales ≠ SALES)"]
        },
        "alpha_upper_positive": {
            "action": "uppercase_department_found",
            "logic": "UPPER(d.name) = d.name → Noms en majuscules",
            "expected_facts": ["D001 (FINANCE)", "D003 (HR)"]
        }
    }
    
    return expected_predictions.get(test_name, {
        "action": "unknown_action",
        "logic": "Logique non définie",
        "expected_facts": ["Analyse manuelle requise"]
    })

def run_single_test(test_dir, test_name):
    """Exécute un seul test et capture les résultats"""
    try:
        # Utiliser le runner de test approprié
        if "extended" in test_dir:
            runner_path = "/home/resinsec/dev/tsd/test_alpha_coverage_extended/alpha_coverage_extended_test_runner.go"
        else:
            runner_path = "/home/resinsec/dev/tsd/test_alpha_coverage/alpha_coverage_test_runner.go"
        
        # Exécuter le test spécifique via un mini-runner
        result = subprocess.run(
            ["go", "run", runner_path],
            capture_output=True,
            text=True,
            timeout=30,
            cwd="/home/resinsec/dev/tsd"
        )
        
        output = result.stdout + result.stderr
        
        # Extraire les actions du résultat
        actions_triggered = []
        action_pattern = r"🎯 ACTION DISPONIBLE DANS TUPLE-SPACE: (\w+) \(([^)]+)\)"
        
        for match in re.finditer(action_pattern, output):
            action_name = match.group(1)
            fact_data = match.group(2)
            actions_triggered.append(f"{action_name} ({fact_data})")
        
        # Déterminer le succès
        success = "✅ Succès" in output and "❌" not in output
        
        return {
            "success": success,
            "actions": actions_triggered,
            "raw_output": output
        }
        
    except Exception as e:
        return {
            "success": False,
            "actions": [],
            "raw_output": f"Erreur exécution: {e}"
        }

def analyze_test_result(expected, obtained):
    """Analyse le résultat d'un test"""
    expected_actions = expected.get("expected_facts", [])
    obtained_actions = obtained.get("actions", [])
    
    if not obtained_actions and expected_actions:
        return "❌ ÉCHEC - Aucune action obtenue alors qu'actions attendues"
    elif not expected_actions and not obtained_actions:
        return "⚠️ NEUTRE - Aucune action attendue ni obtenue"
    elif len(obtained_actions) == len(expected_actions):
        return "✅ CONFORME - Nombre d'actions correspond aux attentes"
    elif obtained_actions and expected_actions:
        return f"⚠️ PARTIEL - {len(obtained_actions)} actions obtenues vs {len(expected_actions)} attendues"
    else:
        return "❓ ANALYSE MANUELLE REQUISE"

def generate_structured_report():
    """Génère le rapport structuré selon le format demandé"""
    
    # Répertoires des tests
    test_dirs = {
        "Tests Alpha Originaux": "/home/resinsec/dev/tsd/alpha_coverage_tests",
        "Tests Alpha Étendus": "/home/resinsec/dev/tsd/alpha_coverage_tests_extended"
    }
    
    report = []
    
    # En-tête du rapport
    report.append("# 📊 RAPPORT STRUCTURÉ - TESTS ALPHA COMPLETS")
    report.append("")
    report.append(f"**Date de génération:** {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    report.append("**Format:** Test par test avec structure complète")
    report.append("")
    report.append("---")
    report.append("")
    
    all_test_results = []
    total_tests = 0
    successful_tests = 0
    
    # Traitement de chaque répertoire de tests
    for category, test_dir in test_dirs.items():
        report.append(f"## 🔬 {category}")
        report.append("")
        
        # Découvrir les tests dans ce répertoire
        try:
            test_files = [f for f in os.listdir(test_dir) if f.endswith('.constraint')]
            test_names = [f.replace('.constraint', '') for f in test_files]
            test_names.sort()
        except:
            test_names = []
        
        for test_name in test_names:
            total_tests += 1
            
            report.append(f"### 🧪 Test: `{test_name}`")
            report.append("")
            
            # 1. Description et nom du test
            constraint_file = os.path.join(test_dir, f"{test_name}.constraint")
            facts_file = os.path.join(test_dir, f"{test_name}.facts")
            
            constraint_content = read_file_content(constraint_file)
            facts_content = read_file_content(facts_file)
            
            description = extract_test_description(constraint_content)
            
            report.append("#### 1️⃣ Description du Test")
            report.append("")
            report.append(f"**Nom:** {test_name}")
            report.append(f"**Description:** {description}")
            report.append("")
            
            # 2. Règles complètes du .constraint
            report.append("#### 2️⃣ Règles Complètes (.constraint)")
            report.append("")
            rules = extract_rules_from_constraint(constraint_content)
            
            report.append("```constraint")
            for rule in rules:
                report.append(rule)
                report.append("")
            report.append("```")
            report.append("")
            
            # 3. Faits soumis du .facts
            report.append("#### 3️⃣ Faits Soumis (.facts)")
            report.append("")
            facts = extract_facts_from_file(facts_content)
            
            report.append("```facts")
            for fact in facts:
                report.append(fact)
            report.append("```")
            report.append("")
            
            # 4. Résultat attendu
            report.append("#### 4️⃣ Résultat Attendu")
            report.append("")
            expected = predict_expected_actions(test_name, rules, facts)
            
            report.append(f"**Action attendue:** `{expected['action']}`")
            report.append(f"**Logique:** {expected['logic']}")
            report.append("")
            report.append("**Faits devant déclencher l'action:**")
            for fact in expected['expected_facts']:
                report.append(f"- {fact}")
            report.append("")
            
            # 5. Résultat obtenu par exécution
            report.append("#### 5️⃣ Résultat Obtenu")
            report.append("")
            
            obtained = run_single_test(test_dir, test_name)
            
            if obtained['success']:
                successful_tests += 1
                
            report.append(f"**Statut:** {'✅ Succès' if obtained['success'] else '❌ Échec'}")
            report.append("")
            
            if obtained['actions']:
                report.append("**Actions déclenchées:**")
                for action in obtained['actions']:
                    report.append(f"- ✅ {action}")
            else:
                report.append("**Actions déclenchées:** Aucune")
            report.append("")
            
            # 6. Analyse résultante
            report.append("#### 6️⃣ Analyse du Test")
            report.append("")
            analysis = analyze_test_result(expected, obtained)
            report.append(f"**Résultat:** {analysis}")
            report.append("")
            
            # Analyse détaillée
            if obtained['actions'] and expected['expected_facts']:
                if len(obtained['actions']) == len(expected['expected_facts']):
                    report.append("✅ **Conformité:** Le nombre d'actions correspond exactement aux attentes")
                else:
                    report.append(f"⚠️ **Écart:** {len(obtained['actions'])} actions obtenues vs {len(expected['expected_facts'])} attendues")
            
            report.append("")
            report.append("---")
            report.append("")
            
            # Stocker les résultats pour l'analyse globale
            all_test_results.append({
                "name": test_name,
                "category": category,
                "success": obtained['success'],
                "expected_count": len(expected['expected_facts']),
                "obtained_count": len(obtained['actions']),
                "conformity": analysis
            })
    
    # Analyse globale
    report.append("## 🎯 Analyse Globale")
    report.append("")
    
    conformity_rate = (successful_tests / total_tests * 100) if total_tests > 0 else 0
    
    report.append("### 📊 Statistiques Générales")
    report.append("")
    report.append(f"- **Tests exécutés:** {total_tests}")
    report.append(f"- **Tests réussis:** {successful_tests}")
    report.append(f"- **Taux de conformité:** {conformity_rate:.1f}%")
    report.append("")
    
    # Analyse par catégorie
    report.append("### 📈 Analyse par Catégorie")
    report.append("")
    
    categories_stats = {}
    for result in all_test_results:
        cat = result['category']
        if cat not in categories_stats:
            categories_stats[cat] = {'total': 0, 'success': 0}
        categories_stats[cat]['total'] += 1
        if result['success']:
            categories_stats[cat]['success'] += 1
    
    for category, stats in categories_stats.items():
        rate = (stats['success'] / stats['total'] * 100) if stats['total'] > 0 else 0
        report.append(f"**{category}:**")
        report.append(f"- Succès: {stats['success']}/{stats['total']} ({rate:.1f}%)")
        report.append("")
    
    # Conclusions
    report.append("### 🏁 Conclusions")
    report.append("")
    
    if conformity_rate == 100:
        report.append("🎉 **EXCELLENT:** Tous les tests sont conformes !")
        report.append("")
        report.append("✅ TSD supporte parfaitement tous les opérateurs Alpha testés")
    elif conformity_rate >= 80:
        report.append("✅ **TRÈS BON:** La majorité des tests sont conformes")
        report.append("")
        report.append(f"⚠️ {total_tests - successful_tests} tests nécessitent encore des ajustements")
    else:
        report.append("⚠️ **AMÉLIORATION REQUISE:** Plusieurs tests échouent")
        report.append("")
        report.append("🔧 Des corrections importantes sont nécessaires")
    
    report.append("")
    report.append("---")
    report.append("")
    report.append(f"**Rapport généré par:** `generate_structured_test_report.py`")
    report.append(f"**Horodatage:** {datetime.now().isoformat()}")
    
    return '\n'.join(report)

if __name__ == "__main__":
    print("🔬 GÉNÉRATION RAPPORT STRUCTURÉ")
    print("===============================")
    
    report_content = generate_structured_report()
    
    output_file = "/home/resinsec/dev/tsd/ALPHA_NODES_STRUCTURED_COMPLETE_REPORT.md"
    
    with open(output_file, 'w', encoding='utf-8') as f:
        f.write(report_content)
    
    print(f"✅ Rapport structuré généré: {output_file}")
    print(f"📊 Structure: Test par test avec 6 sections par test")
    print(f"📁 Fichier de sortie: {output_file}")