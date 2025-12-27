package builder

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/visionik/vAgenda/api/go/pkg/core"
)

func TestTodoListBuilder(t *testing.T) {
	t.Run("creates minimal TodoList", func(t *testing.T) {
		doc := NewTodoList("0.2").Build()

		assert.Equal(t, "0.2", doc.Info.Version)
		assert.NotNil(t, doc.TodoList)
		assert.Nil(t, doc.Plan)
		assert.Empty(t, doc.TodoList.Items)
	})

	t.Run("supports fluent API", func(t *testing.T) {
		doc := NewTodoList("0.2").
			WithAuthor("test-author").
			WithDescription("test description").
			AddPendingItem("Task 1").
			AddInProgressItem("Task 2").
			AddCompletedItem("Task 3").
			Build()

		assert.Equal(t, "test-author", doc.Info.Author)
		assert.Equal(t, "test description", doc.Info.Description)
		assert.Len(t, doc.TodoList.Items, 3)

		assert.Equal(t, "Task 1", doc.TodoList.Items[0].Title)
		assert.Equal(t, core.StatusPending, doc.TodoList.Items[0].Status)

		assert.Equal(t, "Task 2", doc.TodoList.Items[1].Title)
		assert.Equal(t, core.StatusInProgress, doc.TodoList.Items[1].Status)

		assert.Equal(t, "Task 3", doc.TodoList.Items[2].Title)
		assert.Equal(t, core.StatusCompleted, doc.TodoList.Items[2].Status)
	})

	t.Run("supports metadata", func(t *testing.T) {
		doc := NewTodoList("0.2").
			WithMetadata("key1", "value1").
			WithMetadata("key2", 42).
			Build()

		assert.Len(t, doc.Info.Metadata, 2)
		assert.Equal(t, "value1", doc.Info.Metadata["key1"])
		assert.Equal(t, 42, doc.Info.Metadata["key2"])
	})

	t.Run("supports AddItem with custom status", func(t *testing.T) {
		doc := NewTodoList("0.2").
			AddItem("Blocked task", core.StatusBlocked).
			Build()

		assert.Len(t, doc.TodoList.Items, 1)
		assert.Equal(t, core.StatusBlocked, doc.TodoList.Items[0].Status)
	})
}

func TestPlanBuilder(t *testing.T) {
	t.Run("creates minimal Plan", func(t *testing.T) {
		doc := NewPlan("Test Plan", "0.2").Build()

		assert.Equal(t, "0.2", doc.Info.Version)
		assert.Nil(t, doc.TodoList)
		assert.NotNil(t, doc.Plan)
		assert.Equal(t, "Test Plan", doc.Plan.Title)
		assert.Equal(t, core.PlanStatusDraft, doc.Plan.Status)
		assert.Empty(t, doc.Plan.Narratives)
	})

	t.Run("supports fluent API", func(t *testing.T) {
		doc := NewPlan("Auth Plan", "0.2").
			WithAuthor("team-lead").
			WithDescription("Authentication implementation").
			WithStatus(core.PlanStatusApproved).
			WithProposal("Proposal", "Use JWT").
			WithProblem("Problem", "No auth").
			WithContext("Context", "Current state").
			AddPendingPhase("Phase 1").
			AddInProgressPhase("Phase 2").
			AddCompletedPhase("Phase 3").
			Build()

		assert.Equal(t, "team-lead", doc.Info.Author)
		assert.Equal(t, "Authentication implementation", doc.Info.Description)
		assert.Equal(t, core.PlanStatusApproved, doc.Plan.Status)

		assert.Len(t, doc.Plan.Narratives, 3)
		assert.Equal(t, "Use JWT", doc.Plan.Narratives["proposal"].Content)
		assert.Equal(t, "No auth", doc.Plan.Narratives["problem"].Content)
		assert.Equal(t, "Current state", doc.Plan.Narratives["context"].Content)

		assert.Len(t, doc.Plan.Phases, 3)
		assert.Equal(t, core.PhaseStatusPending, doc.Plan.Phases[0].Status)
		assert.Equal(t, core.PhaseStatusInProgress, doc.Plan.Phases[1].Status)
		assert.Equal(t, core.PhaseStatusCompleted, doc.Plan.Phases[2].Status)
	})

	t.Run("supports all narrative types", func(t *testing.T) {
		doc := NewPlan("Full Plan", "0.2").
			WithProposal("P", "proposal content").
			WithProblem("Prob", "problem content").
			WithContext("Ctx", "context content").
			WithAlternatives("Alt", "alternatives content").
			WithRisks("Risk", "risks content").
			WithTesting("Test", "testing content").
			Build()

		assert.Len(t, doc.Plan.Narratives, 6)
		assert.Contains(t, doc.Plan.Narratives, "proposal")
		assert.Contains(t, doc.Plan.Narratives, "problem")
		assert.Contains(t, doc.Plan.Narratives, "context")
		assert.Contains(t, doc.Plan.Narratives, "alternatives")
		assert.Contains(t, doc.Plan.Narratives, "risks")
		assert.Contains(t, doc.Plan.Narratives, "testing")
	})

	t.Run("supports AddPhase with custom status", func(t *testing.T) {
		doc := NewPlan("Plan", "0.2").
			AddPhase("Blocked Phase", core.PhaseStatusBlocked).
			Build()

		assert.Len(t, doc.Plan.Phases, 1)
		assert.Equal(t, core.PhaseStatusBlocked, doc.Plan.Phases[0].Status)
	})
}
