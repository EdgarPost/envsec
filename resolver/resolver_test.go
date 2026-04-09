package resolver

import "testing"

func TestNormalizeRemoteSSH(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"git@github.com:Org/repo.git", "github.com/Org/repo"},
		{"git@github.com:EdgarPost/envsec.git", "github.com/EdgarPost/envsec"},
		{"https://github.com/Org/repo.git", "github.com/Org/repo"},
		{"https://github.com/Org/repo", "github.com/Org/repo"},
		{"ssh://git@github.com/Org/repo.git", "github.com/Org/repo"},
	}

	for _, tt := range tests {
		got := normalizeRemote(tt.input)
		if got != tt.want {
			t.Errorf("normalizeRemote(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
