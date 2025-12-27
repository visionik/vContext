package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTodoListAddItem(t *testing.T) {
	tl := &TodoList{}
	item := TodoItem{Title: "Task 1", Status: StatusPending}

	tl.AddItem(item)

	assert.Len(t, tl.Items, 1)
	assert.Equal(t, "Task 1", tl.Items[0].Title)
}

func TestTodoListRemoveItem(t *testing.T) {
	tl := &TodoList{
		Items: []TodoItem{
			{Title: "Task 1", Status: StatusPending},
			{Title: "Task 2", Status: StatusPending},
			{Title: "Task 3", Status: StatusPending},
		},
	}

	// Remove middle item
	err := tl.RemoveItem(1)
	require.NoError(t, err)
	assert.Len(t, tl.Items, 2)
	assert.Equal(t, "Task 1", tl.Items[0].Title)
	assert.Equal(t, "Task 3", tl.Items[1].Title)

	// Invalid index
	err = tl.RemoveItem(5)
	assert.ErrorIs(t, err, ErrInvalidIndex)
	err = tl.RemoveItem(-1)
	assert.ErrorIs(t, err, ErrInvalidIndex)
}

func TestTodoListUpdateItem(t *testing.T) {
	tl := &TodoList{
		Items: []TodoItem{
			{Title: "Task 1", Status: StatusPending},
		},
	}

	err := tl.UpdateItem(0, func(item *TodoItem) {
		item.Status = StatusCompleted
		item.Title = "Task 1 Updated"
	})

	require.NoError(t, err)
	assert.Equal(t, StatusCompleted, tl.Items[0].Status)
	assert.Equal(t, "Task 1 Updated", tl.Items[0].Title)

	// Invalid index
	err = tl.UpdateItem(5, func(item *TodoItem) {})
	assert.ErrorIs(t, err, ErrInvalidIndex)
}

func TestTodoListFindItem(t *testing.T) {
	tl := &TodoList{
		Items: []TodoItem{
			{Title: "Task 1", Status: StatusPending},
			{Title: "Task 2", Status: StatusCompleted},
			{Title: "Task 3", Status: StatusPending},
		},
	}

	// Find by status
	item := tl.FindItem(func(i *TodoItem) bool {
		return i.Status == StatusCompleted
	})
	require.NotNil(t, item)
	assert.Equal(t, "Task 2", item.Title)

	// Find by title
	item = tl.FindItem(func(i *TodoItem) bool {
		return i.Title == "Task 3"
	})
	require.NotNil(t, item)
	assert.Equal(t, StatusPending, item.Status)

	// Not found
	item = tl.FindItem(func(i *TodoItem) bool {
		return i.Title == "Nonexistent"
	})
	assert.Nil(t, item)
}

func TestPlanAddNarrative(t *testing.T) {
	plan := &Plan{}
	narrative := Narrative{Title: "Overview", Content: "Some content"}

	plan.AddNarrative("overview", narrative)

	assert.Len(t, plan.Narratives, 1)
	assert.Equal(t, "Some content", plan.Narratives["overview"].Content)

	// Update existing
	plan.AddNarrative("overview", Narrative{Title: "Overview", Content: "Updated content"})
	assert.Len(t, plan.Narratives, 1)
	assert.Equal(t, "Updated content", plan.Narratives["overview"].Content)
}

func TestPlanRemoveNarrative(t *testing.T) {
	plan := &Plan{
		Narratives: map[string]Narrative{
			"overview": {Title: "Overview", Content: "Content 1"},
			"details":  {Title: "Details", Content: "Content 2"},
		},
	}

	plan.RemoveNarrative("overview")

	assert.Len(t, plan.Narratives, 1)
	_, exists := plan.Narratives["overview"]
	assert.False(t, exists)
}

func TestPlanUpdateNarrative(t *testing.T) {
	plan := &Plan{
		Narratives: map[string]Narrative{
			"overview": {Title: "Overview", Content: "Original"},
		},
	}

	err := plan.UpdateNarrative("overview", func(n *Narrative) {
		n.Content = "Updated"
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated", plan.Narratives["overview"].Content)

	// Nonexistent key
	err = plan.UpdateNarrative("nonexistent", func(n *Narrative) {})
	assert.ErrorIs(t, err, ErrNarrativeNotFound)
}

func TestPlanAddPhase(t *testing.T) {
	plan := &Plan{}
	phase := Phase{Title: "Phase 1", Status: PhaseStatusPending}

	plan.AddPhase(phase)

	assert.Len(t, plan.Phases, 1)
	assert.Equal(t, "Phase 1", plan.Phases[0].Title)
}

func TestPlanRemovePhase(t *testing.T) {
	plan := &Plan{
		Phases: []Phase{
			{Title: "Phase 1", Status: PhaseStatusPending},
			{Title: "Phase 2", Status: PhaseStatusPending},
		},
	}

	err := plan.RemovePhase(0)
	require.NoError(t, err)
	assert.Len(t, plan.Phases, 1)
	assert.Equal(t, "Phase 2", plan.Phases[0].Title)

	// Invalid index
	err = plan.RemovePhase(5)
	assert.ErrorIs(t, err, ErrInvalidIndex)
}

func TestPlanUpdatePhase(t *testing.T) {
	plan := &Plan{
		Phases: []Phase{
			{Title: "Phase 1", Status: PhaseStatusPending},
		},
	}

	err := plan.UpdatePhase(0, func(p *Phase) {
		p.Status = PhaseStatusCompleted
		p.Title = "Phase 1 Updated"
	})

	require.NoError(t, err)
	assert.Equal(t, PhaseStatusCompleted, plan.Phases[0].Status)
	assert.Equal(t, "Phase 1 Updated", plan.Phases[0].Title)

	// Invalid index
	err = plan.UpdatePhase(5, func(p *Phase) {})
	assert.ErrorIs(t, err, ErrInvalidIndex)
}
