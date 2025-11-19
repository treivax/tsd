package rete

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// FactsParser analyse les fichiers .facts pour extraire des ensembles de faits
type FactsParser struct {
	types    map[string]TypeDefinition
	facts    []*Fact
	metadata map[string]string
}

// NewFactsParser crée un nouveau parseur de faits
func NewFactsParser() *FactsParser {
	return &FactsParser{
		types:    make(map[string]TypeDefinition),
		facts:    make([]*Fact, 0),
		metadata: make(map[string]string),
	}
}

// ParseFactsFile analyse un fichier .facts et retourne les faits parsés
func (fp *FactsParser) ParseFactsFile(filePath string, typeDefinitions map[string]TypeDefinition) ([]*Fact, error) {
	fmt.Printf("📋 PARSING FICHIER FAITS: %s\n", filePath)

	// Copier les définitions de types
	fp.types = typeDefinitions

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("erreur ouverture fichier %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		// Ignorer les lignes vides et les commentaires
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			continue
		}

		// Parser les métadonnées
		if strings.HasPrefix(line, "@") {
			err := fp.parseMetadata(line)
			if err != nil {
				return nil, fmt.Errorf("erreur métadonnée ligne %d: %w", lineNumber, err)
			}
			continue
		}

		// Parser les faits
		fact, err := fp.parseFactLine(line, lineNumber)
		if err != nil {
			return nil, fmt.Errorf("erreur fait ligne %d: %w", lineNumber, err)
		}

		if fact != nil {
			fp.facts = append(fp.facts, fact)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erreur lecture fichier: %w", err)
	}

	// Validation de cohérence
	err = fp.validateFacts()
	if err != nil {
		return nil, fmt.Errorf("erreur validation cohérence: %w", err)
	}

	fmt.Printf("✅ Parsé %d faits avec succès\n", len(fp.facts))
	return fp.facts, nil
}

// parseMetadata analyse les lignes de métadonnées (@author, @version, etc.)
func (fp *FactsParser) parseMetadata(line string) error {
	parts := strings.SplitN(line[1:], ":", 2) // Retirer le @ au début
	if len(parts) != 2 {
		return fmt.Errorf("format métadonnée invalide: %s", line)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])
	fp.metadata[key] = value

	return nil
}

// parseFactLine analyse une ligne de fait
func (fp *FactsParser) parseFactLine(line string, lineNumber int) (*Fact, error) {
	// Format attendu: TypeName(id:value, field:value, ...)
	// Exemple: Utilisateur(id:U001, nom:Martin, prenom:Pierre, age:25)

	// Regex pour extraire le type et les champs
	factRegex := regexp.MustCompile(`^(\w+)\((.+)\)$`)
	matches := factRegex.FindStringSubmatch(line)

	if len(matches) != 3 {
		return nil, fmt.Errorf("format de fait invalide: %s", line)
	}

	typeName := matches[1]
	fieldsStr := matches[2]

	// Vérifier que le type existe
	typeDef, exists := fp.types[typeName]
	if !exists {
		return nil, fmt.Errorf("type non défini: %s", typeName)
	}

	// Parser les champs
	fields, err := fp.parseFields(fieldsStr, typeDef)
	if err != nil {
		return nil, fmt.Errorf("erreur parsing champs: %w", err)
	}

	// Générer un ID unique si pas fourni
	factID, hasID := fields["id"]
	if !hasID {
		factID = fmt.Sprintf("fact_%s_%d", typeName, lineNumber)
		fields["id"] = factID
	}

	fact := &Fact{
		ID:        factID.(string),
		Type:      typeName,
		Fields:    fields,
		Timestamp: time.Now(),
	}

	return fact, nil
}

// parseFields analyse la chaîne de champs et retourne une map
func (fp *FactsParser) parseFields(fieldsStr string, typeDef TypeDefinition) (map[string]interface{}, error) {
	fields := make(map[string]interface{})

	// Diviser par les virgules mais attention aux valeurs entre guillemets
	fieldParts := fp.splitFields(fieldsStr)

	for _, fieldPart := range fieldParts {
		fieldPart = strings.TrimSpace(fieldPart)
		if fieldPart == "" {
			continue
		}

		// Format: nom:valeur
		parts := strings.SplitN(fieldPart, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("format champ invalide: %s", fieldPart)
		}

		fieldName := strings.TrimSpace(parts[0])
		fieldValue := strings.TrimSpace(parts[1])

		// Valider le champ selon le type
		value, err := fp.convertFieldValue(fieldName, fieldValue, typeDef)
		if err != nil {
			return nil, fmt.Errorf("erreur conversion champ %s: %w", fieldName, err)
		}

		fields[fieldName] = value
	}

	return fields, nil
}

