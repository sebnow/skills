import smtplib
from datetime import UTC, datetime, timedelta

import pytest

from pantry import accounts, carts, db, reminders
from pantry.mail import Mailer

T0 = datetime(2026, 5, 1, 12, 0, tzinfo=UTC)


@pytest.fixture
def conn():
    conn = db.connect(":memory:")
    db.migrate(conn)
    conn.execute("INSERT INTO products VALUES ('egg-12', 'Eggs, dozen', 450, 10)")
    return conn


def customer(conn, email="ana@example.com"):
    return accounts.verify(conn, accounts.sign_up(conn, email, "Ana", T0), T0)


def stock(conn):
    return conn.execute("SELECT stock FROM products WHERE sku = 'egg-12'").fetchone()[0]


class FakeMailer(Mailer):
    """Composes real messages and records the recipients instead of sending."""

    def __init__(self):
        super().__init__("relay", "shop@example.com")
        self.sent = []

    def connect(self):
        pass

    def send(self, msg):
        self.sent.append(msg["To"])


def test_idle_cart_is_closed_before_its_stock_is_returned(conn):
    cart = carts.open_cart(conn, customer(conn), T0)
    carts.add_item(conn, cart, "egg-12", 3, T0)

    with pytest.raises(carts.CartClosed):
        carts.checkout(conn, cart, T0 + timedelta(hours=3))
    assert carts.expire_idle(conn, T0 + timedelta(hours=1)) == 0
    assert carts.expire_idle(conn, T0 + timedelta(hours=3)) == 1
    assert stock(conn) == 10


def test_cart_in_use_keeps_its_stock(conn):
    cart = carts.open_cart(conn, customer(conn), T0)
    carts.add_item(conn, cart, "egg-12", 3, T0)
    carts.add_item(conn, cart, "egg-12", 1, T0 + timedelta(hours=1, minutes=59))

    assert carts.expire_idle(conn, T0 + timedelta(hours=3)) == 0
    assert stock(conn) == 6


def test_lapsed_signup_cannot_be_verified(conn):
    token = accounts.sign_up(conn, "late@example.com", "Lee", T0)
    with pytest.raises(accounts.InvalidToken):
        accounts.verify(conn, token, T0 + timedelta(days=8))


def test_unverified_signup_releases_its_address_after_a_week(conn):
    accounts.sign_up(conn, "late@example.com", "Lee", T0)
    with pytest.raises(accounts.EmailTaken):
        accounts.sign_up(conn, "late@example.com", "Lee", T0 + timedelta(days=6))
    accounts.sign_up(conn, "late@example.com", "Lee", T0 + timedelta(days=8))


def test_purge_deletes_only_lapsed_signups(conn):
    accounts.sign_up(conn, "late@example.com", "Lee", T0)
    customer(conn)

    assert accounts.purge_lapsed(conn, T0 + timedelta(days=8)) == 1
    assert conn.execute("SELECT email FROM accounts").fetchall() == [("ana@example.com",)]


def test_dormant_customer_is_reminded_once_per_spell(conn):
    account_id = customer(conn)
    carts.checkout(conn, add_eggs(conn, account_id, T0), T0)
    mailer = FakeMailer()

    reminders.remind_dormant(conn, mailer, T0 + timedelta(days=61))
    reminders.remind_dormant(conn, mailer, T0 + timedelta(days=62))
    assert mailer.sent == ["ana@example.com"]

    later = T0 + timedelta(days=70)
    carts.checkout(conn, add_eggs(conn, account_id, later), later)
    reminders.remind_dormant(conn, mailer, T0 + timedelta(days=131))
    assert len(mailer.sent) == 2


def add_eggs(conn, account_id, when):
    cart = carts.open_cart(conn, account_id, when)
    carts.add_item(conn, cart, "egg-12", 1, when)
    return cart


def test_idle_carts_return_stock_together(conn):
    conn.execute("INSERT INTO products VALUES ('milk-1', 'Milk, 1 l', 120, 5)")
    for email in ("ana@example.com", "ben@example.com"):
        cart = carts.open_cart(conn, customer(conn, email), T0)
        carts.add_item(conn, cart, "egg-12", 2, T0)
        carts.add_item(conn, cart, "milk-1", 1, T0)

    assert carts.expire_idle(conn, T0 + timedelta(hours=3)) == 2
    assert stock(conn) == 10
    assert conn.execute("SELECT stock FROM products WHERE sku = 'milk-1'").fetchone()[0] == 5


