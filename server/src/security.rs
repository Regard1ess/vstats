use std::{
    collections::{HashMap, VecDeque},
    sync::Arc,
    time::{Duration, Instant},
};
use tokio::sync::Mutex;

/// Simple in-memory rate limiter keyed by client IP.
#[derive(Clone)]
pub struct RateLimiter {
    buckets: Arc<Mutex<HashMap<String, VecDeque<Instant>>>>,
    max_requests: usize,
    window: Duration,
}

impl RateLimiter {
    pub fn new(max_requests: usize, window: Duration) -> Self {
        Self {
            buckets: Arc::new(Mutex::new(HashMap::new())),
            max_requests,
            window,
        }
    }

    /// Returns true when the caller is allowed to proceed.
    pub async fn is_allowed(&self, ip: &str) -> bool {
        let mut guard = self.buckets.lock().await;
        let entry = guard.entry(ip.to_string()).or_insert_with(VecDeque::new);
        let now = Instant::now();

        while let Some(&ts) = entry.front() {
            if now.duration_since(ts) > self.window {
                entry.pop_front();
            } else {
                break;
            }
        }

        if entry.len() >= self.max_requests {
            return false;
        }

        entry.push_back(now);
        true
    }
}

/// Sanitize free-form text to mitigate stored XSS vectors.
pub fn sanitize_text(input: String) -> String {
    let trimmed = input.trim();
    let filtered: String = trimmed
        .chars()
        .filter(|c| !matches!(c, '<' | '>' | '`'))
        .take(200)
        .collect();
    filtered
}

/// Restrictive sanitizer for host/url-like strings (no scriptable characters).
pub fn sanitize_host(input: String) -> String {
    let trimmed = input.trim();
    trimmed
        .chars()
        .filter(|c| c.is_ascii_alphanumeric() || *c == '.' || *c == '-' || *c == ':' || *c == '/')
        .take(200)
        .collect()
}
