# vAgenda Extension: Playbooks

> **DRAFT EXTENSION**: This document is a draft and subject to change.

**Extension Name**: Playbooks

**Extension Version**: 0.1

**Last Updated**: 2025-12-27

## Overview

Playbooks are a way to make an agent’s “working context” improve over time **without** changing model weights.

If these concepts feel difficult to understand, it’s partly because they are, and partly because this spec still has room to improve. Playbooks are a very new and evolving concept, and initial implementations are only just starting to be attempted and tested.

The playbook concept in this extension is based on the paper "Agentic Context Engineering: Evolving Contexts for Self-Improving Language Models" (arXiv:2510.04618): https://arxiv.org/abs/2510.04618

In vAgenda terms:
- TodoLists cover **short-term memory** (what to do next).
- Plans cover **medium-term memory** (what/why/how for a piece of work).
- Playbooks cover **long-term memory**: reusable strategies, rules-of-thumb, and warnings that persist across runs.

### What playbooks do

Playbooks capture long-term **lessons learned**.

Where TodoLists and Plans primarily record **what** was (or will be) done and **why**, playbooks record higher-level guidance such as:
- what tends to work,
- what tends to fail,
- what to try next,
- and what to avoid.

Importantly, playbooks are designed to be **log-like and reversible**: old guidance is not silently overwritten. Instead, entries are refined, superseded, deduplicated, or quarantined so that tools (and humans) can see how the playbook evolved over time.

Playbooks treat long-term context as a curated set of **entries** that can be:
- appended (growth),
- refined (revision),
- deduplicated,
- and soft-deprecated/quarantined.

This prevents two common failure modes in iterative prompt editing:
- **Context collapse**: repeated rewriting erodes details over time.
- **Brevity bias**: summaries drop “unimportant” details that later turn out to matter.

### How playbooks work (conceptually)

A typical loop:

1. **Execution**
   - Run a task (coding session, workflow, incident response).
2. **Reflection**
   - Record what helped or hurt (increment counters; attach evidence).
3. **Curation**
   - Add new entries, refine existing ones, deduplicate, and quarantine bad advice.
4. **Reuse**
   - Retrieve a relevant subset of entries for the next run.

vAgenda supports this by:
- representing long-term knowledge as `playbook.entries` (stable IDs, atomic entries), and
- supporting incremental updates via `PlaybookPatch` rather than whole-playbook rewrites.


## Dependencies

- **Requires**: Extension 2 (Identifiers)
- **Recommended**: Extension 10 (Version Control & Sync)

Notes:
- If Extension 10 is present, playbook patching can use `baseDocumentSequence` (document `sequence`) for optimistic concurrency.

## Machine-verifiable schema (JSON)

The playbooks extension schema is provided at `schemas/vagenda-extension-playbooks.schema.json`.

## Data model

This extension adds four core concepts:
- `Playbook`: a container for entries and summary metrics.
- `PlaybookEntry`: a single reusable entry (strategy/rule/warning/etc.).
- `PlaybookPatch`: an incremental update envelope.
- `PlaybookPatchOp`: a single operation inside a patch.

### New Types (reference)

```javascript
Playbook {
  version: number           # Playbook version (monotonic; increments on update)
  created: datetime         # When playbook was created
  updated: datetime         # Last update time
  entries: PlaybookEntry[]  # Itemized entries (playbook-aligned)
  metrics?: PlaybookMetrics # Optional summary stats
}

PlaybookEntry {
  id: string                # Unique identifier within the playbook (stable)
  kind: enum                # "strategy" | "learning" | "rule" | "warning" | "note"
  title?: string            # Optional short label
  text: string              # The entry content (the load-bearing part)
  tags?: string[]
  evidence?: string[]       # Human-readable pointers (links, change ids, outcomes)
  confidence?: number       # 0.0-1.0 (optional; omit when unknown)

  helpfulCount?: number     # Count of positive feedback / successful uses
  harmfulCount?: number     # Count of negative feedback / failures
  feedbackType?: enum       # "humanReview" | "executionOutcome" | "selfReport" | "unknown"

  createdAt?: datetime
  updatedAt?: datetime

  status?: enum             # "active" | "deprecated" | "quarantined"
  deprecatedReason?: string
  supersedes?: string[]     # Entry IDs that this entry supersedes
  supersededBy?: string     # Entry ID that supersedes this entry
  duplicateOf?: string      # Entry ID that this entry duplicates (dedup without erasure)

  metadata?: object         # Extension escape hatch
}

PlaybookMetrics {
  totalEntries: number
  averageConfidence?: number
  lastUpdated?: datetime
}

PlaybookPatch {
  playbookId?: string            # Optional identifier when patch is shipped separately
  baseDocumentSequence?: number  # Optional guard for optimistic concurrency.
                                # When Extension 10 is in use, this MUST refer to the target document's `sequence` value
                                # (e.g., `plan.sequence` or `todoList.sequence`) at the time the patch was generated.
  operations: PlaybookPatchOp[]
}

PlaybookPatchOp {
  op: enum                  # "appendEntry" | "updateEntry" | "incrementCounter" | "deprecateEntry"
  entryId?: string          # Target entry id (when applicable)
  entry?: PlaybookEntry     # Full entry for append/update
  delta?: {
    helpfulCount?: number   # Usually +1
    harmfulCount?: number   # Usually +1
  }
  reason?: string           # Why this op is being applied
}
```

