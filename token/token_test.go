package token

import "testing"

func TestKeywords(t *testing.T) {
	tests := []struct {
		name  string
		word  string
		want  bool
	}{
		{"andalso is keyword", "andalso", true},
		{"foo is not keyword", "foo", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Keywords[tt.word]; got != tt.want {
				t.Errorf("Keywords[%q] = %v, want %v", tt.word, got, tt.want)
			}
		})
	}

	if len(Keywords) != 20 {
		t.Errorf("len(Keywords) = %d, want 20", len(Keywords))
	}
}

func TestKindString(t *testing.T) {
	tests := []struct {
		kind  Kind
		want  string
	}{
		{INT, "INT"},
		{SYMBOL, "SYMBOL"},
		{EOF, "EOF"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.kind.String(); got != tt.want {
				t.Errorf("Kind.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
