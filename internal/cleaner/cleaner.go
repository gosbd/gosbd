package cleaner

type Cleaner struct {
}

func (c *Cleaner) Clean(text string) string {
	panic("not implemented")
}

func NewCleaner() *Cleaner {
	return &Cleaner{}
}
