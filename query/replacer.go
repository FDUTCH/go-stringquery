package query

import (
	"fmt"
	"strings"
)

// Replacer is a string builder with formatting capabilities.
type Replacer[T any] struct {
	builder, keyBuilder strings.Builder
	state               int

	source  T
	queries map[string]Func[T]
}

// grow grows builder if necessary.
func (r *Replacer[T]) grow(size int) {
	toAlloc := size - (r.builder.Cap() - r.builder.Len())
	if toAlloc > 0 {
		r.builder.Grow(toAlloc)
	}
}

// WriteString ...
func (r *Replacer[T]) WriteString(str string) (int, error) {
	r.grow(len(str))
	for _, char := range str {
		r.WriteRune(char)
	}
	return len(str), nil
}

// Write ...
func (r *Replacer[T]) Write(p []byte) (n int, err error) {
	r.grow(len(p))
	for _, b := range p {
		r.WriteRune(rune(b))
	}
	return len(p), nil
}

// WriteRune writes rune.
func (r *Replacer[T]) WriteRune(char rune) {
	switch r.state {
	case stateInit:
		r.handleInit(char)
	case stateAwaitingScope:
		r.handleAwaitScope(char)
	case stateScopeOpened:
		r.handleScopeOpened(char)
	}
}

// SetSource sets data source and queries fields.
func (r *Replacer[T]) SetSource(source T, queries map[string]Func[T]) {
	r.source = source
	r.queries = queries
}

// handleInit handles the initial state.
func (r *Replacer[T]) handleInit(char rune) {
	if char == '$' {
		r.state = stateAwaitingScope
		return
	}
	r.builder.WriteRune(char)
}

// handleAwaitScope handles the scope await state.
func (r *Replacer[T]) handleAwaitScope(char rune) {
	if char == '{' {
		r.state = stateScopeOpened
		return
	}
	r.state = stateInit
	r.builder.WriteRune('$')
	r.builder.WriteRune(char)
}

// handleScopeOpened handles the scope opened state.
func (r *Replacer[T]) handleScopeOpened(char rune) {
	if char == '}' {
		r.state = stateInit
		if fn, ok := r.queries[r.keyBuilder.String()]; ok {
			_, _ = fmt.Fprint(&r.builder, fn(r.source))
		} else {
			r.builder.WriteRune('$')
			r.builder.WriteRune('{')
			r.builder.WriteString(r.keyBuilder.String())
			r.builder.WriteRune('}')
		}
		r.keyBuilder.Reset()
		return
	}
	r.keyBuilder.WriteRune(char)
}

// Finish returns formatted string and resets the state of the Replacer.
// After use, source and queries will also be reset.
func (r *Replacer[T]) Finish() string {
	var none T
	r.source = none
	r.queries = nil
	defer r.builder.Reset()

	if r.state != stateInit {
		r.builder.WriteRune('$')
		r.builder.WriteRune('{')
		r.builder.WriteString(r.keyBuilder.String())
		r.keyBuilder.Reset()
		r.state = stateInit
	}

	return r.builder.String()
}

// NewReplacer creates new Replacer.
func NewReplacer[T any](source T, queries map[string]Func[T]) *Replacer[T] {
	return &Replacer[T]{source: source, queries: queries}
}

const (
	stateInit = iota
	stateAwaitingScope
	stateScopeOpened
)
