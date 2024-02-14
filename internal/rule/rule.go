package rule

import "regexp"

type Rule struct {
	pattern     *regexp.Regexp
	replacement string
}

type Rules []Rule

func NewRule(pattern *regexp.Regexp, replacement string) Rule {
	return Rule{
		pattern:     pattern,
		replacement: replacement,
	}
}

func (r Rule) Apply(text string) string {
	return r.pattern.ReplaceAllString(text, r.replacement)
}

func (r Rule) Pattern() *regexp.Regexp {
	return r.pattern
}

func (r Rules) Apply(text string) string {
	v := text
	for _, rr := range r {
		v = rr.Apply(v)
	}
	return v
}
