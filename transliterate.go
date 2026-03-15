package slug

import (
	"strings"
	"unicode/utf8"
)

// transliterations maps common Unicode characters to their ASCII equivalents.
var transliterations = map[rune]string{
	// Latin accented lowercase
	'à': "a", 'á': "a", 'â': "a", 'ã': "a", 'ä': "a", 'å': "a",
	'æ': "ae", 'ç': "c",
	'è': "e", 'é': "e", 'ê': "e", 'ë': "e",
	'ì': "i", 'í': "i", 'î': "i", 'ï': "i",
	'ñ': "n",
	'ò': "o", 'ó': "o", 'ô': "o", 'õ': "o", 'ö': "o", 'ø': "o",
	'ù': "u", 'ú': "u", 'û': "u", 'ü': "u",
	'ý': "y", 'ÿ': "y",
	'ß': "ss",

	// Latin accented uppercase
	'À': "A", 'Á': "A", 'Â': "A", 'Ã': "A", 'Ä': "A", 'Å': "A",
	'Æ': "AE", 'Ç': "C",
	'È': "E", 'É': "E", 'Ê': "E", 'Ë': "E",
	'Ì': "I", 'Í': "I", 'Î': "I", 'Ï': "I",
	'Ñ': "N",
	'Ò': "O", 'Ó': "O", 'Ô': "O", 'Õ': "O", 'Ö': "O", 'Ø': "O",
	'Ù': "U", 'Ú': "U", 'Û': "U", 'Ü': "U",
	'Ý': "Y", 'Ÿ': "Y",

	// Common symbols
	'&': "and", '@': "at", '©': "c", '®': "r",
	'€': "eur", '£': "gbp", '¥': "yen",

	// Polish
	'ą': "a", 'ć': "c", 'ę': "e", 'ł': "l", 'ń': "n",
	'ś': "s", 'ź': "z", 'ż': "z",
	'Ą': "A", 'Ć': "C", 'Ę': "E", 'Ł': "L", 'Ń': "N",
	'Ś': "S", 'Ź': "Z", 'Ż': "Z",

	// Czech/Slovak
	'ř': "r", 'š': "s", 'ž': "z", 'č': "c", 'ď': "d", 'ť': "t", 'ň': "n", 'ů': "u",
	'Ř': "R", 'Š': "S", 'Ž': "Z", 'Č': "C", 'Ď': "D", 'Ť': "T", 'Ň': "N", 'Ů': "U",

	// Turkish
	'ğ': "g", 'ı': "i", 'İ': "i", 'ş': "s",
	'Ğ': "G", 'Ş': "S",
}

// transliterate replaces known Unicode characters with their ASCII equivalents
// and drops any remaining non-ASCII characters.
func transliterate(s string) string {
	var b strings.Builder
	b.Grow(len(s))

	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if r == utf8.RuneError && size == 1 {
			i++
			continue
		}

		if repl, ok := transliterations[r]; ok {
			b.WriteString(repl)
		} else if r < 128 {
			b.WriteRune(r)
		}
		// Unknown non-ASCII characters are dropped

		i += size
	}

	return b.String()
}
