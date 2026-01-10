package domain

import "time"

// Upload represents a stored user file
type Upload struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    Path      string    `json:"path"`
    MIME      string    `json:"mime"`
    Size      int64     `json:"size"`
    CreatedAt time.Time `json:"created_at"`
}
