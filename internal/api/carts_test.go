package handlers

import (
	"AhmadAbdelrazik/arbun/internal/assert"
	"AhmadAbdelrazik/arbun/internal/domain/admin"
	"AhmadAbdelrazik/arbun/internal/domain/customer"
	"AhmadAbdelrazik/arbun/internal/domain/product"
	"AhmadAbdelrazik/arbun/internal/services"
	"net/http"
	"testing"
)

func TestCart(t *testing.T) {
	ts := NewTestClient()
	customerCookie := InitializeCartTest(t, ts)
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
			Price int           `json:"price"`
		} `json:"cart"`
	}

	err = ts.ReadResponseBody(res, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, len(responseBody.Cart.Items), 0)
	assert.Equal(t, responseBody.Cart.Price, 0)
}

func PostCart(t *testing.T, ts *TestClient, customerCookie *http.Cookie) {
	cartItems := struct {
		Items []services.InputItem `json:"items"`
	}{
		Items: []services.InputItem{
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
				ProductID  int     `json:"product_id"`
				Name       string  `json:"name"`
				Amount     int     `json:"amount"`
				ItemPrice  float32 `json:"item_price"`
				TotalPrice float32 `json:"total_price"`
			} `json:"items"`
			Price float32 `json:"price"`
		} `json:"cart"`
	}

	err = ts.ReadResponseBody(res, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, len(responseBody.Cart.Items), 2)

	assert.Equal(t, responseBody.Cart.Items[0].Amount, 2)
	assert.Equal(t, responseBody.Cart.Items[0].ItemPrice, 5.5)
	assert.Equal(t, responseBody.Cart.Items[0].TotalPrice, 11)

	assert.Equal(t, responseBody.Cart.Items[1].Amount, 5)
	assert.Equal(t, responseBody.Cart.Items[1].ItemPrice, 25.5)
	assert.Equal(t, responseBody.Cart.Items[1].TotalPrice, 127.5)

	assert.Equal(t, responseBody.Cart.Price, 138.5)

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
				ProductID  int     `json:"product_id"`
				Name       string  `json:"name"`
				Amount     int     `json:"amount"`
				ItemPrice  float32 `json:"item_price"`
				TotalPrice float32 `json:"total_price"`
			} `json:"items"`
			Price float32 `json:"price"`
		} `json:"cart"`
	}

	err = ts.ReadResponseBody(res, &responseBody)
	assert.Nil(t, err)

	assert.Equal(t, len(responseBody.Cart.Items), 1)

	assert.Equal(t, responseBody.Cart.Items[0].Amount, 2)
	assert.Equal(t, responseBody.Cart.Items[0].ItemPrice, 5.5)
	assert.Equal(t, responseBody.Cart.Items[0].TotalPrice, 11)

	assert.Equal(t, responseBody.Cart.Price, 11)
}

func InitializeCartTest(t *testing.T, ts *TestClient) *http.Cookie {
	admin1 := admin.Admin{
		FullName: "admin1",
		Email:    "admin1@example.com",
	}
	adminCookie := AddAdmin(t, ts, admin1, "password1")

	product1 := product.Product{
		Name:            "product 1",
		Description:     "description 1",
		Vendor:          "vendor 1",
		AvailableAmount: 10,
		Price:           5.5,
	}
	product2 := product.Product{
		Name:            "product 2",
		Description:     "description 2",
		Vendor:          "vendor 2",
		AvailableAmount: 10,
		Price:           25.5,
	}
	product3 := product.Product{
		Name:            "product 3",
		Description:     "description 3",
		Vendor:          "vendor 3",
		AvailableAmount: 10,
		Price:           15,
	}

	AddProduct(t, ts, product1, adminCookie)
	AddProduct(t, ts, product2, adminCookie)
	AddProduct(t, ts, product3, adminCookie)

	customer1 := customer.Customer{
		FullName: "customer1",
		Email:    "customer1@example.com",
	}
	customerCookie := AddCustomer(t, ts, customer1, "password1")

	return customerCookie
}

func AddCustomer(t *testing.T, ts *TestClient, c customer.Customer, password string) *http.Cookie {
	t.Helper()

	body := struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"type"`
	}{
		FullName: c.FullName,
		Email:    c.Email,
		Password: password,
		UserType: services.TypeCustomer,
	}
	res, err := ts.Post("/signup", body)
	assert.Nil(t, err)
	cookie := ts.GetCookie(res, AuthCookie)

	assert.Nil(t, cookie.Valid())
	assert.True(t, len(cookie.Value) == 26)

	return cookie
}

func AddAdmin(t *testing.T, ts *TestClient, a admin.Admin, password string) *http.Cookie {
	t.Helper()

	body := struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		UserType string `json:"type"`
	}{
		FullName: a.FullName,
		Email:    a.Email,
		Password: password,
		UserType: services.TypeAdmin,
	}
	res, err := ts.Post("/signup", body)
	assert.Nil(t, err)
	cookie := ts.GetCookie(res, AuthCookie)

	assert.Nil(t, cookie.Valid())
	assert.True(t, len(cookie.Value) == 26)

	return cookie
}

func AddProduct(t *testing.T, ts *TestClient, p product.Product, adminCookie *http.Cookie) {
	t.Helper()
	res, err := ts.PostWithCookies("/products", productToPostProductInput(p), adminCookie)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, http.StatusCreated)
}
