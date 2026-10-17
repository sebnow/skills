import os
import smtplib
from email.message import EmailMessage

TEMPLATES = {
    "we-miss-you": (
        "We miss you, {first_name}",
        "Hi {first_name},\n\nIt has been a while. Your favourites are in stock.\n",
    ),
}


class Mailer:
    """Sends templated email through the local relay.

    Connects on connect() or the first send, and reuses the connection until
    the `with` block ends.
    """

    def __init__(self, host: str, sender: str, timeout: float = 30):
        self.host = host
        self.sender = sender
        self.timeout = timeout
        self._smtp: smtplib.SMTP | None = None

    @classmethod
    def from_env(cls) -> "Mailer":
        return cls(os.environ["PANTRY_SMTP_HOST"], os.environ["PANTRY_MAIL_FROM"])

    def __enter__(self) -> "Mailer":
        return self

    def __exit__(self, exc_type, exc, tb) -> None:
        if self._smtp is None:
            return
        try:
            if exc_type is None:
                self._smtp.quit()
        except (smtplib.SMTPException, OSError):
            # Everything was sent; a relay that fumbles QUIT does not undo that.
            pass
        finally:
            self._smtp.close()
            self._smtp = None

    def connect(self) -> None:
        if self._smtp is None:
            self._smtp = smtplib.SMTP(self.host, timeout=self.timeout)

    def compose(self, template: str, to: str, **context: str) -> EmailMessage:
        """Build a message from a template. Raises ValueError if a header would
        hold a line break."""
        subject, body = TEMPLATES[template]
        msg = EmailMessage()
        msg["From"] = self.sender
        msg["To"] = to
        msg["Subject"] = subject.format(**context)
        msg.set_content(body.format(**context))
        return msg

    def send(self, msg: EmailMessage) -> None:
        self.connect()
        self._smtp.send_message(msg)
