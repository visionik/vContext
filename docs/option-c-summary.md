# Option C Implementation Summary

**Date**: 2025-12-27  
**Decision**: Use `Item` as abstract base class for all contained entities

## What Changed

### 1. Added Item Abstraction

**Core atomic class:**
```tron
class Item: Titled, Statused
```

All contained entities now extend Item:
- **TodoItem** - extends Item (unchanged name)
- **PlanItem** - extends Item (renamed from Phase)
- **PlaybookItem** - extends Item (renamed from PlaybookEntry)

### 2. Renamed Types

| Old Name | New Name | Reason |
|----------|----------|--------|
| `Phase` | `PlanItem` | Follow `<Container>Item` convention |
| `phase.subPhases` | `planItem.subItems` | Consistency |
| `plan.phases` | `plan.items` | All containers use `items` |
| `PlaybookEntry` | `PlaybookItem` | Follow `<Container>Item` convention |
| `playbook.entries` | `playbook.items` | All containers use `items` |
| `TodoItem` | `TodoItem` | Already correct (unchanged) |
| `todoList.items` | `todoList.items` | Already correct (unchanged) |

### 3. Container Consistency

**All containers now follow the same pattern:**
```tron
TodoList {
  items: TodoItem[]     # Collection of Items
}

Plan {
  items: PlanItem[]     # Collection of Items
}

Playbook {
  items: PlaybookItem[] # Collection of Items
}
```

## Benefits

### 1. Clear Abstraction
```tron
# Any code that works with Item works with all contained types
class Item: Titled, Statused

function filterActiveItems(items: Item[]): Item[] {
  return items.filter(item => item.status !== "cancelled")
}

# Works for TodoItem, PlanItem, PlaybookItem
```

### 2. Consistent Naming
- `<Container>Item` pattern is explicit
- TodoList → TodoItem ✓
- Plan → PlanItem ✓
- Playbook → PlaybookItem ✓

### 3. Polymorphic Operations
```typescript
// Generic search across any Item type
function searchItems(query: string, items: Item[]): Item[] {
  return items.filter(item => 
    item.title.includes(query) || 
    item.description?.includes(query)
  )
}

// Works for mixed collections
const allItems: Item[] = [
  ...todoList.items,
  ...plan.items,
  ...playbook.items
]
searchItems("security", allItems)
```

### 4. Future-Proof Pattern
New container types follow obvious convention:
```tron
# Future: Workflow container
Workflow {
  items: WorkflowItem[]
}
class WorkflowItem: Item, ...
```

## Breaking Changes

### Phase → PlanItem

**Impact**: All Plans with phases need migration

**Before:**
```json
{
  "plan": {
    "phases": [
      {
        "title": "Foundation",
        "status": "completed",
        "subPhases": [...]
      }
    ]
  }
}
```

**After:**
```json
{
  "plan": {
    "items": [
      {
        "title": "Foundation",
        "status": "completed",
        "subItems": [...]
      }
    ]
  }
}
```

### PlaybookEntry → PlaybookItem

**Impact**: All Playbooks need migration

**Before:**
```json
{
  "playbook": {
    "entries": [
      {
        "eventId": "evt-1",
        "targetId": "entry-1",
        "title": "Test First"
      }
    ]
  }
}
```

**After:**
```json
{
  "playbook": {
    "items": [
      {
        "eventId": "evt-1",
        "targetId": "entry-1",
        "title": "Test First"
      }
    ]
  }
}
```

## Migration Code

```typescript
// Auto-migration function
function migrateToOptionC(doc: any): any {
  // Migrate Plan
  if (doc.plan?.phases) {
    doc.plan.items = doc.plan.phases.map((phase: any) => ({
      ...phase,
      subItems: phase.subPhases
    }))
    delete doc.plan.phases
  }
  
  // Migrate Playbook
  if (doc.playbook?.entries) {
    doc.playbook.items = doc.playbook.entries
    delete doc.playbook.entries
  }
  
  return doc
}
```

## TRON Examples

### Core (minimal)
```tron
class Item: Titled, Statused

# TodoItem already uses Item pattern
class TodoItemCore: Item
TodoItemCore("Fix bug", "pending")

# PlanItem now uses Item pattern
class PlanItemCore: Item
PlanItemCore("Foundation", "completed")

# PlaybookItem now uses Item pattern
class PlaybookItemCore: Item, EventIdentifiable, operation
PlaybookItemCore("Test First", "active", "evt-1", "target-1", "append")
```

### Full (with extensions)
```tron
class Item: Titled, Statused

class TodoItemFull: Item, Identifiable, Timestamped, Tagged, Extensible, priority
class PlanItemFull: Item, Identifiable, Timestamped, Tagged, Extensible, dependencies, subItems, todoList
class PlaybookItemFull: Item, EventIdentifiable, TimestampedSingle, Evidenced, Extensible, operation, ...
```

## Tagged Addition

**Also added Tagged atomic to Extension 3:**

```tron
class Tagged: tags
```

**All entities MAY be Tagged:**
- TodoList, TodoItem
- Plan, PlanItem
- Playbook, PlaybookItem
- ProblemModel

**Example:**
```tron
TodoItem("Fix auth", "inProgress", ["security", "p0", "backend"])
Plan("Migration", "approved", {...}, ["architecture", "q1-2025"])
PlanItem("Database", "completed", ["infrastructure", "postgres"])
Playbook("Patterns", [...], ["backend", "best-practices"])
```

## Documentation Updates

Updated in `atomic-classes-proposal.md`:
1. ✓ Added Item to Core Tier atomics
2. ✓ Renamed all Phase → PlanItem examples
3. ✓ Renamed all PlaybookEntry → PlaybookItem examples
4. ✓ Updated container fields (phases → items, entries → items)
5. ✓ Added Tagged atomic definition
6. ✓ Added migration guide with code examples
7. ✓ Updated field mapping table

## Next Steps

1. **Update main README.md** with Item abstraction
2. **Update JSON Schemas** with new type names
3. **Create migration tools** for existing documents
4. **Update all examples** throughout documentation
5. **Implement deprecation warnings** in v0.3
6. **Full migration** in v0.5

## Rollout Strategy

### Version 0.3 (Current)
- Introduce Item abstraction
- Document PlanItem and PlaybookItem as aliases
- Support both old and new names (Phase/PlanItem, PlaybookEntry/PlaybookItem)
- Add deprecation warnings for old names
- Provide migration utilities

### Version 0.4
- Default to new names in all tooling
- Auto-migrate on read
- Write only new format
- Loud warnings for old names

### Version 0.5
- Drop support for old names
- Phase and PlaybookEntry no longer recognized
- Migration required for old documents
