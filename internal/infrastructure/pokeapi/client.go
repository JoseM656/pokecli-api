package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/JoseM656/pokecli-api/internal/domain"
)

type pokeAPIPokemon struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Types  []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchByName fetches pokemon details from PokeAPI
func (c *Client) FetchByName(ctx context.Context, name string) (*domain.Pokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%s", c.baseURL, strings.ToLower(name))
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, domain.ErrPokemonNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var raw pokeAPIPokemon
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	types := make([]string, len(raw.Types))
	for i, t := range raw.Types {
		types[i] = t.Type.Name
	}

	return &domain.Pokemon{
		ID:     raw.ID,
		Name:   raw.Name,
		Height: raw.Height,
		Weight: raw.Weight,
		Types:  types,
	}, nil
}
