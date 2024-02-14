package processor

import (
	"testing"
)

func Test_processor_replaceContinuousPunctuation(t *testing.T) {
	type fields struct {
		cfg *Config
	}
	type args struct {
		text string
	}
	tests := []struct {
		fields fields
		args   args
		want   string
	}{
		{
			fields: fields{cfg: Standard()},
			args: args{
				text: "Hello!!! How are you!!?",
			},
			want: "Hello&ᓴ&&ᓴ&&ᓴ& How are you&ᓴ&&ᓴ&&ᓷ&",
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.text, func(t *testing.T) {
			p := &Processor{
				cfg: tt.fields.cfg,
			}
			if got := p.replaceContinuousPunctuation(tt.args.text); got != tt.want {
				t.Errorf("processor.replaceContinuousPunctuation() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_processor_replaceNumbers(t *testing.T) {
	type fields struct {
		cfg *Config
	}
	type args struct {
		text string
	}
	tests := []struct {
		fields fields
		args   args
		want   string
	}{
		{
			fields: fields{cfg: Standard()},
			args: args{
				text: "I have 1.5 apples.",
			},
			want: "I have 1∯5 apples.",
		},
		{
			fields: fields{cfg: Standard()},
			args: args{
				text: "I have 1.x apples.",
			},
			want: "I have 1∯x apples.",
		},
		{
			fields: fields{cfg: Standard()},
			args: args{
				text: "\r1. I have 1.5 apples.\r2. I have 2.5 apples.\r3. I have 3.5 apples.",
			},
			want: "\r1∯ I have 1∯5 apples.\r2∯ I have 2∯5 apples.\r3∯ I have 3∯5 apples.",
		},
		{
			fields: fields{cfg: Standard()},
			args: args{
				text: "I have 10.5 apples.",
			},
			want: "I have 10∯5 apples.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.text, func(t *testing.T) {
			p := &Processor{cfg: tt.fields.cfg}
			if got := p.replaceNumbers(tt.args.text); got != tt.want {
				t.Errorf("processor.replaceNumbers() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_processor_replacePeriodsBeforeNumericReferences(t *testing.T) {
	type fields struct {
		cfg *Config
	}
	type args struct {
		text string
	}
	tests := []struct {
		fields fields
		args   args
		want   string
	}{
		{
			fields: fields{cfg: Standard()},
			args: args{
				text: "Saint Maximus (died 250) is a Christian saint and martyr.[1] The emperor Decius published a decree ordering the veneration of busts of the deified emperors.",
			},
			want: "Saint Maximus (died 250) is a Christian saint and martyr∯[1]\rThe emperor Decius published a decree ordering the veneration of busts of the deified emperors.",
		},
		{
			fields: fields{cfg: Standard()},
			args: args{
				text: "Differing agendas can potentially create an understanding gap in a consultation.11 12 Take the example of one of the most common presentations in ill health: the common cold.",
			},
			want: "Differing agendas can potentially create an understanding gap in a consultation∯11 12\rTake the example of one of the most common presentations in ill health: the common cold.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.text, func(t *testing.T) {
			p := &Processor{
				cfg: tt.fields.cfg,
			}
			if got := p.replacePeriodsBeforeNumericReferences(tt.args.text); got != tt.want {
				t.Errorf("processor.replacePeriodsBeforeNumericReferences() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_processor_checkForParens(t *testing.T) {
	type fields struct {
		cfg *Config
	}
	type args struct {
		text string
	}
	tests := []struct {
		fields fields
		args   args
		want   string
	}{
		{
			fields: fields{cfg: Standard()},
			args: args{
				text: `" (Dinah was the cat.) "`,
			},
			want: "\"\r(Dinah was the cat.)\r\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.args.text, func(t *testing.T) {
			p := &Processor{
				cfg: tt.fields.cfg,
			}
			if got := p.checkForParens(tt.args.text); got != tt.want {
				t.Errorf("processor.replacePeriodsBeforeNumericReferences() = %q, want %q", got, tt.want)
			}
		})
	}
}
