package constraint

// API publique du module constraint
// Ce fichier expose les fonctions principales pour l'utilisation externe

// Parse analyse un fichier de contraintes et retourne l'AST
func ParseConstraint(filename string, input []byte) (interface{}, error) {
	return Parse(filename, input)
}

// ValidateConstraintProgram valide un programme d'AST de contraintes
func ValidateConstraintProgram(result interface{}) error {
	return ValidateProgram(result)
}

// ParseConstraintFile analyse un fichier de contraintes depuis le système de fichiers
func ParseConstraintFile(filename string) (interface{}, error) {
	return ParseFile(filename)
}
