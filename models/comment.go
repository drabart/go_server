package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	Contents   string
	UserID     uint
	User       User
	TutorialID uint
	Tutorial   Tutorial
	// comment might be a reply to another comment
	ParentID uint
	Parent   *Comment `gorm:"foreignkey:ParentID"`
}
