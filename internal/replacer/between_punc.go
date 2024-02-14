package replacer

import (
	"regexp"

	"github.com/gosbd/gosbd/internal/processor"
)

var (
	singleQuoteSpaceRegex = regexp.MustCompile(`'\s`)
	// Rubular: http://rubular.com/r/2YFrKWQUYi
	betweenSingleQuotesRegex = regexp.MustCompile(`(\s)'(?:[^']|'[a-zA-Z])*'`)
	// Rubular: http://rubular.com/r/mXf8cW025o
	wordWithLeadingApostropheRegex = regexp.MustCompile(`(\s)'(?:[^']|'[a-zA-Z])*'\S`)

	betweenSingleQuoteSlantedRegex = regexp.MustCompile(`(\s)‘(?:[^’]|’[a-zA-Z])*’`)
	// https://regex101.com/r/r6I1bW/1
	// https://stackoverflow.com/questions/13577372/do-python-regular-expressions-have-an-equivalent-to-rubys-atomic-grouping?noredirect=1&lq=1
	betweenDoubleQuotesRegex = regexp.MustCompile(`"([^"\\]|\\.)*"`)

	// Rubular: http://rubular.com/r/WX4AvnZvlX
	betweenSquareBracketsRegex = regexp.MustCompile(`\[([^]\\]+|\\.)*]`)

	// Rubular: http://rubular.com/r/6tTityPflI
	betweenParensRegex = regexp.MustCompile(`\(([^()\\]+|\\{2}|\\.)*\)`)

	// Rubular: http://rubular.com/r/x6s4PZK8jc
	betweenQuoteArrowRegex = regexp.MustCompile(`«([^»\\]+|\\{2}|\\.)*»`)

	// Rubular: http://rubular.com/r/jTtDKfjxzr
	betweenEmDashesRegex = regexp.MustCompile(`--([^--]*)--`)

	// Rubular: http://rubular.com/r/JbAIpKdlSq
	betweenQuoteSlantedRegex = regexp.MustCompile(`“([^”\\]+|\\{2}|\\.)*”`)
)

type BetweenPunctuation struct {
	punctuationReplacer Punctuation
}

func NewBetweenPunctuation(pr Punctuation) BetweenPunctuation {
	return BetweenPunctuation{
		punctuationReplacer: pr,
	}
}

func (b BetweenPunctuation) Replace(text string) string {
	text = b.punctuationBetweenSingleQuotes(text)
	text = b.punctuationBetweenQuoteSlanted(text)
	text = b.punctuationBetweenDoubleQuotes(text)
	text = b.punctuationBetweenSquareBrackets(text)
	text = b.punctuationBetweenParens(text)
	text = b.punctuationBetweenQuoteArrow(text)
	text = b.punctuationBetweenEmDashes(text)
	text = b.punctuationBetweenSingleQuoteSlanted(text)
	return text
}

func (b BetweenPunctuation) punctuationBetweenQuoteSlanted(text string) string {
	return betweenQuoteSlantedRegex.ReplaceAllStringFunc(
		text, b.punctuationReplacer.ReplaceFunc(processor.PunctuationMatchTypeNone))
}

func (b BetweenPunctuation) punctuationBetweenEmDashes(text string) string {
	return betweenEmDashesRegex.ReplaceAllStringFunc(
		text, b.punctuationReplacer.ReplaceFunc(processor.PunctuationMatchTypeNone))
}

func (b BetweenPunctuation) punctuationBetweenQuoteArrow(text string) string {
	return betweenQuoteArrowRegex.ReplaceAllStringFunc(
		text, b.punctuationReplacer.ReplaceFunc(processor.PunctuationMatchTypeNone))
}

func (b BetweenPunctuation) punctuationBetweenParens(text string) string {
	return betweenParensRegex.ReplaceAllStringFunc(
		text, b.punctuationReplacer.ReplaceFunc(processor.PunctuationMatchTypeNone))
}

func (b BetweenPunctuation) punctuationBetweenSquareBrackets(text string) string {
	return betweenSquareBracketsRegex.ReplaceAllStringFunc(
		text, b.punctuationReplacer.ReplaceFunc(processor.PunctuationMatchTypeNone))
}

func (b BetweenPunctuation) punctuationBetweenDoubleQuotes(text string) string {
	return betweenDoubleQuotesRegex.ReplaceAllStringFunc(
		text, b.punctuationReplacer.ReplaceFunc(processor.PunctuationMatchTypeNone))
}

func (b BetweenPunctuation) punctuationBetweenSingleQuoteSlanted(text string) string {
	return betweenSingleQuoteSlantedRegex.ReplaceAllStringFunc(
		text, b.punctuationReplacer.ReplaceFunc(processor.PunctuationMatchTypeNone))
}

func (b BetweenPunctuation) punctuationBetweenSingleQuotes(text string) string {
	if wordWithLeadingApostropheRegex.MatchString(text) && !singleQuoteSpaceRegex.MatchString(text) {
		return text
	}
	return betweenSingleQuotesRegex.ReplaceAllStringFunc(
		text, b.punctuationReplacer.ReplaceFunc(processor.PunctuationMatchTypeSingle))
}
