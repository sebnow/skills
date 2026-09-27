# Evidence behind the skill-writing rules

Each entry gives the finding, the source, the models it was measured on, and
what it means for a skill file. Findings dated 2025 or 2026 unless noted. Entries marked UNVERIFIED have not been checked against the primary
source.

## Instruction files can cost without helping

- Repository context files (AGENTS.md, CLAUDE.md) did not raise task success
  on SWE-bench or CTXbench for Claude Sonnet 4.5, GPT-5.2, GPT-5.1 Mini, or
  Qwen3-30B; LLM-generated files moved success by about -0.5 to -2 points,
  developer-written files by +2.4 points (not significant), and both raised
  cost by about 20%. Instructions were followed literally: a tool named in the
  file was used 1.6 times per task versus almost never without it.
  Gloaguen et al., "Evaluating AGENTS.md", arXiv:2602.11988.
  Meaning: compliance is not the outcome to grade. Grade the task result, with
  a no-skill arm.
- A two-agent ablation on real repositories bounded the correctness effect of
  context files to at most 10 to 15 points by equivalence testing; the files
  never converted a near-miss into a pass.
  Khatri, arXiv:2607.27250.
  Meaning: state a null result as a bound, not as "no difference".
- Context files cut median runtime by 29% and output tokens by 17% with
  comparable completion. Lulla et al., arXiv:2601.20404.
  Meaning: cost is a legitimate secondary outcome for a skill.
- Adding constraints the model already satisfies lowered task accuracy on
  math, multi-hop QA, and code, including on Claude Sonnet 4.5.
  Qi et al., arXiv:2601.22047.
  Meaning: a rule the baseline already passes can lower accuracy. Delete it.
- Anthropic removed over 80% of Claude Code's system prompt for Claude 5
  generation models with no measurable loss on its coding evaluations, and
  names over-constraint through CLAUDE.md and skills as a cause of worse
  output. Shihipar, "The New Rules of Context Engineering for Claude
  5-generation models", 2026-07-24,
  https://claude.dev/blog/the-new-rules-of-context-engineering-for-claude-5-generation-models/
- Anthropic's Fable 5 prompting guide: "Skills developed for prior models are
  often too prescriptive for Claude Fable 5 and can degrade output quality.
  Review and consider removing older instructions if default performance is
  better." The Opus 5 guide says to remove carried-over verification
  instructions rather than rewrite them.
  https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/prompting-claude-fable-5
  https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/prompting-claude-opus-5
  Meaning: every model release is a trigger to re-run the no-skill arm.

## How many rules a file can carry

- Compliance falls with instruction count long before the context limit.
  Best models held about 68% at 500 instructions; joint satisfaction of eight
  simultaneous constraints was 5.7% even when each alone passed about 41%;
  perfect-response rate fell from 85% at 10 instructions to 15 to 31% at 40 on
  Claude Sonnet 5, Claude Haiku, Gemini Flash and two Qwen models.
  Jaroslawicz et al. (IFScale) arXiv:2507.11538; Vasileva arXiv:2608.12426;
  Eliav arXiv:2607.19257.
  Meaning: a skill is a budget of a few dozen rules. Conflicting rules cost
  more than the count alone.
- Prohibitions decay under context load while requirements persist: "never"
  rules fell from 73% compliance at turn 5 to 33% at turn 16 across 12 models.
  Gamage, arXiv:2604.20911.
  Meaning: state the action to take, not the one to avoid.
- Performance degrades with input length on every frontier model tested, and
  single distractors reduce accuracy. Chroma, "Context Rot", 2025-07-14,
  https://www.trychroma.com/research/context-rot; NoLiMa, arXiv:2502.05167.
  Meaning: each skill token competes with the task.

## Wording

- Uppercase words attract attention without improving accuracy and can
  degrade it. Dillitzer et al., "Attention is Case-Sensitive",
  arXiv:2608.03711. Anthropic: Opus 4.5 and later are more responsive to the
  system prompt, so "CRITICAL: You MUST use this tool when..." over-triggers
  on those models; write "Use this tool when...".
  https://platform.claude.com/docs/en/build-with-claude/prompt-engineering/claude-prompting-best-practices
