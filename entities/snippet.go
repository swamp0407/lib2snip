package entities

type Snippet struct {
	Name        string      `json:"-"`
	Body        string      `json:"body"`
	Prefix      interface{} `json:"prefix"`
	Description string      `json:"description"`
	Scope       string      `json:"scope"`
}

func NewSnippet(name string, body string, prefix string, description string, scope string) *Snippet {
	return &Snippet{
		Name:        name,
		Body:        body,
		Prefix:      prefix,
		Description: description,
		Scope:       scope,
	}
}
