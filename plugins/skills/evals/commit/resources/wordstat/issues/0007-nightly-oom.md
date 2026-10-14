# 0007: nightly docs-stats job killed for running out of memory

The nightly docs-stats job runs `wordstat` over every file under
`exports/` on the shared CI runner, which has 512 MiB of memory. Since
the vendor exports grew past 40 MiB, the job is killed by the kernel's
OOM killer about once a week.

`read_text` reads the whole file and then decodes it, so a file costs
about three times its size in memory while it is counted; with the
current 64 MiB limit that is close to 200 MiB for a single file, on top
of the job's other processes.

Until `wordstat` streams its input, the limit should stay at or below
16 MiB. The exports above that are archives nobody reads; skipping them
is fine.
