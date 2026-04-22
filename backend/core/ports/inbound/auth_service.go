package inbound

import "context"

type AuthService interface {
	Register(ctx context.Context, email, password, inviteCode string) (string, error) // returns JWT
	Login(ctx context.Context, email, password string) (string, error)                // returns JWT
	Logout(ctx context.Context, token string) error
	GetCurrentUser(ctx context.Context, token string) (*UserDTO, error)
}

type UserDTO struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
