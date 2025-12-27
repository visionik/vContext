package main

import (
	"fmt"
	"log"

	"github.com/visionik/vAgenda/api/go/pkg/builder"
	"github.com/visionik/vAgenda/api/go/pkg/convert"
	"github.com/visionik/vAgenda/api/go/pkg/core"
	"github.com/visionik/vAgenda/api/go/pkg/parser"
	"github.com/visionik/vAgenda/api/go/pkg/query"
	"github.com/visionik/vAgenda/api/go/pkg/validator"
)

func main() {
	fmt.Println("=== vAgenda Go Library Examples ===\n")

	// Example 1: Build a TodoList
	fmt.Println("Example 1: Building a TodoList")
	todoDoc := builder.NewTodoList("0.2").
		WithAuthor("agent-alpha").
		AddPendingItem("Implement authentication").
		AddPendingItem("Write API documentation").
		AddInProgressItem("Setup database").
		Build()

	// Convert to JSON
	jsonData, err := convert.ToJSONIndent(todoDoc, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("JSON:")
	fmt.Println(string(jsonData))
	fmt.Println()

	// Convert to TRON
	tronData, err := convert.ToTRON(todoDoc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("TRON:")
	fmt.Println(string(tronData))
	fmt.Println()

	// Example 2: Build a Plan
	fmt.Println("Example 2: Building a Plan")
	planDoc := builder.NewPlan("Add user authentication", "0.2").
		WithAuthor("team-lead").
		WithStatus(core.PlanStatusDraft).
		WithProposal("Proposed Changes", "Implement JWT-based authentication with refresh tokens").
		WithProblem("Problem Statement", "Current system lacks secure authentication").
		AddPendingPhase("Database setup").
		AddInProgressPhase("JWT implementation").
		AddPendingPhase("OAuth integration").
		Build()

	planJSON, err := convert.ToJSONIndent(planDoc, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(planJSON))
	fmt.Println()

	// Example 3: Parse and Query
	fmt.Println("Example 3: Parsing and Querying")
	jsonParser := parser.NewJSONParser()
	parsed, err := jsonParser.ParseBytes(jsonData)
	if err != nil {
		log.Fatal(err)
	}

	// Query pending items
	q := query.NewTodoQuery(parsed.TodoList.Items)
	pendingItems := q.ByStatus(core.StatusPending).All()

	fmt.Printf("Found %d pending items:\n", len(pendingItems))
	for _, item := range pendingItems {
		fmt.Printf("  - %s\n", item.Title)
	}
	fmt.Println()

	// Example 4: Validation
	fmt.Println("Example 4: Validation")
	v := validator.NewValidator()
	if err := v.Validate(todoDoc); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("✓ TodoList document is valid")
	}

	if err := v.Validate(planDoc); err != nil {
		fmt.Printf("Validation failed: %v\n", err)
	} else {
		fmt.Println("✓ Plan document is valid")
	}
}
