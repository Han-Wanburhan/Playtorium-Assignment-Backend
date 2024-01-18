// cart/controller_test.go
package cart

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"shopping/go-api/orm"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCouponFixedAmount(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/total/couponfixdis", CouponFixedAmount)

	requestBody := []byte(`{"total":600,"amount":50}`)
	req, err := http.NewRequest("POST", "/total/couponfixdis", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

assert.JSONEq(t, `{"status":"success","newTotal":550,"discount":50}`, w.Body.String())
}

func TestCouponPercentageDiscount(t *testing.T) {
	// Create a new Gin router
	r := gin.New()

	// Define the input data
	input := couponpercentagediscountbody{
		Total:      600,
		Percentage: 10,
	}

	// Convert input to JSON
	inputJSON, _ := json.Marshal(input)

	// Create a request
	req, err := http.NewRequest("POST", "/coupon", bytes.NewBuffer(inputJSON))
	assert.NoError(t, err)

	// Set the request header
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Serve the request to the handler
	r.POST("/coupon", CouponPercentageDiscount)
	r.ServeHTTP(w, req)

	// Assert the HTTP status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Define the expected response
	expected := gin.H{
		"status":   "ok",
		"newTotal": 540,
		"discount": 60,
	}

	// Parse the response body
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Convert maps to JSON strings and compare
	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(response)
	assert.JSONEq(t, string(expectedJSON), string(actualJSON))
}

func TestOntopDiscountByPoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Mocking the ORM.Db for testing
	orm.Db = MockDB()

	// Create a new Gin router
	r := gin.New()

	// Define the input data
	input := ontopdiscountbypointbody{
		CustomerID: 1,
		Total:      830,
		Point:      68,
	}

	// Convert input to JSON
	inputJSON, _ := json.Marshal(input)

	// Create a request
	req, err := http.NewRequest("POST", "/ontop_discount_by_point", bytes.NewBuffer(inputJSON))
	assert.NoError(t, err)

	// Set the request header
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	w := httptest.NewRecorder()

	// Serve the request to the handler
	r.POST("/ontop_discount_by_point", OntopDiscountByPoint)
	r.ServeHTTP(w, req)

	// Assert the HTTP status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Define the expected response
	expected := gin.H{
		"status":   "success",
		"message":  "Records found",
		"newtotal": float64(830 - 68),
		"discount": 68,
	}

	// Parse the response body
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Convert maps to JSON strings and compare
	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(response)
	assert.JSONEq(t, string(expectedJSON), string(actualJSON))
}

// // MockDB provides a mock implementation of the database for testing
func MockDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	return db
}

// MockDBCleanUp cleans up the mock database after the tests
func MockDBCleanUp() {
	orm.Db = nil
}

