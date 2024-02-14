package replacer

import (
	"testing"

	"github.com/gosbd/gosbd/internal/processor"
)

func TestAbbreviationReplacer_SearchForAbbreviationsInString(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				text: "I can see Mt. Fuji from here.",
			},
			want: "I can see Mt∯ Fuji from here.",
		},
		{
			args: args{
				text: "ext. 13 is my office.",
			},
			want: "ext∯ 13 is my office.",
		},
		{
			args: args{
				text: "Google Inc. is a subsidiary of Alphabet Inc.",
			},
			want: "Google Inc∯ is a subsidiary of Alphabet Inc.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.text, func(t *testing.T) {
			a := NewAbbreviationReplacer(processor.Standard())
			if got := a.SearchForAbbreviationsInString(tt.args.text); got != tt.want {
				t.Errorf("AbbreviationReplacer.SearchForAbbreviationsInString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAbbreviationReplacer_ReplaceAbbreviationAsSentenceBoundary(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				text: "It's U.S.A∯ Why not U.K.?",
			},
			want: "It's U.S.A. Why not U.K.?",
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.text, func(t *testing.T) {
			a := NewAbbreviationReplacer(processor.Standard())
			if got := a.ReplaceAbbreviationAsSentenceBoundary(tt.args.text); got != tt.want {
				t.Errorf("AbbreviationReplacer.ReplaceAbbreviationAsSentenceBoundary() = %q, want %q", got, tt.want)
			}
		})
	}
}
