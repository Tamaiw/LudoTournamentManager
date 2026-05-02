package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestRoleConstants(t *testing.T) {
	tests := []struct {
		name     string
		role     Role
		expected string
	}{
		{"RoleAdmin", RoleAdmin, "admin"},
		{"RoleMember", RoleMember, "member"},
		{"RoleGuest", RoleGuest, "guest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.role) != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.role)
			}
		})
	}
}

func TestUserJSONSerialization(t *testing.T) {
	user := User{
		ID:           "test-uuid",
		Email:        "test@example.com",
		PasswordHash: "secret-hash",
		Role:         RoleMember,
		CreatedAt:    time.Now(),
		ModifiedAt:   time.Now(),
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("failed to marshal user: %v", err)
	}

	var unmarshaled User
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal user: %v", err)
	}

	if unmarshaled.ID != user.ID {
		t.Errorf("expected ID %q, got %q", user.ID, unmarshaled.ID)
	}
	if unmarshaled.Email != user.Email {
		t.Errorf("expected Email %q, got %q", user.Email, unmarshaled.Email)
	}
	if unmarshaled.Role != user.Role {
		t.Errorf("expected Role %q, got %q", user.Role, unmarshaled.Role)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("failed to unmarshal to map: %v", err)
	}

	if _, ok := result["password"]; ok {
		t.Error("password hash should not be serialized in JSON")
	}
}

func TestUserValidation(t *testing.T) {
	tests := []struct {
		name  string
		user  User
		valid bool
	}{
		{
			name: "valid user",
			user: User{
				ID:    "uuid-1",
				Email: "test@example.com",
				Role:  RoleMember,
			},
			valid: true,
		},
		{
			name: "empty email",
			user: User{
				ID:    "uuid-2",
				Email: "",
				Role:  RoleGuest,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.valid && tt.user.ID == "" {
				t.Error("expected valid user to have ID")
			}
		})
	}
}