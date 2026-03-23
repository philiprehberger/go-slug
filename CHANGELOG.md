# Changelog

## 0.2.0

- Add `WithStopWords` option to remove specified words from slugs
- Add `DefaultStopWords` function returning common English stop words
- Add `WithStrict` option for ASCII-only mode without transliteration
- Add `ToKebab` function for kebab-case conversion with camelCase splitting
- Add `ToSnake` function for snake_case conversion with camelCase splitting
- Add `IsSlug` function to validate slug format

## 0.1.3

- Consolidate README badges onto single line, fix CHANGELOG format

## 0.1.2

- Add Development section to README

## 0.1.0

- Initial release
- URL-safe slug generation from strings
- Unicode to ASCII transliteration
- Configurable separator and max length
- Word-boundary-aware truncation
- Uniqueness suffix helper
