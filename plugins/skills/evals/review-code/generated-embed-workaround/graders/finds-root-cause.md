---
type: llm
focus: trace
weight: 1
---
`internal/tax/pl` needs the ruling texts that live in `docs/rulings/`. Because `go:embed` only reaches files in or below the embedding package's directory, the change adds `cmd/gen-rulings`, a `//go:generate` directive, and a committed generated copy of the texts, `internal/tax/pl/rulings_text.go`. The generator is a workaround for a structural mistake: the consumer package reaches for files another directory owns. The Go-native fix is for the owner to expose them — a package in `docs/` (e.g. `docs/docs.go` with `//go:embed rulings/*/*/text.md` and an exported `embed.FS`) that `internal/tax/pl` imports. That deletes the generator, the `go:generate` directive, and the generated file, and removes the stale-copy risk.

PASS only if the review identifies the generator and its generated file as machinery that exists because of where the embedding happens, and recommends restructuring so the texts are embedded directly — preferably by the directory that owns them (a package under `docs/` exposing an `embed.FS` or equivalent, imported by the tax package) — with the generator, `go:generate` directive, and generated file deleted. Offering relocation of the texts as a secondary alternative does not fail the review.

FAIL if the review keeps the generator — for example recommending a CI check that the generated file is up to date, fixing the generator's relative paths, or improving its error handling, without recommending its removal — or loads the texts from disk at runtime. Those treat the symptom and keep the cause.
