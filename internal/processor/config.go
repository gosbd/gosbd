package processor

import (
	"regexp"
	"slices"
	"strings"

	"github.com/yohamta/gosbd/internal/rule"
)

type Abbreviation struct {
	PossessiveAbbreviationRule      rule.Rule
	KommanditgesellschaftRule       rule.Rule
	SingleLetterAbbreviationRules   rule.Rules
	AmPmRules                       rule.Rules
	WithMultiplePeriodsAndEmailRule rule.Rule
	Abbreviations                   []string
	PrePositiveAbbreviations        []string
	NumberAbbreviations             []string
	ReplacePeriodOfAbbrFn           func(text, abbr string) string
}

func (a Abbreviation) IsPrePositive(abbr string) bool {
	return slices.Contains(a.PrePositiveAbbreviations, abbr)
}

func (a Abbreviation) IsNumber(abbr string) bool {
	return slices.Contains(a.NumberAbbreviations, abbr)
}

type ExclamationWords struct {
	Words []string
	Regex *regexp.Regexp
}

type Ellipsis struct {
	All rule.Rules
}

type DoublePunctuationRules struct {
	DoublePunctuationRegex *regexp.Regexp
	All                    rule.Rules
}

type ExclamationPointRules struct {
	All rule.Rules
}

type ReinsertEllipsisRules struct {
	All rule.Rules
}

type SentenceBoundaryRules struct {
	All rule.Rules
}

type Numbers struct {
	PeriodBeforeNumberRule             rule.Rule
	NumberAfterPeriodBeforeLetterRule  rule.Rule
	NewLineNumberPeriodSpaceLetterRule rule.Rule
	StartLineNumberPeriodRule          rule.Rule
	StartLineTwoDigitNumberPeriodRule  rule.Rule
}

func (n Numbers) All() rule.Rules {
	return rule.Rules{
		n.PeriodBeforeNumberRule,
		n.NumberAfterPeriodBeforeLetterRule,
		n.NewLineNumberPeriodSpaceLetterRule,
		n.StartLineNumberPeriodRule,
		n.StartLineTwoDigitNumberPeriodRule,
	}
}

type SubSymbolsRules struct {
	Period                     rule.Rule
	ArabicComma                rule.Rule
	SemiColon                  rule.Rule
	FullWidthPeriod            rule.Rule
	SpecialPeriod              rule.Rule
	FullWidthExclamation       rule.Rule
	ExclamationPoint           rule.Rule
	QuestionMark               rule.Rule
	FullWidthQuestionMark      rule.Rule
	MixedDoubleQE              rule.Rule
	MixedDoubleQQ              rule.Rule
	MixedDoubleEQ              rule.Rule
	MixedDoubleEE              rule.Rule
	LeftParens                 rule.Rule
	RightParens                rule.Rule
	TemporaryEndingPunctuation rule.Rule
	Newline                    rule.Rule
	All                        rule.Rules
}

type PunctuationMatchType string

const (
	PunctuationMatchTypeNone   PunctuationMatchType = "none"
	PunctuationMatchTypeSingle PunctuationMatchType = "single"
)

func newReinsertEllipsisRules() ReinsertEllipsisRules {
	// below rules aren't similar to original rules of pragmatic segmenter
	// modification: symbols replaced with same number of ellipses
	rules := ReinsertEllipsisRules{}
	rules.All = rule.Rules{
		rule.NewRule(regexp.MustCompile(`ƪƪƪ`), "..."),
		rule.NewRule(regexp.MustCompile(`♟♟♟♟♟♟♟`), " . . . "),
		rule.NewRule(regexp.MustCompile(`♝♝♝♝♝♝♝`), ". . . ."),
		rule.NewRule(regexp.MustCompile(`☏☏`), ".."),
		rule.NewRule(regexp.MustCompile(`∮`), "."),
	}
	return rules
}

