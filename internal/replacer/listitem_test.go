package replacer

import (
	"fmt"
	"reflect"
	"testing"
)

func TestListItemReplacer_replaceAlphabetListParens(t *testing.T) {
	type args struct {
		text     string
		alphabet string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				text:     "a) ffegnog (b) fgegkl c)",
				alphabet: "a",
			},
			want: "\ra) ffegnog (b) fgegkl c)",
		},
		{
			args: args{
				text:     "a) ffegnog (b) fgegkl c)",
				alphabet: "b",
			},
			want: "a) ffegnog \r&✂&b) fgegkl c)",
		},
		{
			args: args{
				text:     "a) ffegnog (b) fgegkl c)",
				alphabet: "c",
			},
			want: "a) ffegnog (b) fgegkl \rc)",
		},
		{
			args: args{
				text:     "\ra) ffegnog \r&✂&b) fgegkl c)",
				alphabet: "c",
			},
			want: "\ra) ffegnog \r&✂&b) fgegkl \rc)",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s-%s", tt.args.text, tt.args.alphabet), func(t *testing.T) {
			l := ListItemReplacer{}
			if got := l.replaceAlphabetListParens(tt.args.text, tt.args.alphabet); got != tt.want {
				t.Errorf(`ListItemReplacer.replaceAlphabetListParens() = %q, want %q`, got, tt.want)
			}
		})
	}
}

