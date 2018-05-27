package redis_test

import (
	"encoding"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/gokv/redis"
)

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

func newStore() redis.Store {
	var addr string
	if addr = os.Getenv("REDIS_ADDR"); addr == "" {
		addr = "localhost:6379"
	}

	return redis.New(addr, os.Getenv("REDIS_PASS"))
}

func TestPing(t *testing.T) {
	t.Run("returns nil on a healthy connection", func(t *testing.T) {
		s := newStore()
		defer s.Close()

		if err := s.Ping(); err != nil {
			t.Error(err)
		}
	})

	t.Run("returns non-nil error a failed connection", func(t *testing.T) {
		s := redis.New("fakehost:8888", "")
		defer s.Close()

		if err := s.Ping(); err == nil {
			t.Errorf("expected error, found <nil>")
		}
	})
}

func TestGetSet(t *testing.T) {
	now := time.Now().UTC()
	for _, tc := range [...]struct {
		name string
		in   encoding.BinaryMarshaler
		out  encoding.BinaryUnmarshaler
	}{
		{
			"retrieves a simple string",
			&String{"hey"},
			&String{},
		},
		{
			"retrieves a time.Time",
			&now,
			&time.Time{},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := newStore()
			defer s.Close()

			if err := s.Set("somekey", tc.in); err != nil {
				t.Errorf("setting: %v", err)
			}

			ok, err := s.Get("somekey", tc.out)
			if err != nil {
				t.Errorf("getting: %v", err)
			}
			if !ok {
				t.Errorf("key not found")
			}

			if !reflect.DeepEqual(tc.in, tc.out) {
				t.Errorf("expected value %q, found %q", tc.in, tc.out)
			}
		})
	}
}

func TestSetWithTimeout(t *testing.T) {
	for _, tc := range [...]struct {
		name     string
		ttl      time.Duration
		after    time.Duration
		expected bool
	}{
		{
			"retrieves a value before timeout",
			time.Minute,
			time.Millisecond,
			true,
		},
		{
			"forgets a value after timeout",
			time.Millisecond,
			2 * time.Millisecond,
			false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s := newStore()
			defer s.Close()

			v := &String{"some value"}

			if err := s.SetWithTimeout(tc.name, v, tc.ttl); err != nil {
				t.Errorf("setting: %v", err)
			}

			time.Sleep(tc.after)

			ok, err := s.Get(tc.name, v)
			if err != nil {
				t.Errorf("getting: %v", err)
			}
			if ok != tc.expected {
				t.Errorf("value expected %v, found %v", tc.expected, ok)
			}
		})
	}
}
