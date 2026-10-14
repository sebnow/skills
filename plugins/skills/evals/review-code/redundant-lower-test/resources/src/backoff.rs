use std::time::Duration;

const BASE_MS: u64 = 1_000;
const CAP_MS: u64 = 60_000;

/// Delay after failed send number `attempt` (1-based): 1 s, doubling each
/// time, never more than 60 s.
pub fn backoff(attempt: u32) -> Duration {
    let factor = 1u64.checked_shl(attempt.saturating_sub(1)).unwrap_or(u64::MAX);
    Duration::from_millis(BASE_MS.saturating_mul(factor).min(CAP_MS))
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn backoff_starts_at_one_second_and_doubles() {
        assert_eq!(backoff(1), Duration::from_secs(1));
        assert_eq!(backoff(2), Duration::from_secs(2));
    }

    #[test]
    fn backoff_stays_at_the_cap_however_many_attempts() {
        assert_eq!(backoff(10), Duration::from_secs(60));
        assert_eq!(backoff(64), Duration::from_secs(60));
        assert_eq!(backoff(u32::MAX), Duration::from_secs(60));
    }
}