func TestListItemReplacer_replaceAlphabetListPeriod(t *testing.T) {
	type args struct {
		text     string
		alphabet string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				text:     "a. ffegnog b. fgegkl c.",
				alphabet: "a",
			},
			want: "\ra∯ ffegnog b. fgegkl c.",
		},
		{
			args: args{
				text:     "\ra∯ ffegnog b. fgegkl c.",
				alphabet: "b",
			},
			want: "\ra∯ ffegnog \rb∯ fgegkl c.",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s-%s", tt.args.text, tt.args.alphabet), func(t *testing.T) {
			l := ListItemReplacer{}
			if got := l.replaceAlphabetListPeriod(tt.args.text, tt.args.alphabet); got != tt.want {
				t.Errorf("ListItemReplacer.replaceAlphabetListPeriod() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestListItemReplacer_extractAlphabeticalListLettersWithParens(t *testing.T) {
	type args struct {
		text         string
		romanNumeral bool
	}
	tests := []struct {
		l     ListItemReplacer
		args  args
		want  []string
		want1 map[string]alphabet
	}{
		{
			args: args{
				text: "a) ffegnog (b) fgegkl c)",
			},
			want: []string{"a", "b", "c"},
		},
		{
			args: args{
				text: "1) a) ffegnog (b) fgegkl c) 2)",
			},
			want: []string{"a", "b", "c"},
		},
		{
			args: args{
				text:         "i) ffegnog (ii) fgegkl iii)",
				romanNumeral: true,
			},
			want: []string{"i", "ii", "iii"},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s-%t", tt.args.text, tt.args.romanNumeral), func(t *testing.T) {
			l := ListItemReplacer{}
			got, _ := l.extractAlphabeticalListLettersWithParens(tt.args.text, tt.args.romanNumeral)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListItemReplacer.extractAlphabeticalListLettersWithParens() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListItemReplacer_extractAlphabeticalListLettersWithPeriods(t *testing.T) {
	type args struct {
		text         string
		romanNumeral bool
	}
	tests := []struct {
		l     ListItemReplacer
		args  args
		want  []string
		want1 map[string]alphabet
	}{
		{
			args: args{
				text: "a. ffegnog b. fgegkl c.",
			},
			want: []string{"a", "b", "c"},
		},
		{
			args: args{
				text: "1. a. ffegnog b. fgegkl c. 2.",
			},
			want: []string{"a", "b", "c"},
		},
		{
			args: args{
				text:         "i. ffegnog ii. fgegkl iii.",
				romanNumeral: true,
			},
			want: []string{"i", "ii", "iii"},
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s-%t", tt.args.text, tt.args.romanNumeral), func(t *testing.T) {
			l := ListItemReplacer{}
			got, _ := l.extractAlphabeticalListLettersWithPeriods(tt.args.text, tt.args.romanNumeral)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListItemReplacer.extractAlphabeticalListLettersWithPeriods() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListItemReplacer_iterateAlphabetArray(t *testing.T) {
	type args struct {
		text         string
		romanNumeral bool
		parens       bool
	}
	tests := []struct {
		args args
		want string
	}{
		{
			// a) and b) should be replaced since c and x are not adjacent.
			args: args{
				text:   "a) ffegnog (b) fgegkl c) ekej x) xkdj",
				parens: true,
			},
			want: "\ra) ffegnog \r&✂&b) fgegkl c) ekej x) xkdj",
		},
		{
			// only c) should be replaced since x and b are not adjacent.
			args: args{
				text:   "x) ffegnog (b) fgegkl c) ekej",
				parens: true,
			},
			want: "x) ffegnog (b) fgegkl \rc) ekej",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s-%t-%t", tt.args.text, tt.args.romanNumeral, tt.args.parens), func(t *testing.T) {
			l := ListItemReplacer{}
			if got := l.iterateAlphabetArray(tt.args.text, tt.args.romanNumeral, tt.args.parens); got != tt.want {
				t.Errorf("ListItemReplacer.iterateAlphabetArray() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestListItemReplacer_replacePeriodsInNumberedList(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				text: "1. ffegnog 2. fgegkl 3. ekej",
			},
			want: "1♨ ffegnog 2♨ fgegkl 3♨ ekej",
		},
		{
			args: args{
				text: "1.) ffegnog 2.) fgegkl 3.) ekej",
			},
			want: "1♨) ffegnog 2♨) fgegkl 3♨) ekej",
		},
		{
			args: args{
				text: "5. kajf. 1. ffegnog 2. fgegkl 3. ekej",
			},
			want: "5. kajf. 1♨ ffegnog 2♨ fgegkl 3♨ ekej",
		},
		{
			args: args{
				text: "-1.) ffegnog -2.) fgegkl -3.) ekej",
			},
			want: "-1♨) ffegnog -2♨) fgegkl -3♨) ekej",
		},
		{
			args: args{
				text: "(1.) ffegnog (2.) fgegkl (3.) ekej",
			},
			want: "(1.) ffegnog (2.) fgegkl (3.) ekej",
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.text, func(t *testing.T) {
			l := ListItemReplacer{}
			if got := l.replacePeriodsInNumberedList(tt.args.text); got != tt.want {
				t.Errorf("ListItemReplacer.replacePeriodsInNumberedList() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestListItemReplacer_addLineBreaksForNumberedListWithPeriods(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				text: "1♨ ffegnog 2♨ fgegkl 3♨ ekej",
			},
			want: "1♨ ffegnog\r2♨ fgegkl\r3♨ ekej",
		},
		{
			args: args{
				text: "1♨ 日本語 2♨ ドイツ語 3♨ ekej",
			},
			want: "1♨ 日本語\r2♨ ドイツ語\r3♨ ekej",
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.text, func(t *testing.T) {
			l := ListItemReplacer{}
			if got := l.addLineBreaksForNumberedListWithPeriods(tt.args.text); got != tt.want {
				t.Errorf("ListItemReplacer.addLineBreaksForNumberedListWithPeriods() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestListItemReplacer_formatNumberedListWithPeriods(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				text: "1. ffegnog 2. fgegkl 3. ekej",
			},
			want: "1∯ ffegnog\r2∯ fgegkl\r3∯ ekej",
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.text, func(t *testing.T) {
			l := ListItemReplacer{}
			if got := l.formatNumberedListWithPeriods(tt.args.text); got != tt.want {
				t.Errorf("ListItemReplacer.formatNumberedListWithPeriods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListItemReplacer_replaceParensInNumberedList(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				text: "1) ffegnog 2) fgegkl 3) ekej",
			},
			want: "1☝) ffegnog 2☝) fgegkl 3☝) ekej",
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.text, func(t *testing.T) {
			l := ListItemReplacer{}
			if got := l.replaceParensInNumberedList(tt.args.text); got != tt.want {
				t.Errorf("ListItemReplacer.replaceParensInNumberedList() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestListItemReplacer_addLineBreaksForNumberedListWithParens(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				text: "1☝) ffegnog 2☝) fgegkl 3☝) ekej",
			},
			want: "1☝) ffegnog\r2☝) fgegkl\r3☝) ekej",
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.text, func(t *testing.T) {
			l := ListItemReplacer{}
			if got := l.addLineBreaksForNumberedListWithParens(tt.args.text); got != tt.want {
				t.Errorf("ListItemReplacer.addLineBreaksForNumberedListWithParens() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestListItemReplacer_formatNumberedListWithParens(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{
				text: "1) ffegnog 2) fgegkl 3) ekej",
			},
			want: "1) ffegnog\r2) fgegkl\r3) ekej",
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.text, func(t *testing.T) {
			l := ListItemReplacer{}
			if got := l.formatNumberedListWithParens(tt.args.text); got != tt.want {
				t.Errorf("ListItemReplacer.formatNumberedListWithParens() = %v, want %v", got, tt.want)
			}
		})
	}
}
