package app

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/pkg/assert"
	"net/http"
	"testing"
	"time"

	"github.com/Rhymond/go-money"
)

func TestOrders(t *testing.T) {
	ts := NewTestClient()
	defer ts.Close()

	cookies := InitializeOrderTest(t, ts)
	customerCookie, adminCookie := cookies[0], cookies[1]

	t.Run("place order", func(t *testing.T) {
		t.Run("valid order", func(t *testing.T) {
			validCashOrder(t, ts, customerCookie)
		})
		t.Run("invalid order", func(t *testing.T) {
			invalidOrder(t, ts, customerCookie)
		})
	})

	t.Run("get order", func(t *testing.T) {
		getOrder(t, ts, customerCookie)
	})

	t.Run("get orders", func(t *testing.T) {
		getOrders(t, ts, customerCookie)
	})

	t.Run("change order status", func(t *testing.T) {
		changeOrderStatus(t, ts, customerCookie, adminCookie)
	})
}

func validCashOrder(t *testing.T, ts *TestClient, customerCookie *http.Cookie) {
	var input struct {
		DeliveryAddress domain.Address       `json:"address"`
		MobilePhone     domain.MobilePhone   `json:"mobile_phone"`
		PaymentMethod   domain.PaymentMethod `json:"payment_method"`
	}

	input.DeliveryAddress.Governorate = "Test Governorate"
	input.DeliveryAddress.City = "Test City"
	input.DeliveryAddress.Street = "Test Street"
	input.DeliveryAddress.AdditionalInfo = "Test Additional Info"
	input.MobilePhone = "01234567890"
	input.PaymentMethod = domain.PaymentCash

	res, err := ts.PostWithCookies("/checkout", input, customerCookie)

	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)

	var response struct {
		Order domain.Order `json:"order"`
	}

	err = ts.ReadResponseBody(res, &response)
	assert.Nil(t, err)

	assert.Equal(t, response.Order.Address, input.DeliveryAddress)
	assert.Equal(t, response.Order.MobilePhone, input.MobilePhone)
	assert.Equal(t, response.Order.PaymentType, domain.PaymentCash)
	assert.Equal(t, len(response.Order.Cart.Items), 3)
	equal, err := response.Order.Cart.Price.Equals(money.New(15800, money.EGP))
	assert.True(t, equal)
	assert.Nil(t, err)
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

		res, err := ts.Post("/checkout", input)
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

func InitializeOrderTest(t *testing.T, ts *TestClient) []*http.Cookie {
	cookies := InitializeCartTest(t, ts)
	customerCookie := cookies[0]

	items := []domain.CartItem{
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

	return cookies
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
		equal, err := response.Order.Cart.Price.Equals(money.New(15800, money.EGP))
		assert.True(t, equal)
		assert.Nil(t, err)
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

func getOrders(t *testing.T, ts *TestClient, customerCookie *http.Cookie) {
	res, err := ts.GetWithCookies("/orders", customerCookie)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	var responseBody struct {
		Orders []domain.Order `json:"orders"`
	}

	err = ts.ReadResponseBody(res, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, len(responseBody.Orders), 1)

	order := responseBody.Orders[0]

	var address domain.Address
	address.Governorate = "Test Governorate"
	address.City = "Test City"
	address.Street = "Test Street"
	address.AdditionalInfo = "Test Additional Info"

	mobilePhone := "01234567890"

	assert.Equal(t, order.Address, address)
	assert.Equal(t, string(order.MobilePhone), mobilePhone)
	assert.Equal(t, order.PaymentType, domain.PaymentCash)
	assert.Equal(t, len(order.Cart.Items), 3)
	equal, err := order.Cart.Price.Equals(money.New(15800, money.EGP))
	assert.True(t, equal)
	assert.Nil(t, err)
	assert.Equal(t, order.Status, domain.StatusDispatched)
	assert.True(t, order.CreatedAt.Before(time.Now()))
}

func changeOrderStatus(t *testing.T, ts *TestClient, customerCookie, adminCookie *http.Cookie) {
	t.Run("valid change", func(t *testing.T) {
		validStatusChange(t, ts, customerCookie, adminCookie)
	})

	t.Run("invalid change", func(t *testing.T) {
		invalidStatusChange(t, ts, adminCookie)
	})
}

func validStatusChange(t *testing.T, ts *TestClient, customerCookie, adminCookie *http.Cookie) {
	var input struct {
		Status string `json:"status"`
	}
	input.Status = "canceled"

	res, err := ts.PatchWithCookies("/orders/1", input, adminCookie)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusOK)

	var responseBody struct {
		Status string `json:"status"`
	}

	err = ts.ReadResponseBody(res, &responseBody)
	assert.Nil(t, err)
	assert.Equal(t, responseBody.Status, "order status updated")

	res, err = ts.GetWithCookies("/orders/1", customerCookie)
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
	equal, err := response.Order.Cart.Price.Equals(money.New(15800, money.EGP))
	assert.True(t, equal)
	assert.Nil(t, err)

	assert.Equal(t, response.Order.Status, domain.StatusCanceled)
	assert.True(t, response.Order.CreatedAt.Before(time.Now()))
}

func invalidStatusChange(t *testing.T, ts *TestClient, adminCookie *http.Cookie) {
	t.Run("unknown order", func(t *testing.T) {
		var input struct {
			Status string `json:"status"`
		}
		input.Status = "canceled"

		res, err := ts.PatchWithCookies("/orders/5", input, adminCookie)
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
