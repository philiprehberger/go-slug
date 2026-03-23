// Package slug provides URL-safe slug generation from strings.
//
// It handles Unicode transliteration, configurable separators, max length
// truncation at word boundaries, and unique slug generation.
package slug

import (
	"regexp"
	"strings"
	"unicode"
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
	stopWords  map[string]bool
	strict     bool
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

// WithStopWords sets words to be removed from the slug. Words are removed
// after splitting but before joining. Matching is case-insensitive.
func WithStopWords(words ...string) Option {
	return func(s *Slugger) {
		s.stopWords = make(map[string]bool, len(words))
		for _, w := range words {
			s.stopWords[strings.ToLower(w)] = true
		}
	}
}

// DefaultStopWords returns a built-in list of common English stop words.
func DefaultStopWords() []string {
	return []string{
		"a", "an", "the", "and", "or", "but",
		"in", "on", "at", "to", "for", "of",
		"with", "by", "is", "are", "was", "were",
	}
}

// WithStrict enables strict mode that only allows ASCII lowercase letters,
// digits, and the separator. Any non-matching characters are removed entirely
// (no transliteration).
func WithStrict() Option {
	return func(s *Slugger) {
		s.strict = true
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
//  2. Transliterate Unicode to ASCII (skipped in strict mode)
//  3. Lowercase
//  4. Replace any non [a-z0-9] with separator
//  5. Collapse consecutive separators
//  6. Trim separators from start/end
//  7. Remove stop words
//  8. If maxLen is set, truncate at word boundary
func (sl *Slugger) Make(input string) string {
	s := input

	// 1. Apply custom substitutions
	for old, repl := range sl.customSubs {
		s = strings.ReplaceAll(s, old, repl)
	}

	// 2. Transliterate Unicode to ASCII (skip in strict mode)
	if !sl.strict {
		s = transliterate(s)
	}

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

	// 7. Remove stop words
	if len(sl.stopWords) > 0 {
		parts := strings.Split(s, sl.separator)
		filtered := make([]string, 0, len(parts))
		for _, p := range parts {
			if !sl.stopWords[p] {
				filtered = append(filtered, p)
			}
		}
		s = strings.Join(filtered, sl.separator)
	}

	// 8. Truncate at word boundary if maxLen is set
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

// camelCaseRe matches boundaries for camelCase/PascalCase splitting.
var camelCaseRe = regexp.MustCompile(`([a-z0-9])([A-Z])`)

// upperSequenceRe matches sequences of uppercase letters followed by a lowercase letter.
var upperSequenceRe = regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)

// splitWords splits a string on spaces, underscores, hyphens, and camelCase boundaries.
func splitWords(s string) []string {
	// Insert separator at camelCase boundaries
	s = upperSequenceRe.ReplaceAllString(s, "${1} ${2}")
	s = camelCaseRe.ReplaceAllString(s, "${1} ${2}")

	// Split on spaces, underscores, hyphens
	parts := regexp.MustCompile(`[\s_\-]+`).Split(s, -1)
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// ToKebab converts a string to kebab-case. It splits on spaces, underscores,
// and camelCase boundaries. No Unicode transliteration is performed.
func ToKebab(s string) string {
	words := splitWords(s)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "-")
}

// ToSnake converts a string to snake_case. It splits on spaces, underscores,
// and camelCase boundaries. No Unicode transliteration is performed.
func ToSnake(s string) string {
	words := splitWords(s)
	for i, w := range words {
		words[i] = strings.ToLower(w)
	}
	return strings.Join(words, "_")
}

// IsSlug checks if a string is a valid slug. A valid slug contains only
// lowercase ASCII letters, digits, and hyphens, with no leading, trailing,
// or consecutive hyphens, and is not empty.
func IsSlug(s string) bool {
	if s == "" {
		return false
	}

	prevHyphen := false
	for i, r := range s {
		if r == '-' {
			if i == 0 || i == len(s)-1 || prevHyphen {
				return false
			}
			prevHyphen = true
			continue
		}
		prevHyphen = false
		if !unicode.IsLower(r) && !unicode.IsDigit(r) {
			return false
		}
		// Ensure only ASCII lowercase and digits
		if r > 127 {
			return false
		}
	}
	return true
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
