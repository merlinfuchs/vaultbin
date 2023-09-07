package views

import (
	"embed"
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

//go:embed *.html
var templates embed.FS

type Template struct {
	templates *template.Template
}

func New() *Template {
	t := &Template{
		templates: template.Must(template.ParseFS(templates, "*.html")),
	}

	return t
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
