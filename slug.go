// Package slug provides URL-safe slug generation from strings.
//
// It handles Unicode transliteration, configurable separators, max length
// truncation at word boundaries, and unique slug generation.
package slug

import (
	"strings"
)

// Make generates a URL-safe slug from the given string using default settings.
// It transliterates Unicode to ASCII, lowercases the result, replaces
// non-alphanumeric characters with hyphens, collapses consecutive hyphens,
// and trims hyphens from the start and end.
func Make(s string) string {
	return New().Make(s)
}

// Unique generates a unique slug by appending -2, -3, etc. if the slug
// already exists. The exists function should return true if the slug is
// already taken.
func Unique(input string, exists func(slug string) bool) string {
	return New().Unique(input, exists)
}

// Slugger generates URL-safe slugs with configurable options.
type Slugger struct {
	separator  string
	maxLen     int
	customSubs map[string]string
}

// Option configures a Slugger.
type Option func(*Slugger)

// WithSeparator sets the separator character used between words.
// The default separator is "-".
func WithSeparator(sep string) Option {
	return func(s *Slugger) {
		s.separator = sep
	}
}

// WithMaxLen sets the maximum length of the generated slug.
// The slug is truncated at the last word boundary (separator) before maxLen.
// A value of 0 means no limit.
func WithMaxLen(n int) Option {
	return func(s *Slugger) {
		s.maxLen = n
	}
}

// WithCustomSubs sets custom string substitutions that are applied before
// transliteration. Keys are matched case-sensitively.
func WithCustomSubs(subs map[string]string) Option {
	return func(s *Slugger) {
		s.customSubs = subs
	}
}

// New creates a new Slugger with the given options.
func New(opts ...Option) *Slugger {
	s := &Slugger{
		separator: "-",
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Make generates a URL-safe slug from the given string.
//
// The algorithm is:
//  1. Apply custom substitutions
//  2. Transliterate Unicode to ASCII
//  3. Lowercase
//  4. Replace any non [a-z0-9] with separator
//  5. Collapse consecutive separators
//  6. Trim separators from start/end
//  7. If maxLen is set, truncate at word boundary
func (sl *Slugger) Make(input string) string {
	s := input

	// 1. Apply custom substitutions
	for old, repl := range sl.customSubs {
		s = strings.ReplaceAll(s, old, repl)
	}

	// 2. Transliterate Unicode to ASCII
	s = transliterate(s)

	// 3. Lowercase
	s = strings.ToLower(s)

	// 4. Replace any non [a-z0-9] with separator
	var b strings.Builder
	b.Grow(len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') {
			b.WriteByte(c)
		} else {
			b.WriteString(sl.separator)
		}
	}
	s = b.String()

	// 5. Collapse consecutive separators
	for strings.Contains(s, sl.separator+sl.separator) {
		s = strings.ReplaceAll(s, sl.separator+sl.separator, sl.separator)
	}

	// 6. Trim separators from start/end
	s = strings.Trim(s, sl.separator)

	// 7. Truncate at word boundary if maxLen is set
	if sl.maxLen > 0 && len(s) > sl.maxLen {
		s = s[:sl.maxLen]
		if idx := strings.LastIndex(s, sl.separator); idx > 0 {
			s = s[:idx]
		}
		s = strings.TrimRight(s, sl.separator)
	}

	return s
}

// Unique generates a unique slug by appending -2, -3, etc. if the base slug
// already exists. The exists function should return true if a given slug is
// already taken.
func (sl *Slugger) Unique(input string, exists func(slug string) bool) string {
	base := sl.Make(input)
	if !exists(base) {
		return base
	}

	for i := 2; ; i++ {
		candidate := base + sl.separator + itoa(i)
		if !exists(candidate) {
			return candidate
		}
	}
}

// itoa converts a non-negative integer to its string representation.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
