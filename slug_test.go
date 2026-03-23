package slug

import (
	"testing"
)

func TestMake_Simple(t *testing.T) {
	got := Make("Hello, World!")
	want := "hello-world"
	if got != want {
		t.Errorf("Make(\"Hello, World!\") = %q, want %q", got, want)
	}
}

func TestMake_Accented(t *testing.T) {
	got := Make("Über Café & Naïve")
	want := "uber-cafe-and-naive"
	if got != want {
		t.Errorf("Make(\"Über Café & Naïve\") = %q, want %q", got, want)
	}
}

func TestMake_Numbers(t *testing.T) {
	got := Make("Item 42")
	want := "item-42"
	if got != want {
		t.Errorf("Make(\"Item 42\") = %q, want %q", got, want)
	}
}

func TestMake_MultipleSeparators(t *testing.T) {
	got := Make("Hello---World")
	want := "hello-world"
	if got != want {
		t.Errorf("Make(\"Hello---World\") = %q, want %q", got, want)
	}
}

func TestMake_LeadingTrailing(t *testing.T) {
	got := Make(" -Hello- ")
	want := "hello"
	if got != want {
		t.Errorf("Make(\" -Hello- \") = %q, want %q", got, want)
	}
}

func TestMake_Empty(t *testing.T) {
	got := Make("")
	want := ""
	if got != want {
		t.Errorf("Make(\"\") = %q, want %q", got, want)
	}
}

func TestMake_AlreadySlug(t *testing.T) {
	got := Make("hello-world")
	want := "hello-world"
	if got != want {
		t.Errorf("Make(\"hello-world\") = %q, want %q", got, want)
	}
}

func TestSlugger_MaxLen(t *testing.T) {
	s := New(WithMaxLen(10))
	got := s.Make("this is a long title")
	// "this-is-a-long-title" is 20 chars, truncated to 10 → "this-is-a-" → last sep at 9 → "this-is-a" wait...
	// "this-is-a-long-title"[:10] = "this-is-a-" → lastIndex("-") = 9 → "this-is-a" → trimRight "-" → "this-is-a"
	want := "this-is-a"
	if got != want {
		t.Errorf("Slugger.Make with MaxLen(10) = %q, want %q", got, want)
	}
}

func TestSlugger_MaxLen_NoBreak(t *testing.T) {
	s := New(WithMaxLen(5))
	got := s.Make("abcdefghij")
	want := "abcde"
	if got != want {
		t.Errorf("Slugger.Make with MaxLen(5) on single word = %q, want %q", got, want)
	}
}

func TestSlugger_CustomSeparator(t *testing.T) {
	s := New(WithSeparator("_"))
	got := s.Make("Hello World")
	want := "hello_world"
	if got != want {
		t.Errorf("Slugger.Make with underscore separator = %q, want %q", got, want)
	}
}

func TestSlugger_CustomSubs(t *testing.T) {
	s := New(WithCustomSubs(map[string]string{
		"C++": "cpp",
		"C#":  "csharp",
	}))
	got := s.Make("Learning C++ and C#")
	want := "learning-cpp-and-csharp"
	if got != want {
		t.Errorf("Slugger.Make with custom subs = %q, want %q", got, want)
	}
}

func TestUnique(t *testing.T) {
	existing := map[string]bool{
		"hello-world":   true,
		"hello-world-2": true,
	}

	got := Unique("Hello World", func(s string) bool {
		return existing[s]
	})
	want := "hello-world-3"
	if got != want {
		t.Errorf("Unique(\"Hello World\") = %q, want %q", got, want)
	}
}

func TestUnique_NoConflict(t *testing.T) {
	got := Unique("Hello World", func(s string) bool {
		return false
	})
	want := "hello-world"
	if got != want {
		t.Errorf("Unique with no conflict = %q, want %q", got, want)
	}
}

func TestSlugger_StopWords(t *testing.T) {
	s := New(WithStopWords("the", "a", "and"))
	got := s.Make("The quick and a fox")
	want := "quick-fox"
	if got != want {
		t.Errorf("Slugger.Make with stop words = %q, want %q", got, want)
	}
}

func TestSlugger_StopWords_AllRemoved(t *testing.T) {
	s := New(WithStopWords("hello", "world"))
	got := s.Make("Hello World")
	want := ""
	if got != want {
		t.Errorf("Slugger.Make with all stop words = %q, want %q", got, want)
	}
}

func TestSlugger_StopWords_CaseInsensitive(t *testing.T) {
	s := New(WithStopWords("THE", "And"))
	got := s.Make("The fox and The dog")
	want := "fox-dog"
	if got != want {
		t.Errorf("Slugger.Make with case-insensitive stop words = %q, want %q", got, want)
	}
}

