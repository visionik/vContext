# vAgenda Extension: ACE (Agentic Context Engineering)

> **DRAFT EXTENSION**: This document is a draft and subject to change.

**Extension Name**: ACE (Agentic Context Engineering)

**Extension Version**: 0.1

**Last Updated**: 2025-12-27

## Overview

ACE adds evolving playbooks for **long-term memory** (potentially very long-term). While TodoLists provide short-term memory (hours to days) and Plans provide medium-term memory (days to weeks/months), ACE playbooks accumulate durable entries that persist across projects, sessions, and iterations.

This extension is inspired by the ACE paradigm described in "Agentic Context Engineering" (arXiv:2510.04618).

## Dependencies

- **Requires**: Extension 2 (Identifiers)
- **Recommended**: Extension 10 (Version Control & Sync)

Notes:
- If Extension 10 is present, ACE patching can use `baseDocumentSequence` for optimistic concurrency.

## Machine-verifiable schema (JSON)

The ACE extension schema is provided at `schemas/vagenda-extension-ace.schema.json`.

## Data model

### New Types

```javascript
Playbook {
  version: number           # Playbook version (monotonic; increments on update)
  created: datetime         # When playbook was created
  updated: datetime         # Last update time
  entries: PlaybookEntry[]  # Itemized bullets (ACE-aligned)
  metrics?: PlaybookMetrics # Optional summary stats
}

PlaybookEntry {
  id: string                # Unique identifier within the playbook (stable)
  kind: enum                # "strategy" | "learning" | "rule" | "warning" | "note"
  title?: string            # Optional short label
  text: string              # The bullet content (the load-bearing part)
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

AcePatch {
  playbookId?: string            # Optional identifier when patch is shipped separately
  baseDocumentSequence?: number  # Optional guard for optimistic concurrency.
                                # When Extension 10 is in use, this MUST refer to the target document's `sequence` value
                                # (e.g., `plan.sequence` or `todoList.sequence`) at the time the patch was generated.
  operations: AceOp[]
}

AceOp {
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

### TodoList Extensions

```javascript
TodoList {
  // Prior extensions...
  playbook?: Playbook      # Evolving ACE playbook
}
```

### Plan Extensions

```javascript
Plan {
  // Prior extensions...
  playbook?: Playbook      # Evolving ACE playbook
}
```

## ACE invariants (normative)

- Tools SHOULD update playbooks via **localized operations** (AcePatch) rather than rewriting `entries` wholesale.
- `PlaybookEntry.id` MUST be stable once created.
- `appendEntry` operations MUST NOT modify existing entries.
- `incrementCounter` operations MUST be commutative (safe to merge) and SHOULD be used for feedback tracking.
- Deprecation MUST be soft (retain the entry, mark it non-active).
- If `AcePatch.baseDocumentSequence` is present and Extension 10 is in use, an implementation MUST compare it against the current document `sequence`.
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

- Prefer **AcePatch** updates over rewriting `playbook.entries` wholesale.
- Keep entries **atomic** (one idea per entry) to make merge/dedup easier.
- Add **evidence** as soon as you have it (links to PRs, traces, changeLog entries, benchmark results).
- When revising an entry substantially, create a successor entry and connect them via `supersedes/supersededBy` rather than silently editing history.
- Use `duplicateOf` to deduplicate without erasing; keep older entries for provenance.
- Use `status: quarantined` for potentially bad advice instead of deleting; record `deprecatedReason`/notes.
- Treat low-signal feedback (e.g. `feedbackType: selfReport`) as weaker; avoid inflating `confidence` without corroboration.

## Examples

### TRON

```tron
class vAgendaInfo: version
class Plan: id, title, status, narratives, playbook
class Narrative: title, content
class Playbook: version, created, updated, entries
class PlaybookEntry: id, kind, title, text, confidence, helpfulCount, harmfulCount, status

vAgendaInfo: vAgendaInfo("0.2")
plan: Plan(
  "plan-003",
  "API development patterns",
  "completed",
  {"proposal": Narrative("Overview", "Document learned patterns")},
  Playbook(
    1,
    "2024-12-27T09:00:00Z",
    "2024-12-27T15:00:00Z",
    [
      PlaybookEntry(
        "entry-1",
        "strategy",
        "Test-first development",
        "Write tests before implementation for better coverage",
        0.95,
        12,
        0,
        "active"
      ),
      PlaybookEntry(
        "entry-2",
        "learning",
        null,
        "Early validation prevents late-stage refactoring",
        0.9,
        5,
        0,
        "active"
      )
    ]
  )
)
```

### JSON (AcePatch update example)

```json
{
  "baseDocumentSequence": 12,
  "operations": [
    {
      "op": "incrementCounter",
      "entryId": "entry-1",
      "delta": {"helpfulCount": 1},
      "reason": "Strategy applied successfully on login flow implementation"
    },
    {
      "op": "appendEntry",
      "entry": {
        "id": "entry-3",
        "kind": "warning",
        "text": "Avoid rewriting the entire playbook; prefer incremental entry updates to prevent context collapse.",
        "confidence": 0.8,
        "helpfulCount": 1,
        "harmfulCount": 0,
        "status": "active",
        "feedbackType": "executionOutcome",
        "createdAt": "2025-12-27T00:00:00Z",
        "updatedAt": "2025-12-27T00:00:00Z"
      },
      "reason": "Observed repeated loss of details after summary rewrites"
    }
  ]
}
```

## References

- arXiv:2510.04618: https://arxiv.org/abs/2510.04618
