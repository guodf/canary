package collection

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Set[T comparable] map[T]interface{}


func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

func (s Set[T]) Remove(item T) {
	delete(s, item)
}

func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) String() string {
	var b strings.Builder
	b.WriteByte('[')

	for k := range s {
		fmt.Fprintf(&b, "%v,", k)
	}
	if s.Len() > 0 {
		return strings.TrimRight(b.String(), ",") + "]"
	}
	b.WriteByte(']')
	return b.String()
}

func (s Set[T]) MarshalJSON() ([]byte, error) {
	values := []T{}
	for key := range s {
		values = append(values, key)
	}
	return json.Marshal(&values)
}

func (s Set[T]) UnmarshalJSON(data []byte) error {
	var values []T
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	for _, value := range values {
		s.Add(value)
	}
	return nil
}