### Attaching a playbook to documents

Playbooks can be attached to either a todo list or a plan:

```javascript
TodoList {
  // Prior extensions...
  playbook?: Playbook
}

Plan {
  // Prior extensions...
  playbook?: Playbook
}
```

## Type guide (with examples)

### Playbook

A `Playbook` is the long-lived container. It holds the entry list (`entries`) plus optional summary metrics.

**JSON example:**

```json
{
  "version": 4,
  "created": "2025-01-10T18:00:00Z",
  "updated": "2025-12-27T08:00:00Z",
  "entries": [],
  "metrics": {
    "totalEntries": 0,
    "lastUpdated": "2025-12-27T08:00:00Z"
  }
}
```

### PlaybookEntry

A `PlaybookEntry` is an **atomic** unit of reusable knowledge.

Guidance:
- Keep it to “one idea per entry”.
- Prefer actionable language.
- Add evidence whenever possible.

#### Example: kind = strategy

```json
{
  "id": "entry-test-first",
  "kind": "strategy",
  "title": "Write a failing test first",
  "text": "Before changing code, write a failing test that reproduces the bug; then implement the minimal fix.",
  "tags": ["testing", "debugging"],
  "confidence": 0.95,
  "helpfulCount": 14,
  "harmfulCount": 0,
  "feedbackType": "executionOutcome",
  "status": "active",
  "evidence": ["pr:42", "ci:green-run-2025-12-26"]
}
```

#### Example: kind = rule

```json
{
  "id": "entry-task-first",
  "kind": "rule",
  "text": "For repeatable workflows, add a Task target instead of documenting raw shell commands.",
  "tags": ["workflow", "taskfile"],
  "confidence": 0.9,
  "helpfulCount": 7,
  "harmfulCount": 0,
  "status": "active"
}
```

#### Example: kind = warning

```json
{
  "id": "entry-avoid-blanket-refactors",
  "kind": "warning",
  "text": "Avoid large refactors without a characterization test suite; changes are hard to review and regressions are likely.",
  "tags": ["refactor", "risk"],
  "confidence": 0.85,
  "helpfulCount": 5,
  "harmfulCount": 1,
  "status": "active",
  "evidence": ["incident:2025-09-14-regression"]
}
```

#### Example: kind = learning

```json
{
  "id": "entry-timezone-bugs",
  "kind": "learning",
  "text": "Timezone-related bugs usually come from mixing naive timestamps and offset timestamps; require RFC3339 with offsets everywhere.",
  "tags": ["time", "data-integrity"],
  "confidence": 0.8,
  "helpfulCount": 3,
  "harmfulCount": 0,
  "status": "active"
}
```

#### Example: kind = note

```json
{
  "id": "entry-context",
  "kind": "note",
  "text": "In this repo, the authoritative spec is README.md; extension docs live in vAgenda-extension-*.md.",
  "tags": ["docs"],
  "status": "active"
}
```

#### Example: refinement and dedup

```json
{
  "id": "entry-tests-before-fix-v2",
  "kind": "strategy",
  "text": "Write a failing test first; then implement the minimal change that makes it pass; finally refactor with tests green.",
  "supersedes": ["entry-test-first"],
  "confidence": 0.96,
  "status": "active"
}
```

```json
{
  "id": "entry-test-first-duplicate",
  "kind": "strategy",
  "text": "Always start with a failing test.",
  "duplicateOf": "entry-test-first",
  "status": "deprecated",
  "deprecatedReason": "Duplicate entry; keep canonical entry-test-first"
}
```

### PlaybookMetrics

`PlaybookMetrics` are optional summary fields for UI/tooling.

```json
{
  "totalEntries": 27,
  "averageConfidence": 0.86,
  "lastUpdated": "2025-12-27T08:00:00Z"
}
```

### PlaybookPatch

`PlaybookPatch` is an envelope for incremental updates.

- If Extension 10 is in use, `baseDocumentSequence` SHOULD be set to the current document `sequence`.

```json
{
  "baseDocumentSequence": 12,
  "operations": []
}
```

### PlaybookPatchOp

`PlaybookPatchOp` is a single update operation.

