# pokemon-api

A REST API that acts as a caching proxy for [PokéAPI](https://pokeapi.co). On the first request, data is fetched from PokéAPI and stored in MongoDB. Subsequent requests are served from the local cache.


## Stack

- **Go 1.23**
- **MongoDB 7** — cache layer
- **Docker + Docker Compose**

## Endpoints

```
GET /pokemon/:name?lang=es
```

Returns a projected JSON with the Pokemon's name, types, abilities, base stats, and sprite URL. The `lang` parameter filters translated fields. Defaults to `en`, with English fallback if the requested language is unavailable.

## Response

```json
{
  "id": 6,
  "name":      { "es": "Charizard", "en": "Charizard" },
  "types":     [{ "es": "Fuego", "en": "Fire", "color": "#FD7D24" }],
  "abilities": [{ "es": "Llamarada", "en": "Blaze", "hidden": false }],
  "stats": {
    "hp": 78, "attack": 84, "defense": 78,
    "sp_attack": 109, "sp_defense": 85, "speed": 100
  },
  "sprite_url": "https://...",
  "cached": true,
  "lang": "es"
}
```

## Libraries

| Purpose | Library |
|---|---|
| HTTP router | `github.com/go-chi/chi/v5` |
| MongoDB driver | `go.mongodb.org/mongo-driver` |
| Environment variables | `github.com/joho/godotenv` |
