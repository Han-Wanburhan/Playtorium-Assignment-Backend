// cart/controller_test.go
package cart

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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
	r := gin.New()

	input := couponpercentagediscountbody{
		Total:      600,
		Percentage: 10,
	}

	inputJSON, _ := json.Marshal(input)

	req, err := http.NewRequest("POST", "/coupon", bytes.NewBuffer(inputJSON))
	assert.NoError(t, err)

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	r.POST("/coupon", CouponPercentageDiscount)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	expected := gin.H{
		"status":   "ok",
		"newTotal": 540,
		"discount": 60,
	}

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	expectedJSON, _ := json.Marshal(expected)
	actualJSON, _ := json.Marshal(response)
	assert.JSONEq(t, string(expectedJSON), string(actualJSON))
}


// Coupon Fixed Amount
func TestDisFixedAmount(t *testing.T) {
	total1 := float32(100)
	amount1 := float32(20)
	expectedNewTotal1 := float32(80)

	newTotal1, err1 := DisFixedAmount(total1, amount1)
	if err1 != nil {
		t.Errorf("Unexpected error: %v", err1)
	}

	if newTotal1 != expectedNewTotal1 {
		t.Errorf("Test case 1 failed: expected %.2f, got %.2f", expectedNewTotal1, newTotal1)
	}

	total2 := float32(50)
	amount2 := float32(60)

	_, err2 := DisFixedAmount(total2, amount2)
	expectedError2 := errors.New("Amount exceeds the total, cannot apply discount")

	if err2 == nil || err2.Error() != expectedError2.Error() {
		t.Errorf("Test case 2 failed: expected error '%v', got '%v'", expectedError2, err2)
	}
}

// Coupon Percentage Discount
func TestDisPercentage(t *testing.T) {
	total1 := float32(100)
	percentage1 := float32(10)
	expectedNewTotal1 := float32(90)
	expectedDiscount1 := float32(10)

	newTotal1, discount1, err1 := DisPercentage(total1, percentage1)
	if err1 != nil {
		t.Errorf("Unexpected error: %v", err1)
	}

	if newTotal1 != expectedNewTotal1 || discount1 != expectedDiscount1 {
		t.Errorf("Test case 1 failed: expected (%.2f, %.2f), got (%.2f, %.2f)", expectedNewTotal1, expectedDiscount1, newTotal1, discount1)
	}

	total2 := float32(200)
	percentage2 := float32(15)
	expectedNewTotal2 := float32(170)
	expectedDiscount2 := float32(30)

	newTotal2, discount2, err2 := DisPercentage(total2, percentage2)
	if err2 != nil {
		t.Errorf("Unexpected error: %v", err2)
	}

	if newTotal2 != expectedNewTotal2 || discount2 != expectedDiscount2 {
		t.Errorf("Test case 2 failed: expected (%.2f, %.2f), got (%.2f, %.2f)", expectedNewTotal2, expectedDiscount2, newTotal2, discount2)
	}
}


func TestLimitDis(t *testing.T) {
	total1 := float32(100)
	expectedDiscount1 := float32(20)

	dis1, err1 := limitdis(total1)
	if err1 != nil {
		t.Errorf("Unexpected error: %v", err1)
	}

	if dis1 != expectedDiscount1 {
		t.Errorf("Test case 1 failed: expected %.2f, got %.2f", expectedDiscount1, dis1)
	}

	total2 := float32(200)
	expectedDiscount2 := float32(40)

	dis2, err2 := limitdis(total2)
	if err2 != nil {
		t.Errorf("Unexpected error: %v", err2)
	}

	if dis2 != expectedDiscount2 {
		t.Errorf("Test case 2 failed: expected %.2f, got %.2f", expectedDiscount2, dis2)
	}
}


// test seasonal special campaigns
func TestDis(t *testing.T) {
	total1 := float32(100)
	eve1 := float32(20)
	discount1 := float32(5)

	expectedNewTotal1 := float32(75)
	expectedCalculatedDiscount1 := float32(25)

	newTotal1, calculatedDiscount1, err1 := dis(total1, eve1, discount1)
	if err1 != nil {
		t.Errorf("Unexpected error: %v", err1)
	}

	if newTotal1 != expectedNewTotal1 || calculatedDiscount1 != expectedCalculatedDiscount1 {
		t.Errorf("Test case 1 failed: expected (%.2f, %.2f), got (%.2f, %.2f)", expectedNewTotal1, expectedCalculatedDiscount1, newTotal1, calculatedDiscount1)
	}

	total2 := float32(200)
	eve2 := float32(50)
	discount2 := float32(10)

	expectedNewTotal2 := float32(160)
	expectedCalculatedDiscount2 := float32(40)

	newTotal2, calculatedDiscount2, err2 := dis(total2, eve2, discount2)
	if err2 != nil {
		t.Errorf("Unexpected error: %v", err2)
	}

	if newTotal2 != expectedNewTotal2 || calculatedDiscount2 != expectedCalculatedDiscount2 {
		t.Errorf("Test case 2 failed: expected (%.2f, %.2f), got (%.2f, %.2f)", expectedNewTotal2, expectedCalculatedDiscount2, newTotal2, calculatedDiscount2)
	}
}