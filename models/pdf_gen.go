package models

import (
	"github.com/goravel/framework/database/orm"
)

type PdfGen struct {
	orm.Model
	Name   string  `gorm:"column:name" json:"name"`
	Params JsonMap `gorm:"column:params;type:json" json:"params"`
	Html   string  `gorm:"column:html" json:"html"`
}
type JsonMap string
