// cart/controller_test.go
package cart

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/quick"

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


func TestDisFixedAmount(t *testing.T) {
	err := quick.Check(func(total, amount int) bool {
		if total <= 1 {
			return true
		}

		if amount < 1 || amount > total {
			return true
		}

		newTotal, discountErr := DisFixedAmount(total, amount)

		if discountErr != nil {
			return false
		}

		return newTotal == total-amount
	}, nil)

	if err != nil {
		t.Error("Test failed:", err)
	}

	invalidTotal, invalidAmount := 1, 2
	_, discountErr := DisFixedAmount(invalidTotal, invalidAmount)

	if discountErr == nil {
		t.Error("Expected error but got none")
	}
}





func TestDisPercentage(t *testing.T) {
    total1 := 100
    percentage1 := 10
    expectedNewTotal1 := 90
    expectedDiscount1 := 10

    newTotal1, discount1, err1 := DisPercentage(total1, percentage1)
    if err1 != nil {
        t.Errorf("Unexpected error: %v", err1)
    }

    if newTotal1 != expectedNewTotal1 || discount1 != expectedDiscount1 {
        t.Errorf("Test case 1 failed: expected (%d, %d), got (%d, %d)", expectedNewTotal1, expectedDiscount1, newTotal1, discount1)
    }

    total2 := 200
    percentage2 := 15
    expectedNewTotal2 := 170
    expectedDiscount2 := 30

    newTotal2, discount2, err2 := DisPercentage(total2, percentage2)
    if err2 != nil {
        t.Errorf("Unexpected error: %v", err2)
    }

    if newTotal2 != expectedNewTotal2 || discount2 != expectedDiscount2 {
        t.Errorf("Test case 2 failed: expected (%d, %d), got (%d, %d)", expectedNewTotal2, expectedDiscount2, newTotal2, discount2)
    }

}


func TestLimitDis(t *testing.T) {
    total1 := 100
    expectedDiscount1 := 20

    discount1, err1 := limitdis(total1)
    if err1 != nil {
        t.Errorf("Unexpected error: %v", err1)
    }

    if discount1 != expectedDiscount1 {
        t.Errorf("Test case 1 failed: expected %d, got %d", expectedDiscount1, discount1)
    }

    total2 := 200
    expectedDiscount2 := 40

    discount2, err2 := limitdis(total2)
    if err2 != nil {
        t.Errorf("Unexpected error: %v", err2)
    }

    if discount2 != expectedDiscount2 {
        t.Errorf("Test case 2 failed: expected %d, got %d", expectedDiscount2, discount2)
    }

}

func TestDis(t *testing.T) {
    total1 := 830
    eve1 := 300
    discount1 := 40
    expectedNewTotal1 := 750
    expectedCountdis1 := 80

    newTotal1, countdis1, err1 := dis(total1, eve1, discount1 )
    if err1 != nil {
        t.Errorf("Unexpected error: %v", err1)
    }

    if newTotal1 != expectedNewTotal1 || countdis1 != expectedCountdis1 {
        t.Errorf("Test case 1 failed: expected (%d, %d), got (%d, %d)", expectedNewTotal1, expectedCountdis1, newTotal1, countdis1)
    }

}