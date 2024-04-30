package lang_test

import (
	"reflect"
	"testing"

	"github.com/gosbd/gosbd"
)

func Test_Japanese(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		args args
		want []string
	}{
		{
			args: args{text: "これはペンです。それはマーカーです。"},
			want: []string{"これはペンです。", "それはマーカーです。"},
		},
		{
			args: args{text: "それは何ですか？ペンですか？"},
			want: []string{"それは何ですか？", "ペンですか？"},
		},
		{
			args: args{text: "良かったね！すごい！"},
			want: []string{"良かったね！", "すごい！"},
		},
		{
			args: args{text: "自民党税制調査会の幹部は、「引き下げ幅は３．２９％以上を目指すことになる」と指摘していて、今後、公明党と合意したうえで、３０日に決定する与党税制改正大綱に盛り込むことにしています。２％台後半を目指すとする方向で最終調整に入りました。"},
			want: []string{
				"自民党税制調査会の幹部は、「引き下げ幅は３．２９％以上を目指すことになる」と指摘していて、今後、公明党と合意したうえで、３０日に決定する与党税制改正大綱に盛り込むことにしています。",
				"２％台後半を目指すとする方向で最終調整に入りました。",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.text, func(t *testing.T) {
			sg := gosbd.NewSegmenter("ja")
			if got := sg.Segment(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Segmenter.Segment() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
