package views

import (
	"errors"
	"html/template"
	"net/url"
	"path/filepath"

  "nathanielwheeler.com/models"
  
	"github.com/gorilla/csrf"
)

var (
	templateDir  string = "views/"
	templateExt  string = ".html"
)

// View contains a pointer to a template and the name of a layout.
type View struct {
	Template *template.Template
	Layout   string
}

// NewView takes in a layout name, any number of filename strings, parses them into template, and returns the address of the new view.
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, dirFiles("layouts")...)
	files = append(files, dirFiles("components")...)
	t, err := template.
		New("").
		ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// Render is responsible for rendering the view.  Checks the underlying type of data passed into it.  Then checks cookie for alerts, looks up the user, creates a CSRF field with the request data, and then executes the template.
func (v *View) Render(res http.ResponseWriter, req *http.Request, data interface{}) {
	res.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		// Done so I can access the data in a var with type Data.
		vd = d
	default:
		// If data is NOT of type Data, make one and set the data to the Yield field.
		vd = Data{
			Yield: data,
		}
	}

	// Create CSRF field using current http request and add it onto the template FuncMap.
  csrfField := csrf.TemplateField(req)
  
	tpl := v.Template

	err := tpl.ExecuteTemplate(&buf, v.Layout, vd)
	if err != nil {
		http.Error(res, `Something went wrong, please try again.  If the problem persists, please contact me directly at "nathan@mailftp.com"`, http.StatusInternalServerError)
		return
	}
	io.Copy(res, &buf)
}

// ServeHTTP renders and serves views.
func (v *View) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	v.Render(res, req, nil)
}

func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = templateDir + f
	}
}

func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + templateExt
	}
}

func dirFiles(dir string) []string {
	files, err := filepath.Glob(templateDir + dir + templateExt)
	if err != nil {
		panic(err)
	}
	return files
}
