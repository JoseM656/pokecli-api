package model

// LocalizedName represents a name in multiple languages.
type LocalizedName struct {
	ES string `json:"es,omitempty"`
	EN string `json:"en,omitempty"`
	JA string `json:"ja,omitempty"`
	FR string `json:"fr,omitempty"`
	DE string `json:"de,omitempty"`
	KO string `json:"ko,omitempty"`
	IT string `json:"it,omitempty"`
	ZH string `json:"zh,omitempty"`
}

// PokemonType represents a pokemon type with its localized name and display color.
type PokemonType struct {
	Name  LocalizedName `json:"name" bson:"name"`
	Color string        `json:"color" bson:"color"`
}

// PokemonAbility represents a pokemon ability with its localized name.
type PokemonAbility struct {
	Name   LocalizedName `json:"name" bson:"name"`
	Hidden bool          `json:"hidden" bson:"hidden"`
}

// PokemonStats represents the base stats of a pokemon.
type PokemonStats struct {
	HP        int `json:"hp" bson:"hp"`
	Attack    int `json:"attack" bson:"attack"`
	Defense   int `json:"defense" bson:"defense"`
	SpAttack  int `json:"sp_attack" bson:"sp_attack"`
	SpDefense int `json:"sp_defense" bson:"sp_defense"`
	Speed     int `json:"speed" bson:"speed"`
}

// Pokemon represents the projected model returned by the API.
type Pokemon struct {
	ID        int              `json:"id" bson:"_id"`
	Name      LocalizedName    `json:"name" bson:"name"`
	Types     []PokemonType    `json:"types" bson:"types"`
	Abilities []PokemonAbility `json:"abilities" bson:"abilities"`
	Stats     PokemonStats     `json:"stats" bson:"stats"`
	SpriteURL string           `json:"sprite_url" bson:"sprite_url"`
	Cached    bool             `json:"cached" bson:"-"`
	Lang      string           `json:"lang" bson:"-"`
}