func newSubSymbolsRules() SubSymbolsRules {
	rules := SubSymbolsRules{
		Period:                     rule.NewRule(regexp.MustCompile(`∯`), "."),
		ArabicComma:                rule.NewRule(regexp.MustCompile(`♬`), "،"),
		SemiColon:                  rule.NewRule(regexp.MustCompile(`♭`), ":"),
		FullWidthPeriod:            rule.NewRule(regexp.MustCompile(`&ᓰ&`), "。"),
		SpecialPeriod:              rule.NewRule(regexp.MustCompile(`&ᓱ&`), "．"),
		FullWidthExclamation:       rule.NewRule(regexp.MustCompile(`&ᓳ&`), "！"),
		ExclamationPoint:           rule.NewRule(regexp.MustCompile(`&ᓴ&`), "!"),
		QuestionMark:               rule.NewRule(regexp.MustCompile(`&ᓷ&`), "?"),
		FullWidthQuestionMark:      rule.NewRule(regexp.MustCompile(`&ᓸ&`), "？"),
		MixedDoubleQE:              rule.NewRule(regexp.MustCompile(`☉`), "?!"),
		MixedDoubleQQ:              rule.NewRule(regexp.MustCompile(`☇`), "??"),
		MixedDoubleEQ:              rule.NewRule(regexp.MustCompile(`☈`), "!?"),
		MixedDoubleEE:              rule.NewRule(regexp.MustCompile(`☄`), "!!"),
		LeftParens:                 rule.NewRule(regexp.MustCompile(`&✂&`), "("),
		RightParens:                rule.NewRule(regexp.MustCompile(`&⌬&`), ")"),
		TemporaryEndingPunctuation: rule.NewRule(regexp.MustCompile(`ȸ`), ""),
		Newline:                    rule.NewRule(regexp.MustCompile(`ȹ`), "\n"),
	}
	rules.All = rule.Rules{
		rules.Period,
		rules.ArabicComma,
		rules.SemiColon,
		rules.FullWidthPeriod,
		rules.SpecialPeriod,
		rules.FullWidthExclamation,
		rules.ExclamationPoint,
		rules.QuestionMark,
		rules.FullWidthQuestionMark,
		rules.MixedDoubleQE,
		rules.MixedDoubleQQ,
		rules.MixedDoubleEQ,
		rules.MixedDoubleEE,
		rules.LeftParens,
		rules.RightParens,
		rules.TemporaryEndingPunctuation,
		rules.Newline,
	}
	return rules
}

type Config struct {
	BetweenDoubleQuotes                    *regexp.Regexp
	NumberedReferenceRegex                 *regexp.Regexp
	MultiPeriodAbbreviation                *regexp.Regexp
	QuotationAtEndOfSentenceRegex          *regexp.Regexp
	ContinuousPunctuationRegex             *regexp.Regexp
	SingleNewLineRule                      rule.Rule
	GeoLocationRule                        rule.Rule
	FileFormatRule                         rule.Rule
	SplitSpaceQuotationAtEndOfSentenceRule rule.Rule
	ParensBetweenDoubleQuotesRule          rule.Rule
	QuestionMarkInQuotationRule            rule.Rule
	Abbreviation                           Abbreviation
	Numbers                                Numbers
	Ellipsis                               Ellipsis
	ExclamationWords                       ExclamationWords
	DoublePunctuationRules                 DoublePunctuationRules
	ExclamationPointRules                  ExclamationPointRules
	Punctuations                           []string
	SubSymbolsRules                        SubSymbolsRules
	ReinsertEllipsisRules                  ReinsertEllipsisRules
	SubSingleQuoteRule                     rule.Rule
	SentenceBoundaryRules                  SentenceBoundaryRules
	SentenceStarters                       []string
	BetweenPunctuationReplacer             BetweenPunctuationReplacer
}

