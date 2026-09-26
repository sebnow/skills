---
type: llm
focus: last_message
weight: 1
---
The draft contained the sentence "Beancount syntax is documented elsewhere and is not covered here." It points the reader elsewhere without saying where, so a reader cannot follow it.

PASS if either the revised document names a concrete location for Beancount's documentation (a URL or a named document), or the message explicitly tells the user that this sentence does not say where the documentation is and asks for or suggests the location.

FAIL if the sentence is kept as is, reworded without a location (for example "Beancount syntax is not covered here"), or deleted without being raised to the user.
