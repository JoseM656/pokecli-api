package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Raw structs — PokeAPI response model.

type rawPokemon struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
	Types  []struct {
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Abilities []struct {
		IsHidden bool `json:"is_hidden"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Sprites struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites"`
}

type rawNames struct {
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
		} `json:"language"`
	} `json:"names"`
}

// Client handles communication with PokeAPI.
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

func (c *Client) FetchPokemon(ctx context.Context, name string) (*rawPokemon, error) {
	url := fmt.Sprintf("%s/pokemon/%s", c.baseURL, strings.ToLower(name))
	return fetch[rawPokemon](ctx, c.httpClient, url)
}

func (c *Client) FetchType(ctx context.Context, url string) (*rawNames, error) {
	return fetch[rawNames](ctx, c.httpClient, url)
}

func (c *Client) FetchAbility(ctx context.Context, url string) (*rawNames, error) {
	return fetch[rawNames](ctx, c.httpClient, url)
}

// fetch is a generic helper to avoid repeating the request/decode logic.
func fetch[T any](ctx context.Context, client *http.Client, url string) (*T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("not found: %s", url)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, url)
	}

	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
