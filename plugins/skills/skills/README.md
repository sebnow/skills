# Skills

Each skill is a subdirectory containing a `SKILL.md`:

```
skills/
└── <skill-name>/
    └── SKILL.md
```

`SKILL.md` requires YAML frontmatter with a `description` (and, when the
skill is a standalone file, a `name`). Claude uses the `description` to
decide when to invoke the skill, so make it specific about triggers.

Minimal example:

```markdown
---
name: <skill-name>
description: What the skill does and when to use it.
---

Instructions go here.
```
