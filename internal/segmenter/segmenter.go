package segmenter

import (
	"fmt"
	"regexp"

	"github.com/yohamta/gosbd/internal/processor"
)

type Segmenter struct {
	cfg       *processor.Config
	processor Processor
	cleaner   Cleaner
}

type Processor interface {
	Process(text string) []string
}

type Cleaner interface {
	Clean(text string) string
}

type TextSpan struct {
	Start    int
	End      int
	Sentence string
}

func (sg *Segmenter) Segment(text string) []string {
	if len(text) == 0 {
		return nil
	}
	if sg.cleaner != nil {
		text = sg.cleaner.Clean(text)
	}
	postProcessedSents := sg.processor.Process(text)
	if sg.cleaner != nil {
		return postProcessedSents
	}
	spans := sg.sentencesWithCharSpans(postProcessedSents, text)
	var sentences []string
	for _, span := range spans {
		sentences = append(sentences, span.Sentence)
	}
	return sentences
}

func (sg *Segmenter) TextSpans(text string) []TextSpan {
	if len(text) == 0 {
		return nil
	}
	postProcessedSents := sg.processor.Process(text)
	return sg.sentencesWithCharSpans(postProcessedSents, text)
}

func (sg *Segmenter) sentencesWithCharSpans(sentences []string, original string) []TextSpan {
	var spans []TextSpan
	priorEndCharIdx := 0
	for _, sent := range sentences {
		re := regexp.MustCompile(fmt.Sprintf(`%s\s*`, regexp.QuoteMeta(sent)))
		for _, match := range re.FindAllStringIndex(original, -1) {
			matchStartIdx, matchEndIdx := match[0], match[1]
			if matchEndIdx > priorEndCharIdx {
				// making sure if curren sentence and its span
				// is either first sentence along with its char spans
				// or current sent spans adjacent to prior sentence spans
				spans = append(spans, TextSpan{
					Start:    matchStartIdx,
					End:      matchEndIdx,
					Sentence: sent,
				})
				priorEndCharIdx = matchEndIdx
				break
			}
		}
	}
	return spans
}

type Params struct {
	Config    *processor.Config
	Processor Processor
	Cleaner   Cleaner
}

func NewSegmenter(params *Params) *Segmenter {
	return &Segmenter{
		cfg:       params.Config,
		processor: params.Processor,
		cleaner:   params.Cleaner,
	}
}
