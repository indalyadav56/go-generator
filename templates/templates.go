package templates

import "embed"

//go:embed *.tmpl **/*.tmpl
var TemplateFS embed.FS
