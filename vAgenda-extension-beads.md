# vAgenda Extension Proposal: Beads Interoperability

> **VERY EARLY DRAFT**: This is an initial proposal and subject to significant change. Comments, feedback, and suggestions are strongly encouraged. Please provide input via GitHub issues or discussions.

**Extension Name**: Beads Interop  
**Version**: 0.1 (Draft)  
**Status**: Proposal  
**Author**: Jonathan Taylor (visionik@pobox.com)  
**Date**: 2025-12-27

## Overview

[Beads](https://github.com/steveyegge/beads) is a git-backed issue tracker created by Steve Yegge that provides persistent, structured memory for coding agents. It replaces messy markdown plans with a dependency-aware graph stored in `.beads/` as JSONL, allowing agents to handle long-horizon tasks without losing context between sessions. Beads focuses on "what matters right now" - current work, what just finished, and what's blocked - rather than long-term planning or historical documentation.

This extension enables bidirectional interoperability between vAgenda and Beads. It bridges vAgenda's standardized memory format with Beads' execution-focused task tracking.

## Motivation

**Beads strengths**:
- Real-time task execution tracking
- Dependency-aware graph for current work
- Agent-optimized queries (`bd ready`, `bd blocked`)
- Git-based storage with automatic sync
- Focus on "what matters right now"

**vAgenda strengths**:
- Standardized format for cross-system interop
- Rich narratives for plans (the "why" and "how")
- Long-term memory via playbooks
- Token-efficient TRON encoding
- Separation of short/medium/long-term memory

**Integration goal**: Let Beads handle execution tracking while vAgenda provides knowledge persistence and transfer.

## Dependencies

**Required**:
- Extension 2 (Identifiers) - for beads ID mapping
- Extension 4 (Hierarchical) - for dependency tracking

**Recommended**:
- Extension 10 (Version Control & Sync) - for change tracking
- Extension 12 (Playbooks) - for accumulating learnings

## New Fields

### TodoItem Extensions
```javascript
TodoItem {
  // Prior extensions...
  beadsId?: string            # Beads issue ID (e.g., "bd-a1b2c3d4")
  beadsStatus?: string        # Beads-specific status
  beadsUrl?: string           # Link to bead in viewer
  beadsMetrics?: BeadsMetrics # Graph metrics from beads_viewer
}
```

### Plan Extensions
```javascript
Plan {
  // Prior extensions...
  beadsProject?: string       # Associated Beads project path
  beadsSyncedAt?: datetime    # Last sync timestamp
}
```

### Phase Extensions
```javascript
Phase {
  // Prior extensions...
  beadsId?: string           # Beads issue ID if phase is tracked
}
```

### New Types

```javascript
BeadsMetrics {
  pageRank?: number          # From beads_viewer --robot-insights
  betweenness?: number       # Betweenness centrality
  impactScore?: number       # Impact on project
  isBottleneck?: boolean     # Flagged as bottleneck
  inCycle?: boolean          # Part of dependency cycle
}
```

## Mapping Semantics

### TodoItem ↔ Beads Issue

**vAgenda TodoItem** maps to **Beads issue** when:
- Task is in active execution (not future planning)
- Task needs dependency tracking
- Task will span multiple agent sessions

**Mapping rules**:
```
vAgenda status     →  Beads status
-----------------------------------
pending           →  open
inProgress        →  open (with activity)
blocked           →  blocked
completed         →  closed
cancelled         →  closed (with note)
```

**Sync direction**:
- Beads → vAgenda: After each agent session (capture execution state)
- vAgenda → Beads: When planning new work (populate from plans)

### Plan → Beads Project Context

Plans don't map directly to Beads but provide **context**:
- Plan narratives become AGENTS.md instructions
- Plan phases become milestone-tagged beads
- Plan learnings feed playbooks

## Usage Patterns

### Pattern 1: Session Handoff

**End of agent session** (Beads → vAgenda):
```bash
# Agent "lands the plane"
bd export --format=vagenda > session-2025-12-27.tron

# This creates vAgenda TodoList with:
# - Current beads as TodoItems (with beadsId)
# - Dependency graph preserved via Extension 4
# - BeadsMetrics for prioritization
```

**Start of next session** (vAgenda → Beads):
```bash
# Import updated priorities/context
bd import --format=vagenda session-2025-12-27.tron

# Agent sees:
# - New items added by human during planning
# - Updated priorities from Plan review
# - Context from Plan narratives
```

### Pattern 2: Long-Term Knowledge Extraction

**Weekly/milestone**:
```bash
# Extract completed work into Plan
bd export --format=vagenda --closed-since=7d > week-retro.tron

# This becomes a Plan with:
# - Narrative documenting what was accomplished
# - Reflections on what worked/didn't
# - Feeds into playbooks for future projects
```

### Pattern 3: Multi-Project Learning

**Playbook accumulation**:
```
Project A (beads) → vAgenda Plan + Playbooks
Project B (beads) → vAgenda Plan + Playbooks
Project C (beads) → vAgenda Plan + Playbooks
                          ↓
            Accumulated strategies/learnings
            persist across all projects
```

## Implementation Notes

### For Beads Developers

Add vAgenda export format to `bd export`:
```go
// In bd export command
case "vagenda":
    exporter := vagenda.NewExporter(db)
    return exporter.ExportTodoList(filter)
```

**Minimal export** (just TodoList):
```tron
class vAgendaInfo: version
class TodoList: items
class TodoItem: id, title, status, beadsId, dependencies

vAgendaInfo: vAgendaInfo("0.2")
todoList: TodoList([
  TodoItem("1", "Fix auth bug", "inProgress", "bd-a1b2", []),
  TodoItem("2", "Add tests", "pending", "bd-c3d4", ["1"])
])
```

**Full export** (with metrics):
```tron
class BeadsMetrics: pageRank, impactScore, isBottleneck
class TodoItem: id, title, status, beadsId, dependencies, beadsMetrics

// ... items with BeadsMetrics populated from bv --robot-insights
```

### For vAgenda Tools

Accept Beads-tagged items:
```javascript
// When rendering TodoList, recognize beadsId
if (item.beadsId) {
  // Show link to beads_viewer
  // Highlight if isBottleneck or inCycle
  // Use impactScore for priority ordering
}
```

## Relationship to Existing Extensions

**vs Extension 10 (Version Control)**:
- Extension 10 tracks vAgenda document changes
- Beads Interop tracks sync with Beads execution state
- They complement: Beads tracks task execution, Ext 10 tracks plan evolution

**vs Extension 12 (Playbooks)**:
- Playbooks capture long-term learnings
- Beads provides execution data that feeds playbooks
- Workflow: Beads execution → vAgenda reflection → playbook accumulation

## Open Questions

1. **Should Beads native format adopt vAgenda?**
   - Pro: Standard format, better interop
   - Con: Breaking change, JSONL is working well
   - **Proposal**: Beads keeps JSONL, adds vAgenda as export/import option

2. **How to handle Beads' "memory decay"?**
   - Beads compacts old issues to save context
   - vAgenda Plans could preserve full context
   - **Proposal**: Beads compacts for agents, vAgenda archives for humans

3. **Multi-agent coordination?**
   - Beads uses git for multi-agent sync
   - vAgenda Extension 11 (Forking) handles parallel work
   - **Proposal**: Beads handles real-time, vAgenda handles merge/conflict resolution

## Examples

### Example 1: Simple Export

**Beads state**:
```
bd-a1b2: [open] Implement JWT auth
bd-c3d4: [open] Add auth tests (depends: bd-a1b2)
bd-e5f6: [closed] Setup database
```

**vAgenda TodoList**:
```tron
class vAgendaInfo: version
class TodoList: items
class TodoItem: id, title, status, beadsId, dependencies

vAgendaInfo: vAgendaInfo("0.2")
todoList: TodoList([
  TodoItem("1", "Implement JWT auth", "inProgress", "bd-a1b2", []),
  TodoItem("2", "Add auth tests", "pending", "bd-c3d4", ["1"])
])
```

### Example 2: With Metrics

```tron
class vAgendaInfo: version, author
class TodoList: id, items
class TodoItem: id, title, status, beadsId, dependencies, beadsMetrics
class BeadsMetrics: impactScore, isBottleneck

vAgendaInfo: vAgendaInfo("0.2", "agent-alpha")
todoList: TodoList(
  "session-123",
  [
    TodoItem(
      "1",
      "Refactor login module",
      "inProgress",
      "bd-a1b2",
      [],
      BeadsMetrics(0.85, true)  // High impact, bottleneck
    ),
    TodoItem(
      "2",
      "Update docs",
      "pending",
      "bd-c3d4",
      ["1"],
      BeadsMetrics(0.12, false)  // Low impact, not blocking
    )
  ]
)
```

### Example 3: Learning Extraction

**After completing auth work in Beads**:

```tron
class vAgendaInfo: version
class Plan: title, status, narratives, playbook
class Narrative: title, content
class Playbook: version, strategies, learnings
class Strategy: id, title, description, confidence
class Learning: id, content, confidence, discoveredBy
class Agent: id, type, name

vAgendaInfo: vAgendaInfo("0.2")
plan: Plan(
  "Authentication Implementation Retrospective",
  "completed",
  {
    "proposal": Narrative(
      "Summary",
      "Implemented JWT auth across 8 beads over 3 days"
    ),
    "testing": Narrative(
      "What Worked",
      "Test-first approach caught edge cases early. Beads dependency tracking prevented premature integration."
    ),
    "risks": Narrative(
      "What Didn't",
      "Initial token expiry too short. Had to refactor after testing with real workflows."
    )
  },
  Playbook(
    1,
    [
      Strategy(
        "strat-auth-1",
        "Test token lifecycle early",
        "Don't wait until integration to test full token lifecycle including refresh",
        0.9
      )
    ],
    [
      Learning(
        "learn-auth-1",
        "Beads dependency tracking prevented premature auth integration",
        0.95,
        Agent("agent-alpha", "aiAgent", "Claude")
      )
    ]
  )
)
```

## Migration Path

**Phase 1**: Export only (Beads → vAgenda)
- Add `bd export --format=vagenda`
- Agents can archive completed work to vAgenda Plans
- No breaking changes

**Phase 2**: Import support (vAgenda → Beads)
- Add `bd import --format=vagenda`
- Enable human planning → agent execution workflow
- Still no breaking changes

**Phase 3**: Native integration
- Beads tools read vAgenda Plans for context
- beads_viewer can render vAgenda Plans alongside beads
- vAgenda tools can query live Beads state

## Community Feedback

This is a **draft proposal**. Feedback needed:

1. Is the mapping TodoItem ↔ Bead too simplistic?
2. Should Plans map to Beads "projects" or remain separate?
3. Are BeadsMetrics the right graph properties to expose?
4. Should this be a formal vAgenda extension or a Beads feature?

**Discuss**: https://github.com/visionik/vAgenda/discussions

## References

- vAgenda Specification: https://github.com/visionik/vAgenda
- Beads: https://github.com/steveyegge/beads
- beads_viewer: https://github.com/Dicklesworthstone/beads_viewer
- Steve Yegge's article: https://steve-yegge.medium.com/introducing-beads-a-coding-agent-memory-system-637d7d92514a
