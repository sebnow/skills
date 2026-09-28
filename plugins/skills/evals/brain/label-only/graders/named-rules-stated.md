---
type: llm
focus: last_message
---
The senior's report behind this reply named two groups of existing rules by
id with a few words at most, and the user has read neither group. The
generic rules-stated grader is noisy here because the judge has to decide
alone what counts as stating a rule; this one gives it the rules to look
for.

The five overhaul rules say: (1) the return states each engine's hours since
its last overhaul against its limit; (2) hours run in the seasons before the
logs count towards overhaul at a rate the operator elected; (3) an engine
replaced during the season starts again from zero; (4) the limit follows
the engine's class; (5) an engine past its limit is an open finding.

The three survey rules say: (1) a boat carrying paying passengers is
surveyed every season; (2) hull work at a yard off the approved list counts
towards the survey interval only in part; (3) a season laid up throughout
does not count towards it.

PASS: for each group the reply asks the user a question or decision about,
the reply says what at least all but one of its rules say, in any words (a
clause each is enough). A group the reply asks nothing about needs nothing.

FAIL: the reply asks about a group and states fewer of its rules than
that, for example calling them "the five overhaul rules", "the hours-based
rules" or "how engine hours count", or naming the survey rules only by
subject ("passengers, hull work and laid-up seasons") without what each
does.
