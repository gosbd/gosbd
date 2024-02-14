package gosbd

import (
	"github.com/gosbd/gosbd/internal/cleaner"
	"github.com/gosbd/gosbd/internal/lang"
	"github.com/gosbd/gosbd/internal/processor"
	"github.com/gosbd/gosbd/internal/replacer"
	"github.com/gosbd/gosbd/internal/segmenter"
)

// Segmenter segments a text into an list of sentences.
type Segmenter interface {
	// Segment takes a string of text and returns a slice of sentences.
	Segment(text string) []string
	// TextSpans takes a string of text and returns a slice of TextSpan objects,
	// where each TextSpan represents a sentence and its position in the original text.
	TextSpans(text string) []segmenter.TextSpan
}

// Option is a type that represents a function that modifies the Options struct.
type Option func(*segmenter.Params)

// Clean cleans original tex,by default False
func Clean() Option {
	return func(params *segmenter.Params) {
		params.Cleaner = cleaner.NewCleaner()
	}
}

// NewSegmenter is a factory function that creates a new instance of a Segmenter.
// It takes a language code as input and uses it to configure the Segmenter with
// language-specific settings.
func NewSegmenter(langCode string, option ...Option) Segmenter {
	cfg := lang.Lang(langCode)
	punctuationReplacer := replacer.NewPunctuationReplacer()
	betweenPunctuationReplacer := cfg.BetweenPunctuationReplacer
	if betweenPunctuationReplacer == nil {
		betweenPunctuationReplacer = replacer.NewBetweenPunctuation(punctuationReplacer)
	}
	segmenterParams := &segmenter.Params{
		Config: cfg,
		Processor: processor.NewProcessor(processor.Params{
			Lang:                       cfg,
			ListItemReplacer:           replacer.NewListItemReplacer(),
			AbbrReplacer:               replacer.NewAbbreviationReplacer(cfg),
			PunctuationReplacer:        &punctuationReplacer,
			BetweenPunctuationReplacer: betweenPunctuationReplacer,
		}),
	}
	for _, opt := range option {
		opt(segmenterParams)
	}
	return segmenter.NewSegmenter(segmenterParams)
}
