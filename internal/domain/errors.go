package domain

import "errors"

// Common domain errors
var (
	ErrPokemonNotFound = errors.New("pokemon not found")
	ErrInternal        = errors.New("internal server error")
)
