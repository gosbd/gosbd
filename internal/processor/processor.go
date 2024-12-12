package processor

import (
	"regexp"
	"strings"
)

type ListItemReplacer interface {
	AddLineBreak(text string) string
	ReplaceParens(text string) string
}

type AbbreviationReplacer interface {
	Replace(text string) string
}

type PunctuationReplacer interface {
	ReplaceFunc(matchType PunctuationMatchType) func(string) string
}

type BetweenPunctuationReplacer interface {
	Replace(text string) string
}

type Processor struct {
	cfg                        *Config
	listItemReplacer           ListItemReplacer
	abbrReplacer               AbbreviationReplacer
	punctuationReplacer        PunctuationReplacer
	betweenPunctuationReplacer BetweenPunctuationReplacer
}

// Process implements segmenter.Processor.
func (p *Processor) Process(text string) []string {
	text = strings.ReplaceAll(text, "\n", "\r")
	text = p.listItemReplacer.AddLineBreak(text)
	text = p.abbrReplacer.Replace(text)
	text = p.replaceNumbers(text)
	text = p.replacePeriodsBeforeNumericReferences(text)
	text = p.cfg.Abbreviation.WithMultiplePeriodsAndEmailRule.Apply(text)
	text = p.cfg.GeoLocationRule.Apply(text)
	text = p.cfg.FileFormatRule.Apply(text)
	return p.splitIntoSegments(text)
}

func (p *Processor) splitIntoSegments(text string) []string {
	text = p.checkForParens(text)
	sentences := strings.Split(text, "\r")

	// remove empty values
	sentences = p.filterEmpty(sentences)
	p.applySingleNewLineRule(sentences)
	p.applyEllipsisRule(sentences)

	var sentences2 []string
	for _, s := range sentences {
		sentences2 = append(sentences2, p.checkForPunctuation(s)...)
	}

	var postProcessedSentences []string
	ss := p.filterEmpty(sentences2)
	for _, sent := range ss {
		sent = p.cfg.SubSymbolsRules.All.Apply(sent)
		postProcessedSentences = append(postProcessedSentences, p.postProcessSegments(sent)...)
	}
	for i, s := range postProcessedSentences {
		postProcessedSentences[i] = p.cfg.SubSingleQuoteRule.Apply(s)
	}
	return postProcessedSentences
}

var (
	postProcessRegex  = regexp.MustCompile(`\A[a-zA-Z]*\z`)
	postProcessRegex2 = regexp.MustCompile(`\t`)
)

func (p *Processor) postProcessSegments(text string) []string {
	if len(text) > 2 && postProcessRegex.MatchString(text) {
		return []string{text}
	}
	// below condition present in pragmatic segmenter
	// dont know significance of it yet.
	// if self.consecutive_underscore(txt) or len(txt) < 2:
	//     return txt
	if postProcessRegex2.MatchString(text) {
		return []string{text}
	}

	// TODO: (from pySBD)
	// Decide on keeping or removing Standard.ExtraWhiteSpaceRule
	// removed to retain original text spans
	// txt = Text(txt).apply(*ReinsertEllipsisRules.All,
	//                       Standard.ExtraWhiteSpaceRule)
	text = p.cfg.ReinsertEllipsisRules.All.Apply(text)
	if p.cfg.QuotationAtEndOfSentenceRegex.MatchString(text) {
		return strings.Split(p.cfg.SplitSpaceQuotationAtEndOfSentenceRule.Apply(text), "\r")
	} else {
		text = strings.ReplaceAll(text, "\n", "")
		return []string{strings.TrimSpace(text)}
	}
}

func (p *Processor) applySingleNewLineRule(sents []string) {
	for i, s := range sents {
		sents[i] = p.cfg.SingleNewLineRule.Apply(s)
	}
}

func (p *Processor) applyEllipsisRule(sents []string) {
	for i, s := range sents {
		sents[i] = p.cfg.Ellipsis.All.Apply(s)
	}
}

func (p *Processor) filterEmpty(sents []string) []string {
	var res []string
	for _, s := range sents {
		if s != "" && s != " " {
			res = append(res, s)
		}
	}
	return res
}

func (p *Processor) checkForPunctuation(text string) []string {
	for _, punctuation := range p.cfg.Punctuations {
		if strings.Contains(text, punctuation) {
			return p.processText(text)
		}
	}
	return []string{text}
}

