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



func GetAllProduct(c *gin.Context) {
	var products []orm.PRODUCT

	// ดึงข้อมูลทั้งหมดจากตาราง PRODUCT
	result := orm.Db.Find(&products)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve products", "error": result.Error.Error()})
		return
	}

	for i, product := range products {
		var images []orm.IMAGE

		// ดึงข้อมูลจากตาราง IMAGE ที่มี productID เท่ากับ ID ของสินค้าในลูป
		imageResult := orm.Db.Where("product_id = ?", product.ID).First(&images)
		if imageResult.Error != nil {
			return
		}

		// Update the product with the image information
		if len(images) > 0 {
			selectedImage := images[0]
			products[i].Image = selectedImage.Image
		}
	}

	// ส่งข้อมูลทั้งหมดในรูปแบบ JSON
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Products retrieved successfully", "products": products})
}


func GetProductById(c *gin.Context) {
    // Get the product ID from the request parameters
    productID := c.Param("id")

    // Check if the product ID is provided
    if productID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Product ID is required"})
        return
    }

    var products []orm.PRODUCT

    // ดึงข้อมูลจากตาราง PRODUCT โดยใช้ productID
    result := orm.Db.Where("id = ?", productID).Find(&products)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve product", "error": result.Error.Error()})
        return
    }

    // ส่งข้อมูลในรูปแบบ JSON
    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Product retrieved successfully", "product": products})
}

func GetImageById(c *gin.Context) {
    // Get the product ID from the request parameters
    productID := c.Param("id")

    // Check if the product ID is provided
    if productID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Product ID is required"})
        return
    }

    var images []orm.IMAGE

    // ดึงข้อมูลจากตาราง PRODUCT โดยใช้ productID
    result := orm.Db.Where("product_id = ?", productID).Find(&images)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve product", "error": result.Error.Error()})
        return
    }

    // ส่งข้อมูลในรูปแบบ JSON
    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Product retrieved successfully", "product": images})
}

func GetCategoryById(c *gin.Context) {
    // Get the product ID from the request parameters
    CategoryID := c.Param("id")

    // Check if the product ID is provided
    if  CategoryID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Product ID is required"})
        return
    }

    var category orm.CATEGORY

    // ดึงข้อมูลจากตาราง PRODUCT โดยใช้ productID
    result := orm.Db.Where("id = ?",  CategoryID).Find(&category)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve product", "error": result.Error.Error()})
        return
    }

    // ส่งข้อมูลในรูปแบบ JSON
    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Product retrieved successfully", "product": category})
}
