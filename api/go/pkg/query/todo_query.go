// Package query provides filtering and querying capabilities for vAgenda documents.
package query

import "github.com/visionik/vAgenda/api/go/pkg/core"

// TodoQuery provides filtering for TodoItems.
type TodoQuery struct {
	items []core.TodoItem
}

// NewTodoQuery creates a new query for the given items.
func NewTodoQuery(items []core.TodoItem) *TodoQuery {
	return &TodoQuery{items: items}
}

// ByStatus filters items by status.
func (q *TodoQuery) ByStatus(status core.ItemStatus) *TodoQuery {
	filtered := make([]core.TodoItem, 0)
	for _, item := range q.items {
		if item.Status == status {
			filtered = append(filtered, item)
		}
	}
	return &TodoQuery{items: filtered}
}

// ByTitle filters items by title substring (case-insensitive).
func (q *TodoQuery) ByTitle(substring string) *TodoQuery {
	filtered := make([]core.TodoItem, 0)
	for _, item := range q.items {
		if contains(item.Title, substring) {
			filtered = append(filtered, item)
		}
	}
	return &TodoQuery{items: filtered}
}

// Where filters items using a custom predicate function.
func (q *TodoQuery) Where(predicate func(core.TodoItem) bool) *TodoQuery {
	filtered := make([]core.TodoItem, 0)
	for _, item := range q.items {
		if predicate(item) {
			filtered = append(filtered, item)
		}
	}
	return &TodoQuery{items: filtered}
}

// All returns all matching items.
func (q *TodoQuery) All() []core.TodoItem {
	return q.items
}

// First returns the first matching item, or nil if none match.
func (q *TodoQuery) First() *core.TodoItem {
	if len(q.items) > 0 {
		return &q.items[0]
	}
	return nil
}

// Count returns the number of matching items.
func (q *TodoQuery) Count() int {
	return len(q.items)
}

// Any returns true if there are any matching items.
func (q *TodoQuery) Any() bool {
	return len(q.items) > 0
}

// contains performs case-insensitive substring search.
func contains(s, substr string) bool {
	// Simple case-insensitive contains
	sLower := toLower(s)
	substrLower := toLower(substr)
	return indexOf(sLower, substrLower) >= 0
}

// toLower converts a string to lowercase.
func toLower(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + 32
		} else {
			result[i] = r
		}
	}
	return string(result)
}

// indexOf returns the index of substr in s, or -1 if not found.
func indexOf(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
