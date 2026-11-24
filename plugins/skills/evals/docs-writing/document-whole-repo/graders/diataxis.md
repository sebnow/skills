---
type: llm
focus: trace
weight: 1
---
The agent wrote documentation for a service repository. Judge only the Markdown pages it wrote (the Write and Edit calls).

Diátaxis splits documentation into four types: tutorials (learning by doing), how-to guides (steps for one task), reference (facts to look up), and explanation (understanding and reasons).

PASS if every page is mainly one of these types and links to the others instead of mixing them, and the set covers both readers: people who use or integrate with the service (at least one how-to guide or reference page for them) and engineers who change it (at least one explanation of the architecture).

FAIL if any page mixes types at length (for example an architecture explanation that also holds step-by-step setup instructions and a full API reference), or if the set covers only one of the two readers.