func Standard() *Config {
	return &Config{
		SubSingleQuoteRule:                     subSingleQuoteRule,
		MultiPeriodAbbreviation:                multiPeriodAbbreviation,
		BetweenDoubleQuotes:                    betweenDoubleQuotesRegex,
		NumberedReferenceRegex:                 numberedReferenceRegex,
		ParensBetweenDoubleQuotesRule:          parensBetweenDoubleQuotesRule,
		SingleNewLineRule:                      singleNewLineRule,
		QuestionMarkInQuotationRule:            questionMarkInQuotationRule,
		GeoLocationRule:                        geoLocationRule,
		FileFormatRule:                         fileFormatRule,
		QuotationAtEndOfSentenceRegex:          quotationAtEndOfSentenceRegex,
		ContinuousPunctuationRegex:             continuousPunctuationRegex,
		SplitSpaceQuotationAtEndOfSentenceRule: splitSpaceQuotationAtEndOfSentenceRule,
		Abbreviation: Abbreviation{
			PossessiveAbbreviationRule:      possessiveAbbreviationRule,
			KommanditgesellschaftRule:       kommanditgesellschaftRule,
			SingleLetterAbbreviationRules:   rule.Rules{singleUpperCaseLetterAtStartOfLineRule, singleUpperCaseLetterRule},
			AmPmRules:                       rule.Rules{upperCasePmRule, upperCaseAmRule, lowerCasePmRule, lowerCaseAmRule},
			Abbreviations:                   []string{"adj", "adm", "adv", "al", "ala", "alta", "apr", "arc", "ariz", "ark", "art", "assn", "asst", "attys", "aug", "ave", "bart", "bld", "bldg", "blvd", "brig", "bros", "btw", "cal", "calif", "capt", "cl", "cmdr", "co", "col", "colo", "comdr", "con", "conn", "corp", "cpl", "cres", "ct", "d.phil", "dak", "dec", "del", "dept", "det", "dist", "dr", "dr.phil", "dr.philos", "drs", "e.g", "ens", "esp", "esq", "etc", "exp", "expy", "ext", "feb", "fed", "fla", "ft", "fwy", "fy", "ga", "gen", "gov", "hon", "hosp", "hr", "hway", "hwy", "i.e", "ia", "id", "ida", "ill", "inc", "ind", "ing", "insp", "is", "jan", "jr", "jul", "jun", "kan", "kans", "ken", "ky", "la", "lt", "ltd", "maj", "man", "mar", "mass", "may", "md", "me", "med", "messrs", "mex", "mfg", "mich", "min", "minn", "miss", "mlle", "mm", "mme", "mo", "mont", "mr", "mrs", "ms", "msgr", "mssrs", "mt", "mtn", "neb", "nebr", "nev", "no", "nos", "nov", "nr", "oct", "ok", "okla", "ont", "op", "ord", "ore", "p", "pa", "pd", "pde", "penn", "penna", "pfc", "ph", "ph.d", "pl", "plz", "pp", "prof", "pvt", "que", "rd", "rs", "ref", "rep", "reps", "res", "rev", "rt", "sask", "sec", "sen", "sens", "sep", "sept", "sfc", "sgt", "sr", "st", "supt", "surg", "tce", "tenn", "tex", "univ", "usafa", "u.s", "ut", "va", "v", "ver", "viz", "vs", "vt", "wash", "wis", "wisc", "wy", "wyo", "yuk", "fig"},
			PrePositiveAbbreviations:        []string{"adm", "attys", "brig", "capt", "cmdr", "col", "cpl", "det", "dr", "gen", "gov", "ing", "lt", "maj", "mr", "mrs", "ms", "mt", "messrs", "mssrs", "prof", "ph", "rep", "reps", "rev", "sen", "sens", "sgt", "st", "supt", "v", "vs", "fig"},
			NumberAbbreviations:             []string{"art", "ext", "no", "nos", "p", "pp"},
			WithMultiplePeriodsAndEmailRule: withMultiplePeriodsAndEmailRule,
		},
		SentenceStarters: strings.Split("A Being Did For He How However I In It Millions More She That The There They We What When Where Who Why", " "),
		Numbers: Numbers{
			PeriodBeforeNumberRule:             periodBeforeNumberRule,
			NumberAfterPeriodBeforeLetterRule:  numberAfterPeriodBeforeLetterRule,
			NewLineNumberPeriodSpaceLetterRule: newLineNumberPeriodSpaceLetterRule,
			StartLineNumberPeriodRule:          startLineNumberPeriodRule,
			StartLineTwoDigitNumberPeriodRule:  startLineTwoDigitNumberPeriodRule,
		},
		Punctuations:     []string{"。", "．", ".", "！", "!", "?", "？"},
		ExclamationWords: newExclamationWords(),
		Ellipsis: Ellipsis{
			All: rule.Rules{ellipsisThreeSpaceRule, ellipsisFourSpaceRule, ellipsisFourConsecutiveRule, ellipsisThreeConsecutiveRule, ellipsisOtherThreePeriodRule},
		},
		DoublePunctuationRules: DoublePunctuationRules{
			DoublePunctuationRegex: doublePunctuationRegex,
			All:                    rule.Rules{firstDoublePunctuationRule, secondDoublePunctuationRule, thirdDoublePunctuationRule, forthDoublePunctuationRule},
		},
		ExclamationPointRules: ExclamationPointRules{
			All: rule.Rules{exclamationPointInQuotationRule, exclamationPointBeforeCommaMidSentenceRule, exclamationPointMidSentenceRule},
		},
		SubSymbolsRules:       newSubSymbolsRules(),
		ReinsertEllipsisRules: newReinsertEllipsisRules(),
		SentenceBoundaryRules: SentenceBoundaryRules{
			All: rule.Rules{
				sentenceBoundaryRule1,
				sentenceBoundaryRule2,
				sentenceBoundaryRule3,
				sentenceBoundaryRule4,
				sentenceBoundaryRule5,
				sentenceBoundaryRule6,
				sentenceBoundaryRule7,
				sentenceBoundaryRule8,
				sentenceBoundaryRule9,
			},
		},
	}
}

