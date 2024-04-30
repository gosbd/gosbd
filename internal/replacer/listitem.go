package replacer

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/gosbd/gosbd/internal/processor"
	"github.com/gosbd/gosbd/internal/rule"
)

type ListItemReplacer struct {
}

var (
	// TODO: Make this configurable (e.g., 一, 二, 三, 四, 五, 六, 七, 八, 九, 十)
	isRomanNumerals = map[string]alphabet{}
	isLatinNumerals = map[string]alphabet{}
)

var (
	substitudePeriodRule = rule.NewRule(regexp.MustCompile(`♨`), "∯")
	listMarkerRule       = rule.NewRule(regexp.MustCompile(`☝`), "")
)

var (
	reSpaceBetweenListItems1 = regexp.MustCompile(`\S\S\s\S\s*\d+♨`)
	reSpaceBetweenListItems2 = regexp.MustCompile(`\S\S\s\d{1,2}♨`)
	reSpaceBetweenListItems3 = regexp.MustCompile(`\S\S\s\d{1,2}☝`)
)

var (
	// Rubular: http://rubular.com/r/GcnmQt4a3I
	romanNumeralsInParentheses = regexp.MustCompile(`\((([mdclxvi])m*(c[md]|d?c*)(x[cl]|l?x*)(i[xv]|v?i*))\)(\s[A-Z])`)
)

const (
	whiteSpace = "\t\n\f\r "
)

type alphabet struct {
	val      string
	idx      int
	alphabet bool
}

func init() {
	// TODO: Make this configurable (e.g., 一, 二, 三, 四, 五, 六, 七, 八, 九, 十)
	for i, s := range strings.Split("i ii iii iv v vi vii viii ix x xi xii xiii xiv x xi xii xiii xv xvi xvii xviii xix xx", " ") {
		isRomanNumerals[s] = alphabet{s, i, true}
	}
	for i, s := range strings.Split("abcdefghijklmnopqrstuvwxyz", "") {
		isLatinNumerals[s] = alphabet{s, i, true}
	}
}

func (l ListItemReplacer) AddLineBreak(text string) string {
	text = l.formatAlphabeticalLists(text)
	text = l.formatRomanNumeralLists(text)
	text = l.formatNumberedListWithPeriods(text)
	text = l.formatNumberedListWithParens(text)
	return text
}

func (l ListItemReplacer) ReplaceParens(text string) string {
	return romanNumeralsInParentheses.ReplaceAllString(text, `&✂&$1&⌬&$2`)
}

func (l ListItemReplacer) formatAlphabeticalLists(text string) string {
	text = l.addLineBreaksForAlphabeticalListWithPeriods(text, false)
	text = l.addLineBreaksForAlphabeticalListWithParens(text, false)
	return text
}

func (l ListItemReplacer) formatRomanNumeralLists(text string) string {
	text = l.addLineBreaksForAlphabeticalListWithPeriods(text, true)
	text = l.addLineBreaksForAlphabeticalListWithParens(text, true)
	return text
}

func (l ListItemReplacer) addLineBreaksForAlphabeticalListWithPeriods(text string, romanNumeral bool) string {
	return l.iterateAlphabetArray(text, romanNumeral, false)
}

func (l ListItemReplacer) addLineBreaksForAlphabeticalListWithParens(text string, romanNumeral bool) string {
	return l.iterateAlphabetArray(text, romanNumeral, true)
}

func (l ListItemReplacer) formatNumberedListWithPeriods(text string) string {
	text = l.replacePeriodsInNumberedList(text)
	text = l.addLineBreaksForNumberedListWithPeriods(text)
	text = substitudePeriodRule.Apply(text)
	return text
}

func (l ListItemReplacer) formatNumberedListWithParens(text string) string {
	text = l.replaceParensInNumberedList(text)
	text = l.addLineBreaksForNumberedListWithParens(text)
	text = listMarkerRule.Apply(text)
	return text
}

func (l ListItemReplacer) addLineBreaksForNumberedListWithParens(text string) string {
	if strings.Contains(text, "☝") {
		matched, _ := regexp.MatchString("☝.+\n.+☝|☝.+\r.+☝", text)
		if !matched {
			text = l.spaceBetweenListItems(text, reSpaceBetweenListItems3)
		}
	}
	return text
}

func (l ListItemReplacer) addLineBreaksForNumberedListWithPeriods(text string) string {
	if strings.Contains(text, "♨") {
		matched, _ := regexp.MatchString(`♨.+([\n\r]).+♨`, text)
		matchedFor, _ := regexp.MatchString(`for\s\d{1,2}♨\s[a-z]`, text)
		if !matched && !matchedFor {
			text = l.spaceBetweenListItems(text, reSpaceBetweenListItems1)
			text = l.spaceBetweenListItems(text, reSpaceBetweenListItems2)
		}
	}
	return text
}

