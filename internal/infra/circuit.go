package infra

import (
    "time"

    "github.com/sony/gobreaker"
)

// NewBreaker creates a circuit‑breaker with sensible defaults.
// `timeout` is used both for the internal interval and the open‑state timeout.
func NewBreaker(name string, timeout time.Duration) *gobreaker.CircuitBreaker {
    settings := gobreaker.Settings{
        Name:        name,
        MaxRequests: 5,               // allow a few requests when half‑open
        Interval:    timeout,         // reset counters after this period
        Timeout:     timeout,         // how long to stay open before trying again
        ReadyToTrip: func(counts gobreaker.Counts) bool {
            // Open the circuit after 5 consecutive failures
            return counts.ConsecutiveFailures > 5
        },
    }
    return gobreaker.NewCircuitBreaker(settings)
}