func newExclamationWords() ExclamationWords {
	words := strings.Split("!Xũ !Kung ǃʼOǃKung !Xuun !Kung-Ekoka ǃHu ǃKhung ǃKu ǃung ǃXo ǃXû ǃXung ǃXũ !Xun Yahoo! Y!J Yum!", " ")
	escaped := make([]string, len(words))

	for i, word := range words {
		escaped[i] = regexp.QuoteMeta(word)
	}
	re := regexp.MustCompile(strings.Join(escaped, "|"))
	return ExclamationWords{
		Words: words,
		Regex: re,
	}
}

var (
	// Rubular: http://rubular.com/r/EUbZCNfgei
	// WithMultiplePeriodsAndEmailRule = Rule(r'(\w)(\.)(\w)', '\\1∮\\3')
	// \w in python matches unicode abbreviations also so limit to english alphanumerics
	withMultiplePeriodsAndEmailRule = rule.NewRule(regexp.MustCompile(`([a-zA-Z0-9_])(\.)([a-zA-Z0-9_])`), `$1∮$3`)
	// Rubular: http://rubular.com/r/yqa4Rit8EY
	possessiveAbbreviationRule = rule.NewRule(regexp.MustCompile(`\.('s[\s$])`), "∯$1")
	// Rubular: http://rubular.com/r/xDkpFZ0EgH
	kommanditgesellschaftRule = rule.NewRule(regexp.MustCompile(`(Co)\.(\sKG)`), "$1∯$2")
	// Rubular: http://rubular.com/r/e3H6kwnr6H
	singleUpperCaseLetterAtStartOfLineRule = rule.NewRule(regexp.MustCompile(`(^[A-Z])\.(\s)`), "$1∯$2")
	// Rubular: http://rubular.com/r/gitvf0YWH4
	singleUpperCaseLetterRule = rule.NewRule(regexp.MustCompile(`(\s[A-Z])\.(,?\s)`), "$1∯$2")
	// Rubular: http://rubular.com/r/G2opjedIm9
	geoLocationRule = rule.NewRule(regexp.MustCompile(`([a-zA-z]°)\.(\s*\d+)`), "$1∯$2")
	fileFormatRule  = rule.NewRule(regexp.MustCompile(`(\s)\.((jpe?g|png|gif|tiff?|pdf|ps|docx?|xlsx?|svg|bmp|tga|exif|odt|html?|txt|rtf|bat|sxw|xml|zip|exe|msi|blend|wmv|mp[34]|pptx?|flac|rb|cpp|cs|js)\s)`), "$1∯$2")
	// Rubular: http://rubular.com/r/aXPUGm6fQh
	// QuestionMarkInQuotationRule = Rule(r'\?(?=(\'|\"))', '&ᓷ&')
	questionMarkInQuotationRule = rule.NewRule(regexp.MustCompile(`\?(['"])`), "&ᓷ&")
)

