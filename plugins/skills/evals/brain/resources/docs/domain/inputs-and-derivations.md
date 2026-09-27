# Inputs and derivations

What the return is computed from, and what is computed. The rules in `internal/rules/ccvr` read inputs and
produce an assessment; nothing they produce is stored.

## Inputs

Inputs come in six shapes: imported records, elections, mooring links, operator facts, hand entries and carried balances.

Imported records are voyages read from chartplotter CSV exports and NMEA dumps. They are immutable (ADR-0003); a
wrong import is removed, not edited, and everything that cites it goes with it.

Everything else is stated by the operator. A statement is an input the return relies on, so the return shows each
one with the date it was made and what backs it.

## Shapes

A choice between the options a clause offers is recorded against the clause, and a later choice on the same
clause replaces it in force.

An election picks an option; it never states a quantity. A hand entry states a quantity and never picks an option.

A yes/no fact about how the vessel is run is about the vessel and the season.

A quantity the operator types in carries the document it rests on when there is one.

A hand entry states a quantity. Zero is a quantity.

A mooring link holds for a period. Two links for disjoint periods coexist; a later link for an overlapping period replaces the overlap only.

## Absence

An absent input and an input stated as nothing are different: the first is an open finding, the second is a
statement the return relies on.

## Balances brought forward
A carried balance is stated once per vessel and applies from the first imported season.

It is the starting point for a quantity the imported logs cannot reconstruct, such as engine hours run before the
first imported voyage.

## Derivations

Derived data (hours towards overhaul, survey due date, levy band) is never stored; the return recomputes it from inputs.

A derivation that needs an input nobody has stated uses the assumption the rule declares and raises an open finding
that names the assumption.
