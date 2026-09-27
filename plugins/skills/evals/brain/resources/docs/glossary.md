# Glossary

Terms used in the code, the requirements and the return. Where a term comes from the Regulations, the regulation is named.

- **Approved list**: the yards the Authority lists under reg. 31 as fit for hull work on a vessel the Regulations apply to.

- **Assessment**: what the rules compute for one season from the inputs: hours towards overhaul, survey due date, levy band and the open findings.

- **Blob**: a stored file, kept byte for byte. Imports and documents are blobs.


- **Carried balance**: a quantity that accrued before the period the imported logs cover and that a rule needs to start from, such as engine hours since the last overhaul. Stated by the operator; see hand entry.

- **Clash**: two imported voyages that overlap in time for the same vessel. The later import does not replace the earlier; the clash is a finding.

- **Coastal zone**: the waters named in reg. 3. Whether a vessel's home port lies in it decides whether the Regulations apply.


- **Election**: the option the operator picks where a clause of the regulation offers more than one, recorded so a rule can apply it. An election is a choice, not a fact: mooring links, operator facts and hand entries are not elections.

- **Evidence**: a blob an input rests on, optionally narrowed to one record inside it by a locator.

- **Finding**: something the return has to tell the reader. An open finding blocks submission.



- **Hand entry**: a value the operator types in because no import carries it, backed by a document when one exists.

- **Import**: one voyage log file read into the ledger. Imports are immutable; removing one leaves a tombstone.

- **Ledger**: the SQLite file that holds every input (ADR-0004).


- **Mooring link**: the relation between a vessel and a jurisdiction that decides which rules apply to it, such as home-port registration. A mooring link holds for a period.

- **Open finding**: a finding the operator has to resolve before the return can be submitted.

- **Operating return**: the yearly return an operator makes under reg. 6. tidewright renders it as Markdown.


- **Operator fact**: a yes/no fact about how the operator runs the vessel that changes which rule applies, such as operating for hire or reward.

- **Season**: the calendar year the return covers (reg. 7).

- **Tombstone**: the record that an import was removed, kept so the removal can be explained.


- **Unrecorded input**: an input a rule depends on that nothing in the tool can capture today. Every unrecorded input is an open finding.

- **Voyage**: one passage from casting off to making fast, with its engine hours per engine.