func TestSlugger_StopWords_DefaultList(t *testing.T) {
	words := DefaultStopWords()
	if len(words) != 18 {
		t.Errorf("DefaultStopWords() length = %d, want 18", len(words))
	}
	s := New(WithStopWords(words...))
	got := s.Make("The cat is on the mat")
	want := "cat-mat"
	if got != want {
		t.Errorf("Slugger.Make with default stop words = %q, want %q", got, want)
	}
}

func TestSlugger_Strict(t *testing.T) {
	s := New(WithStrict())
	got := s.Make("Über Café")
	want := "ber-caf"
	if got != want {
		t.Errorf("Slugger.Make with strict mode = %q, want %q", got, want)
	}
}

func TestSlugger_Strict_ASCIIOnly(t *testing.T) {
	s := New(WithStrict())
	got := s.Make("Hello World 123")
	want := "hello-world-123"
	if got != want {
		t.Errorf("Slugger.Make strict with ASCII = %q, want %q", got, want)
	}
}

func TestSlugger_Strict_Symbols(t *testing.T) {
	s := New(WithStrict())
	got := s.Make("rock & roll @ night")
	want := "rock-roll-night"
	if got != want {
		t.Errorf("Slugger.Make strict with symbols = %q, want %q", got, want)
	}
}

func TestToKebab_CamelCase(t *testing.T) {
	got := ToKebab("camelCaseString")
	want := "camel-case-string"
	if got != want {
		t.Errorf("ToKebab(\"camelCaseString\") = %q, want %q", got, want)
	}
}

func TestToKebab_PascalCase(t *testing.T) {
	got := ToKebab("PascalCaseString")
	want := "pascal-case-string"
	if got != want {
		t.Errorf("ToKebab(\"PascalCaseString\") = %q, want %q", got, want)
	}
}

func TestToKebab_Spaces(t *testing.T) {
	got := ToKebab("Hello World")
	want := "hello-world"
	if got != want {
		t.Errorf("ToKebab(\"Hello World\") = %q, want %q", got, want)
	}
}

func TestToKebab_Underscores(t *testing.T) {
	got := ToKebab("snake_case_string")
	want := "snake-case-string"
	if got != want {
		t.Errorf("ToKebab(\"snake_case_string\") = %q, want %q", got, want)
	}
}

func TestToKebab_Mixed(t *testing.T) {
	got := ToKebab("XMLParser_config value")
	want := "xml-parser-config-value"
	if got != want {
		t.Errorf("ToKebab(\"XMLParser_config value\") = %q, want %q", got, want)
	}
}

func TestToKebab_Empty(t *testing.T) {
	got := ToKebab("")
	want := ""
	if got != want {
		t.Errorf("ToKebab(\"\") = %q, want %q", got, want)
	}
}

func TestToSnake_CamelCase(t *testing.T) {
	got := ToSnake("camelCaseString")
	want := "camel_case_string"
	if got != want {
		t.Errorf("ToSnake(\"camelCaseString\") = %q, want %q", got, want)
	}
}

func TestToSnake_PascalCase(t *testing.T) {
	got := ToSnake("PascalCaseString")
	want := "pascal_case_string"
	if got != want {
		t.Errorf("ToSnake(\"PascalCaseString\") = %q, want %q", got, want)
	}
}

func TestToSnake_Spaces(t *testing.T) {
	got := ToSnake("Hello World")
	want := "hello_world"
	if got != want {
		t.Errorf("ToSnake(\"Hello World\") = %q, want %q", got, want)
	}
}

func TestToSnake_Hyphens(t *testing.T) {
	got := ToSnake("kebab-case-string")
	want := "kebab_case_string"
	if got != want {
		t.Errorf("ToSnake(\"kebab-case-string\") = %q, want %q", got, want)
	}
}

func TestToSnake_Mixed(t *testing.T) {
	got := ToSnake("XMLParser_config value")
	want := "xml_parser_config_value"
	if got != want {
		t.Errorf("ToSnake(\"XMLParser_config value\") = %q, want %q", got, want)
	}
}

func TestIsSlug_Valid(t *testing.T) {
	tests := []string{
		"hello-world",
		"hello",
		"a",
		"hello-world-123",
		"123",
		"a-b-c",
	}
	for _, s := range tests {
		if !IsSlug(s) {
			t.Errorf("IsSlug(%q) = false, want true", s)
		}
	}
}

func TestIsSlug_Invalid(t *testing.T) {
	tests := []struct {
		input string
		desc  string
	}{
		{"", "empty string"},
		{"-hello", "leading hyphen"},
		{"hello-", "trailing hyphen"},
		{"hello--world", "consecutive hyphens"},
		{"Hello", "uppercase letters"},
		{"hello world", "spaces"},
		{"hello_world", "underscores"},
		{"héllo", "non-ASCII"},
		{"-", "single hyphen"},
		{"hello--", "trailing consecutive hyphens"},
	}
	for _, tt := range tests {
		if IsSlug(tt.input) {
			t.Errorf("IsSlug(%q) = true, want false (%s)", tt.input, tt.desc)
		}
	}
}
