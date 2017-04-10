package session

import (
	// "net/http"
	"time"
)

type Session struct {
	values  map[string]interface{}
	expires time.Time
}

func (s *Session) Get(key string) interface{} {
	return s.values[key]
}

func (s *Session) Set(key string, value interface{}) {
	s.values[key] = value
}
