package lang

import (
	"fmt"
	"github.com/gosbd/gosbd/internal/processor"
	"regexp"
	"strings"
)

func newRussian() *processor.Config {
	cfg := processor.Standard()
	cfg.Abbreviation.Abbreviations = []string{"y", "y.e", "а", "авт", "адм.-терр", "акад", "в", "вв", "вкз", "вост.-европ", "г", "гг", "гос", "гр", "д", "деп", "дисс", "дол", "долл", "ежедн", "ж", "жен", "з", "зап", "зап.-европ", "заруб", "и", "ин", "иностр", "инст", "к", "канд", "кв", "кг", "куб", "л", "л.h", "л.н", "м", "мин", "моск", "муж", "н", "нед", "о", "п", "пгт", "пер", "пп", "пр", "просп", "проф", "р", "руб", "с", "сек", "см", "спб", "стр", "т", "тел", "тов", "тт", "тыс", "у", "у.е", "ул", "ф", "ч"}
	cfg.Abbreviation.PrePositiveAbbreviations = nil
	cfg.Abbreviation.NumberAbbreviations = nil
	cfg.Abbreviation.ReplacePeriodOfAbbrFn = replacePeriodOfAbbrRu
	cfg.SentenceStarters = nil
	return cfg
}

func replacePeriodOfAbbrRu(text, abbr string) string {
	var (
		re1 = regexp.MustCompile(fmt.Sprintf(`(\s%s)\.`, strings.TrimSpace(abbr)))
		re2 = regexp.MustCompile(fmt.Sprintf(`(\A%s)\.`, strings.TrimSpace(abbr)))
		re3 = regexp.MustCompile(fmt.Sprintf(`(^%s)\.`, strings.TrimSpace(abbr)))
	)
	text = re1.ReplaceAllString(text, "$1∯")
	text = re2.ReplaceAllString(text, "$1∯")
	text = re3.ReplaceAllString(text, "$1∯")
	return text
}
