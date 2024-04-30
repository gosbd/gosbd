package lang

import (
	"regexp"

	"github.com/gosbd/gosbd/internal/processor"
	"github.com/gosbd/gosbd/internal/replacer"
)

func newChinese() *processor.Config {
	cfg := processor.Standard()
	punctuationReplacer := replacer.NewPunctuationReplacer()
	cfg.BetweenPunctuationReplacer = &betweenPunctuationReplacerChinese{
		punctuationReplacer: punctuationReplacer,
		betweenPunctuation:  replacer.NewBetweenPunctuation(punctuationReplacer),
	}
	cfg.SentenceStarters = nil
	return cfg
}

var (
	betweenDoubleAngledQuotationMarksZhRegex = regexp.MustCompile(`《([^》\\]+|\\{2}|\\.)*》`)
	punctuationBetweenLBracketsRegex         = regexp.MustCompile(`「([^」\\]+|\\{2}|\\.)*」`)
)

type betweenPunctuationReplacerChinese struct {
	punctuationReplacer replacer.Punctuation
	betweenPunctuation  replacer.BetweenPunctuation
}

func (b *betweenPunctuationReplacerChinese) Replace(text string) string {
	text = betweenDoubleAngledQuotationMarksZhRegex.ReplaceAllStringFunc(
		text, b.punctuationReplacer.ReplaceFunc(processor.PunctuationMatchTypeNone))
	text = punctuationBetweenLBracketsRegex.ReplaceAllStringFunc(
		text, b.punctuationReplacer.ReplaceFunc(processor.PunctuationMatchTypeNone))
	text = b.betweenPunctuation.Replace(text)
	return text
}

var _ processor.BetweenPunctuationReplacer = (*betweenPunctuationReplacerChinese)(nil)
