package lang

import (
	"fmt"

	"github.com/gosbd/gosbd/internal/processor"
)

var (
	langMap = map[string]*processor.Config{
		"en": processor.Standard(),
		"zh": newChinese(),
		"ja": newJapanese(),
		"ru": newRussian(),
	}
)

func Lang(lang string) *processor.Config {
	if _, ok := langMap[lang]; !ok {
		availableLangs := make([]string, 0, len(langMap))
		for k := range langMap {
			availableLangs = append(availableLangs, k)
		}
		panic(fmt.Errorf("provide valid language ID i.e. ISO code. Available codes are : %v", availableLangs))
	}
	return langMap[lang]
}
