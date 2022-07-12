package msgflag

import (
	"reflect"
	"testing"
)

func TestExtractFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "no flag",
			input: "chore: repo init",
			want:  []string{},
		},
		{
			name:  "simple",
			input: "foo [skip ci]",
			want:  []string{"skip ci"},
		},
		{
			name:  "multiple flags",
			input: "foo [skip ci] more text [img prod]",
			want:  []string{"skip ci", "img prod"},
		},
		{
			name:  "nested flags",
			input: "foo [bar[skip ci][img prod]]",
			want:  []string{"skip ci", "img prod"},
		},
		{
			name:  "mismatched brackets",
			input: "foo [skip ci[img prod][",
			want:  []string{"img prod"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := extractFlags(tt.input)

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("expected %#v, got %#v", tt.want, got)
			}
		})
	}
}

func TestParseFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []string
		want  Flags
	}{
		{
			name:  "empty",
			input: []string{},
			want:  Flags{},
		},
		{
			name:  "valid flag",
			input: []string{"img:test-env"},
			want: Flags{
				"img": []string{"img", "test-env"},
			},
		},
		{
			name:  "overwrite flag",
			input: []string{"img:test-env", "img"},
			want: Flags{
				"img": []string{"img"},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := parseFlags(tt.input)

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("expected %#v, got %#v", tt.want, got)
			}
		})
	}
}

func TestParse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  Flags
	}{
		{"no tags", "chore: add gh-actions", Flags{}},
		{"publish", "feat: add stuff [img]", Flags{"img": []string{"img"}}},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := Parse(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
