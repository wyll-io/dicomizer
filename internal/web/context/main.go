package context

import (
	"html/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	dao "github.com/wyll-io/dicomizer/internal/DAO"
)

type InternalValues struct {
	Cfg aws.Config
	DB  dao.DBActions
}
type TemplatesValues map[string]*template.Template

type HTTP uint8

const (
	User HTTP = iota
	Internal
	Templates
)
