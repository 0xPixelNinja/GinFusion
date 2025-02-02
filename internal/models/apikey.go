package models

// APIKey represents an API key associated with a user.
type APIKey struct {
    Key         string `json:"key"`
    UserID      string `json:"user_id"`
    RateLimit   int    `json:"rate_limit"`   // Requests per minute allowed.
    Concurrency int    `json:"concurrency"`  // Maximum concurrent requests allowed.
    Created     int64  `json:"created_at"`
}
