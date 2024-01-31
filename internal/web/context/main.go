package context

import (
	"html/template"

	dao "github.com/wyll-io/dicomizer/internal/DAO"
)

type InternalValues struct {
	DB dao.DBActions
}
type TemplatesValues map[string]*template.Template

type HTTP uint8

const (
	User HTTP = iota
	Internal
	Templates
)
