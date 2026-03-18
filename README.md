# go-slug

[![CI](https://github.com/philiprehberger/go-slug/actions/workflows/ci.yml/badge.svg)](https://github.com/philiprehberger/go-slug/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/philiprehberger/go-slug.svg)](https://pkg.go.dev/github.com/philiprehberger/go-slug)
[![License](https://img.shields.io/github/license/philiprehberger/go-slug)](LICENSE)

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

## Development

```bash
go test ./...
go vet ./...
```

## License

MIT
