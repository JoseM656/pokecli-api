# pokeproxy

A caching proxy written in Go that sits between a CLI client and [PokéAPI](https://pokeapi.co). It has no business logic of its own — on the first request, data is fetched from PokéAPI and stored in MongoDB. Subsequent requests are served from the local cache without hitting the upstream API.

## Stack

- **Go 1.23**
- **MongoDB 7** — cache layer
- **Docker + Docker Compose**

## Endpoints

``` http
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

## Known Limitations

- **First request must be in English** — pokemon and move names are fetched from PokéAPI using the English name as the canonical identifier. Subsequent requests can use any supported language, as the cache stores all available translations.

- **Supported languages** — translations are available for: Spanish (`es`), English (`en`), Japanese (`ja`), French (`fr`), German (`de`), Korean (`ko`), Italian (`it`), and Chinese Simplified (`zh`).

## Planned for v2.0

- **Redis cache layer** — add Redis as a fast in-memory cache in front of MongoDB to reduce latency on repeated requests.
- **ASCII sprites** — convert pokemon sprite images to ASCII art server-side and include them in the API response, so CLI clients can render them without additional processing.
- **Extended endpoints** — expose additional PokéAPI resources such as items, species, and evolutions.