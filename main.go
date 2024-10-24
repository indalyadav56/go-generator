package main

import (
	"embed"

	"github.com/indalyadav56/go-generator/cmd"
	"github.com/indalyadav56/go-generator/templates"
)

//go:embed templates/*.tmpl templates/**/*.tmpl
var templateFS embed.FS

func main() {
	templates.TemplateFS = templateFS
	cmd.Execute()
}
