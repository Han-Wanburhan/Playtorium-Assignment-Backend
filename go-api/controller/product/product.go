package product

import (
	"net/http"
	"shopping/go-api/orm"

	"github.com/gin-gonic/gin"
)

type addproductbody struct {
	P_Name string `json:"p-name" validate:"require"`
	P_Price int `json:"p-price" validate:"require"`
	P_Detail string `json:"p-detail" validate:"require"`
	P_In_Stock int `json:"p-instock" validate:"require"`
	CategoryID uint `json:"categoryid" validate:"require"`
}

func AddProduct(c *gin.Context){
var json addproductbody

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON", "error": err.Error()})
		return
	}


	// ไม่พบ Category ในฐานข้อมูล, จึงสร้าง Category ใหม่
	newProduct := orm.PRODUCT{
		P_Name: json.P_Name,
		P_Price: json.P_Price,
		P_Detail: json.P_Detail,
		P_In_Stock: json.P_In_Stock,
		CategoryID: json.CategoryID,
	}

	result := orm.Db.Create(&newProduct)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create Category", "error": result.Error.Error()})
		return
	}

	// สร้าง Category สำเร็จ, ส่งข้อมูล Category ที่ถูกสร้างไปด้วย
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Category created successfully", "category": newProduct})

}