func (p *Processor) processText(text string) []string {
	text = p.checkPunctuationAtEnd(text)
	text = p.replaceExclamationWords(text)
	text = p.betweenPunctuationReplacer.Replace(text)
	// handle text having only double punctuations
	if !p.cfg.DoublePunctuationRules.DoublePunctuationRegex.MatchString(text) {
		text = p.cfg.DoublePunctuationRules.All.Apply(text)
	}
	text = p.cfg.QuestionMarkInQuotationRule.Apply(text)
	text = p.cfg.ExclamationPointRules.All.Apply(text)
	text = p.listItemReplacer.ReplaceParens(text)

	return p.sentenceBoundaryPunctuation(text)
}

var (
	exclamationRegex = regexp.MustCompile(`&ᓴ&$`)
)

func (p *Processor) sentenceBoundaryPunctuation(text string) []string {
	// TODO: implement rules below
	// if hasattr(self.lang, 'ReplaceColonBetweenNumbersRule'):
	//    txt = Text(txt).apply(
	//    self.lang.ReplaceColonBetweenNumbersRule)
	// if hasattr(self.lang, 'ReplaceNonSentenceBoundaryCommaRule'):
	// 	txt = Text(txt).apply(
	// 	self.lang.ReplaceNonSentenceBoundaryCommaRule)

	// retain exclamation mark if it is an ending character of a given text
	text = exclamationRegex.ReplaceAllString(text, "!")
	priorIndex := 0
	for _, rule := range p.cfg.SentenceBoundaryRules.All {
		maxIdx := 0
		for _, match := range rule.Pattern().FindAllStringIndex(text, -1) {
			if match[1] > maxIdx {
				maxIdx = match[1]
			}
		}
		if maxIdx == 0 {
			continue
		}
		if priorIndex > len(text) {
			priorIndex = len(text)
		}
		text = text[:priorIndex] + rule.Apply(text[priorIndex:])
		priorIndex = maxIdx - 1
	}
	return p.filterEmpty(strings.Split(text, "\r"))
}

func (p *Processor) checkPunctuationAtEnd(text string) string {
	hasPunctuationAtEnd := false
	for _, punctuation := range p.cfg.Punctuations {
		if strings.HasSuffix(text, punctuation) {
			hasPunctuationAtEnd = true
			break
		}
	}
	if !hasPunctuationAtEnd {
		text += "ȸ"
	}
	return text
}

func (p *Processor) replaceExclamationWords(text string) string {
	return p.cfg.ExclamationWords.Regex.ReplaceAllStringFunc(text,
		p.punctuationReplacer.ReplaceFunc(PunctuationMatchTypeNone))
}

func (p *Processor) checkForParens(text string) string {
	return p.cfg.ParensBetweenDoubleQuotesRule.Apply(text)
}

func (p *Processor) replacePeriodsBeforeNumericReferences(text string) string {
	return p.cfg.NumberedReferenceRegex.ReplaceAllString(text, "$1∯$3\r$9")
}

func (p *Processor) replaceContinuousPunctuation(text string) string {
	replaceFunc := func(match string) string {
		replaced := regexp.MustCompile(`!`).ReplaceAllString(match, "&ᓴ&")
		replaced = regexp.MustCompile(`\?`).ReplaceAllString(replaced, "&ᓷ&")
		return replaced
	}
	return p.cfg.ContinuousPunctuationRegex.ReplaceAllStringFunc(text, replaceFunc)
}

func (p *Processor) replaceNumbers(text string) string {
	return p.cfg.Numbers.All().Apply(text)
}

type Params struct {
	Lang                       *Config
	ListItemReplacer           ListItemReplacer
	AbbrReplacer               AbbreviationReplacer
	PunctuationReplacer        PunctuationReplacer
	BetweenPunctuationReplacer BetweenPunctuationReplacer
}

func NewProcessor(params Params) *Processor {
	return &Processor{
		cfg:                        params.Lang,
		listItemReplacer:           params.ListItemReplacer,
		abbrReplacer:               params.AbbrReplacer,
		punctuationReplacer:        params.PunctuationReplacer,
		betweenPunctuationReplacer: params.BetweenPunctuationReplacer,
	}
}
