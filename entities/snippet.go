package entities

type Snippet struct {
	Name        string
	Body        string
	Prefix      string
	Description string
	Scope       string
}

type VsSnippet struct {
	Snippet
}

type SnippetInterface interface {
	Output() string
}

func NewVS_Snippet(name string, body string, prefix string, description string, scope string) *VsSnippet {

	body = body + ``
	return &VsSnippet{
		Snippet{
			Name:        name,
			Body:        body,
			Prefix:      prefix,
			Description: description,
			Scope:       scope,
		},
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
