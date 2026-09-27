# Research: current (2024-2026) guidance and critique on writing ADRs

Compiled 2026-09-27 by a research worker for the adr-writing skill rewrite.
Every URL cited was retrieved with WebFetch. Quotes marked **[s]** come only
from the fetch tool's summarising model and may be paraphrased; all others
were checked word-for-word against the page text.

## Findings

### 1. What an ADR is for and what belongs in it

Core consensus (unchanged since 2011): an ADR records one significant
decision, the forces or context behind it, and its consequences, good and bad.

- Nygard (2011): "Each record describes a set of forces and a single decision
  in response to those forces." Context is "value-neutral. It is simply
  describing facts." Consequences: "All consequences should be listed here,
  not just the 'positive' ones."
  https://www.cognitect.com/blog/2011/11/15/documenting-architecture-decisions
- Fowler (2026-03-24): "a short document that captures and explains a single
  decision relevant to a product or ecosystem." Two purposes: a record for
  people "months or years later", and "perhaps even more valuable, the act of
  writing them helps to clarify thinking".
  https://martinfowler.com/bliki/ArchitectureDecisionRecord.html
- Microsoft Well-Architected (2026-04-10): "Your architecture is the
  accumulation of its decisions." Capture "alternatives that you ruled out".
  "Always include context and rationale. A record without justification loses
  its value over time."
  https://learn.microsoft.com/en-us/azure/well-architected/architect-role/architecture-decision-record
- AWS Prescriptive Guidance: "it focuses on the reason for the decision rather
  than how the team implemented it."
  https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html
- MADR 4.0.0 (2024-09-17, current) required sections: Context and Problem
  Statement, Considered Options, Decision Outcome. Optional: Decision Drivers,
  Consequences, Confirmation, Pros and Cons of the Options, More Information.
  Front matter: status, date, decision-makers, consulted, informed.
  https://raw.githubusercontent.com/adr/madr/main/template/adr-template.md
  https://adr.github.io/madr/

### 2. Which decisions deserve an ADR

- Nygard: "'architecturally significant' decisions: those that affect the
  structure, non-functional characteristics, dependencies, interfaces, or
  construction techniques." AWS repeats this, citing Richards and Ford 2020.
- Reversibility and cost:
  - Microsoft: "Only include choices that affect the system's structure, key
    quality attributes, or are difficult to reverse."
  - Pureur and Bittner, quoting Booch: "significance is measured by cost of
    change". "Architectural decisions are costly to reverse". A database
    choice hidden behind an abstraction is not architectural **[s]**.
  - Zimmermann 2023: "ADs that are costly to undo simply cannot wait until
    sprint n." "The level of this significance and the severity of the
    consequences of the chosen solution (in terms of cost and risk) determine
    the importance and urgency."
  - Konishi: "ADRs are for decisions that are hard to reverse, that span
    multiple components, or that materially affect operability or security. A
    formatter choice is not an ADR. A datastore choice is."
  - Joel Parker Henderson's Claude Code skill: "affects structure, external
    interfaces, quality attributes, or is expensive/risky to reverse" **[s]**.
    https://github.com/joelparkerhenderson/architecture-decision-record/blob/main/skills/architecture-decision-record-skill/SKILL.md
- Blast radius: GOV.UK framework (2025-11) sorts decisions by level: Team
  ("No impact on other teams or shared platforms"), Programme ("Affects
  multiple teams or shared services"), Departmental ("sets a precedent for
  future work"), Cross-government **[s]**.
  https://www.gov.uk/government/publications/architectural-decision-record-framework/architectural-decision-record-framework
- Zimmermann's Architectural Significance Test (2020): high business value or
  risk; key stakeholder concern; runtime quality needs that "deviate from
  those already satisfied"; external dependencies; "cross-cutting nature";
  "First-of-a-Kind"; "troublesome... on a previous project" **[s]**.
  https://ozimmer.ch/practices/2020/09/24/ASRTestECSADecisions.html
- When to skip (Henderson README): "skip an ADR when a decision is limited in
  scope and time and risk and cost, or is already covered elsewhere".
  https://github.com/joelparkerhenderson/architecture-decision-record
- Dissenting, permissive view: Spotify (2020): "Almost always!" **[s]**.
  https://engineering.atspotify.com/2020/04/when-should-i-write-an-architecture-decision-record
  Google Cloud (reviewed 2024-08-16): write one "When there are two or more
  engineering options and you want to document your thoughts and reasons"
  **[s]**. https://docs.cloud.google.com/architecture/architecture-decision-records
- "Surprise to a future reader": no primary source fetched uses this wording.
  Closest is Nygard's "Nobody is left scratching their heads... 'What were
  they thinking?'" **UNVERIFIED:** that the phrasing is absent from sources
  not fetched.
- Trivial ADRs as a failure mode: Zimmermann 2026 found ADRs "that do not
  even cover architecturally significant design concerns; they... [are]
  descriptions of single implementation-level designs." Konishi names "ADRs
  written for trivial decisions, skipping the load-bearing ones."

### 3. Implementation detail (schemas, DDL, config, code)

Mainstream sources exclude it, with a narrow exception for illustrating
options.

- Zimmermann 2026: "these sections (and the entire ADR) are not the right
  place for implementation details (as-is or to-be); detailed design
  specifications deserve their own place in the documentation."
- Zimmermann 2023, "Mega-ADR" anti-pattern: "component responsibilities and
  collaborations are specified, or multiple diagrams or code snippets appear
  inside the ADR (one per option is ok usually)... Remedy: Move the detail
  design to a separate document." "Blueprint or Policy in Disguise" warns
  against "the amount of details provided and/or a rather commanding,
  authoritative voice."
- Microsoft: "Avoid making decision records design guides. If more
  justification or design ideation is available, provide a link... but the
  decision must be clear and stand alone without that material."
- AWS blog (2025-03): "Separate design from decision – Use a separate design
  document mechanism to explore alternative options thoroughly. Reference
  these design documents within the ADR."
  https://aws.amazon.com/blogs/architecture/master-architecture-decision-records-adrs-best-practices-for-effective-decision-making/
- Internal inconsistency at AWS: the Prescriptive Guidance also says the
  decision log "provides the project context as well as detailed
  implementation and design information."
- MADR lets each option carry "{example | description | pointer to more
  information | …}", which fits the one-snippet-per-option allowance.
