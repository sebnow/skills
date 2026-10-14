//! End to end against the receiving partner's live sandbox over the internet.
//!
//! The sandbox answers 503 to the first two deliveries of a body it has not
//! seen, then 200. It answers 429 to any delivery that arrives sooner after
//! the previous one than our documented backoff, so this test must sleep
//! through the real delays: it is the only check that our real waits satisfy
//! the receiver's rate limit, and that a real receiver accepts what `post`
//! puts on the wire.

use std::thread;
use std::time::{SystemTime, UNIX_EPOCH};

use hookpost::{deliver, post};

#[test]
#[ignore = "costly: real partner sandbox over the internet (needs HOOKPOST_SANDBOX_ADDR), rate-limited, sleeps through the real 1 s + 2 s backoff; run with --ignored"]
fn the_sandbox_accepts_after_two_outages_without_rate_limiting_us() {
    let addr = std::env::var("HOOKPOST_SANDBOX_ADDR").expect("HOOKPOST_SANDBOX_ADDR is not set");
    let nonce = SystemTime::now().duration_since(UNIX_EPOCH).unwrap().as_nanos();
    let body = format!(r#"{{"event":"paid","nonce":{nonce}}}"#);

    let result = deliver(body.as_bytes(), 5, |b| post(&addr, b), thread::sleep);

    assert_eq!(result, Ok(3));
}
