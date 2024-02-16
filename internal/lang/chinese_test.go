package lang_test

import (
	"reflect"
	"testing"

	"github.com/yohamta/gosbd"
)

func Test_Chinese(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		args args
		want []string
	}{
		{
			args: args{text: "安永已聯繫周怡安親屬，協助辦理簽證相關事宜，周怡安家屬1月1日晚間搭乘東方航空班機抵達上海，他們步入入境大廳時神情落寞、不發一語。周怡安來自台中，去年剛從元智大學畢業，同年9月加入安永。"},
			want: []string{"安永已聯繫周怡安親屬，協助辦理簽證相關事宜，周怡安家屬1月1日晚間搭乘東方航空班機抵達上海，他們步入入境大廳時神情落寞、不發一語。", "周怡安來自台中，去年剛從元智大學畢業，同年9月加入安永。"},
		},
		{
			args: args{text: "我们明天一起去看《摔跤吧！爸爸》好吗？好！"},
			want: []string{"我们明天一起去看《摔跤吧！爸爸》好吗？", "好！"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.text, func(t *testing.T) {
			sg := gosbd.NewSegmenter("zh")
			if got := sg.Segment(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Segmenter.Segment() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
