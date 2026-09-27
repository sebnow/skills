---
name: skill-writing
description: "Creates, changes, and retires SKILL.md skills for agent plugins, with an eval that measures each change. Use when asked to create or write a skill, turn notes or a workflow into a skill, add a rule to a skill, or read eval results for a skill and decide what to keep. Covers eval cases with a no-skill arm, what belongs in a skill body, and when a skill or section is obsolete. Not for README or general documentation."
---

# Skill Writing

## Start from a failure, not a topic

A request for "a skill for X" names a topic. Before writing skill text, record
the failure behind it:

1. Name the failure concretely: what the agent did, on what input, and what
   it should have done. If the request does not say, ask, or find it in a
   transcript or a repeated correction.
2. Write an eval case that reproduces it: a realistic prompt of the kind a
   user types, and a grader that checks the outcome the agent got
   wrong. See [eval-runner.md](references/eval-runner.md) for the case format
   and the example cases listed there.
3. Run the case without the skill. If it passes, the model already handles
   this and the skill needs no text for it.

Only then draft the section that addresses the failure, and treat the draft
as provisional until a run with the skill shows the case passing. If you
cannot run the eval in this session, deliver the case and the run command
first and label any draft as unmeasured; a skill written before its failure
is observed usually restates what the model already does.
