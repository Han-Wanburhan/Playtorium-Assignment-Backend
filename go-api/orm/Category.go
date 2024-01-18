package orm

import "gorm.io/gorm"

type CATEGORY struct {
	gorm.Model
	C_Name string
}