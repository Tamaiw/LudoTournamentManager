package domain

import "errors"

var (
	ErrNotFound           = errors.New("entity not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrTournamentActive   = errors.New("tournament is active and cannot be modified")
	ErrGameAlreadyPlayed  = errors.New("game has already been played")
	ErrInvalidAdvancement = errors.New("advancement configuration is invalid")
	ErrNoRematch          = errors.New("players from same source game cannot be seated together")
)
