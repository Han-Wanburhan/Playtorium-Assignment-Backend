package cart

import (
	"errors"
	"fmt"
	"net/http"
	"shopping/go-api/orm"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type addtocartbody struct {
	CustomerID    int `json:"customer_id" validate:"require"`
	ProductID int `json:"product_id" validate:"require"`
	Quantity int `json:"quantity" validate:"require"`

}

func AddToCart(c *gin.Context){
	var json addtocartbody

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON", "error": err.Error()})
		return
	}


	// ไม่พบ Category ในฐานข้อมูล, จึงสร้าง Category ใหม่
	newAddCartItem := orm.CARTITEM{
		CustomerID: json.CustomerID,
		ProductID: json.ProductID,
		Quantity: json.Quantity,

	}

	result := orm.Db.Create(&newAddCartItem)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create Category", "error": result.Error.Error()})
		return
	}

	// สร้าง Category สำเร็จ, ส่งข้อมูล Category ที่ถูกสร้างไปด้วย
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Category created successfully", "category": newAddCartItem})

}

type couponfixedamountbody struct {
	Total   int `json:"total" validate:"require"`
	Amount int `json:"amount" validate:"require"`

}

func CouponFixedAmount(c *gin.Context) {
	var json couponfixedamountbody

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON", "error": err.Error()})
		return
	}

	if json.Total >= json.Amount && json.Amount > 1 && json.Total > 1{
		newTotal := json.Total - json.Amount
		c.JSON(http.StatusOK, gin.H{"status": "success", "newTotal": newTotal})
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Amount exceeds the total, cannot apply discount"})
	}
}

type couponpercentagediscountbody struct {
	Total   int `json:"total" validate:"require"`
	Percentage int `json:"precentage" validate:"require"`

}

func CouponPercentageDiscount(c *gin.Context) {
	var json couponpercentagediscountbody

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON", "error": err.Error()})
		return
	}

	if json.Total < 1 || json.Percentage < 1 || json.Percentage > 100 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Amount exceeds the total, cannot apply discount"})
		return
		} else {
		discount := json.Total * json.Percentage / 100
		newTotal := json.Total - discount
		c.JSON(http.StatusOK, gin.H{"status": "success", "newTotal": newTotal})
	}
}



type OntopPercentageDiscountByCategoryBody struct {
    CustomerID int `json:"customerid" validate:"required"`
    Percentage int `json:"precentage" validate:"required"`
    Total      int `json:"total" validate:"required"`
    Category   int `json:"category" validate:"required"`
}


func OntopPercentageDiscountByCategory(c *gin.Context) {
    var json OntopPercentageDiscountByCategoryBody

    if err := c.ShouldBindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON", "error": err.Error()})
        return
    }

	category := orm.CATEGORY{}
	result := orm.Db.First(&category, "id = ?", json.Category)

if result.Error != nil {
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Category does not exist"})
        return
    }
}


    // ตรวจสอบในฐานข้อมูลว่ามีข้อมูลที่มี customer_id เดียวกันหรือไม่
    var results []orm.CARTITEM

    if err := orm.Db.Preload("Customer").Preload("Product").Where("customer_id = ?", json.CustomerID).Find(&results).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database error", "error": err.Error()})
        return
    }

		if json.Total < 1 || json.Percentage < 1 || json.Percentage > 100 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Amount exceeds the total, cannot apply discount"})
			return
		} else {
		


    // Check if any records were found
    if len(results) > 0 {
		// for _, result := range results {
		// 	fmt.Println(result.Product.P_Price)
		// }
        // c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Records found", "results": results})
		// return
        		var productIDs []int
				var productPrice []int
		for _, result := range results {

			if result.Product.CategoryID == uint(json.Category) {

			productIDs = append(productIDs, result.ProductID)
			productPrice = append(productPrice, result.Product.P_Price)
		}

		
    }
			
	totalincate := sum(productPrice)
	totaldis := float64(totalincate) * float64(json.Percentage) / 100
	newTotal := float64(json.Total) - totaldis


	// c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Records found", "results":productIDs,"newtotal": newTotal})
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Records found","newtotal": newTotal})

		return
}	}

    // No records found, proceed with the rest of your logic here...
    c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "No records found for the specified customer_id"})
}

func sum(numbers []int) int {
    result := 0
    for _, num := range numbers {
        result += num
    }
    return result
}


type ontopdiscountbypointbody struct {
    CustomerID int `json:"customerid" validate:"required"`
    Total      int `json:"total" validate:"required"`
	Point int `json:"point" validate:"required"`
}


func OntopDiscountByPoint(c *gin.Context) {
    var json ontopdiscountbypointbody
		if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON", "error": err.Error()})
		return
	}
	category := orm.CUSTOMER{}
	result := orm.Db.First(&category, "id = ?", json.CustomerID)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "User does not exist"})
				return
			}
		}

	if json.Point < 1 || json.Total < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Point exceeds the limit, cannot apply discount"})
		return
	}


	limitdis := json.Total*20/100
	if json.Point <= limitdis {
		newTotal := json.Total-json.Point
		fmt.Println(newTotal)
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Records found","newtotal": newTotal})
		return
	}else{
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Point exceeds the limit, cannot apply discount"})
		return
	}


}

type seasonalspecialcampaignsbody struct {
    Total int `json:"total" validate:"required"`
    Every      int `json:"every" validate:"required"`
	Discount int `json:"discount" validate:"required"`
}


func SeasonalSpecialCampaigns(c *gin.Context) {
    var json seasonalspecialcampaignsbody
		if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON", "error": err.Error()})
		return
	}
	if json.Total > 0 && json.Every > 0 && json.Discount > 0 {
			countdis := json.Total / json.Every
			newTotal := json.Total - (countdis * json.Discount)
			fmt.Println(newTotal)
		c.JSON(http.StatusOK, gin.H{"status": "ok", "newtotal": newTotal})
		return
	}else{
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
		return
	}

}