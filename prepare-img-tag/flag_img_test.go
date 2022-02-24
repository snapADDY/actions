package main

import (
	"reflect"
	"testing"
)

func TestParseImgFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []string
		want  ImgFlag
	}{
		{
			name:  "empty",
			input: []string{},
			want:  ImgFlag{},
		},
		{
			name:  "publish",
			input: []string{"img"},
			want:  ImgFlag{"", true, false},
		},
		{
			name:  "overwrite env",
			input: []string{"img", "test-env"},
			want:  ImgFlag{"test-env", true, false},
		},
		{
			name:  "blocked overwrite env",
			input: []string{"img", "master"},
			want:  ImgFlag{"", true, false},
		},
		{
			name:  "publish prod",
			input: []string{"img", "", "prod"},
			want:  ImgFlag{"", true, false},
		},
		{
			name:  "publish dev",
			input: []string{"img", "", "dev"},
			want:  ImgFlag{"", true, true},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := parseImgFlag(tt.input)

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("expected %#v, got %#v", tt.want, got)
			}
		})
	}
}
