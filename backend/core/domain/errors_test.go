package domain

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorSentinels(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrNotFound", ErrNotFound},
		{"ErrInvalidInput", ErrInvalidInput},
		{"ErrUnauthorized", ErrUnauthorized},
		{"ErrForbidden", ErrForbidden},
		{"ErrTournamentActive", ErrTournamentActive},
		{"ErrGameAlreadyPlayed", ErrGameAlreadyPlayed},
		{"ErrInvalidAdvancement", ErrInvalidAdvancement},
		{"ErrNoRematch", ErrNoRematch},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("expected %s to be non-nil", tt.name)
			}
			if tt.err.Error() == "" {
				t.Errorf("expected %s to have a message", tt.name)
			}
		})
	}
}

func TestErrorEquality(t *testing.T) {
	if !errors.Is(ErrNotFound, ErrNotFound) {
		t.Error("ErrNotFound should be equal to itself")
	}
	if !errors.Is(ErrInvalidInput, ErrInvalidInput) {
		t.Error("ErrInvalidInput should be equal to itself")
	}
}

func TestErrorWrapping(t *testing.T) {
	// Test using fmt.Errorf for proper error wrapping
	wrapped := fmt.Errorf("wrapped: %w", ErrNotFound)
	if !errors.Is(wrapped, ErrNotFound) {
		t.Error("wrapped error should contain ErrNotFound")
	}
}

func TestErrorMessages(t *testing.T) {
	tests := map[string]error{
		"entity not found":                ErrNotFound,
		"invalid input":                   ErrInvalidInput,
		"unauthorized":                     ErrUnauthorized,
		"forbidden":                       ErrForbidden,
		"tournament is active and cannot be modified": ErrTournamentActive,
		"game has already been played":    ErrGameAlreadyPlayed,
		"advancement configuration is invalid": ErrInvalidAdvancement,
		"players from same source game cannot be seated together": ErrNoRematch,
	}

	for expectedMsg, err := range tests {
		if err.Error() != expectedMsg {
			t.Errorf("expected error message %q, got %q", expectedMsg, err.Error())
		}
	}
}