def test_refused_address_does_not_stop_the_others(conn):
    class FlakyMailer(FakeMailer):
        def send(self, msg):
            if msg["To"] == "ana@example.com":
                raise smtplib.SMTPRecipientsRefused({msg["To"]: (550, b"no such user")})
            super().send(msg)

    for email in ("ana@example.com", "ben@example.com"):
        account_id = customer(conn, email)
        carts.checkout(conn, add_eggs(conn, account_id, T0), T0)
    mailer = FlakyMailer()

    assert reminders.remind_dormant(conn, mailer, T0 + timedelta(days=61)) == 1
    assert mailer.sent == ["ben@example.com"]


@pytest.mark.parametrize(
    "error",
    [
        smtplib.SMTPServerDisconnected("Connection unexpectedly closed"),
        smtplib.SMTPResponseException(421, b"closing connection"),
    ],
)
def test_relay_fumbling_quit_does_not_fail_the_run(monkeypatch, error):
    class FumblingSMTP:
        def __init__(self, host, timeout):
            self.closed = False

        def send_message(self, msg):
            pass

        def quit(self):
            raise error

        def close(self):
            self.closed = True

    monkeypatch.setattr(smtplib, "SMTP", FumblingSMTP)
    with Mailer("relay", "shop@example.com") as mailer:
        mailer.send(mailer.compose("we-miss-you", to="ana@example.com", first_name="Ana"))
        smtp = mailer._smtp
    assert smtp.closed and mailer._smtp is None


def dormant(conn, email):
    account_id = customer(conn, email)
    carts.checkout(conn, add_eggs(conn, account_id, T0), T0)


class FailingMailer(FakeMailer):
    def __init__(self, error):
        super().__init__()
        self.error = error

    def send(self, msg):
        if self.error:
            raise self.error
        super().send(msg)


def test_deferred_reminder_is_retried_on_the_next_run(conn):
    dormant(conn, "ana@example.com")
    mailer = FailingMailer(smtplib.SMTPRecipientsRefused({"ana@example.com": (450, b"try again later")}))

    assert reminders.remind_dormant(conn, mailer, T0 + timedelta(days=61)) == 0
    mailer.error = None
    assert reminders.remind_dormant(conn, mailer, T0 + timedelta(days=62)) == 1


def test_relay_wide_failure_stops_the_run_and_keeps_the_reminder_due(conn):
    dormant(conn, "ana@example.com")
    dormant(conn, "ben@example.com")
    mailer = FailingMailer(smtplib.SMTPSenderRefused(553, b"bad sender", "shop@example.com"))

    with pytest.raises(smtplib.SMTPSenderRefused):
        reminders.remind_dormant(conn, mailer, T0 + timedelta(days=61))
    mailer.error = None
    assert reminders.remind_dormant(conn, mailer, T0 + timedelta(days=62)) == 2


def test_reminder_that_cannot_be_written_is_skipped_unrecorded(conn):
    bad = customer(conn, "lee@example.com")
    conn.execute("UPDATE accounts SET first_name = 'Lee\nBcc: everyone@example.com' WHERE id = ?", (bad,))
    carts.checkout(conn, add_eggs(conn, bad, T0), T0)
    dormant(conn, "ana@example.com")
    mailer = FakeMailer()

    assert reminders.remind_dormant(conn, mailer, T0 + timedelta(days=61)) == 1
    assert mailer.sent == ["ana@example.com"]
    assert conn.execute("SELECT reminded_at FROM accounts WHERE id = ?", (bad,)).fetchone() == (None,)


def test_message_rejected_for_good_does_not_stop_the_run(conn):
    dormant(conn, "ana@example.com")
    dormant(conn, "ben@example.com")

    class RejectingMailer(FakeMailer):
        def send(self, msg):
            if msg["To"] == "ana@example.com":
                raise smtplib.SMTPDataError(554, b"message rejected")
            super().send(msg)

    mailer = RejectingMailer()
    assert reminders.remind_dormant(conn, mailer, T0 + timedelta(days=61)) == 1
    assert reminders.remind_dormant(conn, mailer, T0 + timedelta(days=62)) == 0
    assert mailer.sent == ["ben@example.com"]
