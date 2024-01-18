package orm

import "gorm.io/gorm"

type PRODUCT struct {
	gorm.Model
	P_Name string
	P_Price int
	P_Detail string
	P_In_Stock int
	CategoryID uint // Foreign key referencing Category's ID
	Category   CATEGORY `gorm:"foreignKey:CategoryID"`
	Image string `json:"image"`

	
}
