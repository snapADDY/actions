package msgflag

import (
	"reflect"
	"testing"
)

func TestExtractFlags(t *testing.T) {
	tests := map[string]struct {
		input string
		want  []string
	}{
		"no flag": {
			input: "chore: repo init",
			want:  []string{},
		},
		"simple": {
			input: "foo [skip ci]",
			want:  []string{"skip ci"},
		},
		"multiple flags": {
			input: "foo [skip ci] more text [img prod]",
			want:  []string{"skip ci", "img prod"},
		},
		"nested flags": {
			input: "foo [bar[skip ci][img prod]]",
			want:  []string{"skip ci", "img prod"},
		},
		"missmatched brackets": {
			input: "foo [skip ci[img prod][",
			want:  []string{"img prod"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := extractFlags(test.input)

			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("expected %#v, got %#v", test.want, got)
			}
		})
	}

}

func TestParseFlags(t *testing.T) {
	tests := map[string]struct {
		input []string
		want  Flags
	}{
		"empty": {
			input: []string{},
			want:  Flags{},
		},
		"invalid flag": {
			input: []string{"img test"},
			want:  Flags{},
		},
		"valid flag": {
			input: []string{"img:test-env"},
			want: Flags{
				ImgFlag{"test-env", true, false},
			},
		},
		"overwrite flag": {
			input: []string{"img:test-env", "img"},
			want: Flags{
				ImgFlag{"", true, false},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := parseFlags(test.input)

			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("expected %#v, got %#v", test.want, got)
			}
		})
	}

}

func TestParseImgFlag(t *testing.T) {
	tests := map[string]struct {
		input []string
		want  ImgFlag
	}{
		"empty": {
			input: []string{},
			want:  ImgFlag{},
		},
		"publish": {
			input: []string{"img"},
			want:  ImgFlag{"", true, false},
		},
		"overwrite env": {
			input: []string{"img", "test-env"},
			want:  ImgFlag{"test-env", true, false},
		},
		"blocked overwrite env": {
			input: []string{"img", "master"},
			want:  ImgFlag{"", true, false},
		},
		"publish prod": {
			input: []string{"img", "", "prod"},
			want:  ImgFlag{"", true, false},
		},
		"publish dev": {
			input: []string{"img", "", "dev"},
			want:  ImgFlag{"", true, true},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := parseImgFlag(test.input)

			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("expected %#v, got %#v", test.want, got)
			}
		})
	}

}

func TestParse(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Flags
	}{
		{"no tags", "chore: add gh-actions", Flags{}},
		{"publish", "feat: add stuff [img]", Flags{ImgFlag{"", true, false}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
