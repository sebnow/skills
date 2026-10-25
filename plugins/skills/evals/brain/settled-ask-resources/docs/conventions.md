# Conventions

How decisions and requirements are recorded in this repository.

## Architecture decisions

1. One decision per ADR, so each can be superseded on its own.
2. An accepted ADR is never edited and never amended in part. To change
   anything it decides, write a new ADR that supersedes it whole: restate
   what still holds, say what changes, and set the old one's status to
   `superseded by ADR-NNNN`.
3. When a decision changes, every accepted ADR that states the old decision
   is superseded in the same change, not only the one the work started from.
4. An ADR states the decision and its reasons. The conversation or ticket it
   came out of is not a source and is not cited.

## Requirements

5. One file per rule: `docs/requirements/REQ-NNNN.md`, with a `Status:` of
   draft, accepted or withdrawn.
6. A requirement whose rule changes is never edited. Set its status to
   `withdrawn`, add a `Replaced by: REQ-NNNN` line, and write the new rule in
   a new file.
7. A requirement may cite an ADR. An ADR does not cite requirements.
