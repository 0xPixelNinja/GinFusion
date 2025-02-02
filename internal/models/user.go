package models

// User represents a user in the system.
type User struct {
    ID       string `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"` // Stored as a bcrypt hash.
}
