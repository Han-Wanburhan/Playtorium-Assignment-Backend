package main

import (
	"fmt"

	"shopping/go-api/orm"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"

	CatreToCartController "shopping/go-api/controller/cart"
	CatreGoryController "shopping/go-api/controller/category"
	CatreProController "shopping/go-api/controller/product"

	"github.com/gin-gonic/gin"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	orm.InitDB()

	// สร้าง instance ของ Gin engine
	r := gin.Default()

	// ตั้งค่า CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // ระบุโดเมนของเว็บเบราว์เซอร์
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Content-Length", "Authorization"}

	// เพิ่ม CORS middleware เข้าไปใน Gin engine
	r.Use(cors.New(config))

	// ทดสอบใส่เพิ่มเติม เพื่อทดสอบว่า CORS middleware ทำงานหรือไม่
	r.POST("/addcategory", CatreGoryController.AddCategory)
	r.POST("/addproduct", CatreProController.AddProduct)
	r.POST("/addtocart", CatreToCartController.AddToCart)
	r.POST("/total/couponfixdis", CatreToCartController.CouponFixedAmount)
	r.POST("/total/couponperdis", CatreToCartController.CouponPercentageDiscount)
	r.POST("/total/ontopdisbycate", CatreToCartController.OntopPercentageDiscountByCategory)
	r.POST("/total/ontopdisbypoint", CatreToCartController.OntopDiscountByPoint)
	r.POST("/total/seasonalcampaigns", CatreToCartController.SeasonalSpecialCampaigns)

	// Run the server
	r.Run("localhost:8081")
}
