---
type: llm
focus: trace
weight: 1
---
The user described a reminder system for a clinic booking product and asked to see how it hangs together. The design's parts are: the booking service publishing appointment events to an "appointments" Kafka topic; a reminder planner consuming them and keeping reminder plans in its own Postgres database; a dispatcher that every minute picks due reminders and asks the patient-preferences service for channel and opt-out; channel adapters for SMS (Twilio), email (SES) and push (FCM); provider webhooks to a receipts endpoint that updates reminder status in the planner's database, with SMS falling back to email after two failures; a dead-letter topic with paging past 100; and a staff dashboard reading the reporting warehouse, fed from the planner's tables by CDC.

Judge the agent's explanation, and any picture it made (in the reply, or in a file it wrote), on its level of detail, not on its medium.

PASS if the units of the explanation are these parts and the flows between them, at the level the user described them, so the user can check the design as a whole. Pointing out a gap or a question about the design is fine.

FAIL if the agent goes below the design's level: it invents classes, function signatures, code, database schemas or message field lists and those make up a large share of the answer, or it breaks the design into implementation tasks instead of showing how the parts connect.
