# Architecture visualization evals: notes

## What the cases are for

The user wants a skill that breaks architecture into understandable
concepts and shows it visually, as a building block for designing and for
understanding existing systems. They named three failures:

- ASCII diagrams.
- Prose where a picture would do.
- Going too low: classes, functions, the directory tree.

HTML is the preferred medium, SVG the most customizable, Mermaid
acceptable, canvas possible. How the picture reaches the user is the
harness's concern, not the skill's.

Cases:

| case | occasion | asks for a picture? |
|---|---|---|
| onboarding-overview | existing code (inkwell) | no: "how does inkwell work?" |
| upload-walkthrough | existing code | no: "walk me through what happens when..." |
| parts-fit-together | existing code | "show me how the parts fit together" |
| diagram-asked | existing code | yes: "draw me a diagram" |
| design-breakdown | design in the prompt, no code | no: "break it down for me" |
| design-diagram-asked | design in the prompt | yes: "sketch the architecture" |

`resources/` is inkwell, a stub document conversion service of about
1,000 lines of Python. Its parts are an API, a Postgres job queue,
workers, an external renderer, object storage, a webhook outbox, a usage
meter and a janitor. The cases copy it in with `fixture.sh`.

Graders:

- rendered-graphic (regex on the trace): an `.html` or `.svg` file,
  inline `<svg>` or `<canvas>`, or a Mermaid block appears.
- no-ascii-diagram (regex on the trace): no box-drawing characters and no
  `+---` boxes.
- concept-level (opus judge on the trace): the units are processes,
  stores and external systems, not classes, functions or files.

## Baseline without a skill (2026-11-17)

### The eval runner did not run here

`claude plugin eval` could not run in the session that wrote these cases.
Two things stopped it:

- With Bash granted, every run failed to start with: "a credentials file
  in this environment ... could not be followed (an ssh command uses a
  shell expansion that cannot be resolved here), so the Bash sandbox
  cannot exclude the files it points at". The cause was `GIT_SSH_COMMAND`
  in the environment. The cases now allow Read, Grep, Glob, Edit and
  Write only.
- Without Bash, every run failed with `exit 1: Not logged in · Please run
  /login`, because the session had no credentials.

### Proxy runs

So the no-skill arm was run by hand instead. Each run was an orchestrator
child thread on claude-opus-5-5 with no architecture skill. The child got
the case prompt word for word, plus three things the eval runner would
not add:

- A setup line that fetched inkwell by commit and tree SHA (`350882d`,
  `6f510b7`), so no eval path or name was visible to it.
- "Reply to the user directly, not in a report format".
- "If you write any files, end your final reply with ... each file's path
  and its full contents".

The children could write files and run shell commands. None of them wrote
a file. Each child thread's transcript is kept under its thread id.

Treat these as a proxy for the runner. The orchestrator harness differs
from `claude -p`, and the last instruction may hint that files are
expected. Even with that hint, no run produced a file.

| case | run (thread) | medium | level of detail |
|---|---|---|---|
| onboarding-overview | lbrSE9hCIvwj2bztvumReG | prose, tables, numbered steps | concept-level parts, with SQL, header names, constants and file:line throughout |
| onboarding-overview | UvluMBmtnKS8EE5rPbHS6r | prose, tables, a directory tree in a code block | as above, and restates the package layout ("Code layout") |
| upload-walkthrough | BwQL33Ek4SIYGo17HgEIL7 | prose, numbered steps | function level: `create_document`, `RateLimiter.check`, `auth.py:54`, backoff formula, SQL |
| upload-walkthrough | 2jGHqiMbZul2Qq9thc75aI | prose, numbered steps | function level, same as above |
| parts-fit-together | 60umCvhtZ6AOWscKiBO6JQ | ASCII box diagram, then prose | boxes are processes and stores, labelled with classes (`Worker.run()`, `KeyStore`); a table of the classes each command builds |
| parts-fit-together | 8ZZ5d0vowJzFG3Nfg16gOD | ASCII box diagram, then a table with one row per module | the diagram is concept-level; the table restates the package layout file by file |
| diagram-asked | qTCAw5CGhe8GRuRC5whvt1 | ASCII box diagram, then prose | mostly concept-level, with class names inside boxes |
| diagram-asked | xJ5374AYE5kCCnhhz2m2hL | ASCII box diagram and an ASCII state diagram | low: boxes list file paths, routes, backoff formulas (`10s·4^n`), table names, endpoint paths |
| design-breakdown | ruyrM6OwnA4PqHvCDueCxb | ASCII flow diagram, then prose and a table | concept-level |
| design-breakdown | Q9vJdTnDsaff6bQGmtrTlu | ASCII box diagram, then prose and a table | concept-level |
| design-diagram-asked | ok4rHSy5VrF3mgoiYcjZey | ASCII box diagram, then prose | concept-level |
| design-diagram-asked | DthBVpg6FoWGtdFqb33BCn | ASCII box diagram, then prose | concept-level |

