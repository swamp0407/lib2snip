package entities

type Snippet struct {
	Name        string `json:"name"`
	Body        string `json:"body"`
	Prefix      string `json:"prefix"`
	Description string `json:"description"`
	Scope       string `json:"scope"`
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

func (s *Snippet) Output() string {
	return s.Name + `:{
		"prefix": "` + s.Prefix + `",
		"scope": "` + s.Scope + `"
		"body": "` + s.Body + `",
		"description": "` + s.Description + `",
	}`
}
