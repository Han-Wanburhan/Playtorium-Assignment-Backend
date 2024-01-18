package orm

import "gorm.io/gorm"

type IMAGE struct {
	gorm.Model
	Image string
	ProductID uint // Foreign key referencing Category's ID
	Product  PRODUCT `gorm:"foreignKey:ProductID"`

}