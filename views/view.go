package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var layoutDir = "static/layouts"

func NewView(layout string, files ...string) *View {
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

type View struct {
	Template *template.Template
	Layout   string
}

func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}

func layoutFiles() []string {
	files, err := filepath.Glob(layoutDir + "/*.gohtml")
	if err != nil {
		panic(err)
	}
	return files
}