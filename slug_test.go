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
