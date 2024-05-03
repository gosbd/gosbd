package replacer

import (
	"regexp"

	"github.com/gosbd/gosbd/internal/processor"
	"github.com/gosbd/gosbd/internal/rule"
)

type Punctuation struct {
	escapeRegexReservedCharacters     escapeRegexReservedCharacters
	subEscapedRegexReservedCharacters subEscapedRegexReservedCharacters
}

type escapeRegexReservedCharacters struct {
	leftParen    rule.Rule
	rightParen   rule.Rule
	leftBracket  rule.Rule
	rightBracket rule.Rule
	dash         rule.Rule
	all          rule.Rules
}

func newEscapeRegexReservedCharacters() escapeRegexReservedCharacters {
	leftParen := rule.NewRule(regexp.MustCompile(`\(`), `\\(`)
	rightParen := rule.NewRule(regexp.MustCompile(`\)`), `\\)`)
	leftBracket := rule.NewRule(regexp.MustCompile(`\[`), `\\[`)
	rightBracket := rule.NewRule(regexp.MustCompile(`]`), `\\]`)
	dash := rule.NewRule(regexp.MustCompile(`-`), `\\-`)
	all := rule.Rules{leftParen, rightParen, leftBracket, rightBracket, dash}

	return escapeRegexReservedCharacters{
		leftParen:    leftParen,
		rightParen:   rightParen,
		leftBracket:  leftBracket,
		rightBracket: rightBracket,
		dash:         dash,
		all:          all,
	}
}

type subEscapedRegexReservedCharacters struct {
	subLeftParen    rule.Rule
	subRightParen   rule.Rule
	subLeftBracket  rule.Rule
	subRightBracket rule.Rule
	subDash         rule.Rule
	all             rule.Rules
}

func newSubEscapedRegexReservedCharacters() subEscapedRegexReservedCharacters {
	subLeftParen := rule.NewRule(regexp.MustCompile(`\\\\\(`), `(`)
	subRightParen := rule.NewRule(regexp.MustCompile(`\\\\\)`), `)`)
	subLeftBracket := rule.NewRule(regexp.MustCompile(`\\\\\[`), `[`)
	subRightBracket := rule.NewRule(regexp.MustCompile(`\\\\]`), `]`)
	subDash := rule.NewRule(regexp.MustCompile(`\\\\-`), `-`)
	all := rule.Rules{subLeftParen, subRightParen, subLeftBracket, subRightBracket, subDash}

	return subEscapedRegexReservedCharacters{
		subLeftParen:    subLeftParen,
		subRightParen:   subRightParen,
		subLeftBracket:  subLeftBracket,
		subRightBracket: subRightBracket,
		subDash:         subDash,
		all:             all,
	}
}

func NewPunctuationReplacer() Punctuation {
	return Punctuation{
		escapeRegexReservedCharacters:     newEscapeRegexReservedCharacters(),
		subEscapedRegexReservedCharacters: newSubEscapedRegexReservedCharacters(),
	}
}

func (p *Punctuation) ReplaceFunc(matchType processor.PunctuationMatchType) func(string) string {
	return func(match string) string {
		text := p.escapeRegexReservedCharacters.all.Apply(match)
		sub := regexp.MustCompile(`\.`).ReplaceAllString(text, `∯`)
		sub1 := regexp.MustCompile(`。`).ReplaceAllString(sub, `&ᓰ&`)
		sub2 := regexp.MustCompile(`．`).ReplaceAllString(sub1, `&ᓱ&`)
		sub3 := regexp.MustCompile(`！`).ReplaceAllString(sub2, `&ᓳ&`)
		sub4 := regexp.MustCompile(`!`).ReplaceAllString(sub3, `&ᓴ&`)
		sub5 := regexp.MustCompile(`\?`).ReplaceAllString(sub4, `&ᓷ&`)
		lastSub := regexp.MustCompile(`？`).ReplaceAllString(sub5, `&ᓸ&`)
		if matchType != processor.PunctuationMatchTypeSingle {
			lastSub = regexp.MustCompile(`'`).ReplaceAllString(lastSub, `&⎋&`)
		}
		text = p.subEscapedRegexReservedCharacters.all.Apply(lastSub)
		return text
	}
}
