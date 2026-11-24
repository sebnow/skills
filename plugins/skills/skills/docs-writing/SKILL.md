---
name: docs-writing
description: "Writes and updates the documentation of the repository it runs in: Diátaxis pages for the people who use the system and the engineers who change it, built around the high-level architecture and diagrams, verified against the code, and written to strict plain-English style rules. Use when asked to document this repo, project, service, or system, write docs or a docs site, add an overview, architecture page, how-to guide, or glossary, explain how a part of the system works in the docs, or bring existing docs in line with the code. Not for READMEs of a single module, ADRs (use adr-writing), API docstrings, or answering a question in chat (use architecture-visualization)."
---

# Repository documentation

Document the system that this repository holds, for the people who use it
and for the engineers who change it.

## Inputs

Take these from the user's request and the repository. Ask only for what
neither answers.

- **Location.** Follow the repository's convention: the `docs_dir` of a site
  generator such as MkDocs or Docusaurus, or the directory that already holds
  the docs. Use `docs/` when there is no convention.
- **Audiences.** By default, write for two audiences:
  1. People who use or integrate with the system. They need to know what it
     does for them and how to do each task.
  2. Engineers who change it. They need the parts, how the parts fit, and
     why.

  Add an audience when the user or the existing docs name one, such as
  support staff or product managers. Write for readers without a technical
  background in plain English, without jargon.
- **Style guide.** Look for one in the docs, `CONTRIBUTING.md`, or the agent
  instructions. Where it conflicts with the style rules below, it wins.
- **Glossary.** Use an existing glossary's terms.

## What belongs in the docs

Document what a reader cannot easily get from the code: the parts of the
system, the responsibility of each, the boundaries between them, the flows
and lifecycles, and the reasons behind them. For users, document what the
system does for them and the steps of each task.

Leave out what the code already says to anyone who opens it:

- source file paths, line numbers, and a tour of the directory layout,
- class, function, and table names,
- configuration values and their defaults, constants such as limits and
  retry delays, and schemas.

A user-facing page names what its reader types or sends: a command, an
endpoint, a request field. Name a part of the system as a reader would, such
as "the worker" or "the job queue", not by the class that implements it.

API specifications, such as OpenAPI files, stay where they are. Link to them,
and do not copy them into a page. Leave out details that would help someone
evade an abuse or security check, such as exact detection patterns.

## Structure

Write one Diátaxis type per file, and link to the other types instead of
mixing them:

- **Tutorial**: a first success for a new user, step by step.
- **How-to guide**: the steps of one task, for a reader who knows the goal.
- **Reference**: facts to look up, in a fixed shape per entry.
- **Explanation**: how and why the system works as it does.

A new documentation set gets these pages, named to fit the location's
convention:

- Entry page (`README.md`, or `index.md` for a site generator): what the
  system is in two sentences, then a starting page for each audience.
- `overview.md` (explanation): what the system does and for whom, in plain
  English, with one diagram at most.
- `architecture.md` (explanation): the parts, the responsibility of each, and
  what connects them, with a diagram of the parts.
- `glossary.md` (reference): one short entry per term, in plain English, with
  a stable anchor. List the synonyms to avoid under each term.
- `feature-map.md` (reference): one entry per user-facing feature, with the
  page that explains it and the parts it passes through, in order. List only
  features that a page explains.
- An explanation page for each lifecycle and each main flow, with diagrams.
- A how-to guide for each task a user performs, with an example request and
  response where the task has them, and a link to the specification.
- A tutorial, when a new user needs one to reach a first success.

Leave out a page that would hold only a sentence or two, and say so in the
report. A user can ask for any other page, such as a list of supported
features; give it one Diátaxis type, and add it to the entry page.

## Updating existing pages

- Read every page in the location first, including the entry page, the
  feature map, and the glossary.
- Keep the existing structure, file names, and terms. Propose a structural
  change to the user before making it.
- When you add, rename, or remove a page, update the entry page, the feature
  map, the site generator's navigation, and every link to the page.
- Verify each statement you add or change against the code. Where you touch a
  page, check its other statements against the code too, and report each one
  that is wrong.

