package handlers

import (
	"AhmadAbdelrazik/arbun/internal/assert"
	"AhmadAbdelrazik/arbun/internal/domain/admin"
	"AhmadAbdelrazik/arbun/internal/domain/product"
	"fmt"
	"net/http"
	"testing"
)

func productToPostProductInput(p product.Product) postProductInput {
	return postProductInput{
		Name:        p.Name,
		Description: p.Description,
		Vendor:      p.Vendor,
		Amount:      p.AvailableAmount,
		Properties:  p.Properties,
		Price:       p.Price,
	}
}

func TestProduct(t *testing.T) {
	ts := NewTestClient()
	defer ts.Close()

	a := admin.Admin{
		Email:    "admin@gmail.com",
		FullName: "admin",
	}

	products := []product.Product{
		{
			ID:              1,
			Name:            "product 1",
			Description:     "description of product 1",
			Vendor:          "vendor 1",
			AvailableAmount: 5,
			Version:         1,
			Price:           23.99,
			Properties: map[string]string{
				"size": "12",
			},
		},
		{
			ID:              2,
			Name:            "product 2",
			Description:     "description of product 2",
			Vendor:          "vendor 1",
			AvailableAmount: 8,
			Price:           18.99,
			Version:         1,
			Properties: map[string]string{
				"size": "14",
			},
		},
	}

	adminCookie := AddAdmin(t, ts, a, "admin123")

	t.Run("Post Products", func(t *testing.T) {
		postProduct(t, ts, adminCookie, products)
	})
	t.Run("Get Products", func(t *testing.T) {
		getProduct(t, ts, products)
	})
	t.Run("Patch Products", func(t *testing.T) {
		patchProduct(t, ts, adminCookie, products)
	})
	t.Run("Delete Products", func(t *testing.T) {
		deleteProduct(t, ts, adminCookie)
	})

}

func postProduct(t *testing.T, ts *TestClient, adminCookie *http.Cookie, products []product.Product) {
	t.Run("valid insertion", func(t *testing.T) {
		validPost(t, ts, adminCookie, products)
	})
	t.Run("invalid insertion", func(t *testing.T) {
		invalidPost(t, ts, adminCookie)
	})
}

func getProduct(t *testing.T, ts *TestClient, products []product.Product) {
	t.Run("valid fetching", func(t *testing.T) {
		for i, p := range products {
			var responseBody struct {
				Product product.Product `json:"product"`
			}
			t.Run(fmt.Sprintf("get Product %d", i), func(t *testing.T) {
				res, err := ts.Get(fmt.Sprintf("/products/%v", p.ID))
				assert.Nil(t, err)

				err = ts.ReadResponseBody(res, &responseBody)
				assert.Nil(t, err)
				assert.Equal(t, responseBody.Product.String(), p.String())
			})
		}

	})
	t.Run("invalid fetching", func(t *testing.T) {
		t.Run("no product", func(t *testing.T) {
			res, err := ts.Get("/products/5")
			assert.Nil(t, err)

			var responseBody struct {
				Error string `json:"error"`
			}
			err = ts.ReadResponseBody(res, &responseBody)
			assert.Nil(t, err)
			assert.Equal(t, responseBody.Error, "the requested resource could not be found")
		})
		t.Run("bad value", func(t *testing.T) {
			res, err := ts.Get("/products/x")
			assert.Nil(t, err)

			var responseBody struct {
				Error string `json:"error"`
			}
			err = ts.ReadResponseBody(res, &responseBody)
			assert.Nil(t, err)
			assert.Equal(t, responseBody.Error, "invalid id param")
		})
	})

}

func patchProduct(t *testing.T, ts *TestClient, adminCookie *http.Cookie, products []product.Product) {
	t.Run("valid update", func(t *testing.T) {
		newProduct := products[0]
		newProduct.Name = "product 1 updated"
		newProduct.Description = "product 1 description updated"
		newProduct.AvailableAmount = 53
		newProduct.Version++

		reqBody := patchProductInput{
			Name:            &newProduct.Name,
			Description:     &newProduct.Description,
			AvailableAmount: &newProduct.AvailableAmount,
		}

		res, err := ts.PatchWithCookies("/products/1", reqBody, adminCookie)
		assert.Nil(t, err)

		var responseBody struct {
			Product product.Product `json:"product"`
		}
		err = ts.ReadResponseBody(res, &responseBody)
		assert.Nil(t, err)
		assert.Equal(t, responseBody.Product.String(), newProduct.String())
	})

	t.Run("invalid update", func(t *testing.T) {
		newProduct := products[0]
		newProduct.Name = "product 1 updated"
		newProduct.Description = "product 1 description updated"
		newProduct.AvailableAmount = 53
		newProduct.Version++

		reqBody := patchProductInput{
			Name:            &newProduct.Name,
			Description:     &newProduct.Description,
			AvailableAmount: &newProduct.AvailableAmount,
		}

		res, err := ts.PatchWithCookies("/products/3", reqBody, adminCookie)
		assert.Nil(t, err)

		var responseBody struct {
			Error string `json:"error"`
		}
		err = ts.ReadResponseBody(res, &responseBody)
		assert.Nil(t, err)
		assert.Equal(t, responseBody.Error, "the requested resource could not be found")
	})
}

