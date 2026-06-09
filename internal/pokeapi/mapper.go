package pokeapi

import (
	"context"
	"fmt"

	"github.com/JoseM656/pokecli-api/internal/model"
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

// extractNames builds a LocalizedName from a rawNames response.
func extractNames(raw *rawNames) model.LocalizedName {
	localized := model.LocalizedName{}

	for _, n := range raw.Names {
		switch n.Language.Name {
		case "es":
			localized.ES = n.Name
		case "en":
			localized.EN = n.Name
		case "ja":
			localized.JA = n.Name
		case "fr":
			localized.FR = n.Name
		case "de":
			localized.DE = n.Name
		case "ko":
			localized.KO = n.Name
		case "it":
			localized.IT = n.Name
		case "zh-Hans":
			localized.ZH = n.Name
		}
	}

	return localized
}