## Sources and verification

- The code is the source of truth. Verify every statement against code or
  configuration before writing it. Check existing prose docs against the code
  too, because they can be wrong.
- Where code and a stated requirement disagree, describe what the code does,
  and list the requirement as planned.
- Leave out what you cannot verify, and list each such point in the report.
- Do not state whether a feature or tool is still in use unless the user says
  so.
- Report each defect or contradiction you find in the system. Do not fix code
  as part of documenting it.

## Diagrams

Use the architecture-visualization skill to choose the view and the parts of
each diagram, and put each diagram in the page as a fenced `mermaid` block.

- Draw a state diagram (`stateDiagram-v2`) for every lifecycle or state
  machine. Show each state, and label each transition with the event or
  action that causes it.
- Use a flowchart for the parts and for a process.
- Use a sequence diagram for an exchange between parts or with an outside
  system.

A diagram on a user-facing page shows only what its reader observes from
outside the system. Parse every diagram with the mermaid CLI when it is
installed.

## Style rules

These bind every sentence, except where the repository's style guide says
otherwise.

General clarity:

- Put one instruction or idea in each sentence, and one topic in each
  paragraph.
- Use the active voice, with the actor as the subject. Use the passive only
  when the actor is unknown.
- Use the present tense for behaviour and facts. Do not use "will" for
  behaviour.
- Keep articles, subjects, and verbs. Do not write in telegraphic style.
- Use the glossary's terms, never the synonyms it lists to avoid. An API
  identifier that uses an avoided word appears only in code font, with a list
  that maps it to the glossary term.
- Omit needless words. Prefer specific, concrete words.
- Do not use idioms.
- Write headings in sentence case, without end punctuation. Glossary terms
  keep their capitals.
- Use a list for three or more parallel items, and keep the items parallel in
  form.
- Use a serial comma in an inline series.
- Write link text that names its target.
- Use code font for identifiers, file names, and paths.
- Prefer prose and lists. Use a table only when every line fits in 80 columns
  and most rows are short, single-line entries.

Hard limits:

- Keep an instruction to 20 words and a description to 25. These are
  ceilings, not targets.
- Do not use a noun cluster longer than three words.

House rules:

- Do not use these words: please, simply, just, easily, easy, e.g., i.e., via,
  utilize, and/or.
- Do not use em-dashes.
- State what happens, in positive form. Use MUST NOT only where the
  prohibition is the rule.
- Write an absence with "not", never with "no" in front of the object:
  - "The tests do not cover any edge cases", not "The tests cover no edge
    cases".
  - "The two services do not share an identifier", not "The two services
    share no identifier".
- Do not write "cannot" for a rule that the system enforces. Say what the
  system does, such as "the API rejects the request".

## Checks

Before a page is finished, run these on the documentation directory, and show
the user each command and its output:

- Banned words and em-dashes:
  `grep -rnwiE 'please|simply|just|easily|easy|via|utilize|e\.g\.|i\.e\.|and/or|—' <directory>/`
- Absences written with "no": `grep -rnwi 'no' <directory>/`. Review each hit,
  and rewrite each one that puts "no" in front of an object.
- Sentences over 25 words.
- Code detail that belongs in the code: source paths, class and function
  names, and configuration values. Grep the pages for the names of the
  repository's source directories and main types.
- Broken links and anchors, including links between pages and to glossary
  anchors.
- Mermaid diagrams, parsed with the mermaid CLI when it is installed.

Without a shell, make each check by reading the pages, and say so in the
report.

## Process

1. Analyse the repository. Map the parts, the user-facing features and tasks,
   the lifecycles, what is built and what is only specified, and the existing
   docs.
2. Write the overview, architecture page, and glossary first. Then write the
   explanation pages, how-to guides, and tutorial, and last the feature map
   and entry page.
3. Run the checks on each page.
4. Report the pages you wrote or changed, the points you could not verify,
   and the defects you found.

## Questions

Ask the user only what the repository cannot answer:

- a term that the code does not define,
- a process that happens outside the code,
- a conflict between two instructions.
