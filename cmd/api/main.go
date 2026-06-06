package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/JoseM656/pokecli-api/internal/pokeapi"
)

func main() {
	client := pokeapi.NewClient("https://pokeapi.co/api/v2")

	raw, err := client.FetchPokemon(context.Background(), "pikachu")
	if err != nil {
		panic(err)
	}

	pokemon, err := client.Map(context.Background(), raw)
	if err != nil {
		panic(err)
	}

	out, _ := json.MarshalIndent(pokemon, "", "  ")
	fmt.Println(string(out))
}