- Contrarian, agent-focused vendor view: Actual AI (Kennedy, 2026-06-23):
  "Imperative, not prose. 'MUST use the framework's image component. MUST NOT
  use a raw img tag.'" Each decision should "ship with a check, a grep, a
  lint rule, a command". https://actual.ai/blog/agent-optimized-adrs
  This conflicts with Zimmermann's "Blueprint or Policy in Disguise".

### 4. Immutability and supersession

Consensus: accepted ADRs are append-only; change happens through a new,
linked ADR.

- Nygard: "If a decision is reversed, we will keep the old one around, but
  mark it as superseded."
- Fowler 2026: "Once an ADR is accepted, it should never be reopened or
  changed - instead it should be superseded. That way we have a clear log of
  decisions and how long they governed the work."
- Microsoft 2026: "The ADR serves as an append-only log. Don't go back and
  edit accepted records. If a decision changes, write a new record that
  supersedes the original and link the two together."
- AWS: "When the team accepts an ADR, it becomes immutable." Also after
  rejection: "treat ADRs as immutable documents after the team accepts or
  rejects them." For rejected ADRs the owner "adds a reason for the
  rejection to prevent future discussions on the same topic."

Dissent and nuance:

- Henderson README: "In theory, immutability is ideal. In practice,
  mutability has worked better for our teams. We insert the new info the
  existing ADR, with a date stamp, and a note that the info arrived after the
  decision." His "good ADR" list: "Don't alter existing information... amend
  the ADR by adding new information, or supersede." Dated appends, not
  rewrites.
- MADR's `date:` field means "when the decision was last updated".
- Zimmermann mentions marking ADRs superseded "(or even delete them from
  your log)".
- Konishi: silent edits of conclusions "breaks the audit trail".
- Structural critique (Drosopoulou, Java Code Geeks, 2026-05): ADRs are
  "static, point-in-time documents asked to perform a living-artefact
  function". "ADR maintenance cost is paid by the person with the least
  context"; proposes fitness functions and "Revisit if" clauses **[s]**.
  https://www.javacodegeeks.com/2026/05/the-reason-most-architecture-decision-records-get-written-and-never-read-is-architectural-not-cultural.html

### 5. LLMs and agents writing ADRs (2024-2026)

Empirical studies:

- Dhar, Vaidhyanathan, Varma (ICSA 2024): GPT-4 zero-shot produces "relevant
  and accurate Design Decisions, although they fall short of human-level
  performance" **[s]**. https://arxiv.org/abs/2403.01709
- Zhou, Li, Liang et al. (2025): precision of LLM-generated design rationale
  0.267-0.278, recall 0.627-0.715. Of arguments human experts did not raise,
  "1.59% to 3.24%... are potentially misleading" **[s]**.
  https://arxiv.org/abs/2504.20781
- Gupta, Dhar, Feitosa, Vaidhyanathan (arXiv 2604.03826, 2026): a
  zero-context Gemini-2.5-Pro baseline averaged 1123.41 tokens against 526.74
  for human ADRs: "without historical grounding, LLMs tend to generate
  extensive, verbose text". Supplying the last 3-5 prior ADRs as context
  worked best. "Context strategy matters more than model size" **[s]**.
  Caveat: the baseline prompt was the title only.
- da Silva and Gama, GADR (2026-08), ADRs from meeting transcripts:
  retrieval-based enrichment improves depth while "risking
  transcript-unfaithful content" **[s]**. https://arxiv.org/abs/2608.17694

Practitioner reports on failure modes:

- Equal Experts (2025-10): "The AI tools frequently hallucinated reference
  material, including non-existent APIs, web pages, or entire product
  features." Also "Mismatched justifications". Guardrails: "References MUST
  exist – check each reference to ensure that the link is valid" and "DO NOT
  make up product features that don't exist". A second LLM as "judge".
  https://www.equalexperts.com/blog/our-thinking/accelerating-architectural-decision-records-adrs-with-generative-ai/
