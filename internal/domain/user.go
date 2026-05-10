package domain

import "context"

// User is a minimal representation of an account in the system.
type User struct {
	ID    string
	Email string
	Name  string
}

// UserRepository describes the minimal persistence operation used by LCP.
type UserRepository interface {
	Ensure(ctx context.Context, user *User) error
}
