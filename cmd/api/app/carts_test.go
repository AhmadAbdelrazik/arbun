package app

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/pkg/assert"
	"net/http"
	"testing"

	"github.com/Rhymond/go-money"
)

func TestCart(t *testing.T) {
	ts := NewTestClient()
	defer ts.Close()
	cookies := InitializeCartTest(t, ts)
	customerCookie := cookies[0]
	t.Run("GetEmptyCart", func(t *testing.T) {
		GetEmptyCart(t, ts, customerCookie)
	})
	t.Run("PostToCart", func(t *testing.T) {
		PostCart(t, ts, customerCookie)
	})
	t.Run("DeleteFromCart", func(t *testing.T) {
		DeleteFromCart(t, ts, customerCookie)
	})

}

func GetEmptyCart(t *testing.T, ts *TestClient, customerCookie *http.Cookie) {
	res, err := ts.GetWithCookies("/cart", customerCookie)
	assert.Nil(t, err)

	var responseBody struct {
		Cart struct {
			Items []interface{} `json:"items"`
			Price *money.Money  `json:"price"`
		} `json:"cart"`
	}

	err = ts.ReadResponseBody(res, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, len(responseBody.Cart.Items), 0)
	assert.MoneyEqual(t, responseBody.Cart.Price, money.New(0, money.EGP))

}

func PostCart(t *testing.T, ts *TestClient, customerCookie *http.Cookie) {
	cartItems := struct {
		Items []domain.CartItem `json:"items"`
	}{
		Items: []domain.CartItem{
			{
				ProductID: 1,
				Amount:    2,
			},
			{
				ProductID: 2,
				Amount:    5,
			},
		},
	}

	res, err := ts.PostWithCookies("/cart", cartItems, customerCookie)
	assert.Nil(t, err)

	var responseBody struct {
		Cart struct {
			Items []struct {
				ProductID  int          `json:"product_id"`
				Name       string       `json:"name"`
				Amount     int          `json:"amount"`
				ItemPrice  *money.Money `json:"item_price"`
				TotalPrice *money.Money `json:"total_price"`
			} `json:"items"`
			Price *money.Money `json:"price"`
		} `json:"cart"`
	}

	err = ts.ReadResponseBody(res, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, len(responseBody.Cart.Items), 2)

	assert.Equal(t, responseBody.Cart.Items[0].Amount, 2)
	assert.MoneyEqual(t, responseBody.Cart.Items[0].ItemPrice, money.New(550, money.EGP))
	assert.MoneyEqual(t, responseBody.Cart.Items[0].TotalPrice, money.New(1100, money.EGP))

	assert.Equal(t, responseBody.Cart.Items[1].Amount, 5)
	assert.MoneyEqual(t, responseBody.Cart.Items[1].ItemPrice, money.New(2550, money.EGP))
	assert.MoneyEqual(t, responseBody.Cart.Items[1].TotalPrice, money.New(12750, money.EGP))

	assert.MoneyEqual(t, responseBody.Cart.Price, money.New(13850, money.EGP))

}

func DeleteFromCart(t *testing.T, ts *TestClient, customerCookie *http.Cookie) {
	var input struct {
		ProductID int64 `json:"product_id"`
	}

	input.ProductID = 2

	res, err := ts.DeleteWithCookies("/cart", input, customerCookie)
	assert.Nil(t, err)

	var responseBody struct {
		Cart struct {
			Items []struct {
				ProductID  int          `json:"product_id"`
				Name       string       `json:"name"`
				Amount     int          `json:"amount"`
				ItemPrice  *money.Money `json:"item_price"`
				TotalPrice *money.Money `json:"total_price"`
			} `json:"items"`
			Price *money.Money `json:"price"`
		} `json:"cart"`
	}

	err = ts.ReadResponseBody(res, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, len(responseBody.Cart.Items), 1)

	assert.Equal(t, responseBody.Cart.Items[0].Amount, 2)

	assert.MoneyEqual(t, responseBody.Cart.Items[0].ItemPrice, money.New(550, money.EGP))
	assert.MoneyEqual(t, responseBody.Cart.Items[0].TotalPrice, money.New(1100, money.EGP))

	assert.MoneyEqual(t, responseBody.Cart.Price, money.New(1100, money.EGP))
}

func InitializeCartTest(t *testing.T, ts *TestClient) []*http.Cookie {
	var admin1 domain.Admin
	admin1.Name = "admin1"
	admin1.Email = "admin1@example.com"

	adminCookie := AddAdmin(t, ts, admin1, "password1")

	product1 := domain.Product{
		Name:            "product 1",
		Description:     "description 1",
		Vendor:          "vendor 1",
		AvailableAmount: 10,
		Price:           money.New(550, money.EGP),
	}
	product2 := domain.Product{
		Name:            "product 2",
		Description:     "description 2",
		Vendor:          "vendor 2",
		AvailableAmount: 10,
		Price:           money.New(2550, money.EGP),
	}
	product3 := domain.Product{
		Name:            "product 3",
		Description:     "description 3",
		Vendor:          "vendor 3",
		AvailableAmount: 10,
		Price:           money.New(1500, money.EGP),
	}

	AddProduct(t, ts, product1, adminCookie)
	AddProduct(t, ts, product2, adminCookie)
	AddProduct(t, ts, product3, adminCookie)

	var customer1 domain.Customer
	customer1.Name = "customer1"
	customer1.Email = "customer1@example.com"
	customerCookie := AddCustomer(t, ts, customer1, "password1")

	return []*http.Cookie{customerCookie, adminCookie}
}
