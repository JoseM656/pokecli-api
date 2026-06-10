package model

import "errors"

// Sentinel errors for domain-level error handling.
var (
	ErrPokemonNotFound = errors.New("pokemon not found")
)