### How often each named failure appeared

- **No rendered graphic: 12 of 12.** No run wrote HTML, SVG, canvas or
  Mermaid, in a file or in the reply. One run explained why: "Mermaid
  won't render in a terminal, so this is plain ASCII" (DthBVpg6). The
  agent let its belief about the delivery channel pick the medium.
- **ASCII diagram: 8 of 12.** Every run that drew anything used
  box-drawing characters. That covers both runs of each case that asked
  for a picture ("show me", "draw", "sketch") and both runs of
  design-breakdown, where nothing was asked for.
- **Prose instead of a picture: 4 of 12.** These were both runs of
  onboarding-overview and both runs of upload-walkthrough. They are the
  existing-code cases that do not say "show", "draw" or "sketch". Both
  "big picture" runs answered with tables and numbered lists.
- **Too low-level: 6 of 8 existing-code runs; 0 of 4 design runs.**
  - Both upload-walkthrough runs told the flow as a chain of functions,
    with file:line references.
  - xJ5374AY put file paths, routes and retry formulas inside the boxes.
  - 8ZZ5d0v and UvluMBm restated the package layout.
  - lbrSE9h's parts were concept-level, but SQL, headers and constants
    made up much of the text.
  - 60umCvh and qTCAw5C stayed mostly at concept level, with class names
    as labels.
  - The design runs stayed at the level the user wrote in. They did not
    invent classes or schemas.

Predicted grader results on these replies: rendered-graphic fails 12 of
12. no-ascii-diagram fails 8 of 12; the 4 prose runs pass it and fail
rendered-graphic instead. concept-level was judged by hand above, not by
the opus judge.

### Other defects seen

These were seen in the replies, not graded. No browser was available,
so they come from reading the text as a terminal would show it.

- **Broken box edges.** Labels are wider than their boxes:
  - `║ Kafka: appointments║` and `│ (clinic staff)│` (ok4rHSy5)
  - `│ <tenant>/<id>/│` and `│ LibreOffice│` (60umCvh)
  - `│LibreOffice│  │` (8ZZ5d0v)
- **Diagrams wider than a terminal.** ok4rHSy5 runs to about 180 columns
  and qTCAw5C to about 160, so a narrower terminal wraps them.
- **Edges that cannot be followed.** In ruyrM6O the receipts arrow joins
  a vertical line that serves three other edges, so it cannot be told
  where it ends. In 60umCvh and 8ZZ5d0v several edges cross Postgres's
  column without junctions.
- **Answers to questions not asked.** All 8 existing-code runs ended with
  a list of suspected bugs (unopened pool, lease shorter than the job,
  duplicate webhooks). The list often ran as long as the explanation.
  All 4 design runs added a list of open design questions. That fits
  "check we've got it right", but it took up most of the reply.

## Watch

- Bash was left out so the runner can start where `GIT_SSH_COMMAND` is
  set. In a normal session the agent has Bash. If a skill tells it to
  render or check the picture, consider adding Bash back where the runner
  allows it.
- `focus: trace` gives the judge only the first and last 12 messages. A
  Write early in a long session could fall outside that window.
- Grader frontmatter must not contain `---` anywhere, even in a comment.
  A comment showing `+---+` cut the frontmatter short, and the regex
  failed with "nothing to repeat".
- The regex graders have not run against a real trace. Check that the
  session JSON contains none of the box characters, for example in the
  system prompt, before trusting no-ascii-diagram.
- There is no skill-fired grader yet, because the skill has no name.
  Add one to each case when it does.

## Running

From `plugins/skills`, with the plugin loaded (with-without arms):

    claude plugin eval . --tag architecture-visualization --scaffold \
      --allow-tools Write Edit --model claude-opus-5-5 \
      --judge-model claude-opus-5-5 --no-publish --trust-plugin -j 1 \
      --keep-temp --json /tmp/architecture-visualization.json

For a pure no-skill arm before the skill exists, copy this directory into
an empty plugin and run there with `--ablation none`:

    P=$(mktemp -d); mkdir -p $P/.claude-plugin $P/evals
    echo '{"name":"noskill","version":"0.0.0"}' > $P/.claude-plugin/plugin.json
    cp -R evals/architecture-visualization $P/evals/
    cd $P && claude plugin eval . --ablation none --scaffold \
      --allow-tools Write Edit --model claude-opus-5-5 \
      --judge-model claude-opus-5-5 --no-publish --trust-plugin -j 1 \
      --keep-temp --json /tmp/arch-noskill.json

`--case` keeps only the last value given, so run one case at a time or
use `--tag`. Use `-j 1`.
