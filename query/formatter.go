package query

import (
	"maps"
)

// Formatter replaces queries with values.
type Formatter[T any] struct {
	queries map[string]Func[T]
}

// NewFormatter ...
func NewFormatter[T any]() *Formatter[T] {
	return &Formatter[T]{make(map[string]Func[T])}
}

// Query formats all queries from the str.
func (f *Formatter[T]) Query(str string, source T) string {
	r := NewReplacer[T](source, f.queries)
	_, _ = r.WriteString(str)
	return r.Finish()
}

// WithOption adds new query.Func to the Formatter.
func (f *Formatter[T]) WithOption(key string, q Func[T]) *Formatter[T] {
	mp := maps.Clone(f.queries)
	if mp == nil {
		mp = make(map[string]Func[T])
	}
	mp[key] = q
	return &Formatter[T]{mp}
}

// Func is a function that extracts data from the source.
type Func[T any] func(source T) any
