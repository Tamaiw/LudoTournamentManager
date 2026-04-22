package inbound

import "context"

type UserService interface {
	ListUsers(ctx context.Context) ([]UserDTO, error)
	UpdateUser(ctx context.Context, id string, role string) error
	DeleteUser(ctx context.Context, id string) error
	SendInvite(ctx context.Context, email string, inviterID string, inviteType string) (string, error) // returns code
	AcceptInvite(ctx context.Context, code string, email, password string) (string, error)             // returns JWT
}
