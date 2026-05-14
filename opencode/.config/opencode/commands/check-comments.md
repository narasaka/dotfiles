---
description: Audit comments in current branch diff for AI slop and unnecessary noise
subtask: true
---

Analyze the diff below for comments ADDED on this branch (`+` lines only, ignore `-` lines).

Base branch: `$ARGUMENTS` if provided, otherwise auto-detected.

<diff>
!`BASE=${1:-$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's@refs/remotes/origin/@@' || echo "main")}; git diff "$(git merge-base HEAD "origin/$BASE")"...HEAD`
</diff>

## Classification

For every comment in added lines, classify it into one of these categories.

### SKIP (do not report)
- BDD keywords: `given`, `when`, `then`, `arrange`, `act`, `assert`
- Linter/type directives: `noqa`, `type:`, `pyright:`, `eslint-disable`, `ts-ignore`, `ts-expect-error`, `prettier-ignore`, `clippy:`, `allow`, `deny`, `warn`, `forbid`, `ruff:`, `mypy:`, `pylint:`, `flake8:`, `pyre:`, `pytype:`
- Shebangs: `#!`
- Legal/license headers

### AGENT MEMO (highest severity)
Comments that describe what was changed rather than why. Typical AI agent behavior:
- Change tracking: "Changed from/to", "Modified to", "Updated from/to", "Was changed"
- Action narration: "Refactored", "Replaced", "Removed", "Deleted", "Added", "Implemented"
- Movement: "Moved from/to", "Renamed from/to", "Converted from/to", "Migrated from/to"
- Self-narration: "This implements/adds/removes", "Here we", "Now we/this/it"
- Temporal: "Previously", "Before this", "After this"
- Meta: "Note:", "Implementation of"
- Arrow notation: "oldThing -> newThing"

### UNNECESSARY
- Restates what the code says (`// increment counter` before `counter++`)
- Vague TODOs (`// TODO: fix later`, `// TODO: refactor this`)
- Commented-out code
- Obvious docstrings on self-explanatory functions
- Verbose explanations that the code already makes clear

### JUSTIFIED
- Complex algorithm or business logic explanation
- Security-critical context
- Performance optimization rationale
- Regex or math formula explanation
- Public API docs on non-obvious interfaces
- Links to external specs, RFCs, tickets

## Output

Report by file. For each comment:

```
FILE: <path>
LINE: <number>
CATEGORY: AGENT_MEMO | UNNECESSARY | JUSTIFIED
TEXT: <comment text>
REASON: <one sentence>
ACTION: REMOVE | KEEP
```

End with:

```
SUMMARY:
  Agent memos: N (remove)
  Unnecessary: N (remove)
  Justified: N (keep)
  Total: N
```

If clean, just say: "Clean branch. No unnecessary comments found."

After the summary, ask: "Want me to remove the unnecessary comments?"

Do NOT read files, make edits, or run any tools unless the user confirms removal.
