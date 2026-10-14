//! Delivers webhook payloads, retrying server errors with backoff.

mod backoff;

use std::io::{self, Read, Write};
use std::net::TcpStream;
use std::time::Duration;

pub use backoff::backoff;

#[derive(Debug, PartialEq)]
pub enum Failure {
    /// The receiver refused the payload; sending it again cannot help.
    Rejected(u16),
    /// Every attempt ended in a server error or a connection error.
    GaveUp { attempts: u32 },
}

/// Sends `body` until the receiver accepts it (2xx), refuses it (4xx), or
/// `max_attempts` sends have failed, waiting `backoff(n)` after failed send n.
/// Returns the number of sends it took.
pub fn deliver(
    body: &[u8],
    max_attempts: u32,
    mut send: impl FnMut(&[u8]) -> io::Result<u16>,
    mut sleep: impl FnMut(Duration),
) -> Result<u32, Failure> {
    for attempt in 1..=max_attempts {
        match send(body) {
            Ok(status) if (200..300).contains(&status) => return Ok(attempt),
            Ok(status) if (400..500).contains(&status) => return Err(Failure::Rejected(status)),
            _ => {}
        }
        if attempt < max_attempts {
            sleep(backoff(attempt));
        }
    }
    Err(Failure::GaveUp { attempts: max_attempts })
}

/// Posts `body` to `http://{addr}/hook` and returns the response status.
pub fn post(addr: &str, body: &[u8]) -> io::Result<u16> {
    let mut stream = TcpStream::connect(addr)?;
    let mut request = format!("POST /hook HTTP/1.0\r\nContent-Length: {}\r\n\r\n", body.len()).into_bytes();
    request.extend_from_slice(body);
    stream.write_all(&request)?;
    let mut response = String::new();
    stream.read_to_string(&mut response)?;
    response
        .split_whitespace()
        .nth(1)
        .and_then(|code| code.parse().ok())
        .ok_or_else(|| io::Error::new(io::ErrorKind::InvalidData, "no status line"))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn retries_server_errors_with_backoff_until_accepted() {
        let mut replies = [503, 503, 200].into_iter();
        let mut waits = Vec::new();
        let result = deliver(b"{}", 5, |_| Ok(replies.next().unwrap()), |d| waits.push(d));
        assert_eq!(result, Ok(3));
        assert_eq!(waits, [Duration::from_secs(1), Duration::from_secs(2)]);
    }

    #[test]
    fn a_refused_payload_is_not_sent_again() {
        let mut sends = 0;
        let result = deliver(b"{}", 5, |_| { sends += 1; Ok(422) }, |_| {});
        assert_eq!(result, Err(Failure::Rejected(422)));
        assert_eq!(sends, 1);
    }
}
