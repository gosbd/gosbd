package lang

import (
	"regexp"

	"github.com/gosbd/gosbd/internal/processor"
	"github.com/gosbd/gosbd/internal/replacer"
)

func newJapanese() *processor.Config {
	cfg := processor.Standard()
	punctuationReplacer := replacer.NewPunctuationReplacer()
	cfg.BetweenPunctuationReplacer = &betweenPunctuationReplacerJapanese{
		punctuationReplacer: punctuationReplacer,
		betweenPunctuation:  replacer.NewBetweenPunctuation(punctuationReplacer),
	}
	cfg.SentenceStarters = nil
	return cfg
}

var (
	betweenParensJaRegex = regexp.MustCompile(`（([^（）]+|\\{2}|\\.)*）`)
	betweenQuoteJaRegex  = regexp.MustCompile(`「(([^「」]+|\\{2}|\\.)*)」`)
)

type betweenPunctuationReplacerJapanese struct {
	punctuationReplacer replacer.Punctuation
	betweenPunctuation  replacer.BetweenPunctuation
}

func (b *betweenPunctuationReplacerJapanese) Replace(text string) string {
	text = betweenParensJaRegex.ReplaceAllStringFunc(
		text, b.punctuationReplacer.ReplaceFunc(processor.PunctuationMatchTypeNone))
	text = betweenQuoteJaRegex.ReplaceAllStringFunc(
		text, b.punctuationReplacer.ReplaceFunc(processor.PunctuationMatchTypeNone))
	text = b.betweenPunctuation.Replace(text)
	return text
}

var _ processor.BetweenPunctuationReplacer = (*betweenPunctuationReplacerJapanese)(nil)
