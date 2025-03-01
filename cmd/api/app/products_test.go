package app

import (
	"AhmadAbdelrazik/arbun/internal/domain"
	"AhmadAbdelrazik/arbun/internal/pkg/assert"
	"fmt"
	"net/http"
	"testing"

	"github.com/Rhymond/go-money"
)

func TestProduct(t *testing.T) {
	ts := NewTestClient()
	defer ts.Close()

	var a domain.Admin
	a.Email = "admin@gmail.com"
	a.Name = "admin"

	products := []domain.Product{
		{
			ID:              1,
			Name:            "product 1",
			Description:     "description of product 1",
			Vendor:          "vendor 1",
			AvailableAmount: 5,
			Version:         1,
			Price:           money.New(2399, money.EGP),
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
			Price:           money.New(1899, money.EGP),
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

func postProduct(t *testing.T, ts *TestClient, adminCookie *http.Cookie, products []domain.Product) {
	t.Run("valid insertion", func(t *testing.T) {
		validPost(t, ts, adminCookie, products)
	})
	t.Run("invalid insertion", func(t *testing.T) {
		invalidPost(t, ts, adminCookie)
	})
}

func getProduct(t *testing.T, ts *TestClient, products []domain.Product) {
	t.Run("valid fetching", func(t *testing.T) {
		for i, p := range products {
			var responseBody struct {
				Product domain.Product `json:"product"`
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

func patchProduct(t *testing.T, ts *TestClient, adminCookie *http.Cookie, products []domain.Product) {
	t.Run("valid update", func(t *testing.T) {
		newProduct := products[0]
		newProduct.Name = "product 1 updated"
		newProduct.Description = "product 1 description updated"
		newProduct.AvailableAmount = 53
		newProduct.Version++

		res, err := ts.PatchWithCookies("/products/1", toPostProduct(newProduct), adminCookie)
		assert.Nil(t, err)

		var responseBody struct {
			Product domain.Product `json:"product"`
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

		res, err := ts.PatchWithCookies("/products/3", toPostProduct(newProduct), adminCookie)
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

type PostProduct struct {
	Name        string            `json:"name,omitempty"`
	Description string            `json:"description,omitempty"`
	Properties  map[string]string `json:"properties,omitempty"`
	Vendor      string            `json:"vendor,omitempty"`
	Amount      int               `json:"amount,omitempty"`
	Price       *money.Money      `json:"price,omitempty"`
}

func toPostProduct(p domain.Product) PostProduct {
	return PostProduct{
		Name:        p.Name,
		Description: p.Description,
		Properties:  p.Properties,
		Vendor:      p.Vendor,
		Amount:      p.AvailableAmount,
		Price:       p.Price,
	}
}

func validPost(t *testing.T, ts *TestClient, adminCookie *http.Cookie, products []domain.Product) {
	for i, product := range products {
		t.Run(fmt.Sprintf("product%d", i), func(t *testing.T) {
			res, err := ts.PostWithCookies("/products", toPostProduct(product), adminCookie)
			assert.Nil(t, err)

			var responseBody struct {
				Product domain.Product `json:"product"`
			}

			err = ts.ReadResponseBody(res, &responseBody)
			assert.Nil(t, err)
			assert.Equal(t, res.StatusCode, http.StatusCreated)
			assert.Equal(t, responseBody.Product.String(), product.String())
		})
	}
}

func invalidPost(t *testing.T, ts *TestClient, adminCookie *http.Cookie) {
	product := domain.Product{
		ID:              3,
		Name:            "product 2",
		Description:     "description of product 2",
		Vendor:          "vendor 1",
		AvailableAmount: 4,
		Price:           money.New(2399, money.EGP),
		Version:         1,
		Properties: map[string]string{
			"size": "14",
		},
	}
	t.Run("duplicate product", func(t *testing.T) {
		var responseBody struct {
			Error string `json:"error"`
		}

		res, err := ts.PostWithCookies("/products", toPostProduct(product), adminCookie)
		assert.Nil(t, err)

		err = ts.ReadResponseBody(res, &responseBody)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		assert.Equal(t, responseBody.Error, "product already exists")
	})
	t.Run("invalid product amount", func(t *testing.T) {
		p := product
		p.AvailableAmount = -1

		var responseBody struct {
			Error struct {
				Amount string `json:"amount"`
			} `json:"error"`
		}

		res, err := ts.PostWithCookies("/products", toPostProduct(p), adminCookie)
		assert.Nil(t, err)

		err = ts.ReadResponseBody(res, &responseBody)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
		assert.Equal(t, responseBody.Error.Amount, "must be more than 0")
	})
	t.Run("missing product name and description", func(t *testing.T) {
		p := product
		p.Name = ""
		p.Description = ""

		var responseBody struct {
			Error struct {
				Name        string `json:"name"`
				Description string `json:"description"`
			} `json:"error"`
		}

		res, err := ts.PostWithCookies("/products", toPostProduct(p), adminCookie)
		assert.Nil(t, err)

		ts.ReadResponseBody(res, &responseBody)

		assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
		assert.Equal(t, responseBody.Error.Name, "can't be empty")
		assert.Equal(t, responseBody.Error.Description, "can't be empty")
	})
	t.Run("missing product name", func(t *testing.T) {
		p := product
		p.Name = ""

		var responseBody struct {
			Error struct {
				Name string `json:"name"`
			} `json:"error"`
		}

		res, err := ts.PostWithCookies("/products", toPostProduct(p), adminCookie)
		assert.Nil(t, err)

		ts.ReadResponseBody(res, &responseBody)

		assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
		assert.Equal(t, responseBody.Error.Name, "can't be empty")
	})
	t.Run("invalid product", func(t *testing.T) {
		p := product

		var responseBody struct {
			Error string `json:"error"`
		}

		res, err := ts.PostWithCookies("/products", p, adminCookie)
		assert.Nil(t, err)

		err = ts.ReadResponseBody(res, &responseBody)
		assert.Nil(t, err)

		assert.Equal(t, res.StatusCode, http.StatusBadRequest)
		assert.Equal(t, responseBody.Error, `body contains unknown key "id"`)
	})
}
