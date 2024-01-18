package orm

type CARTITEM struct {
	CustomerID int      // Foreign key referencing Category's ID
	Customer   CUSTOMER `gorm:"foreignKey:CustomerID"`
	ProductID  int      // Foreign key referencing Category's ID
	Product    PRODUCT  `gorm:"foreignKey:ProductID"`
	Quantity   int
}