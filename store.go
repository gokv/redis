package redis

import (
	"encoding"
	"time"

	"github.com/go-redis/redis"
)

type Store struct {
	c *redis.Client
}

func (s Store) Close() error {
	return s.c.Close()
}

func New(address, password string) Store {
	return Store{
		c: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
		}),
	}
}

func (s Store) Ping() (err error) {
	return s.c.Ping().Err()
}

// Get returns the value corresponding the key, and a nil error.
// If no match is found, returns (nil, nil).
func (s Store) Get(key string, v encoding.BinaryUnmarshaler) (bool, error) {
	res, err := s.c.Get(key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, v.UnmarshalBinary([]byte(res))
}

// Set assigns the given value to the given key, possibly overwriting.
func (s Store) Set(key string, v encoding.BinaryMarshaler) error {
	value, err := v.MarshalBinary()
	if err != nil {
		return err
	}
	return s.c.Set(key, value, 0).Err()
}

// SetWithDeadline assigns the given value to the given key, possibly
// overwriting.
// The assigned key will clear after deadline.
func (s Store) SetWithDeadline(key string, v encoding.BinaryMarshaler, deadline time.Time) error {
	return s.SetWithTimeout(key, v, deadline.Sub(time.Now()))
}

// SetWithTimeout assigns the given value to the given key, possibly
// overwriting.
// The assigned key will clear after timeout.
func (s Store) SetWithTimeout(key string, v encoding.BinaryMarshaler, timeout time.Duration) error {
	value, err := v.MarshalBinary()
	if err != nil {
		return err
	}
	return s.c.Set(key, value, timeout).Err()
}
