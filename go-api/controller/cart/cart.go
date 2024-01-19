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
	CustomerID int `json:"customer_id" validate:"required"`
	ProductID  int `json:"product_id" validate:"required"`
	Quantity   int `json:"quantity" validate:"required"`
}


func AddToCart(c *gin.Context) {
	var json addtocartbody

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid JSON", "error": err.Error()})
		return
	}

	existingCartItem := orm.CARTITEM{}
	result := orm.Db.Where("customer_id = ? AND product_id = ?", json.CustomerID, json.ProductID).First(&existingCartItem)
if result.Error == nil { 
	existingCartItem.Quantity = json.Quantity
	result := orm.Db.Model(&existingCartItem).Where("customer_id = ? AND product_id = ?", json.CustomerID, json.ProductID).Update("quantity", json.Quantity)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update Cart Item", "error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Cart Item updated successfully", "cart_item": existingCartItem})
	return
}

	newAddCartItem := orm.CARTITEM{
		CustomerID: json.CustomerID,
		ProductID:  json.ProductID,
		Quantity:   json.Quantity,
	}

	result = orm.Db.Create(&newAddCartItem)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create Cart Item", "error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Cart Item created successfully", "cart_item": newAddCartItem})
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
		// newTotal := json.Total - json.Amount
		newTotal, _ := DisFixedAmount(float32(json.Total), float32(json.Amount))
		c.JSON(http.StatusOK, gin.H{"status": "success", "newTotal": newTotal,"discount": json.Amount})
	} else {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Amount exceeds the total, cannot apply discount"})
	}
}

func DisFixedAmount(total , amount float32) (float32, error) {
	if total >= amount && amount > 1 && total > 1{
		newTotal := total - amount
		return newTotal, nil
	} 
	return 0, errors.New("Amount exceeds the total, cannot apply discount")

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
		// discount := json.Total * json.Percentage / 100
		// newTotal := json.Total - discount
		newTotal, discount, _  := DisPercentage(float32(json.Total), float32(json.Percentage))
		c.JSON(http.StatusOK, gin.H{"status": "ok", "newTotal": newTotal,"discount": discount})
	}
}

func DisPercentage(total, percentage float32) (float32, float32, error) {
    discount := total * percentage / 100
    newTotal := total - discount
    return newTotal, discount, nil
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


    var results []orm.CARTITEM

    if err := orm.Db.Preload("Customer").Preload("Product").Where("customer_id = ?", json.CustomerID).Find(&results).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Database error", "error": err.Error()})
        return
    }

		if json.Total < 1 || json.Percentage < 1 || json.Percentage > 100 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"status": "error", "message": "Amount exceeds the total, cannot apply discount"})
			return
		} else {
		

    if len(results) > 0 {
        		var productIDs []int
				var productPrice []int
				
		for _, result := range results {

			if result.Product.CategoryID == uint(json.Category) {
				if result.Quantity > 1 {
					productIDs = append(productIDs, result.ProductID)
					productPrice = append(productPrice, result.Quantity*result.Product.P_Price)
				}else{
					productIDs = append(productIDs, result.ProductID)
					productPrice = append(productPrice, result.Product.P_Price)}

			
		}


		
    }
			
	totalincate := sum(productPrice)
	totaldis := float32(totalincate) * float32(json.Percentage) / 100
	newTotal := float32(json.Total) - totaldis


	// c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Records found", "results":productIDs,"newtotal": newTotal})
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Records found","newtotal": newTotal,"discount": totaldis})

		return
}	}

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


	// limitdis := json.Total*20/100
	limitdis,_ := limitdis((float32(json.Total)))
	if float32(json.Point) <= limitdis {
		newTotal := json.Total-json.Point
		fmt.Println(newTotal)
		c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Records found","newtotal": newTotal,"discount": json.Point})
		return
	}else{
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Point exceeds the limit, cannot apply discount"})
		return
	}


}
func limitdis(total float32) (float32, error) {
	dis := total*20/100
	return dis, nil
}

type seasonalspecialcampaignsbody struct {
    Total float32 `json:"total" validate:"required"`
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
			// countdis := json.Total / json.Every
			// newTotal := json.Total - (countdis * json.Discount)
			// fmt.Println(newTotal)
			newTotal, countdis, _ := dis(float32(json.Total), float32(json.Every), float32(json.Discount))
		c.JSON(http.StatusOK, gin.H{"status": "ok", "newtotal": newTotal,"discount": countdis})
		return
	}else{
		c.JSON(http.StatusBadRequest, gin.H{"status": "error"})
		return
	}

}

func dis(total, eve, discount float32) (float32, float32, error) {
    countdisInt := int(total / eve)
    countdis := float32(countdisInt)
    newTotal := total - (countdis * discount)
    calculatedDiscount := countdis * discount
    return newTotal,calculatedDiscount, nil
}

func GetAllCartItemById(c *gin.Context) {
    customerID := c.Param("id")

    if customerID == "" {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Customer ID is required"})
        return
    }

    var cartitems []orm.CARTITEM
    result := orm.Db.Where("customer_id = ?", customerID).Find(&cartitems)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve cart items", "error": result.Error.Error()})
        return
    }

    for i, cartitem := range cartitems {
        var product orm.PRODUCT
        result := orm.Db.Where("id = ?", cartitem.ProductID).First(&product)
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve product", "error": result.Error.Error()})
            return
        }

        var images []orm.IMAGE
        imageResult := orm.Db.Where("product_id = ?", product.ID).First(&images)
        if imageResult.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve images", "error": imageResult.Error.Error()})
            return
        }

        if len(images) > 0 {
            selectedImage := images[0]
            product.Image = selectedImage.Image
        }


        cartitems[i].Product = product
    }


    c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Cart items retrieved successfully", "cart_items": cartitems})
}
