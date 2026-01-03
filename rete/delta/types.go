// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package delta

// ChangeType représente le type de changement sur un champ
type ChangeType int

const (
	// ChangeTypeModified indique que le champ a été modifié
	ChangeTypeModified ChangeType = iota
	// ChangeTypeAdded indique que le champ a été ajouté (absent → présent)
	ChangeTypeAdded
	// ChangeTypeRemoved indique que le champ a été supprimé (présent → absent)
	ChangeTypeRemoved
)

// String retourne la représentation string du ChangeType
func (ct ChangeType) String() string {
	switch ct {
	case ChangeTypeModified:
		return "Modified"
	case ChangeTypeAdded:
		return "Added"
	case ChangeTypeRemoved:
		return "Removed"
	default:
		return "Unknown"
	}
}

// ValueType représente le type d'une valeur
type ValueType int

const (
	ValueTypeUnknown ValueType = iota
	ValueTypeString
	ValueTypeNumber
	ValueTypeBool
	ValueTypeObject
	ValueTypeArray
	ValueTypeNull
)

// String retourne la représentation string du ValueType
func (vt ValueType) String() string {
	switch vt {
	case ValueTypeString:
		return "string"
	case ValueTypeNumber:
		return "number"
	case ValueTypeBool:
		return "bool"
	case ValueTypeObject:
		return "object"
	case ValueTypeArray:
		return "array"
	case ValueTypeNull:
		return "null"
	default:
		return "unknown"
	}
}

// inferValueType déduit le ValueType depuis une valeur Go
func inferValueType(value interface{}) ValueType {
	if value == nil {
		return ValueTypeNull
	}

	switch value.(type) {
	case string:
		return ValueTypeString
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return ValueTypeNumber
	case bool:
		return ValueTypeBool
	case map[string]interface{}:
		return ValueTypeObject
	case []interface{}:
		return ValueTypeArray
	default:
		return ValueTypeUnknown
	}
}
