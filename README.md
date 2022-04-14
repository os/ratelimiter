# Rate Limiter

## Design

There are a few main concepts:

* **Identifier**: Provides a way to identify the request (e.g. IPIdentifier, TokenIdentifier)
* **Store**: Stores the counter (e.g. MemoryStore, FileStore, RedisStore)
* **Limiter**: Decides if a call is allowed (e.g. FixedWindowLimiter, SlidingWindowLog)

I've used Fixed Window Counter algorithm for this implementation as it's space efficient and easy to implement.

## Requirements

* [Docker](https://docs.docker.com/get-docker/)
* [Go](https://golang.org/dl/)

## Running

### In Docker

```bash
make run
```

### Locally

```bash
make run-local
```

## Testing

### In Docker

```bash
make test
```

### Locally

```bash
make test-local
```