func deleteProduct(t *testing.T, ts *TestClient, adminCookie *http.Cookie) {
	t.Run("valid delete", func(t *testing.T) {
		res, err := ts.DeleteWithCookies("/products/1", nil, adminCookie)
		assert.Nil(t, err)

		var responseBody struct {
			Message string `json:"message"`
		}
		err = ts.ReadResponseBody(res, &responseBody)
		assert.Nil(t, err)
		assert.Equal(t, responseBody.Message, "deleted successfully")
	})

	t.Run("invalid deletion", func(t *testing.T) {
		res, err := ts.DeleteWithCookies("/products/3", nil, adminCookie)
		assert.Nil(t, err)

		var responseBody struct {
			Error string `json:"error"`
		}
		err = ts.ReadResponseBody(res, &responseBody)
		assert.Nil(t, err)
		assert.Equal(t, responseBody.Error, "the requested resource could not be found")
	})
}

func validPost(t *testing.T, ts *TestClient, adminCookie *http.Cookie, products []product.Product) {
	for i, p := range products {
		t.Run(fmt.Sprintf("product%d", i), func(t *testing.T) {
			res, err := ts.PostWithCookies("/products", productToPostProductInput(p), adminCookie)
			assert.Nil(t, err)

			var responseBody struct {
				Product product.Product `json:"product"`
			}

			err = ts.ReadResponseBody(res, &responseBody)
			assert.Nil(t, err)
			assert.Equal(t, responseBody.Product.String(), p.String())
			assert.Equal(t, res.StatusCode, http.StatusCreated)
		})
	}
}

func invalidPost(t *testing.T, ts *TestClient, adminCookie *http.Cookie) {
	mainProduct := product.Product{
		ID:              3,
		Name:            "product 2",
		Description:     "description of product 2",
		Vendor:          "vendor 1",
		AvailableAmount: 4,
		Price:           23.99,
		Version:         1,
		Properties: map[string]string{
			"size": "14",
		},
	}
	t.Run("duplicate product", func(t *testing.T) {
		var responseBody struct {
			Error string `json:"error"`
		}

		res, err := ts.PostWithCookies("/products", productToPostProductInput(mainProduct), adminCookie)
		assert.Nil(t, err)

		err = ts.ReadResponseBody(res, &responseBody)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		assert.Equal(t, responseBody.Error, "product already exists")
	})
	t.Run("invalid product amount", func(t *testing.T) {
		p := mainProduct
		p.AvailableAmount = 0

		var responseBody struct {
			Error struct {
				Amount string `json:"amount"`
			} `json:"error"`
		}

		res, err := ts.PostWithCookies("/products", productToPostProductInput(p), adminCookie)
		assert.Nil(t, err)

		ts.ReadResponseBody(res, &responseBody)

		assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
		assert.Equal(t, responseBody.Error.Amount, "must be more than 0")
	})
	t.Run("missing product name and description", func(t *testing.T) {
		p := mainProduct
		p.Name = ""
		p.Description = ""

		var responseBody struct {
			Error struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"error"`
		}

		res, err := ts.PostWithCookies("/products", productToPostProductInput(p), adminCookie)
		assert.Nil(t, err)

		ts.ReadResponseBody(res, &responseBody)

		assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
		assert.Equal(t, responseBody.Error.Name, "can't be empty")
		assert.Equal(t, responseBody.Error.Description, "can't be empty")
	})
	t.Run("missing product name", func(t *testing.T) {
		p := mainProduct
		p.Name = ""

		var responseBody struct {
			Error struct {
				Name string `json:"name"`
			} `json:"error"`
		}

		res, err := ts.PostWithCookies("/products", productToPostProductInput(p), adminCookie)
		assert.Nil(t, err)

		ts.ReadResponseBody(res, &responseBody)

		assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
		assert.Equal(t, responseBody.Error.Name, "can't be empty")
	})
	t.Run("invalid product", func(t *testing.T) {
		p := mainProduct

		var responseBody struct {
			Error string `json:"error"`
		}

		res, err := ts.PostWithCookies("/products", p, adminCookie)
		assert.Nil(t, err)

		err = ts.ReadResponseBody(res, &responseBody)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		assert.Equal(t, responseBody.Error, `body contains unknown key "ID"`)
	})
}
