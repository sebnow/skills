---
name: architecture-visualization
description: "Breaks a piece of architecture, in existing code or a proposed design, into concept-level parts (services, workers, stores, outside systems) and shows it as a rendered graphic: a Mermaid or SVG block, or an HTML file. Use when the user asks how a system, service, or repository works, wants the big picture or an overview, asks to explain, walk through, or break down a flow or a design, asks how the parts fit together, or asks for an architecture diagram or sketch. Not for charts of data, and not for reviewing or changing a design (use review-design)."
---

# Architecture Visualization

Answer with a picture of the architecture, and use words only for what the
picture cannot show.

## One picture per question

Start from the question the user asked, and pick the view that answers it:

- What the system is and how its parts fit: the running parts, the stores and
  the outside systems, and what connects them.
- What happens when something occurs ("walk me through an upload"): the parts
  involved in that one scenario, in order.
- What states a key thing moves through: the states and what moves it between
  them.

Each picture stays at one level. A big-picture question can take the parts
view followed by one scenario view.

## Parts are concepts, not code

A part is something that runs, stores data, or sits outside the system: an
API, a worker, a queue, a database, a third-party service. Name it as the
reader would, and give it a one-line responsibility. Label each arrow with
what it is for ("claims the next job"), not with the call that does it.

Read the code to find the parts, then describe the parts, not the code.
Class, function and file names, file:line references, SQL, constants, retry
formulas and the directory layout stay out of the picture and the text,
unless the user's question is about them.

## Draw it in a rendered format

Use the first of these that can show the picture:

1. A fenced `mermaid` block: flowcharts for parts, sequence diagrams for
   scenarios, state diagrams for states.
2. A fenced `svg` block holding an inline SVG element, when the picture needs
   a shape or layout Mermaid cannot draw.
3. An HTML file, only when the reader needs to step through a flow or toggle
   layers.

Put the block in the reply even where you expect it to show as plain text;
the interface decides how to render it. Never
draw with box-drawing characters or ASCII art, and do not let prose, a table,
or a numbered list stand in for the picture.

## Answer the question asked

After the picture, a few sentences cover what it cannot show. Leave out
suspected bugs, risks and open questions the user did not ask about. If one
is serious enough that the user should hear it now, give it one sentence at
the end.
