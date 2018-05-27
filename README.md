# gokv/redis
[![GoDoc](https://godoc.org/github.com/gokv/redis?status.svg)](https://godoc.org/github.com/gokv/redis)
[![Build Status](https://travis-ci.org/gokv/redis.svg?branch=master)](https://travis-ci.org/gokv/redis)

A wrapper around github.com/go-redis/redis that implements the Store interface defined in [gokv/store](https://github.com/gokv/store).

### Test
Test with `go test`. An empty and *disposable* Redis instance must be running at `REDIS_ADDR` (default `localhost:6379`) with password `REDIS_PASS` (default `""`).
