package pokeapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/pokeservicies/pokeproxy/internal/model"
)

// typeColors maps PokeAPI type names to their display hex colors.
var typeColors = map[string]string{
	"normal":   "#A8A878",
	"fire":     "#c3320d",
	"water":    "#4592C4",
	"electric": "#FAE078",
	"grass":    "#9BCC50",
	"ice":      "#51C4E7",
	"fighting": "#dc1010",
	"poison":   "#B97FC9",
	"ground":   "#845910",
	"flying":   "#3DC7EF",
	"psychic":  "#F366B9",
	"bug":      "#729F3F",
	"rock":     "#A38C21",
	"ghost":    "#7B62A3",
	"dragon":   "#F16E57",
	"dark":     "#707070",
	"steel":    "#9EB7B8",
	"fairy":    "#FDB9E9",
}

// categoryColors maps PokeAPI damage class names to display hex colors.
var categoryColors = map[string]string{
	"physical": "#8B1A1A",
	"special":  "#1A4B8C",
	"status":   "#4A4A4A",
}

// Map maps a rawPokemon and its localized data in the model.
func (c *Client) Map(ctx context.Context, raw *rawPokemon) (*model.Pokemon, error) {
	types, err := c.mapTypes(ctx, raw)
	if err != nil {
		return nil, err
	}

	abilities, err := c.mapAbilities(ctx, raw)
	if err != nil {
		return nil, err
	}

	// Fetch localized names from species endpoint
	species, err := c.FetchSpecies(ctx, raw.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch species: %w", err)
	}

	name := extractNames(&rawNames{Names: species.Names})

	return &model.Pokemon{
		ID:        raw.ID,
		Name:      name,
		Types:     types,
		Abilities: abilities,
		Stats:     mapStats(raw),
		SpriteURL: raw.Sprites.FrontDefault,
	}, nil
}

func (c *Client) mapTypes(ctx context.Context, raw *rawPokemon) ([]model.PokemonType, error) {
	types := make([]model.PokemonType, len(raw.Types))

	for i, t := range raw.Types {
		names, err := c.FetchType(ctx, t.Type.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch type %s: %w", t.Type.Name, err)
		}

		color, ok := typeColors[t.Type.Name]
		if !ok {
			color = "#FFFFFF"
		}

		types[i] = model.PokemonType{
			Name:  extractNames(names),
			Color: color,
		}
	}

	return types, nil
}

func (c *Client) mapAbilities(ctx context.Context, raw *rawPokemon) ([]model.PokemonAbility, error) {
	abilities := make([]model.PokemonAbility, len(raw.Abilities))

	for i, a := range raw.Abilities {
		names, err := c.FetchAbility(ctx, a.Ability.URL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch ability %s: %w", a.Ability.Name, err)
		}

		abilities[i] = model.PokemonAbility{
			Name:   extractNames(names),
			Hidden: a.IsHidden,
		}
	}

	return abilities, nil
}

func mapStats(raw *rawPokemon) model.PokemonStats {
	stats := model.PokemonStats{}

	for _, s := range raw.Stats {
		switch s.Stat.Name {
		case "hp":
			stats.HP = s.BaseStat
		case "attack":
			stats.Attack = s.BaseStat
		case "defense":
			stats.Defense = s.BaseStat
		case "special-attack":
			stats.SpAttack = s.BaseStat
		case "special-defense":
			stats.SpDefense = s.BaseStat
		case "speed":
			stats.Speed = s.BaseStat
		}
	}

	return stats
}

// MapMove maps a rawMove into model.
func (c *Client) MapMove(ctx context.Context, raw *rawMove) (*model.Move, error) {
	typeNames, err := c.FetchType(ctx, raw.Type.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch move type: %w", err)
	}

	color, ok := categoryColors[raw.DamageClass.Name]
	if !ok {
		color = "#4A4A4A"
	}

	return &model.Move{
		Name: extractNames(&rawNames{Names: raw.Names}),
		Type: model.PokemonType{
			Name:  extractNames(typeNames),
			Color: typeColors[raw.Type.Name],
		},
		Category: model.MoveCategory{
			Name:  raw.DamageClass.Name,
			Color: color,
		},
		Power:    raw.Power,
		Accuracy: raw.Accuracy,
		PP:       raw.PP,
	}, nil
}

// extractNames builds a LocalizedName from a rawNames response.
func extractNames(raw *rawNames) model.LocalizedName {
	localized := model.LocalizedName{}

	for _, n := range raw.Names {
		name := strings.ToLower(n.Name)
		switch n.Language.Name {
		case "es":
			localized.ES = name
		case "en":
			localized.EN = name
		case "ja":
			localized.JA = name
		case "fr":
			localized.FR = name
		case "de":
			localized.DE = name
		case "ko":
			localized.KO = name
		case "it":
			localized.IT = name
		case "zh-Hans":
			localized.ZH = name
		}
	}

	return localized
}
