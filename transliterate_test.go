package slug

import (
	"testing"
)

func TestTransliterate_Accented(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"café", "cafe"},
		{"über", "uber"},
		{"naïve", "naive"},
		{"résumé", "resume"},
		{"ñoño", "nono"},
		{"straße", "strasse"},
	}

	for _, tt := range tests {
		got := transliterate(tt.input)
		if got != tt.want {
			t.Errorf("transliterate(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestTransliterate_Symbols(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"rock & roll", "rock and roll"},
		{"user@host", "userathost"},
		{"price €50", "price eur50"},
		{"£100", "gbp100"},
		{"¥200", "yen200"},
	}

	for _, tt := range tests {
		got := transliterate(tt.input)
		if got != tt.want {
			t.Errorf("transliterate(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestTransliterate_ASCII(t *testing.T) {
	input := "Hello World 123"
	got := transliterate(input)
	if got != input {
		t.Errorf("transliterate(%q) = %q, want %q", input, got, input)
	}
}

func TestTransliterate_Unknown(t *testing.T) {
	got := transliterate("hello 世界")
	want := "hello "
	if got != want {
		t.Errorf("transliterate(\"hello 世界\") = %q, want %q", got, want)
	}
}

func TestTransliterate_Polish(t *testing.T) {
	got := transliterate("łódź")
	want := "lodz"
	if got != want {
		t.Errorf("transliterate(\"łódź\") = %q, want %q", got, want)
	}
}

func TestTransliterate_Czech(t *testing.T) {
	got := transliterate("řeč")
	want := "rec"
	if got != want {
		t.Errorf("transliterate(\"řeč\") = %q, want %q", got, want)
	}
}

func TestTransliterate_Turkish(t *testing.T) {
	got := transliterate("değişim")
	want := "degisim"
	if got != want {
		t.Errorf("transliterate(\"değişim\") = %q, want %q", got, want)
	}
}