func (l ListItemReplacer) spaceBetweenListItems(text string, pattern *regexp.Regexp) string {
	replaceFunc := func(match string) string {
		leftPart := []rune(match)[:2]
		chompedMatch := strings.TrimSpace(strings.TrimPrefix(match, string(leftPart)))
		return string(leftPart) + "\r" + chompedMatch
	}
	return pattern.ReplaceAllStringFunc(text, replaceFunc)
}

func (l ListItemReplacer) replacePeriodsInNumberedList(text string) string {
	pattern := regexp.MustCompile(`\s\d{1,2}\.\s|^\d{1,2}\.\s|\s\d{1,2}\.\)|^\d{1,2}\.\)|\s-\d{1,2}\.\s|^-\d{1,2}\.\s|\s⁃\d{1,2}\.\s|^⁃\d{1,2}\.\s|\s-\d{1,2}\.\)|^-\d{1,2}\.\)|\s⁃\d{1,2}\.\)|^⁃\d{1,2}\.\)`)
	return l.scanLists(text, pattern, "♨", false)
}

func (l ListItemReplacer) replaceParensInNumberedList(text string) string {
	pattern := regexp.MustCompile(`\d{1,2}\)\s`)
	text = l.scanLists(text, pattern, "☝", true)
	text = l.scanLists(text, pattern, "☝", true)
	return text
}

func (l ListItemReplacer) scanLists(text string, re *regexp.Regexp, replacement string, parens bool) string {
	values := l.extractNumberList(text, re)
	for i, val := range values {
		if i < len(values)-1 && val+1 == values[i+1] {
			text = l.substituteFoundListItems(text, re, val, replacement, parens)
		} else if i > 0 {
			if val-1 == values[i-1] || (val == 0 && values[i-1] == 9) || (val == 9 && values[i-1] == 0) {
				text = l.substituteFoundListItems(text, re, val, replacement, parens)
			}
		}
	}
	return text
}

func (l ListItemReplacer) substituteFoundListItems(text string, re *regexp.Regexp, val int, replacement string, parens bool) string {
	replaceItem := func(match string, repl string) string {
		i := strings.Index(match, strconv.Itoa(val))
		if i == -1 {
			return match
		}
		// extract left whitespaces
		chompedMatch := strings.TrimLeft(match, whiteSpace)
		leadingWS := match[:len(match)-len(chompedMatch)]

		// extract right whitespaces
		chompedMatch = strings.TrimSpace(match)
		trailingWS := match[len(leadingWS)+len(chompedMatch):]

		// extract left of the content
		j := strings.Index(chompedMatch, strconv.Itoa(val))
		left := chompedMatch[:j]

		// extract right of the content
		right := strings.TrimLeft(chompedMatch[j+len(strconv.Itoa(val)):], ".")

		// check if the value is correct for this iteration
		if strconv.Itoa(val) == strings.TrimRight(chompedMatch[j:], ".)] ") {
			return fmt.Sprintf("%s%s%d%s%s%s", leadingWS, left, val, repl, right, trailingWS)
		}
		return match
	}

	return re.ReplaceAllStringFunc(text, func(match string) string {
		return replaceItem(match, replacement)
	})
}

func (l ListItemReplacer) extractAlphabeticalListLettersWithParens(text string, romanNumeral bool) ([]string, map[string]alphabet) {
	isAlphabet := isRomanNumerals
	if !romanNumeral {
		isAlphabet = isLatinNumerals
	}
	pattern := regexp.MustCompile(`(\([a-z]+\)?)|^([a-z]+\))|\A([a-z]+\))|(\s[a-z]+\))`)
	var list []string
	for _, match := range pattern.FindAllString(text, -1) {
		trimmedMatch := strings.Trim(match, " ()")
		if isAlphabet[trimmedMatch].alphabet {
			list = append(list, trimmedMatch)
		}
	}
	return list, isAlphabet
}

func (l ListItemReplacer) extractAlphabeticalListLettersWithPeriods(text string, romanNumeral bool) ([]string, map[string]alphabet) {
	isAlphabet := isRomanNumerals
	if !romanNumeral {
		isAlphabet = isLatinNumerals
	}
	pattern := regexp.MustCompile(`([a-z]+\.)|(\A[a-z]+\.)|(\s[a-z]+\.)`)
	var list []string
	for _, match := range pattern.FindAllString(text, -1) {
		trimmedMatch := strings.Trim(match, " .")
		if isAlphabet[trimmedMatch].alphabet {
			list = append(list, trimmedMatch)
		}
	}
	return list, isAlphabet
}

