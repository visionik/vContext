package updater

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visionik/vAgenda/api/go/pkg/core"
)

func TestNewUpdater(t *testing.T) {
	// With nil validator
	u := New(nil)
	assert.NotNil(t, u)
	assert.NotNil(t, u.validator)
}

func TestAddTodoItem(t *testing.T) {
	u := New(nil)
	doc := &core.Document{
		Info: core.Info{
			Version: "1.0",
		},
		TodoList: &core.TodoList{},
	}

	item := core.TodoItem{Title: "Task 1", Status: core.StatusPending}
	err := u.AddTodoItem(doc, item)

	require.NoError(t, err)
	assert.NotNil(t, doc.TodoList)
	assert.Len(t, doc.TodoList.Items, 1)
	assert.Equal(t, "Task 1", doc.TodoList.Items[0].Title)

	// Add another
	err = u.AddTodoItem(doc, core.TodoItem{Title: "Task 2", Status: core.StatusPending})
	require.NoError(t, err)
	assert.Len(t, doc.TodoList.Items, 2)
}

func TestAddTodoItemValidationError(t *testing.T) {
	u := New(nil)
	doc := &core.Document{
		Info: core.Info{
			Version: "1.0",
		},
		Plan:     &core.Plan{},     // Has Plan, wrong for todo-list
		TodoList: &core.TodoList{}, // Cannot have both
	}

	item := core.TodoItem{Title: "Task 1", Status: core.StatusPending}
	err := u.AddTodoItem(doc, item)

	assert.Error(t, err) // Should fail validation
}

func TestRemoveTodoItem(t *testing.T) {
	u := New(nil)
	doc := &core.Document{
		Info: core.Info{
			Version: "1.0",
		},
		TodoList: &core.TodoList{
			Items: []core.TodoItem{
				{Title: "Task 1", Status: core.StatusPending},
				{Title: "Task 2", Status: core.StatusPending},
			},
		},
	}

	err := u.RemoveTodoItem(doc, 0)
	require.NoError(t, err)
	assert.Len(t, doc.TodoList.Items, 1)
	assert.Equal(t, "Task 2", doc.TodoList.Items[0].Title)

	// Invalid index
	err = u.RemoveTodoItem(doc, 5)
	assert.ErrorIs(t, err, core.ErrInvalidIndex)
}

func TestUpdateTodoItem(t *testing.T) {
	u := New(nil)
	doc := &core.Document{
		Info: core.Info{
			Version: "1.0",
		},
		TodoList: &core.TodoList{
			Items: []core.TodoItem{
				{Title: "Task 1", Status: core.StatusPending},
			},
		},
	}

	err := u.UpdateTodoItem(doc, 0, func(item *core.TodoItem) {
		item.Status = core.StatusCompleted
	})

	require.NoError(t, err)
	assert.Equal(t, core.StatusCompleted, doc.TodoList.Items[0].Status)
}

func TestAddPlanNarrative(t *testing.T) {
	u := New(nil)
	doc := &core.Document{
		Info: core.Info{
			Version: "1.0",
		},
		Plan: &core.Plan{
			Title:      "Test Plan",
			Status:     core.PlanStatusDraft,
			Narratives: map[string]core.Narrative{"proposal": {Title: "Proposal", Content: "Content"}},
		},
	}

	narrative := core.Narrative{Title: "Overview", Content: "Content"}
	err := u.AddPlanNarrative(doc, "overview", narrative)

	require.NoError(t, err)
	assert.NotNil(t, doc.Plan)
	assert.Len(t, doc.Plan.Narratives, 2) // proposal + overview
	assert.Equal(t, "Content", doc.Plan.Narratives["overview"].Content)
}

func TestRemovePlanNarrative(t *testing.T) {
	u := New(nil)
	doc := &core.Document{
		Info: core.Info{
			Version: "1.0",
		},
		Plan: &core.Plan{
			Title:  "Test Plan",
			Status: core.PlanStatusDraft,
			Narratives: map[string]core.Narrative{
				"proposal": {Title: "Proposal", Content: "Content"},
				"overview": {Title: "Overview", Content: "Content"},
			},
		},
	}

	err := u.RemovePlanNarrative(doc, "overview")
	require.NoError(t, err)
	assert.Len(t, doc.Plan.Narratives, 1) // proposal still remains
	_, exists := doc.Plan.Narratives["overview"]
	assert.False(t, exists)
}

func TestUpdatePlanNarrative(t *testing.T) {
	u := New(nil)
	doc := &core.Document{
		Info: core.Info{
			Version: "1.0",
		},
		Plan: &core.Plan{
			Title:  "Test Plan",
			Status: core.PlanStatusDraft,
			Narratives: map[string]core.Narrative{
				"proposal": {Title: "Proposal", Content: "Content"},
				"overview": {Title: "Overview", Content: "Original"},
			},
		},
	}

	err := u.UpdatePlanNarrative(doc, "overview", func(n *core.Narrative) {
		n.Content = "Updated"
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated", doc.Plan.Narratives["overview"].Content)
}

func TestAddPlanPhase(t *testing.T) {
	u := New(nil)
	doc := &core.Document{
		Info: core.Info{
			Version: "1.0",
		},
		Plan: &core.Plan{
			Title:      "Test Plan",
			Status:     core.PlanStatusDraft,
			Narratives: map[string]core.Narrative{"proposal": {Title: "Proposal", Content: "Content"}},
		},
	}

	phase := core.Phase{Title: "Phase 1", Status: core.PhaseStatusPending}
	err := u.AddPlanPhase(doc, phase)

	require.NoError(t, err)
	assert.NotNil(t, doc.Plan)
	assert.Len(t, doc.Plan.Phases, 1)
	assert.Equal(t, "Phase 1", doc.Plan.Phases[0].Title)
}

func TestRemovePlanPhase(t *testing.T) {
	u := New(nil)
	doc := &core.Document{
		Info: core.Info{
			Version: "1.0",
		},
		Plan: &core.Plan{
			Title:      "Test Plan",
			Status:     core.PlanStatusDraft,
			Narratives: map[string]core.Narrative{"proposal": {Title: "Proposal", Content: "Content"}},
			Phases: []core.Phase{
				{Title: "Phase 1", Status: core.PhaseStatusPending},
				{Title: "Phase 2", Status: core.PhaseStatusPending},
			},
		},
	}

	err := u.RemovePlanPhase(doc, 0)
	require.NoError(t, err)
	assert.Len(t, doc.Plan.Phases, 1)
	assert.Equal(t, "Phase 2", doc.Plan.Phases[0].Title)
}

func TestUpdatePlanPhase(t *testing.T) {
	u := New(nil)
	doc := &core.Document{
		Info: core.Info{
			Version: "1.0",
		},
		Plan: &core.Plan{
			Title:      "Test Plan",
			Status:     core.PlanStatusDraft,
			Narratives: map[string]core.Narrative{"proposal": {Title: "Proposal", Content: "Content"}},
			Phases: []core.Phase{
				{Title: "Phase 1", Status: core.PhaseStatusPending},
			},
		},
	}

	err := u.UpdatePlanPhase(doc, 0, func(p *core.Phase) {
		p.Status = core.PhaseStatusCompleted
	})

	require.NoError(t, err)
	assert.Equal(t, core.PhaseStatusCompleted, doc.Plan.Phases[0].Status)
}