- Laskowski and Michalak (2025-03): LLMs "struggle with independently
  capturing accurate context". "AI sometimes introduces incorrect or
  exaggerated details, especially when AI generates pros and cons". "90% of
  the time AI will suggest some 'hybrid' or 'combination' of the previous
  options that needs to be ignored".
  https://handsonarchitects.com/blog/2025/using-generative-ai-as-architect-buddy-for-adrs/

Worker's inferences for the skill rewrite (all **UNVERIFIED**, the worker's
own reading, not stated by the sources):

- Documented agent failure modes line up with Zimmermann's human
  anti-patterns: invented and hybrid options match "Dummy Alternative";
  verbosity matches "Mega-ADR" and "Novel"; exaggerated pros match "Sales
  Pitch"; hallucinated references match "Magic Tricks" and "Pseudo-accuracy".
- The "imperative rules for agents" approach conflicts with human-oriented
  guidance. Keeping the ADR as narrative and putting rules in separate agent
  files may resolve that.

## Sources

Primary and practitioner guidance:

- https://www.cognitect.com/blog/2011/11/15/documenting-architecture-decisions Nygard, original ADR post (2011-11-15)
- https://adr.github.io/madr/ MADR docs and news, 4.0 re-scoping (current 4.0.0)
- https://raw.githubusercontent.com/adr/madr/main/template/adr-template.md MADR full template, main branch
- https://www.ozimmer.ch/practices/2023/04/03/ADRCreation.html Zimmermann, ADR creation practices and anti-patterns (updated 2026-09-03)
- https://ozimmer.ch/practices/2026/09/12/ADRMistakes.html Zimmermann, ten common ADR mistakes (2026-09-12)
- https://ozimmer.ch/practices/2020/09/24/ASRTestECSADecisions.html Zimmermann, Architectural Significance Test (2020-09-24)
- https://martinfowler.com/bliki/ArchitectureDecisionRecord.html Fowler, ADR bliki (2026-03-24)
- https://learn.microsoft.com/en-us/azure/well-architected/architect-role/architecture-decision-record Microsoft Well-Architected ADR page (2026-04-10)
- https://docs.aws.amazon.com/prescriptive-guidance/latest/architectural-decision-records/adr-process.html AWS Prescriptive Guidance, ADR process (undated)
- https://aws.amazon.com/blogs/architecture/master-architecture-decision-records-adrs-best-practices-for-effective-decision-making/ AWS Architecture Blog (2025-03-20)
- https://docs.cloud.google.com/architecture/architecture-decision-records Google Cloud ADR overview (reviewed 2024-08-16)
- https://engineering.atspotify.com/2020/04/when-should-i-write-an-architecture-decision-record Spotify (2020-04-14)
- https://www.gov.uk/government/publications/architectural-decision-record-framework/architectural-decision-record-framework GOV.UK framework text (2025-11-04)
- https://github.com/joelparkerhenderson/architecture-decision-record Henderson ADR repo README (2026)
- https://github.com/joelparkerhenderson/architecture-decision-record/blob/main/skills/architecture-decision-record-skill/SKILL.md Henderson Claude Code ADR skill (2026)

Critiques:

- https://www.infoq.com/articles/architectural-decision-record-purpose/ Pureur and Bittner, "Has Your ADR Lost Its Purpose?" (2023-10-25)
- https://hidekazu-konishi.com/entry/architecture_decision_records_templates_and_operations.html Konishi, templates and operations (updated 2026-09-09)
- https://www.javacodegeeks.com/2026/05/the-reason-most-architecture-decision-records-get-written-and-never-read-is-architectural-not-cultural.html Drosopoulou, structural critique (2026-05-22)

LLM and agent sources:

- https://arxiv.org/abs/2403.01709 Dhar et al. (ICSA 2024)
- https://arxiv.org/abs/2504.20781 Zhou et al. (2025-04, revised 2025-12)
- https://arxiv.org/html/2604.03826v1 Gupta et al., context strategies (2026-04)
- https://arxiv.org/abs/2608.17694 da Silva and Gama, GADR (2026-08-18)
- https://www.equalexperts.com/blog/our-thinking/accelerating-architectural-decision-records-adrs-with-generative-ai/ Equal Experts (2025-10-27)
- https://handsonarchitects.com/blog/2025/using-generative-ai-as-architect-buddy-for-adrs/ Laskowski and Michalak (2025-03-31)
- https://actual.ai/blog/agent-optimized-adrs Kennedy (vendor) (2026-06-23)

Not cited: the AWS guide welcome page (HTTP 404).

## Open questions

1. "Surprise to a future reader": no primary source found; may live in
   Keeling's IEEE Software columns or Harmel-Law's 2024 book, not fetched.
2. Narrative versus imperative: no neutral study compares how agents follow
   each. **UNVERIFIED:** whether one exists. Actual AI wants MUST rules;
   Zimmermann warns against a "commanding, authoritative voice". The skill
   must choose.
3. Quotes marked [s] were not checked word-for-word.
4. AWS Prescriptive Guidance pages show no date.
