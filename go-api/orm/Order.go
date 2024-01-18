package orm

import "gorm.io/gorm"

type ORDER struct {
	gorm.Model
	CustomerID uint     // Foreign key referencing Category's ID
	Customer   CUSTOMER `gorm:"foreignKey:CustomerID"`
}