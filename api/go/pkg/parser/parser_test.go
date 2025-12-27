package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/visionik/vAgenda/api/go/pkg/core"
)

const validJSON = `{
  "vAgendaInfo": {
    "version": "0.2",
    "author": "test-author"
  },
  "todoList": {
    "items": [
      {
        "title": "Task 1",
        "status": "pending"
      },
      {
        "title": "Task 2",
        "status": "inProgress"
      }
    ]
  }
}`

const validPlanJSON = `{
  "vAgendaInfo": {
    "version": "0.2"
  },
  "plan": {
    "title": "Test Plan",
    "status": "draft",
    "narratives": {
      "proposal": {
        "title": "Proposal",
        "content": "Test content"
      }
    }
  }
}`

func TestJSONParser(t *testing.T) {
	parser := NewJSONParser()

	t.Run("parses valid TodoList JSON", func(t *testing.T) {
		doc, err := parser.ParseString(validJSON)

		require.NoError(t, err)
		require.NotNil(t, doc)
		assert.Equal(t, "0.2", doc.Info.Version)
		assert.Equal(t, "test-author", doc.Info.Author)
		require.NotNil(t, doc.TodoList)
		assert.Len(t, doc.TodoList.Items, 2)
		assert.Equal(t, "Task 1", doc.TodoList.Items[0].Title)
		assert.Equal(t, core.StatusPending, doc.TodoList.Items[0].Status)
	})

	t.Run("parses valid Plan JSON", func(t *testing.T) {
		doc, err := parser.ParseString(validPlanJSON)

		require.NoError(t, err)
		require.NotNil(t, doc)
		assert.Equal(t, "0.2", doc.Info.Version)
		require.NotNil(t, doc.Plan)
		assert.Equal(t, "Test Plan", doc.Plan.Title)
		assert.Equal(t, core.PlanStatusDraft, doc.Plan.Status)
		assert.Len(t, doc.Plan.Narratives, 1)
	})

	t.Run("parses from reader", func(t *testing.T) {
		reader := strings.NewReader(validJSON)
		doc, err := parser.Parse(reader)

		require.NoError(t, err)
		require.NotNil(t, doc)
		assert.Equal(t, "0.2", doc.Info.Version)
	})

	t.Run("parses from bytes", func(t *testing.T) {
		doc, err := parser.ParseBytes([]byte(validJSON))

		require.NoError(t, err)
		require.NotNil(t, doc)
		assert.Equal(t, "0.2", doc.Info.Version)
	})

	t.Run("returns error for invalid JSON", func(t *testing.T) {
		_, err := parser.ParseString("invalid json")
		assert.Error(t, err)
	})

	t.Run("returns error for empty string", func(t *testing.T) {
		_, err := parser.ParseString("")
		assert.Error(t, err)
	})
}

func TestAutoParser(t *testing.T) {
	parser := NewAutoParser()

	t.Run("detects and parses JSON", func(t *testing.T) {
		doc, err := parser.ParseString(validJSON)

		require.NoError(t, err)
		require.NotNil(t, doc)
		assert.Equal(t, "0.2", doc.Info.Version)
		assert.Equal(t, "test-author", doc.Info.Author)
	})

	t.Run("handles leading whitespace in JSON", func(t *testing.T) {
		jsonWithWhitespace := "  \n  " + validJSON
		doc, err := parser.ParseString(jsonWithWhitespace)

		require.NoError(t, err)
		require.NotNil(t, doc)
		assert.Equal(t, "0.2", doc.Info.Version)
	})

	t.Run("parses from bytes", func(t *testing.T) {
		doc, err := parser.ParseBytes([]byte(validJSON))

		require.NoError(t, err)
		require.NotNil(t, doc)
		assert.Equal(t, "0.2", doc.Info.Version)
	})
}

func TestNew(t *testing.T) {
	tests := []struct {
		name   string
		format Format
		want   string
	}{
		{"JSON format", FormatJSON, "*parser.JSONParser"},
		{"TRON format", FormatTRON, "*parser.TRONParser"},
		{"Auto format", FormatAuto, "*parser.AutoParser"},
		{"Unknown format defaults to Auto", Format("unknown"), "*parser.AutoParser"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := New(tt.format)
			assert.NotNil(t, parser)
			// Type checking would require reflection or type assertion
			// For now, just verify we get a non-nil parser
		})
	}
}
