package replacer

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/gosbd/gosbd/internal/processor"
)

type AbbreviationReplacer struct {
	cfg *processor.Config
}

func (a AbbreviationReplacer) Replace(text string) string {
	abbr := a.cfg.Abbreviation
	text = abbr.PossessiveAbbreviationRule.Apply(text)
	text = abbr.KommanditgesellschaftRule.Apply(text)
	text = abbr.SingleLetterAbbreviationRules.Apply(text)

	abbrHandledTextBuf := bytes.NewBuffer(nil)
	// split by \r without removing it.
	for i, line := range strings.Split(text, "\r") {
		if i > 0 {
			abbrHandledTextBuf.WriteString("\r")
		}
		s := a.SearchForAbbreviationsInString(line)
		abbrHandledTextBuf.WriteString(s)
	}
	text = abbrHandledTextBuf.String()
	text = a.ReplaceMultiPeriodAbbreviations(text)
	text = abbr.AmPmRules.Apply(text)
	text = a.ReplaceAbbreviationAsSentenceBoundary(text)
	return text
}

func (a AbbreviationReplacer) ReplaceAbbreviationAsSentenceBoundary(text string) string {
	var ss []string
	for _, s := range a.cfg.SentenceStarters {
		ss = append(ss, fmt.Sprintf(`(\s%s\s)`, s))
	}
	sentenceStarters := strings.Join(ss, "|")
	re := regexp.MustCompile(fmt.Sprintf(`(U∯S|U\.S|U∯K|E∯U|E\.U|U∯S∯A|U\.S\.A|I|i\.v|I\.V)∯(%s)`, sentenceStarters))
	text = re.ReplaceAllString(text, "$1.$2")
	return text
}

func (a AbbreviationReplacer) ReplaceMultiPeriodAbbreviations(text string) string {
	mpaReplaceFn := func(match string) string {
		return strings.ReplaceAll(match, ".", "∯")
	}
	text = a.cfg.MultiPeriodAbbreviation.ReplaceAllStringFunc(text, mpaReplaceFn)
	return text
}

func (a AbbreviationReplacer) SearchForAbbreviationsInString(text string) string {
	lowered := strings.ToLower(text)
	for _, abbr := range a.cfg.Abbreviation.Abbreviations {
		stripped := strings.TrimSpace(abbr)
		if !strings.Contains(lowered, stripped) {
			continue
		}

		re := regexp.MustCompile(fmt.Sprintf(`(?i)(^|\s|\r|\n)%s`, stripped))
		matches := re.FindAllString(text, -1)
		if len(matches) == 0 {
			continue
		}

		escaped := regexp.QuoteMeta(stripped)
		reNextWord := regexp.MustCompile(fmt.Sprintf(`(?i)%s[ ](.{1})`, escaped))
		nextWordMatches := reNextWord.FindAllStringSubmatch(text, -1)
		for i, match := range matches {
			text = a.ScanForReplacements(text, match, i, nextWordMatches)
		}
	}
	return text
}

func (a AbbreviationReplacer) ScanForReplacements(text, am string, i int, nextWordMatches [][]string) string {
	var char string
	if len(nextWordMatches) > i && len(nextWordMatches[i]) > 1 {
		char = nextWordMatches[i][1]
	}
	upper := char != "" && strings.ToUpper(char) == char
	loweredMatch := strings.ToLower(strings.TrimSpace(am))
	if !upper || a.cfg.Abbreviation.IsPrePositive(loweredMatch) {
		if a.cfg.Abbreviation.IsPrePositive(loweredMatch) {
			text = a.ReplacePrePositiveAbbr(text, am)
		} else if a.cfg.Abbreviation.IsNumber(loweredMatch) {
			text = a.ReplacePreNumberAbbr(text, am)
		} else {
			text = a.ReplacePeriodOfAbbr(text, am)
		}
	}
	return text
}

func (a AbbreviationReplacer) ReplacePrePositiveAbbr(text, abbr string) string {
	// prepend a space to avoid needing another regex for start of string
	text = " " + text
	stripped := strings.TrimSpace(abbr)
	re := regexp.MustCompile(fmt.Sprintf(`(?i)(\s%s)\.(\s|:\d+)`, stripped))
	text = re.ReplaceAllString(text, "$1∯$2")
	return strings.TrimLeft(text, " ")
}

func (a AbbreviationReplacer) ReplacePreNumberAbbr(text, abbr string) string {
	// prepend a space to avoid needing another regex for start of string
	text = " " + text
	stripped := strings.TrimSpace(abbr)
	re := regexp.MustCompile(fmt.Sprintf(`(?i)(\s%s)\.(\s\d|\s+\()`, stripped))
	text = re.ReplaceAllString(text, "$1∯$2")[1:]
	return strings.TrimLeft(text, " ")
}

func (a AbbreviationReplacer) ReplacePeriodOfAbbr(text, abbr string) string {
	if a.cfg.Abbreviation.ReplacePeriodOfAbbrFn != nil {
		return a.cfg.Abbreviation.ReplacePeriodOfAbbrFn(text, abbr)
	}
	// prepend a space to avoid needing another regex for start of string
	text = " " + text
	stripped := strings.TrimSpace(abbr)
	re := regexp.MustCompile(fmt.Sprintf(`(\s%s)\.((\.|\:|-|\?|,)|(\s([a-z]|I\s|I'm|I'll|\d|\()))`, stripped))
	text = re.ReplaceAllString(text, "$1∯$2")[1:]
	return strings.TrimLeft(text, " ")
}

func NewAbbreviationReplacer(cfg *processor.Config) AbbreviationReplacer {
	return AbbreviationReplacer{cfg: cfg}
}

var _ processor.AbbreviationReplacer = (*AbbreviationReplacer)(nil)
