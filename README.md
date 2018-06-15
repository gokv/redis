# gokv/redis
[![GoDoc](https://godoc.org/github.com/gokv/redis?status.svg)](https://godoc.org/github.com/gokv/redis)
[![Build Status](https://travis-ci.org/gokv/redis.svg?branch=master)](https://travis-ci.org/gokv/redis)

A wrapper around github.com/go-redis/redis that implements the Store interface defined in [gokv/store](https://github.com/gokv/store).

## Intro

This is experimental software. Making it work is not the primary goal.

The idea behind `github.com/gokv/store` is that sometimes, when a developer needs a persistence layer, she only needs a simple one.

## Use

Initialise calling `New` with the address and the (optional) password to Redis.

```Go
type String struct {
	s string
}

func (s *String) UnmarshalBinary(data []byte) error {
	s.s = string(data)
	return nil
}

func (s String) MarshalBinary() ([]byte, error) {
	return []byte(s.s), nil
}

func main() {

	// New instantiates a "github.com/go-redis/redis" connection
	s := redis.New(os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASS"))
	defer s.Close()

	// Call Ping to check readiness
	if err := s.Ping(); err!=nil {
		panic(err)
	}

	given := String{"hello world"}

	if err := s.Add("key", given); err != nil {
		panic(err)
	}

	var found String
	ok, err := s.Get("key", &found)

	if err != nil {
		panic(fmt.Errorf("failure: %s", err))
	}

	if !ok {
		panic(errors.New("key not found"))
	}

	// given == found
}
```


## Test
Test with `go test`. An empty and *disposable* Redis instance must be running at `REDIS_ADDR` (default `localhost:6379`) with password `REDIS_PASS` (empty by default).
