package context

import (
	"html/template"

	dao "github.com/wyll-io/dicomizer/internal/DAO"
	awsStorage "github.com/wyll-io/dicomizer/internal/storage"
)

type InternalValues struct {
	DB     dao.DBActions
	S3     awsStorage.StorageAction
	Center string
}
type TemplatesValues map[string]*template.Template

type HTTP uint8

const (
	User HTTP = iota
	Internal
	Templates
)