#### appendEntry

```json
{
  "op": "appendEntry",
  "entry": {
    "id": "entry-new",
    "kind": "learning",
    "text": "CI flakes were caused by test ordering; randomize tests locally to reproduce.",
    "status": "active",
    "feedbackType": "executionOutcome",
    "helpfulCount": 1,
    "harmfulCount": 0
  },
  "reason": "Repeated flakes found in nightly runs"
}
```

#### updateEntry

```json
{
  "op": "updateEntry",
  "entryId": "entry-avoid-blanket-refactors",
  "entry": {
    "id": "entry-avoid-blanket-refactors",
    "kind": "warning",
    "text": "Avoid large refactors without a characterization test suite and a rollback plan.",
    "status": "active"
  },
  "reason": "Clarify mitigation steps"
}
```

#### incrementCounter

```json
{
  "op": "incrementCounter",
  "entryId": "entry-task-first",
  "delta": {"helpfulCount": 1},
  "reason": "Task target reduced setup time"
}
```

#### deprecateEntry

```json
{
  "op": "deprecateEntry",
  "entryId": "entry-test-first-duplicate",
  "reason": "Canonicalized into entry-test-first"
}
```

## Playbook invariants (normative)

- Tools SHOULD update playbooks via **localized operations** (PlaybookPatch) rather than rewriting `entries` wholesale.
- `PlaybookEntry.id` MUST be stable once created.
- `appendEntry` operations MUST NOT modify existing entries.
- `incrementCounter` operations MUST be commutative (safe to merge) and SHOULD be used for feedback tracking.
- Deprecation MUST be soft (retain the entry, mark it non-active).
- If `PlaybookPatch.baseDocumentSequence` is present and Extension 10 is in use, an implementation MUST compare it against the current document `sequence`.
  - If they differ, the patch MUST be treated as concurrent and applied via merge/conflict rules (or rejected with a conflict error).

## Merge semantics (deterministic, non-LLM)

When merging concurrent updates:

- `appendEntry` commutes (merge-safe).
- `incrementCounter` commutes (merge-safe; counters add).
- Two concurrent `updateEntry` ops on the same `entryId` are a conflict unless they touch disjoint fields.
- `deprecateEntry` wins over `active` unless explicitly reactivated by a later change.

## Grow-and-refine lifecycle

- **Grow**: append new entries (new ids).
- **Refine**: update an entry in place (same id) and link it using `supersedes/supersededBy` when meaningfully revised.
- **Dedup**: mark redundant entries with `duplicateOf` (do not erase history).

## Best practices

- Prefer **PlaybookPatch** updates over rewriting `playbook.entries` wholesale.
- Keep entries **atomic** (one idea per entry) to make merge/dedup easier.
- Add **evidence** as soon as you have it (links to PRs, traces, changeLog entries, benchmark results).
- When revising an entry substantially, create a successor entry and connect them via `supersedes/supersededBy` rather than silently editing history.
- Use `duplicateOf` to deduplicate without erasing; keep older entries for provenance.
- Use `status: quarantined` for potentially bad advice instead of deleting; record `deprecatedReason`/notes.
- Treat low-signal feedback (e.g. `feedbackType: selfReport`) as weaker; avoid inflating `confidence` without corroboration.

## Real-world playbook examples

These examples show how a playbook might look in practice for an agentic software repo.

### TRON: playbook embedded in a Plan

