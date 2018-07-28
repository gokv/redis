package redis

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
)

var ErrDuplicateKey = errors.New("duplicate key")

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
// If no match is found, returns (false, nil).
func (s Store) Get(k string, v json.Unmarshaler) (bool, error) {
	res, err := s.c.Get(k).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, v.UnmarshalJSON([]byte(res))
}

// Set assigns the given value to the given key, possibly overwriting.
func (s Store) Set(k string, v json.Marshaler) error {
	value, err := v.MarshalJSON()
	if err != nil {
		return err
	}
	return s.c.Set(k, value, 0).Err()
}

// Add persists a new object with a new UUIDv4 key.
// Err is non-nil in case of failure.
func (s Store) Add(v json.Marshaler) (string, error) {
	value, err := v.MarshalJSON()
	if err != nil {
		return "", err
	}

	k := uuid.New().String()

	ok, err := s.c.SetNX(k, value, 0).Result()
	if err != nil {
		return "", err
	}

	if !ok {
		return "", ErrDuplicateKey
	}

	return k, nil
}

// SetWithDeadline assigns the given value to the given key, possibly
// overwriting.
// The assigned key will clear after deadline.
func (s Store) SetWithDeadline(k string, v json.Marshaler, deadline time.Time) error {
	return s.SetWithTimeout(k, v, deadline.Sub(time.Now()))
}

// SetWithTimeout assigns the given value to the given key, possibly
// overwriting.
// The assigned key will clear after timeout.
func (s Store) SetWithTimeout(k string, v json.Marshaler, timeout time.Duration) error {
	value, err := v.MarshalJSON()
	if err != nil {
		return err
	}
	return s.c.Set(k, value, timeout).Err()
}
