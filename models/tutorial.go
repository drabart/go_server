package models

import "gorm.io/gorm"

type Tutorial struct {
	gorm.Model
	Name        string
	Description string
	URL         string
	Tags        []Tag `gorm:"many2many:tutorial_tags;"`
}
