package agent

// DocsAgent updates README/docs.
type DocsAgent struct{}

// SuggestDocUpdate returns a suggested doc change for the given context.
func (d *DocsAgent) SuggestDocUpdate(context string) (string, error) {
	return "", nil
}
