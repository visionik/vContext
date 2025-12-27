// Package core provides the core types and interfaces for vAgenda documents.
package core

// Document represents the root vAgenda document.
// A document contains metadata and either a TodoList or a Plan (but not both).
type Document struct {
	Info     Info      `json:"vAgendaInfo" tron:"vAgendaInfo"`
	TodoList *TodoList `json:"todoList,omitempty" tron:"todoList,omitempty"`
	Plan     *Plan     `json:"plan,omitempty" tron:"plan,omitempty"`
}

// Info contains document-level metadata that appears once per file.
type Info struct {
	Version     string                 `json:"version" tron:"version"`
	Author      string                 `json:"author,omitempty" tron:"author,omitempty"`
	Description string                 `json:"description,omitempty" tron:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty" tron:"metadata,omitempty"`
}

// TodoList represents a collection of actionable work items for short-term memory.
type TodoList struct {
	Items []TodoItem `json:"items" tron:"items"`
}

// TodoItem represents a single actionable task with status tracking.
type TodoItem struct {
	Title  string     `json:"title" tron:"title"`
	Status ItemStatus `json:"status" tron:"status"`
}

// ItemStatus represents the status of a todo item.
type ItemStatus string

const (
	// StatusPending indicates the item has not been started.
	StatusPending ItemStatus = "pending"
	// StatusInProgress indicates the item is currently being worked on.
	StatusInProgress ItemStatus = "inProgress"
	// StatusCompleted indicates the item has been finished.
	StatusCompleted ItemStatus = "completed"
	// StatusBlocked indicates the item cannot proceed.
	StatusBlocked ItemStatus = "blocked"
	// StatusCancelled indicates the item has been cancelled.
	StatusCancelled ItemStatus = "cancelled"
)

// IsValid returns true if the ItemStatus is a valid value.
func (s ItemStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusInProgress, StatusCompleted, StatusBlocked, StatusCancelled:
		return true
	default:
		return false
	}
}

// Plan represents a structured design document for medium-term memory.
type Plan struct {
	Title      string               `json:"title" tron:"title"`
	Status     PlanStatus           `json:"status" tron:"status"`
	Narratives map[string]Narrative `json:"narratives" tron:"narratives"`
	Phases     []Phase              `json:"phases,omitempty" tron:"phases,omitempty"`
}

// PlanStatus represents the status of a plan.
type PlanStatus string

const (
	// PlanStatusDraft indicates the plan is being drafted.
	PlanStatusDraft PlanStatus = "draft"
	// PlanStatusProposed indicates the plan has been proposed for review.
	PlanStatusProposed PlanStatus = "proposed"
	// PlanStatusApproved indicates the plan has been approved.
	PlanStatusApproved PlanStatus = "approved"
	// PlanStatusInProgress indicates the plan is being executed.
	PlanStatusInProgress PlanStatus = "inProgress"
	// PlanStatusCompleted indicates the plan has been completed.
	PlanStatusCompleted PlanStatus = "completed"
	// PlanStatusCancelled indicates the plan has been cancelled.
	PlanStatusCancelled PlanStatus = "cancelled"
)

// IsValid returns true if the PlanStatus is a valid value.
func (s PlanStatus) IsValid() bool {
	switch s {
	case PlanStatusDraft, PlanStatusProposed, PlanStatusApproved,
		PlanStatusInProgress, PlanStatusCompleted, PlanStatusCancelled:
		return true
	default:
		return false
	}
}

// Phase represents a stage of work within a plan.
type Phase struct {
	Title  string      `json:"title" tron:"title"`
	Status PhaseStatus `json:"status" tron:"status"`
}

// PhaseStatus represents the status of a phase.
type PhaseStatus string

const (
	// PhaseStatusPending indicates the phase has not been started.
	PhaseStatusPending PhaseStatus = "pending"
	// PhaseStatusInProgress indicates the phase is currently active.
	PhaseStatusInProgress PhaseStatus = "inProgress"
	// PhaseStatusCompleted indicates the phase has been finished.
	PhaseStatusCompleted PhaseStatus = "completed"
	// PhaseStatusBlocked indicates the phase cannot proceed.
	PhaseStatusBlocked PhaseStatus = "blocked"
	// PhaseStatusCancelled indicates the phase has been cancelled.
	PhaseStatusCancelled PhaseStatus = "cancelled"
)

// IsValid returns true if the PhaseStatus is a valid value.
func (s PhaseStatus) IsValid() bool {
	switch s {
	case PhaseStatusPending, PhaseStatusInProgress, PhaseStatusCompleted,
		PhaseStatusBlocked, PhaseStatusCancelled:
		return true
	default:
		return false
	}
}

// Narrative represents a named block of documentation within a plan.
type Narrative struct {
	Title   string `json:"title" tron:"title"`
	Content string `json:"content" tron:"content"`
}
