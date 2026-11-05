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
        ├── agents/               # subagent definitions
        ├── evals/                # claude plugin eval cases, per skill or agent
        └── skills/
            └── <skill-name>/
                └── SKILL.md
```

## Use

Add the marketplace, then install the plugin:

```sh
claude plugin marketplace add <path-or-git-url>
claude plugin install sebnow@sebnow
```

Run the brain agent as the main session:

```sh
claude --agent sebnow:brain
```

The agent asks for Fable, which needs a Max plan or usage credits. On
other plans, override the model at launch:

```sh
claude --agent sebnow:brain --model opus --effort high
```

## Adding a skill

Create `plugins/skills/skills/<skill-name>/SKILL.md`. See
[plugins/skills/skills/README.md](plugins/skills/skills/README.md) for the
required `SKILL.md` frontmatter.
