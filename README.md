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

Run setup after installing or updating the plugin. In the Codex CLI, start a
new session and enter `$sebnow:setup`. If your client shows a setup action for
the plugin, that action invokes the same skill.

Setup needs Python 3.11 or later. It writes `$CODEX_HOME/brain.config.toml` and
installs the brain skill and four worker definitions under `$CODEX_HOME/sebnow/`.
The default `CODEX_HOME` is `~/.codex`.

If setup reports an existing-file conflict, move the named file to a backup
location and rerun setup. Keep the backup so you can recover any custom
settings.

Launch the brain from your vault or project:

```sh
codex --profile brain -C /path/to/vault --add-dir /path/to/project
```

The build defaults are GPT-6.2 Astra for the brain, GPT-6.2 Sol for the senior,
and GPT-6 Luna for the mid, junior and scout. Each role defaults to medium
reasoning effort. To change the default models, edit the `MODELS` mapping in
[scripts/build-codex.py](scripts/build-codex.py). Your Codex account needs access
to the models you configure.
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
