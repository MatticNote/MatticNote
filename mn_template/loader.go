package mn_template

import (
	"embed"
	"text/template"
)

//go:embed *.django */*.django _mail/*.txt
var Templates embed.FS

func LoadTextTemplate() (*template.Template, error) {
	return template.ParseFS(Templates, "_mail/*.txt")
}