// splitFields divise les champs en tenant compte des guillemets
func (fp *FactsParser) splitFields(fieldsStr string) []string {
	var fields []string
	var current strings.Builder
	inQuotes := false

	for _, char := range fieldsStr {
		switch char {
		case '"':
			inQuotes = !inQuotes
			current.WriteRune(char)
		case ',':
			if inQuotes {
				current.WriteRune(char)
			} else {
				fields = append(fields, current.String())
				current.Reset()
			}
		default:
			current.WriteRune(char)
		}
	}

	// Ajouter le dernier champ
	if current.Len() > 0 {
		fields = append(fields, current.String())
	}

	return fields
}

// convertFieldValue convertit une valeur de champ selon son type
func (fp *FactsParser) convertFieldValue(fieldName, value string, typeDef TypeDefinition) (interface{}, error) {
	// Trouver la définition du champ
	var fieldDef *Field
	for _, field := range typeDef.Fields {
		if field.Name == fieldName {
			fieldDef = &field
			break
		}
	}

	if fieldDef == nil {
		// Champ non défini dans le type, essayer de déduire le type
		return fp.guessFieldType(value)
	}

	// Conversion selon le type défini
	switch fieldDef.Type {
	case "string":
		// Retirer les guillemets si présents
		if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
			return value[1 : len(value)-1], nil
		}
		return value, nil

	case "number":
		if strings.Contains(value, ".") {
			return strconv.ParseFloat(value, 64)
		}
		intVal, err := strconv.ParseInt(value, 10, 64)
		return float64(intVal), err

	case "bool", "boolean":
		return strconv.ParseBool(value)

	default:
		return value, nil
	}
}

// guessFieldType devine le type d'un champ d'après sa valeur
func (fp *FactsParser) guessFieldType(value string) (interface{}, error) {
	// Booléen
	if value == "true" || value == "false" {
		return strconv.ParseBool(value)
	}

	// Nombre
	if num, err := strconv.ParseFloat(value, 64); err == nil {
		return num, nil
	}

	// String (retirer les guillemets si présents)
	if strings.HasPrefix(value, `"`) && strings.HasSuffix(value, `"`) {
		return value[1 : len(value)-1], nil
	}

	// String par défaut
	return value, nil
}

// validateFacts effectue une validation de cohérence sur les faits parsés
func (fp *FactsParser) validateFacts() error {
	fmt.Printf("🔍 VALIDATION COHÉRENCE DES FAITS...\n")

	idCounts := make(map[string]int)
	typeCounts := make(map[string]int)

	for _, fact := range fp.facts {
		// Vérifier les IDs dupliqués
		idCounts[fact.ID]++
		if idCounts[fact.ID] > 1 {
			return fmt.Errorf("ID de fait dupliqué: %s", fact.ID)
		}

		// Compter les types
		typeCounts[fact.Type]++

		// Vérifier que le type existe
		if _, exists := fp.types[fact.Type]; !exists {
			return fmt.Errorf("type non défini pour le fait %s: %s", fact.ID, fact.Type)
		}

		// Vérifier les champs requis
		err := fp.validateFactFields(fact)
		if err != nil {
			return fmt.Errorf("validation fait %s: %w", fact.ID, err)
		}
	}

	// Afficher les statistiques
	fmt.Printf("📊 STATISTIQUES FAITS:\n")
	for typeName, count := range typeCounts {
		fmt.Printf("   - %s: %d faits\n", typeName, count)
	}

	fmt.Printf("✅ Validation cohérence réussie\n")
	return nil
}

// validateFactFields valide les champs d'un fait selon sa définition de type
func (fp *FactsParser) validateFactFields(fact *Fact) error {
	typeDef := fp.types[fact.Type]

	// Vérifier que tous les champs requis sont présents
	for _, fieldDef := range typeDef.Fields {
		if _, exists := fact.Fields[fieldDef.Name]; !exists {
			// Pour l'instant, on considère tous les champs comme optionnels sauf 'id'
			if fieldDef.Name == "id" {
				return fmt.Errorf("champ requis manquant: %s", fieldDef.Name)
			}
		}
	}

	return nil
}

// GetMetadata retourne les métadonnées parsées
func (fp *FactsParser) GetMetadata() map[string]string {
	return fp.metadata
}

// GetFactsByType retourne les faits filtrés par type
func (fp *FactsParser) GetFactsByType(typeName string) []*Fact {
	var filtered []*Fact
	for _, fact := range fp.facts {
		if fact.Type == typeName {
			filtered = append(filtered, fact)
		}
	}
	return filtered
}

// GetFactsCount retourne le nombre total de faits
func (fp *FactsParser) GetFactsCount() int {
	return len(fp.facts)
}
