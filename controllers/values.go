package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/nathanielwheeler/go-fullstack/context"
	"github.com/nathanielwheeler/go-fullstack/models"
	"github.com/nathanielwheeler/go-fullstack/views"

	"github.com/gorilla/mux"
)

// Named routes.
const (
	BlogIndexRoute = "blog_index"
	BlogValueRoute = "blog_value"
	EditValue      = "edit_value"
)

const (
	maxMultipartMem = 1 << 20 // 1 megabyte
)

// Values will hold information about views and services
type Values struct {
	New           *views.View
	EditView      *views.View
	BlogValueView *views.View
	BlogIndexView *views.View
	vs            models.ValuesService
	r             *mux.Router
}

// NewValues is a constructor for Values struct
func NewValues(vs models.ValuesService, r *mux.Router) *Values {
	return &Values{
		New:           views.NewView("app", "values/new"),
		EditView:      views.NewView("app", "values/edit"),
		BlogValueView: views.NewView("app", "values/blog/value"),
		BlogIndexView: views.NewView("app", "values/blog/index"),
		vs:            vs,
		r:             r,
	}
}

// ValueForm will hold information for creating a new value
type ValueForm struct {
	Name string `schema:"name"`
}

// GetValue : GET /values/:id
func (p *Values) GetValue(res http.ResponseWriter, req *http.Request) {
	value, err := p.valueByID(res, req)
	if err != nil {
		// valueByID already renders error
		return
	}
	var vd views.Data
	vd.Yield = value
	p.BlogValueView.Render(res, req, vd)
}

// GetAllValues : GET /values
func (p *Values) GetAllValues(res http.ResponseWriter, req *http.Request) {
	values, err := p.vs.GetAll()
	if err != nil {
		log.Println(err)
		http.Error(res, "Something bad happened.", http.StatusInternalServerError)
		return
	}
	var vd views.Data
	vd.Yield = values
	p.BlogIndexView.Render(res, req, vd)
}

// Create : VALUE /values
func (p *Values) Create(res http.ResponseWriter, req *http.Request) {
	var vd views.Data
	var form ValueForm
	if err := parseForm(req, &form); err != nil {
		vd.SetAlert(err)
		p.New.Render(res, req, vd)
		return
	}
	user := context.User(req.Context())
	if user.IsAdmin != true {
		http.Error(res, "You do not have permission to create a value", http.StatusForbidden)
		return
	}
	value := models.Value{
		Name: form.Name,
	}
	if err := p.vs.Create(&value); err != nil {
		vd.SetAlert(err)
		p.New.Render(res, req, vd)
		return
	}
	url, err := p.r.Get(EditValue).URL("id", fmt.Sprintf("%v", value.ID))
	if err != nil {
		log.Println(err)
		http.Redirect(res, req, "/blog", http.StatusFound)
		return
	}
	http.Redirect(res, req, url.Path, http.StatusFound)
}

// Edit : VALUE /values/:id/update
func (p *Values) Edit(res http.ResponseWriter, req *http.Request) {
	value, err := p.valueByID(res, req)
	if err != nil {
		// error handled by valueByID
		return
	}
	user := context.User(req.Context())
	if user.IsAdmin != true {
		http.Error(res, "You do not have permission to edit this value", http.StatusForbidden)
		return
	}
	var vd views.Data
	vd.Yield = value
	p.EditView.Render(res, req, vd)
}

// Update : VALUE /values/:id/update
/*  - This does NOT update the path of the value. */
func (p *Values) Update(res http.ResponseWriter, req *http.Request) {
	value, err := p.valueByID(res, req)
	if err != nil {
		// implemented by valueByID
		return
	}
	user := context.User(req.Context())
	if user.IsAdmin != true {
		http.Error(res, "You do not have permission to edit this value", http.StatusForbidden)
		return
	}

	var vd views.Data
	vd.Yield = value
	var form ValueForm
	if err := parseForm(req, &form); err != nil {
		vd.SetAlert(err)
		p.EditView.Render(res, req, vd)
		return
	}
	value.Name = form.Name
	err = p.vs.Update(value)
	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Value updated successfully!",
		}
	}
	p.EditView.Render(res, req, vd)
}

// Delete : VALUE /values/:id/delete
func (p *Values) Delete(res http.ResponseWriter, req *http.Request) {
	value, err := p.valueByID(res, req)
	if err != nil {
		// valueByID renders error
		return
	}
	user := context.User(req.Context())
	if user.IsAdmin != true {
		http.Error(res, "You do not have permission to edit this value", http.StatusForbidden)
		return
	}
	var vd views.Data
	err = p.vs.Delete(value.ID)
	if err != nil {
		vd.SetAlert(err)
		vd.Yield = value
		p.EditView.Render(res, req, vd)
		return
	}
	url, err := p.r.Get(BlogIndexRoute).URL()
	if err != nil {
		http.Redirect(res, req, "/", http.StatusFound)
		return
	}
	http.Redirect(res, req, url.Path, http.StatusFound)
}

func (p *Values) valueByID(res http.ResponseWriter, req *http.Request) (*models.Value, error) {
	idVar := mux.Vars(req)["id"]
	id, err := strconv.Atoi(idVar)
	if err != nil {
		log.Println(err)
		http.Error(res, "Invalid value ID", http.StatusNotFound)
		return nil, err
	}
	value, err := p.vs.Get(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(res, "Value not found", http.StatusNotFound)
		default:
			log.Println(err)
			http.Error(res, "Something bad happened.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return value, nil
}
