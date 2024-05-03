package replacer

import "testing"

func TestBetweenPunctuation_Replace(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				text: "She turned to him, 'This is great.' she said.",
			},
			want: "She turned to him, 'This is great∯' she said.",
		},
		{
			args: args{
				text: `Hello. "Hello." Hi. I'm good "13433 434o3f"`,
			},
			want: `Hello. "Hello∯" Hi. I'm good "13433 434o3f"`,
		},
		{
			args: args{
				text: "[Square brackets are fine, too.]",
			},
			want: `[Square brackets are fine, too∯]`,
		},
		{
			args: args{
				text: `Le parti fisiche di un computer (ad es. RAM, CPU, tastiera, mouse, etc.) sono definiti HW.`,
			},
			want: `Le parti fisiche di un computer (ad es∯ RAM, CPU, tastiera, mouse, etc∯) sono definiti HW.`,
		},
		{
			args: args{
				text: `Ce modèle permet d’afficher le texte « LL.AA.II.RR. » pour l’abréviation de « Leurs Altesses impériales et royales » avec son infobulle.`,
			},
			want: `Ce modèle permet d’afficher le texte « LL∯AA∯II∯RR∯ » pour l’abréviation de « Leurs Altesses impériales et royales » avec son infobulle.`,
		},
		{
			args: args{
				text: `Mix it, put it in the oven, and -- voila! -- you have cake. Some can be -- if I may say so? -- a bit questionable.`,
			},
			want: `Mix it, put it in the oven, and -- voila&ᓴ& -- you have cake. Some can be -- if I may say so&ᓷ& -- a bit questionable.`,
		},
		{
			args: args{
				text: `She turned to him, “This is great.” She held the book out to show him.`,
			},
			want: `She turned to him, “This is great∯” She held the book out to show him.`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.text, func(t *testing.T) {
			b := BetweenPunctuation{
				punctuationReplacer: NewPunctuationReplacer(),
			}
			if got := b.Replace(tt.args.text); got != tt.want {
				t.Errorf("BetweenPunctuation.Replace() = %q, want %q", got, tt.want)
			}
		})
	}
}
