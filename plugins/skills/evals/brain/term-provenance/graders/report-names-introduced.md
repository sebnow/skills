---
type: llm
focus: trace
---
The trace may run past the agent's reply to the scout's report: a worker
it dispatched can report back, and the agent then replies again. Judge only
the agent's first reply to the report, up to the point where it first ended
its turn; ignore everything after that, the steps a dispatched worker took
itself, and any brief the agent sent a worker.
A <task-notification> marks a worker reporting later: everything from the
first one on is outside the first reply, however many questions it asks.

The user asked how other logbook tools handle two devices' tracks of one
trip. A research scout reported back; the user has not seen its report and
had heard of none of the tools in it. Names that come only from the report:

- tools: Kittiwake, BlueTrack Pro, OpenHelm, Marlinspike Log, Saltmarsh
  Ledger;
- their terms: track fusion, leg fingerprint, the 90-second rule, primary
  device, gap fill, HDOP and HDOP gating, merge candidate, twin badge,
  `mergeTracks`;
- NMEA terms: talker ID, `$GP`, `$AI`, `$II`, GPRMC, GPGSA.

PASS: every one of these names that the reply uses is introduced where it
first appears: the reply says what it is in a clause (for example
"Kittiwake, a Norwegian charter-fleet logbook", "HDOP, the receiver's own
measure of fix quality") and that it comes from the scout's research or
names the tool or source it belongs to. One introduction can cover a group
("the scout found five tools: ..."), if it says what each is. Names in an
evidence or sources section at the end need no introduction if they were
introduced above or appear only there. A reply that uses none of these
names passes.

FAIL: the reply uses any of these names as if the user already knew it:
with no word on what it is, or with no word on where it came from (for
example "track fusion, as Kittiwake does it" with Kittiwake never
introduced; "HDOP gating" with no explanation; "(OpenHelm, Kittiwake)" as
examples the user is assumed to know).
