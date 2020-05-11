package util

import "github.com/jinzhu/gorm"

type (
	Module interface {
		Install(tx *gorm.DB) error
	}

	Migration struct {
		Module string `gorm:"index:module"`
		Name   string
	}
)
