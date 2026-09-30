# sebnow skills marketplace

A Claude Code plugin marketplace that distributes a personal collection of
skills. This repo is both the marketplace (`.claude-plugin/marketplace.json`)
and the host of its one plugin (`plugins/skills/`).

## Layout

```
.
├── .claude-plugin/
│   └── marketplace.json          # marketplace catalog
└── plugins/
    └── skills/                   # the "skills" plugin
        ├── .claude-plugin/
        │   └── plugin.json
        ├── agents/
        │   └── <agent-name>.md
        ├── skills/
        │   └── <skill-name>/
        │       └── SKILL.md
        └── workflows/
            └── <workflow-name>.js
```

## Use

Add the marketplace, then install the plugin:

```sh
claude plugin marketplace add <path-or-git-url>
claude plugin install sebnow@sebnow
```

## Adding a skill

Create `plugins/skills/skills/<skill-name>/SKILL.md`. See
[plugins/skills/skills/README.md](plugins/skills/skills/README.md) for the
required `SKILL.md` frontmatter.

## Workflows

`plugins/skills/workflows/` holds dynamic workflow scripts. Each runs as
`/sebnow:<meta.name>`.

`/sebnow:improve-code` reviews existing code and fixes what it finds. Give
it paths or a revision range as a starting point; reviewers look as far
beyond it as their lens needs. Each round, one reviewer per lens reports
findings, a triage agent decides which to act on, and a fixer commits each
fix on its own with the project's check passing. The loop runs until
nothing is left to act on or `rounds` is reached, then tidies the fix
commits into the history and re-checks every rewritten commit.

```text
/sebnow:improve-code {"paths": ["internal/billing"], "rounds": 2}
```

Arguments: `paths` or `revs` (required, one of them), `rounds` (default
3), `check` (build and test command, otherwise detected), and `model`,
`models`, `efforts` to set the model or effort for every stage or per
stage (`scan`, `review`, `triage`, `fix`, `tidy`). The lenses are the
`reviewer-*` agents in `agents/`.