- Emotional stimuli, tips and threats did not change accuracy overall, but
  single questions swung by up to 36 points either way. Meincke et al.,
  arXiv:2508.00614, 2503.04818.
  Meaning: measure a wording change over several runs before keeping it.
- Personas do not improve factual accuracy; irrelevant persona detail cost up
  to 30 points. Zheng et al. arXiv:2311.10054; Luz de Araujo et al.
  arXiv:2508.19764.
- Anthropic: explain the reason behind an instruction rather than command it;
  the model generalises from the reason.
  https://github.com/anthropics/skills/blob/main/skills/skill-creator/SKILL.md
  and the prompting best practices page above.

## Examples and reasoning scripts

- On strong models few-shot exemplars mostly align output format and are
  otherwise ignored; excess examples degrade results. Cheng et al.
  arXiv:2506.14641; Tang et al. arXiv:2509.13196.
  Meaning: a few canonical examples to fix format, none to teach reasoning.
- Explicit chain-of-thought prompting gained 12 to 14 points on 2024
  non-reasoning models and about +3 or -3 on reasoning models at 20 to 80%
  more latency. Meincke et al., arXiv:2506.07142. It cut accuracy by up to
  36 points on tasks where deliberation hurts. Liu et al., arXiv:2410.21333.
- Self-consistency voting gained under 2 points on Gemini 2.5. Loo,
  arXiv:2511.00751. Self-correction prompts are unreliable; an external check
  (test, validator, fresh-context reviewer) is what works. Tsui
  arXiv:2507.02778; Chen et al. arXiv:2606.05976.
  Meaning: do not script thinking, reflection or voting in a skill. Give it
  something external to check against.

## Where a skill sits in the instruction order

- Across 12 models on coding tasks, rules in a skill description lost
  conflicts to the system prompt, project file and user instruction (mean
  conflict rank 4.56 versus 2.22), and every model did worse on rules that
  oppose its default behaviour. Huang et al., Harness-IF, arXiv:2608.11727.
  Meaning: a skill usually loses a conflict with the system prompt, project
  file, or user instruction. A rule that must always hold belongs in a hook; a
  rule that should win conflicts belongs in CLAUDE.md.
- Tool and skill selection is driven by description wording; ambiguity is the
  main defect. Rewritten descriptions cut the accuracy loss from tool-set
  growth by 29% and raised query success by 61% (Guo et al.
  arXiv:2602.20426); across 10,831 MCP servers, the functionality and
  accuracy parts of a description raised selection by 11.6 and 8.8 points
  (Wang et al. arXiv:2602.18914). Anthropic's skill-creator notes that
  Claude under-triggers skills, while the prompting best practices page says
  pushy descriptions over-trigger on Opus 4.5 and later.
  Meaning: purpose first, concrete triggers, third person; calibrate
  pushiness per model with a trigger eval. UNVERIFIED: no study isolates
  "do not use for" clauses; treat them as weak.

## Judging and sample size

- LLM judges are reproducible but biased: position bias above 0.10 in some
  deployed judges; style bias favouring markdown; verbosity bias varies by
  judge (Claude judges leaned concise). Norman et al. arXiv:2606.19544;
  Soumik arXiv:2604.23178.
- Pairwise preferences flipped in about 35% of cases when a stylistic
  distractor was added, absolute rubric scores in 9%. Tripathi et al.,
  arXiv:2504.14716.
  Meaning: grade with an absolute rubric of concrete PASS and FAIL
  conditions; use pairwise only with order swapping.
- Detecting a 3-point difference at 80% power needs about 970 items;
  resampling the same item helps only until within-item variance is
  negligible. Miller, "Adding Error Bars to Evals", arXiv:2411.00640. Wilson
  intervals are recommended for small binary samples:
  https://inspect.aisi.org.uk/metrics.html
  Meaning: three runs catch gross failures (0/3 to 3/3). A claimed shift of a
  few points needs many more paired runs than a skill eval usually has.
- Anthropic on agent evals: report pass^k where consistency matters; start
  with 20 to 50 tasks drawn from real failures; grade what the agent
  produced, not the path; when the no-skill arm passes every solvable task
  the eval is saturated. Anthropic, "Demystifying evals for AI agents",
  2026-01-09,
  https://www.anthropic.com/engineering/demystifying-evals-for-ai-agents
