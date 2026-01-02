// Copyright (c) 2025 TSD Contributors
// Licensed under the MIT License
// See LICENSE file in the project root for full license text

package actions

import (
	"github.com/treivax/tsd/rete"
)

// UpdateActionHandler est un wrapper pour l'action Update.
// Il implémente l'interface ActionHandler en délèguant au BuiltinActionExecutor.
type UpdateActionHandler struct {
	executor *BuiltinActionExecutor
}

// NewUpdateActionHandler crée un nouveau handler pour l'action Update.
func NewUpdateActionHandler(executor *BuiltinActionExecutor) *UpdateActionHandler {
	return &UpdateActionHandler{executor: executor}
}

// GetName retourne le nom de l'action.
func (h *UpdateActionHandler) GetName() string {
	return ActionUpdate
}

// Execute exécute l'action Update avec les arguments fournis.
func (h *UpdateActionHandler) Execute(args []interface{}, ctx *rete.ExecutionContext) error {
	return h.executor.Execute(ActionUpdate, args, ctx.GetToken())
}

// Validate valide que les arguments sont corrects pour l'action Update.
func (h *UpdateActionHandler) Validate(args []interface{}) error {
	if len(args) != ArgsCountUpdate {
		return NewValidationError(ActionUpdate, ArgsCountUpdate, len(args))
	}

	if _, ok := args[0].(*rete.Fact); !ok {
		return NewTypeError(ActionUpdate, 0, "*rete.Fact", args[0])
	}

	return nil
}

// InsertActionHandler est un wrapper pour l'action Insert.
// Il implémente l'interface ActionHandler en délèguant au BuiltinActionExecutor.
type InsertActionHandler struct {
	executor *BuiltinActionExecutor
}

// NewInsertActionHandler crée un nouveau handler pour l'action Insert.
func NewInsertActionHandler(executor *BuiltinActionExecutor) *InsertActionHandler {
	return &InsertActionHandler{executor: executor}
}

// GetName retourne le nom de l'action.
func (h *InsertActionHandler) GetName() string {
	return ActionInsert
}

// Execute exécute l'action Insert avec les arguments fournis.
func (h *InsertActionHandler) Execute(args []interface{}, ctx *rete.ExecutionContext) error {
	return h.executor.Execute(ActionInsert, args, ctx.GetToken())
}

// Validate valide que les arguments sont corrects pour l'action Insert.
func (h *InsertActionHandler) Validate(args []interface{}) error {
	if len(args) != ArgsCountInsert {
		return NewValidationError(ActionInsert, ArgsCountInsert, len(args))
	}

	if _, ok := args[0].(*rete.Fact); !ok {
		return NewTypeError(ActionInsert, 0, "*rete.Fact", args[0])
	}

	return nil
}

// RetractActionHandler est un wrapper pour l'action Retract.
// Il implémente l'interface ActionHandler en délèguant au BuiltinActionExecutor.
type RetractActionHandler struct {
	executor *BuiltinActionExecutor
}

// NewRetractActionHandler crée un nouveau handler pour l'action Retract.
func NewRetractActionHandler(executor *BuiltinActionExecutor) *RetractActionHandler {
	return &RetractActionHandler{executor: executor}
}

// GetName retourne le nom de l'action.
func (h *RetractActionHandler) GetName() string {
	return ActionRetract
}

// Execute exécute l'action Retract avec les arguments fournis.
func (h *RetractActionHandler) Execute(args []interface{}, ctx *rete.ExecutionContext) error {
	return h.executor.Execute(ActionRetract, args, ctx.GetToken())
}

// Validate valide que les arguments sont corrects pour l'action Retract.
func (h *RetractActionHandler) Validate(args []interface{}) error {
	if len(args) != ArgsCountRetract {
		return NewValidationError(ActionRetract, ArgsCountRetract, len(args))
	}

	if _, ok := args[0].(string); !ok {
		return NewTypeError(ActionRetract, 0, "string", args[0])
	}

	return nil
}

