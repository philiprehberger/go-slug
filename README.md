# go-slug

[![CI](https://github.com/philiprehberger/go-slug/actions/workflows/ci.yml/badge.svg)](https://github.com/philiprehberger/go-slug/actions/workflows/ci.yml) [![Go Reference](https://pkg.go.dev/badge/github.com/philiprehberger/go-slug.svg)](https://pkg.go.dev/github.com/philiprehberger/go-slug) [![License](https://img.shields.io/github/license/philiprehberger/go-slug)](LICENSE)

URL-safe slug generator for Go. Handles Unicode, configurable, zero dependencies

## Installation

```bash
go get github.com/philiprehberger/go-slug
```

## Usage

### Basic

```go
import "github.com/philiprehberger/go-slug"

slug.Make("Hello, World!")        // "hello-world"
slug.Make("Über Café & Naïve")   // "uber-cafe-and-naive"
slug.Make("Item 42")              // "item-42"
```

### Options

```go
// Custom separator
s := slug.New(slug.WithSeparator("_"))
s.Make("Hello World") // "hello_world"

// Max length (truncates at word boundary)
s = slug.New(slug.WithMaxLen(10))
s.Make("this is a long title") // "this-is-a"

// Custom substitutions (applied before transliteration)
s = slug.New(slug.WithCustomSubs(map[string]string{
    "C++": "cpp",
    "C#":  "csharp",
}))
s.Make("Learning C++ and C#") // "learning-cpp-and-csharp"
```

### Stop Words

```go
// Remove specific words
s := slug.New(slug.WithStopWords("the", "a", "and"))
s.Make("The quick and a fox") // "quick-fox"

// Use built-in stop word list
s = slug.New(slug.WithStopWords(slug.DefaultStopWords()...))
s.Make("The cat is on the mat") // "cat-mat"
```

### Strict Mode

```go
// Only allow ASCII letters, digits, and separator (no transliteration)
s := slug.New(slug.WithStrict())
s.Make("Über Café")          // "ber-caf"
s.Make("Hello World 123")    // "hello-world-123"
s.Make("rock & roll @ night") // "rock-roll-night"
```

### Case Conversion

```go
// Kebab-case (splits on spaces, underscores, camelCase boundaries)
slug.ToKebab("camelCaseString")  // "camel-case-string"
slug.ToKebab("PascalCaseString") // "pascal-case-string"
slug.ToKebab("snake_case_string") // "snake-case-string"

// Snake_case (same splitting, underscore separator)
slug.ToSnake("camelCaseString")  // "camel_case_string"
slug.ToSnake("Hello World")     // "hello_world"
slug.ToSnake("kebab-case-string") // "kebab_case_string"
```

### Validation

```go
slug.IsSlug("hello-world")   // true
slug.IsSlug("hello")         // true
slug.IsSlug("Hello-World")   // false (uppercase)
slug.IsSlug("hello--world")  // false (consecutive hyphens)
slug.IsSlug("-hello")        // false (leading hyphen)
slug.IsSlug("")              // false (empty)
```

### Unique Slugs

```go
existing := map[string]bool{
    "hello-world":   true,
    "hello-world-2": true,
}

result := slug.Unique("Hello World", func(s string) bool {
    return existing[s]
})
// result: "hello-world-3"
```

## API

| Function / Type | Description |
|---|---|
| `Make(s string) string` | Generate a slug with default settings |
| `Unique(input string, exists func(string) bool) string` | Generate a unique slug with default settings |
| `New(opts ...Option) *Slugger` | Create a configured slugger |
| `(*Slugger) Make(input string) string` | Generate a slug with configured options |
| `(*Slugger) Unique(input string, exists func(string) bool) string` | Generate a unique slug with configured options |
| `WithSeparator(sep string) Option` | Set the word separator (default: "-") |
| `WithMaxLen(n int) Option` | Set max slug length with word-boundary truncation |
| `WithCustomSubs(subs map[string]string) Option` | Set custom string substitutions |
| `WithStopWords(words ...string) Option` | Remove specified words from the slug |
| `DefaultStopWords() []string` | Returns built-in list of common English stop words |
| `WithStrict() Option` | Strict mode: ASCII only, no transliteration |
| `ToKebab(s string) string` | Convert string to kebab-case |
| `ToSnake(s string) string` | Convert string to snake_case |
| `IsSlug(s string) bool` | Check if string is a valid slug |

## Development

```bash
go test ./...
go vet ./...
```

## License

MIT
