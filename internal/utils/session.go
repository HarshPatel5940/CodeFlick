package utils

import "github.com/gorilla/sessions"

func GetSessionValue[T any](s *sessions.Session, key string, defaultValue T) T {
	if value, ok := s.Values[key]; ok {
		if typedValue, ok := value.(T); ok {
			return typedValue
		}
	}
	return defaultValue
}
