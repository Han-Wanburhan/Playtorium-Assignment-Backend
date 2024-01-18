package orm

type DETAILORDER struct {
	OrderID   uint    // Foreign key referencing Category's ID
	Order     ORDER   `gorm:"foreignKey:OrderID"`
	ProductID uint    // Foreign key referencing Category's ID
	Product   PRODUCT `gorm:"foreignKey:ProductID"`
	Quantity  int
}