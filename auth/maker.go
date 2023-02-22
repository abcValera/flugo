package auth

import "time"

// Interface for managing tokens
type Maker interface {
	CreateToken(UserID int32, username, email string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}
