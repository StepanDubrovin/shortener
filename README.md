# shortener

URL shortening service with click tracking.

## Features

- Shorten long URLs via REST API
- Redirect to original URL by short code
- In-memory storage

## API

### Shorten URL

```bash
curl -X POST localhost:8080/shorten -d '{"url": "https://example.com/very-long-path"}'
```

```json
{"code": "1", "short": "http://localhost:8080/1"}
```

### Follow short link

Open `http://localhost:8080/1` in a browser or:

```bash
curl -i localhost:8080/1
```

Returns `302 Found` with a `Location` header pointing to the original URL.

## Run

```bash
go run ./cmd/shortener
```

Server listens on `:8080`.

## Stack

- Go 1.27, standard library only (`net/http`, `encoding/json`)
