from dataclasses import dataclass
from datetime import datetime


@dataclass
class Job:
    id: int
    doc_id: str
    kind: str
    attempts: int
    run_after: datetime
    last_error: str | None