var (
	// Rubular: http://rubular.com/r/xDkpFZ0EgH
	multiPeriodAbbreviation = regexp.MustCompile(`(?i)\b[a-z](?:\.[a-z])+[.]`)
	// Rubular: http://rubular.com/r/TYzr4qOW1Q
	betweenDoubleQuotesRegex = regexp.MustCompile(`([^"])*[^, ]"|“(?: [ ^”])*[^, ]`)
	// Rubular: http://rubular.com/r/mQ8Es9bxtk
	continuousPunctuationRegex = regexp.MustCompile(`(\S)([!?]){3,}(\s|\z|$)`)
	// Rubular: http://rubular.com/r/NqCqv372Ix
	quotationAtEndOfSentenceRegex = regexp.MustCompile(`[!?.-]["'“”]\s[A-Z]`)
	// Rubular: http://rubular.com/r/JMjlZHAT4g
	splitSpaceQuotationAtEndOfSentenceRule = rule.NewRule(regexp.MustCompile(`([!?.-]["'“”])\s([A-Z])`), "$1\r$2")
	// https://rubular.com/r/UkumQaILKbkeyc
	// https://github.com/diasks2/pragmatic_segmenter/commit/d9ec1a352aff92b91e2e572c30bb9561eb42c703
	numberedReferenceRegex = regexp.MustCompile(`([^\d\s])([.|∯])(\[((\d{1,3},?\s?-?\s?)*\b\d{1,3}])+|((\d{1,3}\s?)?\d{1,3}))(\s)([A-Z])`)
)

var (
	// Rubular: http://rubular.com/r/Vnx3m4Spc8
	upperCasePmRule = rule.NewRule(regexp.MustCompile(`(P∯M)∯(\s[A-Z])`), "$1.$2")
	// Rubular: http://rubular.com/r/AJMCotJVbW
	upperCaseAmRule = rule.NewRule(regexp.MustCompile(`(A∯M)∯(\s[A-Z])`), "$1.$2")
	// Rubular: http://rubular.com/r/13q7SnOhgA
	lowerCasePmRule = rule.NewRule(regexp.MustCompile(`(p∯m)∯(\s[A-Z])`), "$1.$2")
	// Rubular: http://rubular.com/r/DgUDq4mLz5
	lowerCaseAmRule = rule.NewRule(regexp.MustCompile(`(a∯m)∯(\s[A-Z])`), "$1.$2")
	// Rubular: http://rubular.com/r/6flGnUMEVl
	parensBetweenDoubleQuotesRule = rule.NewRule(regexp.MustCompile(`(["”])\s(\(.*\))\s(["“])`), "$1\r$2\r$3")
	singleNewLineRule             = rule.NewRule(regexp.MustCompile("\n"), "ȹ")
	subSingleQuoteRule            = rule.NewRule(regexp.MustCompile(`&⎋&`), "'")
)

var (
	periodBeforeNumberRule = rule.NewRule(regexp.MustCompile(`\.(\d)`), "∯$1")
	// Rubular: http://rubular.com/r/EMk5MpiUzt
	numberAfterPeriodBeforeLetterRule = rule.NewRule(regexp.MustCompile(`(\d)\.(\S)`), "$1∯$2")
	// Rubular: http://rubular.com/r/rf4l1HjtjG
	newLineNumberPeriodSpaceLetterRule = rule.NewRule(regexp.MustCompile(`(\r\d)\.((\s\S)|\))`), "$1∯$2")
	// Rubular: http://rubular.com/r/HPa4sdc6b9
	startLineNumberPeriodRule = rule.NewRule(regexp.MustCompile(`(^\d)\.((\s\S)|\))`), "$1∯$2")
	// Rubular: http://rubular.com/r/NuvWnKleFl
	startLineTwoDigitNumberPeriodRule = rule.NewRule(regexp.MustCompile(`(^\d\d)\.((\s\S)|\))`), "$1∯$2")
)

