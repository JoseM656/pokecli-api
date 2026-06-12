package model

// MoveCategory represents the damage class of a move.
type MoveCategory struct {
	Name  string `json:"name" bson:"name"`
	Color string `json:"color" bson:"color"`
}

// Move represents the projected model returned by the API.
type Move struct {
	Name     LocalizedName `json:"name" bson:"name"`
	Type     PokemonType   `json:"type" bson:"type"`
	Category MoveCategory  `json:"category" bson:"category"`
	Power    int           `json:"power" bson:"power"`
	Accuracy int           `json:"accuracy" bson:"accuracy"`
	PP       int           `json:"pp" bson:"pp"`
	Cached   bool          `json:"cached" bson:"-"`
	Lang     string        `json:"lang" bson:"-"`
}
