# sebnow skills marketplace

A collection of skills and a brain agent for Claude Code and Codex.
The Claude plugin in `plugins/skills/` is the source. The Codex package in
`plugins/codex/` is generated and committed so installation needs no build tools.

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

## Codex

Add this repository as a Codex marketplace and install the plugin:

```sh
codex plugin marketplace add <path-or-git-url>
codex plugin add sebnow@sebnow
```

Start a new Codex session and run the plugin's setup flow, or ask:

```text
Use the setup skill from the sebnow plugin to install the brain profile.
```

Setup needs Python 3.11 or later. It installs `brain.config.toml` and the brain
skill and three worker definitions under `$CODEX_HOME/sebnow/` (default
`~/.codex/sebnow/`). It records hashes so subsequent setup runs can update its
files while refusing to overwrite local edits. Run setup after plugin updates.

Launch the brain from your vault or project:

```sh
codex --profile brain -C /path/to/vault --add-dir /path/to/project
```

The brain uses GPT-6.1 Sol; senior uses GPT-6.1 Sol, and junior and scout use
GPT-6 Luna. These models must be available to your account. Change the `MODELS`
mapping in `scripts/build-codex.py` to choose different build defaults.
You can override the main model with `codex --profile brain --model <model>`.

Workers inherit the session's permissions. Scout's read-only rules are prompt
instructions in Codex; Claude's tool denylist does not carry over. The existing
eval fixtures and eval-runner reference target Claude Code.

## Build the Codex package

After editing the source skills or agents, run:

```sh
python3 scripts/build-codex.py
```

Commit `plugins/codex/` with the source changes. The build uses only Python's
standard library and leaves the Claude package unchanged. Its agent reader
expects the existing frontmatter format: a single JSON-quoted description line.
Codex setup lives only in the generated Codex package.

The Codex manifest declares onboarding, but plugin installation alone does not
register the brain profile. The setup flow performs that separate step.
