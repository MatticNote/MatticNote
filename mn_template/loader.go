package mn_template

import (
	"embed"
	hTemplate "html/template"
	tTemplate "text/template"
)

//go:embed *.django */*.django _mail/*.txt _oauth/*.html
var Templates embed.FS

func LoadTextTemplate() (*tTemplate.Template, error) {
	return tTemplate.ParseFS(Templates, "_mail/*.txt")
}

func LoadOAuthTemplate() (*hTemplate.Template, error) {
	return hTemplate.ParseFS(Templates, "_oauth/*.html")
}
