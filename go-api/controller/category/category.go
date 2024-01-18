package category

import (
	"net/http"
	"shopping/go-api/orm"

	"github.com/gin-gonic/gin"
)

type addCategoryBody struct {
	C_name string `json:"c_name" validate:"required"`
}

func AddCategory(c *gin.Context) {
	var json addCategoryBody

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON", "error": err.Error()})
		return
	}

	// ลองค้นหา Category จากชื่อในฐานข้อมูล
	existingCategory := orm.CATEGORY{}
	result := orm.Db.First(&existingCategory, "c_name = ?", json.C_name)
	if result.Error == nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": "Category already exists"})
		return
	}

	// ไม่พบ Category ในฐานข้อมูล, จึงสร้าง Category ใหม่
	newCategory := orm.CATEGORY{
		C_Name: json.C_name,
	}

	result = orm.Db.Create(&newCategory)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create Category", "error": result.Error.Error()})
		return
	}

	// สร้าง Category สำเร็จ, ส่งข้อมูล Category ที่ถูกสร้างไปด้วย
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Category created successfully", "category": newCategory})
}
