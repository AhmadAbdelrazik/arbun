package app

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/pkg/assert"
	"AhmadAbdelrazik/arbun/internal/services"
	"net/http"
	"testing"
	"time"
)

func TestOrders(t *testing.T) {
	ts := NewTestClient()
	defer ts.Close()

	customerCookie := InitializeOrderTest(t, ts)

	t.Run("place order", func(t *testing.T) {
		t.Run("valid order", func(t *testing.T) {
			validOrder(t, ts, customerCookie)
		})
		t.Run("invalid order", func(t *testing.T) {
			invalidOrder(t, ts, customerCookie)
		})
	})

	t.Run("get order", func(t *testing.T) {
		getOrder(t, ts, customerCookie)
	})
}

func validOrder(t *testing.T, ts *TestClient, customerCookie *http.Cookie) {
	var input struct {
		DeliveryAddress domain.Address `json:"address"`
		MobilePhone     string         `json:"mobile_phone"`
	}

	input.DeliveryAddress.Governorate = "Test Governorate"
	input.DeliveryAddress.City = "Test City"
	input.DeliveryAddress.Street = "Test Street"
	input.DeliveryAddress.AdditionalInfo = "Test Additional Info"
	input.MobilePhone = "01234567890"

	res, err := ts.PostWithCookies("/checkout", input, customerCookie)

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)

	var response struct {
		Order domain.Order `json:"order"`
	}

	err = ts.ReadResponseBody(res, &response)
	assert.Nil(t, err)

	assert.Equal(t, response.Order.Address, input.DeliveryAddress)
	assert.Equal(t, string(response.Order.MobilePhone), input.MobilePhone)
	assert.Equal(t, response.Order.PaymentType, domain.PaymentCash)
	assert.Equal(t, len(response.Order.Cart.Items), 3)
	assert.Equal(t, response.Order.Cart.Price, 158)
	assert.Equal(t, response.Order.Status, domain.StatusDispatched)
	assert.True(t, response.Order.CreatedAt.Before(time.Now()))
}

func invalidOrder(t *testing.T, ts *TestClient, customerCookie *http.Cookie) {
	t.Run("no cookie", func(t *testing.T) {
		var input struct {
			DeliveryAddress domain.Address `json:"address"`
			MobilePhone     string         `json:"mobile_phone"`
		}

		input.DeliveryAddress.Governorate = "Test Governorate"
		input.DeliveryAddress.City = "Test City"
		input.DeliveryAddress.Street = "Test Street"
		input.DeliveryAddress.AdditionalInfo = "Test Additional Info"
		input.MobilePhone = "01234567890"

		res, err := ts.PostWithCookies("/checkout", input, &http.Cookie{})
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusUnauthorized)
		var response struct {
			Error string `json:"error"`
		}

		err = ts.ReadResponseBody(res, &response)
		assert.Nil(t, err)

		assert.Equal(t, response.Error, "Authentication credentials were not provided or are invalid.")
	})

	t.Run("invalid body", func(t *testing.T) {
		var input struct {
			DeliveryAddress domain.Address `json:"address"`
			MobilePhone     string         `json:"mobile_phone"`
			ProductID       int64          `json:"product_id"`
		}

		input.DeliveryAddress.Governorate = "Test Governorate"
		input.DeliveryAddress.City = "Test City"
		input.DeliveryAddress.Street = "Test Street"
		input.DeliveryAddress.AdditionalInfo = "Test Additional Info"
		input.MobilePhone = "01234567890"
		input.ProductID = 2

		res, err := ts.PostWithCookies("/checkout", input, customerCookie)
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		var response struct {
			Error string `json:"error"`
		}

		err = ts.ReadResponseBody(res, &response)
		assert.Nil(t, err)

		assert.Equal(t, response.Error, "body contains unknown key \"product_id\"")
	})
}

func InitializeOrderTest(t *testing.T, ts *TestClient) *http.Cookie {
	customerCookie := InitializeCartTest(t, ts)

	items := []services.InputItem{
		{
			ProductID: 1,
			Amount:    2,
		},
		{
			ProductID: 2,
			Amount:    4,
		},
		{
			ProductID: 3,
			Amount:    3,
		},
	}
	AddToCart(t, ts, items, customerCookie)

	return customerCookie
}

func getOrder(t *testing.T, ts *TestClient, customerCookie *http.Cookie) {
	t.Run("valid get", func(t *testing.T) {
		res, err := ts.GetWithCookies("/orders/1", customerCookie)
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusOK)

		var response struct {
			Order domain.Order `json:"order"`
		}

		err = ts.ReadResponseBody(res, &response)
		assert.Nil(t, err)

		var address domain.Address
		address.Governorate = "Test Governorate"
		address.City = "Test City"
		address.Street = "Test Street"
		address.AdditionalInfo = "Test Additional Info"

		mobilePhone := "01234567890"

		assert.Equal(t, response.Order.Address, address)
		assert.Equal(t, string(response.Order.MobilePhone), mobilePhone)
		assert.Equal(t, response.Order.PaymentType, domain.PaymentCash)
		assert.Equal(t, len(response.Order.Cart.Items), 3)
		assert.Equal(t, response.Order.Cart.Price, 158)
		assert.Equal(t, response.Order.Status, domain.StatusDispatched)
		assert.True(t, response.Order.CreatedAt.Before(time.Now()))
	})
	t.Run("invalid get", func(t *testing.T) {
		res, err := ts.GetWithCookies("/orders/3", customerCookie)
		assert.Nil(t, err)
		assert.Equal(t, res.StatusCode, http.StatusNotFound)

		var response struct {
			Error string `json:"error"`
		}

		err = ts.ReadResponseBody(res, &response)
		assert.Nil(t, err)

		assert.Equal(t, response.Error, "the requested resource could not be found")
	})
}