```tron
class vAgendaInfo: version
class Plan: id, title, status, narratives, playbook
class Narrative: title, content
class Playbook: version, created, updated, entries, metrics
class PlaybookMetrics: totalEntries, averageConfidence, lastUpdated
class PlaybookEntry:
  id, kind, title, text, tags, evidence, confidence,
  helpfulCount, harmfulCount, feedbackType,
  status, deprecatedReason, supersedes, supersededBy, duplicateOf

vAgendaInfo: vAgendaInfo("0.2")
plan: Plan(
  "plan-ace-realworld-001",
  "Agent workflow rules",
  "inProgress",
  {
    "proposal": Narrative(
      "Overview",
      "Curate reusable agent workflow strategies and warnings; update via PlaybookPatch."
    )
  },
  Playbook(
    7,
    "2025-01-10T18:00:00Z",
    "2025-12-27T08:00:00Z",
    [
      PlaybookEntry(
        "entry-task-first",
        "rule",
        "Task-first workflow",
        "For repeatable workflows, add a Task target instead of documenting raw shell commands.",
        ["workflow", "taskfile"],
        ["docs:warp.md"],
        0.9,
        7,
        0,
        "humanReview",
        "active",
        null,
        null,
        null,
        null
      ),
      PlaybookEntry(
        "entry-test-first",
        "strategy",
        "Failing test first",
        "Before changing code, write a failing test that reproduces the bug; then implement the minimal fix.",
        ["testing", "debugging"],
        ["pr:42", "ci:green-run-2025-12-26"],
        0.95,
        14,
        0,
        "executionOutcome",
        "active",
        null,
        null,
        "entry-tests-before-fix-v2",
        null
      ),
      PlaybookEntry(
        "entry-tests-before-fix-v2",
        "strategy",
        "Refined test-first",
        "Write a failing test first; then implement the minimal change that makes it pass; finally refactor with tests green.",
        ["testing"],
        ["pr:57"],
        0.96,
        3,
        0,
        "executionOutcome",
        "active",
        null,
        ["entry-test-first"],
        null,
        null
      ),
      PlaybookEntry(
        "entry-avoid-blanket-refactors",
        "warning",
        null,
        "Avoid large refactors without a characterization test suite and a rollback plan.",
        ["refactor", "risk"],
        ["incident:2025-09-14-regression"],
        0.85,
        5,
        1,
        "executionOutcome",
        "active",
        null,
        null,
        null,
        null
      ),
      PlaybookEntry(
        "entry-timezone-bugs",
        "learning",
        null,
        "Timezone-related bugs usually come from mixing naive timestamps and offset timestamps; require RFC3339 with offsets everywhere.",
        ["time", "data-integrity"],
        ["issue:113"],
        0.8,
        3,
        0,
        "executionOutcome",
        "active",
        null,
        null,
        null,
        null
      )
    ],
    PlaybookMetrics(5, 0.89, "2025-12-27T08:00:00Z")
  )
)
```

### JSON: playbook embedded in a Plan

```json
{
  "vAgendaInfo": {
    "version": "0.2"
  },
  "plan": {
    "id": "plan-ace-realworld-001",
    "title": "Agent workflow rules",
    "status": "inProgress",
    "narratives": {
      "proposal": {
        "title": "Overview",
        "content": "Curate reusable agent workflow strategies and warnings; update via PlaybookPatch."
      }
    },
    "playbook": {
      "version": 7,
      "created": "2025-01-10T18:00:00Z",
      "updated": "2025-12-27T08:00:00Z",
      "entries": [
        {
          "id": "entry-task-first",
          "kind": "rule",
          "title": "Task-first workflow",
          "text": "For repeatable workflows, add a Task target instead of documenting raw shell commands.",
          "tags": ["workflow", "taskfile"],
          "evidence": ["docs:warp.md"],
          "confidence": 0.9,
          "helpfulCount": 7,
          "harmfulCount": 0,
          "feedbackType": "humanReview",
          "status": "active"
        },
        {
          "id": "entry-test-first",
          "kind": "strategy",
          "title": "Failing test first",
          "text": "Before changing code, write a failing test that reproduces the bug; then implement the minimal fix.",
          "tags": ["testing", "debugging"],
          "evidence": ["pr:42", "ci:green-run-2025-12-26"],
          "confidence": 0.95,
          "helpfulCount": 14,
          "harmfulCount": 0,
          "feedbackType": "executionOutcome",
          "status": "active",
          "supersededBy": "entry-tests-before-fix-v2"
        },
        {
          "id": "entry-tests-before-fix-v2",
          "kind": "strategy",
          "title": "Refined test-first",
          "text": "Write a failing test first; then implement the minimal change that makes it pass; finally refactor with tests green.",
          "tags": ["testing"],
          "evidence": ["pr:57"],
          "confidence": 0.96,
          "helpfulCount": 3,
          "harmfulCount": 0,
          "feedbackType": "executionOutcome",
          "status": "active",
          "supersedes": ["entry-test-first"]
        },
        {
          "id": "entry-avoid-blanket-refactors",
          "kind": "warning",
          "text": "Avoid large refactors without a characterization test suite and a rollback plan.",
          "tags": ["refactor", "risk"],
          "evidence": ["incident:2025-09-14-regression"],
          "confidence": 0.85,
          "helpfulCount": 5,
          "harmfulCount": 1,
          "feedbackType": "executionOutcome",
          "status": "active"
        },
        {
          "id": "entry-timezone-bugs",
          "kind": "learning",
          "text": "Timezone-related bugs usually come from mixing naive timestamps and offset timestamps; require RFC3339 with offsets everywhere.",
          "tags": ["time", "data-integrity"],
          "evidence": ["issue:113"],
          "confidence": 0.8,
          "helpfulCount": 3,
          "harmfulCount": 0,
          "feedbackType": "executionOutcome",
          "status": "active"
        }
      ],
      "metrics": {
        "totalEntries": 5,
        "averageConfidence": 0.89,
        "lastUpdated": "2025-12-27T08:00:00Z"
      }
    }
  }
}
```

## References

- https://arxiv.org/abs/2510.04618
