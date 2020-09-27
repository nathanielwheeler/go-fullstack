package controllers

import (
	"github.com/nathanielwheeler/go-fullstack/api/models"
	"github.com/nathanielwheeler/go-fullstack/client/views"
)

// ValuesC will hold information about views and services.
type ValuesC struct {
	ValuesView *views.View
	vs        models.ValuesS
}

// NewValuesC is a constructor for ValuesC struct
func NewValuesC(vs models.ValuesS) *ValuesC {
	return &ValuesC{
		ValuesView: views.NewView("index", "values/values"),
		vs:        vs,
	}
}

// Values : GET /values
// Returns all values.
func (vc *ValuesC) Values(res http.ResponseWriter, req *http.Request) {
  values, err := vc.vs.GetAll()
  if err != nil {
    log.Println(err)
    return
  }

  var vd views.Data
  vd.Yield = values
  vc.ValuesView.Render(res, req, vd)
}