var (
	// Rubular: http://rubular.com/r/YBG1dIHTRu
	ellipsisThreeSpaceRule = rule.NewRule(regexp.MustCompile(`(\s\.){3}\s`), "♟♟♟♟♟♟♟")
	// Rubular: http://rubular.com/r/2VvZ8wRbd8
	ellipsisFourSpaceRule = rule.NewRule(regexp.MustCompile(`([a-z])(\.\s){3}\.($|\\n)`), "${1}♝♝♝♝♝♝♝")
	// Rubular: http://rubular.com/r/Hdqpd90owl
	ellipsisFourConsecutiveRule = rule.NewRule(regexp.MustCompile(`(\S)([.]{3})(\.\s[A-Z])`), "${1}ƪƪƪ${3}")
	// below rules aren't similar to original rules of pragmatic segmenter
	// modification: spaces replaced with same number of symbols
	// Rubular: http://rubular.com/r/i60hCK81fz
	ellipsisThreeConsecutiveRule = rule.NewRule(regexp.MustCompile(`\.\.\.(\s+[A-Z])`), "☏☏.${1}")
	ellipsisOtherThreePeriodRule = rule.NewRule(regexp.MustCompile(`\.\.\.`), "ƪƪƪ")
)

var (
	firstDoublePunctuationRule  = rule.NewRule(regexp.MustCompile(`\?!`), "☉")
	secondDoublePunctuationRule = rule.NewRule(regexp.MustCompile(`!\?`), "☈")
	thirdDoublePunctuationRule  = rule.NewRule(regexp.MustCompile(`\?\?`), "☇")
	forthDoublePunctuationRule  = rule.NewRule(regexp.MustCompile(`!!`), "☄")
	doublePunctuationRegex      = regexp.MustCompile(`^(\?!|!\?|\?\?|!!)`)
)

var (
	// Rubular: http://rubular.com/r/XS1XXFRfM2
	exclamationPointInQuotationRule = rule.NewRule(regexp.MustCompile(`!(['"])`), "&ᓴ&$1")
	// Rubular: http://rubular.com/r/sl57YI8LkA
	exclamationPointBeforeCommaMidSentenceRule = rule.NewRule(regexp.MustCompile(`!(,\s[a-z])`), "&ᓴ&$1")
	// Rubular: http://rubular.com/r/f9zTjmkIPb
	exclamationPointMidSentenceRule = rule.NewRule(regexp.MustCompile(`!(\s[a-z])`), "&ᓴ&$1")
)

var (
	// added special case: r"[。．.！!? ]{2,}" to handle intermittent dots, exclamation, etc.
	// r"[。．.！!?] at end to handle single instances of these symbol inputs
	sentenceBoundaryRule1 = rule.NewRule(regexp.MustCompile(`(（([^）])*）)\s?([A-Z])`), "$1\r$3")
	sentenceBoundaryRule2 = rule.NewRule(regexp.MustCompile(`(「([^」])*」)\s([A-Z])`), "$1\r$3")
	sentenceBoundaryRule3 = rule.NewRule(regexp.MustCompile(`(\(([^)]){2,}\))\s([A-Z])`), "$1\r$3")
	sentenceBoundaryRule4 = rule.NewRule(regexp.MustCompile(`('([^'])*[^,]')\s([A-Z])`), "$1\r$3")
	sentenceBoundaryRule5 = rule.NewRule(regexp.MustCompile(`("([^"])*[^,]")\s([A-Z])`), "$1\r$3")
	sentenceBoundaryRule6 = rule.NewRule(regexp.MustCompile(`(([^”])*[^,]”)\s([A-Z])`), "$1\r$3")
	sentenceBoundaryRule7 = rule.NewRule(regexp.MustCompile(`(\S.*?[。．.！!?？ȸȹ☉☈☇☄])\s*(\S*?)`), "$1\r$2")
	sentenceBoundaryRule8 = rule.NewRule(regexp.MustCompile(`([。．.！!? ]{2,})`), "$1\r")
	sentenceBoundaryRule9 = rule.NewRule(regexp.MustCompile(`([。．.！!?？])`), "$1\r")
)
