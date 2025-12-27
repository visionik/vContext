package updater

import (
	"github.com/visionik/vAgenda/api/go/pkg/core"
	"github.com/visionik/vAgenda/api/go/pkg/validator"
)

// Updater provides validated mutation operations on Documents.
type Updater struct {
	validator validator.Validator
}

// New creates a new Updater with the given validator.
// If validator is nil, a default validator will be used.
func New(v validator.Validator) *Updater {
	if v == nil {
		v = validator.New()
	}
	return &Updater{validator: v}
}

// AddTodoItem adds an item to the document's TodoList and validates the result.
func (u *Updater) AddTodoItem(doc *core.Document, item core.TodoItem) error {
	if doc.TodoList == nil {
		doc.TodoList = &core.TodoList{}
	}
	doc.TodoList.AddItem(item)
	return u.validator.Validate(doc)
}

// RemoveTodoItem removes an item from the document's TodoList and validates the result.
func (u *Updater) RemoveTodoItem(doc *core.Document, index int) error {
	if doc.TodoList == nil {
		return core.ErrInvalidIndex
	}
	if err := doc.TodoList.RemoveItem(index); err != nil {
		return err
	}
	return u.validator.Validate(doc)
}

// UpdateTodoItem updates an item in the document's TodoList and validates the result.
func (u *Updater) UpdateTodoItem(doc *core.Document, index int, updates func(*core.TodoItem)) error {
	if doc.TodoList == nil {
		return core.ErrInvalidIndex
	}
	if err := doc.TodoList.UpdateItem(index, updates); err != nil {
		return err
	}
	return u.validator.Validate(doc)
}

// AddPlanNarrative adds a narrative to the document's Plan and validates the result.
func (u *Updater) AddPlanNarrative(doc *core.Document, key string, narrative core.Narrative) error {
	if doc.Plan == nil {
		doc.Plan = &core.Plan{}
	}
	doc.Plan.AddNarrative(key, narrative)
	return u.validator.Validate(doc)
}

// RemovePlanNarrative removes a narrative from the document's Plan and validates the result.
func (u *Updater) RemovePlanNarrative(doc *core.Document, key string) error {
	if doc.Plan == nil {
		return nil // nothing to remove
	}
	doc.Plan.RemoveNarrative(key)
	return u.validator.Validate(doc)
}

// UpdatePlanNarrative updates a narrative in the document's Plan and validates the result.
func (u *Updater) UpdatePlanNarrative(doc *core.Document, key string, updates func(*core.Narrative)) error {
	if doc.Plan == nil {
		return core.ErrNarrativeNotFound
	}
	if err := doc.Plan.UpdateNarrative(key, updates); err != nil {
		return err
	}
	return u.validator.Validate(doc)
}

// AddPlanPhase adds a phase to the document's Plan and validates the result.
func (u *Updater) AddPlanPhase(doc *core.Document, phase core.Phase) error {
	if doc.Plan == nil {
		doc.Plan = &core.Plan{}
	}
	doc.Plan.AddPhase(phase)
	return u.validator.Validate(doc)
}

// RemovePlanPhase removes a phase from the document's Plan and validates the result.
func (u *Updater) RemovePlanPhase(doc *core.Document, index int) error {
	if doc.Plan == nil {
		return core.ErrInvalidIndex
	}
	if err := doc.Plan.RemovePhase(index); err != nil {
		return err
	}
	return u.validator.Validate(doc)
}

// UpdatePlanPhase updates a phase in the document's Plan and validates the result.
func (u *Updater) UpdatePlanPhase(doc *core.Document, index int, updates func(*core.Phase)) error {
	if doc.Plan == nil {
		return core.ErrInvalidIndex
	}
	if err := doc.Plan.UpdatePhase(index, updates); err != nil {
		return err
	}
	return u.validator.Validate(doc)
}