func (l ListItemReplacer) extractNumberList(text string, re *regexp.Regexp) []int {
	var ns []int
	for _, s := range re.FindAllString(text, -1) {
		ss := strings.Trim(s, " .)⁃-")
		if n, err := strconv.Atoi(ss); err == nil {
			ns = append(ns, n)
		}
	}
	return ns
}

func (l ListItemReplacer) iterateAlphabetArray(text string, romanNumeral bool, parens bool) string {
	var (
		list       []string
		isAlphabet map[string]alphabet
	)
	if parens {
		list, isAlphabet = l.extractAlphabeticalListLettersWithParens(text, romanNumeral)
	} else {
		list, isAlphabet = l.extractAlphabeticalListLettersWithPeriods(text, romanNumeral)
	}
	for i, ss := range list {
		if i == len(list)-1 {
			text = l.lastArrayItemReplacement(text, ss, i, isAlphabet, list, parens)
		} else {
			text = l.otherItemsReplacement(text, ss, i, isAlphabet, list, parens)
		}
	}
	return text
}

func (l ListItemReplacer) otherItemsReplacement(text string, val string, i int, isAlphabet map[string]alphabet, list []string, parens bool) string {
	if len(isAlphabet) == 0 || len(list) == 0 {
		return text
	}
	if !isAlphabet[val].alphabet {
		return text
	}
	if i > 0 {
		if !isAlphabet[list[i-1]].alphabet || math.Abs(float64(isAlphabet[list[i-1]].idx-isAlphabet[val].idx)) != 1 {
			return text
		}
	}
	if !isAlphabet[list[i+1]].alphabet || math.Abs(float64(isAlphabet[list[i+1]].idx-isAlphabet[val].idx)) != 1 {
		return text
	}
	return l.replaceCorrectAlphabetList(text, val, parens)
}

func (l ListItemReplacer) lastArrayItemReplacement(text string, val string, i int, isAlphabet map[string]alphabet, list []string, parens bool) string {
	if len(isAlphabet) == 0 && len(list) == 0 {
		return text
	}
	if i == 0 || !isAlphabet[list[i-1]].alphabet || !isAlphabet[val].alphabet {
		return text
	}
	if math.Abs(float64(isAlphabet[list[i-1]].idx-isAlphabet[val].idx)) != 1 {
		return text
	}
	return l.replaceCorrectAlphabetList(text, val, parens)
}

func (l ListItemReplacer) replaceCorrectAlphabetList(text string, val string, parens bool) string {
	if parens {
		text = l.replaceAlphabetListParens(text, val)
	} else {
		text = l.replaceAlphabetListPeriod(text, val)
	}
	return text
}

var (
	replaceAlphabetListParensRegex = regexp.MustCompile(`(?i)(\([a-z]+\))|^[a-z]+\)|\A[a-z]+\)|\s[a-z]+\)`)
	replaceAlphabetListPeriodRegex = regexp.MustCompile(`(?i)([a-z])\.`)
)

// Input: "a) ffegnog (b) fgegkl c)"
// Output: "\ra) ffegnog \r&✂&b) fgegkl \rc)"
func (l ListItemReplacer) replaceAlphabetListParens(text string, val string) string {
	replaceFunc := func(match string) string {
		if strings.HasPrefix(match, "(") {
			matchWoParen := strings.TrimPrefix(match, "(")
			if matchWoParen == val+")" {
				return "\r&✂&" + matchWoParen
			}
		} else {
			if idx := strings.Index(match, val); idx != -1 && match[idx:] == val+")" {
				return match[:idx] + "\r" + match[idx:]
			}
		}
		return match
	}

	// Use the compiled regex to match and replace in the text
	return replaceAlphabetListParensRegex.ReplaceAllStringFunc(text, func(match string) string {
		return replaceFunc(match)
	})
}

// Input: 'a. ffegnog b. fgegkl c.'
// Output: \ra∯ ffegnog \rb∯ fgegkl \rc∯
func (l ListItemReplacer) replaceAlphabetListPeriod(text string, val string) string {
	replaceFunc := func(match string) string {
		if i := strings.Index(match, val); i != -1 && match[i:] == val+"." {
			s := match[:i] + "\r" + match[i:i+1] + "∯"
			return s
		}
		return match
	}

	// Use the compiled regex to match and replace in the text
	return replaceAlphabetListPeriodRegex.ReplaceAllStringFunc(text, func(match string) string {
		return replaceFunc(match)
	})
}

var _ processor.ListItemReplacer = (*ListItemReplacer)(nil)

func NewListItemReplacer() *ListItemReplacer {
	return &ListItemReplacer{}
}
