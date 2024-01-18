package orm

import "gorm.io/gorm"

type CUSTOMER struct {
	gorm.Model
	Fname string
	Lname string
	UserName string
	Mobile string
	Pass  string
}