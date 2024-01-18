package orm

type POINT struct {
	CustomerID uint     // Foreign key referencing Category's ID
	Customer   CUSTOMER `gorm:"foreignKey:CustomerID"`
	Point      int
}