// PrintActionHandler est un wrapper pour l'action Print.
// Il implémente l'interface ActionHandler en délèguant au BuiltinActionExecutor.
type PrintActionHandler struct {
	executor *BuiltinActionExecutor
}

// NewPrintActionHandler crée un nouveau handler pour l'action Print.
func NewPrintActionHandler(executor *BuiltinActionExecutor) *PrintActionHandler {
	return &PrintActionHandler{executor: executor}
}

// GetName retourne le nom de l'action.
func (h *PrintActionHandler) GetName() string {
	return ActionPrint
}

// Execute exécute l'action Print avec les arguments fournis.
func (h *PrintActionHandler) Execute(args []interface{}, ctx *rete.ExecutionContext) error {
	return h.executor.Execute(ActionPrint, args, ctx.GetToken())
}

// Validate valide que les arguments sont corrects pour l'action Print.
func (h *PrintActionHandler) Validate(args []interface{}) error {
	if len(args) != ArgsCountPrint {
		return NewValidationError(ActionPrint, ArgsCountPrint, len(args))
	}

	if _, ok := args[0].(string); !ok {
		return NewTypeError(ActionPrint, 0, "string", args[0])
	}

	return nil
}

// LogActionHandler est un wrapper pour l'action Log.
// Il implémente l'interface ActionHandler en délèguant au BuiltinActionExecutor.
type LogActionHandler struct {
	executor *BuiltinActionExecutor
}

// NewLogActionHandler crée un nouveau handler pour l'action Log.
func NewLogActionHandler(executor *BuiltinActionExecutor) *LogActionHandler {
	return &LogActionHandler{executor: executor}
}

// GetName retourne le nom de l'action.
func (h *LogActionHandler) GetName() string {
	return ActionLog
}

// Execute exécute l'action Log avec les arguments fournis.
func (h *LogActionHandler) Execute(args []interface{}, ctx *rete.ExecutionContext) error {
	return h.executor.Execute(ActionLog, args, ctx.GetToken())
}

// Validate valide que les arguments sont corrects pour l'action Log.
func (h *LogActionHandler) Validate(args []interface{}) error {
	if len(args) != ArgsCountLog {
		return NewValidationError(ActionLog, ArgsCountLog, len(args))
	}

	if _, ok := args[0].(string); !ok {
		return NewTypeError(ActionLog, 0, "string", args[0])
	}

	return nil
}

// XupleActionHandler est un wrapper pour l'action Xuple.
// Il implémente l'interface ActionHandler en délèguant au BuiltinActionExecutor.
type XupleActionHandler struct {
	executor *BuiltinActionExecutor
}

// NewXupleActionHandler crée un nouveau handler pour l'action Xuple.
func NewXupleActionHandler(executor *BuiltinActionExecutor) *XupleActionHandler {
	return &XupleActionHandler{executor: executor}
}

// GetName retourne le nom de l'action.
func (h *XupleActionHandler) GetName() string {
	return ActionXuple
}

// Execute exécute l'action Xuple avec les arguments fournis.
func (h *XupleActionHandler) Execute(args []interface{}, ctx *rete.ExecutionContext) error {
	return h.executor.Execute(ActionXuple, args, ctx.GetToken())
}

// Validate valide que les arguments sont corrects pour l'action Xuple.
func (h *XupleActionHandler) Validate(args []interface{}) error {
	if len(args) != ArgsCountXuple {
		return NewValidationError(ActionXuple, ArgsCountXuple, len(args))
	}

	if _, ok := args[0].(string); !ok {
		return NewTypeError(ActionXuple, 0, "string", args[0])
	}

	if _, ok := args[1].(*rete.Fact); !ok {
		return NewTypeError(ActionXuple, 1, "*rete.Fact", args[1])
	}

	return nil
}
