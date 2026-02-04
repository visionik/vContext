package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status ItemStatus
		want   bool
	}{
		{"pending is valid", StatusPending, true},
		{"inProgress is valid", StatusInProgress, true},
		{"completed is valid", StatusCompleted, true},
		{"blocked is valid", StatusBlocked, true},
		{"cancelled is valid", StatusCancelled, true},
		{"empty string is invalid", ItemStatus(""), false},
		{"random string is invalid", ItemStatus("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPlanStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status PlanStatus
		want   bool
	}{
		{"draft is valid", PlanStatusDraft, true},
		{"proposed is valid", PlanStatusProposed, true},
		{"approved is valid", PlanStatusApproved, true},
		{"inProgress is valid", PlanStatusInProgress, true},
		{"completed is valid", PlanStatusCompleted, true},
		{"cancelled is valid", PlanStatusCancelled, true},
		{"empty string is invalid", PlanStatus(""), false},
		{"random string is invalid", PlanStatus("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPlanItemStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status PlanItemStatus
		want   bool
	}{
		{"pending is valid", PlanItemStatusPending, true},
		{"inProgress is valid", PlanItemStatusInProgress, true},
		{"completed is valid", PlanItemStatusCompleted, true},
		{"blocked is valid", PlanItemStatusBlocked, true},
		{"cancelled is valid", PlanItemStatusCancelled, true},
		{"empty string is invalid", PlanItemStatus(""), false},
		{"random string is invalid", PlanItemStatus("invalid"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsValid()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDocument_Structure(t *testing.T) {
	t.Run("document with TodoList", func(t *testing.T) {
		doc := Document{
			Info: Info{Version: "0.2"},
			TodoList: &TodoList{
				Items: []TodoItem{
					{Title: "Task 1", Status: StatusPending},
				},
			},
		}

		assert.Equal(t, "0.2", doc.Info.Version)
		assert.NotNil(t, doc.TodoList)
		assert.Nil(t, doc.Plan)
		assert.Len(t, doc.TodoList.Items, 1)
	})

	t.Run("document with Plan", func(t *testing.T) {
		doc := Document{
			Info: Info{Version: "0.2"},
			Plan: &Plan{
				Title:      "Test Plan",
				Status:     PlanStatusDraft,
				Narratives: map[string]string{},
			},
		}

		assert.Equal(t, "0.2", doc.Info.Version)
		assert.Nil(t, doc.TodoList)
		assert.NotNil(t, doc.Plan)
		assert.Equal(t, "Test Plan", doc.Plan.Title)
	})
}
